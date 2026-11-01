package jv1080

import (
	"strings"
	"testing"

	"go-jv-1080/internal/midi/message"
)

func TestResolveParamKeys(t *testing.T) {
	cases := []struct {
		key    string
		addr   message.Address
		kind   ParamKind
		min    int
		maxVal int
	}{
		{"patch/temp/common/level", message.Address{0x03, 0x00, 0x00, 0x2E}, KindRaw, 0, 127},
		{"patch/temp/common/pan", message.Address{0x03, 0x00, 0x00, 0x2F}, KindPan, -64, 63},
		{"patch/temp/common/octave", message.Address{0x03, 0x00, 0x00, 0x41}, KindBipolar, -3, 3},
		{"patch/temp/tone/1/level", message.Address{0x03, 0x00, 0x10, 0x65}, KindRaw, 0, 127},
		{"patch/temp/tone/4/cutoff", message.Address{0x03, 0x00, 0x16, 0x51}, KindRaw, 0, 127},
	}
	for _, c := range cases {
		p, err := resolveParam(c.key)
		if err != nil {
			t.Fatalf("resolveParam(%q): %v", c.key, err)
		}
		if p.area != c.addr {
			t.Errorf("%s: addr = %v, want %v", c.key, p.area, c.addr)
		}
		if p.Kind != c.kind || p.Min != c.min || p.Max != c.maxVal {
			t.Errorf("%s: kind/min/max = %v/%d/%d, want %v/%d/%d", c.key, p.Kind, p.Min, p.Max, c.kind, c.min, c.maxVal)
		}
	}
}

func TestResolveParamRejectsUnknown(t *testing.T) {
	for _, key := range []string{"", "patch/user/1/common/level", "patch/temp/common/nope", "foo/bar/baz/qux", "patch/temp/common"} {
		if _, err := resolveParam(key); err == nil {
			t.Errorf("resolveParam(%q) succeeded, want error", key)
		}
	}
}

func TestEncodeDecodeBipolar(t *testing.T) {
	p, err := resolveParam("patch/temp/common/octave")
	if err != nil {
		t.Fatal(err)
	}
	for ui, wire := range map[int]byte{-3: 0, -1: 2, 0: 3, 2: 5, 3: 6} {
		data, err := p.EncodeValue(ui)
		if err != nil {
			t.Fatalf("EncodeValue(%d): %v", ui, err)
		}
		if data[0] != wire {
			t.Errorf("octave %d: wire = %d, want %d", ui, data[0], wire)
		}
		back, err := p.DecodeValue(data)
		if err != nil {
			t.Fatal(err)
		}
		if back.(int) != ui {
			t.Errorf("octave decode(%d) = %v, want %d", wire, back, ui)
		}
	}
	if _, err := p.EncodeValue(4); err == nil {
		t.Error("EncodeValue(4) succeeded, want range error")
	}
}

func TestEncodeDecodePan(t *testing.T) {
	p, _ := resolveParam("patch/temp/common/pan")
	// UI and wire both start at 0 = center: UI -63..+63 maps to wire 1..127
	cases := []struct {
		ui   any
		wire byte
	}{
		{"C", 64}, {"center", 64}, {"L64", 0}, {"L63", 1}, {"L1", 63}, {"R1", 65}, {"R63", 127},
		{0, 64}, {-1, 63}, {-64, 0}, {1, 65}, {63, 127},
	}
	for _, c := range cases {
		data, err := p.EncodeValue(c.ui)
		if err != nil {
			t.Fatalf("EncodeValue(%v): %v", c.ui, err)
		}
		if data[0] != c.wire {
			t.Errorf("pan %v: wire = %d, want %d", c.ui, data[0], c.wire)
		}
		back, err := p.DecodeValue([]byte{c.wire})
		if err != nil {
			t.Fatal(err)
		}
		wantUI := rawToPan(int(c.wire))
		if back.(int) != wantUI {
			t.Errorf("pan decode(%d) = %v, want %d", c.wire, back, wantUI)
		}
	}
	// one tick left of center in the GUI must be wire 63 (L1), not 0 (L64)
	data, _ := p.EncodeValue(-1)
	if data[0] != 63 {
		t.Errorf("pan -1: wire = %d, want 63 (L1)", data[0])
	}
	if _, err := p.EncodeValue("R64"); err == nil {
		t.Error("EncodeValue(R64) succeeded, want error")
	}
	if _, err := p.EncodeValue("X9"); err == nil {
		t.Error("EncodeValue(X9) succeeded, want error")
	}
	if _, err := p.EncodeValue(64); err == nil {
		t.Error("EncodeValue(64) succeeded, want range error (> +63)")
	}
	if _, err := p.EncodeValue(-65); err == nil {
		t.Error("EncodeValue(-65) succeeded, want range error (< L64)")
	}
}

func TestRawToPan(t *testing.T) {
	cases := map[int]int{0: -64, 1: -63, 63: -1, 64: 0, 65: 1, 127: 63}
	for raw, want := range cases {
		if got := rawToPan(raw); got != want {
			t.Errorf("rawToPan(%d) = %d, want %d", raw, got, want)
		}
	}
}

func TestEncodeRawAndRange(t *testing.T) {
	p, _ := resolveParam("patch/temp/common/level")
	data, err := p.EncodeValue(100)
	if err != nil || data[0] != 100 {
		t.Fatalf("level 100: %v %v", data, err)
	}
	if _, err := p.EncodeValue(128); err == nil {
		t.Error("level 128: want range error")
	}
	if _, err := p.EncodeValue("abc"); err == nil {
		t.Error("level \"abc\": want conversion error")
	}
	back, _ := p.DecodeValue([]byte{100})
	if back.(int) != 100 {
		t.Errorf("level decode = %v, want 100", back)
	}
}

func TestEncodeName(t *testing.T) {
	p, _ := resolveParam("patch/temp/common/name")
	data, err := p.EncodeValue("Piano")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 12 || string(data[:5]) != "Piano" || data[5] != ' ' {
		t.Errorf("name encode = %q, want \"Piano\" padded to 12", data)
	}
	back, err := p.DecodeValue(data)
	if err != nil {
		t.Fatal(err)
	}
	if back.(string) != "Piano" {
		t.Errorf("name decode = %q, want Piano", back)
	}
}

func TestTempPatchKeyAt(t *testing.T) {
	cases := []struct {
		off int
		key string
		ok  bool
	}{
		{0x2E, "patch/temp/common/level", true},
		{0x2F, "patch/temp/common/pan", true},
		{0x41, "patch/temp/common/octave", true},
		{(0x10)*0x80 + 0x65, "patch/temp/tone/1/level", true},
		{0x0D, "", false}, // 0x0C is efxtype, a known key
		{-1, "", false},
	}
	for _, c := range cases {
		key, ok := tempPatchKeyAt(c.off)
		if ok != c.ok || key != c.key {
			t.Errorf("tempPatchKeyAt(%#x) = %q,%v want %q,%v", c.off, key, ok, c.key, c.ok)
		}
	}
}

func TestParamKeys(t *testing.T) {
	keys := ParamKeys()
	joined := strings.Join(keys, ",")
	for _, want := range []string{"patch/temp/common/level", "patch/temp/common/pan", "patch/temp/common/octave"} {
		if !strings.Contains(joined, want) {
			t.Errorf("ParamKeys() missing %q", want)
		}
	}
}
