package widget

import (
	"image"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// TextField is a single-line text input styled like the other widgets
// (dark screen in a metal bezel, mandatory title below), used e.g. for
// the 12-character patch name. It wraps a gioui.org/widget.Editor and
// limits the input to the JV-1080 name character set (ASCII 32-127).
//
// Like the other widgets, a TextField is bound to a device parameter
// through Key, a path-like string resolved solely by the device layer;
// OnChange is called with the full text on every change.
type TextField struct {
	Title  string // mandatory, written below the field
	Key    string // device parameter key this field is bound to
	MaxLen int    // maximum text length (0 = unlimited)

	editor widget.Editor
	mth    *material.Theme
	last   string

	// OnChange is called with the new text whenever the user edits it. It
	// is not called for programmatic Set calls.
	OnChange func(text string)
}

// nameCharset is the character set allowed in a JV-1080 name (ASCII
// 32-127, per doc/patch-edit-mode.md).
const nameCharset = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\x7f"

// fieldDp are the nominal dimensions of the input box.
const (
	fieldWidthDp  = 260
	fieldHeightDp = 22
)

// NewTextField returns a text field with the given title, device
// parameter key and maximum text length (e.g. 12 for a patch name).
func NewTextField(title, key string, maxLen int) *TextField {
	mth := material.NewTheme()
	mth.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	t := &TextField{Title: title, Key: key, MaxLen: maxLen, mth: mth}
	t.editor = widget.Editor{
		SingleLine: true,
		Submit:     true,
		Filter:     nameCharset,
	}
	return t
}

// Text returns the current text.
func (t *TextField) Text() string { return t.editor.Text() }

// Set sets the text programmatically. It does not call OnChange.
func (t *TextField) Set(s string) {
	if t.MaxLen > 0 && len(s) > t.MaxLen {
		s = s[:t.MaxLen]
	}
	t.editor.SetText(s)
	t.last = t.editor.Text()
}

// Layout lays out the text field and handles its input.
func (t *TextField) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	w := gtx.Dp(fieldWidthDp)
	if max := gtx.Constraints.Max.X; w > max && max > 0 {
		w = max
	}
	h := gtx.Dp(fieldHeightDp)
	total := h + gtx.Dp(titleDp)
	gtx.Constraints = layout.Exact(image.Pt(w, total))

	// drain editor events first so changes are detected before the
	// editor layout consumes them
	for {
		ev, ok := t.editor.Update(gtx)
		if !ok {
			break
		}
		if _, isChange := ev.(widget.ChangeEvent); isChange {
			s := t.editor.Text()
			if t.MaxLen > 0 && len(s) > t.MaxLen {
				s = s[:t.MaxLen]
				t.editor.SetText(s)
			}
			if s != t.last {
				t.last = s
				if t.OnChange != nil {
					t.OnChange(s)
				}
			}
		}
	}

	// bezel + dark screen
	fillRect(gtx, 0, 0, float32(w), float32(h), th.PanelEdge)
	fillRect(gtx, 1, 1, float32(w-2), float32(h-2), rgb(0x050705))

	// editor inside the screen area
	es := material.Editor(t.mth, &t.editor, "")
	es.TextSize = unit.Sp(12)
	es.Color = rgb(ledLit)
	sel := rgb(ledLit)
	sel.A = 60
	es.SelectionColor = sel
	pad := gtx.Dp(4)
	child := gtx
	child.Constraints = layout.Exact(image.Pt(w-2*pad, h-2*pad))
	Offset(gtx, pad, pad, func() { es.Layout(child) })

	// mandatory title below the field
	cell := gtx.Dp(FontCellDp)
	ts := PixelTextSize(t.Title, cell)
	Offset(gtx, (w-ts.X)/2, h+1, func() { DrawPixelText(gtx, t.Title, cell) })

	return layout.Dimensions{Size: image.Pt(w, total)}
}
