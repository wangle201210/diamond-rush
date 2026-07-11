package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type demoScript struct {
	ID        int           `json:"id"`
	Sprites   []int         `json:"sprites,omitempty"`
	Commands  []demoCommand `json:"commands"`
	RawLength int           `json:"rawLength"`
}

type demoCommand struct {
	Opcode   int           `json:"opcode"`
	Name     string        `json:"name"`
	X        *int          `json:"x,omitempty"`
	Y        *int          `json:"y,omitempty"`
	Duration *int          `json:"duration,omitempty"`
	Value    *int          `json:"value,omitempty"`
	Frame    *int          `json:"frame,omitempty"`
	Sprite   *int          `json:"sprite,omitempty"`
	Color    *int          `json:"color,omitempty"`
	Side     *int          `json:"side,omitempty"`
	Text     string        `json:"text,omitempty"`
	Commands []demoCommand `json:"commands,omitempty"`
}

type byteReader struct {
	data []byte
	pos  int
}

func main() {
	in := flag.String("in", "/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/demo.f", "original demo.f path")
	out := flag.String("out", "", "optional JSON output path")
	flag.Parse()

	scripts, err := decodeDemoFile(*in)
	if err != nil {
		fatal(err)
	}
	data, err := json.MarshalIndent(scripts, "", "  ")
	if err != nil {
		fatal(err)
	}
	data = append(data, '\n')
	if *out == "" {
		_, err = os.Stdout.Write(data)
	} else {
		err = os.WriteFile(filepath.Clean(*out), data, 0o644)
	}
	if err != nil {
		fatal(err)
	}
}

func decodeDemoFile(path string) ([]demoScript, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	chunk, err := resourceChunk(data, 0)
	if err != nil {
		return nil, err
	}
	r := byteReader{data: chunk}
	count, err := r.u16()
	if err != nil {
		return nil, fmt.Errorf("read script count: %w", err)
	}
	scripts := make([]demoScript, 0, count)
	for index := 0; index < count; index++ {
		id, err := r.u16()
		if err != nil {
			return nil, fmt.Errorf("read script %d id: %w", index, err)
		}
		commandCount, err := r.u16()
		if err != nil {
			return nil, fmt.Errorf("read script %d command count: %w", id, err)
		}
		length, err := r.u32()
		if err != nil {
			return nil, fmt.Errorf("read script %d length: %w", id, err)
		}
		payload, err := r.bytes(length)
		if err != nil {
			return nil, fmt.Errorf("read script %d payload: %w", id, err)
		}
		script, err := decodeDemoScript(id, commandCount, payload)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}
	return scripts, nil
}

func resourceChunk(data []byte, index int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty resource")
	}
	count := int(data[0])
	if index < 0 || index >= count {
		return nil, fmt.Errorf("chunk %d outside 0..%d", index, count-1)
	}
	headerLength := 1 + count*8
	if len(data) < headerLength {
		return nil, fmt.Errorf("resource header exceeds file")
	}
	entry := 1 + index*8
	offset := int(binary.LittleEndian.Uint32(data[entry : entry+4]))
	length := int(binary.LittleEndian.Uint32(data[entry+4 : entry+8]))
	start := headerLength + offset
	end := start + length
	if length < 0 || start < headerLength || end > len(data) {
		return nil, fmt.Errorf("invalid chunk %d range %d..%d", index, start, end)
	}
	return data[start:end], nil
}

func decodeDemoScript(id, commandCount int, payload []byte) (demoScript, error) {
	r := byteReader{data: payload}
	spriteCount, err := r.u16()
	if err != nil {
		return demoScript{}, fmt.Errorf("decode script %d sprites: %w", id, err)
	}
	sprites := make([]int, spriteCount)
	for index := range sprites {
		sprites[index], err = r.u16()
		if err != nil {
			return demoScript{}, fmt.Errorf("decode script %d sprite %d: %w", id, index, err)
		}
	}
	commands, err := parseDemoCommands(&r, commandCount)
	if err != nil {
		return demoScript{}, fmt.Errorf("decode script %d command %d: %w", id, len(commands), err)
	}
	if r.pos != len(payload) {
		return demoScript{}, fmt.Errorf("decode script %d left %d unread bytes", id, len(payload)-r.pos)
	}
	return demoScript{ID: id, Sprites: sprites, Commands: commands, RawLength: len(payload)}, nil
}

func parseDemoCommands(r *byteReader, count int) ([]demoCommand, error) {
	commands := make([]demoCommand, 0, count)
	for index := 0; index < count; index++ {
		opcode, err := r.u8()
		if err != nil {
			return commands, err
		}
		command := demoCommand{Opcode: opcode, Name: demoOpcodeName(opcode)}
		switch opcode {
		case 0:
			parallelCount, err := r.u8()
			if err != nil {
				return commands, err
			}
			command.Commands, err = parseDemoCommands(r, parallelCount)
			if err != nil {
				return commands, err
			}
		case 1:
			x, y, duration, err := r.u16x3()
			if err != nil {
				return commands, err
			}
			command.X, command.Y, command.Duration = intPtr(x), intPtr(y), intPtr(duration)
		case 2:
			side, err := r.u8()
			if err != nil {
				return commands, err
			}
			y, err := r.u16()
			if err != nil {
				return commands, err
			}
			text, err := r.text()
			if err != nil {
				return commands, err
			}
			command.Side, command.Y, command.Text = intPtr(side), intPtr(y), text
		case 4:
			values := make([]int, 7)
			for value := range values {
				values[value], err = r.u16()
				if err != nil {
					return commands, err
				}
			}
			command.X, command.Y = intPtr(values[0]), intPtr(values[1])
			command.Value = intPtr(values[2])
			command.Frame, command.Sprite, command.Duration = intPtr(values[4]), intPtr(values[5]), intPtr(values[6])
		case 5:
			x, y, err := r.u16x2()
			if err != nil {
				return commands, err
			}
			value, err := r.u8()
			if err != nil {
				return commands, err
			}
			command.X, command.Y, command.Value = intPtr(x), intPtr(y), intPtr(value)
		case 6:
			value, err := r.u32()
			if err != nil {
				return commands, err
			}
			command.Duration = intPtr(value)
		case 7, 14, 15:
		case 9:
			x, y, value, err := r.u16x3()
			if err != nil {
				return commands, err
			}
			command.X, command.Y, command.Value = intPtr(x), intPtr(y), intPtr(value)
		case 10:
			value, err := r.u8()
			if err != nil {
				return commands, err
			}
			command.Value = intPtr(value)
		case 11, 12:
			first, second, err := r.u16x2()
			if err != nil {
				return commands, err
			}
			if opcode == 11 {
				command.Frame, command.Sprite = intPtr(first), intPtr(second)
			} else {
				command.X, command.Y = intPtr(first), intPtr(second)
			}
		case 13:
			x, y, duration, err := r.u16x3()
			if err != nil {
				return commands, err
			}
			command.X, command.Y, command.Duration = intPtr(x), intPtr(y), intPtr(duration)
		case 16, 17:
			frame, err := r.u16()
			if err != nil {
				return commands, err
			}
			duration, err := r.u8()
			if err != nil {
				return commands, err
			}
			command.Frame, command.Duration = intPtr(frame), intPtr(duration)
		case 18:
			duration, err := r.u8()
			if err != nil {
				return commands, err
			}
			colorValue, err := r.u24()
			if err != nil {
				return commands, err
			}
			command.Duration, command.Color = intPtr(duration), intPtr(colorValue)
		case 25:
			x, y, err := r.u16x2()
			if err != nil {
				return commands, err
			}
			value, err := r.u16()
			if err != nil {
				return commands, err
			}
			command.X, command.Y, command.Value = intPtr(x), intPtr(y), intPtr(value)
		case 26:
			x, y, err := r.u16x2()
			if err != nil {
				return commands, err
			}
			value, err := r.u32()
			if err != nil {
				return commands, err
			}
			command.X, command.Y, command.Value = intPtr(x), intPtr(y), intPtr(value)
		case 27:
			text, err := r.text()
			if err != nil {
				return commands, err
			}
			command.Text = text
		default:
			return commands, fmt.Errorf("unsupported opcode %d at byte %d", opcode, r.pos-1)
		}
		commands = append(commands, command)
	}
	return commands, nil
}

func demoOpcodeName(opcode int) string {
	names := map[int]string{
		0: "parallel", 1: "camera", 2: "text_bubble", 4: "sprite_move",
		5: "stage_effect", 6: "wait", 7: "wait_for_input", 9: "stage_animation",
		10: "move_player", 11: "portrait_face", 12: "portrait_position",
		13: "portrait_move", 14: "portrait_show", 15: "portrait_hide",
		16: "portrait_mark", 17: "portrait_mark_temporary", 18: "flash",
		25: "set_foreground", 26: "set_object_state", 27: "text_bottom",
	}
	if name := names[opcode]; name != "" {
		return name
	}
	return "unknown"
}

func (r *byteReader) u8() (int, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("read byte at %d beyond %d", r.pos, len(r.data))
	}
	value := int(r.data[r.pos])
	r.pos++
	return value, nil
}

func (r *byteReader) u16() (int, error) {
	data, err := r.bytes(2)
	if err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint16(data)), nil
}

func (r *byteReader) u24() (int, error) {
	data, err := r.bytes(3)
	if err != nil {
		return 0, err
	}
	return int(data[0]) | int(data[1])<<8 | int(data[2])<<16, nil
}

func (r *byteReader) u32() (int, error) {
	data, err := r.bytes(4)
	if err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint32(data)), nil
}

func (r *byteReader) u16x2() (int, int, error) {
	first, err := r.u16()
	if err != nil {
		return 0, 0, err
	}
	second, err := r.u16()
	return first, second, err
}

func (r *byteReader) u16x3() (int, int, int, error) {
	first, second, err := r.u16x2()
	if err != nil {
		return 0, 0, 0, err
	}
	third, err := r.u16()
	return first, second, third, err
}

func (r *byteReader) text() (string, error) {
	length, err := r.u16()
	if err != nil {
		return "", err
	}
	data, err := r.bytes(length)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *byteReader) bytes(count int) ([]byte, error) {
	if count < 0 || r.pos+count > len(r.data) {
		return nil, fmt.Errorf("read %d bytes at %d beyond %d", count, r.pos, len(r.data))
	}
	data := r.data[r.pos : r.pos+count]
	r.pos += count
	return data, nil
}

func intPtr(value int) *int {
	return &value
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
