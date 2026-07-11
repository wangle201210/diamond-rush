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
	if _, ok := worldMap.linkedStage(0, 1, 0, 0); ok {
		t.Fatal("Stage 1 navigated to locked Stage 2")
	}
	if got, ok := worldMap.linkedStage(0, 1, 0, 1); !ok || got != 1 {
		t.Fatalf("right from Stage 1=%d,%v, want Stage 2", got, ok)
	}
	if got, ok := worldMap.linkedStage(1, -1, 0, 2); !ok || got != 0 {
		t.Fatalf("left from Stage 2=%d,%v, want Stage 1", got, ok)
	}
	if got, ok := worldMap.linkedStage(1, 1, 0, 2); !ok || got != 2 {
		t.Fatalf("right from Stage 2=%d,%v, want Stage 3", got, ok)
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
