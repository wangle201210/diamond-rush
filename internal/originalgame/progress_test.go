package originalgame

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/wangle201210/zskc/internal/original"
)

func TestOriginalProgressRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	want := newOriginalProgress()
	want.StageCleared[0] = true
	want.StageAwards[0] = resultAwardVioletGems | resultAwardNoHits
	want.HighestUnlocked = 1
	want.VioletGemBank = 17
	want.TutorialComplete = true
	want.StageConsumedRewards[0] = []original.Point{{X: 19, Y: 2}}
	if err := saveOriginalProgress(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("progress=%+v, want %+v", got, want)
	}
}

func TestOriginalProgressMigratesStage1OnlySave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	if err := os.WriteFile(path, []byte("{\n  \"stage1_cleared\": true,\n  \"stage1_awards\": 20\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != originalProgressVersion || !got.StageCleared[0] || got.StageAwards[0] != 20 {
		t.Fatalf("migrated progress=%+v, want Stage 1 clear/awards 20", got)
	}
	if got.HighestUnlocked != 1 || got.ExtraLives != 5 || got.MaxHealth != 4 {
		t.Fatalf("migrated globals unlocked=%d lives=%d health=%d, want 1/5/4", got.HighestUnlocked, got.ExtraLives, got.MaxHealth)
	}
	if !got.TutorialComplete {
		t.Fatal("legacy Stage-1 save did not migrate tutorial as complete")
	}
}

func TestOriginalProgressMigratesSequentialVersion2Unlocks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	if err := os.WriteFile(path, []byte("{\n  \"version\": 2,\n  \"highest_unlocked\": 5,\n  \"extra_lives\": 7,\n  \"max_health\": 4\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	for stage := 0; stage <= 5; stage++ {
		if !got.StageUnlocked[stage] {
			t.Errorf("migrated Stage %d is locked", stage+1)
		}
	}
	if got.StageUnlocked[6] || got.HighestUnlocked != 5 || got.ExtraLives != 7 {
		t.Fatalf("migrated v2 progress=%+v, want only stages 0..5 unlocked and 7 lives", got)
	}
	if !got.TutorialComplete {
		t.Fatal("version 2 save did not migrate tutorial as complete")
	}
}

func TestOriginalProgressMigratesVersion3WithoutSequentialUnlocks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	data := []byte("{\n  \"version\": 3,\n  \"highest_unlocked\": 10,\n  \"stage_unlocked\": [true, false, false, false, false, false, true, false, false, true, true, false, false, false],\n  \"extra_lives\": 5,\n  \"max_health\": 4\n}\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.StageUnlocked[0] || !got.StageUnlocked[6] || !got.StageUnlocked[9] || !got.StageUnlocked[10] {
		t.Fatalf("version 3 explicit unlocks were lost: %v", got.StageUnlocked)
	}
	for _, stage := range []int{1, 2, 3, 4, 5, 7, 8, 11, 12, 13} {
		if got.StageUnlocked[stage] {
			t.Errorf("version 3 migration unexpectedly unlocked stage %d", stage)
		}
	}
	if !got.TutorialComplete {
		t.Fatal("version 3 save did not migrate tutorial as complete")
	}
}

func TestOriginalProgressMigratesVersion4AngkorSeal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	data := []byte("{\n  \"version\": 4,\n  \"stage_unlocked\": [true, false, false, false, false, false, false, false, true, false, false, false, false, false],\n  \"stage_cleared\": [false, false, false, false, false, false, false, false, true, false, false, false, false, false],\n  \"extra_lives\": 5,\n  \"max_health\": 4,\n  \"tutorial_complete\": true\n}\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != originalProgressVersion || got.RelicMask != 1 || !got.WorldUnlocked[0] {
		t.Fatalf("migrated v4 seal progress=%+v, want v5 Angkor relic and world unlocked", got)
	}
}

func TestConsumedRewardCoordinateReloadsAsOpenedChest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	chest := original.Point{X: 19, Y: 2}
	g.rt.ConsumedRewardCells[chest.X+chest.Y*g.rt.Width()] = true
	g.progress.recordStageCollections(0, g.rt, 0)
	if err := saveOriginalProgress(path, g.progress); err != nil {
		t.Fatal(err)
	}

	reloaded, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := reloaded.enableProgress(path); err != nil {
		t.Fatal(err)
	}
	if id, _ := reloaded.rt.At(original.PlayerLayer, chest.X, chest.Y); id != original.EmptyRawID {
		t.Fatalf("reloaded consumed chest payload=%d, want empty", id)
	}
	if got := reloaded.rt.ObjectState[chest.X+chest.Y*reloaded.rt.Width()]; got != 3 {
		t.Fatalf("reloaded consumed chest state=%d, want open frame 3", got)
	}
	if !reflect.DeepEqual(reloaded.progress.StageConsumedRewards[0], []original.Point{chest}) {
		t.Fatalf("reloaded consumed coordinates=%v, want [%+v]", reloaded.progress.StageConsumedRewards[0], chest)
	}
}

func TestLegacyCompleteRewardsInferOnlyUnambiguousCoordinates(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	legacy := newOriginalProgress()
	legacy.StageRedDiamonds[0] = 1
	if got := persistentRewardsForStage(g.rt, 0, legacy); !reflect.DeepEqual(got, []original.Point{{X: 19, Y: 2}}) {
		t.Fatalf("legacy complete Stage 1 red coordinates=%v, want [(19,2)]", got)
	}

	stageFour, err := original.NewRuntime(g.pack.Stages[3])
	if err != nil {
		t.Fatal(err)
	}
	legacy.StageRedDiamonds[3] = stageFour.TotalRedDiamonds - 1
	if got := persistentRewardsForStage(stageFour, 3, legacy); len(got) != 0 {
		t.Fatalf("legacy partial red collection inferred ambiguous coordinates=%v", got)
	}

	seal, err := original.NewRuntime(g.pack.Stages[angkorSealStage])
	if err != nil {
		t.Fatal(err)
	}
	legacy.RelicMask = 1
	if got := persistentRewardsForStage(seal, angkorSealStage, legacy); !reflect.DeepEqual(got, []original.Point{{X: 27, Y: 6}}) {
		t.Fatalf("legacy Angkor seal coordinates=%v, want [(27,6)]", got)
	}
}

func TestSecretExitUnlocksOnlyItsSourceMapTarget(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	progress := newOriginalProgress()
	for stage := 1; stage <= 6; stage++ {
		progress.unlockStage(stage)
	}
	rt := g.rt
	rt.VioletGems = 3
	rt.RedDiamonds = 1
	progress.recordSecretExit(6, 9, rt)
	if !progress.StageUnlocked[9] || progress.StageUnlocked[7] || progress.StageUnlocked[8] {
		t.Fatalf("secret unlocks=%v, want target 9 without stages 7/8", progress.StageUnlocked)
	}
	if progress.StageCleared[6] {
		t.Fatal("normal Stage 7 was marked cleared by its secret exit")
	}
	if progress.HighestUnlocked != 9 || progress.VioletGemBank != 3 || progress.RedDiamondBank != 1 {
		t.Fatalf("secret progress=%+v, want highest 9 and collected currencies", progress)
	}

	progress.recordSecretStageResult(9, 10, rt)
	if !progress.StageCleared[9] || !progress.StageUnlocked[10] {
		t.Fatalf("secret-stage chain progress=%+v, want stage 9 clear and stage 10 unlocked", progress)
	}
}

func TestStageResultOnlyAnimatesNewPersistentAwards(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.enableProgress(path); err != nil {
		t.Fatal(err)
	}
	g.rt.VioletGems = g.rt.TotalVioletGems
	g.rt.RedDiamonds = g.rt.TotalRedDiamonds
	g.beginStageResults()
	if g.resultNewAwards != g.resultAwards || g.resultAwards == 0 {
		t.Fatalf("first clear new=%#x awards=%#x, want all earned awards new", g.resultNewAwards, g.resultAwards)
	}
	if !g.progress.StageCleared[0] || g.progress.HighestUnlocked != 1 {
		t.Fatalf("first clear progress=%+v, want Stage 1 clear and Stage 2 unlocked", g.progress)
	}
	if got, want := g.progress.ExtraLives, 5+4; got != want {
		t.Fatalf("first perfect clear lives=%d, want %d from four new award lives", got, want)
	}

	g2, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := g2.enableProgress(path); err != nil {
		t.Fatal(err)
	}
	g2.rt.VioletGems = g2.rt.TotalVioletGems
	g2.rt.RedDiamonds = g2.rt.TotalRedDiamonds
	g2.beginStageResults()
	if g2.resultAwards != g.resultAwards || g2.resultNewAwards != 0 {
		t.Fatalf("repeat clear awards=%#x new=%#x, want %#x/0", g2.resultAwards, g2.resultNewAwards, g.resultAwards)
	}
	if g2.progress.ExtraLives != g.progress.ExtraLives {
		t.Fatalf("repeat clear lives=%d, want unchanged %d", g2.progress.ExtraLives, g.progress.ExtraLives)
	}
	if g2.progress.VioletGemBank != g.progress.VioletGemBank+g2.rt.VioletGems {
		t.Fatalf("repeat clear violet bank=%d, want replay-farmable %d", g2.progress.VioletGemBank, g.progress.VioletGemBank+g2.rt.VioletGems)
	}
}

func TestStageCollectionsAccumulatePartialRedDiamondsAcrossRuns(t *testing.T) {
	progress := newOriginalProgress()
	first := &original.Runtime{VioletGems: 12, RedDiamonds: 2, TotalRedDiamonds: 4, ExtraLives: 5, MaxHealth: 4}
	progress.recordStageCollections(3, first, 0)
	second := &original.Runtime{VioletGems: 7, RedDiamonds: 1, TotalRedDiamonds: 4, ExtraLives: 5, MaxHealth: 4}
	progress.recordStageCollections(3, second, 0)
	if progress.VioletGemBank != 19 || progress.StageVioletGems[3] != 12 {
		t.Fatalf("violet progress bank=%d best=%d, want 19/12", progress.VioletGemBank, progress.StageVioletGems[3])
	}
	if progress.RedDiamondBank != 3 || progress.StageRedDiamonds[3] != 3 {
		t.Fatalf("red progress bank=%d stage=%d, want 3/3", progress.RedDiamondBank, progress.StageRedDiamonds[3])
	}
}

func TestBavariaPurpleChestAndMysticHookPersistAcrossStages(t *testing.T) {
	pack, err := original.LoadWorldDir(filepath.Join("..", "..", "decoded", "world1"))
	if err != nil {
		t.Fatal(err)
	}
	progress := newOriginalProgress()
	progress.WorldUnlocked[original.WorldBavaria] = true

	purpleStage, err := original.NewRuntime(pack.Stages[0])
	if err != nil {
		t.Fatal(err)
	}
	openCampaignChestForTest(t, purpleStage, original.Point{X: 25, Y: 14})
	if purpleStage.VioletGems != 10 {
		t.Fatalf("purple chest stage counter=%d, want authored value 10", purpleStage.VioletGems)
	}
	progress.recordStageCollections(0, purpleStage, 0)

	hookStage, err := original.NewRuntime(pack.Stages[2])
	if err != nil {
		t.Fatal(err)
	}
	openCampaignChestForTest(t, hookStage, original.Point{X: 24, Y: 25})
	if hookStage.SpecialItemMask&2 == 0 {
		t.Fatal("displayed Stage 3 source chest did not award the Mystic Hook")
	}
	progress.recordStageCollections(2, hookStage, 0)

	path := filepath.Join(t.TempDir(), "original-progress.json")
	if err := saveOriginalProgress(path, progress); err != nil {
		t.Fatal(err)
	}
	reloaded, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.VioletGemBank != 10 || reloaded.BavariaStageVioletGems[0] != 10 {
		t.Fatalf("saved purple chest bank/stage=%d/%d, want 10/10", reloaded.VioletGemBank, reloaded.BavariaStageVioletGems[0])
	}
	if reloaded.ToolLevel != 2 {
		t.Fatalf("saved tool level=%d, want Mystic Hook level 2", reloaded.ToolLevel)
	}

	stageFive, err := original.NewRuntime(pack.Stages[4])
	if err != nil {
		t.Fatal(err)
	}
	game := &Game{progress: reloaded}
	game.applyCampaignProgress(stageFive, 4)
	if stageFive.SpecialItemMask&2 == 0 {
		t.Fatal("displayed Stage 5 did not inherit the Mystic Hook acquired in displayed Stage 3")
	}
	for _, id := range stageFive.Stage.Player {
		if id == 27 {
			t.Fatal("displayed Stage 5 contains a second Mystic Hook, absent from original w1.bin")
		}
	}
}

func openCampaignChestForTest(t *testing.T, rt *original.Runtime, point original.Point) {
	t.Helper()
	rt.Player = original.Point{X: point.X - 1, Y: point.Y}
	if !rt.TryMove(1, 0) {
		t.Fatalf("failed to enter source chest at %+v", point)
	}
	for rt.PlayerMotion.Remaining > 0 {
		rt.AdvancePlayerMotion()
	}
	if !rt.SettlePlayerMove() {
		t.Fatalf("source chest at %+v did not begin opening", point)
	}
	for rt.ChestOpening {
		rt.TickStatus()
	}
}
