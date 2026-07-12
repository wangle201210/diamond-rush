package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type scriptFile struct {
	Source  string   `json:"source"`
	Scripts []script `json:"scripts"`
}

type script struct {
	ID       int       `json:"id"`
	Sprites  []int     `json:"sprites,omitempty"`
	Commands []command `json:"commands"`
}

type command struct {
	Type     int       `json:"type"`
	Args     []int     `json:"args,omitempty"`
	Text     string    `json:"text,omitempty"`
	Parallel []command `json:"parallel,omitempty"`
}

type reader struct {
	data []byte
	pos  int
}

func main() {
	in := flag.String("in", "", "path to demo.f")
	out := flag.String("out", "", "output JSON path; stdout when empty")
	flag.Parse()
	if *in == "" {
		fmt.Fprintln(os.Stderr, "usage: drdemo -in path/to/demo.f [-out decoded/demo.json]")
		os.Exit(2)
	}
	data, err := os.ReadFile(*in)
	if err != nil {
		fatal(err)
	}
	chunk, err := firstChunk(data)
	if err != nil {
		fatal(err)
	}
	scripts, err := decodeScripts(chunk)
	if err != nil {
		fatal(err)
	}
	encoded, err := json.MarshalIndent(scriptFile{Source: *in, Scripts: scripts}, "", "  ")
	if err != nil {
		fatal(err)
	}
	encoded = append(encoded, '\n')
	if *out == "" {
		_, err = os.Stdout.Write(encoded)
	} else {
		err = os.WriteFile(*out, encoded, 0o644)
	}
	if err != nil {
		fatal(err)
	}
}

func firstChunk(data []byte) ([]byte, error) {
	if len(data) < 9 || data[0] == 0 {
		return nil, fmt.Errorf("invalid chunk container")
	}
	headerSize := 1 + int(data[0])*8
	if len(data) < headerSize {
		return nil, fmt.Errorf("truncated chunk header")
	}
	offset := int(binary.LittleEndian.Uint32(data[1:5]))
	length := int(binary.LittleEndian.Uint32(data[5:9]))
	start := headerSize + offset
	if start < headerSize || length < 0 || start+length > len(data) {
		return nil, fmt.Errorf("invalid first chunk offset/length %d/%d", offset, length)
	}
	return data[start : start+length], nil
}

func decodeScripts(data []byte) ([]script, error) {
	r := &reader{data: data}
	count, err := r.u16()
	if err != nil {
		return nil, err
	}
	scripts := make([]script, 0, count)
	for index := 0; index < count; index++ {
		id, err := r.u16()
		if err != nil {
			return nil, fmt.Errorf("script %d id: %w", index, err)
		}
		commandCount, err := r.u16()
		if err != nil {
			return nil, fmt.Errorf("script %d command count: %w", id, err)
		}
		length, err := r.u32()
		if err != nil {
			return nil, fmt.Errorf("script %d length: %w", id, err)
		}
		payload, err := r.bytes(length)
		if err != nil {
			return nil, fmt.Errorf("script %d payload: %w", id, err)
		}
		decoded, err := decodeScript(id, commandCount, payload)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, decoded)
	}
	return scripts, nil
}

func decodeScript(id, commandCount int, data []byte) (script, error) {
	r := &reader{data: data}
	spriteCount, err := r.u16()
	if err != nil {
		return script{}, fmt.Errorf("script %d sprite count: %w", id, err)
	}
	result := script{ID: id, Sprites: make([]int, 0, spriteCount), Commands: make([]command, 0, commandCount)}
	for index := 0; index < spriteCount; index++ {
		value, err := r.u16()
		if err != nil {
			return script{}, fmt.Errorf("script %d sprite %d: %w", id, index, err)
		}
		result.Sprites = append(result.Sprites, value)
	}
	for index := 0; index < commandCount; index++ {
		value, err := decodeCommand(r)
		if err != nil {
			return script{}, fmt.Errorf("script %d command %d: %w", id, index, err)
		}
		result.Commands = append(result.Commands, value)
	}
	if r.pos != len(r.data) {
		return script{}, fmt.Errorf("script %d left %d undecoded bytes", id, len(r.data)-r.pos)
	}
	return result, nil
}

func decodeCommand(r *reader) (command, error) {
	kind, err := r.u8()
	if err != nil {
		return command{}, err
	}
	result := command{Type: kind}
	read16 := func(count int) error {
		for index := 0; index < count; index++ {
			value, err := r.u16()
			if err != nil {
				return err
			}
			result.Args = append(result.Args, value)
		}
		return nil
	}
	switch kind {
	case 0:
		count, err := r.u8()
		if err != nil {
			return command{}, err
		}
		result.Parallel = make([]command, 0, count)
		for index := 0; index < count; index++ {
			child, err := decodeCommand(r)
			if err != nil {
				return command{}, err
			}
			result.Parallel = append(result.Parallel, child)
		}
	case 1, 9, 13:
		err = read16(3)
	case 2:
		speed, readErr := r.u8()
		if readErr != nil {
			return command{}, readErr
		}
		y, readErr := r.u16()
		if readErr != nil {
			return command{}, readErr
		}
		length, readErr := r.u16()
		if readErr != nil {
			return command{}, readErr
		}
		text, readErr := r.bytes(length)
		if readErr != nil {
			return command{}, readErr
		}
		result.Args = []int{speed, y}
		result.Text = string(text)
	case 4:
		err = read16(7)
	case 5:
		if err = read16(2); err == nil {
			var value int
			value, err = r.u8()
			result.Args = append(result.Args, value)
		}
	case 6:
		var value int
		value, err = r.u32()
		result.Args = []int{value}
	case 7, 14, 15:
	case 10:
		var value int
		value, err = r.u8()
		result.Args = []int{value}
	case 11, 12:
		err = read16(2)
	case 16, 17:
		if err = read16(1); err == nil {
			var value int
			value, err = r.u8()
			result.Args = append(result.Args, value)
		}
	case 18:
		var duration int
		duration, err = r.u8()
		if err == nil {
			var color []byte
			color, err = r.bytes(3)
			if err == nil {
				result.Args = []int{duration, int(color[0]) | int(color[1])<<8 | int(color[2])<<16}
			}
		}
	case 25:
		if err = read16(2); err == nil {
			var first, second int
			first, err = r.u8()
			if err == nil {
				second, err = r.u8()
			}
			result.Args = append(result.Args, first, second)
		}
	case 26:
		if err = read16(2); err == nil {
			var value int
			value, err = r.u32()
			result.Args = append(result.Args, value)
		}
	case 27:
		var length int
		length, err = r.u16()
		if err == nil {
			var text []byte
			text, err = r.bytes(length)
			result.Text = string(text)
		}
	default:
		return command{}, fmt.Errorf("unknown command type %d at offset %d", kind, r.pos-1)
	}
	return result, err
}

func (r *reader) u8() (int, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("unexpected EOF at %d", r.pos)
	}
	value := int(r.data[r.pos])
	r.pos++
	return value, nil
}

func (r *reader) u16() (int, error) {
	if r.pos+2 > len(r.data) {
		return 0, fmt.Errorf("unexpected EOF at %d", r.pos)
	}
	value := int(binary.LittleEndian.Uint16(r.data[r.pos : r.pos+2]))
	r.pos += 2
	return value, nil
}

func (r *reader) u32() (int, error) {
	if r.pos+4 > len(r.data) {
		return 0, fmt.Errorf("unexpected EOF at %d", r.pos)
	}
	value := int(binary.LittleEndian.Uint32(r.data[r.pos : r.pos+4]))
	r.pos += 4
	return value, nil
}

func (r *reader) bytes(length int) ([]byte, error) {
	if length < 0 || r.pos+length > len(r.data) {
		return nil, fmt.Errorf("unexpected EOF at %d reading %d bytes", r.pos, length)
	}
	value := r.data[r.pos : r.pos+length]
	r.pos += length
	return value, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "drdemo:", err)
	os.Exit(1)
}
