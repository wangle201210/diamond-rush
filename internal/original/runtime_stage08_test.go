package original

import (
	"slices"
	"testing"
)

func newStage08Runtime(t *testing.T) *Runtime {
	t.Helper()
	stage := mustLoadOriginalStage(t, "stage08.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	return rt
}

func TestRuntimeStage08InitializesGreatAnacondaAndSourceDoors(t *testing.T) {
	rt := newStage08Runtime(t)
	if !rt.Anaconda.Enabled || rt.Anaconda.Phase != AnacondaPhaseDormant || rt.Anaconda.Health != 3 || rt.Anaconda.Column != 0 {
		t.Fatalf("boss=%+v, want enabled dormant 3-health column 0", rt.Anaconda)
	}
	for _, point := range []Point{{X: 8, Y: 5}, {X: 18, Y: 5}} {
		state, _ := rt.At(BackgroundLayer, point.X, point.Y)
		if state != 0x30 || !rt.IsPassable(point.X, point.Y) {
			t.Errorf("source group door %+v state=%#x passable=%v, want phase 3/open", point, state, rt.IsPassable(point.X, point.Y))
		}
	}
	if marker, _ := rt.At(ForegroundLayer, 8, 4); marker != EmptyRawID {
		t.Fatalf("raw17 above left door=%d, want consumed during source door initialization", marker)
	}
	if marker, _ := rt.At(ForegroundLayer, 1, 13); marker != EmptyRawID {
		t.Fatalf("dummy snake raw17 marker=%d, want consumed during source enemy initialization", marker)
	}
	if rt.EnemyGateCounters[0] != 1 || rt.EnemyGateGroup[rt.index(1, 12)] != 0 {
		t.Fatalf("group counters=%v snake group=%d, want one group-0 dummy enemy", rt.EnemyGateCounters, rt.EnemyGateGroup[rt.index(1, 12)])
	}
}

func TestRuntimeStage08RunsDecodedBossIntroScript33(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Player = Point{X: 5, Y: 5}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter Boss intro event at (6,5)")
	}
	prompts := make([]int, 0, 2)
	started := false
	for sourceTick := 1; sourceTick <= 160; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
		if rt.TutorialScriptActive {
			started = true
			if rt.TutorialScriptID != 33 {
				t.Fatalf("Boss intro script id=%d, want 33", rt.TutorialScriptID)
			}
		}
		if prompt, ok := rt.TutorialPrompt(); ok {
			if rt.AdvanceTutorialPrompt() {
				prompts = append(prompts, prompt.TextIndex)
			}
		}
		if rt.PlayerMotion.Remaining > 0 {
			rt.AdvancePlayerMotion()
		}
		if started && !rt.TutorialScriptActive {
			break
		}
	}
	if !started || rt.TutorialScriptActive || !slices.Equal(prompts, []int{32, 33}) {
		t.Fatalf("Boss intro started=%v active=%v prompts=%v, want completed [32 33]", started, rt.TutorialScriptActive, prompts)
	}
	if rt.Anaconda.Phase != AnacondaPhaseDormant {
		t.Fatalf("Boss phase=%d during room intro, want dormant until hero reaches x=10", rt.Anaconda.Phase)
	}
}

func TestRuntimeDemoBubbleUsesSourceThirtyPixelSlide(t *testing.T) {
	rt := newStage08Runtime(t)
	if !rt.startTutorialScript(33) {
		t.Fatal("failed to start decoded Boss intro script 33")
	}
	for tick := 0; tick < 20; tick++ {
		rt.tickTutorial()
		if prompt, ok := rt.TutorialPrompt(); ok {
			if prompt.TextIndex != 32 || prompt.X != -240 {
				t.Fatalf("initial bubble=%+v, want text 32 at x=-240", prompt)
			}
			break
		}
	}
	for step := 1; step <= 9; step++ {
		rt.tickTutorial()
		prompt, ok := rt.TutorialPrompt()
		if !ok {
			t.Fatalf("bubble disappeared during slide step %d", step)
		}
		wantX := min(7, -240+step*30)
		if prompt.X != wantX {
			t.Fatalf("bubble x at slide step %d=%d, want %d", step, prompt.X, wantX)
		}
	}
	if !rt.AdvanceTutorialPrompt() {
		t.Fatal("failed to acknowledge fully entered source bubble")
	}
	for step := 1; step <= 7; step++ {
		rt.tickTutorial()
		prompt, ok := rt.TutorialPrompt()
		if !ok {
			t.Fatalf("bubble disappeared during exit step %d", step)
		}
		wantX := 7 + step*30
		if prompt.X != wantX {
			t.Fatalf("bubble x at exit step %d=%d, want %d", step, prompt.X, wantX)
		}
	}
	rt.tickTutorial()
	if prompt, ok := rt.TutorialPrompt(); ok && prompt.TextIndex == 32 {
		t.Fatalf("acknowledged bubble remains after passing x=240: %+v", prompt)
	}
}

func TestRuntimeDemoPortraitUsesSourceFiveFrameReveal(t *testing.T) {
	rt := newStage08Runtime(t)
	if !rt.startTutorialScript(33) {
		t.Fatal("failed to start decoded Boss intro script 33")
	}
	for tick := 0; tick < 20 && rt.TutorialPortraitRevealTicks == 0; tick++ {
		rt.tickTutorial()
	}
	for reveal := 1; reveal <= 5; reveal++ {
		if rt.TutorialPortraitRevealTicks != reveal || rt.TutorialPortraitVisible {
			t.Fatalf("portrait reveal tick=%d visible=%v, want reveal %d and hidden", rt.TutorialPortraitRevealTicks, rt.TutorialPortraitVisible, reveal)
		}
		rt.tickTutorial()
	}
	if rt.TutorialPortraitRevealTicks != 0 || !rt.TutorialPortraitVisible {
		t.Fatalf("portrait reveal tick=%d visible=%v after frame 5, want full portrait", rt.TutorialPortraitRevealTicks, rt.TutorialPortraitVisible)
	}
}

func TestRuntimeStage08TriggerClosesEntranceThenArenaDoor(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Player = Point{X: 8, Y: 5}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter source raw26 trigger")
	}
	left, _ := rt.At(BackgroundLayer, 8, 5)
	right, _ := rt.At(BackgroundLayer, 18, 5)
	if left != 0x30 || right != 0x30 || rt.EnemyGateDemoActive {
		t.Fatalf("doors before jInt threshold left=%#x right=%#x demo=%v, want open/open/false", left, right, rt.EnemyGateDemoActive)
	}
	rt.tickPendingForegroundEvent()
	rt.AdvancePlayerMotion()
	rt.tickPendingForegroundEvent()
	if rt.EnemyGateDemoActive {
		t.Fatal("raw26 triggered above source jInt threshold 6")
	}
	rt.AdvancePlayerMotion()
	rt.tickPendingForegroundEvent()
	left, _ = rt.At(BackgroundLayer, 8, 5)
	if left != 0 || right != 0x30 {
		t.Fatalf("trigger door states left=%#x right=%#x at jInt 6, want closed/open", left, right)
	}
	if !rt.EnemyGateDemoActive || rt.EnemyGateDemoTarget != (Point{X: 18, Y: 5}) || rt.ActiveEnemyGateGroup != 0 {
		t.Fatalf("gate demo active=%v target=%+v group=%d", rt.EnemyGateDemoActive, rt.EnemyGateDemoTarget, rt.ActiveEnemyGateGroup)
	}
	for tick := 0; tick < rt.EnemyGateDemoOutboundTicks+9; tick++ {
		rt.tickEnemyGateDemo()
	}
	right, _ = rt.At(BackgroundLayer, 18, 5)
	if right != 0x30 {
		t.Fatalf("right door closed before source demo delay: state=%#x", right)
	}
	rt.tickEnemyGateDemo()
	right, _ = rt.At(BackgroundLayer, 18, 5)
	if right != 0 {
		t.Fatalf("right door state=%#x after pan plus 10 ticks, want closed", right)
	}
	for rt.EnemyGateDemoPhase == 2 {
		rt.tickEnemyGateDemo()
	}
	if index, ticks, ok := rt.EnemyGateMessage(); !ok || index != 51 || ticks != 80 {
		t.Fatalf("boss gate message index/ticks/ok=%d/%d/%v, want 51/80/true", index, ticks, ok)
	}
	for rt.EnemyGateDemoActive {
		rt.tickEnemyGateDemo()
	}
	rt.PlayerMotion = ObjectMotion{}
	if !rt.CanAcceptInput() {
		t.Fatal("input remained locked after group-introduction demo")
	}
}

func TestRuntimeGreatAnacondaEmergesAndPlacesSourceBlocker(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Player = Point{X: 10, Y: 5}
	rt.tickGreatAnaconda(1)
	if rt.Anaconda.Phase != AnacondaPhaseDelay || rt.Anaconda.PhaseTicks != 0 {
		t.Fatalf("activation phase/ticks=%d/%d, want delay/0", rt.Anaconda.Phase, rt.Anaconda.PhaseTicks)
	}
	for tick := 2; tick <= 11; tick++ {
		rt.tickGreatAnaconda(tick)
	}
	if rt.Anaconda.Phase != AnacondaPhaseDelay {
		t.Fatalf("delay ended at tick 10: phase=%d", rt.Anaconda.Phase)
	}
	rt.tickGreatAnaconda(12)
	if rt.Anaconda.Phase != AnacondaPhaseEmerge || rt.Anaconda.Animation != 2 {
		t.Fatalf("phase/animation=%d/%d after source >10 delay, want emerge/2", rt.Anaconda.Phase, rt.Anaconda.Animation)
	}
	for tick := 1; tick <= 20; tick++ {
		rt.tickGreatAnaconda(12 + tick)
	}
	if rt.Anaconda.BlockerSet {
		t.Fatal("raw50 blocker appeared at phase tick 20, want strictly after 20")
	}
	rt.tickGreatAnaconda(33)
	if !rt.Anaconda.BlockerSet || rt.Anaconda.Blocker != (Point{X: 10, Y: 8}) {
		t.Fatalf("blocker=%+v set=%v, want (10,8)", rt.Anaconda.Blocker, rt.Anaconda.BlockerSet)
	}
	for x := 10; x <= 11; x++ {
		if id, _ := rt.At(PlayerLayer, x, 8); id != 50 {
			t.Errorf("blocker cell (%d,8) raw=%d, want raw50", x, id)
		}
	}
	for tick := 34; tick <= 52; tick++ {
		rt.tickGreatAnaconda(tick)
	}
	if rt.Anaconda.Phase != AnacondaPhaseEmerge {
		t.Fatalf("emerge ended before phase tick 41: phase=%d ticks=%d", rt.Anaconda.Phase, rt.Anaconda.PhaseTicks)
	}
	rt.tickGreatAnaconda(53)
	if rt.Anaconda.Phase != AnacondaPhaseVulnerable || rt.Anaconda.PhaseTicks != 0 {
		t.Fatalf("phase/ticks=%d/%d after source >40 emerge, want vulnerable/0", rt.Anaconda.Phase, rt.Anaconda.PhaseTicks)
	}
}

func TestRuntimeGreatAnacondaOnlyTakesOneHitFromRowsSevenAndEight(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Anaconda.Phase = AnacondaPhaseVulnerable
	rt.Anaconda.PhaseTicks = 0
	rt.SetForTest(PlayerLayer, 10, 6, 0)
	if hits, _ := rt.tickGreatAnaconda(1); hits != 0 || rt.Anaconda.Health != 3 {
		t.Fatalf("row-6 strike hits=%d health=%d, want 0/3", hits, rt.Anaconda.Health)
	}
	rt.SetForTest(PlayerLayer, 10, 7, 0)
	rt.SetForTest(PlayerLayer, 11, 8, 0)
	hits, defeated := rt.tickGreatAnaconda(2)
	if hits != 1 || defeated || rt.Anaconda.Health != 2 || rt.Anaconda.Phase != AnacondaPhaseHurt {
		t.Fatalf("strike hits=%d defeated=%v health=%d phase=%d, want 1/false/2/hurt", hits, defeated, rt.Anaconda.Health, rt.Anaconda.Phase)
	}
	for _, point := range []Point{{X: 10, Y: 7}, {X: 11, Y: 8}} {
		if id, _ := rt.At(PlayerLayer, point.X, point.Y); id != EmptyRawID {
			t.Errorf("consumed strike %+v raw=%d, want empty", point, id)
		}
	}
}

func TestRuntimeGreatAnacondaDestroysEarlyStrikeWithoutDamage(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Anaconda.Phase = AnacondaPhaseEmerge
	rt.Anaconda.PhaseTicks = 20
	rt.SetForTest(PlayerLayer, 10, 7, 0)
	hits, _ := rt.tickGreatAnaconda(1)
	id, _ := rt.At(PlayerLayer, 10, 7)
	if hits != 0 || rt.Anaconda.Health != 3 || id != EmptyRawID {
		t.Fatalf("early strike hits=%d health=%d raw=%d, want 0/3/empty", hits, rt.Anaconda.Health, id)
	}
}

func TestRuntimeGreatAnacondaRespawnsMissingShaftBoulders(t *testing.T) {
	rt := newStage08Runtime(t)
	for _, x := range []int{12, 15} {
		rt.SetForTest(PlayerLayer, x, 5, EmptyRawID)
	}
	rt.Anaconda.Phase = AnacondaPhaseRespawn
	rt.Anaconda.PhaseTicks = 0
	rt.Anaconda.Animation = 4
	rt.Anaconda.AnimationTicks = 0
	rt.Anaconda.BodyY = 241
	rt.tickGreatAnaconda(1)
	if rt.Anaconda.Phase != AnacondaPhaseTailCharge || rt.Anaconda.Blocker != (Point{X: 10, Y: 4}) {
		t.Fatalf("phase=%d blocker=%+v, want tail charge and (10,4)", rt.Anaconda.Phase, rt.Anaconda.Blocker)
	}
	for _, x := range []int{12, 15} {
		if id, _ := rt.At(PlayerLayer, x, 2); id != 0 {
			t.Errorf("respawn shaft x=%d raw=%d, want boulder", x, id)
		}
	}
}

func TestRuntimeGreatAnacondaTailUsesFiftyAndTwelveTickWindows(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Anaconda.Phase = AnacondaPhaseTailCharge
	rt.Anaconda.PhaseTicks = 1
	rt.Anaconda.CycleTicks = 0
	rt.Anaconda.BodyY = 223
	rt.Player = Point{X: 3, Y: 5}
	for tick := 1; tick < 50; tick++ {
		rt.tickGreatAnaconda(tick)
	}
	if rt.Anaconda.Phase != AnacondaPhaseTailCharge || rt.Anaconda.CycleTicks != 49 {
		t.Fatalf("tail charge phase/ticks=%d/%d before 50, want charge/49", rt.Anaconda.Phase, rt.Anaconda.CycleTicks)
	}
	rt.tickGreatAnaconda(50)
	if rt.Anaconda.Phase != AnacondaPhaseTailStrike || !rt.Anaconda.TailVisible || rt.Anaconda.CycleTicks != 0 {
		t.Fatalf("tail strike phase=%d visible=%v ticks=%d", rt.Anaconda.Phase, rt.Anaconda.TailVisible, rt.Anaconda.CycleTicks)
	}
	for tick := 51; tick < 62; tick++ {
		rt.tickGreatAnaconda(tick)
	}
	if rt.Anaconda.Phase != AnacondaPhaseTailStrike || rt.Anaconda.CycleTicks != 11 {
		t.Fatalf("tail strike ended before tick 12: phase/ticks=%d/%d", rt.Anaconda.Phase, rt.Anaconda.CycleTicks)
	}
	rt.tickGreatAnaconda(62)
	if rt.Anaconda.Phase != AnacondaPhaseDescend || rt.Anaconda.TailVisible {
		t.Fatalf("tail strike phase=%d visible=%v after 12 ticks, want descend/false", rt.Anaconda.Phase, rt.Anaconda.TailVisible)
	}
}

func TestRuntimeGreatAnacondaThirdHitWaitsThenOpensGroupDoor(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.ActiveEnemyGateGroup = 0
	rt.EnemyGateCounters[0] = 1
	rt.Background[rt.index(18, 5)] = 0
	for hit := 0; hit < 3; hit++ {
		rt.Anaconda.Phase = AnacondaPhaseVulnerable
		rt.Anaconda.PhaseTicks = 0
		rt.SetForTest(PlayerLayer, rt.Anaconda.X(), 8, 0)
		got, _ := rt.tickGreatAnaconda(hit + 1)
		if got != 1 {
			t.Fatalf("boss hit %d result=%d", hit+1, got)
		}
	}
	for tick := 0; rt.Anaconda.Phase != AnacondaPhaseDefeated && tick < 50; tick++ {
		rt.tickGreatAnaconda(10 + tick)
	}
	if rt.Anaconda.Phase != AnacondaPhaseDefeated || !rt.Anaconda.Defeated {
		t.Fatalf("third-hit phase=%d defeated=%v", rt.Anaconda.Phase, rt.Anaconda.Defeated)
	}
	for tick := 1; tick <= 80; tick++ {
		rt.tickGreatAnaconda(100 + tick)
	}
	if rt.Anaconda.Phase != AnacondaPhaseDefeated || rt.EnemyGateCounters[0] != 1 {
		t.Fatalf("death delay phase=%d counter=%d at tick 80", rt.Anaconda.Phase, rt.EnemyGateCounters[0])
	}
	rt.tickGreatAnaconda(181)
	state, _ := rt.At(BackgroundLayer, 18, 5)
	if rt.Anaconda.Phase != AnacondaPhaseComplete || rt.EnemyGateCounters[0] != 0 || state != 0x10 {
		t.Fatalf("death completion phase=%d counter=%d door=%#x, want complete/0/phase1", rt.Anaconda.Phase, rt.EnemyGateCounters[0], state)
	}
}

func TestRuntimeStage08SealCompletesAfterSourceWait(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.ChestRewardID = 53
	rt.applyChestReward()
	if !rt.Anaconda.SealCollected || rt.RelicMask&1 == 0 || !rt.RelicCelebrating {
		t.Fatalf("seal collected=%v relic=%#x, want true/Angkor bit0", rt.Anaconda.SealCollected, rt.RelicMask)
	}
	if rt.CanAcceptInput() {
		t.Fatal("input is enabled during source hBoolean seal transition")
	}
	for tick := 1; tick <= angkorSealCompletionTicks; tick++ {
		rt.tickGreatAnaconda(tick)
	}
	if rt.Anaconda.StageComplete {
		t.Fatal("seal completed at tick 140, want source condition >140")
	}
	rt.tickGreatAnaconda(angkorSealCompletionTicks + 1)
	if !rt.Anaconda.StageComplete {
		t.Fatal("seal did not complete at tick 141")
	}
	if !rt.CanAcceptInput() {
		t.Fatal("input remained locked after seal transition completed")
	}
}

func TestRuntimeStage08CheckpointRestoresBossBeforeFight(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Anaconda.Health = 1
	rt.Anaconda.Phase = AnacondaPhaseTailStrike
	rt.Anaconda.TailVisible = true
	if !rt.RestoreCheckpoint() {
		t.Fatal("RestoreCheckpoint() = false")
	}
	if rt.Anaconda.Health != 3 || rt.Anaconda.Phase != AnacondaPhaseDormant || rt.Anaconda.TailVisible {
		t.Fatalf("restored boss=%+v, want initial source state", rt.Anaconda)
	}
}

func TestRuntimeRaw50IsPassableAndHitsBelowTwelvePixelOffset(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Anaconda.Enabled = false
	rt.Player = Point{X: 9, Y: 5}
	rt.PlayerMotion = ObjectMotion{}
	rt.SetForTest(PlayerLayer, 10, 5, 50)
	if !rt.TryMove(1, 0) {
		t.Fatal("raw50 should be passable")
	}
	rt.TickSourceFrame(8, 1, 0)
	if rt.HitCount != 0 {
		t.Fatalf("raw50 hit at offset %d, want no hit at 18", rt.PlayerMotion.Remaining)
	}
	rt.AdvancePlayerMotion()
	rt.TickSourceFrame(8, 2, 0)
	if rt.HitCount != 0 {
		t.Fatalf("raw50 hit at offset %d, want strict <12", rt.PlayerMotion.Remaining)
	}
	rt.AdvancePlayerMotion()
	rt.TickSourceFrame(8, 3, 0)
	if rt.HitCount != 1 {
		t.Fatalf("raw50 hit count=%d at offset %d, want 1", rt.HitCount, rt.PlayerMotion.Remaining)
	}
}

func TestRuntimeRaw50UsesTurnJIntThreshold(t *testing.T) {
	rt := newStage08Runtime(t)
	rt.Anaconda.Enabled = false
	rt.Player = Point{X: 10, Y: 5}
	rt.PlayerMotion = ObjectMotion{}
	rt.SetForTest(PlayerLayer, 10, 5, 50)
	for tick, offset := range []int{18, 12} {
		rt.SetPlayerTurnOffset(offset)
		rt.TickSourceFrame(8, tick+1, 0)
		if rt.HitCount != 0 {
			t.Fatalf("raw50 hit during turn offset %d, want strict <12", offset)
		}
	}
	rt.SetPlayerTurnOffset(6)
	rt.TickSourceFrame(8, 3, 0)
	if rt.HitCount != 1 {
		t.Fatalf("raw50 turn hit count=%d at offset 6, want 1", rt.HitCount)
	}
}

func TestRuntimeStage08CanDefeatGreatAnacondaAndCollectSealAtSourceCadence(t *testing.T) {
	rt := newStage08Runtime(t)
	route := &stage07Route{t: t, rt: rt}
	route.tick()
	route.walkTo("automatic entrance before door", Point{X: 3, Y: 5})
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 9 source entrance door")
	}
	route.walkTo("automatic entrance marker", Point{X: 4, Y: 5})
	route.walkTo("boss checkpoint", Point{X: 7, Y: 5})
	route.walkTo("cross open left arena door", Point{X: 8, Y: 5})
	route.walkTo("activate boss arena", Point{X: 9, Y: 5})
	route.waitUntil("arena introduction completes", 220, func() bool {
		return !rt.EnemyGateDemoActive
	})
	leftDoor, _ := rt.At(BackgroundLayer, 8, 5)
	rightDoor, _ := rt.At(BackgroundLayer, 18, 5)
	if leftDoor&0xf0 != 0 || rightDoor&0xf0 != 0 {
		t.Fatalf("arena doors left=%#x right=%#x, want both closed after introduction", leftDoor, rightDoor)
	}

	// The first body position spans x=10..11. Push the left shaft boulder
	// from x=12 into x=11 only after aoInt reaches vulnerable state 2.
	route.walkTo("enter arena and start boss", Point{X: 10, Y: 5})
	route.walkTo("first boulder push position", Point{X: 13, Y: 5})
	route.waitUntil("first vulnerable window", 100, func() bool {
		return rt.Anaconda.Phase == AnacondaPhaseVulnerable
	})
	route.push("drop first boulder into left body", -1)
	route.waitUntil("first boss hit", 120, func() bool {
		return rt.Anaconda.Health == 2
	})

	// Stand in source column 1 while the body descends. The left-pushed
	// right-shaft boulder then lands in body columns x=13..14.
	route.walkTo("select middle boss column", Point{X: 13, Y: 5})
	route.waitUntil("middle boss column selected", 180, func() bool {
		return rt.Anaconda.Column == 1 && rt.Anaconda.Phase == AnacondaPhaseDelay
	})
	route.walkTo("second boulder push position", Point{X: 16, Y: 5})
	route.waitUntil("second vulnerable window and respawn", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 15, 5)
		return rt.Anaconda.Phase == AnacondaPhaseVulnerable && id == 0
	})
	route.push("drop second boulder into middle body", -1)
	route.waitUntil("second boss hit", 120, func() bool {
		return rt.Anaconda.Health == 1
	})

	// x=16 selects source column 2. Push the regenerated right-shaft
	// boulder right from x=15 into body columns x=16..17.
	route.walkTo("select right boss column", Point{X: 16, Y: 5})
	route.waitUntil("right boss column selected", 180, func() bool {
		return rt.Anaconda.Column == 2 && rt.Anaconda.Phase == AnacondaPhaseDelay
	})
	route.walkTo("third boulder push position", Point{X: 14, Y: 5})
	route.waitUntil("third vulnerable window and respawn", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 15, 5)
		return rt.Anaconda.Phase == AnacondaPhaseVulnerable && id == 0
	})
	route.push("drop third boulder into right body", 1)
	route.waitUntil("third boss hit", 120, func() bool {
		return rt.Anaconda.Health == 0
	})
	route.waitUntil("boss death opens right arena door", 180, func() bool {
		return rt.Anaconda.Phase == AnacondaPhaseComplete && rt.IsPassable(18, 5)
	})

	route.walkTo("Angkor seal chest", Point{X: 27, Y: 6})
	route.waitUntil("Angkor seal reward", 120, func() bool {
		return rt.Anaconda.SealCollected && !rt.ChestOpening
	})
	route.waitUntil("Angkor seal stage transition", 180, func() bool {
		return rt.Anaconda.StageComplete
	})
	if rt.RelicMask&1 == 0 || rt.SpecialPickups != 1 || rt.Player != (Point{X: 27, Y: 6}) {
		t.Fatalf("seal route relic=%#x pickups=%d player=%+v", rt.RelicMask, rt.SpecialPickups, rt.Player)
	}
}
