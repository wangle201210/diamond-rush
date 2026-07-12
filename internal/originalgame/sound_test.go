package originalgame

import (
	"testing"
	"time"
)

type recordingMIDIBackend struct {
	plays int
	stops int
}

func (b *recordingMIDIBackend) Play([]byte) bool {
	b.plays++
	return true
}

func (b *recordingMIDIBackend) Stop() {
	b.stops++
}

func TestOriginalMIDIBankAndDurations(t *testing.T) {
	sounds, err := loadOriginalSounds(originalAudioDir)
	if err != nil {
		t.Fatal(err)
	}
	for id, data := range sounds.bank {
		duration, ok := midiDuration(data)
		if !ok || duration <= 0 {
			t.Fatalf("sound %d duration=%s ok=%v", id, duration, ok)
		}
	}
	angkorDuration, _ := midiDuration(sounds.bank[16])
	if angkorDuration > 15*time.Second {
		t.Fatalf("Angkor intro MIDI duration=%s, unexpectedly long enough to suppress Stage 1 effects", angkorDuration)
	}
}

func TestOriginalSoundPriorityMatchesJAR(t *testing.T) {
	backend := &recordingMIDIBackend{}
	sounds := &originalSounds{backend: backend, enabled: true, currentID: -1}
	for id := range sounds.bank {
		sounds.bank[id] = []byte("MThd")
	}
	if !sounds.Play(16) {
		t.Fatal("failed to play priority-30 Angkor music")
	}
	if sounds.Play(5) || backend.plays != 1 {
		t.Fatalf("priority-10 hurt interrupted active priority-30 music: plays=%d", backend.plays)
	}
	sounds.startedAt = time.Now().Add(-time.Second)
	sounds.currentUntil = time.Now().Add(-time.Millisecond)
	if !sounds.Play(5) || backend.plays != 2 {
		t.Fatalf("expired music still blocked hurt sound: plays=%d", backend.plays)
	}
	if sounds.Play(14) || backend.plays != 2 {
		t.Fatalf("same-priority sound bypassed the source 50ms guard: plays=%d", backend.plays)
	}
	sounds.startedAt = time.Now().Add(-time.Second)
	if !sounds.Play(14) || backend.plays != 3 {
		t.Fatalf("same-priority sound did not replace after 50ms: plays=%d", backend.plays)
	}
}

func TestOriginalSoundPriorityTable(t *testing.T) {
	for _, id := range []int{1, 2, 4, 15, 16, 17, 18, 19, 20} {
		if got := originalSoundPriority(id); got != 30 {
			t.Errorf("sound %d priority=%d, want 30", id, got)
		}
	}
	for _, id := range []int{3, 7, 8, 9, 11, 12, 13} {
		if got := originalSoundPriority(id); got != 20 {
			t.Errorf("sound %d priority=%d, want 20", id, got)
		}
	}
	for _, id := range []int{0, 5, 6, 10, 14} {
		if got := originalSoundPriority(id); got != 10 {
			t.Errorf("sound %d priority=%d, want 10", id, got)
		}
	}
}

func TestAsyncMIDIBackendNeverBlocksGameThread(t *testing.T) {
	playStarted := make(chan struct{})
	releasePlay := make(chan struct{})
	stopCalled := make(chan struct{})
	backend := newAsyncMIDIBackend(func([]byte) bool {
		close(playStarted)
		<-releasePlay
		return true
	}, func() {
		close(stopCalled)
	})

	if !backend.Play([]byte("first")) {
		t.Fatal("first MIDI request was not queued")
	}
	select {
	case <-playStarted:
	case <-time.After(time.Second):
		t.Fatal("native MIDI worker did not start")
	}

	playReturned := make(chan bool, 1)
	go func() {
		playReturned <- backend.Play([]byte("replacement"))
	}()
	select {
	case queued := <-playReturned:
		if !queued {
			t.Fatal("replacement MIDI request was not queued")
		}
	case <-time.After(time.Second):
		t.Fatal("MIDI Play blocked behind native playback")
	}

	stopReturned := make(chan struct{})
	go func() {
		backend.Stop()
		close(stopReturned)
	}()
	select {
	case <-stopReturned:
	case <-time.After(time.Second):
		t.Fatal("MIDI Stop blocked behind native playback")
	}

	close(releasePlay)
	select {
	case <-backend.done:
	case <-time.After(time.Second):
		t.Fatal("native MIDI worker did not stop")
	}
	select {
	case <-stopCalled:
	default:
		t.Fatal("native MIDI stop was not called")
	}
}
