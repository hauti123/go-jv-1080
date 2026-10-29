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
	// 0x1000 0x0006: hi = 0x1000 -> 0x20 0x00, lo = 0x0006 -> 0x00 0x06.
	addr, err = parseAddress("0x1000", "0x0006")
	if err != nil || addr != (message.Address{0x20, 0x00, 0x00, 0x06}) {
		t.Errorf("parseAddress(0x1000, 0x0006) = %s, %v", addr, err)
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
	if len(tvfEnvOffsets) != 8 || tvfEnvOffsets[0] != 0x55 || tvfEnvOffsets[7] != 0x5C {
		t.Errorf("tvfEnvOffsets = % X", tvfEnvOffsets)
	}
}
