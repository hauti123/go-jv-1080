package message

import "fmt"

// Hex formats a byte slice as upper-case hex, e.g. "F0 41 10 6A".
func Hex(b []byte) string {
	const digits = "0123456789ABCDEF"
	out := make([]byte, 0, len(b)*3)
	for i, v := range b {
		if i > 0 {
			out = append(out, ' ')
		}
		out = append(out, digits[v>>4], digits[v&0x0F])
	}
	return string(out)
}

// ParseHex parses a hex string such as "F0 41 10" or "f0-41-10" into
// bytes. Whitespace, commas, dashes and "0x" prefixes are accepted as
// separators.
func ParseHex(s string) ([]byte, error) {
	var out []byte
	clean := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == ' ' || c == ',' || c == '-' || c == '\t' || c == '\n' || c == '\r':
		case c == 'x' || c == 'X':
		default:
			clean = append(clean, c)
		}
	}
	if len(clean)%2 != 0 {
		return nil, fmt.Errorf("midi: odd number of hex digits in %q", s)
	}
	for i := 0; i < len(clean); i += 2 {
		hi, ok1 := hexVal(clean[i])
		lo, ok2 := hexVal(clean[i+1])
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("midi: invalid hex digit in %q", s)
		}
		out = append(out, hi<<4|lo)
	}
	return out, nil
}

func hexVal(c byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}
