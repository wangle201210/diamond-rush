package original

import "testing"

func newStage11Route(t *testing.T) *stage07Route {
	t.Helper()
	stage := mustLoadOriginalStage(t, "stage11.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 8
	rt.MaxHealth = 8
	rt.Health = 8
	return &stage07Route{t: t, rt: rt}
}

func TestRuntimeStage11EnemyArenaKeyChestsStartLocked(t *testing.T) {
	route := newStage11Route(t)
	rt := route.rt
	for _, key := range []Point{{X: 8, Y: 4}, {X: 20, Y: 2}, {X: 32, Y: 6}, {X: 18, Y: 13}} {
		if !rt.ContainerLockedAt(key.X, key.Y) {
			t.Errorf("arena key chest %+v starts unlocked", key)
		}
	}
	rt.Player = Point{X: 7, Y: 4}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to step onto locked source container")
	}
	settleRuntimePlayerMotion(rt)
	rt.TickSourceFrame(8, 1, 0)
	if rt.pendingChestSet || rt.ChestOpening || rt.KeyForForeground9 != 0 {
		t.Fatalf("locked chest pending=%v opening=%v keys=%d, want inert", rt.pendingChestSet, rt.ChestOpening, rt.KeyForForeground9)
	}
	rt.openEnemyGateGroup(0)
	if rt.ContainerLockedAt(8, 4) {
		t.Fatal("group-0 completion did not unlock its key chest")
	}
	rt.TickSourceFrame(8, 2, 0)
	if !rt.ChestOpening {
		t.Fatal("unlocked chest did not start while hero remained on its cell")
	}
}

func solveStage11EnemyGroup0(route *stage07Route) {
	rt := route.rt
	route.walkTo("left snake freeze position", Point{X: 5, Y: 10})
	route.hammer("freeze left group-0 snake", 1, 0)
	route.push("push left frozen snake under drop", 1)
	route.waitUntil("left frozen snake settles", 80, func() bool {
		id, _ := rt.At(PlayerLayer, 7, 10)
		return id == 9 && rt.ObjectMotion[rt.index(7, 10)].Remaining == 0
	})
	route.walkTo("left boulder push position", Point{X: 5, Y: 7})
	route.push("drop left boulder onto frozen snake", 1)
	route.waitUntil("left boulder rests on frozen snake", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 7, 9)
		return id == 0 && rt.ObjectMotion[rt.index(7, 9)].Remaining == 0
	})
	route.walkTo("left ice thaw position", Point{X: 6, Y: 10})
	route.hammer("thaw left snake under boulder", 1, 0)
	route.waitUntil("left group-0 snake crushed", 80, func() bool {
		return rt.EnemyGateCounters[0] == 1
	})

	route.walkTo("right snake freeze position", Point{X: 11, Y: 10})
	route.hammer("freeze right group-0 snake", -1, 0)
	route.push("push right frozen snake under drop", -1)
	route.waitUntil("right frozen snake settles", 80, func() bool {
		id, _ := rt.At(PlayerLayer, 9, 10)
		return id == 9 && rt.ObjectMotion[rt.index(9, 10)].Remaining == 0
	})
	route.walkTo("right boulder push position", Point{X: 11, Y: 7})
	route.push("drop right boulder onto frozen snake", -1)
	route.waitUntil("right boulder rests on frozen snake", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 9, 9)
		return id == 0 && rt.ObjectMotion[rt.index(9, 9)].Remaining == 0
	})
	route.walkTo("right ice thaw position", Point{X: 10, Y: 10})
	route.hammer("thaw right snake under boulder", -1, 0)
	route.waitUntil("group 0 cleared", 80, func() bool {
		return rt.EnemyGateCounters[0] == 0 && !rt.ContainerLockedAt(8, 4)
	})
}

func solveStage11EnemyGroup1(route *stage07Route) {
	rt := route.rt
	route.walkTo("left group-1 snake freeze position", Point{X: 18, Y: 7})
	route.waitUntil("group-1 snake crosses left drop column", 500, func() bool {
		id, _ := rt.At(PlayerLayer, 19, 7)
		return id == 19 && rt.ObjectMotion[rt.index(19, 7)].Remaining >= 9
	})
	route.hammer("freeze left group-1 snake", 1, 0)
	route.waitUntil("left group-1 ice settles below arena", 80, func() bool {
		id, _ := rt.At(PlayerLayer, 19, 8)
		return id == 9 && rt.ObjectMotion[rt.index(19, 8)].Remaining == 0
	})

	rightAlreadyFrozen, _ := rt.At(PlayerLayer, 21, 8)
	if rightAlreadyFrozen != 9 {
		route.walkTo("right endpoint snake freeze position", Point{X: 22, Y: 6})
		route.waitUntil("remaining snake reaches supported right endpoint", 500, func() bool {
			id, _ := rt.At(PlayerLayer, 23, 7)
			return id == 19 && rt.ObjectMotion[rt.index(23, 7)].Remaining >= 9
		})
		route.hammer("freeze remaining snake at right endpoint", 0, 1)
		if id, _ := rt.At(PlayerLayer, 23, 7); id != 9 {
			route.t.Fatalf("right endpoint freeze raw=%d frozen=%d state=%#x motion=%+v ice=%v snakes=%v, want ice", id, rt.FrozenOriginalAt(23, 7), rt.ObjectState[rt.index(23, 7)], rt.ObjectMotion[rt.index(23, 7)], runtimePointsWithRaw(rt, 9), runtimePointsWithRaw(rt, 19))
		}
	}

	route.walkTo("group-1 central hook position", Point{X: 20, Y: 4})
	route.hook("release left group-1 boulder", -1)
	for _, point := range []Point{
		{X: 21, Y: 4}, {X: 21, Y: 5}, {X: 20, Y: 5}, {X: 19, Y: 5},
		{X: 19, Y: 6}, {X: 18, Y: 6}, {X: 18, Y: 7}, {X: 18, Y: 8},
	} {
		route.walkTo("race left group-1 boulder", point)
	}
	if id, _ := rt.At(PlayerLayer, 19, 7); id != 0 {
		route.t.Fatalf("left group-1 boulder missed thaw window: raw=%d boulders=%v", id, runtimePointsWithRaw(rt, 0))
	}
	route.waitUntil("left group-1 boulder enters crush range", 16, func() bool {
		id, _ := rt.At(PlayerLayer, 19, 7)
		return id == 0 && rt.ObjectMotion[rt.index(19, 7)].Remaining <= 6
	})
	route.hammer("thaw left group-1 snake under falling boulder", 1, 0)
	route.waitUntil("first group-1 snake crushed", 80, func() bool {
		return rt.EnemyGateCounters[1] == 1
	})

	if rightAlreadyFrozen != 9 {
		route.walkTo("right endpoint thaw position", Point{X: 23, Y: 6})
		route.hammer("thaw right endpoint snake toward corridor", 0, 1)
		route.walkTo("right group-1 snake freeze position", Point{X: 22, Y: 6})
		route.waitUntil("remaining snake crosses right drop column", 500, func() bool {
			id, _ := rt.At(PlayerLayer, 21, 7)
			return id == 19 && rt.ObjectMotion[rt.index(21, 7)].Remaining >= 9
		})
		route.hammer("freeze right group-1 snake", 0, 1)
		route.waitUntil("right group-1 ice settles below arena", 80, func() bool {
			id, _ := rt.At(PlayerLayer, 21, 8)
			return id == 9 && rt.ObjectMotion[rt.index(21, 8)].Remaining == 0
		})
	}
	route.walkTo("group-1 central hook position", Point{X: 20, Y: 4})
	route.hook("release right group-1 boulder", 1)
	for _, point := range []Point{{X: 20, Y: 5}, {X: 20, Y: 6}, {X: 20, Y: 7}, {X: 20, Y: 8}} {
		route.walkTo("race right group-1 boulder", point)
	}
	if id, _ := rt.At(PlayerLayer, 21, 7); id != 0 {
		route.t.Fatalf("right group-1 boulder missed thaw window: raw=%d boulders=%v", id, runtimePointsWithRaw(rt, 0))
	}
	route.waitUntil("right group-1 boulder enters crush range", 16, func() bool {
		id, _ := rt.At(PlayerLayer, 21, 7)
		return id == 0 && rt.ObjectMotion[rt.index(21, 7)].Remaining <= 6
	})
	route.hammer("thaw right group-1 snake under falling boulder", 1, 0)
	route.waitUntil("group 1 cleared", 80, func() bool {
		return rt.EnemyGateCounters[1] == 0 && !rt.ContainerLockedAt(20, 2)
	})
}

func solveStage11EnemyGroup2(route *stage07Route) {
	rt := route.rt
	route.walkTo("left group-2 snake freeze position", Point{X: 29, Y: 8})
	route.hammer("freeze group-2 left snake", -1, 0)
	route.walkTo("left group-2 snake hook position", Point{X: 30, Y: 8})
	route.hook("pull left frozen snake inward", -1)
	route.waitUntil("left group-2 ice settles", 80, func() bool {
		id, _ := rt.At(PlayerLayer, 29, 8)
		return id == 9 && rt.ObjectMotion[rt.index(29, 8)].Remaining == 0
	})
	route.walkTo("left group-2 boulder hook position", Point{X: 30, Y: 4})
	route.hook("pull left group-2 boulder inward", -1)
	route.walkTo("follow left group-2 boulder and block its roll", Point{X: 30, Y: 8})
	if id, _ := rt.At(PlayerLayer, 29, 7); id != 0 {
		route.t.Fatalf("left group-2 boulder escaped before roll block: raw=%d boulders=%v", id, runtimePointsWithRaw(rt, 0))
	}
	route.hammer("thaw left group-2 snake", -1, 0)
	route.waitUntil("left group-2 snake crushed", 80, func() bool {
		return rt.EnemyGateCounters[2] == 1
	})

	route.walkTo("right group-2 snake freeze position", Point{X: 31, Y: 8})
	route.hammer("freeze group-2 right snake", 1, 0)
	route.walkTo("right group-2 snake hook position", Point{X: 30, Y: 8})
	route.hook("pull right frozen snake inward", 1)
	route.waitUntil("right group-2 ice settles", 80, func() bool {
		id, _ := rt.At(PlayerLayer, 31, 8)
		return id == 9 && rt.ObjectMotion[rt.index(31, 8)].Remaining == 0
	})
	route.walkTo("right group-2 boulder hook position", Point{X: 30, Y: 4})
	route.hook("pull right group-2 boulder inward", 1)
	route.walkTo("follow right group-2 boulder and block its roll", Point{X: 30, Y: 8})
	if id, _ := rt.At(PlayerLayer, 31, 7); id != 0 {
		route.t.Fatalf("right group-2 boulder escaped before roll block: raw=%d boulders=%v", id, runtimePointsWithRaw(rt, 0))
	}
	route.hammer("thaw right group-2 snake", 1, 0)
	route.waitUntil("group 2 cleared", 80, func() bool {
		return rt.EnemyGateCounters[2] == 0 && !rt.ContainerLockedAt(32, 6)
	})
}

func solveStage11EnemyGroup3(route *stage07Route) {
	rt := route.rt
	route.walkTo("left group-3 snake bottom freeze position", Point{X: 18, Y: 18})
	route.waitUntil("left group-3 snake reaches bottom", 300, func() bool {
		id, _ := rt.At(PlayerLayer, 19, 18)
		return id == 19 && rt.ObjectMotion[rt.index(19, 18)].Remaining == 0
	})
	route.hammer("freeze left group-3 snake", 1, 0)
	route.push("shift left frozen snake into hook range", 1)
	route.walkTo("group-3 bottom hook position", Point{X: 23, Y: 18})
	route.hook("pull first frozen snake into stack column", -1)
	route.waitUntil("first frozen snake settles at stack bottom", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 22, 18)
		return id == 9 && rt.ObjectMotion[rt.index(22, 18)].Remaining == 0
	})

	route.walkTo("group-3 top hook position", Point{X: 23, Y: 13})
	route.waitUntil("near group-3 snake reaches top hook row", 400, func() bool {
		id, _ := rt.At(PlayerLayer, 21, 13)
		return id == 19 && rt.ObjectState[rt.index(21, 13)]&objectDirectionMask == 1
	})
	route.hook("pull first live snake into stack column", -1)
	route.waitUntil("first live snake descends onto bottom ice", 240, func() bool {
		id, _ := rt.At(PlayerLayer, 22, 17)
		return id == 19 && rt.ObjectMotion[rt.index(22, 17)].Remaining == 0
	})
	route.waitUntil("second group-3 snake reaches top hook row", 400, func() bool {
		left, _ := rt.At(PlayerLayer, 20, 13)
		near, _ := rt.At(PlayerLayer, 21, 13)
		return left == 19 || near == 19
	})
	route.hook("pull second live snake into stack column", -1)
	route.waitUntil("second live snake descends onto first", 240, func() bool {
		id, _ := rt.At(PlayerLayer, 22, 16)
		return id == 19 && rt.ObjectMotion[rt.index(22, 16)].Remaining == 0
	})
	route.walkTo("middle stacked snake freeze position", Point{X: 21, Y: 17})
	route.hammer("freeze middle stacked live snake", 1, 0)
	route.walkTo("top stacked snake freeze position", Point{X: 21, Y: 16})
	route.hammer("freeze top stacked live snake", 1, 0)
	route.waitUntil("three-snake ice stack completes", 80, func() bool {
		for y := 16; y <= 18; y++ {
			id, _ := rt.At(PlayerLayer, 22, y)
			if id != 9 {
				return false
			}
		}
		return true
	})

	route.walkTo("middle stacked snake thaw position", Point{X: 21, Y: 17})
	route.hammer("thaw middle stacked snake", 1, 0)
	route.waitUntil("top ice crushes middle stacked snake", 80, func() bool {
		id, _ := rt.At(PlayerLayer, 22, 17)
		return rt.EnemyGateCounters[3] == 2 && id == 9 && rt.ObjectMotion[rt.index(22, 17)].Remaining == 0
	})
	route.walkTo("bottom stack thaw position", Point{X: 21, Y: 18})
	route.hammer("thaw bottom stacked snake", 1, 0)
	route.waitUntil("surviving ice crushes bottom stacked snake", 80, func() bool {
		id, _ := rt.At(PlayerLayer, 22, 18)
		return rt.EnemyGateCounters[3] == 1 && id == 9 && rt.ObjectMotion[rt.index(22, 18)].Remaining == 0
	})

	route.walkTo("group-3 boulder hook position", Point{X: 21, Y: 14})
	route.hook("pull group-3 boulder above surviving ice", 1)
	route.walkTo("follow final boulder and block roll", Point{X: 21, Y: 18})
	if id, _ := rt.At(PlayerLayer, 22, 17); id != 0 {
		route.t.Fatalf("group-3 boulder escaped final ice before thaw: raw=%d boulders=%v ice=%v", id, runtimePointsWithRaw(rt, 0), runtimePointsWithRaw(rt, 9))
	}
	route.hammer("thaw final snake under boulder", 1, 0)
	route.waitUntil("group 3 cleared", 80, func() bool {
		return rt.EnemyGateCounters[3] == 0 && !rt.ContainerLockedAt(18, 13)
	})
}

func TestRuntimeStage11EnemyArenaSolutions(t *testing.T) {
	tests := []struct {
		name    string
		group   int
		trigger Point
		solve   func(*stage07Route)
	}{
		{name: "group0", group: 0, trigger: Point{X: 12, Y: 6}, solve: solveStage11EnemyGroup0},
		{name: "group1", group: 1, trigger: Point{X: 17, Y: 6}, solve: solveStage11EnemyGroup1},
		{name: "group2", group: 2, trigger: Point{X: 28, Y: 6}, solve: solveStage11EnemyGroup2},
		{name: "group3", group: 3, trigger: Point{X: 24, Y: 18}, solve: solveStage11EnemyGroup3},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			route := newStage11Route(t)
			route.rt.Player = test.trigger
			route.rt.PlayerMotion = ObjectMotion{}
			route.rt.activateEnemyGateTriggerAt(test.trigger.X, test.trigger.Y)
			route.waitUntil("enemy-gate camera demo ends", 180, func() bool {
				return !route.rt.EnemyGateDemoActive
			})
			test.solve(route)
			if got := route.rt.EnemyGateCounters[test.group]; got != 0 {
				t.Fatalf("enemy group %d counter=%d, want 0", test.group, got)
			}
		})
	}
}

func TestRuntimeStage11SecretStageCanBeCompletedAtSourceCadence(t *testing.T) {
	route := newStage11Route(t)
	rt := route.rt
	route.tick()
	route.walkTo("automatic entrance before door", Point{X: 2, Y: 27})
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 12 source entrance door")
	}
	route.walkTo("entrance checkpoint", Point{X: 3, Y: 27})
	keysCollected := 0
	completeArena := func(group int, trigger, key Point, solve func(*stage07Route)) {
		route.walkTo("activate enemy-gate arena", trigger)
		if rt.ActiveEnemyGateGroup != group {
			t.Fatalf("arena trigger %+v activated group %d, want %d", trigger, rt.ActiveEnemyGateGroup, group)
		}
		route.waitUntil("enemy-gate camera demo ends", 180, func() bool {
			return !rt.EnemyGateDemoActive
		})
		solve(route)
		route.walkTo("collect arena gold key", key)
		keysCollected++
		route.waitUntil("arena gold key collected", 160, func() bool {
			return rt.KeyForForeground9 == keysCollected && !rt.ChestOpening
		})
	}

	completeArena(1, Point{X: 17, Y: 6}, Point{X: 20, Y: 2}, solveStage11EnemyGroup1)
	completeArena(3, Point{X: 24, Y: 18}, Point{X: 18, Y: 13}, solveStage11EnemyGroup3)
	route.walkTo("lower checkpoint", Point{X: 17, Y: 22})
	for _, point := range []Point{
		{X: 18, Y: 22}, {X: 18, Y: 23}, {X: 18, Y: 24}, {X: 18, Y: 25}, {X: 18, Y: 26},
		{X: 19, Y: 26}, {X: 19, Y: 27}, {X: 20, Y: 27}, {X: 21, Y: 27}, {X: 22, Y: 27},
	} {
		route.walkTo("approach lower boulder-maze support", point)
	}
	route.hammer("clear lower boulder-maze support", 1, 0)
	route.push("push lower boulder-maze rock right once", 1)
	route.push("push lower boulder-maze rock right twice", 1)
	for _, point := range []Point{
		{X: 24, Y: 26}, {X: 23, Y: 26},
		{X: 23, Y: 25}, {X: 23, Y: 24}, {X: 23, Y: 23}, {X: 24, Y: 23},
	} {
		route.walkTo("cross lower boulder maze", point)
	}
	route.walkTo("clear return-side boulder support", Point{X: 22, Y: 23})
	route.push("push return-side boulder left", -1)
	completeArena(0, Point{X: 12, Y: 6}, Point{X: 8, Y: 4}, solveStage11EnemyGroup0)
	completeArena(2, Point{X: 28, Y: 6}, Point{X: 32, Y: 6}, solveStage11EnemyGroup2)
	for _, checkpoint := range []Point{{X: 14, Y: 11}, {X: 27, Y: 11}} {
		route.walkTo("upper checkpoint", checkpoint)
	}
	for index, y := range []int{19, 20, 21, 22} {
		route.walkTo("left of exit gold lock", Point{X: 36, Y: y})
		wantLocks := index + 1
		wantKeys := 3 - index
		route.waitUntil("exit gold lock opens", 180, func() bool {
			return rt.KeyForForeground9 == wantKeys && rt.LocksOpened == wantLocks && !rt.LockOpening
		})
	}
	if !rt.IsPassable(37, 23) {
		t.Fatal("linked door below four gold locks did not open")
	}
	route.walkTo("right-side quota chest", Point{X: 38, Y: 11})
	route.waitUntil("right-side quota bonus collected", 160, func() bool {
		return rt.BonusValue >= 50 && !rt.ChestOpening
	})
	if rt.BonusRemaining > 0 || !rt.BonusGateOpen || !rt.IsPassable(34, 25) {
		t.Fatalf("third secret-stage quota gems=%d bonus=%d remaining=%d open=%v", rt.VioletGems, rt.BonusValue, rt.BonusRemaining, rt.BonusGateOpen)
	}
	route.walkTo("third Angkor secret-stage goal", Point{X: 35, Y: 25})
	if !rt.ReachedGoal || !rt.GoalExitSecret {
		t.Fatalf("third secret-stage finish player=%+v reached=%v secret=%v tick=%d", rt.Player, rt.ReachedGoal, rt.GoalExitSecret, route.sourceTick)
	}
}
