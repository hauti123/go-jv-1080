// Package sysex implements the Roland JV-1080 SysEx parameter address map
// and high level functions to download patches and performances from the
// device and to modify single parameters.
package sysex

import (
	"go-jv-1080/internal/midi/message"
)

// Memory areas (first address byte) of the JV-1080 parameter address map.
const (
	AreaSystem                  = 0x00
	AreaTemporaryPerformance    = 0x01
	AreaPerformanceTempPart     = 0x02 // part number in second byte (0-15), 0x09 = rhythm setup
	AreaPatchModeTemporaryPatch = 0x03
	AreaUserPerformance         = 0x10 // second byte 0x00-0x1F: performance 1-32
	AreaUserRhythmSetup         = 0x10 // second byte 0x40/0x41: rhythm setup 1/2
	AreaUserPatch               = 0x11 // second byte 0x00-0x7F: patch 1-128
	AreaCardPerformance         = 0x20 // JV-1080 data card
	AreaCardRhythmSetup         = 0x20 // second byte 0x40/0x41
	AreaCardPatch               = 0x21 // JV-1080 data card
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

// Counts per bank.
const (
	UserPatchCount       = 128
	UserPerformanceCount = 32
	PatchToneCount       = 4
	PerformancePartCount = 16
)

// toneBaseOffsets are the block offsets (third address byte) of the four
// patch tones within a patch area.
var toneBaseOffsets = [PatchToneCount]int{0x10, 0x12, 0x14, 0x16}

// MaxDataSetBytes is the maximum number of data bytes per DT1 message.
const MaxDataSetBytes = 128

// UserPatchArea returns the address of user patch number (1-128).
func UserPatchArea(number int) message.Address {
	return message.Address{AreaUserPatch, byte(number - 1), 0x00, 0x00}
}

// UserPerformanceArea returns the address of user performance number
// (1-32).
func UserPerformanceArea(number int) message.Address {
	return message.Address{AreaUserPerformance, byte(number - 1), 0x00, 0x00}
}

// CardPatchArea returns the address of card patch number (1-128).
func CardPatchArea(number int) message.Address {
	return message.Address{AreaCardPatch, byte(number - 1), 0x00, 0x00}
}

// TemporaryPatchArea returns the address of the patch-mode temporary patch.
func TemporaryPatchArea() message.Address {
	return message.Address{AreaPatchModeTemporaryPatch, 0x00, 0x00, 0x00}
}

// TemporaryPerformanceArea returns the address of the temporary
// performance.
func TemporaryPerformanceArea() message.Address {
	return message.Address{AreaTemporaryPerformance, 0x00, 0x00, 0x00}
}

// PatchCommonBlock returns the address and size of a patch's common block.
func PatchCommonBlock(area message.Address) (message.Address, int) {
	return area, PatchCommonSize
}

// PatchToneBlock returns the address and size of tone t (1-4) of a patch.
func PatchToneBlock(area message.Address, t int) (message.Address, int) {
	return PatchToneArea(area, t), PatchToneSize
}

// PatchToneArea returns the address of tone t (1-4) within a patch area.
func PatchToneArea(area message.Address, t int) message.Address {
	return message.Address{area[0], area[1], byte(toneBaseOffsets[t-1]), 0x00}
}

// PerformanceCommonBlock returns the address and size of a performance's
// common block.
func PerformanceCommonBlock(area message.Address) (message.Address, int) {
	return area, PerformanceCommonSize
}

// PerformancePartBlock returns the address and size of part p (1-16) of a
// performance. Part p lives at third address byte 0x10 + (p - 1); the
// address value is (0x10 + (p - 1)) * 0x80 relative to the area base.
func PerformancePartBlock(area message.Address, p int) (message.Address, int) {
	return area.Add((0x10 + (p - 1)) * 0x80), PerformancePartSize
}
