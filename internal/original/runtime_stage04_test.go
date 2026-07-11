package original

import (
	"fmt"
	"testing"
)

func TestRuntimeStage04CanBeCompletedAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage04.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	// Angkor Stage 5 is revisited after obtaining the Mystic Hook. Tool level
	// 2 retains the level-1 hammer action and enables the hook puzzles.
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
				if waits < 160 && ((playerID == 0 && dy == 0) || foregroundID == 7 || rt.PlayerDead) {
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
				t.Fatalf("%s step %d/%d failed tick=%d player=%+v target=(%d,%d) raw=%d foreground=%d state=%#x health=%d tool=%d checkpoint=%d remaining=%d nearby=%v remainingGems=%v", label, step, count, sourceTick, rt.Player, targetX, targetY, playerID, foregroundID, rt.objectStateAtForTest(targetX, targetY), rt.Health, rt.SpecialItemMask, rt.CheckpointProgress, rt.BonusRemaining, nearby, runtimePointsWithRaw(rt, 1))
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
		t.Fatalf("%s not reached after %d ticks at tick=%d player=%+v health=%d group=%v remaining=%d gems=%v boulders=%v", label, maxTicks, sourceTick, rt.Player, rt.Health, rt.EnemyGateCounters, rt.BonusRemaining, runtimePointsWithRaw(rt, 1), runtimePointsWithRaw(rt, 0))
	}
	hammer := func(label string, dx, dy int, targets ...Point) {
		t.Helper()
		for busy() {
			tickUpdate()
		}
		if !rt.UseHammer(dx, dy) {
			targetID, _ := rt.At(PlayerLayer, rt.Player.X+dx, rt.Player.Y+dy)
			t.Fatalf("%s failed at tick=%d player=%+v targetRaw=%d", label, sourceTick, rt.Player, targetID)
		}
		waitUntil(label, 120, func() bool {
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

	tickUpdate()
	move("automatic entrance", 1, 0, 3)
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 5 entrance door")
	}
	move("push entrance boulder over drop", 1, 0, 1)
	waitUntil("entrance boulder settles", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 5, 20)
		return id == 0
	})
	move("enter entrance shaft", 1, 0, 1)
	move("climb entrance shaft", 0, -1, 3)
	move("clear first falling-boulder support", 1, 0, 2)
	move("retreat from first falling boulder", -1, 0, 1)
	waitUntil("first shaft boulder settles", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 7, 19)
		return id == 0
	})
	move("enter opened left connector", 1, 0, 1)
	move("climb opened left connector", 0, -1, 3)
	move("collect left vertical gem shelf", -1, 0, 2)
	move("collect upper left vertical gems", 0, -1, 2)
	move("return down left vertical shelf", 0, 1, 2)
	move("return to first hammer wall", 1, 0, 2)
	hammer("break first wall cluster", 1, 0,
		Point{X: 8, Y: 13}, Point{X: 9, Y: 13},
		Point{X: 9, Y: 14}, Point{X: 10, Y: 14}, Point{X: 11, Y: 14},
		Point{X: 9, Y: 15}, Point{X: 10, Y: 15},
	)
	move("enter first wall cluster", 1, 0, 2)
	move("stand left of shaft snake", 0, 1, 2)
	waitUntil("shaft snake enters first hammer range", 240, func() bool {
		id, _ := rt.At(PlayerLayer, 10, 15)
		idx := rt.index(10, 15)
		return id == 43 && rt.ObjectState[idx]&snakeStunMask == 0 && rt.ObjectMotion[idx].Remaining >= 12
	})
	hammer("first shaft-snake hit", 1, 0)
	move("move above shaft snake", 0, -1, 1)
	move("stand above shaft snake", 1, 0, 1)
	for hit := 2; hit <= 3; hit++ {
		waitUntil(fmt.Sprintf("shaft snake reaches source re-hit window %d", hit), 160, func() bool {
			id, _ := rt.At(PlayerLayer, 10, 15)
			idx := rt.index(10, 15)
			return id == 43 && rt.ObjectState[idx]&snakeStunMask == 8 && sourceTick%4 == 2
		})
		hammer(fmt.Sprintf("shaft-snake hit %d", hit), 0, 1)
	}
	move("return left of cleared shaft snake", -1, 0, 1)
	move("return to cleared shaft-snake row", 0, 1, 1)
	move("descend beside shifted stack", 0, 1, 3)
	move("enter shifted-stack shaft", 1, 0, 1)
	move("descend to shifted-stack support", 0, 1, 1)
	move("stand left of shifted-stack support", 1, 0, 1)
	move("dig shifted-stack support", 1, 0, 1)
	move("retreat from shifted stack", -1, 0, 1)
	waitUntil("shifted stack opens row 14", 160, func() bool {
		top, _ := rt.At(PlayerLayer, 12, 14)
		bottom, _ := rt.At(PlayerLayer, 12, 19)
		return top == EmptyRawID && bottom == 0
	})
	move("leave shifted-stack support", -1, 0, 1)
	move("climb back beside first wall cluster", 0, -1, 1)
	move("return left of first wall cluster", -1, 0, 1)
	move("climb first wall cluster", 0, -1, 4)
	move("cross opened first wall cluster", 1, 0, 4)
	move("checkpoint 1", 0, -1, 2)
	if rt.CheckpointProgress != 2 {
		t.Fatalf("checkpoint progress=%d at checkpoint 1, want 2", rt.CheckpointProgress)
	}

	// Collect the two gems sandwiched into the left boulder columns. Removing
	// each gem lets the boulder above settle onto the lower stack.
	move("descend central dirt shaft", 0, 1, 4)
	move("collect right stack gem", 1, 0, 1)
	move("escape right stack", -1, 0, 1)
	move("descend beside left stack gem", 0, 1, 1)
	move("collect shifted left stack gem", -1, 0, 1)
	move("escape shifted left stack", 1, 0, 1)
	move("return to checkpoint 1", 0, -1, 5)

	// Return to the quota-gate corridor and use the hook to pull (8,10) into
	// the open shaft. It falls clear of the upper key-and-gem branch.
	move("return below checkpoint 1", 0, 1, 2)
	move("return through first wall cluster", -1, 0, 4)
	move("climb to quota-gate corridor", 0, -1, 1)
	move("return to hook shaft", -1, 0, 2)
	move("climb hook shaft", 0, -1, 3)
	move("stand left of hook boulder", -1, 0, 1)
	hook("pull quota-corridor boulder", 1)
	waitUntil("hooked quota-corridor boulder settles", 160, func() bool {
		top, _ := rt.At(PlayerLayer, 7, 10)
		bottom, _ := rt.At(PlayerLayer, 7, 18)
		return top == EmptyRawID && bottom == 0
	})
	move("enter upper key branch", 1, 0, 5)
	move("climb upper key branch", 0, -1, 2)
	move("collect left key-branch gem", -1, 0, 1)
	move("collect upper-left key-branch gems", 0, -1, 2)
	move("collect upper key-branch row", 1, 0, 2)
	move("escape falling key-branch boulder", 0, 1, 1)
	move("open gold-key chest", 1, 0, 3)
	waitUntil("gold-key chest reward", 100, func() bool {
		return rt.KeyForForeground9 == 1 && !rt.ChestOpening
	})

	// Open the later return route without crossing into checkpoint 5: drop
	// (17,6), then pull (18,6) left so both settle in the vertical pocket.
	move("dig support below first return boulder", 1, 0, 2)
	move("retreat from first return boulder", -1, 0, 1)
	waitUntil("first return boulder settles", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 17, 8)
		return id == 0
	})
	move("stand left of second return boulder", 0, -1, 1)
	hook("pull second return boulder", 1)
	waitUntil("second return boulder settles", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 17, 7)
		top, _ := rt.At(PlayerLayer, 18, 6)
		return id == 0 && top == EmptyRawID
	})
	move("leave return-boulder pocket", -1, 0, 1)
	move("descend through gold-key chest", 0, 1, 1)
	move("return to fallen-boulder bypass", -1, 0, 2)
	move("descend below fallen key-branch boulder", 0, 1, 1)
	move("cross below fallen key-branch boulder", -1, 0, 2)
	move("descend key-branch shaft", 0, 1, 2)
	move("return to quota corridor", -1, 0, 4)

	// Re-enter checkpoint 1 and solve the two-row boulder gate. Clearing both
	// (15,11)/(15,12) makes (16,11) roll left; (17,11) can then be pushed to
	// x=19, where it drops and opens the route into the central temple.
	move("descend from quota corridor", 0, 1, 3)
	move("enter first wall cluster again", 1, 0, 2)
	move("descend through first wall cluster", 0, 1, 1)
	move("cross to checkpoint shaft", 1, 0, 4)
	move("checkpoint 1 revisit", 0, -1, 2)
	move("clear lower roll support", 1, 0, 2)
	move("retreat from lower roll support", -1, 0, 1)
	move("enter upper roll row", 0, -1, 1)
	move("clear upper roll support", 1, 0, 1)
	move("retreat from upper roll support", -1, 0, 1)
	waitUntil("left roll boulder settles", 160, func() bool {
		rolled, _ := rt.At(PlayerLayer, 15, 12)
		opened, _ := rt.At(PlayerLayer, 16, 11)
		return rolled == 0 && opened == EmptyRawID
	})
	move("cross behind rolled boulder", 1, 0, 2)
	move("push upper gate boulder over drop", 1, 0, 2)
	waitUntil("upper gate boulder settles", 160, func() bool {
		settled, _ := rt.At(PlayerLayer, 19, 13)
		opened, _ := rt.At(PlayerLayer, 19, 11)
		return settled == 0 && opened == EmptyRawID
	})
	move("cross upper gate", 1, 0, 1)
	move("descend into central temple", 0, 1, 1)
	move("enter central temple", 1, 0, 1)
	move("descend central temple", 0, 1, 2)
	move("dig east through central temple", 1, 0, 1)
	move("descend beside checkpoint-2 route", 0, 1, 2)
	move("return to checkpoint-2 row", 0, -1, 1)
	move("cross toward checkpoint 2", 1, 0, 2)
	move("descend to checkpoint 2", 0, 1, 1)
	move("checkpoint 2", 1, 0, 2)
	if rt.CheckpointProgress != 3 {
		t.Fatalf("checkpoint progress=%d at checkpoint 2, want 3", rt.CheckpointProgress)
	}

	// Drop (27,15), push it aside, and enter the central seven-gem grid. Four
	// gems are enough for the source quota; the outer three remain replay loot.
	move("dig below central-grid boulder", 1, 0, 2)
	move("retreat from central-grid boulder", -1, 0, 1)
	waitUntil("central-grid boulder settles", 100, func() bool {
		id, _ := rt.At(PlayerLayer, 27, 16)
		return id == 0
	})
	move("push central-grid boulder right", 1, 0, 1)
	move("climb into central gem grid", 0, -1, 3)
	move("collect left central-grid gem", -1, 0, 2)
	move("collect right central-grid gem", 1, 0, 4)
	move("return to center of central grid", -1, 0, 2)
	move("enter upper central-grid row", 0, -1, 1)
	move("collect near central-grid gem", -1, 0, 1)
	move("escape falling central-grid boulder", 1, 0, 1)
	move("leave central gem grid", 0, 1, 4)
	move("push central-grid boulder to corridor end", 1, 0, 4)
	move("enter checkpoint-3 corridor", 0, -1, 1)
	move("approach side gem shaft", 1, 0, 6)
	move("collect side-shaft gems", 0, -1, 4)
	move("return from side-shaft gems", 0, 1, 4)
	move("stand left of checkpoint-3 snake", 1, 0, 2)
	hammer("stun checkpoint-3 snake", 1, 0)
	move("cross checkpoint-3 snake", 1, 0, 2)
	move("checkpoint 3", 0, 1, 1)
	if rt.CheckpointProgress != 4 {
		t.Fatalf("checkpoint progress=%d at checkpoint 3, want 4", rt.CheckpointProgress)
	}

	// Pull the upper-right boulder onto x=39, circle above it, and push it
	// right once. It falls onto switch 0 and holds the return door open.
	move("descend beside pressure boulder", 0, 1, 1)
	move("stand right of pressure boulder", -1, 0, 1)
	hook("pull pressure boulder", -1)
	waitUntil("pressure boulder settles at x39", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 39, 18)
		return id == 0
	})
	move("circle above pressure boulder", -1, 0, 2)
	move("dig left of pressure boulder", 0, 1, 1)
	move("push pressure boulder right", 1, 0, 1)
	waitUntil("pressure boulder activates switch 0", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 40, 20)
		return id == 0 && rt.IsPassable(42, 19)
	})
	move("return to pressure door", 1, 0, 2)
	move("descend to pressure door", 0, 1, 1)
	move("cross pressure door", 1, 0, 2)

	// The far-right shaft contains the stage's health chest and two quota gems.
	move("climb into far-right shaft", 0, -1, 1)
	move("cross to health chest", 1, 0, 4)
	move("open far-right health chest", 0, 1, 1)
	waitUntil("far-right health reward", 100, func() bool {
		return rt.Health == rt.MaxHealth && !rt.ChestOpening
	})
	move("return to far-right shaft", 0, -1, 2)
	move("collect lower far-right gem", -1, 0, 1)
	move("return right of lower far-right gem", 1, 0, 1)
	move("climb to upper far-right gem", 0, -1, 1)
	move("collect upper far-right gem", -1, 0, 1)
	move("return to far-right climb", 1, 0, 1)
	move("climb far-right shaft", 0, -1, 8)
	move("finish far-right climb after knockback", 0, -1, 1)
	move("checkpoint 4", -1, 0, 1)
	if rt.CheckpointProgress != 5 {
		t.Fatalf("checkpoint progress=%d at checkpoint 4, want 5", rt.CheckpointProgress)
	}

	// Break the six-cell upper wall from the right. The gems settle on row 7
	// with boulders above them, so traverse left continuously while holding each
	// boulder long enough to collect all six.
	move("stand below upper wall cluster", -1, 0, 3)
	upperGemStart := rt.VioletGems
	upperRetries := rt.Retries
	if !rt.UseHammer(0, -1) {
		t.Fatal("failed to start upper wall-cluster hammer action")
	}
	for column := 43; column >= 38; column-- {
		want := upperGemStart + (44 - column)
		waitUntil(fmt.Sprintf("upper wall gem x=%d", column), 80, func() bool {
			return rt.VioletGems >= want
		})
		if column > 38 {
			move(fmt.Sprintf("advance below upper wall from x=%d", column), -1, 0, 1)
		}
	}
	if rt.VioletGems-upperGemStart != 6 {
		t.Fatalf("upper wall gems collected=%d, want 6", rt.VioletGems-upperGemStart)
	}
	if rt.Retries != upperRetries {
		t.Fatalf("upper wall traversal consumed %d retries", rt.Retries-upperRetries)
	}
	move("approach upper divider", -1, 0, 3)
	move("climb above upper divider", 0, -1, 4)
	move("cross above upper divider hazard", -1, 0, 2)
	move("descend left of upper divider", 0, 1, 4)
	move("approach final wall", -1, 0, 2)
	hammer("break final return wall", -1, 0, Point{X: 30, Y: 7})
	move("cross final return wall", -1, 0, 4)
	if rt.CheckpointProgress != 6 || rt.ActiveEnemyGateGroup != 0 || rt.Player != (Point{X: 27, Y: 7}) {
		t.Fatalf("checkpoint 5 player=%+v progress=%d activeGroup=%d, want (27,7)/6/0", rt.Player, rt.CheckpointProgress, rt.ActiveEnemyGateGroup)
	}

	// Drop (25,3) through the horizontal snake corridor. The grouped snake's
	// death opens the left door and releases the route back to the gold lock.
	move("climb to grouped-snake boulder support", 0, -1, 3)
	move("dig first grouped-snake boulder support", -1, 0, 1)
	waitUntil("grouped snake enters boulder timing window", 160, func() bool {
		idx := rt.index(23, 6)
		return rt.EnemyGateGroup[idx] == 0 && rt.ObjectMotion[idx].Remaining == 3 && rt.ObjectState[idx]&objectDirectionMask == 0
	})
	move("dig final grouped-snake boulder support", -1, 0, 1)
	move("retreat from grouped-snake boulder", 1, 0, 1)
	waitUntil("grouped snake crushed", 240, func() bool {
		return rt.EnemyGateCounters[0] == 0 && rt.IsPassable(22, 7)
	})
	move("return toward grouped-snake corridor", 0, 1, 2)
	move("circle right of settled grouped boulder", 1, 0, 1)
	move("return to grouped-snake corridor", 0, 1, 1)
	move("cross opened grouped-snake door", -1, 0, 5)
	if rt.Player.X == 22 {
		move("finish crossing grouped-snake door after snake timing", -1, 0, 1)
	}
	move("enter upper return row", 0, -1, 1)
	move("cross opened return-boulder route", -1, 0, 5)
	move("return below gold-key chest", -1, 0, 1)
	move("descend to key-branch shelf", 0, 1, 1)
	move("approach falling key-branch boulder", -1, 0, 4)
	move("descend ahead of falling key-branch boulder", 0, 1, 3)
	move("escape falling key-branch boulder", -1, 0, 4)
	move("return to gold lock", -1, 0, 2)
	move("stand beside gold lock", 0, -1, 1)
	waitUntil("gold lock and linked door", 120, func() bool {
		return rt.KeyForForeground9 == 0 && rt.LocksOpened == 1 && rt.IsPassable(4, 10)
	})
	move("return to final corridor", 0, 1, 1)
	if rt.VioletGems != 30 || rt.BonusRemaining != 0 || !rt.BonusGateOpen {
		t.Fatalf("Stage 5 quota violet=%d remaining=%d open=%v, want 30/0/true; remaining gems=%v", rt.VioletGems, rt.BonusRemaining, rt.BonusGateOpen, runtimePointsWithRaw(rt, 1))
	}
	move("enter Stage 5 goal", -1, 0, 3)
	if rt.Player != (Point{X: 2, Y: 10}) || !rt.ReachedGoal {
		t.Fatalf("Stage 5 finish player=%+v reached=%v tick=%d", rt.Player, rt.ReachedGoal, sourceTick)
	}
}
