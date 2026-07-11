package originalgame

import (
	"path/filepath"
	"testing"
)

func TestAngkorWorldMapMatchesSourceNodes(t *testing.T) {
	worldMap, err := loadWorldMap(filepath.Join(defaultWorldDir, "map.json"))
	if err != nil {
		t.Fatal(err)
	}
	if worldMap.Source != "map_angkor.out" || worldMap.PayloadLength != 113 || len(worldMap.Nodes) != 13 {
		t.Fatalf("map source=%q payload=%d nodes=%d, want map_angkor.out/113/13", worldMap.Source, worldMap.PayloadLength, len(worldMap.Nodes))
	}
	wants := []worldMapNode{
		{X: 1, Y: 3, Stage: 0},
		{X: 4, Y: 4, Stage: 1},
		{X: 7, Y: 3, Stage: 2},
		{X: 10, Y: 4, Stage: 3},
		{X: 8, Y: 6, Stage: 4},
		{X: 8, Y: 8, Stage: 5},
	}
	for _, want := range wants {
		got, ok := worldMap.nodeForStage(want.Stage)
		if !ok || got.X != want.X || got.Y != want.Y || got.Type != 0 {
			t.Errorf("stage %d node=%+v ok=%v, want (%d,%d) type0", want.Stage, got, ok, want.X, want.Y)
		}
	}
}

func TestAngkorWorldMapNavigationUsesSourceLinksAndUnlocks(t *testing.T) {
	worldMap, err := loadWorldMap(filepath.Join(defaultWorldDir, "map.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := worldMap.linkedStage(0, 1, 0, unlockedThrough(0)); ok {
		t.Fatal("Stage 1 navigated to locked Stage 2")
	}
	if got, ok := worldMap.linkedStage(0, 1, 0, unlockedThrough(1)); !ok || got != 1 {
		t.Fatalf("right from Stage 1=%d,%v, want Stage 2", got, ok)
	}
	if got, ok := worldMap.linkedStage(1, -1, 0, unlockedThrough(2)); !ok || got != 0 {
		t.Fatalf("left from Stage 2=%d,%v, want Stage 1", got, ok)
	}
	if got, ok := worldMap.linkedStage(1, 1, 0, unlockedThrough(2)); !ok || got != 2 {
		t.Fatalf("right from Stage 2=%d,%v, want Stage 3", got, ok)
	}
	if got, ok := worldMap.linkedStage(4, 0, 1, unlockedThrough(5)); !ok || got != 5 {
		t.Fatalf("down from Stage 5=%d,%v, want unlocked Stage 6", got, ok)
	}
}

func TestAngkorWorldMapExitTargetsMatchSourceBranches(t *testing.T) {
	worldMap, err := loadWorldMap(filepath.Join(defaultWorldDir, "map.json"))
	if err != nil {
		t.Fatal(err)
	}
	checks := []struct {
		stage  int
		secret bool
		want   int
	}{
		{stage: 6, want: 7},
		{stage: 6, secret: true, want: 9},
		{stage: 7, want: 8},
		{stage: 7, secret: true, want: 12},
		{stage: 9, secret: true, want: 10},
		{stage: 10, secret: true, want: 11},
		{stage: 11, secret: true, want: 11},
		{stage: 12, secret: true, want: 12},
	}
	for _, check := range checks {
		got, ok := worldMap.exitTarget(check.stage, check.secret)
		if !ok || got != check.want {
			t.Errorf("stage %d secret=%v target=%d,%v, want %d,true", check.stage, check.secret, got, ok, check.want)
		}
	}
}

func TestAngkorWorldMapUsesSecretStageTitles(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if got := g.worldMapStageTitle(6); got != "STAGE 7" {
		t.Fatalf("normal map title=%q, want STAGE 7", got)
	}
	if got := g.worldMapStageTitle(9); got != "SECRET STAGE 1" {
		t.Fatalf("secret map title=%q, want SECRET STAGE 1", got)
	}
	if got := g.worldMapStageTitle(12); got != "SECRET STAGE 4" {
		t.Fatalf("secret map title=%q, want SECRET STAGE 4", got)
	}
}

func TestCompletedResultContinuesToUnlockedWorldMapNode(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.rt.VioletGems = g.rt.TotalVioletGems
	g.beginStageResults()
	g.resultPhase = resultPhaseComplete
	g.updateStageResults(true)
	if g.mode != gameModeWorldMap || g.progress.HighestUnlocked != 1 {
		t.Fatalf("continue mode=%d unlocked=%d, want world map/Stage 2", g.mode, g.progress.HighestUnlocked)
	}
	if g.worldMapTravelFrom != 0 || g.worldMapTravelTo != 1 {
		t.Fatalf("map travel=%d->%d, want Stage 1->Stage 2", g.worldMapTravelFrom, g.worldMapTravelTo)
	}
	for g.worldMapLoadingStep < worldMapLoadingSteps {
		g.updateWorldMap(false)
	}
	for g.worldMapTravelTick < worldMapTravelTicks {
		g.updateWorldMap(false)
	}
	if g.worldMapSelectedStage != 1 {
		t.Fatalf("selected stage=%d, want newly unlocked Stage 2", g.worldMapSelectedStage)
	}
}

func unlockedThrough(highest int) func(int) bool {
	return func(stage int) bool {
		return stage >= 0 && stage <= highest
	}
}
