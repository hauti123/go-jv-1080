package gui

import (
	"testing"

	"go-jv-1080/internal/jv1080"
	"go-jv-1080/internal/midi/message"
)

// buildDT1 builds a raw DT1 message for the given address and data.
func buildDT1(devID uint8, addr message.Address, data []byte) []byte {
	ds := message.DataSet{DeviceID: devID, Address: addr, Data: data}
	return ds.Bytes()
}

func TestHandleDeviceMessageCommon(t *testing.T) {
	c := NewJVController(nil, jv1080.BankUser, 5)
	c.RegisterCommonParams(
		PatchCommonParam("level", 0x2E),
		PatchCommonParam("pan", 0x2F),
	)
	var got []Param
	var vals [][]byte
	c.Update(func(p Param, v []byte) {
		got = append(got, p)
		vals = append(vals, v)
	})

	// unsolicited DT1 writing level (0x2E) of user patch 5:
	// area 0x11, second byte = number-1 = 4
	addr := message.Address{0x11, 0x04, 0x00, 0x2E}
	c.HandleDeviceMessage(buildDT1(0x10, addr, []byte{100}))
	if len(got) != 1 || got[0].Name != "level" || vals[0][0] != 100 {
		t.Fatalf("got %v %v, want [level] [100]", got, vals)
	}

	// a write to an unregistered offset is ignored
	c.HandleDeviceMessage(buildDT1(0x10, message.Address{0x11, 0x04, 0x00, 0x0C}, []byte{1}))
	if len(got) != 1 {
		t.Fatalf("unregistered offset produced update: %v", got)
	}

	// a write to another patch is ignored
	c.HandleDeviceMessage(buildDT1(0x10, message.Address{0x11, 0x05, 0x00, 0x2E}, []byte{1}))
	if len(got) != 1 {
		t.Fatalf("other patch produced update: %v", got)
	}
}

func TestHandleDeviceMessageTemporaryPatch(t *testing.T) {
	c := NewJVController(nil, jv1080.BankPresetA, 1)
	c.RegisterCommonParams(PatchCommonParam("octave", 0x41))
	var got []Param
	c.Update(func(p Param, v []byte) { got = append(got, p) })

	// preset banks live in the patch-mode temporary patch (area 0x03)
	c.HandleDeviceMessage(buildDT1(0x10, message.Address{0x03, 0x00, 0x00, 0x41}, []byte{4}))
	if len(got) != 1 || got[0].Name != "octave" {
		t.Fatalf("got %v, want [octave]", got)
	}
}

func TestHandleDeviceMessageIgnoresNonDataSet(t *testing.T) {
	c := NewJVController(nil, jv1080.BankUser, 1)
	c.RegisterCommonParams(PatchCommonParam("level", 0x2E))
	called := false
	c.Update(func(p Param, v []byte) { called = true })
	c.HandleDeviceMessage([]byte{0x90, 0x3C, 0x40}) // note on
	if called {
		t.Fatal("note on triggered an update")
	}
}

func TestControllerAddressResolution(t *testing.T) {
	c := NewJVController(nil, jv1080.BankUser, 10)
	addr, err := c.address(PatchCommonParam("level", 0x2E))
	if err != nil {
		t.Fatal(err)
	}
	want := message.Address{0x11, 9, 0x00, 0x2E}
	if addr != want {
		t.Errorf("address = %v, want %v", addr, want)
	}
}
