package message

import (
	"bytes"
	"testing"
)

func TestSplitMessages(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
		want [][]byte
	}{
		{
			"single CC",
			[]byte{0xB0, 0x00, 0x50},
			[][]byte{{0xB0, 0x00, 0x50}},
		},
		{
			"bank select pair (CC0 + CC32)",
			[]byte{0xB0, 0x00, 0x50, 0xB0, 0x20, 0x00},
			[][]byte{{0xB0, 0x00, 0x50}, {0xB0, 0x20, 0x00}},
		},
		{
			"RPN select pair",
			[]byte{0xB0, 0x65, 0x00, 0xB0, 0x64, 0x02},
			[][]byte{{0xB0, 0x65, 0x00}, {0xB0, 0x64, 0x02}},
		},
		{
			"program change (one data byte)",
			[]byte{0xC9, 0x28},
			[][]byte{{0xC9, 0x28}},
		},
		{
			"channel pressure (one data byte)",
			[]byte{0xD0, 0x40},
			[][]byte{{0xD0, 0x40}},
		},
		{
			"sysex",
			[]byte{0xF0, 0x7E, 0x10, 0x06, 0x01, 0xF7},
			[][]byte{{0xF0, 0x7E, 0x10, 0x06, 0x01, 0xF7}},
		},
		{
			"sysex followed by CC",
			[]byte{0xF0, 0x7E, 0x10, 0x06, 0x01, 0xF7, 0xB0, 0x07, 0x64},
			[][]byte{{0xF0, 0x7E, 0x10, 0x06, 0x01, 0xF7}, {0xB0, 0x07, 0x64}},
		},
		{
			"realtime clock between channel messages",
			[]byte{0x90, 60, 100, 0xF8, 0xB0, 0x40, 0x7F},
			[][]byte{{0x90, 60, 100}, {0xF8}, {0xB0, 0x40, 0x7F}},
		},
		{
			"unterminated sysex",
			[]byte{0xF0, 0x41, 0x10},
			[][]byte{{0xF0, 0x41, 0x10}},
		},
		{
			"truncated channel message",
			[]byte{0xB0, 0x00},
			[][]byte{{0xB0, 0x00}},
		},
		{
			"stray data bytes skipped",
			[]byte{0x00, 0x50, 0xB0, 0x00, 0x50},
			[][]byte{{0xB0, 0x00, 0x50}},
		},
		{
			"empty",
			nil,
			nil,
		},
	}
	for _, c := range cases {
		got := SplitMessages(c.in)
		if len(got) != len(c.want) {
			t.Errorf("%s: SplitMessages = %d messages, want %d", c.name, len(got), len(c.want))
			continue
		}
		for i := range got {
			if !bytes.Equal(got[i], c.want[i]) {
				t.Errorf("%s: message %d = % X, want % X", c.name, i, got[i], c.want[i])
			}
		}
	}
}

// TestSplitMessagesRoundTrip checks that splitting a message stream and
// concatenating the pieces yields the original stream.
func TestSplitMessagesRoundTrip(t *testing.T) {
	in := []byte{0xB0, 0x00, 0x50, 0xB0, 0x20, 0x00, 0xF0, 0x7E, 0x10, 0x06, 0x01, 0xF7, 0xC9, 0x28}
	var joined []byte
	for _, m := range SplitMessages(in) {
		joined = append(joined, m...)
	}
	if !bytes.Equal(joined, in) {
		t.Errorf("round trip = % X, want % X", joined, in)
	}
}
