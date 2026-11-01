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
	k := NewKnob(KnobMinMax, "LEVEL", "patch/temp/common/level", 0, 127)
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
	k := NewKnob(KnobMinMax, "LEVEL", "patch/temp/common/level", 0, 127)
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
	k := NewKnob(KnobMinMax, "LEVEL", "patch/temp/common/level", 0, 127)
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
	k := NewKnob(KnobBipolar, "OCTAVE", "patch/temp/common/octave", -3, 3)
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
	k := NewKnob(KnobPan, "PAN", "patch/temp/common/pan", -64, 63)
	cases := map[int]string{-64: "L64", -63: "L63", -1: "L1", 0: "C", 1: "R1", 63: "R63"}
	for v, want := range cases {
		k.Set(v)
		if got := k.formatValue(); got != want {
			t.Errorf("pan formatValue(%d) = %q, want %q", v, got, want)
		}
	}
	// C (0) must map to 12 o'clock, R63 to the far right, L63 to the far left
	k.Set(0)
	if math.Abs(k.angle()) > 1e-9 {
		t.Errorf("pan angle at C = %v, want 0 (12 o'clock)", k.angle())
	}
	k.Set(63)
	if math.Abs(k.angle()-MaxAngle) > 1e-9 {
		t.Errorf("pan angle at R63 = %v, want %v", k.angle(), MaxAngle)
	}
	k.Set(-64)
	if math.Abs(k.angle()-MinAngle) > 1e-9 {
		t.Errorf("pan angle at L64 = %v, want %v", k.angle(), MinAngle)
	}
	k.Set(-1)
	if k.angle() >= 0 || k.angle() <= MinAngle {
		t.Errorf("pan angle at L1 = %v, want in (%v, 0)", k.angle(), MinAngle)
	}
}

func TestBipolarAngleSymmetric(t *testing.T) {
	k := NewKnob(KnobBipolar, "OCTAVE", "patch/temp/common/octave", -3, 3)
	for v, want := range map[int]float64{-3: MinAngle, -1: -45, 0: 0, 1: 45, 3: MaxAngle} {
		k.Set(v)
		if math.Abs(k.angle()-want) > 1e-9 {
			t.Errorf("bipolar angle(%d) = %v, want %v", v, k.angle(), want)
		}
	}
}

func TestLitRangeCentered(t *testing.T) {
	// bipolar at 0: lit range is empty (only the 12 o'clock tick lit)
	k := NewKnob(KnobBipolar, "OCTAVE", "patch/temp/common/octave", -3, 3)
	k.Set(0)
	from, to := litRange(k)
	if from != 0 || to != 0 {
		t.Errorf("litRange at 0 = [%v %v], want [0 0]", from, to)
	}
	k.Set(2)
	from, to = litRange(k)
	if from != 0 || math.Abs(to-90) > 1e-9 {
		t.Errorf("litRange at +2 = [%v %v], want [0 90]", from, to)
	}
	k.Set(-2)
	from, to = litRange(k)
	if math.Abs(from+90) > 1e-9 || to != 0 {
		t.Errorf("litRange at -2 = [%v %v], want [-90 0]", from, to)
	}
	// pan at C
	p := NewKnob(KnobPan, "PAN", "patch/temp/common/pan", -64, 63)
	p.Set(0)
	from, to = litRange(p)
	if from != 0 || to != 0 {
		t.Errorf("pan litRange at C = [%v %v], want [0 0]", from, to)
	}
	// min-max knob: lit range always starts at the left end
	m := NewKnob(KnobMinMax, "LEVEL", "patch/temp/common/level", 0, 127)
	m.Set(64)
	from, to = litRange(m)
	if from != MinAngle {
		t.Errorf("minmax litRange from = %v, want %v", from, MinAngle)
	}
}

func TestCenterTickOnlyForCenteredStyles(t *testing.T) {
	// at minimum, a min-max knob must have an empty lit range (no green
	// tick at 12 o'clock); bipolar/pan at center get the 12 o'clock tick
	m := NewKnob(KnobMinMax, "LEVEL", "patch/temp/common/level", 0, 127)
	m.Set(0)
	from, to := litRange(m)
	if from != MinAngle || to != MinAngle {
		t.Errorf("minmax litRange at 0 = [%v %v], want [%v %v] (empty)", from, to, MinAngle, MinAngle)
	}
	b := NewKnob(KnobBipolar, "OCTAVE", "patch/temp/common/octave", -3, 3)
	b.Set(0)
	from, to = litRange(b)
	if from != 0 || to != 0 {
		t.Errorf("bipolar litRange at 0 = [%v %v], want [0 0]", from, to)
	}
}

func TestKnobSetReportsChange(t *testing.T) {
	k := NewKnob(KnobMinMax, "LEVEL", "patch/temp/common/level", 0, 127)
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
