package originalgame

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestStartMenuContinueResumesHighestUnlockedMapNode(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.TutorialComplete = true
	g.progress.unlockStage(4)
	g.enterStartMenu(true)
	if g.mode != gameModeStartMenu || g.startMenuSelection != startMenuContinue {
		t.Fatalf("menu mode=%d selection=%d, want start menu/Continue", g.mode, g.startMenuSelection)
	}

	g.updateStartMenu(true, 0)
	if g.mode != gameModeWorldMap || g.stageIndex != 4 || g.worldMapSelectedStage != 4 {
		t.Fatalf("continued mode=%d stage=%d map=%d, want world map Stage 5", g.mode, g.stageIndex, g.worldMapSelectedStage)
	}
}

func TestStartMenuContinueReturnsToIncompleteTutorial(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.enterStartMenu(true)
	g.updateStartMenu(true, 0)
	if g.mode != gameModeStage || g.stageIndex != angkorTutorialStage || !g.rt.IsTutorialStage() {
		t.Fatalf("continued tutorial mode=%d stage=%d tutorial=%v", g.mode, g.stageIndex, g.rt.IsTutorialStage())
	}
}

func TestStartMenuNewGameOnlyClearsProgressAfterYes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	wantOld := newOriginalProgress()
	wantOld.TutorialComplete = true
	wantOld.unlockStage(4)
	wantOld.VioletGemBank = 139
	wantOld.RedDiamondBank = 7
	if err := saveOriginalProgress(path, wantOld); err != nil {
		t.Fatal(err)
	}

	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.enableProgress(path); err != nil {
		t.Fatal(err)
	}
	g.enterStartMenu(true)
	g.updateStartMenu(false, 1)
	g.updateStartMenu(true, 0)
	if !g.startMenuConfirmNew || g.startMenuConfirmChoice != startMenuNo {
		t.Fatalf("confirmation active=%v choice=%d, want active/No", g.startMenuConfirmNew, g.startMenuConfirmChoice)
	}
	gotOld, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotOld, wantOld.normalized()) {
		t.Fatalf("opening confirmation changed progress=%+v, want %+v", gotOld, wantOld.normalized())
	}

	g.updateStartMenu(true, 0)
	if g.startMenuConfirmNew || g.mode != gameModeStartMenu {
		t.Fatalf("No left confirmation=%v mode=%d, want main menu", g.startMenuConfirmNew, g.mode)
	}
	g.updateStartMenu(true, 0)
	g.updateStartMenu(false, 1)
	g.updateStartMenu(true, 0)
	if g.mode != gameModeStage || g.stageIndex != angkorTutorialStage || !g.rt.IsTutorialStage() {
		t.Fatalf("new game mode=%d stage=%d tutorial=%v", g.mode, g.stageIndex, g.rt.IsTutorialStage())
	}
	wantNew := newOriginalProgress()
	if !reflect.DeepEqual(g.progress, wantNew) {
		t.Fatalf("in-memory new progress=%+v, want %+v", g.progress, wantNew)
	}
	gotNew, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotNew, wantNew) {
		t.Fatalf("saved new progress=%+v, want %+v", gotNew, wantNew)
	}
}

func TestStartMenuWithoutSaveStartsNewGameDirectly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	hasProgress, err := originalProgressExists(path)
	if err != nil {
		t.Fatal(err)
	}
	if hasProgress {
		t.Fatal("nonexistent progress was reported as present")
	}

	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.enableProgress(path); err != nil {
		t.Fatal(err)
	}
	g.enterStartMenu(false)
	if g.startMenuSelection != startMenuNewGame {
		t.Fatalf("selection=%d, want New game when no save exists", g.startMenuSelection)
	}
	g.updateStartMenu(true, 0)
	if g.startMenuConfirmNew || g.mode != gameModeStage || g.stageIndex != angkorTutorialStage {
		t.Fatalf("direct new game confirmation=%v mode=%d stage=%d", g.startMenuConfirmNew, g.mode, g.stageIndex)
	}
	hasProgress, err = originalProgressExists(path)
	if err != nil {
		t.Fatal(err)
	}
	if !hasProgress {
		t.Fatal("new game did not create progress file")
	}
}

func TestNewLoadsOriginalTitleAssets(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if g.splashBackground == nil || g.splashLogo == nil || g.splashCopyright == nil {
		t.Fatal("original spl.f title assets were not loaded")
	}
}
