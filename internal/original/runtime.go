package original

import "fmt"

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
	snakeStunMask               = 0xf8
	snakeStunDuration           = 0x78
)

const (
	SoundSwitch      = 0
	SoundDeath       = 2
	SoundChestOpen   = 3
	SoundChestReward = 4
	SoundHeroHurt    = 5
	SoundCheckpoint  = 9
	SoundEnemyHit    = 10
	SoundBreak       = 11
	SoundHook        = 12
	SoundBoulder     = 14
	SoundDoor        = 8
	SoundStageClear  = 15
	SoundAngkorMusic = 16
)

type ObjectMotion struct {
	DX         int
	DY         int
	Remaining  int
	RollDX     int
	RollOffset int
}

type SourceFrameResult struct {
	GravityMoved  int
	SnakesMoved   int
	CrawlersMoved int
	HazardHits    int
	RockHoldHits  int
	DigCleared    int
	VioletPickups []Point
}

type Runtime struct {
	Stage                *Stage
	Player               Point
	PlayerMotion         ObjectMotion
	EntranceMarker       Point
	EntranceScrollX      int
	EntranceDoor         Point
	EntranceDoorSet      bool
	PlayerLayer          []RawID
	Background           []RawID
	Foreground           []RawID
	ObjectState          []int
	ObjectMotion         []ObjectMotion
	Pushing              bool
	PushDX               int
	PushTicks            int
	pushTarget           Point
	EnemyGateGroup       []int
	EnemyGateCounters    map[int]int
	ActiveEnemyGateGroup int
	Checkpoints          []Point
	GoalMarkers          []Point
	Doors                []Point
	TotalVioletGems      int
	TotalRedDiamonds     int
	VioletGems           int
	RedDiamonds          int
	KeyForForeground9    int
	KeyForForeground8    int
	ExtraLives           int
	HealthRefills        int
	BonusValue           int
	BonusPickups         int
	SpecialItemMask      int
	Hammering            bool
	HammerTicks          int
	HammerAnimation      int
	HammerTarget         Point
	Hooking              bool
	HookTicks            int
	HookAnimation        int
	HookTarget           Point
	hookDX               int
	hookStepsRemaining   int
	hookCollect          bool
	hookReturning        bool
	hookTip              Point
	hookOriginalState    int
	SpecialPickup42      bool
	CompassEnabled       bool
	RelicMask            int
	SpecialPickups       int
	LastForegroundEvent  int
	ForegroundEvents     int
	BonusTarget          Point
	BonusTargetSet       bool
	BonusRemaining       int
	BonusGateOpen        bool
	LocksOpened          int
	BreakableWalls       int
	MaxHealth            int
	Health               int
	DamageTaken          int
	HitCount             int
	Retries              int
	HurtTicks            int
	InvulnerabilityTicks int
	RockHoldTicks        int
	PlayerDead           bool
	DeathTicks           int
	RecallUsed           bool
	RecallPending        bool
	RecallTicks          int
	ExitOpen             bool
	ReachedGoal          bool
	GoalExitDirection    int
	GoalExitComplete     bool
	CheckpointProgress   int
	CheckpointPending    bool
	pendingCheckpoint    Point
	pendingChest         Point
	pendingChestSet      bool
	ChestOpening         bool
	ChestTicks           int
	ChestRewarded        bool
	ChestAnimation       int
	ChestRewardID        RawID
	ChestRewardValue     int
	chestOpeningFresh    bool
	LockOpening          bool
	LockTicks            int
	LockAnimation        int
	LockPoint            Point
	LockForegroundID     RawID
	LockRewarded         bool
	lastPickupTick       int
	lastPickupTickSet    bool
	gravitySourceTick    int
	frameVioletPickups   []Point
	soundEvents          []int
	checkpoint           Snapshot
}

type Snapshot struct {
	Valid                bool
	Player               Point
	PlayerMotion         ObjectMotion
	PlayerLayer          []RawID
	Background           []RawID
	Foreground           []RawID
	ObjectState          []int
	ObjectMotion         []ObjectMotion
	EnemyGateGroup       []int
	EnemyGateCounters    map[int]int
	ActiveEnemyGateGroup int
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
	LastForegroundEvent  int
	ForegroundEvents     int
	BonusTarget          Point
	BonusTargetSet       bool
	BonusRemaining       int
	BonusGateOpen        bool
	LocksOpened          int
	BreakableWalls       int
	ExitOpen             bool
	ReachedGoal          bool
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
		ObjectState:          make([]int, stage.Width*stage.Height),
		ObjectMotion:         make([]ObjectMotion, stage.Width*stage.Height),
		ActiveEnemyGateGroup: -1,
		Checkpoints:          stage.Positions(ForegroundLayer, 4),
		GoalMarkers:          append(stage.Positions(ForegroundLayer, 5), stage.Positions(ForegroundLayer, 28)...),
		Doors:                stage.Positions(ForegroundLayer, 7),
		TotalVioletGems:      countRaw(stage.Player, 1),
		TotalRedDiamonds:     countRaw(stage.Player, 2),
		ExtraLives:           5,
		MaxHealth:            4,
		Health:               4,
		CompassEnabled:       stage.Index != 13,
		RockHoldTicks:        rockHoldDuration,
		ChestRewardID:        EmptyRawID,
	}
	rt.initBonusTarget()
	rt.initObjectState()
	rt.initEnemyGates()
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
	// Java stores (-193 << 8) | 7. Its merged high byte is 0x3f:
	// open animation phase 3 and low-nibble temporary door id 15.
	rt.Background[idx] = 0x3f
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

func (rt *Runtime) initObjectState() {
	for idx, id := range rt.PlayerLayer {
		switch {
		case isSnake(id):
			dir := int(rt.Background[idx]) & 0x7
			if dir >= 1 && dir <= 4 {
				rt.ObjectState[idx] = dir
			}
		case id == 11:
			if rt.Background[idx] == 1 {
				rt.ObjectState[idx] = 16
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
	for y := 1; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 17 {
				continue
			}
			group := int(rt.Background[idx])
			if group == int(EmptyRawID) {
				continue
			}
			aboveIdx := rt.index(x, y-1)
			if isEnemyGateTarget(rt.PlayerLayer[aboveIdx]) {
				rt.EnemyGateCounters[group]++
				rt.EnemyGateGroup[aboveIdx] = group
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
	case BackgroundLayer:
		rt.set(rt.Background, x, y, id)
	case ForegroundLayer:
		rt.set(rt.Foreground, x, y, id)
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
	if !rt.CanAcceptInput() {
		rt.ResetPushAttempt()
		return false
	}
	if dx == 0 && dy == 0 {
		rt.ResetPushAttempt()
		return false
	}
	x := rt.Player.X + dx
	y := rt.Player.Y + dy
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		rt.ResetPushAttempt()
		return false
	}
	if dy == 0 && dx != 0 {
		if playerID, _ := rt.At(PlayerLayer, x, y); playerID == 0 {
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
		rt.pendingChest = Point{X: x, Y: y}
		rt.pendingChestSet = true
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
	if !playerOK || !foregroundOK || !isContainerReward(playerID) || !isPickupContainer(foregroundID) {
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
	rt.ResetPushAttempt()
	rt.RockHoldTicks = rockHoldDuration
	switch foregroundID {
	case 0:
		rt.collectForegroundEventAt(x, y)
	case 1:
		rt.clearForegroundBlob(x, y, 1)
	case 4:
		rt.activateCheckpointAt(x, y)
	case 5, 28:
		rt.ReachedGoal = true
		rt.GoalExitComplete = false
		direction := int(rt.Background[rt.index(x, y)])
		if direction < 1 || direction > 4 {
			direction = 2
		}
		rt.GoalExitDirection = direction
	case 26:
		rt.activateEnemyGateTriggerAt(x, y)
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

// AdvanceGoalExit performs the source xBoolean auto-walk. Raw foreground 5
// stores the exit direction in its background byte, and completion begins only
// after the hero has moved more than five cells beyond the stage boundary.
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
	if !rt.CheckpointPending || rt.PlayerMotion.Remaining > 0 || rt.Player != rt.pendingCheckpoint {
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
	targetX := x + dx
	if !rt.cellEmptyForObject(targetX, y) {
		return false
	}
	rt.moveObject(x, y, targetX, y, 0)
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

func (rt *Runtime) finishGravityMotionAt(x, y int, id RawID) {
	idx := rt.index(x, y)
	if rt.ObjectState[idx]&objectDirectionMask != 3 || y+1 >= rt.Height() {
		return
	}
	if rt.cellEmptyForGravity(x, y+1) && !rt.isPlayerAt(x, y+1) {
		return
	}
	rt.ObjectState[idx] &^= objectDirectionMask
	if id == 0 {
		rt.emitSound(SoundBoulder)
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
	dy := -motion.DY * motion.Remaining
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
	rt.frameVioletPickups = rt.frameVioletPickups[:0]
	chestOpeningFresh := rt.chestOpeningFresh
	rt.TickStatus()
	rt.tickDoorAnimations(sourceTick)
	if rt.PlayerMotion.Remaining <= 0 {
		rt.CommitPendingCheckpoint()
	}
	minX := max(1, rt.Player.X-radius)
	maxX := min(rt.Width()-2, rt.Player.X+radius)
	minY := max(1, rt.Player.Y-radius)
	maxY := min(rt.Height()-2, rt.Player.Y+radius)
	crawlerProcessed := make([]bool, len(rt.PlayerLayer))
	result := SourceFrameResult{}
	if rt.tickRockHold() {
		result.RockHoldHits++
	}
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			idx := rt.index(x, y)
			if isPickupContainer(rt.Foreground[idx]) {
				rt.tickChestForegroundAt(x, y, sourceTick, chestOpeningFresh)
			}
			if rt.Foreground[idx] == 32 && rt.tickDigAnimationAt(idx, sourceTick) {
				result.DigCleared++
			}
			if rt.Hooking && rt.hookReturning && rt.HookTarget == (Point{X: x, Y: y}) {
				continue
			}
			switch id := rt.PlayerLayer[idx]; {
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
				rt.ObjectState[idx] = 0
			case id == 11:
				if rt.tickCrawlerObjectAt(x, y, crawlerProcessed) {
					result.CrawlersMoved++
				}
			case id == 22 || id == 23:
				if rt.horizontalHazardHitsPlayer(x, y, id, hazardReach) && rt.Hurt(1) {
					result.HazardHits++
				}
			}
		}
	}
	rt.updatePressureDoors()
	rt.startAdjacentLockOpening()
	result.VioletPickups = append(result.VioletPickups, rt.frameVioletPickups...)
	return result
}

func (rt *Runtime) tickRockHold() bool {
	if rt.PlayerDead || rt.Player.X < 0 || rt.Player.X >= rt.Width() || rt.Player.Y <= 0 || rt.Player.Y >= rt.Height() {
		return false
	}
	aboveIdx := rt.index(rt.Player.X, rt.Player.Y-1)
	if rt.PlayerLayer[aboveIdx] != 0 || rt.ObjectMotion[aboveIdx].Remaining > 0 {
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
	return rt.PlayerLayer[idx] == 0 && rt.ObjectMotion[idx].Remaining <= 0
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
	state := rt.ObjectState[idx]
	if state <= 0 {
		if !rt.pendingChestSet || rt.PlayerMotion.Remaining > 0 || rt.Player != (Point{X: x, Y: y}) || !isContainerReward(rt.PlayerLayer[idx]) {
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
	state := rt.ObjectState[idx]
	if sourceTick&1 == 0 {
		state++
	}
	if state >= digAnimationFrames {
		rt.Foreground[idx] = EmptyRawID
		rt.ObjectState[idx] = 0
		return true
	}
	rt.ObjectState[idx] = state
	return false
}

func (rt *Runtime) UseHammer(dx, dy int) bool {
	if !rt.CanAcceptInput() || rt.specialToolLevel() < 1 || (dx == 0 && dy == 0) {
		return false
	}
	x := rt.Player.X + dx
	y := rt.Player.Y + dy
	id, ok := rt.At(PlayerLayer, x, y)
	if !ok || (id != 30 && !isSnake(id)) {
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
	switch {
	case id == 30:
		if rt.ObjectState[idx] == 0 {
			rt.ObjectState[idx] = 1
		}
		rt.emitSound(SoundBreak)
		return true
	case id == 43 && rt.ObjectState[idx]&snakeStunMask == 0 && rt.ObjectState[idx]&0x18000 == 0:
		rt.decrementEnemyGateForObjectAt(x, y)
		rt.PlayerLayer[idx] = EmptyRawID
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
		rt.emitSound(SoundEnemyHit)
		return true
	case isSnake(id):
		rt.ObjectState[idx] = (rt.ObjectState[idx] &^ snakeStunMask) | snakeStunDuration
		rt.emitSound(SoundEnemyHit)
		return true
	default:
		return false
	}
}

func (rt *Runtime) UseSpecialBarrier(dx, dy int) bool {
	if !rt.CanAcceptInput() || rt.specialToolLevel() < 2 || (dx == 0 && dy == 0) {
		return false
	}
	x := rt.Player.X + dx
	y := rt.Player.Y + dy
	foregroundID, ok := rt.At(ForegroundLayer, x, y)
	if !ok || foregroundID != 2 {
		return false
	}
	if rt.Background[rt.index(x, y)] != 1 {
		return false
	}
	return rt.clearForegroundBlob(x, y, 2) > 0
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
		rt.hookOriginalState = rt.ObjectState[nextIdx]
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
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.PlayerLayer[idx] != 30 || rt.ObjectState[idx] <= 0 {
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
			rt.hurtFromSnakeAt(x, y, state&0x7)
		}
		return false
	}
	if stunned {
		if rt.gravitySourceTick&3 == 0 {
			remaining := max(0, (state&snakeStunMask)-8)
			rt.ObjectState[idx] = (state &^ snakeStunMask) | remaining
		}
		return false
	}
	dir := state & 0x7
	usedPendingDirection := dir == 0
	if dir == 0 {
		dir = (state & 0x7000) >> 12
		if dir == 0 {
			rt.hurtFromSnakeAt(x, y, 0)
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
			rt.hurtFromSnakeAt(x, y, dir)
			return false
		}
		reverse := reverseSnakeDirection(dir)
		rt.ObjectState[idx] = (state &^ (0x7 | 0x7000)) | reverse<<12
		rt.ObjectMotion[idx] = ObjectMotion{Remaining: 21}
		rt.hurtFromSnakeAt(x, y, dir)
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
	rt.hurtFromSnakeAt(targetX, targetY, dir)
	return true
}

func (rt *Runtime) snakeCrushedAt(x, y int) bool {
	if y <= 0 || rt.Foreground[rt.index(x, y)] == 35 {
		return false
	}
	aboveIdx := rt.index(x, y-1)
	aboveID := rt.PlayerLayer[aboveIdx]
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
	processed := make([]bool, len(rt.PlayerLayer))
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			if rt.tickCrawlerObjectAt(x, y, processed) {
				moved++
			}
		}
	}
	return moved
}

func (rt *Runtime) tickCrawlerObjectAt(x, y int, processed []bool) bool {
	idx := rt.index(x, y)
	if processed[idx] || rt.PlayerLayer[idx] != 11 {
		return false
	}
	processed[idx] = true
	if rt.isPlayerAt(x, y) {
		rt.Hurt(1)
	}
	if rt.advanceObjectMotion(idx, 5, 4) {
		return false
	}
	dir := rt.ObjectState[idx] & 0x7
	if dir == 0 {
		dir = rt.inferCrawlerDirection(x, y, rt.ObjectState[idx]&0x10 != 0)
		rt.ObjectState[idx] = (rt.ObjectState[idx] &^ 0x7) | dir
		if dir == 0 {
			return false
		}
	}
	dx, dy := snakeStep(dir)
	targetX := x + dx
	targetY := y + dy
	if rt.isPlayerAt(x, y) || rt.isPlayerAt(targetX, targetY) {
		rt.Hurt(1)
	}
	if !rt.cellEmptyForSnake(targetX, targetY) {
		rt.ObjectState[idx] = (rt.ObjectState[idx] &^ 0x7) | reverseSnakeDirection(dir)
		return false
	}
	targetIdx := rt.index(targetX, targetY)
	rt.PlayerLayer[targetIdx] = 11
	rt.ObjectState[targetIdx] = rt.ObjectState[idx]
	rt.ObjectMotion[targetIdx] = ObjectMotion{DX: dx, DY: dy, Remaining: 18}
	rt.PlayerLayer[idx] = EmptyRawID
	rt.ObjectState[idx] = 0
	rt.ObjectMotion[idx] = ObjectMotion{}
	processed[targetIdx] = true
	return true
}

func (rt *Runtime) inferCrawlerDirection(x, y int, reversed bool) int {
	if reversed {
		if rt.cellOccupiedForCrawler(x+1, y) {
			return 3
		}
		if !rt.cellOccupiedForCrawler(x, y+1) {
			return 0
		}
		if rt.cellOccupiedForCrawler(x-1, y) {
			return 1
		}
		return 4
	}
	if rt.cellOccupiedForCrawler(x-1, y) {
		return 3
	}
	if !rt.cellOccupiedForCrawler(x, y+1) {
		return 0
	}
	if rt.cellOccupiedForCrawler(x+1, y) {
		return 1
	}
	return 2
}

func (rt *Runtime) cellOccupiedForCrawler(x, y int) bool {
	playerID, ok := rt.At(PlayerLayer, x, y)
	if !ok {
		return true
	}
	if playerID != EmptyRawID {
		return true
	}
	foregroundID, _ := rt.At(ForegroundLayer, x, y)
	return foregroundID != EmptyRawID
}

func (rt *Runtime) propagateBreakableDamage(x, y int) {
	for _, delta := range []Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}} {
		nx := x + delta.X
		ny := y + delta.Y
		if id, ok := rt.At(PlayerLayer, nx, ny); ok && id == 30 {
			idx := rt.index(nx, ny)
			if rt.ObjectState[idx] == 0 {
				rt.ObjectState[idx] = 1
			}
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
	playerID, ok := rt.At(PlayerLayer, x, y)
	if !ok {
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
	case playerID == 53:
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
		return rt.Background[rt.index(x, y)] != 1
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
	if len(rt.EnemyGateGroup) == 0 {
		return
	}
	idx := rt.index(x, y)
	group := rt.EnemyGateGroup[idx]
	if group < 0 {
		return
	}
	rt.EnemyGateGroup[idx] = -1
	if rt.ActiveEnemyGateGroup != group {
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
				rt.Background[aboveIdx] = 0
			}
		}
	}
}

func (rt *Runtime) activateEnemyGateTriggerAt(x, y int) {
	group := int(rt.Background[rt.index(x, y)])
	if group == int(EmptyRawID) {
		rt.ActiveEnemyGateGroup = -1
	} else {
		rt.ActiveEnemyGateGroup = group
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
				rt.openDoorByID(doorID)
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

func (rt *Runtime) openDoorByID(doorID int) {
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] == 7 && rt.doorIDAt(idx) == doorID {
				rt.openDoorAt(x, y)
			}
		}
	}
}

func (rt *Runtime) closeDoorByID(doorID int) {
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] == 7 && rt.doorIDAt(idx) == doorID && !rt.isPlayerAt(x, y) {
				rt.Background[idx] = RawID(doorID)
			}
		}
	}
}

func (rt *Runtime) openDoorAt(x, y int) {
	idx := rt.index(x, y)
	doorID := rt.doorIDAt(idx)
	if doorID < 0 {
		doorID = 0
	}
	if rt.foregroundDoorOpen(x, y) || int(rt.Background[idx])&0xf0 != 0 {
		return
	}
	rt.Background[idx] = RawID(0x10 | doorID)
}

func (rt *Runtime) doorIDAt(idx int) int {
	state := rt.Background[idx]
	if state == EmptyRawID {
		return -1
	}
	return int(state) & 0x0F
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
		if id == 0 && fallingStraightDown {
			rt.Hurt(2)
		}
		return false
	}
	if playerID, ok := rt.At(PlayerLayer, toX, toY); ok && isContactEnemy(playerID) && id == 0 {
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
	if isGravityObject(id) {
		switch {
		case toY > fromY:
			state = (state &^ objectDirectionMask) | 3
		case toX > fromX:
			state = (state &^ (objectDirectionMask | gravityRollPreparing | gravityMoveLeft | gravityMoveRight)) | 2 | gravityMoveRight
		case toX < fromX:
			state = (state &^ (objectDirectionMask | gravityRollPreparing | gravityMoveLeft | gravityMoveRight)) | 4 | gravityMoveLeft
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
	rt.EnemyGateGroup[toIdx] = -1
	rt.EnemyGateGroup[fromIdx] = -1
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
// the source searches clockwise from the incoming direction for the first
// completely empty neighboring cell and applies a fresh jInt=18 knockback.
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
		if playerOK && foregroundOK && playerID == EmptyRawID && foregroundID == EmptyRawID {
			rt.Player = Point{X: x, Y: y}
			rt.PlayerMotion = ObjectMotion{DX: dx, DY: dy, Remaining: playerMoveStartOffset}
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
	rt.tickHammerAction()
	rt.tickHookAction()
	rt.tickChestOpening()
	rt.tickLockOpening()
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
	if rt.LockOpening || !rt.CanAcceptInput() || rt.Player.Y < 0 || rt.Player.Y >= rt.Height() {
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
			rt.openDoorByID(int(rt.Background[idx]) & 0x0f)
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
	case 4:
		rt.KeyForForeground9++
	case 5:
		rt.KeyForForeground8++
	case 6:
		if rt.ExtraLives < 99 {
			rt.ExtraLives++
			rt.persistConsumedExtraLife()
			return
		}
		rt.applyHealthRefillReward()
	case 7:
		rt.applyHealthRefillReward()
	case 24:
		rt.SpecialItemMask |= 1
		rt.SpecialPickups++
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
	case 51, 52, 53:
		bit := int(rt.ChestRewardID - 51)
		rt.RelicMask |= 1 << bit
		rt.SpecialPickups++
	}
}

func (rt *Runtime) applyHealthRefillReward() {
	if rt.Health >= rt.MaxHealth {
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
	rt.checkpoint.ObjectState[idx] = max(rt.checkpoint.ObjectState[idx], 3)
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
	return !rt.PlayerDead && !rt.RecallPending && rt.HurtTicks <= 0 && rt.PlayerMotion.Remaining <= 0 && !rt.pendingChestSet && !rt.ChestOpening && !rt.LockOpening && !rt.Hammering && !rt.Hooking && !rt.ReachedGoal
}

func (rt *Runtime) isPlayerAt(x, y int) bool {
	return rt.Player.X == x && rt.Player.Y == y
}

func isSnake(id RawID) bool {
	return id == 19 || id == 43
}

func isGravityObject(id RawID) bool {
	return id == 0 || id == 1
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
		ObjectState:          append([]int(nil), rt.ObjectState...),
		ObjectMotion:         append([]ObjectMotion(nil), rt.ObjectMotion...),
		EnemyGateGroup:       append([]int(nil), rt.EnemyGateGroup...),
		EnemyGateCounters:    cloneIntMap(rt.EnemyGateCounters),
		ActiveEnemyGateGroup: rt.ActiveEnemyGateGroup,
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
		LastForegroundEvent:  rt.LastForegroundEvent,
		ForegroundEvents:     rt.ForegroundEvents,
		BonusTarget:          rt.BonusTarget,
		BonusTargetSet:       rt.BonusTargetSet,
		BonusRemaining:       rt.BonusRemaining,
		BonusGateOpen:        rt.BonusGateOpen,
		LocksOpened:          rt.LocksOpened,
		BreakableWalls:       rt.BreakableWalls,
		ExitOpen:             rt.ExitOpen,
		ReachedGoal:          rt.ReachedGoal,
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
	if rt.Hooking {
		rt.finishHookAction()
	}
	rt.Player = rt.checkpoint.Player
	rt.PlayerMotion = rt.checkpoint.PlayerMotion
	rt.CheckpointPending = false
	copy(rt.PlayerLayer, rt.checkpoint.PlayerLayer)
	copy(rt.Background, rt.checkpoint.Background)
	copy(rt.Foreground, rt.checkpoint.Foreground)
	copy(rt.ObjectState, rt.checkpoint.ObjectState)
	copy(rt.ObjectMotion, rt.checkpoint.ObjectMotion)
	copy(rt.EnemyGateGroup, rt.checkpoint.EnemyGateGroup)
	rt.EnemyGateCounters = cloneIntMap(rt.checkpoint.EnemyGateCounters)
	rt.ActiveEnemyGateGroup = rt.checkpoint.ActiveEnemyGateGroup
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
	rt.LastForegroundEvent = rt.checkpoint.LastForegroundEvent
	rt.ForegroundEvents = rt.checkpoint.ForegroundEvents
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
	rt.RecallPending = false
	rt.RecallTicks = 0
	rt.ChestOpening = false
	rt.ChestTicks = 0
	rt.ChestRewarded = false
	rt.ChestAnimation = 0
	rt.ChestRewardID = EmptyRawID
	rt.ChestRewardValue = 0
	rt.chestOpeningFresh = false
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
	rt.GoalExitDirection = rt.checkpoint.GoalExitDirection
	rt.GoalExitComplete = rt.checkpoint.GoalExitComplete
	rt.CheckpointProgress = rt.checkpoint.CheckpointProgress
	return true
}

func (rt *Runtime) set(layer []RawID, x, y int, id RawID) {
	layer[rt.index(x, y)] = id
}

func (rt *Runtime) index(x, y int) int {
	return x + y*rt.Width()
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

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
