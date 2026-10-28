package jv1080

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"go-jv-1080/internal/midi"
	"go-jv-1080/internal/midi/message"
)

// The integration tests talk to a real, connected JV-1080. They are skipped
// unless JV1080_INTEGRATION=1 is set. The MIDI ports are selected with the
// environment:
//
//	JV1080_OUT / JV1080_IN         port names (take precedence)
//	JV1080_OUT_NUM / JV1080_IN_NUM port numbers (default 0)
//	JV1080_DEVICE_ID               device ID in hex (default 10 = unit 17)
//	JV1080_TEST_PATCH              user patch slot used by the store test
//	                                   (default 128, restored afterwards)
func skipUnlessIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("JV1080_INTEGRATION") != "1" {
		t.Skip("integration test: set JV1080_INTEGRATION=1 and connect a JV-1080 " +
			"(ports via JV1080_OUT/JV1080_IN or JV1080_OUT_NUM/JV1080_IN_NUM)")
	}
}

func envInt(name string, def int) int {
	if s := os.Getenv(name); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return def
}

func envHex(name string, def uint8) uint8 {
	if s := os.Getenv(name); s != "" {
		if n, err := strconv.ParseUint(s, 16, 8); err == nil {
			return uint8(n)
		}
	}
	return def
}

func newIntegrationClient(t *testing.T) *Client {
	t.Helper()
	dev, err := midi.OpenDevice(
		os.Getenv("JV1080_OUT"), os.Getenv("JV1080_IN"),
		envInt("JV1080_OUT_NUM", 0), envInt("JV1080_IN_NUM", 0),
	)
	if err != nil {
		t.Fatalf("cannot open MIDI device: %v", err)
	}
	t.Cleanup(func() { dev.Close() })
	return NewClient(dev, envHex("JV1080_DEVICE_ID", message.DefaultDeviceID))
}

// TestIntegrationIdentity checks that the connected device answers an
// identity request.
func TestIntegrationIdentity(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	reply, err := c.Identity()
	if err != nil {
		t.Fatalf("no identity reply: %v", err)
	}
	if reply.Manufacturer != message.RolandManufacturerID {
		t.Errorf("manufacturer = %02X, want Roland %02X", reply.Manufacturer, message.RolandManufacturerID)
	}
	t.Logf("identity reply: %s", reply)
}

// TestIntegrationSelectAndReadPatch selects a factory patch and reads it
// back through the temporary patch, checking that the request/response
// round trip returns plausible data.
func TestIntegrationSelectAndReadPatch(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	for _, bank := range []Bank{BankPresetA, BankUser} {
		p, err := c.ReadPatch(bank, 1)
		if err != nil {
			t.Fatalf("ReadPatch(%s, 1): %v", bank, err)
		}
		if len(p.Common) != PatchCommonSize {
			t.Errorf("%s: common size = %d, want %d", bank, len(p.Common), PatchCommonSize)
		}
		for i, tone := range p.Tones {
			if len(tone) != PatchToneSize {
				t.Errorf("%s: tone %d size = %d, want %d", bank, i+1, len(tone), PatchToneSize)
			}
		}
		name := p.Name()
		for i, ch := range []byte(name) {
			if ch < 32 || ch > 127 {
				t.Errorf("%s: name byte %d = %02X, not printable ASCII", bank, i, ch)
			}
		}
		if bank != BankUser {
			if want, ok := PresetPatchName(bank, 1); ok && want != name {
				t.Errorf("%s patch 1 name = %q, factory listing says %q", bank, name, want)
			}
		}
		t.Logf("%s patch 1: %q level=%d pan=%d", bank, name, p.PatchLevel(), p.PatchPan())
	}
}

// TestIntegrationReadUserBank reads a single user patch and a single user
// performance directly from the user memory areas.
func TestIntegrationReadUserBank(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	p, err := c.ReadUserPatch(1)
	if err != nil {
		t.Fatalf("ReadUserPatch(1): %v", err)
	}
	if p.Area != userPatchArea(1) {
		t.Errorf("patch area = %s, want %s", p.Area, userPatchArea(1))
	}
	perf, err := c.ReadUserPerformance(1)
	if err != nil {
		t.Fatalf("ReadUserPerformance(1): %v", err)
	}
	if perf.Area != userPerformanceArea(1) {
		t.Errorf("performance area = %s, want %s", perf.Area, userPerformanceArea(1))
	}
	t.Logf("user patch 1: %q; user performance 1: %q", p.Name(), perf.Name())
}

// TestIntegrationPatchParameterRoundTrip modifies a parameter of a user
// patch via SetPatchParameter, reads the patch back and verifies the
// change; afterwards the original value is restored.
func TestIntegrationPatchParameterRoundTrip(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	slot := envInt("JV1080_TEST_PATCH", 128)

	p, err := c.ReadUserPatch(slot)
	if err != nil {
		t.Fatalf("read user patch %d: %v", slot, err)
	}

	// Patch level, common block offset 0x2E.
	const levelOffset = 0x2E
	origLevel := p.Common[levelOffset]
	newLevel := (origLevel + 37) % 128
	if err := c.SetPatchParameter(BankUser, slot, "common", 0, levelOffset, []byte{newLevel}); err != nil {
		t.Fatalf("SetPatchParameter level: %v", err)
	}
	back, err := c.ReadUserPatch(slot)
	if err != nil {
		t.Fatalf("read back user patch %d: %v", slot, err)
	}
	if got := back.Common[levelOffset]; got != newLevel {
		t.Errorf("patch level = %d, want %d", got, newLevel)
	}

	// Restore the original value.
	if err := c.SetPatchParameter(BankUser, slot, "common", 0, levelOffset, []byte{origLevel}); err != nil {
		t.Fatalf("restore patch level: %v", err)
	}
	restored, err := c.ReadUserPatch(slot)
	if err != nil {
		t.Fatalf("read back restored user patch %d: %v", slot, err)
	}
	if got := restored.Common[levelOffset]; got != origLevel {
		t.Errorf("restored patch level = %d, want %d", got, origLevel)
	}
}

// TestIntegrationStorePatch stores the temporary patch into a user slot,
// verifies it and restores the original patch content.
func TestIntegrationStorePatch(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	slot := envInt("JV1080_TEST_PATCH", 128)

	orig, err := c.ReadUserPatch(slot)
	if err != nil {
		t.Fatalf("read original user patch %d: %v", slot, err)
	}

	tp, err := c.ReadTemporaryPatch()
	if err != nil {
		t.Fatalf("read temporary patch: %v", err)
	}
	tp.SetName("INTTEST")
	tp.SetPatchLevel(111)
	if err := c.StorePatch(slot, tp); err != nil {
		t.Fatalf("StorePatch(%d): %v", slot, err)
	}

	back, err := c.ReadUserPatch(slot)
	if err != nil {
		t.Fatalf("read back stored user patch %d: %v", slot, err)
	}
	if back.Name() != "INTTEST" {
		t.Errorf("stored patch name = %q, want %q", back.Name(), "INTTEST")
	}
	if back.PatchLevel() != 111 {
		t.Errorf("stored patch level = %d, want 111", back.PatchLevel())
	}
	if !bytes.Equal(back.Common, tp.Common) {
		t.Error("stored common block differs from temporary patch")
	}

	// Restore the original patch.
	if err := c.StorePatch(slot, orig); err != nil {
		t.Fatalf("restore user patch %d: %v", slot, err)
	}
	restored, err := c.ReadUserPatch(slot)
	if err != nil {
		t.Fatalf("read back restored user patch %d: %v", slot, err)
	}
	if !bytes.Equal(restored.Common, orig.Common) {
		t.Error("restored common block differs from original")
	}
	for tone := 0; tone < 4; tone++ {
		if !bytes.Equal(restored.Tones[tone], orig.Tones[tone]) {
			t.Errorf("restored tone %d differs from original", tone+1)
		}
	}
}
