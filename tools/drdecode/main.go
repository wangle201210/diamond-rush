package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type manifest struct {
	Source      string         `json:"source"`
	InitialByte int            `json:"initialByte"`
	StageCount  int            `json:"stageCount"`
	Stages      []stageSummary `json:"stages"`
}

type stageSummary struct {
	Index      int                    `json:"index"`
	Group      int                    `json:"group"`
	GroupIndex int                    `json:"groupIndex"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Histograms map[string]map[int]int `json:"histograms"`
	File       string                 `json:"file"`
}

type stageFile struct {
	Index      int                    `json:"index"`
	Group      int                    `json:"group"`
	GroupIndex int                    `json:"groupIndex"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Histograms map[string]map[int]int `json:"histograms"`
	Layers     map[string][]int       `json:"layers"`
	Notes      []string               `json:"notes,omitempty"`
}

type stage struct {
	group      int
	groupIndex int
	width      int
	height     int
	player     []byte
	background []byte
	foreground []byte
}

func main() {
	in := flag.String("in", "", "input w*.bin file")
	out := flag.String("out", "", "output directory")
	flag.Parse()

	if *in == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "usage: drdecode -in path/to/w0.bin -out decoded/world0")
		os.Exit(2)
	}

	data, err := os.ReadFile(*in)
	if err != nil {
		fatal(err)
	}
	if len(data) == 0 {
		fatal(fmt.Errorf("empty input: %s", *in))
	}

	initialByte := int(data[0])
	stages, err := decodeWorld(data)
	if err != nil {
		fatal(err)
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fatal(err)
	}

	m := manifest{
		Source:      *in,
		InitialByte: initialByte,
		StageCount:  len(stages),
		Stages:      make([]stageSummary, 0, len(stages)),
	}

	for i, st := range stages {
		fileName := fmt.Sprintf("stage%02d.json", i)
		histograms := map[string]map[int]int{
			"player":     histogram(st.player),
			"background": histogram(st.background),
			"foreground": histogram(st.foreground),
		}
		sf := stageFile{
			Index:      i,
			Group:      st.group,
			GroupIndex: st.groupIndex,
			Width:      st.width,
			Height:     st.height,
			Histograms: histograms,
			Layers: map[string][]int{
				"player":     byteInts(st.player),
				"background": byteInts(st.background),
				"foreground": byteInts(st.foreground),
			},
		}
		if err := writeJSON(filepath.Join(*out, fileName), sf); err != nil {
			fatal(err)
		}

		m.Stages = append(m.Stages, stageSummary{
			Index:      i,
			Group:      st.group,
			GroupIndex: st.groupIndex,
			Width:      st.width,
			Height:     st.height,
			Histograms: histograms,
			File:       fileName,
		})
	}

	if err := writeJSON(filepath.Join(*out, "manifest.json"), m); err != nil {
		fatal(err)
	}

	fmt.Printf("decoded %d stages from %s into %s\n", len(stages), *in, *out)
	for _, st := range m.Stages {
		fmt.Printf("stage %02d: %dx%d (%s)\n", st.Index, st.Width, st.Height, st.File)
	}
}

func decodeWorld(data []byte) ([]stage, error) {
	var stages []stage
	pos := 1 // The Java loader reads and ignores the first byte before stage groups.
	group := 0
	for pos < len(data) {
		stageCount := int(data[pos])
		pos++
		if stageCount == 0 {
			return nil, fmt.Errorf("group %d has zero stages at offset %d", group, pos-1)
		}
		for groupIndex := 0; groupIndex < stageCount; groupIndex++ {
			if pos+4 > len(data) {
				return nil, io.ErrUnexpectedEOF
			}
			width := little16(data[pos], data[pos+1])
			height := little16(data[pos+2], data[pos+3])
			pos += 4
			if width <= 0 || height <= 0 {
				return nil, fmt.Errorf("invalid stage size %dx%d at group %d stage %d", width, height, group, groupIndex)
			}
			layerSize := width * height
			if pos+layerSize*3 > len(data) {
				return nil, fmt.Errorf("stage %d/%d (%dx%d) exceeds file at offset %d", group, groupIndex, width, height, pos)
			}
			st := stage{
				group:      group,
				groupIndex: groupIndex,
				width:      width,
				height:     height,
				player:     append([]byte(nil), data[pos:pos+layerSize]...),
				background: append([]byte(nil), data[pos+layerSize:pos+layerSize*2]...),
				foreground: append([]byte(nil), data[pos+layerSize*2:pos+layerSize*3]...),
			}
			pos += layerSize * 3
			stages = append(stages, st)
		}
		group++
	}
	return stages, nil
}

func little16(lo, hi byte) int {
	return int(lo) | int(hi)<<8
}

func byteInts(data []byte) []int {
	out := make([]int, len(data))
	for i, b := range data {
		out[i] = int(b)
	}
	return out
}

func histogram(data []byte) map[int]int {
	h := make(map[int]int)
	for _, b := range data {
		h[int(b)]++
	}
	return h
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

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "drdecode:", err)
	os.Exit(1)
}
