package originalgame

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wangle201210/zskc/internal/original"
)

const (
	angkorWorldFrameSheet    = "decoded/sprites/0/chunk02-frames.png"
	angkorBoulderFrameSheet  = "decoded/sprites/0/chunk00-frames.png"
	angkorDiggableFrameSheet = "decoded/sprites/0/chunk01-frames.png"
	angkorFloorSheet         = "decoded/sprites/0/chunk03-modules.png"
	violetGemSheet           = "decoded/sprites/cm/chunk02-frames.png"
	redDiamondSheet          = "decoded/sprites/cm/chunk02-palette01-frames.png"
	checkpointSheet          = "decoded/sprites/cm/chunk06-frames.png"
	quotaModules             = "decoded/sprites/cm/chunk05-modules.png"
	quotaMetadata            = "decoded/sprites/cm/chunk05-animations.json"
	goalSheet                = "decoded/sprites/cm/chunk00-modules.png"
	doorModuleSheet          = "decoded/sprites/cm/chunk01-modules.png"
	doorMetadata             = "decoded/sprites/cm/chunk01-animations.json"
	snakeSheet               = "decoded/sprites/gen1/chunk05-frames.png"
	snakeModuleSheet         = "decoded/sprites/gen1/chunk05-modules.png"
	snakeMetadata            = "decoded/sprites/gen1/chunk05-animations.json"
	redSnakeSheet            = "decoded/sprites/gen1/chunk05-palette01-frames.png"
	redSnakeModuleSheet      = "decoded/sprites/gen1/chunk05-palette01-modules.png"
	crawlerModules           = "decoded/sprites/gen1/chunk04-modules.png"
	crawlerMetadata          = "decoded/sprites/gen1/chunk04-animations.json"
	commonPickupSheet        = "decoded/sprites/gen0/chunk08-frames.png"
	commonPickupModules      = "decoded/sprites/gen0/chunk08-modules.png"
	commonPickupMetadata     = "decoded/sprites/gen0/chunk08-animations.json"
	breakableSheet           = "decoded/sprites/gen0/chunk07-frames.png"
	breakableModules         = "decoded/sprites/gen0/chunk07-modules.png"
	breakableMetadata        = "decoded/sprites/gen0/chunk07-animations.json"
	goldLockSheet            = "decoded/sprites/gen2/chunk08-frames.png"
	goldLockModules          = "decoded/sprites/gen2/chunk08-modules.png"
	goldLockMetadata         = "decoded/sprites/gen2/chunk08-animations.json"
	silverLockSheet          = "decoded/sprites/gen2/chunk08-palette01-frames.png"
	silverLockModules        = "decoded/sprites/gen2/chunk08-palette01-modules.png"
	foregroundEffectSheet    = "decoded/sprites/gen0/chunk04-frames.png"
	foregroundEffectModules  = "decoded/sprites/gen0/chunk04-modules.png"
	foregroundEffectMetadata = "decoded/sprites/gen0/chunk04-animations.json"
	hiddenOverlaySheet       = "decoded/sprites/gen3/chunk03-frames.png"
	hiddenOverlayModules     = "decoded/sprites/gen3/chunk03-modules.png"
	hiddenOverlayMetadata    = "decoded/sprites/gen3/chunk03-animations.json"
	specialContainerSheet    = "decoded/sprites/gen2/chunk02-frames.png"
	specialContainerModules  = "decoded/sprites/gen2/chunk02-modules.png"
	specialContainerMetadata = "decoded/sprites/gen2/chunk02-animations.json"
	goldKeyModules           = "decoded/sprites/gen0/chunk02-modules.png"
	silverKeyModules         = "decoded/sprites/gen0/chunk02-palette01-modules.png"
	keyMetadata              = "decoded/sprites/gen0/chunk02-animations.json"
	toolModules              = "decoded/sprites/gen1/chunk09-modules.png"
	toolMetadata             = "decoded/sprites/gen1/chunk09-animations.json"
	worldMapIconSheet        = "decoded/sprites/ms/chunk00-frames.png"
	worldMapIconModules      = "decoded/sprites/ms/chunk00-modules.png"
	worldMapIconMetadata     = "decoded/sprites/ms/chunk00-animations.json"
	worldMapGroundSheet      = "decoded/sprites/ms/chunk01-frames.png"
	worldMapGroundModules    = "decoded/sprites/ms/chunk01-modules.png"
	worldMapGroundMetadata   = "decoded/sprites/ms/chunk01-animations.json"
	worldMapHeaderSheet      = "decoded/sprites/ms/chunk02-frames.png"
	worldMapHeaderModules    = "decoded/sprites/ms/chunk02-modules.png"
	worldMapHeaderMetadata   = "decoded/sprites/ms/chunk02-animations.json"
	pickupEffectSheet        = "decoded/sprites/cm/chunk07-frames.png"
	pickupEffectModules      = "decoded/sprites/cm/chunk07-modules.png"
	pickupEffectMetadata     = "decoded/sprites/cm/chunk07-animations.json"
	resultSparkModules       = "decoded/sprites/cm/chunk04-modules.png"
	resultSparkMetadata      = "decoded/sprites/cm/chunk04-animations.json"
	resultMedalModules       = "decoded/sprites/ui/chunk04-modules.png"
	resultMedalMetadata      = "decoded/sprites/ui/chunk04-animations.json"
	hazardEmitterSheet       = "decoded/sprites/gen0/chunk09-frames.png"
	hazardFlameSheet         = "decoded/sprites/gen1/chunk00-frames.png"
	hazardFlameModuleSheet   = "decoded/sprites/gen1/chunk00-modules.png"
	hazardFlameMetadata      = "decoded/sprites/gen1/chunk00-animations.json"
	pressureSwitchModules    = "decoded/sprites/gen2/chunk09-modules.png"
	pressureSwitchMetadata   = "decoded/sprites/gen2/chunk09-animations.json"
	hudSheet                 = "decoded/sprites/ui/chunk02-frames.png"
	hudModuleSheet           = "decoded/sprites/ui/chunk02-modules.png"
	hudMetadata              = "decoded/sprites/ui/chunk02-animations.json"
	heroFrameSheet           = "decoded/sprites/o/chunk00-frames.png"
	heroModuleSheet          = "decoded/sprites/o/chunk00-modules.png"
	heroMetadata             = "decoded/sprites/o/chunk00-animations.json"
	fontSmallSheet           = "decoded/fonts/freej2me-small.png"
	fontSmallMetadata        = "decoded/fonts/freej2me-small.json"
	fontMediumSheet          = "decoded/fonts/freej2me-medium.png"
	fontMediumMetadata       = "decoded/fonts/freej2me-medium.json"
	originalAudioDir         = "decoded/audio"
	defaultWorldDir          = "decoded/world0"
	resultLoadingSteps       = 12
	resultTitleTicks         = 40
	resultGemMinimumTicks    = 40
	resultRedDiamondTicks    = 40
	resultHitTicks           = 10
	resultRetryTicks         = 10
	stageIntroDuration       = 60
	deathTransitionTicks     = 80
	chestRewardTick          = 39
	chestRewardSequence      = 13
	chestShortRewardTick     = 23
	chestShortRewardSequence = 6
	sourceTPS                = 20
	framePadding             = 2
	frameCols                = 16
	diggableFrameCellW       = 35
	diggableFrameCellH       = 27
	playfieldHeight          = 240
	playfieldTop             = 40
)

var hookRopeColor = color.RGBA{211, 215, 231, 255}

const resultPhaseLoading = -1

const (
	resultPhaseTitle = iota
	resultPhaseVioletGems
	resultPhaseRedDiamonds
	resultPhaseHits
	resultPhaseRetries
	resultPhaseComplete
)

const (
	resultAwardVioletGems  = 0x04
	resultAwardRedDiamonds = 0x08
	resultAwardNoHits      = 0x10
	resultAwardNoRetries   = 0x20
)

const (
	resultEffectDoubleShort = iota
	resultEffectHalf
	resultEffectDoubleLong
)

var diggableFrameBounds = [...]image.Rectangle{
	image.Rect(0, 0, 24, 24),
	image.Rect(-1, -1, 24, 24),
	image.Rect(-3, -3, 25, 24),
	image.Rect(-4, -3, 27, 24),
	image.Rect(-4, -2, 27, 24),
	image.Rect(-5, 0, 30, 24),
	image.Rect(-6, 2, 28, 25),
	image.Rect(-3, 8, 27, 24),
}

type worldEffect struct {
	Point     original.Point
	Animation int
	Sequence  int
}

type gameMode int

const (
	gameModeStage gameMode = iota
	gameModeWorldMap
)

type Game struct {
	pack                  *original.WorldPack
	rt                    *original.Runtime
	worldFrames           *ebiten.Image
	boulder               *ebiten.Image
	diggable              *ebiten.Image
	floor                 *ebiten.Image
	violetGem             *ebiten.Image
	redDiamond            *ebiten.Image
	checkpoint            *ebiten.Image
	quota                 *spriteSheet
	goal                  *ebiten.Image
	door                  *spriteSheet
	hazard                *ebiten.Image
	snakes                *spriteSheet
	redSnakes             *spriteSheet
	crawler               *spriteSheet
	commonPickups         *spriteSheet
	breakables            *spriteSheet
	goldLock              *spriteSheet
	silverLock            *spriteSheet
	foregroundEffects     *spriteSheet
	hiddenOverlay         *spriteSheet
	specialContainer      *spriteSheet
	goldKey               *spriteSheet
	silverKey             *spriteSheet
	tools                 *spriteSheet
	worldMapIcons         *spriteSheet
	worldMapGround        *spriteSheet
	worldMapHeader        *spriteSheet
	pickupEffects         *spriteSheet
	resultSpark           *spriteSheet
	resultMedal           *spriteSheet
	flames                *spriteSheet
	pressureSwitch        *spriteSheet
	hud                   *spriteSheet
	hero                  *spriteSheet
	fontSmall             *bitmapFont
	fontMedium            *bitmapFont
	worldCanvas           *ebiten.Image
	worldDir              string
	worldMap              *worldMapData
	mode                  gameMode
	stageIndex            int
	worldMapLoadingStep   int
	worldMapSelectedStage int
	worldMapTravelFrom    int
	worldMapTravelTo      int
	worldMapTravelTick    int
	message               string
	tick                  int
	worldDone             bool
	resultPhase           int
	resultPhaseTicks      int
	resultLoadingStep     int
	resultAwards          byte
	resultNewAwards       byte
	introTicks            int
	lastDX                int
	lastDY                int
	heroMoveStart         int
	heroMoveOffset        int
	entranceSteps         int
	checkpointBannerUntil int
	compassDirection      int
	cameraX               int
	cameraY               int
	worldEffects          []worldEffect
	sounds                *originalSounds
	progressPath          string
	progress              originalProgress
}

func Run() error {
	g, err := New(defaultWorldDir)
	if err != nil {
		return err
	}
	if err := g.enableProgress(originalProgressPath()); err != nil {
		return err
	}
	g.sounds.Enable()
	g.sounds.Play(original.SoundAngkorMusic)
	defer g.sounds.Stop()
	ebiten.SetTPS(sourceTPS)
	ebiten.SetWindowTitle("Diamond Rush Original Runtime - Angkor World 0")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	return ebiten.RunGame(g)
}

func New(worldDir string) (*Game, error) {
	resolvedWorldDir := resolvePath(worldDir)
	pack, err := original.LoadWorldDir(resolvedWorldDir)
	if err != nil {
		return nil, err
	}
	frames, err := loadTransparentSheet(angkorWorldFrameSheet)
	if err != nil {
		return nil, fmt.Errorf("load world frame sheet: %w", err)
	}
	boulder, err := loadTransparentSheet(angkorBoulderFrameSheet)
	if err != nil {
		return nil, fmt.Errorf("load boulder frame sheet: %w", err)
	}
	diggable, err := loadTransparentSheet(angkorDiggableFrameSheet)
	if err != nil {
		return nil, fmt.Errorf("load diggable frame sheet: %w", err)
	}
	floor, err := loadTransparentSheet(angkorFloorSheet)
	if err != nil {
		return nil, fmt.Errorf("load floor sheet: %w", err)
	}
	violetGem, err := loadTransparentSheet(violetGemSheet)
	if err != nil {
		return nil, fmt.Errorf("load violet gem sheet: %w", err)
	}
	redDiamond, err := loadTransparentSheet(redDiamondSheet)
	if err != nil {
		return nil, fmt.Errorf("load red diamond sheet: %w", err)
	}
	checkpoint, err := loadTransparentSheet(checkpointSheet)
	if err != nil {
		return nil, fmt.Errorf("load checkpoint sheet: %w", err)
	}
	quota, err := loadModuleSpriteSheet(quotaModules, quotaMetadata)
	if err != nil {
		return nil, fmt.Errorf("load quota marker: %w", err)
	}
	goal, err := loadTransparentSheet(goalSheet)
	if err != nil {
		return nil, fmt.Errorf("load goal sheet: %w", err)
	}
	door, err := loadModuleSpriteSheet(doorModuleSheet, doorMetadata)
	if err != nil {
		return nil, fmt.Errorf("load door sprite: %w", err)
	}
	hazard, err := loadTransparentSheet(hazardEmitterSheet)
	if err != nil {
		return nil, fmt.Errorf("load hazard emitter sheet: %w", err)
	}
	snakes, err := loadSpriteSheetWithModules(snakeSheet, snakeModuleSheet, snakeMetadata)
	if err != nil {
		return nil, fmt.Errorf("load snake sprite: %w", err)
	}
	redSnakes, err := loadSpriteSheetWithModules(redSnakeSheet, redSnakeModuleSheet, snakeMetadata)
	if err != nil {
		return nil, fmt.Errorf("load red snake sprite: %w", err)
	}
	crawler, err := loadModuleSpriteSheet(crawlerModules, crawlerMetadata)
	if err != nil {
		return nil, fmt.Errorf("load crawler modules: %w", err)
	}
	commonPickups, err := loadSpriteSheetWithModules(commonPickupSheet, commonPickupModules, commonPickupMetadata)
	if err != nil {
		return nil, fmt.Errorf("load common pickups: %w", err)
	}
	breakables, err := loadSpriteSheetWithModules(breakableSheet, breakableModules, breakableMetadata)
	if err != nil {
		return nil, fmt.Errorf("load breakable wall sprite: %w", err)
	}
	goldLock, err := loadSpriteSheetWithModules(goldLockSheet, goldLockModules, goldLockMetadata)
	if err != nil {
		return nil, fmt.Errorf("load gold lock sprite: %w", err)
	}
	silverLock, err := loadSpriteSheetWithModules(silverLockSheet, silverLockModules, goldLockMetadata)
	if err != nil {
		return nil, fmt.Errorf("load silver lock sprite: %w", err)
	}
	foregroundEffects, err := loadSpriteSheetWithModules(foregroundEffectSheet, foregroundEffectModules, foregroundEffectMetadata)
	if err != nil {
		return nil, fmt.Errorf("load foreground effects: %w", err)
	}
	hiddenOverlay, err := loadSpriteSheetWithModules(hiddenOverlaySheet, hiddenOverlayModules, hiddenOverlayMetadata)
	if err != nil {
		return nil, fmt.Errorf("load hidden overlay: %w", err)
	}
	specialContainer, err := loadSpriteSheetWithModules(specialContainerSheet, specialContainerModules, specialContainerMetadata)
	if err != nil {
		return nil, fmt.Errorf("load special pickup container: %w", err)
	}
	goldKey, err := loadModuleSpriteSheet(goldKeyModules, keyMetadata)
	if err != nil {
		return nil, fmt.Errorf("load gold key: %w", err)
	}
	silverKey, err := loadModuleSpriteSheet(silverKeyModules, keyMetadata)
	if err != nil {
		return nil, fmt.Errorf("load silver key: %w", err)
	}
	tools, err := loadModuleSpriteSheet(toolModules, toolMetadata)
	if err != nil {
		return nil, fmt.Errorf("load source tool icons: %w", err)
	}
	worldMapIcons, err := loadSpriteSheetWithModules(worldMapIconSheet, worldMapIconModules, worldMapIconMetadata)
	if err != nil {
		return nil, fmt.Errorf("load world-map icons: %w", err)
	}
	worldMapGround, err := loadSpriteSheetWithModules(worldMapGroundSheet, worldMapGroundModules, worldMapGroundMetadata)
	if err != nil {
		return nil, fmt.Errorf("load world-map ground: %w", err)
	}
	worldMapHeader, err := loadSpriteSheetWithModules(worldMapHeaderSheet, worldMapHeaderModules, worldMapHeaderMetadata)
	if err != nil {
		return nil, fmt.Errorf("load world-map header: %w", err)
	}
	pickupEffects, err := loadSpriteSheetWithModules(pickupEffectSheet, pickupEffectModules, pickupEffectMetadata)
	if err != nil {
		return nil, fmt.Errorf("load pickup effects: %w", err)
	}
	resultSpark, err := loadModuleSpriteSheet(resultSparkModules, resultSparkMetadata)
	if err != nil {
		return nil, fmt.Errorf("load result spark: %w", err)
	}
	resultMedal, err := loadModuleSpriteSheet(resultMedalModules, resultMedalMetadata)
	if err != nil {
		return nil, fmt.Errorf("load result medal: %w", err)
	}
	flames, err := loadSpriteSheetWithModules(hazardFlameSheet, hazardFlameModuleSheet, hazardFlameMetadata)
	if err != nil {
		return nil, fmt.Errorf("load hazard flames: %w", err)
	}
	pressureSwitch, err := loadModuleSpriteSheet(pressureSwitchModules, pressureSwitchMetadata)
	if err != nil {
		return nil, fmt.Errorf("load pressure switch: %w", err)
	}
	hud, err := loadSpriteSheetWithModules(hudSheet, hudModuleSheet, hudMetadata)
	if err != nil {
		return nil, fmt.Errorf("load HUD sprite: %w", err)
	}
	hero, err := loadSpriteSheetWithModules(heroFrameSheet, heroModuleSheet, heroMetadata)
	if err != nil {
		return nil, fmt.Errorf("load hero sprite: %w", err)
	}
	fontSmall, err := loadBitmapFont(fontSmallSheet, fontSmallMetadata)
	if err != nil {
		return nil, fmt.Errorf("load FreeJ2ME small font: %w", err)
	}
	fontMedium, err := loadBitmapFont(fontMediumSheet, fontMediumMetadata)
	if err != nil {
		return nil, fmt.Errorf("load FreeJ2ME medium font: %w", err)
	}
	sounds, err := loadOriginalSounds(originalAudioDir)
	if err != nil {
		return nil, fmt.Errorf("load original sound bank: %w", err)
	}
	worldMap, err := loadWorldMap(filepath.Join(resolvedWorldDir, "map.json"))
	if err != nil {
		return nil, fmt.Errorf("load Angkor world map: %w", err)
	}
	rt, err := original.NewRuntime(pack.Stages[0])
	if err != nil {
		return nil, err
	}
	g := &Game{
		pack:              pack,
		rt:                rt,
		worldFrames:       frames,
		boulder:           boulder,
		diggable:          diggable,
		floor:             floor,
		violetGem:         violetGem,
		redDiamond:        redDiamond,
		checkpoint:        checkpoint,
		quota:             quota,
		goal:              goal,
		door:              door,
		hazard:            hazard,
		snakes:            snakes,
		redSnakes:         redSnakes,
		crawler:           crawler,
		commonPickups:     commonPickups,
		breakables:        breakables,
		goldLock:          goldLock,
		silverLock:        silverLock,
		foregroundEffects: foregroundEffects,
		hiddenOverlay:     hiddenOverlay,
		specialContainer:  specialContainer,
		goldKey:           goldKey,
		silverKey:         silverKey,
		tools:             tools,
		worldMapIcons:     worldMapIcons,
		worldMapGround:    worldMapGround,
		worldMapHeader:    worldMapHeader,
		pickupEffects:     pickupEffects,
		resultSpark:       resultSpark,
		resultMedal:       resultMedal,
		flames:            flames,
		pressureSwitch:    pressureSwitch,
		hud:               hud,
		hero:              hero,
		fontSmall:         fontSmall,
		fontMedium:        fontMedium,
		sounds:            sounds,
		worldCanvas:       ebiten.NewImage(original.ScreenWidth, playfieldHeight),
		worldDir:          resolvedWorldDir,
		worldMap:          worldMap,
		message:           "Original Angkor World 0 runtime",
		lastDX:            1,
		entranceSteps:     rt.EntranceScrollX,
		progress:          newOriginalProgress(),
	}
	g.resetCamera()
	g.updateCompass()
	return g, nil
}

func (g *Game) Update() error {
	g.tick++
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if g.mode == gameModeWorldMap {
		g.updateWorldMap(centerActionPressed())
		return nil
	}
	if g.introTicks < stageIntroDuration {
		g.introTicks++
	}
	if g.worldDone {
		g.updateStageResults(centerActionPressed())
		return nil
	}
	// Java updates active foreground/player objects before reducing the
	// hero's jInt movement offset for this source frame.
	recallPendingAtFrameStart := g.rt.RecallPending
	deadAtFrameStart := g.rt.PlayerDead
	chestOpeningAtFrameStart := g.rt.ChestOpening
	lockOpeningAtFrameStart := g.rt.LockOpening
	g.tickWorld()
	g.syncHeroMotion()
	if (recallPendingAtFrameStart && !g.rt.RecallPending) ||
		(deadAtFrameStart && !g.rt.PlayerDead) ||
		(chestOpeningAtFrameStart && !g.rt.ChestOpening) ||
		(lockOpeningAtFrameStart && !g.rt.LockOpening) {
		g.updateCamera()
		return nil
	}
	if g.rt.ReachedGoal {
		if g.rt.PlayerMotion.Remaining > 0 {
			g.advanceHeroMotion()
		} else {
			g.advanceAfterGoal()
		}
		g.updateCamera()
		return nil
	}
	if g.rt.PlayerMotion.Remaining > 0 {
		g.advanceHeroMotion()
		g.updateCamera()
		return nil
	}
	if g.entranceSteps > 0 {
		if g.entranceSteps == 1 {
			g.rt.CloseEntranceDoor()
		}
		if g.startPlayerMove(1, 0) {
			g.entranceSteps--
		}
		g.updateCamera()
		return nil
	}
	if g.rt.CanAcceptInput() && centerActionPressed() {
		if g.rt.IsCheckpoint(g.rt.Player.X, g.rt.Player.Y) {
			if g.rt.ResetCheckpoint() {
				g.message = "checkpoint reset"
			}
		} else if g.rt.UseSpecialBarrier(g.lastDX, g.lastDY) {
			g.message = "special barrier opened"
		} else if g.rt.UseHammer(g.lastDX, g.lastDY) {
			g.message = "hammer hit breakable wall"
		} else if g.rt.UseHook(g.lastDX, g.lastDY) {
			g.message = "hook pulled object"
		}
		g.playPendingSounds()
	}
	if g.rt.CanAcceptInput() && recallPressed() {
		onCheckpoint := g.rt.IsCheckpoint(g.rt.Player.X, g.rt.Player.Y)
		if g.rt.RecallCheckpoint() {
			switch {
			case onCheckpoint:
				g.message = "checkpoint reset"
			case g.rt.RecallPending:
				g.message = "recalling checkpoint"
			case g.rt.PlayerDead:
				g.message = "recall failed: no lives"
			default:
				g.message = "checkpoint recalled"
			}
		}
		g.playPendingSounds()
	}
	dx, dy := heldDirection()
	if dx != 0 || dy != 0 {
		g.startPlayerMove(dx, dy)
	} else {
		g.rt.ResetPushAttempt()
	}
	g.updateCamera()
	return nil
}

func (g *Game) tickWorld() {
	g.advanceWorldEffects()
	flameFrame, _ := g.flames.animationSequenceIndex(0, g.tick)
	flameReach := 0
	switch {
	case flameFrame > 20:
		flameReach = 3
	case flameFrame > 10:
		flameReach = 2
	case flameFrame > 0:
		flameReach = 1
	}
	playerBeforeFrame := g.rt.Player
	result := g.rt.TickSourceFrame(8, g.tick, flameReach)
	for _, point := range result.VioletPickups {
		g.worldEffects = append(g.worldEffects, worldEffect{Point: point, Animation: 3})
	}
	if g.rt.Player != playerBeforeFrame {
		if g.rt.PlayerMotion.Remaining <= 0 {
			g.resetCamera()
		}
	}
	g.syncHeroMotion()
	switch {
	case result.HazardHits > 0:
		g.message = fmt.Sprintf("horizontal hazard hit %d", result.HazardHits)
	case result.GravityMoved > 0:
		g.message = fmt.Sprintf("gravity moved %d object(s)", result.GravityMoved)
	case result.SnakesMoved > 0:
		g.message = fmt.Sprintf("snakes moved %d", result.SnakesMoved)
	case result.CrawlersMoved > 0:
		g.message = fmt.Sprintf("crawlers moved %d", result.CrawlersMoved)
	}
	if broken := g.rt.TickBreakables(); broken > 0 {
		g.message = fmt.Sprintf("hammer broke %d wall(s)", broken)
	}
	if cleared := g.rt.TickForegroundTriggers(); cleared > 0 {
		g.message = fmt.Sprintf("foreground raw2 opened %d cell(s)", cleared)
	}
	if g.tick&0xf == 0 {
		g.updateCompass()
	}
	g.playPendingSounds()
}

func (g *Game) playPendingSounds() {
	if g == nil || g.rt == nil {
		return
	}
	for _, id := range g.rt.DrainSoundEvents() {
		g.sounds.Play(id)
	}
}

func (g *Game) updateCompass() {
	direction, ok := g.rt.CompassDirection()
	if ok {
		g.compassDirection = direction
	}
}

func (g *Game) advanceWorldEffects() {
	if len(g.worldEffects) == 0 || g.pickupEffects == nil {
		return
	}
	kept := g.worldEffects[:0]
	for _, effect := range g.worldEffects {
		effect.Sequence++
		if effect.Animation < 0 || effect.Animation >= len(g.pickupEffects.meta.AnimationCounts) || effect.Sequence >= g.pickupEffects.meta.AnimationCounts[effect.Animation] {
			continue
		}
		kept = append(kept, effect)
	}
	g.worldEffects = kept
}

func (g *Game) startPlayerMove(dx, dy int) bool {
	g.lastDX = dx
	g.lastDY = dy
	checkpointProgress := g.rt.CheckpointProgress
	if !g.rt.TryMove(dx, dy) {
		g.playPendingSounds()
		return false
	}
	g.playPendingSounds()
	g.heroMoveStart = g.tick
	g.syncHeroMotion()
	if g.rt.CheckpointProgress > checkpointProgress {
		g.checkpointBannerUntil = g.tick + 13
	}
	g.message = fmt.Sprintf("pos %d,%d hp %d violet %d red %d bonus %d rem %d key9 %d key8 %d lives %d special %d relic %d locks %d walls %d", g.rt.Player.X, g.rt.Player.Y, g.rt.Health, g.rt.VioletGems, g.rt.RedDiamonds, g.rt.BonusValue, g.rt.BonusRemaining, g.rt.KeyForForeground9, g.rt.KeyForForeground8, g.rt.ExtraLives, g.rt.SpecialPickups, g.rt.RelicMask, g.rt.LocksOpened, g.rt.BreakableWalls)
	if g.rt.ReachedGoal {
		g.message = "goal reached"
	}
	return true
}

func (g *Game) advanceHeroMotion() {
	g.rt.AdvancePlayerMotion()
	g.syncHeroMotion()
}

func (g *Game) advanceAfterGoal() {
	if g.worldDone {
		g.message = fmt.Sprintf("Angkor stage %02d complete", g.stageIndex+1)
		return
	}
	if g.rt == nil || !g.rt.ReachedGoal {
		return
	}
	moved, complete := g.rt.AdvanceGoalExit()
	if moved {
		g.lastDX = g.rt.PlayerMotion.DX
		g.lastDY = g.rt.PlayerMotion.DY
		g.heroMoveStart = g.tick
		// The source auto-exit branch creates jInt=18 and subtracts six in
		// the same update; input-initiated movement keeps the initial 18.
		g.advanceHeroMotion()
	}
	if !complete {
		g.message = fmt.Sprintf("Angkor stage %02d exit", g.stageIndex+1)
		return
	}
	g.worldDone = true
	g.beginStageResults()
	g.message = fmt.Sprintf("Angkor stage %02d complete", g.stageIndex+1)
}

func (g *Game) syncHeroMotion() {
	if g.rt == nil {
		g.heroMoveOffset = 0
		return
	}
	motion := g.rt.PlayerMotion
	g.heroMoveOffset = motion.Remaining
	if motion.Remaining > 0 && (motion.DX != 0 || motion.DY != 0) {
		g.lastDX = motion.DX
		g.lastDY = motion.DY
	}
}

func (g *Game) beginStageResults() {
	g.resultPhase = resultPhaseLoading
	g.resultPhaseTicks = 0
	g.resultLoadingStep = 0
	g.resultAwards = stageResultAwards(g.rt)
	g.resultNewAwards = g.progress.recordStageResult(g.stageIndex, g.rt)
	if g.progressPath != "" {
		if err := saveOriginalProgress(g.progressPath, g.progress); err != nil {
			g.message = err.Error()
		}
	}
}

func (g *Game) enableProgress(path string) error {
	progress, err := loadOriginalProgress(path)
	if err != nil {
		return err
	}
	g.progressPath = path
	g.progress = progress
	g.applyCampaignProgress(g.rt, g.stageIndex)
	return nil
}

func (g *Game) updateStageResults(skip bool) {
	if g.resultPhase == resultPhaseLoading {
		g.resultLoadingStep++
		if g.resultLoadingStep >= resultLoadingSteps {
			g.resultPhase = resultPhaseTitle
			g.resultPhaseTicks = 1
			g.sounds.Play(original.SoundStageClear)
		}
		return
	}
	if skip && g.resultPhase < resultPhaseComplete {
		g.resultPhase++
		g.resultPhaseTicks = 1
		return
	}
	if g.resultPhase == resultPhaseComplete {
		if skip {
			g.enterWorldMap()
			return
		}
		g.resultPhaseTicks++
		return
	}
	if g.resultPhaseTicks > stageResultPhaseDuration(g.resultPhase, g.rt.VioletGems) {
		g.resultPhase++
		g.resultPhaseTicks = 1
		return
	}
	g.resultPhaseTicks++
}

func (g *Game) loadStage(index int) {
	if index < 0 || index >= len(g.pack.Stages) || index >= angkorReplicaStageCount {
		g.message = fmt.Sprintf("invalid Angkor stage %d", index+1)
		return
	}
	rt, err := original.NewRuntime(g.pack.Stages[index])
	if err != nil {
		g.message = err.Error()
		return
	}
	g.applyCampaignProgress(rt, index)
	g.stageIndex = index
	g.rt = rt
	g.worldDone = false
	g.resultPhase = resultPhaseLoading
	g.resultPhaseTicks = 0
	g.resultLoadingStep = 0
	g.resultAwards = 0
	g.resultNewAwards = 0
	g.worldEffects = nil
	g.introTicks = 0
	g.heroMoveStart = 0
	g.heroMoveOffset = 0
	g.entranceSteps = rt.EntranceScrollX
	g.resetCamera()
	g.updateCompass()
	if g.sounds != nil && g.sounds.enabled {
		g.sounds.Play(original.SoundAngkorMusic)
	}
	g.message = fmt.Sprintf("loaded Angkor stage %02d", index+1)
}

func (g *Game) applyCampaignProgress(rt *original.Runtime, stageIndex int) {
	if rt == nil {
		return
	}
	progress := g.progress.normalized()
	rt.ExtraLives = progress.ExtraLives
	rt.MaxHealth = progress.MaxHealth
	rt.Health = rt.MaxHealth
	toolLevel := progress.ToolLevel
	if stageIndex == 4 {
		// Angkor Stage 5 is revisited after the Mystic Hook is obtained in
		// Bavaria. This five-stage slice has no Bavaria node, so load the stage
		// with the same prerequisite state the original route expects.
		toolLevel = maxToolLevel(toolLevel, 2)
	}
	rt.SpecialItemMask = toolLevelSpecialItemMask(toolLevel)
	rt.SaveSnapshot()
}

func toolLevelSpecialItemMask(level int) int {
	switch level {
	case 1:
		return 1
	case 2:
		return 2
	case 8:
		return 8
	default:
		return 0
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{6, 8, 10, 255})
	if g.mode == gameModeWorldMap {
		g.drawWorldMap(screen)
		return
	}
	if g.worldDone {
		g.drawStageResults(screen)
		return
	}
	world := g.worldCanvas
	if world == nil {
		world = ebiten.NewImage(original.ScreenWidth, playfieldHeight)
		g.worldCanvas = world
	}
	world.Fill(color.RGBA{6, 8, 10, 255})
	camX, camY := g.cameraPixels()
	firstX := camX / original.TileSize
	firstY := camY / original.TileSize
	offX := -(camX % original.TileSize)
	offY := -(camY % original.TileSize)
	for y := firstY; y <= firstY+playfieldHeight/original.TileSize+1; y++ {
		for x := firstX; x <= firstX+original.ScreenWidth/original.TileSize+1; x++ {
			if x < 0 || y < 0 || x >= g.rt.Width() || y >= g.rt.Height() {
				continue
			}
			px := offX + (x-firstX)*original.TileSize
			py := offY + (y-firstY)*original.TileSize
			g.drawCell(world, x, y, px, py)
		}
	}
	g.drawWorldEffects(world, camX, camY)
	playerX, playerY := g.renderedPlayerPixels()
	renderedPlayerX := playerX - camX
	renderedPlayerY := playerY - camY
	g.drawPlayer(world, renderedPlayerX, renderedPlayerY)
	g.drawChestRewardEffect(world, renderedPlayerX, renderedPlayerY)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, playfieldTop)
	screen.DrawImage(world, op)
	g.drawHUD(screen)
	if g.introTicks < stageIntroDuration {
		g.drawStageIntro(screen)
	} else if g.rt.ReachedGoal {
		drawSourcePanelLabel(screen, g.fontMedium, "CONGRATULATIONS!", original.ScreenWidth/2, playfieldTop+160)
	} else if g.checkpointBannerUntil > g.tick {
		drawSourcePanelLabel(screen, g.fontMedium, "CHECKPOINT", original.ScreenWidth/2, playfieldTop+playfieldHeight-g.fontMedium.meta.FontHeight-10)
	}
	if curtain := g.deathCurtainHeight(); curtain > 0 {
		drawRect(screen, 0, 0, original.ScreenWidth, curtain, color.Black)
		drawRect(screen, 0, original.ScreenHeight-curtain, original.ScreenWidth, curtain, color.Black)
	}
}

func (g *Game) deathCurtainHeight() int {
	if g.rt == nil || !g.rt.PlayerDead || g.rt.HurtTicks > 0 {
		return 0
	}
	transitionRemaining := min(deathTransitionTicks, g.rt.DeathTicks)
	elapsed := deathTransitionTicks - transitionRemaining
	return min(original.ScreenHeight/2, elapsed*12)
}

func (g *Game) drawStageResults(screen *ebiten.Image) {
	if g.resultPhase == resultPhaseLoading {
		g.drawStageResultLoading(screen)
		return
	}
	screen.Fill(color.RGBA{0x26, 0x17, 0x07, 0xff})
	textColor := color.White
	titleOffset, completeOffset := stageResultTitleOffsets(g.resultPhase, g.resultPhaseTicks)
	g.fontSmall.drawText(screen, fmt.Sprintf("STAGE %d", g.stageIndex+1), original.ScreenWidth/2+titleOffset, 10, true, textColor)
	g.fontSmall.drawText(screen, "COMPLETE!", original.ScreenWidth/2+completeOffset, 25, true, textColor)

	if g.resultPhase >= resultPhaseVioletGems {
		offset := stageResultRowOffset(g.resultPhase, resultPhaseVioletGems, g.resultPhaseTicks)
		drawSpriteFrame(screen, g.violetGem, 0, 7+offset, 69)
		g.fontSmall.drawText(screen, "DIAMONDS", original.ScreenWidth/2, 69, true, textColor)
		count := g.rt.VioletGems
		if g.resultPhase == resultPhaseVioletGems {
			count = min(g.resultPhaseTicks>>1, count)
		}
		g.fontSmall.drawText(screen, fmt.Sprintf("%d/%d", count, g.rt.TotalVioletGems), original.ScreenWidth/2, 81, true, textColor)
	}
	if g.resultPhase >= resultPhaseRedDiamonds {
		offset := stageResultRowOffset(g.resultPhase, resultPhaseRedDiamonds, g.resultPhaseTicks)
		drawSpriteFrame(screen, g.redDiamond, 0, 7+offset, 127)
		g.fontSmall.drawText(screen, "RED DIAMONDS", original.ScreenWidth/2, 127, true, textColor)
		g.fontSmall.drawText(screen, fmt.Sprintf("%d/%d", g.rt.RedDiamonds, g.rt.TotalRedDiamonds), original.ScreenWidth/2, 139, true, textColor)
		g.drawStageResultAward(screen, resultAwardVioletGems, resultPhaseRedDiamonds, 63, 80, resultEffectDoubleShort)
	}
	if g.resultPhase >= resultPhaseHits {
		offset := stageResultRowOffset(g.resultPhase, resultPhaseHits, g.resultPhaseTicks)
		g.drawHeroResultIcon(screen, 10, 7+offset, 189)
		g.fontSmall.drawText(screen, "HITS", original.ScreenWidth/2, 185, true, textColor)
		g.fontSmall.drawText(screen, fmt.Sprintf("%d", g.rt.HitCount), original.ScreenWidth/2, 197, true, textColor)
		g.drawStageResultAward(screen, resultAwardRedDiamonds, resultPhaseHits, 121, 138, resultEffectHalf)
	}
	if g.resultPhase >= resultPhaseRetries {
		offset := stageResultRowOffset(g.resultPhase, resultPhaseRetries, g.resultPhaseTicks)
		g.drawHeroResultIcon(screen, 12, 7+offset, 243)
		g.fontSmall.drawText(screen, "RETRIES", original.ScreenWidth/2, 243, true, textColor)
		g.fontSmall.drawText(screen, fmt.Sprintf("%d", g.rt.Retries), original.ScreenWidth/2, 255, true, textColor)
		g.drawStageResultAward(screen, resultAwardNoHits, resultPhaseRetries, 179, 196, resultEffectHalf)
	}
	if g.resultPhase >= resultPhaseComplete {
		g.drawStageResultAward(screen, resultAwardNoRetries, resultPhaseComplete, 237, 254, resultEffectDoubleLong)
	}
	prompt := "SKIP"
	if g.resultPhase == resultPhaseComplete {
		prompt = "CONTINUE"
	}
	g.fontSmall.drawText(screen, prompt, 5, 318, false, textColor)
}

func (g *Game) drawStageResultLoading(screen *ebiten.Image) {
	screen.Fill(color.Black)
	progress := min(230, (g.resultLoadingStep+1)*230/resultLoadingSteps)
	drawRect(screen, 5, 310, progress, 6, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	drawRect(screen, 4, 309, 231, 1, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	drawRect(screen, 4, 316, 231, 1, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	drawRect(screen, 4, 310, 1, 6, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	drawRect(screen, 234, 310, 1, 6, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	g.fontMedium.drawText(screen, "LOADING", original.ScreenWidth/2, 304, true, color.White)
}

func (g *Game) drawHeroResultIcon(screen *ebiten.Image, animation, x, y int) {
	frame, ok := g.hero.animationFrame(animation, 0)
	if !ok {
		return
	}
	g.hero.drawFrame(screen, frame.Frame, x, y, frame.Flags)
}

func (g *Game) drawStageResultAward(screen *ebiten.Image, award byte, revealPhase, sparkY, medalY, effectMode int) {
	if g.resultAwards&award == 0 {
		return
	}
	if g.resultNewAwards&award != 0 {
		g.resultSpark.drawModule(screen, 0, 200, sparkY)
		if sequence, ok := stageResultEffectSequence(effectMode, g.resultPhaseTicks, g.pickupEffects.meta.AnimationCounts[0]); g.resultPhase == revealPhase && ok {
			g.pickupEffects.drawAnimationRawSequenceFrame(screen, 0, sequence, 200, sparkY, 0)
		}
	}
	g.resultMedal.drawModule(screen, 0, 180, medalY)
}

func stageResultEffectSequence(mode, tick, frameCount int) (int, bool) {
	if tick < 0 || frameCount <= 0 {
		return 0, false
	}
	switch mode {
	case resultEffectDoubleShort:
		return tick << 1, tick < frameCount>>1
	case resultEffectHalf:
		return tick >> 1, tick < frameCount<<1
	case resultEffectDoubleLong:
		return tick << 1, tick < frameCount<<1
	default:
		return 0, false
	}
}

func stageResultPhaseDuration(phase, violetGems int) int {
	switch phase {
	case resultPhaseTitle:
		return resultTitleTicks
	case resultPhaseVioletGems:
		return max(resultGemMinimumTicks, violetGems<<1)
	case resultPhaseRedDiamonds:
		return resultRedDiamondTicks
	case resultPhaseHits:
		return resultHitTicks
	case resultPhaseRetries:
		return resultRetryTicks
	default:
		return 0
	}
}

func stageResultTitleOffsets(phase, tick int) (int, int) {
	if phase != resultPhaseTitle {
		return 0, 0
	}
	raw := -100 + tick*10
	return min(raw, 0), min(raw-240, 0)
}

func stageResultRowOffset(phase, rowPhase, tick int) int {
	if phase != rowPhase {
		return 0
	}
	return min(-100+tick*10, 0)
}

func stageResultAwards(rt *original.Runtime) byte {
	if rt == nil {
		return 0
	}
	var awards byte
	if rt.VioletGems == rt.TotalVioletGems {
		awards |= resultAwardVioletGems
	}
	if rt.RedDiamonds == rt.TotalRedDiamonds {
		awards |= resultAwardRedDiamonds
	}
	if rt.HitCount == 0 {
		awards |= resultAwardNoHits
	}
	if rt.Retries == 0 {
		awards |= resultAwardNoRetries
	}
	return awards
}

func (g *Game) Layout(_, _ int) (int, int) {
	return original.ScreenWidth, original.ScreenHeight
}

func (g *Game) cameraPixels() (int, int) {
	return g.cameraX, g.cameraY
}

func (g *Game) resetCamera() {
	if g.rt == nil {
		return
	}
	g.cameraX = 0
	maxY := max(0, g.rt.Height()*original.TileSize-playfieldHeight)
	g.cameraY = clamp(g.rt.Player.Y*original.TileSize-160, 0, maxY)
}

func (g *Game) updateCamera() {
	if g.rt == nil {
		return
	}
	playerX, playerY := g.renderedPlayerPixels()
	maxX := max(0, g.rt.Width()*original.TileSize-original.ScreenWidth)
	maxY := max(0, g.rt.Height()*original.TileSize-playfieldHeight)
	if playerX < g.cameraX+96 {
		g.cameraX = (g.cameraX - 96 + playerX) >> 1
	} else if playerX > g.cameraX+120 {
		g.cameraX = (g.cameraX - 120 + playerX) >> 1
	}
	screenPlayerY := playerY + playfieldTop
	if screenPlayerY < g.cameraY+96 {
		g.cameraY = (g.cameraY - 96 + screenPlayerY) >> 1
	} else if screenPlayerY > g.cameraY+160 {
		g.cameraY = (g.cameraY - 160 + screenPlayerY) >> 1
	}
	g.cameraX = clamp(g.cameraX, 0, maxX)
	g.cameraY = clamp(g.cameraY, 0, maxY)
}

func (g *Game) renderedPlayerPixels() (int, int) {
	return g.rt.Player.X*original.TileSize - g.lastDX*g.heroMoveOffset,
		g.rt.Player.Y*original.TileSize - g.lastDY*g.heroMoveOffset
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	g.hud.drawFrame(screen, 20, original.ScreenWidth/2, 0, 0)
	g.hud.drawFrame(screen, 0, original.ScreenWidth/2, original.ScreenHeight, 0)
	g.hud.drawFrame(screen, 1, original.ScreenWidth/2, original.ScreenHeight, 0)
	if g.rt.CompassEnabled {
		g.hud.drawFrame(screen, 2, original.ScreenWidth/2+2, original.ScreenHeight, 0)
		g.hud.drawFrame(screen, 3+g.compassDirection, original.ScreenWidth/2+2, original.ScreenHeight, 0)
	}
	g.drawHealth(screen)
	g.hud.drawNumber(screen, g.rt.VioletGems, 190, 308)
	g.hud.drawNumber(screen, g.rt.RedDiamonds, 227, 308)
	g.hud.drawNumber(screen, g.rt.KeyForForeground9, 167, 18)
	g.hud.drawNumber(screen, g.rt.KeyForForeground8, 207, 18)
	g.hud.drawNumber(screen, g.rt.ExtraLives, 91, 18)
}

func (g *Game) drawStageIntro(screen *ebiten.Image) {
	worldX, stageX := stageTitlePositions(g.introTicks)
	drawSourcePanelLabel(screen, g.fontMedium, "ANGKOR WAT", worldX, playfieldTop+15)
	drawSourcePanelLabel(screen, g.fontMedium, fmt.Sprintf("STAGE %d", g.stageIndex+1), stageX, playfieldTop+50)
}

func stageTitlePositions(tick int) (int, int) {
	remaining := stageIntroDuration - clamp(tick, 0, stageIntroDuration)
	worldX := original.ScreenWidth / 2
	switch {
	case remaining < 20:
		worldX = (remaining - 10) * original.ScreenWidth / 20
	case remaining >= 50:
		worldX = (stageIntroDuration - remaining) * original.ScreenWidth / 15
	}
	return worldX, original.ScreenWidth - worldX
}

func (g *Game) drawHealth(screen *ebiten.Image) {
	lowHealthOffset := 0
	if g.rt.Health <= 1 {
		lowHealthOffset = 1
	}
	leftModule := 11 + lowHealthOffset
	emptyModule := 13 + lowHealthOffset
	filledModule := 15 + lowHealthOffset
	rightModule := 17 + lowHealthOffset
	x := original.ScreenWidth/2 - 33 - (g.rt.MaxHealth-4)*g.hud.moduleWidth(filledModule)/2
	y := original.ScreenHeight - 29
	g.hud.drawModule(screen, leftModule, x, y)
	x += g.hud.moduleWidth(leftModule)
	for cell := 0; cell < g.rt.MaxHealth; cell++ {
		module := emptyModule
		if (g.rt.Health > 1 && cell < g.rt.Health) || (g.rt.Health <= 1 && cell == 0 && (g.tick/4)%2 == 0) {
			module = filledModule
		}
		g.hud.drawModule(screen, module, x, y)
		x += g.hud.moduleWidth(module)
	}
	g.hud.drawModule(screen, rightModule, x, y)
}

func (g *Game) drawCell(dst *ebiten.Image, x, y, px, py int) {
	cellPX, cellPY := px, py
	drawRect(dst, px, py, original.TileSize, original.TileSize, color.RGBA{18, 20, 24, 255})
	playerID, _ := g.rt.At(original.PlayerLayer, x, y)
	if playerID == original.EmptyRawID || playerID < 80 {
		drawSpriteFrame(dst, g.floor, 0, px, py)
	} else {
		g.drawWorldFrame(dst, int(playerID-80), px, py)
	}
	foregroundID, _ := g.rt.At(original.ForegroundLayer, x, y)
	if foregroundID != original.EmptyRawID {
		id := foregroundID
		switch {
		case id == 4:
			frame := 7
			orderRaw, _ := g.rt.At(original.BackgroundLayer, x, y)
			order := int(orderRaw)
			if order != int(original.EmptyRawID) && order >= g.rt.CheckpointProgress {
				frame = (g.tick >> 1) % 7
			}
			drawSpriteFrame(dst, g.checkpoint, frame, px, py)
		case id == 5 || id == 28:
			drawSpriteFrame(dst, g.goal, 0, px, py)
		case id == 6:
			g.drawPressureSwitch(dst, x, y, px, py, playerID)
		case id == 7:
			state, _ := g.rt.At(original.BackgroundLayer, x, y)
			frame := max(0, ((int(state)&0xf0)>>4)-1)
			g.door.drawFrame(dst, frame, px, py, 0)
			g.door.drawFrame(dst, frame+3, px, py, 0)
		case id == 8:
			g.silverLock.drawFrame(dst, clamp(g.objectStateAt(x, y), 0, 6), px, py, 0)
		case id == 9:
			g.goldLock.drawFrame(dst, clamp(g.objectStateAt(x, y), 0, 6), px, py, 0)
		case id == 14:
			state := clamp(g.objectStateAt(x, y), 0, 2)
			g.specialContainer.drawAnimationSequenceFrame(dst, 0, state, px, py, 0)
		case id == 32:
			g.drawDiggableFrame(dst, clamp(g.objectStateAt(x, y), 0, len(diggableFrameBounds)-1), px, py)
		case id >= 80:
			g.drawWorldFrame(dst, int(id-80), px, py)
		case id >= 20 && id < 26:
			g.foregroundEffects.drawAnimation(dst, int(id-20), g.tick/2, px, py, 0)
		}
	}
	if playerID < 80 && playerID != original.EmptyRawID {
		idx := x + y*g.rt.Width()
		motion := original.ObjectMotion{}
		if idx >= 0 && idx < len(g.rt.ObjectMotion) {
			motion = g.rt.ObjectMotion[idx]
			if playerID == 0 || playerID == 1 {
				offsetX, offsetY := g.rt.GravityObjectRenderOffset(x, y, g.tick)
				px += offsetX
				py += offsetY
			} else {
				px -= motion.DX * motion.Remaining
				py -= motion.DY * motion.Remaining
			}
		}
		if sourceCellObjectVisible(playerID, foregroundID) {
			switch playerID {
			case 0:
				frame := (g.objectStateAt(x, y) & 0x38) >> 3
				g.drawBoulderFrame(dst, frame, px, py)
			case 1:
				drawSpriteFrame(dst, g.violetGem, sourceGemFrame(g.tick), px, py)
			case 2:
				drawSpriteFrame(dst, g.redDiamond, sourceGemFrame(g.tick), px, py)
			case 4:
				g.drawCenteredModule(dst, g.goldKey, 0, px, py)
			case 5:
				g.drawCenteredModule(dst, g.silverKey, 0, px, py)
			case 6:
				g.commonPickups.drawFrame(dst, 0, px, py, 0)
			case 7:
				g.commonPickups.drawFrame(dst, 1, px, py, 0)
			case 10:
				g.drawDiggableFrame(dst, 0, px, py)
			case 11:
				g.drawCrawler(dst, x, y, px, py)
			case 12:
				g.drawQuotaMarker(dst, px, py)
			case 30:
				state := 0
				idx := x + y*g.rt.Width()
				if idx >= 0 && idx < len(g.rt.ObjectState) {
					state = g.rt.ObjectState[idx]
				}
				frame := clamp((state-1)*7/16, 0, 6)
				g.breakables.drawAnimationSequenceFrame(dst, 0, frame, px, py, 0)
			case 32:
				g.drawHookSegment(dst, x, y, px, py)
			case 24, 31, 33, 41, 79:
				// These are hidden payload/marker objects. Their foreground container,
				// lock, or stage-entry flow owns the visible source sprite.
			case 19:
				g.drawSnake(dst, x, y, px, py)
			case 22:
				drawSpriteFrame(dst, g.hazard, 1, px, py)
				g.flames.drawAnimation(dst, 0, g.tick, px+original.TileSize, py, 0)
			case 23:
				drawSpriteFrame(dst, g.hazard, 0, px, py)
				g.flames.drawAnimation(dst, 0, g.tick, px, py, 1)
			case 43:
				g.drawSnake(dst, x, y, px, py)
			case 48:
				drawRect(dst, px+4, py+5, 16, 14, color.RGBA{110, 112, 118, 245})
				drawRect(dst, px+7, py+3, 10, 4, color.RGBA{178, 184, 190, 245})
				drawRect(dst, px+6, py+10, 12, 2, color.RGBA{70, 74, 82, 245})
			default:
				drawRect(dst, px+7, py+7, 10, 10, color.RGBA{80, 140, 230, 180})
			}
		}
	}
	if foregroundID == 33 {
		frame := hiddenOverlayFrame(g.objectStateAt(x, y))
		g.hiddenOverlay.drawFrame(dst, frame, cellPX, cellPY, 0)
	}
}

func sourceCellObjectVisible(playerID, foregroundID original.RawID) bool {
	if foregroundID != 14 && foregroundID != 33 {
		return true
	}
	switch playerID {
	case 2, 4, 5, 6, 7, 24, 26, 27, 41:
		return false
	default:
		return true
	}
}

func (g *Game) drawSnake(dst *ebiten.Image, x, y, px, py int) {
	state := g.objectStateAt(x, y)
	direction := state & 0x7
	if direction == 0 {
		direction = (state & 0x7000) >> 12
	}
	animation := clamp(direction-1, 0, 3)
	sheet := g.snakes
	if id, _ := g.rt.At(original.PlayerLayer, x, y); id == 43 {
		sheet = g.redSnakes
	}
	sheet.drawAnimationSequenceFrame(dst, animation, g.tick>>1, px, py, 0)
}

func (g *Game) drawCrawler(dst *ebiten.Image, x, y, px, py int) {
	state := g.objectStateAt(x, y)
	frame, offsetX, offsetY, visible := crawlerRenderState(state, g.tick)
	if !visible {
		return
	}
	g.crawler.drawModule(dst, frame, px+offsetX, py+offsetY)
}

func crawlerRenderState(state, sourceTick int) (frame, offsetX, offsetY int, visible bool) {
	phase := (state & 0xf00) >> 8
	if phase >= 4 {
		return 0, 0, 0, false
	}
	frame = (sourceTick >> 1) % 3
	if phase > 0 {
		frame = phase + 2
	}
	direction := state & 0x7
	offsetX, offsetY = 2, 2
	reversed := state&0x10 != 0
	sign := 1
	if reversed {
		sign = -1
	}
	switch direction {
	case 1:
		offsetX += 4 * sign
	case 2:
		offsetY += 4 * sign
	case 3:
		offsetX -= 4 * sign
	case 4:
		offsetY -= 4 * sign
	}
	return clamp(frame, 0, 5), offsetX, offsetY, true
}

func (g *Game) drawQuotaMarker(dst *ebiten.Image, px, py int) {
	g.quota.drawFrame(dst, 0, px, py, 0)
	remaining := max(0, g.rt.BonusRemaining)
	nudge := 0
	if remaining < 10 {
		nudge = g.hud.moduleWidth(0)/2 + 1
	}
	g.hud.drawNumber(dst, remaining, px+19-nudge, py+13)
}

func (g *Game) drawPressureSwitch(dst *ebiten.Image, x, y, px, py int, playerID original.RawID) {
	offset := g.pressureSwitchOffset(x, y, playerID)
	height := 13
	if len(g.pressureSwitch.meta.Modules) > 0 {
		height = g.pressureSwitch.meta.Modules[0].H
	}
	g.pressureSwitch.drawModule(dst, 0, px, py+original.TileSize-height+offset)
}

func (g *Game) pressureSwitchOffset(x, y int, playerID original.RawID) int {
	if playerID == 0 || playerID == 1 || playerID == 8 || playerID == 9 {
		idx := x + y*g.rt.Width()
		remaining := 0
		if idx >= 0 && idx < len(g.rt.ObjectMotion) {
			remaining = g.rt.ObjectMotion[idx].Remaining
		}
		if remaining <= 12 {
			return 12 - remaining
		}
	} else if g.rt.Player == (original.Point{X: x, Y: y}) {
		if g.rt.PlayerMotion.Remaining <= 12 {
			return 12 - g.rt.PlayerMotion.Remaining
		}
	} else if g.rt.Player.Y == y && g.rt.PlayerMotion.Remaining > 12 {
		switch {
		case g.rt.Player.X == x-1 && g.rt.PlayerMotion.DX < 0:
			return g.rt.PlayerMotion.Remaining - 12
		case g.rt.Player.X == x+1 && g.rt.PlayerMotion.DX > 0:
			return g.rt.PlayerMotion.Remaining - 12
		}
	}
	return 0
}

func (g *Game) drawHookSegment(dst *ebiten.Image, x, y, px, py int) {
	idx := x + y*g.rt.Width()
	remaining := 0
	if idx >= 0 && idx < len(g.rt.ObjectMotion) {
		remaining = clamp(g.rt.ObjectMotion[idx].Remaining, 0, original.TileSize)
	}
	state := g.objectStateAt(x, y)
	startX, tipX, module := hookSegmentGeometry(state, remaining, px)
	lineX := min(startX, tipX)
	lineWidth := max(1, max(startX, tipX)-lineX+1)
	drawRect(dst, lineX, py+original.TileSize/2, lineWidth, 1, hookRopeColor)
	if remaining > 0 {
		g.hero.drawModule(dst, module, tipX, py+original.TileSize/2-2)
	}
}

func hookSegmentGeometry(state, remaining, px int) (startX, tipX, module int) {
	remaining = clamp(remaining, 0, original.TileSize)
	startX = px + original.TileSize
	tipX = px + remaining
	module = 1
	if state&1 != 0 {
		startX = px
		tipX = px + original.TileSize - remaining
		module = 0
	}
	return startX, tipX, module
}

func (g *Game) drawCenteredModule(dst *ebiten.Image, sheet *spriteSheet, module, px, py int) {
	if sheet == nil || module < 0 || module >= len(sheet.meta.Modules) {
		return
	}
	bounds := sheet.meta.Modules[module]
	sheet.drawModule(dst, module, px+(original.TileSize-bounds.W)/2, py+(original.TileSize-bounds.H)/2)
}

func sourceGemFrame(sourceTick int) int {
	frame := (sourceTick & 0x3f) >> 1
	if frame >= 4 {
		return 0
	}
	return frame
}

func (g *Game) objectStateAt(x, y int) int {
	idx := x + y*g.rt.Width()
	if idx < 0 || idx >= len(g.rt.ObjectState) {
		return 0
	}
	return g.rt.ObjectState[idx]
}

func (g *Game) objectStateOrBackgroundAt(x, y int) int {
	idx := x + y*g.rt.Width()
	if idx < 0 || idx >= len(g.rt.ObjectState) {
		return 0
	}
	if g.rt.ObjectState[idx] != 0 {
		return g.rt.ObjectState[idx]
	}
	return int(g.rt.Background[idx])
}

func (g *Game) drawWorldFrame(dst *ebiten.Image, frame, px, py int) {
	drawSpriteFrame(dst, g.worldFrames, frame, px, py)
}

func (g *Game) drawBoulderFrame(dst *ebiten.Image, frame, px, py int) {
	drawSpriteFrame(dst, g.boulder, frame, px, py)
}

func (g *Game) drawDiggableFrame(dst *ebiten.Image, frame, px, py int) {
	frame = clamp(frame, 0, len(diggableFrameBounds)-1)
	srcX := framePadding + (frame%frameCols)*(diggableFrameCellW+framePadding)
	srcY := framePadding + (frame/frameCols)*(diggableFrameCellH+framePadding)
	bounds := diggableFrameBounds[frame]
	w := bounds.Dx()
	h := bounds.Dy()
	if g.diggable == nil || srcX+w > g.diggable.Bounds().Dx() || srcY+h > g.diggable.Bounds().Dy() {
		drawRect(dst, px, py, original.TileSize, original.TileSize, color.RGBA{30, 140, 65, 255})
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(px+bounds.Min.X), float64(py+bounds.Min.Y))
	dst.DrawImage(g.diggable.SubImage(image.Rect(srcX, srcY, srcX+w, srcY+h)).(*ebiten.Image), op)
}

func drawSpriteFrame(dst, sheet *ebiten.Image, frame, px, py int) {
	srcX := framePadding + (frame%frameCols)*(original.TileSize+framePadding)
	srcY := framePadding + (frame/frameCols)*(original.TileSize+framePadding)
	if sheet == nil || srcX+original.TileSize > sheet.Bounds().Dx() || srcY+original.TileSize > sheet.Bounds().Dy() {
		drawRect(dst, px, py, original.TileSize, original.TileSize, color.RGBA{120, 40, 160, 255})
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(px), float64(py))
	dst.DrawImage(sheet.SubImage(image.Rect(srcX, srcY, srcX+original.TileSize, srcY+original.TileSize)).(*ebiten.Image), op)
}

func (g *Game) drawPlayer(dst *ebiten.Image, px, py int) {
	if g.rt.InvulnerabilityTicks > 0 && (g.tick>>1)&1 != 0 {
		return
	}
	animation, animationTick := g.heroAnimationState()
	g.hero.drawAnimation(dst, animation, animationTick, px, py, 0)
	if g.rt.ChestOpening && chestRewardIconVisible(g.hero, g.rt.ChestAnimation, g.rt.ChestTicks) {
		g.drawChestRewardIcon(dst, px, py-original.TileSize)
	}
}

func (g *Game) heroAnimationState() (int, int) {
	direction := heroDirection(g.lastDX, g.lastDY)
	animation := direction
	animationTick := g.tick
	if g.rt.HurtTicks > 0 {
		animation = 10
	} else if g.rt.PlayerDead {
		animation = 12
		animationTick = min(14, max(0, deathTransitionTicks-g.rt.DeathTicks))
	} else if g.rt.RecallPending {
		animation = 19
		animationTick = g.rt.RecallTicks
	} else if g.rt.ChestOpening {
		animation = g.rt.ChestAnimation
		if animation != 48 {
			animation = 40
		}
		animationTick = g.rt.ChestTicks
	} else if g.rt.LockOpening {
		animation = g.rt.LockAnimation
		animationTick = g.rt.LockTicks
	} else if g.rt.Hammering {
		animation = g.rt.HammerAnimation
		animationTick = g.rt.HammerTicks
	} else if g.rt.Hooking {
		animation = g.rt.HookAnimation
		animationTick = g.rt.HookTicks
		if duration, ok := g.hero.animationDuration(animation); ok {
			animationTick = min(animationTick, duration-1)
		}
	} else if g.rt.Pushing && g.heroMoveOffset == 0 {
		if g.rt.PushDX > 0 {
			animation = 8
		} else {
			animation = 9
		}
	} else if g.rt.HoldingRock() && g.heroMoveOffset == 0 {
		animation = 11
	} else if g.heroMoveOffset > 0 {
		animation += 4
		animationTick = max(0, g.tick-g.heroMoveStart)
	}
	return animation, animationTick
}

func (g *Game) drawChestRewardEffect(dst *ebiten.Image, px, py int) {
	if !g.rt.ChestOpening {
		return
	}
	animation, ok := chestRewardEffectAnimation(g.rt.ChestRewardID)
	if !ok {
		return
	}
	effectTick, ok := chestRewardEffectTick(g.pickupEffects, animation, g.rt.ChestAnimation, g.rt.ChestTicks)
	if !ok {
		return
	}
	g.pickupEffects.drawAnimation(dst, animation, effectTick, px, py-original.TileSize, 0)
}

func (g *Game) drawChestRewardIcon(dst *ebiten.Image, px, py int) {
	switch g.rt.ChestRewardID {
	case 2:
		drawSpriteFrame(dst, g.redDiamond, 0, px, py)
	case 4:
		g.drawCenteredModule(dst, g.goldKey, 0, px, py)
	case 5:
		g.drawCenteredModule(dst, g.silverKey, 0, px, py)
	case 6:
		g.commonPickups.drawFrame(dst, 0, px, py, 0)
	case 7:
		g.commonPickups.drawFrame(dst, 1, px, py, 0)
	case 24:
		g.tools.drawModule(dst, 0, px, py)
	case 27:
		g.tools.drawModule(dst, 1, px, py)
	case 26:
		g.tools.drawModule(dst, 2, px, py)
	case 41:
		drawSpriteFrame(dst, g.violetGem, 0, px, py)
		g.hud.drawNumber(dst, g.rt.ChestRewardValue, px+original.TileSize, py+14)
	}
}

func (g *Game) drawWorldEffects(dst *ebiten.Image, camX, camY int) {
	for _, effect := range g.worldEffects {
		px := effect.Point.X*original.TileSize - camX
		py := effect.Point.Y*original.TileSize - camY
		g.pickupEffects.drawAnimationSequenceFrame(dst, effect.Animation, effect.Sequence, px, py, 0)
	}
}

func chestRewardIconVisible(hero *spriteSheet, animation, tick int) bool {
	sequence, ok := hero.animationSequenceIndex(animation, tick)
	threshold := chestRewardSequence
	if animation == 48 {
		threshold = chestShortRewardSequence
	}
	return ok && sequence > threshold
}

func chestRewardEffectAnimation(rewardID original.RawID) (int, bool) {
	switch rewardID {
	case 2, 6:
		return 0, true
	case 5:
		return 1, true
	case 4:
		return 2, true
	case 41:
		return 3, true
	case 7:
		return 4, true
	default:
		return 0, false
	}
}

func chestRewardEffectTick(effects *spriteSheet, effectAnimation, heroAnimation, tick int) (int, bool) {
	rewardTick := chestRewardTick
	if heroAnimation == 48 {
		rewardTick = chestShortRewardTick
	}
	effectTick := tick - rewardTick
	duration, ok := effects.animationDuration(effectAnimation)
	return effectTick, ok && effectTick >= 0 && effectTick < duration
}

func drawRect(dst *ebiten.Image, x, y, w, h int, c color.Color) {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(img, op)
}

func heldDirection() (int, int) {
	return heldDirectionWith(ebiten.IsKeyPressed)
}

func heldDirectionWith(pressed func(ebiten.Key) bool) (int, int) {
	shift := shiftPressed(pressed)
	switch {
	case anyKeyPressed(pressed, ebiten.KeyArrowUp, ebiten.KeyNumpad2) || (!shift && pressed(ebiten.KeyDigit2)):
		return 0, -1
	case anyKeyPressed(pressed, ebiten.KeyArrowDown, ebiten.KeyNumpad8) || (!shift && pressed(ebiten.KeyDigit8)):
		return 0, 1
	case anyKeyPressed(pressed, ebiten.KeyArrowLeft, ebiten.KeyNumpad4) || (!shift && pressed(ebiten.KeyDigit4)):
		return -1, 0
	case anyKeyPressed(pressed, ebiten.KeyArrowRight, ebiten.KeyNumpad6) || (!shift && pressed(ebiten.KeyDigit6)):
		return 1, 0
	default:
		return 0, 0
	}
}

func recallPressed() bool {
	return recallPressedWith(inpututil.IsKeyJustPressed, ebiten.IsKeyPressed)
}

func recallPressedWith(justPressed, pressed func(ebiten.Key) bool) bool {
	return justPressed(ebiten.KeyNumpadMultiply) ||
		justPressed(ebiten.KeyBackspace) ||
		justPressed(ebiten.KeyR) ||
		(justPressed(ebiten.KeyDigit8) && shiftPressed(pressed))
}

func shiftPressed(pressed func(ebiten.Key) bool) bool {
	return pressed(ebiten.KeyShift) || pressed(ebiten.KeyShiftLeft) || pressed(ebiten.KeyShiftRight)
}

func centerActionPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDigit5) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpad5) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

func anyKeyPressed(pressed func(ebiten.Key) bool, keys ...ebiten.Key) bool {
	for _, key := range keys {
		if pressed(key) {
			return true
		}
	}
	return false
}

func hiddenOverlayFrame(state int) int {
	return clamp(state, 0, 3)
}

func heroDirection(dx, dy int) int {
	switch {
	case dy < 0:
		return 0
	case dx > 0:
		return 1
	case dy > 0:
		return 2
	case dx < 0:
		return 3
	default:
		return 2
	}
}

func loadTransparentSheet(path string) (*ebiten.Image, error) {
	f, err := os.Open(filepath.Clean(resolvePath(path)))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if uint8(r>>8) == 20 && uint8(g>>8) == 22 && uint8(b>>8) == 28 {
				rgba.SetRGBA(x, y, color.RGBA{})
				continue
			}
			rgba.SetRGBA(x, y, color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)})
		}
	}
	return ebiten.NewImageFromImage(rgba), nil
}

func resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if _, err := os.Stat(path); err == nil {
		return path
	}
	fallback := filepath.Join("..", "..", path)
	if _, err := os.Stat(fallback); err == nil {
		return fallback
	}
	return path
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
