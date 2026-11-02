package gui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"

	"go-jv-1080/internal/gui/widget"
	"go-jv-1080/internal/jv1080"
)

// PoC is the proof-of-concept screen: three patch common knobs (level,
// pan, octave) and the tone 1 pitch envelope, bound to patch/temp
// parameter keys. It is kept separate from the app shell (app.go) so it
// can be replaced by the patch edit screen later.
type PoC struct {
	level  *widget.Knob
	pan    *widget.Knob
	octave *widget.Knob
	env    *widget.Envelope

	send func(key string, value any)
}

// device parameter keys of the temporary patch (patch/temp/...), named
// like in jv1080-cli. The GUI never deals with addresses or offsets; the
// key -> SysEx mapping lives in internal/jv1080.
const (
	keyLevel  = "patch/temp/common/level"
	keyPan    = "patch/temp/common/pan"
	keyOctave = "patch/temp/common/octave"
)

// pitch envelope of tone 1: 4 times + 4 levels (PCH ENVELOPE)
var (
	envTimeKeys  = [4]string{"patch/temp/tone/1/pitch/t1", "patch/temp/tone/1/pitch/t2", "patch/temp/tone/1/pitch/t3", "patch/temp/tone/1/pitch/t4"}
	envLevelKeys = [4]string{"patch/temp/tone/1/pitch/l1", "patch/temp/tone/1/pitch/l2", "patch/temp/tone/1/pitch/l3", "patch/temp/tone/1/pitch/l4"}
)

func newPoC(send func(key string, value any)) *PoC {
	p := &PoC{send: send}
	p.level = widget.NewKnob(widget.KnobMinMax, "LEVEL", keyLevel, 0, 127)
	p.pan = widget.NewKnob(widget.KnobPan, "PAN", keyPan, -64, 63)
	p.octave = widget.NewKnob(widget.KnobBipolar, "OCTAVE", keyOctave, -3, 3)
	p.level.OnChange = func(v int) { p.send(keyLevel, v) }
	p.pan.OnChange = func(v int) { p.send(keyPan, v) }
	p.octave.OnChange = func(v int) { p.send(keyOctave, v) }
	p.env = widget.NewEnvelope("PITCH ENV T1", envTimeKeys, envLevelKeys)
	p.env.OnChange = func(timeKey string, timeVal int, levelKey string, levelVal int) {
		p.send(timeKey, timeVal)
		p.send(levelKey, levelVal)
	}
	return p
}

// Layout renders the PoC screen: the knob row with the envelope below it.
func (p *PoC) Layout(gtx layout.Context, th *widget.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return p.level.Layout(gtx, th) }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return p.pan.Layout(gtx, th) }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return p.octave.Layout(gtx, th) }),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.env.Layout(gtx, th)
			})
		}),
	)
}

// initDevice reads the current values of the bound parameters once.
func (p *PoC) initDevice(c *jv1080.Client) {
	for _, key := range []string{keyLevel, keyPan, keyOctave} {
		if v, err := c.GetParameter(key); err == nil {
			p.applyParam(key, v)
		}
	}
	// read the pitch envelope once (SetAll so the x-order clamp does
	// not interfere with loading)
	var times, levels [4]int
	for i := 0; i < 4; i++ {
		if v, err := c.GetParameter(envTimeKeys[i]); err == nil {
			if n, ok := v.(int); ok {
				times[i] = n
			}
		}
		if v, err := c.GetParameter(envLevelKeys[i]); err == nil {
			if n, ok := v.(int); ok {
				levels[i] = n
			}
		}
	}
	p.env.SetAll(times, levels)
}

// applyParam sets the widget bound to key without re-sending it.
func (p *PoC) applyParam(key string, value any) {
	switch key {
	case keyLevel:
		if n, ok := value.(int); ok {
			p.level.Set(n)
		}
	case keyPan:
		p.pan.Set(panUIToKnob(value))
	case keyOctave:
		if n, ok := value.(int); ok {
			p.octave.Set(n)
		}
	default:
		for i := 0; i < 4; i++ {
			if key == envTimeKeys[i] || key == envLevelKeys[i] {
				tv, lv := p.env.Point(i)
				if n, ok := value.(int); ok {
					if key == envTimeKeys[i] {
						tv = n
					} else {
						lv = n
					}
				}
				p.env.SetPoint(i, tv, lv)
			}
		}
	}
}

// panUIToKnob converts a pan UI value from the device layer ("C", "L13",
// "R13" or signed -64..+63) into the knob's -64..+63 range.
func panUIToKnob(value any) int {
	switch v := value.(type) {
	case int:
		return v // device layer already returns signed UI units
	case string:
		s := v
		if len(s) >= 2 && (s[0] == 'L' || s[0] == 'R') {
			var n int
			fmt.Sscanf(s[1:], "%d", &n)
			if s[0] == 'L' {
				return -n
			}
			return n
		}
		return 0 // "C"
	}
	return 0
}
