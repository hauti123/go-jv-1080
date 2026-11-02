package widget

import (
	"testing"
)

func TestEnvelopePointRoundTrip(t *testing.T) {
	e := NewEnvelope("PITCH ENV T1",
		[4]string{"patch/temp/tone/1/pitch/t1", "patch/temp/tone/1/pitch/t2", "patch/temp/tone/1/pitch/t3", "patch/temp/tone/1/pitch/t4"},
		[4]string{"patch/temp/tone/1/pitch/l1", "patch/temp/tone/1/pitch/l2", "patch/temp/tone/1/pitch/l3", "patch/temp/tone/1/pitch/l4"})
	// ordered values round-trip exactly
	e.SetAll([4]int{10, 40, 80, 120}, [4]int{-63, -10, 20, 63})
	for i := 0; i < 4; i++ {
		gt, gl := e.Point(i)
		wantT := [4]int{10, 40, 80, 120}[i]
		wantL := [4]int{-63, -10, 20, 63}[i]
		if gt != wantT || gl != wantL {
			t.Errorf("point %d: got (%d,%d), want (%d,%d)", i, gt, gl, wantT, wantL)
		}
	}
	// single SetPoint also round-trips when it does not cross a neighbour
	e.SetPoint(2, 90, 30)
	if gt, gl := e.Point(2); gt != 90 || gl != 30 {
		t.Errorf("SetPoint(2,90,30) -> (%d,%d)", gt, gl)
	}
}

func TestEnvelopeXOrderPreserved(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"a", "b", "c", "d"}, [4]string{"e", "f", "g", "h"})
	e.SetAll([4]int{100, 10, 10, 10}, [4]int{0, 0, 0, 0})
	// cumulative x positions must be non-decreasing
	for i := 1; i < 4; i++ {
		if e.points[i].X < e.points[i-1].X {
			t.Errorf("x order violated: p[%d].X=%v < p[%d].X=%v", i, e.points[i].X, i-1, e.points[i-1].X)
		}
	}
	// and the relative times decode to the clamped segments
	got := e.times
	if got[0] != 100 || got[1] != 10 {
		t.Errorf("times() = %v, want first segments [100 10 ...]", got)
	}
}

func TestEnvelopeOnChangeKeys(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"tk0", "tk1", "tk2", "tk3"}, [4]string{"lk0", "lk1", "lk2", "lk3"})
	var gotT, gotL string
	var tv, lv int
	e.OnChange = func(timeKey string, timeVal int, levelKey string, levelVal int) {
		gotT, tv, gotL, lv = timeKey, timeVal, levelKey, levelVal
	}
	e.SetPoint(2, 90, 30)
	e.emit(2)
	if gotT != "tk2" || gotL != "lk2" || tv != 90 || lv != 30 {
		t.Errorf("emit(2) = %q %d %q %d, want tk2 90 lk2 30", gotT, tv, gotL, lv)
	}
}
