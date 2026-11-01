package gui

import (
	"image"
	"image/color"
	"math"
	"time"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

// KnobStyle selects how the value of a Knob maps onto the scale.
type KnobStyle int

const (
	// KnobMinMax has the minimum at the left end (7 o'clock, -135 deg)
	// and the maximum at the right end (5 o'clock, +135 deg).
	KnobMinMax KnobStyle = iota
	// KnobBipolar has 0 at 12 o'clock (0 deg), positive values to the
	// right and negative values to the left.
	KnobBipolar
	// KnobPan has C at 12 o'clock, Lxx to the left and Rxx to the right
	// (raw value range 0..127, 64 = C).
	KnobPan
)

// Knob is a compact hardware-style rotary gadget in the MK-II Tone Stack
// look of cmd/gui-proto: an LED arc plus radial ticks around a metal cap,
// the value as a dot-matrix LED readout in the center of the cap, an
// optional unit in the free space between the scale ends and a mandatory
// title below the scale. The scale spans from 7 o'clock to 5 o'clock and
// lights up as the value moves; the moving scale plus the LED readout
// carry the value, so no separate pointer is needed.
//
// Interaction (same as gui-proto): drag up/down or left/right to dial,
// mouse wheel for steps, Shift = fine, double-click resets to Min.
type Knob struct {
	Style KnobStyle
	Title string // mandatory, written below the scale
	Unit  string // optional, written in the free space between the scale ends

	Min, Max int // value range (for KnobPan: 0..127 raw pan, 64 = C)

	value int

	dragging  bool
	pointerID pointer.ID
	start     f32.Point
	startVal  int
	lastPress time.Duration

	// OnChange is called with the new value whenever the user changes it
	// by dragging, scrolling or double-click reset. It is not called for
	// programmatic Set calls.
	OnChange func(v int)
}

// NewKnob returns a knob with the given style, title and value range.
func NewKnob(style KnobStyle, title string, min, max int) *Knob {
	if max < min {
		min, max = max, min
	}
	return &Knob{Style: style, Title: title, Min: min, Max: max, value: min, lastPress: -time.Hour}
}

// Value returns the current value.
func (k *Knob) Value() int { return k.value }

// Set sets the value programmatically (clamped to the range). It does not
// call OnChange; callers that apply the value to a device should do so
// themselves.
func (k *Knob) Set(v int) {
	k.value = clampInt(v, k.Min, k.Max)
}

// SetRange changes the value range, clamping the current value.
func (k *Knob) SetRange(min, max int) {
	if max < min {
		min, max = max, min
	}
	k.Min, k.Max = min, max
	k.Set(k.value)
}

// angle returns the scale angle of the current value in degrees
// (0 = 12 o'clock, clockwise positive, -135..+135). For bipolar and pan
// knobs the center value maps to 0 deg (12 o'clock), negative values to
// the left (negative angles) and positive values to the right.
func (k *Knob) angle() float64 {
	if k.Style == KnobBipolar || k.Style == KnobPan {
		c := k.center()
		v := clampInt(k.value, k.Min, k.Max) - c
		if v >= 0 {
			span := k.Max - c
			if span <= 0 {
				return 0
			}
			return MaxAngle * float64(v) / float64(span)
		}
		span := c - k.Min
		if span <= 0 {
			return 0
		}
		return MinAngle * float64(-v) / float64(span)
	}
	return angleOf(k.Min, k.Max, k.value)
}

// angleOf maps a value in [min,max] linearly onto [MinAngle,MaxAngle].
func angleOf(min, max, v int) float64 {
	span := max - min
	if span <= 0 {
		return MinAngle
	}
	return MinAngle + (MaxAngle-MinAngle)*float64(clampInt(v, min, max)-min)/float64(span)
}

// formatValue renders the value according to the knob style.
func (k *Knob) formatValue() string {
	switch k.Style {
	case KnobBipolar:
		v := k.value - k.center()
		if v >= 0 {
			return "+" + itoa(v)
		}
		return itoa(v)
	case KnobPan:
		v := k.value - 64
		switch {
		case v == 0:
			return "C"
		case v < 0:
			return "L" + itoa(-v)
		default:
			return "R" + itoa(v)
		}
	default:
		return itoa(k.value)
	}
}

// center returns the value that maps to 12 o'clock: the middle of the
// range for bipolar and pan knobs.
func (k *Knob) center() int {
	switch k.Style {
	case KnobBipolar, KnobPan:
		if k.Style == KnobPan {
			return 64 // raw pan center (C)
		}
		return (k.Min + k.Max) / 2
	default:
		return k.Min
	}
}

// knobDp is the nominal knob diameter in dp; kept compact so that many
// knobs can be stacked next to each other.
const knobDp = 96

// titleDp is the height of the title row below the scale in dp.
const titleDp = 10

// Layout lays out the knob and handles its input.
func (k *Knob) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	size := gtx.Dp(knobDp)
	if max := gtx.Constraints.Max.X; size > max && max > 0 {
		size = max
	}
	total := size + gtx.Dp(titleDp)
	gtx.Constraints = layout.Exact(image.Pt(size, total))

	// declare the knob's tag at its clip area, then drain its events
	cs := clip.Rect{Max: image.Pt(size, total)}.Push(gtx.Ops)
	event.Op(gtx.Ops, k)
	cs.Pop()
	k.processEvents(gtx)

	k.draw(gtx, th, float64(size))
	return layout.Dimensions{Size: image.Pt(size, total)}
}

const dblClick = 400 * time.Millisecond

// scrollAll is a wide scroll range so any wheel amount is delivered.
var scrollAll = pointer.ScrollRange{Min: -(1 << 20), Max: 1 << 20}

// processEvents handles drag / scroll / double-click like gui-proto.
func (k *Knob) processEvents(gtx layout.Context) {
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:  k,
			Kinds:   pointer.Press | pointer.Drag | pointer.Release | pointer.Cancel | pointer.Scroll,
			ScrollX: scrollAll,
			ScrollY: scrollAll,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Kind {
		case pointer.Press:
			if e.Buttons != pointer.ButtonPrimary {
				continue
			}
			if e.Time-k.lastPress < dblClick {
				k.set(k.Min) // double-click resets to minimum
			} else {
				k.dragging = true
				k.pointerID = e.PointerID
				k.start = e.Position
				k.startVal = k.value
				gtx.Execute(pointer.GrabCmd{Tag: k, ID: e.PointerID})
			}
			k.lastPress = e.Time
		case pointer.Drag:
			if !k.dragging || e.PointerID != k.pointerID {
				continue
			}
			span := k.Max - k.Min
			if span <= 0 {
				continue
			}
			sens := float64(span) / 220 // value units per dragged pixel
			if e.Modifiers&key.ModShift != 0 {
				sens /= 4 // fine
			}
			dx := e.Position.X - k.start.X
			dy := e.Position.Y - k.start.Y
			k.set(k.startVal + int(math.Round(float64(dx-dy)*sens))) // up OR right increases
		case pointer.Release, pointer.Cancel:
			k.dragging = false
		case pointer.Scroll:
			// one wheel notch is ~120; Scroll.Y is negative scrolling up
			n := int(math.Round(float64(-e.Scroll.Y) / 120))
			if n == 0 {
				if e.Scroll.Y < 0 {
					n = 1
				} else if e.Scroll.Y > 0 {
					n = -1
				}
			}
			k.set(k.value + n)
		}
	}
}

// set clamps and applies a new value and reports it through OnChange.
func (k *Knob) set(v int) {
	nv := clampInt(v, k.Min, k.Max)
	if nv == k.value {
		return
	}
	k.value = nv
	if k.OnChange != nil {
		k.OnChange(nv)
	}
}

// draw renders the knob in the MK-II Tone Stack style of cmd/gui-proto with
// a green LED arc, radial ticks, metal cap and LED readout in the center.
// For bipolar/pan knobs the lit portion grows from 12 o'clock towards the
// value (left for negative, right for positive); at the center position
// only the 12 o'clock tick is lit.
func (k *Knob) draw(gtx layout.Context, th *Theme, size float64) {
	cx, cy := size/2, size/2

	// lit range of the scale: [litFrom, litTo]
	litFrom, litTo := litRange(k)

	// LED arc: unlit base over the whole scale, lit portion (gui-proto
	// drawAmp, green)
	strokeArc(gtx, cx, cy, size*0.329, MinAngle, MaxAngle, 5, th.ArcBase)
	if litTo > litFrom {
		strokeArc(gtx, cx, cy, size*0.329, litFrom, litTo, 5, th.Lit)
	} else if k.Style == KnobBipolar || k.Style == KnobPan {
		// center position of a bipolar/pan knob: light only the 12
		// o'clock tick itself
		strokeArc(gtx, cx, cy, size*0.329, -0.5, 0.5, 5, th.Lit)
	}
	// radial ticks around the arc; ticks between litFrom and litTo are lit,
	// plus the 12 o'clock tick at the center position of bipolar/pan knobs
	tickFrom, tickTo := litFrom, litTo
	if tickTo <= tickFrom && (k.Style == KnobBipolar || k.Style == KnobPan) {
		tickFrom, tickTo = -0.5, 0.5
	}
	drawTicksRange(gtx, cx, cy, size*0.371, size*0.421, 40, 4, th.TickBase, th.Lit, MinAngle, MaxAngle, tickFrom, tickTo)

	// metal cap
	fillCircleGradient(gtx, cx, cy, size*0.271, th.CapTop, th.CapBot)
	strokeCircle(gtx, cx, cy, size*0.271, 1, th.CapEdge)
	fillCircleGradient(gtx, cx, cy, size*0.207, rgb(0x5D5C60), rgb(0x0C0B0E))
	strokeCircle(gtx, cx, cy, size*0.207, 1.5, rgb(0x000000))

	// value as dot-matrix LED readout in the center of the cap
	val := k.formatValue()
	cell := gtx.Dp(FontCellDp)
	vsz := PixelTextSize(val, cell)
	offset(gtx, int(cx)-vsz.X/2, int(cy)-vsz.Y/2, func() {
		DrawPixelText(gtx, val, cell)
	})

	// optional unit in the free space between the scale ends (bottom gap)
	if k.Unit != "" {
		us := PixelTextSize(k.Unit, cell)
		offset(gtx, int(cx)-us.X/2, int(cy+size*0.30), func() {
			DrawPixelText(gtx, k.Unit, cell)
		})
	}

	// mandatory title below the scale
	ts := PixelTextSize(k.Title, cell)
	offset(gtx, int(cx)-ts.X/2, int(size)+1, func() {
		DrawPixelText(gtx, k.Title, cell)
	})
}

// litRange returns the lit angle range of the scale. Min-max knobs light
// from the left end (7 o'clock) up to the value; bipolar and pan knobs
// light from 12 o'clock towards the value (left for negative values,
// right for positive ones); at the center position the range is empty.
func litRange(k *Knob) (from, to float64) {
	angle := k.angle()
	if k.Style == KnobBipolar || k.Style == KnobPan {
		if angle >= 0 {
			return 0, angle
		}
		return angle, 0
	}
	return MinAngle, angle
}

func strokeCircle(gtx layout.Context, cx, cy, r float64, width float32, col color.NRGBA) {
	strokeArc(gtx, cx, cy, r, 0, 360, width, col)
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
