package gui

import (
	"image"
	"strings"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"go-jv-1080/internal/gui/widget"
	"go-jv-1080/internal/jv1080"
)

// sent captures the (key, value) pairs passed to send.
type sent struct {
	key   string
	value any
}

func newTestEditor() (*patchEditor, *[]sent) {
	var out []sent
	p := newPatchEditor(func(key string, value any) { out = append(out, sent{key, value}) })
	return p, &out
}

// covered reports whether a parameter key is handled by the editor:
// a knob, the name field, an envelope point or the wave selection.
func covered(p *patchEditor, key string) bool {
	if _, ok := p.knobs[key]; ok {
		return true
	}
	if key == nameKey {
		return true
	}
	block := key[:strings.LastIndex(key, "/")]
	if e, ok := p.envs[block]; ok {
		name := key[strings.LastIndex(key, "/")+1:]
		if len(name) == 2 && (name[0] == 't' || name[0] == 'l') && name[1] >= '1' && name[1] <= '4' {
			_ = e
			return true
		}
	}
	if t, ok := toneNumberOf(key); ok && p.waves[t] != nil {
		base := "patch/temp/tone/" + string(rune('0'+t))
		if key == base+"/group" || key == base+"/groupid" || key == base+"/number" {
			return true
		}
	}
	return false
}

func TestPatchEditorCoversAllParams(t *testing.T) {
	p, _ := newTestEditor()
	params := jv1080.ListParams()
	if len(params) == 0 {
		t.Fatal("ListParams is empty")
	}
	for _, param := range params {
		if !covered(p, param.Key) {
			t.Errorf("parameter %q (block %q) is not editable on the screen", param.Key, param.Block)
		}
	}
	// all five tabs have sections
	for i, secs := range p.content {
		if len(secs) == 0 {
			t.Errorf("tab %d has no sections", i)
		}
	}
	// the common tab starts with the NAME section, each tone tab with WAVE
	if p.content[0][0].title != "NAME" {
		t.Errorf("common tab first section = %q, want NAME", p.content[0][0].title)
	}
	for tone := 1; tone <= 4; tone++ {
		if p.content[tone][0].title != "WAVE" {
			t.Errorf("tone %d tab first section = %q, want WAVE", tone, p.content[tone][0].title)
		}
	}
	// 4 wave switch knobs with 6 positions each
	for tone := 1; tone <= 4; tone++ {
		w := p.waves[tone]
		if w == nil || w.knob.Positions() != 6 || len(w.dd.Options) != 255 {
			t.Errorf("tone %d wave selection: knob positions = %d, dropdown options = %d", tone, w.knob.Positions(), len(w.dd.Options))
		}
	}
}

func TestPatchEditorKnobSends(t *testing.T) {
	p, out := newTestEditor()
	key := "patch/temp/common/level"
	k, ok := p.knobs[key]
	if !ok {
		t.Fatalf("no knob for %q", key)
	}
	k.OnChange(77)
	if len(*out) != 1 || (*out)[0].key != key || (*out)[0].value != 77 {
		t.Errorf("send = %v, want [%s 77]", *out, key)
	}
}

func TestPatchEditorApplyParam(t *testing.T) {
	p, _ := newTestEditor()
	p.applyParam("patch/temp/common/level", 100)
	if got := p.knobs["patch/temp/common/level"].Value(); got != 100 {
		t.Errorf("level knob = %d, want 100", got)
	}
	// pan arrives as signed int from the device layer
	p.applyParam("patch/temp/common/pan", -13)
	if got := p.knobs["patch/temp/common/pan"].Value(); got != -13 {
		t.Errorf("pan knob = %d, want -13", got)
	}
	// pan string form
	p.applyParam("patch/temp/common/pan", "R5")
	if got := p.knobs["patch/temp/common/pan"].Value(); got != 5 {
		t.Errorf("pan knob after \"R5\" = %d, want 5", got)
	}
	// name
	p.applyParam(nameKey, "MY PATCH")
	if got := p.name.Text(); got != "MY PATCH" {
		t.Errorf("name = %q, want %q", got, "MY PATCH")
	}
	// envelope point of tone 2 pitch
	p.applyParam("patch/temp/tone/2/pitch/t2", 50)
	env := p.envs["patch/temp/tone/2/pitch"]
	if tv, _ := env.Point(1); tv != 50 {
		t.Errorf("pitch env t2 = %d, want 50", tv)
	}
	p.applyParam("patch/temp/tone/2/pitch/l3", -20)
	if _, lv := env.Point(2); lv != -20 {
		t.Errorf("pitch env l3 = %d, want -20", lv)
	}
	// wave number of tone 1 (device value 1-based)
	p.applyParam("patch/temp/tone/1/number", 42)
	if got := p.waves[1].dd.Value(); got != 41 {
		t.Errorf("wave dropdown = %d, want 41 (number 42)", got)
	}
	// wave group + group id drive the 6-position switch: EXP-C = (2,3) = index 4
	p.applyParam("patch/temp/tone/3/group", 2)
	p.applyParam("patch/temp/tone/3/groupid", 3)
	if got := p.waves[3].knob.Value(); got != 4 {
		t.Errorf("wave switch = %d, want 4 (EXP-C)", got)
	}
	if n := len(p.waves[3].dd.Options); n != 255 {
		t.Errorf("wave dropdown options after EXP-C = %d, want 255", n)
	}
}

func TestPatchEditorWaveSwitchSends(t *testing.T) {
	p, out := newTestEditor()
	p.waves[1].knob.OnChange(3) // EXP-B = group 2, id 2
	want := []sent{
		{"patch/temp/tone/1/group", 2},
		{"patch/temp/tone/1/groupid", 2},
	}
	if len(*out) != len(want) {
		t.Fatalf("send = %v, want %v", *out, want)
	}
	for i, w := range want {
		if (*out)[i].key != w.key || (*out)[i].value != w.value {
			t.Errorf("send[%d] = %v, want %v", i, (*out)[i], w)
		}
	}
	if p.waves[1].group != 2 || p.waves[1].groupID != 2 {
		t.Errorf("cached group/id = %d/%d, want 2/2", p.waves[1].group, p.waves[1].groupID)
	}
	// INT-B options after switching to index 1
	p.waves[2].knob.OnChange(1)
	if n := len(p.waves[2].dd.Options); n != 193 {
		t.Errorf("INT-B options = %d, want 193", n)
	}
}

func TestPatchEditorWaveDropdownSendsNumber(t *testing.T) {
	p, out := newTestEditor()
	p.waves[4].dd.OnChange(99) // index 99 = wave number 100
	if len(*out) != 1 || (*out)[0].key != "patch/temp/tone/4/number" || (*out)[0].value != 100 {
		t.Errorf("send = %v, want [patch/temp/tone/4/number 100]", *out)
	}
}

func TestPatchEditorLayoutSmoke(t *testing.T) {
	p, _ := newTestEditor()
	th := widget.DefaultTheme()
	for tab := 0; tab < 5; tab++ {
		p.tabs.Select(tab)
		// a very tall viewport so the list lays out every section
		// (envelopes, wave dropdown, name field included)
		gtx := layout.Context{
			Ops:         new(op.Ops),
			Constraints: layout.Exact(image.Pt(780, 40000)),
			Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		}
		if d := p.Layout(gtx, &th); d.Size.X <= 0 || d.Size.Y <= 0 {
			t.Errorf("tab %d layout dims = %v", tab, d.Size)
		}
	}
}

func TestParamBytes(t *testing.T) {
	patch := jv1080.NewPatch(jv1080.TemporaryPatchArea())
	patch.SetPatchLevel(100)
	patch.SetWaveNumber(2, 42)
	param, err := jv1080.ResolveParam("patch/temp/common/level")
	if err != nil {
		t.Fatal(err)
	}
	if b := paramBytes(patch, param); len(b) != 1 || b[0] != 100 {
		t.Errorf("common level bytes = %v, want [100]", b)
	}
	param, err = jv1080.ResolveParam("patch/temp/tone/2/number")
	if err != nil {
		t.Fatal(err)
	}
	b := paramBytes(patch, param)
	if len(b) != 2 {
		t.Fatalf("wave number bytes = %v, want 2 nibble bytes", b)
	}
	v, err := param.DecodeValue(b)
	if err != nil || v != 42 {
		t.Errorf("decoded wave number = %v, %v; want 42", v, err)
	}
}
