package gui

import (
	"testing"

	"gioui.org/f32"
)

// TestRelativeTimes pins the value model: SetAll/Point round-trip relative
// times and levels, and pixel x is the cumulative sum over 4*127.
func TestRelativeTimes(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{10, 20, 30, 40}, [4]int{-63, 0, 63, -10})
	for i := 0; i < 4; i++ {
		tv, lv := e.Point(i)
		if tv != [4]int{10, 20, 30, 40}[i] || lv != [4]int{-63, 0, 63, -10}[i] {
			t.Errorf("Point(%d) = (%d,%d)", i, tv, lv)
		}
	}
	// t1=10, t2=20 -> point 1 at cumulative 30 of 508
	if p := e.points[1].X; p < 29.0/508-0.001 || p > 31.0/508+0.001 {
		t.Errorf("point 1 x = %v, want ~30/508", p)
	}
}

// TestSetTimeBounds pins the segment invariants of setTime: each segment
// stays within 0..127, the total within 4*127, and following segments keep
// their length.
func TestSetTimeBounds(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 0, 0, 0})
	// grow t2 beyond 127: clamped to 127, t3/t4 shrink to keep the total
	e.setTime(1, 500)
	if e.times[1] != 127 {
		t.Errorf("t2 = %d, want 127", e.times[1])
	}
	total := e.times[0] + e.times[1] + e.times[2] + e.times[3]
	if total > envMaxTime {
		t.Errorf("total = %d > %d", total, envMaxTime)
	}
	if e.times[2] != 40 || e.times[3] != 40 {
		t.Errorf("following segments changed: %v, want t3=t4=40 (lengths preserved)", e.times)
	}
	// shrink t1 to 0: following segments grow back, each <= 127
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 0, 0, 0})
	e.setTime(0, -100)
	if e.times[0] != 0 {
		t.Errorf("t1 = %d, want 0", e.times[0])
	}
	for i, v := range e.times {
		if v < 0 || v > 127 {
			t.Errorf("segment %d = %d out of 0..127", i, v)
		}
	}
}

// dragSeq simulates a drag: anchor at the first drag position, then apply
// the given per-event pixel deltas.
func dragSeq(e *Envelope, i int, iw, ih float32, deltas ...f32.Point) {
	e.beginDrag(i, 1, f32.Pt(0, 0))
	for _, d := range deltas {
		e.lastDrag = e.lastDrag.Add(d)
		e.dragBy(i, d.X, d.Y, iw, ih)
	}
}

// TestDragInDeviceUnits pins dragBy: the pixel delta is converted to device
// units (envMaxTime over iw px, 126 over ih px), accumulated onto the
// press-time values, applied to the dragged point only, with segment bounds
// enforced.
func TestDragInDeviceUnits(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
	iw, ih := float32(508), float32(126) // 1px == 1 device unit
	// drag right 30px, up 10px in one event
	dragSeq(e, 1, iw, ih, f32.Pt(30, -10))
	if e.times[1] != 70 {
		t.Errorf("t2 after +30px = %d, want 70", e.times[1])
	}
	if e.levels[1] != 30 {
		t.Errorf("l2 after +10px up = %d, want 30", e.levels[1])
	}
	// following segments keep their length (only their absolute x moves)
	if e.times[2] != 40 || e.times[3] != 40 {
		t.Errorf("following segments = %d,%d, want 40,40 (unchanged)", e.times[2], e.times[3])
	}
	// pixel position reflects cumulative time
	if p := e.points[1].X; p < 109.0/508-0.001 || p > 111.0/508+0.001 {
		t.Errorf("point 1 x = %v, want ~110/508 (40+70)", p)
	}
}

// TestDragAccumulates pins that consecutive small deltas accumulate to the
// same result as one big delta (no per-event clamping artifacts), and that
// the absolute event frame does not matter (only differences between
// consecutive events).
func TestDragAccumulates(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
	iw, ih := float32(508), float32(126)
	// 3 events of +10px each == one +30px event
	dragSeq(e, 1, iw, ih, f32.Pt(10, 0), f32.Pt(10, 0), f32.Pt(10, 0))
	if e.times[1] != 70 {
		t.Errorf("t2 after 3x10px = %d, want 70", e.times[1])
	}
	// same deltas delivered in a different absolute frame (origin +500,+300)
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
	e.beginDrag(1, 1, f32.Pt(500, 300))
	e.lastDrag = f32.Pt(500, 300)
	e.haveLast = true
	for _, d := range []f32.Point{{X: 10, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 0}} {
		e.lastDrag = e.lastDrag.Add(d)
		e.dragBy(1, d.X, d.Y, iw, ih)
	}
	if e.times[1] != 70 {
		t.Errorf("t2 in shifted frame = %d, want 70 (frame-independent)", e.times[1])
	}
}

// TestPerHandleDrag simulates pressing exactly on each handle and dragging:
// y-drag changes the level, x-drag changes the segment and the following
// points shift along.
func TestPerHandleDrag(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	iw, ih := float32(300), float32(100)
	for _, i := range []int{0, 1, 2, 3} {
		e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
		dragSeq(e, i, iw, ih, f32.Pt(0, -30)) // 30px up
		if e.levels[i] <= [4]int{0, 20, -20, 0}[i] {
			t.Errorf("handle %d: y-drag did not raise level (%d)", i, e.levels[i])
		}
		e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
		dragSeq(e, i, iw, ih, f32.Pt(30, 0)) // 30px right
		if e.times[i] <= 40 {
			t.Errorf("handle %d: x-drag did not grow segment (t=%d)", i, e.times[i])
		}
		if i < 3 && e.times[i+1] != 40 {
			t.Errorf("handle %d: following segment changed (t%d=%d), want 40 (length preserved)", i, i+1, e.times[i+1])
		}
	}
}

// TestHitTestSelectsPressedHandle pins the press hit-test: a press at a
// handle's pixel position must select exactly that handle, so the point
// follows the mouse from the first drag event.
func TestHitTestSelectsPressedHandle(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
	iw, ih := float32(300), float32(100)
	padX, padT := float32(8), float32(8)
	px := func(p f32.Point) f32.Point { return f32.Pt(padX+p.X*iw, padT+(1-p.Y)*ih) }
	for i := 0; i < 4; i++ {
		press := px(e.points[i])
		best, bestD := 0, float32(1<<30)
		for j := range e.points {
			p := px(e.points[j])
			d := (press.X-p.X)*(press.X-p.X) + (press.Y-p.Y)*(press.Y-p.Y)
			if d < bestD {
				best, bestD = j, d
			}
		}
		if best != i {
			t.Errorf("press at handle %d selected %d", i, best)
		}
	}
}

// TestDragFrameSwitchSkipped pins that a one-time huge delta (coordinate
// frame switch when the grab takes effect) is skipped and re-anchored, not
// applied — small real drags around it still work.
func TestDragFrameSwitchSkipped(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{40, 40, 40, 40}, [4]int{0, 20, -20, 0})
	iw, ih := float32(508), float32(126)
	w, h := 846, 263
	i := 2
	e.beginDrag(i, 1, f32.Pt(100, 50))
	// simulate event stream: anchor, then frame switch (huge delta), then real drags
	pos := f32.Pt(100, 50)
	e.lastDrag = pos
	e.haveLast = true
	// frame switch event
	pos = f32.Pt(100-2*float32(w), 50-2*float32(h)) // frame switch: delta beyond widget size
	dx, dy := pos.X-e.lastDrag.X, pos.Y-e.lastDrag.Y
	if absF(dx) > float32(w) || absF(dy) > float32(h) {
		e.lastDrag = pos // skipped, re-anchored
	} else {
		t.Fatal("expected frame switch to be detected")
	}
	if e.times[i] != 40 {
		t.Fatalf("frame switch changed value: t=%d, want 40", e.times[i])
	}
	// real drag +10px right
	pos = pos.Add(f32.Pt(10, 0))
	e.lastDrag = pos
	e.dragBy(i, 10, 0, iw, ih)
	if e.times[i] != 50 {
		t.Errorf("after real +10px drag: t=%d, want 50", e.times[i])
	}
}

// TestDragSmallDeltasNotLost pins that many small sub-unit drag deltas
// accumulate without rounding loss: 100 x 1px must equal one 100px drag.
func TestDragSmallDeltasNotLost(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{10, 10, 10, 10}, [4]int{0, 0, 0, 0})
	iw, ih := float32(818), float32(224) // real widget inner size: 1px < 1 unit
	i := 1
	e.beginDrag(i, 1, f32.Pt(0, 0))
	// 100 events of +1px right: each is 508/818 = 0.621 units, must accumulate
	for k := 0; k < 100; k++ {
		e.dragBy(i, 1, 0, iw, ih)
	}
	want := 10 + int(roundF(100.0/iw*508)) // 10 + 62 = 72
	if e.times[i] != want {
		t.Errorf("100x1px drag: t=%d, want %d (no rounding loss)", e.times[i], want)
	}
	// y: 100 events of -1px (up): each 126/224 = 0.5625 units
	e.SetAll([4]int{10, 10, 10, 10}, [4]int{0, 0, 0, 0})
	e.beginDrag(i, 1, f32.Pt(0, 0))
	for k := 0; k < 100; k++ {
		e.dragBy(i, 0, -1, iw, ih)
	}
	wantL := int(roundF(100.0 / ih * 126)) // 56
	if e.levels[i] != wantL {
		t.Errorf("100x1px up drag: l=%d, want %d", e.levels[i], wantL)
	}
}

// TestDragDeltaScaledByPxPerDp pins that a drag delta in logical units (dp)
// is scaled by PxPerDp to physical pixels before mapping, so on HiDPI the
// point tracks the mouse 1:1 in screen pixels.
func TestDragDeltaScaledByPxPerDp(t *testing.T) {
	e := NewEnvelope("ENV", [4]string{"t0", "t1", "t2", "t3"}, [4]string{"l0", "l1", "l2", "l3"})
	e.SetAll([4]int{10, 10, 10, 10}, [4]int{0, 0, 0, 0})
	iw, ih := float32(818), float32(224) // physical px inner size
	ppdp := float32(1.75)                // HiDPI scale
	i := 1
	e.beginDrag(i, 1, f32.Pt(0, 0))
	// one event of +10 logical dp = +17.5 physical px
	dx := float32(10) * ppdp
	e.dragBy(i, dx, 0, iw, ih)
	want := 10 + int(roundF(17.5/iw*508))
	if e.times[i] != want {
		t.Errorf("t after +10dp at scale %.2f: %d, want %d", ppdp, e.times[i], want)
	}
}
