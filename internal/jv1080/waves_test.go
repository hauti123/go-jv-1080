package jv1080

import (
	"strings"
	"testing"
)

func TestWaveSources(t *testing.T) {
	want := []struct {
		label   string
		group   int
		groupID int
	}{
		{"INT-A", 0, 1}, {"INT-B", 0, 2},
		{"EXP-A", 2, 1}, {"EXP-B", 2, 2}, {"EXP-C", 2, 3}, {"EXP-D", 2, 4},
	}
	if len(WaveSources) != 6 {
		t.Fatalf("WaveSources has %d entries, want 6", len(WaveSources))
	}
	for i, w := range want {
		s := WaveSources[i]
		if s.Label != w.label || s.Group != w.group || s.GroupID != w.groupID {
			t.Errorf("WaveSources[%d] = %+v, want %+v", i, s, w)
		}
		if idx := WaveSourceIndex(w.group, w.groupID); idx != i {
			t.Errorf("WaveSourceIndex(%d,%d) = %d, want %d", w.group, w.groupID, idx, i)
		}
	}
	if idx := WaveSourceIndex(1, 0); idx != -1 { // PCM card is not selectable
		t.Errorf("WaveSourceIndex(1,0) = %d, want -1", idx)
	}
}

func TestWaveCount(t *testing.T) {
	if got := WaveCount(0, 1); got != 255 {
		t.Errorf("WaveCount(INT-A) = %d, want 255", got)
	}
	if got := WaveCount(0, 2); got != 193 {
		t.Errorf("WaveCount(INT-B) = %d, want 193", got)
	}
	if got := WaveCount(2, 1); got != 255 {
		t.Errorf("WaveCount(EXP-A) = %d, want 255", got)
	}
}

func TestWaveName(t *testing.T) {
	if name, ok := WaveName(0, 1, 1); !ok || name != "Ac Piano1 A" {
		t.Errorf("WaveName(INT-A, 1) = %q, %v; want \"Ac Piano1 A\", true", name, ok)
	}
	if name, ok := WaveName(0, 1, 255); !ok || name != "Vox Noise" {
		t.Errorf("WaveName(INT-A, 255) = %q, %v; want \"Vox Noise\", true", name, ok)
	}
	if name, ok := WaveName(0, 2, 193); !ok || name != "DC" {
		t.Errorf("WaveName(INT-B, 193) = %q, %v; want \"DC\", true", name, ok)
	}
	if _, ok := WaveName(0, 1, 256); ok {
		t.Errorf("WaveName(INT-A, 256) ok = true, want false")
	}
	if _, ok := WaveName(2, 1, 5); ok { // EXP has no documented names
		t.Errorf("WaveName(EXP-A, 5) ok = true, want false")
	}
}

func TestWaveOptions(t *testing.T) {
	opts := WaveOptions(0, 1)
	if len(opts) != 255 {
		t.Fatalf("WaveOptions(INT-A) has %d entries, want 255", len(opts))
	}
	if opts[0] != "001 Ac Piano1 A" {
		t.Errorf("opts[0] = %q, want \"001 Ac Piano1 A\"", opts[0])
	}
	if opts[254] != "255 Vox Noise" {
		t.Errorf("opts[254] = %q, want \"255 Vox Noise\"", opts[254])
	}
	exp := WaveOptions(2, 3)
	if len(exp) != 255 || exp[0] != "001" {
		t.Errorf("WaveOptions(EXP-C) = %d entries starting %q, want 255 entries starting \"001\"", len(exp), exp[0])
	}
	for _, o := range opts {
		if !strings.HasPrefix(o, "0") && !strings.HasPrefix(o, "1") && !strings.HasPrefix(o, "2") {
			t.Errorf("option %q does not start with a zero-padded number", o)
		}
	}
}
