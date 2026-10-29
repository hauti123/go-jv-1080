package sysex

import (
	"bytes"
	"os"
	"strconv"
	"testing"
	"time"

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
//	                                   (default 113, restored afterwards)
func integrationEnabled() bool {
	return os.Getenv("JV1080_INTEGRATION") == "1"
}

func skipUnlessIntegration(t *testing.T) {
	t.Helper()
	if !integrationEnabled() {
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

func openIntegrationDevice(t *testing.T) *midi.Device {
	t.Helper()
	dev, err := midi.OpenDevice(
		os.Getenv("JV1080_OUT"), os.Getenv("JV1080_IN"),
		envInt("JV1080_OUT_NUM", 0), envInt("JV1080_IN_NUM", 0),
	)
	if err != nil {
		t.Fatalf("cannot open MIDI device: %v", err)
	}
	t.Cleanup(func() { dev.Close() })
	return dev
}

func newIntegrationClient(t *testing.T) *Client {
	t.Helper()
	dev := openIntegrationDevice(t)
	return NewClient(dev, envHex("JV1080_DEVICE_ID", message.DefaultDeviceID))
}

// TestIntegrationIdentity sends an identity request and checks that a
// Roland JV answers.
func TestIntegrationIdentity(t *testing.T) {
	skipUnlessIntegration(t)
	dev := openIntegrationDevice(t)
	deviceID := envHex("JV1080_DEVICE_ID", message.DefaultDeviceID)
	reply, err := dev.RequestIdentity(deviceID, 5*time.Second)
	if err != nil {
		t.Fatalf("no identity reply: %v", err)
	}
	if reply.Manufacturer != message.RolandManufacturerID {
		t.Errorf("manufacturer = %02X, want Roland %02X", reply.Manufacturer, message.RolandManufacturerID)
	}
	if reply.FamilyCode != message.JV1080ModelID {
		t.Errorf("family code = %04X, want JV model %04X", reply.FamilyCode, message.JV1080ModelID)
	}
	t.Logf("identity reply: %s", reply)
}

// TestIntegrationDownloadTemporaryPatch requests the temporary patch and
// checks that a complete, plausible patch comes back.
func TestIntegrationDownloadTemporaryPatch(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	p, err := c.DownloadTemporaryPatch()
	if err != nil {
		t.Fatalf("download temporary patch: %v", err)
	}
	if len(p.Common) != PatchCommonSize {
		t.Errorf("common size = %d, want %d", len(p.Common), PatchCommonSize)
	}
	for i, tone := range p.Tones {
		if len(tone) != PatchToneSize {
			t.Errorf("tone %d size = %d, want %d", i+1, len(tone), PatchToneSize)
		}
	}
	for i, ch := range p.Common[0:12] {
		if ch < 32 || ch > 127 {
			t.Errorf("patch name byte %d = %02X, not printable ASCII", i, ch)
		}
	}
	t.Logf("temporary patch name: %q", string(bytes.TrimRight(p.Common[0:12], " ")))
}

// TestIntegrationSetPatchParameter modifies single parameters of the
// temporary patch and reads them back to verify the change took effect.
func TestIntegrationSetPatchParameter(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	area := TemporaryPatchArea()

	p, err := c.DownloadPatch(area)
	if err != nil {
		t.Fatalf("download temporary patch: %v", err)
	}

	// Patch level, common block offset 0x2E.
	const levelOffset = 0x2E
	newLevel := (p.Common[levelOffset] + 37) % 128
	if err := c.SetPatchParameter(area, "common", 0, levelOffset, []byte{newLevel}); err != nil {
		t.Fatalf("set patch level: %v", err)
	}

	// Tone 1 switch, tone block offset 0x00.
	newSwitch := byte(1)
	if p.Tones[0][0x00] != 0 {
		newSwitch = 0
	}
	if err := c.SetPatchParameter(area, "tone", 1, 0x00, []byte{newSwitch}); err != nil {
		t.Fatalf("set tone 1 switch: %v", err)
	}

	// Tone 4 cutoff, tone block offset 0x51.
	const cutoffOffset = 0x51
	newCutoff := (p.Tones[3][cutoffOffset] + 53) % 128
	if err := c.SetPatchParameter(area, "tone", 4, cutoffOffset, []byte{newCutoff}); err != nil {
		t.Fatalf("set tone 4 cutoff: %v", err)
	}

	back, err := c.DownloadPatch(area)
	if err != nil {
		t.Fatalf("read back temporary patch: %v", err)
	}
	if got := back.Common[levelOffset]; got != newLevel {
		t.Errorf("patch level = %d, want %d", got, newLevel)
	}
	if got := back.Tones[0][0x00]; got != newSwitch {
		t.Errorf("tone 1 switch = %d, want %d", got, newSwitch)
	}
	if got := back.Tones[3][cutoffOffset]; got != newCutoff {
		t.Errorf("tone 4 cutoff = %d, want %d", got, newCutoff)
	}
}

// TestIntegrationSetPerformanceParameter modifies single parameters of the
// temporary performance and reads them back.
func TestIntegrationSetPerformanceParameter(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	area := TemporaryPerformanceArea()

	p, err := c.DownloadPerformance(area)
	if err != nil {
		t.Fatalf("download temporary performance: %v", err)
	}

	// Part 1 channel, part block offset 0x01.
	const channelOffset = 0x01
	newChannel := (p.Parts[0][channelOffset] + 5) % 17
	if err := c.SetPerformanceParameter(area, "part", 1, channelOffset, []byte{newChannel}); err != nil {
		t.Fatalf("set part 1 channel: %v", err)
	}

	// Part 2 level, part block offset 0x06.
	const levelOffset = 0x06
	newLevel := (p.Parts[1][levelOffset] + 41) % 128
	if err := c.SetPerformanceParameter(area, "part", 2, levelOffset, []byte{newLevel}); err != nil {
		t.Fatalf("set part 2 level: %v", err)
	}

	back, err := c.DownloadPerformance(area)
	if err != nil {
		t.Fatalf("read back temporary performance: %v", err)
	}
	if got := back.Parts[0][channelOffset]; got != newChannel {
		t.Errorf("part 1 channel = %d, want %d", got, newChannel)
	}
	if got := back.Parts[1][levelOffset]; got != newLevel {
		t.Errorf("part 2 level = %d, want %d", got, newLevel)
	}
}

// TestIntegrationStoreUserPatch writes the temporary patch into a user
// patch slot, verifies the stored data by reading it back and finally
// restores the original patch content.
func TestIntegrationStoreUserPatch(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	slot := envInt("JV1080_TEST_PATCH", 113)

	orig, err := c.DownloadUserPatch(slot)
	if err != nil {
		t.Fatalf("read original user patch %d: %v", slot, err)
	}

	// Build a distinctive temporary patch.
	tp, err := c.DownloadTemporaryPatch()
	if err != nil {
		t.Fatalf("download temporary patch: %v", err)
	}
	copy(tp.Common[0:12], "INTTEST     ")
	tp.Common[0x2E] = 111 // patch level

	// Write it to the temporary area, then store it into the user slot.
	if err := c.WriteBlock(TemporaryPatchArea(), tp.Common); err != nil {
		t.Fatalf("write temp common: %v", err)
	}
	for tone := 1; tone <= PatchToneCount; tone++ {
		addr, _ := PatchToneBlock(TemporaryPatchArea(), tone)
		if err := c.WriteBlock(addr, tp.Tones[tone-1]); err != nil {
			t.Fatalf("write temp tone %d: %v", tone, err)
		}
	}
	if err := c.WriteBlock(UserPatchArea(slot), tp.Common); err != nil {
		t.Fatalf("store common: %v", err)
	}
	for tone := 1; tone <= PatchToneCount; tone++ {
		addr, _ := PatchToneBlock(UserPatchArea(slot), tone)
		if err := c.WriteBlock(addr, tp.Tones[tone-1]); err != nil {
			t.Fatalf("store tone %d: %v", tone, err)
		}
	}

	back, err := c.DownloadUserPatch(slot)
	if err != nil {
		t.Fatalf("read back stored user patch %d: %v", slot, err)
	}
	if !bytes.Equal(back.Common, tp.Common) {
		t.Errorf("stored common differs from written common (name %q, level %d)",
			string(bytes.TrimRight(back.Common[0:12], " ")), back.Common[0x2E])
	}
	for tone := 0; tone < PatchToneCount; tone++ {
		if !bytes.Equal(back.Tones[tone], tp.Tones[tone]) {
			t.Errorf("stored tone %d differs", tone+1)
		}
	}

	// Restore the original patch.
	if err := c.WriteBlock(UserPatchArea(slot), orig.Common); err != nil {
		t.Fatalf("restore common: %v", err)
	}
	for tone := 1; tone <= PatchToneCount; tone++ {
		addr, _ := PatchToneBlock(UserPatchArea(slot), tone)
		if err := c.WriteBlock(addr, orig.Tones[tone-1]); err != nil {
			t.Fatalf("restore tone %d: %v", tone, err)
		}
	}
	restored, err := c.DownloadUserPatch(slot)
	if err != nil {
		t.Fatalf("read back restored user patch %d: %v", slot, err)
	}
	if !bytes.Equal(restored.Common, orig.Common) {
		t.Error("restored common differs from original")
	}
	for tone := 0; tone < PatchToneCount; tone++ {
		if !bytes.Equal(restored.Tones[tone], orig.Tones[tone]) {
			t.Errorf("restored tone %d differs from original", tone+1)
		}
	}
}
