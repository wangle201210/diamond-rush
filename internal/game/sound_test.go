package game

import "testing"

func TestGeneratedSoundBuffersAreStereo16Bit(t *testing.T) {
	buffers := [][]byte{
		tone(440, 50, 0.2),
		sequence([]float64{440, 660}, 50, 0.2),
		noise(50, 0.2),
		templeLoop(),
	}
	for i, buffer := range buffers {
		if len(buffer) == 0 {
			t.Fatalf("buffer %d is empty", i)
		}
		if len(buffer)%4 != 0 {
			t.Fatalf("buffer %d length = %d, want stereo 16-bit frame alignment", i, len(buffer))
		}
	}
}

func TestMuteToggle(t *testing.T) {
	sounds := &Sounds{}
	if !sounds.Muted() {
		t.Fatal("empty sound system should be treated as muted")
	}
	if muted := sounds.ToggleMute(); !muted {
		t.Fatal("first toggle should mute")
	}
	if muted := sounds.ToggleMute(); muted {
		t.Fatal("second toggle should unmute")
	}
}
