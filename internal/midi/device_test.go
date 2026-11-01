package midi

import (
	"sync"
	"testing"

	"go-jv-1080/internal/midi/message"
)

// TestLoggerDirections verifies that the registered logger sees TX for
// sent messages and RX for messages delivered through the listen path.
// It exercises logf directly (the real listen path needs a MIDI port).
func TestLoggerDirections(t *testing.T) {
	d := &Device{done: make(chan struct{})}
	var mu sync.Mutex
	var got []string
	un := d.SetLogger(func(direction string, msg []byte) {
		mu.Lock()
		got = append(got, direction+" "+message.Hex(msg))
		mu.Unlock()
	})
	d.logf("TX", []byte{0x90, 0x3C, 0x40})
	d.logf("RX", []byte{0xF0, 0x41, 0x10, 0xF7})
	un()
	d.logf("TX", []byte{0x80}) // after unregister: not logged
	mu.Lock()
	defer mu.Unlock()
	want := []string{"TX 90 3C 40", "RX F0 41 10 F7"}
	if len(got) != len(want) {
		t.Fatalf("logged %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
