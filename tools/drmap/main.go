package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type point [2]int

type node struct {
	X     int     `json:"x"`
	Y     int     `json:"y"`
	Type  int     `json:"type"`
	Stage int     `json:"stage"`
	Links []point `json:"links"`
}

type worldMap struct {
	Source        string `json:"source"`
	PayloadLength int    `json:"payload_length"`
	Nodes         []node `json:"nodes"`
}

func main() {
	input := flag.String("in", "", "map_*.out input")
	output := flag.String("out", "", "JSON output")
	flag.Parse()
	if *input == "" || *output == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := decode(*input, *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func decode(input, output string) error {
	raw, err := os.ReadFile(filepath.Clean(input))
	if err != nil {
		return fmt.Errorf("read map: %w", err)
	}
	if len(raw) < 3 {
		return fmt.Errorf("map has %d bytes, want at least 3", len(raw))
	}
	payloadLength := int(binary.LittleEndian.Uint16(raw[:2]))
	if payloadLength != len(raw)-3 {
		return fmt.Errorf("payload length is %d, file contains %d", payloadLength, len(raw)-3)
	}
	nodeCount := int(raw[2])
	cursor := 3
	result := worldMap{
		Source:        filepath.Base(input),
		PayloadLength: payloadLength,
		Nodes:         make([]node, 0, nodeCount),
	}
	for index := 0; index < nodeCount; index++ {
		if len(raw)-cursor < 5 {
			return fmt.Errorf("node %d header is truncated", index)
		}
		header := raw[cursor : cursor+5]
		cursor += 5
		linkCount := int(header[4])
		if len(raw)-cursor < linkCount*2 {
			return fmt.Errorf("node %d links are truncated", index)
		}
		n := node{
			X:     int(header[0]),
			Y:     int(header[1]),
			Type:  int(header[2]),
			Stage: int(header[3]),
			Links: make([]point, linkCount),
		}
		for link := range n.Links {
			n.Links[link] = point{int(raw[cursor]), int(raw[cursor+1])}
			cursor += 2
		}
		result.Nodes = append(result.Nodes, n)
	}
	if cursor != len(raw) {
		return fmt.Errorf("map has %d trailing payload bytes", len(raw)-cursor)
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("encode map: %w", err)
	}
	encoded = append(encoded, '\n')
	if err := os.WriteFile(filepath.Clean(output), encoded, 0o644); err != nil {
		return fmt.Errorf("write map: %w", err)
	}
	return nil
}
