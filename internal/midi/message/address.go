package message

import "fmt"

// Address is a Roland 4-byte parameter address (aa bb cc dd), each byte
// holding 7 bits.
type Address [4]byte

// NewAddress returns a Roland 4-byte address from its four bytes.
func NewAddress(b0, b1, b2, b3 uint8) Address {
	return Address{b0, b1, b2, b3}
}

// Addr returns the address as a single number.
func (a Address) Addr() int {
	return ((int(a[0])*128+int(a[1]))*128+int(a[2]))*128 + int(a[3])
}

// AddressFromInt converts a number into a Roland 4-byte address.
func AddressFromInt(n int) Address {
	return Address{byte(n >> 21 & 0x7F), byte(n >> 14 & 0x7F), byte(n >> 7 & 0x7F), byte(n & 0x7F)}
}

// Add returns the address offset by n bytes.
func (a Address) Add(n int) Address {
	return AddressFromInt(a.Addr() + n)
}

func (a Address) String() string {
	return fmt.Sprintf("%02X %02X %02X %02X", a[0], a[1], a[2], a[3])
}

// Checksum calculates the Roland checksum over the given address and
// data (or size) bytes: the low 7 bits of the sum plus the checksum must
// be zero.
func Checksum(b []byte) byte {
	var sum int
	for _, v := range b {
		sum += int(v & 0x7F)
	}
	return byte((128 - sum%128) % 128)
}
