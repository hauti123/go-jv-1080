package gui

import (
	"math"
	"testing"
)

func TestScaleAngles(t *testing.T) {
	// 7 o'clock = 210 degrees, 5 o'clock = -30 degrees (330), sweep 300 deg
	if got := angleAt(0); math.Abs(got-210) > 1e-9 {
		t.Errorf("angleAt(0) = %v, want 210 (7 o'clock)", got)
	}
	if got := angleAt(1); math.Abs(got-510) > 1e-9 {
		t.Errorf("angleAt(1) = %v, want 510 (= -30, 5 o'clock)", got)
	}
	// 12 o'clock is at the middle of the sweep for bipolar knobs
	if got := angleAt(0.5); math.Abs(got-360) > 1e-9 {
		t.Errorf("angleAt(0.5) = %v, want 360 (= 0 deg, 12 o'clock)", got)
	}
}

func TestKnobFraction(t *testing.T) {
	k := NewKnob(KnobMinMax, "LEVEL", 0, 127)
	if k.fraction() != 0 {
		t.Errorf("fraction at min = %v, want 0", k.fraction())
	}
	k.Set(127)
	if k.fraction() != 1 {
		t.Errorf("fraction at max = %v, want 1", k.fraction())
	}
	k.Set(64)
	if f := k.fraction(); math.Abs(f-64.0/127.0) > 1e-9 {
		t.Errorf("fraction at 64 = %v, want %v", f, 64.0/127.0)
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
	if f := k.fraction(); math.Abs(f-0.5) > 1e-9 {
		t.Errorf("bipolar fraction at 0 = %v, want 0.5 (12 o'clock)", f)
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
	// C must map to 12 o'clock
	k.Set(64)
	if f := k.fraction(); math.Abs(f-64.0/127.0) > 1e-9 {
		t.Errorf("pan fraction at C = %v, want %v", f, 64.0/127.0)
	}
}

func TestKnobApplyDeltaReportsChange(t *testing.T) {
	k := NewKnob(KnobMinMax, "LEVEL", 0, 127)
	var got []int
	k.OnChange = func(v int) { got = append(got, v) }
	k.applyDelta(5)
	k.applyDelta(0) // no change, no callback
	k.applyDelta(1000)
	if len(got) != 2 || got[0] != 5 || got[1] != 127 {
		t.Errorf("OnChange calls = %v, want [5 127]", got)
	}
	if !k.ConsumeChanges() {
		t.Error("ConsumeChanges = false after change, want true")
	}
	if k.ConsumeChanges() {
		t.Error("ConsumeChanges = true twice, want false")
	}
}

func TestLEDTextWidth(t *testing.T) {
	if got := LEDTextWidth("C", FontScale); got != 5*FontScale {
		t.Errorf("LEDTextWidth(\"C\") = %v, want %v", got, 5*FontScale)
	}
	if got := LEDTextWidth("L63", FontScale); got != (3*6-1)*FontScale {
		t.Errorf("LEDTextWidth(\"L63\") = %v, want %v", got, (3*6-1)*FontScale)
	}
}

func TestGlyphLookup(t *testing.T) {
	if _, ok := glyph5x7Lookup(' '); !ok {
		t.Error("space glyph missing")
	}
	if _, ok := glyph5x7Lookup('a'); !ok {
		t.Error("lowercase should map to uppercase glyph")
	}
	if _, ok := glyph5x7Lookup('é'); ok {
		t.Error("unsupported rune should not be found")
	}
}
