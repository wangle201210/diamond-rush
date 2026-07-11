package original

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestRuntimeInitializesFromAngkorStageEntrance(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if rt.EntranceMarker != (Point{X: 4, Y: 17}) {
		t.Fatalf("entrance marker = %+v, want (4,17)", rt.EntranceMarker)
	}
	if rt.Player != (Point{X: 0, Y: 17}) {
		t.Fatalf("player = %+v, want Java-style start at x=0,y=17", rt.Player)
	}
	id, ok := rt.At(PlayerLayer, 4, 17)
	if !ok {
		t.Fatal("entrance lookup failed")
	}
	if id != EmptyRawID {
		t.Fatalf("runtime entrance cell = %d, want empty raw 255 after init", id)
	}
	originalID, _ := stage.At(PlayerLayer, 4, 17)
	if originalID != EntranceRawID {
		t.Fatalf("source stage mutated: entrance = %d, want %d", originalID, EntranceRawID)
	}
	if !rt.CompassEnabled {
		t.Fatal("Stage 1 compass is disabled; the original enters it after the intro-stage compass pickup")
	}
	if !rt.EntranceDoorSet || rt.EntranceDoor != (Point{X: 2, Y: 17}) {
		t.Fatalf("entrance door = %+v set=%v, want open temporary door at (2,17)", rt.EntranceDoor, rt.EntranceDoorSet)
	}
	if foreground, _ := rt.At(ForegroundLayer, 2, 17); foreground != 7 {
		t.Fatalf("entrance door foreground = %d, want raw 7", foreground)
	}
	if state, _ := rt.At(BackgroundLayer, 2, 17); state != 0x3f || !rt.IsPassable(2, 17) {
		t.Fatalf("entrance door state = %#x passable=%v, want open state 0x3f", state, rt.IsPassable(2, 17))
	}
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close source entrance door")
	}
	if state, _ := rt.At(BackgroundLayer, 2, 17); state != 0x0f || rt.IsPassable(2, 17) {
		t.Fatalf("closed entrance door state = %#x passable=%v, want blocking state 0x0f", state, rt.IsPassable(2, 17))
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundBoulder {
		t.Fatalf("entrance close sounds=%v, want [%d]", events, SoundBoulder)
	}
	if rt.CloseEntranceDoor() {
		t.Fatal("already-closed entrance door closed twice")
	}
	if events := rt.DrainSoundEvents(); len(events) != 0 {
		t.Fatalf("already-closed entrance sounds=%v, want none", events)
	}
}

func TestCompassDirectionMatchesSourceSixteenWayMapping(t *testing.T) {
	tests := []struct {
		dx, dy int
		want   int
	}{
		{0, -1, 0}, {1, -2, 1}, {1, -1, 2}, {2, -1, 3},
		{1, 0, 4}, {2, 1, 5}, {1, 1, 6}, {1, 2, 7},
		{0, 1, 8}, {-1, 2, 9}, {-1, 1, 10}, {-2, 1, 11},
		{-1, 0, 12}, {-2, -1, 13}, {-1, -1, 14}, {-1, -2, 15},
	}
	for _, tt := range tests {
		if got := CompassDirection(tt.dx, tt.dy); got != tt.want {
			t.Errorf("CompassDirection(%d,%d) = %d, want %d", tt.dx, tt.dy, got, tt.want)
		}
	}
}

func TestRuntimeCompassTargetsOrderedCheckpointsThenRaw5Goal(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		progress int
		want     Point
	}{
		{0, Point{X: 4, Y: 17}},
		{1, Point{X: 11, Y: 9}},
		{2, Point{X: 17, Y: 9}},
		{3, Point{X: 22, Y: 9}},
	}
	for _, tt := range tests {
		rt.CheckpointProgress = tt.progress
		got, ok := rt.NextCompassTarget()
		if !ok || got != tt.want {
			t.Errorf("progress %d target = %+v,%v, want %+v,true", tt.progress, got, ok, tt.want)
		}
	}
}

func TestRuntimeIndexesSourceAnchoredObjects(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if len(rt.Checkpoints) != 5 {
		t.Fatalf("checkpoints = %d, want 5", len(rt.Checkpoints))
	}
	if len(rt.GoalMarkers) != 1 {
		t.Fatalf("goal markers = %d, want 1", len(rt.GoalMarkers))
	}
	if len(rt.Doors) != 3 {
		t.Fatalf("doors = %d, want 3", len(rt.Doors))
	}
	for _, pt := range rt.Checkpoints {
		if !rt.IsCheckpoint(pt.X, pt.Y) {
			t.Fatalf("checkpoint index contains non-checkpoint %+v", pt)
		}
	}
}

func TestRuntimeInitializesBonusQuotaMarker(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if !rt.BonusTargetSet {
		t.Fatal("bonus target set = false, want true")
	}
	if rt.BonusTarget != (Point{X: 20, Y: 9}) {
		t.Fatalf("bonus target = %+v, want (20,9)", rt.BonusTarget)
	}
	if rt.BonusRemaining != 10 || rt.BonusGateOpen {
		t.Fatalf("bonus remaining=%d open=%v, want 10/false", rt.BonusRemaining, rt.BonusGateOpen)
	}
	id, _ := rt.At(PlayerLayer, 20, 9)
	if id != 12 {
		t.Fatalf("bonus marker runtime cell = %d, want raw12 until quota is met", id)
	}
}

func TestRuntimeStage00Raw12QuotaBlocksItsCellUntilMet(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if rt.TotalVioletGems != 21 {
		t.Fatalf("total violet gems = %d, want 21", rt.TotalVioletGems)
	}
	if rt.TotalRedDiamonds != 1 {
		t.Fatalf("total red diamonds = %d, want 1", rt.TotalRedDiamonds)
	}
	if !rt.ExitOpen || !rt.CanExit() {
		t.Fatal("raw12 quota incorrectly closed the source raw5 exit")
	}
	rt.Player = Point{X: rt.BonusTarget.X - 1, Y: rt.BonusTarget.Y}
	if rt.TryMove(1, 0) {
		t.Fatal("moved through raw12 before satisfying its quota")
	}
	for i := 0; i < 9; i++ {
		rt.consumeBonusQuota(1)
	}
	marker, _ := rt.At(PlayerLayer, rt.BonusTarget.X, rt.BonusTarget.Y)
	if marker != 12 || rt.BonusGateOpen {
		t.Fatalf("marker=%d open=%v with one remaining, want raw12/false", marker, rt.BonusGateOpen)
	}
	rt.consumeBonusQuota(1)
	if !rt.BonusGateOpen || !rt.ExitOpen || !rt.CanExit() {
		t.Fatalf("gate=%v exit=%v canExit=%v after ten gems, want all true", rt.BonusGateOpen, rt.ExitOpen, rt.CanExit())
	}
	marker, _ = rt.At(PlayerLayer, rt.BonusTarget.X, rt.BonusTarget.Y)
	if marker != EmptyRawID || !rt.TryMove(1, 0) {
		t.Fatalf("opened raw12 marker=%d move=%v, want empty/passable", marker, rt.Player == rt.BonusTarget)
	}
	settleRuntimePlayerMotion(rt)
	goal := rt.GoalMarkers[0]
	rt.Player = Point{X: goal.X - 1, Y: goal.Y}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.TryMove(1, 0) {
		t.Fatal("move onto goal failed")
	}
	if !rt.ReachedGoal {
		t.Fatal("goal not reached after entering raw5 exit")
	}
	if rt.GoalExitDirection != 2 {
		t.Fatalf("goal exit direction = %d, want source background direction 2", rt.GoalExitDirection)
	}
}

func TestRuntimeTracksStage00VioletGemsAsCollection(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if rt.TotalVioletGems != 21 {
		t.Fatalf("total violet gems = %d, want 21", rt.TotalVioletGems)
	}
	rt.Player = Point{X: 6, Y: 2}
	if !rt.TryMove(1, 0) {
		t.Fatal("move into violet gem failed")
	}
	result := rt.TickSourceFrame(8, 1, 0)
	if rt.VioletGems != 1 {
		t.Fatalf("violet gems = %d, want 1", rt.VioletGems)
	}
	if len(result.VioletPickups) != 1 || result.VioletPickups[0] != (Point{X: 7, Y: 2}) {
		t.Fatalf("violet pickup effects = %+v, want [(7,2)]", result.VioletPickups)
	}
}

func TestRuntimeWorld0VioletTotalsIncludeRaw41Values(t *testing.T) {
	want := []int{21, 19, 37, 30, 97, 17, 67, 65, 0, 187, 116, 233, 398, 8}
	for stageIndex, total := range want {
		stage := mustLoadOriginalStage(t, fmt.Sprintf("stage%02d.json", stageIndex))
		rt, err := NewRuntime(stage)
		if err != nil {
			t.Fatalf("Stage %d: %v", stageIndex+1, err)
		}
		if rt.TotalVioletGems != total {
			t.Errorf("Stage %d total violet=%d, want raw1 + raw41 values = %d", stageIndex+1, rt.TotalVioletGems, total)
		}
	}
}

func TestRuntimeStage00GoalAutoWalksPastSourceBoundary(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 21, Y: 9}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter stage00 raw5 goal")
	}
	settleRuntimePlayerMotion(rt)
	steps := 0
	for !rt.GoalExitComplete && steps < 20 {
		moved, _ := rt.AdvanceGoalExit()
		if !moved {
			t.Fatalf("goal exit did not move at step %d with motion %+v", steps+1, rt.PlayerMotion)
		}
		steps++
		settleRuntimePlayerMotion(rt)
	}
	if steps != 10 || rt.Player != (Point{X: stage.Width + 6, Y: 9}) {
		t.Fatalf("goal exit steps=%d player=%+v, want 10 and (%d,9)", steps, rt.Player, stage.Width+6)
	}
}

func TestRuntimeCheckpointRestoresMutableLayers(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	cp := rt.Checkpoints[0]
	if !rt.SaveCheckpointAt(cp.X, cp.Y) {
		t.Fatalf("SaveCheckpointAt(%+v) failed", cp)
	}
	rt.Player = Point{X: cp.X + 1, Y: cp.Y}
	if ok := rt.SetForTest(PlayerLayer, cp.X, cp.Y, 1); !ok {
		t.Fatal("SetForTest player failed")
	}
	if ok := rt.SetForTest(BackgroundLayer, cp.X, cp.Y, 2); !ok {
		t.Fatal("SetForTest background failed")
	}
	if ok := rt.SetForTest(ForegroundLayer, cp.X, cp.Y, EmptyRawID); !ok {
		t.Fatal("SetForTest foreground failed")
	}
	rt.VioletGems = rt.TotalVioletGems
	rt.RedDiamonds = 2
	rt.KeyForForeground9 = 1
	rt.KeyForForeground8 = 1
	rt.ExtraLives = 1
	rt.HealthRefills = 1
	rt.BonusValue = 50
	rt.BonusPickups = 2
	rt.SpecialItemMask = 7
	rt.SpecialPickup42 = true
	rt.RelicMask = 1
	rt.SpecialPickups = 3
	rt.LastForegroundEvent = 30
	rt.ForegroundEvents = 2
	rt.ActiveEnemyGateGroup = 2
	rt.EnemyGateCounters[0] = 3
	rt.EnemyGateGroup[rt.index(cp.X+1, cp.Y)] = 0
	rt.Health = 1
	rt.DamageTaken = 2
	rt.PlayerDead = false
	rt.ExitOpen = true
	rt.ReachedGoal = true
	if !rt.RestoreCheckpoint() {
		t.Fatal("RestoreCheckpoint failed")
	}
	if rt.Player != cp {
		t.Fatalf("restored player = %+v, want checkpoint %+v", rt.Player, cp)
	}
	foreground, _ := rt.At(ForegroundLayer, cp.X, cp.Y)
	if foreground != 4 {
		t.Fatalf("restored foreground checkpoint = %d, want 4", foreground)
	}
	if rt.VioletGems != 0 || rt.RedDiamonds != 0 || !rt.ExitOpen {
		t.Fatalf("restored counters violet=%d red=%d exit=%v, want zero/zero/open", rt.VioletGems, rt.RedDiamonds, rt.ExitOpen)
	}
	if rt.Health != 1 || rt.DamageTaken != 2 || rt.PlayerDead {
		t.Fatalf("restored health=%d damage=%d dead=%v, want current 1/2/false", rt.Health, rt.DamageTaken, rt.PlayerDead)
	}
	if rt.KeyForForeground9 != 0 || rt.KeyForForeground8 != 0 || rt.HealthRefills != 0 {
		t.Fatalf("restored pickups key9=%d key8=%d refills=%d, want zeros", rt.KeyForForeground9, rt.KeyForForeground8, rt.HealthRefills)
	}
	if rt.ExtraLives != 1 {
		t.Fatalf("restored extra lives = %d, want preserved current value 1", rt.ExtraLives)
	}
	if rt.BonusValue != 0 || rt.BonusPickups != 0 {
		t.Fatalf("restored bonus value=%d pickups=%d, want zeros", rt.BonusValue, rt.BonusPickups)
	}
	if rt.SpecialItemMask != 0 || rt.SpecialPickup42 || rt.RelicMask != 0 || rt.SpecialPickups != 0 {
		t.Fatalf("restored special mask=%d raw42=%v relic=%d pickups=%d, want zero/false/zero/zero", rt.SpecialItemMask, rt.SpecialPickup42, rt.RelicMask, rt.SpecialPickups)
	}
	if rt.LastForegroundEvent != 0 || rt.ForegroundEvents != 0 {
		t.Fatalf("restored foreground event last=%d count=%d, want zero/zero", rt.LastForegroundEvent, rt.ForegroundEvents)
	}
	if rt.ActiveEnemyGateGroup != -1 || rt.EnemyGateCounters[0] != 0 || rt.EnemyGateGroup[rt.index(cp.X+1, cp.Y)] != -1 {
		t.Fatalf("restored active group=%d enemy gate counter=%d group=%d, want -1/0/-1", rt.ActiveEnemyGateGroup, rt.EnemyGateCounters[0], rt.EnemyGateGroup[rt.index(cp.X+1, cp.Y)])
	}
	if rt.ReachedGoal {
		t.Fatal("restored reached goal = true, want false")
	}
}

func TestRuntimeMovesThroughStage00EntranceCorridor(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		if !rt.TryMove(1, 0) {
			t.Fatalf("move %d right failed at %+v", i+1, rt.Player)
		}
		settleRuntimePlayerMotion(rt)
	}
	if rt.Player != (Point{X: 4, Y: 17}) {
		t.Fatalf("player = %+v, want checkpoint cell (4,17)", rt.Player)
	}
	if !rt.IsCheckpoint(rt.Player.X, rt.Player.Y) {
		t.Fatalf("player should be on checkpoint at %+v", rt.Player)
	}
}

func TestRuntimeStage00CanBeCompletedFromEntrance(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}

	move := func(label string, dx, dy, count int) {
		t.Helper()
		for step := 1; step <= count; step++ {
			targetX := rt.Player.X + dx
			targetY := rt.Player.Y + dy
			targetID, _ := rt.At(PlayerLayer, targetX, targetY)
			if !rt.TryMove(dx, dy) {
				playerID, _ := rt.At(PlayerLayer, targetX, targetY)
				foregroundID, _ := rt.At(ForegroundLayer, targetX, targetY)
				t.Fatalf("%s step %d/%d failed at %+v: target (%d,%d) player raw %d foreground raw %d health=%d hurt=%d dead=%v", label, step, count, rt.Player, targetX, targetY, playerID, foregroundID, rt.Health, rt.HurtTicks, rt.PlayerDead)
			}
			if targetID == 1 {
				rt.tickGravityObjectAt(targetX, targetY)
			}
			if targetID == 10 {
				idx := rt.index(targetX, targetY)
				rt.PlayerLayer[idx] = EmptyRawID
				rt.Foreground[idx] = 32
				rt.ObjectState[idx] = 0
			}
			settleRuntimePlayerMotion(rt)
			rt.CommitPendingCheckpoint()
			rt.SettlePlayerMove()
		}
	}
	push := func(label string, dx int) {
		t.Helper()
		for attempt := 1; attempt <= boulderPushAttempts; attempt++ {
			moved := rt.TryMove(dx, 0)
			if attempt < boulderPushAttempts && moved {
				t.Fatalf("%s moved on attempt %d, want attempt %d", label, attempt, boulderPushAttempts)
			}
			if attempt == boulderPushAttempts && !moved {
				t.Fatalf("%s remained blocked after %d attempts", label, attempt)
			}
		}
		settleRuntimePlayerMotion(rt)
	}
	gravity := func(label string, sourceTicks, wantAtLeast int) {
		t.Helper()
		maxMoved := 0
		for sourceTick := 0; sourceTick < sourceTicks; sourceTick++ {
			maxMoved = max(maxMoved, rt.TickGravityNearPlayer(8))
		}
		if maxMoved < wantAtLeast {
			t.Fatalf("%s gravity moved at most %d boulders in one source-cell interval, want at least %d at %+v", label, maxMoved, wantAtLeast, rt.Player)
		}
	}

	// Collect the five gems below the opening rock row, then take the lower
	// corridor around the first snake to the right-hand rock column.
	move("entrance", 1, 0, 6)
	move("first gem row up", 0, -1, 2)
	move("first gem row", 1, 0, 6)
	move("return to lower corridor", 0, 1, 3)
	move("lower corridor right", 1, 0, 9)
	if rt.Player.X > 21 {
		move("correct snake knockback before right shaft", -1, 0, rt.Player.X-21)
	}
	move("right corridor up", 0, -1, rt.Player.Y-11)
	move("dig right of rock column", -1, 0, 2)
	move("dig down rock column", 0, 1, 2)
	if rt.VioletGems != 5 {
		t.Fatalf("violet gems before first rock puzzle = %d, want 5", rt.VioletGems)
	}

	// Standing below the cleared right-hand cells prevents the middle rock
	// from rolling while the top rock starts winding up. Leave before the
	// top rock reaches its source roll threshold and drops into this cell.
	rt.TickGravityNearPlayer(8)
	move("leave rock column pocket", 1, 0, 1)
	gravity("open right rock column", 30, 1)
	move("return above rolled rock", 0, -1, 2)
	move("cross rock column", -1, 0, 4)

	// Dig below the two rocks at x=13. They fall onto the solid ledge and
	// leave the vertical shaft to the first upper checkpoint open.
	move("dig below checkpoint rocks", -1, 0, 3)
	move("dig second cell below rocks", 0, 1, 1)
	move("step clear of falling rocks", 1, 0, 1)
	gravity("checkpoint rocks fall 1", 4, 2)
	gravity("checkpoint rocks fall 2", 4, 2)
	gravity("checkpoint rocks fall 3", 4, 2)
	move("return to opened rock shaft", 0, -1, 1)
	move("climb opened rock shaft", -1, 0, 1)
	move("climb to checkpoint corridor", 0, -1, 2)
	move("first checkpoint", -1, 0, 2)
	if rt.CheckpointProgress != 2 {
		t.Fatalf("checkpoint progress after raw state 1 = %d, want 2", rt.CheckpointProgress)
	}
	move("left detour entry", -1, 0, 7)
	move("left detour climb", 0, -1, 3)
	move("left detour east", 1, 0, 3)
	move("cross detour snake", 0, -1, 2)
	for rt.HurtTicks > 0 {
		rt.TickStatus()
	}
	move("detour around center wall", 1, 0, 2)
	if rt.Player.X < 9 {
		move("correct detour snake knockback", 1, 0, 9-rt.Player.X)
	}
	move("return to center row", 0, 1, 2)
	move("upper connector east", 1, 0, 6)
	move("dig behind corridor rock", 0, 1, 3)
	push("push corridor rock left", -1)
	move("second checkpoint", 1, 0, 3)
	if rt.CheckpointProgress != 3 {
		t.Fatalf("checkpoint progress after raw state 2 = %d, want 3", rt.CheckpointProgress)
	}

	// The visible raw-12 gate requires ten gems. Five more are reachable above
	// the final checkpoint, while the separate raw-5 exit has no quota check.
	move("upper room entry", 1, 0, 2)
	move("upper room climb", 0, -1, 3)
	move("upper room left gem", -1, 0, 1)
	move("upper room right gems", 1, 0, 2)
	move("upper room north", 0, -1, 2)
	move("upper room west", -1, 0, 2)
	move("upper room final gems", 0, -1, 1)
	move("upper room gem pair", -1, 0, 2)
	move("upper room lower gem", 0, 1, 1)
	move("upper room tenth gem", 1, 0, 1)
	if rt.VioletGems != 10 || rt.BonusRemaining != 0 || !rt.BonusGateOpen {
		t.Fatalf("objective state violet=%d remaining=%d open=%v, want 10/0/true", rt.VioletGems, rt.BonusRemaining, rt.BonusGateOpen)
	}

	move("return east", 1, 0, 2)
	move("return south", 0, 1, 5)
	move("exit corridor", 1, 0, 3)
	if rt.Player != (Point{X: 22, Y: 9}) || !rt.ReachedGoal {
		t.Fatalf("stage00 finish player=%+v reached=%v, want goal (22,9)/true", rt.Player, rt.ReachedGoal)
	}
}

func TestRuntimeStage00CanBeCompletedAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	sourceTick := 0
	tickWorld := func() {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	tickUpdate := func() {
		tickWorld()
		if rt.PlayerMotion.Remaining > 0 {
			rt.AdvancePlayerMotion()
		}
	}
	move := func(label string, dx, dy, count int) {
		t.Helper()
		for step := 1; step <= count; step++ {
			pushWaits := 0
			for {
				if rt.PlayerMotion.Remaining > 0 || rt.HurtTicks > 0 {
					tickUpdate()
					continue
				}
				if rt.TryMove(dx, dy) {
					break
				}
				targetX := rt.Player.X + dx
				targetY := rt.Player.Y + dy
				playerID, _ := rt.At(PlayerLayer, targetX, targetY)
				foregroundID, _ := rt.At(ForegroundLayer, targetX, targetY)
				if playerID != 0 || dy != 0 || pushWaits >= boulderPushAttempts-1 {
					t.Fatalf("%s step %d/%d failed at tick %d player=%+v target=(%d,%d) raw=%d foreground=%d health=%d hurt=%d motion=%+v pendingChest=%v chest=%v goal=%v dead=%v", label, step, count, sourceTick, rt.Player, targetX, targetY, playerID, foregroundID, rt.Health, rt.HurtTicks, rt.PlayerMotion, rt.pendingChestSet, rt.ChestOpening, rt.ReachedGoal, rt.PlayerDead)
				}
				pushWaits++
				tickUpdate()
			}
			for motionTick := 0; motionTick < 3; motionTick++ {
				tickUpdate()
			}
			tickUpdate()
		}
	}
	fallWait := func(cells int) {
		for i := 0; i < cells*4; i++ {
			tickUpdate()
		}
	}
	rollWait := func(cells int) {
		for i := 0; i < 26+cells*4; i++ {
			tickUpdate()
		}
	}

	tickWorld()
	move("entrance", 1, 0, 6)
	move("first gem row up", 0, -1, 2)
	move("first gem row", 1, 0, 6)
	move("return to lower corridor", 0, 1, 3)
	move("lower corridor right", 1, 0, 9)
	if rt.Player.X > 21 {
		move("correct snake knockback before right shaft", -1, 0, rt.Player.X-21)
	}
	move("right corridor up", 0, -1, rt.Player.Y-11)
	move("dig right of rock column", -1, 0, 2)
	move("dig down rock column", 0, 1, 2)
	move("leave rock column pocket", 1, 0, 1)
	rollWait(1)
	move("return above rolled rock", 0, -1, 2)
	move("cross rock column", -1, 0, 4)
	move("dig below checkpoint rocks", -1, 0, 3)
	move("dig second cell below rocks", 0, 1, 1)
	move("step clear of falling rocks", 1, 0, 1)
	fallWait(3)
	move("return to opened rock shaft", 0, -1, 1)
	move("climb opened rock shaft", -1, 0, 1)
	move("climb to checkpoint corridor", 0, -1, 2)
	move("first checkpoint", -1, 0, 2)
	move("left detour entry", -1, 0, 7)
	move("left detour climb", 0, -1, 3)
	move("left detour east", 1, 0, 3)
	move("cross detour snake", 0, -1, 2)
	move("detour around center wall", 1, 0, 2)
	if rt.Player.X < 9 {
		move("correct detour snake knockback", 1, 0, 9-rt.Player.X)
	}
	move("return to center row", 0, 1, 2)
	move("upper connector east", 1, 0, 6)
	move("dig behind corridor rock", 0, 1, 3)
	move("push corridor rock left", -1, 0, 1)
	move("second checkpoint", 1, 0, 3)
	move("upper room entry", 1, 0, 2)
	move("upper room climb", 0, -1, 3)
	move("upper room left gem", -1, 0, 1)
	move("upper room right gems", 1, 0, 2)
	move("upper room north", 0, -1, 2)
	move("upper room west", -1, 0, 2)
	move("upper room final gems", 0, -1, 1)
	move("upper room gem pair", -1, 0, 2)
	move("upper room lower gem", 0, 1, 1)
	move("upper room tenth gem", 1, 0, 1)
	if !rt.BonusGateOpen {
		t.Fatalf("bonus quota gate not exhausted at source cadence: violet=%d remaining=%d", rt.VioletGems, rt.BonusRemaining)
	}
	move("return east", 1, 0, 2)
	move("return south", 0, 1, 5)
	move("exit corridor", 1, 0, 3)
	if rt.Player != (Point{X: 22, Y: 9}) || !rt.ReachedGoal {
		t.Fatalf("source-cadence finish player=%+v reached=%v tick=%d", rt.Player, rt.ReachedGoal, sourceTick)
	}
}

func TestRuntimeAutoSavesCheckpointOnEnter(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	probe := Point{X: 10, Y: 10}
	savedProbe, _ := rt.At(PlayerLayer, probe.X, probe.Y)
	for i := 0; i < 4; i++ {
		if !rt.TryMove(1, 0) {
			t.Fatalf("move %d right failed at %+v", i+1, rt.Player)
		}
		settleRuntimePlayerMotion(rt)
	}
	if rt.Player != (Point{X: 4, Y: 17}) {
		t.Fatalf("player = %+v, want checkpoint cell (4,17)", rt.Player)
	}
	if rt.CheckpointProgress != 1 {
		t.Fatalf("checkpoint progress = %d, want 1 after raw state 0 checkpoint", rt.CheckpointProgress)
	}
	if !rt.CheckpointPending {
		t.Fatal("checkpoint pending = false before movement settles")
	}
	if !rt.CommitPendingCheckpoint() {
		t.Fatal("failed to commit checkpoint after movement settled")
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundCheckpoint {
		t.Fatalf("checkpoint commit sounds=%v, want [%d]", events, SoundCheckpoint)
	}
	if !rt.SetForTest(PlayerLayer, probe.X, probe.Y, 1) {
		t.Fatal("failed to mutate probe cell")
	}
	rt.Player = Point{X: 5, Y: 17}
	if !rt.RestoreCheckpoint() {
		t.Fatal("RestoreCheckpoint failed")
	}
	if rt.Player != (Point{X: 4, Y: 17}) {
		t.Fatalf("restored player = %+v, want checkpoint cell (4,17)", rt.Player)
	}
	restoredProbe, _ := rt.At(PlayerLayer, probe.X, probe.Y)
	if restoredProbe != savedProbe {
		t.Fatalf("restored probe = %d, want saved raw %d", restoredProbe, savedProbe)
	}
}

func TestRuntimeStage01CanBeCompletedAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage01.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	sourceTick := 0
	tickUpdate := func() {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
		if rt.PlayerMotion.Remaining > 0 {
			rt.AdvancePlayerMotion()
		}
	}
	move := func(label string, dx, dy, count int) {
		t.Helper()
		for step := 1; step <= count; step++ {
			pushWaits := 0
			for {
				if rt.PlayerMotion.Remaining > 0 || rt.HurtTicks > 0 || rt.ChestOpening {
					tickUpdate()
					continue
				}
				if rt.TryMove(dx, dy) {
					break
				}
				targetX := rt.Player.X + dx
				targetY := rt.Player.Y + dy
				playerID, _ := rt.At(PlayerLayer, targetX, targetY)
				foregroundID, _ := rt.At(ForegroundLayer, targetX, targetY)
				if playerID != 0 || dy != 0 || pushWaits >= boulderPushAttempts-1 {
					nearby := make([]string, 0, 9)
					for ny := targetY - 1; ny <= targetY+1; ny++ {
						for nx := targetX - 1; nx <= targetX+1; nx++ {
							id, _ := rt.At(PlayerLayer, nx, ny)
							nearby = append(nearby, fmt.Sprintf("(%d,%d)=%d", nx, ny, id))
						}
					}
					t.Fatalf("%s step %d/%d failed tick=%d player=%+v target=(%d,%d) raw=%d foreground=%d motion=%+v nearby=%v remainingGems=%v", label, step, count, sourceTick, rt.Player, targetX, targetY, playerID, foregroundID, rt.PlayerMotion, nearby, runtimePointsWithRaw(rt, 1))
				}
				pushWaits++
				tickUpdate()
			}
			for rt.PlayerMotion.Remaining > 0 {
				tickUpdate()
			}
			tickUpdate()
		}
	}
	waitForEmpty := func(label string, x, y, maxTicks int) {
		t.Helper()
		for tick := 0; tick < maxTicks; tick++ {
			id, _ := rt.At(PlayerLayer, x, y)
			if id == EmptyRawID {
				return
			}
			tickUpdate()
		}
		id, _ := rt.At(PlayerLayer, x, y)
		t.Fatalf("%s did not clear (%d,%d) after %d ticks: raw=%d player=%+v", label, x, y, maxTicks, id, rt.Player)
	}

	tickUpdate()
	move("automatic entrance", 1, 0, 2)
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 2 entrance door before final auto-entry step")
	}
	move("automatic entrance checkpoint", 1, 0, 1)
	move("opening corridor", 1, 0, 8)
	move("push center boulder once", 1, 0, 1)
	move("push center boulder over drop", 1, 0, 1)
	move("cross opened center corridor", 1, 0, 6)
	move("first optional checkpoint", 0, 1, 2)
	if rt.CheckpointProgress != 2 {
		t.Fatalf("Stage 2 checkpoint progress=%d, want 2 at (19,7)", rt.CheckpointProgress)
	}
	move("return from first checkpoint", 0, -1, 2)
	move("climb to upper gem shelf", 0, -1, 1)
	move("collect upper gem", 1, 0, 2)
	move("return left from upper gem", -1, 0, 2)
	move("descend through first checkpoint", 0, 1, 3)
	move("descend below first checkpoint", 0, 1, 1)
	move("push first lower boulder left", -1, 0, 1)
	move("step onto cleared lower shelf", -1, 0, 1)
	waitForEmpty("first lower boulder roll", 17, 9, 80)
	move("descend into upper lower-room corridor", 0, 1, 1)
	move("approach corridor boulder", -1, 0, 2)
	move("push corridor boulder left", -1, 0, 1)
	move("follow corridor boulder", -1, 0, 1)
	waitForEmpty("corridor boulder roll", 13, 10, 80)
	move("descend through corridor", 0, 1, 2)
	move("cross left on middle shelf", -1, 0, 4)
	move("descend to diggable shelf", 0, 1, 1)
	move("dig toward first paired boulder", -1, 0, 2)
	move("climb beside first paired boulder", 0, -1, 1)
	move("push first paired boulder right", 1, 0, 1)
	move("return left of first paired boulder", -1, 0, 1)
	move("descend beside second lower boulder", 0, 1, 1)
	move("push second lower boulder left", -1, 0, 1)
	move("dig beneath shifted boulder", 0, 1, 1)
	move("cross under shifted boulder", -1, 0, 2)
	move("dig above shifted boulder", 0, -1, 1)
	move("push shifted boulder right", 1, 0, 1)
	move("climb beside second paired boulder", 0, -1, 1)
	move("push second paired boulder right", 1, 0, 1)
	move("collect left-room gem", -1, 0, 2)
	move("descend into left room", 0, 1, 2)
	move("cross to left-room shaft", -1, 0, 2)
	move("descend to second checkpoint", 0, 1, 2)
	move("descend below second checkpoint", 0, 1, 1)
	move("cross upper left room after snake knockback", 1, 0, 2)
	move("descend to lower-room corridor", 0, 1, 1)
	move("cross lower-room corridor", 1, 0, 9)
	move("descend beside third checkpoint", 0, 1, 1)
	move("cross third checkpoint", 1, 0, 2)
	move("dig under central boulder", 1, 0, 1)
	move("retreat from central boulder", -1, 0, 1)
	move("climb into central corridor", 0, -1, 1)
	waitForEmpty("central boulder fall", 17, 17, 80)
	move("cross central corridor", 1, 0, 3)
	move("climb right-side shaft", 0, -1, 3)
	move("cross to lower goal shaft", 1, 0, 3)
	move("climb lower goal shaft", 0, -1, 3)
	if rt.Player != (Point{X: 22, Y: 11}) {
		t.Fatalf("Stage 2 lower goal shaft player=%+v, want (22,11)", rt.Player)
	}
	move("collect middle goal-shaft gem", 1, 0, 1)
	move("collect lower goal-shaft gem", 0, 1, 1)
	move("follow falling goal-shaft gem", 0, -1, 1)
	move("step left of falling goal-shaft rocks", -1, 0, 1)
	move("dig upper goal-shaft bypass", 0, -1, 1)
	move("enter moving goal-shaft column", 1, 0, 1)
	move("race up goal-shaft column", 0, -1, 3)
	move("enter right goal chamber", 1, 0, 1)
	move("climb right goal chamber", 0, -1, 3)
	move("dig behind final goal boulder", -1, 0, 1)
	move("push final goal boulder left", -1, 0, 1)
	move("collect lower quota gem", 0, 1, 1)
	move("return above lower quota gem", 0, -1, 1)
	move("push upper quota boulder left", -1, 0, 2)
	move("collect second quota gem", 0, 1, 1)
	if rt.BonusRemaining != 0 || !rt.BonusGateOpen {
		t.Fatalf("Stage 2 quota after final gems remaining=%d open=%v", rt.BonusRemaining, rt.BonusGateOpen)
	}
	move("return above second quota gem", 0, -1, 1)
	move("return to quota shaft", 1, 0, 2)
	move("climb through quota marker", 0, -1, 2)
	move("enter Stage 2 goal", 1, 0, 1)
	if rt.Player != (Point{X: 23, Y: 2}) || !rt.ReachedGoal {
		t.Fatalf("Stage 2 finish player=%+v reached=%v, want (23,2)/true", rt.Player, rt.ReachedGoal)
	}
	for !rt.GoalExitComplete {
		if rt.PlayerMotion.Remaining > 0 {
			tickUpdate()
			continue
		}
		rt.AdvanceGoalExit()
	}
	if rt.Player.X != rt.Width()+6 {
		t.Fatalf("Stage 2 exit x=%d, want %d", rt.Player.X, rt.Width()+6)
	}
}

func TestRuntimeCheckpointDoesNotOverwriteWhenRevisited(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		if !rt.TryMove(1, 0) {
			t.Fatalf("move %d right failed at %+v", i+1, rt.Player)
		}
		settleRuntimePlayerMotion(rt)
	}
	probe := Point{X: 10, Y: 10}
	savedProbe, _ := rt.At(PlayerLayer, probe.X, probe.Y)
	if !rt.SetForTest(PlayerLayer, probe.X, probe.Y, 1) {
		t.Fatal("failed to mutate probe cell")
	}
	rt.SetForTest(PlayerLayer, 5, 17, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 5, 17, EmptyRawID)
	if !rt.TryMove(1, 0) {
		t.Fatal("move off checkpoint failed")
	}
	settleRuntimePlayerMotion(rt)
	if !rt.TryMove(-1, 0) {
		t.Fatal("move back onto checkpoint failed")
	}
	settleRuntimePlayerMotion(rt)
	if rt.CheckpointProgress != 1 {
		t.Fatalf("checkpoint progress = %d, want 1 after revisiting same raw state 0 checkpoint", rt.CheckpointProgress)
	}
	if !rt.RestoreCheckpoint() {
		t.Fatal("RestoreCheckpoint failed")
	}
	restoredProbe, _ := rt.At(PlayerLayer, probe.X, probe.Y)
	if restoredProbe != savedProbe {
		t.Fatalf("same checkpoint revisit overwrote snapshot: restored probe = %d, want saved raw %d", restoredProbe, savedProbe)
	}
}

func TestRuntimeBlocksWorldTileIDsAndDigsRaw10(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if rt.IsPassable(0, 18) {
		t.Fatal("raw player ID 110 floor tile should block movement")
	}
	if rt.TryMove(0, 1) {
		t.Fatal("moved down into blocking floor tile")
	}
	if !rt.SetForTest(PlayerLayer, 1, 17, 0) {
		t.Fatal("failed to place boulder")
	}
	if !rt.SetForTest(PlayerLayer, 2, 17, 80) {
		t.Fatal("failed to block boulder target")
	}
	if rt.TryMove(1, 0) {
		t.Fatal("pushed boulder into blocked target")
	}
	if !rt.SetForTest(PlayerLayer, 1, 17, 19) {
		t.Fatal("failed to place green snake")
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("raw 19 green snake should be source-passable contact damage")
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.SetForTest(PlayerLayer, 1, 17, 10) {
		t.Fatal("failed to place raw 10 diggable tile")
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to move into raw 10 diggable tile")
	}
	if rt.Player != (Point{X: 1, Y: 17}) {
		t.Fatalf("player = %+v, want dug cell (1,17)", rt.Player)
	}
	playerID, _ := rt.At(PlayerLayer, 1, 17)
	foregroundID, _ := rt.At(ForegroundLayer, 1, 17)
	if playerID != 10 || foregroundID != EmptyRawID {
		t.Fatalf("move frame raw10 player=%d foreground=%d, want raw10/empty until next object scan", playerID, foregroundID)
	}
	rt.TickSourceFrame(8, 1, 0)
	playerID, _ = rt.At(PlayerLayer, 1, 17)
	foregroundID, _ = rt.At(ForegroundLayer, 1, 17)
	if playerID != EmptyRawID || foregroundID != 32 {
		t.Fatalf("post-scan raw10 player=%d foreground=%d, want empty/raw32", playerID, foregroundID)
	}
}

func TestRuntimeBoulderCanEnterDugRaw10Cell(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 0)
	rt.SetForTest(PlayerLayer, 2, 17, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 17, 32)
	for attempt := 1; attempt < boulderPushAttempts; attempt++ {
		if rt.TryMove(1, 0) {
			t.Fatalf("pushed boulder early on attempt %d", attempt)
		}
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to push boulder into dug raw10 cell after source delay")
	}
	id, _ := rt.At(PlayerLayer, 2, 17)
	if id != 0 {
		t.Fatalf("dug-cell boulder = %d, want raw0", id)
	}
}

func TestRuntimeFallingBoulderKeepsSeparateDigForegroundState(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	x, y := 2, 10
	sourceIdx := rt.index(x, y)
	targetIdx := rt.index(x, y+1)
	rt.SetForTest(PlayerLayer, x, y, 0)
	rt.SetForTest(PlayerLayer, x, y+1, EmptyRawID)
	rt.SetForTest(PlayerLayer, x, y+2, 80)
	rt.SetForTest(ForegroundLayer, x, y+1, 32)
	rt.ObjectState[sourceIdx] = 0x18
	rt.ForegroundState[targetIdx] = 3

	if !rt.tickGravityObjectAt(x, y) {
		t.Fatal("boulder did not fall into the raw32 removal-animation cell")
	}
	if id, _ := rt.At(PlayerLayer, x, y+1); id != 0 {
		t.Fatalf("fallen object raw=%d, want boulder raw0", id)
	}
	wantObjectState := 0x18 | 3
	if got := rt.ObjectState[targetIdx]; got != wantObjectState {
		t.Fatalf("fallen boulder state=%#x, want rotation plus down direction %#x", got, wantObjectState)
	}
	if got := rt.ForegroundStateAt(x, y+1); got != 3 {
		t.Fatalf("dig foreground state=%d after fall, want 3", got)
	}

	if rt.tickDigAnimationAt(targetIdx, 2) {
		t.Fatal("dig foreground cleared before its final frame")
	}
	if got := rt.ForegroundStateAt(x, y+1); got != 4 {
		t.Fatalf("advanced dig foreground state=%d, want 4", got)
	}
	if got := rt.ObjectState[targetIdx]; got != wantObjectState {
		t.Fatalf("dig animation changed boulder state to %#x, want %#x", got, wantObjectState)
	}

	rt.ForegroundState[targetIdx] = digAnimationFrames - 1
	if !rt.tickDigAnimationAt(targetIdx, 4) {
		t.Fatal("final dig foreground frame did not clear")
	}
	if id, _ := rt.At(PlayerLayer, x, y+1); id != 0 || rt.ObjectState[targetIdx] != wantObjectState {
		t.Fatalf("clearing grass changed boulder raw/state to %d/%#x", id, rt.ObjectState[targetIdx])
	}
}

func TestRuntimeCheckpointPreservesForegroundAnimationState(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	idx := rt.index(2, 10)
	rt.Foreground[idx] = 32
	rt.ForegroundState[idx] = 5
	rt.SaveSnapshot()
	rt.ForegroundState[idx] = 1
	if !rt.RestoreCheckpoint() {
		t.Fatal("checkpoint restore failed")
	}
	if got := rt.ForegroundStateAt(2, 10); got != 5 {
		t.Fatalf("restored foreground state=%d, want 5", got)
	}
}

func TestRuntimeRaw10DigAnimationClearsForeground(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 10)
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to dig raw10")
	}
	if got := rt.ObjectState[rt.index(1, 17)]; got != 1 {
		t.Fatalf("activated raw10 state = %d, want 1", got)
	}
	rt.TickSourceFrame(8, 1, 0)
	if got := rt.ObjectState[rt.index(1, 17)]; got != 0 {
		t.Fatalf("new foreground raw32 state = %d, want source frame 0", got)
	}
	for tick := 1; tick < 15; tick++ {
		if got := rt.TickDigAnimations(); got != 0 {
			t.Fatalf("dig animation tick %d cleared %d cells, want 0", tick, got)
		}
	}
	if got := rt.TickDigAnimations(); got != 1 {
		t.Fatalf("final dig animation cleared %d cells, want 1", got)
	}
	foregroundID, _ := rt.At(ForegroundLayer, 1, 17)
	if foregroundID != EmptyRawID {
		t.Fatalf("finished dig foreground = %d, want empty", foregroundID)
	}
}

func TestRuntimeHammerClearsAdjacentRaw10(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 1
	rt.Player = Point{X: 0, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 10)
	if !rt.UseHammer(1, 0) {
		t.Fatal("hammer did not target adjacent raw10")
	}
	for tick := 1; rt.Hammering; tick++ {
		rt.TickSourceFrame(8, tick, 0)
	}
	if id, _ := rt.At(PlayerLayer, 1, 17); id != EmptyRawID {
		t.Fatalf("hammered raw10 player layer=%d, want empty", id)
	}
	if id, _ := rt.At(ForegroundLayer, 1, 17); id != 32 && id != EmptyRawID {
		t.Fatalf("hammered raw10 foreground=%d, want removal animation", id)
	}
}

func TestRuntimePushesBoulderHorizontally(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 0)
	rt.SetForTest(PlayerLayer, 2, 17, EmptyRawID)
	for attempt := 1; attempt < boulderPushAttempts; attempt++ {
		if rt.TryMove(1, 0) {
			t.Fatalf("push move succeeded early on attempt %d", attempt)
		}
		if !rt.Pushing || rt.PushDX != 1 || rt.PushTicks != attempt {
			t.Fatalf("attempt %d push state pushing=%v dx=%d ticks=%d", attempt, rt.Pushing, rt.PushDX, rt.PushTicks)
		}
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("push move failed after source delay")
	}
	if rt.Player != (Point{X: 1, Y: 17}) {
		t.Fatalf("player = %+v, want pushed into boulder source cell", rt.Player)
	}
	id, _ := rt.At(PlayerLayer, 2, 17)
	if id != 0 {
		t.Fatalf("pushed boulder target = %d, want raw 0", id)
	}
	if rt.Pushing || rt.PushTicks != 0 {
		t.Fatalf("push state remained after success: pushing=%v ticks=%d", rt.Pushing, rt.PushTicks)
	}
}

func TestRuntimeBoulderPushDelayResetsWhenInputIsReleased(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 0)
	rt.SetForTest(PlayerLayer, 2, 17, EmptyRawID)
	for attempt := 0; attempt < boulderPushAttempts-1; attempt++ {
		rt.TryMove(1, 0)
	}
	rt.ResetPushAttempt()
	if rt.TryMove(1, 0) {
		t.Fatal("boulder pushed immediately after releasing the direction")
	}
	if rt.PushTicks != 1 {
		t.Fatalf("push ticks after release = %d, want restarted at 1", rt.PushTicks)
	}
}

func TestRuntimeHookPullsBoulderHorizontally(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 2
	rt.Player = Point{X: 4, Y: 10}
	rt.SetForTest(PlayerLayer, 4, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 4, 10, EmptyRawID)
	for x := 5; x <= 7; x++ {
		rt.SetForTest(PlayerLayer, x, 10, EmptyRawID)
		rt.SetForTest(ForegroundLayer, x, 10, EmptyRawID)
		rt.SetForTest(PlayerLayer, x, 11, 80)
	}
	rt.SetForTest(PlayerLayer, 7, 10, 0)
	rt.ObjectState[rt.index(7, 10)] = 0
	if !rt.UseHook(1, 0) {
		t.Fatal("UseHook(right) = false, want boulder pulled")
	}
	if !rt.Hooking || rt.HookAnimation != hookRightCastAnimation || rt.HookTicks != 0 || rt.CanAcceptInput() {
		t.Fatalf("hook start active=%v animation=%d ticks=%d input=%v, want true/20/0/false", rt.Hooking, rt.HookAnimation, rt.HookTicks, rt.CanAcceptInput())
	}
	if id, _ := rt.At(PlayerLayer, 5, 10); id != 32 || rt.ObjectState[rt.index(5, 10)] != 5 || rt.ObjectMotion[rt.index(5, 10)].Remaining != 18 {
		t.Fatalf("first hook segment raw/state/motion=%d/%d/%+v, want 32/5/remaining18", id, rt.ObjectState[rt.index(5, 10)], rt.ObjectMotion[rt.index(5, 10)])
	}
	if rt.Hurt(1) || rt.Health != rt.MaxHealth {
		t.Fatalf("hook action accepted damage: health=%d", rt.Health)
	}
	for sourceTick, want := range []int{12, 6, 0} {
		rt.TickSourceFrame(8, sourceTick+1, 0)
		if got := rt.ObjectMotion[rt.index(5, 10)].Remaining; got != want {
			t.Fatalf("first hook segment tick %d remaining=%d, want %d", sourceTick+1, got, want)
		}
	}
	rt.TickSourceFrame(8, 4, 0)
	if id, _ := rt.At(PlayerLayer, 6, 10); id != 32 || rt.ObjectMotion[rt.index(6, 10)].Remaining != 12 {
		t.Fatalf("right second segment raw/motion=%d/%+v, want 32/remaining12", id, rt.ObjectMotion[rt.index(6, 10)])
	}
	for sourceTick := 5; sourceTick <= 7; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if !rt.hookReturning || rt.HookTarget != (Point{X: 7, Y: 10}) {
		t.Fatalf("hook target after cast returning=%v target=%+v, want true/(7,10)", rt.hookReturning, rt.HookTarget)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundHook {
		t.Fatalf("initial hook impact sounds=%v, want [%d]", events, SoundHook)
	}
	rt.TickSourceFrame(8, 8, 0)
	if id, _ := rt.At(PlayerLayer, 6, 10); id != 0 || rt.ObjectMotion[rt.index(6, 10)].Remaining != 0 {
		t.Fatalf("first pull step raw/motion=%d/%+v, want raw0/remaining0 after rope re-acquire", id, rt.ObjectMotion[rt.index(6, 10)])
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundHook {
		t.Fatalf("rope re-acquire sounds=%v, want [%d]", events, SoundHook)
	}
	for sourceTick := 9; sourceTick <= 13; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	pulled, _ := rt.At(PlayerLayer, 5, 10)
	source, _ := rt.At(PlayerLayer, 7, 10)
	if pulled != 0 || source != EmptyRawID || rt.Hooking {
		t.Fatalf("hook result pulled=%d source=%d active=%v, want raw0/empty/false", pulled, source, rt.Hooking)
	}
	if rt.Player != (Point{X: 4, Y: 10}) {
		t.Fatalf("player moved by hook to %+v, want unchanged", rt.Player)
	}
	if rt.ObjectState[rt.index(5, 10)] != 0 || !rt.CanAcceptInput() {
		t.Fatalf("completed hook state=%d input=%v, want restored0/true", rt.ObjectState[rt.index(5, 10)], rt.CanAcceptInput())
	}
}

func TestRuntimeHookPullsRaw48Candidate(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 2
	rt.Player = Point{X: 4, Y: 10}
	for x := 5; x <= 6; x++ {
		rt.SetForTest(PlayerLayer, x, 10, EmptyRawID)
		rt.SetForTest(ForegroundLayer, x, 10, EmptyRawID)
	}
	rt.SetForTest(PlayerLayer, 6, 10, 48)
	if !rt.UseHook(1, 0) {
		t.Fatal("UseHook(right) = false, want raw48 pulled")
	}
	tickRuntimeHookToCompletion(t, rt, 20)
	pulled, _ := rt.At(PlayerLayer, 5, 10)
	source, _ := rt.At(PlayerLayer, 6, 10)
	if pulled != 48 || source != EmptyRawID {
		t.Fatalf("hook raw48 result pulled=%d source=%d, want raw48/empty", pulled, source)
	}
}

func TestRuntimeHookCollectsVioletGemAtRange(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 2
	rt.Player = Point{X: 4, Y: 10}
	rt.SetForTest(PlayerLayer, 4, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 4, 10, EmptyRawID)
	for x := 5; x <= 7; x++ {
		rt.SetForTest(PlayerLayer, x, 10, EmptyRawID)
		rt.SetForTest(ForegroundLayer, x, 10, EmptyRawID)
		rt.SetForTest(PlayerLayer, x, 11, 80)
	}
	rt.SetForTest(PlayerLayer, 7, 10, 1)
	if !rt.UseHook(1, 0) {
		t.Fatal("UseHook(right) = false, want violet gem collected")
	}
	for sourceTick := 1; sourceTick <= 12; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if rt.VioletGems != 0 {
		t.Fatalf("violet gems before target reaches hero = %d, want 0", rt.VioletGems)
	}
	tickRuntimeHookToCompletionFrom(t, rt, 13, 24)
	source, _ := rt.At(PlayerLayer, 7, 10)
	if source != EmptyRawID {
		t.Fatalf("hook source = %d, want empty after collecting raw1", source)
	}
	if rt.VioletGems != 1 {
		t.Fatalf("violet gems = %d, want 1", rt.VioletGems)
	}
	if rt.BonusRemaining != 9 {
		t.Fatalf("bonus remaining = %d, want 9", rt.BonusRemaining)
	}
	if rt.Player != (Point{X: 4, Y: 10}) {
		t.Fatalf("player moved by hook to %+v, want unchanged", rt.Player)
	}
}

func TestRuntimeHookRequiresToolRangeAndClearPath(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	for x := 1; x <= 3; x++ {
		rt.SetForTest(PlayerLayer, x, 17, EmptyRawID)
		rt.SetForTest(ForegroundLayer, x, 17, EmptyRawID)
	}
	rt.SetForTest(PlayerLayer, 2, 17, 0)
	if rt.UseHook(1, 0) {
		t.Fatal("UseHook without hook item = true, want false")
	}
	rt.SpecialItemMask = 1
	if rt.UseHook(1, 0) {
		t.Fatal("UseHook with only raw24 tool level = true, want false")
	}
	rt.SpecialItemMask = 2
	rt.SetForTest(PlayerLayer, 1, 17, 0)
	if rt.UseHook(1, 0) {
		t.Fatal("UseHook on adjacent target = true, want false")
	}
	rt.SetForTest(PlayerLayer, 1, 17, 1)
	if rt.UseHook(1, 0) {
		t.Fatal("UseHook on adjacent raw1 = true, want false")
	}
	rt.SetForTest(PlayerLayer, 1, 17, EmptyRawID)
	rt.SetForTest(PlayerLayer, 2, 17, 10)
	rt.SetForTest(PlayerLayer, 3, 17, 0)
	if rt.UseHook(1, 0) {
		t.Fatal("UseHook through blocker = true, want false")
	}
}

func TestRuntimeBoulderFallsIntoEmptyCell(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SetForTest(PlayerLayer, 1, 10, 0)
	rt.SetForTest(PlayerLayer, 1, 11, EmptyRawID)
	moved := tickGravityObject(rt, 1, 10)
	if moved == 0 {
		t.Fatal("gravity moved 0 boulders")
	}
	above, _ := rt.At(PlayerLayer, 1, 10)
	below, _ := rt.At(PlayerLayer, 1, 11)
	if above != EmptyRawID || below != 0 {
		t.Fatalf("gravity result above=%d below=%d, want empty/raw0", above, below)
	}
}

func TestRuntimeBoulderUsesSourceMotionTimer(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SetForTest(PlayerLayer, 1, 10, 0)
	rt.SetForTest(PlayerLayer, 1, 11, EmptyRawID)
	rt.SetForTest(PlayerLayer, 1, 12, 80)
	rt.SetForTest(ForegroundLayer, 1, 11, EmptyRawID)
	if moved := tickGravityObject(rt, 1, 10); moved == 0 {
		t.Fatal("first gravity source tick did not start the fall")
	}
	idx := rt.index(1, 11)
	if got := rt.ObjectMotion[idx]; got != (ObjectMotion{DY: 1, Remaining: 18}) {
		t.Fatalf("initial boulder motion = %+v, want down/18", got)
	}
	for step, want := range []int{12, 6, 0} {
		if moved := rt.TickGravity(); moved != 0 {
			t.Fatalf("motion tick %d moved another boulder cell", step+1)
		}
		if got := rt.ObjectMotion[idx].Remaining; got != want {
			t.Fatalf("motion tick %d remaining = %d, want %d", step+1, got, want)
		}
	}
}

func TestRuntimeBoulderLandingClearsDirectionAndEmitsSourceSound(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, 80)
	idx := rt.index(2, 10)
	rt.ObjectState[idx] = 3 | gravityMoveRight | 0x10
	rt.ObjectMotion[idx] = ObjectMotion{DY: 1, Remaining: 6}
	if rt.tickGravityObjectAt(2, 10) {
		t.Fatal("supported landing boulder moved")
	}
	state := rt.ObjectState[idx]
	if state&objectDirectionMask != 0 || state&gravityMoveRight == 0 || state&boulderRotationMask != 0x18 {
		t.Fatalf("landing-frame boulder state=%#x, want direction cleared with final rotation/right marker preserved", state)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundBoulder {
		t.Fatalf("landing sound events=%v, want [%d]", events, SoundBoulder)
	}
	rt.tickGravityObjectAt(2, 10)
	state = rt.ObjectState[idx]
	if state&(gravityMoveRight|gravityMoveLeft) != 0 || state&boulderRotationMask != 0x18 {
		t.Fatalf("settled boulder state=%#x, want side marker cleared and final rotation retained", state)
	}
	if events := rt.DrainSoundEvents(); len(events) != 0 {
		t.Fatalf("stationary boulder repeated landing sound: %v", events)
	}
}

func TestRuntimeGravityObjectRenderOffsetMatchesSourceOVoid(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	x, y := 2, 10
	idx := rt.index(x, y)
	rt.SetForTest(PlayerLayer, x, y, 0)

	rt.ObjectState[idx] = 3
	rt.ObjectMotion[idx] = ObjectMotion{DY: 1, Remaining: 18}
	if dx, dy := rt.GravityObjectRenderOffset(x, y, 2); dx != 0 || dy != -18 {
		t.Fatalf("fall render offset=(%d,%d), want (0,-18)", dx, dy)
	}

	rt.ObjectState[idx] = 4 | gravityMoveLeft | gravityRollPreparing
	rt.ObjectMotion[idx] = ObjectMotion{RollDX: -1, RollOffset: 6}
	if dx, dy := rt.GravityObjectRenderOffset(x, y, 2); dx != -7 || dy != 1 {
		t.Fatalf("left roll-preparation offset=(%d,%d), want (-7,1)", dx, dy)
	}

	rt.ObjectState[idx] = 2 | gravityMoveRight | gravityRollPreparing
	rt.ObjectMotion[idx] = ObjectMotion{RollDX: 1, RollOffset: 6}
	if dx, dy := rt.GravityObjectRenderOffset(x, y, 2); dx != 5 || dy != 1 {
		t.Fatalf("right roll-preparation offset=(%d,%d), want (5,1)", dx, dy)
	}

	rt.SetForTest(PlayerLayer, x, y+1, 0)
	rt.SetForTest(PlayerLayer, x+1, y+1, EmptyRawID)
	rt.ObjectState[rt.index(x, y+1)] = 0
	rt.ObjectState[idx] = 2 | gravityMoveRight
	rt.ObjectMotion[idx] = ObjectMotion{DX: 1, Remaining: 12}
	if dx, dy := rt.GravityObjectRenderOffset(x, y, 2); dx != -11 || dy != 6 {
		t.Fatalf("rounded-support movement offset=(%d,%d), want (-11,6)", dx, dy)
	}
}

func TestRuntimePushedBoulderUsesPackedDirectionAndRotationBits(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.SetForTest(PlayerLayer, 1, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 10, EmptyRawID)
	if !rt.TryPushBoulder(1, 10, 1) {
		t.Fatal("failed to start rightward source boulder push")
	}
	idx := rt.index(2, 10)
	state := rt.ObjectState[idx]
	if state&objectDirectionMask != 2 || state&gravityMoveRight == 0 || state&boulderRotationMask != 0 {
		t.Fatalf("initial pushed boulder state=%#x, want direction2/right-bit/rotation0", state)
	}
	for step, wantRotation := range []int{1, 1, 2} {
		rt.tickGravityObjectAt(2, 10)
		rotation := (rt.ObjectState[idx] & boulderRotationMask) >> 3
		if rotation != wantRotation {
			t.Fatalf("push motion step %d rotation=%d state=%#x, want %d", step+1, rotation, rt.ObjectState[idx], wantRotation)
		}
	}
}

func TestRuntimeVioletGemFallsWithSourceGravity(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.SetForTest(PlayerLayer, 2, 10, 1)
	rt.SetForTest(PlayerLayer, 2, 11, EmptyRawID)
	rt.SetForTest(PlayerLayer, 2, 12, 80)
	rt.SetForTest(ForegroundLayer, 2, 11, EmptyRawID)
	if moved := tickGravityObject(rt, 2, 10); moved != 1 {
		t.Fatalf("falling violet gems = %d, want 1", moved)
	}
	source, _ := rt.At(PlayerLayer, 2, 10)
	target, _ := rt.At(PlayerLayer, 2, 11)
	if source != EmptyRawID || target != 1 {
		t.Fatalf("gem fall source=%d target=%d, want empty/raw1", source, target)
	}
	if motion := rt.ObjectMotion[rt.index(2, 11)]; motion.DY != 1 || motion.Remaining != 18 {
		t.Fatalf("gem motion = %+v, want down/18", motion)
	}
}

func TestRuntimeSourceRollDelayAppliesToBouldersAndGems(t *testing.T) {
	tests := []struct {
		name string
		id   RawID
	}{
		{name: "boulder", id: 0},
		{name: "violet-gem", id: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stage := mustLoadOriginalStage(t, "stage00.json")
			rt, err := NewRuntime(stage)
			if err != nil {
				t.Fatal(err)
			}
			clearRuntimePlayerIDs(rt, 0, 1)
			rt.SetForTest(PlayerLayer, 2, 10, tt.id)
			rt.SetForTest(PlayerLayer, 2, 11, 0)
			rt.SetForTest(PlayerLayer, 2, 12, 80)
			rt.SetForTest(PlayerLayer, 1, 10, EmptyRawID)
			rt.SetForTest(PlayerLayer, 1, 11, EmptyRawID)
			rt.SetForTest(PlayerLayer, 3, 10, 80)
			rt.SetForTest(ForegroundLayer, 1, 10, EmptyRawID)
			rt.SetForTest(ForegroundLayer, 1, 11, EmptyRawID)
			moved, sourceTicks := tickGravityObjectUntilMoved(rt, 2, 10, 30)
			if moved != 1 || sourceTicks != 26 {
				t.Fatalf("source roll moved=%d at tick=%d, want 1 at tick 26", moved, sourceTicks)
			}
			motion := rt.ObjectMotion[rt.index(1, 11)]
			if motion.DX != 0 || motion.DY != 1 || motion.Remaining != 12 {
				t.Fatalf("source roll motion = %+v, want vertical offset 12", motion)
			}
		})
	}
}

func TestRuntimeSourceRollDelayTracksGlobalFramePhase(t *testing.T) {
	tests := []struct {
		startTick int
		wantTicks int
	}{
		{startTick: 1, wantTicks: 26},
		{startTick: 2, wantTicks: 25},
		{startTick: 3, wantTicks: 24},
		{startTick: 4, wantTicks: 27},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("phase-%d", tt.startTick&3), func(t *testing.T) {
			stage := mustLoadOriginalStage(t, "stage00.json")
			rt, err := NewRuntime(stage)
			if err != nil {
				t.Fatal(err)
			}
			clearRuntimePlayerIDs(rt, 0, 1)
			rt.SetForTest(PlayerLayer, 2, 10, 0)
			rt.SetForTest(PlayerLayer, 2, 11, 1)
			rt.SetForTest(PlayerLayer, 2, 12, 80)
			rt.SetForTest(PlayerLayer, 1, 10, EmptyRawID)
			rt.SetForTest(PlayerLayer, 1, 11, EmptyRawID)
			rt.SetForTest(PlayerLayer, 3, 10, 80)
			rt.SetForTest(ForegroundLayer, 1, 10, EmptyRawID)
			rt.SetForTest(ForegroundLayer, 1, 11, EmptyRawID)
			rt.gravitySourceTick = tt.startTick - 1
			moved, sourceTicks := tickGravityObjectUntilMoved(rt, 2, 10, 30)
			if moved != 1 || sourceTicks != tt.wantTicks {
				t.Fatalf("roll from source tick %d moved=%d after %d ticks, want 1 after %d", tt.startTick, moved, sourceTicks, tt.wantTicks)
			}
		})
	}
}

func TestRuntimeNearGravityLeavesDistantObjectsDormant(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 1, Y: 17}
	rt.SetForTest(PlayerLayer, 20, 2, 0)
	rt.SetForTest(PlayerLayer, 20, 3, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 20, 3, EmptyRawID)
	rt.TickGravityNearPlayer(8)
	if id, _ := rt.At(PlayerLayer, 20, 2); id != 0 {
		t.Fatalf("distant boulder = %d, want dormant raw0", id)
	}
}

func TestRuntimeSourceFrameCrushesSnakeBeforeFallingBoulder(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 19, 43)
	rt.Player = Point{X: 4, Y: 12}
	rt.SetForTest(PlayerLayer, 2, 9, 0)
	rt.SetForTest(PlayerLayer, 2, 10, 19)
	rt.SetForTest(PlayerLayer, 3, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 3, 10, EmptyRawID)
	rt.ObjectState[rt.index(2, 10)] = 2
	result := rt.TickSourceFrame(8, 1, 0)
	if result.SnakesMoved != 0 || result.GravityMoved != 1 {
		t.Fatalf("first source frame result = %+v, want snake crushed before immediate vertical fall", result)
	}
	belowRock, _ := rt.At(PlayerLayer, 2, 10)
	snakeTarget, _ := rt.At(PlayerLayer, 3, 10)
	if belowRock != 0 || snakeTarget != EmptyRawID {
		t.Fatalf("ordered scan belowRock=%d snakeTarget=%d, want raw0/empty", belowRock, snakeTarget)
	}
}

func TestRuntimeBoulderRollsLeftWhenBlockedBelow(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, 1)
	rt.SetForTest(PlayerLayer, 2, 12, 80)
	rt.SetForTest(PlayerLayer, 1, 10, EmptyRawID)
	rt.SetForTest(PlayerLayer, 1, 11, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 1, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 1, 11, EmptyRawID)
	moved, _ := tickGravityObjectUntilMoved(rt, 2, 10, 30)
	if moved == 0 {
		t.Fatal("gravity moved 0 boulders")
	}
	source, _ := rt.At(PlayerLayer, 2, 10)
	target, _ := rt.At(PlayerLayer, 1, 11)
	if source != EmptyRawID || target != 0 {
		t.Fatalf("roll left source=%d target=%d, want empty/raw0", source, target)
	}
}

func TestRuntimeBoulderRollsRightWhenLeftBlocked(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, 1)
	rt.SetForTest(PlayerLayer, 2, 12, 80)
	rt.SetForTest(PlayerLayer, 1, 10, 80)
	rt.SetForTest(PlayerLayer, 1, 11, EmptyRawID)
	rt.SetForTest(PlayerLayer, 3, 10, EmptyRawID)
	rt.SetForTest(PlayerLayer, 3, 11, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 3, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 3, 11, EmptyRawID)
	moved, _ := tickGravityObjectUntilMoved(rt, 2, 10, 30)
	if moved == 0 {
		t.Fatal("gravity moved 0 boulders")
	}
	source, _ := rt.At(PlayerLayer, 2, 10)
	target, _ := rt.At(PlayerLayer, 3, 11)
	if source != EmptyRawID || target != 0 {
		t.Fatalf("roll right source=%d target=%d, want empty/raw0", source, target)
	}
}

func TestRuntimeBoulderCrushDamagesPlayer(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.Player = Point{X: 2, Y: 11}
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 11, EmptyRawID)
	rt.ObjectMotion[rt.index(2, 10)] = ObjectMotion{DY: 1}
	moved := rt.TickGravity()
	if moved != 0 {
		t.Fatalf("gravity moved %d boulders, want 0 while crushing player", moved)
	}
	if rt.Health != 2 || rt.DamageTaken != 2 || rt.PlayerDead {
		t.Fatalf("health=%d damage=%d dead=%v, want 2/2/false", rt.Health, rt.DamageTaken, rt.PlayerDead)
	}
	above, _ := rt.At(PlayerLayer, 2, 10)
	if above != 0 {
		t.Fatalf("crushing boulder = %d, want raw0 still above player", above)
	}
}

func TestRuntimeStationaryBoulderIsSupportedByPlayerBelow(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.Player = Point{X: 2, Y: 11}
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 11, EmptyRawID)
	if moved := rt.TickGravity(); moved != 0 {
		t.Fatalf("stationary boulder moved through player: %d", moved)
	}
	if rt.Health != rt.MaxHealth || rt.DamageTaken != 0 {
		t.Fatalf("stationary supported boulder health=%d damage=%d, want full/0", rt.Health, rt.DamageTaken)
	}
}

func TestRuntimeHoldingBoulderTooLongDamagesPlayer(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.Player = Point{X: 2, Y: 11}
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 11, EmptyRawID)
	for tick := 1; tick < rockHoldDuration; tick++ {
		result := rt.TickSourceFrame(8, tick, 0)
		if result.RockHoldHits != 0 || rt.Health != rt.MaxHealth {
			t.Fatalf("tick %d hold hits=%d health=%d, want safe before timeout", tick, result.RockHoldHits, rt.Health)
		}
	}
	result := rt.TickSourceFrame(8, rockHoldDuration, 0)
	if result.RockHoldHits != 1 || rt.Health != 0 || rt.DamageTaken != rt.MaxHealth || !rt.PlayerDead {
		t.Fatalf("timeout hold hits=%d health=%d damage=%d dead=%v, want lethal max-health crush", result.RockHoldHits, rt.Health, rt.DamageTaken, rt.PlayerDead)
	}
}

func TestRuntimeDeadPlayerCannotMove(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.ExtraLives = 0
	rt.Hurt(4)
	if !rt.PlayerDead {
		t.Fatal("player dead = false, want true")
	}
	if rt.TryMove(1, 0) {
		t.Fatal("dead player moved")
	}
	for tick := 0; tick < deathDuration; tick++ {
		rt.TickStatus()
	}
	if !rt.PlayerDead || rt.ExtraLives != 0 || rt.Retries != 1 {
		t.Fatalf("game-over state dead=%v lives=%d retries=%d, want true/0/1", rt.PlayerDead, rt.ExtraLives, rt.Retries)
	}
}

func TestRuntimeExtraLifeRespawnsAtCheckpoint(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	cp := rt.Checkpoints[0]
	if !rt.SaveCheckpointAt(cp.X, cp.Y) {
		t.Fatalf("SaveCheckpointAt(%+v) failed", cp)
	}
	rt.Player = Point{X: cp.X + 2, Y: cp.Y}
	rt.ExtraLives = 2
	rt.Health = 1
	rt.Hurt(2)
	deathPoint := rt.Player
	if !rt.PlayerDead || rt.DeathTicks != deathDuration {
		t.Fatalf("immediate death state dead=%v ticks=%d, want true/%d", rt.PlayerDead, rt.DeathTicks, deathDuration)
	}
	if rt.Player != deathPoint || rt.ExtraLives != 2 || rt.Retries != 0 {
		t.Fatalf("pre-transition player=%+v lives=%d retries=%d, want death point/2/0", rt.Player, rt.ExtraLives, rt.Retries)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundHeroHurt {
		t.Fatalf("lethal hit sounds=%v, want initial [%d]", events, SoundHeroHurt)
	}
	for tick := 1; tick < deathDuration; tick++ {
		rt.TickStatus()
		if !rt.PlayerDead || rt.Player != deathPoint {
			t.Fatalf("death resolved early on tick %d: dead=%v player=%+v", tick, rt.PlayerDead, rt.Player)
		}
		events := rt.DrainSoundEvents()
		if tick == hurtStateDuration {
			if len(events) != 1 || events[0] != SoundDeath {
				t.Fatalf("death transition sounds at tick %d=%v, want [%d]", tick, events, SoundDeath)
			}
		} else if len(events) != 0 {
			t.Fatalf("unexpected death transition sounds at tick %d: %v", tick, events)
		}
	}
	rt.TickStatus()
	if rt.PlayerDead || rt.DeathTicks != 0 {
		t.Fatalf("death did not resolve on tick %d: dead=%v ticks=%d", deathDuration, rt.PlayerDead, rt.DeathTicks)
	}
	if rt.Player != cp {
		t.Fatalf("respawn player = %+v, want checkpoint %+v", rt.Player, cp)
	}
	if rt.ExtraLives != 1 {
		t.Fatalf("extra lives = %d, want 1 after death", rt.ExtraLives)
	}
	if rt.Retries != 1 {
		t.Fatalf("retries = %d, want 1 after death", rt.Retries)
	}
	if rt.Health != rt.MaxHealth {
		t.Fatalf("health = %d, want full %d", rt.Health, rt.MaxHealth)
	}
}

func TestRuntimeRecallCheckpointCostsExtraLife(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	cp := rt.Checkpoints[0]
	if !rt.SaveCheckpointAt(cp.X, cp.Y) {
		t.Fatalf("SaveCheckpointAt(%+v) failed", cp)
	}
	rt.Player = Point{X: cp.X + 2, Y: cp.Y}
	rt.ExtraLives = 2
	rt.Health = 1
	if !rt.RecallCheckpoint() {
		t.Fatal("RecallCheckpoint() = false, want true")
	}
	if !rt.RecallUsed {
		t.Fatal("recall used = false, want true")
	}
	if !rt.RecallPending || rt.RecallTicks != 0 || rt.Player == cp {
		t.Fatalf("initial recall pending=%v ticks=%d player=%+v, want pending at original position", rt.RecallPending, rt.RecallTicks, rt.Player)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundDeath {
		t.Fatalf("recall start sounds=%v, want [%d]", events, SoundDeath)
	}
	if rt.ExtraLives != 2 || rt.Retries != 0 {
		t.Fatalf("pre-animation lives=%d retries=%d, want 2/0", rt.ExtraLives, rt.Retries)
	}
	for tick := 1; tick < recallAnimationDuration; tick++ {
		rt.TickStatus()
		if !rt.RecallPending || rt.Player == cp {
			t.Fatalf("recall resolved early at tick %d: pending=%v player=%+v", tick, rt.RecallPending, rt.Player)
		}
	}
	rt.TickStatus()
	if rt.RecallPending || rt.Player != cp {
		t.Fatalf("recall completion pending=%v player=%+v, want checkpoint %+v", rt.RecallPending, rt.Player, cp)
	}
	if rt.ExtraLives != 1 || rt.Retries != 1 {
		t.Fatalf("completed recall lives=%d retries=%d, want 1/1", rt.ExtraLives, rt.Retries)
	}
	if rt.Health != rt.MaxHealth || rt.PlayerDead {
		t.Fatalf("health=%d dead=%v, want full/false", rt.Health, rt.PlayerDead)
	}
}

func TestRuntimeRecallCheckpointWithoutLifeKillsPlayer(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	cp := rt.Checkpoints[0]
	if !rt.SaveCheckpointAt(cp.X, cp.Y) {
		t.Fatalf("SaveCheckpointAt(%+v) failed", cp)
	}
	rt.Player = Point{X: cp.X + 2, Y: cp.Y}
	rt.ExtraLives = 0
	if !rt.RecallCheckpoint() {
		t.Fatal("RecallCheckpoint() = false, want handled death path")
	}
	if !rt.RecallUsed {
		t.Fatal("recall used = false, want true")
	}
	if !rt.RecallPending || rt.PlayerDead {
		t.Fatalf("initial no-life recall pending=%v dead=%v, want true/false during animation", rt.RecallPending, rt.PlayerDead)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundDeath {
		t.Fatalf("no-life recall start sounds=%v, want [%d]", events, SoundDeath)
	}
	for tick := 0; tick < recallAnimationDuration; tick++ {
		rt.TickStatus()
	}
	if !rt.PlayerDead || rt.Health != 0 {
		t.Fatalf("dead=%v health=%d, want true/0", rt.PlayerDead, rt.Health)
	}
	if rt.ExtraLives != -1 || rt.Retries != 1 {
		t.Fatalf("no-life recall lives=%d retries=%d, want -1/1", rt.ExtraLives, rt.Retries)
	}
	if rt.Player == cp {
		t.Fatalf("player = %+v, want no checkpoint restore without lives", rt.Player)
	}
}

func TestRuntimeRecallOnCheckpointResetsWithoutLifeCost(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	cp := rt.Checkpoints[0]
	if !rt.SaveCheckpointAt(cp.X, cp.Y) {
		t.Fatalf("SaveCheckpointAt(%+v) failed", cp)
	}
	rt.ExtraLives = 1
	savedID, _ := rt.At(PlayerLayer, cp.X+1, cp.Y)
	rt.SetForTest(PlayerLayer, cp.X+1, cp.Y, 1)
	if !rt.RecallCheckpoint() {
		t.Fatal("RecallCheckpoint() on checkpoint = false, want reset")
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundCheckpoint {
		t.Fatalf("checkpoint reset sounds=%v, want [%d]", events, SoundCheckpoint)
	}
	if rt.RecallUsed {
		t.Fatal("recall used = true on checkpoint reset, want false")
	}
	if rt.ExtraLives != 1 {
		t.Fatalf("extra lives = %d, want unchanged 1", rt.ExtraLives)
	}
	id, _ := rt.At(PlayerLayer, cp.X+1, cp.Y)
	if id != savedID {
		t.Fatalf("restored neighbor player layer = %d, want saved %d", id, savedID)
	}
}

func TestRuntimeCheckpointRestorePreservesExtraLives(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	cp := rt.Checkpoints[0]
	if !rt.SaveCheckpointAt(cp.X, cp.Y) {
		t.Fatalf("SaveCheckpointAt(%+v) failed", cp)
	}
	pickup := Point{X: cp.X + 1, Y: cp.Y}
	savedID, _ := rt.At(PlayerLayer, pickup.X, pickup.Y)
	rt.SetForTest(PlayerLayer, pickup.X, pickup.Y, 6)
	rt.SetForTest(ForegroundLayer, pickup.X, pickup.Y, EmptyRawID)
	if !rt.TryMove(1, 0) {
		t.Fatal("move into extra life failed")
	}
	if rt.ExtraLives != 6 {
		t.Fatalf("extra lives after pickup = %d, want 6", rt.ExtraLives)
	}
	if !rt.RestoreCheckpoint() {
		t.Fatal("RestoreCheckpoint failed")
	}
	if rt.ExtraLives != 6 {
		t.Fatalf("extra lives after restore = %d, want preserved 6", rt.ExtraLives)
	}
	restoredID, _ := rt.At(PlayerLayer, pickup.X, pickup.Y)
	if restoredID != savedID {
		t.Fatalf("restored pickup cell = %d, want saved %d", restoredID, savedID)
	}
}

func TestRuntimeSnakeContactDamagesAndAllowsPlayerEntry(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 19)
	rt.SetForTest(ForegroundLayer, 1, 17, EmptyRawID)
	if !rt.TryMove(1, 0) {
		t.Fatal("player did not enter source-passable snake cell")
	}
	if rt.Player != (Point{X: 1, Y: 17}) {
		t.Fatalf("player = %+v, want snake cell", rt.Player)
	}
	if rt.Health != 4 || rt.DamageTaken != 0 {
		t.Fatalf("health=%d damage=%d before object frame, want 4/0", rt.Health, rt.DamageTaken)
	}
	rt.TickSnakes()
	if rt.Health != 3 || rt.DamageTaken != 1 {
		t.Fatalf("health=%d damage=%d after object frame, want 3/1", rt.Health, rt.DamageTaken)
	}
	id, _ := rt.At(PlayerLayer, 1, 17)
	if id != 19 {
		t.Fatalf("snake cell = %d, want raw19", id)
	}
}

func TestRuntimeInitializesSnakeDirectionFromBackground(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := rt.At(PlayerLayer, 7, 4)
	if id != 19 {
		t.Fatalf("stage snake = %d, want raw19", id)
	}
	if got := rt.ObjectState[rt.index(7, 4)] & 0x7; got != 2 {
		t.Fatalf("snake direction = %d, want decoded background direction 2", got)
	}
}

func TestRuntimeSnakeMovesUsingSourceDirection(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimeSnakes(rt)
	rt.SetForTest(PlayerLayer, 1, 17, 19)
	rt.SetForTest(PlayerLayer, 2, 17, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 17, EmptyRawID)
	rt.ObjectState[rt.index(1, 17)] = 2
	moved := rt.TickSnakes()
	if moved != 1 {
		t.Fatalf("moved snakes = %d, want 1", moved)
	}
	source, _ := rt.At(PlayerLayer, 1, 17)
	target, _ := rt.At(PlayerLayer, 2, 17)
	if source != EmptyRawID || target != 19 {
		t.Fatalf("snake move source=%d target=%d, want empty/raw19", source, target)
	}
	if got := rt.ObjectState[rt.index(2, 17)] & 0x7; got != 2 {
		t.Fatalf("moved snake direction = %d, want 2", got)
	}
	if got := rt.ObjectMotion[rt.index(2, 17)].Remaining; got != 18 {
		t.Fatalf("right-moving snake timer = %d, want 18 after same-row rescan", got)
	}
}

func TestRuntimeSnakeReversesWhenBlocked(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimeSnakes(rt)
	rt.SetForTest(PlayerLayer, 1, 17, 19)
	rt.SetForTest(PlayerLayer, 2, 17, 80)
	rt.ObjectState[rt.index(1, 17)] = 2
	moved := rt.TickSnakes()
	if moved != 0 {
		t.Fatalf("moved snakes = %d, want 0", moved)
	}
	id, _ := rt.At(PlayerLayer, 1, 17)
	if id != 19 {
		t.Fatalf("blocked snake cell = %d, want raw19", id)
	}
	state := rt.ObjectState[rt.index(1, 17)]
	if low, pending := state&0x7, (state&0x7000)>>12; low != 0 || pending != 4 {
		t.Fatalf("blocked snake state low=%d pending=%d, want 0/4", low, pending)
	}
	if got := rt.ObjectMotion[rt.index(1, 17)].Remaining; got != 21 {
		t.Fatalf("blocked snake turn timer = %d, want source 21", got)
	}
	rt.Player = Point{X: 5, Y: 17}
	rt.SetForTest(PlayerLayer, 2, 17, EmptyRawID)
	for tick := 0; tick < 7; tick++ {
		if moved := rt.TickSnakes(); moved != 0 {
			t.Fatalf("blocked snake moved during turn tick %d", tick+1)
		}
	}
	if moved := rt.TickSnakes(); moved != 1 {
		t.Fatalf("snake moved after turn timer = %d, want 1", moved)
	}
	target, _ := rt.At(PlayerLayer, 0, 17)
	if target != 19 {
		t.Fatalf("reversed snake target = %d, want raw19 at (0,17)", target)
	}
}

func TestRuntimeSnakeMoveDamagesPlayer(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimeSnakes(rt)
	rt.Player = Point{X: 2, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 19)
	rt.SetForTest(PlayerLayer, 2, 17, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 17, EmptyRawID)
	rt.ObjectState[rt.index(1, 17)] = 2
	moved := rt.TickSnakes()
	if moved != 1 {
		t.Fatalf("moved snakes = %d, want 1 while entering player cell", moved)
	}
	if rt.Health != 3 || rt.DamageTaken != 1 {
		t.Fatalf("health=%d damage=%d, want 3/1", rt.Health, rt.DamageTaken)
	}
	if rt.Player != (Point{X: 3, Y: 17}) || rt.PlayerMotion != (ObjectMotion{DX: 1, Remaining: playerMoveStartOffset}) {
		t.Fatalf("snake knockback player=%+v motion=%+v, want (3,17) and right jInt=18", rt.Player, rt.PlayerMotion)
	}
	id, _ := rt.At(PlayerLayer, 2, 17)
	if id != 19 {
		t.Fatalf("attacking snake target cell = %d, want raw19", id)
	}
}

func TestRuntimeCrawlerContactDamagesAndAllowsPlayerEntry(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 11)
	prepareCrawlerTestArea(rt, 10, 10)
	rt.Player = Point{X: 10, Y: 10}
	rt.SetForTest(PlayerLayer, 11, 10, 11)
	rt.SetForTest(PlayerLayer, 12, 10, 80)
	rt.SetForTest(PlayerLayer, 11, 11, 80)
	rt.ObjectState[rt.index(11, 10)] = 2
	if !rt.TryMove(1, 0) {
		t.Fatal("player did not enter source-passable crawler cell")
	}
	if rt.Player != (Point{X: 11, Y: 10}) {
		t.Fatalf("player = %+v, want crawler cell", rt.Player)
	}
	if rt.Health != 4 || rt.DamageTaken != 0 {
		t.Fatalf("health=%d damage=%d before object frame, want 4/0", rt.Health, rt.DamageTaken)
	}
	rt.TickCrawlers()
	if rt.Health != 3 || rt.DamageTaken != 1 {
		t.Fatalf("health=%d damage=%d after object frame, want 3/1", rt.Health, rt.DamageTaken)
	}
	id, _ := rt.At(PlayerLayer, 11, 10)
	if id != 11 {
		t.Fatalf("crawler cell after contact frame = %d, want raw11 at (11,10)", id)
	}
}

func TestRuntimeCrawlerInfersDirectionFromFloorAndMoves(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 11)
	prepareCrawlerTestArea(rt, 10, 10)
	rt.SetForTest(PlayerLayer, 10, 10, 11)
	rt.SetForTest(PlayerLayer, 10, 11, 80)
	rt.ObjectState[rt.index(10, 10)] = 0
	moved := rt.TickCrawlers()
	if moved != 0 {
		t.Fatalf("initial crawler update moved %d times, want only a direction change", moved)
	}
	if got := rt.ObjectState[rt.index(10, 10)] & objectDirectionMask; got != 2 {
		t.Fatalf("inferred crawler direction = %d, want right/2", got)
	}
	moved = rt.TickCrawlers()
	if moved != 1 {
		t.Fatalf("second crawler update moved %d times, want 1", moved)
	}
	source, _ := rt.At(PlayerLayer, 10, 10)
	target, _ := rt.At(PlayerLayer, 11, 10)
	if source != EmptyRawID || target != 11 {
		t.Fatalf("crawler move source=%d target=%d, want empty/raw11", source, target)
	}
	if got := rt.ObjectState[rt.index(11, 10)] & objectDirectionMask; got != 2 {
		t.Fatalf("crawler direction = %d, want inferred right/2", got)
	}
	if got := rt.ObjectMotion[rt.index(11, 10)].Remaining; got != 13 {
		t.Fatalf("right-moving crawler timer = %d, want source same-scan 13", got)
	}
}

func TestRuntimeCrawlerDamagesPlayerOnTick(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 11)
	prepareCrawlerTestArea(rt, 10, 10)
	rt.Player = Point{X: 11, Y: 10}
	rt.SetForTest(PlayerLayer, 10, 10, 11)
	rt.SetForTest(PlayerLayer, 10, 11, 80)
	rt.ObjectState[rt.index(10, 10)] = 2
	moved := rt.TickCrawlers()
	if moved != 1 {
		t.Fatalf("moved crawlers = %d, want 1 while entering player cell", moved)
	}
	if rt.Health != 3 || rt.DamageTaken != 1 {
		t.Fatalf("health=%d damage=%d, want 3/1", rt.Health, rt.DamageTaken)
	}
	id, _ := rt.At(PlayerLayer, 11, 10)
	if id != 11 {
		t.Fatalf("attacking crawler target cell = %d, want raw11", id)
	}
}

func TestRuntimeHurtStatePreventsRepeatedContactDamage(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if !rt.Hurt(1) {
		t.Fatal("first hurt was rejected")
	}
	if rt.Hurt(1) {
		t.Fatal("second hurt was accepted during source hurt state")
	}
	if rt.TryMove(1, 0) {
		t.Fatal("player moved during the eight-frame hurt animation")
	}
	for i := 0; i < hurtStateDuration; i++ {
		rt.TickStatus()
	}
	if rt.HurtTicks != 0 || rt.InvulnerabilityTicks != hurtInvulnerabilityDuration-hurtStateDuration {
		t.Fatalf("post-animation timers hurt=%d invulnerability=%d", rt.HurtTicks, rt.InvulnerabilityTicks)
	}
	if rt.Hurt(1) {
		t.Fatal("hurt was accepted during the source post-animation protection")
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("player remained input-locked during post-animation invulnerability")
	}
	for rt.InvulnerabilityTicks > 0 {
		rt.TickStatus()
	}
	if !rt.Hurt(1) {
		t.Fatal("hurt was rejected after source invulnerability expired")
	}
	if rt.Health != 2 || rt.DamageTaken != 2 {
		t.Fatalf("health=%d damage=%d, want 2/2", rt.Health, rt.DamageTaken)
	}
	if rt.HitCount != 2 {
		t.Fatalf("hit count = %d, want 2 accepted hits", rt.HitCount)
	}
}

func TestRuntimeCrawlerTurnsAwayWhenForwardAndWallSideAreBlocked(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 11)
	prepareCrawlerTestArea(rt, 10, 10)
	rt.SetForTest(PlayerLayer, 10, 10, 11)
	rt.SetForTest(PlayerLayer, 11, 10, 80)
	rt.SetForTest(PlayerLayer, 10, 11, 80)
	rt.ObjectState[rt.index(10, 10)] = 2
	moved := rt.TickCrawlers()
	if moved != 0 {
		t.Fatalf("moved crawlers = %d, want 0", moved)
	}
	id, _ := rt.At(PlayerLayer, 10, 10)
	if id != 11 {
		t.Fatalf("blocked crawler cell = %d, want raw11", id)
	}
	if got := rt.ObjectState[rt.index(10, 10)] & objectDirectionMask; got != 1 {
		t.Fatalf("blocked crawler direction = %d, want turn away/up 1", got)
	}
}

func TestRuntimeCrawlerTurnsAroundSourceOuterCorner(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 11)
	prepareCrawlerTestArea(rt, 10, 10)
	rt.SetForTest(PlayerLayer, 10, 10, 11)
	// The source's third probe is the diagonal cell behind the wall side.
	rt.SetForTest(PlayerLayer, 9, 11, 80)
	rt.ObjectState[rt.index(10, 10)] = 2
	if moved := rt.TickCrawlers(); moved != 1 {
		t.Fatalf("outer-corner moves = %d, want 1", moved)
	}
	id, _ := rt.At(PlayerLayer, 10, 11)
	if id != 11 {
		t.Fatalf("outer-corner target = %d, want raw11 below", id)
	}
	if got := rt.ObjectState[rt.index(10, 11)] & objectDirectionMask; got != 3 {
		t.Fatalf("outer-corner direction = %d, want down/3", got)
	}
}

func TestRuntimeCrawlerUsesSourceFivePixelMotionCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 11)
	prepareCrawlerTestArea(rt, 10, 10)
	rt.SetForTest(PlayerLayer, 10, 10, 11)
	rt.SetForTest(PlayerLayer, 10, 11, 80)
	rt.ObjectState[rt.index(10, 10)] = 2
	rt.ObjectMotion[rt.index(10, 10)] = ObjectMotion{Remaining: 18}
	for tick, want := range []int{13, 8, 3, -2} {
		if moved := rt.TickCrawlers(); moved != 0 {
			t.Fatalf("crawler moved early on cadence tick %d", tick+1)
		}
		if got := rt.ObjectMotion[rt.index(10, 10)].Remaining; got != want {
			t.Fatalf("crawler timer on tick %d = %d, want %d", tick+1, got, want)
		}
	}
	if moved := rt.TickCrawlers(); moved != 1 {
		t.Fatalf("crawler moves after timer -2 = %d, want 1", moved)
	}
}

func prepareCrawlerTestArea(rt *Runtime, centerX, centerY int) {
	for y := centerY - 2; y <= centerY+2; y++ {
		for x := centerX - 2; x <= centerX+2; x++ {
			rt.SetForTest(PlayerLayer, x, y, EmptyRawID)
			rt.SetForTest(ForegroundLayer, x, y, EmptyRawID)
		}
	}
}

func TestRuntimeRightHorizontalHazardDamagesWithinReach(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 22, 23)
	rt.Player = Point{X: 4, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 22)
	rt.ObjectState[rt.index(1, 17)] = 3
	hits := rt.TickHorizontalHazards()
	if hits != 1 {
		t.Fatalf("hazard hits = %d, want 1", hits)
	}
	if rt.Health != 3 || rt.DamageTaken != 1 {
		t.Fatalf("health=%d damage=%d, want 3/1", rt.Health, rt.DamageTaken)
	}
	if got := rt.ObjectState[rt.index(1, 17)] & 0x3; got != 0 {
		t.Fatalf("hazard phase = %d, want wrapped 0", got)
	}
}

func TestRuntimeLeftHorizontalHazardDamagesWithinReach(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 22, 23)
	rt.Player = Point{X: 1, Y: 17}
	rt.SetForTest(PlayerLayer, 4, 17, 23)
	rt.ObjectState[rt.index(4, 17)] = 3
	hits := rt.TickHorizontalHazards()
	if hits != 1 {
		t.Fatalf("hazard hits = %d, want 1", hits)
	}
	if rt.Health != 3 || rt.DamageTaken != 1 {
		t.Fatalf("health=%d damage=%d, want 3/1", rt.Health, rt.DamageTaken)
	}
}

func TestRuntimeHorizontalHazardDoesNotDamageBeyondCurrentReach(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 22, 23)
	rt.Player = Point{X: 4, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 22)
	rt.ObjectState[rt.index(1, 17)] = 2
	hits := rt.TickHorizontalHazards()
	if hits != 0 {
		t.Fatalf("hazard hits = %d, want 0", hits)
	}
	if rt.Health != 4 || rt.DamageTaken != 0 {
		t.Fatalf("health=%d damage=%d, want 4/0", rt.Health, rt.DamageTaken)
	}
	if got := rt.ObjectState[rt.index(1, 17)] & 0x3; got != 3 {
		t.Fatalf("hazard phase = %d, want 3", got)
	}
}

func TestRuntimeGravityDoesNotEnterGoalMarkers(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 11, 5)
	rt.SetForTest(PlayerLayer, 1, 10, 80)
	rt.SetForTest(PlayerLayer, 3, 10, 80)
	moved := rt.TickGravity()
	if moved != 0 {
		t.Fatalf("gravity moved %d boulders into/around goal marker, want 0", moved)
	}
	source, _ := rt.At(PlayerLayer, 2, 10)
	target, _ := rt.At(PlayerLayer, 2, 11)
	if source != 0 || target != EmptyRawID {
		t.Fatalf("goal-blocked gravity source=%d target=%d, want raw0/empty", source, target)
	}
}

func TestRuntimeBoulderCrushesSnake(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, 43)
	rt.SetForTest(ForegroundLayer, 2, 11, EmptyRawID)
	moved := tickGravityObject(rt, 2, 10)
	if moved != 1 {
		t.Fatalf("gravity moved %d boulders, want 1", moved)
	}
	source, _ := rt.At(PlayerLayer, 2, 10)
	target, _ := rt.At(PlayerLayer, 2, 11)
	if source != EmptyRawID || target != 0 {
		t.Fatalf("snake crush source=%d target=%d, want empty/raw0", source, target)
	}
}

func TestRuntimeBoulderCrushesCrawler(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 11)
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, 11)
	rt.SetForTest(ForegroundLayer, 2, 11, EmptyRawID)
	moved := tickGravityObject(rt, 2, 10)
	if moved != 1 {
		t.Fatalf("gravity moved %d boulders, want 1", moved)
	}
	source, _ := rt.At(PlayerLayer, 2, 10)
	target, _ := rt.At(PlayerLayer, 2, 11)
	if source != EmptyRawID || target != 0 {
		t.Fatalf("crawler crush source=%d target=%d, want empty/raw0", source, target)
	}
}

func TestRuntimeInitializesEnemyGateCountersFromRaw17(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	marker, _ := rt.At(ForegroundLayer, 17, 21)
	if marker != EmptyRawID {
		t.Fatalf("enemy raw17 marker = %d, want cleared after init", marker)
	}
	if got := rt.EnemyGateGroup[rt.index(17, 20)]; got != 0 {
		t.Fatalf("snake gate group = %d, want 0", got)
	}
	if rt.EnemyGateCounters[0] == 0 {
		t.Fatal("enemy gate counter group 0 = 0, want initialized from raw17 enemy markers")
	}
}

func TestRuntimeEnemyGateOpensDoorWhenGroupedSnakeCrushed(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.EnemyGateCounters[0] = 1
	rt.Player = Point{X: 0, Y: 10}
	rt.SetForTest(PlayerLayer, 1, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 1, 10, 26)
	rt.SetForTest(BackgroundLayer, 1, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, 19)
	rt.SetForTest(ForegroundLayer, 2, 11, EmptyRawID)
	rt.EnemyGateGroup[rt.index(2, 11)] = 0
	rt.SetForTest(PlayerLayer, 4, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 4, 10, 7)
	rt.SetForTest(BackgroundLayer, 4, 10, 0)
	rt.SetForTest(ForegroundLayer, 4, 11, 17)
	rt.SetForTest(BackgroundLayer, 4, 11, 0)
	rt.SetForTest(ForegroundLayer, 5, 10, 33)
	rt.SetForTest(BackgroundLayer, 5, 10, EmptyRawID)
	rt.ContainerLocked[rt.index(5, 10)] = true
	rt.SetForTest(ForegroundLayer, 5, 11, 17)
	rt.SetForTest(BackgroundLayer, 5, 11, 0)
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to step onto raw26 enemy gate trigger")
	}
	rt.tickPendingForegroundEvent()
	rt.AdvancePlayerMotion()
	rt.tickPendingForegroundEvent()
	rt.AdvancePlayerMotion()
	rt.tickPendingForegroundEvent()
	if rt.ActiveEnemyGateGroup != 0 {
		t.Fatalf("active enemy gate group = %d, want 0", rt.ActiveEnemyGateGroup)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundRiddle {
		t.Fatalf("enemy-gate trigger sounds=%v, want [%d]", events, SoundRiddle)
	}
	trigger, _ := rt.At(ForegroundLayer, 1, 10)
	if trigger != EmptyRawID {
		t.Fatalf("raw26 trigger foreground = %d, want cleared", trigger)
	}
	if moved := tickGravityObject(rt, 2, 10); moved != 1 {
		t.Fatalf("gravity moved %d boulders, want 1 crushing grouped snake", moved)
	}
	if rt.EnemyGateCounters[0] != 0 {
		t.Fatalf("enemy gate counter = %d, want 0", rt.EnemyGateCounters[0])
	}
	events := rt.DrainSoundEvents()
	foundDoorSound := false
	for _, event := range events {
		foundDoorSound = foundDoorSound || event == SoundDoor
	}
	if !foundDoorSound {
		t.Fatalf("enemy-gate completion sounds=%v, want door sound %d", events, SoundDoor)
	}
	state, _ := rt.At(BackgroundLayer, 4, 10)
	if state != 0x10 || rt.IsPassable(4, 10) {
		t.Fatalf("opening door state=%#x passable=%v, want phase1/blocking", state, rt.IsPassable(4, 10))
	}
	rt.tickDoorAnimations(3)
	state, _ = rt.At(BackgroundLayer, 4, 10)
	if state != 0x20 || !rt.IsPassable(4, 10) {
		t.Fatalf("opened door state=%#x passable=%v, want phase2/passable", state, rt.IsPassable(4, 10))
	}
	overlay, _ := rt.At(ForegroundLayer, 5, 10)
	overlayState, _ := rt.At(BackgroundLayer, 5, 10)
	if overlay != 33 || overlayState != EmptyRawID || rt.ContainerLockedAt(5, 10) {
		t.Fatalf("raw17 overlay side effect foreground=%d payload-state=%d locked=%v, want raw33/unchanged/unlocked", overlay, overlayState, rt.ContainerLockedAt(5, 10))
	}
}

func TestRuntimeHammerBreaksRaw30Wall(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 1
	rt.Player = Point{X: 22, Y: 4}
	rt.SetForTest(PlayerLayer, 23, 4, 30)
	for _, pt := range []Point{{X: 23, Y: 3}, {X: 24, Y: 4}, {X: 23, Y: 5}} {
		rt.SetForTest(PlayerLayer, pt.X, pt.Y, EmptyRawID)
	}
	if !rt.UseHammer(1, 0) {
		t.Fatal("UseHammer() = false, want true against raw30")
	}
	for sourceTick := 1; sourceTick <= 24; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
		rt.TickBreakables()
	}
	id, _ := rt.At(PlayerLayer, 23, 4)
	if id != EmptyRawID {
		t.Fatalf("breakable wall = %d, want empty", id)
	}
	if rt.BreakableWalls != 1 {
		t.Fatalf("breakable walls = %d, want 1", rt.BreakableWalls)
	}
}

func TestRuntimeHammerRequiresToolLevel(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 22, Y: 4}
	rt.SetForTest(PlayerLayer, 23, 4, 30)
	if rt.UseHammer(1, 0) {
		t.Fatal("UseHammer without special tool = true, want false")
	}
	rt.SpecialItemMask = 1
	if !rt.UseHammer(1, 0) {
		t.Fatal("UseHammer with raw24 tool level = false, want true")
	}
}

func TestRuntimeHammerBlockedByBoulderUsesSourceFeedback(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 1
	rt.Player = Point{X: 8, Y: 17}
	rt.SetForTest(PlayerLayer, 9, 17, 0)
	if !rt.UseHammer(1, 0) {
		t.Fatal("blocked boulder did not start the source hammer action")
	}
	for sourceTick := 1; sourceTick <= hammerImpactTick; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundHammerBlock {
		t.Fatalf("blocked hammer sounds=%v, want [%d]", events, SoundHammerBlock)
	}
	if id, _ := rt.At(PlayerLayer, 9, 17); id != 0 {
		t.Fatalf("blocked hammer changed boulder to raw %d", id)
	}
}

func TestRuntimeHammerUsesSourceActionTiming(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 1
	rt.Player = Point{X: 22, Y: 4}
	rt.SetForTest(PlayerLayer, 23, 4, 30)
	if !rt.UseHammer(1, 0) {
		t.Fatal("UseHammer() = false, want source action")
	}
	if !rt.Hammering || rt.HammerAnimation != 14 || rt.HammerTicks != 0 || rt.CanAcceptInput() {
		t.Fatalf("started hammer action active=%v animation=%d ticks=%d input=%v, want true/14/0/false", rt.Hammering, rt.HammerAnimation, rt.HammerTicks, rt.CanAcceptInput())
	}
	for sourceTick := 1; sourceTick < hammerImpactTick; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
		if state := rt.objectStateAtForTest(23, 4); state != 0 {
			t.Fatalf("wall state before impact tick %d=%d, want zero", sourceTick, state)
		}
	}
	rt.TickSourceFrame(8, hammerImpactTick, 0)
	if state := rt.objectStateAtForTest(23, 4); state != 1 {
		t.Fatalf("wall state on impact tick=%d, want 1", state)
	}
	for sourceTick := hammerImpactTick + 1; sourceTick < hammerOtherDuration; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if !rt.Hammering || rt.HammerTicks != hammerOtherDuration-1 {
		t.Fatalf("hammer before final tick active=%v ticks=%d, want true/%d", rt.Hammering, rt.HammerTicks, hammerOtherDuration-1)
	}
	rt.TickSourceFrame(8, hammerOtherDuration, 0)
	if rt.Hammering || !rt.CanAcceptInput() {
		t.Fatalf("completed hammer active=%v input=%v, want false/true", rt.Hammering, rt.CanAcceptInput())
	}
}

func TestRuntimeHammerStunsGreenSnakeAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage03.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 1
	rt.Player = Point{X: 33, Y: 5}
	if !rt.UseHammer(1, 0) {
		t.Fatal("UseHammer() = false, want true against adjacent raw19 snake")
	}
	for sourceTick := 1; sourceTick <= hammerImpactTick; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if id, _ := rt.At(PlayerLayer, 34, 5); id != 19 {
		t.Fatalf("hammered green snake raw=%d, want 19 while stunned", id)
	}
	if state := rt.objectStateAtForTest(34, 5); state&snakeStunMask != snakeStunDuration {
		t.Fatalf("hammered green snake state=%#x, want stun=%#x", state, snakeStunDuration)
	}
	for sourceTick := hammerImpactTick + 1; sourceTick < 68; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if state := rt.objectStateAtForTest(34, 5); state&snakeStunMask == 0 {
		t.Fatalf("green snake recovered before source tick 68: state=%#x", state)
	}
	rt.TickSourceFrame(8, 68, 0)
	if state := rt.objectStateAtForTest(34, 5); state&snakeStunMask != 0 {
		t.Fatalf("green snake stun at source tick 68=%#x, want zero", state&snakeStunMask)
	}
	if rt.Health != rt.MaxHealth {
		t.Fatalf("health while green snake stunned=%d, want %d", rt.Health, rt.MaxHealth)
	}
}

func TestRuntimeHammerDefeatsUnmarkedRedSnake(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage03.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 1
	rt.Player = Point{X: 31, Y: 8}
	rt.ObjectState[rt.index(32, 8)] = 0
	if !rt.UseHammer(1, 0) {
		t.Fatal("UseHammer() = false, want true against adjacent raw43 snake")
	}
	for sourceTick := 1; sourceTick <= hammerImpactTick; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if id, _ := rt.At(PlayerLayer, 32, 8); id != EmptyRawID {
		t.Fatalf("hammered unmarked red snake raw=%d, want empty", id)
	}
}

func TestRuntimeAuthoredRedSnakeUsesSourceThreeHammerHits(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage03.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	point := Point{X: 32, Y: 8}
	idx := rt.index(point.X, point.Y)
	if got := rt.ObjectState[idx] & 0x18000; got != 0x10000 {
		t.Fatalf("authored red-snake durability=%#x, want source %#x", got, 0x10000)
	}
	for hit, wantDurability := range []int{0x8000, 0, 0} {
		rt.HammerTarget = point
		if !rt.applyHammerImpact() {
			t.Fatalf("red-snake hammer hit %d was not handled", hit+1)
		}
		if hit < 2 {
			if id, _ := rt.At(PlayerLayer, point.X, point.Y); id != 43 {
				t.Fatalf("red snake removed on hit %d: raw=%d", hit+1, id)
			}
			if got := rt.ObjectState[idx] & 0x18000; got != wantDurability {
				t.Fatalf("red-snake durability after hit %d=%#x, want %#x", hit+1, got, wantDurability)
			}
			rt.ObjectState[idx] &^= snakeStunMask
		} else if id, _ := rt.At(PlayerLayer, point.X, point.Y); id != EmptyRawID {
			t.Fatalf("red snake raw=%d after third hit, want empty", id)
		}
	}
}

func TestRuntimeRedSnakeStopsSourceChaseAfterBitingHero(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage03.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 19, 43)
	point := Point{X: 10, Y: 10}
	rt.Player = point
	rt.SetForTest(PlayerLayer, point.X, point.Y, 43)
	rt.SetForTest(ForegroundLayer, point.X, point.Y, EmptyRawID)
	rt.ObjectState[rt.index(point.X, point.Y)] = 0xc00 | 0x10000
	rt.TickSnakes()
	if rt.Health != rt.MaxHealth-1 {
		t.Fatalf("red-snake chase bite health=%d, want %d", rt.Health, rt.MaxHealth-1)
	}
	if got := rt.ObjectState[rt.index(point.X, point.Y)] & 0xf00; got != 0 {
		t.Fatalf("red-snake chase state after bite=%#x, want cleared", got)
	}
}

func TestRuntimeBreakableWallDamagePropagatesToAdjacentRaw30(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 1
	rt.Player = Point{X: 18, Y: 5}
	rt.SetForTest(PlayerLayer, 19, 5, 30)
	rt.SetForTest(PlayerLayer, 20, 5, 30)
	for _, pt := range []Point{{X: 19, Y: 4}, {X: 20, Y: 4}, {X: 18, Y: 5}, {X: 21, Y: 5}, {X: 19, Y: 6}, {X: 20, Y: 6}} {
		rt.SetForTest(PlayerLayer, pt.X, pt.Y, EmptyRawID)
	}
	if !rt.UseHammer(1, 0) {
		t.Fatal("UseHammer() = false, want true against raw30 cluster")
	}
	for sourceTick := 1; sourceTick <= 24; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
		rt.TickBreakables()
	}
	left, _ := rt.At(PlayerLayer, 19, 5)
	right, _ := rt.At(PlayerLayer, 20, 5)
	if left != EmptyRawID || right != EmptyRawID {
		t.Fatalf("breakable cluster left=%d right=%d, want empty/empty", left, right)
	}
	if rt.BreakableWalls != 2 {
		t.Fatalf("breakable walls = %d, want 2", rt.BreakableWalls)
	}
}

func TestRuntimeBreakableChainUsesSourceBottomUpLeftRightScan(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	for _, point := range []Point{{X: 10, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 9}} {
		rt.SetForTest(PlayerLayer, point.X, point.Y, 30)
		rt.ObjectState[rt.index(point.X, point.Y)] = 3
	}
	rt.ObjectState[rt.index(10, 10)] = 4
	rt.TickBreakables()
	if got := rt.ObjectState[rt.index(11, 10)]; got != 5 {
		t.Fatalf("right wall state after same-frame propagation = %d, want 5", got)
	}
	if got := rt.ObjectState[rt.index(10, 9)]; got != 5 {
		t.Fatalf("upper wall state after bottom-up propagation = %d, want 5", got)
	}
}

func TestRuntimeBoulderDoesNotRollThroughPlayerSideCell(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1)
	rt.Player = Point{X: 1, Y: 10}
	rt.SetForTest(PlayerLayer, 2, 10, 0)
	rt.SetForTest(PlayerLayer, 2, 11, 80)
	rt.SetForTest(PlayerLayer, 1, 10, EmptyRawID)
	rt.SetForTest(PlayerLayer, 1, 11, EmptyRawID)
	rt.SetForTest(PlayerLayer, 3, 10, 80)
	rt.SetForTest(ForegroundLayer, 1, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 1, 11, EmptyRawID)
	moved := rt.TickGravity()
	if moved != 0 {
		t.Fatalf("gravity moved %d boulders through player side cell, want 0", moved)
	}
	source, _ := rt.At(PlayerLayer, 2, 10)
	if source != 0 {
		t.Fatalf("side-blocked boulder source = %d, want raw0", source)
	}
}

func TestRuntimeCollectsVioletGemCandidate(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if !rt.SetForTest(PlayerLayer, 1, 17, 1) {
		t.Fatal("failed to place raw 1")
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("move into raw 1 failed")
	}
	rt.TickSourceFrame(8, 1, 0)
	if rt.VioletGems != 1 {
		t.Fatalf("violet gems = %d, want 1", rt.VioletGems)
	}
	if rt.BonusRemaining != 9 {
		t.Fatalf("bonus remaining after violet = %d, want 9", rt.BonusRemaining)
	}
	id, _ := rt.At(PlayerLayer, 1, 17)
	if id != EmptyRawID {
		t.Fatalf("collected cell = %d, want empty", id)
	}
}

func TestRuntimeCollectsRaw2RedDiamond(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if !rt.SetForTest(PlayerLayer, 1, 17, 2) {
		t.Fatal("failed to place raw 2")
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("move into raw 2 failed")
	}
	if rt.RedDiamonds != 1 {
		t.Fatalf("red diamonds = %d, want 1", rt.RedDiamonds)
	}
	id, _ := rt.At(PlayerLayer, 1, 17)
	if id != EmptyRawID {
		t.Fatalf("collected raw2 cell = %d, want empty", id)
	}
}

func TestRuntimeHiddenRedDiamondOpensRaw33Overlay(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 18, Y: 2}
	if id, _ := rt.At(PlayerLayer, 19, 2); id != 2 {
		t.Fatalf("hidden pickup = %d, want raw2", id)
	}
	if id, _ := rt.At(ForegroundLayer, 19, 2); id != 33 {
		t.Fatalf("hidden pickup overlay = %d, want raw33", id)
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter hidden red-diamond cell")
	}
	idx := rt.index(19, 2)
	if rt.RedDiamonds != 0 || rt.ObjectState[idx] != 0 || !rt.pendingChestSet || rt.CanAcceptInput() {
		t.Fatalf("moving into chest red=%d overlay=%d pending=%v input=%v, want 0/0/true/false", rt.RedDiamonds, rt.ObjectState[idx], rt.pendingChestSet, rt.CanAcceptInput())
	}
	sourceTick := 0
	for rt.PlayerMotion.Remaining > 0 {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
		rt.AdvancePlayerMotion()
	}
	if rt.ChestOpening {
		t.Fatal("chest opened before the settled-player source scan")
	}
	sourceTick++
	rt.TickSourceFrame(8, sourceTick, 0)
	chestStartTick := sourceTick
	if rt.RedDiamonds != 0 || rt.ObjectState[idx] != 1 || !rt.ChestOpening || rt.ChestTicks != 0 {
		t.Fatalf("started chest red=%d overlay=%d opening=%v ticks=%d, want 0/1/true/0", rt.RedDiamonds, rt.ObjectState[idx], rt.ChestOpening, rt.ChestTicks)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundChestOpen {
		t.Fatalf("chest start sounds=%v, want [%d]", events, SoundChestOpen)
	}
	for rt.ObjectState[idx] < 3 && sourceTick < 8 {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if rt.ObjectState[idx] != 3 {
		t.Fatalf("opened raw33 overlay state = %d, want final frame 3", rt.ObjectState[idx])
	}
	for rt.ChestTicks < chestRewardTick-1 {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if rt.RedDiamonds != 0 {
		t.Fatalf("red diamonds before source reward frame = %d, want 0", rt.RedDiamonds)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundChestReward {
		t.Fatalf("chest sequence-12 sounds=%v, want [%d] at tick %d", events, SoundChestReward, chestRewardSoundTick)
	}
	sourceTick++
	rt.TickSourceFrame(8, sourceTick, 0)
	if rt.RedDiamonds != 1 || !rt.ChestRewarded || !rt.ChestOpening {
		t.Fatalf("reward frame red=%d rewarded=%v opening=%v, want 1/true/true", rt.RedDiamonds, rt.ChestRewarded, rt.ChestOpening)
	}
	for rt.ChestOpening {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if sourceTick != chestStartTick+chestOpenDuration {
		t.Fatalf("chest source completion tick = %d, want %d", sourceTick, chestStartTick+chestOpenDuration)
	}
	if rt.ChestOpening || !rt.CanAcceptInput() {
		t.Fatalf("completed chest opening=%v input=%v, want false/true", rt.ChestOpening, rt.CanAcceptInput())
	}
}

func TestRuntimeStage01HealthChestRewardsOnSourceFrame(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage01.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 20, Y: 15}
	rt.PlayerMotion = ObjectMotion{}
	rt.Health = 2
	if id, _ := rt.At(PlayerLayer, 21, 15); id != 7 {
		t.Fatalf("Stage 2 health chest reward=%d, want raw7", id)
	}
	if foreground, _ := rt.At(ForegroundLayer, 21, 15); foreground != 33 {
		t.Fatalf("Stage 2 health chest foreground=%d, want raw33", foreground)
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter Stage 2 health chest")
	}
	if rt.Health != 2 || !rt.pendingChestSet {
		t.Fatalf("entering health chest health=%d pending=%v, want 2/true before settle", rt.Health, rt.pendingChestSet)
	}
	settleRuntimePlayerMotion(rt)
	rt.TickSourceFrame(8, 1, 0)
	if !rt.ChestOpening || rt.ChestAnimation != 40 || rt.ChestRewardID != 7 || rt.ChestTicks != 0 {
		t.Fatalf("health chest opening=%v animation=%d reward=%d ticks=%d, want true/40/7/0", rt.ChestOpening, rt.ChestAnimation, rt.ChestRewardID, rt.ChestTicks)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundChestOpen {
		t.Fatalf("health chest start sounds=%v, want [%d]", events, SoundChestOpen)
	}
	for sourceTick := 2; sourceTick <= chestRewardTick; sourceTick++ {
		rt.gravitySourceTick = sourceTick
		rt.TickStatus()
	}
	if rt.Health != 2 || rt.HealthRefills != 0 {
		t.Fatalf("health chest rewarded early: health=%d refills=%d", rt.Health, rt.HealthRefills)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundChestReward {
		t.Fatalf("health chest reward sounds=%v, want [%d]", events, SoundChestReward)
	}
	rt.gravitySourceTick = chestRewardTick + 1
	rt.TickStatus()
	if rt.Health != rt.MaxHealth || rt.HealthRefills != 1 || !rt.ChestRewarded {
		t.Fatalf("health chest reward health=%d refills=%d rewarded=%v, want %d/1/true", rt.Health, rt.HealthRefills, rt.ChestRewarded, rt.MaxHealth)
	}
}

func TestRuntimeStage01FullHealthChestBecomesTenPointBonus(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage01.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 20, Y: 15}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter full-health Stage 2 chest")
	}
	settleRuntimePlayerMotion(rt)
	rt.TickSourceFrame(8, 1, 0)
	if rt.ChestRewardID != 41 || rt.ChestRewardValue != 10 {
		t.Fatalf("full-health chest reward=%d value=%d, want source raw41/10 before animation", rt.ChestRewardID, rt.ChestRewardValue)
	}
	for sourceTick := 2; sourceTick <= chestRewardTick+1; sourceTick++ {
		rt.gravitySourceTick = sourceTick
		rt.TickStatus()
	}
	if rt.HealthRefills != 0 || rt.BonusValue != 10 || rt.BonusPickups != 1 || rt.BonusRemaining != 5 {
		t.Fatalf("full-health chest refills=%d bonus=%d pickups=%d remaining=%d, want 0/10/1/5", rt.HealthRefills, rt.BonusValue, rt.BonusPickups, rt.BonusRemaining)
	}
	if rt.VioletGems != 10 || rt.TotalVioletGems != 29 {
		t.Fatalf("full-health chest violet=%d/%d, want source conversion 10/29", rt.VioletGems, rt.TotalVioletGems)
	}
}

func TestRuntimeStage01SecondChestUsesSourceAnimation48(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage01.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 20, Y: 15}
	rt.PlayerMotion = ObjectMotion{}
	rt.Health = 2
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter first Stage 2 health chest")
	}
	settleRuntimePlayerMotion(rt)
	rt.TickSourceFrame(8, 1, 0)
	rt.DrainSoundEvents()
	for sourceTick := 2; sourceTick <= chestOpenDuration+1; sourceTick++ {
		rt.gravitySourceTick = sourceTick
		rt.TickStatus()
	}
	if rt.ChestOpening || !rt.lastPickupTickSet || rt.lastPickupTick != chestOpenDuration+1 {
		t.Fatalf("first chest completion opening=%v lastSet=%v lastTick=%d", rt.ChestOpening, rt.lastPickupTickSet, rt.lastPickupTick)
	}

	rt.Player = Point{X: 5, Y: 18}
	rt.PlayerMotion = ObjectMotion{}
	rt.Health = 2
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter second Stage 2 health chest")
	}
	settleRuntimePlayerMotion(rt)
	startTick := chestOpenDuration + 2
	rt.TickSourceFrame(8, startTick, 0)
	if rt.ChestAnimation != 48 || rt.ChestRewardID != 7 {
		t.Fatalf("second chest animation=%d reward=%d, want 48/7", rt.ChestAnimation, rt.ChestRewardID)
	}
	rt.DrainSoundEvents()
	for sourceTick := startTick + 1; sourceTick < startTick+chestShortRewardTick; sourceTick++ {
		rt.gravitySourceTick = sourceTick
		rt.TickStatus()
	}
	if rt.ChestRewarded || rt.Health != 2 {
		t.Fatalf("short chest rewarded before tick %d: rewarded=%v health=%d", chestShortRewardTick, rt.ChestRewarded, rt.Health)
	}
	rt.gravitySourceTick = startTick + chestShortRewardTick
	rt.TickStatus()
	if !rt.ChestRewarded || rt.Health != rt.MaxHealth || rt.ChestTicks != chestShortRewardTick {
		t.Fatalf("short chest reward rewarded=%v health=%d ticks=%d, want true/%d/%d", rt.ChestRewarded, rt.Health, rt.ChestTicks, rt.MaxHealth, chestShortRewardTick)
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundChestReward {
		t.Fatalf("short chest reward sounds=%v, want [%d]", events, SoundChestReward)
	}
}

func TestRuntimeMaxLifeChestFallsBackToHealthThenBonus(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	point := Point{X: 19, Y: 6}
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = point
	rt.ExtraLives = 99
	rt.Health = 2
	rt.startChestOpening(point, true)
	if rt.ChestRewardID != 7 {
		t.Fatalf("max-life chest reward=%d, want health raw7", rt.ChestRewardID)
	}
	rt.applyChestReward()
	if rt.Health != rt.MaxHealth || len(rt.PersistentRewardCoordinates()) != 0 {
		t.Fatalf("max-life health fallback health=%d coordinates=%v, want full/replayable", rt.Health, rt.PersistentRewardCoordinates())
	}

	rt, err = NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = point
	rt.ExtraLives = 99
	rt.Health = rt.MaxHealth
	rt.startChestOpening(point, true)
	if rt.ChestRewardID != 41 || rt.ChestRewardValue != 10 {
		t.Fatalf("max-life/full-health chest reward=%d/%d, want raw41/10", rt.ChestRewardID, rt.ChestRewardValue)
	}
}

func TestRuntimeRaw33OverlayIsPassableAndPersistent(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage04.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 28, Y: 7}
	rt.SetForTest(PlayerLayer, 29, 7, 33)
	rt.SetForTest(BackgroundLayer, 29, 7, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 29, 7, 7)
	if !rt.TryMove(1, 0) {
		t.Fatal("move onto raw33 overlay failed")
	}
	if rt.Player != (Point{X: 29, Y: 7}) {
		t.Fatalf("player = %+v, want raw33 cell", rt.Player)
	}
	id, _ := rt.At(PlayerLayer, 29, 7)
	if id != 33 {
		t.Fatalf("raw33 overlay cell = %d, want preserved raw33", id)
	}
}

func TestRuntimeForegroundRaw7DoorBlocksUntilStateOpen(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 6, Y: 18}
	rt.SetForTest(PlayerLayer, 7, 18, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 7, 18, 7)
	rt.SetForTest(BackgroundLayer, 7, 18, 0)
	if rt.TryMove(1, 0) {
		t.Fatal("moved through closed foreground raw7 door")
	}
	if rt.Player != (Point{X: 6, Y: 18}) {
		t.Fatalf("player = %+v, want blocked before door", rt.Player)
	}
	rt.SetForTest(BackgroundLayer, 7, 18, EmptyRawID)
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to move through open foreground raw7 door")
	}
	if rt.Player != (Point{X: 7, Y: 18}) {
		t.Fatalf("player = %+v, want open door cell", rt.Player)
	}
}

func TestRuntimeForegroundRaw6PressureSwitchControlsDoor(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 1, 17, 6)
	rt.SetForTest(BackgroundLayer, 1, 17, 0)
	rt.SetForTest(PlayerLayer, 2, 17, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 2, 17, EmptyRawID)
	rt.SetForTest(PlayerLayer, 4, 17, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 4, 17, 7)
	// The source post-load pass stores one linked activator in the low nibble.
	rt.SetForTest(BackgroundLayer, 4, 17, 1)
	rt.DoorGroup[rt.index(4, 17)] = 0
	rt.updatePressureDoors()
	if rt.IsPassable(4, 17) {
		t.Fatal("raw6-linked door is passable while switch is inactive")
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to step onto raw6 pressure switch")
	}
	settleRuntimePlayerMotion(rt)
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundDoor {
		t.Fatalf("pressure-door opening sounds=%v, want [%d]", events, SoundDoor)
	}
	state, _ := rt.At(BackgroundLayer, 4, 17)
	if state != 0x11 || rt.IsPassable(4, 17) {
		t.Fatalf("pressure door phase1 state=%#x passable=%v, want blocking", state, rt.IsPassable(4, 17))
	}
	rt.tickDoorAnimations(3)
	state, _ = rt.At(BackgroundLayer, 4, 17)
	if state != 0x21 || !rt.IsPassable(4, 17) {
		t.Fatalf("pressure door phase2 state=%#x passable=%v, want passable", state, rt.IsPassable(4, 17))
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to step off raw6 pressure switch")
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundBoulder {
		t.Fatalf("pressure-door closing sounds=%v, want [%d]", events, SoundBoulder)
	}
	if rt.IsPassable(4, 17) {
		t.Fatal("raw6-linked door remains passable after switch is released")
	}
	state, _ = rt.At(BackgroundLayer, 4, 17)
	if state != 1 {
		t.Fatalf("closed pressure door state = %d, want retained count 1", state)
	}
}

func TestRuntimeClosingPressureDoorCrushesHeroInside(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 4, Y: 17}
	rt.PlayerMotion = ObjectMotion{}
	rt.SetForTest(PlayerLayer, 4, 17, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 4, 17, 7)
	rt.SetForTest(BackgroundLayer, 4, 17, 0x20)
	rt.closeDoorByID(0)
	if !rt.PlayerDead || rt.Health != 0 {
		t.Fatalf("closing pressure door dead=%v health=%d, want lethal source crush", rt.PlayerDead, rt.Health)
	}
	events := rt.DrainSoundEvents()
	if len(events) != 3 || events[0] != SoundBoulder || events[1] != SoundHeroHurt || events[2] != SoundDeath {
		t.Fatalf("pressure-door crush sounds=%v, want [%d %d %d]", events, SoundBoulder, SoundHeroHurt, SoundDeath)
	}
}

func TestRuntimeForegroundRaw2BarrierClearsAfterAdjacentBreakables(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 18, Y: 4}
	rt.SetForTest(PlayerLayer, 18, 4, EmptyRawID)
	if !rt.TryMove(1, 0) {
		t.Fatal("foreground raw2 must remain player-passable while linked breakables exist")
	}
	settleRuntimePlayerMotion(rt)
	if cleared := rt.TickForegroundTriggers(); cleared != 0 {
		t.Fatalf("raw2 cleared with adjacent raw30 present: %d", cleared)
	}
	for x := 19; x <= 22; x++ {
		rt.SetForTest(PlayerLayer, x, 5, EmptyRawID)
	}
	if cleared := rt.TickForegroundTriggers(); cleared != 4 {
		t.Fatalf("raw2 cleared = %d, want connected blob of 4", cleared)
	}
	for x := 19; x <= 22; x++ {
		id, _ := rt.At(ForegroundLayer, x, 4)
		if id != EmptyRawID {
			t.Fatalf("foreground raw2 at x=%d remains %d, want empty", x, id)
		}
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to move through cleared foreground raw2 barrier")
	}
}

func TestRuntimeForegroundRaw2State1RequiresHookLevelAction(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage04.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 19, Y: 20}
	rt.SetForTest(PlayerLayer, 19, 20, EmptyRawID)
	rt.SetForTest(PlayerLayer, 20, 20, EmptyRawID)
	if foreground, _ := rt.At(ForegroundLayer, 20, 20); foreground != 2 {
		t.Fatalf("test fixture foreground = %d, want raw2", foreground)
	}
	if state, _ := rt.At(BackgroundLayer, 20, 20); state != 1 {
		t.Fatalf("test fixture raw2 state = %d, want 1", state)
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("source movement must allow the hero onto state1 foreground raw2")
	}
	settleRuntimePlayerMotion(rt)
	if rt.UseSpecialBarrier(1, 0) {
		t.Fatal("UseSpecialBarrier without tool = true, want false")
	}
	rt.SpecialItemMask = 1
	if rt.UseSpecialBarrier(1, 0) {
		t.Fatal("UseSpecialBarrier with level 1 tool = true, want false")
	}
	rt.SpecialItemMask = 2
	if !rt.UseSpecialBarrier(1, 0) {
		t.Fatal("UseSpecialBarrier with level 2 tool = false, want true")
	}
	foreground, _ := rt.At(ForegroundLayer, 20, 20)
	if foreground != EmptyRawID {
		t.Fatalf("foreground raw2 after action = %d, want empty", foreground)
	}
	if rt.Player != (Point{X: 20, Y: 20}) {
		t.Fatalf("player=%+v after clearing state1 raw2, want standing on 20,20", rt.Player)
	}
}

func TestRuntimeForegroundRaw2UsesSourceToolPrompt(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage04.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 20, Y: 20}
	rt.SetForTest(PlayerLayer, 20, 20, EmptyRawID)
	if module, available, visible := rt.SpecialBarrierPrompt(); module != 1 || available || !visible {
		t.Fatalf("missing hook prompt=%d/%v/%v, want module1 unavailable visible", module, available, visible)
	}
	rt.SpecialItemMask = 2
	if module, available, visible := rt.SpecialBarrierPrompt(); module != 1 || !available || !visible {
		t.Fatalf("owned hook prompt=%d/%v/%v, want module1 available visible", module, available, visible)
	}
	rt.SetForTest(BackgroundLayer, 20, 20, 0)
	rt.SpecialItemMask = 1
	if module, available, visible := rt.SpecialBarrierPrompt(); module != 0 || !available || !visible {
		t.Fatalf("owned hammer prompt=%d/%v/%v, want module0 available visible", module, available, visible)
	}
	rt.SetForTest(ForegroundLayer, 20, 20, EmptyRawID)
	if _, _, visible := rt.SpecialBarrierPrompt(); visible {
		t.Fatal("tool prompt remained after raw2 cleared")
	}
}

func TestRuntimeForegroundRaw1ClearsConnectedClusterOnEnter(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage05.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 6, Y: 36}
	rt.SetForTest(PlayerLayer, 6, 36, EmptyRawID)
	for x := 7; x <= 16; x++ {
		if id, _ := rt.At(ForegroundLayer, x, 36); id != 1 {
			t.Fatalf("test fixture foreground at x=%d = %d, want raw1", x, id)
		}
		rt.SetForTest(PlayerLayer, x, 36, EmptyRawID)
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter foreground raw1 cluster")
	}
	if rt.Player != (Point{X: 7, Y: 36}) {
		t.Fatalf("player = %+v, want first raw1 cell", rt.Player)
	}
	for x := 7; x <= 16; x++ {
		id, _ := rt.At(ForegroundLayer, x, 36)
		if id != EmptyRawID {
			t.Fatalf("foreground raw1 at x=%d remains %d, want empty", x, id)
		}
	}
}

func TestRuntimeForegroundRaw0RecordsEventAtSourceMotionOffset(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 5, Y: 18}
	rt.SetForTest(PlayerLayer, 5, 18, EmptyRawID)
	rt.SetForTest(PlayerLayer, 6, 18, EmptyRawID)
	if id, _ := rt.At(ForegroundLayer, 6, 18); id != 0 {
		t.Fatalf("test fixture foreground = %d, want raw0", id)
	}
	if state, _ := rt.At(BackgroundLayer, 6, 18); state != 30 {
		t.Fatalf("test fixture background = %d, want 30", state)
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter foreground raw0 event")
	}
	if rt.ForegroundEvents != 0 {
		t.Fatalf("foreground event triggered at movement offset %d, want source delay", rt.PlayerMotion.Remaining)
	}
	rt.TickSourceFrame(8, 1, 0)
	rt.AdvancePlayerMotion()
	rt.TickSourceFrame(8, 2, 0)
	rt.AdvancePlayerMotion()
	if rt.PlayerMotion.Remaining != 6 {
		t.Fatalf("movement offset=%d, want 6 before foreground event scan", rt.PlayerMotion.Remaining)
	}
	rt.TickSourceFrame(8, 3, 0)
	if rt.LastForegroundEvent != 30 || rt.ForegroundEvents != 1 {
		t.Fatalf("foreground event last=%d count=%d, want 30/1", rt.LastForegroundEvent, rt.ForegroundEvents)
	}
	id, _ := rt.At(ForegroundLayer, 6, 18)
	if id != EmptyRawID {
		t.Fatalf("foreground raw0 cell = %d, want empty after trigger", id)
	}
}

func TestRuntimeCollectsKeysLivesAndHealthRefill(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.Health = 1
	pickups := []struct {
		id RawID
		x  int
	}{
		{id: 4, x: 1},
		{id: 5, x: 2},
		{id: 6, x: 3},
		{id: 7, x: 4},
	}
	for _, pickup := range pickups {
		rt.SetForTest(PlayerLayer, pickup.x, 17, pickup.id)
		rt.SetForTest(ForegroundLayer, pickup.x, 17, EmptyRawID)
		if !rt.TryMove(1, 0) {
			t.Fatalf("move into raw %d failed", pickup.id)
		}
		settleRuntimePlayerMotion(rt)
		id, _ := rt.At(PlayerLayer, pickup.x, 17)
		if id != EmptyRawID {
			t.Fatalf("collected raw %d cell = %d, want empty", pickup.id, id)
		}
	}
	if rt.KeyForForeground9 != 1 {
		t.Fatalf("key for foreground 9 = %d, want 1", rt.KeyForForeground9)
	}
	if rt.KeyForForeground8 != 1 {
		t.Fatalf("key for foreground 8 = %d, want 1", rt.KeyForForeground8)
	}
	if rt.ExtraLives != 6 {
		t.Fatalf("extra lives = %d, want 6", rt.ExtraLives)
	}
	if rt.HealthRefills != 1 || rt.Health != rt.MaxHealth {
		t.Fatalf("health refill count=%d health=%d max=%d, want 1/full", rt.HealthRefills, rt.Health, rt.MaxHealth)
	}
}

func TestRuntimeCollectsRaw41BonusValueFromBackground(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage03.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 13, Y: 9}
	rt.SetForTest(PlayerLayer, 14, 9, 41)
	rt.SetForTest(BackgroundLayer, 14, 9, 15)
	rt.SetForTest(ForegroundLayer, 14, 9, 33)
	if !rt.TryMove(1, 0) {
		t.Fatal("move into raw41 bonus failed")
	}
	if rt.BonusValue != 0 || !rt.pendingChestSet {
		t.Fatalf("entering raw41 chest bonus=%d pending=%v, want 0/true", rt.BonusValue, rt.pendingChestSet)
	}
	settleRuntimePlayerMotion(rt)
	rt.TickSourceFrame(8, 1, 0)
	for sourceTick := 2; sourceTick <= chestRewardTick+1; sourceTick++ {
		rt.gravitySourceTick = sourceTick
		rt.TickStatus()
	}
	if rt.BonusValue != 15 || rt.BonusPickups != 1 {
		t.Fatalf("bonus value=%d pickups=%d, want 15/1", rt.BonusValue, rt.BonusPickups)
	}
	if rt.VioletGems != 15 {
		t.Fatalf("raw41 violet gems=%d, want source value 15", rt.VioletGems)
	}
	if rt.BonusRemaining != 10 || rt.BonusGateOpen {
		t.Fatalf("bonus remaining=%d open=%v, want 10/false", rt.BonusRemaining, rt.BonusGateOpen)
	}
	id, _ := rt.At(PlayerLayer, 14, 9)
	if id != EmptyRawID {
		t.Fatalf("bonus cell = %d, want empty", id)
	}
}

func TestRuntimeBonusQuotaOpensAfterEnoughBonusValue(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.SetForTest(PlayerLayer, 1, 17, 41)
	rt.SetForTest(BackgroundLayer, 1, 17, 10)
	rt.SetForTest(ForegroundLayer, 1, 17, EmptyRawID)
	if !rt.TryMove(1, 0) {
		t.Fatal("move into raw41 quota bonus failed")
	}
	if rt.BonusRemaining != 0 || !rt.BonusGateOpen {
		t.Fatalf("bonus remaining=%d open=%v, want 0/true", rt.BonusRemaining, rt.BonusGateOpen)
	}
	if rt.VioletGems != 10 {
		t.Fatalf("raw41 quota violet gems=%d, want 10", rt.VioletGems)
	}
}

func TestRuntimeCollectsSourceSpecialPickups(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	pickups := []struct {
		id RawID
		x  int
	}{
		{id: 24, x: 1},
		{id: 27, x: 2},
		{id: 26, x: 3},
		{id: 42, x: 4},
		{id: 53, x: 5},
	}
	for _, pickup := range pickups {
		rt.SetForTest(PlayerLayer, pickup.x, 17, pickup.id)
		rt.SetForTest(ForegroundLayer, pickup.x, 17, EmptyRawID)
		if !rt.TryMove(1, 0) {
			t.Fatalf("move into raw %d failed", pickup.id)
		}
		settleRuntimePlayerMotion(rt)
		id, _ := rt.At(PlayerLayer, pickup.x, 17)
		if id != EmptyRawID {
			t.Fatalf("collected raw %d cell = %d, want empty", pickup.id, id)
		}
	}
	if rt.SpecialItemMask != 11 {
		t.Fatalf("special item mask = %d, want raw24|27|26 mask 11", rt.SpecialItemMask)
	}
	if !rt.SpecialPickup42 {
		t.Fatal("special raw42 pickup = false, want true")
	}
	if rt.RelicMask != 1 {
		t.Fatalf("relic mask = %d, want 1 from raw53", rt.RelicMask)
	}
	if rt.SpecialPickups != 5 {
		t.Fatalf("special pickups = %d, want 5", rt.SpecialPickups)
	}
}

func TestRuntimeFullHealthRefillBecomesBonusValue(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 0, Y: 17}
	rt.Health = rt.MaxHealth
	rt.SetForTest(PlayerLayer, 1, 17, 7)
	rt.SetForTest(ForegroundLayer, 1, 17, EmptyRawID)
	if !rt.TryMove(1, 0) {
		t.Fatal("move into full-health raw7 failed")
	}
	if rt.HealthRefills != 0 {
		t.Fatalf("health refills = %d, want 0 when converted to bonus", rt.HealthRefills)
	}
	if rt.BonusValue != 10 || rt.BonusPickups != 1 {
		t.Fatalf("bonus value=%d pickups=%d, want 10/1", rt.BonusValue, rt.BonusPickups)
	}
	if rt.VioletGems != 10 || rt.TotalVioletGems != 31 {
		t.Fatalf("converted refill violet=%d/%d, want 10/31", rt.VioletGems, rt.TotalVioletGems)
	}
}

func TestRuntimePersistsOnlySourcePermanentChestRewards(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage00.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	chest := Point{X: 19, Y: 2}
	rt.Player = chest
	rt.ChestRewardID = 2
	rt.applyChestReward()
	if got := rt.PersistentRewardCoordinates(); len(got) != 1 || got[0] != chest {
		t.Fatalf("red-diamond persistent coordinates=%v, want [%+v]", got, chest)
	}

	replay, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	replay.ApplyPersistentRewardCoordinates(rt.PersistentRewardCoordinates())
	if id, _ := replay.At(PlayerLayer, chest.X, chest.Y); id != EmptyRawID {
		t.Fatalf("consumed red-diamond chest payload=%d, want empty", id)
	}
	if replay.ObjectState[replay.index(chest.X, chest.Y)] != 3 {
		t.Fatalf("consumed hidden chest state=%d, want final open frame 3", replay.ObjectState[replay.index(chest.X, chest.Y)])
	}

	rt.ConsumedRewardCells[rt.index(chest.X, chest.Y)] = false
	rt.ChestRewardID = 4
	rt.applyChestReward()
	if got := rt.PersistentRewardCoordinates(); len(got) != 0 {
		t.Fatalf("key chest persisted coordinates=%v, want replayable source key", got)
	}
}

func TestRuntimeConsumedBossRelicBecomesTenPointBonus(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage08.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	seal := Point{X: 27, Y: 6}
	rt.ApplyPersistentRewardCoordinates([]Point{seal})
	if id, _ := rt.At(PlayerLayer, seal.X, seal.Y); id != 41 {
		t.Fatalf("replayed seal chest payload=%d, want raw41", id)
	}
	if value, _ := rt.At(BackgroundLayer, seal.X, seal.Y); value != 10 {
		t.Fatalf("replayed seal chest value=%d, want 10", value)
	}
	if rt.TotalVioletGems != 10 {
		t.Fatalf("replayed seal stage total violet=%d, want 10", rt.TotalVioletGems)
	}
}

func TestRuntimeOpensForeground9LockWithRaw4Key(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 6, Y: 17}
	rt.SetForTest(PlayerLayer, 7, 17, 31)
	rt.SetForTest(ForegroundLayer, 7, 17, 9)
	if rt.startAdjacentLockOpening() {
		t.Fatal("started foreground raw 9 lock without raw 4 key")
	}
	rt.KeyForForeground9 = 1
	rt.SetPlayerTurnOffset(12)
	if rt.startAdjacentLockOpening() {
		t.Fatal("started foreground raw 9 lock before turn jInt reached 6")
	}
	rt.SetPlayerTurnOffset(6)
	if !rt.startAdjacentLockOpening() {
		t.Fatal("failed to start foreground raw 9 lock with raw 4 key")
	}
	if rt.LockAnimation != 18 || !rt.LockOpening {
		t.Fatalf("lock animation=%d opening=%v, want source right-facing animation18", rt.LockAnimation, rt.LockOpening)
	}
	for tick := 1; tick < lockRewardTick; tick++ {
		rt.tickLockOpening()
	}
	if rt.KeyForForeground9 != 1 || rt.LocksOpened != 0 {
		t.Fatalf("before reward key9=%d locks=%d, want 1/0", rt.KeyForForeground9, rt.LocksOpened)
	}
	rt.tickLockOpening()
	if rt.KeyForForeground9 != 0 || rt.LocksOpened != 1 || !rt.LockOpening {
		t.Fatalf("reward tick key9=%d locks=%d opening=%v, want 0/1/true", rt.KeyForForeground9, rt.LocksOpened, rt.LockOpening)
	}
	if rt.Player != (Point{X: 6, Y: 17}) {
		t.Fatalf("player = %+v, want to remain beside lock", rt.Player)
	}
	playerID, _ := rt.At(PlayerLayer, 7, 17)
	foregroundID, _ := rt.At(ForegroundLayer, 7, 17)
	if playerID != 31 || foregroundID != 9 || rt.ObjectState[rt.index(7, 17)] != 1 {
		t.Fatalf("opened lock player=%d foreground=%d state=%d, want raw31/raw9/state1", playerID, foregroundID, rt.ObjectState[rt.index(7, 17)])
	}
	doorState, _ := rt.At(BackgroundLayer, 7, 18)
	if doorState != 0x11 || rt.IsPassable(7, 18) {
		t.Fatalf("linked door state=%#x passable=%v, want phase1/blocking", doorState, rt.IsPassable(7, 18))
	}
	rt.tickDoorAnimations(3)
	if !rt.IsPassable(7, 18) {
		t.Fatal("linked door did not become passable at source phase2")
	}
	for rt.LockOpening {
		rt.tickLockOpening()
	}
	if rt.LockTicks != 0 || rt.CanAcceptInput() == false {
		t.Fatalf("completed lock ticks=%d input=%v, want 0/true", rt.LockTicks, rt.CanAcceptInput())
	}
}

func TestRuntimeSharedDoorWaitsForEveryKeyedLock(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	door := Point{X: 27, Y: 15}
	if state, _ := rt.At(BackgroundLayer, door.X, door.Y); state != 2 {
		t.Fatalf("two-lock door initial count = %#x, want 2", state)
	}
	rt.KeyForForeground8 = 2
	for lockIndex, y := range []int{13, 14} {
		rt.Player = Point{X: 26, Y: y}
		rt.PlayerMotion = ObjectMotion{}
		if !rt.startAdjacentLockOpening() {
			t.Fatalf("failed to start shared lock %d", lockIndex+1)
		}
		for !rt.LockRewarded {
			rt.tickLockOpening()
		}
		state, _ := rt.At(BackgroundLayer, door.X, door.Y)
		if lockIndex == 0 {
			if state != 1 || rt.IsPassable(door.X, door.Y) {
				t.Fatalf("door after first lock state=%#x passable=%v, want count1/closed", state, rt.IsPassable(door.X, door.Y))
			}
		} else if state != 0x11 || rt.IsPassable(door.X, door.Y) {
			t.Fatalf("door after second lock state=%#x passable=%v, want phase1/blocking", state, rt.IsPassable(door.X, door.Y))
		}
		for rt.LockOpening {
			rt.tickLockOpening()
		}
	}
	rt.tickDoorAnimations(3)
	if !rt.IsPassable(door.X, door.Y) {
		t.Fatal("shared door did not become passable after the second lock")
	}
}

func TestRuntimeOpensForeground8LockWithRaw5Key(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 6, Y: 17}
	rt.KeyForForeground8 = 1
	rt.SetForTest(PlayerLayer, 7, 17, 31)
	rt.SetForTest(ForegroundLayer, 7, 17, 8)
	rt.SetForTest(BackgroundLayer, 7, 17, 0)
	if !rt.startAdjacentLockOpening() {
		t.Fatal("failed to start foreground raw 8 lock with raw 5 key")
	}
	for tick := 0; tick < lockRewardTick; tick++ {
		rt.tickLockOpening()
	}
	if rt.KeyForForeground8 != 0 || rt.LocksOpened != 1 {
		t.Fatalf("key8=%d locks=%d, want 0/1", rt.KeyForForeground8, rt.LocksOpened)
	}
	if rt.ObjectState[rt.index(7, 17)] != 1 {
		t.Fatalf("silver lock state=%d, want unlocked frame state1", rt.ObjectState[rt.index(7, 17)])
	}
}

func mustLoadOriginalStage(t *testing.T, name string) *Stage {
	t.Helper()
	stage, err := LoadStageFile(filepath.Join("..", "..", "decoded", "world0", name))
	if err != nil {
		t.Fatal(err)
	}
	return stage
}

func clearRuntimeSnakes(rt *Runtime) {
	clearRuntimePlayerIDs(rt, 19, 43)
}

func runtimePointsWithRaw(rt *Runtime, id RawID) []Point {
	points := make([]Point, 0)
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			if rt.PlayerLayer[rt.index(x, y)] == id {
				points = append(points, Point{X: x, Y: y})
			}
		}
	}
	return points
}

func clearRuntimePlayerIDs(rt *Runtime, ids ...RawID) {
	remove := map[RawID]bool{}
	for _, id := range ids {
		remove[id] = true
	}
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if remove[rt.PlayerLayer[idx]] {
				rt.PlayerLayer[idx] = EmptyRawID
				rt.ObjectState[idx] = 0
			}
		}
	}
}

func tickGravityObject(rt *Runtime, x, y int) int {
	rt.gravitySourceTick++
	if rt.tickGravityObjectAt(x, y) {
		return 1
	}
	return 0
}

func tickGravityObjectUntilMoved(rt *Runtime, x, y, maxTicks int) (int, int) {
	for sourceTick := 1; sourceTick <= maxTicks; sourceTick++ {
		if moved := tickGravityObject(rt, x, y); moved != 0 {
			return moved, sourceTick
		}
	}
	return 0, maxTicks
}

func settleRuntimePlayerMotion(rt *Runtime) {
	for rt.PlayerMotion.Remaining > 0 {
		rt.AdvancePlayerMotion()
	}
}

func TestPlayerTurnOffsetPreservesSourceJIntWithoutMovement(t *testing.T) {
	rt := &Runtime{}
	rt.SetPlayerTurnOffset(18)
	if got := rt.playerSourceOffset(); got != 18 {
		t.Fatalf("turn source offset=%d, want 18", got)
	}
	rt.PlayerMotion = ObjectMotion{DX: 1, Remaining: 12}
	if got := rt.playerSourceOffset(); got != 12 {
		t.Fatalf("movement source offset=%d, want movement to take precedence at 12", got)
	}
	rt.PlayerMotion = ObjectMotion{}
	rt.SetPlayerTurnOffset(-1)
	if got := rt.playerSourceOffset(); got != 0 {
		t.Fatalf("clamped source offset=%d, want 0", got)
	}
}

func tickRuntimeHookToCompletion(t *testing.T, rt *Runtime, maxTicks int) int {
	t.Helper()
	return tickRuntimeHookToCompletionFrom(t, rt, 1, maxTicks)
}

func tickRuntimeHookToCompletionFrom(t *testing.T, rt *Runtime, startTick, maxTick int) int {
	t.Helper()
	for sourceTick := startTick; sourceTick <= maxTick; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
		if !rt.Hooking {
			return sourceTick
		}
	}
	t.Fatalf("hook remained active through source tick %d: target=%+v animation=%d actionTick=%d", maxTick, rt.HookTarget, rt.HookAnimation, rt.HookTicks)
	return maxTick
}
