package message

import "fmt"

// NoteOn is a channel voice note on message (9n kk vv). A note off is
// sent as a note on with velocity 0.
type NoteOn struct {
	Channel  uint8 // 0-15
	Note     uint8 // 0-127
	Velocity uint8 // 0-127
}

// Bytes serializes the message to its wire representation.
func (m NoteOn) Bytes() []byte {
	return []byte{0x90 | (m.Channel & 0x0F), m.Note & 0x7F, m.Velocity & 0x7F}
}

func (m NoteOn) String() string {
	return fmt.Sprintf("NoteOn ch=%d note=%d velocity=%d", m.Channel+1, m.Note, m.Velocity)
}

// ParseNoteOn parses a note on message from its wire bytes.
func ParseNoteOn(b []byte) (NoteOn, error) {
	if len(b) < 3 {
		return NoteOn{}, ErrShortMessage
	}
	if b[0]&0xF0 != 0x90 {
		return NoteOn{}, fmt.Errorf("midi: not a note on message: % X", b)
	}
	return NoteOn{Channel: b[0] & 0x0F, Note: b[1], Velocity: b[2]}, nil
}
