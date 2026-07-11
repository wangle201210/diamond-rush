//go:build darwin && cgo

package originalgame

/*
#cgo LDFLAGS: -framework AVFoundation -framework Foundation
#include <stdint.h>

int dr_play_midi(const uint8_t *bytes, int length);
void dr_stop_midi(void);
*/
import "C"

import "unsafe"

type avMIDIBackend struct{}

func newMIDIBackend() midiBackend {
	return &avMIDIBackend{}
}

func (b *avMIDIBackend) Play(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	return C.dr_play_midi((*C.uint8_t)(unsafe.Pointer(&data[0])), C.int(len(data))) != 0
}

func (b *avMIDIBackend) Stop() {
	C.dr_stop_midi()
}
