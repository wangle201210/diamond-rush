package main

import (
	"encoding/binary"
	"testing"
)

func TestDecodeSoundBank(t *testing.T) {
	headerSize := 1 + 2*8
	data := make([]byte, headerSize)
	data[0] = 2
	first := []byte("MThd-one")
	second := []byte("MThd-two-two")
	binary.LittleEndian.PutUint32(data[1:5], 0)
	binary.LittleEndian.PutUint32(data[5:9], uint32(len(first)))
	binary.LittleEndian.PutUint32(data[9:13], uint32(len(first)))
	binary.LittleEndian.PutUint32(data[13:17], uint32(len(second)))
	data = append(data, first...)
	data = append(data, second...)

	sounds, err := decodeSoundBank(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(sounds) != 2 || string(sounds[0]) != string(first) || string(sounds[1]) != string(second) {
		t.Fatalf("decoded sounds = %q", sounds)
	}
}

func TestDecodeSoundBankRejectsInvalidRange(t *testing.T) {
	data := make([]byte, 9)
	data[0] = 1
	binary.LittleEndian.PutUint32(data[5:9], 100)
	if _, err := decodeSoundBank(data); err == nil {
		t.Fatal("invalid sound range was accepted")
	}
}
