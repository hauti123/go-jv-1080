package widget

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// rgb builds an opaque NRGBA color from a 0xRRGGBB hex value, like the
// gui-proto helper.
func rgb(hex uint32) color.NRGBA {
	return color.NRGBA{R: uint8(hex >> 16), G: uint8(hex >> 8), B: uint8(hex), A: 0xFF}
}

// polarPt returns the point at radius r and angle deg (0 = 12 o'clock,
// increasing clockwise) around (cx, cy) — the gui-proto convention.
func polarPt(cx, cy, r, deg float64) f32.Point {
	a := deg * math.Pi / 180
	return f32.Pt(float32(cx+r*math.Sin(a)), float32(cy-r*math.Cos(a)))
}

func fillCircle(gtx layout.Context, cx, cy, r float64, col color.NRGBA) {
	b := image.Rect(int(cx-r), int(cy-r), int(cx+r+0.5), int(cy+r+0.5))
	paint.FillShape(gtx.Ops, col, clip.Ellipse(b).Op(gtx.Ops))
}

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

// strokeArc draws a stroked arc from angle a0 to a1 (degrees, 0 = 12
// o'clock, clockwise positive) at radius r.
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

// drawTicks draws radial tick marks from a0..a1; ticks up to the value's
// angle use the lit color, the rest the base color. Major ticks every
// majorEvery-th tick are longer and thicker.
func drawTicks(gtx layout.Context, cx, cy, rIn, rOut float64, count, majorEvery int, base, lit color.NRGBA, a0, a1, valueAngle float64) {
	drawTicksRange(gtx, cx, cy, rIn, rOut, count, majorEvery, base, lit, a0, a1, a0, valueAngle)
}

// drawTicksRange draws radial tick marks from a0..a1; ticks inside the lit
// range [litFrom, litTo] use the lit color, the rest the base color.
func drawTicksRange(gtx layout.Context, cx, cy, rIn, rOut float64, count, majorEvery int, base, lit color.NRGBA, a0, a1, litFrom, litTo float64) {
	for i := 0; i <= count; i++ {
		major := majorEvery > 0 && i%majorEvery == 0
		a := a0 + (a1-a0)*float64(i)/float64(count)
		ri, ro, w := rIn, rOut, 1.0
		if !major {
			ri, ro, w = rIn+2, rOut-2, 0.8
		}
		col := base
		if a >= litFrom-0.01 && a <= litTo+0.01 {
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

// Offset runs draw() under a scoped translation.
func Offset(gtx layout.Context, x, y int, draw func()) {
	stack := op.Affine(f32.Affine2D{}.Offset(f32.Pt(float32(x), float32(y)))).Push(gtx.Ops)
	draw()
	stack.Pop()
}
