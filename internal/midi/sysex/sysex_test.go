package sysex

import (
	"testing"

	"go-jv-1080/internal/midi/message"
)

func TestAreas(t *testing.T) {
	cases := []struct {
		name string
		got  message.Address
		want message.Address
	}{
		{"TemporaryPatchArea", TemporaryPatchArea(), message.Address{0x03, 0x00, 0x00, 0x00}},
		{"TemporaryPerformanceArea", TemporaryPerformanceArea(), message.Address{0x01, 0x00, 0x00, 0x00}},
		{"UserPatchArea(1)", UserPatchArea(1), message.Address{0x11, 0x00, 0x00, 0x00}},
		{"UserPatchArea(128)", UserPatchArea(128), message.Address{0x11, 0x7F, 0x00, 0x00}},
		{"UserPerformanceArea(1)", UserPerformanceArea(1), message.Address{0x10, 0x00, 0x00, 0x00}},
		{"UserPerformanceArea(32)", UserPerformanceArea(32), message.Address{0x10, 0x1F, 0x00, 0x00}},
		{"CardPatchArea(1)", CardPatchArea(1), message.Address{0x21, 0x00, 0x00, 0x00}},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %s, want %s", c.name, c.got, c.want)
		}
	}
}

func TestPatchBlocks(t *testing.T) {
	area := TemporaryPatchArea()
	addr, size := PatchCommonBlock(area)
	if addr != area || size != PatchCommonSize {
		t.Errorf("PatchCommonBlock = %s, %d", addr, size)
	}
	want := [4]message.Address{
		{0x03, 0x00, 0x10, 0x00}, {0x03, 0x00, 0x12, 0x00},
		{0x03, 0x00, 0x14, 0x00}, {0x03, 0x00, 0x16, 0x00},
	}
	for tone := 1; tone <= PatchToneCount; tone++ {
		addr, size := PatchToneBlock(area, tone)
		if addr != want[tone-1] || size != PatchToneSize {
			t.Errorf("PatchToneBlock(tone %d) = %s, %d", tone, addr, size)
		}
	}
}

func TestPerformanceBlocks(t *testing.T) {
	area := UserPerformanceArea(1)
	addr, size := PerformanceCommonBlock(area)
	if addr != area || size != PerformanceCommonSize {
		t.Errorf("PerformanceCommonBlock = %s, %d", addr, size)
	}
	// Part p lives at third address byte 0x10 + (p - 1).
	for p := 1; p <= PerformancePartCount; p++ {
		addr, size := PerformancePartBlock(area, p)
		want := message.Address{0x10, 0x00, byte(0x10 + p - 1), 0x00}
		if addr != want || size != PerformancePartSize {
			t.Errorf("PerformancePartBlock(part %d) = %s, %d; want %s, %d", p, addr, size, want, PerformancePartSize)
		}
	}
}

func TestParameterAddress(t *testing.T) {
	// Offset is in the last (fourth) address byte.
	if got := ParameterAddress(TemporaryPatchArea(), 0x2E); got != (message.Address{0x03, 0x00, 0x00, 0x2E}) {
		t.Errorf("ParameterAddress = %s", got)
	}
}

func TestPatchParameterAddress(t *testing.T) {
	area := UserPatchArea(5)
	// Common parameter: patch level at offset 0x2E (fourth address byte).
	addr, err := PatchParameterAddress(area, "common", 0, 0x2E)
	if err != nil || addr != (message.Address{0x11, 0x04, 0x00, 0x2E}) {
		t.Errorf("common = %s, %v", addr, err)
	}
	// Tone parameter: tone 1 switch at offset 0x00 of tone 1.
	addr, err = PatchParameterAddress(area, "tone", 1, 0x00)
	if err != nil || addr != (message.Address{0x11, 0x04, 0x10, 0x00}) {
		t.Errorf("tone 1 = %s, %v", addr, err)
	}
	// Tone parameter: tone 4 cutoff at offset 0x51 of tone 4.
	addr, err = PatchParameterAddress(area, "tone", 4, 0x51)
	if err != nil || addr != (message.Address{0x11, 0x04, 0x16, 0x51}) {
		t.Errorf("tone 4 = %s, %v", addr, err)
	}
	if _, err := PatchParameterAddress(area, "tone", 5, 0); err == nil {
		t.Error("tone 5 should fail")
	}
	if _, err := PatchParameterAddress(area, "bogus", 0, 0); err == nil {
		t.Error("unknown block should fail")
	}
}

func TestPerformanceParameterAddress(t *testing.T) {
	area := UserPerformanceArea(2)
	// Common parameter: part 1 switch at offset 0x0D (fourth address byte).
	addr, err := PerformanceParameterAddress(area, "common", 0, 0x0D)
	if err != nil || addr != (message.Address{0x10, 0x01, 0x00, 0x0D}) {
		t.Errorf("common = %s, %v", addr, err)
	}
	// Part parameter: part 1 channel at offset 0x01 of part 1.
	addr, err = PerformanceParameterAddress(area, "part", 1, 0x01)
	if err != nil || addr != (message.Address{0x10, 0x01, 0x10, 0x01}) {
		t.Errorf("part 1 = %s, %v", addr, err)
	}
	// Part parameter: part 16 pan at offset 0x07 of part 16.
	addr, err = PerformanceParameterAddress(area, "part", 16, 0x07)
	if err != nil || addr != (message.Address{0x10, 0x01, 0x1F, 0x07}) {
		t.Errorf("part 16 = %s, %v", addr, err)
	}
	if _, err := PerformanceParameterAddress(area, "part", 17, 0); err == nil {
		t.Error("part 17 should fail")
	}
	if _, err := PerformanceParameterAddress(area, "bogus", 0, 0); err == nil {
		t.Error("unknown block should fail")
	}
}
