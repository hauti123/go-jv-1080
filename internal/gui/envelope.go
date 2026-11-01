package gui

import (
	"image"
	"log"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

// Envelope is the 4-point envelope editor of the gui-proto ENV-4 GENERATOR
// panel: a green-on-black screen with a grid, the envelope line with area
// fill and four draggable handles. Each handle i is bound to a time key
// TimeKeys[i] and a level key LevelKeys[i] (path-like device parameter keys,
// e.g. "patch/temp/tone/1/pitch/t1" and ".../l1"); dragging a handle sends
// both values through OnChange. The x position of a point encodes the
// time (0-127), the y position the level (-63..+63, displayed 0..126).
//
// Interaction (same as gui-proto): drag the points freely in x/y (x order
// is preserved), double-click a point resets it to its default.
type Envelope struct {
	Title string // mandatory, written below the screen

	// TimeKeys and LevelKeys bind the 4 points to device parameters;
	// both must have length 4.
	TimeKeys  [4]string
	LevelKeys [4]string

	// OnChange is called with (timeKey, timeValue, levelKey, levelValue)
	// whenever the user moves a point. Time values are 0-127, level
	// values -63..+63.
	OnChange func(timeKey string, timeVal int, levelKey string, levelVal int)

	// times/levels are the device values (relative times 0-127, levels
	// -63..+63); points holds the derived pixel positions.
	times    [4]int
	levels   [4]int
	points   [4]f32.Point // x: cumulative time 0..1 (= sum t1..ti over 4*127), y: level 0..1 (-63..+63)
	defaults [4]f32.Point
	active   int

	dragging    bool
	dragIdx     int
	pointerID   pointer.ID
	lastDrag    f32.Point // position of the previous drag event (same frame as current)
	haveLast    bool
	accTF       float32 // accumulated time delta since press, in device units (float; rounded only when applied)
	accLF       float32 // accumulated level delta since press, in device units
	startTimes  [4]int
	startLevels [4]int
}

const (
	envMinX = 0.02
	envMaxX = 0.98
)

var envNames = [4]string{"T1", "T2", "T3", "T4"}

// NewEnvelope returns an envelope widget with the given title and
// parameter key bindings, starting at the proto defaults
// (attack peak, decay to sustain, release).
func NewEnvelope(title string, timeKeys, levelKeys [4]string) *Envelope {
	e := &Envelope{Title: title, TimeKeys: timeKeys, LevelKeys: levelKeys}
	// defaults (proto shape): attack peak, decay to sustain, release
	e.times = [4]int{71, 81, 152, 122} // cumulative 71/152/304/426 of 508
	e.levels = [4]int{63, 19, 19, -63}
	e.relayout()
	e.defaults = e.points
	return e
}

// envMaxTime is the total time span of the envelope: the four time
// parameters are relative to each other, so the screen covers the sum of
// all four (max 4*127).
const envMaxTime = 4 * 127

// SetPoint sets point i programmatically from device values (time 0-127
// relative to the previous point, level -63..+63) without calling
// OnChange.
func (e *Envelope) SetPoint(i int, timeVal, levelVal int) {
	if i < 0 || i > 3 {
		return
	}
	e.times[i] = clampInt(timeVal, 0, 127)
	e.levels[i] = clampInt(levelVal, -63, 63)
	e.relayout()
}

// SetAll sets all four points at once from device values (relative times
// t1..t4, levels l1..l4).
func (e *Envelope) SetAll(times, levels [4]int) {
	for i := 0; i < 4; i++ {
		e.times[i] = clampInt(times[i], 0, 127)
		e.levels[i] = clampInt(levels[i], -63, 63)
	}
	e.relayout()
}

// Point returns the device values of point i (relative time 0-127, level
// -63..+63).
func (e *Envelope) Point(i int) (timeVal, levelVal int) {
	return e.times[i], e.levels[i]
}

// relayout recomputes the pixel positions from times/levels: x is the
// cumulative time (sum of t1..ti) over the full 4*127 span, y is the
// level mapped from -63..+63 to 0..1 (top = +63).
func (e *Envelope) relayout() {
	cum := 0
	for i := 0; i < 4; i++ {
		cum += clampInt(e.times[i], 0, 127)
		e.points[i].X = clampf01(float32(cum) / envMaxTime)
		e.points[i].Y = clampf01((float32(clampInt(e.levels[i], -63, 63)) + 63) / 126)
		log.Printf("ENV relayout: point %d cum=%d -> x=%.4f (%.1f px of 508), level=%d -> y=%.4f",
			i, cum, e.points[i].X, e.points[i].X*envMaxTime, clampInt(e.levels[i], -63, 63), e.points[i].Y)
	}
}

func absF(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func roundF(v float32) float32 {
	if v < 0 {
		return float32(int(v - 0.5))
	}
	return float32(int(v + 0.5))
}

// Active returns the index of the last touched point.
func (e *Envelope) Active() int { return e.active }

func clampf01(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func clampf(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// envDp is the nominal height of the envelope screen in dp (gui-proto 170).
const envDp = 150

// Layout lays out the envelope editor and handles its input.
func (e *Envelope) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	w := gtx.Constraints.Max.X
	h := gtx.Dp(envDp)
	if w < gtx.Dp(120) {
		w = gtx.Dp(120)
	}
	total := h + gtx.Dp(10) // title row below the screen
	gtx.Constraints = layout.Exact(image.Pt(w, total))

	// Everything below (drawing and pointer events) is in physical pixels.
	// pxPerDp > 1 on HiDPI means 1 logical dp is several physical px; pointer
	// deltas arrive in logical units (dp) after the router's inverse scale,
	// so convert drag deltas dp->px with PxPerDp to match the pixel geometry.
	log.Printf("ENV layout: w=%d h=%d total=%d pxPerDp=%.2f", w, h, total, gtx.Metric.PxPerDp)
	padX := float32(gtx.Dp(8))
	padT := float32(gtx.Dp(8))
	padB := float32(gtx.Dp(14))
	iw := float32(w) - 2*padX
	ih := float32(h) - padT - padB
	baseY := padT + ih
	px := func(p f32.Point) f32.Point {
		return f32.Pt(padX+p.X*iw, padT+(1-p.Y)*ih)
	}

	// Screen chrome (gui-proto ENV-4 screen).
	fillRect(gtx, 0, 0, float32(w), float32(h), rgb(0x6F6D66))
	fillRect(gtx, 2, 2, float32(w)-4, float32(h)-4, rgb(0x030503))

	// Grid.
	for i := 0; i <= 8; i++ {
		x := padX + iw*float32(i)/8
		strokeLine(gtx, f32.Pt(x, padT), f32.Pt(x, baseY), 1, rgb(ledDim))
	}
	for i := 0; i <= 4; i++ {
		y := padT + ih*float32(i)/4
		strokeLine(gtx, f32.Pt(padX, y), f32.Pt(padX+iw, y), 1, rgb(ledDim))
	}

	// Area fill + envelope line.
	coords := []f32.Point{f32.Pt(padX, baseY)}
	for _, p := range e.points {
		coords = append(coords, px(p))
	}
	area := append([]f32.Point{}, coords...)
	area = append(area, f32.Pt(coords[len(coords)-1].X, baseY), f32.Pt(padX, baseY))
	areaCol := rgb(ledLit)
	areaCol.A = 18
	fillPolygon(gtx, area, areaCol)
	strokePolyline(gtx, coords, 2.5, rgb(ledLit))

	// Origin marker.
	fillRect(gtx, padX-3, baseY-3, 6, 6, rgb(0x2F7F3A))

	// One tag for the whole screen, declared over the full area. All
	// events (press/drag/release) are received on this single tag, whose
	// coordinate system is the envelope's own, so pe.Position is directly
	// comparable with the handle positions.
	screen := clip.Rect{Max: image.Pt(w, h)}.Push(gtx.Ops)
	event.Op(gtx.Ops, e)
	screen.Pop()

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: e,
			Kinds:  pointer.Press | pointer.Drag | pointer.Release | pointer.Cancel,
		})
		if !ok {
			break
		}
		pe, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		switch pe.Kind {
		case pointer.Press:
			if pe.Buttons != pointer.ButtonPrimary {
				continue
			}
			// pick the nearest handle
			best, bestD := 0, float32(1<<30)
			for i := range e.points {
				p := px(e.points[i])
				d := (pe.Position.X-p.X)*(pe.Position.X-p.X) + (pe.Position.Y-p.Y)*(pe.Position.Y-p.Y)
				if d < bestD {
					best, bestD = i, d
				}
			}
			log.Printf("ENV press: pos=(%.1f,%.1f) -> handle %d (dist %.1f), handle px=%v",
				pe.Position.X, pe.Position.Y, best, bestD, px(e.points[best]))
			e.beginDrag(best, pe.PointerID, pe.Position)
			gtx.Execute(pointer.GrabCmd{Tag: e, ID: pe.PointerID})
			e.active = best
		case pointer.Drag:
			if !e.dragging || pe.PointerID != e.pointerID {
				continue
			}
			if !e.haveLast {
				// first drag event: anchor only, no value change
				e.lastDrag = pe.Position
				e.haveLast = true
				continue
			}
			dx, dy := pe.Position.X-e.lastDrag.X, pe.Position.Y-e.lastDrag.Y
			// event positions are logical (dp); geometry is physical px
			dx *= gtx.Metric.PxPerDp
			dy *= gtx.Metric.PxPerDp
			// A one-time coordinate-frame switch (e.g. when the pointer grab
			// takes effect) shows up as an implausibly large delta; skip it
			// and re-anchor instead of applying it.
			if absF(dx) > float32(w) || absF(dy) > float32(h) {
				log.Printf("ENV drag: handle=%d frame-switch skipped delta=(%.1f,%.1f)", e.dragIdx, dx, dy)
				e.lastDrag = pe.Position
				continue
			}
			e.lastDrag = pe.Position
			e.dragBy(e.dragIdx, dx, dy, iw, ih)
			e.emit(e.dragIdx)
		case pointer.Release, pointer.Cancel:
			log.Printf("ENV release: pos=(%.1f,%.1f) dragging=%v", pe.Position.X, pe.Position.Y, e.dragging)
			if pe.PointerID == e.pointerID {
				e.dragging = false
			}
		}
	}

	// Handle visuals.
	for i := range e.points {
		p := px(e.points[i])
		fillRect(gtx, p.X-6, p.Y-6, 12, 12, rgb(0x0A0F0A))
		col := rgb(ledLit)
		if i == e.active {
			col = rgb(0xB6FFCF)
		}
		fillRect(gtx, p.X-4.5, p.Y-4.5, 9, 9, col)
		// LED label above the point (below when near the top edge).
		cell := gtx.Dp(FontCellDp)
		lbl := envNames[i]
		lsz := PixelTextSize(lbl, cell)
		ty := int(p.Y) - lsz.Y - 3
		if e.points[i].Y > 0.8 {
			ty = int(p.Y) + 8
		}
		Offset(gtx, int(p.X)-lsz.X/2, ty, func() { DrawPixelText(gtx, lbl, cell) })
	}

	// Mandatory title below the screen.
	cell := gtx.Dp(FontCellDp)
	ts := PixelTextSize(e.Title, cell)
	Offset(gtx, (w-ts.X)/2, h+1, func() { DrawPixelText(gtx, e.Title, cell) })

	return layout.Dimensions{Size: image.Pt(w, total)}
}

// beginDrag starts dragging handle i: records the pointer and the press
// position plus the device values at press time.
func (e *Envelope) beginDrag(i int, id pointer.ID, pos f32.Point) {
	e.dragging = true
	e.dragIdx = i
	e.pointerID = id
	e.haveLast = false
	e.accTF, e.accLF = 0, 0
	e.startTimes = e.times
	e.startLevels = e.levels
	log.Printf("ENV drag-begin: handle=%d pointer=%v press=(%.1f,%.1f) times=%v levels=%v",
		i, id, pos.X, pos.Y, e.times, e.levels)
}

// dragBy moves point i by a drag delta (dx, dy) in screen pixels, given
// the inner screen size (iw, ih). The delta is the difference between two
// consecutive drag events, which share the same coordinate frame, so the
// result is independent of the absolute frame (local vs. global after a
// grab). The delta is converted into device units (time: full span
// envMaxTime over iw px; level: 126 over ih px) and accumulated onto the
// values captured at press time, then applied to the relative time/level
// of point i. The segment time is clamped to 0..127 and the total span to
// envMaxTime; dragging in x shifts all following segments along so their
// lengths stay unchanged. Y (level) is clamped to -63..+63 independently.
func (e *Envelope) dragBy(i int, dx, dy float32, iw, ih float32) {
	if iw <= 0 || ih <= 0 {
		return
	}
	// accumulate in float so sub-unit fractions of small per-event deltas
	// are not lost to rounding; round only when applying to the value
	e.accTF += dx / iw * envMaxTime
	e.accLF += -dy / ih * 126
	oldT, oldL := e.times[i], e.levels[i]
	e.setTime(i, e.startTimes[i]+int(roundF(e.accTF)))
	e.levels[i] = clampInt(e.startLevels[i]+int(roundF(e.accLF)), -63, 63)
	e.relayout()
	log.Printf("ENV drag: handle=%d delta=(%.1f,%.1f) acc=(%.2f,%.2f)->(%d,%d) iw=%.1f ih=%.1f | t%d %d->%d l%d %d->%d | times=%v levels=%v | points=%v",
		i, dx, dy, e.accTF, e.accLF, int(roundF(e.accTF)), int(roundF(e.accLF)), iw, ih,
		i, oldT, e.times[i], i, oldL, e.levels[i], e.times, e.levels, e.points)
}

// setTime sets the relative time (segment length) of point i to v,
// clamped to 0..127. The point's absolute position is the cumulative sum
// of t1..ti, so changing segment i shifts every following point by the
// same delta (their own segments keep their length). The total span is
// kept within envMaxTime by clamping v against the sum of the others.
func (e *Envelope) setTime(i, v int) {
	v = clampInt(v, 0, 127)
	others := 0
	for j := 0; j < 4; j++ {
		if j != i {
			others += e.times[j]
		}
	}
	v = clampInt(v, 0, envMaxTime-others)
	log.Printf("ENV setTime: handle=%d requested=%d clamped=%d (others sum=%d) old=%d", i, v, v, others, e.times[i])
	e.times[i] = v
	// following segments are untouched, so relayout() shifts their
	// cumulative pixel positions by the same delta automatically
}

// emit reports the current device values of point i through OnChange.
func (e *Envelope) emit(i int) {
	if e.OnChange == nil {
		return
	}
	t, l := e.Point(i)
	log.Printf("ENV emit: handle=%d -> %s=%d %s=%d", i, e.TimeKeys[i], t, e.LevelKeys[i], l)
	e.OnChange(e.TimeKeys[i], t, e.LevelKeys[i], l)
}
