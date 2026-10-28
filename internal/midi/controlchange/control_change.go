// Package controlchange provides the control change message and typed
// structs for the known control change controllers.
package controlchange

import (
	"fmt"

	"go-jv-1080/internal/midi/message"
)

// ControlChange is a channel voice control change message (Bn cc vv).
type ControlChange struct {
	Channel    uint8      // 0-15
	Controller Controller // 0-127
	Value      uint8      // 0-127
}

// NewControlChange returns a control change message.
func NewControlChange(channel uint8, controller Controller, value uint8) ControlChange {
	return ControlChange{Channel: channel, Controller: controller, Value: value}
}

// Bytes serializes the message to its wire representation.
func (m ControlChange) Bytes() []byte {
	return []byte{0xB0 | (m.Channel & 0x0F), uint8(m.Controller) & 0x7F, m.Value & 0x7F}
}

func (m ControlChange) String() string {
	return fmt.Sprintf("CC ch=%d controller=%s value=%d", m.Channel+1, m.Controller, m.Value)
}

// ParseControlChange parses a control change message from its wire bytes.
func ParseControlChange(b []byte) (ControlChange, error) {
	if len(b) < 3 {
		return ControlChange{}, message.ErrShortMessage
	}
	if b[0]&0xF0 != 0xB0 {
		return ControlChange{}, fmt.Errorf("midi: not a control change message: % X", b)
	}
	return ControlChange{Channel: b[0] & 0x0F, Controller: Controller(b[1]), Value: b[2]}, nil
}

// ccMessage builds the wire bytes of a single control change message.
func ccMessage(channel uint8, controller Controller, value uint8) []byte {
	return ControlChange{Channel: channel, Controller: controller, Value: value}.Bytes()
}
