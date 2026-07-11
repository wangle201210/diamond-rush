package originalgame

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const sourceFontYOffset = 12

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
	image *ebiten.Image
	meta  bitmapFontMetadata
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

func (f *bitmapFont) supports(text string) bool {
	if f == nil {
		return false
	}
	for _, char := range text {
		if int(char) < f.meta.First || int(char) > f.meta.Last {
			return false
		}
	}
	return true
}

func (f *bitmapFont) stringWidth(text string) int {
	if f == nil {
		return 0
	}
	width := 0
	for _, char := range text {
		index := int(char) - f.meta.First
		if index < 0 || index >= len(f.meta.Widths) {
			continue
		}
		width += f.meta.Widths[index]
	}
	return width
}

// drawText mirrors h.drawTextWithFlags and PlatformGraphics.drawString for
// the two FreeJ2ME fonts used by Stage 1.
func (f *bitmapFont) drawText(dst *ebiten.Image, text string, x, y int, centered bool, tint color.Color) {
	if f == nil || f.image == nil {
		return
	}
	if centered {
		x -= f.stringWidth(text) / 2
	}
	strideX := f.meta.CellWidth + f.meta.Padding*2
	strideY := f.meta.CellHeight + f.meta.Padding*2
	cursor := x
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
		op.GeoM.Translate(float64(cursor-f.meta.Padding), float64(y-sourceFontYOffset-f.meta.Padding))
		dst.DrawImage(src, op)
		cursor += f.meta.Widths[index]
	}
}

func drawSourcePanelLabel(dst *ebiten.Image, font *bitmapFont, text string, centerX, panelY int) {
	if font == nil {
		return
	}
	width := font.stringWidth(text)
	x := centerX - width/2 - 5
	y := panelY - 5
	w := width + 10
	h := font.meta.FontHeight + 10
	drawRoundedRect(dst, x, y, w, h, 5, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	drawRoundedRect(dst, x+1, y+1, w-2, h-2, 4, color.RGBA{0x0c, 0x2f, 0x39, 0xff})
	font.drawText(dst, text, centerX, panelY+10, true, color.White)
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
