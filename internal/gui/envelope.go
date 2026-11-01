package gui

import (
	"image"
	"time"

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

	points   [4]f32.Point // normalized x (0..1 = time 0..127), y (0..1 = level -63..+63)
	defaults [4]f32.Point
	tags     [4]envTag
	active   int

	dragging  bool
	dragIdx   int
	pointerID pointer.ID
	start     f32.Point
	startVal  f32.Point
	lastPress time.Duration
}

type envTag struct{ i int }

const (
	envGap  = 0.03 // min distance between neighbouring points
	envMinX = 0.02
	envMaxX = 0.98
)

var envNames = [4]string{"T1", "T2", "T3", "T4"}

// NewEnvelope returns an envelope widget with the given title and
// parameter key bindings, starting at the proto defaults
// (attack peak, decay to sustain, release).
func NewEnvelope(title string, timeKeys, levelKeys [4]string) *Envelope {
	e := &Envelope{Title: title, TimeKeys: timeKeys, LevelKeys: levelKeys, lastPress: -time.Hour}
	e.defaults = [4]f32.Point{
		{X: 0.14, Y: 1.0}, // attack peak
		{X: 0.30, Y: 0.6}, // decay into sustain
		{X: 0.70, Y: 0.6}, // release start
		{X: 0.96, Y: 0.0}, // release end
	}
	e.points = e.defaults
	for i := range e.tags {
		e.tags[i] = envTag{i: i}
	}
	return e
}

// SetPoint sets point i programmatically from device values (time 0-127,
// level -63..+63) without calling OnChange. The x clamp only enforces the
// ordering against the current neighbours, it never moves a point across
// them; to avoid surprises when loading a whole envelope, callers should
// set points left to right (or use SetAll).
func (e *Envelope) SetPoint(i int, timeVal, levelVal int) {
	if i < 0 || i > 3 {
		return
	}
	e.points[i].X = clampf01(float32(timeVal) / 127)
	e.points[i].Y = clampf01((float32(levelVal) + 63) / 126)
	e.clampPoint(i)
}

// SetAll sets all four points at once from device values; the x ordering
// clamp is applied only after all points are placed, so the values are
// taken over verbatim when they are already ordered.
func (e *Envelope) SetAll(times, levels [4]int) {
	for i := 0; i < 4; i++ {
		e.points[i].X = clampf01(float32(times[i]) / 127)
		e.points[i].Y = clampf01((float32(levels[i]) + 63) / 126)
	}
	for i := 0; i < 4; i++ {
		e.clampPoint(i)
	}
}

// Point returns the device values of point i (time 0-127, level -63..+63).
func (e *Envelope) Point(i int) (timeVal, levelVal int) {
	p := e.points[i]
	return int(roundF(clampf01(p.X) * 127)), int(roundF(clampf01(p.Y)*126)) - 63
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

func (e *Envelope) clampPoint(i int) {
	prev := float32(envMinX)
	if i > 0 {
		prev = e.points[i-1].X
	}
	next := float32(envMaxX)
	if i < len(e.points)-1 {
		next = e.points[i+1].X
	}
	e.points[i].X = clampf(e.points[i].X, prev+envGap, next-envGap)
	e.points[i].X = clampf(e.points[i].X, envMinX, envMaxX)
	e.points[i].Y = clampf(e.points[i].Y, 0, 1)
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

	// The whole screen is the envelope's event area; the tags decide
	// which handle a press/drag belongs to (declared per handle below).
	screen := clip.Rect{Max: image.Pt(w, h)}.Push(gtx.Ops)
	event.Op(gtx.Ops, e)
	screen.Pop()

	// Handles: hit area first (so events work), then events, then visuals.
	for i := range e.points {
		p := px(e.points[i])
		cs := clip.Rect{
			Min: image.Pt(int(p.X-14), int(p.Y-14)),
			Max: image.Pt(int(p.X+14), int(p.Y+14)),
		}.Push(gtx.Ops)
		event.Op(gtx.Ops, &e.tags[i])
		cs.Pop()
	}
	// Drain events addressed to the screen tag itself (fallback for
	// presses outside all handle areas).
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: e, Kinds: pointer.Press})
		if !ok {
			break
		}
		if pe, ok := ev.(pointer.Event); ok && pe.Kind == pointer.Press && pe.Buttons == pointer.ButtonPrimary {
			// pick the nearest handle and start dragging it
			best, bestD := 0, float32(1<<30)
			for i := range e.points {
				p := px(e.points[i])
				d := (pe.Position.X-p.X)*(pe.Position.X-p.X) + (pe.Position.Y-p.Y)*(pe.Position.Y-p.Y)
				if d < bestD {
					best, bestD = i, d
				}
			}
			e.dragging = true
			e.dragIdx = best
			e.pointerID = pe.PointerID
			e.start = pe.Position
			e.startVal = e.points[best]
			e.active = best
			gtx.Execute(pointer.GrabCmd{Tag: &e.tags[best], ID: pe.PointerID})
		}
	}
	for i := range e.points {
		tag := &e.tags[i]
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: tag,
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
				if pe.Time-e.lastPress < dblClick {
					e.points[i] = e.defaults[i] // double-click resets
					e.clampPoint(i)
					e.emit(i)
				} else {
					e.dragging = true
					e.dragIdx = i
					e.pointerID = pe.PointerID
					e.start = pe.Position
					e.startVal = e.points[i]
					gtx.Execute(pointer.GrabCmd{Tag: tag, ID: pe.PointerID})
				}
				e.active = i
				e.lastPress = pe.Time
			case pointer.Drag:
				if !e.dragging || e.dragIdx != i || pe.PointerID != e.pointerID {
					continue
				}
				e.dragTo(i, pe.Position, iw, ih)
				e.emit(i)
			case pointer.Release, pointer.Cancel:
				if e.dragging && e.dragIdx == i {
					e.dragging = false
				}
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

// dragTo moves point i by the drag delta (pos - e.start) in screen pixels,
// given the inner screen size (iw, ih). X grows to the right (time), Y
// grows upwards (level), then the point is clamped.
func (e *Envelope) dragTo(i int, pos f32.Point, iw, ih float32) {
	if iw <= 0 || ih <= 0 {
		return
	}
	e.points[i].X = e.startVal.X + (pos.X-e.start.X)/iw
	e.points[i].Y = e.startVal.Y - (pos.Y-e.start.Y)/ih
	e.clampPoint(i)
}

// emit reports the current device values of point i through OnChange.
func (e *Envelope) emit(i int) {
	if e.OnChange == nil {
		return
	}
	t, l := e.Point(i)
	e.OnChange(e.TimeKeys[i], t, e.LevelKeys[i], l)
}
