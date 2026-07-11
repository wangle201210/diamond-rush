package original

import "testing"

func newStage10Route(t *testing.T) *stage07Route {
	t.Helper()
	stage := mustLoadOriginalStage(t, "stage10.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 8
	rt.MaxHealth = 8
	rt.Health = 8
	return &stage07Route{t: t, rt: rt}
}

func TestRuntimeStage10SecretStageCanBeCompletedAtSourceCadence(t *testing.T) {
	route := newStage10Route(t)
	rt := route.rt
	route.tick()
	route.walkTo("automatic entrance before door", Point{X: 1, Y: 19})
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 11 source entrance door")
	}
	route.walkTo("entrance checkpoint", Point{X: 2, Y: 19})
	route.walkTo("first entrance bonus chest", Point{X: 4, Y: 26})
	route.waitUntil("first entrance bonus collected", 160, func() bool {
		return rt.BonusValue >= 20 && !rt.ChestOpening
	})
	route.walkTo("second entrance bonus chest", Point{X: 6, Y: 26})
	route.waitUntil("second entrance bonus collected", 160, func() bool {
		return rt.BonusValue >= 40 && !rt.ChestOpening
	})
	route.walkTo("right of first-switch upper support", Point{X: 16, Y: 26})
	route.waitUntil("first-switch boulder reaches lower support", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 15, 26)
		return id == 0 && rt.ObjectMotion[rt.index(15, 26)].Remaining == 0
	})
	route.walkTo("clear first-switch left gem", Point{X: 16, Y: 28})
	route.walkTo("clear first-switch right gem", Point{X: 17, Y: 28})
	route.walkTo("clear first-switch lower boulder support", Point{X: 15, Y: 27})
	route.walkTo("escape lower first-switch boulder support", Point{X: 16, Y: 27})
	route.waitUntil("first-switch boulder settles above stairs", 180, func() bool {
		id, _ := rt.At(PlayerLayer, 15, 27)
		return id == 0 && rt.ObjectMotion[rt.index(15, 27)].Remaining == 0
	})
	for _, hook := range []struct {
		position Point
		settled  Point
	}{
		{Point{X: 17, Y: 27}, Point{X: 16, Y: 28}},
		{Point{X: 18, Y: 28}, Point{X: 17, Y: 28}},
		{Point{X: 19, Y: 28}, Point{X: 18, Y: 29}},
	} {
		route.walkTo("first-switch hook position", hook.position)
		route.hook("pull boulder toward first switch", -1)
		route.waitUntil("first-switch boulder settles", 160, func() bool {
			id, _ := rt.At(PlayerLayer, hook.settled.X, hook.settled.Y)
			return id == 0 && rt.ObjectMotion[rt.index(hook.settled.X, hook.settled.Y)].Remaining == 0
		})
	}
	route.waitUntil("first pressure door opens", 30, func() bool {
		return rt.pressureSwitchActive(18, 29) && rt.IsPassable(20, 28)
	})
	route.walkTo("first silver-key chest", Point{X: 22, Y: 27})
	route.waitUntil("first silver key collected", 160, func() bool {
		return rt.KeyForForeground8 == 1 && !rt.ChestOpening
	})
	route.walkTo("left of middle passage boulder", Point{X: 18, Y: 17})
	route.push("push middle passage boulder onto support", 1)
	route.push("push middle passage boulder into shaft", 1)
	route.waitUntil("middle passage boulder settles", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 21, 18)
		return id == 0 && rt.ObjectMotion[rt.index(21, 18)].Remaining == 0
	})
	route.walkTo("middle checkpoint", Point{X: 16, Y: 13})
	route.walkTo("upper checkpoint", Point{X: 29, Y: 9})
	route.walkTo("clear second-switch boulder support", Point{X: 13, Y: 10})
	route.walkTo("escape second-switch boulder shaft", Point{X: 14, Y: 10})
	route.waitUntil("second-switch boulder reaches upper shelf", 200, func() bool {
		id, _ := rt.At(PlayerLayer, 13, 13)
		return id == 0 && rt.ObjectMotion[rt.index(13, 13)].Remaining == 0
	})
	route.walkTo("right of second-switch boulder", Point{X: 14, Y: 13})
	route.push("push second-switch boulder left once", -1)
	route.push("push second-switch boulder left twice", -1)
	route.push("push second-switch boulder onto switch", -1)
	route.waitUntil("second pressure door opens", 30, func() bool {
		return rt.pressureSwitchActive(10, 14) && rt.IsPassable(9, 13)
	})
	route.walkTo("second silver-key chest", Point{X: 7, Y: 13})
	route.waitUntil("second silver key collected", 160, func() bool {
		return rt.KeyForForeground8 == 2 && !rt.ChestOpening
	})
	route.walkTo("left of upper passage boulder", Point{X: 30, Y: 9})
	for push := 0; push < 4; push++ {
		route.push("push upper passage boulder right", 1)
	}
	if id, _ := rt.At(PlayerLayer, 35, 9); id != 0 {
		t.Fatalf("upper passage boulder raw=%d at (35,9), want boulder", id)
	}
	route.walkTo("right of upper-shelf boulder", Point{X: 35, Y: 12})
	for push := 0; push < 8; push++ {
		route.push("push upper-shelf boulder left", -1)
	}
	if id, _ := rt.At(PlayerLayer, 26, 12); id != 0 {
		t.Fatalf("upper-shelf boulder raw=%d at (26,12), want boulder", id)
	}
	route.walkTo("descend beside supported upper boulder", Point{X: 27, Y: 14})
	route.walkTo("clear final upper-boulder support from side", Point{X: 26, Y: 14})
	route.walkTo("left of right-corridor boulder", Point{X: 26, Y: 15})
	for push := 0; push < 8; push++ {
		route.push("push right-corridor boulder right", 1)
	}
	if id, _ := rt.At(PlayerLayer, 35, 15); id != 0 {
		t.Fatalf("right-corridor boulder raw=%d at (35,15), want boulder", id)
	}
	route.walkTo("right of lower-passage boulder", Point{X: 35, Y: 18})
	for push := 0; push < 4; push++ {
		route.push("push lower-passage boulder left", -1)
	}
	route.waitUntil("lower-passage boulder settles", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 30, 19)
		return id == 0 && rt.ObjectMotion[rt.index(30, 19)].Remaining == 0
	})
	route.walkTo("gold-key chest", Point{X: 26, Y: 22})
	route.waitUntil("gold key collected", 160, func() bool {
		return rt.KeyForForeground9 == 1 && !rt.ChestOpening
	})
	route.walkTo("left of gold lock", Point{X: 32, Y: 20})
	route.waitUntil("gold lock and linked door open", 180, func() bool {
		return rt.KeyForForeground9 == 0 && rt.LocksOpened == 1 && !rt.LockOpening && rt.IsPassable(33, 21)
	})
	route.walkTo("right-side checkpoint", Point{X: 34, Y: 21})
	waitForKeyColumnBoulder := func() {
		route.waitUntil("third-key column boulder reaches push row", 180, func() bool {
			id, _ := rt.At(PlayerLayer, 43, 15)
			return id == 0 && rt.ObjectMotion[rt.index(43, 15)].Remaining == 0
		})
	}
	for boulder := 0; boulder < 4; boulder++ {
		waitForKeyColumnBoulder()
		route.walkTo("right of third-key column boulder", Point{X: 44, Y: 15})
		route.push("push third-key column boulder left once", -1)
		route.push("push third-key column boulder into side shaft", -1)
	}
	route.push("shift left-shaft top boulder onto side support", -1)
	waitForKeyColumnBoulder()
	route.walkTo("left of fifth key-column boulder", Point{X: 42, Y: 15})
	route.push("push fifth key-column boulder right once", 1)
	route.push("push fifth key-column boulder into right shaft", 1)
	waitForKeyColumnBoulder()
	route.walkTo("right of sixth key-column boulder", Point{X: 44, Y: 15})
	route.push("push sixth key-column boulder left once", -1)
	route.push("push sixth key-column boulder into left shaft", -1)
	waitForKeyColumnBoulder()
	route.push("park final key-column boulder on right support", 1)
	route.waitUntil("third-key side-shaft pile settles", 240, func() bool {
		pile := []Point{{X: 40, Y: 15}, {X: 41, Y: 15}, {X: 41, Y: 16}, {X: 41, Y: 17}, {X: 41, Y: 18}, {X: 44, Y: 15}, {X: 45, Y: 21}}
		for _, point := range pile {
			id, _ := rt.At(PlayerLayer, point.X, point.Y)
			motion := rt.ObjectMotion[rt.index(point.X, point.Y)]
			if id != 0 || motion.Remaining > 0 || motion.RollDX != 0 {
				return false
			}
		}
		for y := 9; y <= 15; y++ {
			id, _ := rt.At(PlayerLayer, 43, y)
			if id == 0 {
				return false
			}
		}
		return true
	})
	route.walkTo("third silver-key chest", Point{X: 43, Y: 16})
	route.waitUntil("third silver key collected", 160, func() bool {
		return rt.KeyForForeground8 == 3 && !rt.ChestOpening
	})
	for index, y := range []int{24, 25, 26} {
		route.walkTo("right of silver lock", Point{X: 37, Y: y})
		wantLocks := index + 2
		wantKeys := 2 - index
		route.waitUntil("silver lock opens", 180, func() bool {
			return rt.KeyForForeground8 == wantKeys && rt.LocksOpened == wantLocks && !rt.LockOpening
		})
	}
	if !rt.IsPassable(36, 27) {
		t.Fatal("linked door below three silver locks did not open")
	}
	if rt.BonusRemaining > 0 || !rt.BonusGateOpen || !rt.IsPassable(42, 27) {
		t.Fatalf("second secret-stage quota gems=%d bonus=%d remaining=%d open=%v", rt.VioletGems, rt.BonusValue, rt.BonusRemaining, rt.BonusGateOpen)
	}
	route.walkTo("second Angkor secret-stage goal", Point{X: 43, Y: 27})
	if !rt.ReachedGoal || !rt.GoalExitSecret {
		t.Fatalf("second secret-stage finish player=%+v reached=%v secret=%v tick=%d", rt.Player, rt.ReachedGoal, rt.GoalExitSecret, route.sourceTick)
	}
}
