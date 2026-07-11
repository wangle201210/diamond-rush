package originalgame

import (
	"image/png"
	"os"
	"slices"
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
	if g.lastDX != 1 || g.lastDY != 0 || g.heroTurnOffset != 0 {
		t.Fatalf("initial facing=%d,%d turn=%d, want source right/settled", g.lastDX, g.lastDY, g.heroTurnOffset)
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

func TestStage07UsesFreezeHammerCrossWorldPrerequisite(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	rt, err := original.NewRuntime(g.pack.Stages[7])
	if err != nil {
		t.Fatal(err)
	}
	g.applyCampaignProgress(rt, 7)
	if rt.SpecialItemMask != 8 {
		t.Fatalf("Stage 8 prerequisite tool=%d, want Freeze Hammer level 8", rt.SpecialItemMask)
	}
	if !rt.RestoreCheckpoint() || rt.SpecialItemMask != 8 {
		t.Fatalf("Stage 8 checkpoint restored tool=%d, want Freeze Hammer level 8", rt.SpecialItemMask)
	}
}

func TestSourceGoalFrameUsesSpecialAngkorStageVariants(t *testing.T) {
	for stageIndex := 0; stageIndex < 14; stageIndex++ {
		want := 0
		if stageIndex == 4 || stageIndex == 7 {
			want = 1
		}
		if got := sourceGoalFrame(stageIndex); got != want {
			t.Errorf("Stage %d goal frame=%d, want %d", stageIndex+1, got, want)
		}
	}
}

func TestNewLoadsFreezeHammerObjectSprites(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for name, sheet := range map[string]*spriteSheet{
		"frozen violet": g.frozenViolet,
		"frozen snake":  g.frozenSnake,
	} {
		if sheet == nil || sheet.moduleImage == nil || len(sheet.meta.FrameCounts) == 0 || sheet.meta.FrameCounts[0] == 0 {
			t.Errorf("%s source sprite is not drawable", name)
		}
	}
}

func TestStage08LoadsAndDrawsOriginalGreatAnacondaAssets(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for name, sheet := range map[string]*spriteSheet{
		"Great Anaconda": g.anaconda,
		"boss platforms": g.anacondaPlatform,
		"Angkor seal":    g.angkorSeal,
		"Bavaria seal":   g.bavariaSeal,
		"Siberia seal":   g.siberiaSeal,
		"seal arrow":     g.sealArrow,
		"softkeys":       g.softkeys,
	} {
		if sheet == nil || sheet.moduleImage == nil || len(sheet.meta.FrameCounts) == 0 {
			t.Errorf("%s source sprite is not drawable", name)
		}
	}
	for animation, want := range map[int]int{4: 12, 7: 22, 8: 2} {
		if duration, ok := g.anaconda.animationDuration(animation); !ok || duration != want {
			t.Errorf("boss animation %d duration=%d,%v, want %d,true", animation, duration, ok, want)
		}
	}
	if duration, ok := g.flames.animationDuration(2); !ok || duration != 11 {
		t.Fatalf("tail animation duration=%d,%v, want source 11,true", duration, ok)
	}

	for stage := 1; stage <= 8; stage++ {
		g.progress.unlockStage(stage)
	}
	g.loadStage(8)
	if g.stageIndex != 8 || !g.rt.Anaconda.Enabled {
		t.Fatalf("loaded stage=%d boss enabled=%v, want Stage 9 boss runtime", g.stageIndex, g.rt.Anaconda.Enabled)
	}
	g.rt.Anaconda.Phase = original.AnacondaPhaseVulnerable
	g.rt.Anaconda.BodyY = 216
	g.rt.Anaconda.Animation = 2
	dst := ebiten.NewImage(original.ScreenWidth, playfieldHeight)
	g.drawGreatAnaconda(dst, 100, 0)
	if len(g.anacondaPlatform.meta.FrameCounts) < 2 || g.anacondaPlatform.meta.FrameCounts[1] != 4 || len(g.angkorSeal.meta.Modules) != 1 {
		t.Fatalf("boss platform frames=%v seal modules=%v, want source frame-1 composition and seal module", g.anacondaPlatform.meta.FrameCounts, g.angkorSeal.meta.Modules)
	}
}

func TestStage08SealUsesDedicatedLoadingTransition(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for stage := 1; stage <= 8; stage++ {
		g.progress.unlockStage(stage)
	}
	g.loadStage(8)
	g.rt.RelicMask = 1
	g.rt.Anaconda.SealCollected = true
	g.rt.Anaconda.StageComplete = true
	if err := g.Update(); err != nil {
		t.Fatal(err)
	}
	if !g.sealExitActive || g.worldDone || !g.progress.StageCleared[8] {
		t.Fatalf("seal transition active=%v worldDone=%v stageCleared=%v", g.sealExitActive, g.worldDone, g.progress.StageCleared[8])
	}
	if g.progress.StageAwards[8] != 0 || g.progress.ExtraLives != 5 {
		t.Fatalf("seal stage awards=%#x lives=%d, want no ordinary result awards or bonus lives", g.progress.StageAwards[8], g.progress.ExtraLives)
	}
	screen := ebiten.NewImage(original.ScreenWidth, original.ScreenHeight)
	g.Draw(screen)
	for step := 0; step < sealLoadingSteps-1; step++ {
		if err := g.Update(); err != nil {
			t.Fatal(err)
		}
	}
	if g.mode != gameModeStage || !g.sealExitActive {
		t.Fatalf("seal transition ended before loading step 11: mode=%d active=%v ticks=%d", g.mode, g.sealExitActive, g.sealExitTicks)
	}
	if err := g.Update(); err != nil {
		t.Fatal(err)
	}
	if g.mode != gameModeWorldSelect || g.sealExitActive || g.worldSelectIncoming != sealPositionAngkor {
		t.Fatalf("seal completion mode=%d active=%v incoming=%d, want world selector/false/Angkor", g.mode, g.sealExitActive, g.worldSelectIncoming)
	}
	if g.progress.RelicMask != 1 {
		t.Fatalf("persisted relic mask=%#x, want Angkor bit", g.progress.RelicMask)
	}
}

func TestStage08SealCelebrationUsesHeroAnimation47(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.rt.RelicCelebrating = true
	g.rt.RelicCelebrationTicks = 17
	g.rt.ChestRewardID = 53
	animation, tick := g.heroAnimationState()
	if animation != 47 || tick != 17 {
		t.Fatalf("seal celebration animation=%d tick=%d, want 47/17", animation, tick)
	}
	if duration, ok := g.hero.animationDuration(47); !ok || duration != 42 {
		t.Fatalf("hero animation 47 duration=%d,%v, want source 42,true", duration, ok)
	}
}

func TestStage05LoadsAndDrawsOriginalFallingTorchAssets(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if g.fallingFire == nil || g.fallingTorches == nil || g.fallingDebris == nil {
		t.Fatal("mm0.f falling-torch resources were not loaded")
	}
	if duration, ok := g.fallingTorches.animationDuration(1); !ok || duration != 38 {
		t.Fatalf("torch collapse duration=%d,%v, want source 38,true", duration, ok)
	}
	if duration, ok := g.fallingFire.animationDuration(2); !ok || duration != 30 {
		t.Fatalf("fire startup duration=%d,%v, want source 30,true", duration, ok)
	}

	g.progress.HighestUnlocked = 5
	g.loadStage(5)
	g.rt.FallingTorchTriggers = 3
	g.rt.FallingTorchAnimation = 2
	g.rt.FallingTorchAnimationTicks = 1
	g.rt.RisingFireAnimation = 0
	g.rt.RisingFireAnimationTicks = 1
	g.rt.RisingFireHeight = 1000
	g.rt.FallingTorchWarningTicks = 60
	g.tick = 80
	dst := ebiten.NewImage(original.ScreenWidth, playfieldHeight)
	g.drawFallingTorchStage(dst, 180, 900)
	if frame, ok := g.fallingFire.animationFrame(g.rt.RisingFireAnimation, g.rt.RisingFireAnimationTicks); !ok || frame.Frame < 0 {
		t.Fatalf("rising-fire frame=%+v ok=%v, want drawable source frame", frame, ok)
	}
}

func TestSourceStageCellViewportMatchesJavaDynamicScan(t *testing.T) {
	view := sourceStageCellViewport(8*original.TileSize+5, 4*original.TileSize+7)
	if view.firstX != 8 || view.firstY != 4 || view.offX != -5 || view.offY != -7 {
		t.Fatalf("viewport origin=%+v, want first (8,4), offset (-5,-7)", view)
	}
	if view.firstRel != -1 || view.lastRelX != 12 || view.lastRelY != 12 {
		t.Fatalf("dynamic scan=%+v, want Java relative range [-1,12)", view)
	}
}

func TestLateForegroundOverlayKeepsEmptyRawIDTransparent(t *testing.T) {
	for _, tt := range []struct {
		id    original.RawID
		frame int
		ok    bool
	}{
		{id: original.EmptyRawID},
		{id: 79},
		{id: 80, frame: 0, ok: true},
		{id: 119, frame: 39, ok: true},
	} {
		frame, ok := sourceWorldOverlayFrame(tt.id)
		if frame != tt.frame || ok != tt.ok {
			t.Errorf("foreground raw %d overlay=%d,%v, want %d,%v", tt.id, frame, ok, tt.frame, tt.ok)
		}
	}
	for tick, want := range map[int]int{0: 0, 3: 0, 4: 1, 7: 1, 8: 2} {
		if got := sourceForegroundEffectSequence(tick); got != want {
			t.Errorf("foreground sequence at tick %d=%d, want %d", tick, got, want)
		}
	}
}

func TestStage05DemoCameraUsesSourcePanTarget(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.HighestUnlocked = 5
	g.loadStage(5)
	g.cameraX = 20
	g.cameraY = 300
	g.rt.ForegroundDemoActive = true
	g.rt.ForegroundDemoID = 3
	g.rt.ForegroundDemoPhase = 1
	g.rt.ForegroundDemoTicks = 30
	g.updateCamera()
	if g.cameraX != 100 || g.cameraY != 600 {
		t.Fatalf("half demo pan camera=%d,%d, want 100,600", g.cameraX, g.cameraY)
	}
	g.rt.ForegroundDemoPhase = 2
	g.rt.ForegroundDemoTicks = 0
	g.updateCamera()
	if g.cameraX != 180 || g.cameraY != 900 {
		t.Fatalf("locked demo camera=%d,%d, want source target 180,900", g.cameraX, g.cameraY)
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

func TestFirstNineStagesUnlockAndCarryOriginalToolPrerequisites(t *testing.T) {
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
	if progress.HighestUnlocked != 5 || !progress.StageCleared[4] || progress.ToolLevel != 2 {
		t.Fatalf("Stage 5 progress=%+v, want Stage 6 unlocked and hook persisted", progress)
	}
	g.progress = progress
	g.loadStage(5)
	if g.stageIndex != 5 || !g.rt.IsFallingTorchStage() || g.rt.SpecialItemMask != 2 {
		t.Fatalf("Stage 6 runtime index=%d falling=%v tool=%d, want 5/true/2", g.stageIndex, g.rt.IsFallingTorchStage(), g.rt.SpecialItemMask)
	}
	progress.recordStageResult(5, g.rt)
	if progress.HighestUnlocked != 6 || !progress.StageCleared[5] || !progress.StageUnlocked[6] {
		t.Fatalf("six-stage progress=%+v, want Stage 6 clear with Stage 7 unlocked", progress)
	}
	g.progress = progress
	g.loadStage(6)
	if g.stageIndex != 6 || g.rt.SpecialItemMask != 2 {
		t.Fatalf("Stage 7 runtime index=%d tool=%d, want 6/2", g.stageIndex, g.rt.SpecialItemMask)
	}
	progress.recordStageResult(6, g.rt)
	if progress.HighestUnlocked != 7 || !progress.StageCleared[6] || !progress.StageUnlocked[7] {
		t.Fatalf("seven-stage progress=%+v, want Stage 7 clear with Stage 8 unlocked", progress)
	}
	g.progress = progress
	g.loadStage(7)
	if g.stageIndex != 7 || g.rt.SpecialItemMask != 8 {
		t.Fatalf("Stage 8 runtime index=%d tool=%d, want 7/8", g.stageIndex, g.rt.SpecialItemMask)
	}
	progress.recordStageResult(7, g.rt)
	if progress.HighestUnlocked != 8 || !progress.StageCleared[7] || !progress.StageUnlocked[8] || progress.ToolLevel != 8 {
		t.Fatalf("eight-stage progress=%+v, want Stage 8 clear, Stage 9 unlocked, and Freeze Hammer persisted", progress)
	}
	g.progress = progress
	g.loadStage(8)
	if g.stageIndex != 8 || !g.rt.Anaconda.Enabled || g.rt.SpecialItemMask != 8 {
		t.Fatalf("Stage 9 runtime index=%d boss=%v tool=%d, want 8/true/8", g.stageIndex, g.rt.Anaconda.Enabled, g.rt.SpecialItemMask)
	}
	progress.recordStageResult(8, g.rt)
	if progress.HighestUnlocked != 8 || !progress.StageCleared[8] {
		t.Fatalf("nine-stage progress=%+v, want final normal Angkor node cleared without sequential secret unlock", progress)
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

func TestSecretStageRaw28UsesResultsAndUnlocksNextSecretNode(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	rt, err := original.NewRuntime(g.pack.Stages[9])
	if err != nil {
		t.Fatal(err)
	}
	g.stageIndex = 9
	g.rt = rt
	g.progress.unlockStage(9)
	rt.ReachedGoal = true
	rt.GoalExitSecret = true
	rt.GoalExitDirection = 2
	rt.Player = original.Point{X: rt.Width() + 5, Y: 11}
	rt.PlayerMotion = original.ObjectMotion{}
	g.advanceAfterGoal()
	if !g.worldDone || g.secretExitActive || g.resultPhase != resultPhaseLoading {
		t.Fatalf("secret-stage exit worldDone=%v secretMessage=%v resultPhase=%d", g.worldDone, g.secretExitActive, g.resultPhase)
	}
	if !g.progress.StageCleared[9] || !g.progress.StageUnlocked[10] || g.pendingMapTarget != 10 {
		t.Fatalf("secret-stage progress cleared=%v next=%v target=%d", g.progress.StageCleared[9], g.progress.StageUnlocked[10], g.pendingMapTarget)
	}
}

func TestStage06SecretExitSkipsResultsAndTravelsToSecretStageOne(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for stage := 1; stage <= 6; stage++ {
		g.progress.unlockStage(stage)
	}
	g.loadStage(6)
	g.introTicks = stageIntroDuration
	g.entranceSteps = 0
	g.rt.Player = original.Point{X: 20, Y: 40}
	g.rt.PlayerMotion = original.ObjectMotion{}
	if !g.rt.TryMove(1, 0) {
		t.Fatal("failed to enter Stage 7 source raw28 exit")
	}
	if !g.rt.GoalExitSecret {
		t.Fatal("raw28 exit was not tagged as secret")
	}
	for tick := 0; tick < 80 && !g.secretExitActive; tick++ {
		if g.rt.PlayerMotion.Remaining > 0 {
			g.advanceHeroMotion()
		} else {
			g.advanceAfterGoal()
		}
	}
	if !g.secretExitActive || g.worldDone || g.resultPhase != resultPhaseLoading {
		t.Fatalf("secret transition active=%v worldDone=%v resultPhase=%d, want true/false/loading", g.secretExitActive, g.worldDone, g.resultPhase)
	}
	if g.pendingMapTarget != 9 || !g.progress.StageUnlocked[9] || g.progress.StageUnlocked[7] || g.progress.StageUnlocked[8] {
		t.Fatalf("secret target=%d unlocks=%v, want target 9 only", g.pendingMapTarget, g.progress.StageUnlocked)
	}
	for g.secretExitActive {
		if err := g.Update(); err != nil {
			t.Fatal(err)
		}
	}
	if g.mode != gameModeWorldMap || g.worldMapTravelFrom != 6 || g.worldMapTravelTo != 9 {
		t.Fatalf("secret map mode=%d travel=%d->%d, want world map 6->9", g.mode, g.worldMapTravelFrom, g.worldMapTravelTo)
	}
}

func TestStage07SecretExitSkipsResultsAndTravelsToSecretStageFour(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for stage := 1; stage <= 7; stage++ {
		g.progress.unlockStage(stage)
	}
	g.loadStage(7)
	g.introTicks = stageIntroDuration
	g.entranceSteps = 0
	g.rt.Player = original.Point{X: 3, Y: 3}
	g.rt.PlayerMotion = original.ObjectMotion{}
	if !g.rt.TryMove(1, 0) || !g.rt.GoalExitSecret {
		t.Fatal("failed to enter Stage 8 source raw28 exit")
	}
	for tick := 0; tick < 100 && !g.secretExitActive; tick++ {
		if g.rt.PlayerMotion.Remaining > 0 {
			g.advanceHeroMotion()
		} else {
			g.advanceAfterGoal()
		}
	}
	if !g.secretExitActive || g.worldDone {
		t.Fatalf("Stage 8 secret transition active=%v worldDone=%v, want true/false", g.secretExitActive, g.worldDone)
	}
	if g.pendingMapTarget != 12 || !g.progress.StageUnlocked[12] || g.progress.StageUnlocked[8] {
		t.Fatalf("Stage 8 secret target=%d unlocks=%v, want target 12 without normal Stage 9", g.pendingMapTarget, g.progress.StageUnlocked)
	}
	for g.secretExitActive {
		if err := g.Update(); err != nil {
			t.Fatal(err)
		}
	}
	if g.mode != gameModeWorldMap || g.worldMapTravelFrom != 7 || g.worldMapTravelTo != 12 {
		t.Fatalf("Stage 8 secret map mode=%d travel=%d->%d, want world map 7->12", g.mode, g.worldMapTravelFrom, g.worldMapTravelTo)
	}
}

func TestSecretExitMessageFitsAndDrawsOnSourceScreen(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	lines := []string{"Congratulations! You have", "unlocked a secret path!"}
	for _, line := range lines {
		if width := g.fontSmall.stringWidth(line) + 10; width > original.ScreenWidth {
			t.Fatalf("secret message line %q panel width=%d, exceeds %d", line, width, original.ScreenWidth)
		}
	}
	g.secretExitActive = true
	screen := ebiten.NewImage(original.ScreenWidth, original.ScreenHeight)
	g.Draw(screen)
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
	if g.hero.moduleImage == nil || g.door.moduleImage == nil || g.snakes.moduleImage == nil || g.crawler.moduleImage == nil || g.flames.moduleImage == nil || g.pressureSwitch.moduleImage == nil || g.quota.moduleImage == nil || g.hiddenOverlay.moduleImage == nil || g.specialContainer.moduleImage == nil || g.goldKey.moduleImage == nil || g.silverKey.moduleImage == nil || g.tools.moduleImage == nil || g.toolPrompt.moduleImage == nil || g.pickupEffects.moduleImage == nil || g.resultSpark.moduleImage == nil || g.resultMedal.moduleImage == nil {
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
	if got := []int{
		resultStageTitleY, resultCompleteY,
		resultVioletLabelY, resultVioletCountY,
		resultRedLabelY, resultRedCountY,
		resultHitsIconY, resultHitsLabelY, resultHitsCountY,
		resultRetriesIconY, resultRetriesLabelY, resultRetriesCountY,
	}; !slices.Equal(got, []int{15, 32, 75, 91, 131, 147, 191, 187, 203, 243, 243, 259}) {
		t.Fatalf("result layout coordinates = %v, want Java case 17 coordinates", got)
	}
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

func TestMacActionAndRecallShortcutsAreUnambiguous(t *testing.T) {
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
	if !recallPressedWith(keySet(ebiten.KeyEnter), keySet()) {
		t.Fatal("Enter did not trigger checkpoint recall")
	}
	if !recallPressedWith(keySet(ebiten.KeyBackspace), keySet()) {
		t.Fatal("Backspace did not trigger checkpoint recall")
	}
	if !centerActionPressedWith(keySet(ebiten.KeySpace)) {
		t.Fatal("Space did not trigger the phone 5 interaction")
	}
	if centerActionPressedWith(keySet(ebiten.KeyEnter)) {
		t.Fatal("Enter still triggered the phone 5 interaction")
	}
	if !centerActionPressedWith(keySet(ebiten.KeyDigit5)) || !centerActionPressedWith(keySet(ebiten.KeyNumpad5)) {
		t.Fatal("phone 5 keyboard aliases did not trigger interaction")
	}
	if !tutorialSkipPressedWith(keySet(ebiten.KeyS)) || tutorialSkipPressedWith(keySet(ebiten.KeySpace)) {
		t.Fatal("tutorial skip must use S without consuming Space")
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
	payloads := map[original.RawID]bool{2: true, 4: true, 5: true, 6: true, 7: true, 24: true, 26: true, 27: true, 40: true, 41: true, 42: true, 51: true, 52: true, 53: true}
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
	if hidden != 120 {
		t.Fatalf("all-Angkor hidden payloads=%d, want 120 audited container cells", hidden)
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

func TestCompassChestRewardUsesOriginalGen3Chunk1Artwork(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if g.compassPickup == nil || g.compassPickup.moduleImage == nil || len(g.compassPickup.meta.Modules) != 1 {
		t.Fatal("gen3.f chunk 1 compass pickup is not loaded")
	}
	if got := g.compassPickup.meta.Modules[0]; got != (spriteModuleMeta{W: 24, H: 24}) {
		t.Fatalf("compass module bounds = %+v, want source 24x24", got)
	}
	f, err := os.Open(resolvePath(compassPickupModules))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	visiblePixels := 0
	for y := framePadding; y < framePadding+original.TileSize; y++ {
		for x := framePadding; x < framePadding+original.TileSize; x++ {
			r, green, b, alpha := img.At(x, y).RGBA()
			isTransparentKey := uint8(r>>8) == 20 && uint8(green>>8) == 22 && uint8(b>>8) == 28
			if alpha != 0 && !isTransparentKey {
				visiblePixels++
			}
		}
	}
	if visiblePixels == 0 {
		t.Fatal("gen3.f chunk 1 compass module contains no visible pixels")
	}

	dst := ebiten.NewImage(original.TileSize, original.TileSize)
	g.rt.ChestRewardID = 42
	if !g.drawChestRewardIcon(dst, 0, 0) {
		t.Fatal("raw 42 chest reward did not select the original compass artwork")
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

func TestDirectionChangeTurnsBeforeMovingAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStageForGame(t)
	rt, err := original.NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	g := &Game{rt: rt, lastDY: -1}
	start := rt.Player
	if g.handlePlayerDirection(1, 0) {
		t.Fatal("first changed-direction input moved instead of turning")
	}
	if rt.Player != start || g.lastDX != 1 || g.lastDY != 0 || g.heroTurnOffset != sourceHeroTurnStartOffset {
		t.Fatalf("turn start player=%+v facing=%d,%d offset=%d", rt.Player, g.lastDX, g.lastDY, g.heroTurnOffset)
	}
	for step, want := range []int{12, 6, 0} {
		if !g.advanceHeroTurn() {
			t.Fatalf("turn cadence stopped at step %d", step)
		}
		if rt.Player != start || g.heroTurnOffset != want {
			t.Fatalf("turn step %d player=%+v offset=%d, want %+v/%d", step, rt.Player, g.heroTurnOffset, start, want)
		}
	}
	if g.advanceHeroTurn() {
		t.Fatal("completed turn consumed an extra source frame")
	}
	if !g.handlePlayerDirection(1, 0) || rt.Player == start {
		t.Fatalf("held direction did not move after turn: player=%+v start=%+v", rt.Player, start)
	}
}

func TestSameFacingDirectionMovesWithoutTurn(t *testing.T) {
	stage := mustLoadOriginalStageForGame(t)
	rt, err := original.NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	g := &Game{rt: rt, lastDX: 1}
	if !g.handlePlayerDirection(1, 0) {
		t.Fatal("same-facing input did not move immediately")
	}
	if g.heroTurnOffset != 0 || rt.Player != (original.Point{X: 1, Y: 17}) {
		t.Fatalf("same-facing move player=%+v turn=%d", rt.Player, g.heroTurnOffset)
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
