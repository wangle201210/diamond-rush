package original

import (
	"fmt"
)

const (
	hurtStateDuration           = 8
	hurtInvulnerabilityDuration = 48
	deathTransitionDuration     = 80
	deathDuration               = hurtStateDuration + deathTransitionDuration
	rockHoldDuration            = 40
	boulderPushAttempts         = 7
	gravityRollFastOffset       = 6
	gravityRollMoveOffset       = 12
	chestOpenDuration           = 67
	chestRewardSoundTick        = 37
	chestRewardTick             = 39
	chestShortOpenDuration      = 46
	chestShortRewardTick        = 23
	pickupComboWindowTicks      = 100
	digAnimationFrames          = 8
	playerMoveStartOffset       = 18
	playerMoveStep              = 6
	recallAnimationDuration     = 42
	lockOpenDuration            = 17
	lockRewardTick              = 10
	hammerImpactTick            = 3
	hammerUpDuration            = 11
	hammerOtherDuration         = 12
	hookRightCastAnimation      = 20
	hookRightPullAnimation      = 21
	hookLeftCastAnimation       = 22
	hookLeftPullAnimation       = 23
	hookRightCastDuration       = 7
	hookRightPullDuration       = 5
	hookLeftCastDuration        = 6
	hookLeftPullDuration        = 4
	objectDirectionMask         = 0x7
	boulderRotationMask         = 0x38
	gravityRollPreparing        = 0x200
	gravityMoveRight            = 0x400
	gravityMoveLeft             = 0x800
	explosiveFallShift          = 17
	explosiveFallMask           = 0x1f << explosiveFallShift
	explosionImpactTick         = 6
	explosionDuration           = 12
	movingHazardCooldownShift   = 8
	movingHazardCooldownMask    = 0xff << movingHazardCooldownShift
	spearTimerShift             = 8
	spearTimerMask              = 0xff << spearTimerShift
	snakeStunMask               = 0xf8
	snakeStunDuration           = 0x78
	fallingTorchStageIndex      = 5
	fallingTorchWarningDuration = 120
	fallingTorchWarningLoop     = 60
	fallingTorchCollapseTicks   = 37
	fallingFireStartTicks       = 29
	fallingFireInitialHeight    = 816
	fallingFireMaximumHeight    = 1704
	foregroundDemoPanTicks      = 60
	foregroundDemoWaitTicks     = 20
	greatAnacondaStageIndex     = 8
	angkorSealCompletionTicks   = 140
	teutonicKnightStageIndex    = 9
	sealCompletionTicks         = 140
)

const (
	SoundSwitch       = 0
	SoundRiddle       = 1
	SoundDeath        = 2
	SoundChestOpen    = 3
	SoundChestReward  = 4
	SoundHeroHurt     = 5
	SoundHammerBlock  = 6
	SoundBossDeath    = 7
	SoundCheckpoint   = 9
	SoundEnemyHit     = 10
	SoundBreak        = 11
	SoundHook         = 12
	SoundWater        = 13
	SoundBoulder      = 14
	SoundDoor         = 8
	SoundStageClear   = 15
	SoundAngkorMusic  = 16
	SoundBavariaMusic = 17
	SoundTitleMusic   = 19
)

type ObjectMotion struct {
	DX         int
	DY         int
	Remaining  int
	RollDX     int
	RollOffset int
}

type SourceFrameResult struct {
	GravityMoved           int
	SnakesMoved            int
	CrawlersMoved          int
	HazardHits             int
	RisingFireHits         int
	RockHoldHits           int
	DigCleared             int
	AnacondaHits           int
	AnacondaDefeated       bool
	TeutonicKnightHits     int
	TeutonicKnightDefeated bool
	TutorialSealActivated  bool
	VioletPickups          []Point
}

type Runtime struct {
	Stage                       *Stage
	Player                      Point
	PlayerMotion                ObjectMotion
	playerFacingDirection       int
	playerTurnOffset            int
	EntranceMarker              Point
	EntranceScrollX             int
	EntranceDoor                Point
	EntranceDoorSet             bool
	PlayerLayer                 []RawID
	Background                  []RawID
	Foreground                  []RawID
	ForegroundState             []int
	ObjectState                 []int
	ObjectMotion                []ObjectMotion
	FrozenOriginal              []RawID
	ContainerLocked             []bool
	ConsumedRewardCells         []bool
	WaterDepth                  []uint8
	WaterSources                []Point
	WaterInitializing           bool
	WaterTicks                  int
	water                       waterRuntimeState
	Pushing                     bool
	PushDX                      int
	PushTicks                   int
	pushTarget                  Point
	EnemyGateGroup              []int
	EnemyGateCounters           map[int]int
	EnemyGateMessages           map[int]int
	ActiveEnemyGateGroup        int
	EnemyGateDemoActive         bool
	EnemyGateDemoPhase          int
	EnemyGateDemoTicks          int
	EnemyGateDemoOutboundTicks  int
	EnemyGateDemoTarget         Point
	EnemyGateDemoTargetSet      bool
	EnemyGateMessageIndex       int
	EnemyGateMessageTicks       int
	Anaconda                    GreatAnaconda
	TeutonicKnight              EvilTeutonicKnight
	Checkpoints                 []Point
	GoalMarkers                 []Point
	Doors                       []Point
	DoorGroup                   []int
	TotalVioletGems             int
	TotalRedDiamonds            int
	VioletGems                  int
	RedDiamonds                 int
	KeyForForeground9           int
	KeyForForeground8           int
	ExtraLives                  int
	HealthRefills               int
	BonusValue                  int
	BonusPickups                int
	SpecialItemMask             int
	Hammering                   bool
	HammerTicks                 int
	HammerAnimation             int
	HammerTarget                Point
	Hooking                     bool
	HookTicks                   int
	HookAnimation               int
	HookTarget                  Point
	hookDX                      int
	hookStepsRemaining          int
	hookCollect                 bool
	hookReturning               bool
	hookTip                     Point
	hookOriginalState           int
	SpecialPickup42             bool
	CompassEnabled              bool
	RelicMask                   int
	SpecialPickups              int
	RelicCelebrating            bool
	RelicCelebrationTicks       int
	SealCollected               bool
	SealTicks                   int
	SealStageComplete           bool
	LastForegroundEvent         int
	ForegroundEvents            int
	ForegroundDemoActive        bool
	ForegroundDemoID            int
	ForegroundDemoPhase         int
	ForegroundDemoTicks         int
	foregroundDemoMoved         bool
	pendingForegroundEvent      Point
	pendingForegroundEventSet   bool
	TutorialScriptActive        bool
	TutorialScriptID            int
	TutorialTextIndex           int
	TutorialTextPlacement       int
	TutorialTextY               int
	TutorialTextSide            int
	TutorialPromptX             int
	TutorialComplete            bool
	TutorialRecallHintVisible   bool
	TutorialSealActivated       bool
	TutorialCameraActive        bool
	TutorialCameraTarget        Point
	TutorialCameraTicks         int
	TutorialCameraDuration      int
	TutorialCameraPhase         int
	TutorialPortraitVisible     bool
	TutorialPortraitX           int
	TutorialPortraitY           int
	TutorialPortraitFace        int
	TutorialPortraitMark        int
	TutorialPortraitRevealTicks int
	TutorialFlashVisible        bool
	tutorialCommandIndex        int
	tutorialCommandTicks        int
	tutorialCommandStarted      bool
	tutorialCommandMoveDone     bool
	tutorialMoveStarted         bool
	tutorialMoveAttempts        int
	tutorialPromptAcknowledged  bool
	tutorialSkipping            bool
	tutorialQueuedScript        int
	tutorialResetFirst          bool
	tutorialResetSecond         bool
	tutorialRecallHintTriggered bool
	FallingTorchTriggers        int
	FallingTorchWarningTicks    int
	FallingTorchAnimation       int
	FallingTorchAnimationTicks  int
	SpikeSlowExtent             int
	SpikeFastExtent             int
	FanPhase                    int
	FanDirection                int
	RisingFireHeight            int
	RisingFireAnimation         int
	RisingFireAnimationTicks    int
	viewportX                   int
	viewportY                   int
	viewportSet                 bool
	BonusTarget                 Point
	BonusTargetSet              bool
	BonusRemaining              int
	BonusGateOpen               bool
	LocksOpened                 int
	BreakableWalls              int
	MaxHealth                   int
	Health                      int
	DamageTaken                 int
	HitCount                    int
	Retries                     int
	HurtTicks                   int
	InvulnerabilityTicks        int
	RockHoldTicks               int
	PlayerDead                  bool
	DeathTicks                  int
	RecallUsed                  bool
	RecallPending               bool
	RecallTicks                 int
	ExitOpen                    bool
	ReachedGoal                 bool
	GoalExitSecret              bool
	GoalExitDirection           int
	GoalExitComplete            bool
	CheckpointProgress          int
	CheckpointPending           bool
	pendingCheckpoint           Point
	pendingChest                Point
	pendingChestSet             bool
	ChestOpening                bool
	ChestTicks                  int
	ChestRewarded               bool
	ChestAnimation              int
	ChestRewardID               RawID
	ChestRewardValue            int
	chestOpeningFresh           bool
	LockOpening                 bool
	LockTicks                   int
	LockAnimation               int
	LockPoint                   Point
	LockForegroundID            RawID
	LockRewarded                bool
	lastPickupTick              int
	lastPickupTickSet           bool
	gravitySourceTick           int
	frameVioletPickups          []Point
	soundEvents                 []int
	checkpoint                  Snapshot
}

type Snapshot struct {
	Valid                bool
	Player               Point
	PlayerMotion         ObjectMotion
	PlayerLayer          []RawID
	Background           []RawID
	Foreground           []RawID
	ForegroundState      []int
	ObjectState          []int
	ObjectMotion         []ObjectMotion
	FrozenOriginal       []RawID
	ContainerLocked      []bool
	ConsumedRewardCells  []bool
	WaterDepth           []uint8
	WaterInitializing    bool
	WaterTicks           int
	Water                waterRuntimeState
	EnemyGateGroup       []int
	EnemyGateCounters    map[int]int
	ActiveEnemyGateGroup int
	Anaconda             GreatAnaconda
	TeutonicKnight       EvilTeutonicKnight
	VioletGems           int
	RedDiamonds          int
	KeyForForeground9    int
	KeyForForeground8    int
	HealthRefills        int
	BonusValue           int
	BonusPickups         int
	SpecialItemMask      int
	SpecialPickup42      bool
	CompassEnabled       bool
	RelicMask            int
	SpecialPickups       int
	SealCollected        bool
	SealTicks            int
	SealStageComplete    bool
	LastForegroundEvent  int
	ForegroundEvents     int
	FallingTorchTriggers int
	RisingFireHeight     int
	FanPhase             int
	FanDirection         int
	BonusTarget          Point
	BonusTargetSet       bool
	BonusRemaining       int
	BonusGateOpen        bool
	LocksOpened          int
	BreakableWalls       int
	ExitOpen             bool
	ReachedGoal          bool
	GoalExitSecret       bool
	GoalExitDirection    int
	GoalExitComplete     bool
	CheckpointProgress   int
}

func NewRuntime(stage *Stage) (*Runtime, error) {
	if err := stage.Validate(); err != nil {
		return nil, err
	}
	entrances := stage.EntranceMarkers()
	if len(entrances) != 1 {
		return nil, fmt.Errorf("stage %02d has %d entrance markers, want 1", stage.Index, len(entrances))
	}
	entrance := entrances[0]
	rt := &Runtime{
		Stage:                stage,
		Player:               Point{X: 0, Y: entrance.Y},
		EntranceMarker:       entrance,
		EntranceScrollX:      entrance.X,
		PlayerLayer:          append([]RawID(nil), stage.Player...),
		Background:           append([]RawID(nil), stage.Background...),
		Foreground:           append([]RawID(nil), stage.Foreground...),
		ForegroundState:      make([]int, stage.Width*stage.Height),
		ObjectState:          make([]int, stage.Width*stage.Height),
		ObjectMotion:         make([]ObjectMotion, stage.Width*stage.Height),
		FrozenOriginal:       make([]RawID, stage.Width*stage.Height),
		ContainerLocked:      make([]bool, stage.Width*stage.Height),
		ConsumedRewardCells:  make([]bool, stage.Width*stage.Height),
		WaterDepth:           make([]uint8, stage.Width*stage.Height),
		DoorGroup:            make([]int, stage.Width*stage.Height),
		ActiveEnemyGateGroup: -1,
		Checkpoints:          stage.Positions(ForegroundLayer, 4),
		GoalMarkers:          append(stage.Positions(ForegroundLayer, 5), stage.Positions(ForegroundLayer, 28)...),
		Doors:                stage.Positions(ForegroundLayer, 7),
		TotalVioletGems:      stageVioletTotal(stage),
		TotalRedDiamonds:     countRaw(stage.Player, 2),
		ExtraLives:           5,
		MaxHealth:            4,
		Health:               4,
		CompassEnabled:       !(stage.World == WorldAngkor && stage.Index == tutorialStageIndex),
		RockHoldTicks:        rockHoldDuration,
		ChestRewardID:        EmptyRawID,
	}
	rt.playerFacingDirection = 2
	for idx := range rt.FrozenOriginal {
		rt.FrozenOriginal[idx] = EmptyRawID
		rt.DoorGroup[idx] = -1
	}
	if stage.World == WorldAngkor && stage.Index == fallingTorchStageIndex {
		rt.RisingFireHeight = fallingFireInitialHeight
		rt.RisingFireAnimation = 2
	}
	if stage.World == WorldAngkor && stage.Index == greatAnacondaStageIndex {
		rt.Anaconda = newGreatAnaconda()
	}
	if stage.World == WorldBavaria && stage.Index == teutonicKnightStageIndex {
		rt.TeutonicKnight = newEvilTeutonicKnight()
	}
	rt.initBonusTarget()
	rt.initBavariaObjects()
	rt.initObjectState()
	rt.initWater()
	rt.initDoorStates()
	rt.initEnemyGates()
	rt.initEnemyGateDoors()
	rt.initTutorial()
	rt.set(rt.PlayerLayer, entrance.X, entrance.Y, EmptyRawID)
	rt.initEntranceDoor()
	rt.TickForegroundTriggers()
	rt.updatePressureDoors()
	rt.updateExitOpen()
	rt.SaveSnapshot()
	return rt, nil
}

func (rt *Runtime) initEntranceDoor() {
	point := Point{X: rt.EntranceMarker.X - 2, Y: rt.EntranceMarker.Y}
	if point.X < 0 || point.Y < 0 || point.X >= rt.Width() || point.Y >= rt.Height() {
		return
	}
	idx := rt.index(point.X, point.Y)
	rt.EntranceDoor = point
	rt.EntranceDoorSet = true
	rt.Foreground[idx] = 7
	rt.DoorGroup[idx] = -1
	// Java stores (-193 << 8) | 7. Its merged high byte is 0x3f:
	// open animation phase 3 and low-nibble temporary door id 15.
	rt.Background[idx] = 0x3f
}

// initDoorStates mirrors the source's post-load Hashtable pass. The authored
// background remains the immutable group identity, while the door's runtime
// low nibble becomes the number of switches and keyed locks still to fire.
func (rt *Runtime) initDoorStates() {
	activators := map[int]int{}
	for idx, foregroundID := range rt.Foreground {
		if foregroundID != 6 && foregroundID != 8 && foregroundID != 9 {
			continue
		}
		group := int(rt.Background[idx])
		if group != int(EmptyRawID) {
			activators[group]++
		}
	}
	for idx, foregroundID := range rt.Foreground {
		if foregroundID != 7 {
			continue
		}
		group := int(rt.Background[idx])
		if group == int(EmptyRawID) {
			continue
		}
		rt.DoorGroup[idx] = group
		if count := activators[group]; count > 0 {
			rt.Background[idx] = RawID(count)
		}
	}
}

// CloseEntranceDoor mirrors doorHeadClose(playerX-1, playerY) on the final
// raw-79 auto-entry step.
func (rt *Runtime) CloseEntranceDoor() bool {
	if !rt.EntranceDoorSet {
		return false
	}
	idx := rt.index(rt.EntranceDoor.X, rt.EntranceDoor.Y)
	if rt.Foreground[idx] != 7 || rt.Background[idx]&0xf0 == 0 || rt.isPlayerAt(rt.EntranceDoor.X, rt.EntranceDoor.Y) {
		return false
	}
	rt.Background[idx] &= 0x0f
	rt.emitSound(SoundBoulder)
	return true
}

func (rt *Runtime) DrainSoundEvents() []int {
	if len(rt.soundEvents) == 0 {
		return nil
	}
	events := append([]int(nil), rt.soundEvents...)
	rt.soundEvents = rt.soundEvents[:0]
	return events
}

func (rt *Runtime) emitSound(id int) {
	rt.soundEvents = append(rt.soundEvents, id)
}

func (rt *Runtime) initBonusTarget() {
	for idx, id := range rt.PlayerLayer {
		if id != 12 {
			continue
		}
		x := idx % rt.Width()
		y := idx / rt.Width()
		value := int(rt.Background[idx])
		if value == int(EmptyRawID) || value < 0 {
			value = 0
		}
		rt.BonusTarget = Point{X: x, Y: y}
		rt.BonusTargetSet = true
		rt.BonusRemaining = value
		rt.BonusGateOpen = value <= 0
		if rt.BonusGateOpen {
			rt.PlayerLayer[idx] = EmptyRawID
		}
	}
}

func (rt *Runtime) initBavariaObjects() {
	if rt.Stage.World != WorldBavaria {
		return
	}
	for idx, id := range rt.Stage.Player {
		x, y := idx%rt.Width(), idx/rt.Width()
		switch id {
		case 16:
			if y > 0 && y+1 < rt.Height() && rt.Stage.Player[rt.index(x, y+1)] != 16 {
				above := rt.index(x, y-1)
				rt.PlayerLayer[above] = 16
				rt.Background[above] = rt.Background[idx]
			}
		case 34:
			rt.PlayerLayer[idx] = EmptyRawID
			rt.Foreground[idx] = 15
		case 35:
			rt.Foreground[idx] = EmptyRawID
		case 38:
			rt.Foreground[idx] = 27
		}
	}
}

func (rt *Runtime) initObjectState() {
	for idx, id := range rt.PlayerLayer {
		switch {
		case id == 43:
			dir := int(rt.Background[idx]) & objectDirectionMask
			rt.ObjectState[idx] = dir | 0x10000
		case id == 19:
			dir := int(rt.Background[idx]) & 0x7
			if dir >= 1 && dir <= 4 {
				rt.ObjectState[idx] = dir
			}
		case id == 11:
			if rt.Background[idx] == 1 {
				rt.ObjectState[idx] = 16
			}
		case id == 28:
			state := int(rt.Background[idx])
			if state > 10 {
				state = state/11 | 0x8
			}
			rt.ObjectState[idx] = state
		case id == 14:
			if rt.Background[idx] == 4 {
				rt.ObjectState[idx] = 0x8
			}
		case id == 16:
			direction := int(rt.Background[idx]) & objectDirectionMask
			if direction != 4 {
				direction = 2
			}
			rt.ObjectState[idx] = direction
		case id == 36:
			if rt.Background[idx] == 1 {
				rt.ObjectState[idx] = 1
			}
		}
	}
}

func (rt *Runtime) initEnemyGates() {
	rt.EnemyGateGroup = make([]int, len(rt.PlayerLayer))
	for i := range rt.EnemyGateGroup {
		rt.EnemyGateGroup[i] = -1
	}
	rt.EnemyGateCounters = map[int]int{}
	rt.EnemyGateMessages = map[int]int{}
	for y := 1; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 17 {
				continue
			}
			aboveIdx := rt.index(x, y-1)
			if isPickupContainer(rt.Foreground[aboveIdx]) {
				rt.ContainerLocked[aboveIdx] = true
			}
			group := int(rt.Background[idx])
			if group == int(EmptyRawID) {
				continue
			}
			if isEnemyGateTarget(rt.PlayerLayer[aboveIdx]) {
				rt.EnemyGateCounters[group]++
				rt.EnemyGateGroup[aboveIdx] = group
				message := 56
				if rt.PlayerLayer[aboveIdx] == 36 {
					message = 58
				}
				if rt.Stage.World == WorldAngkor && rt.Stage.Index == greatAnacondaStageIndex {
					message = 51
				}
				rt.EnemyGateMessages[group] = message
				rt.Foreground[idx] = EmptyRawID
			}
		}
	}
}

func (rt *Runtime) Width() int {
	return rt.Stage.Width
}

func (rt *Runtime) Height() int {
	return rt.Stage.Height
}

func (rt *Runtime) At(layer Layer, x, y int) (RawID, bool) {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return 0, false
	}
	idx := x + y*rt.Width()
	switch layer {
	case PlayerLayer:
		return rt.PlayerLayer[idx], true
	case BackgroundLayer:
		return rt.Background[idx], true
	case ForegroundLayer:
		return rt.Foreground[idx], true
	default:
		return 0, false
	}
}

func (rt *Runtime) SetForTest(layer Layer, x, y int, id RawID) bool {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return false
	}
	switch layer {
	case PlayerLayer:
		rt.set(rt.PlayerLayer, x, y, id)
		rt.ObjectMotion[rt.index(x, y)] = ObjectMotion{}
		if id != 9 {
			rt.FrozenOriginal[rt.index(x, y)] = EmptyRawID
		}
	case BackgroundLayer:
		rt.set(rt.Background, x, y, id)
		idx := rt.index(x, y)
		if rt.Foreground[idx] == 7 && id != EmptyRawID && int(id)&0xf0 == 0 {
			rt.DoorGroup[idx] = int(id)
		}
	case ForegroundLayer:
		rt.set(rt.Foreground, x, y, id)
		idx := rt.index(x, y)
		if id != 7 {
			rt.DoorGroup[idx] = -1
		} else if rt.Background[idx] != EmptyRawID && int(rt.Background[idx])&0xf0 == 0 {
			rt.DoorGroup[idx] = int(rt.Background[idx])
		}
	default:
		return false
	}
	return true
}

func (rt *Runtime) IsCheckpoint(x, y int) bool {
	id, ok := rt.At(ForegroundLayer, x, y)
	return ok && id == 4
}

func (rt *Runtime) TryMove(dx, dy int) bool {
	return rt.tryMove(dx, dy, false)
}

func (rt *Runtime) tryMove(dx, dy int, scripted bool) bool {
	if (!scripted && !rt.CanAcceptInput()) || (scripted && !rt.canStartPlayerMove()) {
		rt.ResetPushAttempt()
		return false
	}
	if dx == 0 && dy == 0 {
		rt.ResetPushAttempt()
		return false
	}
	rt.SetPlayerFacing(dx, dy)
	x := rt.Player.X + dx
	y := rt.Player.Y + dy
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		rt.ResetPushAttempt()
		return false
	}
	if dy == 0 && dx != 0 {
		if playerID, _ := rt.At(PlayerLayer, x, y); playerID == 0 || playerID == 8 || playerID == 9 {
			if !rt.pushAttemptReady(x, y, dx) {
				return false
			}
			if !rt.TryPushBoulder(x, y, dx) {
				return false
			}
			rt.ResetPushAttempt()
		} else {
			rt.ResetPushAttempt()
		}
	} else {
		rt.ResetPushAttempt()
	}
	if !rt.IsPassable(x, y) {
		return false
	}
	playerID, _ := rt.At(PlayerLayer, x, y)
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	if isPickupContainer(foregroundID) && isContainerReward(playerID) {
		if !rt.ContainerLocked[rt.index(x, y)] {
			rt.pendingChest = Point{X: x, Y: y}
			rt.pendingChestSet = true
		}
		return rt.finishMove(x, y, dx, dy)
	}
	switch playerID {
	case 10:
		idx := rt.index(x, y)
		rt.ObjectState[idx] = 1
	case 2:
		rt.RedDiamonds++
		rt.set(rt.PlayerLayer, x, y, EmptyRawID)
	case 4:
		rt.KeyForForeground9++
		rt.set(rt.PlayerLayer, x, y, EmptyRawID)
	case 5:
		rt.KeyForForeground8++
		rt.set(rt.PlayerLayer, x, y, EmptyRawID)
	case 6:
		rt.ExtraLives++
		rt.set(rt.PlayerLayer, x, y, EmptyRawID)
	case 7:
		if rt.Health >= rt.MaxHealth {
			rt.TotalVioletGems += 10
			rt.collectBonusAt(x, y, 10)
			return rt.finishMove(x, y, dx, dy)
		}
		rt.HealthRefills++
		rt.HealFull()
		rt.set(rt.PlayerLayer, x, y, EmptyRawID)
	case 24:
		rt.collectSpecialItemAt(x, y, 1)
	case 26:
		rt.collectSpecialItemAt(x, y, 8)
	case 27:
		rt.collectSpecialItemAt(x, y, 2)
	case 41:
		rt.collectBonusAt(x, y, int(rt.Background[rt.index(x, y)]))
	case 42:
		rt.SpecialPickup42 = true
		rt.CompassEnabled = true
		rt.SpecialPickups++
		rt.set(rt.PlayerLayer, x, y, EmptyRawID)
		rt.queueTutorialScript(11)
	case 53:
		rt.RelicMask |= 1
		rt.SpecialPickups++
		rt.set(rt.PlayerLayer, x, y, EmptyRawID)
	}
	return rt.finishMove(x, y, dx, dy)
}

// SettlePlayerMove applies interactions that wait for the hero's movement
// interpolation to complete.
func (rt *Runtime) SettlePlayerMove() bool {
	if !rt.pendingChestSet || rt.PlayerMotion.Remaining > 0 {
		return false
	}
	if rt.Player != rt.pendingChest {
		rt.pendingChestSet = false
		return false
	}
	point := rt.pendingChest
	rt.pendingChestSet = false
	playerID, playerOK := rt.At(PlayerLayer, point.X, point.Y)
	foregroundID, foregroundOK := rt.At(ForegroundLayer, point.X, point.Y)
	if !playerOK || !foregroundOK || !isContainerReward(playerID) || !isPickupContainer(foregroundID) || rt.ContainerLocked[rt.index(point.X, point.Y)] {
		return false
	}
	rt.startChestOpening(point, true)
	return true
}

func (rt *Runtime) startChestOpening(point Point, fresh bool) {
	idx := rt.index(point.X, point.Y)
	rewardID := rt.PlayerLayer[idx]
	if !isContainerReward(rewardID) {
		return
	}
	rewardValue := int(rt.Background[idx])
	if rewardValue == int(EmptyRawID) {
		rewardValue = 0
	}
	if rewardID == 6 && rt.ExtraLives >= 99 {
		rewardID = 7
	}
	if rewardID == 7 && rt.Health >= rt.MaxHealth {
		rewardID = 41
		rewardValue = 10
		rt.TotalVioletGems += 10
	}
	rt.PlayerLayer[idx] = EmptyRawID
	rt.ObjectState[idx] = 1
	rt.ChestOpening = true
	rt.ChestTicks = 0
	rt.ChestRewarded = false
	rt.ChestAnimation = 40
	if rt.lastPickupTickSet && rt.gravitySourceTick-rt.lastPickupTick >= 0 && rt.gravitySourceTick-rt.lastPickupTick < pickupComboWindowTicks {
		rt.ChestAnimation = 48
	}
	rt.ChestRewardID = rewardID
	rt.ChestRewardValue = rewardValue
	rt.chestOpeningFresh = fresh
	rt.emitSound(SoundChestOpen)
}

func (rt *Runtime) finishMove(x, y, dx, dy int) bool {
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	rt.Player = Point{X: x, Y: y}
	rt.PlayerMotion = ObjectMotion{DX: dx, DY: dy, Remaining: playerMoveStartOffset}
	rt.playerTurnOffset = 0
	rt.ResetPushAttempt()
	rt.RockHoldTicks = rockHoldDuration
	switch foregroundID {
	case 0, 26:
		rt.pendingForegroundEvent = Point{X: x, Y: y}
		rt.pendingForegroundEventSet = true
	case 1:
		if rt.Stage.World == WorldAngkor && rt.Stage.Index == fallingTorchStageIndex {
			rt.FallingTorchWarningTicks = fallingTorchWarningDuration
			rt.FallingTorchTriggers++
		}
		rt.clearForegroundBlob(x, y, 1)
	case 4:
		rt.activateCheckpointAt(x, y)
	case 5, 28:
		rt.ReachedGoal = true
		rt.GoalExitSecret = foregroundID == 28
		rt.GoalExitComplete = false
		direction := int(rt.Background[rt.index(x, y)])
		if direction < 1 || direction > 4 {
			direction = 2
		}
		rt.GoalExitDirection = direction
	}
	rt.updatePressureDoors()
	return true
}

// AdvancePlayerMotion applies the source jInt cadence after the object scan.
// It returns true on the frame the interpolation reaches zero.
func (rt *Runtime) AdvancePlayerMotion() bool {
	if rt.PlayerMotion.Remaining <= 0 {
		rt.PlayerMotion.Remaining = 0
		return false
	}
	rt.PlayerMotion.Remaining = max(0, rt.PlayerMotion.Remaining-playerMoveStep)
	rt.updatePressureDoors()
	return rt.PlayerMotion.Remaining == 0
}

func (rt *Runtime) SetPlayerTurnOffset(offset int) {
	if rt == nil {
		return
	}
	rt.playerTurnOffset = max(0, offset)
}

// SetPlayerFacing stores the source player direction bits:
// up=1, right=2, down=3, left=4.
func (rt *Runtime) SetPlayerFacing(dx, dy int) {
	if rt == nil {
		return
	}
	switch {
	case dy < 0:
		rt.playerFacingDirection = 1
	case dx > 0:
		rt.playerFacingDirection = 2
	case dy > 0:
		rt.playerFacingDirection = 3
	case dx < 0:
		rt.playerFacingDirection = 4
	}
}

func (rt *Runtime) playerSourceOffset() int {
	if rt == nil {
		return 0
	}
	if rt.PlayerMotion.Remaining > 0 {
		return rt.PlayerMotion.Remaining
	}
	return rt.playerTurnOffset
}

// AdvanceGoalExit performs the source xBoolean auto-walk. Raw foreground 5 and
// 28 store separate normal/secret directions; both use this movement cadence
// before their completion branches diverge.
func (rt *Runtime) AdvanceGoalExit() (moved, complete bool) {
	if !rt.ReachedGoal || rt.GoalExitComplete || rt.PlayerMotion.Remaining > 0 {
		return false, rt.GoalExitComplete
	}
	dx, dy := snakeStep(rt.GoalExitDirection)
	if dx == 0 && dy == 0 {
		dx = 1
		rt.GoalExitDirection = 2
	}
	rt.Player.X += dx
	rt.Player.Y += dy
	rt.PlayerMotion = ObjectMotion{DX: dx, DY: dy, Remaining: playerMoveStartOffset}
	rt.GoalExitComplete = rt.Player.X < -5 || rt.Player.X > rt.Width()+5 || rt.Player.Y < -5 || rt.Player.Y > rt.Height()+5
	return true, rt.GoalExitComplete
}

func (rt *Runtime) pushAttemptReady(x, y, dx int) bool {
	target := Point{X: x, Y: y}
	if !rt.Pushing || rt.PushDX != dx || rt.pushTarget != target {
		rt.Pushing = true
		rt.PushDX = dx
		rt.PushTicks = 0
		rt.pushTarget = target
	}
	rt.PushTicks++
	return rt.PushTicks >= boulderPushAttempts
}

func (rt *Runtime) ResetPushAttempt() {
	rt.Pushing = false
	rt.PushDX = 0
	rt.PushTicks = 0
	rt.pushTarget = Point{}
}

func (rt *Runtime) activateCheckpointAt(x, y int) bool {
	if !rt.IsCheckpoint(x, y) {
		return false
	}
	order := rt.checkpointOrderAt(x, y)
	if order < rt.CheckpointProgress {
		return false
	}
	rt.CheckpointProgress = order + 1
	rt.CheckpointPending = true
	rt.pendingCheckpoint = Point{X: x, Y: y}
	return true
}

func (rt *Runtime) CommitPendingCheckpoint() bool {
	if !rt.CheckpointPending || rt.playerSourceOffset() > 0 || rt.Player != rt.pendingCheckpoint {
		return false
	}
	rt.CheckpointPending = false
	rt.SaveSnapshot()
	rt.emitSound(SoundCheckpoint)
	return true
}

func (rt *Runtime) collectForegroundEventAt(x, y int) {
	idx := rt.index(x, y)
	rt.LastForegroundEvent = int(rt.Background[idx])
	rt.ForegroundEvents++
	rt.set(rt.Foreground, x, y, EmptyRawID)
}

func (rt *Runtime) collectBonusAt(x, y, value int) {
	rt.collectBonusValue(value)
	rt.set(rt.PlayerLayer, x, y, EmptyRawID)
}

func (rt *Runtime) collectBonusValue(value int) {
	if value == int(EmptyRawID) || value < 0 {
		value = 0
	}
	rt.VioletGems += value
	rt.BonusValue += value
	rt.BonusPickups++
	rt.consumeBonusQuota(value)
	if rt.BonusRemaining < 0 {
		rt.BonusRemaining = 0
	}
}

func (rt *Runtime) collectSpecialItemAt(x, y, mask int) {
	rt.SpecialItemMask |= mask
	rt.SpecialPickups++
	rt.set(rt.PlayerLayer, x, y, EmptyRawID)
	if mask == 1 {
		rt.queueTutorialScript(22)
	}
}

func (rt *Runtime) consumeBonusQuota(value int) {
	if value <= 0 || !rt.BonusTargetSet {
		return
	}
	rt.BonusRemaining -= value
	if rt.BonusRemaining <= 0 {
		rt.BonusGateOpen = true
		if rt.BonusTargetSet {
			idx := rt.index(rt.BonusTarget.X, rt.BonusTarget.Y)
			rt.PlayerLayer[idx] = EmptyRawID
			rt.ObjectState[idx] = 0
			rt.ObjectMotion[idx] = ObjectMotion{}
		}
		rt.updateExitOpen()
	}
}

func (rt *Runtime) TryPushBoulder(x, y, dx int) bool {
	if dx == 0 {
		return false
	}
	id, ok := rt.At(PlayerLayer, x, y)
	if !ok || id != 0 && id != 8 && id != 9 {
		return false
	}
	targetX := x + dx
	if !rt.cellEmptyForObject(targetX, y) {
		return false
	}
	rt.moveObject(x, y, targetX, y, id)
	return true
}

func (rt *Runtime) TickGravity() int {
	rt.gravitySourceTick++
	return rt.tickGravityBounds(0, rt.Width()-1, 0, rt.Height()-1)
}

func (rt *Runtime) TickGravityNearPlayer(radius int) int {
	rt.gravitySourceTick++
	return rt.tickGravityBounds(
		max(1, rt.Player.X-radius),
		min(rt.Width()-2, rt.Player.X+radius),
		max(1, rt.Player.Y-radius),
		min(rt.Height()-2, rt.Player.Y+radius),
	)
}

func (rt *Runtime) tickGravityBounds(minX, maxX, minY, maxY int) int {
	moved := 0
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			if rt.tickGravityObjectAt(x, y) {
				moved++
			}
		}
	}
	if moved > 0 {
		rt.updatePressureDoors()
	}
	return moved
}

func (rt *Runtime) tickGravityObjectAt(x, y int) bool {
	idx := rt.index(x, y)
	id := rt.PlayerLayer[idx]
	if !isGravityObject(id) {
		return false
	}
	if id == 1 && rt.objectOverlapsPlayer(x, y) {
		rt.collectVioletAt(x, y)
		return false
	}
	if rt.advanceObjectMotion(idx, 6, 0) {
		if id == 0 && (rt.ObjectMotion[idx].Remaining == 12 || rt.ObjectMotion[idx].Remaining == 0) {
			rt.advanceBoulderRotation(idx)
		}
		if rt.ObjectMotion[idx].Remaining == 0 {
			rt.finishGravityMotionAt(x, y, id)
		}
		return false
	}
	motion := &rt.ObjectMotion[idx]
	if rt.gravityObjectBuoyantAt(x, y) {
		motion.RollDX = 0
		motion.RollOffset = 0
		if y > 0 && rt.tryMoveGravityObject(x, y, x, y-1, id) {
			return true
		}
		rt.ObjectState[idx] &^= objectDirectionMask | gravityRollPreparing | gravityMoveRight | gravityMoveLeft | explosiveFallMask
		rt.ObjectMotion[idx] = ObjectMotion{}
		return false
	}
	if motion.RollDX != 0 {
		return rt.tickGravityRollAt(x, y, id)
	}
	if rt.tryMoveGravityObject(x, y, x, y+1, id) {
		return true
	}
	if rt.canStartGravityRoll(x, y, -1) {
		rt.ObjectMotion[idx] = ObjectMotion{RollDX: -1, RollOffset: 1}
		rt.ObjectState[idx] = (rt.ObjectState[idx] &^ (objectDirectionMask | gravityMoveRight | gravityMoveLeft)) | 4 | gravityMoveLeft | gravityRollPreparing
		return false
	}
	if rt.canStartGravityRoll(x, y, 1) {
		rt.ObjectMotion[idx] = ObjectMotion{RollDX: 1, RollOffset: 1}
		rt.ObjectState[idx] = (rt.ObjectState[idx] &^ (objectDirectionMask | gravityMoveRight | gravityMoveLeft)) | 2 | gravityMoveRight | gravityRollPreparing
		return false
	}
	rt.ObjectState[idx] &^= objectDirectionMask
	belowID, _ := rt.At(PlayerLayer, x, y+1)
	if !isRoundedGravitySupport(belowID) {
		rt.ObjectState[idx] &^= gravityMoveRight | gravityMoveLeft
	}
	rt.ObjectMotion[idx] = ObjectMotion{}
	return false
}

func (rt *Runtime) gravityObjectBuoyantAt(x, y int) bool {
	cell := rt.waterCellAt(x, y)
	return cell != 0 && cell != 3
}

func (rt *Runtime) finishGravityMotionAt(x, y int, id RawID) {
	idx := rt.index(x, y)
	if rt.ObjectState[idx]&objectDirectionMask != 3 || y+1 >= rt.Height() {
		return
	}
	if rt.cellEmptyForGravity(x, y+1) && !rt.isPlayerAt(x, y+1) {
		return
	}
	fallDistance := (rt.ObjectState[idx] & explosiveFallMask) >> explosiveFallShift
	if fallDistance >= 2 {
		if id == 8 {
			rt.startExplosionAt(x, y)
			return
		}
		if belowID, ok := rt.At(PlayerLayer, x, y+1); ok && belowID == 8 {
			rt.startExplosionAt(x, y+1)
		}
	}
	rt.ObjectState[idx] &^= objectDirectionMask
	rt.ObjectState[idx] &^= explosiveFallMask
	if id == 0 || id == 9 {
		rt.emitSound(SoundBoulder)
	}
}

func (rt *Runtime) startExplosionAt(x, y int) bool {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return false
	}
	idx := rt.index(x, y)
	if rt.PlayerLayer[idx] != 8 {
		return false
	}
	rt.PlayerLayer[idx] = 54
	rt.ObjectState[idx] = 0
	rt.ObjectMotion[idx] = ObjectMotion{}
	return true
}

func (rt *Runtime) tickExplosionAt(x, y int) {
	idx := rt.index(x, y)
	if rt.PlayerLayer[idx] != 54 {
		return
	}
	rt.ObjectState[idx]++
	tick := rt.ObjectState[idx]
	if tick == 1 {
		rt.emitSound(SoundBossDeath)
	}
	if tick == explosionImpactTick {
		rt.applyExplosionImpact(x, y)
	}
	if tick >= explosionDuration && rt.PlayerLayer[idx] == 54 {
		rt.PlayerLayer[idx] = EmptyRawID
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
	}
}

func (rt *Runtime) applyExplosionImpact(x, y int) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx, ny := x+dx, y+dy
			if nx < 0 || ny < 0 || nx >= rt.Width() || ny >= rt.Height() {
				continue
			}
			idx := rt.index(nx, ny)
			switch rt.PlayerLayer[idx] {
			case 10:
				// JAR bN() only primes vegetation while the packed water
				// state is stable (xByte == 3).
				if rt.waterStable() && rt.ObjectState[idx] < 1 {
					rt.ObjectState[idx] = 1
				}
			case 8:
				rt.PlayerLayer[idx] = 54
				rt.ObjectState[idx] = 0
				rt.ObjectMotion[idx] = ObjectMotion{}
			case 30, 37:
				if rt.ObjectState[idx] < 1 {
					rt.ObjectState[idx] = 1
				}
			case 16, 19, 43, 49:
				rt.decrementEnemyGateForObjectAt(nx, ny)
				rt.PlayerLayer[idx] = EmptyRawID
				rt.ObjectState[idx] = 0
				rt.ObjectMotion[idx] = ObjectMotion{}
			}
			if rt.isPlayerAt(nx, ny) {
				rt.Hurt(1)
			}
		}
	}
}

func sourceSpikeExtent(sourceTick, period, closedTicks, extendTicks, openTicks, retractTicks int) int {
	phase := sourceTick % period
	switch {
	case phase < closedTicks:
		return 0
	case phase < closedTicks+extendTicks:
		return 48 * (phase - closedTicks) / extendTicks
	case phase < closedTicks+extendTicks+openTicks:
		return 48
	default:
		return 48 - 48*(phase-closedTicks-extendTicks-openTicks)/retractTicks
	}
}

func (rt *Runtime) spikeExtentAt(x, y int) int {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return 0
	}
	if rt.ObjectState[rt.index(x, y)]&0x8 != 0 {
		return rt.SpikeFastExtent
	}
	return rt.SpikeSlowExtent
}

func (rt *Runtime) SpikeExtentAt(x, y int) int {
	return rt.spikeExtentAt(x, y)
}

func (rt *Runtime) spikeTipAt(x, y int) (Point, bool) {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() || rt.PlayerLayer[rt.index(x, y)] != 28 {
		return Point{}, false
	}
	extent := rt.spikeExtentAt(x, y)
	segments := 1
	if extent > 0 {
		segments = (extent-1)/TileSize + 2
	}
	direction := -1
	if rt.ObjectState[rt.index(x, y)]&objectDirectionMask == 3 {
		direction = 1
	}
	tip := Point{X: x, Y: y + (segments-1)*direction}
	if tip.Y < 0 || tip.Y >= rt.Height() {
		return Point{}, false
	}
	return tip, true
}

func (rt *Runtime) spikeOccupies(x, y int) bool {
	for _, baseY := range []int{y - 1, y + 1} {
		if baseY < 0 || baseY >= rt.Height() || rt.PlayerLayer[rt.index(x, baseY)] != 28 {
			continue
		}
		state := rt.ObjectState[rt.index(x, baseY)]
		if rt.SpikeFastExtent >= TileSize || state&0x8 == 0 && rt.SpikeSlowExtent >= TileSize {
			return true
		}
	}
	return false
}

func (rt *Runtime) tickSpikeColumnAt(x, y int) {
	tip, ok := rt.spikeTipAt(x, y)
	if !ok {
		return
	}
	if rt.isPlayerAt(tip.X, tip.Y) {
		rt.HurtFromDirection(2, rt.playerCollisionDirection())
	}
	idx := rt.index(tip.X, tip.Y)
	id := rt.PlayerLayer[idx]
	if id == EmptyRawID || id == 28 || id == 32 {
		return
	}
	rt.decrementEnemyGateForObjectAt(tip.X, tip.Y)
	rt.PlayerLayer[idx] = EmptyRawID
	rt.ObjectState[idx] = 0
	rt.ObjectMotion[idx] = ObjectMotion{}
}

func (rt *Runtime) tickFanPhase(sourceTick int) {
	if rt.Stage.World != WorldBavaria || rt.FanDirection == 0 || (sourceTick>>1)&1 != 0 {
		return
	}
	rt.FanPhase += rt.FanDirection
	if rt.FanPhase == 5 {
		rt.swapFanPods()
	}
	if rt.FanPhase <= 0 {
		rt.FanPhase = 0
		rt.FanDirection = 0
	} else if rt.FanPhase >= 9 {
		rt.FanPhase = 9
		rt.FanDirection = 0
	}
}

func (rt *Runtime) swapFanPods() {
	for idx := range rt.PlayerLayer {
		x, y := idx%rt.Width(), idx/rt.Width()
		switch {
		case rt.Foreground[idx] == 15:
			rt.Foreground[idx] = EmptyRawID
			rt.PlayerLayer[idx] = 34
		case rt.Foreground[idx] == 16:
			rt.Foreground[idx] = EmptyRawID
			rt.PlayerLayer[idx] = 35
		case rt.PlayerLayer[idx] == 34:
			rt.PlayerLayer[idx] = EmptyRawID
			rt.Foreground[idx] = 15
			rt.triggerWaterReflow(x, y)
		case rt.PlayerLayer[idx] == 35:
			rt.PlayerLayer[idx] = EmptyRawID
			rt.Foreground[idx] = 16
			rt.triggerWaterReflow(x, y)
		}
	}
}

func (rt *Runtime) tickMovingHazardAt(x, y int) {
	idx := rt.index(x, y)
	if rt.PlayerLayer[idx] != 14 {
		return
	}
	state := rt.ObjectState[idx]
	direction := 2
	if state&0x8 != 0 {
		direction = 4
	}
	if rt.isPlayerAt(x, y) {
		rt.HurtFromDirection(1, direction)
	}
	if rt.advanceObjectMotion(idx, 6, 0) {
		return
	}
	cooldown := (state & movingHazardCooldownMask) >> movingHazardCooldownShift
	dx := 1
	if state&0x8 != 0 {
		dx = -1
	}
	if cooldown >= 20 {
		blockedID, _ := rt.At(PlayerLayer, x+dx, y)
		pathReady := rt.hazardCellOpen(x, y+1) || rt.hazardCellOpen(x+dx, y) || blockedID == 16 || blockedID == 19 || blockedID == 43
		if pathReady {
			rt.ObjectState[idx] = (state &^ movingHazardCooldownMask) | 19<<movingHazardCooldownShift
		}
		return
	}
	if cooldown > 0 {
		cooldown--
		rt.ObjectState[idx] = (state &^ movingHazardCooldownMask) | cooldown<<movingHazardCooldownShift
		return
	}
	if rt.hazardCellOpen(x, y+1) {
		rt.moveStatefulObject(x, y, x, y+1, 14, (state&0x8)|3)
		return
	}
	if rt.hazardCellOpen(x+dx, y) {
		rt.moveStatefulObject(x, y, x+dx, y, 14, (state&0x8)|direction)
		return
	}
	blockedID, _ := rt.At(PlayerLayer, x+dx, y)
	if blockedID == 16 || blockedID == 19 || blockedID == 43 {
		rt.ObjectState[idx] = state &^ objectDirectionMask
		return
	}
	rt.ObjectState[idx] = (state &^ (movingHazardCooldownMask | objectDirectionMask)) | 20<<movingHazardCooldownShift
}

func (rt *Runtime) hazardCellOpen(x, y int) bool {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() || rt.PlayerLayer[rt.index(x, y)] != EmptyRawID {
		return false
	}
	foreground := rt.Foreground[rt.index(x, y)]
	if foreground == 14 || foreground == 33 || foreground == 4 || foreground == 32 {
		return false
	}
	return foreground != 7 || rt.foregroundDoorOpen(x, y)
}

func (rt *Runtime) moveStatefulObject(fromX, fromY, toX, toY int, id RawID, state int) {
	fromIdx := rt.index(fromX, fromY)
	toIdx := rt.index(toX, toY)
	rt.PlayerLayer[toIdx] = id
	rt.PlayerLayer[fromIdx] = EmptyRawID
	rt.ObjectState[toIdx] = state
	rt.ObjectState[fromIdx] = 0
	rt.ObjectMotion[toIdx] = ObjectMotion{DX: toX - fromX, DY: toY - fromY, Remaining: 18}
	rt.ObjectMotion[fromIdx] = ObjectMotion{}
}

func (rt *Runtime) gravityObjectStrikesAt(x, y int) bool {
	if y <= 0 {
		return false
	}
	above := rt.index(x, y-1)
	id := rt.PlayerLayer[above]
	return isGravityObject(id) && rt.ObjectState[above]&objectDirectionMask == 3 && rt.ObjectMotion[above].Remaining <= 6
}

func (rt *Runtime) tickSpearPairAt(x, y int) {
	if y <= 0 || y+1 >= rt.Height() || rt.PlayerLayer[rt.index(x, y)] != 16 || rt.PlayerLayer[rt.index(x, y-1)] != 16 || rt.PlayerLayer[rt.index(x, y+1)] == 16 {
		return
	}
	baseIdx := rt.index(x, y)
	topIdx := rt.index(x, y-1)
	if rt.gravityObjectStrikesAt(x, y-1) || rt.movingHazardStrikesAt(x, y-1) || rt.movingHazardStrikesAt(x, y) {
		rt.emitSound(SoundBoulder)
		rt.decrementEnemyGateForObjectAt(x, y)
		rt.decrementEnemyGateForObjectAt(x, y-1)
		rt.PlayerLayer[baseIdx] = EmptyRawID
		rt.PlayerLayer[topIdx] = EmptyRawID
		rt.ObjectState[baseIdx] = 0
		rt.ObjectState[topIdx] = 0
		return
	}
	state := rt.ObjectState[baseIdx]
	direction := state & objectDirectionMask
	attackX := x - 1
	if direction == 4 {
		attackX = x + 1
	}
	adjacent := rt.isPlayerAt(attackX, y) || rt.isPlayerAt(attackX, y-1)
	timer := (state & spearTimerMask) >> spearTimerShift
	if timer <= 0 && adjacent && rt.playerSourceOffset() <= 12 {
		timer = 36
	} else if timer > 0 {
		timer--
	}
	if timer <= 11 && timer > 0 && adjacent {
		rt.HurtFromDirection(1, direction)
	}
	state = (state &^ spearTimerMask) | timer<<spearTimerShift
	rt.ObjectState[baseIdx] = state
	rt.ObjectState[topIdx] = state
}

func (rt *Runtime) movingHazardStrikesAt(x, y int) bool {
	if x <= 0 || x+1 >= rt.Width() || y < 0 || y >= rt.Height() {
		return false
	}
	if y > 0 {
		aboveIdx := rt.index(x, y-1)
		if rt.PlayerLayer[aboveIdx] == 14 && rt.ObjectMotion[aboveIdx].Remaining <= 6 {
			return true
		}
	}
	leftIdx := rt.index(x-1, y)
	if rt.PlayerLayer[leftIdx] == 14 && rt.ObjectMotion[leftIdx].Remaining <= 0 {
		state := rt.ObjectState[leftIdx]
		if state&0x8 == 0 && state&objectDirectionMask != 3 {
			return true
		}
	}
	rightIdx := rt.index(x+1, y)
	if rt.PlayerLayer[rightIdx] == 14 && rt.ObjectMotion[rightIdx].Remaining <= 0 {
		state := rt.ObjectState[rightIdx]
		if state&0x8 != 0 && state&objectDirectionMask != 3 {
			return true
		}
	}
	return false
}

func (rt *Runtime) tickCrawlerTrapAt(x, y int) {
	idx := rt.index(x, y)
	if rt.PlayerLayer[idx] != 36 || y <= 0 {
		return
	}
	if rt.ObjectState[idx] == 0 && rt.PlayerLayer[rt.index(x, y-1)] == 11 {
		rt.ObjectState[idx] = 1
		rt.decrementEnemyGateForObjectAt(x, y)
	}
	if rt.ObjectState[idx] == 1 && rt.isPlayerAt(x, y-1) {
		rt.Hurt(1)
	}
}

// GravityObjectRenderOffset mirrors the positional part of Java OVoid().
func (rt *Runtime) GravityObjectRenderOffset(x, y, sourceTick int) (int, int) {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return 0, 0
	}
	idx := rt.index(x, y)
	motion := rt.ObjectMotion[idx]
	dx := -motion.DX * motion.Remaining
	dy := -motion.DY*motion.Remaining + rt.waterBobOffsetAt(x, y, sourceTick)
	state := rt.ObjectState[idx]
	if state&gravityRollPreparing != 0 && motion.RollDX != 0 {
		offset := motion.RollOffset
		dx += motion.RollDX*offset + 1 - sourceTick%3
		dy += offset * offset / TileSize
		return dx, dy
	}
	if motion.Remaining > 0 && rt.gravityObjectUsesSupportArc(x, y, state) {
		offset := motion.Remaining
		dx += -1 + sourceTick%3
		dy += offset * offset / TileSize
	}
	return dx, dy
}

func (rt *Runtime) gravityObjectUsesSupportArc(x, y, state int) bool {
	direction := state & objectDirectionMask
	aheadX := x
	switch direction {
	case 2:
		aheadX++
	case 4:
		aheadX--
	default:
		return false
	}
	if aheadX < 0 || aheadX >= rt.Width() || y+1 >= rt.Height() {
		return false
	}
	aheadBelow, _ := rt.At(PlayerLayer, aheadX, y+1)
	below, _ := rt.At(PlayerLayer, x, y+1)
	belowIdx := rt.index(x, y+1)
	if aheadBelow != EmptyRawID || !isRoundedGravitySupport(below) || rt.ObjectState[belowIdx]&objectDirectionMask != 0 {
		return false
	}
	return !rt.Hooking || (x != rt.HookTarget.X && y != rt.HookTarget.Y)
}

func (rt *Runtime) canStartGravityRoll(x, y, dx int) bool {
	if y+1 >= rt.Height() {
		return false
	}
	belowIdx := rt.index(x, y+1)
	if !isGravityObject(rt.PlayerLayer[belowIdx]) {
		return false
	}
	belowMotion := rt.ObjectMotion[belowIdx]
	if belowMotion.Remaining > 0 || belowMotion.RollDX != 0 {
		return false
	}
	return !rt.isPlayerAt(x+dx, y) &&
		rt.cellEmptyForGravity(x+dx, y) &&
		rt.cellEmptyForGravity(x+dx, y+1)
}

func (rt *Runtime) tickGravityRollAt(x, y int, id RawID) bool {
	idx := rt.index(x, y)
	motion := &rt.ObjectMotion[idx]
	if y+1 >= rt.Height() {
		*motion = ObjectMotion{}
		return false
	}
	if rt.cellEmptyForGravity(x, y+1) && !rt.isPlayerAt(x, y+1) {
		motion.RollOffset = max(0, motion.RollOffset-6)
		if motion.RollOffset == 0 {
			*motion = ObjectMotion{}
			rt.ObjectState[idx] &^= gravityRollPreparing | objectDirectionMask
		}
		return false
	}

	belowMotion := rt.ObjectMotion[rt.index(x, y+1)]
	if rt.canRollBoulder(x, y, motion.RollDX) && belowMotion.RollDX == 0 {
		if motion.RollOffset >= gravityRollFastOffset || rt.gravitySourceTick&3 == 0 {
			motion.RollOffset++
		}
		if motion.RollOffset < gravityRollMoveOffset {
			return false
		}
		dx := motion.RollDX
		rt.ObjectState[idx] = (rt.ObjectState[idx] &^ (gravityRollPreparing | objectDirectionMask)) | 3
		rt.moveObjectWithMotion(x, y, x+dx, y+1, id, ObjectMotion{DY: 1, Remaining: 12})
		return true
	}

	motion.RollOffset = max(0, motion.RollOffset-6)
	if motion.RollOffset == 0 {
		*motion = ObjectMotion{}
		rt.ObjectState[idx] &^= gravityRollPreparing | objectDirectionMask
	}
	return false
}

func (rt *Runtime) TickDigAnimations() int {
	rt.gravitySourceTick++
	cleared := 0
	for idx, id := range rt.Foreground {
		if id != 32 {
			continue
		}
		if rt.tickDigAnimationAt(idx, rt.gravitySourceTick) {
			cleared++
		}
	}
	return cleared
}

func (rt *Runtime) TickSourceFrame(radius, sourceTick, hazardReach int) SourceFrameResult {
	rt.gravitySourceTick = sourceTick
	rt.tickWaterSourceStart()
	rt.SpikeSlowExtent = sourceSpikeExtent(sourceTick, 89, 15, 30, 15, 30)
	rt.SpikeFastExtent = sourceSpikeExtent(sourceTick, 44, 7, 15, 8, 15)
	rt.tickFanPhase(sourceTick)
	rt.frameVioletPickups = rt.frameVioletPickups[:0]
	chestOpeningFresh := rt.chestOpeningFresh
	tutorialSealWasActivated := rt.TutorialSealActivated
	rt.TickStatus()
	result := SourceFrameResult{}
	if rt.tickFallingTorchStage(sourceTick) {
		result.RisingFireHits++
	}
	rt.tickForegroundDemo()
	rt.tickEnemyGateDemo()
	result.AnacondaHits, result.AnacondaDefeated = rt.tickGreatAnaconda(sourceTick)
	result.TeutonicKnightHits, result.TeutonicKnightDefeated = rt.tickEvilTeutonicKnight(sourceTick)
	if !rt.Anaconda.Enabled && !rt.TeutonicKnight.Enabled {
		rt.tickSealTransition()
	}
	rt.tickDoorAnimations(sourceTick)
	if rt.playerSourceOffset() <= 0 {
		rt.CommitPendingCheckpoint()
	}
	minX := max(1, rt.Player.X-radius)
	maxX := min(rt.Width()-2, rt.Player.X+radius)
	minY := max(1, rt.Player.Y-radius)
	maxY := min(rt.Height()-2, rt.Player.Y+radius)
	if rt.tickRockHold() {
		result.RockHoldHits++
	}
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			idx := rt.index(x, y)
			if isPickupContainer(rt.Foreground[idx]) {
				rt.tickChestForegroundAt(x, y, sourceTick, chestOpeningFresh)
			}
			if rt.Foreground[idx] == 35 || rt.Foreground[idx] == 37 {
				rt.tickWindForegroundAt(x, y)
			}
			if rt.Foreground[idx] == 32 && rt.tickDigAnimationAt(idx, sourceTick) {
				result.DigCleared++
			}
			if rt.Hooking && rt.hookReturning && rt.HookTarget == (Point{X: x, Y: y}) {
				continue
			}
			switch id := rt.PlayerLayer[idx]; {
			case id == 47:
				if rt.tickGravityObjectAt(x, y) {
					result.GravityMoved++
				}
				rt.tickWindPodAt(x, y)
			case isGravityObject(id):
				if rt.tickGravityObjectAt(x, y) {
					result.GravityMoved++
				}
			case isSnake(id):
				if rt.tickSnakeObjectAt(x, y) {
					result.SnakesMoved++
				}
			case id == 10 && rt.ObjectState[idx] > 0:
				rt.PlayerLayer[idx] = EmptyRawID
				rt.Foreground[idx] = 32
				rt.ForegroundState[idx] = 0
				rt.ObjectState[idx] = 0
				rt.triggerWaterReflow(x, y)
			case id == 11:
				if rt.tickCrawlerObjectAt(x, y) {
					result.CrawlersMoved++
				}
			case id == 22 || id == 23:
				if rt.horizontalHazardHitsPlayer(x, y, id, hazardReach) && rt.Hurt(1) {
					result.HazardHits++
				}
			case id == 50:
				if rt.playerSourceOffset() < gravityRollMoveOffset && rt.isPlayerAt(x, y) && rt.HurtFromDirection(1, rt.playerCollisionDirection()) {
					result.HazardHits++
				}
			case id == 54:
				rt.tickExplosionAt(x, y)
			case id == 28:
				rt.tickSpikeColumnAt(x, y)
			case id == 14:
				rt.tickMovingHazardAt(x, y)
			case id == 16:
				rt.tickSpearPairAt(x, y)
			case id == 36:
				rt.tickCrawlerTrapAt(x, y)
			}
		}
	}
	// Java advances the packed water state in asVoid(), after ahVoid() has
	// scanned gravity objects and hazards for the current source frame.
	rt.tickWater()
	rt.updatePressureDoors()
	rt.startAdjacentLockOpening()
	rt.tickPendingForegroundEvent()
	rt.tickTutorial()
	result.TutorialSealActivated = !tutorialSealWasActivated && rt.TutorialSealActivated
	result.VioletPickups = append(result.VioletPickups, rt.frameVioletPickups...)
	return result
}

func (rt *Runtime) tickPendingForegroundEvent() {
	if !rt.pendingForegroundEventSet || rt.playerSourceOffset() > 6 {
		return
	}
	point := rt.pendingForegroundEvent
	rt.pendingForegroundEventSet = false
	if rt.Player != point {
		return
	}
	id, ok := rt.At(ForegroundLayer, point.X, point.Y)
	if !ok {
		return
	}
	switch id {
	case 0:
		eventID := int(rt.Background[rt.index(point.X, point.Y)])
		rt.collectForegroundEventAt(point.X, point.Y)
		rt.startTutorialForegroundEvent(eventID)
		if rt.Stage.World == WorldAngkor && rt.Stage.Index == fallingTorchStageIndex && eventID == 3 {
			rt.ForegroundDemoActive = true
			rt.ForegroundDemoID = eventID
			rt.ForegroundDemoPhase = 0
			rt.ForegroundDemoTicks = 0
			rt.foregroundDemoMoved = false
		}
	case 26:
		rt.activateEnemyGateTriggerAt(point.X, point.Y)
	}
}

func (rt *Runtime) tickForegroundDemo() {
	if !rt.ForegroundDemoActive || rt.ForegroundDemoID != 3 {
		return
	}
	switch rt.ForegroundDemoPhase {
	case 0:
		if !rt.foregroundDemoMoved {
			if rt.tryMove(0, -1, true) {
				rt.foregroundDemoMoved = true
			}
			return
		}
		if rt.PlayerMotion.Remaining > 0 {
			return
		}
		rt.ForegroundDemoPhase = 1
		rt.ForegroundDemoTicks = 0
	case 1:
		rt.ForegroundDemoTicks++
		if rt.ForegroundDemoTicks > foregroundDemoPanTicks {
			rt.ForegroundDemoPhase = 2
			rt.ForegroundDemoTicks = 0
		}
	case 2:
		rt.ForegroundDemoTicks++
		if rt.ForegroundDemoTicks > foregroundDemoWaitTicks {
			rt.ForegroundDemoActive = false
			rt.ForegroundDemoPhase = 0
			rt.ForegroundDemoTicks = 0
			rt.foregroundDemoMoved = false
		}
	}
}

func (rt *Runtime) tickFallingTorchStage(sourceTick int) bool {
	if rt.Stage.World != WorldAngkor || rt.Stage.Index != fallingTorchStageIndex {
		return false
	}
	if rt.FallingTorchWarningTicks > 0 {
		rt.FallingTorchWarningTicks--
	}
	trigger := Point{X: 18, Y: 63}
	if rt.FallingTorchTriggers == 0 && trigger.X < rt.Width() && trigger.Y < rt.Height() {
		idx := rt.index(trigger.X, trigger.Y)
		motion := rt.ObjectMotion[idx]
		if rt.PlayerLayer[idx] == 0 && motion.Remaining <= 0 && motion.RollDX == 0 {
			rt.FallingTorchWarningTicks = fallingTorchWarningDuration
			rt.FallingTorchTriggers = 1
		}
	}
	if rt.FallingTorchTriggers == 3 {
		switch rt.FallingTorchAnimation {
		case 0:
			rt.FallingTorchAnimation = 1
			rt.FallingTorchAnimationTicks = 0
		case 1:
			if rt.FallingTorchAnimationTicks >= fallingTorchCollapseTicks {
				rt.FallingTorchAnimation = 2
				rt.FallingTorchAnimationTicks = 0
			}
		}
	} else if rt.FallingTorchAnimation != 0 {
		rt.FallingTorchAnimation = 0
		rt.FallingTorchAnimationTicks = 0
	}
	rt.FallingTorchAnimationTicks++
	if rt.FallingTorchAnimation != 2 {
		return false
	}
	if rt.FallingTorchWarningTicks == 10 {
		rt.FallingTorchWarningTicks = fallingTorchWarningLoop
	}
	if rt.RisingFireAnimation != 0 {
		rt.RisingFireAnimationTicks++
		if rt.RisingFireAnimation == 2 && rt.RisingFireAnimationTicks >= fallingFireStartTicks {
			rt.RisingFireAnimation = 0
			rt.RisingFireAnimationTicks = 0
		}
		return false
	}
	rt.RisingFireAnimationTicks++
	if rt.ForegroundDemoActive {
		return false
	}
	if rt.RisingFireHeight < fallingFireMaximumHeight {
		rt.RisingFireHeight++
		viewportY := rt.viewportY
		if !rt.viewportSet {
			viewportY = clampRuntime(rt.Player.Y*TileSize-160, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
		}
		minimumHeight := rt.Height()*TileSize - (viewportY + ScreenHeight - 80)
		if rt.RisingFireHeight < minimumHeight {
			rt.RisingFireHeight = minimumHeight
		}
		if rt.RisingFireHeight > fallingFireMaximumHeight {
			rt.RisingFireHeight = fallingFireMaximumHeight
		}
	}
	if rt.Height()*TileSize-rt.RisingFireHeight <= rt.Player.Y*TileSize+18 && rt.Player.X < 17 {
		return rt.HurtFromDirection(rt.MaxHealth, 1)
	}
	return false
}

func (rt *Runtime) SetViewport(x, y int) {
	rt.viewportX = clampRuntime(x, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	rt.viewportY = clampRuntime(y, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	rt.viewportSet = true
}

func (rt *Runtime) SetViewportY(y int) {
	rt.SetViewport(rt.viewportX, y)
}

func (rt *Runtime) ForegroundDemoCamera() (x, y, elapsed, duration int, ok bool) {
	if !rt.ForegroundDemoActive || rt.ForegroundDemoID != 3 || rt.ForegroundDemoPhase < 1 {
		return 0, 0, 0, 0, false
	}
	x = clampRuntime(12*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	y = clampRuntime(42*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	duration = foregroundDemoPanTicks
	if rt.ForegroundDemoPhase == 1 {
		elapsed = clampRuntime(rt.ForegroundDemoTicks, 0, duration)
	} else {
		elapsed = duration
	}
	return x, y, elapsed, duration, true
}

func (rt *Runtime) RisingFireWorldY() (int, bool) {
	if rt.Stage.World != WorldAngkor || rt.Stage.Index != fallingTorchStageIndex || rt.FallingTorchAnimation != 2 {
		return 0, false
	}
	return rt.Height()*TileSize - rt.RisingFireHeight, true
}

func (rt *Runtime) IsFallingTorchStage() bool {
	return rt != nil && rt.Stage != nil && rt.Stage.World == WorldAngkor && rt.Stage.Index == fallingTorchStageIndex
}

func (rt *Runtime) RisingFireFillVisible() bool {
	return rt.IsFallingTorchStage() && rt.FallingTorchAnimation == 2 && rt.RisingFireHeight > fallingFireInitialHeight
}

func (rt *Runtime) FallingTorchShake(sourceTick int) int {
	if rt.Stage.World != WorldAngkor || rt.Stage.Index != fallingTorchStageIndex || rt.FallingTorchWarningTicks <= 0 {
		return 0
	}
	warning := rt.FallingTorchWarningTicks
	return warning * sourceTick % ((warning >> 1) + 1) % 12
}

func (rt *Runtime) tickRockHold() bool {
	if rt.PlayerDead || rt.Player.X < 0 || rt.Player.X >= rt.Width() || rt.Player.Y <= 0 || rt.Player.Y >= rt.Height() {
		return false
	}
	aboveIdx := rt.index(rt.Player.X, rt.Player.Y-1)
	if !isRockHoldObject(rt.PlayerLayer[aboveIdx]) || rt.ObjectMotion[aboveIdx].Remaining > 0 {
		return false
	}
	if waterCellGet(rt.waterCellAt(rt.Player.X, rt.Player.Y), 0, 0, 3) != 0 {
		return false
	}
	if rt.RockHoldTicks <= 0 {
		rt.RockHoldTicks = rockHoldDuration
	}
	rt.RockHoldTicks--
	if rt.RockHoldTicks > 0 {
		return false
	}
	rt.RockHoldTicks = rockHoldDuration
	return rt.Hurt(rt.MaxHealth)
}

func (rt *Runtime) HoldingRock() bool {
	if rt.PlayerDead || rt.Player.X < 0 || rt.Player.X >= rt.Width() || rt.Player.Y <= 0 || rt.Player.Y >= rt.Height() {
		return false
	}
	idx := rt.index(rt.Player.X, rt.Player.Y-1)
	return isRockHoldObject(rt.PlayerLayer[idx]) && rt.ObjectMotion[idx].Remaining <= 0
}

func isRockHoldObject(id RawID) bool {
	return id == 0 || id == 8 || id == 9 || id == 48
}

func (rt *Runtime) NextCompassTarget() (Point, bool) {
	if !rt.CompassEnabled {
		return Point{}, false
	}
	for _, checkpoint := range rt.Checkpoints {
		if rt.checkpointOrderAt(checkpoint.X, checkpoint.Y) == rt.CheckpointProgress {
			return checkpoint, true
		}
	}
	for _, goal := range rt.GoalMarkers {
		id, _ := rt.At(ForegroundLayer, goal.X, goal.Y)
		if id == 5 {
			return goal, true
		}
	}
	return Point{}, false
}

func (rt *Runtime) CompassDirection() (int, bool) {
	target, ok := rt.NextCompassTarget()
	if !ok {
		return 0, false
	}
	return CompassDirection(target.X-rt.Player.X, target.Y-rt.Player.Y), true
}

func CompassDirection(dx, dy int) int {
	if dy == 0 {
		if dx < 0 {
			return 12
		}
		return 4
	}
	if dx == 0 {
		if dy < 0 {
			return 0
		}
		return 8
	}
	ratio := (dx << 7) / dy
	if ratio > 0 {
		switch {
		case ratio < 128:
			if dx > 0 {
				return 7
			}
			return 15
		case ratio > 128:
			if dx > 0 {
				return 5
			}
			return 13
		default:
			if dx > 0 {
				return 6
			}
			return 14
		}
	}
	if ratio > -128 {
		if dx < 0 {
			return 9
		}
		return 1
	}
	if ratio < -128 {
		if dx < 0 {
			return 11
		}
		return 3
	}
	if dx < 0 {
		return 10
	}
	return 2
}

func (rt *Runtime) tickChestForegroundAt(x, y, sourceTick int, skipAdvance bool) {
	idx := rt.index(x, y)
	if rt.ContainerLocked[idx] {
		return
	}
	state := rt.ObjectState[idx]
	if state <= 0 {
		if rt.playerSourceOffset() > 0 || rt.Player != (Point{X: x, Y: y}) || !isContainerReward(rt.PlayerLayer[idx]) {
			return
		}
		rt.pendingChestSet = false
		rt.startChestOpening(Point{X: x, Y: y}, false)
		return
	}
	maxState := 3
	if rt.Foreground[idx] == 14 {
		maxState = 2
	}
	if !skipAdvance && state < maxState && (sourceTick>>1)&1 == 0 {
		rt.ObjectState[idx] = state + 1
	}
}

func (rt *Runtime) tickDigAnimationAt(idx, sourceTick int) bool {
	state := rt.ForegroundState[idx]
	if sourceTick&1 == 0 {
		state++
	}
	if state >= digAnimationFrames {
		rt.Foreground[idx] = EmptyRawID
		rt.ForegroundState[idx] = 0
		return true
	}
	rt.ForegroundState[idx] = state
	return false
}

func (rt *Runtime) UseHammer(dx, dy int) bool {
	if !rt.CanAcceptInput() || rt.specialToolLevel() < 1 || (dx == 0 && dy == 0) {
		return false
	}
	x := rt.Player.X + dx
	y := rt.Player.Y + dy
	id, ok := rt.At(PlayerLayer, x, y)
	if !ok {
		return false
	}
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	blockedTarget := id == 0 || id != EmptyRawID && id >= 80 || foregroundID == 7 && !rt.foregroundDoorOpen(x, y)
	canHit := id == 10 || id == 18 || id == 30 || isSnake(id) || blockedTarget
	if rt.specialToolLevel() >= 8 {
		canHit = canHit || id == 1 || id == 9 || rt.hasFreezeHammerTargetAt(x, y)
	}
	if !canHit {
		return false
	}
	rt.Hammering = true
	rt.HammerTicks = 0
	rt.HammerAnimation = hammerAnimation(dx, dy)
	rt.HammerTarget = Point{X: x, Y: y}
	return true
}

func hammerAnimation(dx, dy int) int {
	switch {
	case dy < 0:
		return 13
	case dx > 0:
		return 14
	case dy > 0:
		return 15
	default:
		return 16
	}
}

func hammerAnimationDuration(animation int) int {
	if animation == 13 {
		return hammerUpDuration
	}
	return hammerOtherDuration
}

func (rt *Runtime) tickHammerAction() {
	if !rt.Hammering {
		return
	}
	rt.HammerTicks++
	if rt.HammerTicks == hammerImpactTick {
		rt.applyHammerImpact()
	}
	if rt.HammerTicks < hammerAnimationDuration(rt.HammerAnimation) {
		return
	}
	rt.Hammering = false
	rt.HammerTicks = 0
	rt.HammerAnimation = 0
	rt.HammerTarget = Point{}
}

func (rt *Runtime) applyHammerImpact() bool {
	x := rt.HammerTarget.X
	y := rt.HammerTarget.Y
	id, ok := rt.At(PlayerLayer, x, y)
	if !ok {
		return false
	}
	idx := rt.index(x, y)
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	if id == 0 {
		rt.emitSound(SoundHammerBlock)
		if rt.specialToolLevel() < 8 {
			return true
		}
		if rt.freezeHammerTargetsAt(x, y) {
			rt.emitSound(SoundEnemyHit)
		}
		return true
	}
	if id != EmptyRawID && id >= 80 || foregroundID == 7 && !rt.foregroundDoorOpen(x, y) {
		rt.emitSound(SoundHammerBlock)
		return true
	}
	handled := false
	if id == 10 {
		// Water state 3 is the source's stable state. During a fill or
		// redistribution the hero can still move, but vegetation cannot be
		// disturbed until the water solver settles.
		if rt.waterStable() && rt.ObjectState[idx] <= 0 {
			rt.ObjectState[idx] = 1
		}
		return true
	}
	if id == 18 {
		playerForeground := rt.Foreground[rt.index(rt.Player.X, rt.Player.Y)]
		if rt.Stage.World == WorldBavaria && rt.waterStable() && rt.FanDirection == 0 && playerForeground != 15 && playerForeground != 16 {
			if rt.FanPhase <= 0 {
				rt.FanDirection = 1
			} else {
				rt.FanDirection = -1
			}
			rt.emitSound(SoundSwitch)
		}
		return true
	}
	if id == 30 {
		if rt.ObjectState[idx] == 0 {
			rt.ObjectState[idx] = 1
		}
		rt.emitSound(SoundBreak)
		// The source sets n13 for a breakable hit and skips the enhanced
		// hammer's five-cell freeze scan for this swing.
		return true
	}
	if rt.specialToolLevel() >= 8 {
		if id == 9 {
			if rt.thawFrozenAt(x, y) {
				rt.emitSound(SoundBreak)
				return true
			}
		}
		if rt.freezeHammerTargetsAt(x, y) {
			rt.emitSound(SoundEnemyHit)
			handled = true
		}
	}
	if handled {
		return true
	}
	switch {
	case id == 43 && rt.ObjectState[idx]&snakeStunMask == 0:
		if rt.ObjectState[idx]&0x18000 == 0 {
			rt.decrementEnemyGateForObjectAt(x, y)
			rt.PlayerLayer[idx] = EmptyRawID
			rt.ObjectState[idx] = 0
			rt.ObjectMotion[idx] = ObjectMotion{}
			rt.emitSound(SoundEnemyHit)
			return true
		}
		rt.ObjectState[idx] = packRedSnakeHammerTarget(rt.ObjectState[idx], x, y)
		fallthrough
	case isSnake(id):
		rt.ObjectState[idx] = (rt.ObjectState[idx] &^ snakeStunMask) | snakeStunDuration
		rt.emitSound(SoundEnemyHit)
		return true
	default:
		return false
	}
}

func packRedSnakeHammerTarget(state, x, y int) int {
	packed := uint32(state)
	packed = ((packed - 0x8000) & 0xff01ffff) | uint32(x<<17)
	packed = (packed & 0x80ffffff) | uint32(y<<24)
	direction := state & objectDirectionMask
	if direction == 1 || direction == 3 {
		packed |= 0x80000000
	} else {
		packed &^= 0x80000000
	}
	return int(packed)
}

func (rt *Runtime) hasFreezeHammerTargetAt(x, y int) bool {
	for _, delta := range []Point{{}, {X: -1}, {X: 1}, {Y: -1}, {Y: 1}} {
		nx, ny := x+delta.X, y+delta.Y
		id, ok := rt.At(PlayerLayer, nx, ny)
		if !ok {
			continue
		}
		if rt.freezeHammerEligibleAt(x, y, nx, ny, id) {
			return true
		}
	}
	return false
}

func (rt *Runtime) freezeHammerTargetsAt(x, y int) bool {
	frozen := false
	for _, delta := range []Point{{}, {X: -1}, {X: 1}, {Y: -1}, {Y: 1}} {
		nx, ny := x+delta.X, y+delta.Y
		id, ok := rt.At(PlayerLayer, nx, ny)
		if !ok {
			continue
		}
		if !rt.freezeHammerEligibleAt(x, y, nx, ny, id) {
			continue
		}
		if rt.freezeObjectAt(nx, ny) {
			frozen = true
		}
	}
	return frozen
}

func (rt *Runtime) freezeHammerEligibleAt(impactX, impactY, objectX, objectY int, id RawID) bool {
	if id == 1 {
		return objectX == impactX && objectY == impactY
	}
	if !isSnake(id) {
		return false
	}
	idx := rt.index(objectX, objectY)
	direction := rt.ObjectState[idx] & objectDirectionMask
	timer := 0
	if direction != 0 {
		timer = rt.ObjectMotion[idx].Remaining
	}
	dx, dy := snakeStep(direction)
	pixelX := objectX*TileSize - dx*timer
	pixelY := objectY*TileSize - dy*timer
	return absInt(impactX*TileSize-pixelX) < TileSize && absInt(impactY*TileSize-pixelY) < TileSize
}

func (rt *Runtime) freezeObjectAt(x, y int) bool {
	idx := rt.index(x, y)
	id := rt.PlayerLayer[idx]
	if id != 1 && !isSnake(id) {
		return false
	}
	rt.FrozenOriginal[idx] = id
	rt.PlayerLayer[idx] = 9
	return true
}

func (rt *Runtime) thawFrozenAt(x, y int) bool {
	idx := rt.index(x, y)
	if rt.PlayerLayer[idx] != 9 {
		return false
	}
	original := rt.FrozenOriginal[idx]
	if original != 1 && !isSnake(original) {
		return false
	}
	rt.PlayerLayer[idx] = original
	rt.FrozenOriginal[idx] = EmptyRawID
	rt.ObjectMotion[idx] = ObjectMotion{}
	if isSnake(original) {
		direction := 1
		if rt.Player.X == x && rt.Player.Y == y-1 {
			direction = 2
		}
		// The thaw branch restores raw19/raw43, then deliberately calls
		// bVoid(19,...), which applies the full snake stun without reducing
		// red-snake durability.
		rt.ObjectState[idx] = direction | snakeStunDuration
	} else {
		rt.ObjectState[idx] = 0
	}
	return true
}

func (rt *Runtime) FrozenOriginalAt(x, y int) RawID {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return EmptyRawID
	}
	return rt.FrozenOriginal[rt.index(x, y)]
}

func (rt *Runtime) ForegroundStateAt(x, y int) int {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return 0
	}
	return rt.ForegroundState[rt.index(x, y)]
}

func (rt *Runtime) UseSpecialBarrier(dx, dy int) bool {
	if !rt.CanAcceptInput() || rt.specialToolLevel() < 2 {
		return false
	}
	// Source movement allows the hero onto raw 2. Pressing 5 then clears the
	// connected state-1 cluster under the hero; direction is irrelevant.
	x := rt.Player.X
	y := rt.Player.Y
	foregroundID, ok := rt.At(ForegroundLayer, x, y)
	if !ok || foregroundID != 2 {
		return false
	}
	if rt.Background[rt.index(x, y)] != 1 {
		return false
	}
	return rt.clearForegroundBlob(x, y, 2) > 0
}

func (rt *Runtime) SpecialBarrierPrompt() (toolModule int, available, visible bool) {
	if rt == nil || rt.Player.X < 0 || rt.Player.Y < 0 || rt.Player.X >= rt.Width() || rt.Player.Y >= rt.Height() {
		return 0, false, false
	}
	idx := rt.index(rt.Player.X, rt.Player.Y)
	if rt.Foreground[idx] != 2 {
		return 0, false, false
	}
	state := int(rt.Background[idx])
	switch state {
	case 0:
		return 0, rt.specialToolLevel() >= 1, true
	case 1:
		return 1, rt.specialToolLevel() >= 2, true
	default:
		return 0, false, false
	}
}

func (rt *Runtime) ContainerLockedAt(x, y int) bool {
	if rt == nil || x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return false
	}
	return rt.ContainerLocked[rt.index(x, y)]
}

func (rt *Runtime) UseHook(dx, dy int) bool {
	if !rt.CanAcceptInput() || rt.specialToolLevel() < 2 || dx == 0 || dy != 0 {
		return false
	}
	step := 1
	if dx < 0 {
		step = -1
	}
	for distance := 1; distance <= 3; distance++ {
		x := rt.Player.X + step*distance
		y := rt.Player.Y
		if x < 0 || x >= rt.Width() {
			return false
		}
		playerID, _ := rt.At(PlayerLayer, x, y)
		foregroundID, _ := rt.At(ForegroundLayer, x, y)
		if foregroundID == 7 && !rt.foregroundDoorOpen(x, y) {
			return false
		}
		if playerID == EmptyRawID {
			continue
		}
		if distance <= 1 {
			return false
		}
		if hookSelectable(playerID) && (playerID != 48 || rt.ObjectState[rt.index(x, y)]&0x8 == 0) {
			rt.startHookAction(Point{X: x, Y: y}, step)
			return true
		}
		return false
	}
	return false
}

func (rt *Runtime) startHookAction(target Point, dx int) {
	rt.Hooking = true
	rt.HookTicks = 0
	rt.HookTarget = target
	rt.hookDX = dx
	rt.hookStepsRemaining = 0
	rt.hookCollect = false
	rt.hookReturning = false
	rt.hookOriginalState = 0
	rt.hookTip = Point{X: rt.Player.X + dx, Y: rt.Player.Y}
	if dx > 0 {
		rt.HookAnimation = hookRightCastAnimation
	} else {
		rt.HookAnimation = hookLeftCastAnimation
	}
	idx := rt.index(rt.hookTip.X, rt.hookTip.Y)
	rt.PlayerLayer[idx] = 32
	rt.ObjectState[idx] = 4
	if dx > 0 {
		rt.ObjectState[idx] |= 1
	}
	rt.ObjectMotion[idx] = ObjectMotion{Remaining: 18}
}

func hookCastDuration(animation int) int {
	if animation == hookLeftCastAnimation {
		return hookLeftCastDuration
	}
	return hookRightCastDuration
}

func hookPullDuration(animation int) int {
	if animation == hookLeftPullAnimation {
		return hookLeftPullDuration
	}
	return hookRightPullDuration
}

func (rt *Runtime) tickHookAction() {
	if !rt.Hooking {
		return
	}
	rt.HookTicks++
	if !rt.hookReturning {
		rt.tickHookExtension()
		return
	}
	rt.tickHookPull()
}

func (rt *Runtime) tickHookExtension() {
	idx := rt.index(rt.hookTip.X, rt.hookTip.Y)
	if rt.PlayerLayer[idx] != 32 {
		rt.finishHookAction()
		return
	}
	if rt.ObjectMotion[idx].Remaining > 0 {
		rt.ObjectMotion[idx].Remaining = max(0, rt.ObjectMotion[idx].Remaining-6)
		return
	}

	next := Point{X: rt.hookTip.X + rt.hookDX, Y: rt.hookTip.Y}
	if next.X < 0 || next.X >= rt.Width() {
		rt.finishHookAction()
		return
	}
	nextIdx := rt.index(next.X, next.Y)
	nextID := rt.PlayerLayer[nextIdx]
	if hookSelectable(nextID) && (nextID != 48 || rt.ObjectState[nextIdx]&0x8 == 0) {
		rt.HookTarget = next
		rt.hookOriginalState = sourceHookRestoreState(nextID, rt.ObjectState[nextIdx])
		rt.hookCollect = nextID == 1
		rt.hookStepsRemaining = absInt(next.X-rt.Player.X) - 1
		if rt.hookCollect {
			rt.hookStepsRemaining++
		}
		rt.hookReturning = true
		rt.emitSound(SoundHook)
		return
	}

	remaining := rt.ObjectState[idx] >> 1
	foregroundID := rt.Foreground[nextIdx]
	if nextID != EmptyRawID || remaining <= 0 || foregroundID == 14 || foregroundID == 33 {
		rt.finishHookAction()
		return
	}
	state := (remaining - 1) << 1
	if rt.hookDX > 0 {
		state |= 1
	}
	rt.PlayerLayer[nextIdx] = 32
	rt.ObjectState[nextIdx] = state
	timer := 18
	// The source scans columns left-to-right. A new right-facing segment is
	// visited later in the same scan and immediately loses its first 6 pixels.
	if rt.hookDX > 0 {
		timer = 12
	}
	rt.ObjectMotion[nextIdx] = ObjectMotion{Remaining: timer}
	rt.hookTip = next
}

func sourceHookRestoreState(id RawID, state int) int {
	switch id {
	case 0, 8, 9, 47:
		return state &^ (0x7000 | gravityRollPreparing)
	case 14:
		return state
	default:
		return -1
	}
}

func (rt *Runtime) tickHookPull() {
	idx := rt.index(rt.HookTarget.X, rt.HookTarget.Y)
	id := rt.PlayerLayer[idx]
	if !hookSelectable(id) {
		rt.finishHookAction()
		return
	}
	if rt.ObjectMotion[idx].Remaining > 0 {
		rt.ObjectMotion[idx].Remaining = max(0, rt.ObjectMotion[idx].Remaining-6)
		if rt.hookStepsRemaining == 0 && rt.ObjectMotion[idx].Remaining <= 6 {
			rt.startHookPullAnimation()
		}
		return
	}
	if rt.hookStepsRemaining <= 0 {
		rt.ObjectState[idx] = rt.hookOriginalState
		rt.finishHookAction()
		return
	}

	from := rt.HookTarget
	to := Point{X: from.X - rt.hookDX, Y: from.Y}
	toIdx := rt.index(to.X, to.Y)
	if rt.PlayerLayer[toIdx] != EmptyRawID && rt.PlayerLayer[toIdx] != 32 {
		rt.finishHookAction()
		return
	}
	rt.moveHookTarget(from, to, id)
	rt.HookTarget = to
	rt.hookStepsRemaining--
	if rt.hasHookRope() {
		// The remaining segment sees the newly moved target later in the same
		// object scan, re-acquires it, and resets the first pull step to zero.
		rt.ObjectMotion[toIdx].Remaining = 0
		rt.emitSound(SoundHook)
	}
	rt.updatePressureDoors()
}

func (rt *Runtime) startHookPullAnimation() {
	animation := hookRightPullAnimation
	if rt.hookDX < 0 {
		animation = hookLeftPullAnimation
	}
	if rt.HookAnimation == animation {
		return
	}
	rt.HookAnimation = animation
	rt.HookTicks = 0
}

func (rt *Runtime) moveHookTarget(from, to Point, id RawID) {
	fromIdx := rt.index(from.X, from.Y)
	toIdx := rt.index(to.X, to.Y)
	state := rt.ObjectState[fromIdx]
	group := rt.EnemyGateGroup[fromIdx]
	frozenOriginal := rt.FrozenOriginal[fromIdx]
	direction := 4
	if to.X > from.X {
		direction = 2
	}
	rt.PlayerLayer[toIdx] = id
	rt.PlayerLayer[fromIdx] = EmptyRawID
	rt.ObjectState[toIdx] = (state &^ objectDirectionMask) | direction
	rt.ObjectState[fromIdx] = 0
	rt.ObjectMotion[toIdx] = ObjectMotion{DX: to.X - from.X, Remaining: 18}
	rt.ObjectMotion[fromIdx] = ObjectMotion{}
	rt.EnemyGateGroup[toIdx] = group
	rt.EnemyGateGroup[fromIdx] = -1
	rt.FrozenOriginal[toIdx] = frozenOriginal
	rt.FrozenOriginal[fromIdx] = EmptyRawID
}

func (rt *Runtime) hasHookRope() bool {
	for distance := 1; distance <= 3; distance++ {
		x := rt.Player.X + rt.hookDX*distance
		if x < 0 || x >= rt.Width() {
			break
		}
		if rt.PlayerLayer[rt.index(x, rt.Player.Y)] == 32 {
			return true
		}
	}
	return false
}

func (rt *Runtime) clearHookRope() {
	for distance := 1; distance <= 3; distance++ {
		x := rt.Player.X + rt.hookDX*distance
		if x < 0 || x >= rt.Width() {
			break
		}
		idx := rt.index(x, rt.Player.Y)
		if rt.PlayerLayer[idx] != 32 {
			continue
		}
		rt.PlayerLayer[idx] = EmptyRawID
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
	}
}

func (rt *Runtime) finishHookAction() {
	if rt.hookDX != 0 {
		rt.clearHookRope()
	}
	rt.Hooking = false
	rt.HookTicks = 0
	rt.HookAnimation = 0
	rt.HookTarget = Point{}
	rt.hookDX = 0
	rt.hookStepsRemaining = 0
	rt.hookCollect = false
	rt.hookReturning = false
	rt.hookTip = Point{}
	rt.hookOriginalState = 0
}

func (rt *Runtime) collectVioletAt(x, y int) {
	idx := rt.index(x, y)
	rt.VioletGems++
	rt.PlayerLayer[idx] = EmptyRawID
	rt.ObjectState[idx] = 0
	rt.ObjectMotion[idx] = ObjectMotion{}
	rt.consumeBonusQuota(1)
	rt.updateExitOpen()
	rt.frameVioletPickups = append(rt.frameVioletPickups, Point{X: x, Y: y})
}

func (rt *Runtime) objectOverlapsPlayer(x, y int) bool {
	if absInt(x-rt.Player.X) > 1 || absInt(y-rt.Player.Y) > 1 {
		return false
	}
	objectMotion := rt.ObjectMotion[rt.index(x, y)]
	objectX := x*TileSize - objectMotion.DX*objectMotion.Remaining
	objectY := y*TileSize - objectMotion.DY*objectMotion.Remaining
	playerX := rt.Player.X*TileSize - rt.PlayerMotion.DX*rt.PlayerMotion.Remaining
	playerY := rt.Player.Y*TileSize - rt.PlayerMotion.DY*rt.PlayerMotion.Remaining
	return absInt(objectX-playerX) < TileSize && absInt(objectY-playerY) < TileSize
}

func (rt *Runtime) TickBreakables() int {
	broken := 0
	for y := rt.Height() - 1; y >= 0; y-- {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			id := rt.PlayerLayer[idx]
			if id != 30 && id != 37 || rt.ObjectState[idx] <= 0 {
				continue
			}
			if id == 37 {
				if rt.ObjectState[idx] >= 8 {
					rt.triggerWaterReflow(x, y)
					rt.PlayerLayer[idx] = EmptyRawID
					rt.ObjectState[idx] = 0
					rt.BreakableWalls++
					broken++
					continue
				}
				rt.ObjectState[idx]++
				continue
			}
			if rt.ObjectState[idx] == 4 {
				rt.propagateBreakableDamage(x, y)
			}
			if rt.ObjectState[idx] >= 16 {
				rt.PlayerLayer[idx] = EmptyRawID
				rt.ObjectState[idx] = 0
				rt.BreakableWalls++
				broken++
				continue
			}
			rt.ObjectState[idx]++
		}
	}
	if broken > 0 {
		rt.TickForegroundTriggers()
	}
	return broken
}

func (rt *Runtime) TickForegroundTriggers() int {
	cleared := 0
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 2 || rt.Background[idx] != 0 {
				continue
			}
			if rt.adjacentPlayerRaw(x, y, 30) {
				continue
			}
			cleared += rt.clearForegroundBlob(x, y, 2)
		}
	}
	return cleared
}

func (rt *Runtime) TickSnakes() int {
	return rt.tickSnakesBounds(0, rt.Width()-1, 0, rt.Height()-1)
}

func (rt *Runtime) TickSnakesNearPlayer(radius int) int {
	return rt.tickSnakesBounds(
		max(1, rt.Player.X-radius),
		min(rt.Width()-2, rt.Player.X+radius),
		max(1, rt.Player.Y-radius),
		min(rt.Height()-2, rt.Player.Y+radius),
	)
}

func (rt *Runtime) tickSnakesBounds(minX, maxX, minY, maxY int) int {
	moved := 0
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			if rt.tickSnakeObjectAt(x, y) {
				moved++
			}
		}
	}
	return moved
}

func (rt *Runtime) tickSnakeObjectAt(x, y int) bool {
	idx := rt.index(x, y)
	id := rt.PlayerLayer[idx]
	if !isSnake(id) {
		return false
	}
	if rt.snakeCrushedAt(x, y) {
		rt.decrementEnemyGateForObjectAt(x, y)
		rt.PlayerLayer[idx] = EmptyRawID
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
		return false
	}
	state := rt.ObjectState[idx]
	stunned := state&snakeStunMask != 0
	if rt.advanceObjectMotion(idx, 3, 0) {
		if !stunned {
			if rt.hurtFromSnakeAt(x, y, state&0x7) {
				rt.clearRedSnakeChase(idx, id)
			}
		}
		return false
	}
	if stunned {
		if rt.gravitySourceTick&3 == 0 {
			remaining := max(0, (state&snakeStunMask)-8)
			state = (state &^ snakeStunMask) | remaining
			if id == 43 && remaining == 0 {
				state = (state &^ 0xf00) | 0xc00
			}
			rt.ObjectState[idx] = state
		}
		return false
	}
	if id == 43 && state&0xf00 != 0 {
		dir := rt.snakeDirectionToward(x, y, rt.Player.X, rt.Player.Y)
		state = (state &^ objectDirectionMask) | dir
		state -= 0x100
		rt.ObjectState[idx] = state
	} else if uint32(state)&0x00fe0000 != 0 {
		targetX := int(uint32(state)&0x00fe0000) >> 17
		targetY := int(uint32(state)&0x7f000000) >> 24
		if x == targetX && y == targetY {
			dir := 2
			if uint32(state)&0x80000000 != 0 {
				dir = 1
			}
			state = (state & 0xff01ffff &^ objectDirectionMask) | dir
			rt.ObjectState[idx] = state
		} else {
			dir := rt.snakeDirectionToward(x, y, targetX, targetY)
			state = (state &^ objectDirectionMask) | dir
			rt.ObjectState[idx] = state
		}
	}
	dir := state & 0x7
	usedPendingDirection := dir == 0
	if dir == 0 {
		dir = (state & 0x7000) >> 12
		if dir == 0 {
			if rt.hurtFromSnakeAt(x, y, 0) {
				rt.clearRedSnakeChase(idx, id)
			}
			return false
		}
		state = (state &^ 0x7) | dir
		rt.ObjectState[idx] = state
	}
	dx, dy := snakeStep(dir)
	targetX := x + dx
	targetY := y + dy
	if !rt.cellEmptyForSnake(targetX, targetY) {
		if usedPendingDirection {
			// A pending turn gets one target probe before the regular blocked
			// branch schedules another turn cycle.
			if rt.hurtFromSnakeAt(x, y, dir) {
				rt.clearRedSnakeChase(idx, id)
			}
			return false
		}
		reverse := reverseSnakeDirection(dir)
		rt.ObjectState[idx] = (state &^ (0x7 | 0x7000)) | reverse<<12
		rt.ObjectMotion[idx] = ObjectMotion{Remaining: 21}
		if rt.hurtFromSnakeAt(x, y, dir) {
			rt.clearRedSnakeChase(idx, id)
		}
		return false
	}
	targetIdx := rt.index(targetX, targetY)
	rt.PlayerLayer[targetIdx] = id
	rt.ObjectState[targetIdx] = (state &^ 0x7) | dir
	rt.ObjectMotion[targetIdx] = ObjectMotion{DX: dx, DY: dy, Remaining: 21}
	rt.PlayerLayer[idx] = EmptyRawID
	rt.ObjectState[idx] = 0
	rt.ObjectMotion[idx] = ObjectMotion{}
	rt.transferEnemyGateGroup(idx, targetIdx)
	if rt.hurtFromSnakeAt(targetX, targetY, dir) {
		rt.clearRedSnakeChase(targetIdx, id)
	}
	return true
}

func (rt *Runtime) snakeDirectionToward(x, y, targetX, targetY int) int {
	dx := targetX - x
	dy := targetY - y
	dir := 0
	if absInt(dx) > absInt(dy) {
		if dx < 0 {
			dir = 4
		} else if dx > 0 {
			dir = 2
		}
		if dir != 0 {
			stepX, stepY := snakeStep(dir)
			if !rt.cellEmptyForSnake(x+stepX, y+stepY) {
				dir = 0
			}
		}
	}
	if dir != 0 {
		return dir
	}
	if dy < 0 {
		dir = 1
	} else if dy > 0 {
		dir = 3
	}
	if dir == 0 {
		return 0
	}
	stepX, stepY := snakeStep(dir)
	if rt.cellEmptyForSnake(x+stepX, y+stepY) && rt.waterCellAt(x+stepX, y+stepY) == 0 {
		return dir
	}
	if dx < 0 {
		dir = 4
	} else if dx > 0 {
		dir = 2
	} else {
		return 0
	}
	stepX, stepY = snakeStep(dir)
	playerID, ok := rt.At(PlayerLayer, x+stepX, y+stepY)
	if !ok || playerID != EmptyRawID {
		return 0
	}
	return dir
}

func (rt *Runtime) clearRedSnakeChase(idx int, id RawID) {
	if id == 43 && idx >= 0 && idx < len(rt.ObjectState) {
		rt.ObjectState[idx] &^= 0xf00
	}
}

func (rt *Runtime) snakeCrushedAt(x, y int) bool {
	if rt.Foreground[rt.index(x, y)] == 35 {
		return false
	}
	if rt.movingHazardStrikesAt(x, y) {
		return true
	}
	if y <= 0 {
		return false
	}
	aboveIdx := rt.index(x, y-1)
	aboveID := rt.PlayerLayer[aboveIdx]
	if aboveID == 14 && rt.ObjectMotion[aboveIdx].Remaining <= 6 {
		return true
	}
	if !isGravityObject(aboveID) || rt.ObjectMotion[aboveIdx].Remaining > 6 {
		return false
	}
	falling := rt.ObjectMotion[aboveIdx].DY == 1
	return falling || aboveID != 1
}

func (rt *Runtime) hurtFromSnakeAt(x, y, direction int) bool {
	if rt.Hammering {
		return false
	}
	if !rt.objectOverlapsPlayer(x, y) {
		return false
	}
	return rt.HurtFromDirection(1, direction)
}

func (rt *Runtime) TickHorizontalHazards() int {
	hits := rt.tickHorizontalHazardsBounds(0, rt.Width()-1, 0, rt.Height()-1, -1)
	return hits
}

func (rt *Runtime) TickHorizontalHazardsNearPlayer(radius, reach int) int {
	return rt.tickHorizontalHazardsBounds(
		max(1, rt.Player.X-radius),
		min(rt.Width()-2, rt.Player.X+radius),
		max(1, rt.Player.Y-radius),
		min(rt.Height()-2, rt.Player.Y+radius),
		reach,
	)
}

func (rt *Runtime) tickHorizontalHazardsBounds(minX, maxX, minY, maxY, fixedReach int) int {
	hits := 0
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			idx := rt.index(x, y)
			id := rt.PlayerLayer[idx]
			if id != 22 && id != 23 {
				continue
			}
			reach := fixedReach
			if reach < 0 {
				reach = rt.ObjectState[idx] & 0x3
			}
			if rt.horizontalHazardHitsPlayer(x, y, id, reach) && rt.Hurt(1) {
				hits++
			}
			if fixedReach < 0 {
				rt.ObjectState[idx] = (reach + 1) & 0x3
			}
		}
	}
	return hits
}

func (rt *Runtime) horizontalHazardHitsPlayer(x, y int, id RawID, reach int) bool {
	if rt.Player.Y != y {
		return false
	}
	dx := 1
	if id == 23 {
		dx = -1
	}
	for step := 0; step <= reach; step++ {
		if rt.Player.X == x+step*dx {
			return true
		}
	}
	return false
}

func (rt *Runtime) TickCrawlers() int {
	return rt.tickCrawlersBounds(0, rt.Width()-1, 0, rt.Height()-1)
}

func (rt *Runtime) TickCrawlersNearPlayer(radius int) int {
	return rt.tickCrawlersBounds(
		max(1, rt.Player.X-radius),
		min(rt.Width()-2, rt.Player.X+radius),
		max(1, rt.Player.Y-radius),
		min(rt.Height()-2, rt.Player.Y+radius),
	)
}

func (rt *Runtime) tickCrawlersBounds(minX, maxX, minY, maxY int) int {
	moved := 0
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			if rt.tickCrawlerObjectAt(x, y) {
				moved++
			}
		}
	}
	return moved
}

func (rt *Runtime) tickCrawlerObjectAt(x, y int) bool {
	idx := rt.index(x, y)
	if rt.PlayerLayer[idx] != 11 {
		return false
	}
	state := rt.ObjectState[idx]
	deathPhase := (state & 0xf00) >> 8
	moved := false
	if deathPhase != 0 {
		if deathPhase >= 4 {
			rt.PlayerLayer[idx] = EmptyRawID
			rt.ObjectState[idx] = 0
			rt.ObjectMotion[idx] = ObjectMotion{}
		} else if (rt.gravitySourceTick>>1)&1 == 0 {
			rt.ObjectState[idx] += 0x100
		}
	} else if rt.waterCellAt(x, y) != 0 {
		rt.ObjectState[idx] = (state &^ 0xf00) | 0x100
	} else if rt.ObjectMotion[idx].Remaining <= 4 {
		dir := state & objectDirectionMask
		if dir == 0 {
			dir = rt.initialCrawlerDirection(x, y, state&0x10 != 0)
			rt.ObjectState[idx] = (state &^ objectDirectionMask) | dir
		} else {
			moved = rt.advanceCrawlerAt(x, y, dir, state)
		}
	}

	// amVoid() tests contact at the crawler's source cell. A crawler moving
	// right or up is visited again later in the same source scan, which is
	// what makes those directions hit immediately and start at timer 13.
	if rt.isPlayerAt(x, y) {
		rt.Hurt(1)
	}
	if !moved && rt.PlayerLayer[idx] == 11 && rt.ObjectMotion[idx].Remaining > 0 {
		rt.ObjectMotion[idx].Remaining -= 5
	}
	return moved
}

func (rt *Runtime) advanceCrawlerAt(x, y, dir, state int) bool {
	forwardX, forwardY := snakeStep(dir)
	sideX, sideY := -forwardY, forwardX
	if state&0x10 != 0 {
		sideX, sideY = forwardY, -forwardX
	}
	forwardEmpty := rt.cellEmptyForSnake(x+forwardX, y+forwardY)
	sideEmpty := rt.cellEmptyForSnake(x+sideX, y+sideY)
	behindSideEmpty := rt.cellEmptyForSnake(x+sideX-forwardX, y+sideY-forwardY)

	if forwardEmpty && sideEmpty && behindSideEmpty {
		if rt.ObjectMotion[rt.index(x, y)].Remaining <= 0 {
			return rt.moveCrawler(x, y, x+forwardX, y+forwardY, state)
		}
		return false
	}
	if sideEmpty {
		turnDir := directionForStep(sideX, sideY)
		state = (state &^ objectDirectionMask) | turnDir
		return rt.moveCrawler(x, y, x+sideX, y+sideY, state)
	}
	if forwardEmpty {
		if rt.ObjectMotion[rt.index(x, y)].Remaining <= 0 {
			return rt.moveCrawler(x, y, x+forwardX, y+forwardY, state)
		}
		return false
	}
	turnDir := directionForStep(-sideX, -sideY)
	rt.ObjectState[rt.index(x, y)] = (state &^ objectDirectionMask) | turnDir
	return false
}

func (rt *Runtime) moveCrawler(fromX, fromY, toX, toY int, state int) bool {
	fromIdx := rt.index(fromX, fromY)
	targetIdx := rt.index(toX, toY)
	rt.PlayerLayer[targetIdx] = 11
	rt.ObjectState[targetIdx] = state
	rt.ObjectMotion[targetIdx] = ObjectMotion{DX: toX - fromX, DY: toY - fromY, Remaining: 18}
	rt.PlayerLayer[fromIdx] = EmptyRawID
	rt.ObjectState[fromIdx] = 0
	rt.ObjectMotion[fromIdx] = ObjectMotion{}
	return true
}

func (rt *Runtime) initialCrawlerDirection(x, y int, reversed bool) int {
	dir := 0
	leftOccupied := rt.playerCellOccupied(x-1, y)
	rightOccupied := rt.playerCellOccupied(x+1, y)
	belowOccupied := rt.playerCellOccupied(x, y+1)
	if leftOccupied {
		dir = 3
		if reversed {
			dir = 1
		}
	} else if belowOccupied {
		dir = 4
		if reversed {
			dir = 2
		}
	}
	if rightOccupied {
		dir = 1
		if reversed {
			dir = 3
		}
	}
	if !belowOccupied {
		return dir
	}
	if reversed {
		return 4
	}
	return 2
}

func (rt *Runtime) playerCellOccupied(x, y int) bool {
	playerID, ok := rt.At(PlayerLayer, x, y)
	return !ok || playerID != EmptyRawID
}

func directionForStep(dx, dy int) int {
	switch {
	case dy < 0:
		return 1
	case dx > 0:
		return 2
	case dy > 0:
		return 3
	case dx < 0:
		return 4
	default:
		return 0
	}
}

func (rt *Runtime) propagateBreakableDamage(x, y int) {
	for _, delta := range []Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}} {
		nx := x + delta.X
		ny := y + delta.Y
		if id, ok := rt.At(PlayerLayer, nx, ny); ok && id == 30 {
			idx := rt.index(nx, ny)
			rt.ObjectState[idx]++
		}
	}
}

func (rt *Runtime) CanExit() bool {
	return len(rt.GoalMarkers) > 0
}

func (rt *Runtime) updateExitOpen() {
	rt.ExitOpen = rt.CanExit()
}

func (rt *Runtime) IsPassable(x, y int) bool {
	if rt.spikeOccupies(x, y) {
		return false
	}
	if rt.WaterAt(x, y) > 0 && rt.SpecialItemMask&4 == 0 {
		return false
	}
	playerID, ok := rt.At(PlayerLayer, x, y)
	if !ok {
		return false
	}
	if playerID == 10 && !rt.waterStable() {
		return false
	}
	if !playerLayerPassable(playerID) {
		return false
	}
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	return rt.foregroundPassable(x, y, foregroundID)
}

func playerLayerPassable(playerID RawID) bool {
	switch {
	case playerID == EmptyRawID:
		return true
	case playerID == 1:
		return true
	case playerID == 2:
		return true
	case playerID == 4:
		return true
	case playerID == 5:
		return true
	case playerID == 6:
		return true
	case playerID == 7:
		return true
	case playerID == 10:
		return true
	case playerID == 14:
		return true
	case playerID == 24:
		return true
	case playerID == 26:
		return true
	case playerID == 27:
		return true
	case playerID == 33:
		return true
	case playerID == 41:
		return true
	case playerID == 42:
		return true
	case playerID == 40:
		return true
	case playerID == 51:
		return true
	case playerID == 52:
		return true
	case playerID == 53:
		return true
	case playerID == 50:
		return true
	case isContactEnemy(playerID):
		return true
	case playerID == 31:
		return false
	case playerID == EntranceRawID:
		return true
	case playerID >= 80:
		return false
	default:
		return false
	}
}

func (rt *Runtime) foregroundPassable(x, y int, id RawID) bool {
	switch id {
	case 2:
		return true
	case 7:
		return rt.foregroundDoorOpen(x, y)
	case 8, 9:
		return false
	default:
		return true
	}
}

func (rt *Runtime) foregroundDoorOpen(x, y int) bool {
	state, ok := rt.At(BackgroundLayer, x, y)
	if !ok {
		return false
	}
	return state == EmptyRawID || (int(state)&0xF0)>>4 >= 2
}

func (rt *Runtime) decrementEnemyGateForObjectAt(x, y int) {
	idx := rt.index(x, y)
	if idx >= 0 && idx < len(rt.EnemyGateGroup) {
		rt.EnemyGateGroup[idx] = -1
	}
	// Java alVoid() does not consume the arena counter while a special-stage
	// boss still has health. The boss death sequence performs the one final
	// decrement that opens the exit.
	if rt.Anaconda.Enabled && rt.Anaconda.Health > 0 || rt.TeutonicKnight.Enabled && rt.TeutonicKnight.Health > 0 {
		return
	}
	group := rt.ActiveEnemyGateGroup
	if group < 0 {
		return
	}
	if rt.EnemyGateCounters[group] <= 0 {
		return
	}
	rt.EnemyGateCounters[group]--
	if rt.EnemyGateCounters[group] == 0 {
		rt.openEnemyGateGroup(group)
	}
}

func (rt *Runtime) openEnemyGateGroup(group int) {
	rt.emitSound(SoundDoor)
	for y := 1; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 17 || int(rt.Background[idx]) != group {
				continue
			}
			aboveIdx := rt.index(x, y-1)
			switch rt.Foreground[aboveIdx] {
			case 7:
				if !rt.foregroundDoorOpen(x, y-1) {
					rt.openDoorAt(x, y-1)
				}
			case 14, 33:
				rt.ContainerLocked[aboveIdx] = false
			}
		}
	}
}

func (rt *Runtime) enemyGateGroupValid(group int) bool {
	if group < 0 || group == int(EmptyRawID) {
		return false
	}
	if _, ok := rt.EnemyGateCounters[group]; ok {
		return true
	}
	for idx, id := range rt.Foreground {
		if id == 17 && int(rt.Background[idx]) == group {
			return true
		}
	}
	return false
}

func (rt *Runtime) activateEnemyGateTriggerAt(x, y int) {
	group := int(rt.Background[rt.index(x, y)])
	validGroup := rt.enemyGateGroupValid(group)
	if !validGroup {
		rt.ActiveEnemyGateGroup = -1
	} else {
		rt.ActiveEnemyGateGroup = group
	}
	rt.startEnemyGateDemo(x, y, group)
	if validGroup {
		rt.emitSound(SoundRiddle)
	}
	rt.set(rt.Foreground, x, y, EmptyRawID)
}

func (rt *Runtime) transferEnemyGateGroup(fromIdx, toIdx int) {
	if len(rt.EnemyGateGroup) == 0 {
		return
	}
	rt.EnemyGateGroup[toIdx] = rt.EnemyGateGroup[fromIdx]
	rt.EnemyGateGroup[fromIdx] = -1
}

func (rt *Runtime) updatePressureDoors() {
	seen := map[int]bool{}
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 6 {
				continue
			}
			doorID := int(rt.Background[idx])
			if doorID == int(EmptyRawID) {
				continue
			}
			doorID &= 0x0F
			if rt.pressureSwitchActive(x, y) {
				if rt.activateDoorByID(doorID) {
					rt.emitSound(SoundDoor)
				}
				seen[doorID] = true
			} else if !seen[doorID] {
				rt.closeDoorByID(doorID)
			}
		}
	}
}

func (rt *Runtime) pressureSwitchActive(x, y int) bool {
	if rt.isPlayerAt(x, y) {
		return rt.PlayerMotion.Remaining < 12
	}
	id, ok := rt.At(PlayerLayer, x, y)
	if !ok {
		return false
	}
	switch id {
	case 0, 1, 8, 9, 47, 48:
		return rt.ObjectMotion[rt.index(x, y)].Remaining < 12
	default:
		return false
	}
}

func (rt *Runtime) activateDoorByID(doorID int) bool {
	activated := false
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 7 || rt.doorGroupAt(idx) != doorID || int(rt.Background[idx])&0xf0 != 0 {
				continue
			}
			remaining := int(rt.Background[idx]) & 0x0f
			if remaining <= 0 {
				continue
			}
			remaining--
			if remaining == 0 {
				// hVoid() starts phase 1 while retaining the low nibble from
				// before the decrement. Closing later therefore restores 1.
				rt.Background[idx] = RawID((int(rt.Background[idx]) & 0x0f) | 0x10)
			} else {
				rt.Background[idx] = RawID((int(rt.Background[idx]) & 0xf0) | remaining)
			}
			activated = true
		}
	}
	return activated
}

func (rt *Runtime) closeDoorByID(doorID int) {
	closed := false
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 7 || rt.doorGroupAt(idx) != doorID || int(rt.Background[idx])&0xf0 == 0 || rt.PlayerLayer[idx] == 32 {
				continue
			}
			rt.Background[idx] &= 0x0f
			if !closed {
				rt.emitSound(SoundBoulder)
				closed = true
			}
			rt.resolveClosedDoorOccupant(x, y)
		}
	}
}

func (rt *Runtime) resolveClosedDoorOccupant(x, y int) {
	if rt.isPlayerAt(x, y) {
		rt.HurtFromDirection(rt.MaxHealth, 0)
		rt.emitSound(SoundDeath)
		return
	}
	idx := rt.index(x, y)
	id := rt.PlayerLayer[idx]
	switch id {
	case 0, 1, 19, 43, 45:
		// doorHeadClose() calls jVoid() for every listed object, including
		// boulders and violet gems; jVoid() decrements the active cmInt group.
		rt.decrementEnemyGateForObjectAt(x, y)
		rt.PlayerLayer[idx] = EmptyRawID
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
		rt.FrozenOriginal[idx] = EmptyRawID
	}
}

func (rt *Runtime) openDoorAt(x, y int) bool {
	idx := rt.index(x, y)
	if rt.foregroundDoorOpen(x, y) || int(rt.Background[idx])&0xf0 != 0 {
		return false
	}
	rt.Background[idx] = RawID(0x10 | int(rt.Background[idx])&0x0f)
	return true
}

func (rt *Runtime) doorGroupAt(idx int) int {
	if idx < 0 || idx >= len(rt.DoorGroup) {
		return -1
	}
	return rt.DoorGroup[idx]
}

func (rt *Runtime) adjacentPlayerRaw(x, y int, id RawID) bool {
	for _, delta := range []Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}} {
		if playerID, ok := rt.At(PlayerLayer, x+delta.X, y+delta.Y); ok && playerID == id {
			return true
		}
	}
	return false
}

func (rt *Runtime) clearForegroundBlob(x, y int, id RawID) int {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return 0
	}
	idx := rt.index(x, y)
	if rt.Foreground[idx] != id {
		return 0
	}
	rt.Foreground[idx] = EmptyRawID
	cleared := 1
	cleared += rt.clearForegroundBlob(x-1, y, id)
	cleared += rt.clearForegroundBlob(x+1, y, id)
	cleared += rt.clearForegroundBlob(x, y-1, id)
	cleared += rt.clearForegroundBlob(x, y+1, id)
	return cleared
}

func (rt *Runtime) cellEmptyForObject(x, y int) bool {
	playerID, ok := rt.At(PlayerLayer, x, y)
	if !ok || playerID != EmptyRawID {
		return false
	}
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	if foregroundID == 7 {
		return rt.foregroundDoorOpen(x, y)
	}
	return foregroundAllowsFallingObject(foregroundID)
}

func (rt *Runtime) cellEmptyForHookTarget(x, y int) bool {
	playerID, ok := rt.At(PlayerLayer, x, y)
	if !ok || playerID != EmptyRawID {
		return false
	}
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	return rt.foregroundPassable(x, y, foregroundID)
}

func (rt *Runtime) cellEmptyForSnake(x, y int) bool {
	playerID, ok := rt.At(PlayerLayer, x, y)
	if !ok || playerID != EmptyRawID {
		return false
	}
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	switch foregroundID {
	case 4, 14, 32, 33:
		return false
	case 7:
		return rt.foregroundDoorOpen(x, y)
	default:
		return true
	}
}

func (rt *Runtime) canRollBoulder(x, y, dx int) bool {
	return !rt.isPlayerAt(x+dx, y) &&
		!rt.isPlayerAt(x+dx, y+1) &&
		rt.cellEmptyForGravity(x+dx, y) &&
		rt.cellEmptyForGravity(x+dx, y+1)
}

func (rt *Runtime) tryMoveBoulder(fromX, fromY, toX, toY int) bool {
	return rt.tryMoveGravityObject(fromX, fromY, toX, toY, 0)
}

func (rt *Runtime) tryMoveGravityObject(fromX, fromY, toX, toY int, id RawID) bool {
	if rt.isPlayerAt(toX, toY) {
		motion := rt.ObjectMotion[rt.index(fromX, fromY)]
		fallingStraightDown := motion.DX == 0 && motion.DY == 1
		if id == 1 && fallingStraightDown {
			rt.VioletGems++
			rt.consumeBonusQuota(1)
			rt.set(rt.PlayerLayer, fromX, fromY, EmptyRawID)
			rt.ObjectMotion[rt.index(fromX, fromY)] = ObjectMotion{}
			return true
		}
		if (id == 0 || id == 9) && fallingStraightDown {
			rt.Hurt(2)
		}
		return false
	}
	if playerID, ok := rt.At(PlayerLayer, toX, toY); ok && isContactEnemy(playerID) && (id == 0 || id == 9) {
		rt.decrementEnemyGateForObjectAt(toX, toY)
		rt.moveObject(fromX, fromY, toX, toY, id)
		return true
	}
	if !rt.cellEmptyForGravity(toX, toY) {
		return false
	}
	rt.moveObject(fromX, fromY, toX, toY, id)
	return true
}

func (rt *Runtime) cellEmptyForGravity(x, y int) bool {
	playerID, ok := rt.At(PlayerLayer, x, y)
	if !ok || playerID != EmptyRawID {
		return false
	}
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	return foregroundAllowsFallingObject(foregroundID)
}

func foregroundAllowsFallingObject(id RawID) bool {
	switch id {
	case 5, 14, 28, 33:
		return false
	default:
		return true
	}
}

func (rt *Runtime) moveObject(fromX, fromY, toX, toY int, id RawID) {
	rt.moveObjectWithMotion(fromX, fromY, toX, toY, id, ObjectMotion{
		DX:        toX - fromX,
		DY:        toY - fromY,
		Remaining: 18,
	})
}

func (rt *Runtime) moveObjectWithMotion(fromX, fromY, toX, toY int, id RawID, motion ObjectMotion) {
	fromIdx := rt.index(fromX, fromY)
	toIdx := rt.index(toX, toY)
	state := rt.ObjectState[fromIdx]
	group := rt.EnemyGateGroup[fromIdx]
	frozenOriginal := rt.FrozenOriginal[fromIdx]
	if isGravityObject(id) {
		switch {
		case toY > fromY:
			fallDistance := min(0x1f, ((state&explosiveFallMask)>>explosiveFallShift)+1)
			state = (state &^ explosiveFallMask) | fallDistance<<explosiveFallShift
			state = (state &^ objectDirectionMask) | 3
		case toY < fromY:
			state = (state &^ (objectDirectionMask | gravityRollPreparing | gravityMoveLeft | gravityMoveRight | explosiveFallMask)) | 1
		case toX > fromX:
			state = (state &^ (objectDirectionMask | gravityRollPreparing | gravityMoveLeft | gravityMoveRight | explosiveFallMask)) | 2 | gravityMoveRight
		case toX < fromX:
			state = (state &^ (objectDirectionMask | gravityRollPreparing | gravityMoveLeft | gravityMoveRight | explosiveFallMask)) | 4 | gravityMoveLeft
		}
	} else if id != 48 {
		state = 0
	}
	rt.set(rt.PlayerLayer, toX, toY, id)
	rt.set(rt.PlayerLayer, fromX, fromY, EmptyRawID)
	rt.ObjectState[toIdx] = state
	rt.ObjectState[fromIdx] = 0
	rt.ObjectMotion[toIdx] = motion
	rt.ObjectMotion[fromIdx] = ObjectMotion{}
	if id == 9 {
		rt.EnemyGateGroup[toIdx] = group
		rt.FrozenOriginal[toIdx] = frozenOriginal
	} else {
		rt.EnemyGateGroup[toIdx] = -1
		rt.FrozenOriginal[toIdx] = EmptyRawID
	}
	rt.EnemyGateGroup[fromIdx] = -1
	rt.FrozenOriginal[fromIdx] = EmptyRawID
}

func (rt *Runtime) advanceBoulderRotation(idx int) {
	state := rt.ObjectState[idx]
	rotation := (state & boulderRotationMask) >> 3
	switch state & (gravityMoveLeft | gravityMoveRight) {
	case gravityMoveLeft:
		rotation = (rotation + 7) & 0x7
	case gravityMoveRight:
		rotation = (rotation + 1) & 0x7
	default:
		return
	}
	rt.ObjectState[idx] = (state &^ boulderRotationMask) | rotation<<3
}

func (rt *Runtime) advanceObjectMotion(idx, step, processAt int) bool {
	motion := &rt.ObjectMotion[idx]
	if motion.Remaining <= processAt {
		motion.Remaining = 0
		return false
	}
	motion.Remaining -= step
	if motion.Remaining < 0 {
		motion.Remaining = 0
	}
	return true
}

func (rt *Runtime) Hurt(amount int) bool {
	if amount <= 0 || rt.PlayerDead || rt.InvulnerabilityTicks > 0 || rt.Hooking {
		return false
	}
	rt.HurtTicks = hurtStateDuration
	rt.playerTurnOffset = 0
	rt.InvulnerabilityTicks = hurtInvulnerabilityDuration
	rt.DamageTaken += amount
	rt.HitCount++
	rt.Health -= amount
	rt.emitSound(SoundHeroHurt)
	if rt.Health <= 0 {
		rt.Health = 0
		rt.PlayerDead = true
		rt.DeathTicks = deathDuration
	}
	return true
}

// HurtFromDirection mirrors hurtHero(..., direction): after accepting damage,
// the source searches clockwise for the first neighboring cell whose player
// and foreground bytes are both negative, then applies a fresh jInt=18 move.
func (rt *Runtime) HurtFromDirection(amount, direction int) bool {
	if !rt.Hurt(amount) {
		return false
	}
	if direction < 1 || direction > 4 {
		return true
	}
	for candidate := direction; ; candidate = candidate%4 + 1 {
		dx, dy := snakeStep(candidate)
		x := rt.Player.X + dx
		y := rt.Player.Y + dy
		playerID, playerOK := rt.At(PlayerLayer, x, y)
		foregroundID, foregroundOK := rt.At(ForegroundLayer, x, y)
		if playerOK && foregroundOK && playerID.Signed() < 0 && foregroundID.Signed() < 0 {
			rt.Player = Point{X: x, Y: y}
			rt.PlayerMotion = ObjectMotion{DX: dx, DY: dy, Remaining: playerMoveStartOffset}
			rt.playerFacingDirection = candidate
			rt.ResetPushAttempt()
			break
		}
		if candidate%4+1 == direction {
			break
		}
	}
	return true
}

func (rt *Runtime) TickStatus() {
	hurtWasActive := rt.HurtTicks > 0
	if rt.EnemyGateMessageTicks > 0 {
		rt.EnemyGateMessageTicks--
		if rt.EnemyGateMessageTicks == 0 {
			rt.EnemyGateMessageIndex = 0
		}
	}
	rt.tickHammerAction()
	rt.tickHookAction()
	rt.tickChestOpening()
	rt.tickLockOpening()
	if rt.RelicCelebrating {
		rt.RelicCelebrationTicks++
		if rt.RelicCelebrationTicks > 42 {
			rt.RelicCelebrating = false
		}
	}
	if rt.RecallPending {
		rt.RecallTicks++
		if rt.RecallTicks >= recallAnimationDuration {
			rt.finishRecall()
		}
	}
	if rt.HurtTicks > 0 {
		rt.HurtTicks--
	}
	if hurtWasActive && rt.HurtTicks == 0 && rt.PlayerDead {
		rt.emitSound(SoundDeath)
	}
	if rt.InvulnerabilityTicks > 0 {
		rt.InvulnerabilityTicks--
	}
	if rt.PlayerDead && rt.DeathTicks > 0 {
		rt.DeathTicks--
		if rt.DeathTicks == 0 {
			rt.finishDeath()
		}
	}
}

func (rt *Runtime) startAdjacentLockOpening() bool {
	if rt.LockOpening || !rt.CanAcceptInput() || rt.playerSourceOffset() > 6 || rt.Player.Y < 0 || rt.Player.Y >= rt.Height() {
		return false
	}
	for x := rt.Player.X - 1; x <= rt.Player.X+1; x += 2 {
		if x < 0 || x >= rt.Width() {
			continue
		}
		idx := rt.index(x, rt.Player.Y)
		foregroundID := rt.Foreground[idx]
		if foregroundID != 8 && foregroundID != 9 || rt.ObjectState[idx] != 0 {
			continue
		}
		if foregroundID == 8 && rt.KeyForForeground8 <= 0 || foregroundID == 9 && rt.KeyForForeground9 <= 0 {
			continue
		}
		rt.LockOpening = true
		rt.LockTicks = 0
		rt.LockPoint = Point{X: x, Y: rt.Player.Y}
		rt.LockForegroundID = foregroundID
		rt.LockRewarded = false
		if rt.Player.X < x {
			rt.LockAnimation = 18
		} else {
			rt.LockAnimation = 17
		}
		rt.ResetPushAttempt()
		return true
	}
	return false
}

func (rt *Runtime) tickLockOpening() {
	if !rt.LockOpening {
		return
	}
	rt.LockTicks++
	if !rt.LockRewarded && rt.LockTicks >= lockRewardTick {
		idx := rt.index(rt.LockPoint.X, rt.LockPoint.Y)
		if rt.Foreground[idx] == rt.LockForegroundID && rt.ObjectState[idx] == 0 {
			switch rt.LockForegroundID {
			case 8:
				if rt.KeyForForeground8 > 0 {
					rt.KeyForForeground8--
				}
			case 9:
				if rt.KeyForForeground9 > 0 {
					rt.KeyForForeground9--
				}
			}
			rt.ObjectState[idx] = 1
			rt.activateDoorByID(int(rt.Background[idx]) & 0x0f)
			rt.LocksOpened++
			rt.emitSound(SoundDoor)
		}
		rt.LockRewarded = true
	}
	if rt.LockTicks >= lockOpenDuration {
		rt.LockOpening = false
		rt.LockTicks = 0
		rt.LockAnimation = 0
		rt.LockForegroundID = EmptyRawID
		rt.LockRewarded = false
	}
}

func (rt *Runtime) tickDoorAnimations(sourceTick int) {
	if sourceTick%3 != 0 {
		return
	}
	for idx, foregroundID := range rt.Foreground {
		if foregroundID != 7 || rt.Background[idx] == EmptyRawID {
			continue
		}
		state := int(rt.Background[idx])
		phase := (state & 0xf0) >> 4
		if phase < 1 || phase >= 3 {
			continue
		}
		rt.Background[idx] = RawID((state & 0x0f) | (phase+1)<<4)
	}
}

func (rt *Runtime) tickChestOpening() {
	if !rt.ChestOpening {
		return
	}
	if rt.chestOpeningFresh {
		rt.chestOpeningFresh = false
		return
	}
	rt.ChestTicks++
	rewardTick, soundTick, duration := chestAnimationTiming(rt.ChestAnimation)
	if rt.ChestTicks == soundTick {
		rt.emitSound(SoundChestReward)
	}
	if !rt.ChestRewarded && rt.ChestTicks >= rewardTick {
		rt.applyChestReward()
		rt.ChestRewarded = true
		if isRelicReward(rt.ChestRewardID) {
			rt.lastPickupTick = rt.gravitySourceTick
			rt.lastPickupTickSet = true
			rt.ChestOpening = false
			rt.ChestTicks = 0
			rt.chestOpeningFresh = false
			return
		}
	}
	if rt.ChestTicks >= duration {
		rt.lastPickupTick = rt.gravitySourceTick
		rt.lastPickupTickSet = true
		rt.ChestOpening = false
		rt.ChestTicks = 0
		rt.chestOpeningFresh = false
	}
}

func chestAnimationTiming(animation int) (rewardTick, soundTick, duration int) {
	if animation == 48 {
		return chestShortRewardTick, chestShortRewardTick, chestShortOpenDuration
	}
	return chestRewardTick, chestRewardSoundTick, chestOpenDuration
}

func (rt *Runtime) applyChestReward() {
	switch rt.ChestRewardID {
	case 2:
		rt.RedDiamonds++
		rt.markPersistentRewardAt(rt.Player)
	case 4:
		rt.KeyForForeground9++
	case 5:
		rt.KeyForForeground8++
	case 6:
		if rt.ExtraLives < 99 {
			rt.ExtraLives++
			rt.markPersistentRewardAt(rt.Player)
			rt.persistConsumedExtraLife()
			return
		}
		rt.applyHealthRefillReward()
	case 7:
		rt.applyHealthRefillReward()
	case 24:
		rt.SpecialItemMask |= 1
		rt.SpecialPickups++
		rt.queueTutorialScript(22)
	case 26:
		rt.SpecialItemMask |= 8
		rt.SpecialPickups++
	case 27:
		rt.SpecialItemMask |= 2
		rt.SpecialPickups++
	case 40:
		rt.SpecialItemMask |= 4
		rt.SpecialPickups++
	case 41:
		rt.collectBonusValue(rt.ChestRewardValue)
	case 42:
		rt.SpecialPickup42 = true
		rt.CompassEnabled = true
		rt.SpecialPickups++
		rt.queueTutorialScript(11)
	case 51, 52, 53:
		bit := 2
		if rt.ChestRewardID == 53 {
			bit = 0
		} else if rt.ChestRewardID == 51 {
			bit = 1
		}
		rt.RelicMask |= 1 << bit
		rt.SpecialPickups++
		rt.markPersistentRewardAt(rt.Player)
		rt.RelicCelebrating = true
		rt.RelicCelebrationTicks = 0
		rt.SealCollected = true
		rt.SealTicks = 0
		rt.SealStageComplete = false
		if rt.ChestRewardID == 53 && rt.Anaconda.Enabled {
			rt.Anaconda.SealCollected = true
			rt.Anaconda.SealTicks = 0
		}
	}
}

func (rt *Runtime) applyHealthRefillReward() {
	if rt.Health >= rt.MaxHealth {
		rt.TotalVioletGems += 10
		rt.collectBonusValue(10)
		return
	}
	rt.HealthRefills++
	rt.HealFull()
}

func (rt *Runtime) persistConsumedExtraLife() {
	if !rt.checkpoint.Valid {
		return
	}
	idx := rt.index(rt.Player.X, rt.Player.Y)
	if idx < 0 || idx >= len(rt.checkpoint.PlayerLayer) {
		return
	}
	rt.checkpoint.PlayerLayer[idx] = EmptyRawID
	rt.checkpoint.ObjectState[idx] = containerOpenState(rt.Foreground[idx])
	rt.checkpoint.ConsumedRewardCells[idx] = true
}

func (rt *Runtime) markPersistentRewardAt(point Point) {
	if point.X < 0 || point.Y < 0 || point.X >= rt.Width() || point.Y >= rt.Height() {
		return
	}
	idx := rt.index(point.X, point.Y)
	if isPickupContainer(rt.Foreground[idx]) {
		rt.ConsumedRewardCells[idx] = true
	}
}

// PersistentRewardCoordinates mirrors the mutable coordinate list in the
// original RMS save. In Angkor, red diamonds, awarded extra lives, and relics
// remove their container coordinate; keys, healing, and bonus gems do not.
func (rt *Runtime) PersistentRewardCoordinates() []Point {
	if rt == nil {
		return nil
	}
	points := make([]Point, 0)
	for idx, consumed := range rt.ConsumedRewardCells {
		if consumed {
			points = append(points, Point{X: idx % rt.Width(), Y: idx / rt.Width()})
		}
	}
	return points
}

// ApplyPersistentRewardCoordinates applies the source stage-init aBoolean
// check before the initial checkpoint is saved.
func (rt *Runtime) ApplyPersistentRewardCoordinates(points []Point) {
	if rt == nil {
		return
	}
	for _, point := range points {
		if point.X < 0 || point.Y < 0 || point.X >= rt.Width() || point.Y >= rt.Height() {
			continue
		}
		idx := rt.index(point.X, point.Y)
		if !isPickupContainer(rt.Foreground[idx]) {
			continue
		}
		rt.ConsumedRewardCells[idx] = true
		if rt.Anaconda.Enabled {
			rt.PlayerLayer[idx] = 41
			rt.Background[idx] = 10
			rt.ObjectState[idx] = 0
			rt.TotalVioletGems += 10
			continue
		}
		rt.PlayerLayer[idx] = EmptyRawID
		rt.ObjectState[idx] = containerOpenState(rt.Foreground[idx])
	}
}

func containerOpenState(foreground RawID) int {
	if foreground == 14 {
		return 2
	}
	return 3
}

func (rt *Runtime) finishDeath() bool {
	rt.Retries++
	if rt.ExtraLives <= 0 || !rt.checkpoint.Valid {
		return false
	}
	livesAfterDeath := rt.ExtraLives - 1
	if !rt.RestoreCheckpoint() {
		return false
	}
	rt.ExtraLives = livesAfterDeath
	rt.Health = rt.MaxHealth
	rt.PlayerDead = false
	rt.HurtTicks = 0
	rt.InvulnerabilityTicks = 0
	return true
}

func (rt *Runtime) RecallCheckpoint() bool {
	if !rt.CanAcceptInput() || !rt.checkpoint.Valid {
		return false
	}
	if rt.IsCheckpoint(rt.Player.X, rt.Player.Y) {
		return rt.ResetCheckpoint()
	}
	rt.emitSound(SoundDeath)
	rt.RecallUsed = true
	rt.RecallPending = true
	rt.RecallTicks = 0
	rt.PlayerMotion = ObjectMotion{}
	rt.ResetPushAttempt()
	return true
}

func (rt *Runtime) finishRecall() bool {
	rt.RecallPending = false
	rt.RecallTicks = 0
	rt.Retries++
	livesAfterRecall := rt.ExtraLives - 1
	if livesAfterRecall < 0 || !rt.checkpoint.Valid {
		rt.ExtraLives = livesAfterRecall
		rt.Health = 0
		rt.PlayerDead = true
		return false
	}
	if !rt.RestoreCheckpoint() {
		return false
	}
	rt.ExtraLives = livesAfterRecall
	rt.Health = rt.MaxHealth
	rt.PlayerDead = false
	return true
}

func (rt *Runtime) ResetCheckpoint() bool {
	if !rt.CanAcceptInput() {
		return false
	}
	extraLives := rt.ExtraLives
	if !rt.RestoreCheckpoint() {
		return false
	}
	rt.ExtraLives = extraLives
	rt.emitSound(SoundCheckpoint)
	return true
}

func (rt *Runtime) HealFull() {
	if rt.PlayerDead {
		return
	}
	rt.Health = rt.MaxHealth
}

func (rt *Runtime) CanAcceptInput() bool {
	sealTransition := rt.SealCollected && !rt.SealStageComplete
	return !sealTransition && !rt.ForegroundDemoActive && !rt.EnemyGateDemoActive && !rt.TutorialScriptActive && rt.canStartPlayerMove()
}

func (rt *Runtime) canStartPlayerMove() bool {
	return !rt.PlayerDead && !rt.RecallPending && rt.HurtTicks <= 0 && rt.PlayerMotion.Remaining <= 0 && !rt.pendingChestSet && !rt.ChestOpening && !rt.LockOpening && !rt.Hammering && !rt.Hooking && !rt.ReachedGoal
}

func (rt *Runtime) isPlayerAt(x, y int) bool {
	return rt.Player.X == x && rt.Player.Y == y
}

func isSnake(id RawID) bool {
	return id == 19 || id == 43
}

func isGravityObject(id RawID) bool {
	return id == 0 || id == 1 || id == 8 || id == 9 || id == 47
}

func isPickupContainer(id RawID) bool {
	return id == 14 || id == 33
}

func isContainerReward(id RawID) bool {
	switch id {
	case 2, 4, 5, 6, 7, 24, 26, 27, 40, 41, 42, 51, 52, 53:
		return true
	default:
		return false
	}
}

func isRelicReward(id RawID) bool {
	return id == 51 || id == 52 || id == 53
}

func isEnemyGateTarget(id RawID) bool {
	return isSnake(id) || id == 36 || id == 45 || id == 46 || id == 49
}

func isContactEnemy(id RawID) bool {
	return isSnake(id) || id == 11
}

func isRoundedGravitySupport(id RawID) bool {
	switch id {
	case 0, 1, 8, 9:
		return true
	default:
		return false
	}
}

func hookSelectable(id RawID) bool {
	switch id {
	case 0, 1, 8, 9, 11, 14, 19, 43, 47, 48:
		return true
	default:
		return false
	}
}

func (rt *Runtime) specialToolLevel() int {
	switch {
	case rt.SpecialItemMask&8 != 0:
		return 8
	case rt.SpecialItemMask&2 != 0:
		return 2
	case rt.SpecialItemMask&1 != 0:
		return 1
	default:
		return 0
	}
}

func snakeStep(dir int) (int, int) {
	switch dir {
	case 1:
		return 0, -1
	case 2:
		return 1, 0
	case 3:
		return 0, 1
	case 4:
		return -1, 0
	default:
		return 0, 0
	}
}

func reverseSnakeDirection(dir int) int {
	switch dir {
	case 1:
		return 3
	case 2:
		return 4
	case 3:
		return 1
	case 4:
		return 2
	default:
		return 0
	}
}

func (rt *Runtime) SaveSnapshot() {
	rt.checkpoint = Snapshot{
		Valid:                true,
		Player:               rt.Player,
		PlayerMotion:         rt.PlayerMotion,
		PlayerLayer:          append([]RawID(nil), rt.PlayerLayer...),
		Background:           append([]RawID(nil), rt.Background...),
		Foreground:           append([]RawID(nil), rt.Foreground...),
		ForegroundState:      append([]int(nil), rt.ForegroundState...),
		ObjectState:          append([]int(nil), rt.ObjectState...),
		ObjectMotion:         append([]ObjectMotion(nil), rt.ObjectMotion...),
		FrozenOriginal:       append([]RawID(nil), rt.FrozenOriginal...),
		ContainerLocked:      append([]bool(nil), rt.ContainerLocked...),
		ConsumedRewardCells:  append([]bool(nil), rt.ConsumedRewardCells...),
		WaterDepth:           append([]uint8(nil), rt.WaterDepth...),
		WaterInitializing:    rt.WaterInitializing,
		WaterTicks:           rt.WaterTicks,
		Water:                rt.water.clone(),
		EnemyGateGroup:       append([]int(nil), rt.EnemyGateGroup...),
		EnemyGateCounters:    cloneIntMap(rt.EnemyGateCounters),
		ActiveEnemyGateGroup: rt.ActiveEnemyGateGroup,
		Anaconda:             rt.Anaconda,
		TeutonicKnight:       rt.TeutonicKnight,
		VioletGems:           rt.VioletGems,
		RedDiamonds:          rt.RedDiamonds,
		KeyForForeground9:    rt.KeyForForeground9,
		KeyForForeground8:    rt.KeyForForeground8,
		HealthRefills:        rt.HealthRefills,
		BonusValue:           rt.BonusValue,
		BonusPickups:         rt.BonusPickups,
		SpecialItemMask:      rt.SpecialItemMask,
		SpecialPickup42:      rt.SpecialPickup42,
		CompassEnabled:       rt.CompassEnabled,
		RelicMask:            rt.RelicMask,
		SpecialPickups:       rt.SpecialPickups,
		SealCollected:        rt.SealCollected,
		SealTicks:            rt.SealTicks,
		SealStageComplete:    rt.SealStageComplete,
		LastForegroundEvent:  rt.LastForegroundEvent,
		ForegroundEvents:     rt.ForegroundEvents,
		FallingTorchTriggers: rt.FallingTorchTriggers,
		RisingFireHeight:     rt.RisingFireHeight,
		FanPhase:             rt.FanPhase,
		FanDirection:         rt.FanDirection,
		BonusTarget:          rt.BonusTarget,
		BonusTargetSet:       rt.BonusTargetSet,
		BonusRemaining:       rt.BonusRemaining,
		BonusGateOpen:        rt.BonusGateOpen,
		LocksOpened:          rt.LocksOpened,
		BreakableWalls:       rt.BreakableWalls,
		ExitOpen:             rt.ExitOpen,
		ReachedGoal:          rt.ReachedGoal,
		GoalExitSecret:       rt.GoalExitSecret,
		GoalExitDirection:    rt.GoalExitDirection,
		GoalExitComplete:     rt.GoalExitComplete,
		CheckpointProgress:   rt.CheckpointProgress,
	}
}

func (rt *Runtime) SaveCheckpointAt(x, y int) bool {
	if !rt.IsCheckpoint(x, y) {
		return false
	}
	rt.Player = Point{X: x, Y: y}
	rt.PlayerMotion = ObjectMotion{}
	order := rt.checkpointOrderAt(x, y)
	if order >= rt.CheckpointProgress {
		rt.CheckpointProgress = order + 1
	}
	rt.CheckpointPending = false
	rt.SaveSnapshot()
	return true
}

func (rt *Runtime) checkpointOrderAt(x, y int) int {
	state, ok := rt.At(BackgroundLayer, x, y)
	if !ok || state == EmptyRawID {
		return 0
	}
	return int(state)
}

func (rt *Runtime) RestoreCheckpoint() bool {
	if !rt.checkpoint.Valid {
		return false
	}
	tutorialResetFirst := rt.tutorialResetFirst
	tutorialResetSecond := rt.tutorialResetSecond
	if rt.Hooking {
		rt.finishHookAction()
	}
	rt.Player = rt.checkpoint.Player
	rt.PlayerMotion = rt.checkpoint.PlayerMotion
	rt.playerFacingDirection = 2
	rt.playerTurnOffset = 0
	rt.CheckpointPending = false
	copy(rt.PlayerLayer, rt.checkpoint.PlayerLayer)
	copy(rt.Background, rt.checkpoint.Background)
	copy(rt.Foreground, rt.checkpoint.Foreground)
	copy(rt.ForegroundState, rt.checkpoint.ForegroundState)
	copy(rt.ObjectState, rt.checkpoint.ObjectState)
	copy(rt.ObjectMotion, rt.checkpoint.ObjectMotion)
	copy(rt.FrozenOriginal, rt.checkpoint.FrozenOriginal)
	copy(rt.ContainerLocked, rt.checkpoint.ContainerLocked)
	copy(rt.ConsumedRewardCells, rt.checkpoint.ConsumedRewardCells)
	copy(rt.WaterDepth, rt.checkpoint.WaterDepth)
	rt.WaterInitializing = rt.checkpoint.WaterInitializing
	rt.WaterTicks = rt.checkpoint.WaterTicks
	rt.water = rt.checkpoint.Water.clone()
	copy(rt.EnemyGateGroup, rt.checkpoint.EnemyGateGroup)
	rt.EnemyGateCounters = cloneIntMap(rt.checkpoint.EnemyGateCounters)
	rt.ActiveEnemyGateGroup = rt.checkpoint.ActiveEnemyGateGroup
	rt.Anaconda = rt.checkpoint.Anaconda
	rt.TeutonicKnight = rt.checkpoint.TeutonicKnight
	rt.VioletGems = rt.checkpoint.VioletGems
	rt.RedDiamonds = rt.checkpoint.RedDiamonds
	rt.KeyForForeground9 = rt.checkpoint.KeyForForeground9
	rt.KeyForForeground8 = rt.checkpoint.KeyForForeground8
	rt.HealthRefills = rt.checkpoint.HealthRefills
	rt.BonusValue = rt.checkpoint.BonusValue
	rt.BonusPickups = rt.checkpoint.BonusPickups
	rt.SpecialItemMask = rt.checkpoint.SpecialItemMask
	rt.SpecialPickup42 = rt.checkpoint.SpecialPickup42
	rt.CompassEnabled = rt.checkpoint.CompassEnabled
	rt.RelicMask = rt.checkpoint.RelicMask
	rt.SpecialPickups = rt.checkpoint.SpecialPickups
	rt.SealCollected = rt.checkpoint.SealCollected
	rt.SealTicks = rt.checkpoint.SealTicks
	rt.SealStageComplete = rt.checkpoint.SealStageComplete
	rt.LastForegroundEvent = rt.checkpoint.LastForegroundEvent
	rt.ForegroundEvents = rt.checkpoint.ForegroundEvents
	rt.FallingTorchTriggers = rt.checkpoint.FallingTorchTriggers
	rt.RisingFireHeight = rt.checkpoint.RisingFireHeight
	rt.FanPhase = rt.checkpoint.FanPhase
	rt.FanDirection = rt.checkpoint.FanDirection
	rt.BonusTarget = rt.checkpoint.BonusTarget
	rt.BonusTargetSet = rt.checkpoint.BonusTargetSet
	rt.BonusRemaining = rt.checkpoint.BonusRemaining
	rt.BonusGateOpen = rt.checkpoint.BonusGateOpen
	rt.LocksOpened = rt.checkpoint.LocksOpened
	rt.BreakableWalls = rt.checkpoint.BreakableWalls
	rt.HurtTicks = 0
	rt.InvulnerabilityTicks = 0
	rt.RockHoldTicks = rockHoldDuration
	rt.DeathTicks = 0
	rt.pendingChestSet = false
	rt.pendingForegroundEventSet = false
	rt.ForegroundDemoActive = false
	rt.ForegroundDemoID = 0
	rt.ForegroundDemoPhase = 0
	rt.ForegroundDemoTicks = 0
	rt.foregroundDemoMoved = false
	rt.EnemyGateDemoActive = false
	rt.EnemyGateDemoPhase = 0
	rt.EnemyGateDemoTicks = 0
	rt.EnemyGateDemoOutboundTicks = 0
	rt.EnemyGateDemoTarget = Point{}
	rt.EnemyGateDemoTargetSet = false
	rt.EnemyGateMessageIndex = 0
	rt.EnemyGateMessageTicks = 0
	rt.FallingTorchWarningTicks = 0
	rt.FallingTorchAnimation = 0
	rt.FallingTorchAnimationTicks = 0
	rt.RisingFireAnimation = 2
	rt.RisingFireAnimationTicks = 0
	rt.RecallPending = false
	rt.RecallTicks = 0
	rt.ChestOpening = false
	rt.ChestTicks = 0
	rt.ChestRewarded = false
	rt.ChestAnimation = 0
	rt.ChestRewardID = EmptyRawID
	rt.ChestRewardValue = 0
	rt.chestOpeningFresh = false
	rt.RelicCelebrating = false
	rt.RelicCelebrationTicks = 0
	rt.LockOpening = false
	rt.LockTicks = 0
	rt.LockAnimation = 0
	rt.LockPoint = Point{}
	rt.LockForegroundID = EmptyRawID
	rt.LockRewarded = false
	rt.Hammering = false
	rt.HammerTicks = 0
	rt.HammerAnimation = 0
	rt.HammerTarget = Point{}
	rt.finishHookAction()
	rt.ResetPushAttempt()
	rt.ExitOpen = rt.checkpoint.ExitOpen
	rt.ReachedGoal = rt.checkpoint.ReachedGoal
	rt.GoalExitSecret = rt.checkpoint.GoalExitSecret
	rt.GoalExitDirection = rt.checkpoint.GoalExitDirection
	rt.GoalExitComplete = rt.checkpoint.GoalExitComplete
	rt.CheckpointProgress = rt.checkpoint.CheckpointProgress
	if rt.Stage.World == WorldAngkor && rt.Stage.Index == fallingTorchStageIndex && 18 < rt.Width() && 63 < rt.Height() {
		idx := rt.index(18, 63)
		rt.PlayerLayer[idx] = EmptyRawID
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
		rt.updatePressureDoors()
	}
	rt.restoreTutorialCheckpoint(tutorialResetFirst, tutorialResetSecond)
	return true
}

func (rt *Runtime) set(layer []RawID, x, y int, id RawID) {
	layer[rt.index(x, y)] = id
}

func (rt *Runtime) index(x, y int) int {
	return x + y*rt.Width()
}

func clampRuntime(value, low, high int) int {
	return min(max(value, low), high)
}

func cloneIntMap(in map[int]int) map[int]int {
	if in == nil {
		return nil
	}
	out := make(map[int]int, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func countRaw(layer []RawID, id RawID) int {
	count := 0
	for _, cell := range layer {
		if cell == id {
			count++
		}
	}
	return count
}

func stageVioletTotal(stage *Stage) int {
	if stage == nil {
		return 0
	}
	total := countRaw(stage.Player, 1)
	for idx, id := range stage.Player {
		if id != 41 {
			continue
		}
		value := int(stage.Background[idx])
		if value == int(EmptyRawID) || value <= 0 {
			value = 1
		}
		total += value
	}
	return total
}

func StageVioletTotal(stage *Stage) int {
	return stageVioletTotal(stage)
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
