// Package jv1080 implements the Roland JV-1080 parameter address map and a
// client on top of the internal/midi transport layer.
package jv1080

import "go-jv-1080/internal/midi"

// Memory areas (first address byte) of the JV-1080 parameter address map.
const (
	areaSystem                  = 0x00
	areaTemporaryPerformance    = 0x01
	areaPerformanceTempPart     = 0x02 // part number in second byte (0-15), 0x09 = rhythm setup
	areaPatchModeTemporaryPatch = 0x03
	areaUserPerformance         = 0x10 // second byte 0x00-0x1F: performance 1-32
	areaUserRhythmSetup         = 0x10 // second byte 0x40/0x41: rhythm setup 1/2
	areaUserPatch               = 0x11 // second byte 0x00-0x7F: patch 1-128
	areaCardPerformance         = 0x20 // JV-1080 data card
	areaCardRhythmSetup         = 0x20 // second byte 0x40/0x41
	areaCardPatch               = 0x21 // JV-1080 data card
)

// Block sizes in bytes of the JV-1080 parameter blocks.
const (
	PatchCommonSize       = 0x48 // 72 bytes
	PatchToneSize         = 0x81 // 129 bytes
	PerformanceCommonSize = 0x40
	PerformancePartSize   = 0x13
	RhythmCommonSize      = 0x0C
	RhythmNoteSize        = 0x3A
)

// toneBaseOffsets are the block offsets (third address byte) of the four
// patch tones within a patch area.
var toneBaseOffsets = [4]int{0x10, 0x12, 0x14, 0x16}

// Bank identifies a patch or performance bank of the JV-1080.
type Bank int

// Patch and performance banks of the JV-1080.
const (
	BankUser Bank = iota
	BankPresetA
	BankPresetB
	BankPresetC
	BankGM
	BankCardA
	BankCardB
	BankExpA
	BankExpB
	BankExpC
	BankExpD
)

// Banks lists all selectable banks in display order.
var Banks = []Bank{
	BankUser, BankPresetA, BankPresetB, BankPresetC, BankGM,
	BankCardA, BankCardB, BankExpA, BankExpB, BankExpC, BankExpD,
}

func (b Bank) String() string {
	switch b {
	case BankUser:
		return "USER"
	case BankPresetA:
		return "PRESET A"
	case BankPresetB:
		return "PRESET B"
	case BankPresetC:
		return "PRESET C"
	case BankGM:
		return "GM"
	case BankCardA:
		return "CARD A"
	case BankCardB:
		return "CARD B"
	case BankExpA:
		return "EXP A"
	case BankExpB:
		return "EXP B"
	case BankExpC:
		return "EXP C"
	case BankExpD:
		return "EXP D"
	}
	return "UNKNOWN"
}

// PatchCount is the number of patches or performances per bank.
const PatchCount = 128

// PerformanceCount is the number of performances per performance bank.
const PerformanceCount = 32

// bankSelect returns the CC0/CC32 bank select values for a patch bank.
func (b Bank) bankSelect() (msb, lsb uint8) {
	switch b {
	case BankUser:
		return 0x50, 0x00
	case BankPresetA:
		return 0x51, 0x00
	case BankPresetB:
		return 0x51, 0x01
	case BankPresetC:
		return 0x51, 0x02
	case BankGM:
		return 0x51, 0x03
	case BankCardA:
		return 0x52, 0x00
	case BankCardB:
		return 0x53, 0x00
	case BankExpA:
		return 0x54, 0x00
	case BankExpB:
		return 0x54, 0x02
	case BankExpC:
		return 0x54, 0x04
	case BankExpD:
		return 0x54, 0x06
	}
	return 0, 0
}

// performanceBankSelect returns the CC0/CC32 bank select values for a
// performance bank.
func (b Bank) performanceBankSelect() (msb, lsb uint8) {
	switch b {
	case BankUser:
		return 0x50, 0x00
	case BankPresetA:
		return 0x51, 0x00
	case BankPresetB:
		return 0x51, 0x01
	case BankPresetC:
		return 0x51, 0x02
	case BankGM:
		return 0x51, 0x03
	case BankCardA:
		return 0x52, 0x00
	case BankCardB:
		return 0x53, 0x00
	case BankExpA:
		return 0x54, 0x00
	case BankExpB:
		return 0x54, 0x02
	case BankExpC:
		return 0x54, 0x04
	case BankExpD:
		return 0x54, 0x06
	}
	return 0, 0
}

// patchArea returns the 4-byte address of the patch with the given number
// (1-128) in the given bank.
func (b Bank) patchArea(number int) midi.Address {
	n := number - 1
	switch b {
	case BankUser:
		return midi.Address{areaUserPatch, byte(n), 0x00, 0x00}
	case BankCardA, BankCardB:
		return midi.Address{areaCardPatch, byte(n), 0x00, 0x00}
	}
	// Preset, GM and expansion banks are read-only on the device and are
	// not part of the writable address map; they are addressed through the
	// temporary patch after a program change. Return the temporary patch.
	return TemporaryPatchArea()
}

// userPatchArea returns the address of a user patch (number 1-128).
func userPatchArea(number int) midi.Address {
	return midi.Address{areaUserPatch, byte(number - 1), 0x00, 0x00}
}

// userPerformanceArea returns the address of a user performance (1-32).
func userPerformanceArea(number int) midi.Address {
	return midi.Address{areaUserPerformance, byte(number - 1), 0x00, 0x00}
}

// TemporaryPatchArea returns the address of the patch-mode temporary patch.
func TemporaryPatchArea() midi.Address {
	return midi.Address{areaPatchModeTemporaryPatch, 0x00, 0x00, 0x00}
}

// TemporaryPerformanceArea returns the address of the temporary performance.
func TemporaryPerformanceArea() midi.Address {
	return midi.Address{areaTemporaryPerformance, 0x00, 0x00, 0x00}
}

// performancePartArea returns the address of performance part p (1-16) of
// the temporary performance.
func performancePartArea(p int) midi.Address {
	return midi.Address{areaTemporaryPerformance, 0x00, byte(0x10 + (p - 1)), 0x00}
}

// patchToneArea returns the address of tone t (1-4) within a patch area.
func patchToneArea(area midi.Address, t int) midi.Address {
	return midi.Address{area[0], area[1], byte(toneBaseOffsets[t-1]), 0x00}
}

// patchCommonBlock returns the address and size of a patch's common block.
func patchCommonBlock(area midi.Address) (midi.Address, int) {
	return area, PatchCommonSize
}

// patchToneBlock returns the address and size of tone t (1-4) of a patch.
func patchToneBlock(area midi.Address, t int) (midi.Address, int) {
	return patchToneArea(area, t), PatchToneSize
}

// performanceCommonBlock returns the address and size of a performance's
// common block.
func performanceCommonBlock(area midi.Address) (midi.Address, int) {
	return area, PerformanceCommonSize
}

// performancePartBlock returns the address and size of part p (1-16) of a
// performance. Part p lives at third address byte 0x10 + (p - 1); the
// address value is (0x10 + (p - 1)) * 0x80 relative to the area base.
func performancePartBlock(area midi.Address, p int) (midi.Address, int) {
	return area.Add((0x10 + (p - 1)) * 0x80), PerformancePartSize
}
