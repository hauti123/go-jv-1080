package jv1080

import (
	"bytes"
	"fmt"
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
//	                                   (default 113, restored afterwards)
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
// identity request. Some JV-1080 units do not implement the identity
// reply; in that case the test falls back to a data request as a liveness
// check.
func TestIntegrationIdentity(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	reply, err := c.Identity()
	if err != nil {
		t.Logf("no identity reply (%v); checking device liveness with a data request instead", err)
		if _, err := c.ReadUserPatch(1); err != nil {
			t.Fatalf("device does not answer data requests either: %v", err)
		}
		return
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
	slot := envInt("JV1080_TEST_PATCH", 113)

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

// patchParameterTest describes one patch parameter: where it lives
// (block, tone, offset), how many raw bytes it occupies and how to turn
// its raw bytes into a value for comparison.
type patchParameterTest struct {
	name   string
	block  string // "common" or "tone"
	tone   int    // 1-4, 0 for common
	offset int
	size   int // raw bytes: 1, or 2 for nibblized values
	read   func(*Patch) int
}

// patchParameterTests lists every patch parameter the Patch model exposes.
func patchParameterTests() []patchParameterTest {
	common := []patchParameterTest{
		{"name[0]", "common", 0, 0x00, 1, func(p *Patch) int { return int(p.Common[0x00]) }},
		{"name[11]", "common", 0, 0x0B, 1, func(p *Patch) int { return int(p.Common[0x0B]) }},
		{"EFX type", "common", 0, 0x0C, 1, func(p *Patch) int { return p.EFXType() }},
		{"patch level", "common", 0, 0x2E, 1, func(p *Patch) int { return p.PatchLevel() }},
		{"patch pan", "common", 0, 0x2F, 1, func(p *Patch) int { return p.PatchPan() }},
	}
	var tones []patchParameterTest
	for t := 1; t <= 4; t++ {
		tones = append(tones,
			patchParameterTest{fmt.Sprintf("tone %d switch", t), "tone", t, 0x00, 1,
				func(p *Patch) int {
					b := 0
					if p.ToneSwitch(t) {
						b = 1
					}
					return b
				}},
			patchParameterTest{fmt.Sprintf("tone %d wave group type", t), "tone", t, 0x01, 1,
				func(p *Patch) int { return p.WaveGroupType(t) }},
			patchParameterTest{fmt.Sprintf("tone %d wave group ID", t), "tone", t, 0x02, 1,
				func(p *Patch) int { return p.WaveGroupID(t) }},
			patchParameterTest{fmt.Sprintf("tone %d wave number", t), "tone", t, 0x03, 2,
				func(p *Patch) int { return p.WaveNumber(t) }},
			patchParameterTest{fmt.Sprintf("tone %d coarse tune", t), "tone", t, 0x3D, 1,
				func(p *Patch) int { return p.CoarseTune(t) + 48 }},
			patchParameterTest{fmt.Sprintf("tone %d fine tune", t), "tone", t, 0x3E, 1,
				func(p *Patch) int { return p.FineTune(t) + 50 }},
			patchParameterTest{fmt.Sprintf("tone %d cutoff", t), "tone", t, 0x51, 1,
				func(p *Patch) int { return p.Cutoff(t) }},
			patchParameterTest{fmt.Sprintf("tone %d resonance", t), "tone", t, 0x53, 1,
				func(p *Patch) int { return p.Resonance(t) }},
			patchParameterTest{fmt.Sprintf("tone %d level", t), "tone", t, 0x65, 1,
				func(p *Patch) int { return p.ToneLevel(t) }},
			patchParameterTest{fmt.Sprintf("tone %d pan", t), "tone", t, 0x77, 1,
				func(p *Patch) int { return p.TonePan(t) }},
			patchParameterTest{fmt.Sprintf("tone %d chorus send", t), "tone", t, 0x7F, 1,
				func(p *Patch) int { return p.ChorusSend(t) }},
			patchParameterTest{fmt.Sprintf("tone %d reverb send", t), "tone", t, 0x80, 1,
				func(p *Patch) int { return p.ReverbSend(t) }},
		)
	}
	return append(common, tones...)
}

// newParameterValue derives a new raw value for a parameter from its
// current raw value, so the test never writes the value that is already
// there. For nibblized (2-byte) parameters the raw value is a 0-254 number
// split into a high nibble byte and a low nibble byte (high nibble first);
// the shift is applied to that number and re-nibblized, so both bytes
// stay valid nibbles.
func newParameterValue(current []byte, size int) []byte {
	if size == 2 {
		raw := int(current[0])<<4 | int(current[1])
		raw = (raw + 37) % 255
		return []byte{byte(raw >> 4), byte(raw & 0x0F)}
	}
	out := make([]byte, size)
	for i := range out {
		out[i] = (current[i] + 37) % 128
	}
	return out
}

// TestPatchParameterTable checks the parameter table itself: every entry
// must address a valid offset within its block, sizes must be 1 or 2, and
// the raw value transformation must actually change the value.
func TestPatchParameterTable(t *testing.T) {
	tests := patchParameterTests()
	if len(tests) != 5+4*12 {
		t.Errorf("table has %d entries, want %d", len(tests), 5+4*12)
	}
	for _, pt := range tests {
		blockSize := PatchCommonSize
		if pt.block == "tone" {
			if pt.tone < 1 || pt.tone > 4 {
				t.Errorf("%s: tone %d out of range", pt.name, pt.tone)
			}
			blockSize = PatchToneSize
		}
		if pt.offset < 0 || pt.offset+pt.size > blockSize {
			t.Errorf("%s: offset %d size %d out of block bounds %d", pt.name, pt.offset, pt.size, blockSize)
		}
		if pt.size != 1 && pt.size != 2 {
			t.Errorf("%s: size %d, want 1 or 2", pt.name, pt.size)
		}
		current := make([]byte, pt.size)
		if got := newParameterValue(current, pt.size); bytes.Equal(got, current) {
			t.Errorf("%s: newParameterValue did not change the value", pt.name)
		}
	}
}

// TestIntegrationPatchParameters is a table test that writes every patch
// parameter of the temporary patch through SetPatchParameter, reads the
// whole patch back and verifies that the parameter changed and that no
// other byte of the patch was touched.
func TestIntegrationPatchParameters(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)

	before, err := c.ReadTemporaryPatch()
	if err != nil {
		t.Fatalf("read temporary patch: %v", err)
	}

	for _, pt := range patchParameterTests() {
		t.Run(pt.name, func(t *testing.T) {
			// Current raw value.
			var block []byte
			if pt.block == "common" {
				block = before.Common
			} else {
				block = before.Tones[pt.tone-1]
			}
			current := append([]byte{}, block[pt.offset:pt.offset+pt.size]...)
			value := newParameterValue(current, pt.size)

			// Non-user banks are written to the temporary patch.
			if err := c.SetPatchParameter(BankPresetA, 0, pt.block, pt.tone, pt.offset, value); err != nil {
				t.Fatalf("SetPatchParameter: %v", err)
			}

			after, err := c.ReadTemporaryPatch()
			if err != nil {
				t.Fatalf("read back temporary patch: %v", err)
			}

			// The parameter must hold the new value.
			if got := pt.read(after); got != pt.read(mutateCopy(before, pt, value)) {
				t.Errorf("parameter value = %d, want %d", got, pt.read(mutateCopy(before, pt, value)))
			}

			// No other byte may have changed: compare against the
			// expected patch with exactly this parameter modified.
			expected := clonePatch(before)
			var eblock []byte
			if pt.block == "common" {
				eblock = expected.Common
			} else {
				eblock = expected.Tones[pt.tone-1]
			}
			copy(eblock[pt.offset:pt.offset+pt.size], value)
			if !bytes.Equal(after.Common, expected.Common) {
				t.Error("common block changed beyond the parameter")
			}
			for i := range after.Tones {
				if !bytes.Equal(after.Tones[i], expected.Tones[i]) {
					t.Errorf("tone %d changed beyond the parameter", i+1)
				}
			}

			// Restore the original value.
			if err := c.SetPatchParameter(BankPresetA, 0, pt.block, pt.tone, pt.offset, current); err != nil {
				t.Fatalf("restore parameter: %v", err)
			}
		})
	}
}

// mutateCopy returns a copy of p with the parameter applied, for computing
// the expected value.
func mutateCopy(p *Patch, pt patchParameterTest, value []byte) *Patch {
	q := clonePatch(p)
	var block []byte
	if pt.block == "common" {
		block = q.Common
	} else {
		block = q.Tones[pt.tone-1]
	}
	copy(block[pt.offset:pt.offset+pt.size], value)
	return q
}

func clonePatch(p *Patch) *Patch {
	q := NewPatch(p.Area)
	copy(q.Common, p.Common)
	for i := range p.Tones {
		copy(q.Tones[i], p.Tones[i])
	}
	return q
}

// TestIntegrationStorePatch stores the temporary patch into a user slot,
// verifies it and restores the original patch content.
func TestIntegrationStorePatch(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)
	slot := envInt("JV1080_TEST_PATCH", 113)

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
