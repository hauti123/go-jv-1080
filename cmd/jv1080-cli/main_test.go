package main

import (
	"testing"

	"go-jv-1080/internal/jv1080"
	"go-jv-1080/internal/midi/message"
)

func TestParseNote(t *testing.T) {
	cases := map[string]int{
		"c-1": 0, "c0": 12, "c1": 24, "c4": 60, "c5": 72, "g9": 127,
		"d3": 50, "e3": 52, "f#2": 42, "b3": 59,
		"bb3": 58, "ab3": 56, "eb3": 51, "C4": 60, "F#2": 42,
	}
	for name, want := range cases {
		got, err := parseNote(name)
		if err != nil {
			t.Errorf("parseNote(%q): %v", name, err)
			continue
		}
		if int(got) != want {
			t.Errorf("parseNote(%q) = %d, want %d", name, got, want)
		}
	}
	for _, bad := range []string{"", "c", "h1", "c#", "c10x", "c-2"} {
		if _, err := parseNote(bad); err == nil {
			t.Errorf("parseNote(%q) should fail", bad)
		}
	}
}

func TestParseAddress(t *testing.T) {
	addr, err := parseAddress("0x0000", "0x0006")
	if err != nil || addr != (message.Address{0x00, 0x00, 0x00, 0x06}) {
		t.Errorf("parseAddress(0x0000, 0x0006) = %s, %v", addr, err)
	}
	// 0x1000 0x0006: hi holds the first two address bytes (0x10 0x00),
	// lo the last two (0x00 0x06).
	addr, err = parseAddress("0x1000", "0x0006")
	if err != nil || addr != (message.Address{0x10, 0x00, 0x00, 0x06}) {
		t.Errorf("parseAddress(0x1000, 0x0006) = %s, %v", addr, err)
	}
	// 0x0300 0x002E -> 03 00 00 2E (temporary patch, patch level).
	addr, err = parseAddress("0x0300", "0x002E")
	if err != nil || addr != (message.Address{0x03, 0x00, 0x00, 0x2E}) {
		t.Errorf("parseAddress(0x0300, 0x002E) = %s, %v", addr, err)
	}
	// 0x1106 0x0000 -> 11 06 00 00 (user patch 7).
	addr, err = parseAddress("0x1106", "0x0000")
	if err != nil || addr != (message.Address{0x11, 0x06, 0x00, 0x00}) {
		t.Errorf("parseAddress(0x1106, 0x0000) = %s, %v", addr, err)
	}
	if _, err := parseAddress("0xzzzz", "0x0006"); err == nil {
		t.Error("parseAddress with invalid hi should fail")
	}
	if _, err := parseAddress("0x0000", "0xzzzz"); err == nil {
		t.Error("parseAddress with invalid lo should fail")
	}
}

func TestParseValues(t *testing.T) {
	v, err := parseValues([]string{"42"}, 1)
	if err != nil || len(v) != 1 || v[0] != 42 {
		t.Errorf("parseValues([42]) = %v, %v", v, err)
	}
	// Out-of-documentation-range values are accepted deliberately.
	v, err = parseValues([]string{"200"}, 1)
	if err != nil || v[0] != 200 {
		t.Errorf("parseValues([200]) = %v, %v (must be accepted)", v, err)
	}
	v, err = parseValues([]string{"0x0c", "0x01"}, 2)
	if err != nil || v[0] != 0x0C || v[1] != 0x01 {
		t.Errorf("parseValues nibble pair = % X, %v", v, err)
	}
	if _, err := parseValues([]string{"1", "2"}, 1); err == nil {
		t.Error("parseValues with wrong count should fail")
	}
	if _, err := parseValues([]string{"xyz"}, 1); err == nil {
		t.Error("parseValues with invalid value should fail")
	}
}

func TestParseBank(t *testing.T) {
	for _, s := range []string{"user", "pr-a", "pr-b", "pr-c", "gm", "card-a", "card-b", "exp-a", "exp-d", "USER", "PR-A"} {
		if _, err := parseBank(s); err != nil {
			t.Errorf("parseBank(%q): %v", s, err)
		}
	}
	if _, err := parseBank("bogus"); err == nil {
		t.Error("parseBank(bogus) should fail")
	}
}

func TestParsePatchParams(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"common", "level", "100"}, "patch/temp/common/level"},
		{[]string{"common", "pan", "l13"}, "patch/temp/common/pan"},
		{[]string{"common", "struct12", "3"}, "patch/temp/common/struct12"},
		{[]string{"common", "tempo", "120"}, "patch/temp/common/tempo"},
		{[]string{"common", "name", "Piano"}, "patch/temp/common/name"},
		{[]string{"efx", "reverbtype", "3"}, "patch/temp/common/reverbtype"},
		{[]string{"efx", "tone", "3", "output", "2"}, "patch/temp/tone/3/efx/output"},
		{[]string{"efx", "tone", "2", "reverb", "64"}, "patch/temp/tone/2/efx/reverb"},
		{[]string{"ctrl", "assign", "1"}, "patch/temp/common/assign"},
		{[]string{"ctrl", "3", "src", "4"}, "patch/temp/common/ctrl3src"},
		{[]string{"ctrl", "2", "hold", "1"}, "patch/temp/common/ctrl2hold"},
		{[]string{"ctrl", "tone", "4", "rxpan", "2"}, "patch/temp/tone/4/rxpan"},
		{[]string{"ctrl", "2", "tone", "3", "dst3", "7"}, "patch/temp/tone/3/ctrl2/dst3"},
		{[]string{"ctrl", "3", "tone", "1", "dep2", "70"}, "patch/temp/tone/1/ctrl3/dep2"},
		// WAVE, LFO, PITCH, TVF and TVA categories.
		{[]string{"wave", "gain", "2"}, "patch/temp/tone/1/gain"},
		{[]string{"wave", "tone", "2", "number", "194"}, "patch/temp/tone/2/number"},
		{[]string{"lfo", "1", "tone", "3", "rate", "60"}, "patch/temp/tone/3/lfo1/rate"},
		{[]string{"lfo", "2", "tone", "1", "fadetime", "40"}, "patch/temp/tone/1/lfo2/fadetime"},
		{[]string{"lfo", "2", "pitch", "70"}, "patch/temp/tone/1/lfo2/pitch"},
		{[]string{"lfo", "1", "filter", "64"}, "patch/temp/tone/1/lfo1/filter"},
		{[]string{"lfo", "2", "pan", "80"}, "patch/temp/tone/1/lfo2/pan"},
		{[]string{"pitch", "tone", "2", "coarse", "-12"}, "patch/temp/tone/2/pitch/coarse"},
		{[]string{"pitch", "tone", "2", "env", "depth", "12"}, "patch/temp/tone/2/pitch/depth"},
		{[]string{"pitch", "env", "l4", "64"}, "patch/temp/tone/1/pitch/l4"},
		{[]string{"pitch", "env", "keyfollow", "5"}, "patch/temp/tone/1/pitch/envkeyfollow"},
		{[]string{"tvf", "tone", "1", "cutoff", "100"}, "patch/temp/tone/1/tvf/cutoff"},
		{[]string{"tvf", "tone", "1", "envdepth", "64"}, "patch/temp/tone/1/tvf/envdepth"},
		{[]string{"tvf", "tone", "4", "env", "t2", "50"}, "patch/temp/tone/4/tvf/t2"},
		{[]string{"tvf", "env", "l1", "60"}, "patch/temp/tone/1/tvf/l1"},
		{[]string{"tva", "tone", "2", "level", "110"}, "patch/temp/tone/2/tva/level"},
		{[]string{"tva", "tone", "2", "biaspoint", "60"}, "patch/temp/tone/2/tva/biaspoint"},
		{[]string{"tva", "tone", "3", "env", "t1", "40"}, "patch/temp/tone/3/tva/t1"},
		{[]string{"tva", "env", "l3", "70"}, "patch/temp/tone/1/tva/l3"},
		{[]string{"tone", "tone", "3", "keyupper", "100"}, "patch/temp/tone/3/keyupper"},
	}
	for _, c := range cases {
		cmds, err := parsePatchParams(c.args[0], c.args[1:])
		if err != nil {
			t.Errorf("parsePatchParams(%v): %v", c.args, err)
			continue
		}
		if len(cmds) != 1 || cmds[0].key != c.want {
			t.Errorf("parsePatchParams(%v) = %+v, want key %q", c.args, cmds, c.want)
			continue
		}
		// every produced key must resolve in the device layer
		if _, err := jv1080.ResolveParam(cmds[0].key); err != nil {
			t.Errorf("parsePatchParams(%v) produced unresolvable key %q: %v", c.args, cmds[0].key, err)
		}
	}
	// Errors: unknown names, wrong blocks, missing markers.
	for _, bad := range [][]string{
		{"common", "bogus", "1"},
		{"efx", "level", "1"},               // level is not an efx tone param
		{"efx", "tone", "9", "output", "1"}, // bad tone
		{"ctrl", "4", "hold", "1"},          // hold index 1-3
		{"ctrl", "1", "src", "1"},           // src index 2-3
		{"ctrl", "tone", "1", "dst1", "1"},  // dst needs control index
		{"lfo", "3", "rate", "1"},           // LFO 1-2
		{"pitch", "tone", "1", "t1", "1"},   // needs env marker
		{"pitch", "env", "bogus", "1"},
		{"tvf", "bogus", "1"},
		{"tva", "env", "l4", "1"}, // TVA env has l1-l3 only
		{"tone", "tone", "1", "bogus", "1"},
		{"tone", "tone", "9", "keyupper", "1"},
	} {
		if _, err := parsePatchParams(bad[0], bad[1:]); err == nil {
			t.Errorf("parsePatchParams(%v) should fail", bad)
		}
	}
}

func TestParseParamValue(t *testing.T) {
	p, err := jv1080.ResolveParam("patch/temp/common/pan")
	if err != nil {
		t.Fatal(err)
	}
	v, err := parseParamValue(p, []string{"l13"})
	if err != nil || v != "l13" {
		t.Errorf("parseParamValue(pan, l13) = %v, %v", v, err)
	}
	v, err = parseParamValue(p, []string{"-13"})
	if err != nil || v != -13 {
		t.Errorf("parseParamValue(pan, -13) = %v, %v", v, err)
	}
	p, _ = jv1080.ResolveParam("patch/temp/tone/2/pitch/coarse")
	v, err = parseParamValue(p, []string{"-12"})
	if err != nil || v != -12 {
		t.Errorf("parseParamValue(coarse, -12) = %v, %v", v, err)
	}
	p, _ = jv1080.ResolveParam("patch/temp/tone/2/number")
	v, err = parseParamValue(p, []string{"0xC2"})
	if err != nil || v != 194 {
		t.Errorf("parseParamValue(number, 0xC2) = %v, %v", v, err)
	}
	p, _ = jv1080.ResolveParam("patch/temp/common/name")
	v, err = parseParamValue(p, []string{"My", "Patch"})
	if err != nil || v != "My Patch" {
		t.Errorf("parseParamValue(name) = %v, %v", v, err)
	}
	// wrong value counts fail
	p, _ = jv1080.ResolveParam("patch/temp/common/level")
	if _, err := parseParamValue(p, []string{"1", "2"}); err == nil {
		t.Error("parseParamValue(level, 2 values) should fail")
	}
}

func TestFormatParamValue(t *testing.T) {
	pan, _ := jv1080.ResolveParam("patch/temp/common/pan")
	if got := formatParamValue(pan, -13); got != "l13 (-13)" {
		t.Errorf("formatParamValue(pan, -13) = %q", got)
	}
	if got := formatParamValue(pan, 0); got != "center (0)" {
		t.Errorf("formatParamValue(pan, 0) = %q", got)
	}
	if got := formatParamValue(pan, 13); got != "r13 (13)" {
		t.Errorf("formatParamValue(pan, 13) = %q", got)
	}
	level, _ := jv1080.ResolveParam("patch/temp/common/level")
	if got := formatParamValue(level, 100); got != "100" {
		t.Errorf("formatParamValue(level, 100) = %q", got)
	}
	name, _ := jv1080.ResolveParam("patch/temp/common/name")
	if got := formatParamValue(name, "Piano"); got != `"Piano"` {
		t.Errorf("formatParamValue(name, Piano) = %q", got)
	}
}

func TestScriptReadCommandsResolve(t *testing.T) {
	if len(scriptReadCommands) == 0 {
		t.Fatal("no script commands")
	}
	for _, cmd := range scriptReadCommands {
		cmds, err := parsePatchParams(cmd[0], cmd[1:])
		if err != nil {
			t.Errorf("parsePatchParams(%v): %v", cmd, err)
			continue
		}
		if len(cmds) != 1 {
			t.Errorf("parsePatchParams(%v) = %d commands, want 1", cmd, len(cmds))
			continue
		}
		if _, err := jv1080.ResolveParam(cmds[0].key); err != nil {
			t.Errorf("parsePatchParams(%v) key %q does not resolve: %v", cmd, cmds[0].key, err)
		}
	}
}
