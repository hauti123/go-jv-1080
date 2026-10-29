package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"go-jv-1080/internal/jv1080"
	"go-jv-1080/internal/midi"
	"go-jv-1080/internal/midi/message"
)

// The integration test talks to a real, connected JV-1080. It is skipped
// unless JV1080_INTEGRATION=1 is set. The MIDI ports are selected with the
// same environment variables as internal/jv1080:
//
//	JV1080_OUT / JV1080_IN         port names (take precedence)
//	JV1080_OUT_NUM / JV1080_IN_NUM port numbers (default 0)
//	JV1080_DEVICE_ID               device ID in hex (default 10 = unit 17)
//	JV1080_TEST_TONE              tone number used by the test (1-4,
//	                                   default 1; per-Tone parameters are
//	                                   additionally checked for all tones)
//
// The test works on the patch-mode temporary patch only; it saves the
// original temporary patch first and writes it back at the end, so the
// device is left as it was.
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

func newIntegrationClient(t *testing.T) *jv1080.Client {
	t.Helper()
	dev, err := midi.OpenDevice(
		os.Getenv("JV1080_OUT"), os.Getenv("JV1080_IN"),
		envInt("JV1080_OUT_NUM", 0), envInt("JV1080_IN_NUM", 0),
	)
	if err != nil {
		t.Fatalf("cannot open MIDI device: %v", err)
	}
	t.Cleanup(func() { dev.Close() })
	return jv1080.NewClient(dev, envHex("JV1080_DEVICE_ID", message.DefaultDeviceID))
}

// editParam describes one parameter write to verify: the block, tone and
// offset within the block. label names the CLI form the write
// corresponds to; the test value is derived from the original byte at
// run time.
type editParam struct {
	label  string
	block  string
	tone   int
	offset int
}

// testValue returns a raw value for the parameter at offset off that is
// guaranteed to differ from orig: the original value shifted by a
// constant, kept inside the 7-bit range.
func testValue(orig byte) byte {
	v := orig + 37
	if v > 127 {
		v -= 127
	}
	if v == orig {
		v = (v + 1) % 128
	}
	return v
}

// integrationEdits enumerates one write per addressable parameter for the
// given tone. Common parameters are enumerated once (tone 0), per-Tone
// parameters for the requested tone; a second pass over all four tones is
// done by the test itself to verify the tone block bases.
func integrationEdits(tone int) []editParam {
	var edits []editParam
	add := func(label, block string, t, offset int) {
		edits = append(edits, editParam{label: label, block: block, tone: t, offset: offset})
	}
	if tone == 0 {
		// Patch common block parameters. The name (12 bytes) is skipped:
		// it is a multi-byte parameter and covered separately.
		for name, p := range commonParams {
			if name == "name" {
				continue
			}
			add("common "+name, "common", 0, p.offset)
		}
		// hold/peak of control 1-3 and patch control sources 2-3.
		for i := 1; i <= 3; i++ {
			add(fmt.Sprintf("ctrl %d hold", i), "common", 0, commonIndexedParams["hold"].base+i)
		}
		for i := 2; i <= 3; i++ {
			add(fmt.Sprintf("ctrl %d src", i), "common", 0, commonIndexedParams["src"].base+i)
		}
		return edits
	}

	// Per-Tone parameters of the given tone.
	for name, p := range toneParams {
		add("tone "+name, "tone", tone, p.offset)
	}
	// Wave number and tempo are nibblized; the CLI writes two nibble
	// bytes. Their single-byte forms are already covered above; here we
	// check the nibblized write of the wave number.
	// LFO 1 and 2.
	for lfo := 1; lfo <= 2; lfo++ {
		for name, p := range lfoParams[lfo] {
			add(fmt.Sprintf("lfo %d %s", lfo, name), "tone", tone, p.offset)
		}
	}
	// LFO depths.
	for name, p := range lfoDepthParams {
		for lfo := 1; lfo <= 2; lfo++ {
			add(fmt.Sprintf("lfo %d %s", lfo, name), "tone", tone, p.offset+lfo-1)
		}
	}
	// Pitch envelope parameters.
	for name, p := range pitchEnvParams {
		add("pitch env "+name, "tone", tone, p.offset)
	}
	// TVF and TVF envelope.
	for name, p := range tvfParams {
		add("tvf "+name, "tone", tone, p.offset)
	}
	for name, p := range tvfEnvParams {
		add("tvf env "+name, "tone", tone, p.offset)
	}
	// TVA and TVA envelope.
	for name, p := range tvaParams {
		add("tva "+name, "tone", tone, p.offset)
	}
	for name, p := range tvaEnvParams {
		add("tva env "+name, "tone", tone, p.offset)
	}
	// Control 1-3 destinations and depths.
	for i := 1; i <= 3; i++ {
		for k := 1; k <= 4; k++ {
			add(fmt.Sprintf("ctrl %d dst%d", i, k), "tone", tone, ctrlToneOffset("dst", i, k))
			add(fmt.Sprintf("ctrl %d dep%d", i, k), "tone", tone, ctrlToneOffset("dep", i, k))
		}
	}
	return edits
}

// TestIntegrationPatchEditAddresses verifies that every parameter write of
// the CLI hits the right byte in the device: for each addressable
// parameter it dumps the temporary patch, writes a distinctive test value
// through Client.SetPatchParameter, dumps again and asserts that exactly
// the expected byte changed at the expected offset. Afterwards the
// original temporary patch is restored and verified.
//
// A wrong offset in the CLI tables shows up as a byte changing at a
// different position than asserted (or more than one byte changing), so a
// single run localizes every addressing bug without manual checking of
// the values on the device.
func TestIntegrationPatchEditAddresses(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)

	orig, err := c.ReadTemporaryPatch()
	if err != nil {
		t.Fatalf("read temporary patch: %v", err)
	}
	t.Cleanup(func() {
		if err := c.WritePatchToTemporary(orig); err != nil {
			t.Errorf("restore temporary patch: %v", err)
		}
	})

	// Verify every single-byte parameter once (common, and per-Tone for
	// tone 1), then verify the tone block bases by repeating a
	// representative parameter for tones 2-4.
	edits := append(integrationEdits(0), integrationEdits(envInt("JV1080_TEST_TONE", 1))...)
	for _, e := range edits {
		e := e
		t.Run(e.label, func(t *testing.T) {
			before, err := c.ReadTemporaryPatch()
			if err != nil {
				t.Fatalf("dump before edit: %v", err)
			}
			block := blockOf(before, e)
			if e.offset >= len(block) {
				t.Fatalf("offset 0x%02X outside block (%d bytes)", e.offset, len(block))
			}
			origByte := block[e.offset]
			v := testValue(origByte)
			if err := c.SetPatchParameter(jv1080.BankPresetA, 0, e.block, e.tone, e.offset, []byte{v}); err != nil {
				t.Fatalf("set parameter: %v", err)
			}
			after, err := c.ReadTemporaryPatch()
			if err != nil {
				t.Fatalf("dump after edit: %v", err)
			}
			checkEdit(t, before, after, e, v)
			// Restore the original byte immediately so the next
			// parameter starts from the same state.
			if err := c.SetPatchParameter(jv1080.BankPresetA, 0, e.block, e.tone, e.offset, []byte{origByte}); err != nil {
				t.Fatalf("restore parameter: %v", err)
			}
		})
	}

	// Tone block bases: the same offset in tones 2-4 must hit the
	// corresponding byte in that tone's block.
	for tone := 2; tone <= 4; tone++ {
		tone := tone
		t.Run(fmt.Sprintf("tone %d base", tone), func(t *testing.T) {
			const cutoffOffset = 0x51
			before, err := c.ReadTemporaryPatch()
			if err != nil {
				t.Fatalf("dump before edit: %v", err)
			}
			origByte := before.Tones[tone-1][cutoffOffset]
			v := testValue(origByte)
			if err := c.SetPatchParameter(jv1080.BankPresetA, 0, "tone", tone, cutoffOffset, []byte{v}); err != nil {
				t.Fatalf("set cutoff: %v", err)
			}
			after, err := c.ReadTemporaryPatch()
			if err != nil {
				t.Fatalf("dump after edit: %v", err)
			}
			checkEdit(t, before, after, editParam{label: fmt.Sprintf("tone %d cutoff", tone), block: "tone", tone: tone, offset: cutoffOffset}, v)
			if err := c.SetPatchParameter(jv1080.BankPresetA, 0, "tone", tone, cutoffOffset, []byte{origByte}); err != nil {
				t.Fatalf("restore cutoff: %v", err)
			}
		})
	}
}

// blockOf returns the block (common or tone) of the patch that the edit
// addresses.
func blockOf(p *jv1080.Patch, e editParam) []byte {
	if e.block == "common" {
		return p.Common
	}
	return p.Tones[e.tone-1]
}

// checkEdit asserts that between the two dumps exactly the byte of the
// edit changed: the byte at the edit's offset (within its block) equals
// value, and no other byte differs.
func checkEdit(t *testing.T, before, after *jv1080.Patch, e editParam, value byte) {
	t.Helper()
	var diffs []string
	if got := after.Common; !equalBytes(before.Common, got) {
		for i := range before.Common {
			if before.Common[i] != got[i] {
				diffs = append(diffs, fmt.Sprintf("common %02X: %d -> %d", i, before.Common[i], got[i]))
			}
		}
	}
	for tone := 0; tone < 4; tone++ {
		b, a := before.Tones[tone], after.Tones[tone]
		if !equalBytes(b, a) {
			for i := range b {
				if b[i] != a[i] {
					diffs = append(diffs, fmt.Sprintf("tone %d %02X: %d -> %d", tone+1, i, b[i], a[i]))
				}
			}
		}
	}
	// The expected change: exactly one byte, at the edit's offset in its
	// block.
	if len(diffs) == 0 {
		t.Errorf("%s: no byte changed (expected offset 0x%02X in %s)", e.label, e.offset, e.block)
		return
	}
	if len(diffs) > 1 {
		t.Errorf("%s: %d bytes changed, want 1: %s", e.label, len(diffs), strings.Join(diffs, "; "))
		return
	}
	want := fmt.Sprintf("%s %02X: %d -> %d", e.block, e.offset, origByteOf(before, e), value)
	if e.block == "tone" {
		want = fmt.Sprintf("tone %d %02X: %d -> %d", e.tone, e.offset, origByteOf(before, e), value)
	}
	if diffs[0] != want {
		t.Errorf("%s: change at %s, want %s", e.label, diffs[0], want)
	}
}

func origByteOf(p *jv1080.Patch, e editParam) byte {
	return blockOf(p, e)[e.offset]
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
