package message

// ParseRoland parses a raw MIDI message into the matching struct. It
// returns an error for messages that are not Roland JV RQ1/DT1 messages.
func ParseRoland(b []byte) (Message, error) {
	if len(b) >= 5 && b[0] == SysExStart && b[1] == RolandManufacturerID && b[3] == JV1080ModelID {
		switch b[4] {
		case CommandRQ1:
			m, err := ParseDataRequest(b)
			return m, err
		case CommandDT1:
			m, err := ParseDataSet(b)
			return m, err
		}
	}
	return nil, ErrNotRolandJV
}
