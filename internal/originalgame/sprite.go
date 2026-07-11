package originalgame

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type spriteFrameMeta struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type spriteModuleMeta struct {
	W int `json:"w"`
	H int `json:"h"`
}

type spriteAnimationFrame struct {
	Frame int `json:"frame"`
	Time  int `json:"time"`
	X     int `json:"x"`
	Y     int `json:"y"`
	Flags int `json:"flags"`
}

type spriteFrameModule struct {
	Module int `json:"module"`
	X      int `json:"x"`
	Y      int `json:"y"`
	Flags  int `json:"flags"`
}

type spriteMetadata struct {
	Modules         []spriteModuleMeta     `json:"modules"`
	Frames          []spriteFrameMeta      `json:"frames"`
	FrameModules    []spriteFrameModule    `json:"frameModules"`
	FrameFirst      []int                  `json:"frameFirst"`
	FrameCounts     []int                  `json:"frameCounts"`
	AnimationFrames []spriteAnimationFrame `json:"animationFrames"`
	AnimationFirst  []int                  `json:"animationFirst"`
	AnimationCounts []int                  `json:"animationCounts"`
}

type spriteSheet struct {
	image       *ebiten.Image
	moduleImage *ebiten.Image
	meta        spriteMetadata
	cellW       int
	cellH       int
	moduleCellW int
	moduleCellH int
}

func loadSpriteSheet(imagePath, metadataPath string) (*spriteSheet, error) {
	return loadSpriteSheetAssets(imagePath, "", metadataPath)
}

func loadSpriteSheetWithModules(imagePath, modulePath, metadataPath string) (*spriteSheet, error) {
	return loadSpriteSheetAssets(imagePath, modulePath, metadataPath)
}

func loadModuleSpriteSheet(modulePath, metadataPath string) (*spriteSheet, error) {
	return loadSpriteSheetAssets("", modulePath, metadataPath)
}

func loadSpriteSheetAssets(imagePath, modulePath, metadataPath string) (*spriteSheet, error) {
	var img *ebiten.Image
	var err error
	if imagePath != "" {
		img, err = loadTransparentSheet(imagePath)
		if err != nil {
			return nil, err
		}
	}
	var moduleImage *ebiten.Image
	if modulePath != "" {
		moduleImage, err = loadTransparentSheet(modulePath)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.Open(filepath.Clean(resolvePath(metadataPath)))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var metadata spriteMetadata
	if err := json.NewDecoder(f).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decode sprite metadata: %w", err)
	}
	if len(metadata.Frames) == 0 && (moduleImage == nil || len(metadata.Modules) == 0) {
		return nil, fmt.Errorf("sprite metadata has no drawable frames or modules: %s", metadataPath)
	}
	cellW := 1
	cellH := 1
	for _, frame := range metadata.Frames {
		cellW = max(cellW, frame.W)
		cellH = max(cellH, frame.H)
	}
	moduleCellW := 1
	moduleCellH := 1
	for _, module := range metadata.Modules {
		moduleCellW = max(moduleCellW, module.W)
		moduleCellH = max(moduleCellH, module.H)
	}
	return &spriteSheet{
		image:       img,
		moduleImage: moduleImage,
		meta:        metadata,
		cellW:       cellW,
		cellH:       cellH,
		moduleCellW: moduleCellW,
		moduleCellH: moduleCellH,
	}, nil
}

func (s *spriteSheet) drawAnimation(dst *ebiten.Image, animation, tick, px, py, flags int) {
	frame, ok := s.animationFrame(animation, tick)
	if !ok {
		return
	}
	s.drawFrame(dst, frame.Frame, px+frame.X, py+frame.Y, flags^frame.Flags)
}

func (s *spriteSheet) drawAnimationSequenceFrame(dst *ebiten.Image, animation, sequence, px, py, flags int) {
	frame, ok := s.animationFrameAtSequence(animation, sequence)
	if !ok {
		return
	}
	s.drawFrame(dst, frame.Frame, px+frame.X, py+frame.Y, flags^frame.Flags)
}

func (s *spriteSheet) drawAnimationRawSequenceFrame(dst *ebiten.Image, animation, sequence, px, py, flags int) {
	frame, ok := s.animationFrameAtRawSequence(animation, sequence)
	if !ok {
		return
	}
	s.drawFrame(dst, frame.Frame, px+frame.X, py+frame.Y, flags^frame.Flags)
}

func (s *spriteSheet) animationFrameAtSequence(animation, sequence int) (spriteAnimationFrame, bool) {
	if s == nil || animation < 0 || animation >= len(s.meta.AnimationCounts) || animation >= len(s.meta.AnimationFirst) {
		return spriteAnimationFrame{}, false
	}
	first := s.meta.AnimationFirst[animation]
	count := s.meta.AnimationCounts[animation]
	if first < 0 || count <= 0 || first+count > len(s.meta.AnimationFrames) {
		return spriteAnimationFrame{}, false
	}
	sequence %= count
	if sequence < 0 {
		sequence += count
	}
	return s.meta.AnimationFrames[first+sequence], true
}

// animationFrameAtRawSequence preserves f_Sprite.drawAnimationFrame's flat
// animation-frame indexing. The original result effects intentionally index
// past animation 0 into the following sequences.
func (s *spriteSheet) animationFrameAtRawSequence(animation, sequence int) (spriteAnimationFrame, bool) {
	if s == nil || animation < 0 || animation >= len(s.meta.AnimationFirst) || sequence < 0 {
		return spriteAnimationFrame{}, false
	}
	index := s.meta.AnimationFirst[animation] + sequence
	if index < 0 || index >= len(s.meta.AnimationFrames) {
		return spriteAnimationFrame{}, false
	}
	return s.meta.AnimationFrames[index], true
}

func (s *spriteSheet) animationFrame(animation, tick int) (spriteAnimationFrame, bool) {
	index, ok := s.animationSequenceIndex(animation, tick)
	if !ok {
		return spriteAnimationFrame{}, false
	}
	return s.meta.AnimationFrames[s.meta.AnimationFirst[animation]+index], true
}

func (s *spriteSheet) animationSequenceIndex(animation, tick int) (int, bool) {
	if s == nil || animation < 0 || animation >= len(s.meta.AnimationCounts) || animation >= len(s.meta.AnimationFirst) {
		return 0, false
	}
	first := s.meta.AnimationFirst[animation]
	count := s.meta.AnimationCounts[animation]
	if first < 0 || count <= 0 || first+count > len(s.meta.AnimationFrames) {
		return 0, false
	}
	total := 0
	for _, frame := range s.meta.AnimationFrames[first : first+count] {
		total += max(1, frame.Time)
	}
	if total <= 0 {
		return 0, false
	}
	remaining := tick % total
	if remaining < 0 {
		remaining += total
	}
	for index, frame := range s.meta.AnimationFrames[first : first+count] {
		duration := max(1, frame.Time)
		if remaining < duration {
			return index, true
		}
		remaining -= duration
	}
	return count - 1, true
}

func (s *spriteSheet) animationDuration(animation int) (int, bool) {
	if s == nil || animation < 0 || animation >= len(s.meta.AnimationCounts) || animation >= len(s.meta.AnimationFirst) {
		return 0, false
	}
	first := s.meta.AnimationFirst[animation]
	count := s.meta.AnimationCounts[animation]
	if first < 0 || count <= 0 || first+count > len(s.meta.AnimationFrames) {
		return 0, false
	}
	total := 0
	for _, frame := range s.meta.AnimationFrames[first : first+count] {
		total += max(1, frame.Time)
	}
	return total, total > 0
}

func (s *spriteSheet) drawFrame(dst *ebiten.Image, frameIndex, px, py, flags int) {
	if s != nil && s.moduleImage != nil && s.drawComposedFrame(dst, frameIndex, px, py, flags) {
		return
	}
	if s == nil || s.image == nil || frameIndex < 0 || frameIndex >= len(s.meta.Frames) {
		return
	}
	frame := s.meta.Frames[frameIndex]
	if frame.W <= 0 || frame.H <= 0 {
		return
	}
	srcX := framePadding + (frameIndex%frameCols)*(s.cellW+framePadding)
	srcY := framePadding + (frameIndex/frameCols)*(s.cellH+framePadding)
	if srcX+frame.W > s.image.Bounds().Dx() || srcY+frame.H > s.image.Bounds().Dy() {
		return
	}
	op := &ebiten.DrawImageOptions{}
	translateX := px + frame.X
	translateY := py + frame.Y
	if flags&1 != 0 {
		op.GeoM.Scale(-1, 1)
		translateX = px - frame.X
	}
	if flags&2 != 0 {
		op.GeoM.Scale(1, -1)
		translateY = py - frame.Y
	}
	op.GeoM.Translate(float64(translateX), float64(translateY))
	src := s.image.SubImage(image.Rect(srcX, srcY, srcX+frame.W, srcY+frame.H)).(*ebiten.Image)
	dst.DrawImage(src, op)
}

func (s *spriteSheet) drawComposedFrame(dst *ebiten.Image, frameIndex, px, py, flags int) bool {
	if frameIndex < 0 || frameIndex >= len(s.meta.FrameFirst) || frameIndex >= len(s.meta.FrameCounts) {
		return false
	}
	first := s.meta.FrameFirst[frameIndex]
	count := s.meta.FrameCounts[frameIndex]
	if first < 0 || count < 0 || first+count > len(s.meta.FrameModules) {
		return false
	}
	for _, frameModule := range s.meta.FrameModules[first : first+count] {
		if frameModule.Module < 0 || frameModule.Module >= len(s.meta.Modules) {
			continue
		}
		module := s.meta.Modules[frameModule.Module]
		x := frameModule.X
		y := frameModule.Y
		moduleFlags := frameModule.Flags ^ flags
		if flags&1 != 0 {
			x = -x - module.W
		}
		if flags&2 != 0 {
			y = -y - module.H
		}
		s.drawModuleWithFlags(dst, frameModule.Module, px+x, py+y, moduleFlags)
	}
	return true
}

func (s *spriteSheet) drawModule(dst *ebiten.Image, moduleIndex, px, py int) {
	s.drawModuleWithFlags(dst, moduleIndex, px, py, 0)
}

func (s *spriteSheet) drawModuleWithFlags(dst *ebiten.Image, moduleIndex, px, py, flags int) {
	if s == nil || s.moduleImage == nil || moduleIndex < 0 || moduleIndex >= len(s.meta.Modules) {
		return
	}
	module := s.meta.Modules[moduleIndex]
	if module.W <= 0 || module.H <= 0 {
		return
	}
	srcX := framePadding + (moduleIndex%frameCols)*(s.moduleCellW+framePadding)
	srcY := framePadding + (moduleIndex/frameCols)*(s.moduleCellH+framePadding)
	if srcX+module.W > s.moduleImage.Bounds().Dx() || srcY+module.H > s.moduleImage.Bounds().Dy() {
		return
	}
	op := &ebiten.DrawImageOptions{}
	translateX := px
	translateY := py
	if flags&1 != 0 {
		op.GeoM.Scale(-1, 1)
		translateX += module.W
	}
	if flags&2 != 0 {
		op.GeoM.Scale(1, -1)
		translateY += module.H
	}
	op.GeoM.Translate(float64(translateX), float64(translateY))
	src := s.moduleImage.SubImage(image.Rect(srcX, srcY, srcX+module.W, srcY+module.H)).(*ebiten.Image)
	dst.DrawImage(src, op)
}

func (s *spriteSheet) drawNumber(dst *ebiten.Image, value, rightX, y int) {
	if value < 0 {
		value = 0
	}
	if value == 0 {
		rightX -= s.moduleWidth(0)
		s.drawModule(dst, 0, rightX, y)
		return
	}
	for value > 0 {
		digit := value % 10
		value /= 10
		rightX -= s.moduleWidth(digit)
		s.drawModule(dst, digit, rightX, y)
	}
}

func (s *spriteSheet) moduleWidth(moduleIndex int) int {
	if s == nil || moduleIndex < 0 || moduleIndex >= len(s.meta.Modules) {
		return 0
	}
	return s.meta.Modules[moduleIndex].W
}
