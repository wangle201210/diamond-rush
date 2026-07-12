package originalgame

import (
	"runtime"
	"sync"
)

// asyncMIDIBackend keeps native MIDI setup and teardown off the game loop.
// Only the newest pending sound matters because originalSounds has already
// applied the source priority and replacement rules.
type asyncMIDIBackend struct {
	playNative func([]byte) bool
	stopNative func()

	mu       sync.Mutex
	requests chan []byte
	stop     chan struct{}
	done     chan struct{}
	started  bool
	stopped  bool
}

func newAsyncMIDIBackend(play func([]byte) bool, stop func()) *asyncMIDIBackend {
	return &asyncMIDIBackend{
		playNative: play,
		stopNative: stop,
		requests:   make(chan []byte, 1),
		stop:       make(chan struct{}),
		done:       make(chan struct{}),
	}
}

func (b *asyncMIDIBackend) Play(data []byte) bool {
	if b == nil || len(data) == 0 || b.playNative == nil {
		return false
	}
	request := append([]byte(nil), data...)

	b.mu.Lock()
	defer b.mu.Unlock()
	if b.stopped {
		return false
	}
	if !b.started {
		b.started = true
		go b.run()
	}

	select {
	case b.requests <- request:
		return true
	default:
	}
	select {
	case <-b.requests:
	default:
	}
	select {
	case b.requests <- request:
		return true
	default:
		return false
	}
}

func (b *asyncMIDIBackend) Stop() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.stopped {
		return
	}
	b.stopped = true
	if b.started {
		close(b.stop)
	}
}

func (b *asyncMIDIBackend) run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(b.done)

	for {
		select {
		case <-b.stop:
			b.stopPlayback()
			return
		default:
		}
		select {
		case <-b.stop:
			b.stopPlayback()
			return
		case data := <-b.requests:
			b.playNative(data)
		}
	}
}

func (b *asyncMIDIBackend) stopPlayback() {
	if b.stopNative != nil {
		b.stopNative()
	}
}
