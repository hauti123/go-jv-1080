package widget

import (
	"image"
	"math"
	"time"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"

	"go-jv-1080/internal/util"
)

// SwitchKnob is a stepped rotary switch with detents, e.g. for selecting
// the wave group (INT-A, INT-B, EXP-A..D) of a tone. The positions are
// spread evenly over the same scale as the Knob (7 o'clock to 5 o'clock,
// -135..+135 deg), but ticks are drawn only at the detent positions. The
// center readout shows the label of the current position.
//
// Like the Knob, a SwitchKnob is bound to a device parameter through Key,
// a path-like string resolved solely by the device layer; OnChange is
// called with the position index (0..len(Labels)-1).
//
// Interaction: drag towards a position (the switch snaps to the nearest
// detent), mouse wheel steps one position, double-click resets to the
// first position.
type SwitchKnob struct {
	Title  string   // mandatory, written below the scale
	Key    string   // device parameter key this switch is bound to
	Labels []string // one label per position, e.g. "INT-A".."EXP-D"

	value int

	dragging  bool
	pointerID pointer.ID
	cx, cy    float64 // center in local coordinates, set during Layout
	lastPress time.Duration

	// OnChange is called with the new position whenever the user changes
	// it by dragging, scrolling or double-click reset. It is not called
	// for programmatic Set calls.
	OnChange func(pos int)
}

// NewSwitchKnob returns a switch knob with the given title, device
// parameter key and one label per detent position.
func NewSwitchKnob(title, key string, labels []string) *SwitchKnob {
	return &SwitchKnob{Title: title, Key: key, Labels: labels, lastPress: -time.Hour}
}

// Positions returns the number of detent positions.
func (s *SwitchKnob) Positions() int { return len(s.Labels) }

// Value returns the current position index.
func (s *SwitchKnob) Value() int { return s.value }

// Set sets the position programmatically (clamped to the range). It does
// not call OnChange.
func (s *SwitchKnob) Set(pos int) {
	s.value = util.Clamp(pos, 0, max(len(s.Labels)-1, 0))
}

// switchAngle maps position pos of count positions linearly onto the
// scale [MinAngle, MaxAngle] (7 o'clock to 5 o'clock).
func switchAngle(pos, count int) float64 {
	if count <= 1 {
		return MinAngle
	}
	return MinAngle + (MaxAngle-MinAngle)*float64(util.Clamp(pos, 0, count-1))/float64(count-1)
}

// snapAngle returns the detent position nearest to angle a (clamped into
// the scale range).
func snapAngle(a float64, count int) int {
	if count <= 1 {
		return 0
	}
	a = util.Clamp(a, MinAngle, MaxAngle)
	return int(math.Round((a - MinAngle) / (MaxAngle - MinAngle) * float64(count-1)))
}

// detentAngles returns the scale angles of all count detent positions.
func detentAngles(count int) []float64 {
	angles := make([]float64, count)
	for i := range angles {
		angles[i] = switchAngle(i, count)
	}
	return angles
}

// angle returns the scale angle of the current position.
func (s *SwitchKnob) angle() float64 { return switchAngle(s.value, len(s.Labels)) }

// Layout lays out the switch knob and handles its input.
func (s *SwitchKnob) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	size := gtx.Dp(knobDp)
	if max := gtx.Constraints.Max.X; size > max && max > 0 {
		size = max
	}
	total := size + gtx.Dp(titleDp)
	gtx.Constraints = layout.Exact(image.Pt(size, total))
	s.cx, s.cy = float64(size)/2, float64(size)/2

	cs := clip.Rect{Max: image.Pt(size, total)}.Push(gtx.Ops)
	event.Op(gtx.Ops, s)
	cs.Pop()
	s.processEvents(gtx)

	s.draw(gtx, th, float64(size))
	return layout.Dimensions{Size: image.Pt(size, total)}
}

// processEvents handles drag (snap to nearest detent), scroll (step one
// position) and double-click (reset to the first position).
func (s *SwitchKnob) processEvents(gtx layout.Context) {
	n := len(s.Labels)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:  s,
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
			if e.Time-s.lastPress < dblClick {
				s.set(0) // double-click resets to the first position
			} else {
				s.dragging = true
				s.pointerID = e.PointerID
				gtx.Execute(pointer.GrabCmd{Tag: s, ID: e.PointerID})
			}
			s.lastPress = e.Time
		case pointer.Drag:
			if !s.dragging || e.PointerID != s.pointerID || n == 0 {
				continue
			}
			// angle of the pointer around the center (0 = 12 o'clock,
			// clockwise positive), snapped to the nearest detent
			dx := float64(e.Position.X) - s.cx
			dy := float64(e.Position.Y) - s.cy
			a := math.Atan2(dx, -dy) * 180 / math.Pi
			s.set(snapAngle(a, n))
		case pointer.Release, pointer.Cancel:
			s.dragging = false
		case pointer.Scroll:
			n := int(math.Round(float64(-e.Scroll.Y) / 120))
			if n == 0 {
				if e.Scroll.Y < 0 {
					n = 1
				} else if e.Scroll.Y > 0 {
					n = -1
				}
			}
			s.set(s.value + n)
		}
	}
}

// set clamps and applies a new position and reports it through OnChange.
func (s *SwitchKnob) set(pos int) {
	nv := util.Clamp(pos, 0, max(len(s.Labels)-1, 0))
	if nv == s.value {
		return
	}
	s.value = nv
	if s.OnChange != nil {
		s.OnChange(nv)
	}
}

// draw renders the switch knob like the Knob (LED arc, metal cap, LED
// readout, title), but with ticks only at the detent positions and the
// position label as the readout.
func (s *SwitchKnob) draw(gtx layout.Context, th *Theme, size float64) {
	cx, cy := size/2, size/2
	n := len(s.Labels)
	angle := s.angle()

	// LED arc: unlit base over the whole scale, lit portion up to the
	// current position
	strokeArc(gtx, cx, cy, size*0.329, MinAngle, MaxAngle, 5, th.ArcBase)
	if angle > MinAngle {
		strokeArc(gtx, cx, cy, size*0.329, MinAngle, angle, 5, th.Lit)
	}
	// ticks only at the detent positions; detents up to the current
	// position are lit
	for i, a := range detentAngles(n) {
		col := th.TickBase
		if i <= s.value {
			col = th.Lit
		}
		strokeLine(gtx, polarPt(cx, cy, size*0.36, a), polarPt(cx, cy, size*0.43, a), 2, col)
	}

	// metal cap
	fillCircleGradient(gtx, cx, cy, size*0.271, th.CapTop, th.CapBot)
	strokeCircle(gtx, cx, cy, size*0.271, 1, th.CapEdge)
	fillCircleGradient(gtx, cx, cy, size*0.207, rgb(0x5D5C60), rgb(0x0C0B0E))
	strokeCircle(gtx, cx, cy, size*0.207, 1.5, rgb(0x000000))

	// position label as dot-matrix LED readout in the center of the cap,
	// shrunk to fit the cap when necessary
	label := ""
	if s.value >= 0 && s.value < n {
		label = s.Labels[s.value]
	}
	cell := gtx.Dp(FontCellDp)
	if w1 := PixelTextSize(label, 1).X; w1 > 0 && PixelTextSize(label, cell).X > int(size*0.38) {
		cell = max(1, int(size*0.38)/w1)
	}
	vsz := PixelTextSize(label, cell)
	Offset(gtx, int(cx)-vsz.X/2, int(cy)-vsz.Y/2, func() {
		DrawPixelText(gtx, label, cell)
	})

	// mandatory title below the scale
	cell = gtx.Dp(FontCellDp)
	ts := PixelTextSize(s.Title, cell)
	Offset(gtx, int(cx)-ts.X/2, int(size)+1, func() {
		DrawPixelText(gtx, s.Title, cell)
	})
}
