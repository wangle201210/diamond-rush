package original

import "testing"

func newStage12Route(t *testing.T) *stage07Route {
	t.Helper()
	stage := mustLoadOriginalStage(t, "stage12.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 8
	rt.MaxHealth = 8
	rt.Health = 8
	return &stage07Route{t: t, rt: rt}
}

func collectStage12RewardChamber(route *stage07Route) {
	route.t.Helper()
	rt := route.rt
	for {
		found := false
		for y := 0; y < rt.Height() && !found; y++ {
			for x := 40; x <= 44; x++ {
				id, _ := rt.At(PlayerLayer, x, y)
				if id != 1 {
					continue
				}
				target := Point{X: x, Y: y}
				if _, _, reachable := stage07RouteStep(rt, target); !reachable {
					continue
				}
				route.walkTo("collect right reward-chamber gem", target)
				found = true
				break
			}
		}
		if !found {
			return
		}
	}
}

func collectStage12TopQuota(route *stage07Route) {
	route.t.Helper()
	rt := route.rt
	for rt.BonusRemaining > 72 {
		found := false
		for y := 2; y <= 3 && !found; y++ {
			for x := 6; x <= 18; x++ {
				id, _ := rt.At(PlayerLayer, x, y)
				if id != 1 {
					continue
				}
				target := Point{X: x, Y: y}
				if _, _, reachable := stage07RouteStep(rt, target); !reachable {
					continue
				}
				route.walkTo("collect safe top-corridor gem", target)
				found = true
				break
			}
		}
		if !found {
			return
		}
	}
}

func TestRuntimeStage12SecretStageCanBeCompletedAtSourceCadence(t *testing.T) {
	route := newStage12Route(t)
	rt := route.rt
	route.tick()
	route.walkTo("automatic entrance before door", Point{X: 3, Y: 28})
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close final secret-stage entrance door")
	}
	route.walkTo("entrance checkpoint", Point{X: 4, Y: 28})
	route.walkTo("below second shaft boulder", Point{X: 17, Y: 22})
	route.hammer("break second shaft wall cluster", 1, 0)
	route.waitUntil("second shaft wall cluster clears", 200, func() bool {
		first, _ := rt.At(PlayerLayer, 18, 21)
		second, _ := rt.At(PlayerLayer, 19, 21)
		return first == EmptyRawID && second == EmptyRawID
	})
	route.walkTo("third shaft hook position", Point{X: 19, Y: 17})
	route.hammer("clear third shaft hook path", -1, 0)
	route.walkTo("right of fallen third shaft boulder", Point{X: 21, Y: 17})
	route.hook("pull fallen third shaft boulder right", -1)
	route.walkTo("left of fallen third shaft boulder", Point{X: 19, Y: 17})
	route.push("park fallen third shaft boulder", 1)
	route.walkTo("third shaft boulder hook position", Point{X: 19, Y: 17})
	route.hook("pull third shaft boulder right", -1)
	route.walkTo("below upper shaft wall", Point{X: 17, Y: 9})
	route.hammer("break upper shaft wall cluster", 0, -1)
	route.waitUntil("upper shaft wall cluster clears", 220, func() bool {
		first, _ := rt.At(PlayerLayer, 17, 8)
		second, _ := rt.At(PlayerLayer, 18, 8)
		return first == EmptyRawID && second == EmptyRawID
	})
	route.waitUntil("upper shaft boulder falls clear", 220, func() bool {
		id, _ := rt.At(PlayerLayer, 19, 5)
		return id == EmptyRawID
	})
	route.walkTo("first upper checkpoint", Point{X: 19, Y: 4})
	route.walkTo("second upper checkpoint", Point{X: 25, Y: 3})
	collectStage12TopQuota(route)
	route.walkTo("return to second upper checkpoint", Point{X: 25, Y: 3})
	route.walkTo("bottom of right shaft", Point{X: 32, Y: 27})
	route.walkTo("top of right reward chamber", Point{X: 40, Y: 3})
	collectStage12RewardChamber(route)
	if rt.BonusRemaining > 0 || !rt.BonusGateOpen || !rt.IsPassable(38, 27) {
		t.Fatalf("final secret-stage quota remaining=%d open=%v passable=%v", rt.BonusRemaining, rt.BonusGateOpen, rt.IsPassable(38, 27))
	}
	route.walkTo("final Angkor secret-stage goal", Point{X: 39, Y: 27})
	if !rt.ReachedGoal || !rt.GoalExitSecret {
		t.Fatalf("final secret-stage finish player=%+v reached=%v secret=%v tick=%d", rt.Player, rt.ReachedGoal, rt.GoalExitSecret, route.sourceTick)
	}
}
