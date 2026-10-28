package message

import "fmt"

// ControlChange is a channel voice control change message (Bn cc vv).
type ControlChange struct {
	Channel    uint8 // 0-15
	Controller uint8 // 0-127
	Value      uint8 // 0-127
}

// Bytes serializes the message to its wire representation.
func (m ControlChange) Bytes() []byte {
	return []byte{0xB0 | (m.Channel & 0x0F), m.Controller & 0x7F, m.Value & 0x7F}
}

func (m ControlChange) String() string {
	return fmt.Sprintf("CC ch=%d controller=%d value=%d", m.Channel+1, m.Controller, m.Value)
}

// ParseControlChange parses a control change message from its wire bytes.
func ParseControlChange(b []byte) (ControlChange, error) {
	if len(b) < 3 {
		return ControlChange{}, ErrShortMessage
	}
	if b[0]&0xF0 != 0xB0 {
		return ControlChange{}, fmt.Errorf("midi: not a control change message: % X", b)
	}
	return ControlChange{Channel: b[0] & 0x0F, Controller: b[1], Value: b[2]}, nil
}
