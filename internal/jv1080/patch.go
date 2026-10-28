package jv1080

import (
	"fmt"
	"strings"
)

// Patch is a JV-1080 patch: a common block plus four tones. The raw block
// data is kept so that every single parameter can be updated individually
// and written back to the device.
type Patch struct {
	Area   [4]byte // memory area the patch was read from / is written to
	Common []byte  // PatchCommonSize bytes
	Tones  [4][]byte
}

// NewPatch returns an empty patch with correctly sized blocks.
func NewPatch(area [4]byte) *Patch {
	p := &Patch{Area: area, Common: make([]byte, PatchCommonSize)}
	for i := range p.Tones {
		p.Tones[i] = make([]byte, PatchToneSize)
	}
	return p
}

// Name returns the 12 character patch name.
func (p *Patch) Name() string {
	return decodeName(p.Common[0:12])
}

// SetName sets the 12 character patch name (padded with spaces).
func (p *Patch) SetName(name string) {
	encodeName(p.Common[0:12], name)
}

// ToneSwitch returns whether tone t (1-4) is on.
func (p *Patch) ToneSwitch(t int) bool {
	return p.Tones[t-1][0x00] != 0
}

// SetToneSwitch turns tone t (1-4) on or off.
func (p *Patch) SetToneSwitch(t int, on bool) {
	if on {
		p.Tones[t-1][0x00] = 1
	} else {
		p.Tones[t-1][0x00] = 0
	}
}

// PatchLevel returns the patch level (0-127), offset 0x2E in common.
func (p *Patch) PatchLevel() int { return int(p.Common[0x2E]) }

// SetPatchLevel sets the patch level (0-127).
func (p *Patch) SetPatchLevel(v int) { p.Common[0x2E] = clampByte(v, 0, 127) }

// PatchPan returns the patch pan (0-127, 64 = center), offset 0x2F.
func (p *Patch) PatchPan() int { return int(p.Common[0x2F]) }

// SetPatchPan sets the patch pan (0-127, 64 = center).
func (p *Patch) SetPatchPan(v int) { p.Common[0x2F] = clampByte(v, 0, 127) }

// ToneLevel returns the level of tone t (1-4) (0-127), offset 0x65.
func (p *Patch) ToneLevel(t int) int { return int(p.Tones[t-1][0x65]) }

// SetToneLevel sets the level of tone t (1-4) (0-127).
func (p *Patch) SetToneLevel(t, v int) { p.Tones[t-1][0x65] = clampByte(v, 0, 127) }

// TonePan returns the pan of tone t (1-4) (0-127, 64 = center), offset 0x77.
func (p *Patch) TonePan(t int) int { return int(p.Tones[t-1][0x77]) }

// SetTonePan sets the pan of tone t (1-4) (0-64 = L63..R63).
func (p *Patch) SetTonePan(t, v int) { p.Tones[t-1][0x77] = clampByte(v, 0, 127) }

// Cutoff returns the TVF cutoff frequency of tone t (1-4) (0-127),
// offset 0x51.
func (p *Patch) Cutoff(t int) int { return int(p.Tones[t-1][0x51]) }

// SetCutoff sets the TVF cutoff frequency of tone t (1-4) (0-127).
func (p *Patch) SetCutoff(t, v int) { p.Tones[t-1][0x51] = clampByte(v, 0, 127) }

// Resonance returns the TVF resonance of tone t (1-4) (0-127), offset 0x53.
func (p *Patch) Resonance(t int) int { return int(p.Tones[t-1][0x53]) }

// SetResonance sets the TVF resonance of tone t (1-4) (0-127).
func (p *Patch) SetResonance(t, v int) { p.Tones[t-1][0x53] = clampByte(v, 0, 127) }

// CoarseTune returns the coarse tune of tone t (1-4) in semitones
// (-48..+48), offset 0x3D.
func (p *Patch) CoarseTune(t int) int { return int(p.Tones[t-1][0x3D]) - 48 }

// SetCoarseTune sets the coarse tune of tone t (1-4) in semitones.
func (p *Patch) SetCoarseTune(t, v int) { p.Tones[t-1][0x3D] = clampByte(v+48, 0, 96) }

// FineTune returns the fine tune of tone t (1-4) in cents (-50..+50),
// offset 0x3E.
func (p *Patch) FineTune(t int) int { return int(p.Tones[t-1][0x3E]) - 50 }

// SetFineTune sets the fine tune of tone t (1-4) in cents.
func (p *Patch) SetFineTune(t, v int) { p.Tones[t-1][0x3E] = clampByte(v+50, 0, 100) }

// WaveGroupType returns the wave group type of tone t (1-4): 0 = INT,
// 1 = PCM card, 2 = EXP (SR-JV80 board), offset 0x01.
func (p *Patch) WaveGroupType(t int) int { return int(p.Tones[t-1][0x01]) }

// SetWaveGroupType sets the wave group type of tone t (1-4).
func (p *Patch) SetWaveGroupType(t, v int) { p.Tones[t-1][0x01] = clampByte(v, 0, 2) }

// WaveGroupID returns the wave group ID of tone t (1-4), offset 0x02.
// For INT: 1 = INT-A, 2 = INT-B. For EXP: SR-JV80 board number.
func (p *Patch) WaveGroupID(t int) int { return int(p.Tones[t-1][0x02]) }

// SetWaveGroupID sets the wave group ID of tone t (1-4).
func (p *Patch) SetWaveGroupID(t, v int) { p.Tones[t-1][0x02] = clampByte(v, 0, 127) }

// WaveNumber returns the wave number of tone t (1-4) (1-255). The raw
// value 0-254 is nibblized over two bytes, high nibble first, at offset
// 0x03-0x04.
func (p *Patch) WaveNumber(t int) int {
	return (int(p.Tones[t-1][0x03])<<4 | int(p.Tones[t-1][0x04])) + 1
}

// SetWaveNumber sets the wave number of tone t (1-4) (1-255). The raw
// value 0-254 is nibblized over two bytes, high nibble first.
func (p *Patch) SetWaveNumber(t, v int) {
	raw := clamp(v, 1, 255) - 1
	p.Tones[t-1][0x03] = byte(raw >> 4)
	p.Tones[t-1][0x04] = byte(raw & 0x0F)
}

// EFXType returns the EFX type of the patch (0-39), offset 0x0C.
func (p *Patch) EFXType() int { return int(p.Common[0x0C]) }

// SetEFXType sets the EFX type of the patch (0-39).
func (p *Patch) SetEFXType(v int) { p.Common[0x0C] = clampByte(v, 0, 39) }

// ReverbSend returns the reverb send level of tone t (1-4) (0-127),
// offset 0x80 within the tone block (address 01 00).
func (p *Patch) ReverbSend(t int) int { return int(p.Tones[t-1][0x80]) }

// SetReverbSend sets the reverb send level of tone t (1-4) (0-127).
func (p *Patch) SetReverbSend(t, v int) { p.Tones[t-1][0x80] = clampByte(v, 0, 127) }

// ChorusSend returns the chorus send level of tone t (1-4) (0-127),
// offset 0x7F.
func (p *Patch) ChorusSend(t int) int { return int(p.Tones[t-1][0x7F]) }

// SetChorusSend sets the chorus send level of tone t (1-4) (0-127).
func (p *Patch) SetChorusSend(t, v int) { p.Tones[t-1][0x7F] = clampByte(v, 0, 127) }

// String returns a human readable summary of the patch.
func (p *Patch) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Patch %q level=%d pan=%d EFX=%d", p.Name(), p.PatchLevel(), p.PatchPan(), p.EFXType())
	for t := 1; t <= 4; t++ {
		fmt.Fprintf(&sb, "\n  Tone %d: switch=%v wave=%d/%d/%d level=%d pan=%d cutoff=%d res=%d",
			t, p.ToneSwitch(t), p.WaveGroupType(t), p.WaveGroupID(t), p.WaveNumber(t),
			p.ToneLevel(t), p.TonePan(t), p.Cutoff(t), p.Resonance(t))
	}
	return sb.String()
}

// Performance is a JV-1080 performance: a common block plus 16 parts.
type Performance struct {
	Area   [4]byte
	Common []byte // PerformanceCommonSize bytes
	Parts  [16][]byte
}

// NewPerformance returns an empty performance with correctly sized blocks.
func NewPerformance(area [4]byte) *Performance {
	p := &Performance{Area: area, Common: make([]byte, PerformanceCommonSize)}
	for i := range p.Parts {
		p.Parts[i] = make([]byte, PerformancePartSize)
	}
	return p
}

// Name returns the 12 character performance name.
func (p *Performance) Name() string {
	return decodeName(p.Common[0:12])
}

// SetName sets the 12 character performance name (padded with spaces).
func (p *Performance) SetName(name string) {
	encodeName(p.Common[0:12], name)
}

// PartSwitch returns whether part n (1-16) is on, offset 0x0D + n - 1.
func (p *Performance) PartSwitch(n int) bool {
	return p.Common[0x0D+(n-1)] != 0
}

// SetPartSwitch turns part n (1-16) on or off.
func (p *Performance) SetPartSwitch(n int, on bool) {
	if on {
		p.Common[0x0D+(n-1)] = 1
	} else {
		p.Common[0x0D+(n-1)] = 0
	}
}

// PartLevel returns the level of part n (1-16) (0-127), offset 0x06.
func (p *Performance) PartLevel(n int) int { return int(p.Parts[n-1][0x06]) }

// SetPartLevel sets the level of part n (1-16) (0-127).
func (p *Performance) SetPartLevel(n, v int) { p.Parts[n-1][0x06] = clampByte(v, 0, 127) }

// PartPan returns the pan of part n (1-16) (0-127, 64 = center), offset 0x07.
func (p *Performance) PartPan(n int) int { return int(p.Parts[n-1][0x07]) }

// SetPartPan sets the pan of part n (1-16).
func (p *Performance) SetPartPan(n, v int) { p.Parts[n-1][0x07] = clampByte(v, 0, 127) }

// PartChannel returns the MIDI channel of part n (1-16) (0-15 = ch 1-16,
// 16 = RX channel from the performance common), offset 0x01.
func (p *Performance) PartChannel(n int) int { return int(p.Parts[n-1][0x01]) }

// SetPartChannel sets the MIDI channel of part n (1-16).
func (p *Performance) SetPartChannel(n, ch int) { p.Parts[n-1][0x01] = clampByte(ch, 0, 16) }

// PartPatchGroup returns the patch group type of part n (1-16):
// 0 = USER&PRESET, 1 = CARD, 2 = EXP, offset 0x02.
func (p *Performance) PartPatchGroup(n int) int { return int(p.Parts[n-1][0x02]) }

// PartPatchGroupID returns the patch group ID of part n (1-16), offset 0x03.
// Group IDs: 1 = USER, 2 = CARD, 3 = PR-A, 4 = PR-B, 5 = PR-C, 6 = PR-D (GM),
// 7-8 = EXP boards.
func (p *Performance) PartPatchGroupID(n int) int { return int(p.Parts[n-1][0x03]) }

// PartPatchNumber returns the patch number of part n (1-16) (1-128). The
// raw value 0-127 is nibblized over two bytes, high nibble first, at
// offset 0x04-0x05.
func (p *Performance) PartPatchNumber(n int) int {
	return (int(p.Parts[n-1][0x04])<<4 | int(p.Parts[n-1][0x05])) + 1
}

// SetPartPatchNumber sets the patch number of part n (1-16) (1-128).
func (p *Performance) SetPartPatchNumber(n, v int) {
	raw := clamp(v, 1, 128) - 1
	p.Parts[n-1][0x04] = byte(raw >> 4)
	p.Parts[n-1][0x05] = byte(raw & 0x0F)
}

// String returns a human readable summary of the performance.
func (p *Performance) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Performance %q", p.Name())
	for n := 1; n <= 16; n++ {
		if p.PartSwitch(n) {
			fmt.Fprintf(&sb, "\n  Part %2d: ch=%2d group=%d id=%d patch=%3d level=%d pan=%d",
				n, p.PartChannel(n)+1, p.PartPatchGroup(n), p.PartPatchGroupID(n),
				p.PartPatchNumber(n), p.PartLevel(n), p.PartPan(n))
		}
	}
	return sb.String()
}

func decodeName(b []byte) string {
	s := make([]byte, 0, len(b))
	for _, c := range b {
		if c < 32 || c > 127 {
			c = ' '
		}
		s = append(s, c)
	}
	return strings.TrimRight(string(s), " ")
}

func encodeName(dst []byte, name string) {
	for i := range dst {
		if i < len(name) {
			dst[i] = clampByte(int(name[i]), 32, 127)
		} else {
			dst[i] = ' '
		}
	}
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampByte(v, min, max int) byte {
	return byte(clamp(v, min, max))
}
