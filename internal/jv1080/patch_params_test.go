package jv1080

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

// patchParameterTest describes one patch parameter: where it lives
// (block, tone, offset), how many raw bytes it occupies, the valid raw
// value range and how to read its value back from a patch.
type patchParameterTest struct {
	name   string
	block  string // "common" or "tone"
	tone   int    // 1-4, 0 for common
	offset int
	size   int // raw bytes: 1, or 2 for nibblized values
	min    int // smallest valid raw value
	max    int // largest valid raw value (for nibblized parameters: the raw 0-254 range)
	read   func(*Patch) int
}

// patchParameterTests lists every patch parameter the Patch model exposes.
func patchParameterTests() []patchParameterTest {
	common := []patchParameterTest{
		{"name[0]", "common", 0, 0x00, 1, 32, 127, func(p *Patch) int { return int(p.Common[0x00]) }},
		{"name[11]", "common", 0, 0x0B, 1, 32, 127, func(p *Patch) int { return int(p.Common[0x0B]) }},
		{"EFX type", "common", 0, 0x0C, 1, 0, 39, func(p *Patch) int { return p.EFXType() }},
		{"patch level", "common", 0, 0x2E, 1, 0, 127, func(p *Patch) int { return p.PatchLevel() }},
		{"patch pan", "common", 0, 0x2F, 1, 0, 127, func(p *Patch) int { return p.PatchPan() }},
	}
	var tones []patchParameterTest
	for t := 1; t <= 4; t++ {
		tones = append(tones,
			patchParameterTest{fmt.Sprintf("tone %d switch", t), "tone", t, 0x00, 1, 0, 1,
				func(p *Patch) int {
					b := 0
					if p.ToneSwitch(t) {
						b = 1
					}
					return b
				}},
			patchParameterTest{fmt.Sprintf("tone %d wave group type", t), "tone", t, 0x01, 1, 0, 2,
				func(p *Patch) int { return p.WaveGroupType(t) }},
			patchParameterTest{fmt.Sprintf("tone %d wave group ID", t), "tone", t, 0x02, 1, 0, 127,
				func(p *Patch) int { return p.WaveGroupID(t) }},
			patchParameterTest{fmt.Sprintf("tone %d wave number", t), "tone", t, 0x03, 2, 0, 254,
				func(p *Patch) int { return p.WaveNumber(t) }},
			patchParameterTest{fmt.Sprintf("tone %d coarse tune", t), "tone", t, 0x3D, 1, 0, 96,
				func(p *Patch) int { return p.CoarseTune(t) + 48 }},
			patchParameterTest{fmt.Sprintf("tone %d fine tune", t), "tone", t, 0x3E, 1, 0, 100,
				func(p *Patch) int { return p.FineTune(t) + 50 }},
			patchParameterTest{fmt.Sprintf("tone %d cutoff", t), "tone", t, 0x51, 1, 0, 127,
				func(p *Patch) int { return p.Cutoff(t) }},
			patchParameterTest{fmt.Sprintf("tone %d resonance", t), "tone", t, 0x53, 1, 0, 127,
				func(p *Patch) int { return p.Resonance(t) }},
			patchParameterTest{fmt.Sprintf("tone %d level", t), "tone", t, 0x65, 1, 0, 127,
				func(p *Patch) int { return p.ToneLevel(t) }},
			patchParameterTest{fmt.Sprintf("tone %d pan", t), "tone", t, 0x77, 1, 0, 127,
				func(p *Patch) int { return p.TonePan(t) }},
			patchParameterTest{fmt.Sprintf("tone %d chorus send", t), "tone", t, 0x7F, 1, 0, 127,
				func(p *Patch) int { return p.ChorusSend(t) }},
			patchParameterTest{fmt.Sprintf("tone %d reverb send", t), "tone", t, 0x80, 1, 0, 127,
				func(p *Patch) int { return p.ReverbSend(t) }},
		)
	}
	return append(common, tones...)
}

// rawValue reads the raw value of a parameter: the plain byte, or the
// combined 0-254 number for nibblized parameters.
func rawValue(current []byte, size int) int {
	if size == 2 {
		return int(current[0])<<4 | int(current[1])
	}
	return int(current[0])
}

// newParameterValue derives a new raw value for a parameter from its
// current raw value, so the test never writes the value that is already
// there. The result is always within the parameter's valid range: for
// single-byte parameters the value is rotated within min..max, for
// nibblized (2-byte) parameters the raw number is rotated and
// re-nibblized (high nibble first).
func newParameterValue(current []byte, size, min, max int) []byte {
	span := max - min + 1
	raw := (rawValue(current, size)-min+37)%span + min
	if size == 2 {
		return []byte{byte(raw >> 4), byte(raw & 0x0F)}
	}
	return []byte{byte(raw)}
}

// blockOf returns the block of p the parameter lives in.
func (pt patchParameterTest) blockOf(p *Patch) []byte {
	if pt.block == "common" {
		return p.Common
	}
	return p.Tones[pt.tone-1]
}

// TestPatchParameterTable checks the parameter table itself: every entry
// must address a valid offset within its block, sizes must be 1 or 2, the
// range must be non-empty, and the raw value transformation must produce
// a different, in-range value.
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
		if pt.min > pt.max {
			t.Errorf("%s: empty range %d..%d", pt.name, pt.min, pt.max)
		}
		current := make([]byte, pt.size)
		for i := range current {
			current[i] = byte(pt.min)
		}
		got := newParameterValue(current, pt.size, pt.min, pt.max)
		if bytes.Equal(got, current) {
			t.Errorf("%s: newParameterValue did not change the value", pt.name)
		}
		if raw := rawValue(got, pt.size); raw < pt.min || raw > pt.max {
			t.Errorf("%s: newParameterValue = %d out of range %d..%d", pt.name, raw, pt.min, pt.max)
		}
		if pt.size == 2 && (got[0] > 0x0F || got[1] > 0x0F) {
			t.Errorf("%s: newParameterValue = % X not valid nibbles", pt.name, got)
		}
	}
}

// TestIntegrationPatchParameters is a table test that writes every patch
// parameter of the temporary patch through SetPatchParameter, reads it
// back and verifies the change.
//
// The temporary patch is loaded from PRESET A patch 1 first: selecting a
// factory patch copies it into the temporary patch (volatile RAM, never
// written back to ROM), so the tests can modify all parameters freely
// without touching any user patch. Each subtest reads the affected block
// fresh before and after the write and leaves a settle delay for the
// device to apply the write.
func TestIntegrationPatchParameters(t *testing.T) {
	skipUnlessIntegration(t)
	c := newIntegrationClient(t)

	// Load PRESET A patch 1 into the temporary patch (read-only source).
	if _, err := c.ReadPatch(BankPresetA, 1); err != nil {
		t.Fatalf("load PRESET A patch 1 into temporary patch: %v", err)
	}

	for _, pt := range patchParameterTests() {
		t.Run(pt.name, func(t *testing.T) {
			readBlock := func() []byte {
				p, err := c.ReadTemporaryPatch()
				if err != nil {
					t.Fatalf("read temporary patch: %v", err)
				}
				return append([]byte{}, pt.blockOf(p)...)
			}

			block := readBlock()
			current := append([]byte{}, block[pt.offset:pt.offset+pt.size]...)
			value := newParameterValue(current, pt.size, pt.min, pt.max)

			if err := c.SetPatchParameter(BankPresetA, 0, pt.block, pt.tone, pt.offset, value); err != nil {
				t.Fatalf("SetPatchParameter: %v", err)
			}
			time.Sleep(50 * time.Millisecond)

			after := readBlock()
			if got := after[pt.offset : pt.offset+pt.size]; !bytes.Equal(got, value) {
				t.Errorf("parameter = % X, want % X", got, value)
			}
			// No other byte of the block may have changed.
			for i := range after {
				if i >= pt.offset && i < pt.offset+pt.size {
					continue
				}
				if after[i] != block[i] {
					t.Errorf("block byte %d changed: %d -> %d", i, block[i], after[i])
				}
			}

			// Restore the original value.
			if err := c.SetPatchParameter(BankPresetA, 0, pt.block, pt.tone, pt.offset, current); err != nil {
				t.Fatalf("restore parameter: %v", err)
			}
			time.Sleep(50 * time.Millisecond)
		})
	}
}
