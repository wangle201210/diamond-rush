//go:build !darwin || !cgo

package originalgame

type silentMIDIBackend struct{}

func newMIDIBackend() midiBackend {
	return &silentMIDIBackend{}
}

func (b *silentMIDIBackend) Play([]byte) bool {
	return false
}

func (b *silentMIDIBackend) Stop() {}
