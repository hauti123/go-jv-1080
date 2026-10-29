package message

// SplitMessages splits a byte stream of concatenated MIDI messages into
// the individual messages. SysEx messages (F0 ... F7), channel voice
// messages (status byte plus one or two data bytes) and system common and
// realtime messages are recognized. The stream must use explicit status
// bytes (no running status). Stray data bytes without a status byte are
// skipped.
func SplitMessages(b []byte) [][]byte {
	var msgs [][]byte
	for i := 0; i < len(b); {
		start := i
		status := b[i]
		switch {
		case status == SysExStart:
			// SysEx message: up to and including EOX.
			i++
			for i < len(b) && b[i] != SysExEnd {
				i++
			}
			if i < len(b) {
				i++
			}
		case status >= 0xF8:
			// Realtime message: no data bytes.
			i++
		case status >= 0xF0:
			// System common message.
			i++
			switch status {
			case 0xF1, 0xF3:
				i++
			case 0xF2:
				i += 2
			}
			if i > len(b) {
				i = len(b)
			}
		case status >= 0x80:
			// Channel voice message: program change and channel
			// pressure have one data byte, the others two.
			i++
			if status&0xF0 == 0xC0 || status&0xF0 == 0xD0 {
				i++
			} else {
				i += 2
			}
			if i > len(b) {
				i = len(b)
			}
		default:
			// Stray data byte without a status byte.
			i++
			continue
		}
		msgs = append(msgs, b[start:i])
	}
	return msgs
}
