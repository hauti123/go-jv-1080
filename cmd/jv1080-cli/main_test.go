package main

import (
	"testing"

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

func TestPanValue(t *testing.T) {
	cases := map[string]int{
		"l1": 63, "l13": 51, "l63": 1, "r1": 65, "r13": 77, "r63": 127,
		"center": 64, "c": 64,
	}
	for in, want := range cases {
		got, ok := panValue(in)
		if !ok || got != want {
			t.Errorf("panValue(%q) = %d, %v; want %d", in, got, ok, want)
		}
	}
	for _, bad := range []string{"l0", "l64", "x5", ""} {
		if _, ok := panValue(bad); ok {
			t.Errorf("panValue(%q) should fail", bad)
		}
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

func TestCommonParam(t *testing.T) {
	off, size, err := commonParam("level")
	if err != nil || off != 0x2E || size != 1 {
		t.Errorf("commonParam(level) = %d, %d, %v", off, size, err)
	}
	if _, _, err := commonParam("bogus"); err == nil {
		t.Error("commonParam(bogus) should fail")
	}
}

func TestToneParam(t *testing.T) {
	off, size, err := toneParam("wave")
	if err != nil || off != 0x03 || size != 2 {
		t.Errorf("toneParam(wave) = %d, %d, %v (nibblized)", off, size, err)
	}
	if _, _, err := toneParam("bogus"); err == nil {
		t.Error("toneParam(bogus) should fail")
	}
}

func TestFormatValue(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		want string
	}{
		{"level", []byte{44}, "44"},
		{"level", []byte{0x7F}, "127"},
		{"pan", []byte{51}, "l13 (51)"},
		{"pan", []byte{77}, "r13 (77)"},
		{"pan", []byte{64}, "center (64)"},
		{"coarse", []byte{36}, "-12 (36)"},
		{"fine", []byte{75}, "25 (75)"},
		{"wave", []byte{0x0C, 0x01}, "193"}, // nibblized raw 0xC1 = 193 -> wave 194
		{"switch", []byte{1}, "1"},
	}
	for _, c := range cases {
		if got := formatValue(c.name, c.data); got != c.want {
			t.Errorf("formatValue(%q, % X) = %q, want %q", c.name, c.data, got, c.want)
		}
	}
}

func TestWaveAndTVFParam(t *testing.T) {
	if off, _, err := waveParam("gain"); err != nil || off != 0x05 {
		t.Errorf("waveParam(gain) = %d, %v", off, err)
	}
	if off, _, err := tvfParam("cutoff"); err != nil || off != 0x51 {
		t.Errorf("tvfParam(cutoff) = %d, %v", off, err)
	}
	if len(tvfEnvOffsets) != 8 || tvfEnvOffsets[0] != 0x5B || tvfEnvOffsets[7] != 0x62 {
		t.Errorf("tvfEnvOffsets = % X", tvfEnvOffsets)
	}
}

func TestParsePatchParams(t *testing.T) {
	cases := []struct {
		args []string
		want patchParam
	}{
		{[]string{"common", "level", "100"}, patchParam{"common", 0, 0x2E, 1, "level"}},
		{[]string{"common", "pan", "l13"}, patchParam{"common", 0, 0x2F, 1, "pan"}},
		{[]string{"common", "struct12", "3"}, patchParam{"common", 0, 0x44, 1, "struct12"}},
		{[]string{"common", "tempo", "0x01", "0x0A"}, patchParam{"common", 0, 0x2C, 2, "tempo"}},
		{[]string{"efx", "reverbtype", "3"}, patchParam{"common", 0, 0x27, 1, "reverbtype"}},
		{[]string{"efx", "tone", "3", "output", "2"}, patchParam{"tone", 3, 0x7D, 1, "output"}},
		{[]string{"efx", "tone", "2", "reverb", "64"}, patchParam{"tone", 2, 0x80, 1, "reverb"}},
		{[]string{"ctrl", "assign", "1"}, patchParam{"common", 0, 0x33, 1, "assign"}},
		{[]string{"ctrl", "3", "src", "4"}, patchParam{"common", 0, 0x3B, 1, "src"}},
		{[]string{"ctrl", "2", "hold", "1"}, patchParam{"common", 0, 0x3E, 1, "hold"}},
		{[]string{"ctrl", "tone", "4", "rxpan", "2"}, patchParam{"tone", 4, 0x14, 1, "rxpan"}},
		{[]string{"ctrl", "2", "tone", "3", "dst3", "7"}, patchParam{"tone", 3, 0x21, 1, "dst3"}},
		{[]string{"ctrl", "3", "tone", "1", "dep2", "70"}, patchParam{"tone", 1, 0x28, 1, "dep2"}},
	}
	for _, c := range cases {
		cmds, err := parsePatchParams(c.args[0], c.args[1:])
		if err != nil {
			t.Errorf("parsePatchParams(%v): %v", c.args, err)
			continue
		}
		if len(cmds) != 1 || cmds[0].param != c.want {
			t.Errorf("parsePatchParams(%v) = %+v, want %+v", c.args, cmds, c.want)
		}
	}
	// WAVE, LFO, PITCH, TVF and TVA categories.
	special := []struct {
		args []string
		want patchParam
	}{
		{[]string{"wave", "gain", "2"}, patchParam{"tone", 1, 0x05, 1, "gain"}},
		{[]string{"wave", "tone", "2", "number", "0x0C", "0x01"}, patchParam{"tone", 2, 0x03, 2, "number"}},
		{[]string{"lfo", "1", "tone", "3", "rate", "60"}, patchParam{"tone", 3, 0x2F, 1, "rate"}},
		{[]string{"lfo", "2", "tone", "1", "fadetime", "40"}, patchParam{"tone", 1, 0x3B, 1, "fadetime"}},
		{[]string{"lfo", "2", "pitch", "70"}, patchParam{"tone", 1, 0x4F, 1, "pitch"}},
		{[]string{"lfo", "1", "filter", "64"}, patchParam{"tone", 1, 0x63, 1, "filter"}},
		{[]string{"lfo", "2", "pan", "80"}, patchParam{"tone", 1, 0x7C, 1, "pan"}},
		{[]string{"pitch", "tone", "2", "coarse", "-12"}, patchParam{"tone", 2, 0x3D, 1, "coarse"}},
		{[]string{"pitch", "tone", "2", "env", "depth", "12"}, patchParam{"tone", 2, 0x41, 1, "depth"}},
		{[]string{"pitch", "env", "l4", "64"}, patchParam{"tone", 1, 0x4D, 1, "l4"}},
		{[]string{"tvf", "tone", "1", "cutoff", "100"}, patchParam{"tone", 1, 0x51, 1, "cutoff"}},
		{[]string{"tvf", "tone", "1", "envdepth", "64"}, patchParam{"tone", 1, 0x55, 1, "envdepth"}},
		{[]string{"tvf", "tone", "4", "env", "t2", "50"}, patchParam{"tone", 4, 0x5C, 1, "t2"}},
		{[]string{"tvf", "env", "l1", "60"}, patchParam{"tone", 1, 0x5F, 1, "l1"}},
		{[]string{"tva", "tone", "2", "level", "110"}, patchParam{"tone", 2, 0x65, 1, "level"}},
		{[]string{"tva", "tone", "2", "biaspoint", "60"}, patchParam{"tone", 2, 0x67, 1, "biaspoint"}},
		{[]string{"tva", "tone", "3", "env", "t1", "40"}, patchParam{"tone", 3, 0x6E, 1, "t1"}},
		{[]string{"tva", "env", "l3", "70"}, patchParam{"tone", 1, 0x74, 1, "l3"}},
		{[]string{"tone", "tone", "3", "keyupper", "100"}, patchParam{"tone", 3, 0x0F, 1, "keyupper"}},
	}
	for _, c := range special {
		cmds, err := parsePatchParams(c.args[0], c.args[1:])
		if err != nil {
			t.Errorf("parsePatchParams(%v): %v", c.args, err)
			continue
		}
		if len(cmds) != 1 || cmds[0].param != c.want {
			t.Errorf("parsePatchParams(%v) = %+v, want %+v", c.args, cmds, c.want)
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

func TestParsePatchValues(t *testing.T) {
	v, err := parsePatchValues(patchParam{"common", 0, 0x2F, 1, "pan"}, []string{"l13"})
	if err != nil || v[0] != 51 {
		t.Errorf("parsePatchValues(pan l13) = % X, %v", v, err)
	}
	v, err = parsePatchValues(patchParam{"tone", 1, 0x3D, 1, "coarse"}, []string{"-12"})
	if err != nil || v[0] != 36 {
		t.Errorf("parsePatchValues(coarse -12) = % X, %v", v, err)
	}
	v, err = parsePatchValues(patchParam{"tone", 1, 0x03, 2, "number"}, []string{"0x0C", "0x01"})
	if err != nil || len(v) != 2 || v[0] != 0x0C || v[1] != 0x01 {
		t.Errorf("parsePatchValues(number) = % X, %v", v, err)
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
		}
	}
}
