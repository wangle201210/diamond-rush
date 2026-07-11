package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	pixel8888 = 0x8888
	pixel4444 = 0x4444
	pixel1555 = 0x5515
	pixel0565 = 0x6505

	encI256    = 0x5602
	encI16     = 0x1600
	encI4      = 0x0400
	encI2      = 0x0200
	encI256RLE = 0x56F2
	encI127RLE = 0x27F1

	flagFlipX = 0x1
	flagFlipY = 0x2
)

type chunkInfo struct {
	Index  int    `json:"index"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	File   string `json:"file,omitempty"`
	Error  string `json:"error,omitempty"`
	Sprite *meta  `json:"sprite,omitempty"`
}

type meta struct {
	Modules       int            `json:"modules"`
	Frames        int            `json:"frames"`
	Animations    int            `json:"animations"`
	Palettes      int            `json:"palettes"`
	Colors        int            `json:"colors"`
	DataFormat    int            `json:"dataFormat"`
	ModuleSheet   string         `json:"moduleSheet,omitempty"`
	FrameSheet    string         `json:"frameSheet,omitempty"`
	PaletteSheets []paletteSheet `json:"paletteSheets,omitempty"`
	AnimationJSON string         `json:"animationJson,omitempty"`
}

type paletteSheet struct {
	Palette     int    `json:"palette"`
	ModuleSheet string `json:"moduleSheet,omitempty"`
	FrameSheet  string `json:"frameSheet,omitempty"`
}

type sprite struct {
	modules             []module
	frameModules        []frameModule
	frameFirst          []int
	frameCounts         []int
	frameBBoxes         []bbox
	animationFrames     []animationFrame
	animationFirst      []int
	animationCounts     []int
	palettes            [][]uint32
	dataFormat          int
	moduleDataPointers  []int
	moduleData          []byte
	colorsPerPalette    int
	hasModuleData       bool
	hasRenderableFrames bool
}

type module struct {
	W int `json:"w"`
	H int `json:"h"`
}

type frameModule struct {
	Module int `json:"module"`
	X      int `json:"x"`
	Y      int `json:"y"`
	Flags  int `json:"flags"`
}

type bbox struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type animationFrame struct {
	Frame int `json:"frame"`
	Time  int `json:"time"`
	X     int `json:"x"`
	Y     int `json:"y"`
	Flags int `json:"flags"`
}

func main() {
	in := flag.String("in", "/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources", "resource directory containing .f files")
	out := flag.String("out", "decoded/sprites", "output directory")
	flag.Parse()

	paths, err := filepath.Glob(filepath.Join(*in, "*.f"))
	if err != nil {
		fatal(err)
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		fatal(fmt.Errorf("no .f files found under %s", *in))
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fatal(err)
	}

	var all []chunkInfo
	for _, path := range paths {
		infos, err := exportFile(path, *out)
		if err != nil {
			fatal(err)
		}
		all = append(all, infos...)
	}
	if err := writeJSON(filepath.Join(*out, "manifest.json"), all); err != nil {
		fatal(err)
	}
	if err := writeHTMLIndex(filepath.Join(*out, "index.html"), all); err != nil {
		fatal(err)
	}

	ok := 0
	for _, info := range all {
		if info.Sprite != nil {
			ok++
		}
	}
	fmt.Printf("exported %d sprite chunks from %d chunks into %s\n", ok, len(all), *out)
}

func exportFile(path, outRoot string) ([]chunkInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty resource: %s", path)
	}
	chunkCount := int(data[0])
	headerLen := 1 + chunkCount*8
	if len(data) < headerLen {
		return nil, fmt.Errorf("%s header exceeds file", path)
	}

	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	outDir := filepath.Join(outRoot, base)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}

	infos := make([]chunkInfo, 0, chunkCount)
	for i := 0; i < chunkCount; i++ {
		entry := 1 + i*8
		offset := le32(data[entry : entry+4])
		length := le32(data[entry+4 : entry+8])
		info := chunkInfo{
			Index:  i,
			Offset: offset,
			Length: length,
			File:   filepath.Base(path),
		}
		start := headerLen + offset
		end := start + length
		if start < headerLen || end > len(data) || length <= 0 {
			info.Error = "invalid chunk range"
			infos = append(infos, info)
			continue
		}
		chunk := data[start:end]
		sp, err := parseSprite(chunk)
		if err != nil {
			info.Error = err.Error()
			infos = append(infos, info)
			continue
		}
		prefix := fmt.Sprintf("chunk%02d", i)
		m := &meta{
			Modules:       len(sp.modules),
			Frames:        len(sp.frameCounts),
			Animations:    len(sp.animationCounts),
			Palettes:      len(sp.palettes),
			Colors:        sp.colorsPerPalette,
			DataFormat:    sp.dataFormat,
			AnimationJSON: prefix + "-animations.json",
		}
		if len(sp.modules) > 0 && sp.hasModuleData {
			name := prefix + "-modules.png"
			if err := writePNG(filepath.Join(outDir, name), sp.renderModuleSheet(0)); err != nil {
				return nil, err
			}
			m.ModuleSheet = name
		}
		if len(sp.frameCounts) > 0 && sp.hasModuleData {
			name := prefix + "-frames.png"
			if err := writePNG(filepath.Join(outDir, name), sp.renderFrameSheet(0)); err != nil {
				return nil, err
			}
			m.FrameSheet = name
		}
		for palette := 1; palette < len(sp.palettes); palette++ {
			sheets := paletteSheet{Palette: palette}
			if len(sp.modules) > 0 && sp.hasModuleData {
				name := fmt.Sprintf("%s-palette%02d-modules.png", prefix, palette)
				if err := writePNG(filepath.Join(outDir, name), sp.renderModuleSheet(palette)); err != nil {
					return nil, err
				}
				sheets.ModuleSheet = name
			}
			if len(sp.frameCounts) > 0 && sp.hasModuleData {
				name := fmt.Sprintf("%s-palette%02d-frames.png", prefix, palette)
				if err := writePNG(filepath.Join(outDir, name), sp.renderFrameSheet(palette)); err != nil {
					return nil, err
				}
				sheets.FrameSheet = name
			}
			m.PaletteSheets = append(m.PaletteSheets, sheets)
		}
		if err := writeJSON(filepath.Join(outDir, m.AnimationJSON), sp.animationSummary()); err != nil {
			return nil, err
		}
		info.Sprite = m
		infos = append(infos, info)
	}
	return infos, nil
}

func parseSprite(data []byte) (*sprite, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("chunk too small for sprite")
	}
	ptr := 6
	moduleCount := le16(data[ptr : ptr+2])
	ptr += 2

	sp := &sprite{}
	if moduleCount > 0 {
		if ptr+moduleCount*2 > len(data) {
			return nil, fmt.Errorf("module table exceeds chunk")
		}
		sp.modules = make([]module, moduleCount)
		for i := 0; i < moduleCount; i++ {
			sp.modules[i] = module{W: int(data[ptr+i*2]), H: int(data[ptr+i*2+1])}
		}
		ptr += moduleCount * 2
	}

	if ptr+2 > len(data) {
		return nil, fmt.Errorf("missing frame module count")
	}
	frameModuleCount := le16(data[ptr : ptr+2])
	ptr += 2
	if frameModuleCount > 0 {
		if ptr+frameModuleCount*4 > len(data) {
			return nil, fmt.Errorf("frame module table exceeds chunk")
		}
		sp.frameModules = make([]frameModule, frameModuleCount)
		for i := 0; i < frameModuleCount; i++ {
			sp.frameModules[i] = frameModule{
				Module: int(data[ptr+i*4]),
				X:      int(int8(data[ptr+i*4+1])),
				Y:      int(int8(data[ptr+i*4+2])),
				Flags:  int(data[ptr+i*4+3]),
			}
		}
		ptr += frameModuleCount * 4
	}

	if ptr+2 > len(data) {
		return nil, fmt.Errorf("missing frame count")
	}
	frameCount := le16(data[ptr : ptr+2])
	ptr += 2
	if frameCount > 0 {
		if ptr+frameCount*4+frameCount*4 > len(data) {
			return nil, fmt.Errorf("frame tables exceed chunk")
		}
		sp.frameCounts = make([]int, frameCount)
		sp.frameFirst = make([]int, frameCount)
		for i := 0; i < frameCount; i++ {
			sp.frameCounts[i] = int(data[ptr])
			ptr += 2
			sp.frameFirst[i] = le16(data[ptr : ptr+2])
			ptr += 2
		}
		sp.frameBBoxes = make([]bbox, frameCount)
		for i := 0; i < frameCount; i++ {
			sp.frameBBoxes[i] = bbox{
				X: int(int8(data[ptr])),
				Y: int(int8(data[ptr+1])),
				W: int(data[ptr+2]),
				H: int(data[ptr+3]),
			}
			ptr += 4
		}
	}

	if ptr+2 > len(data) {
		return nil, fmt.Errorf("missing animation frame count")
	}
	animationFrameCount := le16(data[ptr : ptr+2])
	ptr += 2
	if animationFrameCount > 0 {
		if ptr+animationFrameCount*5 > len(data) {
			return nil, fmt.Errorf("animation frame table exceeds chunk")
		}
		sp.animationFrames = make([]animationFrame, animationFrameCount)
		for i := 0; i < animationFrameCount; i++ {
			sp.animationFrames[i] = animationFrame{
				Frame: int(data[ptr+i*5]),
				Time:  int(data[ptr+i*5+1]),
				X:     int(int8(data[ptr+i*5+2])),
				Y:     int(int8(data[ptr+i*5+3])),
				Flags: int(data[ptr+i*5+4]),
			}
		}
		ptr += animationFrameCount * 5
	}

	if ptr+2 > len(data) {
		return nil, fmt.Errorf("missing animation count")
	}
	animationCount := le16(data[ptr : ptr+2])
	ptr += 2
	if animationCount > 0 {
		if ptr+animationCount*4 > len(data) {
			return nil, fmt.Errorf("animation tables exceed chunk")
		}
		sp.animationCounts = make([]int, animationCount)
		sp.animationFirst = make([]int, animationCount)
		for i := 0; i < animationCount; i++ {
			sp.animationCounts[i] = int(data[ptr])
			ptr += 2
			sp.animationFirst[i] = le16(data[ptr : ptr+2])
			ptr += 2
		}
	}

	if moduleCount == 0 {
		return sp, nil
	}
	if ptr+4 > len(data) {
		return nil, fmt.Errorf("missing palette header")
	}
	pixelFormat := le16(data[ptr : ptr+2])
	ptr += 2
	paletteCount := int(data[ptr])
	ptr++
	colorsPerPalette := int(data[ptr])
	ptr++
	if paletteCount <= 0 || colorsPerPalette <= 0 {
		return nil, fmt.Errorf("empty palette table")
	}
	sp.colorsPerPalette = colorsPerPalette
	sp.palettes = make([][]uint32, paletteCount)
	for p := 0; p < paletteCount; p++ {
		palette := make([]uint32, colorsPerPalette)
		for c := 0; c < colorsPerPalette; c++ {
			col, next, err := readColor(data, ptr, pixelFormat)
			if err != nil {
				return nil, err
			}
			ptr = next
			palette[c] = col
		}
		sp.palettes[p] = palette
	}
	if ptr+2 > len(data) {
		return nil, fmt.Errorf("missing data format")
	}
	sp.dataFormat = le16(data[ptr : ptr+2])
	ptr += 2

	sp.moduleDataPointers = make([]int, moduleCount)
	sizePtr := ptr
	total := 0
	for i := 0; i < moduleCount; i++ {
		if sizePtr+2 > len(data) {
			return nil, fmt.Errorf("module data size table exceeds chunk")
		}
		size := le16(data[sizePtr : sizePtr+2])
		sizePtr += 2 + size
		if sizePtr > len(data) {
			return nil, fmt.Errorf("module data exceeds chunk")
		}
		sp.moduleDataPointers[i] = total
		total += size
	}
	sp.moduleData = make([]byte, total)
	dst := 0
	for i := 0; i < moduleCount; i++ {
		size := le16(data[ptr : ptr+2])
		ptr += 2
		copy(sp.moduleData[dst:dst+size], data[ptr:ptr+size])
		ptr += size
		dst += size
	}
	sp.hasModuleData = true
	return sp, nil
}

func readColor(data []byte, ptr, pixelFormat int) (uint32, int, error) {
	switch pixelFormat {
	case pixel8888:
		if ptr+4 > len(data) {
			return 0, ptr, fmt.Errorf("8888 palette exceeds chunk")
		}
		v := uint32(data[ptr]) | uint32(data[ptr+1])<<8 | uint32(data[ptr+2])<<16 | uint32(data[ptr+3])<<24
		return v, ptr + 4, nil
	case pixel4444:
		if ptr+2 > len(data) {
			return 0, ptr, fmt.Errorf("4444 palette exceeds chunk")
		}
		v := uint32(le16(data[ptr : ptr+2]))
		argb := (v&0xF000)<<16 | (v&0xF000)<<12 | (v&0x0F00)<<12 | (v&0x0F00)<<8 | (v&0x00F0)<<8 | (v&0x00F0)<<4 | (v&0x000F)<<4 | (v & 0x000F)
		return argb, ptr + 2, nil
	case pixel1555:
		if ptr+2 > len(data) {
			return 0, ptr, fmt.Errorf("1555 palette exceeds chunk")
		}
		v := uint32(le16(data[ptr : ptr+2]))
		alpha := uint32(0xFF000000)
		if v&0x8000 == 0 {
			alpha = 0
		}
		return alpha | (v&0x7C00)<<9 | (v&0x03E0)<<6 | (v&0x001F)<<3, ptr + 2, nil
	case pixel0565:
		if ptr+2 > len(data) {
			return 0, ptr, fmt.Errorf("0565 palette exceeds chunk")
		}
		v := uint32(le16(data[ptr : ptr+2]))
		alpha := uint32(0xFF000000)
		if v == 0xF81F {
			alpha = 0
		}
		return alpha | (v&0xF800)<<8 | (v&0x07E0)<<5 | (v&0x001F)<<3, ptr + 2, nil
	default:
		return 0, ptr, fmt.Errorf("unsupported pixel format %#x", pixelFormat)
	}
}

func (sp *sprite) renderModuleSheet(palette int) image.Image {
	images := make([]*image.RGBA, len(sp.modules))
	maxW := 1
	maxH := 1
	for i := range sp.modules {
		img, err := sp.renderModule(i, palette)
		if err == nil {
			images[i] = img
			maxW = max(maxW, img.Bounds().Dx())
			maxH = max(maxH, img.Bounds().Dy())
		}
	}
	return packImages(images, maxW, maxH)
}

func (sp *sprite) renderFrameSheet(palette int) image.Image {
	images := make([]*image.RGBA, len(sp.frameCounts))
	maxW := 1
	maxH := 1
	for i := range sp.frameCounts {
		img, err := sp.renderFrame(i, palette)
		if err == nil {
			images[i] = img
			maxW = max(maxW, img.Bounds().Dx())
			maxH = max(maxH, img.Bounds().Dy())
		}
	}
	return packImages(images, maxW, maxH)
}

func (sp *sprite) renderModule(i, palette int) (*image.RGBA, error) {
	if i < 0 || i >= len(sp.modules) {
		return nil, fmt.Errorf("module index out of range")
	}
	mod := sp.modules[i]
	if mod.W <= 0 || mod.H <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
	}
	pixels, err := sp.parseImageData(i, palette)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, mod.W, mod.H))
	for y := 0; y < mod.H; y++ {
		for x := 0; x < mod.W; x++ {
			img.SetRGBA(x, y, argb(pixels[x+y*mod.W]))
		}
	}
	return img, nil
}

func (sp *sprite) renderFrame(frame, palette int) (*image.RGBA, error) {
	if frame < 0 || frame >= len(sp.frameCounts) {
		return nil, fmt.Errorf("frame index out of range")
	}
	minX, minY, maxX, maxY := 0, 0, 1, 1
	first := sp.frameFirst[frame]
	count := sp.frameCounts[frame]
	for i := 0; i < count; i++ {
		fm := sp.frameModules[first+i]
		if fm.Module < 0 || fm.Module >= len(sp.modules) {
			continue
		}
		mod := sp.modules[fm.Module]
		minX = min(minX, fm.X)
		minY = min(minY, fm.Y)
		maxX = max(maxX, fm.X+mod.W)
		maxY = max(maxY, fm.Y+mod.H)
	}
	img := image.NewRGBA(image.Rect(0, 0, maxX-minX, maxY-minY))
	for i := 0; i < count; i++ {
		fm := sp.frameModules[first+i]
		modImg, err := sp.renderModule(fm.Module, palette)
		if err != nil {
			continue
		}
		if fm.Flags&flagFlipX != 0 || fm.Flags&flagFlipY != 0 {
			modImg = flip(modImg, fm.Flags)
		}
		draw.Draw(img, image.Rect(fm.X-minX, fm.Y-minY, fm.X-minX+modImg.Bounds().Dx(), fm.Y-minY+modImg.Bounds().Dy()), modImg, image.Point{}, draw.Over)
	}
	return img, nil
}

func (sp *sprite) parseImageData(moduleIndex, paletteIndex int) ([]uint32, error) {
	if paletteIndex < 0 || paletteIndex >= len(sp.palettes) {
		return nil, fmt.Errorf("palette index out of range")
	}
	if moduleIndex < 0 || moduleIndex >= len(sp.modules) {
		return nil, fmt.Errorf("module index out of range")
	}
	mod := sp.modules[moduleIndex]
	total := mod.W * mod.H
	out := make([]uint32, 0, total)
	palette := sp.palettes[paletteIndex]
	ptr := sp.moduleDataPointers[moduleIndex]
	emit := func(idx int) error {
		if idx < 0 || idx >= len(palette) {
			return fmt.Errorf("palette color index %d out of range", idx)
		}
		out = append(out, palette[idx])
		return nil
	}
	for len(out) < total {
		if ptr >= len(sp.moduleData) {
			return nil, fmt.Errorf("module data underrun")
		}
		switch sp.dataFormat {
		case encI127RLE:
			n := int(sp.moduleData[ptr])
			ptr++
			if n > 127 {
				if ptr >= len(sp.moduleData) {
					return nil, fmt.Errorf("i127 rle underrun")
				}
				idx := int(sp.moduleData[ptr])
				ptr++
				for c := 0; c < n-128 && len(out) < total; c++ {
					if err := emit(idx); err != nil {
						return nil, err
					}
				}
			} else if err := emit(n); err != nil {
				return nil, err
			}
		case encI16:
			b := sp.moduleData[ptr]
			ptr++
			if err := emit(int(b >> 4 & 0xF)); err != nil {
				return nil, err
			}
			if len(out) < total {
				if err := emit(int(b & 0xF)); err != nil {
					return nil, err
				}
			}
		case encI4:
			b := sp.moduleData[ptr]
			ptr++
			for shift := 6; shift >= 0 && len(out) < total; shift -= 2 {
				if err := emit(int(b >> shift & 0x3)); err != nil {
					return nil, err
				}
			}
		case encI2:
			b := sp.moduleData[ptr]
			ptr++
			for shift := 7; shift >= 0 && len(out) < total; shift-- {
				if err := emit(int(b >> shift & 0x1)); err != nil {
					return nil, err
				}
			}
		case encI256:
			idx := int(sp.moduleData[ptr])
			ptr++
			if err := emit(idx); err != nil {
				return nil, err
			}
		case encI256RLE:
			n := int(sp.moduleData[ptr])
			ptr++
			if n > 127 {
				for c := 0; c < n-128 && len(out) < total; c++ {
					if ptr >= len(sp.moduleData) {
						return nil, fmt.Errorf("i256 rle underrun")
					}
					if err := emit(int(sp.moduleData[ptr])); err != nil {
						return nil, err
					}
					ptr++
				}
			} else {
				if ptr >= len(sp.moduleData) {
					return nil, fmt.Errorf("i256 rle underrun")
				}
				idx := int(sp.moduleData[ptr])
				ptr++
				for c := 0; c < n && len(out) < total; c++ {
					if err := emit(idx); err != nil {
						return nil, err
					}
				}
			}
		default:
			return nil, fmt.Errorf("unsupported data format %#x", sp.dataFormat)
		}
	}
	return out, nil
}

func (sp *sprite) animationSummary() any {
	return struct {
		Modules         []module         `json:"modules"`
		Frames          []bbox           `json:"frames"`
		FrameModules    []frameModule    `json:"frameModules"`
		FrameFirst      []int            `json:"frameFirst"`
		FrameCounts     []int            `json:"frameCounts"`
		AnimationFrames []animationFrame `json:"animationFrames"`
		AnimationFirst  []int            `json:"animationFirst"`
		AnimationCounts []int            `json:"animationCounts"`
	}{sp.modules, sp.frameBBoxes, sp.frameModules, sp.frameFirst, sp.frameCounts, sp.animationFrames, sp.animationFirst, sp.animationCounts}
}

func packImages(images []*image.RGBA, cellW, cellH int) image.Image {
	padding := 2
	cols := 16
	rows := (len(images) + cols - 1) / cols
	img := image.NewRGBA(image.Rect(0, 0, cols*(cellW+padding)+padding, rows*(cellH+padding)+padding))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{20, 22, 28, 255}), image.Point{}, draw.Src)
	for i, src := range images {
		if src == nil {
			continue
		}
		x := padding + (i%cols)*(cellW+padding)
		y := padding + (i/cols)*(cellH+padding)
		draw.Draw(img, image.Rect(x, y, x+src.Bounds().Dx(), y+src.Bounds().Dy()), src, image.Point{}, draw.Over)
	}
	return img
}

func flip(src *image.RGBA, flags int) *image.RGBA {
	w := src.Bounds().Dx()
	h := src.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx, sy := x, y
			if flags&flagFlipX != 0 {
				sx = w - 1 - x
			}
			if flags&flagFlipY != 0 {
				sy = h - 1 - y
			}
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}

func argb(v uint32) color.RGBA {
	return color.RGBA{
		R: uint8(v >> 16),
		G: uint8(v >> 8),
		B: uint8(v),
		A: uint8(v >> 24),
	}
}

func le16(b []byte) int {
	return int(b[0]) | int(b[1])<<8
}

func le32(b []byte) int {
	return int(b[0]) | int(b[1])<<8 | int(b[2])<<16 | int(b[3])<<24
}

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func writeHTMLIndex(path string, infos []chunkInfo) error {
	var b strings.Builder
	b.WriteString("<!doctype html><meta charset=\"utf-8\"><title>Diamond Rush sprites</title>")
	b.WriteString("<style>body{font:14px system-ui,sans-serif;background:#111;color:#eee;margin:24px}a{color:#8fd3ff}.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(320px,1fr));gap:18px}.card{border:1px solid #333;background:#1a1d24;padding:10px}.meta{color:#aaa}img{max-width:100%;image-rendering:pixelated;background:#141820}</style>")
	b.WriteString("<h1>Diamond Rush sprite exports</h1>")
	b.WriteString("<p>Generated by <code>tools/drsprite</code>. Sheets are raw sprite chunks; stage IDs still need Java render mapping before full-stage visual previews are exact.</p>")
	b.WriteString("<p><a href=\"manifest.json\">manifest.json</a></p><div class=\"grid\">")
	for _, info := range infos {
		if info.Sprite == nil {
			continue
		}
		base := strings.TrimSuffix(info.File, filepath.Ext(info.File))
		dir := base + "/"
		title := fmt.Sprintf("%s chunk %02d", info.File, info.Index)
		b.WriteString("<section class=\"card\"><h2>")
		b.WriteString(title)
		b.WriteString("</h2><p class=\"meta\">")
		b.WriteString(fmt.Sprintf("%d modules, %d frames, %d animations, %d palettes", info.Sprite.Modules, info.Sprite.Frames, info.Sprite.Animations, info.Sprite.Palettes))
		b.WriteString("</p>")
		if info.Sprite.FrameSheet != "" {
			src := dir + info.Sprite.FrameSheet
			b.WriteString("<p><a href=\"")
			b.WriteString(src)
			b.WriteString("\">frames</a></p><img src=\"")
			b.WriteString(src)
			b.WriteString("\">")
		} else if info.Sprite.ModuleSheet != "" {
			src := dir + info.Sprite.ModuleSheet
			b.WriteString("<p><a href=\"")
			b.WriteString(src)
			b.WriteString("\">modules</a></p><img src=\"")
			b.WriteString(src)
			b.WriteString("\">")
		}
		b.WriteString("</section>")
	}
	b.WriteString("</div>")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "drsprite:", err)
	os.Exit(1)
}
