package jv1080

import (
	"testing"
)

func TestBankNames(t *testing.T) {
	want := []string{
		"USER", "PRESET A", "PRESET B", "PRESET C", "GM",
		"CARD A", "CARD B", "EXP A", "EXP B", "EXP C", "EXP D",
	}
	if len(Banks) != len(want) {
		t.Fatalf("len(Banks) = %d, want %d", len(Banks), len(want))
	}
	for i, b := range Banks {
		if b.String() != want[i] {
			t.Errorf("Banks[%d].String() = %q, want %q", i, b.String(), want[i])
		}
	}
}

func TestBankSelect(t *testing.T) {
	cases := []struct {
		bank     Bank
		msb, lsb uint8
	}{
		{BankUser, 0x50, 0x00},
		{BankPresetA, 0x51, 0x00},
		{BankPresetB, 0x51, 0x01},
		{BankPresetC, 0x51, 0x02},
		{BankGM, 0x51, 0x03},
		{BankCardA, 0x52, 0x00},
		{BankCardB, 0x53, 0x00},
		{BankExpA, 0x54, 0x00},
		{BankExpB, 0x54, 0x02},
		{BankExpC, 0x54, 0x04},
		{BankExpD, 0x54, 0x06},
	}
	for _, c := range cases {
		msb, lsb := c.bank.bankSelect()
		if msb != c.msb || lsb != c.lsb {
			t.Errorf("%s bankSelect = %02X %02X, want %02X %02X", c.bank, msb, lsb, c.msb, c.lsb)
		}
	}
}

func TestPatchAreas(t *testing.T) {
	if got := TemporaryPatchArea(); got != [4]byte{0x03, 0x00, 0x00, 0x00} {
		t.Errorf("TemporaryPatchArea = %s", got)
	}
	if got := userPatchArea(1); got != [4]byte{0x11, 0x00, 0x00, 0x00} {
		t.Errorf("userPatchArea(1) = %s", got)
	}
	if got := userPatchArea(128); got != [4]byte{0x11, 0x7F, 0x00, 0x00} {
		t.Errorf("userPatchArea(128) = %s", got)
	}
	if got := BankCardA.patchArea(1); got != [4]byte{0x21, 0x00, 0x00, 0x00} {
		t.Errorf("CardA.patchArea(1) = %s", got)
	}
	if got := userPerformanceArea(32); got != [4]byte{0x10, 0x1F, 0x00, 0x00} {
		t.Errorf("userPerformanceArea(32) = %s", got)
	}
	if got := patchToneArea(TemporaryPatchArea(), 4); got != [4]byte{0x03, 0x00, 0x16, 0x00} {
		t.Errorf("patchToneArea(temp, 4) = %s", got)
	}
	// Tone block spans 129 bytes: tone 1 ends at 03 00 10 00 + 0x80 = 03 00 11 00.
	addr, size := patchToneBlock(TemporaryPatchArea(), 1)
	if addr != [4]byte{0x03, 0x00, 0x10, 0x00} || size != 0x81 {
		t.Errorf("patchToneBlock(temp, 1) = %s, %d", addr, size)
	}
	// Performance parts live at third address byte 0x10 + (p - 1).
	if addr, size := performancePartBlock(TemporaryPerformanceArea(), 1); addr != [4]byte{0x01, 0x00, 0x10, 0x00} || size != PerformancePartSize {
		t.Errorf("performancePartBlock(temp, 1) = %s, %d", addr, size)
	}
	if addr, size := performancePartBlock(userPerformanceArea(1), 16); addr != [4]byte{0x10, 0x00, 0x1F, 0x00} || size != PerformancePartSize {
		t.Errorf("performancePartBlock(user perf 1, 16) = %s, %d", addr, size)
	}
}

func TestPatchNameRoundTrip(t *testing.T) {
	p := NewPatch(TemporaryPatchArea())
	p.SetName("My Patch")
	if got := p.Name(); got != "My Patch" {
		t.Errorf("Name = %q, want %q", got, "My Patch")
	}
	p.SetName("A very long patch name")
	if got := p.Name(); got != "A very long" {
		t.Errorf("Name = %q, want truncated to 12 chars", got)
	}
}

func TestPatchParameters(t *testing.T) {
	p := NewPatch(TemporaryPatchArea())
	p.SetPatchLevel(100)
	p.SetPatchPan(32)
	p.SetToneSwitch(2, true)
	p.SetToneLevel(3, 90)
	p.SetTonePan(1, 127)
	p.SetCutoff(1, 70)
	p.SetResonance(4, 12)
	p.SetCoarseTune(1, -12)
	p.SetFineTune(2, 25)
	p.SetWaveGroupType(1, 0)
	p.SetWaveGroupID(1, 2)
	p.SetWaveNumber(1, 193)
	p.SetEFXType(39)
	p.SetReverbSend(1, 64)
	p.SetChorusSend(2, 32)

	if p.PatchLevel() != 100 || p.PatchPan() != 32 {
		t.Errorf("level/pan = %d/%d", p.PatchLevel(), p.PatchPan())
	}
	if !p.ToneSwitch(2) || p.ToneSwitch(1) {
		t.Error("tone switch wrong")
	}
	if p.ToneLevel(3) != 90 || p.TonePan(1) != 127 {
		t.Error("tone level/pan wrong")
	}
	if p.Cutoff(1) != 70 || p.Resonance(4) != 12 {
		t.Error("cutoff/resonance wrong")
	}
	if p.CoarseTune(1) != -12 || p.FineTune(2) != 25 {
		t.Errorf("coarse/fine = %d/%d", p.CoarseTune(1), p.FineTune(2))
	}
	if p.WaveGroupType(1) != 0 || p.WaveGroupID(1) != 2 || p.WaveNumber(1) != 193 {
		t.Errorf("wave = %d/%d/%d", p.WaveGroupType(1), p.WaveGroupID(1), p.WaveNumber(1))
	}
	if p.EFXType() != 39 || p.ReverbSend(1) != 64 || p.ChorusSend(2) != 32 {
		t.Error("efx/sends wrong")
	}

	// Out of range values are clamped.
	p.SetPatchLevel(200)
	if p.PatchLevel() != 127 {
		t.Errorf("PatchLevel not clamped: %d", p.PatchLevel())
	}
	p.SetWaveNumber(1, 300)
	if p.WaveNumber(1) != 255 {
		t.Errorf("WaveNumber not clamped: %d", p.WaveNumber(1))
	}
}

func TestPerformanceParameters(t *testing.T) {
	p := NewPerformance(TemporaryPerformanceArea())
	p.SetName("Live Set")
	p.SetPartSwitch(1, true)
	p.SetPartChannel(1, 0)
	p.SetPartPatchNumber(1, 64)
	p.SetPartLevel(1, 100)
	p.SetPartPan(2, 96)

	if p.Name() != "Live Set" {
		t.Errorf("Name = %q", p.Name())
	}
	if !p.PartSwitch(1) || p.PartSwitch(2) {
		t.Error("part switch wrong")
	}
	if p.PartChannel(1) != 0 || p.PartPatchNumber(1) != 64 {
		t.Error("channel/patch wrong")
	}
	if p.PartLevel(1) != 100 || p.PartPan(2) != 96 {
		t.Error("level/pan wrong")
	}
}

func TestPresetPatchNames(t *testing.T) {
	if name, ok := PresetPatchName(BankPresetA, 1); !ok || name != "64voicePiano" {
		t.Errorf("PresetA 1 = %q, %v", name, ok)
	}
	if name, ok := PresetPatchName(BankGM, 1); !ok || name != "Piano 1" {
		t.Errorf("GM 1 = %q, %v", name, ok)
	}
	if _, ok := PresetPatchName(BankUser, 1); ok {
		t.Error("user bank should have no factory names")
	}
	if _, ok := PresetPatchName(BankPresetA, 129); ok {
		t.Error("out of range number should fail")
	}
	for _, bank := range []Bank{BankPresetA, BankPresetB, BankPresetC, BankGM} {
		count := 0
		for n := 1; n <= 128; n++ {
			if _, ok := PresetPatchName(bank, n); ok {
				count++
			}
		}
		if count != 128 {
			t.Errorf("%s: %d names, want 128", bank, count)
		}
	}
}
