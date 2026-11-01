package gui

import (
	"testing"

	"gioui.org/f32"
)

// TestRelativeTimes pins the cumulative-time model: points are placed at
// the SUM of the relative t values over the full 4*127 span, and Point
// decodes the relative values back.
func TestRelativeTimes(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	// t1=10, t2=20 => point 2 sits at absolute 30 of 508
	e.SetAll([4]int{10, 20, 30, 40}, [4]int{0, 10, -10, 0})
	if got := e.times(); got != [4]int{10, 20, 30, 40} {
		t.Errorf("times() = %v, want [10 20 30 40]", got)
	}
	if p := e.points[1]; p.X < 29.0/508-0.01 || p.X > 31.0/508+0.01 {
		t.Errorf("point 1 x = %v, want ~30/508=%v", p.X, 30.0/508)
	}
	for i := 0; i < 4; i++ {
		tv, lv := e.Point(i)
		want := [4]int{10, 20, 30, 40}[i]
		if tv != want {
			t.Errorf("Point(%d) time = %d, want %d", i, tv, want)
		}
		_ = lv
	}
	// levels round-trip
	e.SetAll([4]int{10, 20, 30, 40}, [4]int{-63, 0, 63, -10})
	wantL := [4]int{-63, 0, 63, -10}
	for i := 0; i < 4; i++ {
		_, lv := e.Point(i)
		if lv != wantL[i] {
			t.Errorf("Point(%d) level = %d, want %d", i, lv, wantL[i])
		}
	}
}

// TestDragMathRegression pins the drag behavior via dragTo: dragging up
// raises the level, dragging right raises the relative time of the segment,
// far drags are clamped (max level +63, segment 0..127).
func TestDragMathRegression(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
	iw, ih := float32(300), float32(100)
	i := 1

	// press at the point and drag up 30px: level must rise
	start := e.points[i]
	e.dragging, e.dragIdx = true, i
	e.start = f32.Pt(100, 50)
	e.startVal = start
	e.dragTo(i, f32.Pt(100, 20), iw, ih) // 30px up
	if e.points[i].Y <= start.Y {
		t.Errorf("drag up did not raise level: %v -> %v", start, e.points[i])
	}

	// drag right 60px: the segment t2 grows (point 1 moves right)
	before := e.points[i].X
	e.start = f32.Pt(100, 50)
	e.startVal = e.points[i]
	e.dragTo(i, f32.Pt(160, 50), iw, ih)
	if e.points[i].X <= before {
		t.Errorf("drag right did not move point right: %v -> %v", before, e.points[i].X)
	}
	if got := e.times(); got[1] <= 40 {
		t.Errorf("segment t2 = %d, want > 40 after dragging right", got[1])
	}

	// drag far up: clamped at max level +63, stays there but does not exceed
	e.startVal = e.points[i]
	e.start = f32.Pt(160, 50)
	e.dragTo(i, f32.Pt(160, -1000), iw, ih)
	if e.points[i].Y > 1 {
		t.Errorf("Y not clamped to 1: %v", e.points[i].Y)
	}
	if _, lv := e.Point(i); lv != 63 {
		t.Errorf("max level should be +63, got %d", lv)
	}

	// dragging right beyond the next point clamps the segment at the neighbour
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 0, 0, 0})
	e.dragging, e.dragIdx = true, 1
	e.start = f32.Pt(100, 50)
	e.startVal = e.points[1]
	e.dragTo(1, f32.Pt(100+10000, 50), iw, ih) // far right
	got := e.times()
	if got[1] > 127 || got[2] < 0 {
		t.Errorf("segments after far-right drag = %v, want t2<=127, t3>=0", got)
	}
	if e.points[1].X > e.points[2].X+1e-6 {
		t.Errorf("x order violated: p1.X=%v > p2.X=%v", e.points[1].X, e.points[2].X)
	}
}
