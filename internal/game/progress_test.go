package game

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/wangle201210/zskc/internal/world"
)

func TestProgressRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	progress := defaultProgress(3)
	progress.UnlockedLevel = 2
	progress.BestSteps[0] = 42
	progress.BestScores[0] = 1200
	progress.RedDiamonds[0] = 1
	progress.PurpleDiamonds[0] = 6
	progress.SecretExits[0] = true
	progress.AllPurpleClears[0] = true
	progress.AllRedClears[0] = true
	progress.NoDamageClears[0] = true
	progress.NoRecallClears[0] = true
	progress.NoRestartClears[0] = true
	progress.PurpleBank = 10
	progress.MaxHealthUpgrades = 1
	progress.ArmorUpgrades = 2
	progress.LifeUpgrades = 1
	progress.HasCompass = true
	progress.HasHammer = true
	progress.HasHook = true
	progress.AncientSealOpen = true
	if err := saveProgress(path, progress); err != nil {
		t.Fatal(err)
	}
	got, err := loadProgress(path, 3)
	if err != nil {
		t.Fatal(err)
	}
	if got.UnlockedLevel != 2 || got.BestSteps[0] != 42 || got.BestScores[0] != 1200 || got.RedDiamonds[0] != 1 || got.PurpleDiamonds[0] != 6 || !got.SecretExits[0] || !got.AllPurpleClears[0] || !got.AllRedClears[0] || !got.NoDamageClears[0] || !got.NoRecallClears[0] || !got.NoRestartClears[0] || got.PurpleBank != 10 || got.MaxHealthUpgrades != 1 || got.ArmorUpgrades != 2 || got.LifeUpgrades != 1 || !got.HasCompass || !got.HasHammer || !got.HasHook || !got.AncientSealOpen || len(got.BestSteps) != 3 || len(got.BestScores) != 3 || len(got.RedDiamonds) != 3 || len(got.PurpleDiamonds) != 3 || len(got.SecretExits) != 3 || len(got.AllPurpleClears) != 3 || len(got.AllRedClears) != 3 || len(got.NoDamageClears) != 3 || len(got.NoRecallClears) != 3 || len(got.NoRestartClears) != 3 {
		t.Fatalf("progress = %+v", got)
	}
}

func TestProgressNormalizesCorruptValues(t *testing.T) {
	progress := Progress{UnlockedLevel: 99, BestSteps: []int{10, -5}, BestScores: []int{500, -20}, RedDiamonds: []int{1, -1}, PurpleDiamonds: []int{2, -1}, PurpleBank: -5, MaxHealthUpgrades: -1, ArmorUpgrades: -1, LifeUpgrades: -1}
	normalizeProgress(&progress, 3)
	if progress.UnlockedLevel != 3 {
		t.Fatalf("unlocked = %d, want 3", progress.UnlockedLevel)
	}
	if len(progress.BestSteps) != 3 {
		t.Fatalf("best steps len = %d, want 3", len(progress.BestSteps))
	}
	if progress.BestSteps[1] != 0 {
		t.Fatalf("negative best step was not normalized: %+v", progress.BestSteps)
	}
	if len(progress.BestScores) != 3 {
		t.Fatalf("best scores len = %d, want 3", len(progress.BestScores))
	}
	if progress.BestScores[1] != 0 {
		t.Fatalf("negative best score was not normalized: %+v", progress.BestScores)
	}
	if len(progress.RedDiamonds) != 3 {
		t.Fatalf("red diamonds len = %d, want 3", len(progress.RedDiamonds))
	}
	if progress.RedDiamonds[1] != 0 {
		t.Fatalf("negative red diamonds were not normalized: %+v", progress.RedDiamonds)
	}
	if len(progress.PurpleDiamonds) != 3 {
		t.Fatalf("purple diamonds len = %d, want 3", len(progress.PurpleDiamonds))
	}
	if len(progress.SecretExits) != 3 {
		t.Fatalf("secret exits len = %d, want 3", len(progress.SecretExits))
	}
	if len(progress.AllPurpleClears) != 3 || len(progress.AllRedClears) != 3 {
		t.Fatalf("all-collection flags len = purple %d red %d, want 3/3", len(progress.AllPurpleClears), len(progress.AllRedClears))
	}
	if len(progress.NoDamageClears) != 3 || len(progress.NoRecallClears) != 3 {
		t.Fatalf("perfect flags len = no-damage %d no-recall %d, want 3/3", len(progress.NoDamageClears), len(progress.NoRecallClears))
	}
	if len(progress.NoRestartClears) != 3 {
		t.Fatalf("no-restart flags len = %d, want 3", len(progress.NoRestartClears))
	}
	if progress.PurpleDiamonds[1] != 0 || progress.PurpleBank != 0 || progress.MaxHealthUpgrades != 0 || progress.ArmorUpgrades != 0 || progress.LifeUpgrades != 0 {
		t.Fatalf("purple/upgrade values were not normalized: %+v", progress)
	}
}

func TestCompletionScoreAddsParBonus(t *testing.T) {
	if got := completionScore(700, 80, 100); got != 900 {
		t.Fatalf("completionScore() = %d, want 900", got)
	}
	if got := completionScore(700, 120, 100); got != 700 {
		t.Fatalf("completionScore() = %d, want 700", got)
	}
}

func TestStarRating(t *testing.T) {
	tests := []struct {
		name string
		best int
		par  int
		want int
	}{
		{name: "none", best: 0, par: 100, want: 0},
		{name: "three", best: 100, par: 100, want: 3},
		{name: "two", best: 125, par: 100, want: 2},
		{name: "one", best: 140, par: 100, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := starRating(tt.best, tt.par); got != tt.want {
				t.Fatalf("starRating(%d, %d) = %d, want %d", tt.best, tt.par, got, tt.want)
			}
		})
	}
}

func TestRecordWinUpdatesProgressAndRating(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	g := &Game{
		world:      &world.World{Steps: 80, Score: 600, RedDiamonds: 1, Diamonds: 6, TotalRedDiamonds: 1, TotalDiamonds: 6, SecretExitFound: true},
		levels:     []string{"one", "two"},
		parSteps:   []int{100, 120},
		levelIndex: 0,
		progress:   defaultProgress(2),
		savePath:   path,
	}
	stars, newBest, newBestScore, err := g.recordWin()
	if err != nil {
		t.Fatal(err)
	}
	if stars != 3 {
		t.Fatalf("stars = %d, want 3", stars)
	}
	if !newBest {
		t.Fatal("newBest = false, want true")
	}
	if !newBestScore {
		t.Fatal("newBestScore = false, want true")
	}
	if g.progress.BestSteps[0] != 80 {
		t.Fatalf("best steps = %d, want 80", g.progress.BestSteps[0])
	}
	if g.progress.BestScores[0] != 800 {
		t.Fatalf("best score = %d, want 800", g.progress.BestScores[0])
	}
	if g.progress.RedDiamonds[0] != 1 {
		t.Fatalf("red diamonds = %d, want 1", g.progress.RedDiamonds[0])
	}
	if g.progress.PurpleDiamonds[0] != 6 || g.progress.PurpleBank != 6 {
		t.Fatalf("purple progress = per-level %v bank %d, want level0=6 bank=6", g.progress.PurpleDiamonds, g.progress.PurpleBank)
	}
	if !g.progress.SecretExits[0] {
		t.Fatalf("secret exits = %v, want level0=true", g.progress.SecretExits)
	}
	if !g.progress.AllPurpleClears[0] || !g.progress.AllRedClears[0] {
		t.Fatalf("all-collection flags = purple %v red %v, want level0 true/true", g.progress.AllPurpleClears, g.progress.AllRedClears)
	}
	if !g.progress.NoDamageClears[0] || !g.progress.NoRecallClears[0] || !g.progress.NoRestartClears[0] {
		t.Fatalf("perfect flags = no-damage %v no-recall %v no-restart %v, want level0 true/true/true", g.progress.NoDamageClears, g.progress.NoRecallClears, g.progress.NoRestartClears)
	}
	if g.progress.UnlockedLevel != 2 {
		t.Fatalf("unlocked level = %d, want 2", g.progress.UnlockedLevel)
	}
}

func TestRecordWinKeepsExistingBetterBest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	progress := defaultProgress(2)
	progress.UnlockedLevel = 2
	progress.BestSteps[0] = 70
	progress.BestScores[0] = 1000
	g := &Game{
		world:      &world.World{Steps: 90, Score: 500},
		levels:     []string{"one", "two"},
		parSteps:   []int{100, 120},
		levelIndex: 0,
		progress:   progress,
		savePath:   path,
	}
	stars, newBest, newBestScore, err := g.recordWin()
	if err != nil {
		t.Fatal(err)
	}
	if stars != 3 {
		t.Fatalf("stars = %d, want existing 3-star best", stars)
	}
	if newBest {
		t.Fatal("newBest = true, want false")
	}
	if newBestScore {
		t.Fatal("newBestScore = true, want false")
	}
	if g.progress.BestSteps[0] != 70 {
		t.Fatalf("best steps = %d, want unchanged 70", g.progress.BestSteps[0])
	}
	if g.progress.BestScores[0] != 1000 {
		t.Fatalf("best score = %d, want unchanged 1000", g.progress.BestScores[0])
	}
}

func TestRecordWinDoesNotAwardCleanFlagsAfterDamageOrRecall(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	g := &Game{
		world:      &world.World{Steps: 90, Score: 500, Damaged: true, RecallUsed: true},
		levels:     []string{"one"},
		parSteps:   []int{100},
		levelIndex: 0,
		progress:   defaultProgress(1),
		savePath:   path,
	}
	if _, _, _, err := g.recordWin(); err != nil {
		t.Fatal(err)
	}
	if g.progress.NoDamageClears[0] || g.progress.NoRecallClears[0] {
		t.Fatalf("perfect flags = no-damage %v no-recall %v, want false/false", g.progress.NoDamageClears, g.progress.NoRecallClears)
	}
}

func TestRecordWinDoesNotAwardNoRestartAfterRetry(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	g := &Game{
		world:       &world.World{Steps: 90, Score: 500},
		levels:      []string{"one"},
		parSteps:    []int{100},
		levelIndex:  0,
		progress:    defaultProgress(1),
		savePath:    path,
		restartUsed: true,
	}
	if _, _, _, err := g.recordWin(); err != nil {
		t.Fatal(err)
	}
	if g.progress.NoRestartClears[0] {
		t.Fatalf("no-restart flags = %v, want false", g.progress.NoRestartClears)
	}
}

func TestRecordWinDoesNotAwardAllCollectionFlagsWhenMissingGems(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	g := &Game{
		world:      &world.World{Steps: 90, Score: 500, Diamonds: 2, TotalDiamonds: 3, RedDiamonds: 0, TotalRedDiamonds: 1},
		levels:     []string{"one"},
		parSteps:   []int{100},
		levelIndex: 0,
		progress:   defaultProgress(1),
		savePath:   path,
	}
	if _, _, _, err := g.recordWin(); err != nil {
		t.Fatal(err)
	}
	if g.progress.AllPurpleClears[0] || g.progress.AllRedClears[0] {
		t.Fatalf("all-collection flags = purple %v red %v, want false/false", g.progress.AllPurpleClears, g.progress.AllRedClears)
	}
}

func TestRecordFinalWinOpensAncientSealWhenRedDiamondsAreComplete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	progress := defaultProgress(5)
	progress.UnlockedLevel = 5
	progress.RedDiamonds[3] = 1
	g := &Game{
		world:      &world.World{Steps: 150, Score: 900, RedDiamonds: 2, TotalRedDiamonds: 2},
		levels:     []string{"one", "two", "three", "four", "five"},
		parSteps:   []int{85, 115, 130, 145, 165},
		levelIndex: 4,
		progress:   progress,
		savePath:   path,
	}
	if _, _, _, err := g.recordWin(); err != nil {
		t.Fatal(err)
	}
	if !g.progress.AncientSealOpen {
		t.Fatalf("ancient seal open = false, progress = %+v", g.progress)
	}
	if got := sealStatusText(g.progress, len(g.levels)); got != "Seal OPEN Red 3/3" {
		t.Fatalf("sealStatusText() = %q", got)
	}
}

func TestRecordFinalWinKeepsAncientSealClosedWhenRedDiamondsAreMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	progress := defaultProgress(5)
	progress.UnlockedLevel = 5
	g := &Game{
		world:      &world.World{Steps: 150, Score: 900, RedDiamonds: 2, TotalRedDiamonds: 2},
		levels:     []string{"one", "two", "three", "four", "five"},
		parSteps:   []int{85, 115, 130, 145, 165},
		levelIndex: 4,
		progress:   progress,
		savePath:   path,
	}
	if _, _, _, err := g.recordWin(); err != nil {
		t.Fatal(err)
	}
	if g.progress.AncientSealOpen {
		t.Fatalf("ancient seal open = true, progress = %+v", g.progress)
	}
	if got := sealStatusText(g.progress, len(g.levels)); got != "Seal sealed Red 2/3" {
		t.Fatalf("sealStatusText() = %q", got)
	}
}

func TestClearResultLinesShowCollectionAndMarks(t *testing.T) {
	g := &Game{
		world: &world.World{
			Diamonds:         2,
			TotalDiamonds:    3,
			RedDiamonds:      1,
			TotalRedDiamonds: 1,
			SecretExitFound:  true,
		},
	}
	if got := g.clearGemLine(); got != "Gems P 2/3  R 1/1" {
		t.Fatalf("clearGemLine() = %q", got)
	}
	if got := g.clearMarkLine(); got != "Marks Sec Y  All -/Y  Clean Y" {
		t.Fatalf("clearMarkLine() = %q", got)
	}
	g.world.Damaged = true
	if got := g.clearMarkLine(); got != "Marks Sec Y  All -/Y  Clean -" {
		t.Fatalf("clearMarkLine() after damage = %q", got)
	}
}

func TestLevelSelectLineShowsBestStepAndScore(t *testing.T) {
	progress := defaultProgress(2)
	progress.UnlockedLevel = 1
	progress.BestSteps[0] = 80
	progress.BestScores[0] = 1600
	progress.SecretExits[0] = true
	progress.AllPurpleClears[0] = true
	progress.AllRedClears[0] = true
	progress.NoDamageClears[0] = true
	progress.NoRecallClears[0] = true
	progress.NoRestartClears[0] = true
	g := &Game{
		titles:   []string{"Angkor Gate", "Rolling Stones"},
		parSteps: []int{100, 120},
		progress: progress,
		selected: 0,
	}
	line := g.levelSelectLine(0)
	for _, want := range []string{"> 1", "Angkor Gate", "open", "st80", "sc1600", "sY", "aY", "cY", "***"} {
		if !strings.Contains(line, want) {
			t.Fatalf("levelSelectLine(0) = %q, missing %q", line, want)
		}
	}
	locked := g.levelSelectLine(1)
	if !strings.Contains(locked, "locked") || !strings.Contains(locked, "sc-") {
		t.Fatalf("levelSelectLine(1) = %q, want locked empty-score state", locked)
	}
}

func TestRedDiamondGateBlocksFinalLevel(t *testing.T) {
	progress := defaultProgress(5)
	progress.UnlockedLevel = 5
	g := &Game{
		levels:   []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"},
		titles:   []string{"One", "Two", "Three", "Four", "Five"},
		parSteps: []int{85, 115, 130, 145, 165},
		progress: progress,
		selected: 4,
		savePath: filepath.Join(t.TempDir(), "progress.json"),
	}
	line := g.levelSelectLine(4)
	if !strings.Contains(line, "red1") {
		t.Fatalf("levelSelectLine(4) = %q, want red1 gate", line)
	}
	if err := g.startSelectedLevel(); err != nil {
		t.Fatal(err)
	}
	if g.world != nil {
		t.Fatal("world loaded despite missing red diamond gate")
	}
	if g.message != "Need 1 red diamonds." {
		t.Fatalf("message = %q, want red gate message", g.message)
	}
}

func TestRedDiamondGateAllowsFinalLevelWhenMet(t *testing.T) {
	progress := defaultProgress(5)
	progress.UnlockedLevel = 5
	progress.RedDiamonds[3] = 1
	g := &Game{
		levels:   []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"},
		titles:   []string{"One", "Two", "Three", "Four", "Five"},
		parSteps: []int{85, 115, 130, 145, 165},
		progress: progress,
		selected: 4,
		savePath: filepath.Join(t.TempDir(), "progress.json"),
	}
	if !strings.Contains(g.levelSelectLine(4), "open") {
		t.Fatalf("levelSelectLine(4) = %q, want open", g.levelSelectLine(4))
	}
	if err := g.startSelectedLevel(); err != nil {
		t.Fatal(err)
	}
	if g.world == nil || g.levelIndex != 4 {
		t.Fatalf("world/level = %v/%d, want level 5 loaded", g.world != nil, g.levelIndex)
	}
}

func TestSecretRouteBypassesFinalRedDiamondGate(t *testing.T) {
	progress := defaultProgress(5)
	progress.UnlockedLevel = 5
	progress.SecretExits[3] = true
	g := &Game{
		levels:   []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"},
		titles:   []string{"One", "Two", "Three", "Four", "Five"},
		parSteps: []int{85, 115, 130, 145, 165},
		progress: progress,
		selected: 4,
		savePath: filepath.Join(t.TempDir(), "progress.json"),
	}
	if got := g.levelMapState(4); got != "secret" {
		t.Fatalf("levelMapState(4) = %q, want secret", got)
	}
	if !strings.Contains(g.levelSelectLine(4), "secret") {
		t.Fatalf("levelSelectLine(4) = %q, want secret state", g.levelSelectLine(4))
	}
	if err := g.startSelectedLevel(); err != nil {
		t.Fatal(err)
	}
	if g.world == nil || g.levelIndex != 4 {
		t.Fatalf("world/level = %v/%d, want level 5 loaded via secret route", g.world != nil, g.levelIndex)
	}
}

func TestWorldMapNodesAndStates(t *testing.T) {
	nodes := levelMapNodes(5)
	if len(nodes) != 5 {
		t.Fatalf("levelMapNodes(5) len = %d, want 5", len(nodes))
	}
	progress := defaultProgress(5)
	progress.UnlockedLevel = 5
	g := &Game{
		levels:   []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"},
		titles:   []string{"One", "Two", "Three", "Four", "Five"},
		parSteps: []int{85, 115, 130, 145, 165},
		progress: progress,
		selected: 4,
	}
	if got := g.levelMapState(4); got != "red1" {
		t.Fatalf("levelMapState(4) = %q, want red1", got)
	}
	g.progress.RedDiamonds[3] = 1
	if got := g.levelMapState(4); got != "open" {
		t.Fatalf("levelMapState(4) = %q, want open", got)
	}
}

func TestNextLevelAfterClearUsesSecretRoute(t *testing.T) {
	target, ok := secretRouteTarget(3)
	if !ok || target != 4 {
		t.Fatalf("secretRouteTarget(3) = %d/%v, want 4/true", target, ok)
	}
	g := &Game{
		levels:     []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"},
		levelIndex: 3,
		world:      &world.World{SecretExitFound: true},
	}
	if got := g.nextLevelAfterClear(); got != 4 {
		t.Fatalf("nextLevelAfterClear() = %d, want secret target 4", got)
	}
	g.world.SecretExitFound = false
	if got := g.nextLevelAfterClear(); got != 4 {
		t.Fatalf("nextLevelAfterClear() = %d, want normal next 4", got)
	}
}

func TestLoadLevelAppliesPersistentTools(t *testing.T) {
	progress := defaultProgress(5)
	progress.HasCompass = true
	progress.HasHammer = true
	progress.HasHook = true
	progress.LifeUpgrades = 2
	g := &Game{
		levels:   []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"},
		parSteps: []int{85, 115, 130, 145, 165},
		progress: progress,
		savePath: filepath.Join(t.TempDir(), "progress.json"),
	}
	if err := g.loadLevel(0); err != nil {
		t.Fatal(err)
	}
	if !g.world.HasCompass || !g.world.HasHammer || !g.world.HasHook {
		t.Fatalf("world tools = compass %v hammer %v hook %v, want all persistent tools", g.world.HasCompass, g.world.HasHammer, g.world.HasHook)
	}
	if g.world.Lives != 5 {
		t.Fatalf("world lives = %d, want 5 from persistent life upgrade", g.world.Lives)
	}
}

func TestRestartLevelKeepsNoRestartPenalty(t *testing.T) {
	progress := defaultProgress(5)
	g := &Game{
		levels:   []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"},
		parSteps: []int{85, 115, 130, 145, 165},
		progress: progress,
		savePath: filepath.Join(t.TempDir(), "progress.json"),
	}
	if err := g.loadLevel(0); err != nil {
		t.Fatal(err)
	}
	if err := g.restartLevel(); err != nil {
		t.Fatal(err)
	}
	if !g.restartUsed {
		t.Fatal("restartUsed = false, want true after retry")
	}
}

func TestTotalRedDiamonds(t *testing.T) {
	progress := Progress{RedDiamonds: []int{1, 0, 2, -1}}
	if got := totalRedDiamonds(progress); got != 3 {
		t.Fatalf("totalRedDiamonds() = %d, want 3", got)
	}
}

func TestBuyHealthUpgrade(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	progress := defaultProgress(1)
	progress.PurpleBank = 8
	g := &Game{
		world:    &world.World{MaxHealth: 3, Health: 1},
		levels:   []string{"one"},
		progress: progress,
		savePath: path,
	}
	if err := g.buyHealthUpgrade(); err != nil {
		t.Fatal(err)
	}
	if g.progress.PurpleBank != 0 || g.progress.MaxHealthUpgrades != 1 {
		t.Fatalf("progress after upgrade = %+v", g.progress)
	}
	if g.world.MaxHealth != 4 || g.world.Health != 4 {
		t.Fatalf("world health = %d/%d, want 4/4", g.world.Health, g.world.MaxHealth)
	}
}

func TestBuyArmorUpgrade(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	progress := defaultProgress(1)
	progress.PurpleBank = 6
	g := &Game{
		world:    &world.World{MaxArmor: 0, Armor: 0},
		levels:   []string{"one"},
		progress: progress,
		savePath: path,
	}
	if err := g.buyArmorUpgrade(); err != nil {
		t.Fatal(err)
	}
	if g.progress.PurpleBank != 0 || g.progress.ArmorUpgrades != 1 {
		t.Fatalf("progress after armor upgrade = %+v", g.progress)
	}
	if g.world.MaxArmor != 1 || g.world.Armor != 1 {
		t.Fatalf("world armor = %d/%d, want 1/1", g.world.Armor, g.world.MaxArmor)
	}
}

func TestBuyLifeUpgrade(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	progress := defaultProgress(1)
	progress.PurpleBank = 10
	g := &Game{
		world:    &world.World{Lives: 1},
		levels:   []string{"one"},
		progress: progress,
		savePath: path,
	}
	if err := g.buyLifeUpgrade(); err != nil {
		t.Fatal(err)
	}
	if g.progress.PurpleBank != 0 || g.progress.LifeUpgrades != 1 {
		t.Fatalf("progress after life upgrade = %+v", g.progress)
	}
	if g.world.Lives != 4 {
		t.Fatalf("world lives = %d, want 4", g.world.Lives)
	}
}

func TestMaxHealthUpgradeCost(t *testing.T) {
	progress := Progress{MaxHealthUpgrades: 2}
	if got := maxHealthForProgress(progress); got != 5 {
		t.Fatalf("maxHealthForProgress() = %d, want 5", got)
	}
	if got := maxHealthUpgradeCost(progress); got != 16 {
		t.Fatalf("maxHealthUpgradeCost() = %d, want 16", got)
	}
	progress = Progress{ArmorUpgrades: 2}
	if got := maxArmorForProgress(progress); got != 2 {
		t.Fatalf("maxArmorForProgress() = %d, want 2", got)
	}
	if got := armorUpgradeCost(progress); got != 12 {
		t.Fatalf("armorUpgradeCost() = %d, want 12", got)
	}
	progress = Progress{LifeUpgrades: 2}
	if got := maxLivesForProgress(progress); got != 5 {
		t.Fatalf("maxLivesForProgress() = %d, want 5", got)
	}
	if got := lifeUpgradeCost(progress); got != 20 {
		t.Fatalf("lifeUpgradeCost() = %d, want 20", got)
	}
}
