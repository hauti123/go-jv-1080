package gui

import (
	"testing"

	"gioui.org/f32"
)

// TestDragMathRegression pins the drag behavior of Envelope.Layout via
// dragTo: dragging up raises the level, dragging right raises the time,
// and far drags are clamped (max level +63, x order preserved).
func TestDragMathRegression(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{10, 40, 80, 120}, [4]int{0, 20, -20, 0})
	iw, ih := float32(300), float32(100)
	i := 1

	// press at the current point position and drag up 30px
	start := e.points[i]
	e.dragging, e.dragIdx = true, i
	e.start = f32.Pt(100, 50)
	e.startVal = start
	e.dragTo(i, f32.Pt(100, 50-30), iw, ih) // 30px up
	if e.points[i].Y <= start.Y {
		t.Errorf("drag up did not raise level: %v -> %v", start, e.points[i])
	}

	// drag right 60px from there
	e.start = f32.Pt(100, 20)
	e.startVal = e.points[i]
	e.dragTo(i, f32.Pt(160, 20), iw, ih)
	if e.points[i].X <= start.X {
		t.Errorf("drag right did not move time: %v -> %v", start, e.points[i])
	}

	// drag far up: clamped at max level +63
	e.startVal = e.points[i]
	e.start = f32.Pt(160, 20)
	e.dragTo(i, f32.Pt(160, -1000), iw, ih)
	if e.points[i].Y > 1 {
		t.Errorf("Y not clamped to 1: %v", e.points[i].Y)
	}
	if _, lv := e.Point(i); lv != 63 {
		t.Errorf("max level should be +63, got %d", lv)
	}

	// x order is preserved when dragging towards a neighbour
	e.SetAll([4]int{10, 40, 80, 120}, [4]int{0, 20, -20, 0})
	e.dragging, e.dragIdx = true, 1
	e.start = f32.Pt(100, 50)
	e.startVal = e.points[1]
	e.dragTo(1, f32.Pt(-1000, 50), iw, ih) // far left, would cross point 0
	t0, _ := e.Point(0)
	t1, _ := e.Point(1)
	if t1 <= t0 {
		t.Errorf("x order violated: t1=%d <= t0=%d", t1, t0)
	}
}
