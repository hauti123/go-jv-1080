package widget

import (
	"image"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"

	"go-jv-1080/internal/util"
)

// Tabs is a row of tab buttons (e.g. COMMON, TONE 1-4 on the patch
// editor screen). The selected tab is highlighted with a lit underline;
// clicking a tab selects it and reports it through OnSelect.
type Tabs struct {
	Labels   []string
	Selected int

	tags []tabTag

	// OnSelect is called with the index of the newly selected tab.
	OnSelect func(index int)
}

// tabTag is the event tag of tab i.
type tabTag struct{ idx int }

// tabDp is the height of the tab row.
const tabDp = 24

// NewTabs returns a tab bar with the given labels, initially with the
// first tab selected.
func NewTabs(labels []string) *Tabs {
	t := &Tabs{Labels: labels}
	t.tags = make([]tabTag, len(labels))
	for i := range t.tags {
		t.tags[i].idx = i
	}
	return t
}

// Select selects a tab programmatically. It does not call OnSelect.
func (t *Tabs) Select(i int) {
	t.Selected = util.Clamp(i, 0, max(len(t.Labels)-1, 0))
}

// selectTab applies a user click on tab i.
func (t *Tabs) selectTab(i int) {
	if i < 0 || i >= len(t.Labels) || i == t.Selected {
		return
	}
	t.Selected = i
	if t.OnSelect != nil {
		t.OnSelect(i)
	}
}

// Layout lays out the tab row and handles its input.
func (t *Tabs) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	w := gtx.Constraints.Max.X
	if w <= 0 {
		w = gtx.Dp(unit.Dp(tabDp * len(t.Labels)))
	}
	h := gtx.Dp(tabDp)
	n := len(t.Labels)
	if n == 0 {
		return layout.Dimensions{Size: image.Pt(w, h)}
	}
	tw := w / n
	for i := range t.Labels {
		x := i * tw
		cs := clip.Rect{Min: image.Pt(x, 0), Max: image.Pt(x+tw, h)}.Push(gtx.Ops)
		event.Op(gtx.Ops, &t.tags[i])
		cs.Pop()
		for {
			ev, ok := gtx.Event(pointer.Filter{Target: &t.tags[i], Kinds: pointer.Press})
			if !ok {
				break
			}
			if e, ok := ev.(pointer.Event); ok && e.Kind == pointer.Press && e.Buttons == pointer.ButtonPrimary {
				t.selectTab(i)
			}
		}
		// tab background: selected tab raised with a lit underline
		if i == t.Selected {
			fillRect(gtx, float32(x), 0, float32(tw), float32(h), th.PanelTop)
			fillRect(gtx, float32(x), float32(h-2), float32(tw), 2, th.Lit)
		} else {
			fillRect(gtx, float32(x), 0, float32(tw), float32(h), th.PanelBot)
			fillRect(gtx, float32(x), float32(h-1), float32(tw), 1, th.PanelEdge)
		}
		// label, centered, shrunk to fit the tab
		lbl := t.Labels[i]
		cell := fitCell(gtx, lbl, tw-10)
		lsz := PixelTextSize(lbl, cell)
		Offset(gtx, x+(tw-lsz.X)/2, (h-lsz.Y)/2, func() { DrawPixelText(gtx, lbl, cell) })
	}
	fillRect(gtx, 0, float32(h-1), float32(w), 1, th.PanelEdge)
	return layout.Dimensions{Size: image.Pt(w, h)}
}
