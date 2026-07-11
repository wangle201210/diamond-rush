package originalgame

import (
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

func TestAngkorTutorialTextMappingMatchesSource(t *testing.T) {
	want := map[int]string{
		0:  "I should check that chest first.",
		4:  "by going back to the last circle and pressing SPACE.",
		7:  "you can press ENTER at any time to go back to the last circle",
		8:  "but it will cost you a life.",
		9:  "Is this a kind of seal?",
		12: "The great temple of Angkor Wat...",
		19: "You found a compass! It will help you find your way out.",
		20: "You found the mystic mallet!",
		21: "Press SPACE to use it.",
		22: "Great! Now I can crush those weak walls.",
		32: "The final chamber in Angkor Wat! The fire crystal is supposed to be hidden here...",
		33: "But I have a bad feeling about this...",
	}
	for index, text := range want {
		if got := tutorialText(index); got != text {
			t.Errorf("tutorial text %d=%q, want %q", index, got, text)
		}
	}
}

func TestDesktopControlLabelsMatchKeyboardBindings(t *testing.T) {
	if desktopActionKeyLabel != "SPACE" || desktopRecallKeyLabel != "ENTER" || desktopSkipKeyLabel != "S" {
		t.Fatalf("desktop labels action=%q recall=%q skip=%q", desktopActionKeyLabel, desktopRecallKeyLabel, desktopSkipKeyLabel)
	}
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, label := range []string{desktopActionKeyLabel, desktopRecallKeyLabel, desktopSkipKeyLabel} {
		if !g.fontSmall.supports(label) {
			t.Errorf("desktop control label is not renderable: %q", label)
		}
	}
}

func TestAngkorTutorialPromptsFitSourcePanel(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for index, text := range angkorTutorialTexts {
		if text == "" {
			continue
		}
		if !g.fontSmall.supports(text) {
			t.Errorf("tutorial text %d contains unsupported glyphs: %q", index, text)
		}
		lines := wrapTutorialText(g.fontSmall, text, 214)
		if len(lines) == 0 || len(lines) > 3 {
			t.Errorf("tutorial text %d wraps to %d lines, want 1..3", index, len(lines))
		}
		for _, line := range lines {
			if width := g.fontSmall.stringWidth(line); width > 214 {
				t.Errorf("tutorial text %d line width=%d: %q", index, width, line)
			}
		}
	}
}

func TestAngkorTutorialUsesSourceSealFramesAndDrawsPrompt(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.loadStage(angkorTutorialStage)
	if g.tutorialSeal == nil {
		t.Fatal("tutorial seal sprite is nil")
	}
	if g.tutorialRecallHint == nil || len(g.tutorialRecallHint.meta.AnimationCounts) != 3 || g.tutorialRecallHint.meta.AnimationCounts[0] != 4 {
		t.Fatal("gen3.f chunk 0 recall-key cue does not expose its source animation cadence")
	}
	if g.demoUI == nil || g.demoUIBlue == nil {
		t.Fatal("demoui.f prompt borders are not loaded")
	}
	if len(g.tutorialSeal.meta.FrameCounts) != 29 {
		t.Fatalf("tutorial seal frames=%d, want 29", len(g.tutorialSeal.meta.FrameCounts))
	}
	if g.tutorialFaces == nil || len(g.tutorialFaces.meta.FrameCounts) != 9 || g.tutorialMarks == nil || len(g.tutorialMarks.meta.FrameCounts) != 9 || g.tutorialPortrait == nil || len(g.tutorialPortrait.meta.FrameCounts) != 2 {
		t.Fatal("tutorial portrait source sprites do not expose 9 face, 9 punctuation, and 2 base frames")
	}
	g.rt.TutorialScriptActive = true
	g.rt.TutorialScriptID = 29
	g.rt.TutorialTextIndex = 12
	g.rt.TutorialTextPlacement = original.TutorialTextBubble
	g.rt.TutorialTextY = 90
	g.rt.TutorialTextSide = 2
	screen := ebiten.NewImage(original.ScreenWidth, original.ScreenHeight)
	g.Draw(screen)
}

func TestTutorialRecallHintUsesDesktopEnterKeycapAtRenderedHeroPosition(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.loadStage(angkorTutorialStage)
	g.rt.TutorialRecallHintVisible = true
	g.rt.Player = original.Point{X: 47, Y: 7}
	g.rt.PlayerMotion = original.ObjectMotion{DX: 1, Remaining: 12}
	g.lastDX = 1
	g.lastDY = 0
	g.heroMoveOffset = 12
	g.tick = 7
	sequence, px, py, ok := g.tutorialRecallHintRenderState(10, 20)
	if !ok || sequence != 3 || px != 47*original.TileSize-12-10+original.TileSize/2 || py != 7*original.TileSize-20-original.TileSize {
		t.Fatalf("recall hint ok=%v sequence=%d position=%d,%d", ok, sequence, px, py)
	}
	g.rt.TutorialScriptActive = true
	if _, _, _, ok := g.tutorialRecallHintRenderState(10, 20); ok {
		t.Fatal("recall hint remained visible while a source demo script was active")
	}
}

func TestTutorialSealKeepsOverlayAndStartsSourceEffectAtOffsetSix(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	g.loadStage(angkorTutorialStage)
	g.rt.Player = original.Point{X: 61, Y: 3}
	g.rt.PlayerMotion = original.ObjectMotion{DX: 1, Remaining: 6}
	g.tickWorld()
	if !g.rt.TutorialSealActivated || !g.tutorialSealOverlayVisible() {
		t.Fatalf("seal active=%v overlay=%v, want hidden hero with visible source overlay", g.rt.TutorialSealActivated, g.tutorialSealOverlayVisible())
	}
	if len(g.worldEffects) != 1 || g.worldEffects[0].Animation != 5 || g.worldEffects[0].Point != (original.Point{X: 61, Y: 3}) {
		t.Fatalf("seal effects=%+v, want source animation 5 at (61,3)", g.worldEffects)
	}
	g.tickWorld()
	if len(g.worldEffects) != 1 || g.worldEffects[0].Sequence != 1 {
		t.Fatalf("seal effect after second tick=%+v, want one advancing effect", g.worldEffects)
	}
}

func TestAngkorTutorialCompletionPersistsAndLoadsStageOne(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original-progress.json")
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.enableProgress(path); err != nil {
		t.Fatal(err)
	}
	g.loadStage(angkorTutorialStage)
	g.rt.TutorialComplete = true
	g.finishTutorial()
	if g.stageIndex != 0 || g.rt.Stage.Index != 0 || !g.progress.TutorialComplete {
		t.Fatalf("tutorial transition stage=%d source=%d complete=%v", g.stageIndex, g.rt.Stage.Index, g.progress.TutorialComplete)
	}
	loaded, err := loadOriginalProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.TutorialComplete {
		t.Fatal("tutorial completion was not persisted")
	}
}

func TestAngkorSecretStagesUseSourceCampaignPrerequisites(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for stage := angkorFirstSecretStage; stage < angkorTutorialStage; stage++ {
		g.loadStage(stage)
		if g.rt.MaxHealth != 8 || g.rt.Health != 8 || g.rt.SpecialItemMask != 8 {
			t.Errorf("secret stage %d health=%d/%d tool=%d, want 8/8/8", stage, g.rt.Health, g.rt.MaxHealth, g.rt.SpecialItemMask)
		}
	}
}
