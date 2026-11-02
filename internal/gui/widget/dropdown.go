package widget

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"

	"go-jv-1080/internal/util"
)

// Dropdown is a custom-drawn drop-down list in the MK-II style of the
// other widgets: a box showing the current entry (dot-matrix text) that
// opens on click; while open, a scrollable list of all entries is laid
// out below the box in flow. Clicking an entry selects it and closes the
// list. The mandatory Title is written below the box.
//
// A Dropdown is bound to a device parameter through Key, a path-like
// string resolved solely by the device layer (e.g. the wave number
// "patch/temp/tone/1/number"); OnChange is called with the selected
// index (0..len(Options)-1).
type Dropdown struct {
	Title   string   // mandatory, written below the box
	Key     string   // device parameter key this dropdown is bound to
	Options []string // entries of the list

	value int  // selected index
	open  bool // list is open
	list  layout.List
	tags  []ddTag // tags[0] = button, tags[i+1] = option i

	// OnChange is called with the new index whenever the user selects an
	// entry. It is not called for programmatic Set calls.
	OnChange func(index int)
}

// ddTag is the event tag of the dropdown button (idx -1) or of option i.
type ddTag struct{ idx int }

const (
	ddWidthDp  = 260 // nominal width of the box
	ddButtonDp = 22  // height of the box
	ddRowDp    = 16  // height of one list row
	ddListDp   = 168 // max height of the open list
)

// NewDropdown returns a dropdown with the given title, device parameter
// key and options, initially closed with the first entry selected.
func NewDropdown(title, key string, options []string) *Dropdown {
	d := &Dropdown{Title: title, Key: key, list: layout.List{Axis: layout.Vertical}}
	d.SetOptions(options)
	return d
}

// SetOptions replaces the options, keeping the current selection clamped
// into the new range. It does not call OnChange.
func (d *Dropdown) SetOptions(options []string) {
	d.Options = options
	d.tags = make([]ddTag, len(options)+1)
	for i := range d.tags {
		d.tags[i].idx = i - 1
	}
	d.Set(d.value)
}

// Value returns the selected index.
func (d *Dropdown) Value() int { return d.value }

// Set selects an index programmatically (clamped to the range). It does
// not call OnChange.
func (d *Dropdown) Set(i int) {
	if len(d.Options) == 0 {
		d.value = 0
		return
	}
	d.value = util.Clamp(i, 0, len(d.Options)-1)
}

// choose applies a user selection: select, close, report through OnChange.
func (d *Dropdown) choose(i int) {
	d.Set(i)
	d.open = false
	if d.OnChange != nil {
		d.OnChange(d.value)
	}
}

// Layout lays out the dropdown box (and the open list below it) and
// handles its input.
func (d *Dropdown) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	w := gtx.Dp(ddWidthDp)
	if max := gtx.Constraints.Max.X; w > max && max > 0 {
		w = max
	}
	bh := gtx.Dp(ddButtonDp)
	boxTotal := bh + gtx.Dp(titleDp) // box + title row

	// button hit area + events
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

	// box chrome: border, background, arrow
	fillRect(gtx, 0, 0, float32(w), float32(bh), th.PanelEdge)
	fillRect(gtx, 1, 1, float32(w-2), float32(bh-2), th.PanelTop)
	ay := float32(bh) / 2
	fillPolygon(gtx, []f32.Point{
		f32.Pt(float32(w-18), ay-3),
		f32.Pt(float32(w-8), ay-3),
		f32.Pt(float32(w-13), ay+3),
	}, th.Label)

	// current entry, centered, shrunk to fit the box
	cur := ""
	if d.value >= 0 && d.value < len(d.Options) {
		cur = d.Options[d.value]
	}
	cell := fitCell(gtx, cur, w-28)
	csz := PixelTextSize(cur, cell)
	Offset(gtx, (w-csz.X)/2, (bh-csz.Y)/2, func() { DrawPixelText(gtx, cur, cell) })

	// mandatory title below the box
	cell = gtx.Dp(FontCellDp)
	ts := PixelTextSize(d.Title, cell)
	Offset(gtx, (w-ts.X)/2, bh+1, func() { DrawPixelText(gtx, d.Title, cell) })

	height := boxTotal
	if d.open && len(d.Options) > 0 {
		lh := gtx.Dp(ddListDp)
		if rem := gtx.Constraints.Max.Y - boxTotal; rem > 0 && lh > rem {
			lh = rem
		}
		child := gtx
		child.Constraints.Min = image.Point{}
		child.Constraints.Max.Y = lh
		d.list.Layout(child, len(d.Options), func(gtx layout.Context, i int) layout.Dimensions {
			return d.layoutRow(gtx, th, w, i)
		})
		height = boxTotal + lh
	}
	return layout.Dimensions{Size: image.Pt(w, height)}
}

// layoutRow lays out one option row of the open list: hit area, events,
// background and entry text.
func (d *Dropdown) layoutRow(gtx layout.Context, th *Theme, w, i int) layout.Dimensions {
	rh := gtx.Dp(ddRowDp)
	cs := clip.Rect{Max: image.Pt(w, rh)}.Push(gtx.Ops)
	event.Op(gtx.Ops, &d.tags[i+1])
	cs.Pop()
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: &d.tags[i+1], Kinds: pointer.Press})
		if !ok {
			break
		}
		if e, ok := ev.(pointer.Event); ok && e.Kind == pointer.Press && e.Buttons == pointer.ButtonPrimary {
			d.choose(i)
		}
	}
	// row background: the selected entry is highlighted
	if i == d.value {
		sel := rgb(ledLit)
		sel.A = 40
		fillRect(gtx, 0, 0, float32(w), float32(rh), sel)
	} else {
		fillRect(gtx, 0, 0, float32(w), float32(rh), rgb(0x101210))
	}
	fillRect(gtx, 0, float32(rh)-1, float32(w), 1, rgb(0x1C1E1C))
	opt := d.Options[i]
	cell := fitCell(gtx, opt, w-12)
	osz := PixelTextSize(opt, cell)
	Offset(gtx, 6, (rh-osz.Y)/2, func() { DrawPixelText(gtx, opt, cell) })
	return layout.Dimensions{Size: image.Pt(w, rh)}
}

// fitCell returns the dot-matrix cell size for s that fits into maxW px,
// at most gtx.Dp(FontCellDp) and at least 1 px.
func fitCell(gtx layout.Context, s string, maxW int) int {
	cell := gtx.Dp(FontCellDp)
	if w1 := PixelTextSize(s, 1).X; w1 > 0 && PixelTextSize(s, cell).X > maxW {
		cell = max(1, maxW/w1)
	}
	return cell
}
