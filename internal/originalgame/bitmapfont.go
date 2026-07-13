package originalgame

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

const sourceFontYOffset = 12

const localizedPanelOpticalYOffset = 2

// Localized text is rasterized below native resolution, then gently scaled up
// on the final screen so it sits naturally beside the source bitmap font.
const localizedTextRasterScale = 0.72

type bitmapFontMetadata struct {
	First      int   `json:"first"`
	Last       int   `json:"last"`
	Columns    int   `json:"columns"`
	Padding    int   `json:"padding"`
	CellWidth  int   `json:"cellWidth"`
	CellHeight int   `json:"cellHeight"`
	FontHeight int   `json:"fontHeight"`
	Ascent     int   `json:"ascent"`
	Widths     []int `json:"widths"`
}

type bitmapFont struct {
	image         *ebiten.Image
	meta          bitmapFontMetadata
	cjkSource     *opentype.Font
	cjkPointSize  float64
	cjkFace       *rasterCJKFace
	cjkFinalFaces map[int]*rasterCJKFace
	cjkDraws      []cjkDrawCommand
}

type cjkDrawCommand struct {
	text     string
	x        float64
	y        int
	centered bool
	tint     color.Color
}

type rasterCJKFace struct {
	source      *opentype.Font
	face        xfont.Face
	metrics     xfont.Metrics
	mu          sync.Mutex
	glyphBuffer sfnt.Buffer
	cache       map[string]rasterizedCJKText
}

type rasterizedCJKText struct {
	image      *ebiten.Image
	lineHeight int
	advance    float64
}

func loadBitmapFont(imagePath, metadataPath string) (*bitmapFont, error) {
	metadataFile, err := os.Open(filepath.Clean(resolvePath(metadataPath)))
	if err != nil {
		return nil, err
	}
	defer metadataFile.Close()
	var metadata bitmapFontMetadata
	if err := json.NewDecoder(metadataFile).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decode bitmap font metadata: %w", err)
	}
	if metadata.First < 0 || metadata.Last < metadata.First || metadata.Columns <= 0 || metadata.CellWidth <= 0 || metadata.CellHeight <= 0 || len(metadata.Widths) != metadata.Last-metadata.First+1 {
		return nil, fmt.Errorf("invalid bitmap font metadata: %s", metadataPath)
	}

	imageFile, err := os.Open(filepath.Clean(resolvePath(imagePath)))
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()
	atlas, err := png.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("decode bitmap font atlas: %w", err)
	}
	return &bitmapFont{image: ebiten.NewImageFromImage(atlas), meta: metadata}, nil
}

func (f *bitmapFont) useCJKSource(source *opentype.Font) error {
	if f == nil || source == nil {
		return fmt.Errorf("nil Chinese font source")
	}
	pointSize := float64(f.meta.FontHeight + 2)
	face, err := newRasterCJKFace(source, pointSize)
	if err != nil {
		return err
	}
	f.cjkSource = source
	f.cjkPointSize = pointSize
	f.cjkFace = face
	f.cjkFinalFaces = make(map[int]*rasterCJKFace)
	return nil
}

func newRasterCJKFace(source *opentype.Font, pointSize float64) (*rasterCJKFace, error) {
	face, err := opentype.NewFace(source, &opentype.FaceOptions{
		Size:    pointSize,
		DPI:     72,
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	return &rasterCJKFace{
		source:  source,
		face:    face,
		metrics: face.Metrics(),
		cache:   make(map[string]rasterizedCJKText),
	}, nil
}

func (f *bitmapFont) lineHeight() int {
	if f == nil {
		return 0
	}
	height := f.meta.FontHeight
	if f.cjkFace != nil {
		height = max(height, f.cjkFace.metrics.Height.Ceil())
	}
	return height
}

func (f *bitmapFont) bitmapSupports(char rune) bool {
	return f != nil && int(char) >= f.meta.First && int(char) <= f.meta.Last
}

func (f *bitmapFont) supports(text string) bool {
	if f == nil {
		return false
	}
	for _, char := range text {
		if f.bitmapSupports(char) {
			continue
		}
		if f.cjkFace == nil || char < 0x80 || !f.cjkFace.hasGlyph(char) {
			return false
		}
	}
	return true
}

func (f *bitmapFont) stringWidth(text string) int {
	if f == nil {
		return 0
	}
	return int(math.Ceil(f.stringAdvance(text)))
}

func (f *bitmapFont) stringAdvance(text string) float64 {
	if f.usesCJK(text) && f.cjkFace != nil {
		return f.cjkFace.advance(text)
	}
	width := 0.0
	f.forEachRun(text, func(run string, bitmap bool) {
		if bitmap {
			for _, char := range run {
				width += float64(f.meta.Widths[int(char)-f.meta.First])
			}
			return
		}
		if f.cjkFace != nil {
			width += f.cjkFace.advance(run)
		}
	})
	return width
}

// drawText mirrors h.drawTextWithFlags and PlatformGraphics.drawString for
// the two FreeJ2ME fonts used by Stage 1.
func (f *bitmapFont) drawText(dst *ebiten.Image, text string, x, y int, centered bool, tint color.Color) {
	if f == nil || f.image == nil {
		return
	}
	if f.usesCJK(text) && f.cjkFace != nil {
		f.cjkDraws = append(f.cjkDraws, cjkDrawCommand{
			text:     text,
			x:        float64(x),
			y:        y,
			centered: centered,
			tint:     tint,
		})
		return
	}
	cursor := float64(x)
	if centered {
		cursor -= f.stringAdvance(text) / 2
	}
	f.forEachRun(text, func(run string, bitmap bool) {
		cursor = math.Round(cursor)
		if bitmap {
			cursor = f.drawBitmapRun(dst, run, cursor, y, tint)
			return
		}
		if f.cjkFace == nil {
			return
		}
		f.cjkDraws = append(f.cjkDraws, cjkDrawCommand{text: run, x: cursor, y: y, tint: tint})
		cursor += f.cjkFace.advance(run)
	})
}

func (f *bitmapFont) usesCJK(text string) bool {
	for _, char := range text {
		if !f.bitmapSupports(char) {
			return true
		}
	}
	return false
}

func (f *bitmapFont) beginFrame() {
	if f != nil {
		f.cjkDraws = f.cjkDraws[:0]
	}
}

func (f *bitmapFont) drawFinalCJK(screen ebiten.FinalScreen, geoM ebiten.GeoM) {
	if f == nil || f.cjkFace == nil || len(f.cjkDraws) == 0 {
		return
	}
	scale := math.Hypot(geoM.Element(0, 0), geoM.Element(0, 1))
	if scale <= 0 {
		return
	}
	face := f.finalCJKFace(scale * localizedTextRasterScale)
	if face == nil {
		return
	}
	displayScale := 1 / localizedTextRasterScale
	logicalLineHeight := float64(f.lineHeight())
	for _, command := range f.cjkDraws {
		rasterized := face.rasterize(command.text)
		x, boxTop := geoM.Apply(command.x, float64(command.y-sourceFontYOffset))
		_, boxBottom := geoM.Apply(command.x, float64(command.y-sourceFontYOffset)+logicalLineHeight)
		if command.centered {
			x -= rasterized.advance * displayScale / 2
		}
		displayHeight := float64(rasterized.lineHeight) * displayScale
		top := boxTop + (boxBottom-boxTop-displayHeight)/2 - displayScale
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleWithColor(command.tint)
		op.Filter = ebiten.FilterLinear
		op.GeoM.Scale(displayScale, displayScale)
		op.GeoM.Translate(math.Round(x)-displayScale, math.Round(top))
		screen.DrawImage(rasterized.image, op)
	}
}

func (f *bitmapFont) finalCJKFace(scale float64) *rasterCJKFace {
	key := max(1, int(math.Round(f.cjkPointSize*scale*8)))
	if face := f.cjkFinalFaces[key]; face != nil {
		return face
	}
	if len(f.cjkFinalFaces) >= 8 {
		clear(f.cjkFinalFaces)
	}
	face, err := newRasterCJKFace(f.cjkSource, float64(key)/8)
	if err != nil {
		return nil
	}
	f.cjkFinalFaces[key] = face
	return face
}

func (f *rasterCJKFace) hasGlyph(char rune) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	glyph, err := f.source.GlyphIndex(&f.glyphBuffer, char)
	return err == nil && glyph != 0
}

func (f *rasterCJKFace) advance(text string) float64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return fixedToFloat(xfont.MeasureString(f.face, text))
}

func (f *rasterCJKFace) rasterize(text string) rasterizedCJKText {
	f.mu.Lock()
	defer f.mu.Unlock()
	if cached, ok := f.cache[text]; ok {
		return cached
	}

	advance := xfont.MeasureString(f.face, text)
	lineHeight := max(1, f.metrics.Height.Ceil())
	width := max(1, advance.Ceil()+2)
	alpha := image.NewAlpha(image.Rect(0, 0, width, lineHeight+2))
	drawer := xfont.Drawer{
		Dst:  alpha,
		Src:  image.NewUniform(color.Alpha{A: 0xff}),
		Face: f.face,
		Dot:  fixed.P(1, 1+f.metrics.Ascent.Ceil()),
	}
	drawer.DrawString(text)

	rasterized := rasterizedCJKText{
		image:      ebiten.NewImageFromImage(alpha),
		lineHeight: lineHeight,
		advance:    fixedToFloat(advance),
	}
	f.cache[text] = rasterized
	return rasterized
}

func fixedToFloat(value fixed.Int26_6) float64 {
	return float64(value) / 64
}

func (f *bitmapFont) drawBitmapRun(dst *ebiten.Image, text string, cursor float64, y int, tint color.Color) float64 {
	strideX := f.meta.CellWidth + f.meta.Padding*2
	strideY := f.meta.CellHeight + f.meta.Padding*2
	for _, char := range text {
		index := int(char) - f.meta.First
		if index < 0 || index >= len(f.meta.Widths) {
			continue
		}
		srcX := (index % f.meta.Columns) * strideX
		srcY := (index / f.meta.Columns) * strideY
		src := f.image.SubImage(image.Rect(srcX, srcY, srcX+strideX, srcY+strideY)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleWithColor(tint)
		op.GeoM.Translate(cursor-float64(f.meta.Padding), float64(y-sourceFontYOffset-f.meta.Padding))
		dst.DrawImage(src, op)
		cursor += float64(f.meta.Widths[index])
	}
	return cursor
}

func (f *bitmapFont) forEachRun(text string, visit func(run string, bitmap bool)) {
	if text == "" {
		return
	}
	start := 0
	bitmap := false
	initialized := false
	for index, char := range text {
		charBitmap := f.bitmapSupports(char)
		if !initialized {
			bitmap = charBitmap
			initialized = true
			continue
		}
		if charBitmap == bitmap {
			continue
		}
		visit(text[start:index], bitmap)
		start = index
		bitmap = charBitmap
	}
	visit(text[start:], bitmap)
}

func drawSourcePanelLabel(dst *ebiten.Image, font *bitmapFont, text string, centerX, panelY int) {
	if font == nil {
		return
	}
	width := font.stringWidth(text)
	x := centerX - width/2 - 5
	y := panelY - 5
	w := width + 10
	h := font.lineHeight() + 10
	drawRoundedRect(dst, x, y, w, h, 5, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	drawRoundedRect(dst, x+1, y+1, w-2, h-2, 4, color.RGBA{0x0c, 0x2f, 0x39, 0xff})
	font.drawText(dst, text, centerX, panelY+sourceFontYOffset+localizedPanelOpticalYOffset, true, color.White)
}

func drawSourcePanelLines(dst *ebiten.Image, font *bitmapFont, lines []string, centerX, panelY int) {
	if font == nil || len(lines) == 0 {
		return
	}
	width := 0
	for _, line := range lines {
		width = max(width, font.stringWidth(line))
	}
	lineHeight := font.lineHeight() + 2
	x := centerX - width/2 - 5
	y := panelY - 5
	w := width + 10
	h := lineHeight*len(lines) + 8
	drawRoundedRect(dst, x, y, w, h, 5, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	drawRoundedRect(dst, x+1, y+1, w-2, h-2, 4, color.RGBA{0x0c, 0x2f, 0x39, 0xff})
	for index, line := range lines {
		font.drawText(dst, line, centerX, panelY+sourceFontYOffset+localizedPanelOpticalYOffset+index*lineHeight, true, color.White)
	}
}

func drawControlKeycap(dst *ebiten.Image, font *bitmapFont, label string, centerX, top int) {
	if dst == nil || font == nil || label == "" {
		return
	}
	const height = 14
	width := font.stringWidth(label) + 8
	bounds := dst.Bounds()
	x := clamp(centerX-width/2, bounds.Min.X+1, bounds.Max.X-width-1)
	y := clamp(top, bounds.Min.Y+1, bounds.Max.Y-height-1)
	drawRect(dst, x, y, width, height, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	drawRect(dst, x+1, y+1, width-2, height-2, color.RGBA{0x0c, 0x2f, 0x39, 0xff})
	font.drawText(dst, label, x+width/2, y+sourceFontYOffset, true, color.White)
}

func drawRoundedRect(dst *ebiten.Image, x, y, width, height, radius int, fill color.Color) {
	if width <= 0 || height <= 0 || radius <= 0 {
		return
	}
	radius = min(radius, min(width, height)/2)
	vector.DrawFilledRect(dst, float32(x+radius), float32(y), float32(width-radius*2), float32(height), fill, false)
	vector.DrawFilledRect(dst, float32(x), float32(y+radius), float32(width), float32(height-radius*2), fill, false)
	for _, center := range [][2]int{{x + radius, y + radius}, {x + width - radius, y + radius}, {x + radius, y + height - radius}, {x + width - radius, y + height - radius}} {
		vector.DrawFilledCircle(dst, float32(center[0]), float32(center[1]), float32(radius), fill, false)
	}
}
