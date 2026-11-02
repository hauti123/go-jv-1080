package widget

import (
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

// smokeGtx returns a context with no input source (no events) at 1 px/dp.
func smokeGtx(w, h int) layout.Context {
	return layout.Context{
		Ops:         new(op.Ops),
		Constraints: layout.Exact(image.Pt(w, h)),
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
	}
}

func TestLayoutSmoke(t *testing.T) {
	th := DefaultTheme()
	gtx := smokeGtx(780, 600)

	k := NewKnob(KnobPan, "PAN", "patch/temp/common/pan", -64, 63)
	if d := k.Layout(gtx, &th); d.Size.X <= 0 || d.Size.Y <= 0 {
		t.Errorf("knob dims = %v", d.Size)
	}
	sw := NewSwitchKnob("WAVE GRP", "patch/temp/tone/1", []string{"INT-A", "INT-B", "EXP-A", "EXP-B", "EXP-C", "EXP-D"})
	sw.Set(4)
	if d := sw.Layout(gtx, &th); d.Size.X <= 0 || d.Size.Y <= 0 {
		t.Errorf("switch knob dims = %v", d.Size)
	}
	dd := NewDropdown("WAVE", "patch/temp/tone/1/number", []string{"001 A", "002 B", "003 C"})
	if d := dd.Layout(gtx, &th); d.Size.X <= 0 || d.Size.Y <= 0 {
		t.Errorf("dropdown (closed) dims = %v", d.Size)
	}
	dd.open = true
	dd.Set(1)
	if d := dd.Layout(gtx, &th); d.Size.Y <= gtx.Dp(ddButtonDp+titleDp) {
		t.Errorf("dropdown (open) dims = %v, want taller than the box", d.Size)
	}
	dd.open = false
	tb := NewTabs([]string{"COMMON", "TONE 1", "TONE 2", "TONE 3", "TONE 4"})
	if d := tb.Layout(gtx, &th); d.Size.X != 780 || d.Size.Y <= 0 {
		t.Errorf("tabs dims = %v, want width 780", d.Size)
	}
	tf := NewTextField("NAME", "patch/temp/common/name", 12)
	tf.Set("TEST PATCH")
	if d := tf.Layout(gtx, &th); d.Size.X <= 0 || d.Size.Y <= 0 {
		t.Errorf("text field dims = %v", d.Size)
	}
	var tk, lk [4]string
	for i := 0; i < 4; i++ {
		tk[i] = "patch/temp/tone/1/pitch/t" + string(rune('1'+i))
		lk[i] = "patch/temp/tone/1/pitch/l" + string(rune('1'+i))
	}
	env := NewEnvelope("PITCH ENV", tk, lk)
	if d := env.Layout(gtx, &th); d.Size.X <= 0 || d.Size.Y <= 0 {
		t.Errorf("envelope dims = %v", d.Size)
	}
}
