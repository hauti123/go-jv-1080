package controlchange

// SwitchValue is the value of a MIDI switch controller (pedals, portamento
// switch, ...): 0-63 = off, 64-127 = on.
type SwitchValue uint8

// Switch values.
const (
	SwitchOff SwitchValue = 0
	SwitchOn  SwitchValue = 127
)

// OnOff returns the switch value as a bool.
func (s SwitchValue) OnOff() bool { return s >= 64 }

// Value returns the switch value as a plain byte (0 or 127).
func (s SwitchValue) Value() uint8 { return uint8(s) }

func (s SwitchValue) String() string {
	if s.OnOff() {
		return "on"
	}
	return "off"
}
