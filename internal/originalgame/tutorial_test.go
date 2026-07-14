package originalgame

import (
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

func TestAngkorTutorialTextIDsUseChineseLocalization(t *testing.T) {
	want := [40]string{
		0:  "我得先看看那个宝箱。",
		1:  "推动石头时，",
		2:  "别堵住自己的路。",
		3:  "你可以让所有物体恢复原位，",
		4:  "只需回到最近的复活点并按 SPACE。",
		5:  "如果你回不到复活点，",
		6:  "而道路又被堵住，",
		7:  "可随时按 ENTER 返回最近的复活点，",
		8:  "但会失去一条命。",
		9:  "这是一种封印吗？",
		10: "啊！封印有反应了！",
		11: "踩上去看看会发生什么……",
		12: "吴哥窟的宏伟神庙……",
		13: "我终于进来了！",
		14: "出发吧！",
		15: "看，前面有一把魔法锁！",
		16: "收集指定数量的钻石即可开锁。",
		17: "糟了，这扇门上锁了！",
		18: "钥匙一定就在附近……",
		19: "找到罗盘了！它会帮我找到出口。",
		20: "找到神秘锤了！",
		21: "按 SPACE 使用。",
		22: "太好了！现在能砸碎那些脆弱的墙了。",
		23: "找到神秘钩索了！",
		24: "隔着一段距离按 SPACE，可将物体拉向自己。",
		25: "有意思……也许可以用这个……",
		26: "找到神秘药水了！",
		27: "现在可以在水下呼吸了。",
		28: "找到冰冻锤了！",
		29: "按 SPACE 冻结物体。",
		30: "把东西冻住？",
		31: "试试看吧！",
		32: "吴哥窟的最后一间密室！据说火焰水晶就藏在这里……",
		33: "但我有种不祥的预感……",
		34: "银色钻石就在这里……我敢肯定！",
		35: "嗯……附近有一股黑暗力量……",
		36: "大干一场吧！",
		37: "寒冰钻石一定就在西伯利亚的最后一间密室里！",
		38: "成败在此一举！",
		39: "完成本隐藏关需要神秘药水。请先在巴伐利亚第8关取得，再返回这里。",
	}
	for index, text := range want {
		if got := tutorialText(index); got != text {
			t.Errorf("tutorial text %d=%q, want %q", index, got, text)
		}
	}
}

func TestDesktopControlLabelsMatchKeyboardBindings(t *testing.T) {
	if desktopActionKeyLabel != "SPACE" || desktopRecallKeyLabel != "ENTER" || desktopNavigationKeyLabel != "TAB" || desktopSkipKeyLabel != "S" {
		t.Fatalf("desktop labels action=%q recall=%q navigation=%q skip=%q", desktopActionKeyLabel, desktopRecallKeyLabel, desktopNavigationKeyLabel, desktopSkipKeyLabel)
	}
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, label := range []string{desktopActionKeyLabel, desktopRecallKeyLabel, desktopNavigationKeyLabel, desktopSkipKeyLabel} {
		if !g.fontSmall.supports(label) {
			t.Errorf("desktop control label is not renderable: %q", label)
		}
	}
}

func TestBavariaSecretStageEntryHintFollowsIntroAndPotionState(t *testing.T) {
	g, err := New("decoded/world1")
	if err != nil {
		t.Fatal(err)
	}
	g.loadStage(bavariaFirstSecretStage)
	if !g.stageEntryHintPending || g.rt.TutorialScriptActive {
		t.Fatalf("entry hint pending/active=%v/%v, want pending until title completes", g.stageEntryHintPending, g.rt.TutorialScriptActive)
	}
	for tick := 0; tick < stageIntroDuration-1; tick++ {
		if err := g.updateSource(sourceInput{}); err != nil {
			t.Fatal(err)
		}
	}
	if g.rt.TutorialScriptActive {
		t.Fatal("entry hint replaced the stage title before its source duration elapsed")
	}
	if err := g.updateSource(sourceInput{}); err != nil {
		t.Fatal(err)
	}
	prompt, ok := g.rt.TutorialPrompt()
	if g.stageEntryHintPending || !g.rt.TutorialScriptActive || !ok || prompt.TextIndex != 39 {
		t.Fatalf("entry hint pending/active/prompt=%v/%v/%+v,%v", g.stageEntryHintPending, g.rt.TutorialScriptActive, prompt, ok)
	}

	g.progress.WaterBreathingPotion = true
	g.loadStage(bavariaFirstSecretStage)
	if g.stageEntryHintPending || g.rt.TutorialScriptActive {
		t.Fatal("entry hint was scheduled after the campaign already owned the potion")
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
