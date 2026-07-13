package originalgame

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wangle201210/zskc/internal/original"
)

const (
	worldMapLoadingSteps = 15
	worldMapTravelTicks  = 4
	worldMapOriginX      = 37
	worldMapOriginY      = 73
	worldMapGridSize     = 13
)

type worldMapPoint struct {
	X int
	Y int
}

func (p *worldMapPoint) UnmarshalJSON(data []byte) error {
	var pair [2]int
	if err := json.Unmarshal(data, &pair); err != nil {
		return fmt.Errorf("decode world-map point: %w", err)
	}
	p.X = pair[0]
	p.Y = pair[1]
	return nil
}

type worldMapNode struct {
	X     int             `json:"x"`
	Y     int             `json:"y"`
	Type  int             `json:"type"`
	Stage int             `json:"stage"`
	Links []worldMapPoint `json:"links"`
}

type worldMapData struct {
	Source        string         `json:"source"`
	PayloadLength int            `json:"payload_length"`
	Nodes         []worldMapNode `json:"nodes"`
	byStage       map[int]int
	byPoint       map[worldMapPoint]int
}

func loadWorldMap(path string) (*worldMapData, error) {
	data, err := os.ReadFile(filepath.Clean(resolvePath(path)))
	if err != nil {
		return nil, err
	}
	var result worldMapData
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode world map: %w", err)
	}
	result.byStage = make(map[int]int, len(result.Nodes))
	result.byPoint = make(map[worldMapPoint]int, len(result.Nodes))
	for index, node := range result.Nodes {
		if node.X < 0 || node.X >= 12 || node.Y < 0 || node.Y >= 12 {
			return nil, fmt.Errorf("world map node %d has invalid position %d,%d", index, node.X, node.Y)
		}
		if _, exists := result.byStage[node.Stage]; exists {
			return nil, fmt.Errorf("world map stage %d is duplicated", node.Stage)
		}
		point := worldMapPoint{X: node.X, Y: node.Y}
		if _, exists := result.byPoint[point]; exists {
			return nil, fmt.Errorf("world map position %d,%d is duplicated", node.X, node.Y)
		}
		result.byStage[node.Stage] = index
		result.byPoint[point] = index
	}
	for _, node := range result.Nodes {
		for _, link := range node.Links {
			if _, ok := result.byPoint[link]; !ok {
				return nil, fmt.Errorf("world map stage %d links to missing node %d,%d", node.Stage, link.X, link.Y)
			}
		}
	}
	return &result, nil
}

func (m *worldMapData) nodeForStage(stage int) (worldMapNode, bool) {
	if m == nil {
		return worldMapNode{}, false
	}
	index, ok := m.byStage[stage]
	if !ok || index < 0 || index >= len(m.Nodes) {
		return worldMapNode{}, false
	}
	return m.Nodes[index], true
}

func (m *worldMapData) linkedStage(stage, dx, dy int, unlocked func(int) bool) (int, bool) {
	node, ok := m.nodeForStage(stage)
	if !ok {
		return 0, false
	}
	bestStage := -1
	bestPrimary := 1 << 30
	bestSecondary := 1 << 30
	for _, point := range node.Links {
		index, exists := m.byPoint[point]
		if !exists {
			continue
		}
		candidate := m.Nodes[index]
		// Java checks the node's runtime lock state here, not its authored map
		// type. An unlocked secret node is therefore navigable in both directions.
		if unlocked != nil && !unlocked(candidate.Stage) {
			continue
		}
		deltaX := candidate.X - node.X
		deltaY := candidate.Y - node.Y
		if dx != 0 && deltaX*dx <= 0 || dy != 0 && deltaY*dy <= 0 {
			continue
		}
		primary, secondary := abs(deltaX), abs(deltaY)
		if dy != 0 {
			primary, secondary = abs(deltaY), abs(deltaX)
		}
		if primary < bestPrimary || primary == bestPrimary && secondary < bestSecondary {
			bestStage = candidate.Stage
			bestPrimary = primary
			bestSecondary = secondary
		}
	}
	return bestStage, bestStage >= 0
}

func (m *worldMapData) exitTarget(stage int, secret bool) (int, bool) {
	node, ok := m.nodeForStage(stage)
	if !ok {
		return 0, false
	}
	if secret && len(node.Links) <= 1 {
		return stage, true
	}
	target := -1
	for _, point := range node.Links {
		index, exists := m.byPoint[point]
		if !exists {
			continue
		}
		candidate := m.Nodes[index]
		if secret {
			// zVoid selects a connected type-1 node whose stage index is
			// greater than the current node. Angkor has at most one such edge.
			if candidate.Type == 1 && candidate.Stage > node.Stage {
				target = candidate.Stage
			}
			continue
		}
		if candidate.Type == 0 && candidate.Stage > node.Stage && (target < 0 || candidate.Stage < target) {
			target = candidate.Stage
		}
	}
	return target, target >= 0
}

func (m *worldMapData) stagesLinked(fromStage, toStage int) bool {
	from, ok := m.nodeForStage(fromStage)
	if !ok {
		return false
	}
	to, ok := m.nodeForStage(toStage)
	if !ok {
		return false
	}
	for _, point := range from.Links {
		if point == (worldMapPoint{X: to.X, Y: to.Y}) {
			return true
		}
	}
	return fromStage == toStage
}

func (g *Game) enterWorldMap() {
	g.mode = gameModeWorldMap
	g.worldMapLoadingStep = 0
	g.worldMapSelectedStage = clamp(g.stageIndex, 0, worldStageCount(g.worldIndex)-1)
	g.worldMapTravelFrom = g.worldMapSelectedStage
	g.worldMapTravelTo = g.worldMapSelectedStage
	g.worldMapTravelTick = 0
	if g.pendingMapTarget >= 0 && g.progress.stageUnlockedForWorld(g.worldIndex, g.pendingMapTarget) && g.worldMap.stagesLinked(g.stageIndex, g.pendingMapTarget) {
		g.worldMapTravelTo = g.pendingMapTarget
	}
	g.pendingMapTarget = -1
	g.message = worldName(g.worldIndex) + " world map"
}

func (g *Game) updateWorldMap(action bool) {
	if g.worldMapLoadingStep < worldMapLoadingSteps {
		g.worldMapLoadingStep++
		return
	}
	if g.worldMapTravelTick < worldMapTravelTicks {
		g.worldMapTravelTick++
		if g.worldMapTravelTick == worldMapTravelTicks {
			g.worldMapSelectedStage = g.worldMapTravelTo
			g.worldMapTravelFrom = g.worldMapTravelTo
		}
		return
	}
	if g.sourceInput.Navigate {
		g.pendingMapTarget = -1
		g.sealExitActive = true
		g.sealExitTicks = 0
		g.sealExitIncoming = -1
		g.message = "world selection"
		return
	}
	if action {
		stage := g.worldMapSelectedStage
		if !g.progress.stageUnlockedForWorld(g.worldIndex, stage) {
			g.message = fmt.Sprintf("%s stage %d is locked", worldName(g.worldIndex), stage+1)
		} else if worldStageImplemented(g.worldIndex, stage) {
			g.loadStage(stage)
			g.mode = gameModeStage
		} else {
			g.message = fmt.Sprintf("%s stage %d is not replicated yet", worldName(g.worldIndex), stage+1)
		}
		return
	}
	dx, dy := g.sourceInput.DirectionDX, g.sourceInput.DirectionDY
	if dx == 0 && dy == 0 {
		return
	}
	unlocked := func(stage int) bool {
		return g.progress.stageUnlockedForWorld(g.worldIndex, stage)
	}
	if stage, ok := g.worldMap.linkedStage(g.worldMapSelectedStage, dx, dy, unlocked); ok {
		g.worldMapTravelFrom = g.worldMapSelectedStage
		g.worldMapTravelTo = stage
		g.worldMapTravelTick = 0
	}
}

func justPressedDirection() (int, int) {
	return heldDirectionWith(inpututil.IsKeyJustPressed)
}

func (g *Game) drawWorldMap(screen *ebiten.Image) {
	if g.worldMapLoadingStep < worldMapLoadingSteps {
		g.drawWorldMapLoading(screen)
		return
	}
	screen.Fill(color.RGBA{0x0e, 0x55, 0x12, 0xff})
	if g.worldMapHeader != nil {
		g.worldMapHeader.drawFrame(screen, 0, original.ScreenWidth/2, 0, 0)
	}
	g.fontMedium.drawText(screen, worldDisplayName(g.worldIndex), original.ScreenWidth/2, 15, true, color.White)
	g.fontMedium.drawText(screen, g.worldMapStageTitle(g.worldMapDisplayStage()), 8, 61, false, color.White)
	if g.worldMapGround != nil {
		g.worldMapGround.drawFrame(screen, 0, 120, 171, 0)
	}
	g.drawWorldMapPaths(screen)
	g.drawWorldMapNodes(screen)
	g.drawWorldMapHero(screen)
	g.drawWorldMapStageProgress(screen)

	drawRoundedRect(screen, 2, 275, 236, 16, 4, color.RGBA{0x2f, 0x7b, 0x46, 0xff})
	if g.worldMapIcons != nil {
		g.worldMapIcons.drawFrame(screen, 12, 10, 278, 0)
		g.worldMapIcons.drawFrame(screen, 11, 80, 280, 0)
		g.worldMapIcons.drawFrame(screen, 10, 155, 280, 0)
	}
	g.fontSmall.drawText(screen, fmt.Sprintf("%d", g.progress.ExtraLives), 41, 290, true, color.White)
	g.fontSmall.drawText(screen, fmt.Sprintf("%d", g.progress.VioletGemBank), 99, 290, true, color.White)
	g.fontSmall.drawText(screen, fmt.Sprintf("%d", g.progress.RedDiamondBank), 174, 290, true, color.White)
	selectPrompt := tr(textPromptSelect, desktopActionKeyLabel)
	worldsPrompt := tr(textPromptWorlds, desktopNavigationKeyLabel)
	g.fontSmall.drawText(screen, worldsPrompt, 2, 314, false, color.White)
	g.fontSmall.drawText(screen, selectPrompt, 236-g.fontSmall.stringWidth(selectPrompt), 314, false, color.White)
}

func (g *Game) drawWorldMapStageProgress(screen *ebiten.Image) {
	if g.worldMap == nil || g.fontSmall == nil || g.worldMapTravelTick < worldMapTravelTicks {
		return
	}
	stage := g.worldMapSelectedStage
	node, ok := g.worldMap.nodeForStage(stage)
	if !ok {
		return
	}
	total := g.stageTotalRed(stage)
	collected := min(total, g.progress.stageRedDiamondsForWorld(g.worldIndex, stage))
	label := fmt.Sprintf("%d/%d", collected, total)
	width := g.fontSmall.stringWidth(label) + 20
	centerX, centerY := worldMapScreenPoint(node)
	x := centerX - width/2
	y := centerY - 43
	if y <= 63 {
		y = 63
		x = centerX + 20
		if x+width >= 220 {
			x = centerX - width - 20
		}
	}
	x = clamp(x, 25, 220-width)
	drawRoundedRect(screen, x, y, width, 17, 4, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	drawRoundedRect(screen, x+1, y+1, width-2, 15, 3, color.RGBA{0x00, 0x90, 0xb2, 0xff})
	g.fontSmall.drawText(screen, label, x+2, y+10, false, color.White)
	if g.worldMapIcons != nil {
		g.worldMapIcons.drawFrame(screen, 10, x+width-16, y+1, 0)
	}
}

func (g *Game) drawWorldMapLoading(screen *ebiten.Image) {
	screen.Fill(color.Black)
	progress := min(230, (g.worldMapLoadingStep+1)*230/worldMapLoadingSteps)
	drawRect(screen, 5, 310, progress, 6, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	drawRect(screen, 4, 309, 231, 1, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	drawRect(screen, 4, 316, 231, 1, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	g.fontMedium.drawText(screen, tr(textLoading), original.ScreenWidth/2, 304, true, color.White)
}

func (g *Game) drawWorldMapPaths(screen *ebiten.Image) {
	if g.worldMap == nil || g.worldMapIcons == nil {
		return
	}
	drawn := map[[2]int]bool{}
	for _, node := range g.worldMap.Nodes {
		for _, linkedPoint := range node.Links {
			linkedIndex := g.worldMap.byPoint[linkedPoint]
			linked := g.worldMap.Nodes[linkedIndex]
			key := [2]int{min(node.Stage, linked.Stage), max(node.Stage, linked.Stage)}
			if drawn[key] {
				continue
			}
			drawn[key] = true
			unlocked := g.progress.stageUnlockedForWorld(g.worldIndex, node.Stage) && g.progress.stageUnlockedForWorld(g.worldIndex, linked.Stage)
			frame := 2
			if !unlocked {
				frame = 3
			}
			drawMapLine(screen, g.worldMapIcons, frame, node, linked)
		}
	}
}

func drawMapLine(screen *ebiten.Image, icons *spriteSheet, frame int, from, to worldMapNode) {
	x0, y0 := worldMapScreenPoint(from)
	x1, y1 := worldMapScreenPoint(to)
	dx, dy := x1-x0, y1-y0
	steps := max(abs(dx), abs(dy)) / 8
	if steps < 1 {
		steps = 1
	}
	for step := 1; step < steps; step++ {
		x := x0 + dx*step/steps
		y := y0 + dy*step/steps
		icons.drawFrame(screen, frame, x, y, 0)
	}
}

func (g *Game) drawWorldMapNodes(screen *ebiten.Image) {
	if g.worldMap == nil || g.worldMapIcons == nil {
		return
	}
	for _, node := range g.worldMap.Nodes {
		frame := 1
		if g.progress.stageUnlockedForWorld(g.worldIndex, node.Stage) {
			frame = 0
		}
		x, y := worldMapScreenPoint(node)
		g.worldMapIcons.drawFrame(screen, frame, x, y, 0)
		if g.progress.stageClearedForWorld(g.worldIndex, node.Stage) {
			g.worldMapIcons.drawFrame(screen, 17, x, y, 0)
		}
	}
}

func (g *Game) drawWorldMapHero(screen *ebiten.Image) {
	if g.worldMap == nil || g.worldMapIcons == nil {
		return
	}
	from, fromOK := g.worldMap.nodeForStage(g.worldMapTravelFrom)
	to, toOK := g.worldMap.nodeForStage(g.worldMapTravelTo)
	if !fromOK || !toOK {
		return
	}
	x0, y0 := worldMapScreenPoint(from)
	x1, y1 := worldMapScreenPoint(to)
	tick := clamp(g.worldMapTravelTick*renderStepsPerSource+g.renderPhase, 0, worldMapTravelTicks*renderStepsPerSource)
	duration := worldMapTravelTicks * renderStepsPerSource
	x := x0 + (x1-x0)*tick/duration
	y := y0 + (y1-y0)*tick/duration
	frame := 6
	if x1 < x0 {
		frame = 7
	}
	g.worldMapIcons.drawFrame(screen, frame, x, y, 0)
}

func (g *Game) worldMapDisplayStage() int {
	if g.worldMapTravelTick < worldMapTravelTicks {
		return clamp(g.worldMapTravelTo, 0, worldStageCount(g.worldIndex)-1)
	}
	return clamp(g.worldMapSelectedStage, 0, worldStageCount(g.worldIndex)-1)
}

func (g *Game) worldMapStageTitle(stage int) string {
	if node, ok := g.worldMap.nodeForStage(stage); ok && node.Type == 1 {
		return tr(textSecretStage, stage-worldFirstSecretStage(g.worldIndex)+1)
	}
	return tr(textStage, stage+1)
}

func angkorStageImplemented(stage int) bool {
	return stage >= 0 && stage < angkorReplicaStageCount
}

func (g *Game) stageTotalViolet(stage int) int {
	if g.pack == nil || stage < 0 || stage >= len(g.pack.Stages) {
		return 0
	}
	return original.StageVioletTotal(g.pack.Stages[stage])
}

func (g *Game) stageTotalRed(stage int) int {
	if g.pack == nil || stage < 0 || stage >= len(g.pack.Stages) || g.pack.Stages[stage] == nil {
		return 0
	}
	total := 0
	for _, id := range g.pack.Stages[stage].Player {
		if id == 2 {
			total++
		}
	}
	return total
}

func worldMapScreenPoint(node worldMapNode) (int, int) {
	return worldMapOriginX + node.X*worldMapGridSize + 6,
		worldMapOriginY + node.Y*worldMapGridSize + 6
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
