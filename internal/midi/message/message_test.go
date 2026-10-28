package message

import (
	"bytes"
	"testing"
)

func TestChecksum(t *testing.T) {
	// Example from the JV MIDI implementation: F0 41 10 6A 11 03 00 10 00 00 00 07 01 65 F7
	if got := Checksum([]byte{0x03, 0x00, 0x10, 0x00, 0x00, 0x00, 0x07, 0x01}); got != 0x65 {
		t.Errorf("Checksum = %02X, want 65", got)
	}
	// Sum already a multiple of 128 -> checksum 0.
	if got := Checksum([]byte{0x40, 0x40}); got != 0x00 {
		t.Errorf("Checksum = %02X, want 00", got)
	}
}

func TestAddressRoundTrip(t *testing.T) {
	for _, n := range []int{0, 1, 0x7F, 0x80, 0x3FFF, 0x100000, 0x0FFFFFFF} {
		if got := AddressFromInt(n).Addr(); got != n {
			t.Errorf("AddressFromInt(%d).Addr() = %d", n, got)
		}
	}
	// Adding crosses 7-bit boundaries: 03 00 10 7F + 1 = 03 00 11 00.
	a := Address{0x03, 0x00, 0x10, 0x7F}.Add(1)
	if a != (Address{0x03, 0x00, 0x11, 0x00}) {
		t.Errorf("Add = %s, want 03 00 11 00", a)
	}
}

func TestDataRequestRoundTrip(t *testing.T) {
	// Request the four tones of the temporary patch (560 bytes = 0x230).
	req := DataRequest{DeviceID: 0x10, Address: Address{0x03, 0x00, 0x10, 0x00}, Size: 0x230}
	want := []byte{0xF0, 0x41, 0x10, 0x6A, 0x11, 0x03, 0x00, 0x10, 0x00, 0x00, 0x00, 0x04, 0x30, 0x39, 0xF7}
	if got := req.Bytes(); !bytes.Equal(got, want) {
		t.Errorf("Bytes = % X, want % X", got, want)
	}
	back, err := ParseDataRequest(req.Bytes())
	if err != nil {
		t.Fatalf("ParseDataRequest: %v", err)
	}
	if back.DeviceID != req.DeviceID || back.Address != req.Address || back.Size != req.Size {
		t.Errorf("round trip = %+v, want %+v", back, req)
	}
}

func TestDataRequestSingleParameter(t *testing.T) {
	// Turn off Tone 1 of the temporary patch (documented example from the web):
	// F0 41 10 6A 12 03 00 10 00 00 6D F7 is a DT1; the matching RQ1 for one
	// byte at 03 00 10 00 is F0 41 10 6A 11 03 00 10 00 00 00 00 01 <sum> F7.
	req := DataRequest{DeviceID: 0x10, Address: Address{0x03, 0x00, 0x10, 0x00}, Size: 1}
	want := []byte{0xF0, 0x41, 0x10, 0x6A, 0x11, 0x03, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x01, 0x6C, 0xF7}
	if got := req.Bytes(); !bytes.Equal(got, want) {
		t.Errorf("Bytes = % X, want % X", got, want)
	}
}

func TestDataSetRoundTrip(t *testing.T) {
	// Documented example: switch off Tone 1 of the temporary patch.
	ds := DataSet{DeviceID: 0x10, Address: Address{0x03, 0x00, 0x10, 0x00}, Data: []byte{0x00}}
	want := []byte{0xF0, 0x41, 0x10, 0x6A, 0x12, 0x03, 0x00, 0x10, 0x00, 0x00, 0x6D, 0xF7}
	if got := ds.Bytes(); !bytes.Equal(got, want) {
		t.Errorf("Bytes = % X, want % X", got, want)
	}
	back, err := ParseDataSet(ds.Bytes())
	if err != nil {
		t.Fatalf("ParseDataSet: %v", err)
	}
	if back.DeviceID != ds.DeviceID || back.Address != ds.Address || !bytes.Equal(back.Data, ds.Data) {
		t.Errorf("round trip = %+v, want %+v", back, ds)
	}
}

func TestParseDataSetBadChecksum(t *testing.T) {
	b := []byte{0xF0, 0x41, 0x10, 0x6A, 0x12, 0x03, 0x00, 0x10, 0x00, 0x00, 0x00, 0xF7}
	if _, err := ParseDataSet(b); err != ErrBadChecksum {
		t.Errorf("err = %v, want ErrBadChecksum", err)
	}
}

func TestParseRolandRejectsForeign(t *testing.T) {
	// GS data set (model ID 0x42) must not be accepted.
	b := []byte{0xF0, 0x41, 0x10, 0x42, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF7}
	if _, err := ParseRoland(b); err != ErrNotRolandJV {
		t.Errorf("err = %v, want ErrNotRolandJV", err)
	}
}

func TestIdentityRequestReply(t *testing.T) {
	req := IdentityRequest{DeviceID: 0x10}
	want := []byte{0xF0, 0x7E, 0x10, 0x06, 0x01, 0xF7}
	if got := req.Bytes(); !bytes.Equal(got, want) {
		t.Errorf("Bytes = % X, want % X", got, want)
	}

	// JV-1010 identity reply example from the MIDI implementation:
	// F0 7E 10 06 02 41 6A 00 05 00 00 00 00 00 00 F7
	raw := []byte{0xF0, 0x7E, 0x10, 0x06, 0x02, 0x41, 0x6A, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF7}
	reply, err := ParseIdentityReply(raw)
	if err != nil {
		t.Fatalf("ParseIdentityReply: %v", err)
	}
	if reply.Manufacturer != 0x41 || reply.FamilyCode != 0x006A || reply.FamilyNumber != 0x0005 {
		t.Errorf("reply = %+v", reply)
	}
	if got := reply.Bytes(); !bytes.Equal(got, raw) {
		t.Errorf("Bytes = % X, want % X", got, raw)
	}
}

func TestChannelMessagesRoundTrip(t *testing.T) {
	cc := ControlChange{Channel: 0, Controller: 0, Value: 0x50}
	if got := cc.Bytes(); !bytes.Equal(got, []byte{0xB0, 0x00, 0x50}) {
		t.Errorf("CC Bytes = % X", got)
	}
	cc2, err := ParseControlChange(cc.Bytes())
	if err != nil || cc2 != cc {
		t.Errorf("CC round trip = %+v, %v", cc2, err)
	}

	pc := ProgramChange{Channel: 9, Program: 40}
	if got := pc.Bytes(); !bytes.Equal(got, []byte{0xC9, 0x28}) {
		t.Errorf("PC Bytes = % X", got)
	}
	pc2, err := ParseProgramChange(pc.Bytes())
	if err != nil || pc2 != pc {
		t.Errorf("PC round trip = %+v, %v", pc2, err)
	}

	on := NoteOn{Channel: 0, Note: 60, Velocity: 100}
	if got := on.Bytes(); !bytes.Equal(got, []byte{0x90, 60, 100}) {
		t.Errorf("NoteOn Bytes = % X", got)
	}
	on2, err := ParseNoteOn(on.Bytes())
	if err != nil || on2 != on {
		t.Errorf("NoteOn round trip = %+v, %v", on2, err)
	}
}

func TestHexRoundTrip(t *testing.T) {
	b, err := ParseHex("F0 41 10 6A 11")
	if err != nil || !bytes.Equal(b, []byte{0xF0, 0x41, 0x10, 0x6A, 0x11}) {
		t.Errorf("ParseHex = % X, %v", b, err)
	}
	if got := Hex(b); got != "F0 41 10 6A 11" {
		t.Errorf("Hex = %q", got)
	}
	if _, err := ParseHex("F0 4"); err == nil {
		t.Error("ParseHex with odd digits should fail")
	}
	if _, err := ParseHex("F0 XY"); err == nil {
		t.Error("ParseHex with invalid digits should fail")
	}
}
