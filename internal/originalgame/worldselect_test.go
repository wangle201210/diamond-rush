package originalgame

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

func TestWorldSelectUsesSourceEightTickArrowMovement(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.enterWorldSelect(-1)
	g.updateWorldSelect(false, 1, 0)
	if g.worldSelectPosition != sealPositionBavaria || g.worldSelectMoveTick != 0 {
		t.Fatalf("right move selected=%d tick=%d, want Bavaria and pending interpolation", g.worldSelectPosition, g.worldSelectMoveTick)
	}
	for tick := 0; tick < sealMoveTicks-1; tick++ {
		g.updateWorldSelect(false, 0, 0)
	}
	if g.worldSelectMoveTick != sealMoveTicks-1 || g.worldSelectArrowX == sealArrowOffsets[sealPositionBavaria][0] {
		t.Fatalf("arrow completed early at tick=%d x=%d", g.worldSelectMoveTick, g.worldSelectArrowX)
	}
	g.updateWorldSelect(false, 0, 0)
	if g.worldSelectMoveTick != sealMoveTicks || g.worldSelectArrowX != 14 || g.worldSelectArrowY != -54 {
		t.Fatalf("arrow after 8 ticks=(%d,%d) tick=%d, want (14,-54)/8", g.worldSelectArrowX, g.worldSelectArrowY, g.worldSelectMoveTick)
	}
}

func TestWorldSelectLocksInputDuringIncomingAngkorRelic(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.RelicMask = 1
	g.enterWorldSelect(sealPositionAngkor)
	g.updateWorldSelect(true, 1, 0)
	if g.mode != gameModeWorldSelect || g.worldSelectPosition != sealPositionAngkor {
		t.Fatalf("input escaped incoming animation: mode=%d position=%d", g.mode, g.worldSelectPosition)
	}
	for tick := 0; tick < 120 && g.worldSelectIncoming >= 0; tick++ {
		g.updateWorldSelect(false, 0, 0)
	}
	if g.worldSelectIncoming != -1 {
		t.Fatal("Angkor relic did not finish the source move/flash/effect sequence")
	}
	target := sealItemOffsets[sealPositionAngkor]
	if g.worldSelectRelicX != original.ScreenWidth/2+target[0] || g.worldSelectRelicY != 136+target[1] {
		t.Fatalf("relic stopped at (%d,%d), want source socket center", g.worldSelectRelicX, g.worldSelectRelicY)
	}
	g.updateWorldSelect(true, 0, 0)
	if g.mode != gameModeWorldMap || g.worldMapSelectedStage != 0 {
		t.Fatalf("Angkor selection mode=%d stage=%d, want world map at current stage", g.mode, g.worldMapSelectedStage)
	}
}

func TestWorldSelectUnlocksAtSourceRedDiamondThresholds(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.RedDiamondBank = sealWorldPrices[sealPositionSiberia]
	g.enterWorldSelect(-1)
	if !g.progress.WorldUnlocked[sealPositionBavaria] || !g.progress.WorldUnlocked[sealPositionSiberia] {
		t.Fatalf("world unlocks=%v, want Bavaria and Siberia", g.progress.WorldUnlocked)
	}
	if g.worldSelectUnlocking != sealPositionSiberia || g.worldSelectPosition != sealPositionSiberia {
		t.Fatalf("unlock animation=%d position=%d, want Siberia", g.worldSelectUnlocking, g.worldSelectPosition)
	}
	for tick := 0; tick < 100 && g.worldSelectUnlocking != 0; tick++ {
		g.updateWorldSelect(false, 0, 0)
	}
	if g.worldSelectUnlocking != 0 {
		t.Fatal("source unlock effect and 15-frame overlay flash did not finish")
	}
}

func TestWorldSelectDrawsOriginalSelectorAssets(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.progress.RelicMask = 1
	g.enterWorldSelect(-1)
	screen := ebiten.NewImage(original.ScreenWidth, original.ScreenHeight)
	g.drawWorldSelect(screen)
	for name, sheet := range map[string]*spriteSheet{
		"central seal": g.tutorialSeal,
		"Angkor relic": g.angkorSeal,
		"arrow":        g.sealArrow,
		"softkeys":     g.softkeys,
		"panel":        g.demoUI,
	} {
		if sheet == nil || sheet.moduleImage == nil {
			t.Errorf("%s source asset is not drawable", name)
		}
	}
}

func TestWorldSelectCancelReturnsToMainMenu(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.enterWorldSelect(-1)
	if err := g.updateSource(sourceInput{Recall: true}); err != nil {
		t.Fatal(err)
	}
	if g.mode != gameModeStartMenu || !g.startMenuHasProgress {
		t.Fatalf("cancel mode/has-progress=%d/%v, want saved-game main menu", g.mode, g.startMenuHasProgress)
	}
}

func TestWorldSelectCannotEnterLockedWorld(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.enterWorldSelect(-1)
	g.worldSelectPosition = sealPositionBavaria
	g.updateWorldSelect(true, 0, 0)
	if g.mode != gameModeWorldSelect || g.worldIndex != sealPositionAngkor {
		t.Fatalf("locked-world activation mode/world=%d/%d, want selector/Angkor", g.mode, g.worldIndex)
	}
}

func TestSealMoveTableMatchesOriginalConfig(t *testing.T) {
	want := [4][4]int{
		{-1, -1, 0, -1},
		{1, -1, 3, -1},
		{2, 2, -1, -1},
		{-1, 0, -1, 2},
	}
	if sealMoveTargets != want {
		t.Fatalf("seal movement table=%v, want source %v", sealMoveTargets, want)
	}
}
