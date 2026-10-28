package message

import "fmt"

// ProgramChange is a channel voice program change message (Cn pp).
type ProgramChange struct {
	Channel uint8 // 0-15
	Program uint8 // 0-127
}

// NewProgramChange returns a program change message.
func NewProgramChange(channel, program uint8) ProgramChange {
	return ProgramChange{Channel: channel, Program: program}
}

// Bytes serializes the message to its wire representation.
func (m ProgramChange) Bytes() []byte {
	return []byte{0xC0 | (m.Channel & 0x0F), m.Program & 0x7F}
}

func (m ProgramChange) String() string {
	return fmt.Sprintf("PC ch=%d program=%d", m.Channel+1, m.Program+1)
}

// ParseProgramChange parses a program change message from its wire bytes.
func ParseProgramChange(b []byte) (ProgramChange, error) {
	if len(b) < 2 {
		return ProgramChange{}, ErrShortMessage
	}
	if b[0]&0xF0 != 0xC0 {
		return ProgramChange{}, fmt.Errorf("midi: not a program change message: % X", b)
	}
	return ProgramChange{Channel: b[0] & 0x0F, Program: b[1]}, nil
}
