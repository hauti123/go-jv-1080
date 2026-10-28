package controlchange

// Controller is a MIDI control change controller number (0-127).
type Controller uint8

// Defined and undefined control change controller numbers.
const (
	ControllerBankSelectMSB      Controller = 0  // also CC32 in some docs
	ControllerModulation         Controller = 1  // modulation wheel / lever
	ControllerBreath             Controller = 2  // breath controller
	ControllerFoot               Controller = 4  // foot controller
	ControllerPortamentoTime     Controller = 5  // portamento time
	ControllerDataEntryMSB       Controller = 6  // data entry MSB
	ControllerChannelVolume      Controller = 7  // channel volume
	ControllerBalance            Controller = 8  // balance
	ControllerPan                Controller = 10 // pan
	ControllerExpression         Controller = 11 // expression
	ControllerEffectControl1     Controller = 12 // effect control 1
	ControllerEffectControl2     Controller = 13 // effect control 2
	ControllerGeneralPurpose1    Controller = 16 // general purpose controller 1
	ControllerGeneralPurpose2    Controller = 17 // general purpose controller 2
	ControllerGeneralPurpose3    Controller = 18 // general purpose controller 3
	ControllerGeneralPurpose4    Controller = 19 // general purpose controller 4
	ControllerBankSelectLSB      Controller = 32 // bank select LSB
	ControllerModulationLSB      Controller = 33 // modulation wheel LSB
	ControllerBreathLSB          Controller = 34 // breath controller LSB
	ControllerFootLSB            Controller = 36 // foot controller LSB
	ControllerPortamentoTimeLSB  Controller = 37 // portamento time LSB
	ControllerDataEntryLSB       Controller = 38 // data entry LSB
	ControllerChannelVolumeLSB   Controller = 39 // channel volume LSB
	ControllerBalanceLSB         Controller = 40 // balance LSB
	ControllerPanLSB             Controller = 42 // pan LSB
	ControllerExpressionLSB      Controller = 43 // expression LSB
	ControllerEffectControl1LSB  Controller = 44 // effect control 1 LSB
	ControllerEffectControl2LSB  Controller = 45 // effect control 2 LSB
	ControllerGeneralPurpose1LSB Controller = 48 // general purpose controller 1 LSB
	ControllerGeneralPurpose2LSB Controller = 49 // general purpose controller 2 LSB
	ControllerGeneralPurpose3LSB Controller = 50 // general purpose controller 3 LSB
	ControllerGeneralPurpose4LSB Controller = 51 // general purpose controller 4 LSB

	ControllerHold1             Controller = 64 // damper pedal / hold 1 (sustain)
	ControllerPortamento        Controller = 65 // portamento on/off
	ControllerSostenuto         Controller = 66 // sostenuto
	ControllerSoft              Controller = 67 // soft pedal
	ControllerLegato            Controller = 68 // legato footswitch
	ControllerHold2             Controller = 69 // hold 2
	ControllerSoundVariation    Controller = 70 // sound controller 1: sound variation
	ControllerTimbre            Controller = 71 // sound controller 2: timbre / harmonic content
	ControllerReleaseTime       Controller = 72 // sound controller 3: release time
	ControllerAttackTime        Controller = 73 // sound controller 4: attack time
	ControllerBrightness        Controller = 74 // sound controller 5: brightness (cutoff)
	ControllerDecayTime         Controller = 75 // sound controller 6: decay time
	ControllerVibratoRate       Controller = 76 // sound controller 7: vibrato rate
	ControllerVibratoDepth      Controller = 77 // sound controller 8: vibrato depth
	ControllerVibratoDelay      Controller = 78 // sound controller 9: vibrato delay
	ControllerSoundController10 Controller = 79 // sound controller 10
	ControllerGeneralPurpose5   Controller = 80 // general purpose controller 5
	ControllerGeneralPurpose6   Controller = 81 // general purpose controller 6
	ControllerGeneralPurpose7   Controller = 82 // general purpose controller 7
	ControllerGeneralPurpose8   Controller = 83 // general purpose controller 8
	ControllerPortamentoCtrl    Controller = 84 // portamento control

	ControllerEffect1Depth Controller = 91 // effect 1 depth (reverb send)
	ControllerEffect2Depth Controller = 92 // effect 2 depth (tremolo)
	ControllerEffect3Depth Controller = 93 // effect 3 depth (chorus send)
	ControllerEffect4Depth Controller = 94 // effect 4 depth (celeste)
	ControllerEffect5Depth Controller = 95 // effect 5 depth (phaser)

	ControllerDataIncrement Controller = 96 // data increment (+1)
	ControllerDataDecrement Controller = 97 // data decrement (-1)

	ControllerNRPNSelectLSB Controller = 98  // non-registered parameter number LSB
	ControllerNRPNSelectMSB Controller = 99  // non-registered parameter number MSB
	ControllerRPNSelectLSB  Controller = 100 // registered parameter number LSB
	ControllerRPNSelectMSB  Controller = 101 // registered parameter number MSB

	ControllerAllSoundOff  Controller = 120 // all sound off
	ControllerResetAllCtrl Controller = 121 // reset all controllers
	ControllerLocalControl Controller = 122 // local control on/off
	ControllerAllNotesOff  Controller = 123 // all notes off
	ControllerOmniModeOff  Controller = 124 // omni mode off (+ all notes off)
	ControllerOmniModeOn   Controller = 125 // omni mode on (+ all notes off)
	ControllerMonoModeOn   Controller = 126 // mono mode on (+ poly off, + all notes off)
	ControllerPolyModeOn   Controller = 127 // poly mode on (+ mono off, + all notes off)
)

func (c Controller) String() string {
	switch c {
	case ControllerBankSelectMSB:
		return "BankSelectMSB"
	case ControllerModulation:
		return "Modulation"
	case ControllerBreath:
		return "Breath"
	case ControllerFoot:
		return "Foot"
	case ControllerPortamentoTime:
		return "PortamentoTime"
	case ControllerDataEntryMSB:
		return "DataEntryMSB"
	case ControllerChannelVolume:
		return "ChannelVolume"
	case ControllerBalance:
		return "Balance"
	case ControllerPan:
		return "Pan"
	case ControllerExpression:
		return "Expression"
	case ControllerEffectControl1:
		return "EffectControl1"
	case ControllerEffectControl2:
		return "EffectControl2"
	case ControllerGeneralPurpose1:
		return "GeneralPurpose1"
	case ControllerGeneralPurpose2:
		return "GeneralPurpose2"
	case ControllerGeneralPurpose3:
		return "GeneralPurpose3"
	case ControllerGeneralPurpose4:
		return "GeneralPurpose4"
	case ControllerBankSelectLSB:
		return "BankSelectLSB"
	case ControllerModulationLSB:
		return "ModulationLSB"
	case ControllerBreathLSB:
		return "BreathLSB"
	case ControllerFootLSB:
		return "FootLSB"
	case ControllerPortamentoTimeLSB:
		return "PortamentoTimeLSB"
	case ControllerDataEntryLSB:
		return "DataEntryLSB"
	case ControllerChannelVolumeLSB:
		return "ChannelVolumeLSB"
	case ControllerBalanceLSB:
		return "BalanceLSB"
	case ControllerPanLSB:
		return "PanLSB"
	case ControllerExpressionLSB:
		return "ExpressionLSB"
	case ControllerEffectControl1LSB:
		return "EffectControl1LSB"
	case ControllerEffectControl2LSB:
		return "EffectControl2LSB"
	case ControllerGeneralPurpose1LSB:
		return "GeneralPurpose1LSB"
	case ControllerGeneralPurpose2LSB:
		return "GeneralPurpose2LSB"
	case ControllerGeneralPurpose3LSB:
		return "GeneralPurpose3LSB"
	case ControllerGeneralPurpose4LSB:
		return "GeneralPurpose4LSB"
	case ControllerHold1:
		return "Hold1"
	case ControllerPortamento:
		return "Portamento"
	case ControllerSostenuto:
		return "Sostenuto"
	case ControllerSoft:
		return "Soft"
	case ControllerLegato:
		return "Legato"
	case ControllerHold2:
		return "Hold2"
	case ControllerSoundVariation:
		return "SoundVariation"
	case ControllerTimbre:
		return "Timbre"
	case ControllerReleaseTime:
		return "ReleaseTime"
	case ControllerAttackTime:
		return "AttackTime"
	case ControllerBrightness:
		return "Brightness"
	case ControllerDecayTime:
		return "DecayTime"
	case ControllerVibratoRate:
		return "VibratoRate"
	case ControllerVibratoDepth:
		return "VibratoDepth"
	case ControllerVibratoDelay:
		return "VibratoDelay"
	case ControllerSoundController10:
		return "SoundController10"
	case ControllerGeneralPurpose5:
		return "GeneralPurpose5"
	case ControllerGeneralPurpose6:
		return "GeneralPurpose6"
	case ControllerGeneralPurpose7:
		return "GeneralPurpose7"
	case ControllerGeneralPurpose8:
		return "GeneralPurpose8"
	case ControllerPortamentoCtrl:
		return "PortamentoControl"
	case ControllerEffect1Depth:
		return "Effect1Depth"
	case ControllerEffect2Depth:
		return "Effect2Depth"
	case ControllerEffect3Depth:
		return "Effect3Depth"
	case ControllerEffect4Depth:
		return "Effect4Depth"
	case ControllerEffect5Depth:
		return "Effect5Depth"
	case ControllerDataIncrement:
		return "DataIncrement"
	case ControllerDataDecrement:
		return "DataDecrement"
	case ControllerNRPNSelectLSB:
		return "NRPNSelectLSB"
	case ControllerNRPNSelectMSB:
		return "NRPNSelectMSB"
	case ControllerRPNSelectLSB:
		return "RPNSelectLSB"
	case ControllerRPNSelectMSB:
		return "RPNSelectMSB"
	case ControllerAllSoundOff:
		return "AllSoundOff"
	case ControllerResetAllCtrl:
		return "ResetAllControllers"
	case ControllerLocalControl:
		return "LocalControl"
	case ControllerAllNotesOff:
		return "AllNotesOff"
	case ControllerOmniModeOff:
		return "OmniModeOff"
	case ControllerOmniModeOn:
		return "OmniModeOn"
	case ControllerMonoModeOn:
		return "MonoModeOn"
	case ControllerPolyModeOn:
		return "PolyModeOn"
	}
	return "Controller(" + itoa(int(c)) + ")"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [4]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
