package originalgame

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const originalSoundCount = 21

type midiBackend interface {
	Play([]byte) bool
	Stop()
}

type originalSounds struct {
	bank            [originalSoundCount][]byte
	backend         midiBackend
	enabled         bool
	currentID       int
	currentPriority int
	startedAt       time.Time
	currentUntil    time.Time
}

func loadOriginalSounds(dir string) (*originalSounds, error) {
	sounds := &originalSounds{backend: newMIDIBackend(), currentID: -1}
	for id := range sounds.bank {
		path := filepath.Join(resolvePath(dir), fmt.Sprintf("sound%02d.mid", id))
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
		if len(data) < 14 || string(data[:4]) != "MThd" {
			return nil, fmt.Errorf("invalid original MIDI sound %d: %s", id, path)
		}
		sounds.bank[id] = data
	}
	return sounds, nil
}

func (s *originalSounds) Enable() {
	if s != nil {
		s.enabled = true
	}
}

func (s *originalSounds) Play(id int) bool {
	if s == nil || !s.enabled || s.backend == nil || id < 0 || id >= len(s.bank) {
		return false
	}
	now := time.Now()
	if s.currentID >= 0 && !s.currentUntil.IsZero() && !now.Before(s.currentUntil) {
		s.currentID = -1
		s.currentPriority = 0
	}
	priority := originalSoundPriority(id)
	if s.currentID >= 0 && (s.currentPriority > priority || (s.currentPriority == priority && now.Sub(s.startedAt) <= 50*time.Millisecond)) {
		return false
	}
	if !s.backend.Play(s.bank[id]) {
		return false
	}
	s.currentID = id
	s.currentPriority = priority
	s.startedAt = now
	duration, ok := midiDuration(s.bank[id])
	if ok {
		s.currentUntil = now.Add(duration)
	} else {
		s.currentUntil = time.Time{}
	}
	return true
}

func (s *originalSounds) Stop() {
	if s == nil || s.backend == nil {
		return
	}
	s.backend.Stop()
	s.currentID = -1
	s.currentPriority = 0
	s.currentUntil = time.Time{}
}

func originalSoundPriority(id int) int {
	switch id {
	case 1, 2, 4, 15, 16, 17, 18, 19, 20:
		return 30
	case 3, 7, 8, 9, 11, 12, 13:
		return 20
	case 0, 5, 6, 10, 14:
		return 10
	default:
		return 0
	}
}

func midiDuration(data []byte) (time.Duration, bool) {
	if len(data) < 14 || string(data[:4]) != "MThd" {
		return 0, false
	}
	headerLength := int(binary.BigEndian.Uint32(data[4:8]))
	if headerLength < 6 || 8+headerLength > len(data) {
		return 0, false
	}
	division := int(binary.BigEndian.Uint16(data[12:14]))
	if division <= 0 || division&0x8000 != 0 {
		return 0, false
	}
	position := 8 + headerLength
	if position+8 > len(data) || string(data[position:position+4]) != "MTrk" {
		return 0, false
	}
	trackLength := int(binary.BigEndian.Uint32(data[position+4 : position+8]))
	position += 8
	end := position + trackLength
	if trackLength < 0 || end > len(data) {
		return 0, false
	}
	tempo := int64(500000)
	microseconds := int64(0)
	runningStatus := byte(0)
	for position < end {
		delta, next, ok := readVariableLength(data, position, end)
		if !ok {
			return 0, false
		}
		position = next
		microseconds += int64(delta) * tempo / int64(division)
		if position >= end {
			return 0, false
		}
		status := data[position]
		if status&0x80 != 0 {
			position++
			if status < 0xf0 {
				runningStatus = status
			}
		} else {
			if runningStatus == 0 {
				return 0, false
			}
			status = runningStatus
		}
		switch {
		case status == 0xff:
			if position >= end {
				return 0, false
			}
			typeByte := data[position]
			position++
			length, next, ok := readVariableLength(data, position, end)
			if !ok || next+length > end {
				return 0, false
			}
			position = next
			if typeByte == 0x51 && length == 3 {
				tempo = int64(data[position])<<16 | int64(data[position+1])<<8 | int64(data[position+2])
			}
			position += length
		case status == 0xf0 || status == 0xf7:
			length, next, ok := readVariableLength(data, position, end)
			if !ok || next+length > end {
				return 0, false
			}
			position = next + length
		case status&0xf0 == 0xc0 || status&0xf0 == 0xd0:
			position++
		default:
			position += 2
		}
		if position > end {
			return 0, false
		}
	}
	return time.Duration(microseconds) * time.Microsecond, microseconds > 0
}

func readVariableLength(data []byte, position, end int) (int, int, bool) {
	value := 0
	for count := 0; count < 4 && position < end; count++ {
		b := data[position]
		position++
		value = value<<7 | int(b&0x7f)
		if b&0x80 == 0 {
			return value, position, true
		}
	}
	return 0, position, false
}
