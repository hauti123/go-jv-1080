package jv1080

import (
	"fmt"
	"strconv"
	"strings"

	"go-jv-1080/internal/midi/message"
	"go-jv-1080/internal/util"
)

// Parameter keys are path-like strings addressed to the device layer, e.g.
// "patch/temp/common/level". The first segment selects the memory area
// ("patch"), the second the patch ("temp" = patch-mode temporary patch,
// later "user/1".."user/128" etc.), then the block ("common", "tone/<n>"
// or a tone sub-block "tone/<n>/<sub>"), then the parameter name.
//
// The mapping from key to SysEx address and the conversion between UI
// values (int, string, bool) and wire bytes live solely in this package;
// callers never deal with addresses or offsets.

// ParamKind describes how a parameter's UI value maps onto wire bytes.
type ParamKind int

const (
	// KindRaw is a plain 0-127 byte parameter (e.g. level).
	KindRaw ParamKind = iota
	// KindBipolar is a signed value stored with an offset (e.g. octave
	// -3..+3 stored as 0-6, or envelope levels -63..+63 stored as 0-126).
	KindBipolar
	// KindPan is a pan value: UI value is "C", "L1".."L64", "R1".."R63"
	// (or signed -64..+63); wire value is 0-127 with 64 = center.
	KindPan
	// KindNibblized is a value wider than 7 bits transmitted as two
	// nibble bytes, high nibble first (e.g. wave number, tempo).
	KindNibblized
	// KindName is the 12 character patch name (ASCII bytes 32-127).
	KindName
)

// Param describes one addressable parameter of the device.
type Param struct {
	Key    string    // path-like key, e.g. "patch/temp/common/level"
	Kind   ParamKind // how UI values convert to wire bytes
	Min    int       // UI value minimum (for KindBipolar/KindPan: signed)
	Max    int       // UI value maximum
	Offset int       // byte offset within the block (device-internal)
	Block  string    // "common" or "tone/N" (device-internal)
	area   message.Address
}

// paramSpec is the static definition of a parameter: its key template and
// wire location within the temporary patch. block is "common" or a tone
// sub-block path "tone/N" / "tone/N/<sub>"; tone 0 means common.
type paramSpec struct {
	name   string
	block  string // "common" or "tone/N" or "tone/N/<sub>"
	tone   int    // 0 for common, else 1-4
	offset int
	kind   ParamKind
	min    int
	max    int
}

// paramSize is the number of wire bytes a parameter occupies.
func (s paramSpec) paramSize() int {
	switch s.kind {
	case KindName:
		return 12
	case KindNibblized:
		return 2
	default:
		return 1
	}
}

// toneSpec builds the parameter spec of one tone (1-4) from a template
// with tone 0; the tone number is filled in and the block paths get the
// tone prefix.
func toneSpec(t int, tmpl paramSpec) paramSpec {
	s := tmpl
	s.tone = t
	s.block = strings.Replace(s.block, "tone/0", fmt.Sprintf("tone/%d", t), 1)
	return s
}

// tempPatchParams are the parameters of the patch-mode temporary patch
// ("patch/temp/...") exposed through the key system. The offsets and
// ranges follow doc/jv-1080-midi.md tables *1-3-1 (patch common) and
// *1-3-2 (patch tone); the CLI names follow
// doc/patch-edit-mode-structure.md.
var tempPatchParams = buildTempPatchParams()

func buildTempPatchParams() []paramSpec {
	var specs []paramSpec
	add := func(s paramSpec) { specs = append(specs, s) }

	// --- Patch common (*1-3-1) ---
	add(paramSpec{name: "name", block: "common", offset: 0x00, kind: KindName})
	add(paramSpec{name: "efxtype", block: "common", offset: 0x0C, kind: KindRaw, min: 0, max: 39})
	for i := 1; i <= 12; i++ {
		add(paramSpec{name: fmt.Sprintf("efxprm%d", i), block: "common", offset: 0x0C + i, kind: KindRaw, min: 0, max: 127})
	}
	add(paramSpec{name: "efxout", block: "common", offset: 0x19, kind: KindRaw, min: 0, max: 2})
	add(paramSpec{name: "efxlevel", block: "common", offset: 0x1A, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "efxchorus", block: "common", offset: 0x1B, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "efxreverb", block: "common", offset: 0x1C, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "efxctrl1src", block: "common", offset: 0x1D, kind: KindRaw, min: 0, max: 10})
	add(paramSpec{name: "efxctrl1depth", block: "common", offset: 0x1E, kind: KindBipolar, min: -63, max: 63})
	add(paramSpec{name: "efxctrl2src", block: "common", offset: 0x1F, kind: KindRaw, min: 0, max: 10})
	add(paramSpec{name: "efxctrl2depth", block: "common", offset: 0x20, kind: KindBipolar, min: -63, max: 63})
	add(paramSpec{name: "choruslevel", block: "common", offset: 0x21, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "chorusrate", block: "common", offset: 0x22, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "chorusdepth", block: "common", offset: 0x23, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "choruspredelay", block: "common", offset: 0x24, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "chorusfeedback", block: "common", offset: 0x25, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "chorusout", block: "common", offset: 0x26, kind: KindRaw, min: 0, max: 2})
	add(paramSpec{name: "reverbtype", block: "common", offset: 0x27, kind: KindRaw, min: 0, max: 7})
	add(paramSpec{name: "reverblevel", block: "common", offset: 0x28, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "reverbtime", block: "common", offset: 0x29, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "hfdamp", block: "common", offset: 0x2A, kind: KindRaw, min: 0, max: 17})
	add(paramSpec{name: "reverbfeedback", block: "common", offset: 0x2B, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "tempo", block: "common", offset: 0x2C, kind: KindNibblized, min: 20, max: 250})
	add(paramSpec{name: "level", block: "common", offset: 0x2E, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "pan", block: "common", offset: 0x2F, kind: KindPan, min: -64, max: 63})
	add(paramSpec{name: "analogfeel", block: "common", offset: 0x30, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "bendup", block: "common", offset: 0x31, kind: KindRaw, min: 0, max: 12})
	add(paramSpec{name: "benddown", block: "common", offset: 0x32, kind: KindRaw, min: 0, max: 48})
	add(paramSpec{name: "assign", block: "common", offset: 0x33, kind: KindRaw, min: 0, max: 1})
	add(paramSpec{name: "legato", block: "common", offset: 0x34, kind: KindRaw, min: 0, max: 1})
	add(paramSpec{name: "portamento", block: "common", offset: 0x35, kind: KindRaw, min: 0, max: 1})
	add(paramSpec{name: "portmode", block: "common", offset: 0x36, kind: KindRaw, min: 0, max: 1})
	add(paramSpec{name: "porttype", block: "common", offset: 0x37, kind: KindRaw, min: 0, max: 1})
	add(paramSpec{name: "portstart", block: "common", offset: 0x38, kind: KindRaw, min: 0, max: 1})
	add(paramSpec{name: "porttime", block: "common", offset: 0x39, kind: KindRaw, min: 0, max: 127})
	add(paramSpec{name: "ctrl2src", block: "common", offset: 0x3A, kind: KindRaw, min: 0, max: 15})
	add(paramSpec{name: "ctrl3src", block: "common", offset: 0x3B, kind: KindRaw, min: 0, max: 15})
	add(paramSpec{name: "efxhold", block: "common", offset: 0x3C, kind: KindRaw, min: 0, max: 2})
	add(paramSpec{name: "ctrl1hold", block: "common", offset: 0x3D, kind: KindRaw, min: 0, max: 2})
	add(paramSpec{name: "ctrl2hold", block: "common", offset: 0x3E, kind: KindRaw, min: 0, max: 2})
	add(paramSpec{name: "ctrl3hold", block: "common", offset: 0x3F, kind: KindRaw, min: 0, max: 2})
	add(paramSpec{name: "velrange", block: "common", offset: 0x40, kind: KindRaw, min: 0, max: 1})
	add(paramSpec{name: "octave", block: "common", offset: 0x41, kind: KindBipolar, min: -3, max: 3})
	add(paramSpec{name: "stretch", block: "common", offset: 0x42, kind: KindRaw, min: 0, max: 3})
	add(paramSpec{name: "priority", block: "common", offset: 0x43, kind: KindRaw, min: 0, max: 1})
	add(paramSpec{name: "struct12", block: "common", offset: 0x44, kind: KindRaw, min: 0, max: 9})
	add(paramSpec{name: "booster12", block: "common", offset: 0x45, kind: KindRaw, min: 0, max: 3})
	add(paramSpec{name: "struct34", block: "common", offset: 0x46, kind: KindRaw, min: 0, max: 9})
	add(paramSpec{name: "booster34", block: "common", offset: 0x47, kind: KindRaw, min: 0, max: 3})

	// --- Patch tone (*1-3-2), template with tone 0, instantiated 1-4 ---
	tone := []paramSpec{
		{name: "switch", block: "tone/0", offset: 0x00, kind: KindRaw, min: 0, max: 1},
		{name: "group", block: "tone/0", offset: 0x01, kind: KindRaw, min: 0, max: 2},
		{name: "groupid", block: "tone/0", offset: 0x02, kind: KindRaw, min: 0, max: 127},
		{name: "number", block: "tone/0", offset: 0x03, kind: KindNibblized, min: 1, max: 255},
		{name: "gain", block: "tone/0", offset: 0x05, kind: KindRaw, min: 0, max: 3},
		{name: "fxmswitch", block: "tone/0", offset: 0x06, kind: KindRaw, min: 0, max: 1},
		{name: "fxmcolor", block: "tone/0", offset: 0x07, kind: KindRaw, min: 0, max: 3},
		{name: "fxmdepth", block: "tone/0", offset: 0x08, kind: KindRaw, min: 0, max: 15},
		{name: "delaymode", block: "tone/0", offset: 0x09, kind: KindRaw, min: 0, max: 6},
		{name: "delaytime", block: "tone/0", offset: 0x0A, kind: KindRaw, min: 0, max: 127},
		{name: "velfade", block: "tone/0", offset: 0x0B, kind: KindRaw, min: 0, max: 127},
		{name: "vellower", block: "tone/0", offset: 0x0C, kind: KindRaw, min: 1, max: 127},
		{name: "velupper", block: "tone/0", offset: 0x0D, kind: KindRaw, min: 1, max: 127},
		{name: "keylower", block: "tone/0", offset: 0x0E, kind: KindRaw, min: 0, max: 127},
		{name: "keyupper", block: "tone/0", offset: 0x0F, kind: KindRaw, min: 0, max: 127},
		{name: "redamper", block: "tone/0", offset: 0x10, kind: KindRaw, min: 0, max: 1},
		{name: "rxvolume", block: "tone/0", offset: 0x11, kind: KindRaw, min: 0, max: 1},
		{name: "rxhold", block: "tone/0", offset: 0x12, kind: KindRaw, min: 0, max: 1},
		{name: "rxbend", block: "tone/0", offset: 0x13, kind: KindRaw, min: 0, max: 1},
		{name: "rxpan", block: "tone/0", offset: 0x14, kind: KindRaw, min: 0, max: 2},
	}
	// Controller 1-3 destinations (dst1-4) and depths (dep1-4).
	for i := 1; i <= 3; i++ {
		for k := 1; k <= 4; k++ {
			base := 0x15 + (i-1)*8 + (k-1)*2
			tone = append(tone,
				paramSpec{name: fmt.Sprintf("dst%d", k), block: fmt.Sprintf("tone/0/ctrl%d", i), offset: base, kind: KindRaw, min: 0, max: 18},
				paramSpec{name: fmt.Sprintf("dep%d", k), block: fmt.Sprintf("tone/0/ctrl%d", i), offset: base + 1, kind: KindBipolar, min: -63, max: 63},
			)
		}
	}
	// LFO1/LFO2.
	lfoNames := []string{"wave", "keytrig", "rate", "offset", "delay", "fade", "fadetime", "sync"}
	for l := 1; l <= 2; l++ {
		for j, n := range lfoNames {
			tone = append(tone, paramSpec{name: n, block: fmt.Sprintf("tone/0/lfo%d", l), offset: 0x2D + (l-1)*8 + j, kind: KindRaw, min: 0, max: 127})
		}
	}
	// LFO depth 1:2 (pitch, filter, amp, pan); LFO2 sits one byte up.
	lfoDepth := []struct {
		name   string
		offset int
	}{
		{"pitch", 0x4E}, {"filter", 0x63}, {"amp", 0x75}, {"pan", 0x7B},
	}
	for l := 1; l <= 2; l++ {
		for _, d := range lfoDepth {
			tone = append(tone, paramSpec{name: d.name, block: fmt.Sprintf("tone/0/lfo%d", l), offset: d.offset + l - 1, kind: KindBipolar, min: -63, max: 63})
		}
	}
	// PITCH.
	tone = append(tone,
		paramSpec{name: "coarse", block: "tone/0/pitch", offset: 0x3D, kind: KindBipolar, min: -48, max: 48},
		paramSpec{name: "fine", block: "tone/0/pitch", offset: 0x3E, kind: KindBipolar, min: -50, max: 50},
		paramSpec{name: "random", block: "tone/0/pitch", offset: 0x3F, kind: KindRaw, min: 0, max: 30},
		paramSpec{name: "keyfollow", block: "tone/0/pitch", offset: 0x40, kind: KindRaw, min: 0, max: 15},
		paramSpec{name: "depth", block: "tone/0/pitch", offset: 0x41, kind: KindBipolar, min: -12, max: 12},
		paramSpec{name: "velsense", block: "tone/0/pitch", offset: 0x42, kind: KindRaw, min: 0, max: 125},
		paramSpec{name: "vt1", block: "tone/0/pitch", offset: 0x43, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "vt4", block: "tone/0/pitch", offset: 0x44, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "envkeyfollow", block: "tone/0/pitch", offset: 0x45, kind: KindRaw, min: 0, max: 14},
	)
	for i := 1; i <= 4; i++ {
		tone = append(tone, paramSpec{name: fmt.Sprintf("t%d", i), block: "tone/0/pitch", offset: 0x45 + i, kind: KindRaw, min: 0, max: 127})
		tone = append(tone, paramSpec{name: fmt.Sprintf("l%d", i), block: "tone/0/pitch", offset: 0x49 + i, kind: KindBipolar, min: -63, max: 63})
	}
	// TVF.
	tone = append(tone,
		paramSpec{name: "filtertype", block: "tone/0/tvf", offset: 0x50, kind: KindRaw, min: 0, max: 4},
		paramSpec{name: "cutoff", block: "tone/0/tvf", offset: 0x51, kind: KindRaw, min: 0, max: 127},
		paramSpec{name: "cutoffkeyfollow", block: "tone/0/tvf", offset: 0x52, kind: KindRaw, min: 0, max: 15},
		paramSpec{name: "resonance", block: "tone/0/tvf", offset: 0x53, kind: KindRaw, min: 0, max: 127},
		paramSpec{name: "velres", block: "tone/0/tvf", offset: 0x54, kind: KindRaw, min: 0, max: 125},
		paramSpec{name: "envdepth", block: "tone/0/tvf", offset: 0x55, kind: KindBipolar, min: -63, max: 63},
		paramSpec{name: "velcurve", block: "tone/0/tvf", offset: 0x56, kind: KindRaw, min: 0, max: 6},
		paramSpec{name: "velsens", block: "tone/0/tvf", offset: 0x57, kind: KindRaw, min: 0, max: 125},
		paramSpec{name: "vt1", block: "tone/0/tvf", offset: 0x58, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "vt4", block: "tone/0/tvf", offset: 0x59, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "envkeyfollow", block: "tone/0/tvf", offset: 0x5A, kind: KindRaw, min: 0, max: 14},
	)
	for i := 1; i <= 4; i++ {
		tone = append(tone, paramSpec{name: fmt.Sprintf("t%d", i), block: "tone/0/tvf", offset: 0x5A + i, kind: KindRaw, min: 0, max: 127})
		tone = append(tone, paramSpec{name: fmt.Sprintf("l%d", i), block: "tone/0/tvf", offset: 0x5E + i, kind: KindRaw, min: 0, max: 127})
	}
	// TVA.
	tone = append(tone,
		paramSpec{name: "level", block: "tone/0/tva", offset: 0x65, kind: KindRaw, min: 0, max: 127},
		paramSpec{name: "biasdir", block: "tone/0/tva", offset: 0x66, kind: KindRaw, min: 0, max: 3},
		paramSpec{name: "biaspoint", block: "tone/0/tva", offset: 0x67, kind: KindRaw, min: 0, max: 127},
		paramSpec{name: "bias", block: "tone/0/tva", offset: 0x68, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "envcurve", block: "tone/0/tva", offset: 0x69, kind: KindRaw, min: 0, max: 6},
		paramSpec{name: "envvel", block: "tone/0/tva", offset: 0x6A, kind: KindRaw, min: 0, max: 125},
		paramSpec{name: "vt1", block: "tone/0/tva", offset: 0x6B, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "vt4", block: "tone/0/tva", offset: 0x6C, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "envkeyfollow", block: "tone/0/tva", offset: 0x6D, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "pan", block: "tone/0/tva", offset: 0x77, kind: KindPan, min: -64, max: 63},
		paramSpec{name: "pankeyfollow", block: "tone/0/tva", offset: 0x78, kind: KindRaw, min: 0, max: 14},
		paramSpec{name: "randompan", block: "tone/0/tva", offset: 0x79, kind: KindRaw, min: 0, max: 63},
		paramSpec{name: "altpan", block: "tone/0/tva", offset: 0x7A, kind: KindRaw, min: 1, max: 127},
	)
	for i := 1; i <= 4; i++ {
		tone = append(tone, paramSpec{name: fmt.Sprintf("t%d", i), block: "tone/0/tva", offset: 0x6D + i, kind: KindRaw, min: 0, max: 127})
	}
	for i := 1; i <= 3; i++ {
		tone = append(tone, paramSpec{name: fmt.Sprintf("l%d", i), block: "tone/0/tva", offset: 0x71 + i, kind: KindRaw, min: 0, max: 127})
	}
	// Per-Tone output routing (EFFECTS category).
	tone = append(tone,
		paramSpec{name: "output", block: "tone/0/efx", offset: 0x7D, kind: KindRaw, min: 0, max: 3},
		paramSpec{name: "outputlevel", block: "tone/0/efx", offset: 0x7E, kind: KindRaw, min: 0, max: 127},
		paramSpec{name: "chorus", block: "tone/0/efx", offset: 0x7F, kind: KindRaw, min: 0, max: 127},
		paramSpec{name: "reverb", block: "tone/0/efx", offset: 0x80, kind: KindRaw, min: 0, max: 127},
	)

	for _, tmpl := range tone {
		for t := 1; t <= 4; t++ {
			add(toneSpec(t, tmpl))
		}
	}
	return specs
}

// resolveParam maps a path-like key to its parameter definition and the
// SysEx address of its first byte. Only "patch/temp/..." keys are
// supported for now.
func resolveParam(key string) (Param, error) {
	segs := strings.Split(strings.Trim(key, "/"), "/")
	if len(segs) < 4 || segs[0] != "patch" || segs[1] != "temp" {
		return Param{}, fmt.Errorf("jv1080: unsupported parameter key %q (want patch/temp/<block>/<name> or patch/temp/tone/<n>/<name>)", key)
	}
	block := strings.Join(segs[2:len(segs)-1], "/")
	name := segs[len(segs)-1]
	for _, spec := range tempPatchParams {
		if spec.block != block || spec.name != name {
			continue
		}
		area := TemporaryPatchArea()
		var addr message.Address
		if spec.tone == 0 {
			addr = area.Add(spec.offset)
		} else {
			// tone blocks are 0x80 bytes apart; offset 0x80 and above
			// (e.g. the reverb send at 01 00) wraps into the next 0x80
			// step of the address, which Address.Add handles
			addr = PatchToneArea(area, spec.tone).Add(spec.offset)
		}
		return Param{
			Key:    paramKeyOf(spec),
			Kind:   spec.kind,
			Min:    spec.min,
			Max:    spec.max,
			Offset: spec.offset,
			Block:  spec.block,
			area:   addr,
		}, nil
	}
	return Param{}, fmt.Errorf("jv1080: unknown parameter key %q", key)
}

// EncodeValue converts a UI value (int, string, bool) into the raw wire
// bytes for the parameter. For KindBipolar the signed value is shifted by
// -Min; for KindPan strings like "L13"/"R13"/"C" (or signed ints) are
// converted to the 0-127 wire value with 64 = center; for KindNibblized
// the value is split into two nibble bytes, high nibble first.
func (p Param) EncodeValue(v any) ([]byte, error) {
	switch p.Kind {
	case KindBipolar:
		n, err := toInt(v)
		if err != nil {
			return nil, fmt.Errorf("jv1080: %s: %w", p.Key, err)
		}
		if n < p.Min || n > p.Max {
			return nil, fmt.Errorf("jv1080: %s: value %d out of range %d..%d", p.Key, n, p.Min, p.Max)
		}
		return []byte{byte(n - p.Min)}, nil
	case KindPan:
		raw, err := panToRaw(v)
		if err != nil {
			return nil, fmt.Errorf("jv1080: %s: %w", p.Key, err)
		}
		return []byte{byte(raw)}, nil
	case KindNibblized:
		n, err := toInt(v)
		if err != nil {
			return nil, fmt.Errorf("jv1080: %s: %w", p.Key, err)
		}
		if n < p.Min || n > p.Max {
			return nil, fmt.Errorf("jv1080: %s: value %d out of range %d..%d", p.Key, n, p.Min, p.Max)
		}
		w := n - p.Min
		return []byte{byte(w >> 4), byte(w & 0x0F)}, nil
	case KindName:
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("jv1080: %s: want string value", p.Key)
		}
		b := make([]byte, 12)
		for i := range b {
			if i < len(s) {
				b[i] = byte(util.Clamp(int(s[i]), 32, 127))
			} else {
				b[i] = ' '
			}
		}
		return b, nil
	default:
		n, err := toInt(v)
		if err != nil {
			return nil, fmt.Errorf("jv1080: %s: %w", p.Key, err)
		}
		if n < p.Min || n > p.Max {
			return nil, fmt.Errorf("jv1080: %s: value %d out of range %d..%d", p.Key, n, p.Min, p.Max)
		}
		return []byte{byte(n)}, nil
	}
}

// DecodeValue converts raw wire bytes into the UI value for the parameter
// (int for KindRaw/KindBipolar/KindNibblized, string for KindName, signed
// UI units for KindPan).
func (p Param) DecodeValue(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("jv1080: %s: empty data", p.Key)
	}
	switch p.Kind {
	case KindBipolar:
		return int(data[0]) + p.Min, nil
	case KindPan:
		return rawToPan(int(data[0])), nil // signed UI units -64..+63, 0 = C
	case KindNibblized:
		if len(data) < 2 {
			return nil, fmt.Errorf("jv1080: %s: want 2 nibble bytes, got %d", p.Key, len(data))
		}
		return (int(data[0])<<4 | int(data[1])) + p.Min, nil
	case KindName:
		return decodeName(data), nil
	default:
		return int(data[0]), nil
	}
}

// panToRaw converts a pan UI value to the raw 0-127 wire byte. UI and
// wire both start at 0 = center; the JV-1080 scale is L64..L1, C,
// R1..R63, so a plain int is the signed UI value -64..+63 (0 = C,
// -64 = L64) and maps to 64+ui. Strings "C"/"center" (0), "L1".."L64"
// (negative) and "R1".."R63" (positive) are also accepted.
func panToRaw(v any) (int, error) {
	ui, err := panToUI(v)
	if err != nil {
		return 0, err
	}
	if ui < -64 || ui > 63 {
		return 0, fmt.Errorf("pan value %d out of range -64..+63 (L64..R63)", ui)
	}
	return 64 + ui, nil
}

// rawToPan converts a raw 0-127 wire byte to the signed UI pan value
// -64..+63 (0 = C, -64 = L64).
func rawToPan(raw int) int {
	return util.Clamp(raw, 0, 127) - 64
}

// panToUI converts a pan UI value (signed int -64..+63, or "C"/"Lxx"/
// "Rxx" string) to the signed UI value.
func panToUI(v any) (int, error) {
	switch t := v.(type) {
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		switch s {
		case "c", "center":
			return 0, nil
		}
		if len(s) >= 2 && (s[0] == 'l' || s[0] == 'r') {
			n, err := strconv.Atoi(s[1:])
			max := 63
			if s[0] == 'l' {
				max = 64 // left side reaches L64
			}
			if err != nil || n < 1 || n > max {
				return 0, fmt.Errorf("invalid pan %q (want C, L1-L64 or R1-R63)", t)
			}
			if s[0] == 'l' {
				return -n, nil
			}
			return n, nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid pan %q", t)
		}
		return n, nil
	default:
		n, err := toInt(v)
		if err != nil {
			return 0, err
		}
		return n, nil
	}
}

func toInt(v any) (int, error) {
	switch t := v.(type) {
	case int:
		return t, nil
	case int32:
		return int(t), nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		return strconv.Atoi(strings.TrimSpace(t))
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("cannot convert %T to int", v)
}
