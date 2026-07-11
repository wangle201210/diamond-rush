package original

import (
	"fmt"
	"slices"
	"testing"
)

func TestRuntimeStage02CanBeCompletedAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage02.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	sourceTick := 0
	demoPrompts := make([]int, 0, 2)
	tickUpdate := func() {
		sourceTick++
		rt.TickSourceFrame(8, sourceTick, 0)
		if prompt, ok := rt.TutorialPrompt(); ok {
			if rt.AdvanceTutorialPrompt() {
				demoPrompts = append(demoPrompts, prompt.TextIndex)
			}
		}
		if rt.PlayerMotion.Remaining > 0 {
			rt.AdvancePlayerMotion()
		}
	}
	busy := func() bool {
		return rt.PlayerMotion.Remaining > 0 || rt.HurtTicks > 0 || rt.ChestOpening || rt.LockOpening || rt.RecallPending || rt.EnemyGateDemoActive || rt.TutorialScriptActive
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
				if waits < 100 && ((playerID == 0 && dy == 0) || foregroundID == 7 || rt.PlayerDead) {
					waits++
					tickUpdate()
					continue
				}
				nearby := make([]string, 0, 9)
				for ny := targetY - 1; ny <= targetY+1; ny++ {
					for nx := targetX - 1; nx <= targetX+1; nx++ {
						id, _ := rt.At(PlayerLayer, nx, ny)
						fg, _ := rt.At(ForegroundLayer, nx, ny)
						nearby = append(nearby, fmt.Sprintf("(%d,%d)=%d/%d", nx, ny, id, fg))
					}
				}
				t.Fatalf("%s step %d/%d failed tick=%d player=%+v target=(%d,%d) raw=%d foreground=%d state=%d health=%d hurt=%d locks=%v chest=%v nearby=%v remainingGems=%v", label, step, count, sourceTick, rt.Player, targetX, targetY, playerID, foregroundID, rt.objectStateAtForTest(targetX, targetY), rt.Health, rt.HurtTicks, rt.LockOpening, rt.ChestOpening, nearby, runtimePointsWithRaw(rt, 1))
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
		t.Fatalf("%s not reached after %d ticks at tick=%d player=%+v counters=%v", label, maxTicks, sourceTick, rt.Player, rt.EnemyGateCounters)
	}

	tickUpdate()
	move("automatic entrance first step", 1, 0, 1)
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 3 entrance door")
	}
	move("automatic entrance checkpoint", 1, 0, 1)
	move("enter lower room", 1, 0, 1)
	move("cross lower-room top", 0, 1, 1)
	move("cross lower-room top right", 1, 0, 3)
	move("descend key room", 0, 1, 3)
	move("cross above lower boulder", -1, 0, 4)
	move("open silver-key chest", 0, 1, 1)
	waitUntil("silver key reward", 100, func() bool { return rt.KeyForForeground8 == 1 && !rt.ChestOpening })
	move("return above silver-key chest", 0, -1, 1)
	move("return to silver lock shaft", 1, 0, 4)
	move("climb silver lock shaft", 0, -1, 4)
	waitUntil("silver lock and linked door", 100, func() bool { return rt.LocksOpened == 1 && !rt.LockOpening && rt.IsPassable(7, 18) })
	move("descend beside silver door", 0, 1, 1)
	move("cross silver door", 1, 0, 2)

	// Circle under the sealed gold-key corridor. Digging (13,16) removes
	// the support below the first boulder stack; the boulder settles at
	// (13,17), and the violet gem above it settles at (13,16).
	move("descend into central room", 0, 1, 2)
	move("cross below central snake", 1, 0, 4)
	move("climb into right gem pocket", 0, -1, 1)
	move("collect right-pocket gems", 1, 0, 2)
	move("climb to first support", 0, -1, 3)
	move("dig first boulder support", -1, 0, 1)
	move("retreat from first boulder stack", 1, 0, 1)
	waitUntil("first boulder stack settles", 100, func() bool {
		boulder, _ := rt.At(PlayerLayer, 13, 17)
		gem, _ := rt.At(PlayerLayer, 13, 16)
		return boulder == 0 && gem == 1
	})
	move("collect first falling gem", -1, 0, 1)
	move("climb first opened shaft", 0, -1, 2)

	// The second stack must be approached from its right side. Clear both
	// dirt cells, step left before the boulder drops, then push it right to
	// expose the vertical route into the upper temple.
	move("cross second stack pocket", 1, 0, 1)
	move("climb beside second stack", 0, -1, 2)
	move("dig second boulder support", -1, 0, 1)
	move("step left of falling second stack", -1, 0, 1)
	waitUntil("second boulder stack settles", 100, func() bool {
		boulder, _ := rt.At(PlayerLayer, 13, 12)
		gem, _ := rt.At(PlayerLayer, 13, 11)
		return boulder == 0 && gem == 1
	})
	move("push second boulder right", 1, 0, 1)
	move("collect second falling gem", 0, -1, 1)
	move("enter upper temple", 0, -1, 1)

	// Route around the fixed horizontal hazard at (11,6), then descend to
	// the double-boulder support puzzle at (8,14)/(8,15).
	move("upper route north", 0, -1, 1)
	move("upper route west", -1, 0, 1)
	move("upper route climb", 0, -1, 3)
	move("detour right of horizontal hazard", 1, 0, 1)
	move("detour above horizontal hazard", 0, -1, 2)
	move("collect upper quota gem row", 1, 0, 5)
	move("return from upper quota gem row", -1, 0, 5)
	move("cross above horizontal hazard", -1, 0, 3)
	move("descend left of horizontal hazard", 0, 1, 2)
	move("enter lower quota pocket", -1, 0, 1)
	move("collect lower quota pocket gem", 0, 1, 2)
	move("leave lower quota pocket", 0, -1, 2)
	move("return right of lower quota pocket", 1, 0, 1)
	move("cross upper dirt row", -1, 0, 8)
	move("approach upper-left quota route", 1, 0, 2)
	move("climb upper-left quota route", 0, -1, 2)
	move("collect upper-left quota gem", -1, 0, 1)
	move("cross to second upper-left quota gem", 1, 0, 3)
	move("collect second upper-left quota gem", 0, 1, 1)
	move("return above second upper-left quota gem", 0, -1, 1)
	move("return to far-left shaft top", -1, 0, 4)
	move("return to far-left shaft", 0, 1, 2)
	move("descend far-left shaft", 0, 1, 4)
	move("enter lower-left shaft", 1, 0, 1)
	move("descend lower-left shaft", 0, 1, 5)
	move("approach double-boulder support", 1, 0, 4)
	move("push lower support boulder right", 1, 0, 1)
	move("retreat from double-boulder stack", -1, 0, 1)
	waitUntil("double-boulder stack settles", 100, func() bool {
		upper, _ := rt.At(PlayerLayer, 8, 14)
		lower, _ := rt.At(PlayerLayer, 8, 15)
		shifted, _ := rt.At(PlayerLayer, 9, 15)
		return upper == EmptyRawID && lower == 0 && shifted == 0
	})
	move("climb around fallen support boulder", 0, -1, 1)
	move("cross into gold-key shaft", 1, 0, 4)
	move("descend gold-key shaft", 0, 1, 3)
	waitUntil("gold key reward", 100, func() bool { return rt.KeyForForeground9 == 1 && !rt.ChestOpening })

	// Approach the original (9,17) boulder from its right, push it into the
	// unsupported (8,17), and follow it back down to the gold lock room.
	move("dig behind entrance boulder", -1, 0, 1)
	move("push entrance boulder left", -1, 0, 1)
	waitUntil("entrance boulder falls clear", 100, func() bool {
		id, _ := rt.At(PlayerLayer, 8, 17)
		return id == EmptyRawID
	})
	move("follow entrance boulder", -1, 0, 1)
	move("descend to gold lock room", 0, 1, 3)
	move("stand left of gold lock", 1, 0, 4)
	waitUntil("gold lock and linked door", 100, func() bool { return rt.LocksOpened == 2 && !rt.LockOpening && rt.IsPassable(13, 21) })
	move("descend beside gold door", 0, 1, 1)
	move("cross gold door and trigger", 1, 0, 2)
	if rt.ActiveEnemyGateGroup != 0 {
		t.Fatalf("enemy gate group=%d, want 0 after trigger", rt.ActiveEnemyGateGroup)
	}
	move("final-route checkpoint", 1, 0, 1)
	move("enter boulder chamber", 1, 0, 1)
	move("climb below boulder row", 0, -1, 2)
	move("race under three boulders", 1, 0, 4)
	waitUntil("two grouped snakes crushed", 160, func() bool { return rt.EnemyGateCounters[0] == 0 && rt.IsPassable(21, 21) })
	move("descend to opened enemy door", 0, 1, 2)
	move("enter Stage 3 goal", 1, 0, 3)
	if rt.Player != (Point{X: 23, Y: 21}) || !rt.ReachedGoal {
		t.Fatalf("Stage 3 finish player=%+v reached=%v tick=%d", rt.Player, rt.ReachedGoal, sourceTick)
	}
	if !slices.Equal(demoPrompts, []int{17, 18}) {
		t.Fatalf("Stage 3 demo prompts=%v, want [17 18]", demoPrompts)
	}
}

func (rt *Runtime) objectStateAtForTest(x, y int) int {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return 0
	}
	return rt.ObjectState[rt.index(x, y)]
}
