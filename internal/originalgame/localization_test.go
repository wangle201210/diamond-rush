package originalgame

import (
	"image/color"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

func TestSimplifiedChineseCatalogIsCompleteAndRenderable(t *testing.T) {
	g, err := New(defaultWorldDir)
	if err != nil {
		t.Fatal(err)
	}
	if g.fontSmall.cjkFace == nil || g.fontMedium.cjkFace == nil {
		t.Fatal("system Chinese font was not attached to both source font sizes")
	}
	for _, key := range allUITextKeys {
		value, ok := simplifiedChinese[key]
		if !ok || strings.TrimSpace(value) == "" {
			t.Errorf("missing Chinese UI text %q", key)
			continue
		}
		if !g.fontSmall.supports(value) {
			t.Errorf("Chinese UI text is not renderable: %q", value)
		}
	}
	for index, value := range angkorTutorialTexts {
		if strings.TrimSpace(value) == "" {
			t.Errorf("missing Chinese tutorial text %d", index)
			continue
		}
		if !g.fontSmall.supports(value) {
			t.Errorf("tutorial text %d is not renderable: %q", index, value)
		}
	}

	screen := ebiten.NewImage(original.ScreenWidth, original.ScreenHeight)
	g.fontSmall.beginFrame()
	g.fontSmall.drawText(screen, "按 SPACE 使用。", 8, 24, false, color.White)
	if len(g.fontSmall.cjkDraws) != 1 || g.fontSmall.cjkDraws[0].text != "按 SPACE 使用。" {
		t.Fatalf("mixed Chinese/ASCII text queued %+v, want one system-font run", g.fontSmall.cjkDraws)
	}
	if bounds := screen.Bounds(); bounds.Dx() != original.ScreenWidth || bounds.Dy() != original.ScreenHeight {
		t.Fatalf("localized render target=%v", bounds)
	}
}

func TestWindowTitleUsesCurrentWorldAndPlayerFacingNumber(t *testing.T) {
	if got := windowTitleForWorld(original.WorldAngkor); got != "钻石狂潮原作运行版 - 吴哥窟（世界1）" {
		t.Fatalf("Angkor window title=%q", got)
	}
	if got := windowTitleForWorld(original.WorldBavaria); got != "钻石狂潮原作运行版 - 巴伐利亚（世界2）" {
		t.Fatalf("Bavaria window title=%q", got)
	}
}

func TestSystemChineseFontDiscoveryLoadsUsableFace(t *testing.T) {
	font, err := loadSystemCJKFont()
	if err != nil {
		t.Fatal(err)
	}
	if font.source == nil || font.path == "" || font.family == "" {
		t.Fatalf("system Chinese font=%+v", font)
	}
}
