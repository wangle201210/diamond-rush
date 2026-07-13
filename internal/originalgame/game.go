package originalgame

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wangle201210/zskc/internal/original"
)

const (
	desktopActionKeyLabel = "SPACE"
	desktopRecallKeyLabel = "ENTER"
	desktopSkipKeyLabel   = "S"

	angkorWorldFrameSheet      = "decoded/sprites/0/chunk02-frames.png"
	angkorWorldFrameModules    = "decoded/sprites/0/chunk02-modules.png"
	angkorWorldFrameMetadata   = "decoded/sprites/0/chunk02-animations.json"
	angkorBoulderFrameSheet    = "decoded/sprites/0/chunk00-frames.png"
	angkorBoulderModules       = "decoded/sprites/0/chunk00-modules.png"
	angkorBoulderMetadata      = "decoded/sprites/0/chunk00-animations.json"
	angkorDiggableFrameSheet   = "decoded/sprites/0/chunk01-frames.png"
	angkorDiggableModules      = "decoded/sprites/0/chunk01-modules.png"
	angkorDiggableMetadata     = "decoded/sprites/0/chunk01-animations.json"
	angkorFloorSheet           = "decoded/sprites/0/chunk03-modules.png"
	angkorFloorMetadata        = "decoded/sprites/0/chunk03-animations.json"
	bavariaWorldFrameSheet     = "decoded/sprites/1/chunk02-frames.png"
	bavariaWorldFrameModules   = "decoded/sprites/1/chunk02-modules.png"
	bavariaWorldFrameMetadata  = "decoded/sprites/1/chunk02-animations.json"
	bavariaBoulderFrameSheet   = "decoded/sprites/1/chunk00-frames.png"
	bavariaBoulderModules      = "decoded/sprites/1/chunk00-modules.png"
	bavariaBoulderMetadata     = "decoded/sprites/1/chunk00-animations.json"
	bavariaDiggableFrameSheet  = "decoded/sprites/1/chunk01-frames.png"
	bavariaDiggableModules     = "decoded/sprites/1/chunk01-modules.png"
	bavariaDiggableMetadata    = "decoded/sprites/1/chunk01-animations.json"
	bavariaFloorSheet          = "decoded/sprites/1/chunk03-modules.png"
	bavariaFloorMetadata       = "decoded/sprites/1/chunk03-animations.json"
	violetGemSheet             = "decoded/sprites/cm/chunk02-frames.png"
	violetGemModules           = "decoded/sprites/cm/chunk02-modules.png"
	redDiamondSheet            = "decoded/sprites/cm/chunk02-palette01-frames.png"
	redDiamondModules          = "decoded/sprites/cm/chunk02-palette01-modules.png"
	gemMetadata                = "decoded/sprites/cm/chunk02-animations.json"
	checkpointSheet            = "decoded/sprites/cm/chunk06-frames.png"
	checkpointModules          = "decoded/sprites/cm/chunk06-modules.png"
	checkpointMetadata         = "decoded/sprites/cm/chunk06-animations.json"
	quotaModules               = "decoded/sprites/cm/chunk05-modules.png"
	quotaMetadata              = "decoded/sprites/cm/chunk05-animations.json"
	goalSheet                  = "decoded/sprites/cm/chunk00-modules.png"
	goalMetadata               = "decoded/sprites/cm/chunk00-animations.json"
	doorModuleSheet            = "decoded/sprites/cm/chunk01-modules.png"
	doorMetadata               = "decoded/sprites/cm/chunk01-animations.json"
	snakeSheet                 = "decoded/sprites/gen1/chunk05-frames.png"
	snakeModuleSheet           = "decoded/sprites/gen1/chunk05-modules.png"
	snakeMetadata              = "decoded/sprites/gen1/chunk05-animations.json"
	redSnakeSheet              = "decoded/sprites/gen1/chunk05-palette01-frames.png"
	redSnakeModuleSheet        = "decoded/sprites/gen1/chunk05-palette01-modules.png"
	crawlerModules             = "decoded/sprites/gen1/chunk04-modules.png"
	crawlerMetadata            = "decoded/sprites/gen1/chunk04-animations.json"
	commonPickupModules        = "decoded/sprites/cm/chunk04-modules.png"
	commonPickupMetadata       = "decoded/sprites/cm/chunk04-animations.json"
	frozenVioletSheet          = "decoded/sprites/gen0/chunk01-frames.png"
	frozenVioletModules        = "decoded/sprites/gen0/chunk01-modules.png"
	frozenVioletMetadata       = "decoded/sprites/gen0/chunk01-animations.json"
	frozenSnakeSheet           = "decoded/sprites/gen1/chunk06-frames.png"
	frozenSnakeModules         = "decoded/sprites/gen1/chunk06-modules.png"
	frozenSnakeMetadata        = "decoded/sprites/gen1/chunk06-animations.json"
	breakableSheet             = "decoded/sprites/gen0/chunk07-frames.png"
	breakableModules           = "decoded/sprites/gen0/chunk07-modules.png"
	breakableMetadata          = "decoded/sprites/gen0/chunk07-animations.json"
	goldLockSheet              = "decoded/sprites/gen2/chunk08-frames.png"
	goldLockModules            = "decoded/sprites/gen2/chunk08-modules.png"
	goldLockMetadata           = "decoded/sprites/gen2/chunk08-animations.json"
	silverLockSheet            = "decoded/sprites/gen2/chunk08-palette01-frames.png"
	silverLockModules          = "decoded/sprites/gen2/chunk08-palette01-modules.png"
	foregroundEffectSheet      = "decoded/sprites/gen0/chunk04-frames.png"
	foregroundEffectModules    = "decoded/sprites/gen0/chunk04-modules.png"
	foregroundEffectMetadata   = "decoded/sprites/gen0/chunk04-animations.json"
	hiddenOverlaySheet         = "decoded/sprites/gen3/chunk03-frames.png"
	hiddenOverlayModules       = "decoded/sprites/gen3/chunk03-modules.png"
	hiddenOverlayMetadata      = "decoded/sprites/gen3/chunk03-animations.json"
	specialContainerSheet      = "decoded/sprites/gen2/chunk02-frames.png"
	specialContainerModules    = "decoded/sprites/gen2/chunk02-modules.png"
	specialContainerMetadata   = "decoded/sprites/gen2/chunk02-animations.json"
	goldKeyModules             = "decoded/sprites/gen0/chunk02-modules.png"
	silverKeyModules           = "decoded/sprites/gen0/chunk02-palette01-modules.png"
	keyMetadata                = "decoded/sprites/gen0/chunk02-animations.json"
	toolModules                = "decoded/sprites/gen1/chunk09-modules.png"
	toolMetadata               = "decoded/sprites/gen1/chunk09-animations.json"
	toolPromptModules          = "decoded/sprites/gen2/chunk00-modules.png"
	toolPromptMetadata         = "decoded/sprites/gen2/chunk00-animations.json"
	compassPickupModules       = "decoded/sprites/gen3/chunk01-modules.png"
	compassPickupMetadata      = "decoded/sprites/gen3/chunk01-animations.json"
	tutorialRecallHintSheet    = "decoded/sprites/gen3/chunk00-frames.png"
	tutorialRecallHintModules  = "decoded/sprites/gen3/chunk00-modules.png"
	tutorialRecallHintMetadata = "decoded/sprites/gen3/chunk00-animations.json"
	worldMapIconSheet          = "decoded/sprites/ms/chunk00-frames.png"
	worldMapIconModules        = "decoded/sprites/ms/chunk00-modules.png"
	worldMapIconMetadata       = "decoded/sprites/ms/chunk00-animations.json"
	worldMapGroundSheet        = "decoded/sprites/ms/chunk01-frames.png"
	worldMapGroundModules      = "decoded/sprites/ms/chunk01-modules.png"
	worldMapGroundMetadata     = "decoded/sprites/ms/chunk01-animations.json"
	worldMapHeaderSheet        = "decoded/sprites/ms/chunk02-frames.png"
	worldMapHeaderModules      = "decoded/sprites/ms/chunk02-modules.png"
	worldMapHeaderMetadata     = "decoded/sprites/ms/chunk02-animations.json"
	bavariaMapHeaderSheet      = "decoded/sprites/ms/chunk03-frames.png"
	bavariaMapHeaderModules    = "decoded/sprites/ms/chunk03-modules.png"
	bavariaMapHeaderMetadata   = "decoded/sprites/ms/chunk03-animations.json"
	pickupEffectSheet          = "decoded/sprites/cm/chunk07-frames.png"
	pickupEffectModules        = "decoded/sprites/cm/chunk07-modules.png"
	pickupEffectMetadata       = "decoded/sprites/cm/chunk07-animations.json"
	resultSparkModules         = "decoded/sprites/cm/chunk04-modules.png"
	resultSparkMetadata        = "decoded/sprites/cm/chunk04-animations.json"
	resultMedalModules         = "decoded/sprites/ui/chunk04-modules.png"
	resultMedalMetadata        = "decoded/sprites/ui/chunk04-animations.json"
	hazardEmitterSheet         = "decoded/sprites/gen0/chunk09-frames.png"
	hazardEmitterModules       = "decoded/sprites/gen0/chunk09-modules.png"
	hazardEmitterMetadata      = "decoded/sprites/gen0/chunk09-animations.json"
	hazardFlameSheet           = "decoded/sprites/gen1/chunk00-frames.png"
	hazardFlameModuleSheet     = "decoded/sprites/gen1/chunk00-modules.png"
	hazardFlameMetadata        = "decoded/sprites/gen1/chunk00-animations.json"
	pressureSwitchModules      = "decoded/sprites/gen2/chunk09-modules.png"
	pressureSwitchMetadata     = "decoded/sprites/gen2/chunk09-animations.json"
	fallingFireSheet           = "decoded/sprites/mm0/chunk00-frames.png"
	fallingFireModules         = "decoded/sprites/mm0/chunk00-modules.png"
	fallingFireMetadata        = "decoded/sprites/mm0/chunk00-animations.json"
	fallingTorchSheet          = "decoded/sprites/mm0/chunk01-frames.png"
	fallingTorchModules        = "decoded/sprites/mm0/chunk01-modules.png"
	fallingTorchMetadata       = "decoded/sprites/mm0/chunk01-animations.json"
	fallingDebrisModules       = "decoded/sprites/mm0/chunk02-modules.png"
	fallingDebrisMetadata      = "decoded/sprites/mm0/chunk02-animations.json"
	anacondaModules            = "decoded/sprites/b0/chunk00-modules.png"
	anacondaMetadata           = "decoded/sprites/b0/chunk00-animations.json"
	anacondaPlatformModules    = "decoded/sprites/b0/chunk01-modules.png"
	anacondaPlatformMetadata   = "decoded/sprites/b0/chunk01-animations.json"
	angkorSealModules          = "decoded/sprites/mmv/chunk03-modules.png"
	angkorSealMetadata         = "decoded/sprites/mmv/chunk03-animations.json"
	tutorialSealModules        = "decoded/sprites/mmv/chunk00-modules.png"
	tutorialSealMetadata       = "decoded/sprites/mmv/chunk00-animations.json"
	siberiaSealModules         = "decoded/sprites/mmv/chunk01-modules.png"
	siberiaSealMetadata        = "decoded/sprites/mmv/chunk01-animations.json"
	bavariaSealModules         = "decoded/sprites/mmv/chunk02-modules.png"
	bavariaSealMetadata        = "decoded/sprites/mmv/chunk02-animations.json"
	sealArrowModules           = "decoded/sprites/mmv/chunk04-modules.png"
	sealArrowMetadata          = "decoded/sprites/mmv/chunk04-animations.json"
	softkeyModules             = "decoded/sprites/ui/chunk03-modules.png"
	softkeyMetadata            = "decoded/sprites/ui/chunk03-animations.json"
	splashBackgroundImage      = "decoded/sprites/splash/background.png"
	splashLogoImage            = "decoded/sprites/splash/logo.png"
	splashCopyrightImage       = "decoded/sprites/splash/copyright.png"
	demoUIModules              = "decoded/sprites/demoui/chunk00-modules.png"
	demoUIBlueModules          = "decoded/sprites/demoui/chunk00-palette01-modules.png"
	demoUIMetadata             = "decoded/sprites/demoui/chunk00-animations.json"
	tutorialFaceModules        = "decoded/sprites/tutorial/demoSpr/sprite00-modules.png"
	tutorialFaceMetadata       = "decoded/sprites/tutorial/demoSpr/sprite00-animations.json"
	tutorialMarkModules        = "decoded/sprites/tutorial/demoSpr/sprite01-modules.png"
	tutorialMarkMetadata       = "decoded/sprites/tutorial/demoSpr/sprite01-animations.json"
	tutorialPortraitModules    = "decoded/sprites/tutorial/demoSpr/sprite02-modules.png"
	tutorialPortraitMetadata   = "decoded/sprites/tutorial/demoSpr/sprite02-animations.json"
	hudSheet                   = "decoded/sprites/ui/chunk02-frames.png"
	hudModuleSheet             = "decoded/sprites/ui/chunk02-modules.png"
	hudMetadata                = "decoded/sprites/ui/chunk02-animations.json"
	heroFrameSheet             = "decoded/sprites/o/chunk00-frames.png"
	heroModuleSheet            = "decoded/sprites/o/chunk00-modules.png"
	heroMetadata               = "decoded/sprites/o/chunk00-animations.json"
	fontSmallSheet             = "decoded/fonts/freej2me-small.png"
	fontSmallMetadata          = "decoded/fonts/freej2me-small.json"
	fontMediumSheet            = "decoded/fonts/freej2me-medium.png"
	fontMediumMetadata         = "decoded/fonts/freej2me-medium.json"
	originalAudioDir           = "decoded/audio"
	defaultWorldDir            = "decoded/world0"
	resultLoadingSteps         = 12
	sealLoadingSteps           = 11
	resultTitleTicks           = 40
	resultGemMinimumTicks      = 40
	resultRedDiamondTicks      = 40
	resultHitTicks             = 10
	resultRetryTicks           = 10
	resultStageTitleY          = 23
	resultCompleteY            = 42
	resultVioletLabelY         = 75
	resultVioletCountY         = 91
	resultRedLabelY            = 131
	resultRedCountY            = 147
	resultHitsIconY            = 191
	resultHitsLabelY           = 187
	resultHitsCountY           = 203
	resultRetriesIconY         = 243
	resultRetriesLabelY        = 243
	resultRetriesCountY        = 259
	stageIntroDuration         = 60
	secretExitDuration         = 30
	deathTransitionTicks       = 80
	chestRewardTick            = 39
	chestRewardSequence        = 13
	chestShortRewardTick       = 23
	chestShortRewardSequence   = 6
	sourceTPS                  = 20
	sourceHeroTurnStartOffset  = 18
	sourceHeroTurnStep         = 6
	framePadding               = 2
	frameCols                  = 16
	playfieldHeight            = 240
	playfieldTop               = 40
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

type worldEffect struct {
	Point     original.Point
	Animation int
	Sequence  int
}

type gameMode int

const (
	gameModeStage gameMode = iota
	gameModeWorldMap
	gameModeWorldSelect
	gameModeStartMenu
)

type Game struct {
	pack                   *original.WorldPack
	rt                     *original.Runtime
	worldFrames            *spriteSheet
	boulder                *spriteSheet
	diggable               *spriteSheet
	floor                  *spriteSheet
	violetGem              *spriteSheet
	redDiamond             *spriteSheet
	checkpoint             *spriteSheet
	quota                  *spriteSheet
	goal                   *spriteSheet
	door                   *spriteSheet
	hazard                 *spriteSheet
	snakes                 *spriteSheet
	redSnakes              *spriteSheet
	crawler                *spriteSheet
	commonPickups          *spriteSheet
	frozenViolet           *spriteSheet
	frozenSnake            *spriteSheet
	breakables             *spriteSheet
	goldLock               *spriteSheet
	silverLock             *spriteSheet
	foregroundEffects      *spriteSheet
	hiddenOverlay          *spriteSheet
	specialContainer       *spriteSheet
	goldKey                *spriteSheet
	silverKey              *spriteSheet
	tools                  *spriteSheet
	toolPrompt             *spriteSheet
	compassPickup          *spriteSheet
	tutorialRecallHint     *spriteSheet
	worldMapIcons          *spriteSheet
	worldMapGround         *spriteSheet
	worldMapHeader         *spriteSheet
	pickupEffects          *spriteSheet
	resultSpark            *spriteSheet
	resultMedal            *spriteSheet
	flames                 *spriteSheet
	pressureSwitch         *spriteSheet
	fallingFire            *spriteSheet
	fallingTorches         *spriteSheet
	fallingDebris          *spriteSheet
	anaconda               *spriteSheet
	anacondaPlatform       *spriteSheet
	angkorSeal             *spriteSheet
	tutorialSeal           *spriteSheet
	siberiaSeal            *spriteSheet
	bavariaSeal            *spriteSheet
	bavaria                bavariaSpriteSet
	sealArrow              *spriteSheet
	softkeys               *spriteSheet
	splashBackground       *ebiten.Image
	splashLogo             *ebiten.Image
	splashCopyright        *ebiten.Image
	demoUI                 *spriteSheet
	demoUIBlue             *spriteSheet
	tutorialFaces          *spriteSheet
	tutorialMarks          *spriteSheet
	tutorialPortrait       *spriteSheet
	hud                    *spriteSheet
	hero                   *spriteSheet
	fontSmall              *bitmapFont
	fontMedium             *bitmapFont
	worldCanvas            *ebiten.Image
	worldDir               string
	worldRoot              string
	worldIndex             int
	worldMap               *worldMapData
	mode                   gameMode
	stageIndex             int
	worldMapLoadingStep    int
	worldMapSelectedStage  int
	worldMapTravelFrom     int
	worldMapTravelTo       int
	worldMapTravelTick     int
	pendingMapTarget       int
	message                string
	tick                   int
	worldDone              bool
	secretExitActive       bool
	secretExitTicks        int
	sealExitActive         bool
	sealExitTicks          int
	worldSelectPosition    int
	worldSelectArrowX      int
	worldSelectArrowY      int
	worldSelectTargetX     int
	worldSelectTargetY     int
	worldSelectMoveTick    int
	worldSelectArrowTick   int
	worldSelectIncoming    int
	worldSelectRelicX      int
	worldSelectRelicY      int
	worldSelectFlashTick   int
	worldSelectEffectTick  int
	worldSelectUnlocking   int
	worldSelectUnlockTick  int
	worldSelectUnlockFlash int
	startMenuHasProgress   bool
	startMenuSelection     int
	startMenuConfirmNew    bool
	startMenuConfirmChoice int
	resultPhase            int
	resultPhaseTicks       int
	resultLoadingStep      int
	resultAwards           byte
	resultNewAwards        byte
	introTicks             int
	lastDX                 int
	lastDY                 int
	heroMoveStart          int
	heroMoveOffset         int
	heroTurnOffset         int
	entranceSteps          int
	checkpointBannerUntil  int
	compassDirection       int
	waterSpecialFrame      int
	cameraX                int
	cameraY                int
	demoCameraLocked       bool
	demoCameraStartX       int
	demoCameraStartY       int
	demoCameraPhase        int
	worldEffects           []worldEffect
	sounds                 *originalSounds
	progressPath           string
	progress               originalProgress
	startupWorldOverridden bool
}

func Run() error {
	startupWorld, worldOverridden, err := requestedStartupWorld()
	if err != nil {
		return err
	}
	worldDir := filepath.Join(filepath.Dir(defaultWorldDir), fmt.Sprintf("world%d", startupWorld))
	g, err := New(worldDir)
	if err != nil {
		return err
	}
	progressPath := originalProgressPath()
	hasProgress, err := originalProgressExists(progressPath)
	if err != nil {
		return err
	}
	if err := g.enableProgress(progressPath); err != nil {
		return err
	}
	if worldOverridden {
		g.startupWorldOverridden = true
		g.progress.WorldUnlocked[startupWorld] = true
		g.progress.unlockStageForWorld(startupWorld, 0)
		g.progress.LastWorld = startupWorld
		g.progress = g.progress.normalized()
	}
	if stageText := os.Getenv("ORIGINALRUSH_STAGE"); stageText != "" {
		stageNumber, err := strconv.Atoi(stageText)
		if err != nil || stageNumber < 1 || stageNumber > len(g.pack.Stages) || !worldStageImplemented(g.worldIndex, stageNumber-1) {
			return fmt.Errorf("invalid ORIGINALRUSH_STAGE %q", stageText)
		}
		g.progress.unlockStageForWorld(g.worldIndex, stageNumber-1)
		g.loadStage(stageNumber - 1)
	} else if worldOverridden {
		g.stageIndex = g.highestUnlockedMapStageForWorld(startupWorld)
		g.enterWorldMap()
	} else {
		g.enterStartMenu(hasProgress)
	}
	g.sounds.Enable()
	if g.mode == gameModeStartMenu {
		g.sounds.Play(original.SoundTitleMusic)
	} else {
		g.sounds.Play(worldMusic(g.worldIndex))
	}
	defer g.sounds.Stop()
	ebiten.SetTPS(sourceTPS)
	ebiten.SetWindowTitle(fmt.Sprintf("Diamond Rush Original Runtime - %s World %d", worldName(g.worldIndex), g.worldIndex))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	options := &ebiten.RunGameOptions{}
	if runtime.GOOS == "darwin" {
		// Metal can occasionally fail to acquire a CAMetalLayer drawable after
		// a window transition, leaving the native game window permanently black.
		options.GraphicsLibrary = ebiten.GraphicsLibraryOpenGL
	}
	return ebiten.RunGameWithOptions(g, options)
}

func requestedStartupWorld() (world int, overridden bool, err error) {
	value := strings.TrimSpace(strings.ToLower(os.Getenv("ORIGINALRUSH_WORLD")))
	if value == "" {
		return original.WorldAngkor, false, nil
	}
	switch value {
	case "0", "angkor", "angkor-wat":
		return original.WorldAngkor, true, nil
	case "1", "bavaria":
		return original.WorldBavaria, true, nil
	default:
		return 0, false, fmt.Errorf("invalid ORIGINALRUSH_WORLD %q", value)
	}
}

func New(worldDir string) (*Game, error) {
	resolvedWorldDir := resolvePath(worldDir)
	pack, err := original.LoadWorldDir(resolvedWorldDir)
	if err != nil {
		return nil, err
	}
	visuals := worldVisualDefinitionFor(pack.World)
	frames, err := loadSpriteSheetWithModules(visuals.frames, visuals.frameModules, visuals.frameMetadata)
	if err != nil {
		return nil, fmt.Errorf("load world frame sheet: %w", err)
	}
	boulder, err := loadSpriteSheetWithModules(visuals.boulder, visuals.boulderModules, visuals.boulderMetadata)
	if err != nil {
		return nil, fmt.Errorf("load boulder frame sheet: %w", err)
	}
	diggable, err := loadSpriteSheetWithModules(visuals.diggable, visuals.diggableModules, visuals.diggableMetadata)
	if err != nil {
		return nil, fmt.Errorf("load diggable frame sheet: %w", err)
	}
	floor, err := loadModuleSpriteSheet(visuals.floor, visuals.floorMetadata)
	if err != nil {
		return nil, fmt.Errorf("load floor sheet: %w", err)
	}
	violetGem, err := loadSpriteSheetWithModules(violetGemSheet, violetGemModules, gemMetadata)
	if err != nil {
		return nil, fmt.Errorf("load violet gem sheet: %w", err)
	}
	redDiamond, err := loadSpriteSheetWithModules(redDiamondSheet, redDiamondModules, gemMetadata)
	if err != nil {
		return nil, fmt.Errorf("load red diamond sheet: %w", err)
	}
	checkpoint, err := loadSpriteSheetWithModules(checkpointSheet, checkpointModules, checkpointMetadata)
	if err != nil {
		return nil, fmt.Errorf("load checkpoint sheet: %w", err)
	}
	quota, err := loadModuleSpriteSheet(quotaModules, quotaMetadata)
	if err != nil {
		return nil, fmt.Errorf("load quota marker: %w", err)
	}
	goal, err := loadModuleSpriteSheet(goalSheet, goalMetadata)
	if err != nil {
		return nil, fmt.Errorf("load goal sheet: %w", err)
	}
	door, err := loadModuleSpriteSheet(doorModuleSheet, doorMetadata)
	if err != nil {
		return nil, fmt.Errorf("load door sprite: %w", err)
	}
	hazard, err := loadSpriteSheetWithModules(hazardEmitterSheet, hazardEmitterModules, hazardEmitterMetadata)
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
	commonPickups, err := loadModuleSpriteSheet(commonPickupModules, commonPickupMetadata)
	if err != nil {
		return nil, fmt.Errorf("load common pickups: %w", err)
	}
	frozenViolet, err := loadSpriteSheetWithModules(frozenVioletSheet, frozenVioletModules, frozenVioletMetadata)
	if err != nil {
		return nil, fmt.Errorf("load frozen violet sprite: %w", err)
	}
	frozenSnake, err := loadSpriteSheetWithModules(frozenSnakeSheet, frozenSnakeModules, frozenSnakeMetadata)
	if err != nil {
		return nil, fmt.Errorf("load frozen snake sprite: %w", err)
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
	toolPrompt, err := loadModuleSpriteSheet(toolPromptModules, toolPromptMetadata)
	if err != nil {
		return nil, fmt.Errorf("load source tool prompt: %w", err)
	}
	compassPickup, err := loadModuleSpriteSheet(compassPickupModules, compassPickupMetadata)
	if err != nil {
		return nil, fmt.Errorf("load compass pickup: %w", err)
	}
	tutorialRecallHint, err := loadSpriteSheetWithModules(tutorialRecallHintSheet, tutorialRecallHintModules, tutorialRecallHintMetadata)
	if err != nil {
		return nil, fmt.Errorf("load tutorial recall hint: %w", err)
	}
	worldMapIcons, err := loadSpriteSheetWithModules(worldMapIconSheet, worldMapIconModules, worldMapIconMetadata)
	if err != nil {
		return nil, fmt.Errorf("load world-map icons: %w", err)
	}
	worldMapGround, err := loadSpriteSheetWithModules(worldMapGroundSheet, worldMapGroundModules, worldMapGroundMetadata)
	if err != nil {
		return nil, fmt.Errorf("load world-map ground: %w", err)
	}
	worldMapHeader, err := loadSpriteSheetWithModules(visuals.mapHeaderSheet, visuals.mapHeaderModules, visuals.mapHeaderMetadata)
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
	fallingFire, err := loadSpriteSheetWithModules(fallingFireSheet, fallingFireModules, fallingFireMetadata)
	if err != nil {
		return nil, fmt.Errorf("load falling-torches fire: %w", err)
	}
	fallingTorches, err := loadSpriteSheetWithModules(fallingTorchSheet, fallingTorchModules, fallingTorchMetadata)
	if err != nil {
		return nil, fmt.Errorf("load falling torches: %w", err)
	}
	fallingDebris, err := loadModuleSpriteSheet(fallingDebrisModules, fallingDebrisMetadata)
	if err != nil {
		return nil, fmt.Errorf("load falling-torches debris: %w", err)
	}
	anaconda, err := loadModuleSpriteSheet(anacondaModules, anacondaMetadata)
	if err != nil {
		return nil, fmt.Errorf("load Great Anaconda sprite: %w", err)
	}
	anacondaPlatform, err := loadModuleSpriteSheet(anacondaPlatformModules, anacondaPlatformMetadata)
	if err != nil {
		return nil, fmt.Errorf("load Great Anaconda platforms: %w", err)
	}
	angkorSeal, err := loadModuleSpriteSheet(angkorSealModules, angkorSealMetadata)
	if err != nil {
		return nil, fmt.Errorf("load Angkor seal: %w", err)
	}
	tutorialSeal, err := loadModuleSpriteSheet(tutorialSealModules, tutorialSealMetadata)
	if err != nil {
		return nil, fmt.Errorf("load tutorial seal: %w", err)
	}
	siberiaSeal, err := loadModuleSpriteSheet(siberiaSealModules, siberiaSealMetadata)
	if err != nil {
		return nil, fmt.Errorf("load Siberia seal: %w", err)
	}
	bavariaSeal, err := loadModuleSpriteSheet(bavariaSealModules, bavariaSealMetadata)
	if err != nil {
		return nil, fmt.Errorf("load Bavaria seal: %w", err)
	}
	bavaria, err := loadBavariaSpriteSet()
	if err != nil {
		return nil, err
	}
	sealArrow, err := loadModuleSpriteSheet(sealArrowModules, sealArrowMetadata)
	if err != nil {
		return nil, fmt.Errorf("load seal selector arrow: %w", err)
	}
	softkeys, err := loadModuleSpriteSheet(softkeyModules, softkeyMetadata)
	if err != nil {
		return nil, fmt.Errorf("load seal selector softkeys: %w", err)
	}
	splashBackground, err := loadTransparentSheet(splashBackgroundImage)
	if err != nil {
		return nil, fmt.Errorf("load title background: %w", err)
	}
	splashLogo, err := loadTransparentSheet(splashLogoImage)
	if err != nil {
		return nil, fmt.Errorf("load title logo: %w", err)
	}
	splashCopyright, err := loadTransparentSheet(splashCopyrightImage)
	if err != nil {
		return nil, fmt.Errorf("load title copyright: %w", err)
	}
	demoUI, err := loadModuleSpriteSheet(demoUIModules, demoUIMetadata)
	if err != nil {
		return nil, fmt.Errorf("load source panel border: %w", err)
	}
	demoUIBlue, err := loadModuleSpriteSheet(demoUIBlueModules, demoUIMetadata)
	if err != nil {
		return nil, fmt.Errorf("load blue source panel border: %w", err)
	}
	tutorialFaces, err := loadModuleSpriteSheet(tutorialFaceModules, tutorialFaceMetadata)
	if err != nil {
		return nil, fmt.Errorf("load tutorial faces: %w", err)
	}
	tutorialMarks, err := loadModuleSpriteSheet(tutorialMarkModules, tutorialMarkMetadata)
	if err != nil {
		return nil, fmt.Errorf("load tutorial punctuation: %w", err)
	}
	tutorialPortrait, err := loadModuleSpriteSheet(tutorialPortraitModules, tutorialPortraitMetadata)
	if err != nil {
		return nil, fmt.Errorf("load tutorial portrait: %w", err)
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
		pack:                pack,
		rt:                  rt,
		worldFrames:         frames,
		boulder:             boulder,
		diggable:            diggable,
		floor:               floor,
		violetGem:           violetGem,
		redDiamond:          redDiamond,
		checkpoint:          checkpoint,
		quota:               quota,
		goal:                goal,
		door:                door,
		hazard:              hazard,
		snakes:              snakes,
		redSnakes:           redSnakes,
		crawler:             crawler,
		commonPickups:       commonPickups,
		frozenViolet:        frozenViolet,
		frozenSnake:         frozenSnake,
		breakables:          breakables,
		goldLock:            goldLock,
		silverLock:          silverLock,
		foregroundEffects:   foregroundEffects,
		hiddenOverlay:       hiddenOverlay,
		specialContainer:    specialContainer,
		goldKey:             goldKey,
		silverKey:           silverKey,
		tools:               tools,
		toolPrompt:          toolPrompt,
		compassPickup:       compassPickup,
		tutorialRecallHint:  tutorialRecallHint,
		worldMapIcons:       worldMapIcons,
		worldMapGround:      worldMapGround,
		worldMapHeader:      worldMapHeader,
		pickupEffects:       pickupEffects,
		resultSpark:         resultSpark,
		resultMedal:         resultMedal,
		flames:              flames,
		pressureSwitch:      pressureSwitch,
		fallingFire:         fallingFire,
		fallingTorches:      fallingTorches,
		fallingDebris:       fallingDebris,
		anaconda:            anaconda,
		anacondaPlatform:    anacondaPlatform,
		angkorSeal:          angkorSeal,
		tutorialSeal:        tutorialSeal,
		siberiaSeal:         siberiaSeal,
		bavariaSeal:         bavariaSeal,
		bavaria:             bavaria,
		sealArrow:           sealArrow,
		softkeys:            softkeys,
		splashBackground:    splashBackground,
		splashLogo:          splashLogo,
		splashCopyright:     splashCopyright,
		demoUI:              demoUI,
		demoUIBlue:          demoUIBlue,
		tutorialFaces:       tutorialFaces,
		tutorialMarks:       tutorialMarks,
		tutorialPortrait:    tutorialPortrait,
		hud:                 hud,
		hero:                hero,
		fontSmall:           fontSmall,
		fontMedium:          fontMedium,
		sounds:              sounds,
		worldCanvas:         ebiten.NewImage(original.ScreenWidth, playfieldHeight),
		worldDir:            resolvedWorldDir,
		worldRoot:           filepath.Dir(resolvedWorldDir),
		worldIndex:          pack.World,
		worldMap:            worldMap,
		pendingMapTarget:    -1,
		worldSelectIncoming: -1,
		message:             fmt.Sprintf("Original %s World %d runtime", worldName(pack.World), pack.World),
		lastDX:              1,
		entranceSteps:       rt.EntranceScrollX,
		progress:            newOriginalProgress(),
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
	if g.mode == gameModeStartMenu {
		_, dy := justPressedDirection()
		g.updateStartMenu(centerActionPressed(), dy)
		return nil
	}
	if g.mode == gameModeWorldMap {
		g.updateWorldMap(centerActionPressed())
		return nil
	}
	if g.mode == gameModeWorldSelect {
		dx, dy := justPressedDirection()
		g.updateWorldSelect(centerActionPressed(), dx, dy)
		return nil
	}
	if g.sealExitActive {
		g.sealExitTicks++
		if g.sealExitTicks >= sealLoadingSteps {
			g.sealExitActive = false
			g.enterWorldSelect(g.worldIndex)
		}
		return nil
	}
	if g.secretExitActive {
		g.secretExitTicks++
		if g.secretExitTicks > secretExitDuration {
			g.secretExitActive = false
			g.enterWorldMap()
		}
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
	if g.heroTurnOffset > 0 && !g.rt.CanAcceptInput() {
		g.setHeroTurnOffset(0)
	}
	if g.rt.TutorialComplete {
		g.finishTutorial()
		return nil
	}
	if g.rt.SealStageComplete {
		g.beginSealExit()
		g.updateCamera()
		return nil
	}
	checkpointRestored := (recallPendingAtFrameStart && !g.rt.RecallPending) ||
		(deadAtFrameStart && !g.rt.PlayerDead)
	if checkpointRestored {
		g.resetHeroFacing()
	}
	if checkpointRestored ||
		(chestOpeningAtFrameStart && !g.rt.ChestOpening) ||
		(lockOpeningAtFrameStart && !g.rt.LockOpening) {
		g.updateCamera()
		return nil
	}
	if g.rt.TutorialScriptActive {
		if tutorialSkipPressed() {
			g.rt.SkipTutorialScript()
		} else if _, ok := g.rt.TutorialPrompt(); ok && centerActionPressed() {
			g.rt.AdvanceTutorialPrompt()
		}
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
	if g.rt.TutorialScriptActive {
		g.updateCamera()
		return nil
	}
	if g.rt.CanAcceptInput() && g.heroTurnOffset == 0 && centerActionPressed() {
		if g.rt.IsCheckpoint(g.rt.Player.X, g.rt.Player.Y) {
			if g.rt.ResetCheckpoint() {
				g.resetHeroFacing()
				g.message = "checkpoint reset"
			}
		} else if g.rt.UsesSwimmingAnimationAt(g.rt.Player.X, g.rt.Player.Y) {
			// The JAR ignores tool input while layer-zero water uses its
			// surface/swimming shapes (7 or 8).
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
			g.setHeroTurnOffset(0)
			if onCheckpoint {
				g.resetHeroFacing()
			}
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
	if g.heroTurnOffset > 0 {
		g.advanceHeroTurn()
		g.updateCamera()
		return nil
	}
	dx, dy := heldDirection()
	if dx != 0 || dy != 0 {
		if g.rt.CanAcceptInput() {
			g.handlePlayerDirection(dx, dy)
		}
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
	g.rt.SetViewport(g.cameraX, g.cameraY)
	result := g.rt.TickSourceFrame(8, g.tick, flameReach)
	for _, point := range result.VioletPickups {
		g.worldEffects = append(g.worldEffects, worldEffect{Point: point, Animation: 3})
	}
	if result.TutorialSealActivated {
		g.worldEffects = append(g.worldEffects, worldEffect{Point: g.rt.Player, Animation: 5})
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
	case result.RisingFireHits > 0:
		g.message = "rising fire reached hero"
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

func (g *Game) handlePlayerDirection(dx, dy int) bool {
	if g == nil || g.rt == nil || (dx == 0 && dy == 0) {
		return false
	}
	if dx != g.lastDX || dy != g.lastDY {
		g.lastDX = dx
		g.lastDY = dy
		g.rt.SetPlayerFacing(dx, dy)
		g.setHeroTurnOffset(sourceHeroTurnStartOffset)
		g.rt.ResetPushAttempt()
		return false
	}
	return g.startPlayerMove(dx, dy)
}

func (g *Game) advanceHeroTurn() bool {
	if g == nil || g.heroTurnOffset <= 0 {
		return false
	}
	g.setHeroTurnOffset(g.heroTurnOffset - sourceHeroTurnStep)
	return true
}

func (g *Game) setHeroTurnOffset(offset int) {
	if g == nil {
		return
	}
	g.heroTurnOffset = max(0, offset)
	if g.rt != nil {
		g.rt.SetPlayerTurnOffset(g.heroTurnOffset)
	}
}

func (g *Game) resetHeroFacing() {
	if g == nil {
		return
	}
	g.lastDX = 1
	g.lastDY = 0
	if g.rt != nil {
		g.rt.SetPlayerFacing(1, 0)
	}
	g.setHeroTurnOffset(0)
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
	if g.worldDone || g.secretExitActive {
		g.message = fmt.Sprintf("%s stage %02d complete", worldName(g.worldIndex), g.stageIndex+1)
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
		g.message = fmt.Sprintf("%s stage %02d exit", worldName(g.worldIndex), g.stageIndex+1)
		return
	}
	if g.rt.GoalExitSecret && g.stageIndex < worldFirstSecretStage(g.worldIndex) {
		target := g.stageIndex
		if next, ok := g.worldMap.exitTarget(g.stageIndex, true); ok {
			target = next
		}
		g.pendingMapTarget = target
		g.progress.recordSecretExit(g.stageIndex, target, g.rt)
		if g.progressPath != "" {
			if err := saveOriginalProgress(g.progressPath, g.progress); err != nil {
				g.message = err.Error()
			}
		}
		g.secretExitActive = true
		g.secretExitTicks = 0
		g.message = "secret path unlocked"
		return
	}
	g.worldDone = true
	g.beginStageResults()
	g.message = fmt.Sprintf("%s stage %02d complete", worldName(g.worldIndex), g.stageIndex+1)
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
		g.rt.SetPlayerFacing(motion.DX, motion.DY)
	}
}

func (g *Game) beginStageResults() {
	g.resultPhase = resultPhaseLoading
	g.resultPhaseTicks = 0
	g.resultLoadingStep = 0
	g.resultAwards = stageResultAwards(g.rt)
	g.pendingMapTarget = -1
	if g.stageIndex >= worldFirstSecretStage(g.worldIndex) && g.rt.GoalExitSecret {
		target := g.stageIndex
		if next, ok := g.worldMap.exitTarget(g.stageIndex, true); ok {
			target = next
		}
		g.resultNewAwards = g.progress.recordSecretStageResult(g.stageIndex, target, g.rt)
		if g.progress.stageUnlockedForWorld(g.worldIndex, target) {
			g.pendingMapTarget = target
		}
	} else {
		g.resultNewAwards = g.progress.recordStageResult(g.stageIndex, g.rt)
		if target, ok := g.worldMap.exitTarget(g.stageIndex, false); ok && g.progress.stageUnlockedForWorld(g.worldIndex, target) {
			g.pendingMapTarget = target
		}
	}
	if g.progressPath != "" {
		if err := saveOriginalProgress(g.progressPath, g.progress); err != nil {
			g.message = err.Error()
		}
	}
}

func (g *Game) beginSealExit() {
	if g.sealExitActive || g.rt == nil || !g.rt.SealStageComplete {
		return
	}
	g.progress.recordSealStageCompletion(g.stageIndex, g.rt)
	if g.progressPath != "" {
		if err := saveOriginalProgress(g.progressPath, g.progress); err != nil {
			g.message = err.Error()
		}
	}
	g.pendingMapTarget = -1
	g.sealExitActive = true
	g.sealExitTicks = 0
	g.message = fmt.Sprintf("%s seal recovered", worldName(g.worldIndex))
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
	if index < 0 || index >= len(g.pack.Stages) || !worldStageImplemented(g.worldIndex, index) {
		g.message = fmt.Sprintf("invalid %s stage %d", worldName(g.worldIndex), index+1)
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
	g.secretExitActive = false
	g.secretExitTicks = 0
	g.sealExitActive = false
	g.sealExitTicks = 0
	g.pendingMapTarget = -1
	g.resultPhase = resultPhaseLoading
	g.resultPhaseTicks = 0
	g.resultLoadingStep = 0
	g.resultAwards = 0
	g.resultNewAwards = 0
	g.worldEffects = nil
	g.introTicks = 0
	if rt.IsTutorialStage() {
		g.introTicks = stageIntroDuration
	}
	g.heroMoveStart = 0
	g.heroMoveOffset = 0
	g.resetHeroFacing()
	g.entranceSteps = rt.EntranceScrollX
	g.demoCameraLocked = false
	g.demoCameraPhase = 0
	g.resetCamera()
	g.updateCompass()
	if g.sounds != nil && g.sounds.enabled {
		g.sounds.Play(worldMusic(g.worldIndex))
	}
	g.message = fmt.Sprintf("loaded %s stage %02d", worldName(g.worldIndex), index+1)
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
	waterBreathingPotion := progress.WaterBreathingPotion
	if g.startupWorldOverridden && rt.Stage.World == original.WorldBavaria {
		// A forced Bavaria start skips the preceding Angkor journey. Preserve
		// the minimum source campaign state at each directly reachable stage.
		toolLevel = maxToolLevel(toolLevel, 1)
		if stageIndex >= 3 {
			toolLevel = maxToolLevel(toolLevel, 2)
		}
		if stageIndex >= 8 && stageIndex <= bavariaSealStage {
			waterBreathingPotion = true
		}
	}
	if rt.Stage.World == original.WorldAngkor && stageIndex == 4 {
		// Angkor Stage 5 is revisited after the Mystic Hook is obtained in
		// Bavaria. Keep direct-stage and legacy saves in that source-valid state.
		toolLevel = maxToolLevel(toolLevel, 2)
	}
	if rt.Stage.World == original.WorldAngkor && stageIndex == 7 {
		// Angkor Stage 8's secret route is revisited with the Freeze Hammer.
		// The intervening world is outside this Angkor-only campaign slice.
		toolLevel = maxToolLevel(toolLevel, 8)
	}
	if rt.Stage.World == original.WorldAngkor && stageIndex >= angkorFirstSecretStage && stageIndex < angkorTutorialStage {
		// The four Angkor secret stages are revisited after the intervening
		// worlds and shop upgrades. Supply that source-valid campaign state in
		// this Angkor-only slice.
		toolLevel = maxToolLevel(toolLevel, 8)
		rt.MaxHealth = max(rt.MaxHealth, 8)
		rt.Health = rt.MaxHealth
	}
	rt.SpecialItemMask = toolLevelSpecialItemMask(toolLevel)
	if waterBreathingPotion {
		rt.SpecialItemMask |= 4
	}
	rt.ApplyPersistentRewardCoordinates(persistentRewardsForStage(rt, stageIndex, progress))
	rt.SaveSnapshot()
}

func persistentRewardsForStage(rt *original.Runtime, stageIndex int, progress originalProgress) []original.Point {
	if rt == nil || rt.Stage == nil || stageIndex < 0 || stageIndex >= worldStageCount(rt.Stage.World) {
		return nil
	}
	points := append([]original.Point(nil), progress.consumedRewardsForWorld(rt.Stage.World, stageIndex)...)
	redComplete := rt.TotalRedDiamonds > 0 && progress.stageRedDiamondsForWorld(rt.Stage.World, stageIndex) >= rt.TotalRedDiamonds
	relicBit := 1 << rt.Stage.World
	relicComplete := stageIndex == worldSealStage(rt.Stage.World) && progress.RelicMask&relicBit != 0
	if !redComplete && !relicComplete {
		return points
	}
	for idx, id := range rt.Stage.Player {
		if id != 2 && id != 53 {
			continue
		}
		if id == 2 && !redComplete || id == 53 && !relicComplete {
			continue
		}
		foreground := rt.Stage.Foreground[idx]
		if foreground != 14 && foreground != 33 {
			continue
		}
		point := original.Point{X: idx % rt.Width(), Y: idx / rt.Width()}
		if !containsPoint(points, point) {
			points = append(points, point)
		}
	}
	return points
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
	if g.mode == gameModeStartMenu {
		g.drawStartMenu(screen)
		return
	}
	if g.mode == gameModeWorldMap {
		g.drawWorldMap(screen)
		return
	}
	if g.mode == gameModeWorldSelect {
		g.drawWorldSelect(screen)
		return
	}
	if g.sealExitActive {
		g.drawSealLoading(screen)
		return
	}
	if g.secretExitActive {
		g.drawSecretExit(screen)
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
	g.drawStageCells(world, camX, camY)
	g.drawGreatAnaconda(world, camX, camY)
	g.drawEvilTeutonicKnight(world, camX, camY)
	g.drawFallingTorchStage(world, camX, camY)
	g.drawTutorialSealOverlay(world, camX, camY)
	g.drawTutorialRecallHint(world, camX, camY)
	playerX, playerY := g.renderedPlayerPixels()
	renderedPlayerX := playerX - camX
	renderedPlayerY := playerY - camY
	if !g.rt.TutorialSealActivated {
		g.drawPlayer(world, renderedPlayerX, renderedPlayerY)
	}
	g.drawBavariaWater(world, camX, camY)
	g.drawStageForegroundOverlays(world, camX, camY)
	g.drawWorldEffects(world, camX, camY)
	if !g.rt.TutorialSealActivated {
		g.drawChestRewardEffect(world, renderedPlayerX, renderedPlayerY)
		g.drawSpecialBarrierPrompt(world, renderedPlayerX, renderedPlayerY)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, playfieldTop)
	screen.DrawImage(world, op)
	g.drawHUD(screen)
	g.drawGreatAnacondaHealth(screen)
	g.drawEvilTeutonicKnightHealth(screen)
	g.drawTutorialFlash(screen)
	if index, _, ok := g.rt.EnemyGateMessage(); ok {
		drawSourcePanelLabel(screen, g.fontSmall, enemyGateMessageText(index), original.ScreenWidth/2, 223)
	}
	if _, ok := g.rt.TutorialPrompt(); ok {
		g.drawTutorialPrompt(screen)
	} else if g.rt.TutorialScriptActive {
		g.drawTutorialChrome(screen)
	} else if g.introTicks < stageIntroDuration {
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

func enemyGateMessageText(index int) string {
	switch index {
	case 51:
		return "Defeat the Great Anaconda!"
	case 58:
		return "Light the torch!"
	default:
		return "Defeat everyone!"
	}
}

// Java keeps the scrolling tile map in a background buffer, then scans dynamic
// cells from -1 through 11 around the 240x240 viewport. Keeping those passes
// separate prevents a later floor tile from clipping a moving object or flame.
type stageCellViewport struct {
	firstX, firstY     int
	offX, offY         int
	firstRel           int
	lastRelX, lastRelY int
}

func sourceStageCellViewport(camX, camY int) stageCellViewport {
	return stageCellViewport{
		firstX:   camX / original.TileSize,
		firstY:   camY / original.TileSize,
		offX:     -(camX % original.TileSize),
		offY:     -(camY % original.TileSize),
		firstRel: -1,
		lastRelX: original.ScreenWidth/original.TileSize + 2,
		lastRelY: playfieldHeight/original.TileSize + 2,
	}
}

func (g *Game) drawStageCells(dst *ebiten.Image, camX, camY int) {
	view := sourceStageCellViewport(camX, camY)

	for relY := view.firstRel; relY < view.lastRelY; relY++ {
		for relX := view.firstRel; relX < view.lastRelX; relX++ {
			x := view.firstX + relX
			y := view.firstY + relY
			if x < 0 || y < 0 || x >= g.rt.Width() || y >= g.rt.Height() {
				continue
			}
			g.drawCellBackground(dst, x, y, view.offX+relX*original.TileSize, view.offY+relY*original.TileSize)
		}
	}
	for relY := view.firstRel; relY < view.lastRelY; relY++ {
		for relX := view.firstRel; relX < view.lastRelX; relX++ {
			x := view.firstX + relX
			y := view.firstY + relY
			if x < 0 || y < 0 || x >= g.rt.Width() || y >= g.rt.Height() {
				continue
			}
			g.drawCellObjects(dst, x, y, view.offX+relX*original.TileSize, view.offY+relY*original.TileSize)
		}
	}
}

func (g *Game) drawStageForegroundOverlays(dst *ebiten.Image, camX, camY int) {
	view := sourceStageCellViewport(camX, camY)
	for relY := view.firstRel; relY < view.lastRelY; relY++ {
		for relX := view.firstRel; relX < view.lastRelX; relX++ {
			x := view.firstX + relX
			y := view.firstY + relY
			if x < 0 || y < 0 || x >= g.rt.Width() || y >= g.rt.Height() {
				continue
			}
			g.drawCellForegroundOverlay(dst, x, y, view.offX+relX*original.TileSize, view.offY+relY*original.TileSize)
		}
	}
}

func (g *Game) drawSealLoading(screen *ebiten.Image) {
	screen.Fill(color.Black)
	progress := min(230, (g.sealExitTicks+1)*230/sealLoadingSteps)
	drawRect(screen, 5, 310, progress, 6, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	drawRect(screen, 4, 309, 231, 1, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	drawRect(screen, 4, 316, 231, 1, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	drawRect(screen, 4, 310, 1, 6, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	drawRect(screen, 234, 310, 1, 6, color.RGBA{0xfc, 0x9a, 0x04, 0xff})
	g.fontMedium.drawText(screen, "LOADING", original.ScreenWidth/2, 304, true, color.White)
}

func (g *Game) drawSecretExit(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawSourcePanelLines(screen, g.fontSmall, []string{
		"Congratulations! You have",
		"unlocked a secret path!",
	}, original.ScreenWidth/2, 164)
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
	g.fontSmall.drawText(screen, fmt.Sprintf("STAGE %d", g.stageIndex+1), original.ScreenWidth/2+titleOffset, resultStageTitleY, true, textColor)
	g.fontSmall.drawText(screen, "COMPLETE!", original.ScreenWidth/2+completeOffset, resultCompleteY, true, textColor)

	if g.resultPhase >= resultPhaseVioletGems {
		offset := stageResultRowOffset(g.resultPhase, resultPhaseVioletGems, g.resultPhaseTicks)
		g.violetGem.drawFrame(screen, 0, 7+offset, resultVioletLabelY, 0)
		g.fontSmall.drawText(screen, "DIAMONDS", original.ScreenWidth/2, resultVioletLabelY, true, textColor)
		count := g.rt.VioletGems
		if g.resultPhase == resultPhaseVioletGems {
			count = min(g.resultPhaseTicks>>1, count)
		}
		g.fontSmall.drawText(screen, fmt.Sprintf("%d/%d", count, g.rt.TotalVioletGems), original.ScreenWidth/2, resultVioletCountY, true, textColor)
	}
	if g.resultPhase >= resultPhaseRedDiamonds {
		offset := stageResultRowOffset(g.resultPhase, resultPhaseRedDiamonds, g.resultPhaseTicks)
		g.redDiamond.drawFrame(screen, 0, 7+offset, resultRedLabelY, 0)
		g.fontSmall.drawText(screen, "RED DIAMONDS", original.ScreenWidth/2, resultRedLabelY, true, textColor)
		g.fontSmall.drawText(screen, fmt.Sprintf("%d/%d", g.rt.RedDiamonds, g.rt.TotalRedDiamonds), original.ScreenWidth/2, resultRedCountY, true, textColor)
		g.drawStageResultAward(screen, resultAwardVioletGems, resultPhaseRedDiamonds, 69, 86, resultEffectDoubleShort)
	}
	if g.resultPhase >= resultPhaseHits {
		offset := stageResultRowOffset(g.resultPhase, resultPhaseHits, g.resultPhaseTicks)
		g.drawHeroResultIcon(screen, 10, 7+offset, resultHitsIconY)
		g.fontSmall.drawText(screen, "HITS", original.ScreenWidth/2, resultHitsLabelY, true, textColor)
		g.fontSmall.drawText(screen, fmt.Sprintf("%d", g.rt.HitCount), original.ScreenWidth/2, resultHitsCountY, true, textColor)
		g.drawStageResultAward(screen, resultAwardRedDiamonds, resultPhaseHits, 125, 142, resultEffectHalf)
	}
	if g.resultPhase >= resultPhaseRetries {
		offset := stageResultRowOffset(g.resultPhase, resultPhaseRetries, g.resultPhaseTicks)
		g.drawHeroResultIcon(screen, 12, 7+offset, resultRetriesIconY)
		g.fontSmall.drawText(screen, "RETRIES", original.ScreenWidth/2, resultRetriesLabelY, true, textColor)
		g.fontSmall.drawText(screen, fmt.Sprintf("%d", g.rt.Retries), original.ScreenWidth/2, resultRetriesCountY, true, textColor)
		g.drawStageResultAward(screen, resultAwardNoHits, resultPhaseRetries, 181, 198, resultEffectHalf)
	}
	if g.resultPhase >= resultPhaseComplete {
		g.drawStageResultAward(screen, resultAwardNoRetries, resultPhaseComplete, 237, 254, resultEffectDoubleLong)
	}
	prompt := desktopActionKeyLabel + ": SKIP"
	if g.resultPhase == resultPhaseComplete {
		prompt = desktopActionKeyLabel + ": CONTINUE"
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
	if g.rt == nil {
		return g.cameraX, g.cameraY
	}
	shake := g.rt.FallingTorchShake(g.tick)
	if g.rt.Anaconda.Enabled && g.rt.Anaconda.RumbleTicks > 0 {
		rumble := g.rt.Anaconda.RumbleTicks
		shake = max(shake, rumble*g.tick%((rumble>>1)+1)%12)
	}
	if g.rt.TeutonicKnight.Enabled && g.rt.TeutonicKnight.RumbleTicks > 0 {
		rumble := g.rt.TeutonicKnight.RumbleTicks
		shake = max(shake, rumble*g.tick%((rumble>>1)+1)%12)
	}
	return g.cameraX, max(0, g.cameraY-shake)
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
	if targetX, targetY, phase, elapsed, duration, ok := g.rt.TutorialCamera(); ok {
		phase += 200
		if !g.demoCameraLocked || g.demoCameraPhase != phase {
			g.demoCameraLocked = true
			g.demoCameraPhase = phase
			g.demoCameraStartX = g.cameraX
			g.demoCameraStartY = g.cameraY
		}
		g.cameraX = (targetX*elapsed + g.demoCameraStartX*(duration-elapsed)) / duration
		g.cameraY = (targetY*elapsed + g.demoCameraStartY*(duration-elapsed)) / duration
		return
	}
	if targetX, targetY, phase, elapsed, duration, ok := g.rt.EnemyGateDemoCamera(); ok {
		if !g.demoCameraLocked || g.demoCameraPhase != phase {
			g.demoCameraLocked = true
			g.demoCameraPhase = phase
			g.demoCameraStartX = g.cameraX
			g.demoCameraStartY = g.cameraY
		}
		if duration <= 0 {
			duration = 1
		}
		g.cameraX = (targetX*elapsed + g.demoCameraStartX*(duration-elapsed)) / duration
		g.cameraY = (targetY*elapsed + g.demoCameraStartY*(duration-elapsed)) / duration
		return
	}
	if targetX, targetY, elapsed, duration, ok := g.rt.ForegroundDemoCamera(); ok {
		phase := 100 + g.rt.ForegroundDemoPhase
		if !g.demoCameraLocked || g.demoCameraPhase != phase {
			g.demoCameraLocked = true
			g.demoCameraPhase = phase
			g.demoCameraStartX = g.cameraX
			g.demoCameraStartY = g.cameraY
		}
		if duration <= 0 {
			duration = 1
		}
		g.cameraX = (targetX*elapsed + g.demoCameraStartX*(duration-elapsed)) / duration
		g.cameraY = (targetY*elapsed + g.demoCameraStartY*(duration-elapsed)) / duration
		return
	}
	g.demoCameraLocked = false
	g.demoCameraPhase = 0
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

type hudCounterValues struct {
	Violet     int
	Red        int
	GoldKeys   int
	SilverKeys int
	Lives      int
}

func (g *Game) currentHUDCounterValues() hudCounterValues {
	if g == nil || g.rt == nil {
		return hudCounterValues{}
	}
	return hudCounterValues{
		Violet:     g.rt.VioletGems,
		Red:        g.rt.RedDiamonds,
		GoldKeys:   g.rt.KeyForForeground9,
		SilverKeys: g.rt.KeyForForeground8,
		Lives:      g.rt.ExtraLives,
	}
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	counters := g.currentHUDCounterValues()
	g.hud.drawFrame(screen, 20, original.ScreenWidth/2, 0, 0)
	g.hud.drawFrame(screen, 0, original.ScreenWidth/2, original.ScreenHeight, 0)
	g.hud.drawFrame(screen, 1, original.ScreenWidth/2, original.ScreenHeight, 0)
	if g.rt.CompassEnabled {
		g.hud.drawFrame(screen, 2, original.ScreenWidth/2+2, original.ScreenHeight, 0)
		g.hud.drawFrame(screen, 3+g.compassDirection, original.ScreenWidth/2+2, original.ScreenHeight, 0)
	}
	g.drawHealth(screen)
	g.hud.drawNumber(screen, counters.Violet, 190, 308)
	g.hud.drawNumber(screen, counters.Red, 227, 308)
	g.hud.drawNumber(screen, counters.GoldKeys, 167, 18)
	g.hud.drawNumber(screen, counters.SilverKeys, 207, 18)
	g.hud.drawNumber(screen, counters.Lives, 91, 18)
}

func (g *Game) drawStageIntro(screen *ebiten.Image) {
	worldX, stageX := stageTitlePositions(g.introTicks)
	drawSourcePanelLabel(screen, g.fontMedium, strings.ToUpper(worldName(g.worldIndex)), worldX, playfieldTop+15)
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

func (g *Game) drawFallingTorchStage(dst *ebiten.Image, camX, camY int) {
	if g.rt == nil || !g.rt.IsFallingTorchStage() || g.fallingTorches == nil || g.fallingFire == nil {
		return
	}
	torchY := 42*original.TileSize - camY
	g.fallingTorches.drawAnimationWithFrameOffset(dst, g.rt.FallingTorchAnimation, g.rt.FallingTorchAnimationTicks, 10*original.TileSize-camX, torchY, 0)
	g.fallingTorches.drawAnimationWithFrameOffset(dst, g.rt.FallingTorchAnimation, g.rt.FallingTorchAnimationTicks, 14*original.TileSize-camX, torchY, 1)

	if g.rt.FallingTorchWarningTicks > 10 && g.fallingDebris != nil {
		for particle := 3; particle < 13; particle += 2 {
			period := 10 * (particle*2/5 + 1)
			x := (period + g.tick/period) * particle % original.ScreenWidth
			y := (original.ScreenHeight / period * g.tick) % original.ScreenHeight
			g.fallingDebris.drawModule(dst, particle&1, x, y)
		}
	}

	fireWorldY, active := g.rt.RisingFireWorldY()
	if !active {
		return
	}
	fireY := fireWorldY - camY
	left := 168 - camX
	for left <= -original.TileSize {
		left += original.TileSize
	}
	if g.rt.RisingFireFillVisible() {
		for y := fireY + 20; y < playfieldHeight; y += original.TileSize {
			for x := left; x < original.ScreenWidth; x += original.TileSize {
				sequence := ((g.tick >> 1) + x + y) & 1
				g.fallingFire.drawAnimationSequenceFrame(dst, 1, sequence, x, y, 0)
			}
		}
	}
	topX := 168 + original.ScreenWidth/2 - camX
	g.fallingFire.drawAnimationWithFrameOffset(dst, g.rt.RisingFireAnimation, g.rt.RisingFireAnimationTicks, topX, fireY, 0)
	g.fallingFire.drawAnimationWithFrameOffset(dst, g.rt.RisingFireAnimation, g.rt.RisingFireAnimationTicks, topX, fireY, 1)

	flash := (g.tick << 3) % 160
	green := 160 - flash
	if (g.tick/160)&1 != 0 {
		green = flash
	}
	border := color.RGBA{255, uint8(clamp(green, 0, 255)), 0, 255}
	drawRect(dst, 0, 0, original.ScreenWidth, 1, border)
	drawRect(dst, 0, 0, 1, playfieldHeight, border)
	drawRect(dst, original.ScreenWidth-1, 0, 1, playfieldHeight, border)
}

func (g *Game) drawGreatAnaconda(dst *ebiten.Image, camX, camY int) {
	if g.rt == nil || !g.rt.Anaconda.Enabled || g.anaconda == nil || g.anacondaPlatform == nil {
		return
	}
	boss := g.rt.Anaconda
	g.anaconda.drawAnimation(dst, boss.Animation, boss.AnimationTicks, boss.X()*original.TileSize-camX, boss.BodyY-camY, 0)
	if boss.TailVisible && g.flames != nil {
		g.flames.drawAnimation(dst, boss.TailAnimation, boss.TailAnimationTicks, (boss.X()+1)*original.TileSize-camX, 96-camY, 0)
	}
	for column := 0; column < 3; column++ {
		x := 10 + column*(2+boolToInt(column > 0))
		g.anacondaPlatform.drawFrame(dst, 1, x*original.TileSize-camX, 216-camY, 0)
	}
}

func (g *Game) drawGreatAnacondaHealth(screen *ebiten.Image) {
	if g.rt == nil || !g.rt.Anaconda.Enabled || g.rt.Anaconda.Health <= 0 {
		return
	}
	phase := g.rt.Anaconda.Phase
	if phase == original.AnacondaPhaseDormant || phase == original.AnacondaPhaseComplete {
		return
	}
	x := (original.ScreenWidth - 44) / 2
	drawRect(screen, x, 5, 44, 12, color.Black)
	for segment := 0; segment < g.rt.Anaconda.Health; segment++ {
		drawRect(screen, x+2+segment*14, 7, 12, 8, color.RGBA{0x3b, 0xb7, 0x8f, 0xff})
	}
}

func (g *Game) drawCellBackground(dst *ebiten.Image, x, y, px, py int) {
	drawRect(dst, px, py, original.TileSize, original.TileSize, color.RGBA{18, 20, 24, 255})
	playerID, _ := g.rt.At(original.PlayerLayer, x, y)
	if playerID == original.EmptyRawID || playerID < 80 {
		g.floor.drawModule(dst, 0, px, py)
	} else {
		g.drawWorldFrame(dst, int(playerID-80), px, py)
	}
}

func (g *Game) drawCellObjects(dst *ebiten.Image, x, y, px, py int) {
	cellPX, cellPY := px, py
	playerID, _ := g.rt.At(original.PlayerLayer, x, y)
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
			g.checkpoint.drawFrame(dst, frame, px, py, 0)
		case id == 5 || id == 28:
			g.goal.drawModule(dst, sourceGoalFrameForWorld(g.worldIndex, g.stageIndex), px, py)
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
		case id == 15 || id == 16 || id == 27 || id == 34 || id == 35 || id == 37:
			g.drawBavariaForeground(dst, id, px, py)
		case id == 14:
			state := clamp(g.objectStateAt(x, y), 0, 2)
			g.specialContainer.drawAnimationSequenceFrame(dst, 0, state, px, py, 0)
		case id == 33:
			g.hiddenOverlay.drawFrame(dst, hiddenOverlayFrame(g.objectStateAt(x, y)), px, py, 0)
		}
	}
	if playerID < 80 && playerID != original.EmptyRawID {
		idx := x + y*g.rt.Width()
		motion := original.ObjectMotion{}
		if idx >= 0 && idx < len(g.rt.ObjectMotion) {
			motion = g.rt.ObjectMotion[idx]
			if playerID == 0 || playerID == 1 || playerID == 8 || playerID == 9 || playerID == 47 {
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
				g.violetGem.drawFrame(dst, sourceGemFrame(g.tick), px, py, 0)
			case 2:
				g.redDiamond.drawFrame(dst, sourceGemFrame(g.tick), px, py, 0)
			case 4:
				g.drawCenteredModule(dst, g.goldKey, 0, px, py)
			case 5:
				g.drawCenteredModule(dst, g.silverKey, 0, px, py)
			case 6:
				g.commonPickups.drawModule(dst, 0, px, py)
			case 7:
				g.commonPickups.drawModule(dst, 1, px, py)
			case 9:
				switch g.rt.FrozenOriginalAt(x, y) {
				case 1:
					g.frozenViolet.drawFrame(dst, 0, px, py, 0)
				default:
					g.frozenSnake.drawFrame(dst, 0, px, py, 0)
				}
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
			case 24, 31, 33, 41, 50, 79:
				// These are hidden payload/marker objects. Their foreground container,
				// lock, or stage-entry flow owns the visible source sprite.
			case 19:
				g.drawSnake(dst, x, y, px, py)
			case 22:
				g.hazard.drawFrame(dst, 1, px, py, 0)
				g.flames.drawAnimationWithFrameOffset(dst, 0, g.tick, px+original.TileSize, py, 0)
			case 23:
				g.hazard.drawFrame(dst, 0, px, py, 0)
				g.flames.drawAnimationWithFrameOffset(dst, 0, g.tick, px, py, 1)
			case 43:
				g.drawSnake(dst, x, y, px, py)
			case 48:
				drawRect(dst, px+4, py+5, 16, 14, color.RGBA{110, 112, 118, 245})
				drawRect(dst, px+7, py+3, 10, 4, color.RGBA{178, 184, 190, 245})
				drawRect(dst, px+6, py+10, 12, 2, color.RGBA{70, 74, 82, 245})
			default:
				if g.drawBavariaObject(dst, playerID, x, y, px, py) {
					break
				}
				drawRect(dst, px+7, py+7, 10, 10, color.RGBA{80, 140, 230, 180})
			}
		}
	}
	g.drawTutorialSealCell(dst, x, y, cellPX, cellPY)
}

func (g *Game) drawCellForegroundOverlay(dst *ebiten.Image, x, y, px, py int) {
	foregroundID, _ := g.rt.At(original.ForegroundLayer, x, y)
	switch {
	case foregroundID == 32:
		g.drawDiggableFrame(dst, g.rt.ForegroundStateAt(x, y), px, py)
	case foregroundID >= 20 && foregroundID < 26:
		g.foregroundEffects.drawAnimationSequenceFrame(dst, sourceForegroundEffectAnimation(foregroundID), sourceForegroundEffectSequence(g.tick), px, py, 0)
	default:
		if frame, ok := sourceWorldOverlayFrame(foregroundID); ok {
			g.drawWorldFrame(dst, frame, px, py)
		}
	}
}

func sourceForegroundEffectSequence(sourceTick int) int {
	return sourceTick >> 2
}

func sourceForegroundEffectAnimation(id original.RawID) int {
	animation := int(id - 20)
	if animation >= 0 && animation < 4 {
		animation ^= 2
	}
	return animation
}

func sourceWorldOverlayFrame(id original.RawID) (int, bool) {
	if id == original.EmptyRawID || id < 80 {
		return 0, false
	}
	return int(id - 80), true
}

func sourceCellObjectVisible(playerID, foregroundID original.RawID) bool {
	if foregroundID != 14 && foregroundID != 33 {
		return true
	}
	switch playerID {
	case 2, 4, 5, 6, 7, 24, 26, 27, 40, 41, 42, 51, 52, 53:
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
	clip := image.Rect(px, py, px+original.TileSize, py+original.TileSize).Intersect(dst.Bounds())
	if clip.Empty() {
		return
	}
	clipped := dst.SubImage(clip).(*ebiten.Image)
	g.pressureSwitch.drawModule(clipped, 0, px, py+original.TileSize-height+offset)
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

func (g *Game) drawSpecialBarrierPrompt(dst *ebiten.Image, playerX, playerY int) {
	toolModule, available, visible := g.rt.SpecialBarrierPrompt()
	if !visible || g.toolPrompt == nil || g.tools == nil {
		return
	}
	g.drawCenteredAt(dst, g.toolPrompt, 0, playerX-12, playerY-22)
	if available {
		g.drawCenteredAt(dst, g.tools, toolModule, playerX-12, playerY-24)
		drawControlKeycap(dst, g.fontSmall, desktopActionKeyLabel, playerX+original.TileSize/2, playerY-54)
		return
	}
	g.drawCenteredAt(dst, g.toolPrompt, 1, playerX-12, playerY-24)
}

func (g *Game) drawCenteredAt(dst *ebiten.Image, sheet *spriteSheet, module, centerX, centerY int) {
	if sheet == nil || module < 0 || module >= len(sheet.meta.Modules) {
		return
	}
	bounds := sheet.meta.Modules[module]
	sheet.drawModule(dst, module, centerX-bounds.W/2, centerY-bounds.H/2)
}

func sourceGemFrame(sourceTick int) int {
	frame := (sourceTick & 0x3f) >> 1
	if frame >= 4 {
		return 0
	}
	return frame
}

func sourceGoalFrame(stageIndex int) int {
	return sourceGoalFrameForWorld(original.WorldAngkor, stageIndex)
}

func sourceGoalFrameForWorld(world, stageIndex int) int {
	if world == original.WorldAngkor && (stageIndex == 4 || stageIndex == 7) || world == original.WorldBavaria && stageIndex == 8 {
		return 1
	}
	return 0
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
	g.worldFrames.drawFrame(dst, frame, px, py, 0)
}

func (g *Game) drawBoulderFrame(dst *ebiten.Image, frame, px, py int) {
	g.boulder.drawFrame(dst, frame, px, py, 0)
}

func (g *Game) drawDiggableFrame(dst *ebiten.Image, frame, px, py int) {
	if g.diggable == nil || len(g.diggable.meta.FrameCounts) == 0 {
		return
	}
	g.diggable.drawFrame(dst, clamp(frame, 0, len(g.diggable.meta.FrameCounts)-1), px, py, 0)
}

func (g *Game) drawPlayer(dst *ebiten.Image, px, py int) {
	if g.rt.InvulnerabilityTicks > 0 && (g.tick>>1)&1 != 0 {
		return
	}
	if g.worldIndex == original.WorldBavaria && g.lastDY == 0 && g.rt.Player.Y+1 < g.rt.Height() && g.rt.WaterAt(g.rt.Player.X, g.rt.Player.Y+1) > 0 {
		below, _ := g.rt.At(original.PlayerLayer, g.rt.Player.X, g.rt.Player.Y+1)
		if below == 0 || below == 1 || below == 8 || below == 9 {
			phase := (g.tick >> 1) + g.rt.Player.X
			bob := phase % 4
			if phase/4&1 != 0 {
				bob = 4 - bob
			}
			py += bob
		}
	}
	animation, animationTick := g.heroAnimationState()
	g.hero.drawAnimationWithFrameOffset(dst, animation, animationTick, px, py, 0)
	if (g.rt.ChestOpening && chestRewardIconVisible(g.hero, g.rt.ChestAnimation, g.rt.ChestTicks)) || g.rt.RelicCelebrating {
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
	} else if g.rt.RelicCelebrating {
		animation = 47
		animationTick = g.rt.RelicCelebrationTicks
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
	} else if g.worldIndex == original.WorldBavaria && g.rt.UsesSwimmingAnimationAt(g.rt.Player.X, g.rt.Player.Y) {
		animation = 36 + direction
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

func (g *Game) drawChestRewardIcon(dst *ebiten.Image, px, py int) bool {
	switch g.rt.ChestRewardID {
	case 2:
		g.redDiamond.drawFrame(dst, 0, px, py, 0)
	case 4:
		g.drawCenteredModule(dst, g.goldKey, 0, px, py)
	case 5:
		g.drawCenteredModule(dst, g.silverKey, 0, px, py)
	case 6:
		g.commonPickups.drawModule(dst, 0, px, py)
	case 7:
		g.commonPickups.drawModule(dst, 1, px, py)
	case 24:
		g.tools.drawModule(dst, 0, px, py)
	case 27:
		g.tools.drawModule(dst, 1, px, py)
	case 26:
		g.tools.drawModule(dst, 2, px, py)
	case 41:
		g.violetGem.drawFrame(dst, 0, px, py, 0)
		g.hud.drawNumber(dst, g.rt.ChestRewardValue, px+original.TileSize, py+14)
	case 42:
		g.compassPickup.drawModule(dst, 0, px, py)
	case 40:
		g.drawCenteredModule(dst, g.bavaria.waterPotion, 0, px, py)
	case 51:
		g.drawCenteredModule(dst, g.bavariaSeal, 0, px, py)
	case 52:
		g.drawCenteredModule(dst, g.siberiaSeal, 0, px, py)
	case 53:
		g.drawCenteredModule(dst, g.angkorSeal, 0, px, py)
	default:
		return false
	}
	return true
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

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
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
		justPressed(ebiten.KeyEnter) ||
		justPressed(ebiten.KeyR) ||
		(justPressed(ebiten.KeyDigit8) && shiftPressed(pressed))
}

func shiftPressed(pressed func(ebiten.Key) bool) bool {
	return pressed(ebiten.KeyShift) || pressed(ebiten.KeyShiftLeft) || pressed(ebiten.KeyShiftRight)
}

func centerActionPressed() bool {
	return centerActionPressedWith(inpututil.IsKeyJustPressed)
}

func centerActionPressedWith(justPressed func(ebiten.Key) bool) bool {
	return justPressed(ebiten.KeyDigit5) ||
		justPressed(ebiten.KeyNumpad5) ||
		justPressed(ebiten.KeySpace)
}

func tutorialSkipPressed() bool {
	return tutorialSkipPressedWith(inpututil.IsKeyJustPressed)
}

func tutorialSkipPressedWith(justPressed func(ebiten.Key) bool) bool {
	return justPressed(ebiten.KeyS)
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
