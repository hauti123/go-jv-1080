package controlchange

import "fmt"

// BankSelect is a bank select control change pair (CC0/CC32). The JV-1080
// uses the MSB/LSB pair to select a patch bank.
type BankSelect struct {
	Channel uint8 // 0-15
	MSB     uint8 // 0-127
	LSB     uint8 // 0-127
}

// NewBankSelect returns a bank select for the given channel.
func NewBankSelect(channel, msb, lsb uint8) BankSelect {
	return BankSelect{Channel: channel, MSB: msb, LSB: lsb}
}

// Bytes serializes the message to its wire representation: the CC0 (MSB)
// message followed by the CC32 (LSB) message.
func (m BankSelect) Bytes() []byte {
	return append(ccMessage(m.Channel, ControllerBankSelectMSB, m.MSB),
		ccMessage(m.Channel, ControllerBankSelectLSB, m.LSB)...)
}

func (m BankSelect) String() string {
	return fmt.Sprintf("BankSelect ch=%d msb=%d lsb=%d", m.Channel+1, m.MSB, m.LSB)
}

// Modulation is the modulation wheel control change (CC1).
type Modulation struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewModulation returns a modulation wheel message.
func NewModulation(channel, value uint8) Modulation {
	return Modulation{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m Modulation) Bytes() []byte {
	return ccMessage(m.Channel, ControllerModulation, m.Value)
}

func (m Modulation) String() string {
	return fmt.Sprintf("Modulation ch=%d value=%d", m.Channel+1, m.Value)
}

// Breath is the breath controller control change (CC2).
type Breath struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewBreath returns a breath controller message.
func NewBreath(channel, value uint8) Breath {
	return Breath{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m Breath) Bytes() []byte {
	return ccMessage(m.Channel, ControllerBreath, m.Value)
}

func (m Breath) String() string {
	return fmt.Sprintf("Breath ch=%d value=%d", m.Channel+1, m.Value)
}

// Foot is the foot controller control change (CC4).
type Foot struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewFoot returns a foot controller message.
func NewFoot(channel, value uint8) Foot {
	return Foot{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m Foot) Bytes() []byte {
	return ccMessage(m.Channel, ControllerFoot, m.Value)
}

func (m Foot) String() string {
	return fmt.Sprintf("Foot ch=%d value=%d", m.Channel+1, m.Value)
}

// PortamentoTime is the portamento time control change (CC5).
type PortamentoTime struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewPortamentoTime returns a portamento time message.
func NewPortamentoTime(channel, value uint8) PortamentoTime {
	return PortamentoTime{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m PortamentoTime) Bytes() []byte {
	return ccMessage(m.Channel, ControllerPortamentoTime, m.Value)
}

func (m PortamentoTime) String() string {
	return fmt.Sprintf("PortamentoTime ch=%d value=%d", m.Channel+1, m.Value)
}

// DataEntry is the data entry control change (CC6, MSB only).
type DataEntry struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewDataEntry returns a data entry message.
func NewDataEntry(channel, value uint8) DataEntry {
	return DataEntry{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m DataEntry) Bytes() []byte {
	return ccMessage(m.Channel, ControllerDataEntryMSB, m.Value)
}

func (m DataEntry) String() string {
	return fmt.Sprintf("DataEntry ch=%d value=%d", m.Channel+1, m.Value)
}

// ChannelVolume is the channel volume control change (CC7).
type ChannelVolume struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewChannelVolume returns a channel volume message.
func NewChannelVolume(channel, value uint8) ChannelVolume {
	return ChannelVolume{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m ChannelVolume) Bytes() []byte {
	return ccMessage(m.Channel, ControllerChannelVolume, m.Value)
}

func (m ChannelVolume) String() string {
	return fmt.Sprintf("ChannelVolume ch=%d value=%d", m.Channel+1, m.Value)
}

// Balance is the balance control change (CC8).
type Balance struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127, 64 = center
}

// NewBalance returns a balance message.
func NewBalance(channel, value uint8) Balance {
	return Balance{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m Balance) Bytes() []byte {
	return ccMessage(m.Channel, ControllerBalance, m.Value)
}

func (m Balance) String() string {
	return fmt.Sprintf("Balance ch=%d value=%d", m.Channel+1, m.Value)
}

// Pan is the pan control change (CC10).
type Pan struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127, 64 = center
}

// NewPan returns a pan message.
func NewPan(channel, value uint8) Pan {
	return Pan{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m Pan) Bytes() []byte {
	return ccMessage(m.Channel, ControllerPan, m.Value)
}

func (m Pan) String() string {
	return fmt.Sprintf("Pan ch=%d value=%d", m.Channel+1, m.Value)
}

// Expression is the expression control change (CC11).
type Expression struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewExpression returns an expression message.
func NewExpression(channel, value uint8) Expression {
	return Expression{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m Expression) Bytes() []byte {
	return ccMessage(m.Channel, ControllerExpression, m.Value)
}

func (m Expression) String() string {
	return fmt.Sprintf("Expression ch=%d value=%d", m.Channel+1, m.Value)
}

// EffectControl1 is the effect control 1 control change (CC12).
type EffectControl1 struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewEffectControl1 returns an effect control 1 message.
func NewEffectControl1(channel, value uint8) EffectControl1 {
	return EffectControl1{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m EffectControl1) Bytes() []byte {
	return ccMessage(m.Channel, ControllerEffectControl1, m.Value)
}

func (m EffectControl1) String() string {
	return fmt.Sprintf("EffectControl1 ch=%d value=%d", m.Channel+1, m.Value)
}

// EffectControl2 is the effect control 2 control change (CC13).
type EffectControl2 struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewEffectControl2 returns an effect control 2 message.
func NewEffectControl2(channel, value uint8) EffectControl2 {
	return EffectControl2{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m EffectControl2) Bytes() []byte {
	return ccMessage(m.Channel, ControllerEffectControl2, m.Value)
}

func (m EffectControl2) String() string {
	return fmt.Sprintf("EffectControl2 ch=%d value=%d", m.Channel+1, m.Value)
}

// GeneralPurpose is a general purpose controller control change (CC16-19).
type GeneralPurpose struct {
	Channel uint8 // 0-15
	Number  uint8 // 1-4 selects CC16-19
	Value   uint8 // 0-127
}

// NewGeneralPurpose returns a general purpose controller message.
func NewGeneralPurpose(channel, number, value uint8) GeneralPurpose {
	return GeneralPurpose{Channel: channel, Number: number, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m GeneralPurpose) Bytes() []byte {
	return ccMessage(m.Channel, ControllerGeneralPurpose1+Controller(m.Number-1), m.Value)
}

func (m GeneralPurpose) String() string {
	return fmt.Sprintf("GeneralPurpose%d ch=%d value=%d", m.Number, m.Channel+1, m.Value)
}

// Hold1 is the damper pedal / hold 1 control change (CC64). It is a
// switch controller.
type Hold1 struct {
	Channel uint8 // 0-15
	Value   SwitchValue
}

// NewHold1 returns a hold 1 (damper pedal) message.
func NewHold1(channel uint8, on bool) Hold1 {
	if on {
		return Hold1{Channel: channel, Value: SwitchOn}
	}
	return Hold1{Channel: channel, Value: SwitchOff}
}

// Bytes serializes the message to its wire representation.
func (m Hold1) Bytes() []byte {
	return ccMessage(m.Channel, ControllerHold1, m.Value.Value())
}

func (m Hold1) String() string {
	return fmt.Sprintf("Hold1 ch=%d %s", m.Channel+1, m.Value)
}

// Portamento is the portamento switch control change (CC65).
type Portamento struct {
	Channel uint8       // 0-15
	Value   SwitchValue // on/off
}

// NewPortamento returns a portamento switch message.
func NewPortamento(channel uint8, on bool) Portamento {
	if on {
		return Portamento{Channel: channel, Value: SwitchOn}
	}
	return Portamento{Channel: channel, Value: SwitchOff}
}

// Bytes serializes the message to its wire representation.
func (m Portamento) Bytes() []byte {
	return ccMessage(m.Channel, ControllerPortamento, m.Value.Value())
}

func (m Portamento) String() string {
	return fmt.Sprintf("Portamento ch=%d %s", m.Channel+1, m.Value)
}

// Sostenuto is the sostenuto pedal control change (CC66).
type Sostenuto struct {
	Channel uint8       // 0-15
	Value   SwitchValue // on/off
}

// NewSostenuto returns a sostenuto message.
func NewSostenuto(channel uint8, on bool) Sostenuto {
	if on {
		return Sostenuto{Channel: channel, Value: SwitchOn}
	}
	return Sostenuto{Channel: channel, Value: SwitchOff}
}

// Bytes serializes the message to its wire representation.
func (m Sostenuto) Bytes() []byte {
	return ccMessage(m.Channel, ControllerSostenuto, m.Value.Value())
}

func (m Sostenuto) String() string {
	return fmt.Sprintf("Sostenuto ch=%d %s", m.Channel+1, m.Value)
}

// Soft is the soft pedal control change (CC67).
type Soft struct {
	Channel uint8       // 0-15
	Value   SwitchValue // on/off
}

// NewSoft returns a soft pedal message.
func NewSoft(channel uint8, on bool) Soft {
	if on {
		return Soft{Channel: channel, Value: SwitchOn}
	}
	return Soft{Channel: channel, Value: SwitchOff}
}

// Bytes serializes the message to its wire representation.
func (m Soft) Bytes() []byte {
	return ccMessage(m.Channel, ControllerSoft, m.Value.Value())
}

func (m Soft) String() string {
	return fmt.Sprintf("Soft ch=%d %s", m.Channel+1, m.Value)
}

// Legato is the legato footswitch control change (CC68).
type Legato struct {
	Channel uint8       // 0-15
	Value   SwitchValue // on/off
}

// NewLegato returns a legato footswitch message.
func NewLegato(channel uint8, on bool) Legato {
	if on {
		return Legato{Channel: channel, Value: SwitchOn}
	}
	return Legato{Channel: channel, Value: SwitchOff}
}

// Bytes serializes the message to its wire representation.
func (m Legato) Bytes() []byte {
	return ccMessage(m.Channel, ControllerLegato, m.Value.Value())
}

func (m Legato) String() string {
	return fmt.Sprintf("Legato ch=%d %s", m.Channel+1, m.Value)
}

// Hold2 is the hold 2 (freeze) control change (CC69).
type Hold2 struct {
	Channel uint8       // 0-15
	Value   SwitchValue // on/off
}

// NewHold2 returns a hold 2 message.
func NewHold2(channel uint8, on bool) Hold2 {
	if on {
		return Hold2{Channel: channel, Value: SwitchOn}
	}
	return Hold2{Channel: channel, Value: SwitchOff}
}

// Bytes serializes the message to its wire representation.
func (m Hold2) Bytes() []byte {
	return ccMessage(m.Channel, ControllerHold2, m.Value.Value())
}

func (m Hold2) String() string {
	return fmt.Sprintf("Hold2 ch=%d %s", m.Channel+1, m.Value)
}

// SoundController is one of the ten sound controllers (CC70-79).
type SoundController struct {
	Channel uint8 // 0-15
	Number  uint8 // 1-10 selects CC70-79
	Value   uint8 // 0-127
}

// NewSoundController returns a sound controller message. Number 1 selects
// CC70 (sound variation), 5 selects CC74 (brightness), and so on.
func NewSoundController(channel, number, value uint8) SoundController {
	return SoundController{Channel: channel, Number: number, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m SoundController) Bytes() []byte {
	return ccMessage(m.Channel, ControllerSoundVariation+Controller(m.Number-1), m.Value)
}

func (m SoundController) String() string {
	return fmt.Sprintf("SoundController%d ch=%d value=%d", m.Number, m.Channel+1, m.Value)
}

// PortamentoControl is the portamento control control change (CC84): the
// note number of the source note for a portamento.
type PortamentoControl struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127, MIDI note number
}

// NewPortamentoControl returns a portamento control message.
func NewPortamentoControl(channel, note uint8) PortamentoControl {
	return PortamentoControl{Channel: channel, Value: note}
}

// Bytes serializes the message to its wire representation.
func (m PortamentoControl) Bytes() []byte {
	return ccMessage(m.Channel, ControllerPortamentoCtrl, m.Value)
}

func (m PortamentoControl) String() string {
	return fmt.Sprintf("PortamentoControl ch=%d note=%d", m.Channel+1, m.Value)
}

// EffectDepth is one of the five effect depth control changes (CC91-95).
// On the JV-1080 CC91 is reverb send and CC93 is chorus send.
type EffectDepth struct {
	Channel uint8 // 0-15
	Number  uint8 // 1-5 selects CC91-95
	Value   uint8 // 0-127
}

// NewEffectDepth returns an effect depth message. Number 1 selects CC91
// (reverb send), 3 selects CC93 (chorus send).
func NewEffectDepth(channel, number, value uint8) EffectDepth {
	return EffectDepth{Channel: channel, Number: number, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m EffectDepth) Bytes() []byte {
	return ccMessage(m.Channel, ControllerEffect1Depth+Controller(m.Number-1), m.Value)
}

func (m EffectDepth) String() string {
	return fmt.Sprintf("Effect%dDepth ch=%d value=%d", m.Number, m.Channel+1, m.Value)
}

// ReverbSend is the reverb send level control change (CC91).
type ReverbSend struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewReverbSend returns a reverb send message.
func NewReverbSend(channel, value uint8) ReverbSend {
	return ReverbSend{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m ReverbSend) Bytes() []byte {
	return ccMessage(m.Channel, ControllerEffect1Depth, m.Value)
}

func (m ReverbSend) String() string {
	return fmt.Sprintf("ReverbSend ch=%d value=%d", m.Channel+1, m.Value)
}

// ChorusSend is the chorus send level control change (CC93).
type ChorusSend struct {
	Channel uint8 // 0-15
	Value   uint8 // 0-127
}

// NewChorusSend returns a chorus send message.
func NewChorusSend(channel, value uint8) ChorusSend {
	return ChorusSend{Channel: channel, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m ChorusSend) Bytes() []byte {
	return ccMessage(m.Channel, ControllerEffect3Depth, m.Value)
}

func (m ChorusSend) String() string {
	return fmt.Sprintf("ChorusSend ch=%d value=%d", m.Channel+1, m.Value)
}

// DataIncrement is the data increment control change (CC96, value ignored).
type DataIncrement struct {
	Channel uint8 // 0-15
}

// NewDataIncrement returns a data increment message.
func NewDataIncrement(channel uint8) DataIncrement {
	return DataIncrement{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m DataIncrement) Bytes() []byte {
	return ccMessage(m.Channel, ControllerDataIncrement, 0)
}

func (m DataIncrement) String() string {
	return fmt.Sprintf("DataIncrement ch=%d", m.Channel+1)
}

// DataDecrement is the data decrement control change (CC97, value ignored).
type DataDecrement struct {
	Channel uint8 // 0-15
}

// NewDataDecrement returns a data decrement message.
func NewDataDecrement(channel uint8) DataDecrement {
	return DataDecrement{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m DataDecrement) Bytes() []byte {
	return ccMessage(m.Channel, ControllerDataDecrement, 0)
}

func (m DataDecrement) String() string {
	return fmt.Sprintf("DataDecrement ch=%d", m.Channel+1)
}

// RPNParameter identifies a registered parameter number (RPN) selected via
// CC101 (MSB) and CC100 (LSB).
type RPNParameter uint16

// Registered parameter numbers (RPN).
const (
	RPNPitchBendSensitivity RPNParameter = 0x0000
	RPNFineTune             RPNParameter = 0x0001
	RPNCoarseTune           RPNParameter = 0x0002
	RPMTuningProgramChange  RPNParameter = 0x0003
	RPMTuningBankSelect     RPNParameter = 0x0004
	RPNModulationDepthRange RPNParameter = 0x0005
)

// RPNSelect is an RPN parameter selection (CC101/CC100).
type RPNSelect struct {
	Channel   uint8 // 0-15
	Parameter RPNParameter
}

// NewRPNSelect returns an RPN selection message.
func NewRPNSelect(channel uint8, parameter RPNParameter) RPNSelect {
	return RPNSelect{Channel: channel, Parameter: parameter}
}

// Bytes serializes the message to its wire representation: the CC101 (MSB)
// message followed by the CC100 (LSB) message.
func (m RPNSelect) Bytes() []byte {
	return append(ccMessage(m.Channel, ControllerRPNSelectMSB, uint8(m.Parameter>>7)),
		ccMessage(m.Channel, ControllerRPNSelectLSB, uint8(m.Parameter&0x7F))...)
}

func (m RPNSelect) String() string {
	return fmt.Sprintf("RPNSelect ch=%d parameter=%04X", m.Channel+1, uint16(m.Parameter))
}

// NRPNSelect is a non-registered parameter number selection (CC99/CC98).
type NRPNSelect struct {
	Channel uint8 // 0-15
	MSB     uint8 // 0-127
	LSB     uint8 // 0-127
}

// NewNRPNSelect returns an NRPN selection message.
func NewNRPNSelect(channel, msb, lsb uint8) NRPNSelect {
	return NRPNSelect{Channel: channel, MSB: msb, LSB: lsb}
}

// Bytes serializes the message to its wire representation: the CC99 (MSB)
// message followed by the CC98 (LSB) message.
func (m NRPNSelect) Bytes() []byte {
	return append(ccMessage(m.Channel, ControllerNRPNSelectMSB, m.MSB),
		ccMessage(m.Channel, ControllerNRPNSelectLSB, m.LSB)...)
}

func (m NRPNSelect) String() string {
	return fmt.Sprintf("NRPNSelect ch=%d msb=%d lsb=%d", m.Channel+1, m.MSB, m.LSB)
}

// AllSoundOff is the all sound off channel mode message (CC120).
type AllSoundOff struct {
	Channel uint8 // 0-15
}

// NewAllSoundOff returns an all sound off message.
func NewAllSoundOff(channel uint8) AllSoundOff {
	return AllSoundOff{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m AllSoundOff) Bytes() []byte {
	return ccMessage(m.Channel, ControllerAllSoundOff, 0)
}

func (m AllSoundOff) String() string {
	return fmt.Sprintf("AllSoundOff ch=%d", m.Channel+1)
}

// ResetAllControllers is the reset all controllers channel mode message
// (CC121).
type ResetAllControllers struct {
	Channel uint8 // 0-15
}

// NewResetAllControllers returns a reset all controllers message.
func NewResetAllControllers(channel uint8) ResetAllControllers {
	return ResetAllControllers{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m ResetAllControllers) Bytes() []byte {
	return ccMessage(m.Channel, ControllerResetAllCtrl, 0)
}

func (m ResetAllControllers) String() string {
	return fmt.Sprintf("ResetAllControllers ch=%d", m.Channel+1)
}

// LocalControl is the local control channel mode message (CC122).
type LocalControl struct {
	Channel uint8       // 0-15
	Value   SwitchValue // on/off
}

// NewLocalControl returns a local control message.
func NewLocalControl(channel uint8, on bool) LocalControl {
	if on {
		return LocalControl{Channel: channel, Value: SwitchOn}
	}
	return LocalControl{Channel: channel, Value: SwitchOff}
}

// Bytes serializes the message to its wire representation.
func (m LocalControl) Bytes() []byte {
	return ccMessage(m.Channel, ControllerLocalControl, m.Value.Value())
}

func (m LocalControl) String() string {
	return fmt.Sprintf("LocalControl ch=%d %s", m.Channel+1, m.Value)
}

// AllNotesOff is the all notes off channel mode message (CC123).
type AllNotesOff struct {
	Channel uint8 // 0-15
}

// NewAllNotesOff returns an all notes off message.
func NewAllNotesOff(channel uint8) AllNotesOff {
	return AllNotesOff{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m AllNotesOff) Bytes() []byte {
	return ccMessage(m.Channel, ControllerAllNotesOff, 0)
}

func (m AllNotesOff) String() string {
	return fmt.Sprintf("AllNotesOff ch=%d", m.Channel+1)
}

// OmniModeOff is the omni mode off channel mode message (CC124).
type OmniModeOff struct {
	Channel uint8 // 0-15
}

// NewOmniModeOff returns an omni mode off message.
func NewOmniModeOff(channel uint8) OmniModeOff {
	return OmniModeOff{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m OmniModeOff) Bytes() []byte {
	return ccMessage(m.Channel, ControllerOmniModeOff, 0)
}

func (m OmniModeOff) String() string {
	return fmt.Sprintf("OmniModeOff ch=%d", m.Channel+1)
}

// OmniModeOn is the omni mode on channel mode message (CC125).
type OmniModeOn struct {
	Channel uint8 // 0-15
}

// NewOmniModeOn returns an omni mode on message.
func NewOmniModeOn(channel uint8) OmniModeOn {
	return OmniModeOn{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m OmniModeOn) Bytes() []byte {
	return ccMessage(m.Channel, ControllerOmniModeOn, 0)
}

func (m OmniModeOn) String() string {
	return fmt.Sprintf("OmniModeOn ch=%d", m.Channel+1)
}

// MonoModeOn is the mono mode on channel mode message (CC126).
type MonoModeOn struct {
	Channel uint8 // 0-15
}

// NewMonoModeOn returns a mono mode on message.
func NewMonoModeOn(channel uint8) MonoModeOn {
	return MonoModeOn{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m MonoModeOn) Bytes() []byte {
	return ccMessage(m.Channel, ControllerMonoModeOn, 0)
}

func (m MonoModeOn) String() string {
	return fmt.Sprintf("MonoModeOn ch=%d", m.Channel+1)
}

// PolyModeOn is the poly mode on channel mode message (CC127).
type PolyModeOn struct {
	Channel uint8 // 0-15
}

// NewPolyModeOn returns a poly mode on message.
func NewPolyModeOn(channel uint8) PolyModeOn {
	return PolyModeOn{Channel: channel}
}

// Bytes serializes the message to its wire representation.
func (m PolyModeOn) Bytes() []byte {
	return ccMessage(m.Channel, ControllerPolyModeOn, 0)
}

func (m PolyModeOn) String() string {
	return fmt.Sprintf("PolyModeOn ch=%d", m.Channel+1)
}
