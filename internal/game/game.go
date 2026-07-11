package game

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wangle201210/zskc/internal/level"
	"github.com/wangle201210/zskc/internal/world"
)

const (
	tileSize       = 32
	viewportTilesX = 15
	viewportTilesY = 11
	moveCooldown   = 8
	animFrames     = 7
	hudHeight      = 96
	dpadSize       = 72
	dpadMargin     = 12
)

type screenMode int

const (
	modeTitle screenMode = iota
	modeSelect
	modePlay
	modePause
)

type Game struct {
	world         *world.World
	assets        *Assets
	sounds        *Sounds
	levels        []string
	titles        []string
	themes        []string
	parSteps      []int
	levelIndex    int
	selected      int
	titleSelected int
	cooldown      int
	facing        world.Direction
	message       string
	hint          string
	playerAnim    tileAnim
	mode          screenMode
	progress      Progress
	savePath      string
	winSaved      bool
	clearStars    int
	clearScore    int
	newBest       bool
	newBestScore  bool
	restartUsed   bool
	tick          int
}

type tileAnim struct {
	active bool
	fromX  int
	fromY  int
	toX    int
	toY    int
	frame  int
	frames int
}

func Run() error {
	levels := []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"}
	titles, themes, parSteps, err := loadLevelMeta(levels)
	if err != nil {
		return err
	}
	savePath := progressPath()
	progress, err := loadProgress(savePath, len(levels))
	if err != nil {
		return err
	}
	g := &Game{
		assets:   NewAssets(tileSize),
		sounds:   NewSounds(),
		levels:   levels,
		titles:   titles,
		themes:   themes,
		parSteps: parSteps,
		selected: progress.UnlockedLevel - 1,
		message:  "Collect diamonds, open the exit.",
		mode:     modeTitle,
		progress: progress,
		savePath: savePath,
	}
	if err := g.loadLevel(0); err != nil {
		return err
	}
	ebiten.SetWindowTitle("Diamond Rush - Go/Ebitengine")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	return ebiten.RunGame(g)
}

func loadLevelMeta(files []string) ([]string, []string, []int, error) {
	titles := make([]string, len(files))
	themes := make([]string, len(files))
	parSteps := make([]int, len(files))
	for i, file := range files {
		def, err := level.LoadFile(file)
		if err != nil {
			return nil, nil, nil, err
		}
		titles[i] = def.Title
		themes[i] = def.Theme
		parSteps[i] = def.ParSteps
	}
	return titles, themes, parSteps, nil
}

func (g *Game) Update() error {
	g.tick++
	if isAudioPressed() {
		if g.sounds.ToggleMute() {
			g.message = "Audio muted."
		} else {
			g.message = "Audio on."
		}
	}
	switch g.mode {
	case modeTitle:
		return g.updateTitle()
	case modeSelect:
		return g.updateSelect()
	case modePause:
		return g.updatePause()
	default:
		return g.updatePlay()
	}
}

func (g *Game) updateTitle() error {
	if isMenuUpPressed() {
		g.titleSelected = max(0, g.titleSelected-1)
	}
	if isDownPressed() {
		g.titleSelected = min(len(titleMenuItems())-1, g.titleSelected+1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if isConfirmPressed() {
		return g.activateTitleMenu()
	}
	return nil
}

func (g *Game) activateTitleMenu() error {
	switch titleMenuItems()[clamp(g.titleSelected, 0, len(titleMenuItems())-1)] {
	case "Continue":
		g.selected = clamp(g.progress.UnlockedLevel-1, 0, max(0, len(g.levels)-1))
		return g.startSelectedLevel()
	case "New Game":
		g.selected = 0
		return g.startSelectedLevel()
	case "Level Map":
		g.mode = modeSelect
	case "Options":
		if g.sounds.ToggleMute() {
			g.message = "Audio muted."
		} else {
			g.message = "Audio on."
		}
	case "Help":
		g.message = "Help: 2/4/6/8 move, 5 action, * recall."
	case "About":
		g.message = "About: Angkor five-stage Diamond Rush remake."
	case "Exit":
		return ebiten.Termination
	}
	return nil
}

func (g *Game) updateSelect() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = modeTitle
		return nil
	}
	if isBuyHealthPressed() {
		return g.buyHealthUpgrade()
	}
	if isBuyArmorPressed() {
		return g.buyArmorUpgrade()
	}
	if isBuyLifePressed() {
		return g.buyLifeUpgrade()
	}
	if isLeftPressed() {
		g.selected = max(0, g.selected-1)
	}
	if isRightPressed() {
		g.selected = min(g.progress.UnlockedLevel-1, g.selected+1)
	}
	if isConfirmPressed() {
		return g.startSelectedLevel()
	}
	return nil
}

func (g *Game) startSelectedLevel() error {
	if !g.canEnterLevel(g.selected) {
		g.message = g.levelGateMessage(g.selected)
		return nil
	}
	if err := g.loadLevel(g.selected); err != nil {
		return err
	}
	g.mode = modePlay
	return nil
}

func (g *Game) updatePlay() error {
	if g.world.Status == world.Playing && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = modePause
		return nil
	}
	if g.world.Status != world.Playing && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = modeSelect
		g.selected = g.levelIndex
		return nil
	}
	if isRetryPressed() {
		return g.restartLevel()
	}
	if g.world.Status == world.Won && isConfirmPressed() {
		next := g.nextLevelAfterClear()
		if next >= len(g.levels) {
			g.mode = modeSelect
			g.selected = g.levelIndex
			g.message = "All stages clear."
			return nil
		}
		if g.canEnterLevel(next) {
			return g.loadLevel(next)
		}
		g.mode = modeSelect
		g.selected = min(g.progress.UnlockedLevel-1, next)
		g.message = g.levelGateMessage(g.selected)
		return nil
	}
	var dir world.Direction
	var actionDir world.Direction
	recall := false
	if g.cooldown > 0 {
		g.cooldown--
	} else if isRecallPressed() {
		recall = true
	} else if isActionPressed() || g.pointerActionPressed() {
		actionDir = g.facing
	} else {
		dir = g.inputDirection()
		if dir != (world.Direction{}) {
			g.facing = dir
			g.cooldown = moveCooldown
		}
	}

	oldX, oldY := g.world.Player.X, g.world.Player.Y
	if recall {
		if g.world.RecallCheckpoint() {
			g.cooldown = moveCooldown
		}
	} else if actionDir != (world.Direction{}) {
		if g.world.UpdateAction(actionDir) {
			g.cooldown = moveCooldown
		}
	} else {
		g.world.Update(dir)
	}
	if oldX != g.world.Player.X || oldY != g.world.Player.Y {
		g.playerAnim = tileAnim{
			active: true,
			fromX:  oldX,
			fromY:  oldY,
			toX:    g.world.Player.X,
			toY:    g.world.Player.Y,
			frames: animFrames,
		}
	}
	if g.playerAnim.active {
		g.playerAnim.frame++
		if g.playerAnim.frame >= g.playerAnim.frames {
			g.playerAnim.active = false
		}
	}
	fatalSoundPlayed := false
	for _, event := range g.world.Events() {
		switch event {
		case world.EventStep, world.EventDig:
			g.sounds.Play(soundStep)
		case world.EventDiamond:
			g.sounds.Play(soundDiamond)
			g.message = "Diamond collected."
		case world.EventRedDiamond:
			g.sounds.Play(soundChest)
			g.message = "Red diamond found."
		case world.EventKey:
			g.sounds.Play(soundKey)
			g.message = "Key collected."
		case world.EventChest:
			g.sounds.Play(soundChest)
			g.message = "Treasure found."
		case world.EventDoor:
			g.sounds.Play(soundDoor)
			g.message = "Door opened."
		case world.EventExitOpen:
			g.sounds.Play(soundSwitch)
			g.message = "Exit is open."
		case world.EventDamage:
			g.sounds.Play(soundBreak)
			g.message = "Energy lost."
		case world.EventPotion:
			g.sounds.Play(soundKey)
			g.message = "Energy restored."
		case world.EventExtraLife:
			g.sounds.Play(soundKey)
			g.message = "Extra life gained."
		case world.EventCompass:
			g.sounds.Play(soundKey)
			g.progress.HasCompass = true
			if err := saveProgress(g.savePath, g.progress); err != nil {
				return err
			}
			g.message = "Compass acquired."
		case world.EventSecretExit:
			g.sounds.Play(soundTeleport)
			g.message = "Secret exit found."
		case world.EventBossHit:
			g.sounds.Play(soundBreak)
			g.message = "Guardian wounded."
		case world.EventBossDefeat:
			g.sounds.Play(soundWin)
			g.message = "Guardian defeated."
		case world.EventDeath:
			g.sounds.Play(soundDeath)
			fatalSoundPlayed = true
			if g.world.Status == world.Lost {
				g.message = "No lives left. Press 9 to restart."
			} else {
				g.message = "Lost a life. Returned to checkpoint."
			}
		case world.EventWin:
			stars, newBest, newBestScore, err := g.recordWin()
			if err != nil {
				return err
			}
			g.clearStars = stars
			g.newBest = newBest
			g.newBestScore = newBestScore
			g.sounds.Play(soundWin)
			if g.world.SecretExitFound {
				g.message = "Secret route clear. Press 5 for next."
			} else {
				g.message = "Stage clear. Press 5 for next."
			}
		case world.EventTrap:
			if g.world.Health > 0 && !fatalSoundPlayed {
				g.message = "Trap. Energy lost."
				break
			}
			if !fatalSoundPlayed {
				g.sounds.Play(soundDeath)
				fatalSoundPlayed = true
			}
			if g.world.Status == world.Lost {
				g.message = "Trap. Press 9 to restart."
			} else {
				g.message = "Trap. Returned to checkpoint."
			}
		case world.EventSwitch:
			g.sounds.Play(soundSwitch)
			g.message = "Switch opened a path."
		case world.EventBreak:
			g.sounds.Play(soundBreak)
			g.message = "Cracked wall broke."
		case world.EventTeleport:
			g.sounds.Play(soundTeleport)
			g.message = "Teleported."
		case world.EventBurn:
			if g.world.Status == world.Lost {
				if !fatalSoundPlayed {
					g.sounds.Play(soundDeath)
					fatalSoundPlayed = true
				}
				g.message = "Burned. Press 9 to restart."
			} else if fatalSoundPlayed {
				g.message = "Burned. Returned to checkpoint."
			} else if g.world.Health > 0 {
				g.message = "Burned. Energy lost."
			} else {
				g.sounds.Play(soundBreak)
				g.message = "Lava swallowed a falling object."
			}
		case world.EventHook:
			g.sounds.Play(soundSwitch)
			g.message = "Hook pulled."
		case world.EventHammer:
			g.sounds.Play(soundBreak)
			g.message = "Hammer struck."
		case world.EventToolHammer:
			g.sounds.Play(soundKey)
			g.progress.HasHammer = true
			if err := saveProgress(g.savePath, g.progress); err != nil {
				return err
			}
			g.message = "Mystic Hammer acquired."
		case world.EventToolHook:
			g.sounds.Play(soundKey)
			g.progress.HasHook = true
			if err := saveProgress(g.savePath, g.progress); err != nil {
				return err
			}
			g.message = "Mystic Hook acquired."
		case world.EventCheckpoint:
			g.sounds.Play(soundSwitch)
			g.message = "Checkpoint reached."
		case world.EventRecall:
			g.sounds.Play(soundTeleport)
			g.message = "Returned to checkpoint."
		case world.EventReset:
			g.sounds.Play(soundTeleport)
			g.message = "Checkpoint room reset."
		case world.EventReveal:
			g.sounds.Play(soundSwitch)
			g.message = "Secret passage revealed."
		}
	}
	return nil
}

func (g *Game) updatePause() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || isConfirmPressed() {
		g.mode = modePlay
		return nil
	}
	if isRetryPressed() {
		g.mode = modePlay
		return g.restartLevel()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.mode = modeSelect
		g.selected = g.levelIndex
		return nil
	}
	return nil
}

func (g *Game) recordWin() (int, bool, bool, error) {
	if g.winSaved {
		return g.clearStars, g.newBest, g.newBestScore, nil
	}
	g.winSaved = true
	normalizeProgress(&g.progress, len(g.levels))
	best := g.progress.BestSteps[g.levelIndex]
	newBest := best == 0 || g.world.Steps < best
	if best == 0 || g.world.Steps < best {
		g.progress.BestSteps[g.levelIndex] = g.world.Steps
	}
	g.clearScore = completionScore(g.world.Score, g.world.Steps, g.parSteps[g.levelIndex])
	newBestScore := g.clearScore > g.progress.BestScores[g.levelIndex]
	if newBestScore {
		g.progress.BestScores[g.levelIndex] = g.clearScore
	}
	if g.world.RedDiamonds > g.progress.RedDiamonds[g.levelIndex] {
		g.progress.RedDiamonds[g.levelIndex] = g.world.RedDiamonds
	}
	if g.world.Diamonds > g.progress.PurpleDiamonds[g.levelIndex] {
		g.progress.PurpleBank += g.world.Diamonds - g.progress.PurpleDiamonds[g.levelIndex]
		g.progress.PurpleDiamonds[g.levelIndex] = g.world.Diamonds
	}
	if g.world.SecretExitFound {
		g.progress.SecretExits[g.levelIndex] = true
	}
	if g.world.TotalDiamonds == 0 || g.world.Diamonds >= g.world.TotalDiamonds {
		g.progress.AllPurpleClears[g.levelIndex] = true
	}
	if g.world.TotalRedDiamonds == 0 || g.world.RedDiamonds >= g.world.TotalRedDiamonds {
		g.progress.AllRedClears[g.levelIndex] = true
	}
	if !g.world.Damaged {
		g.progress.NoDamageClears[g.levelIndex] = true
	}
	if !g.world.RecallUsed {
		g.progress.NoRecallClears[g.levelIndex] = true
	}
	if !g.restartUsed {
		g.progress.NoRestartClears[g.levelIndex] = true
	}
	if g.levelIndex+1 < len(g.levels) && g.progress.UnlockedLevel < g.levelIndex+2 {
		g.progress.UnlockedLevel = g.levelIndex + 2
	}
	if g.world.SecretExitFound {
		if target, ok := secretRouteTarget(g.levelIndex); ok && target < len(g.levels) && g.progress.UnlockedLevel < target+1 {
			g.progress.UnlockedLevel = target + 1
		}
	}
	if g.levelIndex == len(g.levels)-1 && totalRedDiamonds(g.progress) >= ancientSealRequirement(len(g.levels)) {
		g.progress.AncientSealOpen = true
	}
	normalizeProgress(&g.progress, len(g.levels))
	stars := starRating(g.progress.BestSteps[g.levelIndex], g.parSteps[g.levelIndex])
	return stars, newBest, newBestScore, saveProgress(g.savePath, g.progress)
}

func (g *Game) nextLevelAfterClear() int {
	if g.world != nil && g.world.SecretExitFound {
		if target, ok := secretRouteTarget(g.levelIndex); ok {
			return target
		}
	}
	return g.levelIndex + 1
}

func (g *Game) loadLevel(index int) error {
	if index < 0 || index >= len(g.levels) {
		return fmt.Errorf("level index %d out of range", index)
	}
	def, err := level.LoadFile(g.levels[index])
	if err != nil {
		return err
	}
	w, err := world.New(def)
	if err != nil {
		return err
	}
	w.Lives = maxLivesForProgress(g.progress)
	w.MaxHealth = maxHealthForProgress(g.progress)
	w.Health = w.MaxHealth
	w.MaxArmor = maxArmorForProgress(g.progress)
	w.Armor = w.MaxArmor
	w.HasCompass = w.HasCompass || g.progress.HasCompass
	w.HasHammer = w.HasHammer || g.progress.HasHammer
	w.HasHook = w.HasHook || g.progress.HasHook
	g.world = w
	g.levelIndex = index
	g.selected = index
	g.cooldown = 0
	g.facing = world.Right
	g.playerAnim = tileAnim{}
	g.winSaved = false
	g.clearStars = 0
	g.clearScore = 0
	g.newBest = false
	g.newBestScore = false
	g.restartUsed = false
	g.hint = def.Hint
	g.message = "Collect diamonds, open the exit."
	return nil
}

func (g *Game) restartLevel() error {
	index := g.levelIndex
	if err := g.loadLevel(index); err != nil {
		return err
	}
	g.restartUsed = true
	return nil
}

func (g *Game) buyHealthUpgrade() error {
	normalizeProgress(&g.progress, len(g.levels))
	cost := maxHealthUpgradeCost(g.progress)
	if g.progress.PurpleBank < cost {
		g.message = fmt.Sprintf("Need %d purple diamonds for HP upgrade.", cost)
		return nil
	}
	g.progress.PurpleBank -= cost
	g.progress.MaxHealthUpgrades++
	if g.world != nil {
		g.world.MaxHealth = maxHealthForProgress(g.progress)
		g.world.Health = g.world.MaxHealth
	}
	g.message = fmt.Sprintf("Max HP upgraded to %d.", maxHealthForProgress(g.progress))
	return saveProgress(g.savePath, g.progress)
}

func (g *Game) buyArmorUpgrade() error {
	normalizeProgress(&g.progress, len(g.levels))
	cost := armorUpgradeCost(g.progress)
	if g.progress.PurpleBank < cost {
		g.message = fmt.Sprintf("Need %d purple diamonds for armor.", cost)
		return nil
	}
	g.progress.PurpleBank -= cost
	g.progress.ArmorUpgrades++
	if g.world != nil {
		g.world.MaxArmor = maxArmorForProgress(g.progress)
		g.world.Armor = g.world.MaxArmor
	}
	g.message = fmt.Sprintf("Armor upgraded to %d.", maxArmorForProgress(g.progress))
	return saveProgress(g.savePath, g.progress)
}

func (g *Game) buyLifeUpgrade() error {
	normalizeProgress(&g.progress, len(g.levels))
	cost := lifeUpgradeCost(g.progress)
	if g.progress.PurpleBank < cost {
		g.message = fmt.Sprintf("Need %d purple diamonds for life.", cost)
		return nil
	}
	g.progress.PurpleBank -= cost
	g.progress.LifeUpgrades++
	if g.world != nil {
		g.world.Lives = maxLivesForProgress(g.progress)
	}
	g.message = fmt.Sprintf("Lives upgraded to %d.", maxLivesForProgress(g.progress))
	return saveProgress(g.savePath, g.progress)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 18, 26, 255})
	switch g.mode {
	case modeTitle:
		g.drawMenuBackdrop(screen)
		g.drawTitle(screen)
	case modeSelect:
		g.drawMenuBackdrop(screen)
		g.drawLevelSelect(screen)
	case modePause:
		g.drawWorld(screen)
		g.drawHUD(screen)
		g.drawPauseOverlay(screen)
	default:
		g.drawWorld(screen)
		g.drawHUD(screen)
		g.drawVirtualDPad(screen)
		g.drawStatusOverlay(screen)
	}
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	originX, originY := g.cameraOrigin()
	for y := originY; y < min(g.world.Height, originY+viewportTilesY); y++ {
		for x := originX; x < min(g.world.Width, originX+viewportTilesX); x++ {
			tile := g.world.TileAt(x, y)
			if tile == world.HiddenWall {
				tile = world.Wall
			}
			if tile == world.TimedSpike && g.world.TimedSpikeActive() {
				tile = world.Spike
			}
			if tile == world.FireTrap && g.world.FireTrapActive() {
				tile = world.Lava
			}
			img := g.assets.Tile(tile)
			if tile == world.Diamond {
				img = g.assets.DiamondFrame(g.tick)
			}
			drawImage(screen, img, (x-originX)*tileSize, (y-originY)*tileSize)
		}
	}
	for _, enemy := range g.world.Enemies {
		if enemy.X < originX || enemy.Y < originY || enemy.X >= originX+viewportTilesX || enemy.Y >= originY+viewportTilesY {
			continue
		}
		drawImage(screen, g.assets.EnemyImage(enemy.Type, g.tick), (enemy.X-originX)*tileSize, (enemy.Y-originY)*tileSize)
	}
	for _, boss := range g.world.Bosses {
		if boss.X < originX || boss.Y < originY || boss.X >= originX+viewportTilesX || boss.Y >= originY+viewportTilesY {
			continue
		}
		drawImage(screen, g.assets.EnemyImage(world.EnemyChaser, g.tick), (boss.X-originX)*tileSize, (boss.Y-originY)*tileSize)
	}
	playerX, playerY := g.playerDrawPosition()
	drawImage(screen, g.assets.PlayerFrame(g.tick, g.playerAnim.active), playerX-originX*tileSize, playerY-originY*tileSize)
}

func (g *Game) drawMenuBackdrop(screen *ebiten.Image) {
	width, height := g.Layout(0, 0)
	if g.assets.MenuBackdrop != nil {
		drawCoverImage(screen, g.assets.MenuBackdrop, width, height)
		ebitenutil.DrawRect(screen, 0, 0, float64(width), 86, color.RGBA{5, 8, 14, 122})
		ebitenutil.DrawRect(screen, 0, float64(height-86), float64(width), 86, color.RGBA{5, 8, 14, 156})
		return
	}
	playfieldHeight := viewportPixelHeight()
	screen.Fill(color.RGBA{13, 18, 27, 255})
	ebitenutil.DrawRect(screen, 0, float64(playfieldHeight), float64(width), float64(height-playfieldHeight), color.RGBA{22, 25, 31, 255})
	ebitenutil.DrawRect(screen, 0, float64(playfieldHeight-34), float64(width), 34, color.RGBA{26, 32, 41, 255})
	ebitenutil.DrawRect(screen, 0, float64(playfieldHeight-4), float64(width), 4, color.RGBA{172, 130, 62, 255})

	moon := color.RGBA{226, 210, 156, 255}
	for y := 0; y < 22; y++ {
		for x := 0; x < 22; x++ {
			dx := x - 11
			dy := y - 11
			if dx*dx+dy*dy <= 100 {
				ebitenutil.DrawRect(screen, float64(365+x), float64(38+y), 1, 1, moon)
			}
		}
	}

	stone := color.RGBA{76, 71, 61, 255}
	shadow := color.RGBA{42, 40, 38, 255}
	trim := color.RGBA{162, 125, 64, 255}
	ebitenutil.DrawRect(screen, 92, 106, 296, 18, trim)
	ebitenutil.DrawRect(screen, 110, 124, 260, 132, stone)
	ebitenutil.DrawRect(screen, 132, 146, 44, 110, shadow)
	ebitenutil.DrawRect(screen, 218, 146, 44, 110, shadow)
	ebitenutil.DrawRect(screen, 304, 146, 44, 110, shadow)
	ebitenutil.DrawRect(screen, 82, 256, 316, 20, trim)
	ebitenutil.DrawRect(screen, 150, 188, 180, 68, color.RGBA{18, 21, 27, 255})

	for x := 68; x < width; x += 82 {
		drawImage(screen, g.assets.Tile(world.Diamond), x, 290+(x/82)%2*12)
	}
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	hudY := viewportPixelHeight() + 6
	best := "-"
	rating := "---"
	if g.progress.BestSteps[g.levelIndex] > 0 {
		best = fmt.Sprintf("%d", g.progress.BestSteps[g.levelIndex])
		rating = ratingText(starRating(g.progress.BestSteps[g.levelIndex], g.parSteps[g.levelIndex]))
	}
	hook := "-"
	if g.world.HasHook {
		hook = "Y"
	}
	hammer := "-"
	if g.world.HasHammer {
		hammer = "Y"
	}
	compass := "--"
	if dir, distance, ok := g.world.CompassToCheckpoint(); ok && g.world.HasCompass {
		compass = fmt.Sprintf("%s%d", compassDirectionText(dir), distance)
	}
	status := fmt.Sprintf("L%d/%d Life %d HP %d/%d Ar %d/%d Dia %d/%d R %d K %d G %d", g.levelIndex+1, len(g.levels), g.world.Lives, g.world.Health, g.world.MaxHealth, g.world.Armor, g.world.MaxArmor, g.world.Diamonds, g.world.RequiredDiamonds, g.currentRedDiamonds(), g.world.Keys, g.world.GoldKeys)
	if hp, maxHP, ok := g.world.BossHealth(); ok {
		status = fmt.Sprintf("%s Boss %d/%d", status, hp, maxHP)
	}
	tools := fmt.Sprintf("Cmp %s Hmr %s Hook %s Score %d Step %d Best %s %s", compass, hammer, hook, g.world.Score, g.world.Steps, best, rating)
	ebitenutil.DebugPrintAt(screen, status, 8, hudY)
	ebitenutil.DebugPrintAt(screen, tools, 8, hudY+16)
	ebitenutil.DebugPrintAt(screen, g.message, 8, hudY+32)
	if g.world.Steps == 0 && g.hint != "" {
		ebitenutil.DebugPrintAt(screen, g.hint, 8, hudY+48)
		return
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Move 2/4/6/8  Action 5  Recall Back/*  Audio 0 %s", audioState(g.sounds)), 8, hudY+48)
}

func (g *Game) drawVirtualDPad(screen *ebiten.Image) {
	width, height := g.Layout(0, 0)
	rect := virtualDPadRect(width, height)
	button := rect.Dx() / 3
	base := color.RGBA{34, 39, 50, 210}
	lit := color.RGBA{78, 89, 110, 230}
	outline := color.RGBA{170, 184, 202, 255}
	cells := []struct {
		col   int
		row   int
		label string
	}{
		{1, 0, "2"},
		{0, 1, "4"},
		{1, 1, "5"},
		{2, 1, "6"},
		{1, 2, "8"},
	}
	for _, cell := range cells {
		x := float64(rect.Min.X + cell.col*button)
		y := float64(rect.Min.Y + cell.row*button)
		ebitenutil.DrawRect(screen, x, y, float64(button-1), float64(button-1), base)
		ebitenutil.DrawRect(screen, x, y, float64(button-1), 1, outline)
		ebitenutil.DrawRect(screen, x, y+float64(button-2), float64(button-1), 1, lit)
		ebitenutil.DebugPrintAt(screen, cell.label, int(x)+8, int(y)+6)
	}
}

func (g *Game) drawStatusOverlay(screen *ebiten.Image) {
	if g.world.Status == world.Playing {
		return
	}
	screenWidth, _ := g.Layout(0, 0)
	y := 146.0
	w := 296.0
	h := 142.0
	if g.world.Status == world.Won && g.levelIndex == len(g.levels)-1 {
		h = 158
	}
	x := float64((screenWidth - int(w)) / 2)
	ebitenutil.DrawRect(screen, x, y, w, h, color.RGBA{7, 9, 12, 214})
	ebitenutil.DrawRect(screen, x, y, w, 2, color.RGBA{200, 210, 220, 255})
	ebitenutil.DrawRect(screen, x, y+h-2, w, 2, color.RGBA{80, 90, 105, 255})
	if g.world.Status == world.Won {
		best := g.progress.BestSteps[g.levelIndex]
		bestScore := g.progress.BestScores[g.levelIndex]
		ebitenutil.DebugPrintAt(screen, "STAGE CLEAR", int(x)+24, int(y)+18)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rating %s  Score %d  Best %d", ratingText(g.clearStars), g.clearScore, bestScore), int(x)+24, int(y)+42)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Steps %d  Best %d", g.world.Steps, best), int(x)+24, int(y)+62)
		if g.newBest || g.newBestScore {
			ebitenutil.DebugPrintAt(screen, "New record", int(x)+176, int(y)+62)
		}
		ebitenutil.DebugPrintAt(screen, g.clearGemLine(), int(x)+24, int(y)+82)
		ebitenutil.DebugPrintAt(screen, g.clearMarkLine(), int(x)+24, int(y)+102)
		if g.levelIndex == len(g.levels)-1 {
			ebitenutil.DebugPrintAt(screen, g.sealLine(), int(x)+24, int(y)+122)
			ebitenutil.DebugPrintAt(screen, "5: map  9: retry  Esc: select", int(x)+24, int(y)+140)
			return
		}
		ebitenutil.DebugPrintAt(screen, "5: next  9: retry  Esc: select", int(x)+24, int(y)+122)
		return
	}
	ebitenutil.DebugPrintAt(screen, "TRY AGAIN", int(x)+24, int(y)+22)
	ebitenutil.DebugPrintAt(screen, g.message, int(x)+24, int(y)+48)
	ebitenutil.DebugPrintAt(screen, "9: retry  Esc: level select", int(x)+24, int(y)+82)
}

func (g *Game) clearGemLine() string {
	if g.world == nil {
		return "Gems P -/-  R -/-"
	}
	return fmt.Sprintf("Gems P %d/%d  R %d/%d", g.world.Diamonds, g.world.TotalDiamonds, g.world.RedDiamonds, g.world.TotalRedDiamonds)
}

func (g *Game) clearMarkLine() string {
	if g.world == nil {
		return "Marks Sec -  All -/-  Clean -"
	}
	secret := markText(g.world.SecretExitFound)
	allPurple := markText(g.world.TotalDiamonds == 0 || g.world.Diamonds >= g.world.TotalDiamonds)
	allRed := markText(g.world.TotalRedDiamonds == 0 || g.world.RedDiamonds >= g.world.TotalRedDiamonds)
	clean := markText(!g.world.Damaged && !g.world.RecallUsed && !g.restartUsed)
	return fmt.Sprintf("Marks Sec %s  All %s/%s  Clean %s", secret, allPurple, allRed, clean)
}

func markText(ok bool) string {
	if ok {
		return "Y"
	}
	return "-"
}

func (g *Game) sealLine() string {
	return sealStatusText(g.progress, len(g.levels))
}

func (g *Game) drawPauseOverlay(screen *ebiten.Image) {
	screenWidth, _ := g.Layout(0, 0)
	y := 144.0
	w := 296.0
	h := 104.0
	x := float64((screenWidth - int(w)) / 2)
	ebitenutil.DrawRect(screen, x, y, w, h, color.RGBA{7, 9, 12, 220})
	ebitenutil.DrawRect(screen, x, y, w, 2, color.RGBA{200, 210, 220, 255})
	ebitenutil.DrawRect(screen, x, y+h-2, w, 2, color.RGBA{80, 90, 105, 255})
	ebitenutil.DebugPrintAt(screen, "PAUSED", int(x)+24, int(y)+18)
	ebitenutil.DebugPrintAt(screen, "5/Enter/Esc: resume", int(x)+24, int(y)+44)
	ebitenutil.DebugPrintAt(screen, "9: retry   Backspace: select", int(x)+24, int(y)+68)
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	x, y := 122, 48
	ebitenutil.DebugPrintAt(screen, "DIAMOND RUSH", x, y)
	ebitenutil.DebugPrintAt(screen, "ANGKOR TRIAL", x+16, y+22)
	menuX, menuY := 150, 144
	for i, label := range titleMenuItems() {
		prefix := "  "
		if i == g.titleSelected {
			prefix = "> "
		}
		ebitenutil.DebugPrintAt(screen, prefix+label, menuX, menuY+i*18)
	}
	ebitenutil.DebugPrintAt(screen, "2/8 choose  5 select  Esc exit", 106, 286)
	ebitenutil.DebugPrintAt(screen, g.message, 74, 306)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Purple %d  Red %d  Life %d  HP %d  Ar %d  Tools %s", g.progress.PurpleBank, totalRedDiamonds(g.progress), maxLivesForProgress(g.progress), maxHealthForProgress(g.progress), maxArmorForProgress(g.progress), progressToolsText(g.progress)), 68, 382)
	ebitenutil.DebugPrintAt(screen, sealStatusText(g.progress, len(g.levels)), 148, 400)
}

func titleMenuItems() []string {
	return []string{"Continue", "New Game", "Level Map", "Options", "Help", "About", "Exit"}
}

func (g *Game) drawLevelSelect(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "ANGKOR MAP", 42, 34)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Purple %d  Red %d  Tools %s", g.progress.PurpleBank, totalRedDiamonds(g.progress), progressToolsText(g.progress)), 42, 52)
	ebitenutil.DebugPrintAt(screen, sealStatusText(g.progress, len(g.levels)), 292, 52)
	g.drawWorldMap(screen)
	g.drawLevelSelectDetails(screen)
}

func (g *Game) drawWorldMap(screen *ebiten.Image) {
	nodes := levelMapNodes(len(g.levels))
	for i := 1; i < len(nodes); i++ {
		pathColor := color.RGBA{88, 76, 55, 255}
		if g.canEnterLevel(i) {
			pathColor = color.RGBA{214, 174, 86, 255}
		}
		if i < g.progress.UnlockedLevel && g.hasSecretRouteForLevel(i) && !g.hasRedDiamondsForLevel(i) {
			pathColor = color.RGBA{86, 194, 180, 255}
		} else if i < g.progress.UnlockedLevel && !g.hasRedDiamondsForLevel(i) {
			pathColor = color.RGBA{126, 42, 46, 255}
		}
		drawThickLine(screen, nodes[i-1].X+16, nodes[i-1].Y+16, nodes[i].X+16, nodes[i].Y+16, pathColor)
	}
	g.drawSecretRoutes(screen, nodes)
	for i, node := range nodes {
		g.drawWorldMapNode(screen, i, node)
	}
}

func (g *Game) drawSecretRoutes(screen *ebiten.Image, nodes []image.Point) {
	for from := range g.levels {
		to, ok := secretRouteTarget(from)
		if !ok || from >= len(nodes) || to >= len(nodes) || from >= len(g.progress.SecretExits) || !g.progress.SecretExits[from] {
			continue
		}
		start := nodes[from]
		end := nodes[to]
		drawDashedLine(screen, start.X+16, start.Y+16, end.X+16, end.Y-10, color.RGBA{86, 194, 180, 255})
		ebitenutil.DebugPrintAt(screen, "SECRET", (start.X+end.X)/2-16, min(start.Y, end.Y)-18)
	}
}

func (g *Game) drawWorldMapNode(screen *ebiten.Image, i int, node image.Point) {
	state := g.levelMapState(i)
	border := color.RGBA{85, 73, 55, 255}
	if i == g.selected {
		border = color.RGBA{255, 232, 132, 255}
	}
	ebitenutil.DrawRect(screen, float64(node.X-4), float64(node.Y-4), 40, 40, border)
	ebitenutil.DrawRect(screen, float64(node.X-2), float64(node.Y-2), 36, 36, color.RGBA{25, 23, 19, 235})
	switch state {
	case "open":
		drawImage(screen, g.assets.Tile(world.ExitOpen), node.X, node.Y)
	case "secret":
		drawImage(screen, g.assets.Tile(world.SecretExit), node.X, node.Y)
		ebitenutil.DrawRect(screen, float64(node.X+21), float64(node.Y+3), 8, 8, color.RGBA{86, 194, 180, 255})
	case "locked":
		drawImage(screen, g.assets.Tile(world.Door), node.X, node.Y)
	default:
		drawImage(screen, g.assets.Tile(world.Door), node.X, node.Y)
		ebitenutil.DrawRect(screen, float64(node.X+21), float64(node.Y+3), 8, 8, color.RGBA{210, 44, 54, 255})
	}
	label := fmt.Sprintf("%d", i+1)
	if i < len(g.progress.SecretExits) && g.progress.SecretExits[i] {
		label += "S"
	}
	ebitenutil.DebugPrintAt(screen, label, node.X+10, node.Y+38)
}

func (g *Game) drawLevelSelectDetails(screen *ebiten.Image) {
	x, y := 42, 322
	i := clamp(g.selected, 0, max(0, len(g.levels)-1))
	state := g.levelMapState(i)
	title := ""
	if i < len(g.titles) {
		title = g.titles[i]
	}
	ebitenutil.DrawRect(screen, float64(x-8), float64(y-10), 406, 88, color.RGBA{8, 9, 12, 210})
	ebitenutil.DrawRect(screen, float64(x-8), float64(y-10), 406, 2, color.RGBA{202, 166, 86, 255})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Stage %d  %s  %s", i+1, title, state), x, y)
	ebitenutil.DebugPrintAt(screen, g.levelSelectLine(i), x, y+16)
	shop := fmt.Sprintf("Shop 1:+HP %d  3:+Ar %d  7:+Life %d  0 audio %s", maxHealthUpgradeCost(g.progress), armorUpgradeCost(g.progress), lifeUpgradeCost(g.progress), audioState(g.sounds))
	ebitenutil.DebugPrintAt(screen, shop, x, y+34)
	ebitenutil.DebugPrintAt(screen, "4/6 choose  5 play  Esc title", x, y+50)
	ebitenutil.DebugPrintAt(screen, g.message, x, y+66)
}

func levelMapNodes(levelCount int) []image.Point {
	base := []image.Point{
		image.Pt(62, 244),
		image.Pt(138, 194),
		image.Pt(226, 226),
		image.Pt(302, 164),
		image.Pt(382, 208),
	}
	if levelCount <= len(base) {
		return base[:levelCount]
	}
	nodes := append([]image.Point{}, base...)
	for len(nodes) < levelCount {
		i := len(nodes)
		nodes = append(nodes, image.Pt(58+(i%5)*80, 150+(i/5)*58))
	}
	return nodes
}

func drawThickLine(screen *ebiten.Image, x1, y1, x2, y2 int, c color.Color) {
	for offset := -1; offset <= 1; offset++ {
		ebitenutil.DrawLine(screen, float64(x1), float64(y1+offset), float64(x2), float64(y2+offset), c)
	}
}

func drawDashedLine(screen *ebiten.Image, x1, y1, x2, y2 int, c color.Color) {
	segments := max(abs(x2-x1), abs(y2-y1)) / 8
	if segments < 1 {
		segments = 1
	}
	for i := 0; i < segments; i += 2 {
		t1 := float64(i) / float64(segments)
		t2 := float64(min(i+1, segments)) / float64(segments)
		sx := float64(x1) + float64(x2-x1)*t1
		sy := float64(y1) + float64(y2-y1)*t1
		ex := float64(x1) + float64(x2-x1)*t2
		ey := float64(y1) + float64(y2-y1)*t2
		ebitenutil.DrawLine(screen, sx, sy, ex, ey, c)
	}
}

func (g *Game) levelSelectLine(i int) string {
	prefix := "  "
	if i == g.selected {
		prefix = "> "
	}
	state := g.levelMapState(i)
	bestStep := "-"
	bestScore := "-"
	rating := "---"
	if i < len(g.progress.BestSteps) && g.progress.BestSteps[i] > 0 {
		bestStep = fmt.Sprintf("%d", g.progress.BestSteps[i])
		rating = ratingText(starRating(g.progress.BestSteps[i], g.parSteps[i]))
	}
	if i < len(g.progress.BestScores) && g.progress.BestScores[i] > 0 {
		bestScore = fmt.Sprintf("%d", g.progress.BestScores[i])
	}
	red := "-"
	if i < len(g.progress.RedDiamonds) && g.progress.RedDiamonds[i] > 0 {
		red = fmt.Sprintf("%d", g.progress.RedDiamonds[i])
	}
	secret := "-"
	if i < len(g.progress.SecretExits) && g.progress.SecretExits[i] {
		secret = "Y"
	}
	clean := "-"
	if i < len(g.progress.NoDamageClears) && i < len(g.progress.NoRecallClears) && i < len(g.progress.NoRestartClears) && g.progress.NoDamageClears[i] && g.progress.NoRecallClears[i] && g.progress.NoRestartClears[i] {
		clean = "Y"
	}
	all := "-"
	if i < len(g.progress.AllPurpleClears) && i < len(g.progress.AllRedClears) && g.progress.AllPurpleClears[i] && g.progress.AllRedClears[i] {
		all = "Y"
	}
	return fmt.Sprintf("%s%d %s %s st%s sc%s r%s s%s a%s c%s %s", prefix, i+1, g.titles[i], state, bestStep, bestScore, red, secret, all, clean, rating)
}

func (g *Game) levelMapState(i int) string {
	if i >= g.progress.UnlockedLevel {
		return "locked"
	}
	if !g.hasRedDiamondsForLevel(i) {
		if g.hasSecretRouteForLevel(i) {
			return "secret"
		}
		return fmt.Sprintf("red%d", redDiamondRequirement(i))
	}
	return "open"
}

func secretRouteTarget(levelIndex int) (int, bool) {
	switch levelIndex {
	case 3:
		return 4, true
	default:
		return 0, false
	}
}

func (g *Game) canEnterLevel(index int) bool {
	return index >= 0 && index < len(g.levels) && index < g.progress.UnlockedLevel && g.hasEntryRequirementForLevel(index)
}

func (g *Game) hasRedDiamondsForLevel(index int) bool {
	return totalRedDiamonds(g.progress) >= redDiamondRequirement(index)
}

func (g *Game) hasEntryRequirementForLevel(index int) bool {
	return g.hasRedDiamondsForLevel(index) || g.hasSecretRouteForLevel(index)
}

func (g *Game) hasSecretRouteForLevel(index int) bool {
	for from, found := range g.progress.SecretExits {
		if !found {
			continue
		}
		if target, ok := secretRouteTarget(from); ok && target == index {
			return true
		}
	}
	return false
}

func (g *Game) levelGateMessage(index int) string {
	if index < 0 || index >= len(g.levels) || index >= g.progress.UnlockedLevel {
		return "Level locked."
	}
	need := redDiamondRequirement(index)
	have := totalRedDiamonds(g.progress)
	if have < need {
		return fmt.Sprintf("Need %d red diamonds.", need)
	}
	return "Level locked."
}

func (g *Game) currentRedDiamonds() int {
	if g.world == nil {
		return totalRedDiamonds(g.progress)
	}
	total := totalRedDiamonds(g.progress)
	if g.levelIndex >= 0 && g.levelIndex < len(g.progress.RedDiamonds) && g.world.RedDiamonds > g.progress.RedDiamonds[g.levelIndex] {
		total += g.world.RedDiamonds - g.progress.RedDiamonds[g.levelIndex]
	}
	return total
}

func audioState(sounds *Sounds) string {
	if sounds.Muted() {
		return "off"
	}
	return "on"
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return viewportPixelWidth(), viewportPixelHeight() + hudHeight
}

func (g *Game) inputDirection() world.Direction {
	if dir := keyboardDirection(); dir != (world.Direction{}) {
		return dir
	}
	width, height := g.Layout(0, 0)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if dir := directionFromVirtualPad(x, y, width, height); dir != (world.Direction{}) {
			return dir
		}
	}
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if dir := directionFromVirtualPad(x, y, width, height); dir != (world.Direction{}) {
			return dir
		}
	}
	return world.Direction{}
}

func (g *Game) pointerActionPressed() bool {
	width, height := g.Layout(0, 0)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return hookFromVirtualPad(x, y, width, height)
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if hookFromVirtualPad(x, y, width, height) {
			return true
		}
	}
	return false
}

func keyboardDirection() world.Direction {
	switch {
	case isUpPressed():
		return world.Up
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown), ebiten.IsKeyPressed(ebiten.KeyDigit8), ebiten.IsKeyPressed(ebiten.KeyNumpad8):
		return world.Down
	case ebiten.IsKeyPressed(ebiten.KeyArrowLeft), ebiten.IsKeyPressed(ebiten.KeyDigit4), ebiten.IsKeyPressed(ebiten.KeyNumpad4):
		return world.Left
	case ebiten.IsKeyPressed(ebiten.KeyArrowRight), ebiten.IsKeyPressed(ebiten.KeyDigit6), ebiten.IsKeyPressed(ebiten.KeyNumpad6):
		return world.Right
	default:
		return world.Direction{}
	}
}

func directionFromVirtualPad(x, y, width, height int) world.Direction {
	rect := virtualDPadRect(width, height)
	if !image.Pt(x, y).In(rect) {
		return world.Direction{}
	}
	button := rect.Dx() / 3
	col := (x - rect.Min.X) / button
	row := (y - rect.Min.Y) / button
	switch {
	case col == 1 && row == 0:
		return world.Up
	case col == 1 && row == 2:
		return world.Down
	case col == 0 && row == 1:
		return world.Left
	case col == 2 && row == 1:
		return world.Right
	default:
		return world.Direction{}
	}
}

func hookFromVirtualPad(x, y, width, height int) bool {
	rect := virtualDPadRect(width, height)
	if !image.Pt(x, y).In(rect) {
		return false
	}
	button := rect.Dx() / 3
	col := (x - rect.Min.X) / button
	row := (y - rect.Min.Y) / button
	return col == 1 && row == 1
}

func virtualDPadRect(width, height int) image.Rectangle {
	return image.Rect(width-dpadSize-dpadMargin, height-dpadSize-dpadMargin, width-dpadMargin, height-dpadMargin)
}

func (g *Game) cameraOrigin() (int, int) {
	maxX := max(0, g.world.Width-viewportTilesX)
	maxY := max(0, g.world.Height-viewportTilesY)
	x := clamp(g.world.Player.X-viewportTilesX/2, 0, maxX)
	y := clamp(g.world.Player.Y-viewportTilesY/2, 0, maxY)
	return x, y
}

func viewportPixelWidth() int {
	return viewportTilesX * tileSize
}

func viewportPixelHeight() int {
	return viewportTilesY * tileSize
}

func isConfirmPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyDigit5) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad5)
}

func isActionPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyShift) ||
		inpututil.IsKeyJustPressed(ebiten.KeyDigit5) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad5)
}

func compassDirectionText(dir world.Direction) string {
	switch {
	case dir.DY < 0 && dir.DX < 0:
		return "NW"
	case dir.DY < 0 && dir.DX > 0:
		return "NE"
	case dir.DY > 0 && dir.DX < 0:
		return "SW"
	case dir.DY > 0 && dir.DX > 0:
		return "SE"
	case dir.DY < 0:
		return "N"
	case dir.DY > 0:
		return "S"
	case dir.DX < 0:
		return "W"
	case dir.DX > 0:
		return "E"
	default:
		return "-"
	}
}

func progressToolsText(progress Progress) string {
	compass := "-"
	if progress.HasCompass {
		compass = "C"
	}
	hammer := "-"
	if progress.HasHammer {
		hammer = "H"
	}
	hook := "-"
	if progress.HasHook {
		hook = "K"
	}
	return compass + hammer + hook
}

func isRecallPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyBackspace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpadMultiply)
}

func isMenuUpPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyDigit2) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad2)
}

func isUpPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
		ebiten.IsKeyPressed(ebiten.KeyDigit2) ||
		ebiten.IsKeyPressed(ebiten.KeyNumpad2)
}

func isDownPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyDigit8) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad8)
}

func isLeftPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyDigit4) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad4)
}

func isRightPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyDigit6) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad6)
}

func isBuyHealthPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDigit1) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad1)
}

func isBuyArmorPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDigit3) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad3)
}

func isBuyLifePressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDigit7) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad7)
}

func isRetryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDigit9) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad9)
}

func isAudioPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDigit0) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad0)
}

func (g *Game) playerDrawPosition() (int, int) {
	if !g.playerAnim.active || g.playerAnim.frames <= 0 {
		return g.world.Player.X * tileSize, g.world.Player.Y * tileSize
	}
	t := float64(g.playerAnim.frame) / float64(g.playerAnim.frames)
	x := lerp(float64(g.playerAnim.fromX*tileSize), float64(g.playerAnim.toX*tileSize), t)
	y := lerp(float64(g.playerAnim.fromY*tileSize), float64(g.playerAnim.toY*tileSize), t)
	return int(x + 0.5), int(y + 0.5)
}

type Assets struct {
	tiles         map[world.Tile]*ebiten.Image
	MenuBackdrop  *ebiten.Image
	Player        *ebiten.Image
	PlayerFrames  []*ebiten.Image
	Enemy         *ebiten.Image
	EnemyFrames   map[world.EnemyType][]*ebiten.Image
	DiamondFrames []*ebiten.Image
	EnemyVertical *ebiten.Image
	EnemyChaser   *ebiten.Image
}

func NewAssets(size int) *Assets {
	a := &Assets{tiles: make(map[world.Tile]*ebiten.Image)}
	if img, err := loadPNG("assets/sprites/tileset.png"); err == nil {
		a.loadTileset(img, size)
	} else {
		a.makeFallback(size)
	}
	if img, err := loadPNG("assets/sprites/menu-background.png"); err == nil {
		a.MenuBackdrop = ebiten.NewImageFromImage(img)
	}
	return a
}

func (a *Assets) Tile(tile world.Tile) *ebiten.Image {
	if img := a.tiles[tile]; img != nil {
		return img
	}
	return a.tiles[world.Empty]
}

func (a *Assets) PlayerFrame(tick int, moving bool) *ebiten.Image {
	if moving && len(a.PlayerFrames) > 0 {
		return a.PlayerFrames[(tick/5)%len(a.PlayerFrames)]
	}
	return a.Player
}

func (a *Assets) DiamondFrame(tick int) *ebiten.Image {
	if len(a.DiamondFrames) == 0 {
		return a.Tile(world.Diamond)
	}
	return a.DiamondFrames[(tick/10)%len(a.DiamondFrames)]
}

func (a *Assets) EnemyImage(enemyType world.EnemyType, tick int) *ebiten.Image {
	if frames := a.EnemyFrames[enemyType]; len(frames) > 0 {
		return frames[(tick/8)%len(frames)]
	}
	return a.Enemy
}

func (a *Assets) loadTileset(src image.Image, size int) {
	source := ebiten.NewImageFromImage(src)
	for i, tile := range []world.Tile{world.Wall, world.Dirt, world.Diamond, world.Rock, world.ExitClosed, world.Key, world.Door, world.ExitOpen, world.Spike, world.Switch, world.Bridge, world.CrackedWall, world.Teleporter, world.Lava, world.Chest, world.Checkpoint, world.Potion, world.HammerPickup, world.HookPickup, world.CompassPickup, world.SecretExit, world.HiddenWall, world.TimedSpike, world.FireTrap, world.GoldKey, world.GoldDoor} {
		rect := image.Rect(i*size, 0, (i+1)*size, size)
		if rect.Max.X <= source.Bounds().Dx() {
			a.tiles[tile] = source.SubImage(rect).(*ebiten.Image)
		}
	}
	if a.tiles[world.HiddenWall] == nil {
		a.tiles[world.HiddenWall] = a.tiles[world.Wall]
	}
	if a.tiles[world.TimedSpike] == nil {
		a.tiles[world.TimedSpike] = timedSpike(size)
	}
	if a.tiles[world.FireTrap] == nil {
		a.tiles[world.FireTrap] = fireTrap(size)
	}
	if a.tiles[world.GoldKey] == nil {
		a.tiles[world.GoldKey] = goldKey(size)
	}
	if a.tiles[world.GoldDoor] == nil {
		a.tiles[world.GoldDoor] = block(size, color.RGBA{132, 96, 24, 255}, color.RGBA{232, 184, 70, 255})
	}
	a.tiles[world.Empty] = solid(size, color.RGBA{20, 24, 31, 255})
	if a.tiles[world.Chest] == nil {
		a.tiles[world.Chest] = chest(size)
	}
	if a.tiles[world.Checkpoint] == nil {
		a.tiles[world.Checkpoint] = checkpoint(size)
	}
	if a.tiles[world.Potion] == nil {
		a.tiles[world.Potion] = potion(size)
	}
	if a.tiles[world.HammerPickup] == nil {
		a.tiles[world.HammerPickup] = hammerPickup(size)
	}
	if a.tiles[world.HookPickup] == nil {
		a.tiles[world.HookPickup] = hookPickup(size)
	}
	if a.tiles[world.CompassPickup] == nil {
		a.tiles[world.CompassPickup] = compassPickup(size)
	}
	if a.tiles[world.SecretExit] == nil {
		a.tiles[world.SecretExit] = secretExit(size)
	}
	a.Player = source.SubImage(image.Rect(0, size, size, size*2)).(*ebiten.Image)
	a.Enemy = source.SubImage(image.Rect(size, size, size*2, size*2)).(*ebiten.Image)
	a.EnemyVertical = tintImage(a.Enemy, color.RGBA{55, 176, 220, 255})
	a.EnemyChaser = tintImage(a.Enemy, color.RGBA{226, 69, 54, 255})
	a.buildAnimationFrames()
}

func (a *Assets) makeFallback(size int) {
	a.tiles[world.Empty] = solid(size, color.RGBA{20, 24, 31, 255})
	a.tiles[world.Wall] = block(size, color.RGBA{75, 80, 93, 255}, color.RGBA{50, 54, 65, 255})
	a.tiles[world.Dirt] = block(size, color.RGBA{118, 84, 48, 255}, color.RGBA{82, 59, 38, 255})
	a.tiles[world.Diamond] = diamond(size, color.RGBA{70, 220, 255, 255}, color.RGBA{235, 255, 255, 255})
	a.tiles[world.Rock] = rock(size)
	a.tiles[world.ExitClosed] = block(size, color.RGBA{36, 64, 92, 255}, color.RGBA{20, 34, 50, 255})
	a.tiles[world.ExitOpen] = block(size, color.RGBA{54, 160, 112, 255}, color.RGBA{190, 255, 192, 255})
	a.tiles[world.Key] = key(size)
	a.tiles[world.Door] = block(size, color.RGBA{132, 85, 38, 255}, color.RGBA{204, 145, 68, 255})
	a.tiles[world.GoldKey] = goldKey(size)
	a.tiles[world.GoldDoor] = block(size, color.RGBA{132, 96, 24, 255}, color.RGBA{232, 184, 70, 255})
	a.tiles[world.Spike] = spike(size)
	a.tiles[world.Switch] = lever(size)
	a.tiles[world.Bridge] = block(size, color.RGBA{54, 88, 105, 255}, color.RGBA{105, 151, 169, 255})
	a.tiles[world.CrackedWall] = crackedWall(size)
	a.tiles[world.Teleporter] = teleporter(size)
	a.tiles[world.Lava] = lava(size)
	a.tiles[world.Chest] = chest(size)
	a.tiles[world.Checkpoint] = checkpoint(size)
	a.tiles[world.Potion] = potion(size)
	a.tiles[world.HammerPickup] = hammerPickup(size)
	a.tiles[world.HookPickup] = hookPickup(size)
	a.tiles[world.CompassPickup] = compassPickup(size)
	a.tiles[world.SecretExit] = secretExit(size)
	a.tiles[world.HiddenWall] = a.tiles[world.Wall]
	a.tiles[world.TimedSpike] = timedSpike(size)
	a.tiles[world.FireTrap] = fireTrap(size)
	a.Player = player(size)
	a.Enemy = enemy(size)
	a.EnemyVertical = tintImage(a.Enemy, color.RGBA{55, 176, 220, 255})
	a.EnemyChaser = tintImage(a.Enemy, color.RGBA{226, 69, 54, 255})
	a.buildAnimationFrames()
}

func (a *Assets) buildAnimationFrames() {
	a.PlayerFrames = []*ebiten.Image{
		a.Player,
		shiftImage(a.Player, -1, 0),
		a.Player,
		shiftImage(a.Player, 1, 0),
	}
	a.DiamondFrames = []*ebiten.Image{
		a.Tile(world.Diamond),
		brightnessImage(a.Tile(world.Diamond), 1.18),
		brightnessImage(a.Tile(world.Diamond), 0.92),
	}
	a.EnemyFrames = map[world.EnemyType][]*ebiten.Image{
		world.EnemyHorizontal: {
			a.Enemy,
			shiftImage(a.Enemy, -1, 0),
			a.Enemy,
			shiftImage(a.Enemy, 1, 0),
		},
		world.EnemyVertical: {
			a.EnemyVertical,
			shiftImage(a.EnemyVertical, 0, -1),
			a.EnemyVertical,
			shiftImage(a.EnemyVertical, 0, 1),
		},
		world.EnemyChaser: {
			a.EnemyChaser,
			brightnessImage(a.EnemyChaser, 1.2),
			a.EnemyChaser,
			shiftImage(a.EnemyChaser, 1, 0),
		},
	}
}

func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		f, err = os.Open(filepath.Join("..", "..", path))
		if err != nil {
			return nil, err
		}
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func drawImage(dst *ebiten.Image, src *ebiten.Image, x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(src, op)
}

func drawCoverImage(dst *ebiten.Image, src *ebiten.Image, width, height int) {
	bounds := src.Bounds()
	scaleX := float64(width) / float64(bounds.Dx())
	scaleY := float64(height) / float64(bounds.Dy())
	scale := scaleX
	if scaleY > scale {
		scale = scaleY
	}
	drawW := float64(bounds.Dx()) * scale
	drawH := float64(bounds.Dy()) * scale
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((float64(width)-drawW)/2, (float64(height)-drawH)/2)
	dst.DrawImage(src, op)
}

func tintImage(src *ebiten.Image, tint color.RGBA) *ebiten.Image {
	bounds := src.Bounds()
	img := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.Scale(float32(tint.R)/255, float32(tint.G)/255, float32(tint.B)/255, 1)
	img.DrawImage(src, op)
	return img
}

func shiftImage(src *ebiten.Image, dx, dy int) *ebiten.Image {
	bounds := src.Bounds()
	img := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dx), float64(dy))
	img.DrawImage(src, op)
	return img
}

func brightnessImage(src *ebiten.Image, multiplier float64) *ebiten.Image {
	bounds := src.Bounds()
	img := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	op := &ebiten.DrawImageOptions{}
	scale := float32(multiplier)
	op.ColorScale.Scale(scale, scale, scale, 1)
	img.DrawImage(src, op)
	return img
}

func solid(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	img.Fill(c)
	return img
}

func block(size int, base color.RGBA, edge color.RGBA) *ebiten.Image {
	img := solid(size, base)
	for x := 0; x < size; x++ {
		img.Set(x, 0, edge)
		img.Set(x, size-1, edge)
	}
	for y := 0; y < size; y++ {
		img.Set(0, y, edge)
		img.Set(size-1, y, edge)
	}
	return img
}

func diamond(size int, base color.RGBA, shine color.RGBA) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	mid := size / 2
	for y := 5; y < size-5; y++ {
		half := mid - abs(mid-y)
		for x := mid - half; x <= mid+half; x++ {
			if x >= 0 && x < size {
				img.Set(x, y, base)
			}
		}
	}
	for i := 0; i < 6; i++ {
		img.Set(mid-3+i, 8, shine)
		img.Set(mid-5+i, 9, shine)
	}
	return img
}

func rock(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	center := float64(size) / 2
	radius := float64(size) * 0.36
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx, dy := float64(x)-center, float64(y)-center
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x, y, color.RGBA{124, 128, 138, 255})
			}
		}
	}
	return img
}

func key(size int) *ebiten.Image {
	return keyWithColor(size, color.RGBA{206, 214, 224, 255})
}

func goldKey(size int) *ebiten.Image {
	return keyWithColor(size, color.RGBA{242, 191, 77, 255})
}

func keyWithColor(size int, keyColor color.RGBA) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	for y := 12; y < 18; y++ {
		for x := 9; x < 24; x++ {
			img.Set(x, y, keyColor)
		}
	}
	for y := 9; y < 21; y++ {
		for x := 5; x < 13; x++ {
			if (x-9)*(x-9)+(y-15)*(y-15) < 22 {
				img.Set(x, y, keyColor)
			}
		}
	}
	return img
}

func chest(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	wood := color.RGBA{143, 82, 35, 255}
	dark := color.RGBA{81, 48, 28, 255}
	gold := color.RGBA{238, 185, 66, 255}
	shine := color.RGBA{255, 224, 109, 255}
	for y := 11; y < 27; y++ {
		for x := 5; x < 27; x++ {
			img.Set(x, y, wood)
		}
	}
	for x := 5; x < 27; x++ {
		img.Set(x, 11, dark)
		img.Set(x, 26, dark)
		img.Set(x, 17, dark)
	}
	for y := 11; y < 27; y++ {
		img.Set(5, y, dark)
		img.Set(26, y, dark)
		img.Set(15, y, gold)
		img.Set(16, y, gold)
	}
	for y := 18; y < 23; y++ {
		for x := 13; x < 19; x++ {
			img.Set(x, y, gold)
		}
	}
	img.Set(15, 19, shine)
	img.Set(16, 19, shine)
	return img
}

func checkpoint(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	outer := color.RGBA{98, 218, 170, 255}
	inner := color.RGBA{224, 255, 205, 255}
	stone := color.RGBA{83, 94, 88, 255}
	for y := 22; y < 27; y++ {
		for x := 8; x < 24; x++ {
			img.Set(x, y, stone)
		}
	}
	for y := 6; y < 23; y++ {
		for x := 7; x < 25; x++ {
			dx := float64(x - size/2)
			dy := float64(y - 14)
			d := dx*dx + dy*dy
			if d > 70 && d < 105 {
				img.Set(x, y, outer)
			}
			if d > 22 && d < 38 {
				img.Set(x, y, inner)
			}
		}
	}
	return img
}

func potion(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	glass := color.RGBA{170, 235, 240, 255}
	liquid := color.RGBA{224, 57, 91, 255}
	shine := color.RGBA{248, 250, 255, 255}
	for y := 8; y < 13; y++ {
		for x := 13; x < 19; x++ {
			img.Set(x, y, glass)
		}
	}
	for y := 12; y < 27; y++ {
		for x := 9; x < 23; x++ {
			dx := float64(x - size/2)
			dy := float64(y - 20)
			if dx*dx/36+dy*dy/49 <= 1 {
				img.Set(x, y, glass)
			}
		}
	}
	for y := 18; y < 26; y++ {
		for x := 10; x < 22; x++ {
			dx := float64(x - size/2)
			dy := float64(y - 20)
			if dx*dx/36+dy*dy/49 <= 1 {
				img.Set(x, y, liquid)
			}
		}
	}
	img.Set(13, 14, shine)
	img.Set(12, 15, shine)
	return img
}

func hammerPickup(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	handle := color.RGBA{132, 82, 42, 255}
	metal := color.RGBA{182, 191, 204, 255}
	shine := color.RGBA{245, 248, 252, 255}
	drawLine(img, 9, 24, 22, 11, handle)
	drawLine(img, 10, 25, 23, 12, handle)
	for y := 7; y < 14; y++ {
		for x := 16; x < 28; x++ {
			img.Set(x, y, metal)
		}
	}
	for y := 9; y < 17; y++ {
		for x := 13; x < 19; x++ {
			img.Set(x, y, metal)
		}
	}
	for x := 18; x < 26; x++ {
		img.Set(x, 8, shine)
	}
	return img
}

func hookPickup(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	metal := color.RGBA{196, 205, 214, 255}
	rope := color.RGBA{190, 142, 72, 255}
	drawLine(img, 7, 8, 19, 20, rope)
	drawLine(img, 8, 7, 20, 19, rope)
	for y := 9; y < 25; y++ {
		for x := 15; x < 26; x++ {
			dx := x - 20
			dy := y - 16
			d := dx*dx + dy*dy
			if d > 32 && d < 62 && !(x < 20 && y < 15) {
				img.Set(x, y, metal)
			}
		}
	}
	drawLine(img, 21, 22, 27, 19, metal)
	drawLine(img, 21, 23, 27, 20, metal)
	return img
}

func compassPickup(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	rim := color.RGBA{229, 181, 78, 255}
	face := color.RGBA{232, 237, 214, 255}
	needle := color.RGBA{213, 58, 58, 255}
	center := size / 2
	for y := 6; y < size-6; y++ {
		for x := 6; x < size-6; x++ {
			dx := x - center
			dy := y - center
			d := dx*dx + dy*dy
			if d <= 100 {
				img.Set(x, y, rim)
			}
			if d <= 64 {
				img.Set(x, y, face)
			}
		}
	}
	drawLine(img, center, center, center+5, center-8, needle)
	drawLine(img, center, center, center-5, center+8, color.RGBA{62, 87, 158, 255})
	img.Set(center, center, color.RGBA{40, 42, 48, 255})
	return img
}

func secretExit(size int) *ebiten.Image {
	img := block(size, color.RGBA{48, 35, 78, 255}, color.RGBA{116, 82, 160, 255})
	glow := color.RGBA{236, 69, 118, 255}
	shine := color.RGBA{255, 208, 224, 255}
	for y := 8; y < size-5; y++ {
		for x := 8; x < size-8; x++ {
			if x == 8 || x == size-9 || y == 8 || y == size-6 {
				img.Set(x, y, glow)
			}
		}
	}
	mid := size / 2
	for y := 11; y < 22; y++ {
		half := 5 - abs(mid-y)
		if half < 0 {
			continue
		}
		for x := mid - half; x <= mid+half; x++ {
			img.Set(x, y, glow)
		}
	}
	for x := mid - 2; x <= mid+2; x++ {
		img.Set(x, 12, shine)
	}
	return img
}

func player(size int) *ebiten.Image {
	img := solid(size, color.RGBA{0, 0, 0, 0})
	body := color.RGBA{212, 58, 54, 255}
	face := color.RGBA{232, 188, 130, 255}
	for y := 10; y < 28; y++ {
		for x := 9; x < 23; x++ {
			img.Set(x, y, body)
		}
	}
	for y := 4; y < 13; y++ {
		for x := 10; x < 22; x++ {
			img.Set(x, y, face)
		}
	}
	return img
}

func enemy(size int) *ebiten.Image {
	img := solid(size, color.RGBA{0, 0, 0, 0})
	body := color.RGBA{116, 86, 184, 255}
	for y := 7; y < 27; y++ {
		for x := 7; x < 25; x++ {
			if (x-16)*(x-16)+(y-17)*(y-17) < 95 {
				img.Set(x, y, body)
			}
		}
	}
	img.Set(12, 14, color.White)
	img.Set(20, 14, color.White)
	return img
}

func spike(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	steel := color.RGBA{186, 193, 204, 255}
	shadow := color.RGBA{88, 94, 105, 255}
	for i := 0; i < 4; i++ {
		baseX := 3 + i*7
		for y := 10; y < 26; y++ {
			half := (y - 10) / 2
			for x := baseX - half; x <= baseX+half; x++ {
				if x >= 0 && x < size {
					img.Set(x, y, steel)
				}
			}
		}
	}
	for x := 1; x < size-1; x++ {
		img.Set(x, 26, shadow)
		img.Set(x, 27, shadow)
	}
	return img
}

func timedSpike(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	base := color.RGBA{86, 91, 103, 255}
	edge := color.RGBA{144, 150, 162, 255}
	for x := 2; x < size-2; x++ {
		img.Set(x, size-7, edge)
		img.Set(x, size-6, base)
	}
	for x := 5; x < size-5; x += 7 {
		for y := size - 10; y < size-7; y++ {
			img.Set(x, y, edge)
			img.Set(x+1, y, edge)
		}
	}
	return img
}

func fireTrap(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	stone := color.RGBA{74, 66, 56, 255}
	metal := color.RGBA{142, 138, 130, 255}
	ember := color.RGBA{220, 58, 28, 255}
	for y := 20; y < 27; y++ {
		for x := 4; x < size-4; x++ {
			img.Set(x, y, stone)
		}
	}
	for x := 5; x < size-5; x += 5 {
		for y := 10; y < 27; y++ {
			img.Set(x, y, metal)
		}
	}
	for x := 4; x < size-4; x++ {
		img.Set(x, 10, metal)
		img.Set(x, 26, metal)
	}
	for x := 9; x < size-9; x += 5 {
		img.Set(x, 22, ember)
		img.Set(x+1, 21, ember)
	}
	return img
}

func lever(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	base := color.RGBA{100, 88, 72, 255}
	metal := color.RGBA{190, 196, 206, 255}
	red := color.RGBA{220, 58, 50, 255}
	for y := 21; y < 27; y++ {
		for x := 8; x < 24; x++ {
			img.Set(x, y, base)
		}
	}
	for y := 8; y < 22; y++ {
		x := 17 - (y-8)/3
		img.Set(x, y, metal)
		img.Set(x+1, y, metal)
	}
	for y := 5; y < 10; y++ {
		for x := 14; x < 19; x++ {
			img.Set(x, y, red)
		}
	}
	return img
}

func crackedWall(size int) *ebiten.Image {
	img := block(size, color.RGBA{101, 86, 74, 255}, color.RGBA{54, 48, 44, 255})
	crack := color.RGBA{25, 24, 25, 255}
	points := [][2]int{{15, 3}, {13, 8}, {17, 12}, {12, 16}, {16, 20}, {14, 27}}
	for i := 0; i < len(points)-1; i++ {
		drawLine(img, points[i][0], points[i][1], points[i+1][0], points[i+1][1], crack)
	}
	drawLine(img, 17, 12, 24, 10, crack)
	drawLine(img, 12, 16, 6, 20, crack)
	return img
}

func teleporter(size int) *ebiten.Image {
	img := solid(size, color.RGBA{20, 24, 31, 255})
	outer := color.RGBA{78, 213, 240, 255}
	inner := color.RGBA{189, 248, 255, 255}
	for y := 5; y < size-5; y++ {
		for x := 5; x < size-5; x++ {
			dx := float64(x - size/2)
			dy := float64(y - size/2)
			d := dx*dx + dy*dy
			if d > 85 && d < 124 {
				img.Set(x, y, outer)
			}
			if d > 24 && d < 42 {
				img.Set(x, y, inner)
			}
		}
	}
	return img
}

func lava(size int) *ebiten.Image {
	img := solid(size, color.RGBA{102, 26, 18, 255})
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if (x+y)%9 < 4 {
				img.Set(x, y, color.RGBA{230, 74, 28, 255})
			}
			if (x*3+y*2)%17 < 3 {
				img.Set(x, y, color.RGBA{255, 190, 68, 255})
			}
		}
	}
	return img
}

func drawLine(img *ebiten.Image, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func clamp(v, low, high int) int {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}
