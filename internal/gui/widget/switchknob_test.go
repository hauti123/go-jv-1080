package widget

import (
	"math"
	"testing"
)

func TestSwitchAngleEndpoints(t *testing.T) {
	// 6 positions over the same 7 o'clock..5 o'clock scale as the knob
	if got := switchAngle(0, 6); math.Abs(got-MinAngle) > 1e-9 {
		t.Errorf("switchAngle(0) = %v, want %v (7 o'clock)", got, MinAngle)
	}
	if got := switchAngle(5, 6); math.Abs(got-MaxAngle) > 1e-9 {
		t.Errorf("switchAngle(5) = %v, want %v (5 o'clock)", got, MaxAngle)
	}
	// evenly spread: position 2 of 6 is at -27 deg
	if got := switchAngle(2, 6); math.Abs(got-(-27)) > 1e-9 {
		t.Errorf("switchAngle(2) = %v, want -27", got)
	}
}

func TestSnapAngleRoundTrip(t *testing.T) {
	for pos := 0; pos < 6; pos++ {
		a := switchAngle(pos, 6)
		if got := snapAngle(a, 6); got != pos {
			t.Errorf("snapAngle(switchAngle(%d)) = %d, want %d", pos, got, pos)
		}
		// angles slightly off the detent still snap to it
		if got := snapAngle(a+10, 6); got != pos {
			t.Errorf("snapAngle(%v+10) = %d, want %d", a, got, pos)
		}
		if got := snapAngle(a-10, 6); got != pos {
			t.Errorf("snapAngle(%v-10) = %d, want %d", a, got, pos)
		}
	}
	// out of scale clamps to the end positions
	if got := snapAngle(200, 6); got != 5 {
		t.Errorf("snapAngle(200) = %d, want 5", got)
	}
	if got := snapAngle(-200, 6); got != 0 {
		t.Errorf("snapAngle(-200) = %d, want 0", got)
	}
}

func TestDetentAngles(t *testing.T) {
	angles := detentAngles(6)
	if len(angles) != 6 {
		t.Fatalf("detentAngles(6) has %d entries, want 6 (ticks only at the positions)", len(angles))
	}
	if angles[0] != MinAngle || angles[5] != MaxAngle {
		t.Errorf("detentAngles endpoints = [%v %v], want [%v %v]", angles[0], angles[5], MinAngle, MaxAngle)
	}
	for i := 1; i < len(angles); i++ {
		if angles[i] <= angles[i-1] {
			t.Errorf("detentAngles not increasing at %d: %v <= %v", i, angles[i], angles[i-1])
		}
	}
}

func TestSwitchKnobSetClamp(t *testing.T) {
	s := NewSwitchKnob("WAVE GRP", "patch/temp/tone/1", []string{"INT-A", "INT-B", "EXP-A", "EXP-B", "EXP-C", "EXP-D"})
	if s.Value() != 0 {
		t.Errorf("initial value = %d, want 0", s.Value())
	}
	if s.Positions() != 6 {
		t.Errorf("Positions() = %d, want 6", s.Positions())
	}
	s.Set(9)
	if s.Value() != 5 {
		t.Errorf("Set(9) = %d, want 5", s.Value())
	}
	s.Set(-3)
	if s.Value() != 0 {
		t.Errorf("Set(-3) = %d, want 0", s.Value())
	}
}

func TestSwitchKnobAngle(t *testing.T) {
	s := NewSwitchKnob("WAVE GRP", "patch/temp/tone/1", []string{"INT-A", "INT-B", "EXP-A", "EXP-B", "EXP-C", "EXP-D"})
	s.Set(5)
	if math.Abs(s.angle()-MaxAngle) > 1e-9 {
		t.Errorf("angle at position 5 = %v, want %v", s.angle(), MaxAngle)
	}
	s.Set(0)
	if math.Abs(s.angle()-MinAngle) > 1e-9 {
		t.Errorf("angle at position 0 = %v, want %v", s.angle(), MinAngle)
	}
}

func TestSwitchKnobSetReportsChange(t *testing.T) {
	s := NewSwitchKnob("WAVE GRP", "patch/temp/tone/1", []string{"INT-A", "INT-B", "EXP-A", "EXP-B", "EXP-C", "EXP-D"})
	var got []int
	s.OnChange = func(pos int) { got = append(got, pos) }
	s.set(3)
	s.set(3) // no change, no callback
	s.set(99)
	if len(got) != 2 || got[0] != 3 || got[1] != 5 {
		t.Errorf("OnChange calls = %v, want [3 5]", got)
	}
	// programmatic Set must not call OnChange
	s.Set(1)
	if len(got) != 2 {
		t.Errorf("Set triggered OnChange: %v", got)
	}
}
