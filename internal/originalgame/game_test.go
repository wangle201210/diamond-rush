package originalgame

import (
	"image/png"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

func TestLayoutUsesOriginalPhoneScreen(t *testing.T) {
	g := &Game{}
	gotW, gotH := g.Layout(0, 0)
	if gotW != original.ScreenWidth || gotH != original.ScreenHeight {
		t.Fatalf("Layout() = %dx%d, want %dx%d", gotW, gotH, original.ScreenWidth, original.ScreenHeight)
	}
}

func TestNewLoadsAngkorWorldPack(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.pack.Stages) != 14 {
		t.Fatalf("stages = %d, want 14", len(g.pack.Stages))
	}
	if g.stageIndex != 0 {
		t.Fatalf("stage index = %d, want 0", g.stageIndex)
	}
	if g.rt.Stage.Index != 0 {
		t.Fatalf("runtime stage = %d, want 0", g.rt.Stage.Index)
	}
	if g.rt.Health != 4 || g.rt.MaxHealth != 4 {
		t.Fatalf("initial health = %d/%d, want source default 4/4", g.rt.Health, g.rt.MaxHealth)
	}
	if g.rt.ExtraLives != 5 {
		t.Fatalf("initial extra lives = %d, want source default 5", g.rt.ExtraLives)
	}
	if g.introTicks != 0 {
		t.Fatalf("intro ticks = %d, want first-stage intro pending", g.introTicks)
	}
	if g.entranceSteps != 4 {
		t.Fatalf("entrance steps = %d, want source raw79 auto-walk distance 4", g.entranceSteps)
	}
	if !g.rt.CompassEnabled || g.compassDirection != 4 {
		t.Fatalf("Stage 1 compass enabled=%v direction=%d, want true/east(4)", g.rt.CompassEnabled, g.compassDirection)
	}
	if g.cameraY != 248 {
		t.Fatalf("initial camera y = %d, want source entrance camera 248", g.cameraY)
	}
}

func TestLoadStageSwitchesRuntime(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.HighestUnlocked = 4
	g.loadStage(4)
	if g.stageIndex != 4 {
		t.Fatalf("stage index = %d, want 4", g.stageIndex)
	}
	if g.rt.Stage.Index != 4 {
		t.Fatalf("runtime stage = %d, want 4", g.rt.Stage.Index)
	}
	if got := g.rt.EntranceMarker; got != (original.Point{X: 3, Y: 19}) {
		t.Fatalf("stage04 entrance = %+v, want (3,19)", got)
	}
	if g.rt.SpecialItemMask != 2 {
		t.Fatalf("stage04 standalone prerequisite tool=%d, want Mystic Hook level 2", g.rt.SpecialItemMask)
	}
}

func TestLoadStageCarriesSourceGlobalState(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.HighestUnlocked = 4
	g.progress.ExtraLives = 8
	g.progress.MaxHealth = 6
	g.progress.ToolLevel = 1
	g.loadStage(3)
	if g.rt.ExtraLives != 8 || g.rt.Health != 6 || g.rt.MaxHealth != 6 || g.rt.SpecialItemMask != 1 {
		t.Fatalf("carried runtime lives=%d health=%d/%d tool=%d, want 8, 6/6, 1", g.rt.ExtraLives, g.rt.Health, g.rt.MaxHealth, g.rt.SpecialItemMask)
	}
	g.rt.Health = 2
	if !g.rt.ResetCheckpoint() {
		t.Fatal("failed to reset initial checkpoint snapshot")
	}
	if g.rt.Health != 2 || g.rt.MaxHealth != 6 || g.rt.ExtraLives != 8 || g.rt.SpecialItemMask != 1 {
		t.Fatalf("restored campaign snapshot health=%d/%d lives=%d tool=%d, want 2/6/8/1", g.rt.Health, g.rt.MaxHealth, g.rt.ExtraLives, g.rt.SpecialItemMask)
	}
}

func TestFirstFiveStagesUnlockAndCarryOriginalToolPrerequisites(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	progress := newOriginalProgress()
	for stageIndex := 0; stageIndex < 4; stageIndex++ {
		rt, err := original.NewRuntime(g.pack.Stages[stageIndex])
		if err != nil {
			t.Fatal(err)
		}
		if stageIndex == 3 {
			rt.SpecialItemMask = 1
		}
		progress.recordStageResult(stageIndex, rt)
		if progress.HighestUnlocked != stageIndex+1 {
			t.Fatalf("after Stage %d highest unlocked=%d, want %d", stageIndex+1, progress.HighestUnlocked, stageIndex+1)
		}
	}
	if progress.ToolLevel != 1 {
		t.Fatalf("tool after Stage 4=%d, want source hammer level 1", progress.ToolLevel)
	}
	g.progress = progress
	g.loadStage(4)
	if g.rt.SpecialItemMask != 2 {
		t.Fatalf("Stage 5 runtime tool=%d, want original cross-world hook prerequisite level 2", g.rt.SpecialItemMask)
	}
	progress.recordStageResult(4, g.rt)
	if progress.HighestUnlocked != 4 || !progress.StageCleared[4] || progress.ToolLevel != 2 {
		t.Fatalf("final five-stage progress=%+v, want Stage 5 clear, no Stage 6 unlock, hook persisted", progress)
	}
}

func TestAdvanceAfterGoalAutoWalksBeyondStageBeforeResults(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.rt.Player = original.Point{X: 21, Y: 9}
	g.rt.PlayerMotion = original.ObjectMotion{}
	if !g.rt.TryMove(1, 0) {
		t.Fatal("failed to enter source raw5 goal")
	}
	g.syncHeroMotion()
	for tick := 0; tick < 64 && !g.worldDone; tick++ {
		if g.rt.PlayerMotion.Remaining > 0 {
			g.advanceHeroMotion()
		} else {
			g.advanceAfterGoal()
		}
	}
	if g.stageIndex != 0 {
		t.Fatalf("stage index = %d, want first stage to remain loaded", g.stageIndex)
	}
	if g.rt.Stage.Index != 0 || !g.rt.ReachedGoal {
		t.Fatalf("runtime stage=%d reached=%v, want cleared first-stage runtime", g.rt.Stage.Index, g.rt.ReachedGoal)
	}
	if !g.worldDone {
		t.Fatal("world done = false after source exit auto-walk")
	}
	if g.rt.Player.X != g.rt.Width()+6 || !g.rt.GoalExitComplete {
		t.Fatalf("exit player=%+v complete=%v, want x=%d and complete", g.rt.Player, g.rt.GoalExitComplete, g.rt.Width()+6)
	}
	if g.resultPhase != resultPhaseLoading || g.resultLoadingStep != 0 {
		t.Fatalf("result phase=%d loading step=%d, want source loading phase at step 0", g.resultPhase, g.resultLoadingStep)
	}
}

func TestAdvanceAfterGoalDoesNotCompleteAtGoalCell(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.rt.Player = original.Point{X: 21, Y: 9}
	g.rt.PlayerMotion = original.ObjectMotion{}
	if !g.rt.TryMove(1, 0) {
		t.Fatal("failed to enter source raw5 goal")
	}
	g.syncHeroMotion()
	for g.rt.PlayerMotion.Remaining > 0 {
		g.advanceHeroMotion()
	}
	g.advanceAfterGoal()
	if g.worldDone {
		t.Fatal("world done = true before hero left the stage boundary")
	}
	if g.rt.Player != (original.Point{X: 23, Y: 9}) || g.heroMoveOffset != 12 {
		t.Fatalf("first exit step player=%+v offset=%d, want (23,9)/12", g.rt.Player, g.heroMoveOffset)
	}
}

func TestHeroUsesExtractedAnimationMetadata(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	frame, ok := g.hero.animationFrame(5, 0)
	if !ok {
		t.Fatal("source hero walking animation 5 is unavailable")
	}
	if frame.Frame != 0 {
		t.Fatalf("hero walking animation first frame = %d, want source frame 0", frame.Frame)
	}
	if got := heroDirection(1, 0); got != 1 {
		t.Fatalf("right-facing hero animation = %d, want 1", got)
	}
}

func TestHeroUsesSourceHammerAnimation(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.HighestUnlocked = 3
	g.loadStage(3)
	g.rt.SpecialItemMask = 1
	g.rt.Player = original.Point{X: 33, Y: 5}
	if !g.rt.UseHammer(1, 0) {
		t.Fatal("failed to start right-facing hammer action")
	}
	if animation, tick := g.heroAnimationState(); animation != 14 || tick != 0 {
		t.Fatalf("hammer hero animation=%d tick=%d, want 14/0", animation, tick)
	}
	g.rt.TickSourceFrame(8, 1, 0)
	if animation, tick := g.heroAnimationState(); animation != 14 || tick != 1 {
		t.Fatalf("advanced hammer hero animation=%d tick=%d, want 14/1", animation, tick)
	}
}

func TestHeroAndRopeUseSourceHookAnimation(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.HighestUnlocked = 4
	g.loadStage(4)
	g.rt.SpecialItemMask = 2
	g.rt.Player = original.Point{X: 4, Y: 10}
	for x := 5; x <= 7; x++ {
		g.rt.SetForTest(original.PlayerLayer, x, 10, original.EmptyRawID)
		g.rt.SetForTest(original.ForegroundLayer, x, 10, original.EmptyRawID)
		g.rt.SetForTest(original.PlayerLayer, x, 11, 80)
	}
	g.rt.SetForTest(original.PlayerLayer, 7, 10, 0)
	g.rt.ObjectState[7+10*g.rt.Width()] = 0
	if !g.rt.UseHook(1, 0) {
		t.Fatal("failed to start right-facing hook action")
	}
	if animation, tick := g.heroAnimationState(); animation != 20 || tick != 0 {
		t.Fatalf("hook cast hero animation=%d tick=%d, want 20/0", animation, tick)
	}
	g.rt.HookTicks = 100
	if animation, tick := g.heroAnimationState(); animation != 20 || tick != 6 {
		t.Fatalf("held hook cast animation=%d tick=%d, want final cast frame tick 6", animation, tick)
	}
	g.rt.HookTicks = 0

	ropeIdx := 5 + 10*g.rt.Width()
	startX, tipX, module := hookSegmentGeometry(g.rt.ObjectState[ropeIdx], g.rt.ObjectMotion[ropeIdx].Remaining, 0)
	if startX != 0 || tipX != 6 || module != 0 {
		t.Fatalf("right hook geometry start=%d tip=%d module=%d, want 0/6/0", startX, tipX, module)
	}
	if hookRopeColor.R != 0xd3 || hookRopeColor.G != 0xd7 || hookRopeColor.B != 0xe7 || hookRopeColor.A != 0xff {
		t.Fatalf("hook rope color=%v, want source #d3d7e7", hookRopeColor)
	}

	for sourceTick := 1; sourceTick <= 16 && g.rt.Hooking && g.rt.HookAnimation != 21; sourceTick++ {
		g.rt.TickSourceFrame(8, sourceTick, 0)
	}
	if animation, _ := g.heroAnimationState(); animation != 21 {
		t.Fatalf("hook pull hero animation=%d, want source animation 21", animation)
	}
}

func TestSnakeAnimationUsesSourceSequenceCadence(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	first, ok := g.snakes.animationFrameAtSequence(0, 0)
	if !ok || first.Frame != 11 {
		t.Fatalf("snake sequence 0 frame=%d ok=%v, want source frame 11", first.Frame, ok)
	}
	second, ok := g.snakes.animationFrameAtSequence(0, 1)
	if !ok || second.Frame != 12 {
		t.Fatalf("snake sequence 1 frame=%d ok=%v, want source frame 12", second.Frame, ok)
	}
	if timed, ok := g.snakes.animationFrame(0, 1); !ok || timed.Frame != 11 {
		t.Fatalf("generic timed snake frame=%d ok=%v, fixture must demonstrate why source sequence indexing is required", timed.Frame, ok)
	}
}

func TestCrawlerRenderStateUsesSourceModulesAndHiddenPhase(t *testing.T) {
	tests := []struct {
		name                    string
		state, tick             int
		frame, offsetX, offsetY int
		visible                 bool
	}{
		{name: "normal frame", state: 0, tick: 4, frame: 2, offsetX: 2, offsetY: 2, visible: true},
		{name: "death phase", state: 2 << 8, frame: 4, offsetX: 2, offsetY: 2, visible: true},
		{name: "clockwise offset", state: 1, frame: 0, offsetX: 6, offsetY: 2, visible: true},
		{name: "reversed offset", state: 1 | 0x10, frame: 0, offsetX: -2, offsetY: 2, visible: true},
		{name: "hidden", state: 4 << 8, visible: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame, offsetX, offsetY, visible := crawlerRenderState(tt.state, tt.tick)
			if frame != tt.frame || offsetX != tt.offsetX || offsetY != tt.offsetY || visible != tt.visible {
				t.Fatalf("crawler state=%#x tick=%d gives %d,(%d,%d),%v; want %d,(%d,%d),%v", tt.state, tt.tick, frame, offsetX, offsetY, visible, tt.frame, tt.offsetX, tt.offsetY, tt.visible)
			}
		})
	}
}

func TestPressureSwitchOffsetMatchesSourceDepression(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	x, y := 10, 10
	idx := x + y*g.rt.Width()
	g.rt.SetForTest(original.PlayerLayer, x, y, 8)
	g.rt.ObjectMotion[idx] = original.ObjectMotion{Remaining: 6}
	if got := g.pressureSwitchOffset(x, y, 8); got != 6 {
		t.Fatalf("raw8 pressure offset=%d, want 6", got)
	}

	g.rt.SetForTest(original.PlayerLayer, x, y, original.EmptyRawID)
	g.rt.Player = original.Point{X: x, Y: y}
	g.rt.PlayerMotion = original.ObjectMotion{Remaining: 3}
	if got := g.pressureSwitchOffset(x, y, original.EmptyRawID); got != 9 {
		t.Fatalf("entering-player pressure offset=%d, want 9", got)
	}

	g.rt.Player = original.Point{X: x - 1, Y: y}
	g.rt.PlayerMotion = original.ObjectMotion{DX: -1, Remaining: 18}
	if got := g.pressureSwitchOffset(x, y, original.EmptyRawID); got != 6 {
		t.Fatalf("leaving-player pressure offset=%d, want 6", got)
	}
}

func TestGemAnimationUsesSourceBurstThenRestCadence(t *testing.T) {
	wants := map[int]int{
		0:  0,
		2:  1,
		4:  2,
		6:  3,
		8:  0,
		63: 0,
		64: 0,
		66: 1,
	}
	for tick, want := range wants {
		if got := sourceGemFrame(tick); got != want {
			t.Errorf("sourceGemFrame(%d)=%d, want %d", tick, got, want)
		}
	}
}

func TestVioletPickupQueuesSourceAnimation3(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.rt.Player = original.Point{X: 6, Y: 2}
	g.rt.PlayerMotion = original.ObjectMotion{}
	if !g.rt.TryMove(1, 0) {
		t.Fatal("failed to enter stage00 violet gem")
	}
	g.tick = 1
	g.tickWorld()
	if len(g.worldEffects) != 1 {
		t.Fatalf("world effects=%d, want one violet pickup", len(g.worldEffects))
	}
	effect := g.worldEffects[0]
	if effect.Point != (original.Point{X: 7, Y: 2}) || effect.Animation != 3 || effect.Sequence != 0 {
		t.Fatalf("violet effect=%+v, want point (7,2), animation3, sequence0", effect)
	}
	g.advanceWorldEffects()
	if g.worldEffects[0].Sequence != 1 {
		t.Fatalf("advanced violet effect sequence=%d, want 1", g.worldEffects[0].Sequence)
	}
}

func TestHUDUsesExtractedDigitAndHealthModules(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if g.hud.moduleImage == nil {
		t.Fatal("HUD module image is nil")
	}
	if got := g.hud.moduleWidth(0); got != 7 {
		t.Fatalf("HUD zero width = %d, want 7", got)
	}
	if got := g.hud.moduleWidth(15); got != 5 {
		t.Fatalf("HUD health segment width = %d, want 5", got)
	}
	if got := g.hud.meta.FrameCounts[1]; got != 2 {
		t.Fatalf("HUD gem-vessel frame modules = %d, want purple and red modules", got)
	}
	if g.hero.moduleImage == nil || g.door.moduleImage == nil || g.snakes.moduleImage == nil || g.crawler.moduleImage == nil || g.flames.moduleImage == nil || g.pressureSwitch.moduleImage == nil || g.quota.moduleImage == nil || g.hiddenOverlay.moduleImage == nil || g.specialContainer.moduleImage == nil || g.goldKey.moduleImage == nil || g.silverKey.moduleImage == nil || g.tools.moduleImage == nil || g.pickupEffects.moduleImage == nil || g.resultSpark.moduleImage == nil || g.resultMedal.moduleImage == nil {
		t.Fatal("source module sheets are required for gameplay and result-screen composite sprites")
	}
	if len(g.quota.meta.FrameCounts) < 1 || g.quota.meta.FrameCounts[0] != 2 {
		t.Fatalf("quota marker frame modules=%v, want source two-module composition", g.quota.meta.FrameCounts)
	}
	if len(g.crawler.meta.Modules) != 6 || g.crawler.meta.Modules[0].W != original.TileSize {
		t.Fatalf("crawler modules=%v, want six source modules with a 24px base frame", g.crawler.meta.Modules)
	}
	if len(g.pressureSwitch.meta.Modules) != 1 || g.pressureSwitch.meta.Modules[0].H != 13 {
		t.Fatalf("pressure switch modules=%v, want one 24x13 source module", g.pressureSwitch.meta.Modules)
	}
	if g.fontSmall == nil || g.fontMedium == nil || g.fontSmall.meta.FontHeight != 10 || g.fontMedium.meta.FontHeight != 12 {
		t.Fatal("FreeJ2ME 10px and 12px source font atlases are required")
	}
	first := g.hud.meta.FrameFirst[1]
	modules := g.hud.meta.FrameModules[first : first+g.hud.meta.FrameCounts[1]]
	if modules[0].Module != 24 || modules[1].Module != 25 {
		t.Fatalf("HUD vessel modules = %d,%d, want source red/purple modules 24,25", modules[0].Module, modules[1].Module)
	}
}

func TestStageResultUsesSourceLoadingAndPhaseTiming(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.rt.VioletGems = 21
	g.beginStageResults()
	for step := 1; step < resultLoadingSteps; step++ {
		g.updateStageResults(false)
		if g.resultPhase != resultPhaseLoading || g.resultLoadingStep != step {
			t.Fatalf("loading step %d phase=%d step=%d, want loading/%d", step, g.resultPhase, g.resultLoadingStep, step)
		}
	}
	g.updateStageResults(false)
	if g.resultPhase != resultPhaseTitle || g.resultPhaseTicks != 1 {
		t.Fatalf("after loading phase=%d tick=%d, want title/1", g.resultPhase, g.resultPhaseTicks)
	}

	advanceResultPhaseToLastTick(t, g, resultPhaseTitle, resultTitleTicks+1)
	g.updateStageResults(false)
	if g.resultPhase != resultPhaseVioletGems || g.resultPhaseTicks != 1 {
		t.Fatalf("after title phase=%d tick=%d, want violet/1", g.resultPhase, g.resultPhaseTicks)
	}
	advanceResultPhaseToLastTick(t, g, resultPhaseVioletGems, 43)
	g.updateStageResults(false)
	if g.resultPhase != resultPhaseRedDiamonds || g.resultPhaseTicks != 1 {
		t.Fatalf("after violet count phase=%d tick=%d, want red/1", g.resultPhase, g.resultPhaseTicks)
	}
	for _, wantTicks := range []int{resultRedDiamondTicks, resultHitTicks, resultRetryTicks} {
		phase := g.resultPhase
		advanceResultPhaseToLastTick(t, g, phase, wantTicks+1)
		g.updateStageResults(false)
		if g.resultPhase != phase+1 || g.resultPhaseTicks != 1 {
			t.Fatalf("after phase %d result phase=%d tick=%d, want %d/1", phase, g.resultPhase, g.resultPhaseTicks, phase+1)
		}
	}
	if g.resultPhase != resultPhaseComplete {
		t.Fatalf("final result phase=%d, want complete", g.resultPhase)
	}
}

func TestStageResultSourceSlideAndCountFrames(t *testing.T) {
	tests := []struct {
		tick                 int
		title, complete, row int
	}{
		{1, -90, -330, -90},
		{10, 0, -240, 0},
		{34, 0, 0, 0},
	}
	for _, tt := range tests {
		title, complete := stageResultTitleOffsets(resultPhaseTitle, tt.tick)
		row := stageResultRowOffset(resultPhaseVioletGems, resultPhaseVioletGems, tt.tick)
		if title != tt.title || complete != tt.complete || row != tt.row {
			t.Errorf("tick %d offsets=%d,%d,%d, want %d,%d,%d", tt.tick, title, complete, row, tt.title, tt.complete, tt.row)
		}
	}
	if got := min(1>>1, 21); got != 0 {
		t.Fatalf("violet count at phase tick 1 = %d, want 0", got)
	}
	if got := min(2>>1, 21); got != 1 {
		t.Fatalf("violet count at phase tick 2 = %d, want 1", got)
	}
	if got := stageResultPhaseDuration(resultPhaseVioletGems, 21); got != 42 {
		t.Fatalf("21-gem count phase duration = %d, want 42", got)
	}
}

func TestStageResultAwardsMatchSourceConditions(t *testing.T) {
	rt := &original.Runtime{TotalVioletGems: 21, VioletGems: 21, TotalRedDiamonds: 1, RedDiamonds: 1}
	want := byte(resultAwardVioletGems | resultAwardRedDiamonds | resultAwardNoHits | resultAwardNoRetries)
	if got := stageResultAwards(rt); got != want {
		t.Fatalf("perfect result awards = %#x, want %#x", got, want)
	}
	rt.VioletGems--
	rt.HitCount = 1
	want &^= resultAwardVioletGems | resultAwardNoHits
	if got := stageResultAwards(rt); got != want {
		t.Fatalf("imperfect result awards = %#x, want %#x", got, want)
	}
}

func TestStageResultEffectsUseOriginalJARFlatSequenceShifts(t *testing.T) {
	tests := []struct {
		mode       int
		tick       int
		sequence   int
		shouldDraw bool
	}{
		{resultEffectDoubleShort, 1, 2, true},
		{resultEffectDoubleShort, 4, 8, true},
		{resultEffectDoubleShort, 5, 10, false},
		{resultEffectHalf, 1, 0, true},
		{resultEffectHalf, 19, 9, true},
		{resultEffectHalf, 20, 10, false},
		{resultEffectDoubleLong, 1, 2, true},
		{resultEffectDoubleLong, 19, 38, true},
		{resultEffectDoubleLong, 20, 40, false},
	}
	for _, tt := range tests {
		sequence, ok := stageResultEffectSequence(tt.mode, tt.tick, 10)
		if sequence != tt.sequence || ok != tt.shouldDraw {
			t.Errorf("mode=%d tick=%d sequence=%d draw=%v, want %d/%v", tt.mode, tt.tick, sequence, ok, tt.sequence, tt.shouldDraw)
		}
	}

	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	frame, ok := g.pickupEffects.animationFrameAtRawSequence(0, 38)
	if !ok || frame.Frame != 38 {
		t.Fatalf("flat result effect frame=%+v ok=%v, want source frame 38 across animation boundaries", frame, ok)
	}
}

func advanceResultPhaseToLastTick(t *testing.T, g *Game, phase, lastTick int) {
	t.Helper()
	if g.resultPhase != phase || g.resultPhaseTicks != 1 {
		t.Fatalf("phase start=%d/%d, want %d/1", g.resultPhase, g.resultPhaseTicks, phase)
	}
	for g.resultPhaseTicks < lastTick {
		g.updateStageResults(false)
	}
	if g.resultPhase != phase || g.resultPhaseTicks != lastTick {
		t.Fatalf("phase end=%d/%d, want %d/%d", g.resultPhase, g.resultPhaseTicks, phase, lastTick)
	}
}

func TestStageIntroUsesSourceSlidePositionsAndText(t *testing.T) {
	font, err := loadBitmapFont(fontMediumSheet, fontMediumMetadata)
	if err != nil {
		t.Fatal(err)
	}
	if !font.supports("ANGKOR WATSTAGE 1CONGRATULATIONS!CHECKPOINT") {
		t.Fatal("FreeJ2ME medium atlas is missing a Stage 1 presentation glyph")
	}
	tests := []struct {
		tick           int
		worldX, stageX int
	}{
		{0, 0, 240},
		{11, 120, 120},
		{40, 120, 120},
		{50, 0, 240},
	}
	for _, tt := range tests {
		worldX, stageX := stageTitlePositions(tt.tick)
		if worldX != tt.worldX || stageX != tt.stageX {
			t.Errorf("tick %d title positions = %d,%d, want %d,%d", tt.tick, worldX, stageX, tt.worldX, tt.stageX)
		}
	}
}

func TestMacCheckpointRecallShortcutsDoNotMoveDown(t *testing.T) {
	pressed := keySet(ebiten.KeyShiftLeft, ebiten.KeyDigit8)
	justPressed := keySet(ebiten.KeyDigit8)
	if !recallPressedWith(justPressed, pressed) {
		t.Fatal("Shift+8 (*) did not trigger checkpoint recall")
	}
	if dx, dy := heldDirectionWith(pressed); dx != 0 || dy != 0 {
		t.Fatalf("Shift+8 movement = %d,%d, want no movement", dx, dy)
	}
	if !recallPressedWith(keySet(ebiten.KeyR), keySet()) {
		t.Fatal("R did not trigger checkpoint recall")
	}
	if !recallPressedWith(keySet(ebiten.KeyBackspace), keySet()) {
		t.Fatal("Backspace did not trigger checkpoint recall")
	}
	if dx, dy := heldDirectionWith(keySet(ebiten.KeyDigit8)); dx != 0 || dy != 1 {
		t.Fatalf("plain 8 movement = %d,%d, want down", dx, dy)
	}
}

func TestHiddenChestUsesClosedFrameBeforeOpening(t *testing.T) {
	if got := hiddenOverlayFrame(0); got != 0 {
		t.Fatalf("closed chest frame = %d, want 0", got)
	}
	if got := hiddenOverlayFrame(int(original.EmptyRawID)); got != 3 {
		t.Fatalf("invalid high chest state frame = %d, want clamped final frame 3", got)
	}
}

func TestClosedContainersDoNotRevealTheirPayloads(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	payloads := map[original.RawID]bool{2: true, 4: true, 5: true, 6: true, 7: true, 24: true, 26: true, 27: true, 41: true}
	hidden := 0
	for stageIndex := 0; stageIndex < angkorReplicaStageCount; stageIndex++ {
		stage := g.pack.Stages[stageIndex]
		for idx, foregroundID := range stage.Foreground {
			playerID := stage.Player[idx]
			if (foregroundID != 14 && foregroundID != 33) || !payloads[playerID] {
				continue
			}
			hidden++
			if sourceCellObjectVisible(playerID, foregroundID) {
				t.Errorf("stage %d container (%d,%d) foreground raw%d reveals payload raw%d", stageIndex+1, idx%stage.Width, idx/stage.Width, foregroundID, playerID)
			}
		}
	}
	if hidden != 28 {
		t.Fatalf("first-five hidden payloads=%d, want 28 audited container cells", hidden)
	}
	if !sourceCellObjectVisible(2, original.EmptyRawID) || !sourceCellObjectVisible(1, 33) {
		t.Fatal("container visibility filter hid a loose reward or non-payload object")
	}
}

func TestChestRewardTickMatchesSourceHeroAnimation40(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if got := animationDuration(g.hero, 40); got != chestOpenDurationForTest {
		t.Fatalf("hero animation 40 duration = %d, want %d", got, chestOpenDurationForTest)
	}
	atSound, ok := g.hero.animationSequenceIndex(40, chestRewardSoundTick)
	if !ok || atSound != 12 {
		t.Fatalf("animation 40 reward sound tick = %d,%v, want frame index 12", atSound, ok)
	}
	before, ok := g.hero.animationSequenceIndex(40, chestRewardTickForTest-1)
	if !ok || before != 12 {
		t.Fatalf("animation 40 before reward = %d,%v, want frame index 12", before, ok)
	}
	atReward, ok := g.hero.animationSequenceIndex(40, chestRewardTickForTest)
	if !ok || atReward != 13 {
		t.Fatalf("animation 40 reward tick = %d,%v, want frame index 13", atReward, ok)
	}
	if got := g.hiddenOverlay.meta.FrameCounts[0]; got != 2 {
		t.Fatalf("closed chest frame modules = %d, want source lid/body pair", got)
	}
	if got, ok := g.pickupEffects.animationDuration(0); !ok || got != 10 {
		t.Fatalf("chest pickup effect duration = %d,%v, want 10,true", got, ok)
	}
	if chestRewardIconVisible(g.hero, 40, chestRewardTickForTest) {
		t.Fatal("red-diamond overhead icon is visible on source sequence frame 13")
	}
	if !chestRewardIconVisible(g.hero, 40, chestRewardTickForTest+2) {
		t.Fatal("red-diamond overhead icon is hidden after source sequence frame 13")
	}
	for _, tt := range []struct {
		chestTick int
		wantTick  int
		wantOK    bool
	}{
		{chestTick: chestRewardTickForTest - 1, wantTick: -1, wantOK: false},
		{chestTick: chestRewardTickForTest, wantTick: 0, wantOK: true},
		{chestTick: chestRewardTickForTest + 9, wantTick: 9, wantOK: true},
		{chestTick: chestRewardTickForTest + 10, wantTick: 10, wantOK: false},
	} {
		gotTick, gotOK := chestRewardEffectTick(g.pickupEffects, 0, 40, tt.chestTick)
		if gotTick != tt.wantTick || gotOK != tt.wantOK {
			t.Errorf("chest tick %d effect = %d,%v, want %d,%v", tt.chestTick, gotTick, gotOK, tt.wantTick, tt.wantOK)
		}
	}
	g.rt.ChestOpening = true
	g.rt.ChestAnimation = 40
	g.rt.ChestTicks = 20
	if animation, tick := g.heroAnimationState(); animation != 40 || tick != 20 {
		t.Fatalf("opening chest hero animation = %d,%d, want 40,20", animation, tick)
	}
}

func TestChestAnimation48UsesSourceRewardFrame(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if got := animationDuration(g.hero, 48); got != 46 {
		t.Fatalf("hero animation 48 duration=%d, want 46", got)
	}
	sequence, ok := g.hero.animationSequenceIndex(48, chestShortRewardTick)
	if !ok || sequence != chestShortRewardSequence {
		t.Fatalf("animation 48 reward tick sequence=%d,%v, want %d,true", sequence, ok, chestShortRewardSequence)
	}
	if chestRewardIconVisible(g.hero, 48, chestShortRewardTick) {
		t.Fatal("short-pickup icon is visible on reward sequence 6")
	}
	if !chestRewardIconVisible(g.hero, 48, chestShortRewardTick+2) {
		t.Fatal("short-pickup icon is hidden after reward sequence 6")
	}
	for rewardID, wantAnimation := range map[original.RawID]int{2: 0, 6: 0, 5: 1, 4: 2, 41: 3, 7: 4} {
		got, ok := chestRewardEffectAnimation(rewardID)
		if !ok || got != wantAnimation {
			t.Errorf("reward raw%d effect=%d,%v, want %d,true", rewardID, got, ok, wantAnimation)
		}
	}
	if _, ok := chestRewardEffectAnimation(24); ok {
		t.Fatal("special tool pickup unexpectedly uses a common reward effect")
	}
	g.rt.ChestOpening = true
	g.rt.ChestAnimation = 48
	g.rt.ChestTicks = 20
	if animation, tick := g.heroAnimationState(); animation != 48 || tick != 20 {
		t.Fatalf("short chest hero animation=%d tick=%d, want 48/20", animation, tick)
	}
}

func TestHeroUsesSourceRecallAnimation19(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.rt.RecallPending = true
	g.rt.RecallTicks = 17
	animation, tick := g.heroAnimationState()
	if animation != 19 || tick != 17 {
		t.Fatalf("recall hero animation=%d tick=%d, want 19/17", animation, tick)
	}
	if got := animationDuration(g.hero, 19); got != 42 {
		t.Fatalf("hero animation 19 duration = %d, want source 42", got)
	}
}

func TestDiggableFrameUsesExtractedAngkorPixels(t *testing.T) {
	f, err := os.Open(resolvePath(angkorDiggableFrameSheet))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	greenPixels := 0
	for y := framePadding; y < framePadding+original.TileSize; y++ {
		for x := framePadding; x < framePadding+original.TileSize; x++ {
			r, green, b, a := img.At(x, y).RGBA()
			if a > 0 && green > r && green > b {
				greenPixels++
			}
		}
	}
	if greenPixels < 40 {
		t.Fatalf("diggable frame green pixels = %d, want extracted vegetation art", greenPixels)
	}
}

func keySet(keys ...ebiten.Key) func(ebiten.Key) bool {
	set := make(map[ebiten.Key]bool, len(keys))
	for _, key := range keys {
		set[key] = true
	}
	return func(key ebiten.Key) bool {
		return set[key]
	}
}

const (
	chestOpenDurationForTest = 67
	chestRewardSoundTick     = 37
	chestRewardTickForTest   = 39
)

func animationDuration(sheet *spriteSheet, animation int) int {
	first := sheet.meta.AnimationFirst[animation]
	count := sheet.meta.AnimationCounts[animation]
	total := 0
	for _, frame := range sheet.meta.AnimationFrames[first : first+count] {
		total += max(1, frame.Time)
	}
	return total
}

func TestCameraPixelsClampToStage(t *testing.T) {
	stage := mustLoadOriginalStageForGame(t)
	rt, err := original.NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	g := &Game{rt: rt}
	g.resetCamera()
	x, y := g.cameraPixels()
	if x != 0 {
		t.Fatalf("initial camera x = %d, want 0", x)
	}
	if y <= 0 {
		t.Fatalf("initial camera y = %d, want positive vertical follow", y)
	}
	rt.Player = original.Point{X: stage.Width - 1, Y: stage.Height - 1}
	for i := 0; i < 16; i++ {
		g.updateCamera()
	}
	x, y = g.cameraPixels()
	wantX := stage.Width*original.TileSize - original.ScreenWidth
	wantY := stage.Height*original.TileSize - playfieldHeight
	if x != wantX || y != wantY {
		t.Fatalf("bottom-right camera = %d,%d, want %d,%d", x, y, wantX, wantY)
	}
}

func TestHeroMovementUsesSourceSubTileOffsets(t *testing.T) {
	stage := mustLoadOriginalStageForGame(t)
	rt, err := original.NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	g := &Game{rt: rt, lastDX: 1}
	if !g.startPlayerMove(1, 0) {
		t.Fatal("source corridor first move failed")
	}
	if rt.Player != (original.Point{X: 1, Y: 17}) {
		t.Fatalf("logical player = %+v, want target cell (1,17)", rt.Player)
	}
	if x, y := g.renderedPlayerPixels(); x != 6 || y != 17*original.TileSize {
		t.Fatalf("first rendered position = %d,%d, want 6,%d", x, y, 17*original.TileSize)
	}
	for step, wantX := range []int{12, 18, 24} {
		g.advanceHeroMotion()
		if x, _ := g.renderedPlayerPixels(); x != wantX {
			t.Fatalf("interpolation step %d x = %d, want %d", step+1, x, wantX)
		}
	}
	if g.heroMoveOffset != 0 {
		t.Fatalf("final hero offset = %d, want 0", g.heroMoveOffset)
	}
}

func TestUpdateAutoWalksRaw79EntranceAndCommitsCheckpoint(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for sourceTick := 0; sourceTick < 17; sourceTick++ {
		if err := g.Update(); err != nil {
			t.Fatal(err)
		}
	}
	if g.rt.Player != (original.Point{X: 4, Y: 17}) {
		t.Fatalf("auto-walk player = %+v, want entrance checkpoint (4,17)", g.rt.Player)
	}
	if g.entranceSteps != 0 || g.heroMoveOffset != 0 {
		t.Fatalf("entrance steps=%d offset=%d, want settled zeroes", g.entranceSteps, g.heroMoveOffset)
	}
	if state, _ := g.rt.At(original.BackgroundLayer, 2, 17); state != 0x0f || g.rt.IsPassable(2, 17) {
		t.Fatalf("entrance door state=%#x passable=%v, want closed source door", state, g.rt.IsPassable(2, 17))
	}
	if g.rt.CheckpointPending {
		t.Fatal("entrance checkpoint remained pending after movement settled")
	}
	if !g.rt.RestoreCheckpoint() || g.rt.Player != (original.Point{X: 4, Y: 17}) {
		t.Fatalf("restored entrance checkpoint player = %+v", g.rt.Player)
	}
	if g.compassDirection != 1 {
		t.Fatalf("compass after entrance = %d, want source north-east sector 1 toward checkpoint (11,9)", g.compassDirection)
	}
}

func TestCompassCacheTracksSourceDirectionWhileStageIdles(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for sourceTick := 1; sourceTick <= 320; sourceTick++ {
		if err := g.Update(); err != nil {
			t.Fatal(err)
		}
		if sourceTick&0xf != 0 {
			continue
		}
		want, ok := g.rt.CompassDirection()
		if !ok || g.compassDirection != want {
			t.Fatalf("tick %d cached compass=%d, runtime=%d,%v player=%+v progress=%d", sourceTick, g.compassDirection, want, ok, g.rt.Player, g.rt.CheckpointProgress)
		}
	}
}

func mustLoadOriginalStageForGame(t *testing.T) *original.Stage {
	t.Helper()
	stage, err := original.LoadStageFile("../../decoded/world0/stage00.json")
	if err != nil {
		t.Fatal(err)
	}
	return stage
}
