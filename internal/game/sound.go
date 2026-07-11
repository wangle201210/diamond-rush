package game

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

type soundID string

const (
	soundStep     soundID = "step"
	soundDiamond  soundID = "diamond"
	soundKey      soundID = "key"
	soundDoor     soundID = "door"
	soundSwitch   soundID = "switch"
	soundTeleport soundID = "teleport"
	soundBreak    soundID = "break"
	soundDeath    soundID = "death"
	soundWin      soundID = "win"
	soundChest    soundID = "chest"
)

type Sounds struct {
	context     *audio.Context
	effects     map[soundID][]byte
	music       []byte
	musicPlayer *audio.Player
	muted       bool
}

func NewSounds() *Sounds {
	s := &Sounds{
		context: audio.NewContext(sampleRate),
		effects: make(map[soundID][]byte),
	}
	s.effects[soundStep] = tone(120, 35, 0.18)
	s.effects[soundDiamond] = sequence([]float64{640, 880}, 42, 0.24)
	s.effects[soundKey] = sequence([]float64{520, 1040}, 55, 0.22)
	s.effects[soundDoor] = sequence([]float64{180, 140}, 80, 0.24)
	s.effects[soundSwitch] = sequence([]float64{360, 520}, 45, 0.22)
	s.effects[soundTeleport] = sequence([]float64{760, 540, 920}, 36, 0.20)
	s.effects[soundBreak] = noise(120, 0.18)
	s.effects[soundDeath] = sequence([]float64{220, 150, 90}, 80, 0.25)
	s.effects[soundWin] = sequence([]float64{520, 660, 880, 1040}, 70, 0.22)
	s.effects[soundChest] = sequence([]float64{420, 630, 840, 1260}, 45, 0.24)
	s.music = templeLoop()
	s.startMusic()
	return s
}

func (s *Sounds) Play(id soundID) {
	if s == nil || s.context == nil || s.muted {
		return
	}
	data := s.effects[id]
	if len(data) == 0 {
		return
	}
	player := audio.NewPlayerFromBytes(s.context, data)
	player.Play()
}

func (s *Sounds) ToggleMute() bool {
	if s == nil {
		return true
	}
	s.muted = !s.muted
	if s.musicPlayer != nil {
		if s.muted {
			s.musicPlayer.Pause()
		} else {
			s.musicPlayer.Play()
		}
	}
	return s.muted
}

func (s *Sounds) Muted() bool {
	return s == nil || s.context == nil || s.muted
}

func (s *Sounds) startMusic() {
	if s == nil || s.context == nil || len(s.music) == 0 {
		return
	}
	loop := audio.NewInfiniteLoop(bytes.NewReader(s.music), int64(len(s.music)))
	player, err := audio.NewPlayer(s.context, loop)
	if err != nil {
		return
	}
	player.SetVolume(0.18)
	player.Play()
	s.musicPlayer = player
}

func templeLoop() []byte {
	notes := []float64{146.83, 0, 196.00, 0, 174.61, 0, 220.00, 0, 146.83, 164.81, 196.00, 0, 130.81, 0, 174.61, 0}
	var out []byte
	for i, note := range notes {
		out = append(out, musicBeat(note, 260, i)...)
	}
	return out
}

func musicBeat(freq float64, milliseconds int, step int) []byte {
	samples := sampleRate * milliseconds / 1000
	out := make([]byte, samples*4)
	drone := 73.42
	for i := 0; i < samples; i++ {
		t := float64(i) / sampleRate
		pulse := 0.55 + 0.45*math.Sin(2*math.Pi*float64(step+1)*t/2)
		value := math.Sin(2*math.Pi*drone*t) * 0.13
		if freq > 0 {
			value += math.Sin(2*math.Pi*freq*t) * 0.10 * pulse
			value += math.Sin(2*math.Pi*freq*2*t) * 0.035 * pulse
		}
		env := envelope(i, samples)
		sample := int16(value * env * math.MaxInt16)
		binary.LittleEndian.PutUint16(out[i*4:], uint16(sample))
		binary.LittleEndian.PutUint16(out[i*4+2:], uint16(sample))
	}
	return out
}

func tone(freq float64, milliseconds int, volume float64) []byte {
	samples := sampleRate * milliseconds / 1000
	out := make([]byte, samples*4)
	for i := 0; i < samples; i++ {
		t := float64(i) / sampleRate
		env := envelope(i, samples)
		value := int16(math.Sin(2*math.Pi*freq*t) * volume * env * math.MaxInt16)
		binary.LittleEndian.PutUint16(out[i*4:], uint16(value))
		binary.LittleEndian.PutUint16(out[i*4+2:], uint16(value))
	}
	return out
}

func sequence(freqs []float64, milliseconds int, volume float64) []byte {
	var out []byte
	for _, freq := range freqs {
		out = append(out, tone(freq, milliseconds, volume)...)
	}
	return out
}

func noise(milliseconds int, volume float64) []byte {
	samples := sampleRate * milliseconds / 1000
	out := make([]byte, samples*4)
	seed := uint32(1)
	for i := 0; i < samples; i++ {
		seed = seed*1664525 + 1013904223
		raw := int16((seed>>16)&0xffff) - math.MaxInt16
		value := int16(float64(raw) * volume * envelope(i, samples))
		binary.LittleEndian.PutUint16(out[i*4:], uint16(value))
		binary.LittleEndian.PutUint16(out[i*4+2:], uint16(value))
	}
	return out
}

func envelope(i, samples int) float64 {
	if samples <= 1 {
		return 1
	}
	attack := samples / 5
	release := samples / 3
	if attack > 0 && i < attack {
		return float64(i) / float64(attack)
	}
	if release > 0 && i > samples-release {
		return float64(samples-i) / float64(release)
	}
	return 1
}
