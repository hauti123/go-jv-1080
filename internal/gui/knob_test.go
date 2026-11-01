package gui

import (
	"math"
	"testing"
)

func TestScaleAngles(t *testing.T) {
	// 7 o'clock = -135 deg, 5 o'clock = +135 deg, 12 o'clock = 0 deg
	if got := angleOf(0, 127, 0); math.Abs(got-MinAngle) > 1e-9 {
		t.Errorf("angleOf(min) = %v, want %v (7 o'clock)", got, MinAngle)
	}
	if got := angleOf(0, 127, 127); math.Abs(got-MaxAngle) > 1e-9 {
		t.Errorf("angleOf(max) = %v, want %v (5 o'clock)", got, MaxAngle)
	}
	if got := angleOf(0, 126, 63); math.Abs(got) > 1e-9 {
		t.Errorf("angleOf(mid) = %v, want 0 (12 o'clock)", got)
	}
}

func TestKnobAngle(t *testing.T) {
	k := NewKnob(KnobMinMax, "LEVEL", 0, 127)
	if k.angle() != MinAngle {
		t.Errorf("angle at min = %v, want %v", k.angle(), MinAngle)
	}
	k.Set(127)
	if k.angle() != MaxAngle {
		t.Errorf("angle at max = %v, want %v", k.angle(), MaxAngle)
	}
	k.Set(64)
	want := MinAngle + (MaxAngle-MinAngle)*64.0/127.0
	if math.Abs(k.angle()-want) > 1e-9 {
		t.Errorf("angle at 64 = %v, want %v", k.angle(), want)
	}
}

func TestKnobClamp(t *testing.T) {
	k := NewKnob(KnobMinMax, "LEVEL", 0, 127)
	k.Set(200)
	if k.Value() != 127 {
		t.Errorf("Set(200) = %d, want 127", k.Value())
	}
	k.Set(-5)
	if k.Value() != 0 {
		t.Errorf("Set(-5) = %d, want 0", k.Value())
	}
}

func TestKnobFormatMinMax(t *testing.T) {
	k := NewKnob(KnobMinMax, "LEVEL", 0, 127)
	k.Set(0)
	if got := k.formatValue(); got != "0" {
		t.Errorf("formatValue at 0 = %q, want %q", got, "0")
	}
	k.Set(127)
	if got := k.formatValue(); got != "127" {
		t.Errorf("formatValue at 127 = %q, want %q", got, "127")
	}
}

func TestKnobFormatBipolar(t *testing.T) {
	k := NewKnob(KnobBipolar, "OCTAVE", -3, 3)
	cases := map[int]string{-3: "-3", -1: "-1", 0: "+0", 1: "+1", 3: "+3"}
	for v, want := range cases {
		k.Set(v)
		if got := k.formatValue(); got != want {
			t.Errorf("bipolar formatValue(%d) = %q, want %q", v, got, want)
		}
	}
	// 0 must map to 12 o'clock
	if k.center() != 0 {
		t.Errorf("center() = %d, want 0", k.center())
	}
	k.Set(0)
	if math.Abs(k.angle()) > 1e-9 {
		t.Errorf("bipolar angle at 0 = %v, want 0 (12 o'clock)", k.angle())
	}
	// positive to the right, negative to the left
	k.Set(3)
	if k.angle() != MaxAngle {
		t.Errorf("bipolar angle at +3 = %v, want %v", k.angle(), MaxAngle)
	}
	k.Set(-3)
	if k.angle() != MinAngle {
		t.Errorf("bipolar angle at -3 = %v, want %v", k.angle(), MinAngle)
	}
}

func TestKnobFormatPan(t *testing.T) {
	k := NewKnob(KnobPan, "PAN", 0, 127)
	cases := map[int]string{0: "L64", 1: "L63", 63: "L1", 64: "C", 65: "R1", 126: "R62", 127: "R63"}
	for v, want := range cases {
		k.Set(v)
		if got := k.formatValue(); got != want {
			t.Errorf("pan formatValue(%d) = %q, want %q", v, got, want)
		}
	}
	// C must map to 12 o'clock area (64 of 0..127)
	k.Set(64)
	if math.Abs(k.angle()-angleOf(0, 127, 64)) > 1e-9 {
		t.Errorf("pan angle at C = %v, want %v", k.angle(), angleOf(0, 127, 64))
	}
}

func TestKnobSetReportsChange(t *testing.T) {
	k := NewKnob(KnobMinMax, "LEVEL", 0, 127)
	var got []int
	k.OnChange = func(v int) { got = append(got, v) }
	k.set(5)
	k.set(5) // no change, no callback
	k.set(1000)
	if len(got) != 2 || got[0] != 5 || got[1] != 127 {
		t.Errorf("OnChange calls = %v, want [5 127]", got)
	}
	// programmatic Set must not call OnChange
	k.Set(10)
	if len(got) != 2 {
		t.Errorf("Set triggered OnChange: %v", got)
	}
}

func TestPixelTextSize(t *testing.T) {
	if got := PixelTextSize("C", 2); got.X != 5*2 || got.Y != 7*2 {
		t.Errorf("PixelTextSize(\"C\", 2) = %v, want (10, 14)", got)
	}
	if got := PixelTextSize("L63", 2); got.X != (3*6-1)*2 || got.Y != 7*2 {
		t.Errorf("PixelTextSize(\"L63\", 2) = %v, want (%d, 14)", got, (3*6-1)*2)
	}
}

func TestPixFontCoverage(t *testing.T) {
	for _, r := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 -+." {
		if _, ok := pixFont[string(r)]; !ok {
			t.Errorf("glyph %q missing from pixFont", r)
		}
	}
}
