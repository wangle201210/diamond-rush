package originalgame

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

const (
	sealPositionAngkor = iota
	sealPositionBavaria
	sealPositionSiberia
	sealPositionShop
)

const (
	sealMoveUp = iota
	sealMoveRight
	sealMoveDown
	sealMoveLeft
)

const (
	sealMoveTicks          = 8
	sealRelicFlashTicks    = 20
	sealRelicEffectTicks   = 22
	sealUnlockEffectTicks  = 22
	sealUnlockVisibleTicks = 15
)

var sealWorldPrices = [...]int{0, 10, 25}

var sealItemOffsets = [...][2]int{
	{-24, -23},
	{24, -23},
	{0, 23},
	{24, 23},
}

var sealArrowOffsets = [...][2]int{
	{-33, -54},
	{14, -54},
	{-8, -8},
	{22, 2},
}

var sealMoveTargets = [...][4]int{
	{-1, -1, sealPositionAngkor, -1},
	{sealPositionBavaria, -1, sealPositionShop, -1},
	{sealPositionSiberia, sealPositionSiberia, -1, -1},
	{-1, sealPositionAngkor, -1, sealPositionSiberia},
}

func (p *originalProgress) unlockEligibleWorlds() int {
	if p == nil {
		return 0
	}
	p.WorldUnlocked[sealPositionAngkor] = true
	newest := 0
	for world := sealPositionBavaria; world <= sealPositionSiberia; world++ {
		if !p.WorldUnlocked[world] && p.RedDiamondBank >= sealWorldPrices[world] {
			p.WorldUnlocked[world] = true
			newest = world
		}
	}
	return newest
}

func (g *Game) enterWorldSelect(incomingRelic int) {
	g.mode = gameModeWorldSelect
	g.worldSelectUnlocking = g.progress.unlockEligibleWorlds()
	g.worldSelectPosition = g.worldIndex
	if g.worldSelectPosition < sealPositionAngkor || g.worldSelectPosition > sealPositionSiberia {
		g.worldSelectPosition = sealPositionAngkor
	}
	if g.worldSelectUnlocking != 0 {
		g.worldSelectPosition = g.worldSelectUnlocking
	}
	offset := sealArrowOffsets[g.worldSelectPosition]
	g.worldSelectArrowX = offset[0]
	g.worldSelectArrowY = offset[1]
	g.worldSelectTargetX = offset[0]
	g.worldSelectTargetY = offset[1]
	g.worldSelectMoveTick = sealMoveTicks
	g.worldSelectArrowTick = 0
	g.worldSelectIncoming = incomingRelic
	if incomingRelic < sealPositionAngkor || incomingRelic > sealPositionSiberia {
		g.worldSelectIncoming = -1
	}
	g.worldSelectRelicX = 10
	g.worldSelectRelicY = 10
	g.worldSelectFlashTick = 0
	g.worldSelectEffectTick = 0
	g.worldSelectUnlockTick = 0
	g.worldSelectUnlockFlash = 0
	g.pendingMapTarget = -1
	g.message = "seal world selection"
	if g.sounds != nil {
		g.sounds.Stop()
	}
	if g.progressPath != "" {
		if err := saveOriginalProgress(g.progressPath, g.progress); err != nil {
			g.message = err.Error()
		}
	}
}

func (g *Game) updateWorldSelect(action bool, dx, dy int) {
	if g.worldSelectMoveTick < sealMoveTicks {
		remaining := sealMoveTicks - g.worldSelectMoveTick
		g.worldSelectArrowX += (g.worldSelectTargetX - g.worldSelectArrowX) / remaining
		g.worldSelectArrowY += (g.worldSelectTargetY - g.worldSelectArrowY) / remaining
		g.worldSelectMoveTick++
		if g.worldSelectMoveTick == sealMoveTicks {
			g.worldSelectArrowX = g.worldSelectTargetX
			g.worldSelectArrowY = g.worldSelectTargetY
		}
		g.worldSelectArrowTick++
		return
	}

	advancedUnlockEffect := false
	if g.worldSelectUnlocking != 0 && g.worldSelectUnlockTick < sealUnlockEffectTicks {
		g.worldSelectUnlockTick++
		advancedUnlockEffect = true
	}
	if g.worldSelectIncoming >= 0 {
		g.updateIncomingSealRelic()
		return
	}
	if g.worldSelectUnlocking != 0 {
		if !advancedUnlockEffect {
			g.updateWorldUnlockAnimation()
		}
		g.worldSelectArrowTick++
		return
	}

	g.worldSelectArrowTick++
	if g.sourceInput.Recall {
		g.enterStartMenu(true)
		return
	}
	if action {
		g.activateWorldSelectPosition()
		return
	}
	move := -1
	switch {
	case dy < 0:
		move = sealMoveUp
	case dx > 0:
		move = sealMoveRight
	case dy > 0:
		move = sealMoveDown
	case dx < 0:
		move = sealMoveLeft
	}
	if move < 0 {
		return
	}
	target := sealMoveTargets[move][g.worldSelectPosition]
	if target < 0 {
		return
	}
	g.worldSelectPosition = target
	g.worldSelectTargetX = sealArrowOffsets[target][0]
	g.worldSelectTargetY = sealArrowOffsets[target][1]
	g.worldSelectMoveTick = 0
}

func (g *Game) updateIncomingSealRelic() {
	target := sealItemOffsets[g.worldSelectIncoming]
	targetX := original.ScreenWidth/2 + target[0]
	targetY := 136 + target[1]
	g.worldSelectRelicX = approach(g.worldSelectRelicX, targetX, 3)
	g.worldSelectRelicY = approach(g.worldSelectRelicY, targetY, 2)
	if g.worldSelectRelicX != targetX || g.worldSelectRelicY != targetY {
		return
	}
	if g.worldSelectFlashTick < sealRelicFlashTicks {
		g.worldSelectFlashTick++
		return
	}
	if g.worldSelectEffectTick < sealRelicEffectTicks {
		g.worldSelectEffectTick++
		return
	}
	g.worldSelectIncoming = -1
	g.worldSelectEffectTick = 0
}

func (g *Game) updateWorldUnlockAnimation() {
	if g.worldSelectUnlockTick < sealUnlockEffectTicks {
		return
	}
	if g.worldSelectArrowTick%8 >= 4 {
		g.worldSelectUnlockFlash++
	}
	if g.worldSelectUnlockFlash >= sealUnlockVisibleTicks {
		g.worldSelectUnlocking = 0
		g.worldSelectUnlockTick = 0
		g.worldSelectUnlockFlash = 0
	}
}

func (g *Game) activateWorldSelectPosition() {
	switch g.worldSelectPosition {
	case sealPositionAngkor, sealPositionBavaria:
		world := g.worldSelectPosition
		if !g.progress.WorldUnlocked[world] {
			return
		}
		if err := g.switchWorld(world); err != nil {
			g.message = err.Error()
			return
		}
		g.progress.LastWorld = world
		g.progress = g.progress.normalized()
		if g.progressPath != "" {
			if err := saveOriginalProgress(g.progressPath, g.progress); err != nil {
				g.message = err.Error()
				return
			}
		}
		g.stageIndex = g.highestUnlockedMapStageForWorld(world)
		g.enterWorldMap()
		if g.sounds != nil && g.sounds.enabled {
			g.sounds.Play(worldMusic(world))
		}
	case sealPositionSiberia:
		if g.progress.WorldUnlocked[g.worldSelectPosition] {
			g.message = fmt.Sprintf("%s world is not replicated yet", sealWorldName(g.worldSelectPosition))
		}
	case sealPositionShop:
		g.message = "Shop is not replicated yet"
	}
}

func (g *Game) drawWorldSelect(screen *ebiten.Image) {
	for y := 0; y < original.ScreenHeight; y += original.TileSize {
		for x := 0; x < original.ScreenWidth; x += original.TileSize {
			g.floor.drawModule(screen, 0, x, y)
		}
	}
	if g.tutorialSeal != nil {
		g.tutorialSeal.drawFrame(screen, 0, 60, 76, 0)
	}
	if g.worldMapIcons != nil {
		g.worldMapIcons.drawFrame(screen, 11, 144, 159, 0)
	}
	if g.softkeys != nil {
		g.softkeys.drawFrame(screen, 0, 223, 308, 0)
		g.softkeys.drawFrame(screen, 3, 2, 308, 0)
	}
	if g.fontSmall != nil {
		text := desktopRecallKeyLabel + ": MAIN MENU"
		g.fontSmall.drawText(screen, text, 218-g.fontSmall.stringWidth(text), 314, false, color.White)
	}

	for world := sealPositionAngkor; world <= sealPositionSiberia; world++ {
		if g.worldSelectOverlayVisible(world) && g.tutorialSeal != nil {
			g.tutorialSeal.drawFrame(screen, world+1, 60, 76, 0)
		}
	}
	if g.worldSelectUnlocking != 0 && g.worldSelectUnlockTick < sealUnlockEffectTicks && g.pickupEffects != nil {
		offset := sealItemOffsets[g.worldSelectUnlocking]
		g.pickupEffects.drawAnimationSequenceFrame(screen, 5, g.worldSelectUnlockTick, original.ScreenWidth/2+offset[0]-12, 124+offset[1], 0)
	}
	for relic := sealPositionAngkor; relic <= sealPositionSiberia; relic++ {
		if g.progress.RelicMask&(1<<relic) == 0 || g.worldSelectIncoming == relic {
			continue
		}
		g.drawSealRelic(screen, relic, original.ScreenWidth/2+sealItemOffsets[relic][0], 136+sealItemOffsets[relic][1])
	}
	if g.worldSelectIncoming >= 0 {
		g.drawIncomingSealRelic(screen)
		return
	}
	if g.sealArrow != nil {
		x, y := g.renderedWorldSelectArrow()
		g.sealArrow.drawAnimationSequenceFrame(screen, 0, g.worldSelectArrowTick, original.ScreenWidth/2+x, 136+y, 0)
	}
	if g.worldSelectUnlocking == 0 {
		g.drawWorldSelectPrompt(screen)
	}
}

func (g *Game) worldSelectOverlayVisible(world int) bool {
	if world < sealPositionAngkor || world > sealPositionSiberia || !g.progress.WorldUnlocked[world] {
		return false
	}
	if g.worldSelectUnlocking == 0 || world < g.worldSelectUnlocking {
		return true
	}
	return world == g.worldSelectUnlocking &&
		g.worldSelectUnlockTick >= sealUnlockEffectTicks &&
		g.worldSelectArrowTick%8 >= 4
}

func (g *Game) drawIncomingSealRelic(screen *ebiten.Image) {
	relic := g.worldSelectIncoming
	x, y := g.renderedIncomingSealRelic()
	g.drawSealRelic(screen, relic, x, y)
	if g.worldSelectFlashTick > 0 && g.worldSelectFlashTick < sealRelicFlashTicks && g.worldSelectFlashTick%2 == 0 {
		rgb := 838860 * (g.worldSelectFlashTick - 1)
		screen.Fill(color.RGBA{uint8(rgb >> 16), uint8(rgb >> 8), uint8(rgb), 0xff})
	}
	if g.worldSelectFlashTick >= sealRelicFlashTicks && g.worldSelectEffectTick > 0 && g.worldSelectEffectTick <= sealRelicEffectTicks && g.pickupEffects != nil {
		offset := sealItemOffsets[relic]
		g.pickupEffects.drawAnimationSequenceFrame(screen, 5, g.worldSelectEffectTick-1, original.ScreenWidth/2+offset[0]-12, 124+offset[1], 0)
	}
}

func (g *Game) renderedWorldSelectArrow() (int, int) {
	x, y := g.worldSelectArrowX, g.worldSelectArrowY
	if g.renderPhase <= 0 || g.worldSelectMoveTick >= sealMoveTicks {
		return x, y
	}
	remaining := max(1, sealMoveTicks-g.worldSelectMoveTick)
	nextX := x + (g.worldSelectTargetX-x)/remaining
	nextY := y + (g.worldSelectTargetY-y)/remaining
	return x + (nextX-x)*g.renderPhase/renderStepsPerSource,
		y + (nextY-y)*g.renderPhase/renderStepsPerSource
}

func (g *Game) renderedIncomingSealRelic() (int, int) {
	x, y := g.worldSelectRelicX, g.worldSelectRelicY
	if g.renderPhase <= 0 || g.worldSelectIncoming < 0 || g.worldSelectIncoming >= len(sealItemOffsets) {
		return x, y
	}
	target := sealItemOffsets[g.worldSelectIncoming]
	targetX := original.ScreenWidth/2 + target[0]
	targetY := 136 + target[1]
	nextX := approach(x, targetX, 3)
	nextY := approach(y, targetY, 2)
	return x + (nextX-x)*g.renderPhase/renderStepsPerSource,
		y + (nextY-y)*g.renderPhase/renderStepsPerSource
}

func (g *Game) drawSealRelic(screen *ebiten.Image, relic, centerX, centerY int) {
	var sheet *spriteSheet
	switch relic {
	case sealPositionAngkor:
		sheet = g.angkorSeal
	case sealPositionBavaria:
		sheet = g.bavariaSeal
	case sealPositionSiberia:
		sheet = g.siberiaSeal
	}
	if sheet == nil || len(sheet.meta.Modules) == 0 {
		return
	}
	module := sheet.meta.Modules[0]
	sheet.drawFrame(screen, 0, centerX-module.W/2, centerY-module.H/2, 0)
}

func (g *Game) drawWorldSelectPrompt(screen *ebiten.Image) {
	if g.fontSmall == nil {
		return
	}
	action := "Press " + desktopActionKeyLabel + " to go to"
	if g.worldSelectPosition != sealPositionShop && !g.progress.WorldUnlocked[g.worldSelectPosition] {
		action = fmt.Sprintf("%d red diamonds to go to", sealWorldPrices[g.worldSelectPosition])
	}
	lines := []string{action, sealWorldName(g.worldSelectPosition)}
	width := max(g.fontSmall.stringWidth(lines[0]), g.fontSmall.stringWidth(lines[1]))
	lineHeight := g.fontSmall.meta.FontHeight
	height := lineHeight * len(lines)
	x := (original.ScreenWidth-width)/2 - 3
	y := 240 - g.fontSmall.meta.FontHeight/2 - 3
	g.drawDemoPanel(screen, x, y, width+6, height+5, color.RGBA{0x6c, 0x49, 0x0b, 0xff})
	top := 240 - height/2
	for index, line := range lines {
		g.fontSmall.drawText(screen, line, original.ScreenWidth/2, top+sourceFontYOffset+index*lineHeight, true, color.White)
	}
}

func (g *Game) drawDemoPanel(screen *ebiten.Image, x, y, width, height int, fill color.Color) {
	g.drawDemoPanelWithSheet(screen, g.demoUI, x, y, width, height, fill)
}

func (g *Game) drawDemoPanelWithSheet(screen *ebiten.Image, sheet *spriteSheet, x, y, width, height int, fill color.Color) {
	drawRect(screen, x, y, width, height, fill)
	if sheet == nil || len(sheet.meta.Modules) < 8 {
		return
	}
	for py := y; py < y+height; py += 8 {
		sheet.drawModule(screen, 7, x-3, py)
		sheet.drawModule(screen, 5, x+width, py)
	}
	for px := x; px < x+width; px += 8 {
		sheet.drawModule(screen, 4, px, y-3)
		sheet.drawModule(screen, 6, px, y+height)
	}
	sheet.drawModule(screen, 0, x-3, y-3)
	sheet.drawModule(screen, 1, x+width, y-3)
	sheet.drawModule(screen, 2, x-3, y+height)
	sheet.drawModule(screen, 3, x+width, y+height)
}

func sealWorldName(position int) string {
	switch position {
	case sealPositionAngkor:
		return "Angkor Wat"
	case sealPositionBavaria:
		return "Bavaria"
	case sealPositionSiberia:
		return "Siberia"
	case sealPositionShop:
		return "Shop"
	default:
		return ""
	}
}

func approach(value, target, step int) int {
	switch {
	case value < target:
		return min(target, value+step)
	case value > target:
		return max(target, value-step)
	default:
		return value
	}
}
