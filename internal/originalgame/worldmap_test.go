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

func TestBavariaWorldMapMatchesSourceNodes(t *testing.T) {
	worldMap, err := loadWorldMap(filepath.Join("decoded", "world1", "map.json"))
	if err != nil {
		t.Fatal(err)
	}
	if worldMap.Source != "map_scotland.out" || worldMap.PayloadLength != 113 || len(worldMap.Nodes) != 13 {
		t.Fatalf("map source=%q payload=%d nodes=%d, want map_scotland.out/113/13", worldMap.Source, worldMap.PayloadLength, len(worldMap.Nodes))
	}
	wantStages := []int{0, 1, 2, 3, 10, 4, 5, 6, 11, 12, 7, 8, 9}
	for index, want := range wantStages {
		if got := worldMap.Nodes[index].Stage; got != want {
			t.Fatalf("node %d stage=%d, want %d", index, got, want)
		}
	}
	if got, ok := worldMap.exitTarget(3, true); !ok || got != 10 {
		t.Fatalf("Stage 4 secret target=%d/%v, want 10", got, ok)
	}
	if got, ok := worldMap.exitTarget(6, true); !ok || got != 11 {
		t.Fatalf("Stage 7 secret target=%d/%v, want 11", got, ok)
	}
	for _, branch := range []struct {
		stage, want int
	}{
		{stage: 10, want: 10},
		{stage: 11, want: 12},
		{stage: 12, want: 12},
	} {
		if got, ok := worldMap.exitTarget(branch.stage, true); !ok || got != branch.want {
			t.Errorf("secret stage %d target=%d/%v, want %d", branch.stage, got, ok, branch.want)
		}
	}
}

func TestWorldMapNavigationUsesUnlockedSourceLinksIncludingSecretNodes(t *testing.T) {
	worldMap, err := loadWorldMap(filepath.Join(defaultWorldDir, "map.json"))
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := worldMap.linkedStage(0, 1, 0, unlockedThrough(0)); ok {
		t.Fatalf("right from Stage 1=%d,%v, want locked Stage 2 skipped", got, ok)
	}
	if got, ok := worldMap.linkedStage(0, 1, 0, unlockedThrough(1)); !ok || got != 1 {
		t.Fatalf("right from Stage 1=%d,%v, want unlocked Stage 2", got, ok)
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
	// Angkor Stage 7 links to normal Stage 8 and Secret Stage 1. Once the
	// secret node is unlocked, the source runtime lock state becomes zero and
	// directional navigation can traverse it like any other node.
	if got, ok := worldMap.linkedStage(6, -1, 0, func(int) bool { return true }); !ok || got != 7 {
		t.Fatalf("left from Stage 7=%d,%v, want normal Stage 8", got, ok)
	}
	if got, ok := worldMap.linkedStage(6, 1, 0, func(int) bool { return true }); !ok || got != 9 {
		t.Fatalf("right from Stage 7=%d,%v, want unlocked Secret Stage 1", got, ok)
	}
	for _, step := range []struct {
		from, dx, dy, want int
	}{
		{from: 11, dy: 1, want: 10},
		{from: 10, dy: 1, want: 9},
		{from: 9, dx: -1, want: 6},
		{from: 6, dx: -1, want: 7},
		{from: 7, dx: -1, want: 8},
	} {
		if got, ok := worldMap.linkedStage(step.from, step.dx, step.dy, func(int) bool { return true }); !ok || got != step.want {
			t.Fatalf("return route from stage %d direction=%d,%d got=%d,%v, want %d", step.from, step.dx, step.dy, got, ok, step.want)
		}
	}
}

func TestBavariaWorldMapRequiresSourceUnlockState(t *testing.T) {
	worldMap, err := loadWorldMap(filepath.Join("decoded", "world1", "map.json"))
	if err != nil {
		t.Fatal(err)
	}
	lockedAfterStageZero := unlockedThrough(0)
	if stageTwo, ok := worldMap.linkedStage(0, 0, 1, lockedAfterStageZero); ok {
		t.Fatalf("down from Bavaria Stage 1=%d,%v, want locked Stage 2 skipped", stageTwo, ok)
	}
	stageTwo, ok := worldMap.linkedStage(0, 0, 1, unlockedThrough(1))
	if !ok || stageTwo != 1 {
		t.Fatalf("down from Bavaria Stage 1=%d,%v, want unlocked Stage 2", stageTwo, ok)
	}
	stageThree, ok := worldMap.linkedStage(stageTwo, 1, 0, lockedAfterStageZero)
	if ok {
		t.Fatalf("right from Bavaria Stage 2=%d,%v, want locked Mystic Hook Stage 3 skipped", stageThree, ok)
	}
}

func TestWorldMapCannotEnterLockedNormalStage(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.mode = gameModeWorldMap
	g.worldMapLoadingStep = worldMapLoadingSteps
	g.worldMapSelectedStage = 1
	g.worldMapTravelTick = worldMapTravelTicks
	g.progress.HighestUnlocked = 0
	g.updateWorldMap(true)
	if g.mode != gameModeWorldMap || g.stageIndex == 1 {
		t.Fatalf("locked normal-stage selection mode/stage=%d/%d, want map to remain active", g.mode, g.stageIndex)
	}
}

func TestTabNavigatesFromStageToUnlockedWorld(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.WorldUnlocked[sealPositionBavaria] = true
	g.progress.unlockStageForWorld(sealPositionAngkor, 3)
	g.progress.unlockStageForWorld(sealPositionBavaria, 2)
	if err := g.switchWorld(sealPositionBavaria); err != nil {
		t.Fatal(err)
	}
	g.loadStage(2)
	g.mode = gameModeStage
	if err := g.updateSource(sourceInput{Navigate: true}); err != nil {
		t.Fatal(err)
	}
	if g.mode != gameModeWorldMap || g.worldMapSelectedStage != 2 {
		t.Fatalf("first Tab mode/stage=%d/%d, want Bavaria map/Stage 3", g.mode, g.worldMapSelectedStage)
	}
	g.worldMapLoadingStep = worldMapLoadingSteps
	g.worldMapTravelTick = worldMapTravelTicks
	if err := g.updateSource(sourceInput{Recall: true}); err != nil {
		t.Fatal(err)
	}
	if g.sealExitActive || g.mode != gameModeWorldMap {
		t.Fatal("Enter must not start world switching from the world map")
	}
	if err := g.updateSource(sourceInput{Navigate: true}); err != nil {
		t.Fatal(err)
	}
	for step := 0; step < sealLoadingSteps; step++ {
		if err := g.updateSource(sourceInput{}); err != nil {
			t.Fatal(err)
		}
	}
	if g.mode != gameModeWorldSelect || g.worldSelectPosition != sealPositionBavaria {
		t.Fatalf("selector mode/position=%d/%d, want current Bavaria world", g.mode, g.worldSelectPosition)
	}
	g.worldSelectPosition = sealPositionAngkor
	g.activateWorldSelectPosition()
	if g.mode != gameModeWorldMap || g.worldIndex != sealPositionAngkor || g.stageIndex != 3 {
		t.Fatalf("Angkor return mode/world/stage=%d/%d/%d, want world map/Angkor/Stage 4", g.mode, g.worldIndex, g.stageIndex)
	}
	if !g.progress.BavariaStageUnlocked[2] {
		t.Fatal("returning to Angkor discarded Bavaria map progress")
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
	if got := g.worldMapStageTitle(6); got != "第7关" {
		t.Fatalf("normal map title=%q, want 第7关", got)
	}
	if got := g.worldMapStageTitle(9); got != "隐藏关卡1" {
		t.Fatalf("secret map title=%q, want 隐藏关卡1", got)
	}
	if got := g.worldMapStageTitle(12); got != "隐藏关卡4" {
		t.Fatalf("secret map title=%q, want 隐藏关卡4", got)
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
