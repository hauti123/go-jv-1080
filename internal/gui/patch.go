package gui

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"

	"go-jv-1080/internal/gui/widget"
	"go-jv-1080/internal/jv1080"
)

// patchEditor is the patch edit screen: a tab bar with one COMMON tab
// (all patch-wide parameters) and one tab per tone (TONE 1-4, all
// per-tone parameters). Each tab shows its parameters as knob grids in
// sections, the wave selection as a 6-position switch knob plus a drop
// down list of the waves, the pitch/TVF envelopes as envelope editors
// and the patch name as a text field. All widgets are bound to
// "patch/temp/..." parameter keys; the key -> SysEx mapping lives in
// internal/jv1080.
type patchEditor struct {
	send func(key string, value any)

	tabs *widget.Tabs
	list layout.List

	knobs map[string]*widget.Knob
	envs  map[string]*widget.Envelope // block key -> envelope, e.g. "patch/temp/tone/1/pitch"
	waves map[int]*waveSel            // tone 1-4 -> wave group switch + wave dropdown
	name  *widget.TextField

	content [5][]section // section list per tab: 0 = common, 1-4 = tone 1-4
}

// waveSel holds the wave group switch knob and the wave number dropdown
// of one tone, plus the last known group/groupid for applyParam.
type waveSel struct {
	knob           *widget.SwitchKnob
	dd             *widget.Dropdown
	group, groupID int
}

// itemKind classifies a layout item of a section.
type itemKind int

const (
	itemKnob itemKind = iota // a knob, looked up in knobs by key
	itemWave                 // the wave group switch + wave dropdown of a tone (key = "patch/temp/tone/N")
	itemEnv                  // an envelope editor, looked up in envs by block key
	itemText                 // the patch name field (key = "patch/temp/common/name")
)

// item is one entry of a section: rows of consecutive knob items are
// packed into a knob row at layout time; all other kinds get a row of
// their own.
type item struct {
	kind itemKind
	key  string
}

// section is a titled group of items within a tab.
type section struct {
	title string
	items []item
}

// nameKey is the parameter key of the patch name.
const nameKey = "patch/temp/common/name"

// newPatchEditor returns the patch edit screen with all widgets created
// and wired to the send callback.
func newPatchEditor(send func(key string, value any)) *patchEditor {
	p := &patchEditor{
		send:  send,
		knobs: make(map[string]*widget.Knob),
		envs:  make(map[string]*widget.Envelope),
		waves: make(map[int]*waveSel),
		list:  layout.List{Axis: layout.Vertical},
	}
	p.tabs = widget.NewTabs([]string{"COMMON", "TONE 1", "TONE 2", "TONE 3", "TONE 4"})
	p.name = widget.NewTextField("NAME", nameKey, 12)
	p.name.OnChange = func(s string) { p.send(nameKey, s) }
	p.buildCommon()
	for t := 1; t <= 4; t++ {
		p.buildTone(t)
	}
	return p
}

// toneKey renders the parameter key of a tone parameter: block "" is
// the tone block itself, otherwise a sub-block ("pitch", "tvf", ...).
func toneKey(t int, block, name string) string {
	if block == "" {
		return fmt.Sprintf("patch/temp/tone/%d/%s", t, name)
	}
	return fmt.Sprintf("patch/temp/tone/%d/%s/%s", t, block, name)
}

// commonKey renders the parameter key of a patch common parameter.
func commonKey(name string) string { return "patch/temp/common/" + name }

// knobTitle derives the knob title from the parameter block and name,
// e.g. block "tone/2/lfo1" name "rate" -> "LFO1 RATE".
func knobTitle(block, name string) string {
	segs := strings.Split(block, "/")
	if len(segs) >= 3 {
		return strings.ToUpper(segs[2] + " " + name)
	}
	return strings.ToUpper(name)
}

// addKnob creates a knob for a raw/bipolar/pan parameter and wires its
// OnChange to send.
func (p *patchEditor) addKnob(param jv1080.Param) {
	key := param.Key
	title := knobTitle(param.Block, key[strings.LastIndex(key, "/")+1:])
	var k *widget.Knob
	switch param.Kind {
	case jv1080.KindPan:
		k = widget.NewKnob(widget.KnobPan, title, key, param.Min, param.Max)
	case jv1080.KindBipolar:
		k = widget.NewKnob(widget.KnobBipolar, title, key, param.Min, param.Max)
	default:
		k = widget.NewKnob(widget.KnobMinMax, title, key, param.Min, param.Max)
	}
	k.OnChange = func(v int) { p.send(key, v) }
	p.knobs[key] = k
}

// addEnvelope creates the 4-point envelope editor for a tone sub-block
// (pitch or tvf) and wires its OnChange to send.
func (p *patchEditor) addEnvelope(t int, block string) {
	base := toneKey(t, block, "")
	base = strings.TrimSuffix(base, "/")
	var tk, lk [4]string
	for i := 0; i < 4; i++ {
		tk[i] = fmt.Sprintf("%s/t%d", base, i+1)
		lk[i] = fmt.Sprintf("%s/l%d", base, i+1)
	}
	e := widget.NewEnvelope(strings.ToUpper(block)+" ENV", tk, lk)
	e.OnChange = func(timeKey string, timeVal int, levelKey string, levelVal int) {
		p.send(timeKey, timeVal)
		p.send(levelKey, levelVal)
	}
	p.envs[base] = e
}

// addWaveSel creates the wave group switch knob (6 positions INT-A,
// INT-B, EXP-A..D) and the wave number dropdown for a tone.
func (p *patchEditor) addWaveSel(t int) {
	base := fmt.Sprintf("patch/temp/tone/%d", t)
	keyGroup := base + "/group"
	keyGroupID := base + "/groupid"
	keyNumber := base + "/number"
	labels := make([]string, len(jv1080.WaveSources))
	for i, s := range jv1080.WaveSources {
		labels[i] = s.Label
	}
	w := &waveSel{}
	w.knob = widget.NewSwitchKnob("WAVE GRP", base, labels)
	w.knob.OnChange = func(pos int) {
		src := jv1080.WaveSourceAt(pos)
		w.group, w.groupID = src.Group, src.GroupID
		p.send(keyGroup, src.Group)
		p.send(keyGroupID, src.GroupID)
		opts := jv1080.WaveOptions(src.Group, src.GroupID)
		w.dd.SetOptions(opts)
		// clamp the wave number into the new group's range
		if n := w.dd.Value() + 1; n > len(opts) {
			w.dd.Set(len(opts) - 1)
			p.send(keyNumber, len(opts))
		}
	}
	w.dd = widget.NewDropdown("WAVE", keyNumber, jv1080.WaveOptions(0, 1))
	w.dd.OnChange = func(i int) { p.send(keyNumber, i+1) }
	p.waves[t] = w
}

// numbered returns the names name1..nameN, e.g. numbered("efxprm", 12).
func numbered(name string, n int) []string {
	out := make([]string, n)
	for i := 1; i <= n; i++ {
		out[i-1] = fmt.Sprintf("%s%d", name, i)
	}
	return out
}

// buildCommon creates the widgets and the COMMON tab sections.
func (p *patchEditor) buildCommon() {
	// all common parameters except the name become knobs
	for _, param := range jv1080.ListParams() {
		if param.Block != "common" || param.Key == nameKey {
			continue
		}
		p.addKnob(param)
	}
	efx := append([]string{"efxtype"}, numbered("efxprm", 12)...)
	efx = append(efx, "efxout", "efxlevel", "efxchorus", "efxreverb",
		"efxctrl1src", "efxctrl1depth", "efxctrl2src", "efxctrl2depth")
	ctrl := []string{"assign", "legato", "bendup", "benddown", "portamento",
		"portmode", "porttype", "portstart", "porttime", "ctrl2src", "ctrl3src",
		"efxhold", "ctrl1hold", "ctrl2hold", "ctrl3hold"}
	defs := []struct {
		title string
		names []string
	}{
		{"NAME", []string{"__name__"}},
		{"PATCH", []string{"level", "pan", "analogfeel", "octave", "stretch", "priority", "velrange", "tempo"}},
		{"STRUCT", []string{"struct12", "booster12", "struct34", "booster34"}},
		{"EFX", efx},
		{"CHORUS", []string{"choruslevel", "chorusrate", "chorusdepth", "choruspredelay", "chorusfeedback", "chorusout"}},
		{"REVERB", []string{"reverbtype", "reverblevel", "reverbtime", "hfdamp", "reverbfeedback"}},
		{"CONTROL", ctrl},
	}
	var secs []section
	for _, d := range defs {
		var items []item
		for _, n := range d.names {
			if n == "__name__" {
				items = append(items, item{kind: itemText, key: nameKey})
				continue
			}
			items = append(items, item{kind: itemKnob, key: commonKey(n)})
		}
		secs = append(secs, section{title: d.title, items: items})
	}
	p.content[0] = p.ensureAllParams(secs, "common")
}

// buildTone creates the widgets and the tab sections of one tone.
func (p *patchEditor) buildTone(t int) {
	p.addWaveSel(t)
	p.addEnvelope(t, "pitch")
	p.addEnvelope(t, "tvf")
	// all tone parameters except the wave group/number and the envelope
	// times/levels become knobs
	skip := map[string]bool{
		toneKey(t, "", "group"): true, toneKey(t, "", "groupid"): true, toneKey(t, "", "number"): true,
	}
	for _, param := range jv1080.ListParams() {
		if !strings.HasPrefix(param.Block, fmt.Sprintf("tone/%d", t)) {
			continue
		}
		name := param.Key[strings.LastIndex(param.Key, "/")+1:]
		if skip[param.Key] || isEnvPoint(param.Block, name) {
			continue
		}
		p.addKnob(param)
	}
	lfo := []string{"wave", "keytrig", "rate", "offset", "delay", "fade", "fadetime", "sync", "pitch", "filter", "amp", "pan"}
	ctrlNames := func() []string {
		var out []string
		for k := 1; k <= 4; k++ {
			out = append(out, fmt.Sprintf("dst%d", k), fmt.Sprintf("dep%d", k))
		}
		return out
	}()
	waveItems := []item{
		{kind: itemKnob, key: toneKey(t, "", "switch")},
		{kind: itemWave, key: fmt.Sprintf("patch/temp/tone/%d", t)},
	}
	for _, n := range []string{"gain", "fxmswitch", "fxmcolor", "fxmdepth", "delaymode", "delaytime"} {
		waveItems = append(waveItems, item{kind: itemKnob, key: toneKey(t, "", n)})
	}
	defs := []struct {
		title string
		block string
		names []string
		env   bool
	}{
		{"WAVE", "", nil, false},
		{"RANGES", "", []string{"velfade", "vellower", "velupper", "keylower", "keyupper"}, false},
		{"RX", "", []string{"redamper", "rxvolume", "rxhold", "rxbend", "rxpan"}, false},
		{"CTRL 1", "ctrl1", ctrlNames, false},
		{"CTRL 2", "ctrl2", ctrlNames, false},
		{"CTRL 3", "ctrl3", ctrlNames, false},
		{"LFO 1", "lfo1", lfo, false},
		{"LFO 2", "lfo2", lfo, false},
		{"PITCH", "pitch", []string{"coarse", "fine", "random", "keyfollow", "depth", "velsense", "vt1", "vt4", "envkeyfollow"}, true},
		{"TVF", "tvf", []string{"filtertype", "cutoff", "cutoffkeyfollow", "resonance", "velres", "envdepth", "velcurve", "velsens", "vt1", "vt4", "envkeyfollow"}, true},
		{"TVA", "tva", []string{"level", "biasdir", "biaspoint", "bias", "envcurve", "envvel", "vt1", "vt4", "envkeyfollow", "pan", "pankeyfollow", "randompan", "altpan", "t1", "t2", "t3", "t4", "l1", "l2", "l3"}, false},
		{"EFX", "efx", []string{"output", "outputlevel", "chorus", "reverb"}, false},
	}
	var secs []section
	for _, d := range defs {
		var items []item
		if d.block == "" && d.names == nil {
			items = waveItems
		} else {
			for _, n := range d.names {
				items = append(items, item{kind: itemKnob, key: toneKey(t, d.block, n)})
			}
		}
		if d.env {
			items = append(items, item{kind: itemEnv, key: strings.TrimSuffix(toneKey(t, d.block, ""), "/")})
		}
		secs = append(secs, section{title: d.title, items: items})
	}
	p.content[t] = p.ensureAllParams(secs, fmt.Sprintf("tone/%d", t))
}

// isEnvPoint reports whether name is an envelope time/level key (t1-t4,
// l1-l4) of an envelope sub-block (pitch/tvf).
func isEnvPoint(block, name string) bool {
	if len(name) != 2 || (name[0] != 't' && name[0] != 'l') || name[1] < '1' || name[1] > '4' {
		return false
	}
	segs := strings.Split(block, "/")
	return len(segs) == 3 && (segs[2] == "pitch" || segs[2] == "tvf")
}

// ensureAllParams appends an "OTHER" section with knobs for all
// parameters of the given block (and its sub-blocks) that are not
// referenced by the sections yet (so that every parameter is editable).
func (p *patchEditor) ensureAllParams(secs []section, block string) []section {
	seen := make(map[string]bool)
	for _, s := range secs {
		for _, it := range s.items {
			seen[it.key] = true
		}
	}
	inBlock := func(b string) bool {
		return b == block || strings.HasPrefix(b, block+"/")
	}
	var missing []item
	for _, param := range jv1080.ListParams() {
		if !inBlock(param.Block) || seen[param.Key] {
			continue
		}
		name := param.Key[strings.LastIndex(param.Key, "/")+1:]
		if param.Key == nameKey || isEnvPoint(param.Block, name) || isWaveParam(param.Block, name) {
			continue
		}
		if _, ok := p.knobs[param.Key]; !ok {
			p.addKnob(param)
		}
		missing = append(missing, item{kind: itemKnob, key: param.Key})
	}
	if len(missing) > 0 {
		secs = append(secs, section{title: "OTHER", items: missing})
	}
	return secs
}

// isWaveParam reports whether name is one of the wave selection
// parameters (group, groupid, number) of a tone root block; they are
// handled by the wave switch knob and dropdown.
func isWaveParam(block, name string) bool {
	if name != "group" && name != "groupid" && name != "number" {
		return false
	}
	segs := strings.Split(block, "/")
	return len(segs) == 2 && segs[0] == "tone"
}

// Layout renders the tab bar and the scrollable content of the selected
// tab.
func (p *patchEditor) Layout(gtx layout.Context, th *widget.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.tabs.Layout(gtx, th)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.layoutContent(gtx, th)
			})
		}),
	)
}

// layoutContent lays out the sections of the selected tab in a scrollable
// list.
func (p *patchEditor) layoutContent(gtx layout.Context, th *widget.Theme) layout.Dimensions {
	secs := p.content[p.tabs.Selected]
	return p.list.Layout(gtx, len(secs), func(gtx layout.Context, i int) layout.Dimensions {
		return p.layoutSection(gtx, th, secs[i])
	})
}

// layoutSection lays out one section: a title bar, then the items as
// knob rows (consecutive knobs packed per row) and full-width rows for
// the other item kinds.
func (p *patchEditor) layoutSection(gtx layout.Context, th *widget.Theme, sec section) layout.Dimensions {
	children := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			cell := gtx.Dp(widget.FontCellDp)
			sz := widget.PixelTextSize(sec.title, cell)
			h := sz.Y + gtx.Dp(unit.Dp(8))
			widget.Offset(gtx, 0, gtx.Dp(unit.Dp(2)), func() { widget.DrawPixelText(gtx, sec.title, cell) })
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, h)}
		}),
	}
	perRow := max(1, gtx.Constraints.Max.X/(gtx.Dp(widget.KnobDp)+gtx.Dp(unit.Dp(12))))
	var row []item
	flush := func() {
		if len(row) == 0 {
			return
		}
		r := append([]item(nil), row...)
		row = nil
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			kids := make([]layout.FlexChild, len(r))
			for i, it := range r {
				it := it
				kids[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.layoutItem(gtx, th, it)
				})
			}
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx, kids...)
		}))
	}
	for _, it := range sec.items {
		if it.kind == itemKnob && len(row) < perRow {
			row = append(row, it)
			continue
		}
		flush()
		it := it
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.layoutItem(gtx, th, it)
			})
		}))
	}
	flush()
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

// layoutItem lays out a single item.
func (p *patchEditor) layoutItem(gtx layout.Context, th *widget.Theme, it item) layout.Dimensions {
	switch it.kind {
	case itemKnob:
		return p.knobs[it.key].Layout(gtx, th)
	case itemText:
		return p.name.Layout(gtx, th)
	case itemEnv:
		return p.envs[it.key].Layout(gtx, th)
	case itemWave:
		t, _ := strconv.Atoi(strings.TrimPrefix(it.key, "patch/temp/tone/"))
		w := p.waves[t]
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return w.knob.Layout(gtx, th) }),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return w.dd.Layout(gtx, th) }),
		)
	}
	return layout.Dimensions{}
}

// initDevice reads the whole temporary patch once and applies every
// parameter to its widget.
func (p *patchEditor) initDevice(c *jv1080.Client) {
	patch, err := c.ReadTemporaryPatch()
	if err != nil {
		return
	}
	for _, param := range jv1080.ListParams() {
		data := paramBytes(patch, param)
		if data == nil {
			continue
		}
		if v, err := param.DecodeValue(data); err == nil {
			p.applyParam(param.Key, v)
		}
	}
}

// paramBytes extracts the raw bytes of a parameter from a patch block.
func paramBytes(patch *jv1080.Patch, param jv1080.Param) []byte {
	size := param.Size()
	if param.Block == "common" {
		if param.Offset+size <= len(patch.Common) {
			return patch.Common[param.Offset : param.Offset+size]
		}
		return nil
	}
	segs := strings.Split(param.Block, "/")
	if len(segs) < 2 || segs[0] != "tone" {
		return nil
	}
	t, err := strconv.Atoi(segs[1])
	if err != nil || t < 1 || t > 4 {
		return nil
	}
	tone := patch.Tones[t-1]
	if param.Offset+size <= len(tone) {
		return tone[param.Offset : param.Offset+size]
	}
	return nil
}

// applyParam applies a device-originated parameter update to the widget
// bound to key, without re-sending it.
func (p *patchEditor) applyParam(key string, value any) {
	if k, ok := p.knobs[key]; ok {
		k.Set(asKnobValue(value))
		return
	}
	if key == nameKey {
		if s, ok := value.(string); ok {
			p.name.Set(s)
		}
		return
	}
	// envelope time/level keys: "<block>/tN" or "<block>/lN"
	if block := key[:strings.LastIndex(key, "/")]; block != "" {
		if e, ok := p.envs[block]; ok {
			name := key[strings.LastIndex(key, "/")+1:]
			if len(name) == 2 && (name[0] == 't' || name[0] == 'l') && name[1] >= '1' && name[1] <= '4' {
				i := int(name[1] - '1')
				tv, lv := e.Point(i)
				if n, ok := asInt(value); ok {
					if name[0] == 't' {
						tv = n
					} else {
						lv = n
					}
				}
				e.SetPoint(i, tv, lv)
			}
			return
		}
	}
	// wave group / group id / number of a tone
	if t, ok := toneNumberOf(key); ok {
		w := p.waves[t]
		base := fmt.Sprintf("patch/temp/tone/%d", t)
		switch key {
		case base + "/number":
			if n, ok := asInt(value); ok {
				w.dd.Set(n - 1)
			}
			return
		case base + "/group":
			if n, ok := asInt(value); ok {
				w.group = n
			}
		case base + "/groupid":
			if n, ok := asInt(value); ok {
				w.groupID = n
			}
		}
		if key == base+"/group" || key == base+"/groupid" {
			if idx := jv1080.WaveSourceIndex(w.group, w.groupID); idx >= 0 {
				w.knob.Set(idx)
				w.dd.SetOptions(jv1080.WaveOptions(w.group, w.groupID))
			}
		}
	}
}

// toneNumberOf extracts the tone number (1-4) from a tone parameter key
// "patch/temp/tone/N/...".
func toneNumberOf(key string) (int, bool) {
	rest, ok := strings.CutPrefix(key, "patch/temp/tone/")
	if !ok {
		return 0, false
	}
	seg, _, _ := strings.Cut(rest, "/")
	t, err := strconv.Atoi(seg)
	if err != nil || t < 1 || t > 4 {
		return 0, false
	}
	return t, true
}

// asInt converts a device UI value to an int.
func asInt(v any) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case int32:
		return int(t), true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(t))
		return n, err == nil
	}
	return 0, false
}

// asKnobValue converts a device UI value to the knob value, handling the
// pan string forms ("C", "L13", "R13") in addition to signed ints.
func asKnobValue(v any) int {
	if n, ok := asInt(v); ok {
		return n
	}
	if s, ok := v.(string); ok {
		s := strings.ToUpper(strings.TrimSpace(s))
		switch {
		case s == "C" || s == "CENTER" || s == "":
			return 0
		case len(s) >= 2 && (s[0] == 'L' || s[0] == 'R'):
			n, err := strconv.Atoi(s[1:])
			if err != nil {
				return 0
			}
			if s[0] == 'L' {
				return -n
			}
			return n
		}
	}
	return 0
}
