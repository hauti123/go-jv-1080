package message

import (
	"bytes"
	"testing"
)

func TestControlChangeConstructor(t *testing.T) {
	cc := NewControlChange(2, ControllerPan, 100)
	want := []byte{0xB2, 0x0A, 0x64}
	if got := cc.Bytes(); !bytes.Equal(got, want) {
		t.Errorf("Bytes = % X, want % X", got, want)
	}
	back, err := ParseControlChange(cc.Bytes())
	if err != nil || back != cc {
		t.Errorf("round trip = %+v, %v", back, err)
	}
	if got := cc.String(); got != "CC ch=3 controller=Pan value=100" {
		t.Errorf("String = %q", got)
	}
}

func TestTypedControlChanges(t *testing.T) {
	cases := []struct {
		name string
		msg  Message
		want []byte
	}{
		{"BankSelect", NewBankSelect(0, 0x50, 0x00), []byte{0xB0, 0x00, 0x50, 0xB0, 0x20, 0x00}},
		{"Modulation", NewModulation(0, 64), []byte{0xB0, 0x01, 0x40}},
		{"Breath", NewBreath(1, 127), []byte{0xB1, 0x02, 0x7F}},
		{"Foot", NewFoot(0, 10), []byte{0xB0, 0x04, 0x0A}},
		{"PortamentoTime", NewPortamentoTime(0, 100), []byte{0xB0, 0x05, 0x64}},
		{"DataEntry", NewDataEntry(0, 127), []byte{0xB0, 0x06, 0x7F}},
		{"ChannelVolume", NewChannelVolume(0, 100), []byte{0xB0, 0x07, 0x64}},
		{"Balance", NewBalance(0, 32), []byte{0xB0, 0x08, 0x20}},
		{"Pan", NewPan(0, 64), []byte{0xB0, 0x0A, 0x40}},
		{"Expression", NewExpression(0, 127), []byte{0xB0, 0x0B, 0x7F}},
		{"EffectControl1", NewEffectControl1(0, 1), []byte{0xB0, 0x0C, 0x01}},
		{"EffectControl2", NewEffectControl2(0, 2), []byte{0xB0, 0x0D, 0x02}},
		{"GeneralPurpose1", NewGeneralPurpose(0, 1, 5), []byte{0xB0, 0x10, 0x05}},
		{"GeneralPurpose4", NewGeneralPurpose(0, 4, 5), []byte{0xB0, 0x13, 0x05}},
		{"Hold1 on", NewHold1(0, true), []byte{0xB0, 0x40, 0x7F}},
		{"Hold1 off", NewHold1(0, false), []byte{0xB0, 0x40, 0x00}},
		{"Portamento on", NewPortamento(0, true), []byte{0xB0, 0x41, 0x7F}},
		{"Sostenuto on", NewSostenuto(0, true), []byte{0xB0, 0x42, 0x7F}},
		{"Soft on", NewSoft(0, true), []byte{0xB0, 0x43, 0x7F}},
		{"Legato on", NewLegato(0, true), []byte{0xB0, 0x44, 0x7F}},
		{"Hold2 on", NewHold2(0, true), []byte{0xB0, 0x45, 0x7F}},
		{"SoundController1", NewSoundController(0, 1, 100), []byte{0xB0, 0x46, 0x64}},
		{"SoundController5 (brightness)", NewSoundController(0, 5, 100), []byte{0xB0, 0x4A, 0x64}},
		{"SoundController10", NewSoundController(0, 10, 100), []byte{0xB0, 0x4F, 0x64}},
		{"PortamentoControl", NewPortamentoControl(0, 60), []byte{0xB0, 0x54, 0x3C}},
		{"EffectDepth1", NewEffectDepth(0, 1, 64), []byte{0xB0, 0x5B, 0x40}},
		{"EffectDepth5", NewEffectDepth(0, 5, 64), []byte{0xB0, 0x5F, 0x40}},
		{"ReverbSend", NewReverbSend(0, 64), []byte{0xB0, 0x5B, 0x40}},
		{"ChorusSend", NewChorusSend(0, 32), []byte{0xB0, 0x5D, 0x20}},
		{"DataIncrement", NewDataIncrement(0), []byte{0xB0, 0x60, 0x00}},
		{"DataDecrement", NewDataDecrement(0), []byte{0xB0, 0x61, 0x00}},
		{"RPNSelect pitch bend", NewRPNSelect(0, RPNPitchBendSensitivity), []byte{0xB0, 0x65, 0x00, 0xB0, 0x64, 0x00}},
		{"RPNSelect fine tune", NewRPNSelect(0, RPNFineTune), []byte{0xB0, 0x65, 0x00, 0xB0, 0x64, 0x01}},
		{"RPNSelect coarse tune", NewRPNSelect(0, RPNCoarseTune), []byte{0xB0, 0x65, 0x00, 0xB0, 0x64, 0x02}},
		{"NRPNSelect", NewNRPNSelect(0, 0x12, 0x34), []byte{0xB0, 0x63, 0x12, 0xB0, 0x62, 0x34}},
		{"AllSoundOff", NewAllSoundOff(0), []byte{0xB0, 0x78, 0x00}},
		{"ResetAllControllers", NewResetAllControllers(0), []byte{0xB0, 0x79, 0x00}},
		{"LocalControl on", NewLocalControl(0, true), []byte{0xB0, 0x7A, 0x7F}},
		{"LocalControl off", NewLocalControl(0, false), []byte{0xB0, 0x7A, 0x00}},
		{"AllNotesOff", NewAllNotesOff(0), []byte{0xB0, 0x7B, 0x00}},
		{"OmniModeOff", NewOmniModeOff(0), []byte{0xB0, 0x7C, 0x00}},
		{"OmniModeOn", NewOmniModeOn(0), []byte{0xB0, 0x7D, 0x00}},
		{"MonoModeOn", NewMonoModeOn(0), []byte{0xB0, 0x7E, 0x00}},
		{"PolyModeOn", NewPolyModeOn(0), []byte{0xB0, 0x7F, 0x00}},
	}
	for _, c := range cases {
		if got := c.msg.Bytes(); !bytes.Equal(got, c.want) {
			t.Errorf("%s: Bytes = % X, want % X", c.name, got, c.want)
		}
		if c.msg.String() == "" {
			t.Errorf("%s: empty String", c.name)
		}
	}
}

func TestSwitchValue(t *testing.T) {
	if SwitchOff.OnOff() || !SwitchOn.OnOff() {
		t.Error("SwitchValue.OnOff wrong")
	}
	if SwitchOff.Value() != 0 || SwitchOn.Value() != 127 {
		t.Error("SwitchValue.Value wrong")
	}
	if SwitchOff.String() != "off" || SwitchOn.String() != "on" {
		t.Error("SwitchValue.String wrong")
	}
}

func TestControllerString(t *testing.T) {
	cases := map[Controller]string{
		ControllerBankSelectMSB: "BankSelectMSB",
		ControllerModulation:    "Modulation",
		ControllerBreath:        "Breath",
		ControllerChannelVolume: "ChannelVolume",
		ControllerPan:           "Pan",
		ControllerHold1:         "Hold1",
		ControllerBrightness:    "Brightness",
		ControllerEffect1Depth:  "Effect1Depth",
		ControllerRPNSelectMSB:  "RPNSelectMSB",
		ControllerAllNotesOff:   "AllNotesOff",
		ControllerPolyModeOn:    "PolyModeOn",
		Controller(3):           "Controller(3)",
	}
	for c, want := range cases {
		if got := c.String(); got != want {
			t.Errorf("Controller(%d).String() = %q, want %q", uint8(c), got, want)
		}
	}
}

func TestExistingConstructors(t *testing.T) {
	if got := NewProgramChange(9, 40).Bytes(); !bytes.Equal(got, []byte{0xC9, 0x28}) {
		t.Errorf("NewProgramChange Bytes = % X", got)
	}
	if got := NewNoteOn(0, 60, 100).Bytes(); !bytes.Equal(got, []byte{0x90, 60, 100}) {
		t.Errorf("NewNoteOn Bytes = % X", got)
	}
	if got := NewNoteOff(0, 60).Bytes(); !bytes.Equal(got, []byte{0x90, 60, 0}) {
		t.Errorf("NewNoteOff Bytes = % X", got)
	}
	if got := NewIdentityRequest(0x10).Bytes(); !bytes.Equal(got, []byte{0xF0, 0x7E, 0x10, 0x06, 0x01, 0xF7}) {
		t.Errorf("NewIdentityRequest Bytes = % X", got)
	}
	addr := NewAddress(0x03, 0x00, 0x10, 0x00)
	if got := NewDataRequest(0x10, addr, 1).Bytes(); !bytes.Equal(got, []byte{0xF0, 0x41, 0x10, 0x6A, 0x11, 0x03, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x01, 0x6C, 0xF7}) {
		t.Errorf("NewDataRequest Bytes = % X", got)
	}
	ds := NewDataSet(0x10, addr, []byte{0x00})
	if got := ds.Bytes(); !bytes.Equal(got, []byte{0xF0, 0x41, 0x10, 0x6A, 0x12, 0x03, 0x00, 0x10, 0x00, 0x00, 0x6D, 0xF7}) {
		t.Errorf("NewDataSet Bytes = % X", got)
	}
	// NewDataSet copies the data slice.
	data := []byte{0x01}
	ds = NewDataSet(0x10, addr, data)
	data[0] = 0x02
	if ds.Data[0] != 0x01 {
		t.Error("NewDataSet should copy the data slice")
	}
}
