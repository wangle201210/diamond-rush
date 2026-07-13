package originalgame

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const cjkFontEnvironment = "ORIGINALRUSH_CJK_FONT"

type cachedCJKFont struct {
	source    *opentype.Font
	path      string
	family    string
	subfamily string
}

var (
	cjkFontCacheMu sync.Mutex
	cjkFontCache   = make(map[string]*cachedCJKFont)
)

func loadSystemCJKFont() (*cachedCJKFont, error) {
	configured := strings.TrimSpace(os.Getenv(cjkFontEnvironment))
	if configured != "" {
		path := resolvePath(configured)
		font, err := loadCachedCJKFont(path)
		if err != nil {
			return nil, fmt.Errorf("load %s=%q: %w", cjkFontEnvironment, configured, err)
		}
		return font, nil
	}

	var failures []string
	for _, path := range systemCJKFontCandidates(runtime.GOOS) {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		font, err := loadCachedCJKFont(path)
		if err == nil {
			return font, nil
		}
		failures = append(failures, fmt.Sprintf("%s: %v", path, err))
	}
	if len(failures) > 0 {
		return nil, fmt.Errorf("no usable system Chinese font (%s); set %s to a .ttf, .otf, .ttc, or .otc file", strings.Join(failures, "; "), cjkFontEnvironment)
	}
	return nil, fmt.Errorf("no system Chinese font found; set %s to a .ttf, .otf, .ttc, or .otc file", cjkFontEnvironment)
}

func systemCJKFontCandidates(goos string) []string {
	switch goos {
	case "darwin":
		return []string{
			"/System/Library/Fonts/PingFang.ttc",
			"/System/Library/Fonts/Hiragino Sans GB.ttc",
			"/System/Library/Fonts/STHeiti Medium.ttc",
			"/System/Library/Fonts/STHeiti Light.ttc",
		}
	case "windows":
		root := os.Getenv("WINDIR")
		if root == "" {
			root = `C:\Windows`
		}
		return []string{
			filepath.Join(root, "Fonts", "msyh.ttc"),
			filepath.Join(root, "Fonts", "simhei.ttf"),
			filepath.Join(root, "Fonts", "simsun.ttc"),
		}
	default:
		return []string{
			"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/opentype/noto/NotoSansCJKsc-Regular.otf",
			"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
			"/usr/share/fonts/truetype/arphic/uming.ttc",
		}
	}
}

func loadCachedCJKFont(path string) (*cachedCJKFont, error) {
	cleanPath := filepath.Clean(path)
	cjkFontCacheMu.Lock()
	defer cjkFontCacheMu.Unlock()
	if font := cjkFontCache[cleanPath]; font != nil {
		return font, nil
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}
	font, family, subfamily, err := parseCJKFont(data)
	if err != nil {
		return nil, err
	}
	result := &cachedCJKFont{
		source:    font,
		path:      cleanPath,
		family:    family,
		subfamily: subfamily,
	}
	cjkFontCache[cleanPath] = result
	return result, nil
}

func parseCJKFont(data []byte) (*opentype.Font, string, string, error) {
	collection, err := opentype.ParseCollection(data)
	if err != nil {
		return nil, "", "", err
	}
	if collection.NumFonts() == 0 {
		return nil, "", "", fmt.Errorf("font collection has no faces")
	}

	var best *opentype.Font
	bestFamily, bestSubfamily := "", ""
	bestScore := -1
	for index := 0; index < collection.NumFonts(); index++ {
		font, err := collection.Font(index)
		if err != nil {
			return nil, "", "", err
		}
		family, _ := font.Name(nil, sfnt.NameIDFamily)
		subfamily, _ := font.Name(nil, sfnt.NameIDSubfamily)
		score := cjkFontScore(font, family, subfamily)
		if score > bestScore {
			best = font
			bestFamily = family
			bestSubfamily = subfamily
			bestScore = score
		}
	}
	if best == nil || bestScore < 10000 {
		return nil, "", "", fmt.Errorf("font has no Simplified Chinese glyphs")
	}
	return best, bestFamily, bestSubfamily, nil
}

func cjkFontScore(font *opentype.Font, family, subfamily string) int {
	var buffer sfnt.Buffer
	for _, char := range "中文，。！？" {
		glyph, err := font.GlyphIndex(&buffer, char)
		if err != nil || glyph == 0 {
			return -1
		}
	}

	name := strings.ToLower(family + " " + subfamily)
	score := 10000
	for _, marker := range []string{"simplified", " sans gb", " sc", "简"} {
		if strings.Contains(name, marker) {
			score += 1000
			break
		}
	}
	for _, marker := range []string{"semibold", "demibold", "medium", "bold", "w6", "w7"} {
		if strings.Contains(name, marker) {
			score += 500
			break
		}
	}
	for _, marker := range []string{"ultralight", "thin", "light", "w1", "w2", "w3"} {
		if strings.Contains(name, marker) {
			score -= 250
			break
		}
	}
	return score
}
