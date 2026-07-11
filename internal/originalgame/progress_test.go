package originalgame

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOriginalProgressRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	want := newOriginalProgress()
	want.StageCleared[0] = true
	want.StageAwards[0] = resultAwardVioletGems | resultAwardNoHits
	want.HighestUnlocked = 1
	want.VioletGemBank = 17
	if err := saveOriginalProgress(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
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
}
