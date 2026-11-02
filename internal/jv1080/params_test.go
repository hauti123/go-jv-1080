package jv1080

import (
	"fmt"
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
		{"patch/temp/common/tempo", message.Address{0x03, 0x00, 0x00, 0x2C}, KindNibblized, 20, 250},
		{"patch/temp/tone/1/tva/level", message.Address{0x03, 0x00, 0x10, 0x65}, KindRaw, 0, 127},
		{"patch/temp/tone/2/tva/level", message.Address{0x03, 0x00, 0x12, 0x65}, KindRaw, 0, 127},
		{"patch/temp/tone/4/tvf/cutoff", message.Address{0x03, 0x00, 0x16, 0x51}, KindRaw, 0, 127},
		{"patch/temp/tone/3/number", message.Address{0x03, 0x00, 0x14, 0x03}, KindNibblized, 1, 255},
		{"patch/temp/tone/2/lfo1/rate", message.Address{0x03, 0x00, 0x12, 0x2F}, KindRaw, 0, 127},
		{"patch/temp/tone/4/lfo2/sync", message.Address{0x03, 0x00, 0x16, 0x3C}, KindRaw, 0, 127},
		{"patch/temp/tone/2/lfo2/filter", message.Address{0x03, 0x00, 0x12, 0x64}, KindBipolar, -63, 63},
		{"patch/temp/tone/1/ctrl2/dst3", message.Address{0x03, 0x00, 0x10, 0x21}, KindRaw, 0, 18},
		{"patch/temp/tone/1/ctrl3/dep2", message.Address{0x03, 0x00, 0x10, 0x28}, KindBipolar, -63, 63},
		{"patch/temp/common/ctrl2src", message.Address{0x03, 0x00, 0x00, 0x3A}, KindRaw, 0, 15},
		{"patch/temp/common/ctrl3hold", message.Address{0x03, 0x00, 0x00, 0x3F}, KindRaw, 0, 2},
		{"patch/temp/tone/4/efx/reverb", message.Address{0x03, 0x00, 0x17, 0x00}, KindRaw, 0, 127}, // offset 0x80 of tone 4 wraps to 17 00
		{"patch/temp/tone/2/tva/pan", message.Address{0x03, 0x00, 0x12, 0x77}, KindPan, -64, 63},
		// pitch envelope of tone 1 (third address byte 0x10)
		{"patch/temp/tone/1/pitch/t1", message.Address{0x03, 0x00, 0x10, 0x46}, KindRaw, 0, 127},
		{"patch/temp/tone/1/pitch/t4", message.Address{0x03, 0x00, 0x10, 0x49}, KindRaw, 0, 127},
		{"patch/temp/tone/1/pitch/l1", message.Address{0x03, 0x00, 0x10, 0x4A}, KindBipolar, -63, 63},
		{"patch/temp/tone/1/pitch/l4", message.Address{0x03, 0x00, 0x10, 0x4D}, KindBipolar, -63, 63},
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
	for _, key := range []string{"", "patch/user/1/common/level", "patch/temp/common/nope", "foo/bar/baz/qux", "patch/temp/common", "patch/temp/tone/5/tvf/cutoff", "patch/temp/tone/1/tvf/nope", "patch/temp/tone/1/lfo1/nope"} {
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

func TestEncodeDecodePitchEnvLevels(t *testing.T) {
	p, err := resolveParam("patch/temp/tone/1/pitch/l1")
	if err != nil {
		t.Fatal(err)
	}
	for ui, wire := range map[int]byte{-63: 0, -1: 62, 0: 63, 63: 126} {
		data, err := p.EncodeValue(ui)
		if err != nil {
			t.Fatalf("l1 EncodeValue(%d): %v", ui, err)
		}
		if data[0] != wire {
			t.Errorf("l1 %d: wire = %d, want %d", ui, data[0], wire)
		}
		back, _ := p.DecodeValue(data)
		if back.(int) != ui {
			t.Errorf("l1 decode(%d) = %v, want %d", wire, back, ui)
		}
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
		{(0x10)*0x80 + 0x65, "patch/temp/tone/1/tva/level", true},
		{0x0D, "patch/temp/common/efxprm1", true},
		{-1, "", false},
		{0x10*0x80 + 0x46, "patch/temp/tone/1/pitch/t1", true},
		{0x10*0x80 + 0x4D, "patch/temp/tone/1/pitch/l4", true},
		{0x16*0x80 + 0x46, "patch/temp/tone/4/pitch/t1", true},
		{0x12*0x80 + 0x51, "patch/temp/tone/2/tvf/cutoff", true},
		{0x10*0x80 + 0x2F, "patch/temp/tone/1/lfo1/rate", true},
		{0x10*0x80 + 0x21, "patch/temp/tone/1/ctrl2/dst3", true},
	}
	for _, c := range cases {
		key, ok := tempPatchKeyAt(c.off)
		if ok != c.ok || key != c.key {
			t.Errorf("tempPatchKeyAt(%#x) = %q,%v want %q,%v", c.off, key, ok, c.key, c.ok)
		}
	}
}

func TestHandleDeviceMessageDispatch(t *testing.T) {
	c := &Client{deviceID: 0x10}
	var gotKey string
	var gotVal any
	c.OnParameterChange(func(key string, value any) { gotKey, gotVal = key, value })

	// DT1 writing level (0x2E) of the temporary patch
	ds := message.DataSet{DeviceID: 0x10, Address: message.Address{0x03, 0x00, 0x00, 0x2E}, Data: []byte{100}}
	c.handleDeviceMessage(ds.Bytes())
	if gotKey != "patch/temp/common/level" || gotVal.(int) != 100 {
		t.Errorf("got %q %v, want patch/temp/common/level 100", gotKey, gotVal)
	}

	// DT1 for another device ID is ignored
	gotKey = ""
	ds.DeviceID = 0x11
	c.handleDeviceMessage(ds.Bytes())
	if gotKey != "" {
		t.Errorf("foreign device ID produced update %q", gotKey)
	}

	// DT1 at an offset with no registered parameter is ignored (0x3D is
	// between ctrl1hold 0x3D... use an offset in a gap, e.g. tone area of
	// common block byte 0x48+ is past the common block; use 0x03 00 00 47+1)
	ds.DeviceID = 0x10
	ds.Address = message.Address{0x03, 0x00, 0x00, 0x48} // just past the common block
	c.handleDeviceMessage(ds.Bytes())
	if gotKey != "" {
		t.Errorf("unknown offset produced update %q", gotKey)
	}

	// non-sysex message is ignored
	c.handleDeviceMessage([]byte{0x90, 0x3C, 0x40})
	if gotKey != "" {
		t.Errorf("note on produced update %q", gotKey)
	}

	// pan decode: wire 51 -> UI -13 (L13)
	gotKey, gotVal = "", nil
	ds.Address = message.Address{0x03, 0x00, 0x00, 0x2F}
	ds.Data = []byte{51}
	c.handleDeviceMessage(ds.Bytes())
	if gotKey != "patch/temp/common/pan" || gotVal.(int) != -13 {
		t.Errorf("pan: got %q %v, want patch/temp/common/pan -13", gotKey, gotVal)
	}
}

func TestListParams(t *testing.T) {
	params := ListParams()
	if len(params) == 0 {
		t.Fatal("ListParams() returned no parameters")
	}
	joined := make([]string, 0, len(params))
	for _, p := range params {
		joined = append(joined, p.Key)
	}
	s := strings.Join(joined, ",")
	for _, want := range []string{
		"patch/temp/common/level", "patch/temp/common/pan", "patch/temp/common/octave",
		"patch/temp/common/name", "patch/temp/common/tempo",
		"patch/temp/tone/1/tva/level", "patch/temp/tone/4/tvf/cutoff",
		"patch/temp/tone/2/lfo2/rate", "patch/temp/tone/3/ctrl1/dep4",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("ListParams() missing %q", want)
		}
	}
	for _, p := range params {
		q, err := ResolveParam(p.Key)
		if err != nil || q.Key != p.Key || q.Kind != p.Kind || q.Offset != p.Offset {
			t.Errorf("ResolveParam(%q) = %+v, %v; want same as listed %+v", p.Key, q, err, p)
		}
	}
}

func TestEncodeDecodeNibblized(t *testing.T) {
	p, err := resolveParam("patch/temp/tone/1/number")
	if err != nil {
		t.Fatal(err)
	}
	for ui, wire := range map[int][2]byte{1: {0x00, 0x00}, 2: {0x00, 0x01}, 194: {0x0C, 0x01}, 255: {0x0F, 0x0E}} {
		data, err := p.EncodeValue(ui)
		if err != nil {
			t.Fatalf("number EncodeValue(%d): %v", ui, err)
		}
		if data[0] != wire[0] || data[1] != wire[1] {
			t.Errorf("number %d: wire = % X, want % X", ui, data, wire)
		}
		back, err := p.DecodeValue(data)
		if err != nil {
			t.Fatal(err)
		}
		if back.(int) != ui {
			t.Errorf("number decode(% X) = %v, want %d", data, back, ui)
		}
	}
	if _, err := p.EncodeValue(0); err == nil {
		t.Error("number 0: want range error")
	}
	if _, err := p.EncodeValue(256); err == nil {
		t.Error("number 256: want range error")
	}
	p, err = resolveParam("patch/temp/common/tempo")
	if err != nil {
		t.Fatal(err)
	}
	data, err := p.EncodeValue(120)
	if err != nil || data[0] != 0x06 || data[1] != 0x04 {
		t.Errorf("tempo 120: wire = % X, %v; want 06 04", data, err)
	}
	back, _ := p.DecodeValue(data)
	if back.(int) != 120 {
		t.Errorf("tempo decode = %v, want 120", back)
	}
}

func TestToneParamsAllTones(t *testing.T) {
	bases := map[int]int{1: 0x10, 2: 0x12, 3: 0x14, 4: 0x16}
	for key, off := range map[string]int{
		"patch/temp/tone/%d/tvf/cutoff": 0x51,
		"patch/temp/tone/%d/tva/level":  0x65,
		"patch/temp/tone/%d/pitch/l4":   0x4D,
		"patch/temp/tone/%d/efx/reverb": 0x80,
		"patch/temp/tone/%d/lfo1/wave":  0x2D,
		"patch/temp/tone/%d/ctrl1/dst1": 0x15,
		"patch/temp/tone/%d/velfade":    0x0B,
		"patch/temp/tone/%d/rxpan":      0x14,
		"patch/temp/tone/%d/switch":     0x00,
		"patch/temp/tone/%d/number":     0x03,
		"patch/temp/tone/%d/tvf/l3":     0x61,
		"patch/temp/tone/%d/tva/l3":     0x74,
		"patch/temp/tone/%d/lfo2/pan":   0x7C,
		"patch/temp/tone/%d/ctrl3/dep4": 0x2C,
	} {
		for tone, base := range bases {
			k := fmt.Sprintf(key, tone)
			p, err := resolveParam(k)
			if err != nil {
				t.Errorf("resolveParam(%q): %v", k, err)
				continue
			}
			want := TemporaryPatchArea().Add(base*0x80 + off)
			if p.area != want {
				t.Errorf("%s: addr = %v, want %v", p.Key, p.area, want)
			}
		}
	}
}
