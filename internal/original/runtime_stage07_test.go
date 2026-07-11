package original

import "testing"

type stage07Route struct {
	t          *testing.T
	rt         *Runtime
	sourceTick int
	action     string
}

func newStage07Route(t *testing.T, toolLevel int) *stage07Route {
	t.Helper()
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = toolLevel
	return &stage07Route{t: t, rt: rt}
}

func (route *stage07Route) tick() {
	route.sourceTick++
	route.rt.TickSourceFrame(8, route.sourceTick, 0)
	if !route.rt.IsTutorialStage() {
		if _, ok := route.rt.TutorialPrompt(); ok {
			route.rt.AdvanceTutorialPrompt()
		}
	}
	route.rt.TickBreakables()
	route.rt.TickForegroundTriggers()
	if route.rt.PlayerMotion.Remaining > 0 {
		route.rt.AdvancePlayerMotion()
	}
	if route.rt.PlayerDead {
		route.t.Fatalf("hero died during %s: tick=%d player=%+v health=%d hits=%d boulders=%v", route.action, route.sourceTick, route.rt.Player, route.rt.Health, route.rt.HitCount, runtimePointsWithRaw(route.rt, 0))
	}
}

func (route *stage07Route) busy() bool {
	rt := route.rt
	return rt.PlayerMotion.Remaining > 0 || rt.HurtTicks > 0 || rt.ChestOpening || rt.LockOpening || rt.Hammering || rt.Hooking || rt.RecallPending || rt.EnemyGateDemoActive || rt.TutorialScriptActive
}

func (route *stage07Route) waitUntil(label string, maxTicks int, condition func() bool) {
	route.t.Helper()
	route.action = label
	for tick := 0; tick < maxTicks; tick++ {
		if condition() {
			return
		}
		route.tick()
	}
	route.t.Fatalf("%s not reached after %d ticks: tick=%d player=%+v health=%d gems=%d bonus=%d keys=%d/%d locks=%d groups=%v boulders=%v snakes=%v ice=%v", label, maxTicks, route.sourceTick, route.rt.Player, route.rt.Health, route.rt.VioletGems, route.rt.BonusValue, route.rt.KeyForForeground8, route.rt.KeyForForeground9, route.rt.LocksOpened, route.rt.EnemyGateCounters, runtimePointsWithRaw(route.rt, 0), runtimePointsWithRaw(route.rt, 19), runtimePointsWithRaw(route.rt, 9))
}

func (route *stage07Route) walkTo(label string, target Point) {
	route.t.Helper()
	route.action = label
	blockedTicks := 0
	for steps := 0; route.rt.Player != target; steps++ {
		if steps >= route.rt.Width()*route.rt.Height()*2 {
			route.t.Fatalf("%s exceeded route limit: player=%+v target=%+v", label, route.rt.Player, target)
		}
		for route.busy() {
			route.tick()
		}
		dx, dy, ok := stage07RouteStep(route.rt, target)
		if !ok {
			if blockedTicks < 240 {
				blockedTicks++
				route.tick()
				continue
			}
			route.t.Fatalf("%s has no passable route: tick=%d player=%+v target=%+v keys=%d/%d locks=%d bonus=%d boulders=%v", label, route.sourceTick, route.rt.Player, target, route.rt.KeyForForeground8, route.rt.KeyForForeground9, route.rt.LocksOpened, route.rt.BonusValue, runtimePointsWithRaw(route.rt, 0))
		}
		blockedTicks = 0
		if !route.rt.TryMove(dx, dy) {
			route.tick()
			continue
		}
		for route.rt.PlayerMotion.Remaining > 0 {
			route.tick()
		}
		route.tick()
		if route.rt.PlayerDead {
			route.t.Fatalf("%s killed the hero at tick=%d player=%+v", label, route.sourceTick, route.rt.Player)
		}
	}
}

func stage07RouteStep(rt *Runtime, target Point) (int, int, bool) {
	if target.X < 0 || target.Y < 0 || target.X >= rt.Width() || target.Y >= rt.Height() {
		return 0, 0, false
	}
	start := rt.Player
	queue := []Point{start}
	seen := map[Point]bool{start: true}
	previous := map[Point]Point{}
	directions := [...]Point{{Y: -1}, {X: -1}, {X: 1}, {Y: 1}}
	found := false
	for len(queue) > 0 && !found {
		point := queue[0]
		queue = queue[1:]
		for _, direction := range directions {
			next := Point{X: point.X + direction.X, Y: point.Y + direction.Y}
			playerID, _ := rt.At(PlayerLayer, next.X, next.Y)
			if seen[next] || isContactEnemy(playerID) || !rt.IsPassable(next.X, next.Y) {
				continue
			}
			seen[next] = true
			previous[next] = point
			if next == target {
				found = true
				break
			}
			queue = append(queue, next)
		}
	}
	if !found {
		return 0, 0, false
	}
	step := target
	for previous[step] != start {
		step = previous[step]
	}
	return step.X - start.X, step.Y - start.Y, true
}

func (route *stage07Route) hammer(label string, dx, dy int) {
	route.t.Helper()
	route.action = label
	for route.busy() {
		route.tick()
	}
	if !route.rt.UseHammer(dx, dy) {
		target := Point{X: route.rt.Player.X + dx, Y: route.rt.Player.Y + dy}
		id, _ := route.rt.At(PlayerLayer, target.X, target.Y)
		route.t.Fatalf("%s failed: tick=%d player=%+v target=%+v raw=%d tool=%d", label, route.sourceTick, route.rt.Player, target, id, route.rt.SpecialItemMask)
	}
	route.waitUntil(label, 120, func() bool { return !route.rt.Hammering })
}

func (route *stage07Route) hook(label string, dx int) {
	route.t.Helper()
	route.action = label
	for route.busy() {
		route.tick()
	}
	if !route.rt.UseHook(dx, 0) {
		route.t.Fatalf("%s failed: tick=%d player=%+v tool=%d boulders=%v", label, route.sourceTick, route.rt.Player, route.rt.SpecialItemMask, runtimePointsWithRaw(route.rt, 0))
	}
	route.waitUntil(label, 160, func() bool { return !route.rt.Hooking })
}

func (route *stage07Route) push(label string, dx int) {
	route.t.Helper()
	route.action = label
	for attempt := 0; attempt < boulderPushAttempts+12; attempt++ {
		for route.busy() {
			route.tick()
		}
		if route.rt.TryMove(dx, 0) {
			for route.rt.PlayerMotion.Remaining > 0 {
				route.tick()
			}
			route.tick()
			return
		}
		route.tick()
	}
	route.t.Fatalf("%s failed: tick=%d player=%+v boulders=%v", label, route.sourceTick, route.rt.Player, runtimePointsWithRaw(route.rt, 0))
}

func TestRuntimeStage07NormalExitCanBeCompletedAtSourceCadence(t *testing.T) {
	route := newStage07Route(t, 2)
	rt := route.rt
	route.tick()
	route.walkTo("automatic entrance before door", Point{X: 4, Y: 8})
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 8 entrance door")
	}
	route.walkTo("automatic entrance checkpoint", Point{X: 5, Y: 8})
	route.walkTo("entrance-side weak wall", Point{X: 9, Y: 15})
	route.hammer("break entrance wall", 1, 0)
	route.waitUntil("entrance wall cluster clears", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 10, 13)
		return id == EmptyRawID
	})

	// Ten fixed gems plus the three 15/10/5-point chests meet raw12's
	// 40-point source quota without relying on incidental enemy movement.
	for index, point := range []Point{
		{X: 12, Y: 15}, {X: 13, Y: 15},
		{X: 20, Y: 15}, {X: 21, Y: 15}, {X: 22, Y: 15}, {X: 23, Y: 15}, {X: 24, Y: 15}, {X: 25, Y: 15},
		{X: 19, Y: 17}, {X: 20, Y: 17},
	} {
		route.walkTo("collect normal-route quota gem", point)
		if rt.VioletGems < index+1 {
			t.Fatalf("quota gem %d at %+v was not collected: gems=%d", index+1, point, rt.VioletGems)
		}
	}

	route.walkTo("first silver-key chest", Point{X: 20, Y: 20})
	route.waitUntil("first silver key reward", 120, func() bool {
		return rt.KeyForForeground8 == 1 && !rt.ChestOpening
	})
	// Pull the boulder off (16,23) before digging its support at (16,24),
	// otherwise it falls into the one-cell doorway and traps the hero.
	route.walkTo("lower bonus-room hook position", Point{X: 18, Y: 23})
	route.hook("pull lower bonus-room boulder", -1)
	route.waitUntil("lower bonus-room boulder settles", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 17, 23)
		return id == 0 && rt.ObjectMotion[rt.index(17, 23)].Remaining == 0
	})
	route.walkTo("lower bonus-room wall", Point{X: 16, Y: 24})
	route.hammer("break lower bonus-room wall", -1, 0)
	route.waitUntil("lower bonus-room wall clears", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 15, 24)
		return id == EmptyRawID
	})
	bonusBeforeLowerChest := rt.BonusValue
	route.walkTo("collect 10-point chest", Point{X: 10, Y: 24})
	route.waitUntil("10-point chest reward", 120, func() bool {
		return rt.BonusValue >= bonusBeforeLowerChest+10 && !rt.ChestOpening
	})

	route.walkTo("upper weak-wall cluster", Point{X: 17, Y: 5})
	route.hammer("break upper weak-wall cluster", -1, 0)
	route.waitUntil("upper weak-wall cluster clears", 240, func() bool {
		left, _ := rt.At(PlayerLayer, 16, 4)
		right, _ := rt.At(PlayerLayer, 22, 4)
		return left == EmptyRawID && right == EmptyRawID
	})
	// Digging (26,6) first lets the boulder at (25,5) roll over the key
	// chest and repeatedly crush the hero during its long opening animation.
	route.walkTo("right of second-key boulder", Point{X: 26, Y: 5})
	route.push("push second-key boulder onto the left stack", -1)
	route.walkTo("second silver-key chest", Point{X: 26, Y: 7})
	route.waitUntil("second silver key reward", 120, func() bool {
		return rt.KeyForForeground8 == 2 && !rt.ChestOpening
	})
	route.walkTo("upper secret-room wall", Point{X: 16, Y: 3})
	route.hammer("break upper secret-room wall", 0, -1)
	route.waitUntil("upper secret-room wall clears", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 16, 2)
		return id == EmptyRawID
	})
	route.walkTo("collect 15-point chest", Point{X: 11, Y: 5})
	route.waitUntil("15-point chest reward", 120, func() bool {
		return rt.BonusValue >= 25 && !rt.ChestOpening
	})

	route.walkTo("upper silver lock", Point{X: 26, Y: 13})
	route.waitUntil("upper silver lock opens", 160, func() bool {
		return rt.LocksOpened >= 1 && !rt.LockOpening
	})
	route.walkTo("lower silver lock", Point{X: 26, Y: 14})
	route.waitUntil("lower silver lock opens", 160, func() bool {
		return rt.LocksOpened >= 2 && !rt.LockOpening && rt.IsPassable(27, 15)
	})

	// Crossing raw26 at (32,9) closes both arena doors. Kill the red snake
	// from the safe cell above it, then time the center boulder for the green
	// snake before continuing into the pressure-switch room.
	redHammerPosition := Point{X: 33, Y: 6}
	redSnakePosition := Point{X: 33, Y: 7}
	route.walkTo("red arena snake hammer position", redHammerPosition)
	route.waitUntil("red arena snake enters first hammer cell", 240, func() bool {
		idx := rt.index(redSnakePosition.X, redSnakePosition.Y)
		return rt.PlayerLayer[idx] == 43 && rt.ObjectState[idx]&snakeStunMask == 0 && rt.ObjectMotion[idx].Remaining >= 12
	})
	route.hammer("hit red arena snake 1", 0, 1)
	for hit := 2; hit <= 3; hit++ {
		route.waitUntil("red arena snake reaches source re-hit window", 160, func() bool {
			idx := rt.index(redSnakePosition.X, redSnakePosition.Y)
			return rt.PlayerLayer[idx] == 43 && rt.ObjectState[idx]&snakeStunMask == 8 && route.sourceTick%4 == 2
		})
		route.hammer("hit red arena snake", 0, 1)
	}
	if rt.EnemyGateCounters[0] != 1 {
		t.Fatalf("red arena snake counter=%d after hammer, want 1", rt.EnemyGateCounters[0])
	}
	route.walkTo("return to arena center support", Point{X: 34, Y: 6})
	route.waitUntil("green arena snake enters boulder timing window", 240, func() bool {
		id, _ := rt.At(PlayerLayer, 37, 8)
		return id == 19 && rt.ObjectMotion[rt.index(37, 8)].Remaining == 6 && rt.ObjectState[rt.index(37, 8)]&objectDirectionMask == 0
	})
	route.hammer("release enemy-arena boulder", 1, 0)
	route.waitUntil("enemy arena opens after both snakes", 360, func() bool {
		return rt.EnemyGateCounters[0] == 0 && rt.IsPassable(38, 9)
	})

	// Move the arena boulder left off its supported column. It settles at
	// (31,15); pushing it right once then drops it onto switch 1 at (32,17).
	route.walkTo("right of arena boulder", Point{X: 34, Y: 12})
	route.push("push arena boulder to supported column", -1)
	route.push("push arena boulder into drop shaft", -1)
	route.waitUntil("arena boulder settles at lower shelf", 240, func() bool {
		id, _ := rt.At(PlayerLayer, 31, 15)
		return id == 0 && rt.ObjectMotion[rt.index(31, 15)].Remaining == 0
	})
	route.walkTo("left of lower arena boulder", Point{X: 30, Y: 15})
	route.push("push arena boulder over pressure shaft", 1)
	route.waitUntil("arena boulder activates pressure switch", 240, func() bool {
		id, _ := rt.At(PlayerLayer, 32, 17)
		return id == 0 && rt.pressureSwitchActive(32, 17) && rt.IsPassable(34, 16)
	})

	route.walkTo("collect 5-point chest", Point{X: 41, Y: 12})
	route.waitUntil("5-point chest reward", 120, func() bool {
		return rt.BonusValue >= 30 && !rt.ChestOpening
	})
	if rt.BonusRemaining != 0 || !rt.BonusGateOpen {
		t.Fatalf("normal-route quota remaining=%d open=%v, want 0/true", rt.BonusRemaining, rt.BonusGateOpen)
	}
	route.walkTo("gold-key area checkpoint", Point{X: 36, Y: 16})
	route.walkTo("gold-key chest", Point{X: 37, Y: 26})
	route.waitUntil("gold key reward", 120, func() bool {
		return rt.KeyForForeground9 == 1 && !rt.ChestOpening
	})

	route.walkTo("stand beside normal-exit gold lock", Point{X: 8, Y: 18})
	route.waitUntil("normal-exit gold lock opens", 180, func() bool {
		return rt.KeyForForeground9 == 0 && !rt.LockOpening && rt.IsPassable(7, 19)
	})
	route.walkTo("Stage 8 normal goal", Point{X: 5, Y: 19})
	if rt.Player != (Point{X: 5, Y: 19}) || !rt.ReachedGoal || rt.GoalExitSecret {
		t.Fatalf("Stage 8 normal finish player=%+v reached=%v secret=%v tick=%d", rt.Player, rt.ReachedGoal, rt.GoalExitSecret, route.sourceTick)
	}
}

func TestRuntimeStage07DistinguishesNormalAndSecretExits(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage07.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 4, Y: 19}
	if !rt.TryMove(1, 0) || !rt.ReachedGoal || rt.GoalExitSecret {
		t.Fatalf("normal exit reached=%v secret=%v player=%+v", rt.ReachedGoal, rt.GoalExitSecret, rt.Player)
	}

	rt, err = NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 3, Y: 3}
	if !rt.TryMove(1, 0) || !rt.ReachedGoal || !rt.GoalExitSecret {
		t.Fatalf("secret exit reached=%v secret=%v player=%+v", rt.ReachedGoal, rt.GoalExitSecret, rt.Player)
	}
}

func TestRuntimeStage07SecretExitRequiresFreezeHammerAtSourceCadence(t *testing.T) {
	route := newStage07Route(t, 8)
	rt := route.rt
	route.tick()
	route.walkTo("automatic secret-route entrance before door", Point{X: 4, Y: 8})
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 8 secret-route entrance door")
	}
	route.walkTo("automatic secret-route entrance checkpoint", Point{X: 5, Y: 8})
	route.walkTo("secret-route entrance weak wall", Point{X: 9, Y: 15})
	route.hammer("break secret-route entrance wall", 1, 0)
	route.waitUntil("secret-route entrance wall clears", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 10, 13)
		return id == EmptyRawID
	})

	route.walkTo("secret-route upper weak-wall cluster", Point{X: 17, Y: 5})
	route.hammer("break secret-route upper weak-wall cluster", -1, 0)
	route.waitUntil("secret-route upper weak-wall cluster clears", 240, func() bool {
		left, _ := rt.At(PlayerLayer, 16, 4)
		right, _ := rt.At(PlayerLayer, 22, 4)
		return left == EmptyRawID && right == EmptyRawID
	})
	route.walkTo("secret-room wall hammer position", Point{X: 16, Y: 3})
	route.hammer("break secret-room wall", 0, -1)
	route.waitUntil("secret-room wall clears", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 16, 2)
		return id == EmptyRawID
	})

	route.walkTo("right of secret-room snake", Point{X: 9, Y: 2})
	route.waitUntil("secret-room snake returns to upper cell", 240, func() bool {
		id, _ := rt.At(PlayerLayer, 8, 2)
		return id == 19
	})
	route.hammer("freeze secret-room snake", -1, 0)
	route.waitUntil("frozen snake settles above pressure switch", 200, func() bool {
		id, _ := rt.At(PlayerLayer, 8, 3)
		return id == 9 && rt.FrozenOriginalAt(8, 3) == 19 && rt.ObjectMotion[rt.index(8, 3)].Remaining == 0
	})
	route.walkTo("right of frozen secret-room snake", Point{X: 9, Y: 3})
	route.push("push frozen snake over secret pressure switch", -1)
	route.waitUntil("frozen snake activates secret door", 200, func() bool {
		id, _ := rt.At(PlayerLayer, 7, 4)
		return id == 9 && rt.FrozenOriginalAt(7, 4) == 19 && rt.pressureSwitchActive(7, 4) && rt.IsPassable(6, 3)
	})
	route.walkTo("Stage 8 secret goal", Point{X: 4, Y: 3})
	if rt.Player != (Point{X: 4, Y: 3}) || !rt.ReachedGoal || !rt.GoalExitSecret {
		t.Fatalf("Stage 8 secret finish player=%+v reached=%v secret=%v tick=%d", rt.Player, rt.ReachedGoal, rt.GoalExitSecret, route.sourceTick)
	}
}
