package original

import "testing"

func TestRuntimeStage06DistinguishesNormalAndSecretExits(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage06.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 18, Y: 4}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.TryMove(1, 0) || !rt.ReachedGoal || rt.GoalExitSecret {
		t.Fatalf("normal exit reached=%v secret=%v player=%+v", rt.ReachedGoal, rt.GoalExitSecret, rt.Player)
	}

	rt, err = NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.Player = Point{X: 20, Y: 40}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.TryMove(1, 0) || !rt.ReachedGoal || !rt.GoalExitSecret {
		t.Fatalf("secret exit reached=%v secret=%v player=%+v", rt.ReachedGoal, rt.GoalExitSecret, rt.Player)
	}
	rt.SaveSnapshot()
	rt.GoalExitSecret = false
	if !rt.RestoreCheckpoint() || !rt.GoalExitSecret {
		t.Fatal("checkpoint snapshot did not preserve secret-exit state")
	}
}

func TestRuntimeStage06NormalExitCanBeCompletedAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage06.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 2
	sourceTick := 0
	tickUpdate := func() {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
		rt.TickBreakables()
		rt.TickForegroundTriggers()
		if rt.PlayerMotion.Remaining > 0 {
			rt.AdvancePlayerMotion()
		}
	}
	busy := func() bool {
		return rt.PlayerMotion.Remaining > 0 || rt.HurtTicks > 0 || rt.ChestOpening || rt.LockOpening || rt.Hammering || rt.Hooking || rt.RecallPending || rt.EnemyGateDemoActive
	}
	move := func(label string, dx, dy, count int) {
		t.Helper()
		for step := 1; step <= count; step++ {
			waits := 0
			for {
				if busy() {
					tickUpdate()
					continue
				}
				if rt.TryMove(dx, dy) {
					break
				}
				targetX, targetY := rt.Player.X+dx, rt.Player.Y+dy
				playerID, _ := rt.At(PlayerLayer, targetX, targetY)
				foregroundID, _ := rt.At(ForegroundLayer, targetX, targetY)
				if waits < 200 && ((playerID == 0 && dy == 0) || foregroundID == 7 || foregroundID == 9) {
					waits++
					tickUpdate()
					continue
				}
				t.Fatalf("%s step %d/%d failed tick=%d player=%+v target=(%d,%d) raw=%d foreground=%d health=%d key=%d gems=%d motion=%+v", label, step, count, sourceTick, rt.Player, targetX, targetY, playerID, foregroundID, rt.Health, rt.KeyForForeground9, rt.VioletGems, rt.PlayerMotion)
			}
			for rt.PlayerMotion.Remaining > 0 {
				tickUpdate()
			}
			tickUpdate()
		}
	}
	waitUntil := func(label string, maxTicks int, condition func() bool) {
		t.Helper()
		for tick := 0; tick < maxTicks; tick++ {
			if condition() {
				return
			}
			tickUpdate()
		}
		t.Fatalf("%s not reached after %d ticks at tick=%d player=%+v key=%d gems=%d boulders=%v", label, maxTicks, sourceTick, rt.Player, rt.KeyForForeground9, rt.VioletGems, runtimePointsWithRaw(rt, 0))
	}
	hook := func(label string, dx int) {
		t.Helper()
		for busy() {
			tickUpdate()
		}
		if !rt.UseHook(dx, 0) {
			t.Fatalf("%s failed at tick=%d player=%+v tool=%d", label, sourceTick, rt.Player, rt.SpecialItemMask)
		}
		waitUntil(label, 120, func() bool { return !rt.Hooking })
	}
	hammer := func(label string, dx, dy int) {
		t.Helper()
		for busy() {
			tickUpdate()
		}
		target := Point{X: rt.Player.X + dx, Y: rt.Player.Y + dy}
		if !rt.UseHammer(dx, dy) {
			t.Fatalf("%s failed at tick=%d player=%+v tool=%d", label, sourceTick, rt.Player, rt.SpecialItemMask)
		}
		waitUntil(label, 120, func() bool {
			id, _ := rt.At(PlayerLayer, target.X, target.Y)
			return !rt.Hammering && id != 30
		})
	}
	tickUpdate()
	move("automatic entrance before door", 1, 0, 5)
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 7 entrance door")
	}
	move("automatic entrance checkpoint", 1, 0, 1)
	move("cross top corridor", 1, 0, 9)
	move("dig entrance shaft", 0, 1, 4)

	// The quota needs only one violet gem. Route around the solid divider into
	// the upper-left room, collect one, then return to the hook position.
	move("descend around upper divider", 0, 1, 2)
	move("cross below upper divider", -1, 0, 2)
	move("climb around upper divider", 0, -1, 2)
	move("cross upper-left shelf", -1, 0, 4)
	move("enter upper-left gem row", 0, 1, 1)
	move("collect quota gem", -1, 0, 2)
	if rt.VioletGems != 1 || rt.BonusRemaining != 34 || rt.BonusGateOpen {
		t.Fatalf("quota state gems=%d remaining=%d open=%v, want 1/34/false", rt.VioletGems, rt.BonusRemaining, rt.BonusGateOpen)
	}

	// The source target is 35 points rather than 35 loose gems. Open the
	// 20-point chest behind the weak wall, then later the 15-point key-room
	// chest. The extra loose gem demonstrates that values clamp at zero.
	move("cross from upper-left gem row", 1, 0, 3)
	move("descend to weak wall", 0, 1, 8)
	move("stand left of weak wall", 1, 0, 1)
	hammer("break bonus-room wall", 1, 0)
	move("enter bonus-room shaft", 1, 0, 5)
	move("dig down to 20-point chest", 0, 1, 2)
	move("open 20-point chest", 1, 0, 1)
	waitUntil("20-point chest reward", 100, func() bool {
		return rt.BonusValue == 20 && !rt.ChestOpening
	})
	move("leave 20-point chest", -1, 0, 1)
	move("climb bonus-room shaft", 0, -1, 2)
	move("return from weak wall", -1, 0, 6)
	move("climb back to upper shelf row", 0, -1, 8)
	move("return to upper-left shelf", -1, 0, 1)
	move("climb onto upper-left shelf", 0, -1, 1)
	move("cross upper-left shelf again", 1, 0, 4)
	move("descend around upper divider again", 0, 1, 2)
	move("return to entrance shaft below divider", 1, 0, 2)
	move("climb to hook position", 0, -1, 2)

	// Pull the lower boulder out of the two-cell stack. It drops to (16,11).
	// Circle to its right and pull it once more onto pressure switch 1.
	hook("pull pressure-switch boulder", 1)
	waitUntil("hooked boulder settles", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 16, 11)
		return id == 0
	})
	move("descend above hooked boulder", 0, 1, 2)
	move("circle right of hooked boulder", 1, 0, 3)
	move("descend right of hooked boulder", 0, 1, 1)
	hook("pull boulder onto pressure switch", -1)
	waitUntil("pressure door opens", 160, func() bool {
		boulder, _ := rt.At(PlayerLayer, 17, 11)
		return boulder == 0 && rt.IsPassable(19, 11)
	})
	move("cross pressure door and open key chest", 1, 0, 2)
	waitUntil("gold-key chest reward", 100, func() bool {
		return rt.KeyForForeground9 == 1 && !rt.ChestOpening
	})
	move("open 15-point chest", 1, 0, 1)
	waitUntil("15-point chest reward", 100, func() bool {
		return rt.BonusValue == 35 && rt.BonusGateOpen && !rt.ChestOpening
	})

	move("return through pressure door", -1, 0, 3)
	move("climb from pressure-door row", 0, -1, 1)
	move("return to entrance shaft", -1, 0, 3)
	move("climb entrance shaft", 0, -1, 6)
	move("stand beside upper gold lock", 0, -1, 1)
	waitUntil("gold lock opens top door", 120, func() bool {
		return rt.KeyForForeground9 == 0 && rt.LocksOpened == 1 && rt.IsPassable(16, 4)
	})
	move("return to top corridor", 0, 1, 1)
	move("cross lock, quota, and normal goal", 1, 0, 4)
	if rt.Player != (Point{X: 19, Y: 4}) || !rt.ReachedGoal || rt.GoalExitSecret {
		t.Fatalf("Stage 7 normal finish player=%+v reached=%v secret=%v tick=%d", rt.Player, rt.ReachedGoal, rt.GoalExitSecret, sourceTick)
	}
}

func TestRuntimeStage06SecretExitCanBeCompletedAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage06.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.SpecialItemMask = 2
	sourceTick := 0
	tickUpdate := func() {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
		rt.TickBreakables()
		rt.TickForegroundTriggers()
		if rt.PlayerMotion.Remaining > 0 {
			rt.AdvancePlayerMotion()
		}
	}
	busy := func() bool {
		return rt.PlayerMotion.Remaining > 0 || rt.HurtTicks > 0 || rt.ChestOpening || rt.LockOpening || rt.Hammering || rt.Hooking || rt.RecallPending || rt.EnemyGateDemoActive
	}
	move := func(label string, dx, dy, count int) {
		t.Helper()
		for step := 1; step <= count; step++ {
			waits := 0
			for {
				if busy() {
					tickUpdate()
					continue
				}
				if rt.TryMove(dx, dy) {
					break
				}
				targetX, targetY := rt.Player.X+dx, rt.Player.Y+dy
				playerID, _ := rt.At(PlayerLayer, targetX, targetY)
				foregroundID, _ := rt.At(ForegroundLayer, targetX, targetY)
				if waits < 240 && ((playerID == 0 && dy == 0) || foregroundID == 7 || rt.PlayerDead) {
					waits++
					tickUpdate()
					continue
				}
				t.Fatalf("%s step %d/%d failed tick=%d player=%+v target=(%d,%d) raw=%d foreground=%d health=%d checkpoint=%d boulders=%v redSnakes=%v", label, step, count, sourceTick, rt.Player, targetX, targetY, playerID, foregroundID, rt.Health, rt.CheckpointProgress, runtimePointsWithRaw(rt, 0), runtimePointsWithRaw(rt, 43))
			}
			for rt.PlayerMotion.Remaining > 0 {
				tickUpdate()
			}
			tickUpdate()
		}
	}
	waitUntil := func(label string, maxTicks int, condition func() bool) {
		t.Helper()
		for tick := 0; tick < maxTicks; tick++ {
			if condition() {
				return
			}
			tickUpdate()
		}
		t.Fatalf("%s not reached after %d ticks at tick=%d player=%+v health=%d boulders=%v", label, maxTicks, sourceTick, rt.Player, rt.Health, runtimePointsWithRaw(rt, 0))
	}
	hammer := func(label string, dx, dy int) {
		t.Helper()
		for busy() {
			tickUpdate()
		}
		target := Point{X: rt.Player.X + dx, Y: rt.Player.Y + dy}
		if !rt.UseHammer(dx, dy) {
			id, _ := rt.At(PlayerLayer, target.X, target.Y)
			t.Fatalf("%s failed at tick=%d player=%+v target=%+v raw=%d", label, sourceTick, rt.Player, target, id)
		}
		waitUntil(label, 160, func() bool {
			id, _ := rt.At(PlayerLayer, target.X, target.Y)
			return !rt.Hammering && id != 30
		})
	}
	hook := func(label string, dx int) {
		t.Helper()
		for busy() {
			tickUpdate()
		}
		if !rt.UseHook(dx, 0) {
			t.Fatalf("%s failed at tick=%d player=%+v boulders=%v", label, sourceTick, rt.Player, runtimePointsWithRaw(rt, 0))
		}
		waitUntil(label, 120, func() bool { return !rt.Hooking })
	}

	tickUpdate()
	move("automatic entrance before door", 1, 0, 5)
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 7 entrance door")
	}
	move("automatic entrance checkpoint", 1, 0, 1)
	move("cross top corridor", 1, 0, 9)
	move("dig entrance shaft", 0, 1, 4)
	move("descend around upper divider", 0, 1, 2)
	move("cross below upper divider", -1, 0, 2)
	move("climb around upper divider", 0, -1, 2)
	move("enter left descent", -1, 0, 3)
	move("descend through left rooms", 0, 1, 12)
	move("cross to lower checkpoint shaft", -1, 0, 3)
	move("descend to weak-wall cluster", 0, 1, 3)
	hammer("break lower weak-wall cluster", 0, 1)
	move("shift off enclosed snake column", 1, 0, 1)
	move("descend through opened weak walls", 0, 1, 4)
	move("enter lower gem room", 0, 1, 1)
	if rt.Player.Y < 28 {
		move("correct lower-snake knockback", 0, 1, 28-rt.Player.Y)
	}

	// The secret route branches at x=10 before the raw26 arena entrance.
	// Entering x=12 closes the door at x=11, so crossing the trigger and then
	// walking back left was only possible while the arena logic was missing.
	move("approach enemy gate without triggering it", 1, 0, 2)
	if rt.Player != (Point{X: 10, Y: 28}) || rt.ActiveEnemyGateGroup != -1 {
		t.Fatalf("secret branch player=%+v activeGroup=%d, want (10,28)/-1", rt.Player, rt.ActiveEnemyGateGroup)
	}
	move("descend to lower red-diamond row", 0, 1, 4)
	move("cross to lower red-diamond chest", -1, 0, 6)
	waitUntil("lower red-diamond chest", 100, func() bool {
		return rt.RedDiamonds == 1 && !rt.ChestOpening
	})
	hammer("break wall into bottom route", -1, 0)
	move("enter bottom route shaft", -1, 0, 2)
	move("descend to animated corridor", 0, 1, 3)
	move("cross animated corridor", 1, 0, 8)
	move("descend to bottom boulder bypass", 0, 1, 3)

	// Cross both supports before either boulder can pin the hero in the narrow
	// corridor. From the right, pull the second boulder four cells onto switch 0.
	move("cross bottom bypass to left shaft", -1, 0, 8)
	move("descend left shaft", 0, 1, 4)
	move("stand left of boulder supports", 1, 0, 1)
	move("dig through both boulder supports", 1, 0, 3)
	waitUntil("bottom boulder pair settles", 120, func() bool {
		left, _ := rt.At(PlayerLayer, 4, 42)
		right, _ := rt.At(PlayerLayer, 5, 42)
		return left == 0 && right == 0
	})
	move("set first hook distance", 1, 0, 1)
	hook("pull switch boulder to x6", -1)
	move("set second hook distance", 1, 0, 1)
	hook("pull switch boulder to x7", -1)
	move("set third hook distance", 1, 0, 1)
	hook("pull switch boulder to x8", -1)
	move("set fourth hook distance", 1, 0, 1)
	hook("pull switch boulder onto pressure plate", -1)
	waitUntil("secret-route pressure door opens", 160, func() bool {
		boulder, _ := rt.At(PlayerLayer, 9, 43)
		return boulder == 0 && rt.IsPassable(13, 42)
	})
	move("return right of remaining boulder", -1, 0, 5)
	move("push remaining boulder left clear", -1, 0, 3)

	move("climb left shaft after switch", 0, -1, 4)
	move("cross bottom bypass to hazard", 1, 0, 9)
	move("descend around left-facing hazard", 0, 1, 1)
	move("step behind left-facing hazard", 1, 0, 1)
	move("descend to pressure-door row", 0, 1, 3)
	move("cross secret-route pressure door", 1, 0, 7)
	move("climb to secret-exit row", 0, -1, 2)
	move("enter Stage 7 secret exit", 1, 0, 2)
	if rt.Player != (Point{X: 21, Y: 40}) || !rt.ReachedGoal || !rt.GoalExitSecret {
		t.Fatalf("Stage 7 secret finish player=%+v reached=%v secret=%v tick=%d", rt.Player, rt.ReachedGoal, rt.GoalExitSecret, sourceTick)
	}
}
