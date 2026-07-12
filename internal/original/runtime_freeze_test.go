package original

import "testing"

func TestRuntimeFreezeHammerUsesSourceFiveCellTargets(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 19, 43)
	rt.SpecialItemMask = 8
	rt.Player = Point{X: 10, Y: 10}
	rt.SetForTest(PlayerLayer, 11, 10, 1)
	rt.SetForTest(PlayerLayer, 12, 10, 19)
	rt.SetForTest(PlayerLayer, 11, 9, 43)
	rt.SetForTest(PlayerLayer, 11, 11, 19)
	rt.SetForTest(PlayerLayer, 12, 9, 43)
	for point, direction := range map[Point]int{
		{X: 12, Y: 10}: 2,
		{X: 11, Y: 9}:  1,
		{X: 11, Y: 11}: 3,
	} {
		idx := rt.index(point.X, point.Y)
		rt.ObjectState[idx] = direction
		rt.ObjectMotion[idx] = ObjectMotion{Remaining: 12}
	}

	if !rt.UseHammer(1, 0) {
		t.Fatal("Freeze Hammer did not start against the centered violet gem")
	}
	for tick := 1; tick < hammerImpactTick; tick++ {
		rt.tickHammerAction()
		if id, _ := rt.At(PlayerLayer, 11, 10); id != 1 {
			t.Fatalf("center gem froze before source impact tick %d: raw=%d", hammerImpactTick, id)
		}
	}
	rt.tickHammerAction()

	wants := map[Point]RawID{
		{X: 11, Y: 10}: 1,
		{X: 12, Y: 10}: 19,
		{X: 11, Y: 9}:  43,
		{X: 11, Y: 11}: 19,
	}
	for point, sourceID := range wants {
		id, _ := rt.At(PlayerLayer, point.X, point.Y)
		if id != 9 || rt.FrozenOriginalAt(point.X, point.Y) != sourceID {
			t.Errorf("frozen cell %+v raw/source=%d/%d, want 9/%d", point, id, rt.FrozenOriginalAt(point.X, point.Y), sourceID)
		}
	}
	if id, _ := rt.At(PlayerLayer, 12, 9); id != 43 || rt.FrozenOriginalAt(12, 9) != EmptyRawID {
		t.Fatalf("diagonal snake raw/source=%d/%d, want untouched raw43/empty source", id, rt.FrozenOriginalAt(12, 9))
	}
	if events := rt.DrainSoundEvents(); len(events) != 1 || events[0] != SoundEnemyHit {
		t.Fatalf("freeze impact sounds=%v, want [%d]", events, SoundEnemyHit)
	}
}

func TestRuntimeFreezeHammerDoesNotFreezeStaticAdjacentSnake(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 19, 43)
	rt.SpecialItemMask = 8
	rt.Player = Point{X: 10, Y: 10}
	rt.SetForTest(PlayerLayer, 11, 10, 1)
	rt.SetForTest(PlayerLayer, 12, 10, 19)
	if !rt.UseHammer(1, 0) {
		t.Fatal("Freeze Hammer did not start against centered gem")
	}
	for tick := 0; tick < hammerImpactTick; tick++ {
		rt.tickHammerAction()
	}
	if id, _ := rt.At(PlayerLayer, 12, 10); id != 19 {
		t.Fatalf("static adjacent snake raw=%d, want unchanged raw19", id)
	}
}

func TestRuntimeFreezeHammerScansMovingSnakeAfterBoulderBlock(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 19, 43)
	rt.SpecialItemMask = 8
	rt.Player = Point{X: 10, Y: 10}
	rt.SetForTest(PlayerLayer, 11, 10, 0)
	rt.SetForTest(PlayerLayer, 12, 10, 19)
	idx := rt.index(12, 10)
	rt.ObjectState[idx] = 2
	rt.ObjectMotion[idx] = ObjectMotion{Remaining: 12}
	if !rt.UseHammer(1, 0) {
		t.Fatal("Freeze Hammer did not start against boulder")
	}
	for tick := 0; tick < hammerImpactTick; tick++ {
		rt.tickHammerAction()
	}
	if id, _ := rt.At(PlayerLayer, 11, 10); id != 0 {
		t.Fatalf("blocked boulder raw=%d, want unchanged raw0", id)
	}
	if id, _ := rt.At(PlayerLayer, 12, 10); id != 9 || rt.FrozenOriginalAt(12, 10) != 19 {
		t.Fatalf("moving snake after boulder block raw/source=%d/%d, want 9/19", id, rt.FrozenOriginalAt(12, 10))
	}
	if events := rt.DrainSoundEvents(); len(events) != 2 || events[0] != SoundHammerBlock || events[1] != SoundEnemyHit {
		t.Fatalf("boulder-plus-freeze sounds=%v, want [%d %d]", events, SoundHammerBlock, SoundEnemyHit)
	}
}

func TestRuntimeFreezeHammerRequiresLevelEightAndSkipsScanOnBreakableHit(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 19, 43)
	rt.Player = Point{X: 10, Y: 10}
	rt.SetForTest(PlayerLayer, 11, 10, 1)
	rt.SpecialItemMask = 2
	if rt.UseHammer(1, 0) {
		t.Fatal("Mystic Hook tool level started Freeze Hammer against a gem")
	}

	rt.SpecialItemMask = 8
	rt.SetForTest(PlayerLayer, 11, 10, 30)
	rt.SetForTest(PlayerLayer, 12, 10, 19)
	if !rt.UseHammer(1, 0) {
		t.Fatal("Freeze Hammer did not start against raw30")
	}
	for tick := 0; tick < hammerImpactTick; tick++ {
		rt.tickHammerAction()
	}
	if state := rt.ObjectState[rt.index(11, 10)]; state != 1 {
		t.Fatalf("breakable state=%d, want source impact state 1", state)
	}
	if id, _ := rt.At(PlayerLayer, 12, 10); id != 19 || rt.FrozenOriginalAt(12, 10) != EmptyRawID {
		t.Fatalf("snake beside breakable raw/source=%d/%d, want untouched raw19/empty source", id, rt.FrozenOriginalAt(12, 10))
	}
}

func TestRuntimeFreezeHammerThawsSnakeWithSourceDirection(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 19, 43)
	rt.SpecialItemMask = 8
	rt.Player = Point{X: 11, Y: 9}
	rt.SetForTest(PlayerLayer, 11, 10, 43)
	if !rt.freezeObjectAt(11, 10) {
		t.Fatal("failed to prepare source red-snake frozen block")
	}

	if !rt.UseHammer(0, 1) {
		t.Fatal("Freeze Hammer did not start against raw9")
	}
	for tick := 0; tick < hammerImpactTick; tick++ {
		rt.tickHammerAction()
	}
	id, _ := rt.At(PlayerLayer, 11, 10)
	if id != 43 || rt.FrozenOriginalAt(11, 10) != EmptyRawID {
		t.Fatalf("thawed snake raw/source=%d/%d, want raw43/empty source", id, rt.FrozenOriginalAt(11, 10))
	}
	if state := rt.ObjectState[rt.index(11, 10)]; state != 2|snakeStunDuration {
		t.Fatalf("thawed snake state=%#x, want direction 2 plus source stun %#x", state, 2|snakeStunDuration)
	}
}

func TestRuntimeFrozenObjectsPreserveSourceAcrossPhysicsHookAndCheckpoint(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 19, 43)
	rt.SetForTest(PlayerLayer, 5, 10, 1)
	if !rt.freezeObjectAt(5, 10) {
		t.Fatal("failed to freeze violet gem")
	}
	rt.SetForTest(PlayerLayer, 6, 10, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 6, 10, EmptyRawID)
	if !rt.TryPushBoulder(5, 10, 1) {
		t.Fatal("raw9 frozen gem was not pushable")
	}
	if id, _ := rt.At(PlayerLayer, 6, 10); id != 9 || rt.FrozenOriginalAt(6, 10) != 1 {
		t.Fatalf("pushed frozen gem raw/source=%d/%d, want 9/1", id, rt.FrozenOriginalAt(6, 10))
	}
	rt.ObjectMotion[rt.index(6, 10)] = ObjectMotion{}
	rt.SetForTest(PlayerLayer, 6, 11, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 6, 11, EmptyRawID)
	if moved := tickGravityObject(rt, 6, 10); moved != 1 {
		t.Fatalf("frozen gem gravity moved=%d, want 1", moved)
	}
	if id, _ := rt.At(PlayerLayer, 6, 11); id != 9 || rt.FrozenOriginalAt(6, 11) != 1 {
		t.Fatalf("fallen frozen gem raw/source=%d/%d, want 9/1", id, rt.FrozenOriginalAt(6, 11))
	}
	rt.ObjectMotion[rt.index(6, 11)] = ObjectMotion{}
	if !rt.pressureSwitchActive(6, 11) {
		t.Fatal("settled raw9 frozen object did not activate a pressure switch")
	}
	rt.SaveSnapshot()
	if !rt.thawFrozenAt(6, 11) {
		t.Fatal("failed to mutate frozen object after checkpoint")
	}
	if !rt.RestoreCheckpoint() {
		t.Fatal("failed to restore frozen-object checkpoint")
	}
	if id, _ := rt.At(PlayerLayer, 6, 11); id != 9 || rt.FrozenOriginalAt(6, 11) != 1 {
		t.Fatalf("restored frozen gem raw/source=%d/%d, want 9/1", id, rt.FrozenOriginalAt(6, 11))
	}

	rt, err = NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 0, 1, 19, 43)
	rt.SpecialItemMask = 8
	rt.Player = Point{X: 4, Y: 10}
	for x := 4; x <= 7; x++ {
		rt.SetForTest(PlayerLayer, x, 10, EmptyRawID)
		rt.SetForTest(ForegroundLayer, x, 10, EmptyRawID)
		rt.SetForTest(PlayerLayer, x, 11, 80)
	}
	rt.SetForTest(PlayerLayer, 7, 10, 43)
	if !rt.freezeObjectAt(7, 10) {
		t.Fatal("failed to prepare frozen snake for hook")
	}
	if !rt.UseHook(1, 0) {
		t.Fatal("hook did not select raw9 frozen snake")
	}
	tickRuntimeHookToCompletion(t, rt, 30)
	if id, _ := rt.At(PlayerLayer, 5, 10); id != 9 || rt.FrozenOriginalAt(5, 10) != 43 {
		t.Fatalf("hooked frozen snake raw/source=%d/%d, want 9/43", id, rt.FrozenOriginalAt(5, 10))
	}
	if rt.FrozenOriginalAt(7, 10) != EmptyRawID {
		t.Fatalf("hook source retained frozen metadata %d", rt.FrozenOriginalAt(7, 10))
	}
}

func TestRuntimeHookUsesSourceRestoreStateByObjectType(t *testing.T) {
	if got := sourceHookRestoreState(9, 0x7000|gravityRollPreparing|3); got != 3 {
		t.Fatalf("frozen-object hook restore state=%#x, want direction-only %#x", got, 3)
	}
	if got := sourceHookRestoreState(14, 3); got != 3 {
		t.Errorf("moving-hazard hook restore state=%#x, want original state 3", got)
	}
	for _, id := range []RawID{11, 19, 43, 48} {
		if got := sourceHookRestoreState(id, 3); got != -1 {
			t.Errorf("raw %d hook restore state=%#x, want source sentinel -1", id, got)
		}
	}
}

func TestRuntimeHookReleasedSnakeResumesTowardPackedSourceTarget(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage11.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	clearRuntimePlayerIDs(rt, 19, 43)
	rt.Player = Point{X: 5, Y: 5}
	rt.SetForTest(PlayerLayer, 10, 10, 19)
	rt.SetForTest(PlayerLayer, 10, 11, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 10, 11, EmptyRawID)
	idx := rt.index(10, 10)
	rt.ObjectState[idx] = -1
	for sourceTick := 1; sourceTick <= 124; sourceTick++ {
		rt.gravitySourceTick = sourceTick
		rt.tickSnakeObjectAt(10, 10)
	}
	if rt.ObjectState[idx]&snakeStunMask != 0 {
		t.Fatalf("hook-release stun=%#x after source countdown, want zero", rt.ObjectState[idx]&snakeStunMask)
	}
	rt.gravitySourceTick = 125
	if !rt.tickSnakeObjectAt(10, 10) {
		t.Fatal("hook-released snake did not resume toward packed (127,127) target")
	}
	if id, _ := rt.At(PlayerLayer, 10, 11); id != 19 {
		t.Fatalf("hook-released snake target raw=%d, want raw19 one cell down", id)
	}
}
