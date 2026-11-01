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

// editParam describes one parameter write to verify, addressed by its
// device-layer key ("patch/temp/..."). label names the CLI form the write
// corresponds to.
type editParam struct {
	label string
	key   string
	// aux is an optional range [start, start+len) within the same block
	// where additional byte changes are allowed as a side effect of the
	// write (the device resets the 12 EFX parameters 0x0D-0x18 when the
	// EFX type at 0x0C changes).
	aux [2]int
	// setup holds raw writes applied before the edit and reverted after
	// it, used to relax coupled parameters (the device ignores a velocity
	// range upper below the velocity range lower).
	setup []setupWrite
}

// setupWrite is a raw parameter write applied before an edit, addressed
// by a device-layer key.
type setupWrite struct {
	key   string
	value []byte
}

// blockToneOffset resolves an editParam's key into its block, tone and
// offset within the block for the patch-dump comparison.
func blockToneOffset(e editParam) (block string, tone, offset, size int, err error) {
	p, err := jv1080.ResolveParam(e.key)
	if err != nil {
		return "", 0, 0, 0, err
	}
	block = p.Block
	if p.Block == "common" {
		tone = 0
	} else {
		var t int
		if _, err := fmt.Sscanf(p.Block, "tone/%d", &t); err != nil {
			return "", 0, 0, 0, fmt.Errorf("bad block %q", p.Block)
		}
		tone = t
	}
	return block, tone, p.Offset, p.Size(), nil
}

// candidateValues returns test values for a parameter, derived from its
// original raw bytes so that they always differ from them. The JV-1080
// silently ignores a DT1 whose value is outside the valid range of the
// addressed parameter (which depends on other parameters, e.g. the wave
// number must not exceed the wave count of the wave group), so the test
// tries the candidates in order and uses the first one the device
// accepts. Single-byte parameters get orig^1 first, then a full scan;
// nibblized (two-nibble) parameters get bit-flipped combined values
// first, then a scan of all 256 combined values.
func candidateValues(orig []byte) [][]byte {
	var out [][]byte
	seen := map[string]bool{}
	add := func(v ...byte) {
		key := string(v)
		if seen[key] {
			return
		}
		seen[key] = true
		out = append(out, append([]byte(nil), v...))
	}
	if len(orig) == 1 {
		o := orig[0]
		add(o ^ 1)
		for v := 0; v < 128; v++ {
			if byte(v) != o {
				add(byte(v))
			}
		}
		return out
	}
	// Nibblized parameter: the two bytes are the high and low nibble of
	// the combined value.
	h, l := orig[0], orig[1]
	add(h^1, l)
	add(h, l^1)
	add(h^1, l^1)
	c := int(h)<<4 | int(l)
	if c > 0 {
		add(byte((c-1)>>4), byte((c-1)&0xF))
	}
	if c < 0xFF {
		add(byte((c+1)>>4), byte((c+1)&0xF))
	}
	for v := 0; v < 0x100; v++ {
		if v != c {
			add(byte(v>>4), byte(v&0xF))
		}
	}
	return out
}

// findAcceptedValue writes the candidate values for the parameter at addr
// one by one and returns the first one the device accepts (reads back
// unchanged at exactly that address). Candidates the device modifies
// instead of accepting are reverted before the next try.
func findAcceptedValue(t *testing.T, c *jv1080.Client, e editParam, orig []byte) []byte {
	t.Helper()
	p, err := jv1080.ResolveParam(e.key)
	if err != nil {
		t.Fatalf("resolve %q: %v", e.key, err)
	}
	for _, cand := range candidateValues(orig) {
		ui, err := p.DecodeValue(cand)
		if err != nil {
			t.Fatalf("decode candidate: %v", err)
		}
		if err := c.SetParameter(e.key, ui); err != nil {
			t.Fatalf("set parameter %s: %v", e.key, err)
		}
		got, err := c.GetParameter(e.key)
		if err != nil {
			t.Fatalf("read back parameter: %v", err)
		}
		gotBytes, err := p.EncodeValue(got)
		if err != nil {
			t.Fatalf("re-encode: %v", err)
		}
		if equalBytes(gotBytes, cand) {
			return cand
		}
		if !equalBytes(gotBytes, orig) {
			// The device stored something else (e.g. clamped): undo
			// before trying the next candidate.
			origUI, _ := p.DecodeValue(orig)
			if err := c.SetParameter(e.key, origUI); err != nil {
				t.Fatalf("revert parameter: %v", err)
			}
		}
	}
	t.Fatalf("%s: device rejected every candidate value for %s (orig %v)",
		e.label, e.key, orig)
	return nil
}

// verifyEdit performs one parameter write and checks it against a full
// patch dump: exactly the parameter's bytes change, at exactly the
// expected block/tone/offset, to exactly the written value (plus allowed
// side effects in e.aux). The original bytes are restored afterwards.
func verifyEdit(t *testing.T, c *jv1080.Client, e editParam) {
	t.Helper()
	block, tone, offset, size, err := blockToneOffset(e)
	if err != nil {
		t.Fatalf("resolve edit key: %v", err)
	}
	// Optional setup writes, reverted at the end.
	var setupOrig [][]byte
	if len(e.setup) > 0 {
		pre, err := c.ReadTemporaryPatch()
		if err != nil {
			t.Fatalf("dump before setup: %v", err)
		}
		blk := blockOf(pre, block, tone)
		for _, s := range e.setup {
			sp, err := jv1080.ResolveParam(s.key)
			if err != nil {
				t.Fatalf("resolve setup key: %v", err)
			}
			if sp.Offset+len(s.value) > len(blk) {
				t.Fatalf("setup offset 0x%02X outside block", sp.Offset)
			}
			setupOrig = append(setupOrig, append([]byte(nil), blk[sp.Offset:sp.Offset+len(s.value)]...))
			sui, err := sp.DecodeValue(s.value)
			if err != nil {
				t.Fatalf("decode setup value: %v", err)
			}
			if err := c.SetParameter(s.key, sui); err != nil {
				t.Fatalf("setup write: %v", err)
			}
			got, err := c.GetParameter(s.key)
			if err != nil {
				t.Fatalf("read back setup: %v", err)
			}
			gotBytes, _ := sp.EncodeValue(got)
			if !equalBytes(gotBytes, s.value) {
				t.Fatalf("setup write at %s offset 0x%02X not accepted (wrote %v, reads %v)",
					s.key, sp.Offset, s.value, gotBytes)
			}
		}
	}
	before, err := c.ReadTemporaryPatch()
	if err != nil {
		t.Fatalf("dump before edit: %v", err)
	}
	blk := blockOf(before, block, tone)
	if offset+size > len(blk) {
		t.Fatalf("offset 0x%02X (+%d bytes) outside block (%d bytes)", offset, size, len(blk))
	}
	orig := append([]byte(nil), blk[offset:offset+size]...)
	value := findAcceptedValue(t, c, e, orig)
	after, err := c.ReadTemporaryPatch()
	if err != nil {
		t.Fatalf("dump after edit: %v", err)
	}
	checkEdit(t, before, after, e, orig, value)
	// Restore the original bytes immediately so the next parameter
	// starts from the same state, then revert the setup writes.
	p, err := jv1080.ResolveParam(e.key)
	if err != nil {
		t.Fatalf("resolve edit key: %v", err)
	}
	origUI, _ := p.DecodeValue(orig)
	if err := c.SetParameter(e.key, origUI); err != nil {
		t.Fatalf("restore parameter: %v", err)
	}
	for i, s := range e.setup {
		sp, _ := jv1080.ResolveParam(s.key)
		sui, _ := sp.DecodeValue(setupOrig[i])
		if err := c.SetParameter(s.key, sui); err != nil {
			t.Fatalf("revert setup: %v", err)
		}
	}
}

// integrationEdits enumerates one write per addressable parameter for the
// given tone. Common parameters are enumerated once (tone 0), per-Tone
// parameters for the requested tone; a second pass over all four tones is
// done by the test itself to verify the tone block bases.
func integrationEdits(tone int) []editParam {
	var edits []editParam
	if tone == 0 {
		// Patch common block parameters, enumerated from the device
		// layer's parameter list. The name (12 bytes) is skipped: it is
		// a multi-byte parameter and covered separately.
		for _, p := range jv1080.ListParams() {
			if p.Block != "common" || p.Key == "patch/temp/common/name" {
				continue
			}
			e := editParam{label: "common " + p.Key[len("patch/temp/common/"):], key: p.Key}
			if p.Offset == 0x0C {
				// Changing the EFX type resets the 12 EFX parameters.
				e.aux = [2]int{0x0D, 12}
			}
			edits = append(edits, e)
		}
		return edits
	}

	// Per-Tone parameters of the given tone, enumerated from the device
	// layer's parameter list.
	prefix := fmt.Sprintf("patch/temp/tone/%d/", tone)
	for _, p := range jv1080.ListParams() {
		if !strings.HasPrefix(p.Key, prefix) {
			continue
		}
		e := editParam{label: p.Key[len("patch/temp/"):], key: p.Key}
		if p.Key == fmt.Sprintf("patch/temp/tone/%d/velupper", tone) {
			// The device ignores a velocity range upper below the
			// velocity range lower; relax the lower value first.
			e.setup = []setupWrite{{key: fmt.Sprintf("patch/temp/tone/%d/vellower", tone), value: []byte{1}}}
		}
		edits = append(edits, e)
	}
	return edits
}

// TestIntegrationPatchEditAddresses verifies that every parameter write of
// the CLI hits the right byte in the device: for each addressable
// parameter it dumps the temporary patch, writes a distinctive test value
// through Client.SetPatchParameter, dumps again and asserts that exactly
// the parameter's byte(s) changed at the expected offset. Afterwards the
// original temporary patch is restored and verified.
//
// The JV-1080 silently ignores writes whose value is outside the valid
// range of the addressed parameter (e.g. a wave number beyond the wave
// count of the wave group, or a one-byte write to a nibblized parameter),
// so the test value is chosen by probing candidate values derived from
// the original bytes and using the first one the device accepts.
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
			verifyEdit(t, c, e)
		})
	}

	// Tone block bases: the same offset in tones 2-4 must hit the
	// corresponding byte in that tone's block.
	for tone := 2; tone <= 4; tone++ {
		tone := tone
		t.Run(fmt.Sprintf("tone %d base", tone), func(t *testing.T) {
			verifyEdit(t, c, editParam{
				label: fmt.Sprintf("tone %d cutoff", tone),
				key:   fmt.Sprintf("patch/temp/tone/%d/tvf/cutoff", tone),
			})
		})
	}
}

// blockOf returns the block (common or tone) of the patch that the edit
// addresses.
func blockOf(p *jv1080.Patch, block string, tone int) []byte {
	if block == "common" {
		return p.Common
	}
	return p.Tones[tone-1]
}

// byteDiff is one changed byte between two patch dumps.
type byteDiff struct {
	block string
	tone  int // 1-4, 0 for common
	off   int
	from  byte
	to    byte
}

func (d byteDiff) String() string {
	if d.block == "common" {
		return fmt.Sprintf("common %02X: %d -> %d", d.off, d.from, d.to)
	}
	return fmt.Sprintf("tone %d %02X: %d -> %d", d.tone, d.off, d.from, d.to)
}

func joinDiffs(ds []byteDiff) string {
	parts := make([]string, len(ds))
	for i, d := range ds {
		parts[i] = d.String()
	}
	return strings.Join(parts, "; ")
}

// collectDiffs returns all bytes that differ between two patch dumps.
func collectDiffs(before, after *jv1080.Patch) []byteDiff {
	var diffs []byteDiff
	for i := range before.Common {
		if before.Common[i] != after.Common[i] {
			diffs = append(diffs, byteDiff{block: "common", off: i, from: before.Common[i], to: after.Common[i]})
		}
	}
	for t := 0; t < 4; t++ {
		for i := range before.Tones[t] {
			if before.Tones[t][i] != after.Tones[t][i] {
				diffs = append(diffs, byteDiff{block: "tone", tone: t + 1, off: i,
					from: before.Tones[t][i], to: after.Tones[t][i]})
			}
		}
	}
	return diffs
}

// checkEdit asserts that between the two dumps exactly the bytes of the
// edited parameter changed: each byte at the edit's offset (within its
// block) that differs between orig and value equals the written value,
// and no other byte differs except in the allowed side-effect range
// e.aux.
func checkEdit(t *testing.T, before, after *jv1080.Patch, e editParam, orig, value []byte) {
	t.Helper()
	block, tone, offset, _, err := blockToneOffset(e)
	if err != nil {
		t.Fatalf("resolve edit key: %v", err)
	}
	var want []byteDiff
	for i, v := range value {
		if v != orig[i] {
			want = append(want, byteDiff{block: block, tone: tone, off: offset + i, from: orig[i], to: v})
		}
	}
	var unexpected []byteDiff
	for _, d := range collectDiffs(before, after) {
		if e.aux[1] > 0 && d.block == block && d.tone == tone &&
			d.off >= e.aux[0] && d.off < e.aux[0]+e.aux[1] {
			continue
		}
		unexpected = append(unexpected, d)
	}
	if len(unexpected) != len(want) {
		t.Errorf("%s: %d byte(s) changed, want %d: %s (want %s)",
			e.label, len(unexpected), len(want), joinDiffs(unexpected), joinDiffs(want))
		return
	}
	for i := range unexpected {
		if unexpected[i] != want[i] {
			t.Errorf("%s: change at %s, want %s", e.label, unexpected[i], want[i])
		}
	}
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
