// Knob Lab — a Gio port of the HTML knob/envelope style prototype.
//
// It demonstrates heavily customized widgets drawn from scratch with Gio
// primitives (no stock widgets): hardware-style knobs in five designs,
// a draggable 4-point envelope editor, and matching dropdown boxes.
//
// Setup and run (Go 1.22+):
//
//	mkdir knoblab && cd knoblab
//	go mod init knoblab
//	go get gioui.org@latest
//	go mod tidy
//	go run .
//
// Linux additionally needs Wayland/X11 development packages, see
// https://gioui.org/doc/install for per-platform requirements.
//
// Interactions (same as the HTML prototype):
//   - Knobs: drag up/down OR left/right to dial, mouse wheel for steps,
//     Shift = fine, double-click resets. The stepped switch snaps to 12 detents.
//   - Envelope: drag the four points freely in x/y (order is preserved),
//     double-click a point to reset it.
//
// Everything maps to the immediate-mode model you know from iced: each
// widget is state + a draw function; events arrive through gtx.Event.
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const (
	minVal   float32 = 0
	maxVal   float32 = 10
	minAngle         = -135.0
	maxAngle         = 135.0
	knobSize         = 140 // dp
	dblClick         = 400 * time.Millisecond
)

// Wide scroll range so any wheel amount is delivered to the knobs.
var scrollAll = pointer.ScrollRange{Min: -(1 << 20), Max: 1 << 20}

func rgb(hex uint32) color.NRGBA {
	return color.NRGBA{R: uint8(hex >> 16), G: uint8(hex >> 8), B: uint8(hex), A: 0xFF}
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

func angleOf(v float32) float64 {
	return minAngle + (maxAngle-minAngle)*float64((v-minVal)/(maxVal-minVal))
}

// ---------------------------------------------------------------------------
// Low-level drawing helpers (the Gio equivalents of the SVG primitives)
// ---------------------------------------------------------------------------

func polarPt(cx, cy, r, deg float64) f32.Point {
	a := deg * math.Pi / 180
	return f32.Pt(float32(cx+r*math.Sin(a)), float32(cy-r*math.Cos(a)))
}

func fillCircle(gtx layout.Context, cx, cy, r float64, col color.NRGBA) {
	b := image.Rect(int(cx-r), int(cy-r), int(cx+r+0.5), int(cy+r+0.5))
	paint.FillShape(gtx.Ops, col, clip.Ellipse(b).Op(gtx.Ops))
}

// fillCircleGradient fakes the 3D shading of a knob cap with a linear
// gradient (Gio has no radial gradient primitive).
func fillCircleGradient(gtx layout.Context, cx, cy, r float64, c1, c2 color.NRGBA) {
	b := image.Rect(int(cx-r), int(cy-r), int(cx+r+0.5), int(cy+r+0.5))
	stack := clip.Ellipse(b).Push(gtx.Ops)
	paint.LinearGradientOp{
		Stop1:  f32.Pt(float32(cx-r*0.8), float32(cy-r*0.9)),
		Color1: c1,
		Stop2:  f32.Pt(float32(cx+r), float32(cy+r)),
		Color2: c2,
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Pop()
}

func strokeCircle(gtx layout.Context, cx, cy, r float64, width float32, col color.NRGBA) {
	strokeArc(gtx, cx, cy, r, 0, 360, width, col)
}

func strokeArc(gtx layout.Context, cx, cy, r, a0, a1 float64, width float32, col color.NRGBA) {
	if a1-a0 < 0.5 {
		return
	}
	const segments = 48
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(polarPt(cx, cy, r, a0))
	sweep := a1 - a0
	for i := 1; i <= segments; i++ {
		path.LineTo(polarPt(cx, cy, r, a0+sweep*float64(i)/segments))
	}
	paint.FillShape(gtx.Ops, col, clip.Stroke{Path: path.End(), Width: width}.Op())
}

func strokeLine(gtx layout.Context, p1, p2 f32.Point, width float32, col color.NRGBA) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(p1)
	path.LineTo(p2)
	paint.FillShape(gtx.Ops, col, clip.Stroke{Path: path.End(), Width: width}.Op())
}

func fillPolygon(gtx layout.Context, pts []f32.Point, col color.NRGBA) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(pts[0])
	for _, p := range pts[1:] {
		path.LineTo(p)
	}
	path.Close()
	paint.FillShape(gtx.Ops, col, clip.Outline{Path: path.End()}.Op())
}

func strokePolyline(gtx layout.Context, pts []f32.Point, width float32, col color.NRGBA) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(pts[0])
	for _, p := range pts[1:] {
		path.LineTo(p)
	}
	paint.FillShape(gtx.Ops, col, clip.Stroke{Path: path.End(), Width: width}.Op())
}

func fillRect(gtx layout.Context, x, y, w, h float32, col color.NRGBA) {
	paint.FillShape(gtx.Ops, col, clip.Rect{
		Min: image.Pt(int(x), int(y)),
		Max: image.Pt(int(x+w+0.5), int(y+h+0.5)),
	}.Op())
}

// drawTicks draws radial tick marks; ticks up to the current angle use "lit".
func drawTicks(gtx layout.Context, cx, cy, rIn, rOut float64, count, majorEvery int, base, lit color.NRGBA, value float32) {
	angle := angleOf(value)
	for i := 0; i <= count; i++ {
		major := i%majorEvery == 0
		a := minAngle + (maxAngle-minAngle)*float64(i)/float64(count)
		ri, ro, w := rIn, rOut, 1.0
		if !major {
			ri, ro, w = rIn+2, rOut-2, 0.8
		}
		col := base
		if a <= angle+0.01 {
			col = lit
		}
		strokeLine(gtx, polarPt(cx, cy, ri, a), polarPt(cx, cy, ro, a), float32(w), col)
	}
}

// rotated runs draw() under a rotation of deg around (cx, cy).
func rotated(gtx layout.Context, cx, cy, deg float64, draw func()) {
	rad := float32(deg * math.Pi / 180)
	stack := op.Affine(f32.Affine2D{}.Rotate(f32.Pt(float32(cx), float32(cy)), rad)).Push(gtx.Ops)
	draw()
	stack.Pop()
}

// offset runs draw() under a scoped translation (unlike a raw
// op.Offset(...).Add, it cannot leak into subsequent drawing).
func offset(gtx layout.Context, x, y int, draw func()) {
	stack := op.Affine(f32.Affine2D{}.Offset(f32.Pt(float32(x), float32(y)))).Push(gtx.Ops)
	draw()
	stack.Pop()
}

// ---------------------------------------------------------------------------
// 5x7 dot-matrix font (the 90s LED display; drawn pixel by pixel)
// ---------------------------------------------------------------------------

var pixFont = map[string][7]string{
	"A": {"01110", "10001", "10001", "11111", "10001", "10001", "10001"},
	"B": {"11110", "10001", "10001", "11110", "10001", "10001", "11110"},
	"C": {"01110", "10001", "10000", "10000", "10000", "10001", "01110"},
	"D": {"11110", "10001", "10001", "10001", "10001", "10001", "11110"},
	"E": {"11111", "10000", "10000", "11110", "10000", "10000", "11111"},
	"F": {"11111", "10000", "10000", "11110", "10000", "10000", "10000"},
	"G": {"01110", "10001", "10000", "10111", "10001", "10001", "01111"},
	"H": {"10001", "10001", "10001", "11111", "10001", "10001", "10001"},
	"I": {"01110", "00100", "00100", "00100", "00100", "00100", "01110"},
	"K": {"10001", "10010", "10100", "11000", "10100", "10010", "10001"},
	"L": {"10000", "10000", "10000", "10000", "10000", "10000", "11111"},
	"M": {"10001", "11011", "10101", "10101", "10001", "10001", "10001"},
	"N": {"10001", "10001", "11001", "10101", "10011", "10001", "10001"},
	"O": {"01110", "10001", "10001", "10001", "10001", "10001", "01110"},
	"P": {"11110", "10001", "10001", "11110", "10000", "10000", "10000"},
	"R": {"11110", "10001", "10001", "11110", "10100", "10010", "10001"},
	"S": {"01111", "10000", "10000", "01110", "00001", "00001", "11110"},
	"T": {"11111", "00100", "00100", "00100", "00100", "00100", "00100"},
	"U": {"10001", "10001", "10001", "10001", "10001", "10001", "01110"},
	"V": {"10001", "10001", "10001", "10001", "10001", "01010", "00100"},
	"X": {"10001", "10001", "01010", "00100", "01010", "10001", "10001"},
	"Y": {"10001", "10001", "01010", "00100", "00100", "00100", "00100"},
	"Z": {"11111", "00001", "00010", "00100", "01000", "10000", "11111"},
	"0": {"01110", "10001", "10011", "10101", "11001", "10001", "01110"},
	"1": {"00100", "01100", "00100", "00100", "00100", "00100", "01110"},
	"2": {"01110", "10001", "00001", "00010", "00100", "01000", "11111"},
	"3": {"11111", "00010", "00100", "00010", "00001", "10001", "01110"},
	"4": {"00010", "00110", "01010", "10010", "11111", "00010", "00010"},
	"5": {"11111", "10000", "11110", "00001", "00001", "10001", "01110"},
	"6": {"00110", "01000", "10000", "11110", "10001", "10001", "01110"},
	"7": {"11111", "00001", "00010", "00100", "01000", "01000", "01000"},
	"8": {"01110", "10001", "10001", "01110", "10001", "10001", "01110"},
	"9": {"01110", "10001", "10001", "01111", "00001", "00010", "01100"},
	" ": {"00000", "00000", "00000", "00000", "00000", "00000", "00000"},
	"-": {"00000", "00000", "00000", "01110", "00000", "00000", "00000"},
	".": {"00000", "00000", "00000", "00000", "00000", "00110", "00110"},
}

const (
	ledLit = 0x39FF70 // bright green
	ledDim = 0x0B3315 // unlit dot grid
)

func pixelTextSize(s string, cellPx int) image.Point {
	n := len([]rune(s))
	return image.Pt((n*6-1)*cellPx, 7*cellPx)
}

// drawPixelText draws s as a green-on-black dot matrix at the current origin.
func drawPixelText(gtx layout.Context, s string, cellPx int) layout.Dimensions {
	s = strings.ToUpper(s)
	for i, r := range []rune(s) {
		glyph, ok := pixFont[string(r)]
		if !ok {
			glyph = pixFont[" "]
		}
		for row := 0; row < 7; row++ {
			for col := 0; col < 5; col++ {
				colr := rgb(ledDim)
				colr.A = 140
				if glyph[row][col] == '1' {
					colr = rgb(ledLit)
				}
				x := (i*6 + col) * cellPx
				y := row * cellPx
				fillRect(gtx, float32(x), float32(y), float32(cellPx), float32(cellPx), colr)
			}
		}
	}
	return layout.Dimensions{Size: pixelTextSize(s, cellPx)}
}

// ledFrame draws a black bezel + dot-matrix display, like a 90s rack unit.
func ledFrame(gtx layout.Context, s string, cell unit.Dp) layout.Dimensions {
	cellPx := gtx.Dp(cell)
	pad := gtx.Dp(4)
	sz := pixelTextSize(strings.ToUpper(s), cellPx)
	total := image.Pt(sz.X+2*pad, sz.Y+2*pad)
	fillRect(gtx, 0, 0, float32(total.X), float32(total.Y), rgb(0x6F6D66))
	fillRect(gtx, 2, 2, float32(total.X-4), float32(total.Y-4), rgb(0x050705))
	offset(gtx, pad, pad, func() { drawPixelText(gtx, s, cellPx) })
	return layout.Dimensions{Size: total}
}

// ---------------------------------------------------------------------------
// Knob widget — shared interaction engine, five visual styles
// ---------------------------------------------------------------------------

type knobStyle int

const (
	styleAmp knobStyle = iota
	styleVintage
	style90s
	styleSwitch
	styleStudio
)

type knob struct {
	style  knobStyle
	label  string
	value  float32
	def    float32
	steps  int // > 0: snap to this many positions (stepped switch)
	active bool

	dragging  bool
	pointerID pointer.ID
	start     f32.Point
	startVal  float32
	lastPress time.Duration
}

func newKnob(style knobStyle, label string, val, def float32, steps int) *knob {
	return &knob{style: style, label: label, value: val, def: def, steps: steps, lastPress: -time.Hour}
}

func (k *knob) stepSize() float32 {
	if k.steps <= 1 {
		return 0
	}
	return (maxVal - minVal) / float32(k.steps-1)
}

func (k *knob) set(v float32) {
	v = clampf(v, minVal, maxVal)
	if s := k.stepSize(); s > 0 {
		v = float32(math.Round(float64(v/s))) * s
	}
	k.value = v
}

func (k *knob) processEvents(gtx layout.Context) {
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
				k.set(k.def) // double-click resets
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
			sens := (maxVal - minVal) / 220 // value units per dragged pixel
			if e.Modifiers&key.ModShift != 0 {
				sens /= 4 // fine
			}
			dx := e.Position.X - k.start.X
			dy := e.Position.Y - k.start.Y
			k.set(k.startVal + (dx-dy)*sens) // up OR right increases
		case pointer.Release, pointer.Cancel:
			k.dragging = false
		case pointer.Scroll:
			// One wheel notch is ~120 on Windows (and Scroll.Y is negative
			// when scrolling up, see Gio's platform drivers).
			step := float32(0.004)
			if e.Modifiers&key.ModShift != 0 {
				step = 0.0012 // fine
			}
			k.set(k.value - e.Scroll.Y*step)
		}
	}
}

func (k *knob) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	size := gtx.Dp(knobSize)
	gtx.Constraints = layout.Exact(image.Pt(size, size))

	// Declare the knob's tag at its clip area, then drain its events.
	cs := clip.Rect{Max: image.Pt(size, size)}.Push(gtx.Ops)
	event.Op(gtx.Ops, k)
	cs.Pop()
	k.processEvents(gtx)

	cx, cy := float64(size)/2, float64(size)/2
	switch k.style {
	case styleAmp:
		k.drawAmp(gtx, cx, cy)
	case styleVintage:
		k.drawVintage(gtx, cx, cy)
	case style90s:
		k.draw90s(gtx, cx, cy)
	case styleSwitch:
		k.drawSwitch(gtx, cx, cy)
	case styleStudio:
		k.drawStudio(gtx, cx, cy)
	}
	return layout.Dimensions{Size: image.Pt(size, size)}
}

// ---- style 1: MK-II metal cap + amber LED arc ----

func (k *knob) drawAmp(gtx layout.Context, cx, cy float64) {
	angle := angleOf(k.value)
	amber := rgb(0xFFB14D)
	strokeArc(gtx, cx, cy, 46, minAngle, maxAngle, 5, rgb(0x34302A))
	strokeArc(gtx, cx, cy, 46, minAngle, angle, 5, amber)
	drawTicks(gtx, cx, cy, 52, 59, 40, 4, rgb(0x6B655C), amber, k.value)
	fillCircleGradient(gtx, cx, cy, 38, rgb(0x4A4640), rgb(0x141310))
	strokeCircle(gtx, cx, cy, 38, 1, rgb(0x0A0A08))
	fillCircleGradient(gtx, cx, cy, 29, rgb(0x5D5C60), rgb(0x0C0B0E))
	strokeCircle(gtx, cx, cy, 29, 1.5, rgb(0x000000))
	rotated(gtx, cx, cy, angle, func() {
		strokeLine(gtx, f32.Pt(float32(cx), float32(cy-12)), f32.Pt(float32(cx), float32(cy-25)), 3.5, rgb(0xF4F1EA))
	})
}

// ---- style 2: vintage chicken-head (cream) ----

func (k *knob) drawVintage(gtx layout.Context, cx, cy float64) {
	angle := angleOf(k.value)
	gold := rgb(0xC98F2E)
	strokeArc(gtx, cx, cy, 40, minAngle, maxAngle, 4, rgb(0xCDBF9D))
	strokeArc(gtx, cx, cy, 40, minAngle, angle, 4, gold)
	drawTicks(gtx, cx, cy, 46, 52, 21, 4, rgb(0x8A6F42), gold, k.value)
	fillCircleGradient(gtx, cx, cy, 27, rgb(0xFDF6E0), rgb(0xCBB98D))
	strokeCircle(gtx, cx, cy, 27, 1.5, rgb(0xB9A678))
	rotated(gtx, cx, cy, angle, func() {
		fillPolygon(gtx, []f32.Point{
			f32.Pt(float32(cx-7), float32(cy+6)),
			f32.Pt(float32(cx), float32(cy-24)),
			f32.Pt(float32(cx+7), float32(cy+6)),
		}, rgb(0x6B4A2A))
	})
	fillCircle(gtx, cx, cy, 8, rgb(0x4A3118))
	fillCircle(gtx, cx, cy, 3, rgb(0x8A6F42))
}

// ---- style 3: DC-90 flat plastic (90s low-tech) ----

func (k *knob) draw90s(gtx layout.Context, cx, cy float64) {
	angle := angleOf(k.value)
	drawTicks(gtx, cx, cy, 40, 46, 11, 1, rgb(0x3C3B38), rgb(0x2F7F3A), k.value)
	fillCircleGradient(gtx, cx, cy, 33, rgb(0x56575C), rgb(0x26272B))
	strokeCircle(gtx, cx, cy, 33, 2, rgb(0x1C1D20))
	strokeCircle(gtx, cx, cy, 25, 1.5, rgb(0x525359))
	rotated(gtx, cx, cy, angle, func() {
		fillRect(gtx, float32(cx-3), float32(cy-31), 6, 16, rgb(0x101113))
	})
	fillCircle(gtx, cx, cy, 4, rgb(0x101113))
}

// ---- style 4: SELECT-12 stepped aluminum switch ----

func (k *knob) drawSwitch(gtx layout.Context, cx, cy float64) {
	angle := angleOf(k.value)
	// 12 detent ticks
	for i := 0; i < 12; i++ {
		a := minAngle + (maxAngle-minAngle)*float64(i)/11
		col := rgb(0xC6CAD0)
		if a <= angle+0.01 {
			col = rgb(0xFFD98A)
		}
		strokeLine(gtx, polarPt(cx, cy, 44, a), polarPt(cx, cy, 50, a), 2.5, col)
	}
	// knurled rim
	for i := 0; i < 28; i++ {
		a := float64(i) * 360 / 28
		strokeLine(gtx, polarPt(cx, cy, 30, a), polarPt(cx, cy, 35, a), 1.4, rgb(0x4C4F55))
	}
	fillCircleGradient(gtx, cx, cy, 30, rgb(0xE2E5E9), rgb(0x585B60))
	strokeCircle(gtx, cx, cy, 30, 1.5, rgb(0x3A3D42))
	strokeCircle(gtx, cx, cy, 21, 1, rgb(0x8F9298))
	rotated(gtx, cx, cy, angle, func() {
		strokeLine(gtx, f32.Pt(float32(cx), float32(cy-8)), f32.Pt(float32(cx), float32(cy-27)), 4, rgb(0x23252A))
	})
	fillCircle(gtx, cx, cy, 5, rgb(0x2B2D32))
}

// ---- style 5: studio matte black + cyan LED dot ring ----

func (k *knob) drawStudio(gtx layout.Context, cx, cy float64) {
	angle := angleOf(k.value)
	for i := 0; i < 24; i++ {
		a := minAngle + (maxAngle-minAngle)*float64(i)/23
		col := rgb(0x20242A)
		if a <= angle+0.01 {
			col = rgb(0x7FE8FF)
		}
		p := polarPt(cx, cy, 40, a)
		fillCircle(gtx, float64(p.X), float64(p.Y), 2.6, col)
	}
	fillCircleGradient(gtx, cx, cy, 27, rgb(0x23262B), rgb(0x0A0B0D))
	strokeCircle(gtx, cx, cy, 27, 1, rgb(0x000000))
	strokeCircle(gtx, cx, cy, 21, 0.8, rgb(0x2C3036))
	rotated(gtx, cx, cy, angle, func() {
		strokeLine(gtx, f32.Pt(float32(cx), float32(cy-9)), f32.Pt(float32(cx), float32(cy-23)), 2.5, rgb(0xF0F4F7))
	})
}

// ---------------------------------------------------------------------------
// Envelope editor — four points, movable in x and y
// ---------------------------------------------------------------------------

type envTag struct{ i int }

type envelope struct {
	points   [4]f32.Point // normalized x (0..1), y (0..1)
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

const (
	envGap  = 0.03 // min distance between neighbouring points
	envMinX = 0.02
	envMaxX = 0.98
)

func newEnvelope() *envelope {
	e := &envelope{lastPress: -time.Hour}
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

func (e *envelope) clampPoint(i int) {
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

func (e *envelope) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	w := gtx.Constraints.Max.X
	h := gtx.Dp(170)
	if w < gtx.Dp(120) {
		w = gtx.Dp(120)
	}
	gtx.Constraints = layout.Exact(image.Pt(w, h))

	padX := float32(gtx.Dp(8))
	padT := float32(gtx.Dp(8))
	padB := float32(gtx.Dp(14))
	iw := float32(w) - 2*padX
	ih := float32(h) - padT - padB
	baseY := padT + ih
	px := func(p f32.Point) f32.Point {
		return f32.Pt(padX+p.X*iw, padT+(1-p.Y)*ih)
	}

	// Screen chrome.
	fillRect(gtx, 0, 0, float32(w), float32(h), rgb(0x6F6D66))
	fillRect(gtx, 2, 2, float32(w)-4, float32(h)-4, rgb(0x030503))

	// Grid.
	for i := 0; i <= 8; i++ {
		x := padX + iw*float32(i)/8
		strokeLine(gtx, f32.Pt(x, padT), f32.Pt(x, baseY), 1, rgb(0x0B3315))
	}
	for i := 0; i <= 4; i++ {
		y := padT + ih*float32(i)/4
		strokeLine(gtx, f32.Pt(padX, y), f32.Pt(padX+iw, y), 1, rgb(0x0B3315))
	}

	// Area fill + envelope line.
	coords := []f32.Point{f32.Pt(padX, baseY)}
	for _, p := range e.points {
		coords = append(coords, px(p))
	}
	area := append([]f32.Point{}, coords...)
	area = append(area, f32.Pt(coords[len(coords)-1].X, baseY), f32.Pt(padX, baseY))
	areaCol := rgb(0x39FF70)
	areaCol.A = 18
	fillPolygon(gtx, area, areaCol)
	strokePolyline(gtx, coords, 2.5, rgb(0x39FF70))

	// Origin marker.
	fillRect(gtx, padX-3, baseY-3, 6, 6, rgb(0x2F7F3A))

	// Handles: hit area first (so events work), then visuals + labels.
	names := [4]string{"A", "D", "S", "R"}
	for i := range e.points {
		p := px(e.points[i])
		cs := clip.Rect{
			Min: image.Pt(int(p.X-14), int(p.Y-14)),
			Max: image.Pt(int(p.X+14), int(p.Y+14)),
		}.Push(gtx.Ops)
		event.Op(gtx.Ops, &e.tags[i])
		cs.Pop()
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
				e.points[i].X = e.startVal.X + (pe.Position.X-e.start.X)/iw
				e.points[i].Y = e.startVal.Y - (pe.Position.Y-e.start.Y)/ih
				e.clampPoint(i)
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
		// Letter label, above (below when near the top edge).
		ty := p.Y - 14
		if e.points[i].Y > 0.8 {
			ty = p.Y + 8
		}
		offset(gtx, int(p.X-3), int(ty), func() {
			lbl := material.Label(th, unit.Sp(9), names[i])
			lbl.Color = rgb(0x2F7F3A)
			lbl.Layout(gtx)
		})
	}
	return layout.Dimensions{Size: image.Pt(w, h)}
}

func (e *envelope) readout() string {
	p := e.points[e.active]
	name := names[e.active]
	return fmt.Sprintf("%s %03d %03d", name,
		int(math.Round(float64(p.X*99))), int(math.Round(float64(p.Y*99))))
}

var names = [4]string{"A", "D", "S", "R"}

// ---------------------------------------------------------------------------
// Dropdown — custom-drawn, styled per unit (press opens, press selects)
// ---------------------------------------------------------------------------

type ddStyle struct {
	bg, border, text, selBG, selText color.NRGBA
	pixel                            bool // render text with the dot-matrix font
}

type ddTag struct{}

type dropdown struct {
	name    string
	options []string
	style   ddStyle
	value   int
	open    bool
	tags    []ddTag
}

func newDropdown(name string, opts []string, st ddStyle) *dropdown {
	return &dropdown{name: name, options: opts, style: st, tags: make([]ddTag, len(opts)+1)}
}

func (d *dropdown) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	st := d.style
	w := gtx.Dp(200)
	if max := gtx.Constraints.Max.X; w > max {
		w = max
	}
	bh := gtx.Dp(30)

	// --- button events + visuals ---
	cs := clip.Rect{Max: image.Pt(w, bh)}.Push(gtx.Ops)
	event.Op(gtx.Ops, &d.tags[0])
	cs.Pop()
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: &d.tags[0], Kinds: pointer.Press})
		if !ok {
			break
		}
		if e, ok := ev.(pointer.Event); ok && e.Kind == pointer.Press && e.Buttons == pointer.ButtonPrimary {
			d.open = !d.open
		}
	}
	paint.FillShape(gtx.Ops, st.border, clip.UniformRRect(image.Rect(0, 0, w, bh), 5).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, st.bg, clip.UniformRRect(image.Rect(1, 1, w-2, bh-2), 4).Op(gtx.Ops))
	ay := bh / 2
	fillPolygon(gtx, []f32.Point{
		f32.Pt(float32(w-20), float32(ay-3)),
		f32.Pt(float32(w-8), float32(ay-3)),
		f32.Pt(float32(w-14), float32(ay+3)),
	}, st.text)

	cur := d.options[d.value]
	if st.pixel {
		cell := gtx.Dp(2)
		tw := pixelTextSize(cur, cell).X
		offset(gtx, (w-tw)/2, (bh-7*cell)/2, func() { drawPixelText(gtx, cur, cell) })
	} else {
		offset(gtx, 0, (bh-gtx.Dp(12))/2, func() {
			labelCentered(gtx, th, strings.ToUpper(d.name+": "+cur), st.text, unit.Sp(12))
		})
	}

	height := bh
	// --- open list (laid out in flow, so nothing is ever covered) ---
	if d.open {
		ih := gtx.Dp(26)
		for i, opt := range d.options {
			y := height + gtx.Dp(2)
			// row events
			cs := clip.Rect{Min: image.Pt(0, y), Max: image.Pt(w, y+ih)}.Push(gtx.Ops)
			event.Op(gtx.Ops, &d.tags[i+1])
			cs.Pop()
			for {
				ev, ok := gtx.Event(pointer.Filter{Target: &d.tags[i+1], Kinds: pointer.Press})
				if !ok {
					break
				}
				if e, ok := ev.(pointer.Event); ok && e.Kind == pointer.Press && e.Buttons == pointer.ButtonPrimary {
					d.value = i
					d.open = false
				}
			}
			rowCol := st.bg
			if i == d.value {
				rowCol = st.selBG
			}
			paint.FillShape(gtx.Ops, rowCol, clip.Rect{Min: image.Pt(0, y), Max: image.Pt(w, y+ih)}.Op())
			offset(gtx, 0, y, func() {
				if st.pixel {
					cell := gtx.Dp(2)
					tw := pixelTextSize(opt, cell).X
					offset(gtx, (w-tw)/2, (ih-7*cell)/2, func() { drawPixelText(gtx, opt, cell) })
				} else {
					col := st.text
					if i == d.value {
						col = st.selText
					}
					offset(gtx, 0, (ih-gtx.Dp(12))/2, func() {
						labelCentered(gtx, th, opt, col, unit.Sp(12))
					})
				}
			})
			height = y + ih
		}
	}
	return layout.Dimensions{Size: image.Pt(w, height)}
}

// labelCentered draws a label horizontally centered within the available width.
func labelCentered(gtx layout.Context, th *material.Theme, txt string, col color.NRGBA, size unit.Sp) layout.Dimensions {
	lbl := material.Label(th, size, txt)
	lbl.Color = col
	macro := op.Record(gtx.Ops)
	dims := lbl.Layout(gtx)
	macro.Stop() // discard measurement pass
	dx := (gtx.Constraints.Max.X - dims.Size.X) / 2
	if dx < 0 {
		dx = 0
	}
	var out layout.Dimensions
	offset(gtx, dx, 0, func() { out = lbl.Layout(gtx) })
	return out
}

// ---------------------------------------------------------------------------
// Layout helpers
// ---------------------------------------------------------------------------

// vstack lays out widgets vertically with a fixed gap between them.
func vstack(gtx layout.Context, gap unit.Dp, ws ...layout.Widget) layout.Dimensions {
	children := make([]layout.FlexChild, len(ws))
	for i, w := range ws {
		w := w
		if i > 0 {
			inner := w
			w = func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: gap}.Layout(gtx, inner)
			}
		}
		children[i] = layout.Rigid(w)
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func centered(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Center.Layout(gtx, w)
}

// panel paints a rounded unit panel behind the given content and gives the
// content the full available width.
func panel(gtx layout.Context, th *material.Theme, bgTop, bgBot, border, titleCol, subCol color.NRGBA, title, sub string, content layout.Widget) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	inset := layout.UniformInset(unit.Dp(14))
	macro := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		header := layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(13), title)
				lbl.Color = titleCol
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(9), sub)
				lbl.Color = subCol
				return lbl.Layout(gtx)
			}),
		)
		return vstack(gtx, unit.Dp(12), func(layout.Context) layout.Dimensions { return header }, content)
	})
	call := macro.Stop()

	r := gtx.Dp(12)
	paint.FillShape(gtx.Ops, border, clip.UniformRRect(image.Rect(0, 0, dims.Size.X, dims.Size.Y), r).Op(gtx.Ops))
	stack := clip.UniformRRect(image.Rect(1, 1, dims.Size.X-1, dims.Size.Y-1), r-1).Push(gtx.Ops)
	paint.LinearGradientOp{
		Stop1:  f32.Pt(0, 0),
		Color1: bgTop,
		Stop2:  f32.Pt(0, float32(dims.Size.Y)),
		Color2: bgBot,
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Pop()
	call.Add(gtx.Ops)
	return dims
}

// ---------------------------------------------------------------------------
// Application UI
// ---------------------------------------------------------------------------

type ui struct {
	tubeVol, tubeTone  *knob
	ampGain, ampMaster *knob
	dcIn, dcMix        *knob
	selPreset          *knob
	stMon, stBlend     *knob
	env                *envelope

	ddChannel, ddMode   *dropdown
	ddSource, ddEnvMode *dropdown
	ddBank, ddOutput    *dropdown
}

func newUI() *ui {
	return &ui{
		tubeVol:   newKnob(styleVintage, "Volume", 6.0, 5.0, 0),
		tubeTone:  newKnob(styleVintage, "Tone", 4.0, 5.0, 0),
		ampGain:   newKnob(styleAmp, "Gain", 6.5, 6.5, 0),
		ampMaster: newKnob(styleAmp, "Master", 4.5, 5.0, 0),
		dcIn:      newKnob(style90s, "Input", 5.0, 5.0, 0),
		dcMix:     newKnob(style90s, "Mix", 3.0, 5.0, 0),
		selPreset: newKnob(styleSwitch, "Preset", 2.73, 2.73, 12),
		stMon:     newKnob(styleStudio, "Monitor", 7.0, 7.0, 0),
		stBlend:   newKnob(styleStudio, "Blend", 5.0, 5.0, 0),
		env:       newEnvelope(),
		ddChannel: newDropdown("Channel", []string{"Normal", "Bright", "Vibrato"}, ddStyle{
			bg: rgb(0xF3ECDA), border: rgb(0x9C8C64), text: rgb(0x4D3A26),
			selBG: rgb(0xDCCFAE), selText: rgb(0x4D3A26)}),
		ddMode: newDropdown("Mode", []string{"Clean", "Crunch", "Lead", "Burn"}, ddStyle{
			bg: rgb(0x23211C), border: rgb(0x0E0D0B), text: rgb(0xFFB14D),
			selBG: rgb(0x37332C), selText: rgb(0xFFB14D)}),
		ddSource: newDropdown("Source", []string{"ANALOG", "DIGITAL", "SPDIF", "TOSLINK"}, ddStyle{
			bg: rgb(0x050705), border: rgb(0x6F6D66), text: rgb(ledLit),
			selBG: rgb(0x0D2413), selText: rgb(ledLit), pixel: true}),
		ddEnvMode: newDropdown("Mode", []string{"GATE", "TRIG", "LOOP", "ONCE"}, ddStyle{
			bg: rgb(0x050705), border: rgb(0x6F6D66), text: rgb(ledLit),
			selBG: rgb(0x0D2413), selText: rgb(ledLit), pixel: true}),
		ddBank: newDropdown("Bank", []string{"Bank A", "Bank B", "Bank C", "Bank D"}, ddStyle{
			bg: rgb(0x3F4247), border: rgb(0x24262B), text: rgb(0xE8EAEE),
			selBG: rgb(0x565A60), selText: rgb(0xFFD98A)}),
		ddOutput: newDropdown("Output", []string{"Monitor A", "Monitor B", "Phones", "Both"}, ddStyle{
			bg: rgb(0x0C0D0F), border: rgb(0x2C2F34), text: rgb(0xDFE4E9),
			selBG: rgb(0x16181C), selText: rgb(0x7FE8FF)}),
	}
}

// knobColumn composes a knob with its label and readout.
func knobColumn(gtx layout.Context, th *material.Theme, k *knob, labelCol, roCol color.NRGBA, roKind int) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Dp(knobSize)
	gtx.Constraints.Max.X = gtx.Dp(knobSize)

	var readout layout.Widget
	switch roKind {
	case roText:
		readout = func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(13), fmt.Sprintf("%.1f", k.value))
			lbl.Color = roCol
			return lbl.Layout(gtx)
		}
	case roLed:
		readout = func(gtx layout.Context) layout.Dimensions {
			return centered(gtx, func(gtx layout.Context) layout.Dimensions {
				return ledFrame(gtx, fmt.Sprintf("%03d", int(math.Round(float64(k.value*10)))), unit.Dp(2.6))
			})
		}
	case roStep:
		readout = func(gtx layout.Context) layout.Dimensions {
			idx := int(math.Round(float64(k.value/k.stepSize()))) + 1
			lbl := material.Label(th, unit.Sp(13), fmt.Sprintf("STEP %02d", idx))
			lbl.Color = roCol
			return lbl.Layout(gtx)
		}
	default:
		readout = func(layout.Context) layout.Dimensions { return layout.Dimensions{} }
	}

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions { return k.Layout(gtx, th) }),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(10), strings.ToUpper(k.label))
			lbl.Color = labelCol
			return lbl.Layout(gtx)
		}),
		layout.Rigid(readout),
	)
}

const (
	roNone = iota
	roText
	roLed
	roStep
)

func knobRow(gtx layout.Context, cols ...layout.Widget) layout.Dimensions {
	children := make([]layout.FlexChild, len(cols))
	for i, c := range cols {
		children[i] = layout.Rigid(c)
	}
	return layout.Flex{Spacing: layout.SpaceAround, Alignment: layout.Middle}.Layout(gtx, children...)
}

func (u *ui) layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.Fill(gtx.Ops, rgb(0x0A0A0D))
	return layout.UniformInset(unit.Dp(18)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		colA := func(gtx layout.Context) layout.Dimensions {
			return u.colA(gtx, th)
		}
		colB := func(gtx layout.Context) layout.Dimensions {
			return u.colB(gtx, th)
		}
		grid := func(gtx layout.Context) layout.Dimensions {
			if gtx.Constraints.Max.X >= gtx.Dp(780) {
				return layout.Flex{Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Flexed(0.46, colA),
					layout.Flexed(0.46, colB),
				)
			}
			return vstack(gtx, unit.Dp(16), colA, colB)
		}
		return vstack(gtx, unit.Dp(18),
			func(gtx layout.Context) layout.Dimensions {
				return vstack(gtx, unit.Dp(2),
					func(gtx layout.Context) layout.Dimensions {
						return centered(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(18), "KNOB-LAB")
							lbl.Color = rgb(0xCFC4AC)
							return lbl.Layout(gtx)
						})
					},
					func(gtx layout.Context) layout.Dimensions {
						return centered(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(9), "GIO PORT — KNOBS · ENVELOPE · DROPDOWNS")
							lbl.Color = rgb(0x7A746A)
							return lbl.Layout(gtx)
						})
					},
				)
			},
			grid,
			func(gtx layout.Context) layout.Dimensions {
				return centered(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(10),
						"Drag knobs up/down or left/right · wheel = steps · Shift = fine · double-click resets · envelope points move in x/y")
					lbl.Color = rgb(0x827A6D)
					return lbl.Layout(gtx)
				})
			},
		)
	})
}

func (u *ui) colA(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return vstack(gtx, unit.Dp(16),
		func(gtx layout.Context) layout.Dimensions {
			return panel(gtx, th, rgb(0xEFE7D2), rgb(0xD9CEA8), rgb(0xA89870), rgb(0x5A4632), rgb(0x7C6A4D),
				"TUBE COMBO", "CHICKEN-HEAD · 1959", func(gtx layout.Context) layout.Dimensions {
					return vstack(gtx, unit.Dp(12),
						func(gtx layout.Context) layout.Dimensions {
							return knobRow(gtx,
								func(gtx layout.Context) layout.Dimensions {
									return knobColumn(gtx, th, u.tubeVol, rgb(0x5A4632), rgb(0x8A5A20), roText)
								},
								func(gtx layout.Context) layout.Dimensions {
									return knobColumn(gtx, th, u.tubeTone, rgb(0x5A4632), rgb(0x8A5A20), roText)
								},
							)
						},
						func(gtx layout.Context) layout.Dimensions {
							return centered(gtx, func(gtx layout.Context) layout.Dimensions {
								return u.ddChannel.Layout(gtx, th)
							})
						},
					)
				})
		},
		func(gtx layout.Context) layout.Dimensions {
			return panel(gtx, th, rgb(0x26231F), rgb(0x1B1916), rgb(0x0E0D0B), rgb(0xD8CDB2), rgb(0x8D8474),
				"MK-II TONE STACK", "METAL CAP · LED ARC", func(gtx layout.Context) layout.Dimensions {
					return vstack(gtx, unit.Dp(12),
						func(gtx layout.Context) layout.Dimensions {
							return knobRow(gtx,
								func(gtx layout.Context) layout.Dimensions {
									return knobColumn(gtx, th, u.ampGain, rgb(0xB3A893), rgb(0xFFB14D), roText)
								},
								func(gtx layout.Context) layout.Dimensions {
									return knobColumn(gtx, th, u.ampMaster, rgb(0xB3A893), rgb(0xFFB14D), roText)
								},
							)
						},
						func(gtx layout.Context) layout.Dimensions {
							return centered(gtx, func(gtx layout.Context) layout.Dimensions {
								return u.ddMode.Layout(gtx, th)
							})
						},
					)
				})
		},
		func(gtx layout.Context) layout.Dimensions {
			return panel(gtx, th, rgb(0x53565C), rgb(0x35383D), rgb(0x22242A), rgb(0xE8EAEE), rgb(0xA7ABB3),
				"SELECT-12", "STEPPED SWITCH · DETENTS", func(gtx layout.Context) layout.Dimensions {
					return vstack(gtx, unit.Dp(12),
						func(gtx layout.Context) layout.Dimensions {
							return knobRow(gtx, func(gtx layout.Context) layout.Dimensions {
								return knobColumn(gtx, th, u.selPreset, rgb(0xD5D8DE), rgb(0xFFD98A), roStep)
							})
						},
						func(gtx layout.Context) layout.Dimensions {
							return centered(gtx, func(gtx layout.Context) layout.Dimensions {
								return u.ddBank.Layout(gtx, th)
							})
						},
					)
				})
		},
	)
}

func (u *ui) colB(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return vstack(gtx, unit.Dp(16),
		func(gtx layout.Context) layout.Dimensions {
			return panel(gtx, th, rgb(0xD3D0C6), rgb(0xB5B2A9), rgb(0x8F8D86), rgb(0x2E2E2C), rgb(0x55544F),
				"DC-90 SIGNAL PATH", "90S LOW-TECH · DOT MATRIX", func(gtx layout.Context) layout.Dimensions {
					return vstack(gtx, unit.Dp(12),
						func(gtx layout.Context) layout.Dimensions {
							return knobRow(gtx,
								func(gtx layout.Context) layout.Dimensions {
									return knobColumn(gtx, th, u.dcIn, rgb(0x2E2E2C), rgb(0x39FF70), roLed)
								},
								func(gtx layout.Context) layout.Dimensions {
									return knobColumn(gtx, th, u.dcMix, rgb(0x2E2E2C), rgb(0x39FF70), roLed)
								},
							)
						},
						func(gtx layout.Context) layout.Dimensions {
							return centered(gtx, func(gtx layout.Context) layout.Dimensions {
								return u.ddSource.Layout(gtx, th)
							})
						},
					)
				})
		},
		func(gtx layout.Context) layout.Dimensions {
			return panel(gtx, th, rgb(0xD3D0C6), rgb(0xB5B2A9), rgb(0x8F8D86), rgb(0x2E2E2C), rgb(0x55544F),
				"ENV-4 GENERATOR", "90S LOW-TECH · ENVELOPE", func(gtx layout.Context) layout.Dimensions {
					return vstack(gtx, unit.Dp(10),
						func(gtx layout.Context) layout.Dimensions {
							return u.env.Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return centered(gtx, func(gtx layout.Context) layout.Dimensions {
								return ledFrame(gtx, u.env.readout(), unit.Dp(2.4))
							})
						},
						func(gtx layout.Context) layout.Dimensions {
							return centered(gtx, func(gtx layout.Context) layout.Dimensions {
								return u.ddEnvMode.Layout(gtx, th)
							})
						},
					)
				})
		},
		func(gtx layout.Context) layout.Dimensions {
			return panel(gtx, th, rgb(0x1C1E22), rgb(0x131518), rgb(0x000000), rgb(0xEEF1F4), rgb(0x6F7680),
				"STUDIO ONE", "MATTE · LED RING", func(gtx layout.Context) layout.Dimensions {
					return vstack(gtx, unit.Dp(12),
						func(gtx layout.Context) layout.Dimensions {
							return knobRow(gtx,
								func(gtx layout.Context) layout.Dimensions {
									return knobColumn(gtx, th, u.stMon, rgb(0x9AA1AB), rgb(0x7FE8FF), roText)
								},
								func(gtx layout.Context) layout.Dimensions {
									return knobColumn(gtx, th, u.stBlend, rgb(0x9AA1AB), rgb(0x7FE8FF), roText)
								},
							)
						},
						func(gtx layout.Context) layout.Dimensions {
							return centered(gtx, func(gtx layout.Context) layout.Dimensions {
								return u.ddOutput.Layout(gtx, th)
							})
						},
					)
				})
		},
	)
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Knob Lab — Gio Port"), app.Size(unit.Dp(1180), unit.Dp(920)))
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	u := newUI()
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			u.layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}
