package main

import (
	"fmt"
	"strconv"
	"strings"

	"go-jv-1080/internal/jv1080"
)

// This file maps the CLI's category syntax ("set/get patch common <param>
// ...", "... efx tone 2 ...", ...) onto the device layer's path-like
// parameter keys ("patch/temp/..."). It contains no offsets and no value
// conversion: all addressing and value conversion happens in
// internal/jv1080, and every key built here is validated through
// jv1080.ResolveParam.

// patchCommand is a parsed "set/get patch <category> ..." command: the
// device-layer parameter key plus the index-free value arguments.
type patchCommand struct {
	key     string
	valArgs []string
}

// normalize lowercases a parameter name and strips '-', '_' and spaces.
func normalize(s string) string {
	s = strings.ToLower(s)
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' || c == '_' || c == ' ' {
			continue
		}
		out = append(out, c)
	}
	return string(out)
}

// paramAliases maps legacy CLI names to the canonical device-layer names.
var paramAliases = map[string]string{
	"efx":        "efxtype",
	"patchlevel": "level",
	"patchpan":   "pan",
	"wavenumber": "number",
	"wave":       "number",
	"wavegroup":  "group",
}

// aliasName resolves legacy aliases.
func aliasName(name string) string {
	if a, ok := paramAliases[name]; ok {
		return a
	}
	return name
}

// resolveKey validates a parameter key through the device layer.
func resolveKey(key string) error {
	_, err := jv1080.ResolveParam(key)
	return err
}

// parseToneArg parses a "tone <1-4>" selector prefix.
func parseToneArg(rest []string) (tone int, restOut []string, selected bool, err error) {
	if len(rest) >= 2 && rest[0] == "tone" {
		n, e := strconv.Atoi(rest[1])
		if e != nil || n < 1 || n > 4 {
			return 0, nil, false, fmt.Errorf("invalid tone %q (want 1-4)", rest[1])
		}
		return n, rest[2:], true, nil
	}
	return 1, rest, false, nil
}

// parseCommon parses: <param> <values...>
func parseCommon(args []string) ([]patchCommand, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("usage: set/get patch common <param> <values...>")
	}
	name := aliasName(normalize(args[0]))
	key := "patch/temp/common/" + name
	if err := resolveKey(key); err != nil {
		return nil, fmt.Errorf("unknown common parameter %q", args[0])
	}
	return []patchCommand{{key: key, valArgs: args[1:]}}, nil
}

// efxToneNames are the per-Tone OUTPUT parameter names of the EFFECTS
// category.
var efxToneNames = map[string]bool{
	"output": true, "outputlevel": true, "chorus": true, "reverb": true,
}

// efxCommonExcluded are common-block names that are not EFX parameters
// and must not be accepted without a tone selector in the efx category
// (e.g. level lives in the common block but belongs to the patch common,
// not the effects).
var efxCommonExcluded = map[string]bool{
	"level": true, "pan": true, "name": true, "tempo": true,
	"analogfeel": true, "bendup": true, "benddown": true, "assign": true,
	"legato": true, "portamento": true, "portmode": true, "porttype": true,
	"portstart": true, "porttime": true, "velrange": true, "octave": true,
	"stretch": true, "priority": true, "struct12": true, "booster12": true,
	"struct34": true, "booster34": true,
	"ctrl1hold": true, "ctrl2hold": true, "ctrl3hold": true,
	"ctrl2src": true, "ctrl3src": true,
}

const efxUsage = "usage: set/get patch efx [tone <1-4>] <param> <val>"

// parseEfx parses: [tone <1-4>] <param> <val>. Without a tone selector the
// parameter lives in the patch common block (EFX/chorus/reverb), with a
// tone selector in the per-Tone output routing (tone/<n>/efx/...).
func parseEfx(args []string) ([]patchCommand, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%s", efxUsage)
	}
	tone, rest, toneSelected, err := parseToneArg(args)
	if err != nil {
		return nil, err
	}
	if len(rest) < 1 {
		return nil, fmt.Errorf("%s", efxUsage)
	}
	name := aliasName(normalize(rest[0]))
	if !toneSelected {
		if efxCommonExcluded[name] {
			return nil, fmt.Errorf("unknown efx parameter %q (common EFX parameters take no tone; per-Tone: output, outputlevel, chorus, reverb)", rest[0])
		}
		key := "patch/temp/common/" + name
		if err := resolveKey(key); err != nil {
			return nil, fmt.Errorf("unknown efx parameter %q (common EFX parameters take no tone; per-Tone: output, outputlevel, chorus, reverb)", rest[0])
		}
		return []patchCommand{{key: key, valArgs: rest[1:]}}, nil
	}
	if !efxToneNames[name] {
		return nil, fmt.Errorf("unknown efx tone parameter %q (want output, outputlevel, chorus, reverb; common EFX parameters take no tone: set patch efx <param> <val>)", rest[0])
	}
	key := fmt.Sprintf("patch/temp/tone/%d/efx/%s", tone, name)
	if err := resolveKey(key); err != nil {
		return nil, fmt.Errorf("unknown efx tone parameter %q", rest[0])
	}
	return []patchCommand{{key: key, valArgs: rest[1:]}}, nil
}

// pitchToneNames are the PITCH parameters that need no env marker.
var pitchToneNames = map[string]bool{
	"coarse": true, "fine": true, "random": true, "keyfollow": true,
}

// pitchEnvNames are the PCH ENV DPT / PCH TIME ENV / PCH ENVELOPE
// parameters (env marker).
var pitchEnvNames = map[string]bool{
	"depth": true, "velsense": true, "vt1": true, "vt4": true, "keyfollow": true,
	"t1": true, "t2": true, "t3": true, "t4": true,
	"l1": true, "l2": true, "l3": true, "l4": true,
}

// tvfNames are the FILTER / TVF VELOCITY parameters (no env marker).
var tvfNames = map[string]bool{
	"filtertype": true, "cutoff": true, "cutoffkeyfollow": true, "resonance": true,
	"velres": true, "envdepth": true, "velcurve": true, "velsens": true,
}

// tvfEnvNames are the TVF TIME ENV / TVF ENVELOPE parameters (env marker).
var tvfEnvNames = map[string]bool{
	"vt1": true, "vt4": true, "keyfollow": true,
	"t1": true, "t2": true, "t3": true, "t4": true,
	"l1": true, "l2": true, "l3": true, "l4": true,
}

// tvaNames are the TVA / BIAS / PAN MODULATION parameters (no env marker).
var tvaNames = map[string]bool{
	"level": true, "pan": true, "envvel": true, "envcurve": true,
	"bias": true, "biaspoint": true, "biasdir": true,
	"pankeyfollow": true, "randompan": true, "altpan": true,
}

// tvaEnvNames are the TVA TIME ENV / TVA ENVELOPE parameters (env marker).
var tvaEnvNames = map[string]bool{
	"vt1": true, "vt4": true, "keyfollow": true,
	"t1": true, "t2": true, "t3": true, "t4": true,
	"l1": true, "l2": true, "l3": true,
}

// parseToneCategory parses the [tone <n>] [env] <param> <val> form shared
// by the pitch, tvf and tva categories and maps it onto
// patch/temp/tone/<n>/<cat>/<param>.
func parseToneCategory(cat string, args []string) ([]patchCommand, error) {
	usage := fmt.Sprintf("usage: set/get patch %s [tone <1-4>] [env] <param> <val>", cat)
	if len(args) < 1 {
		return nil, fmt.Errorf("%s", usage)
	}
	tone, rest, _, err := parseToneArg(args)
	if err != nil {
		return nil, err
	}
	env := false
	if len(rest) >= 1 && rest[0] == "env" {
		env = true
		rest = rest[1:]
	}
	if len(rest) < 1 {
		return nil, fmt.Errorf("%s", usage)
	}
	name := normalize(rest[0])
	var allowed map[string]bool
	switch cat {
	case "pitch":
		if env {
			allowed = pitchEnvNames
		} else {
			allowed = pitchToneNames
		}
	case "tvf":
		if env {
			allowed = tvfEnvNames
		} else {
			allowed = tvfNames
		}
	case "tva":
		if env {
			allowed = tvaEnvNames
		} else {
			allowed = tvaNames
		}
	}
	if env && name == "keyfollow" && allowed["keyfollow"] {
		// env keyfollow is the time keyfollow of the envelope, distinct
		// from the plain keyfollow of the same category
		name = "envkeyfollow"
	}
	if !allowed[name] && !(env && name == "envkeyfollow" && allowed["keyfollow"]) {
		if env {
			return nil, fmt.Errorf("unknown %s env parameter %q", cat, rest[0])
		}
		return nil, fmt.Errorf("unknown %s parameter %q (env parameters need the env marker: set/get patch %s [tone <1-4>] env <param> <val>)", cat, rest[0], cat)
	}
	key := fmt.Sprintf("patch/temp/tone/%d/%s/%s", tone, cat, name)
	if err := resolveKey(key); err != nil {
		return nil, fmt.Errorf("unknown %s parameter %q", cat, rest[0])
	}
	return []patchCommand{{key: key, valArgs: rest[1:]}}, nil
}

// ctrlToneNames are the per-Tone receive-switch names of the CONTROL
// category (RxSWITCH / DAMPER), living directly in the tone block.
var ctrlToneNames = map[string]bool{
	"redamper": true, "rxvolume": true, "rxhold": true,
	"rxbend": true, "rxpan": true,
}

// ctrlToneIndexedNames are the per-Tone CONTROL 1-3 destination and depth
// names that need a control index (1-3) and a tone number.
var ctrlToneIndexedNames = map[string]bool{
	"dst1": true, "dst2": true, "dst3": true, "dst4": true,
	"dep1": true, "dep2": true, "dep3": true, "dep4": true,
}

// parseCtrl parses the CONTROL category forms:
//
//	set/get patch ctrl <param> <val>                     (common, no index)
//	set/get patch ctrl <2-3> src <val>                   (common, index)
//	set/get patch ctrl <1-3> hold <val>                  (common, index)
//	set/get patch ctrl tone <1-4> <param> <val>          (per-Tone RxSWITCH)
//	set/get patch ctrl <1-3> tone <1-4> <dst|dep><1-4> <val>
func parseCtrl(args []string) ([]patchCommand, error) {
	const usage = "usage: set/get patch ctrl [<1-3>] [tone <1-4>] <param> <val>"
	if len(args) < 2 {
		return nil, fmt.Errorf("%s", usage)
	}
	// Control 1-3 per-Tone form: <1-3> tone <1-4> <dst|dep><1-4> [<val>].
	if len(args) >= 4 && args[1] == "tone" {
		i, err := strconv.Atoi(args[0])
		if err != nil || i < 1 || i > 3 {
			return nil, fmt.Errorf("invalid control index %q (want 1-3)", args[0])
		}
		tone, err := strconv.Atoi(args[2])
		if err != nil || tone < 1 || tone > 4 {
			return nil, fmt.Errorf("invalid tone %q (want 1-4)", args[2])
		}
		name := normalize(args[3])
		if !ctrlToneIndexedNames[name] {
			return nil, fmt.Errorf("unknown ctrl parameter %q (want dst1-dst4 or dep1-dep4)", args[3])
		}
		key := fmt.Sprintf("patch/temp/tone/%d/ctrl%d/%s", tone, i, name)
		if err := resolveKey(key); err != nil {
			return nil, fmt.Errorf("unknown ctrl parameter %q", args[3])
		}
		return []patchCommand{{key: key, valArgs: args[4:]}}, nil
	}
	// Common forms first: plain name or <index> hold|src.
	if args[0] != "tone" {
		if len(args) >= 2 && (args[1] == "hold" || args[1] == "src") {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, fmt.Errorf("invalid control index %q", args[0])
			}
			var key string
			if args[1] == "hold" {
				if n < 1 || n > 3 {
					return nil, fmt.Errorf("control index %d out of range for %q (want 1-3; usage: set patch ctrl <1-3> hold <val>)", n, args[1])
				}
				key = fmt.Sprintf("patch/temp/common/ctrl%dhold", n)
			} else {
				if n < 2 || n > 3 {
					return nil, fmt.Errorf("control index %d out of range for %q (want 2-3; usage: set patch ctrl <2-3> src <val>)", n, args[1])
				}
				key = fmt.Sprintf("patch/temp/common/ctrl%dsrc", n)
			}
			if err := resolveKey(key); err != nil {
				return nil, err
			}
			return []patchCommand{{key: key, valArgs: args[2:]}}, nil
		}
		name := normalize(args[0])
		if ctrlToneNames[name] || ctrlToneIndexedNames[name] {
			return nil, fmt.Errorf("parameter %q needs a tone selector: set patch ctrl [ctrl <1-3>] tone <1-4> %s <val>", args[0], args[0])
		}
		key := "patch/temp/common/" + aliasName(name)
		if err := resolveKey(key); err != nil {
			return nil, fmt.Errorf("unknown common parameter %q", args[0])
		}
		return []patchCommand{{key: key, valArgs: args[1:]}}, nil
	}
	// Per-Tone RxSWITCH form: tone <1-4> <param> [<val>].
	if len(args) < 3 {
		return nil, fmt.Errorf("usage: set/get patch ctrl tone <1-4> <param> [<val>]")
	}
	tone, err := strconv.Atoi(args[1])
	if err != nil || tone < 1 || tone > 4 {
		return nil, fmt.Errorf("invalid tone %q (want 1-4)", args[1])
	}
	name := normalize(args[2])
	if !ctrlToneNames[name] {
		if ctrlToneIndexedNames[name] {
			return nil, fmt.Errorf("parameter %q needs a control index: set patch ctrl <1-3> tone <1-4> %s <val>", args[2], args[2])
		}
		return nil, fmt.Errorf("unknown ctrl tone parameter %q (want redamper, rxvolume, rxhold, rxbend, rxpan or dst1-4, dep1-4)", args[2])
	}
	key := fmt.Sprintf("patch/temp/tone/%d/%s", tone, name)
	if err := resolveKey(key); err != nil {
		return nil, fmt.Errorf("unknown ctrl tone parameter %q", args[2])
	}
	return []patchCommand{{key: key, valArgs: args[3:]}}, nil
}

// waveNames are the WAVE category parameter names (tone switch, wave
// group/number/gain, FXM, tone delay), living directly in the tone block.
var waveNames = map[string]bool{
	"switch": true, "group": true, "groupid": true, "number": true,
	"gain": true, "fxmswitch": true,
	"fxmcolor": true, "fxmdepth": true, "delaymode": true, "delaytime": true,
}

// parseWave parses: [tone <1-4>] <param> <val>.
func parseWave(args []string) ([]patchCommand, error) {
	const usage = "usage: set/get patch wave [tone <1-4>] <param> <val>"
	if len(args) < 1 {
		return nil, fmt.Errorf("%s", usage)
	}
	tone, rest, _, err := parseToneArg(args)
	if err != nil {
		return nil, err
	}
	if len(rest) < 1 {
		return nil, fmt.Errorf("%s", usage)
	}
	name := aliasName(normalize(rest[0]))
	if !waveNames[name] {
		return nil, fmt.Errorf("unknown wave parameter %q (want switch, group, groupid, number, gain, fxmswitch, fxmcolor, fxmdepth, delaymode, delaytime)", rest[0])
	}
	key := fmt.Sprintf("patch/temp/tone/%d/%s", tone, name)
	if err := resolveKey(key); err != nil {
		return nil, fmt.Errorf("unknown wave parameter %q", rest[0])
	}
	return []patchCommand{{key: key, valArgs: rest[1:]}}, nil
}

// lfoParamNames are the LFO1/LFO2 parameter names.
var lfoParamNames = map[string]bool{
	"wave": true, "keytrig": true, "rate": true, "offset": true,
	"delay": true, "fade": true, "fadetime": true, "sync": true,
}

// lfoDepthNames are the LFO depth 1:2 destination names.
var lfoDepthNames = map[string]bool{
	"pitch": true, "filter": true, "amp": true, "pan": true,
}

// parseLFO parses the LFO category forms:
//
//	set/get patch lfo <1-2> tone <1-4> <param> <val>   (LFO1/LFO2)
//	set/get patch lfo <1-2> <param> <val>              (LFO depth 1:2, tone 1)
func parseLFO(args []string) ([]patchCommand, error) {
	const usage = "usage: set/get patch lfo <1-2> [tone <1-4>] <param> <val>"
	if len(args) < 1 {
		return nil, fmt.Errorf("%s", usage)
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 1 || n > 2 {
		return nil, fmt.Errorf("invalid LFO number %q (want 1 or 2)", args[0])
	}
	rest := args[1:]
	if len(rest) >= 2 && rest[0] == "tone" {
		t, err := strconv.Atoi(rest[1])
		if err != nil || t < 1 || t > 4 {
			return nil, fmt.Errorf("invalid tone %q (want 1-4)", rest[1])
		}
		if len(rest) < 3 {
			return nil, fmt.Errorf("usage: set/get patch lfo <1-2> tone <1-4> <param> <val>")
		}
		name := normalize(rest[2])
		if !lfoParamNames[name] {
			return nil, fmt.Errorf("unknown lfo parameter %q (want wave, keytrig, rate, offset, delay, fade, fadetime, sync)", rest[2])
		}
		key := fmt.Sprintf("patch/temp/tone/%d/lfo%d/%s", t, n, name)
		if err := resolveKey(key); err != nil {
			return nil, fmt.Errorf("unknown lfo parameter %q", rest[2])
		}
		return []patchCommand{{key: key, valArgs: rest[3:]}}, nil
	}
	if len(rest) < 1 {
		return nil, fmt.Errorf("usage: set/get patch lfo <1-2> <param> <val>")
	}
	name := normalize(rest[0])
	if !lfoDepthNames[name] {
		return nil, fmt.Errorf("unknown lfo depth parameter %q (want pitch, filter, amp, pan)", rest[0])
	}
	key := fmt.Sprintf("patch/temp/tone/1/lfo%d/%s", n, name)
	if err := resolveKey(key); err != nil {
		return nil, fmt.Errorf("unknown lfo depth parameter %q", rest[0])
	}
	return []patchCommand{{key: key, valArgs: rest[1:]}}, nil
}

// toneRangeNames are the per-Tone Velocity range and Key range parameter
// names, living directly in the tone block.
var toneRangeNames = map[string]bool{
	"vellower": true, "velupper": true, "velfade": true,
	"keylower": true, "keyupper": true,
}

// parsePatchTone parses: tone <1-4> <param> <val>.
func parsePatchTone(args []string) ([]patchCommand, error) {
	const usage = "usage: set/get patch tone <1-4> <param> <val>"
	if len(args) < 3 || args[0] != "tone" {
		return nil, fmt.Errorf("%s", usage)
	}
	tone, err := strconv.Atoi(args[1])
	if err != nil || tone < 1 || tone > 4 {
		return nil, fmt.Errorf("invalid tone %q (want 1-4)", args[1])
	}
	name := normalize(args[2])
	if !toneRangeNames[name] {
		return nil, fmt.Errorf("unknown tone parameter %q (want vellower, velupper, velfade, keylower, keyupper)", args[2])
	}
	key := fmt.Sprintf("patch/temp/tone/%d/%s", tone, name)
	if err := resolveKey(key); err != nil {
		return nil, fmt.Errorf("unknown tone parameter %q", args[2])
	}
	return []patchCommand{{key: key, valArgs: args[3:]}}, nil
}

// parsePatchParams resolves the category-specific argument list (args are
// the tokens after the category) to the device-layer parameter key it
// addresses together with the value arguments.
func parsePatchParams(cat string, args []string) ([]patchCommand, error) {
	switch cat {
	case "common":
		return parseCommon(args)
	case "efx":
		return parseEfx(args)
	case "pitch", "tvf", "tva":
		return parseToneCategory(cat, args)
	case "ctrl":
		return parseCtrl(args)
	case "wave":
		return parseWave(args)
	case "lfo":
		return parseLFO(args)
	case "tone":
		return parsePatchTone(args)
	}
	return nil, fmt.Errorf("unknown patch category %q", cat)
}

// parseParamValue converts the CLI value arguments into the UI value for
// the device layer: a string for the patch name, a pan string ("l13",
// "r13", "center") or signed int for pan parameters, and a single int for
// all other parameters (the device layer converts to wire bytes, including
// nibblized two-byte values).
func parseParamValue(p jv1080.Param, args []string) (any, error) {
	switch p.Kind {
	case jv1080.KindName:
		if len(args) == 0 {
			return nil, fmt.Errorf("%s: want the patch name as value", p.Key)
		}
		return strings.Join(args, " "), nil
	case jv1080.KindPan:
		if len(args) != 1 {
			return nil, fmt.Errorf("%s: want 1 value, got %d", p.Key, len(args))
		}
		s := strings.ToLower(args[0])
		if s == "center" || s == "c" || (len(s) >= 2 && (s[0] == 'l' || s[0] == 'r')) {
			return s, nil
		}
		return parseValue(args[0])
	default:
		if len(args) != 1 {
			return nil, fmt.Errorf("%s: want 1 value, got %d", p.Key, len(args))
		}
		return parseValue(args[0])
	}
}

// parseValue parses a decimal or hex (0x..) value, signed.
func parseValue(s string) (int, error) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		n, err := strconv.ParseInt(s[2:], 16, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid value %q", s)
		}
		return int(n), nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid value %q", s)
	}
	return n, nil
}

// formatParamValue formats a UI value returned by the device layer for
// printing: pan as l/r/center notation, the name as string, everything
// else as number.
func formatParamValue(p jv1080.Param, v any) string {
	switch p.Kind {
	case jv1080.KindName:
		if s, ok := v.(string); ok {
			return strconv.Quote(s)
		}
	case jv1080.KindPan:
		if n, ok := v.(int); ok {
			switch {
			case n == 0:
				return "center (0)"
			case n < 0:
				return fmt.Sprintf("l%d (%d)", -n, n)
			default:
				return fmt.Sprintf("r%d (%d)", n, n)
			}
		}
	}
	return fmt.Sprintf("%v", v)
}
