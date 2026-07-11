package original

import (
	"fmt"
	"testing"
)

func TestRuntimeStage03CanBeCompletedAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage03.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
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
		return rt.PlayerMotion.Remaining > 0 || rt.HurtTicks > 0 || rt.ChestOpening || rt.LockOpening || rt.Hammering || rt.Hooking || rt.RecallPending
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
				if waits < 120 && ((playerID == 0 && dy == 0) || foregroundID == 7 || rt.PlayerDead) {
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
				t.Fatalf("%s step %d/%d failed tick=%d player=%+v target=(%d,%d) raw=%d foreground=%d state=%d health=%d hurt=%d tool=%d checkpoint=%d remaining=%d nearby=%v remainingGems=%v", label, step, count, sourceTick, rt.Player, targetX, targetY, playerID, foregroundID, rt.objectStateAtForTest(targetX, targetY), rt.Health, rt.HurtTicks, rt.SpecialItemMask, rt.CheckpointProgress, rt.BonusRemaining, nearby, runtimePointsWithRaw(rt, 1))
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
		t.Fatalf("%s not reached after %d ticks at tick=%d player=%+v health=%d group=%v remaining=%d", label, maxTicks, sourceTick, rt.Player, rt.Health, rt.EnemyGateCounters, rt.BonusRemaining)
	}
	hammer := func(label string, dx, dy int, targets ...Point) {
		t.Helper()
		for busy() {
			tickUpdate()
		}
		if !rt.UseHammer(dx, dy) {
			targetID, _ := rt.At(PlayerLayer, rt.Player.X+dx, rt.Player.Y+dy)
			t.Fatalf("%s failed at tick=%d player=%+v targetRaw=%d tool=%d", label, sourceTick, rt.Player, targetID, rt.SpecialItemMask)
		}
		waitUntil(label, 100, func() bool {
			if rt.Hammering {
				return false
			}
			for _, target := range targets {
				id, _ := rt.At(PlayerLayer, target.X, target.Y)
				if id == 30 {
					return false
				}
			}
			return true
		})
	}

	tickUpdate()
	move("automatic entrance", 1, 0, 3)
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 4 entrance door")
	}

	// Move the entrance boulder onto pressure switch 0. It opens the first
	// indexed door and stays on the switch for the rest of the stage.
	move("approach entrance boulder", 1, 0, 2)
	move("descend above entrance boulder", 0, 1, 2)
	move("circle left of entrance boulder", -1, 0, 2)
	move("collect entrance gem", 0, 1, 1)
	move("push entrance boulder onto switch 0", 1, 0, 1)
	waitUntil("entrance boulder settles on switch 0", 100, func() bool {
		id, _ := rt.At(PlayerLayer, 5, 21)
		return id == 0 && rt.IsPassable(7, 20)
	})
	move("cross pressure door 0", 1, 0, 5)

	// The left shaft leads to checkpoint 1, then the upper temple. Collect
	// every gem in the two stacked upper pockets before continuing east.
	move("climb lower-left connector", 0, -1, 6)
	move("enter checkpoint shaft", -1, 0, 1)
	move("dig into checkpoint shaft", -1, 0, 2)
	move("collect checkpoint-shaft gems", 0, -1, 2)
	move("cross checkpoint-shaft gems", 1, 0, 1)
	move("climb remaining checkpoint-shaft gems", 0, -1, 3)
	move("checkpoint 1", 1, 0, 1)
	if rt.CheckpointProgress != 2 {
		t.Fatalf("checkpoint progress=%d at checkpoint 1, want 2", rt.CheckpointProgress)
	}
	move("climb to upper snake corridor", 0, -1, 3)
	move("cross upper snake corridor", 1, 0, 2)
	move("enter upper gem pocket", 0, -1, 1)
	move("collect first upper support gem", 1, 0, 1)
	move("dig first upper support", 1, 0, 1)
	move("collect first lower upper-pocket gem", 0, 1, 1)
	move("return to upper support row", 0, -1, 1)
	move("dig second upper support", 1, 0, 1)
	move("collect second lower upper-pocket gem", 0, 1, 1)
	move("return beside upper boulders", 0, -1, 1)
	move("cross below upper boulders", 1, 0, 3)
	move("collect right lower upper-pocket gems", 0, 1, 1)
	move("collect second right lower gem", -1, 0, 1)
	move("return above right lower gems", 1, 0, 1)
	move("leave upper boulder pocket", 1, 0, 1)
	move("checkpoint 2", 1, 0, 1)
	move("collect checkpoint-2 gem row", 1, 0, 3)
	if rt.CheckpointProgress != 3 {
		t.Fatalf("checkpoint progress=%d at checkpoint 2, want 3", rt.CheckpointProgress)
	}

	// The health chest beside the right shaft restores the damage taken while
	// crossing the upper snake corridor. Continue down past the horizontal
	// hazard line to checkpoint 3 after its opening animation completes.
	move("descend to health-chest row", 0, 1, 3)
	move("approach health chest", -1, 0, 2)
	move("open health chest", 0, 1, 1)
	waitUntil("health chest reward", 100, func() bool {
		return rt.Health == rt.MaxHealth && !rt.ChestOpening
	})
	move("leave health chest", 0, -1, 1)
	move("return to right temple shaft", 1, 0, 2)
	move("descend right temple shaft", 0, 1, 6)
	move("checkpoint 3", -1, 0, 3)
	if rt.CheckpointProgress != 4 {
		t.Fatalf("checkpoint progress=%d at checkpoint 3, want 4", rt.CheckpointProgress)
	}
	move("activate enemy gate group 0", 1, 0, 1)
	if rt.ActiveEnemyGateGroup != 0 {
		t.Fatalf("active enemy gate group=%d, want 0", rt.ActiveEnemyGateGroup)
	}
	move("return left of trigger", -1, 0, 1)
	move("climb beside grouped boulders", 0, -1, 2)
	move("push right grouped boulder left", -1, 0, 1)
	move("return to grouped supports", 1, 0, 1)
	move("descend to grouped supports", 0, 1, 1)
	move("dig grouped supports", -1, 0, 4)
	waitUntil("two grouped snakes crushed", 180, func() bool {
		return rt.EnemyGateCounters[0] == 0 && rt.IsPassable(20, 14)
	})

	// Push the isolated boulder left, let it fall to row 19, then push it
	// right twice so it drops onto switch 1 and holds the second door open.
	move("return to isolated boulder", 0, 1, 1)
	move("cross grouped gate", -1, 0, 2)
	move("descend right of isolated boulder", 0, 1, 2)
	move("push isolated boulder left", -1, 0, 1)
	waitUntil("isolated boulder settles on row 19", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 10, 19)
		return id == 0
	})
	move("circle left of isolated boulder", -1, 0, 2)
	move("descend left of isolated boulder", 0, 1, 3)
	move("push isolated boulder to switch 1", 1, 0, 2)
	waitUntil("isolated boulder settles on switch 1", 100, func() bool {
		id, _ := rt.At(PlayerLayer, 12, 20)
		return id == 0 && rt.IsPassable(14, 19)
	})
	move("cross pressure door 1", 1, 0, 7)
	if rt.CheckpointProgress != 5 || rt.Player != (Point{X: 18, Y: 19}) {
		t.Fatalf("checkpoint 4 player=%+v progress=%d, want (18,19)/5", rt.Player, rt.CheckpointProgress)
	}

	// The one-way trigger closes the return door. The lower-right pocket
	// contains raw 24, which raises the source tool level from 0 to 1.
	move("cross hammer-room trigger", 1, 0, 7)
	move("climb beside hammer", 0, -1, 1)
	move("collect Mystic Hammer", 1, 0, 1)
	waitUntil("Mystic Hammer chest reward", 100, func() bool {
		return rt.SpecialItemMask == 1 && !rt.ChestOpening
	})
	if rt.SpecialItemMask != 1 {
		t.Fatalf("tool mask=%d after raw24, want 1", rt.SpecialItemMask)
	}

	// Three hammer walls form the route to checkpoint 5.
	move("approach lower hammer wall", 1, 0, 3)
	move("descend to lower hammer wall", 0, 1, 1)
	move("enter lower raw2 overlay", 1, 0, 1)
	hammer("break lower hammer wall", 1, 0, Point{X: 31, Y: 19}, Point{X: 32, Y: 19})
	move("cross lower hammer wall", 1, 0, 2)
	move("climb below vertical hammer wall", 0, -1, 1)
	hammer("break vertical hammer wall", 0, -1, Point{X: 32, Y: 17}, Point{X: 32, Y: 16})
	move("climb vertical hammer shaft", 0, -1, 5)
	hammer("break upper hammer wall", 0, -1, Point{X: 32, Y: 12}, Point{X: 31, Y: 12})
	move("enter upper hammer corridor", 0, -1, 1)
	move("checkpoint 5", -1, 0, 3)
	if rt.CheckpointProgress != 6 {
		t.Fatalf("checkpoint progress=%d at checkpoint 5, want 6", rt.CheckpointProgress)
	}

	// Break west through two wall groups, climb the long shaft, then open
	// the final raw2/raw30 gate into checkpoint 6.
	move("stand right of first west wall", -1, 0, 1)
	hammer("break first west wall", -1, 0, Point{X: 27, Y: 12}, Point{X: 26, Y: 12})
	move("collect west corridor gem", -1, 0, 3)
	hammer("break second west wall", -1, 0, Point{X: 24, Y: 12}, Point{X: 23, Y: 12}, Point{X: 22, Y: 12})
	move("enter long west shaft", -1, 0, 4)
	move("climb long west shaft", 0, -1, 7)
	move("enter final raw2 overlay", 1, 0, 1)
	hammer("break checkpoint-6 wall", 1, 0, Point{X: 23, Y: 5}, Point{X: 24, Y: 5}, Point{X: 24, Y: 6})
	move("cross checkpoint-6 upper gate", 1, 0, 2)
	move("descend to checkpoint-6 row", 0, 1, 1)
	move("checkpoint 6", 1, 0, 1)
	if rt.CheckpointProgress != 7 {
		t.Fatalf("checkpoint progress=%d at checkpoint 6, want 7", rt.CheckpointProgress)
	}

	// Breaking the final cluster drops the (27,6) boulder onto row 8.
	// Collect the two gems above it, then push the boulder left to clear the
	// lower route into the goal chamber.
	move("descend beside final wall cluster", 0, 1, 1)
	move("stand left of final wall cluster", 1, 0, 1)
	hammer("break final wall cluster", 1, 0, Point{X: 27, Y: 7}, Point{X: 28, Y: 7}, Point{X: 29, Y: 7}, Point{X: 27, Y: 8})
	waitUntil("final objects settle on row 8", 120, func() bool {
		boulder, _ := rt.At(PlayerLayer, 27, 8)
		leftGem, _ := rt.At(PlayerLayer, 28, 8)
		rightGem, _ := rt.At(PlayerLayer, 29, 8)
		return boulder == 0 && leftGem == 1 && rightGem == 1
	})
	move("enter final gem pocket", 1, 0, 2)
	move("collect first final gem", 0, 1, 1)
	move("collect second final gem", 1, 0, 1)
	move("stand right of final boulder", -1, 0, 1)
	move("push final boulder left", -1, 0, 1)
	if rt.BonusRemaining != 0 || !rt.BonusGateOpen || rt.VioletGems != 25 {
		t.Fatalf("final quota violet=%d remaining=%d open=%v, want 25/0/true; remaining=%v", rt.VioletGems, rt.BonusRemaining, rt.BonusGateOpen, runtimePointsWithRaw(rt, 1))
	}
	move("approach final snake", 1, 0, 4)
	waitUntil("final snake enters upper cell", 120, func() bool {
		idx := rt.index(32, 7)
		return rt.PlayerLayer[idx] == 43 && rt.ObjectMotion[idx].Remaining == 0
	})
	move("pass below final snake", 1, 0, 2)
	move("climb into goal chamber", 0, -1, 3)
	hammer("stun goal snake", 1, 0)
	if state := rt.objectStateAtForTest(34, 5); state&snakeStunMask == 0 {
		t.Fatalf("goal snake state=%#x, want active stun", state)
	}
	move("enter Stage 4 goal", 1, 0, 4)
	if rt.Player != (Point{X: 37, Y: 5}) || !rt.ReachedGoal {
		t.Fatalf("Stage 4 finish player=%+v reached=%v tick=%d", rt.Player, rt.ReachedGoal, sourceTick)
	}
}
