package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
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
	cellSize  = 12
	layerGap  = 8
	labelGap  = 16
	maxIDByte = 255
)

type stageFile struct {
	Index  int              `json:"index"`
	Width  int              `json:"width"`
	Height int              `json:"height"`
	Layers map[string][]int `json:"layers"`
}

type idStats struct {
	Count  int
	Stages map[string]int
}

func main() {
	in := flag.String("in", "decoded", "decoded directory containing world*/stage*.json")
	out := flag.String("out", "decoded/preview", "output preview directory")
	flag.Parse()

	stagePaths, err := filepath.Glob(filepath.Join(*in, "world*", "stage*.json"))
	if err != nil {
		fatal(err)
	}
	sort.Strings(stagePaths)
	if len(stagePaths) == 0 {
		fatal(fmt.Errorf("no decoded stage json files found under %s", *in))
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fatal(err)
	}

	index := map[string]map[int]*idStats{
		"player":     {},
		"background": {},
		"foreground": {},
	}
	var htmlPages []string

	for _, path := range stagePaths {
		stage, err := readStage(path)
		if err != nil {
			fatal(err)
		}
		rel := relName(*in, path)
		stageID := strings.TrimSuffix(rel, ".json")
		for layerName, layer := range stage.Layers {
			for _, id := range layer {
				stats := index[layerName][id]
				if stats == nil {
					stats = &idStats{Stages: map[string]int{}}
					index[layerName][id] = stats
				}
				stats.Count++
				stats.Stages[stageID]++
			}
		}

		outDir := filepath.Join(*out, filepath.Dir(rel))
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			fatal(err)
		}
		if err := writeStagePreviews(stage, outDir, strings.TrimSuffix(filepath.Base(path), ".json")); err != nil {
			fatal(err)
		}
		htmlName := strings.TrimSuffix(filepath.Base(path), ".json") + ".html"
		if err := writeStageHTML(filepath.Join(outDir, htmlName), stage, stageID); err != nil {
			fatal(err)
		}
		htmlPages = append(htmlPages, filepath.ToSlash(filepath.Join(filepath.Dir(rel), htmlName)))
	}

	if err := writeIndex(filepath.Join(*out, "id-index.md"), index); err != nil {
		fatal(err)
	}
	if err := writeHTMLIndex(filepath.Join(*out, "index.html"), htmlPages); err != nil {
		fatal(err)
	}
	if err := writePalette(filepath.Join(*out, "palette.png")); err != nil {
		fatal(err)
	}

	fmt.Printf("inspected %d stages from %s into %s\n", len(stagePaths), *in, *out)
}

func readStage(path string) (stageFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return stageFile{}, err
	}
	var stage stageFile
	if err := json.Unmarshal(data, &stage); err != nil {
		return stageFile{}, err
	}
	for _, name := range []string{"player", "background", "foreground"} {
		layer := stage.Layers[name]
		if len(layer) != stage.Width*stage.Height {
			return stageFile{}, fmt.Errorf("%s %s layer has %d cells, want %d", path, name, len(layer), stage.Width*stage.Height)
		}
	}
	return stage, nil
}

func writeStagePreviews(stage stageFile, outDir, base string) error {
	for _, layerName := range []string{"background", "foreground", "player"} {
		img := renderLayer(stage.Width, stage.Height, stage.Layers[layerName], false)
		if err := writePNG(filepath.Join(outDir, base+"-"+layerName+".png"), img); err != nil {
			return err
		}
	}
	if err := writePNG(filepath.Join(outDir, base+"-composite.png"), renderComposite(stage)); err != nil {
		return err
	}
	if err := writePNG(filepath.Join(outDir, base+"-contactsheet.png"), renderContactSheet(stage)); err != nil {
		return err
	}
	return nil
}

func renderLayer(width, height int, layer []int, grid bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width*cellSize, height*cellSize))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{12, 13, 18, 255}), image.Point{}, draw.Src)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			id := layer[x+y*width]
			fillCell(img, x, y, id)
		}
	}
	if grid {
		drawGrid(img, width, height)
	}
	return img
}

func renderComposite(stage stageFile) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, stage.Width*cellSize, stage.Height*cellSize))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{12, 13, 18, 255}), image.Point{}, draw.Src)
	for y := 0; y < stage.Height; y++ {
		for x := 0; x < stage.Width; x++ {
			id := firstVisibleID(stage, x, y)
			fillCell(img, x, y, id)
		}
	}
	drawGrid(img, stage.Width, stage.Height)
	return img
}

func renderContactSheet(stage stageFile) *image.RGBA {
	layerW := stage.Width * cellSize
	layerH := stage.Height * cellSize
	totalW := layerW*2 + layerGap
	totalH := (layerH+labelGap)*2 + layerGap
	img := image.NewRGBA(image.Rect(0, 0, totalW, totalH))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{20, 22, 28, 255}), image.Point{}, draw.Src)

	entries := []struct {
		name string
		img  *image.RGBA
		x    int
		y    int
	}{
		{"background", renderLayer(stage.Width, stage.Height, stage.Layers["background"], true), 0, labelGap},
		{"foreground", renderLayer(stage.Width, stage.Height, stage.Layers["foreground"], true), layerW + layerGap, labelGap},
		{"player", renderLayer(stage.Width, stage.Height, stage.Layers["player"], true), 0, layerH + labelGap*2 + layerGap},
		{"composite", renderComposite(stage), layerW + layerGap, layerH + labelGap*2 + layerGap},
	}
	for _, e := range entries {
		drawTinyText(img, e.x, e.y-labelGap+2, e.name, color.RGBA{232, 236, 244, 255})
		draw.Draw(img, image.Rect(e.x, e.y, e.x+layerW, e.y+layerH), e.img, image.Point{}, draw.Src)
	}
	return img
}

func firstVisibleID(stage stageFile, x, y int) int {
	idx := x + y*stage.Width
	for _, layerName := range []string{"player", "foreground", "background"} {
		id := stage.Layers[layerName][idx]
		if id != maxIDByte {
			return id
		}
	}
	return maxIDByte
}

func fillCell(img *image.RGBA, x, y, id int) {
	c := idColor(id)
	r := image.Rect(x*cellSize, y*cellSize, (x+1)*cellSize, (y+1)*cellSize)
	draw.Draw(img, r, image.NewUniform(c), image.Point{}, draw.Src)
}

func idColor(id int) color.RGBA {
	if id == maxIDByte {
		return color.RGBA{8, 9, 12, 255}
	}
	if id < 0 {
		id = -id
	}
	h := uint32(id*1103515245 + 12345)
	r := uint8(55 + (h>>16)%180)
	g := uint8(55 + (h>>8)%180)
	b := uint8(55 + h%180)
	return color.RGBA{r, g, b, 255}
}

func drawGrid(img *image.RGBA, width, height int) {
	gridColor := color.RGBA{0, 0, 0, 65}
	for x := 0; x <= width; x++ {
		px := x * cellSize
		for y := 0; y < height*cellSize; y++ {
			img.Set(px, y, gridColor)
		}
	}
	for y := 0; y <= height; y++ {
		py := y * cellSize
		for x := 0; x < width*cellSize; x++ {
			img.Set(x, py, gridColor)
		}
	}
}

func writeIndex(path string, index map[string]map[int]*idStats) error {
	var b strings.Builder
	b.WriteString("# Diamond Rush Raw ID Index\n\n")
	b.WriteString("Generated by `tools/drinspect` from decoded stage JSON. The PNG previews are diagnostic raw-ID color maps, not original sprite art; use `decoded/sprites/index.html` for extracted original pixel assets.\n\n")
	for _, layerName := range []string{"player", "background", "foreground"} {
		b.WriteString("## " + layerName + "\n\n")
		b.WriteString("| ID | Count | Stage appearances |\n")
		b.WriteString("| ---: | ---: | --- |\n")
		ids := sortedIDs(index[layerName])
		for _, id := range ids {
			stats := index[layerName][id]
			b.WriteString(fmt.Sprintf("| %d | %d | %s |\n", id, stats.Count, formatStages(stats.Stages)))
		}
		b.WriteString("\n")
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeHTMLIndex(path string, pages []string) error {
	var b strings.Builder
	b.WriteString("<!doctype html><meta charset=\"utf-8\"><title>Diamond Rush decoded stages</title>")
	b.WriteString("<style>body{font:14px system-ui,sans-serif;background:#111;color:#eee;margin:24px}a{color:#8fd3ff}.grid{columns:3 240px}li{break-inside:avoid;margin:4px 0}</style>")
	b.WriteString("<h1>Diamond Rush decoded stages</h1>")
	b.WriteString("<p>Open a stage to inspect raw byte IDs. These PNG/HTML previews are diagnostic color maps, not original sprite art. Hover cells for coordinates and layer IDs.</p>")
	b.WriteString("<p>Extracted original pixel assets are under <code>../sprites/</code>, especially <code>../sprites/index.html</code>.</p>")
	b.WriteString("<p><a href=\"id-index.md\">Raw ID index</a> | <a href=\"palette.png\">Palette PNG</a></p>")
	b.WriteString("<ul class=\"grid\">")
	for _, page := range pages {
		b.WriteString("<li><a href=\"")
		b.WriteString(html.EscapeString(page))
		b.WriteString("\">")
		b.WriteString(html.EscapeString(strings.TrimSuffix(page, ".html")))
		b.WriteString("</a></li>")
	}
	b.WriteString("</ul>")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeStageHTML(path string, stage stageFile, stageID string) error {
	var b strings.Builder
	b.WriteString("<!doctype html><meta charset=\"utf-8\">")
	b.WriteString("<title>")
	b.WriteString(html.EscapeString(stageID))
	b.WriteString("</title>")
	b.WriteString(`<style>
body{font:12px system-ui,sans-serif;background:#111;color:#eee;margin:18px}
a{color:#8fd3ff}
.meta{color:#aaa}
.wrap{display:flex;gap:18px;flex-wrap:wrap;align-items:flex-start}
.panel{background:#1a1d24;padding:10px;border:1px solid #333}
.layer{display:grid;grid-template-columns:repeat(var(--w),24px);grid-auto-rows:24px;gap:1px;background:#050609;padding:1px}
.cell{width:24px;height:24px;display:flex;align-items:center;justify-content:center;font-size:9px;color:rgba(255,255,255,.9);text-shadow:0 1px 1px #000;box-sizing:border-box}
.empty{color:transparent}
h2{margin:0 0 8px;font-size:14px}
table{border-collapse:collapse;margin-top:18px}
td,th{border:1px solid #333;padding:4px 6px;text-align:right}
th{text-align:left;background:#20242d}
</style>`)
	b.WriteString("<p><a href=\"../index.html\">Index</a></p>")
	b.WriteString("<h1>")
	b.WriteString(html.EscapeString(stageID))
	b.WriteString("</h1>")
	b.WriteString(fmt.Sprintf("<p class=\"meta\">%dx%d cells. Colors are raw byte ID diagnostics, not original art. Cell text is raw byte ID; hover for x/y/layer.</p>", stage.Width, stage.Height))
	b.WriteString("<div class=\"wrap\">")
	for _, layerName := range []string{"background", "foreground", "player", "composite"} {
		b.WriteString("<section class=\"panel\"><h2>")
		b.WriteString(layerName)
		b.WriteString("</h2>")
		b.WriteString(fmt.Sprintf("<div class=\"layer\" style=\"--w:%d\">", stage.Width))
		for y := 0; y < stage.Height; y++ {
			for x := 0; x < stage.Width; x++ {
				id := 255
				titleLayer := layerName
				if layerName == "composite" {
					id = firstVisibleID(stage, x, y)
					titleLayer = fmt.Sprintf("p:%d f:%d b:%d", stage.Layers["player"][x+y*stage.Width], stage.Layers["foreground"][x+y*stage.Width], stage.Layers["background"][x+y*stage.Width])
				} else {
					id = stage.Layers[layerName][x+y*stage.Width]
				}
				className := "cell"
				text := fmt.Sprintf("%d", id)
				if id == 255 {
					className += " empty"
					text = "."
				}
				c := idColor(id)
				b.WriteString(fmt.Sprintf("<div class=\"%s\" style=\"background:#%02x%02x%02x\" title=\"x:%d y:%d %s id:%d\">%s</div>",
					className, c.R, c.G, c.B, x, y, html.EscapeString(titleLayer), id, html.EscapeString(text)))
			}
		}
		b.WriteString("</div></section>")
	}
	b.WriteString("</div>")
	b.WriteString("<h2>Histograms</h2>")
	for _, layerName := range []string{"background", "foreground", "player"} {
		b.WriteString("<h3>")
		b.WriteString(layerName)
		b.WriteString("</h3><table><tr><th>ID</th><th>Count</th></tr>")
		h := map[int]int{}
		for _, id := range stage.Layers[layerName] {
			h[id]++
		}
		ids := make([]int, 0, len(h))
		for id := range h {
			ids = append(ids, id)
		}
		sort.Ints(ids)
		for _, id := range ids {
			b.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%d</td></tr>", id, h[id]))
		}
		b.WriteString("</table>")
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writePalette(path string) error {
	cols := 16
	rows := 16
	img := image.NewRGBA(image.Rect(0, 0, cols*cellSize, rows*cellSize))
	for id := 0; id <= maxIDByte; id++ {
		x := id % cols
		y := id / cols
		fillCell(img, x, y, id)
	}
	drawGrid(img, cols, rows)
	return writePNG(path, img)
}

func sortedIDs(m map[int]*idStats) []int {
	ids := make([]int, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	return ids
}

func formatStages(stages map[string]int) string {
	keys := make([]string, 0, len(stages))
	for key := range stages {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) > 18 {
		keys = append(keys[:18], fmt.Sprintf("... +%d more", len(keys)-18))
	}
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		if strings.HasPrefix(key, "...") {
			parts = append(parts, key)
			continue
		}
		parts = append(parts, fmt.Sprintf("%s:%d", key, stages[key]))
	}
	return strings.Join(parts, ", ")
}

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func relName(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.Base(path)
	}
	return rel
}

// drawTinyText renders a tiny subset of ASCII for labels. Unknown glyphs become blanks.
func drawTinyText(img *image.RGBA, x, y int, text string, c color.RGBA) {
	for _, r := range strings.ToLower(text) {
		if r == ' ' {
			x += 4
			continue
		}
		glyph, ok := tinyFont[r]
		if !ok {
			x += 4
			continue
		}
		for gy, row := range glyph {
			for gx, bit := range row {
				if bit == '1' {
					img.Set(x+gx, y+gy, c)
				}
			}
		}
		x += 4
	}
}

var tinyFont = map[rune][]string{
	'a': {"010", "101", "111", "101", "101"},
	'b': {"110", "101", "110", "101", "110"},
	'c': {"011", "100", "100", "100", "011"},
	'd': {"110", "101", "101", "101", "110"},
	'e': {"111", "100", "110", "100", "111"},
	'f': {"111", "100", "110", "100", "100"},
	'g': {"011", "100", "101", "101", "011"},
	'h': {"101", "101", "111", "101", "101"},
	'i': {"111", "010", "010", "010", "111"},
	'j': {"001", "001", "001", "101", "010"},
	'k': {"101", "101", "110", "101", "101"},
	'l': {"100", "100", "100", "100", "111"},
	'm': {"101", "111", "111", "101", "101"},
	'n': {"101", "111", "111", "111", "101"},
	'o': {"010", "101", "101", "101", "010"},
	'p': {"110", "101", "110", "100", "100"},
	'q': {"010", "101", "101", "111", "011"},
	'r': {"110", "101", "110", "101", "101"},
	's': {"011", "100", "010", "001", "110"},
	't': {"111", "010", "010", "010", "010"},
	'u': {"101", "101", "101", "101", "111"},
	'v': {"101", "101", "101", "101", "010"},
	'w': {"101", "101", "111", "111", "101"},
	'x': {"101", "101", "010", "101", "101"},
	'y': {"101", "101", "010", "010", "010"},
	'z': {"111", "001", "010", "100", "111"},
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "drinspect:", err)
	os.Exit(1)
}
