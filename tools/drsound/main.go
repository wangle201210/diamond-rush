package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type soundEntry struct {
	ID     int    `json:"id"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	File   string `json:"file"`
}

type soundManifest struct {
	Source string       `json:"source"`
	Count  int          `json:"count"`
	Sounds []soundEntry `json:"sounds"`
}

func main() {
	in := flag.String("in", "", "path to Diamond Rush snd.f")
	out := flag.String("out", "decoded/audio", "output directory")
	flag.Parse()
	if *in == "" {
		fatal(fmt.Errorf("-in is required"))
	}
	data, err := os.ReadFile(filepath.Clean(*in))
	if err != nil {
		fatal(err)
	}
	sounds, err := decodeSoundBank(data)
	if err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fatal(err)
	}
	manifest := soundManifest{Source: *in, Count: len(sounds)}
	headerSize := 1 + len(sounds)*8
	for id, sound := range sounds {
		name := fmt.Sprintf("sound%02d.mid", id)
		if err := os.WriteFile(filepath.Join(*out, name), sound, 0o644); err != nil {
			fatal(err)
		}
		offset := int(binary.LittleEndian.Uint32(data[1+id*8:]))
		manifest.Sounds = append(manifest.Sounds, soundEntry{ID: id, Offset: headerSize + offset, Length: len(sound), File: name})
	}
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fatal(err)
	}
	manifestData = append(manifestData, '\n')
	if err := os.WriteFile(filepath.Join(*out, "manifest.json"), manifestData, 0o644); err != nil {
		fatal(err)
	}
}

func decodeSoundBank(data []byte) ([][]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty sound bank")
	}
	count := int(data[0])
	headerSize := 1 + count*8
	if count <= 0 || headerSize > len(data) {
		return nil, fmt.Errorf("invalid sound count %d for %d-byte bank", count, len(data))
	}
	sounds := make([][]byte, count)
	for id := 0; id < count; id++ {
		header := data[1+id*8 : 1+(id+1)*8]
		offset := int(binary.LittleEndian.Uint32(header[0:4]))
		length := int(binary.LittleEndian.Uint32(header[4:8]))
		start := headerSize + offset
		end := start + length
		if length <= 0 || start < headerSize || end < start || end > len(data) {
			return nil, fmt.Errorf("sound %d range offset=%d length=%d exceeds %d-byte bank", id, offset, length, len(data))
		}
		if length < 4 || string(data[start:start+4]) != "MThd" {
			return nil, fmt.Errorf("sound %d is not a standard MIDI file", id)
		}
		sounds[id] = append([]byte(nil), data[start:end]...)
	}
	return sounds, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
