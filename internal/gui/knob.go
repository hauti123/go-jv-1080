package gui

import (
	"image"
	"math"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// KnobStyle selects how the value of a Knob maps onto the scale.
type KnobStyle int

const (
	// KnobMinMax has 0 at the left end (7 o'clock) and the maximum at the
	// right end (5 o'clock).
	KnobMinMax KnobStyle = iota
	// KnobBipolar has 0 at 12 o'clock, positive values to the right and
	// negative values to the left.
	KnobBipolar
	// KnobPan has C at 12 o'clock, Lxx to the left and Rxx to the right.
	KnobPan
)

// scale geometry: the scale spans from 7 o'clock to 5 o'clock, i.e. from
// 210 degrees to -30 degrees (330 degrees) going clockwise through 12
// o'clock. In screen coordinates (y down) 0 degrees points right and
// angles increase clockwise.
const (
	scaleStartAngle = 210.0 // 7 o'clock
	scaleSweep      = 300.0 // degrees, clockwise through 12 o'clock to 5 o'clock
)

// Knob is a compact rotary gadget: a colored arc scale from 7 to 5
// o'clock, the value in the center, an optional unit below the value and a
// mandatory title below the scale. There is no knob in the middle; the
// scale itself moves. Drag vertically (or use the scroll wheel) to change
// the value; every change is reported through OnChange.
type Knob struct {
	Style KnobStyle
	Title string // mandatory, written below the scale
	Unit  string // optional, written between the scale ends

	Min, Max int // value range (for KnobPan: 0..127 raw pan, 64 = C)

	value int

	drag    gesture.Drag
	scroll  gesture.Scroll
	changed bool

	// OnChange is called with the new value whenever the user changes it
	// by dragging or scrolling. It is not called for programmatic Set
	// calls.
	OnChange func(v int)
}

// NewKnob returns a knob with the given style, title and value range.
func NewKnob(style KnobStyle, title string, min, max int) *Knob {
	if max < min {
		min, max = max, min
	}
	return &Knob{Style: style, Title: title, Min: min, Max: max, value: min}
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

// center returns the value that maps to 12 o'clock.
func (k *Knob) center() int {
	switch k.Style {
	case KnobBipolar, KnobPan:
		return (k.Min + k.Max) / 2
	default:
		return k.Min
	}
}

// fraction returns the position of the value on the scale in [0,1].
func (k *Knob) fraction() float64 {
	span := k.Max - k.Min
	if span <= 0 {
		return 0
	}
	return float64(k.value-k.Min) / float64(span)
}

// angleAt returns the scale angle (degrees, screen coordinates) for a
// fraction in [0,1].
func angleAt(frac float64) float64 {
	return scaleStartAngle + frac*scaleSweep
}

// Layout lays out the knob and handles its input. The knob is as compact
// as the constraints allow.
func (k *Knob) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	for {
		ev, ok := k.drag.Update(gtx.Metric, gtx.Source, gesture.Vertical)
		if !ok {
			break
		}
		switch ev.Kind {
		case pointer.Drag, pointer.Release:
			k.applyDelta(-int(ev.Position.Y) / dragPixelsPerStep)
		}
	}
	sdist := k.scroll.Update(gtx.Metric, gtx.Source, gtx.Now, gesture.Vertical,
		pointer.ScrollRange{},
		pointer.ScrollRange{Min: math.MinInt32, Max: math.MaxInt32})
	if sdist != 0 {
		k.applyDelta(sdist / scrollPixelsPerStep)
	}

	size := knobSize(gtx)
	macro := op.Record(gtx.Ops)
	k.layout(gtx, th, size)
	call := macro.Stop()

	area := clip.Rect(image.Rect(0, 0, size.X, size.Y)).Push(gtx.Ops)
	k.drag.Add(gtx.Ops)
	k.scroll.Add(gtx.Ops)
	event.Op(gtx.Ops, k)
	area.Pop()

	call.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

const (
	dragPixelsPerStep   = 4
	scrollPixelsPerStep = 12
)

// applyDelta changes the value by steps and reports the change.
func (k *Knob) applyDelta(steps int) {
	if steps == 0 {
		return
	}
	nv := clampInt(k.value+steps, k.Min, k.Max)
	if nv == k.value {
		return
	}
	k.value = nv
	k.changed = true
	if k.OnChange != nil {
		k.OnChange(nv)
	}
}

// ConsumeChanges reports whether the user changed the value since the last
// call, and resets the flag.
func (k *Knob) ConsumeChanges() bool {
	c := k.changed
	k.changed = false
	return c
}

// knobSize computes the compact knob size: square scale area plus the title
// row below.
func knobSize(gtx layout.Context) image.Point {
	d := gtx.Constraints
	s := d.Max.X
	if d.Max.Y < s {
		s = d.Max.Y - titleHeight(gtx)
	}
	min := gtx.Dp(unit.Dp(48))
	if s < min {
		s = min
	}
	// keep room for the title below the scale
	return image.Point{X: s, Y: s + titleHeight(gtx)}
}

func titleHeight(gtx layout.Context) int {
	return int(float32(glyph5x7Height+1)*FontScale) + gtx.Dp(unit.Dp(2))
}

// layout draws the knob into a size.X x size.Y box (the title occupies the
// bottom titleHeight pixels).
func (k *Knob) layout(gtx layout.Context, th *Theme, size image.Point) {
	scaleH := size.Y - titleHeight(gtx)
	cx := float32(size.X) / 2
	cy := float32(scaleH) / 2
	radius := float32(minInt(size.X, scaleH)) / 2
	thickness := float32(gtx.Dp(unit.Dp(3)))
	if thickness < 2 {
		thickness = 2
	}

	// background surface
	paint.FillShape(gtx.Ops, th.Surface, clip.Rect(image.Rect(0, 0, size.X, size.Y)).Op())

	// scale ticks: unlit base arc, then the lit portion up to the value
	const ticks = 24
	for i := 0; i <= ticks; i++ {
		frac := float64(i) / float64(ticks)
		ang := angleAt(frac) * math.Pi / 180
		// tick from inner to outer radius, perpendicular (radial) segment
		dx := float32(math.Cos(ang))
		dy := float32(math.Sin(ang))
		inner := radius - thickness
		p0 := f32.Point{X: cx + dx*inner, Y: cy + dy*inner}
		p1 := f32.Point{X: cx + dx*radius, Y: cy + dy*radius}
		col := th.LEDOff
		if frac <= k.fraction()+1e-9 {
			col = th.LEDOn
		}
		paint.FillShape(gtx.Ops, col, tickShape(gtx.Ops, p0, p1, thickness*0.6))
	}

	// value in the center of the knob
	val := k.formatValue()
	vw := LEDTextWidth(val, FontScale)
	LEDText(gtx.Ops, f32.Point{X: cx - vw/2, Y: cy - (glyph5x7Height*FontScale)/2}, val, FontScale, th.Text)

	// optional unit in the free space between the scale ends (below center)
	if k.Unit != "" {
		uw := LEDTextWidth(k.Unit, FontScale)
		unitY := cy + radius*float32(math.Sin(math.Pi/2))*0.35
		LEDText(gtx.Ops, f32.Point{X: cx - uw/2, Y: unitY}, k.Unit, FontScale, th.TextDim)
	}

	// mandatory title below the scale
	tw := LEDTextWidth(k.Title, FontScale)
	LEDText(gtx.Ops, f32.Point{X: cx - tw/2, Y: float32(scaleH) + 1}, k.Title, FontScale, th.TextDim)
}

// tickShape builds a small quad (thick line) path for a scale tick.
func tickShape(ops *op.Ops, p0, p1 f32.Point, width float32) clip.Op {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	l := float32(math.Hypot(float64(dx), float64(dy)))
	if l == 0 {
		return clip.Rect(image.Rectangle{}).Op()
	}
	nx := -dy / l * width / 2
	ny := dx / l * width / 2
	var p clip.Path
	p.Begin(ops)
	p.MoveTo(f32.Point{X: p0.X + nx, Y: p0.Y + ny})
	p.LineTo(f32.Point{X: p1.X + nx, Y: p1.Y + ny})
	p.LineTo(f32.Point{X: p1.X - nx, Y: p1.Y - ny})
	p.LineTo(f32.Point{X: p0.X - nx, Y: p0.Y - ny})
	p.Close()
	return clip.Outline{Path: p.End()}.Op()
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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
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
