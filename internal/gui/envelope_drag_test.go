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

// TestPressDragSequence simulates the full event flow of Layout with a
// single screen tag: press near a handle -> grab -> drags move exactly
// that point relative to the press position (no jumping to fixed values).
func TestPressDragSequence(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
	iw, ih := float32(300), float32(100)

	// press exactly at point 1's position (px mapping: padX=8, padT=8)
	padX, padT := float32(8), float32(8)
	p1 := e.points[1]
	press := f32.Pt(padX+p1.X*iw, padT+(1-p1.Y)*ih)

	// emulate the Press handler: nearest handle, record start
	best, bestD := 0, float32(1<<30)
	for i := range e.points {
		p := e.points[i]
		pp := f32.Pt(padX+p.X*iw, padT+(1-p.Y)*ih)
		d := (press.X-pp.X)*(press.X-pp.X) + (press.Y-pp.Y)*(press.Y-pp.Y)
		if d < bestD {
			best, bestD = i, d
		}
	}
	if best != 1 {
		t.Fatalf("press hit handle %d, want 1", best)
	}
	e.dragging, e.dragIdx = true, best
	e.start = press
	e.startVal = e.points[best]

	// drag 30px up and 60px right: point 1 moves by the delta, but is
	// clamped at the next point minus the gap; the following points then
	// shift along by the applied (clamped) delta.
	before := e.points[1]
	p2before := e.points[2].X
	e.dragTo(1, f32.Pt(press.X+60, press.Y-30), iw, ih)
	wantX := before.X + 60/iw
	if wantX > p2before-envGap {
		wantX = p2before - envGap // clamped at the next point (before the shift)
	}
	wantY := before.Y + 30/ih
	if d := e.points[1].X - wantX; d < -1e-4 || d > 1e-4 {
		t.Errorf("after drag x = %v, want %v", e.points[1].X, wantX)
	}
	if d := e.points[1].Y - wantY; d < -1e-4 || d > 1e-4 {
		t.Errorf("after drag y = %v, want %v", e.points[1].Y, wantY)
	}
	// no jump to a fixed position: the segment time grew by the drag amount
	got := e.times()
	wantSeg := 40 + int(roundF((wantX-before.X)*508))
	if got[1] < wantSeg-2 || got[1] > wantSeg+2 {
		t.Errorf("segment t2 = %d, want ~%d", got[1], wantSeg)
	}
	// following points moved along by the same (clamped) delta
	if d := (e.points[2].X - p2before) - (wantX - before.X); d < -1e-4 || d > 1e-4 {
		t.Errorf("point 2 did not shift along: moved %v, want %v", d, wantX-before.X)
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

	// dragging right shifts the following points along (segments keep
	// their length); a far drag clamps at the screen edge
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 0, 0, 0})
	e.dragging, e.dragIdx = true, 1
	e.start = f32.Pt(100, 50)
	e.startVal = e.points[1]
	e.dragTo(1, f32.Pt(100+30, 50), iw, ih) // 30px right
	got := e.times()
	if got[1] <= 40 {
		t.Errorf("segment t2 = %d, want > 40 after dragging right", got[1])
	}
	if got[2] != 40 || got[3] != 40 {
		t.Errorf("following segments changed: %v, want t3=t4=40 (shifted along)", got)
	}
	for i := 1; i < 4; i++ {
		if e.points[i].X < e.points[i-1].X {
			t.Errorf("x order violated at %d", i)
		}
	}
	// far right: point 1 stops at the last point minus gap, last point at edge
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 0, 0, 0})
	e.dragging, e.dragIdx = true, 1
	e.start = f32.Pt(100, 50)
	e.startVal = e.points[1]
	e.dragTo(1, f32.Pt(100+10000, 50), iw, ih)
	if e.points[3].X > envMaxX+1e-6 {
		t.Errorf("last point beyond screen: %v", e.points[3].X)
	}
	if e.points[1].X > e.points[2].X+1e-6 {
		t.Errorf("x order violated: p1.X=%v > p2.X=%v", e.points[1].X, e.points[2].X)
	}
}
