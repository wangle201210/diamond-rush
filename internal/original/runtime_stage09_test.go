package original

import "testing"

func newStage09Route(t *testing.T) *stage07Route {
	t.Helper()
	stage := mustLoadOriginalStage(t, "stage09.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 2
	return &stage07Route{t: t, rt: rt}
}

func TestRuntimeStage09SecretStageCanBeCompletedAtSourceCadence(t *testing.T) {
	route := newStage09Route(t)
	rt := route.rt
	route.tick()
	route.walkTo("automatic entrance before door", Point{X: 3, Y: 6})
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 10 source entrance door")
	}
	route.walkTo("entrance checkpoint", Point{X: 4, Y: 6})
	for _, point := range []Point{{X: 6, Y: 6}, {X: 7, Y: 6}, {X: 6, Y: 7}, {X: 7, Y: 7}, {X: 8, Y: 7}} {
		route.walkTo("clear first-switch boulder shaft", point)
	}
	route.walkTo("first-switch hook position", Point{X: 8, Y: 4})
	route.hook("pull first-switch boulder into shaft", -1)
	route.waitUntil("first-switch boulder settles", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 7, 7)
		return id == 0 && rt.ObjectMotion[rt.index(7, 7)].Remaining == 0
	})
	route.walkTo("left of first-switch boulder", Point{X: 6, Y: 7})
	route.push("push first-switch boulder right once", 1)
	route.push("push first-switch boulder onto switch", 1)
	if id, _ := rt.At(PlayerLayer, 9, 7); id != 0 {
		t.Fatalf("first pressure switch raw=%d, want boulder", id)
	}
	route.walkTo("checkpoint beyond first pressure door", Point{X: 12, Y: 7})
	route.walkTo("clear second-switch boulder support", Point{X: 9, Y: 11})
	route.walkTo("escape from under second-switch boulder", Point{X: 10, Y: 11})
	route.waitUntil("second-switch boulder reaches first stair", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 9, 13)
		return id == 0 && rt.ObjectMotion[rt.index(9, 13)].Remaining == 0
	})
	route.walkTo("first stair hook position", Point{X: 11, Y: 13})
	route.hook("pull boulder onto second stair", -1)
	route.waitUntil("second-switch boulder reaches second stair", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 10, 14)
		return id == 0 && rt.ObjectMotion[rt.index(10, 14)].Remaining == 0
	})
	route.walkTo("second stair hook position", Point{X: 12, Y: 14})
	route.hook("pull boulder into second-switch shaft", -1)
	route.waitUntil("second-switch boulder settles left of switch", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 11, 17)
		return id == 0 && rt.ObjectMotion[rt.index(11, 17)].Remaining == 0
	})
	route.walkTo("left of second-switch boulder", Point{X: 10, Y: 17})
	route.push("push boulder onto second switch", 1)
	route.walkTo("second checkpoint", Point{X: 16, Y: 17})
	for _, gem := range []Point{{X: 4, Y: 10}, {X: 4, Y: 11}, {X: 5, Y: 11}, {X: 4, Y: 12}} {
		route.walkTo("collect return-route quota gem", gem)
	}
	route.walkTo("third-switch first hook position", Point{X: 22, Y: 18})
	route.hook("pull third-switch boulder onto first stair", -1)
	route.waitUntil("third-switch boulder reaches first stair", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 21, 19)
		return id == 0 && rt.ObjectMotion[rt.index(21, 19)].Remaining == 0
	})
	route.walkTo("third-switch second hook position", Point{X: 23, Y: 19})
	route.hook("pull third-switch boulder onto second stair", -1)
	route.waitUntil("third-switch boulder reaches second stair", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 22, 20)
		return id == 0 && rt.ObjectMotion[rt.index(22, 20)].Remaining == 0
	})
	route.walkTo("third-switch third hook position", Point{X: 24, Y: 20})
	route.hook("pull third-switch boulder onto third stair", -1)
	route.waitUntil("third-switch boulder reaches third stair", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 23, 21)
		return id == 0 && rt.ObjectMotion[rt.index(23, 21)].Remaining == 0
	})
	route.walkTo("third-switch final hook position", Point{X: 25, Y: 21})
	route.hook("pull third-switch boulder onto switch", -1)
	route.waitUntil("third-switch boulder settles on switch", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 24, 23)
		return id == 0 && rt.ObjectMotion[rt.index(24, 23)].Remaining == 0
	})
	route.walkTo("silver-key chest", Point{X: 27, Y: 23})
	route.waitUntil("silver key collected", 160, func() bool {
		return rt.KeyForForeground8 == 1 && !rt.ChestOpening
	})
	route.walkTo("left of silver lock", Point{X: 30, Y: 24})
	route.waitUntil("silver lock and linked door open", 180, func() bool {
		return rt.KeyForForeground8 == 0 && rt.LocksOpened == 1 && !rt.LockOpening && rt.IsPassable(31, 25)
	})
	route.walkTo("right-shaft checkpoint", Point{X: 32, Y: 12})
	for _, hook := range []struct {
		position Point
		label    string
	}{
		{Point{X: 24, Y: 7}, "pull gold-room boulder left once"},
		{Point{X: 23, Y: 7}, "pull gold-room boulder left twice"},
		{Point{X: 22, Y: 7}, "pull gold-room boulder left three times"},
		{Point{X: 21, Y: 7}, "pull gold-room boulder into stair shaft"},
	} {
		route.walkTo(hook.label+" position", hook.position)
		route.hook(hook.label, 1)
	}
	route.waitUntil("gold-room boulder reaches first stair", 180, func() bool {
		id, _ := rt.At(PlayerLayer, 22, 10)
		return id == 0 && rt.ObjectMotion[rt.index(22, 10)].Remaining == 0
	})
	for _, hook := range []struct {
		position Point
		settled  Point
		label    string
	}{
		{Point{X: 24, Y: 10}, Point{X: 23, Y: 11}, "pull gold-room boulder onto second stair"},
		{Point{X: 25, Y: 11}, Point{X: 24, Y: 12}, "pull gold-room boulder onto third stair"},
		{Point{X: 26, Y: 12}, Point{X: 25, Y: 13}, "pull gold-room boulder onto fourth stair"},
		{Point{X: 27, Y: 13}, Point{X: 26, Y: 15}, "pull gold-room boulder onto switch"},
	} {
		route.walkTo(hook.label+" position", hook.position)
		route.hook(hook.label, -1)
		route.waitUntil(hook.label+" settles", 180, func() bool {
			id, _ := rt.At(PlayerLayer, hook.settled.X, hook.settled.Y)
			return id == 0 && rt.ObjectMotion[rt.index(hook.settled.X, hook.settled.Y)].Remaining == 0
		})
	}
	if !rt.pressureSwitchActive(26, 15) || !rt.IsPassable(28, 15) {
		t.Fatalf("gold-room pressure route active=%v door-passable=%v", rt.pressureSwitchActive(26, 15), rt.IsPassable(28, 15))
	}
	route.walkTo("gold-key chest", Point{X: 29, Y: 15})
	route.waitUntil("gold key collected", 160, func() bool {
		return rt.KeyForForeground9 == 1 && !rt.ChestOpening
	})
	for _, gem := range []Point{
		{X: 29, Y: 11}, {X: 30, Y: 11},
		{X: 18, Y: 12}, {X: 19, Y: 12}, {X: 20, Y: 12},
		{X: 28, Y: 12}, {X: 29, Y: 12}, {X: 30, Y: 12}, {X: 31, Y: 12},
		{X: 18, Y: 13}, {X: 20, Y: 13},
	} {
		route.walkTo("collect safe quota gem", gem)
	}
	if rt.BonusRemaining != 0 || !rt.BonusGateOpen || !rt.IsPassable(36, 11) {
		t.Fatalf("secret-stage quota remaining=%d open=%v passable=%v", rt.BonusRemaining, rt.BonusGateOpen, rt.IsPassable(36, 11))
	}
	route.walkTo("left of exit gold lock", Point{X: 34, Y: 10})
	route.waitUntil("exit gold lock and linked door open", 180, func() bool {
		return rt.KeyForForeground9 == 0 && rt.LocksOpened == 2 && !rt.LockOpening && rt.IsPassable(35, 11)
	})
	route.walkTo("first Angkor secret-stage goal", Point{X: 42, Y: 11})
	if !rt.ReachedGoal || !rt.GoalExitSecret {
		t.Fatalf("secret-stage finish player=%+v reached=%v secret=%v tick=%d", rt.Player, rt.ReachedGoal, rt.GoalExitSecret, route.sourceTick)
	}
}
