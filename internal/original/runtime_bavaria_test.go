package original

import (
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

func TestBavariaExplosiveBoulderUsesSourceFallAndBlastCadence(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 0)
	for y := 6; y <= 10; y++ {
		for x := 4; x <= 6; x++ {
			rt.SetForTest(PlayerLayer, x, y, EmptyRawID)
			rt.SetForTest(ForegroundLayer, x, y, EmptyRawID)
		}
	}
	rt.SetForTest(PlayerLayer, 5, 7, 8)
	rt.SetForTest(PlayerLayer, 4, 10, 80)
	rt.SetForTest(PlayerLayer, 5, 10, 80)
	rt.SetForTest(PlayerLayer, 6, 10, 80)
	rt.SetForTest(PlayerLayer, 6, 9, 37)
	rt.SetForTest(PlayerLayer, 4, 9, 8)
	rt.Player = Point{X: 6, Y: 8}
	rt.Health = 4

	for sourceTick := 1; sourceTick <= 20 && rt.PlayerLayer[rt.index(5, 9)] != 54; sourceTick++ {
		rt.TickSourceFrame(20, sourceTick, 4)
	}
	if got := rt.PlayerLayer[rt.index(5, 9)]; got != 54 {
		t.Fatalf("fallen explosive raw=%d, want active explosion raw54", got)
	}
	for sourceTick := 21; sourceTick <= 26; sourceTick++ {
		rt.TickSourceFrame(20, sourceTick, 4)
	}
	if got := rt.PlayerLayer[rt.index(4, 9)]; got != 54 {
		t.Fatalf("chain explosive raw=%d, want raw54 at impact tick", got)
	}
	if got := rt.ObjectState[rt.index(6, 9)]; got != 1 {
		t.Fatalf("blast wall state=%d, want destruction state1", got)
	}
	if rt.Health != 3 {
		t.Fatalf("blast health=%d, want 3", rt.Health)
	}
	for tick := 0; tick < 8; tick++ {
		rt.TickBreakables()
	}
	if got := rt.PlayerLayer[rt.index(6, 9)]; got != EmptyRawID {
		t.Fatalf("blast wall raw=%d after 8 ticks, want empty", got)
	}
}

func TestBavariaHeroCanPushSourceExplosiveBoulder(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 5)
	rt.Player = Point{X: 3, Y: 6}
	rt.SetForTest(PlayerLayer, 4, 6, 8)
	rt.SetForTest(PlayerLayer, 5, 6, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 4, 6, EmptyRawID)
	rt.SetForTest(ForegroundLayer, 5, 6, EmptyRawID)
	for attempt := 1; attempt < boulderPushAttempts; attempt++ {
		if rt.TryMove(1, 0) {
			t.Fatalf("pushed explosive early on attempt %d", attempt)
		}
	}
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to push source raw8 explosive after push delay")
	}
	if rt.Player != (Point{X: 4, Y: 6}) {
		t.Fatalf("player=%+v, want (4,6) after explosive push", rt.Player)
	}
	if id, _ := rt.At(PlayerLayer, 5, 6); id != 8 {
		t.Fatalf("pushed explosive raw=%d, want 8 at (5,6)", id)
	}
}

func TestBavariaFallingBoulderDetonatesSourceExplosive(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 0)
	for y := 4; y <= 9; y++ {
		rt.SetForTest(PlayerLayer, 5, y, EmptyRawID)
		rt.SetForTest(ForegroundLayer, 5, y, EmptyRawID)
	}
	rt.SetForTest(PlayerLayer, 5, 5, 0)
	rt.SetForTest(PlayerLayer, 5, 8, 8)
	rt.SetForTest(PlayerLayer, 5, 9, 80)
	rt.Player = Point{X: 10, Y: 10}

	for sourceTick := 1; sourceTick <= 40 && rt.PlayerLayer[rt.index(5, 8)] != 54; sourceTick++ {
		rt.TickSourceFrame(20, sourceTick, 4)
	}
	if got := rt.PlayerLayer[rt.index(5, 8)]; got != 54 {
		t.Fatalf("boulder-triggered explosive raw=%d, want active explosion raw54", got)
	}
	if got := rt.PlayerLayer[rt.index(5, 7)]; got != 0 {
		t.Fatalf("impact boulder raw=%d, want raw0 at (5,7)", got)
	}
}

func TestBavariaSpikeColumnsUseSourceSlowAndFastCycles(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 0)
	for y := 6; y <= 10; y++ {
		rt.SetForTest(PlayerLayer, 5, y, EmptyRawID)
		rt.SetForTest(ForegroundLayer, 5, y, EmptyRawID)
	}
	rt.SetForTest(PlayerLayer, 5, 7, 28)
	rt.ObjectState[rt.index(5, 7)] = 3
	rt.Player = Point{X: 5, Y: 9}
	rt.Health = 4
	rt.TickSourceFrame(20, 45, 4)
	if rt.SpikeSlowExtent != 48 || rt.Health != 2 {
		t.Fatalf("slow spike extent/health=%d/%d, want 48/2", rt.SpikeSlowExtent, rt.Health)
	}

	rt.Player = Point{X: 10, Y: 10}
	rt.HurtTicks = 0
	rt.InvulnerabilityTicks = 0
	rt.TickSourceFrame(20, 30, 4)
	if rt.SpikeSlowExtent != 24 || rt.IsPassable(5, 8) {
		t.Fatalf("half spike extent/passable=%d/%v, want 24/false", rt.SpikeSlowExtent, rt.IsPassable(5, 8))
	}

	rt.ObjectState[rt.index(5, 7)] = 3 | 0x8
	rt.TickSourceFrame(20, 22, 4)
	if rt.SpikeFastExtent != 48 {
		t.Fatalf("fast spike extent=%d, want 48", rt.SpikeFastExtent)
	}
}

func TestBavariaFanSwitchSwapsAuthoredPodGroupsAtPhaseFive(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 4)
	closedAtZero := Point{X: 19, Y: 4}
	openAtZero := Point{X: 19, Y: 7}
	if rt.PlayerLayer[rt.index(openAtZero.X, openAtZero.Y)] != EmptyRawID || rt.Foreground[rt.index(openAtZero.X, openAtZero.Y)] != 15 {
		t.Fatal("raw34 was not initialized as the phase-zero open pod")
	}
	if rt.PlayerLayer[rt.index(closedAtZero.X, closedAtZero.Y)] != 35 || rt.Foreground[rt.index(closedAtZero.X, closedAtZero.Y)] != EmptyRawID {
		t.Fatal("raw35 was not initialized as the phase-zero closed pod")
	}
	if rt.PlayerLayer[rt.index(7, 9)] != 16 || rt.PlayerLayer[rt.index(7, 10)] != 16 {
		t.Fatal("raw16 did not expand into its source two-cell spear pair")
	}

	rt.SpecialItemMask = 1
	rt.SetForTest(PlayerLayer, 5, 24, EmptyRawID)
	rt.Player = Point{X: 5, Y: 24}
	rt.SetForTest(ForegroundLayer, 5, 24, 15)
	if !rt.UseHammer(1, 0) {
		t.Fatal("hammer did not target raw18 while hero stood on a fan pod")
	}
	for rt.Hammering {
		rt.TickStatus()
	}
	if rt.FanDirection != 0 {
		t.Fatalf("fan direction while standing on pod=%d, want source no-op", rt.FanDirection)
	}
	rt.SetForTest(ForegroundLayer, 5, 24, EmptyRawID)
	if !rt.UseHammer(1, 0) {
		t.Fatal("hammer did not target raw18 fan switch")
	}
	for tick := 0; tick < hammerImpactTick; tick++ {
		rt.TickStatus()
	}
	if rt.FanDirection != 1 {
		t.Fatalf("fan direction=%d, want opening direction +1", rt.FanDirection)
	}
	for _, sourceTick := range []int{0, 1, 4, 5, 8} {
		rt.tickFanPhase(sourceTick)
	}
	if rt.FanPhase != 5 || rt.PlayerLayer[rt.index(openAtZero.X, openAtZero.Y)] != 34 || rt.Foreground[rt.index(closedAtZero.X, closedAtZero.Y)] != 16 {
		t.Fatalf("phase-five pods phase=%d openRaw=%d closedForeground=%d", rt.FanPhase, rt.PlayerLayer[rt.index(openAtZero.X, openAtZero.Y)], rt.Foreground[rt.index(closedAtZero.X, closedAtZero.Y)])
	}
	for _, sourceTick := range []int{9, 12, 13, 16} {
		rt.tickFanPhase(sourceTick)
	}
	if rt.FanPhase != 9 || rt.FanDirection != 0 {
		t.Fatalf("fan end phase/direction=%d/%d, want 9/0", rt.FanPhase, rt.FanDirection)
	}
}

func TestBavariaStageFiveSharedSilverDoorStartsLockedWithoutKeys(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 4)
	door := Point{X: 15, Y: 17}
	trigger := Point{X: 16, Y: 17}

	if got := rt.Foreground[rt.index(trigger.X, trigger.Y)]; got != 26 {
		t.Fatalf("adjacent trigger raw=%d, want authored raw26", got)
	}
	state, _ := rt.At(BackgroundLayer, door.X, door.Y)
	if state != 2 || rt.IsPassable(door.X, door.Y) {
		t.Fatalf("shared silver door state/passable=%#x/%v, want count2/false", state, rt.IsPassable(door.X, door.Y))
	}
	if rt.KeyForForeground8 != 0 {
		t.Fatalf("initial silver keys=%d, want 0", rt.KeyForForeground8)
	}

	rt.Player = Point{X: door.X - 1, Y: door.Y}
	rt.PlayerMotion = ObjectMotion{}
	if rt.TryMove(1, 0) {
		t.Fatal("hero entered the shared silver door without a key")
	}
	if rt.Player != (Point{X: door.X - 1, Y: door.Y}) {
		t.Fatalf("blocked hero moved to %+v", rt.Player)
	}
}

func TestBavariaMovingHazardDestroysSourceSpearPair(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 0)
	for y := 5; y <= 7; y++ {
		for x := 4; x <= 6; x++ {
			rt.SetForTest(PlayerLayer, x, y, EmptyRawID)
			rt.SetForTest(ForegroundLayer, x, y, EmptyRawID)
		}
	}
	top, base := rt.index(5, 5), rt.index(5, 6)
	rt.SetForTest(PlayerLayer, 5, 5, 16)
	rt.SetForTest(PlayerLayer, 5, 6, 16)
	rt.ObjectState[top] = 2
	rt.ObjectState[base] = 2
	rt.SetForTest(PlayerLayer, 4, 5, 14)
	rt.ObjectState[rt.index(4, 5)] = 2
	rt.DrainSoundEvents()
	rt.tickSpearPairAt(5, 6)
	if rt.PlayerLayer[top] != EmptyRawID || rt.PlayerLayer[base] != EmptyRawID {
		t.Fatalf("moving-hazard spear pair raw=%d/%d, want both removed", rt.PlayerLayer[top], rt.PlayerLayer[base])
	}
	if sounds := rt.DrainSoundEvents(); !slices.Contains(sounds, SoundBoulder) {
		t.Fatalf("moving-hazard spear sounds=%v, want source boulder sound %d", sounds, SoundBoulder)
	}
}

func TestBavariaFallenMovingHazardDestroysSourceSpearPair(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 4)
	for y := 4; y <= 7; y++ {
		for x := 4; x <= 6; x++ {
			rt.SetForTest(PlayerLayer, x, y, EmptyRawID)
			rt.SetForTest(ForegroundLayer, x, y, EmptyRawID)
		}
	}
	top, base := rt.index(5, 6), rt.index(5, 7)
	rt.SetForTest(PlayerLayer, 5, 6, 16)
	rt.SetForTest(PlayerLayer, 5, 7, 16)
	rt.ObjectState[top] = 2
	rt.ObjectState[base] = 2
	rt.SetForTest(PlayerLayer, 5, 5, 14)
	rt.ObjectMotion[rt.index(5, 5)] = ObjectMotion{DY: 1, Remaining: 6}
	rt.tickSpearPairAt(5, 7)
	if rt.PlayerLayer[top] != EmptyRawID || rt.PlayerLayer[base] != EmptyRawID {
		t.Fatalf("fallen moving-hazard spear pair raw=%d/%d, want both removed", rt.PlayerLayer[top], rt.PlayerLayer[base])
	}
}

func TestBavariaMovingHazardCrushesSnakeInTravelDirection(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 3)
	for x := 20; x <= 22; x++ {
		rt.SetForTest(PlayerLayer, x, 10, EmptyRawID)
		rt.SetForTest(ForegroundLayer, x, 10, EmptyRawID)
	}
	rt.SetForTest(PlayerLayer, 20, 10, 14)
	rt.ObjectState[rt.index(20, 10)] = 2
	rt.SetForTest(PlayerLayer, 21, 10, 43)
	group := 7
	rt.EnemyGateGroup[rt.index(21, 10)] = group
	rt.EnemyGateCounters[group] = 1
	rt.ActiveEnemyGateGroup = group

	rt.tickSnakeObjectAt(21, 10)
	if id, _ := rt.At(PlayerLayer, 21, 10); id != EmptyRawID || rt.EnemyGateCounters[group] != 0 {
		t.Fatalf("moving-hazard snake raw/count=%d/%d, want empty/0", id, rt.EnemyGateCounters[group])
	}
}

func TestBavariaFallingBoulderDoesNotDeleteMovingHazard(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 2)
	for y := 23; y <= 27; y++ {
		for x := 9; x <= 11; x++ {
			rt.SetForTest(PlayerLayer, x, y, EmptyRawID)
			rt.SetForTest(ForegroundLayer, x, y, EmptyRawID)
		}
	}
	rt.SetForTest(PlayerLayer, 10, 25, 0)
	rt.ObjectState[rt.index(10, 25)] = 3
	rt.ObjectMotion[rt.index(10, 25)] = ObjectMotion{DY: 1, Remaining: 6}
	rt.SetForTest(PlayerLayer, 10, 26, 14)
	rt.ObjectState[rt.index(10, 26)] = 2
	rt.ObjectMotion[rt.index(10, 26)] = ObjectMotion{DX: 1, Remaining: 12}

	rt.tickMovingHazardAt(10, 26)
	if id, _ := rt.At(PlayerLayer, 10, 26); id != 14 {
		t.Fatalf("moving hazard below falling boulder raw=%d, want source raw14 retained", id)
	}
	rt.tickGravityObjectAt(10, 25)
	if boulder, _ := rt.At(PlayerLayer, 10, 25); boulder != 0 {
		t.Fatalf("boulder above moving hazard raw=%d, want boulder retained", boulder)
	}
}

func TestBavariaCrawlerTrapCountsEnemyGroupAndHurtsAboveCell(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 0)
	for y := 6; y <= 9; y++ {
		rt.SetForTest(PlayerLayer, 5, y, EmptyRawID)
		rt.SetForTest(ForegroundLayer, 5, y, EmptyRawID)
	}
	trap := rt.index(5, 8)
	rt.SetForTest(PlayerLayer, 5, 8, 36)
	rt.SetForTest(PlayerLayer, 5, 7, 11)
	rt.EnemyGateGroup[trap] = 2
	rt.EnemyGateCounters[2] = 1
	rt.ActiveEnemyGateGroup = 2
	rt.tickCrawlerTrapAt(5, 8)
	if rt.ObjectState[trap] != 1 || rt.EnemyGateCounters[2] != 0 {
		t.Fatalf("trap state/count=%d/%d, want 1/0", rt.ObjectState[trap], rt.EnemyGateCounters[2])
	}
	rt.SetForTest(PlayerLayer, 5, 7, EmptyRawID)
	rt.Player = Point{X: 5, Y: 7}
	rt.Health = 4
	rt.tickCrawlerTrapAt(5, 8)
	if rt.Health != 3 {
		t.Fatalf("active crawler trap health=%d, want 3", rt.Health)
	}
}

func TestBavariaPurpleChestAddsAuthoredValueToVioletCount(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 0)
	chest := Point{X: 25, Y: 14}
	if rt.PlayerLayer[rt.index(chest.X, chest.Y)] != 41 || rt.Foreground[rt.index(chest.X, chest.Y)] != 33 {
		t.Fatal("Bavaria Stage 1 authored purple chest moved")
	}
	rt.Player = chest
	rt.startChestOpening(chest, false)
	for tick := 0; tick < chestRewardTick; tick++ {
		rt.TickStatus()
	}
	if rt.VioletGems != 10 || rt.BonusValue != 10 || !rt.ChestRewarded {
		t.Fatalf("purple chest violet/bonus/rewarded=%d/%d/%v, want 10/10/true", rt.VioletGems, rt.BonusValue, rt.ChestRewarded)
	}
}

func TestBavariaStageThreeAwardsMysticHookFromSourceChest(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 2)
	chest := Point{X: 24, Y: 25}
	if rt.PlayerLayer[rt.index(chest.X, chest.Y)] != 27 || rt.Foreground[rt.index(chest.X, chest.Y)] != 14 {
		t.Fatal("Bavaria Stage 3 Mystic Hook source chest moved")
	}
	rt.Player = chest
	rt.startChestOpening(chest, false)
	for tick := 0; tick < chestRewardTick; tick++ {
		rt.TickStatus()
	}
	if rt.SpecialItemMask&2 == 0 || rt.SpecialPickups != 1 || !rt.ChestRewarded {
		t.Fatalf("hook mask/pickups/rewarded=%d/%d/%v, want hook/1/true", rt.SpecialItemMask, rt.SpecialPickups, rt.ChestRewarded)
	}
}

func TestBavariaWaterPotionRewardRunsSourceExplanation(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 7)
	chest := Point{X: 20, Y: 12}
	idx := rt.index(chest.X, chest.Y)
	if rt.PlayerLayer[idx] != 40 || rt.Foreground[idx] != 14 {
		t.Fatal("Bavaria Stage 8 Mystic Potion source chest moved")
	}
	rt.Player = chest
	rt.startChestOpening(chest, false)
	for tick := 0; tick < chestRewardTick; tick++ {
		rt.TickStatus()
	}
	if rt.SpecialItemMask&4 == 0 || !rt.TutorialScriptActive || rt.TutorialScriptID != bavariaWaterPotionScriptID {
		t.Fatalf("potion mask/script=%d/%v/%d, want water breathing and script 24", rt.SpecialItemMask, rt.TutorialScriptActive, rt.TutorialScriptID)
	}

	prompts := make([]int, 0, 2)
	for tick := 0; tick < 20 && rt.TutorialScriptActive; tick++ {
		rt.tickTutorial()
		if prompt, ok := rt.TutorialPrompt(); ok && rt.AdvanceTutorialPrompt() {
			prompts = append(prompts, prompt.TextIndex)
		}
	}
	if rt.TutorialScriptActive || !slices.Equal(prompts, []int{26, 27}) {
		t.Fatalf("potion explanation active/prompts=%v/%v, want completed [26 27]", rt.TutorialScriptActive, prompts)
	}
}

func TestBavariaSecretStageEntryHintRequiresMissingPotion(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 10)
	if !rt.HasStageEntryHint() || !rt.StartStageEntryHint() {
		t.Fatal("Bavaria Secret Stage 1 did not expose its missing-potion entry hint")
	}
	rt.tickTutorial()
	prompt, ok := rt.TutorialPrompt()
	if !ok || prompt.TextIndex != bavariaSecretPotionHintTextIndex || prompt.Placement != TutorialTextBottom {
		t.Fatalf("entry hint prompt=%+v,%v, want bottom text %d", prompt, ok, bavariaSecretPotionHintTextIndex)
	}

	rt = mustLoadBavariaRuntime(t, 10)
	rt.SpecialItemMask |= 4
	if rt.HasStageEntryHint() || rt.StartStageEntryHint() {
		t.Fatal("Bavaria Secret Stage 1 still showed the prerequisite after collecting the potion")
	}
}

func TestBavariaMysticHookExistsOnlyInDisplayedStageThree(t *testing.T) {
	pack, err := LoadWorldDir(filepath.Join("..", "..", "decoded", "world1"))
	if err != nil {
		t.Fatal(err)
	}
	var locations []struct {
		stage int
		point Point
	}
	for stageIndex, stage := range pack.Stages {
		for idx, id := range stage.Player {
			if id == 27 {
				locations = append(locations, struct {
					stage int
					point Point
				}{stage: stageIndex, point: Point{X: idx % stage.Width, Y: idx / stage.Width}})
			}
		}
	}
	if len(locations) != 1 || locations[0].stage != 2 || locations[0].point != (Point{X: 24, Y: 25}) {
		t.Fatalf("Bavaria hook locations=%v, want only displayed Stage 3 at (24,25)", locations)
	}
}

func TestBavariaDemoScriptsMatchAuthoredForegroundEvents(t *testing.T) {
	tests := []struct {
		stage    int
		point    Point
		scriptID int
		commands []tutorialCommand
	}{
		{stage: 3, point: Point{X: 20, Y: 16}, scriptID: 4, commands: []tutorialCommand{
			tutorialCamera(13, 16, 30), tutorialWait(20), tutorialCamera(19, 16, 30),
		}},
		{stage: 8, point: Point{X: 30, Y: 10}, scriptID: 6, commands: []tutorialCommand{
			tutorialCamera(28, 18, 40), tutorialWait(20), tutorialCamera(26, 11, 40), tutorialWait(20),
		}},
		{stage: 9, point: Point{X: 9, Y: 20}, scriptID: 34, commands: []tutorialCommand{
			tutorialPortraitFace(1),
			tutorialPortraitPosition(17, 50),
			tutorialPrompt(34, TutorialTextBubble, 90, 2),
			tutorialPortraitFace(3),
			tutorialFlash(1),
			tutorialPortraitMark(0, 3, true),
			tutorialPrompt(35, TutorialTextBubble, 90, 2),
			tutorialPortraitFace(0),
			tutorialPortraitMark(4, 3, true),
			tutorialPrompt(36, TutorialTextBubble, 90, 2),
		}},
		{stage: 12, point: Point{X: 6, Y: 47}, scriptID: 19, commands: []tutorialCommand{
			tutorialCamera(7, 42, 20), tutorialWait(20), tutorialCamera(13, 56, 45), tutorialWait(20),
		}},
	}

	for _, test := range tests {
		rt := mustLoadBavariaRuntime(t, test.stage)
		idx := rt.index(test.point.X, test.point.Y)
		if rt.Foreground[idx] != 0 || int(rt.Background[idx]) != test.scriptID {
			t.Errorf("Stage %d event at %+v raw/background=%d/%d, want 0/%d", test.stage+1, test.point, rt.Foreground[idx], rt.Background[idx], test.scriptID)
			continue
		}
		commands, ok := rt.demoScriptCommands(test.scriptID)
		if !ok || !slices.Equal(commands, test.commands) {
			t.Errorf("Stage %d demo %d commands=%v, want %v", test.stage+1, test.scriptID, commands, test.commands)
		}
		rt.Player = test.point
		rt.pendingForegroundEvent = test.point
		rt.pendingForegroundEventSet = true
		rt.tickPendingForegroundEvent()
		if !rt.TutorialScriptActive || rt.TutorialScriptID != test.scriptID || rt.CanAcceptInput() || rt.Foreground[idx] != EmptyRawID {
			t.Errorf("Stage %d demo trigger active/id/input/foreground=%v/%d/%v/%d", test.stage+1, rt.TutorialScriptActive, rt.TutorialScriptID, rt.CanAcceptInput(), rt.Foreground[idx])
		}
	}
}

func TestBavariaKnightIntroRunsSourceDialogueSequence(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, teutonicKnightStageIndex)
	point := Point{X: 9, Y: 20}
	rt.Player = point
	rt.pendingForegroundEvent = point
	rt.pendingForegroundEventSet = true
	rt.tickPendingForegroundEvent()

	prompts := make([]int, 0, 3)
	seenFlash := false
	seenMarks := map[int]bool{}
	for tick := 0; tick < 400 && rt.TutorialScriptActive; tick++ {
		rt.tickTutorial()
		seenFlash = seenFlash || rt.TutorialFlashVisible
		if rt.TutorialPortraitMark >= 0 {
			seenMarks[rt.TutorialPortraitMark] = true
		}
		if prompt, ok := rt.TutorialPrompt(); ok && rt.AdvanceTutorialPrompt() {
			prompts = append(prompts, prompt.TextIndex)
		}
	}
	if rt.TutorialScriptActive || !slices.Equal(prompts, []int{34, 35, 36}) {
		t.Fatalf("knight intro active/prompts=%v/%v, want completed [34 35 36]", rt.TutorialScriptActive, prompts)
	}
	if !seenFlash || !seenMarks[0] || !seenMarks[4] || rt.TutorialPortraitFace != 0 {
		t.Fatalf("knight intro flash/marks/face=%v/%v/%d, want flash, marks 0+4, face 0", seenFlash, seenMarks, rt.TutorialPortraitFace)
	}
}

func TestBavariaEnemyGateTriggersUseAllAuthoredGroupCounts(t *testing.T) {
	tests := []struct {
		stage     int
		point     Point
		group     int
		remaining int
	}{
		{stage: 0, point: Point{X: 7, Y: 3}, group: 0, remaining: 2},
		{stage: 3, point: Point{X: 30, Y: 25}, group: 0, remaining: 2},
		{stage: 4, point: Point{X: 38, Y: 14}, group: 1, remaining: 1},
		{stage: 4, point: Point{X: 16, Y: 17}, group: 0, remaining: 1},
		{stage: 5, point: Point{X: 35, Y: 15}, group: 0, remaining: 2},
		{stage: 5, point: Point{X: 8, Y: 17}, group: 1, remaining: 2},
		{stage: 6, point: Point{X: 26, Y: 4}, group: 0, remaining: 1},
		{stage: 6, point: Point{X: 23, Y: 15}, group: 1, remaining: 2},
		{stage: 6, point: Point{X: 26, Y: 15}, group: 2, remaining: 2},
		{stage: 7, point: Point{X: 40, Y: 4}, group: 2, remaining: 2},
		{stage: 7, point: Point{X: 45, Y: 9}, group: 1, remaining: 2},
		{stage: 7, point: Point{X: 12, Y: 11}, group: 0, remaining: 2},
		{stage: 8, point: Point{X: 11, Y: 5}, group: 0, remaining: 2},
		{stage: 8, point: Point{X: 10, Y: 29}, group: 1, remaining: 2},
		{stage: 9, point: Point{X: 13, Y: 20}, group: 0, remaining: 1},
		{stage: 12, point: Point{X: 17, Y: 27}, group: 0, remaining: 1},
	}
	for _, test := range tests {
		rt := mustLoadBavariaRuntime(t, test.stage)
		idx := rt.index(test.point.X, test.point.Y)
		if rt.Foreground[idx] != 26 || int(rt.Background[idx]) != test.group {
			t.Errorf("Stage %d gate trigger %+v raw/group=%d/%d, want 26/%d", test.stage+1, test.point, rt.Foreground[idx], rt.Background[idx], test.group)
			continue
		}
		if got := rt.EnemyGateCounters[test.group]; got != test.remaining {
			t.Errorf("Stage %d group %d counter=%d, want authored %d", test.stage+1, test.group, got, test.remaining)
		}
	}
}

func TestBavariaWaterBlocksHeroWithoutPotionAndBuoysGravityObjects(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 7)
	if len(rt.WaterSources) != 3 {
		t.Fatalf("water sources=%v, want the three authored Stage 8 sources", rt.WaterSources)
	}
	if !rt.WaterInitializing || !rt.CanAcceptInput() {
		t.Fatalf("initial water state initializing/input=%v/%v, want true/true", rt.WaterInitializing, rt.CanAcceptInput())
	}
	var vegetation Point
	for idx, id := range rt.PlayerLayer {
		if id == 10 {
			vegetation = Point{X: idx % rt.Width(), Y: idx / rt.Width()}
			break
		}
	}
	rt.TickSourceFrame(20, 1, 4)
	if rt.IsPassable(vegetation.X, vegetation.Y) {
		t.Fatalf("vegetation %+v is passable while source water state is unstable", vegetation)
	}
	if got := rt.WaterAt(7, 12); got != 1 {
		t.Fatalf("water at source frame 1=%d, want first packed sub-cell", got)
	}
	for sourceTick := 2; sourceTick <= 240 && rt.WaterInitializing; sourceTick++ {
		rt.TickSourceFrame(20, sourceTick, 4)
	}
	if rt.WaterInitializing || !rt.CanAcceptInput() {
		t.Fatalf("settled water initializing/input=%v/%v, want false/true", rt.WaterInitializing, rt.CanAcceptInput())
	}
	if !rt.IsPassable(vegetation.X, vegetation.Y) {
		t.Fatalf("vegetation %+v remained blocked after source water settled", vegetation)
	}
	totalDepth := 0
	passableWater := Point{X: -1, Y: -1}
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			depth := int(rt.WaterAt(x, y))
			totalDepth += depth
			idx := rt.index(x, y)
			if depth > 0 && rt.PlayerLayer[idx] == EmptyRawID && rt.foregroundPassable(x, y, rt.Foreground[idx]) {
				passableWater = Point{X: x, Y: y}
			}
		}
	}
	if totalDepth != 99 {
		t.Fatalf("settled packed water depth=%d, want source-JAR shape total 99", totalDepth)
	}
	for _, point := range []Point{{X: 7, Y: 12}, {X: 25, Y: 14}, {X: 33, Y: 5}} {
		if got := rt.WaterAt(point.X, point.Y); got != 0 {
			t.Fatalf("cleaned source outlet %+v depth=%d, want 0", point, got)
		}
	}
	for _, point := range []Point{{X: 6, Y: 16}, {X: 10, Y: 16}, {X: 28, Y: 16}, {X: 30, Y: 7}, {X: 34, Y: 7}} {
		if got := rt.WaterAt(point.X, point.Y); got != 3 {
			t.Fatalf("settled pool %+v depth=%d, want 3", point, got)
		}
	}
	if passableWater.X < 0 || rt.IsPassable(passableWater.X, passableWater.Y) {
		t.Fatalf("water cell %+v should block the hero before the potion", passableWater)
	}
	rt.SpecialItemMask |= 4
	if !rt.IsPassable(passableWater.X, passableWater.Y) {
		t.Fatalf("water cell %+v should be passable after the potion", passableWater)
	}

	for y := 14; y <= 15; y++ {
		rt.SetForTest(PlayerLayer, 15, y, EmptyRawID)
		rt.SetForTest(ForegroundLayer, 15, y, EmptyRawID)
		rt.water.Cells[rt.index(15, y)] = 0
	}
	rt.SetForTest(PlayerLayer, 15, 15, 0)
	rt.water.Cells[rt.index(15, 15)] = waterCellSet(0, 0, 1, 0, 3)
	rt.syncWaterState()
	if !rt.tickGravityObjectAt(15, 15) || rt.PlayerLayer[rt.index(15, 14)] != 0 {
		t.Fatal("a submerged boulder did not move upward with source buoyancy")
	}
}

func TestBavariaStageEightWaterSourcesStartInSourceColumnOrder(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 7)
	secondSource := rt.index(25, 13)
	thirdSource := rt.index(33, 4)
	for sourceTick := 1; sourceTick <= 82; sourceTick++ {
		rt.TickSourceFrame(20, sourceTick, 4)
	}
	if got := rt.WaterAt(25, 14); got != 0 || rt.PlayerLayer[secondSource] != 38 {
		t.Fatalf("second source before frame 83 outlet/raw=%d/%d, want 0/38", got, rt.PlayerLayer[secondSource])
	}
	rt.TickSourceFrame(20, 83, 4)
	if got := rt.WaterAt(25, 14); got != 1 || rt.PlayerLayer[secondSource] != EmptyRawID {
		t.Fatalf("second source at frame 83 outlet/raw=%d/%d, want 1/empty", got, rt.PlayerLayer[secondSource])
	}
	for sourceTick := 84; sourceTick <= 129; sourceTick++ {
		rt.TickSourceFrame(20, sourceTick, 4)
	}
	if got := rt.WaterAt(33, 5); got != 0 || rt.PlayerLayer[thirdSource] != 38 {
		t.Fatalf("third source before frame 130 outlet/raw=%d/%d, want 0/38", got, rt.PlayerLayer[thirdSource])
	}
	rt.TickSourceFrame(20, 130, 4)
	if got := rt.WaterAt(33, 5); got != 1 || rt.PlayerLayer[thirdSource] != EmptyRawID {
		t.Fatalf("third source at frame 130 outlet/raw=%d/%d, want 1/empty", got, rt.PlayerLayer[thirdSource])
	}
}

func TestBavariaAuthoredWaterSettlesToPackedSourceShapes(t *testing.T) {
	expectedDepths := map[int]int{7: 99, 8: 72, 10: 48, 11: 48, 12: 129}
	for stageIndex, expectedDepth := range expectedDepths {
		rt := mustLoadBavariaRuntime(t, stageIndex)
		for sourceTick := 1; sourceTick <= 2000 && rt.WaterInitializing; sourceTick++ {
			rt.TickSourceFrame(20, sourceTick, 4)
		}
		actualDepth := 0
		for _, depth := range rt.WaterDepth {
			actualDepth += int(depth)
		}
		if rt.WaterInitializing || rt.water.Phase != 3 || actualDepth != expectedDepth {
			t.Fatalf("Stage %d settled initializing/depth=%v/%d, want false/%d", stageIndex+1, rt.WaterInitializing, actualDepth, expectedDepth)
		}
	}
}

func TestBavariaCheckpointRestoresPackedWaterState(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 7)
	for sourceTick := 1; sourceTick <= 40; sourceTick++ {
		rt.TickSourceFrame(20, sourceTick, 4)
	}
	rt.SaveSnapshot()
	wantWater := rt.water.clone()
	wantDepth := append([]uint8(nil), rt.WaterDepth...)
	for sourceTick := 41; sourceTick <= 100; sourceTick++ {
		rt.TickSourceFrame(20, sourceTick, 4)
	}
	if !rt.RestoreCheckpoint() {
		t.Fatal("packed-water checkpoint restore failed")
	}
	if !reflect.DeepEqual(rt.water, wantWater) || !slices.Equal(rt.WaterDepth, wantDepth) {
		t.Fatal("packed water records or query depth did not restore exactly")
	}
	rt.water.Cells[0] ^= 1
	if reflect.DeepEqual(rt.water.Cells, rt.checkpoint.Water.Cells) {
		t.Fatal("restored water cells still alias the checkpoint snapshot")
	}
}

func TestBavariaDestroyedWaterBarriersTriggerSourceReflow(t *testing.T) {
	for _, test := range []struct {
		name string
		raw  RawID
		act  func(*Runtime, int, int)
	}{
		{name: "vegetation", raw: 10, act: func(rt *Runtime, _, _ int) { rt.TickSourceFrame(20, 1, 4) }},
		{name: "blast wall", raw: 37, act: func(rt *Runtime, _, _ int) { rt.TickBreakables() }},
	} {
		t.Run(test.name, func(t *testing.T) {
			rt := mustLoadBavariaRuntime(t, 0)
			x, y := 10, 10
			rt.water.StartNext = false
			rt.water.SourceCount = 0
			rt.water.SourceIndex = 0
			rt.water.Phase = 3
			rt.water.Cells[rt.index(x-1, y)] = waterCellSet(0, 0, 1, 0, 3)
			rt.SetForTest(PlayerLayer, x, y, test.raw)
			rt.SetForTest(ForegroundLayer, x, y, EmptyRawID)
			rt.ObjectState[rt.index(x, y)] = 8
			rt.DrainSoundEvents()
			test.act(rt, x, y)
			if sounds := rt.DrainSoundEvents(); !slices.Contains(sounds, SoundWater) {
				t.Fatalf("destroyed raw %d sounds=%v, want water reflow sound", test.raw, sounds)
			}
		})
	}
}

func TestBavariaWindPodBuildsAndAdvancesSourceForegroundColumn(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, 0)
	x, y := 10, 10
	for _, point := range []Point{{X: x, Y: y - 2}, {X: x, Y: y - 1}, {X: x, Y: y}, {X: x - 1, Y: y}, {X: x + 1, Y: y}, {X: x, Y: y + 1}} {
		rt.SetForTest(PlayerLayer, point.X, point.Y, EmptyRawID)
		rt.SetForTest(ForegroundLayer, point.X, point.Y, EmptyRawID)
	}
	rt.SetForTest(PlayerLayer, x, y, 47)
	rt.SetForTest(PlayerLayer, x-1, y, 80)
	rt.SetForTest(PlayerLayer, x+1, y, 80)
	rt.SetForTest(PlayerLayer, x, y+1, 80)
	rt.tickWindPodAt(x, y)
	above := rt.index(x, y-1)
	if rt.Foreground[above] != 35 || rt.ForegroundState[above] != 18 {
		t.Fatalf("wind foreground/timer=%d/%d, want 35/18", rt.Foreground[above], rt.ForegroundState[above])
	}
	for range 3 {
		rt.tickWindForegroundAt(x, y-1)
	}
	rt.tickWindForegroundAt(x, y-1)
	if rt.Foreground[rt.index(x, y-2)] != 35 {
		t.Fatal("wind foreground did not advance one cell after the source 18/6 cadence")
	}
}

func TestBavariaStageTenInitializesSourceKnightArena(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, teutonicKnightStageIndex)
	boss := rt.TeutonicKnight
	if !boss.Enabled || boss.State != TeutonicKnightStateDormant || boss.Health != 4 || boss.X != 408 || boss.Animation != 10 {
		t.Fatalf("knight=%+v, want enabled state13 health4 x408 animation10", boss)
	}
	if rt.EnemyGateCounters[0] != 1 {
		t.Fatalf("arena group-0 counter=%d, want the authored sentinel count 1", rt.EnemyGateCounters[0])
	}
	for _, point := range []Point{{X: 12, Y: 20}, {X: 22, Y: 20}} {
		if !rt.foregroundDoorOpen(point.X, point.Y) {
			t.Fatalf("arena door %+v starts closed, want source-open", point)
		}
	}
	for _, point := range []Point{{X: 16, Y: 18}, {X: 19, Y: 18}} {
		if id, _ := rt.At(PlayerLayer, point.X, point.Y); id != 0 {
			t.Fatalf("authored arena boulder %+v raw=%d, want raw0", point, id)
		}
	}
}

func TestBavariaKnightDormantAnimationStartsOnlyFromRightSide(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, teutonicKnightStageIndex)
	rt.Player = Point{X: 17, Y: 20}
	rt.tickEvilTeutonicKnight(1)
	if rt.TeutonicKnight.Animation != TeutonicKnightStateDormant || rt.TeutonicKnight.AnimationTicks != 0 || rt.TeutonicKnight.IntroActivated {
		t.Fatalf("left-side intro=%+v, want dormant animation13 paused", rt.TeutonicKnight)
	}
	rt.Player = Point{X: 18, Y: 20}
	wakeTick := 0
	for tick := 2; tick <= 91; tick++ {
		rt.tickEvilTeutonicKnight(tick)
		if rt.TeutonicKnight.State != TeutonicKnightStateDormant {
			wakeTick = tick
			break
		}
	}
	if wakeTick == 0 || rt.TeutonicKnight.State != TeutonicKnightStateWalkLeft || rt.TeutonicKnight.Animation != TeutonicKnightStateWalkLeft {
		t.Fatalf("awakened knight at source tick %d: %+v, want state/animation0 after animation13", wakeTick, rt.TeutonicKnight)
	}
}

func TestBavariaKnightSlashFrameRespawnsSourceBoulders(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, teutonicKnightStageIndex)
	for _, x := range []int{16, 19} {
		rt.SetForTest(PlayerLayer, x, teutonicKnightRestY, EmptyRawID)
		rt.SetForTest(PlayerLayer, x, teutonicKnightSpawnY, EmptyRawID)
	}
	rt.Player = Point{X: 5, Y: 20}
	rt.TeutonicKnight.State = TeutonicKnightStateSlashLeft
	rt.TeutonicKnight.Animation = TeutonicKnightStateSlashLeft
	rt.TeutonicKnight.AnimationTicks = 7
	rt.tickEvilTeutonicKnight(3)
	for _, x := range []int{16, 19} {
		if id, _ := rt.At(PlayerLayer, x, teutonicKnightSpawnY); id != 0 {
			t.Fatalf("respawn at (%d,%d) raw=%d, want raw0 on slash frame 5", x, teutonicKnightSpawnY, id)
		}
	}
}

func TestBavariaKnightOnlyFallingBouldersDamageAndAllBouldersBreak(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, teutonicKnightStageIndex)
	rt.Player = Point{X: 5, Y: 20}
	rt.TeutonicKnight.State = TeutonicKnightStateWalkLeft
	rt.TeutonicKnight.Animation = TeutonicKnightStateWalkLeft
	centerX := (rt.TeutonicKnight.X + TileSize) / TileSize

	rt.SetForTest(PlayerLayer, centerX, teutonicKnightBoulderTopY, 0)
	rt.ObjectState[rt.index(centerX, teutonicKnightBoulderTopY)] = 0
	hits, _ := rt.tickEvilTeutonicKnight(1)
	if hits != 0 || rt.TeutonicKnight.Health != 4 || rt.PlayerLayer[rt.index(centerX, teutonicKnightBoulderTopY)] != 30 {
		t.Fatalf("stationary impact hits/health/raw=%d/%d/%d, want 0/4/raw30", hits, rt.TeutonicKnight.Health, rt.PlayerLayer[rt.index(centerX, teutonicKnightBoulderTopY)])
	}

	rt.SetForTest(PlayerLayer, centerX, teutonicKnightBoulderTopY, 0)
	rt.ObjectState[rt.index(centerX, teutonicKnightBoulderTopY)] = 3
	hits, _ = rt.tickEvilTeutonicKnight(2)
	if hits != 1 || rt.TeutonicKnight.Health != 3 || rt.TeutonicKnight.State != TeutonicKnightStateHurtLeft {
		t.Fatalf("falling impact hits/health/state=%d/%d/%d, want 1/3/hurt-left", hits, rt.TeutonicKnight.Health, rt.TeutonicKnight.State)
	}
}

func TestBavariaKnightFourHitsRunsDeathDelayAndOpensArena(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, teutonicKnightStageIndex)
	rt.Player = Point{X: 5, Y: 20}
	rt.PlayerMotion = ObjectMotion{DX: 1}
	rt.activateEnemyGateTriggerAt(13, 20)
	for _, point := range []Point{{X: 12, Y: 20}, {X: 22, Y: 20}} {
		rt.closeDoorAt(point.X, point.Y)
		if rt.foregroundDoorOpen(point.X, point.Y) {
			t.Fatalf("arena door %+v did not close", point)
		}
	}
	rt.EnemyGateDemoActive = false
	rt.TeutonicKnight.State = TeutonicKnightStateWalkLeft
	rt.TeutonicKnight.Animation = TeutonicKnightStateWalkLeft

	for hit := 1; hit <= 4; hit++ {
		centerX := (rt.TeutonicKnight.X + TileSize) / TileSize
		rt.SetForTest(PlayerLayer, centerX, teutonicKnightBoulderTopY, 0)
		rt.ObjectState[rt.index(centerX, teutonicKnightBoulderTopY)] = 3
		got, defeated := rt.tickEvilTeutonicKnight(hit)
		if got != 1 || defeated != (hit == 4) {
			t.Fatalf("hit %d result=%d defeated=%v", hit, got, defeated)
		}
	}
	if rt.TeutonicKnight.State != TeutonicKnightStateDefeated || rt.EnemyGateCounters[0] != 1 {
		t.Fatalf("lethal state/counter=%d/%d, want defeated/1 during explosion delay", rt.TeutonicKnight.State, rt.EnemyGateCounters[0])
	}
	for tick := 1; tick <= 101; tick++ {
		rt.tickEvilTeutonicKnight(100 + tick)
	}
	if rt.TeutonicKnight.State != TeutonicKnightStateDefeated || rt.EnemyGateCounters[0] != 1 {
		t.Fatalf("death ended at tick 101: state/counter=%d/%d", rt.TeutonicKnight.State, rt.EnemyGateCounters[0])
	}
	rt.tickEvilTeutonicKnight(202)
	if rt.TeutonicKnight.State != TeutonicKnightStateComplete || rt.EnemyGateCounters[0] != 0 {
		t.Fatalf("death completion state/counter=%d/%d, want complete/0", rt.TeutonicKnight.State, rt.EnemyGateCounters[0])
	}
	for _, point := range []Point{{X: 12, Y: 20}, {X: 22, Y: 20}} {
		state, _ := rt.At(BackgroundLayer, point.X, point.Y)
		if int(state)&0xf0 != 0x10 {
			t.Fatalf("arena door %+v state=%#x, want source opening phase1", point, state)
		}
	}
}

func TestBavariaSealChestUsesGenericSourceTransition(t *testing.T) {
	rt := mustLoadBavariaRuntime(t, teutonicKnightStageIndex)
	chest := Point{X: 29, Y: 13}
	if rt.PlayerLayer[rt.index(chest.X, chest.Y)] != 51 || rt.Foreground[rt.index(chest.X, chest.Y)] != 14 {
		t.Fatal("Bavaria seal source chest moved")
	}
	rt.Player = chest
	rt.startChestOpening(chest, false)
	for tick := 0; tick < chestRewardTick; tick++ {
		rt.TickStatus()
	}
	if rt.RelicMask&2 == 0 || !rt.SealCollected || rt.SealStageComplete || rt.CanAcceptInput() {
		t.Fatalf("seal relic/collected/complete/input=%#x/%v/%v/%v", rt.RelicMask, rt.SealCollected, rt.SealStageComplete, rt.CanAcceptInput())
	}
	rt.TeutonicKnight.State = TeutonicKnightStateComplete
	for tick := 1; tick <= sealCompletionTicks; tick++ {
		rt.tickEvilTeutonicKnight(tick)
	}
	if rt.SealStageComplete {
		t.Fatal("Bavaria seal completed at tick 140, want source condition >140")
	}
	rt.tickEvilTeutonicKnight(sealCompletionTicks + 1)
	if !rt.SealStageComplete {
		t.Fatal("Bavaria seal did not complete at tick 141")
	}
}

func mustLoadBavariaRuntime(t *testing.T, stageIndex int) *Runtime {
	t.Helper()
	pack, err := LoadWorldDir(filepath.Join("..", "..", "decoded", "world1"))
	if err != nil {
		t.Fatal(err)
	}
	rt, err := NewRuntime(pack.Stages[stageIndex])
	if err != nil {
		t.Fatal(err)
	}
	return rt
}
