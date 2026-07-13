package originalgame

import (
	"image/color"
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

const angkorTutorialStage = 13

func (g *Game) drawTutorialPrompt(screen *ebiten.Image) {
	if g == nil || g.rt == nil || g.fontSmall == nil {
		return
	}
	prompt, ok := g.rt.TutorialPrompt()
	if !ok {
		return
	}
	maxWidth := 222
	if prompt.Placement == original.TutorialTextBottom {
		maxWidth = 220
	}
	lines := wrapTutorialText(g.fontSmall, tutorialText(prompt.TextIndex), maxWidth)
	if len(lines) == 0 {
		return
	}

	g.drawTutorialChrome(screen)
	if prompt.Placement == original.TutorialTextBottom {
		g.drawTutorialBottomPrompt(screen, lines)
		return
	}
	g.drawTutorialBubblePrompt(screen, prompt, lines)
}

func (g *Game) drawTutorialBubblePrompt(screen *ebiten.Image, prompt original.TutorialPrompt, lines []string) {
	lineHeight := g.fontSmall.lineHeight()
	panelY := prompt.Y + 4
	panelHeight := len(lines)*lineHeight + 4
	fill := color.RGBA{0x00, 0x00, 0x49, 0xff}
	g.drawDemoPanelWithSheet(screen, g.demoUI, prompt.X, panelY, 226, panelHeight, fill)
	for index, line := range lines {
		g.fontSmall.drawText(screen, line, prompt.X+2, prompt.Y+16+index*lineHeight, false, color.White)
	}
	g.drawTutorialPromptIndicator(screen, prompt.X+216, prompt.Y+len(lines)*lineHeight+4)
}

func (g *Game) drawTutorialBottomPrompt(screen *ebiten.Image, lines []string) {
	lineHeight := g.fontSmall.lineHeight()
	fill := color.RGBA{0x00, 0x00, 0x49, 0xff}
	bodyHeight := len(lines)*lineHeight + 8
	g.drawDemoPanelWithSheet(screen, g.demoUIBlue, 6, 212, 226, bodyHeight, fill)

	label := tr(textHint)
	labelWidth := g.fontSmall.stringWidth(label) + 10
	g.drawDemoPanelWithSheet(screen, g.demoUIBlue, 16, 193, labelWidth, 16, fill)
	drawRect(screen, 13, 195, labelWidth+6, 3, fill)
	g.fontSmall.drawText(screen, label, 19, 200, false, color.White)
	for index, line := range lines {
		g.fontSmall.drawText(screen, line, 8, 224+index*lineHeight, false, color.White)
	}
	g.drawTutorialPromptIndicator(screen, 223, 197)
}

func (g *Game) drawTutorialPromptIndicator(screen *ebiten.Image, centerX, y int) {
	if g.fontSmall == nil || (g.tick/2)%4 >= 2 {
		return
	}
	drawControlKeycap(screen, g.fontSmall, desktopActionKeyLabel, centerX, y)
}

func (g *Game) drawTutorialChrome(screen *ebiten.Image) {
	if g == nil || g.rt == nil || g.fontSmall == nil || !g.rt.TutorialScriptActive {
		return
	}
	drawRect(screen, 0, 0, original.ScreenWidth, 42, color.Black)
	drawRect(screen, 0, 278, original.ScreenWidth, 42, color.Black)
	g.drawTutorialPortrait(screen)
	g.fontSmall.drawText(screen, tr(textPromptSkip, desktopSkipKeyLabel), 2, 318, false, color.White)
}

func (g *Game) drawTutorialPortrait(screen *ebiten.Image) {
	if g == nil || g.rt == nil || !g.rt.TutorialScriptActive {
		return
	}
	x := g.rt.TutorialPortraitX
	y := g.rt.TutorialPortraitY
	if reveal := min(5, g.rt.TutorialPortraitRevealTicks); reveal > 0 {
		camX, camY := g.cameraPixels()
		playerX := g.rt.Player.X*original.TileSize - camX
		playerY := g.rt.Player.Y*original.TileSize - camY
		revealX := (playerX*(5-reveal) + x*reveal) / 5
		revealY := (playerY*(5-reveal) + y*reveal) / 5
		drawRect(screen, revealX, revealY, reveal*102/5, reveal*38/5, color.White)
		return
	}
	if !g.rt.TutorialPortraitVisible || g.tutorialFaces == nil || g.tutorialPortrait == nil {
		return
	}
	drawRect(screen, x-3, y-3, 109, 45, color.Black)
	g.tutorialPortrait.drawAnimation(screen, 0, g.tick, x, y, 0)
	g.tutorialFaces.drawFrame(screen, g.rt.TutorialPortraitFace, x, y, 0)
	if g.rt.TutorialPortraitMark >= 0 && g.tutorialMarks != nil {
		g.tutorialMarks.drawFrame(screen, g.rt.TutorialPortraitMark, x+90, y-6, 0)
	}
}

func (g *Game) drawTutorialFlash(screen *ebiten.Image) {
	if g != nil && g.rt != nil && g.rt.TutorialFlashVisible {
		screen.Fill(color.White)
	}
}

func wrapTutorialText(font *bitmapFont, text string, maxWidth int) []string {
	if strings.TrimSpace(text) == "" || font == nil || maxWidth <= 0 {
		return nil
	}
	lines := make([]string, 0, 3)
	line := ""
	pendingSpace := false
	flush := func() {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
		line = ""
		pendingSpace = false
	}
	for _, token := range tutorialWrapTokens(text) {
		if token == "\n" {
			flush()
			continue
		}
		if token == " " {
			pendingSpace = line != ""
			continue
		}
		candidate := line
		if pendingSpace && candidate != "" {
			candidate += " "
		}
		candidate += token
		if line == "" || font.stringWidth(candidate) <= maxWidth {
			line = candidate
			pendingSpace = false
			continue
		}
		flush()
		line = token
	}
	flush()
	return lines
}

func tutorialWrapTokens(text string) []string {
	tokens := make([]string, 0, len([]rune(text)))
	var ascii strings.Builder
	flushASCII := func() {
		if ascii.Len() == 0 {
			return
		}
		tokens = append(tokens, ascii.String())
		ascii.Reset()
	}
	for _, char := range text {
		switch {
		case char == '\n':
			flushASCII()
			tokens = append(tokens, "\n")
		case unicode.IsSpace(char):
			flushASCII()
			if len(tokens) == 0 || tokens[len(tokens)-1] != " " {
				tokens = append(tokens, " ")
			}
		case char <= unicode.MaxASCII:
			ascii.WriteRune(char)
		default:
			flushASCII()
			tokens = append(tokens, string(char))
		}
	}
	flushASCII()
	return tokens
}

func (g *Game) drawTutorialSealCell(dst *ebiten.Image, x, y, px, py int) {
	if g == nil || g.rt == nil || !g.rt.IsTutorialStage() || g.tutorialSeal == nil {
		return
	}
	if x < 60 || x >= 65 || y < 2 || y >= 7 {
		return
	}
	frame := 4 + (y-2)*5 + x - 60
	g.tutorialSeal.drawFrame(dst, frame, px, py, 0)
}

func (g *Game) drawTutorialSealOverlay(dst *ebiten.Image, camX, camY int) {
	if !g.tutorialSealOverlayVisible() {
		return
	}
	g.tutorialSeal.drawFrame(dst, 1, 60*original.TileSize-camX, 2*original.TileSize-camY, 0)
}

func (g *Game) tutorialSealOverlayVisible() bool {
	return g != nil && g.rt != nil && g.rt.IsTutorialStage() && g.tutorialSeal != nil &&
		(g.rt.Player == (original.Point{X: 60, Y: 3}) || g.rt.Player == (original.Point{X: 61, Y: 3}))
}

func (g *Game) drawTutorialRecallHint(dst *ebiten.Image, camX, camY int) {
	sequence, centerX, top, ok := g.tutorialRecallHintRenderState(camX, camY)
	if !ok {
		return
	}
	drawControlKeycap(dst, g.fontSmall, desktopRecallKeyLabel, centerX, top+(sequence&1))
}

func (g *Game) tutorialRecallHintRenderState(camX, camY int) (sequence, centerX, top int, ok bool) {
	if g == nil || g.rt == nil || !g.rt.TutorialRecallHintVisible || g.rt.TutorialScriptActive || g.fontSmall == nil || g.tutorialRecallHint == nil {
		return 0, 0, 0, false
	}
	playerX, playerY := g.renderedPlayerPixels()
	return (g.tick >> 1) & 0x3, playerX - camX + original.TileSize/2, playerY - camY - original.TileSize, true
}

func (g *Game) finishTutorial() {
	if g == nil || g.rt == nil || !g.rt.TutorialComplete {
		return
	}
	g.progress.TutorialComplete = true
	if g.progressPath != "" {
		if err := saveOriginalProgress(g.progressPath, g.progress); err != nil {
			g.message = err.Error()
			return
		}
	}
	g.loadStage(0)
	g.mode = gameModeStage
	g.message = "Angkor tutorial complete"
}
