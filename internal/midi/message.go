// Package midi provides easy-to-understand structs for the individual MIDI
// messages used to communicate with the Roland JV-1080, together with
// serialization from and to raw byte streams.
package midi

import (
	"errors"
	"fmt"
)

const (
	// SysExStart and SysExEnd delimit a system exclusive message on the wire.
	SysExStart = 0xF0
	SysExEnd   = 0xF7

	// RolandManufacturerID is Roland's MIDI manufacturer ID.
	RolandManufacturerID = 0x41

	// JV1080ModelID is the model ID of the JV-1080 (also shared by the
	// JV-2080, JV-1010 and XP series).
	JV1080ModelID = 0x6A

	// DefaultDeviceID is the factory device ID of the JV-1080 (unit 17).
	DefaultDeviceID = 0x10

	commandRQ1 = 0x11 // Data Request 1
	commandDT1 = 0x12 // Data Set 1
)

// ErrShortMessage is returned when a byte stream is too short to hold the
// message it is parsed into.
var ErrShortMessage = errors.New("midi: message too short")

// ErrBadChecksum is returned when the checksum of a parsed RQ1/DT1 message
// does not match its address and data/size bytes.
var ErrBadChecksum = errors.New("midi: bad checksum")

// ErrNotRolandJV is returned when a message is not a Roland JV (model 0x6A)
// RQ1 or DT1 message.
var ErrNotRolandJV = errors.New("midi: not a Roland JV RQ1/DT1 message")

// Message is any MIDI message that can be serialized to and parsed from a
// raw byte stream.
type Message interface {
	Bytes() []byte
	String() string
}

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

// ProgramChange is a channel voice program change message (Cn pp).
type ProgramChange struct {
	Channel uint8 // 0-15
	Program uint8 // 0-127
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

// NoteOn is a channel voice note on message (9n kk vv).
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

// IdentityRequest is the universal non-realtime identity request
// (F0 7E dev 06 01 F7).
type IdentityRequest struct {
	DeviceID uint8 // 0x10 for the JV-1080, 0x7F for broadcast
}

// Bytes serializes the message to its wire representation.
func (m IdentityRequest) Bytes() []byte {
	return []byte{SysExStart, 0x7E, m.DeviceID & 0x7F, 0x06, 0x01, SysExEnd}
}

func (m IdentityRequest) String() string {
	return fmt.Sprintf("IdentityRequest dev=%02X", m.DeviceID)
}

// IdentityReply is the universal non-realtime identity reply
// (F0 7E dev 06 02 41 6A famMSB famLSB numMSB numLSB rev[4] F7) sent by the
// device in answer to an IdentityRequest. FamilyCode and FamilyNumber are
// each transmitted as two bytes, MSB first.
type IdentityReply struct {
	DeviceID     uint8
	Manufacturer uint8
	FamilyCode   uint16
	FamilyNumber uint16
	Revision     [4]uint8
}

// Bytes serializes the message to its wire representation.
func (m IdentityReply) Bytes() []byte {
	b := []byte{SysExStart, 0x7E, m.DeviceID & 0x7F, 0x06, 0x02, m.Manufacturer & 0x7F,
		byte(m.FamilyCode), byte(m.FamilyCode >> 8), byte(m.FamilyNumber), byte(m.FamilyNumber >> 8)}
	b = append(b, m.Revision[:]...)
	return append(b, 0x00, SysExEnd)
}

func (m IdentityReply) String() string {
	return fmt.Sprintf("IdentityReply dev=%02X manufacturer=%02X family=%04X number=%04X rev=%v",
		m.DeviceID, m.Manufacturer, m.FamilyCode, m.FamilyNumber, m.Revision)
}

// ParseIdentityReply parses an identity reply from its wire bytes. It
// returns ErrNotRolandJV for replies of other manufacturers and
// ErrShortMessage for truncated messages.
func ParseIdentityReply(b []byte) (IdentityReply, error) {
	if len(b) < 16 {
		return IdentityReply{}, ErrShortMessage
	}
	if b[0] != SysExStart || b[len(b)-1] != SysExEnd {
		return IdentityReply{}, fmt.Errorf("midi: not a sysex message: % X", b)
	}
	if b[1] != 0x7E || b[3] != 0x06 || b[4] != 0x02 {
		return IdentityReply{}, fmt.Errorf("midi: not an identity reply: % X", b)
	}
	var m IdentityReply
	m.DeviceID = b[2]
	m.Manufacturer = b[5]
	if m.Manufacturer != RolandManufacturerID {
		return m, ErrNotRolandJV
	}
	m.FamilyCode = uint16(b[6]) | uint16(b[7])<<8
	m.FamilyNumber = uint16(b[8]) | uint16(b[9])<<8
	copy(m.Revision[:], b[10:14])
	return m, nil
}

// Address is a Roland 4-byte parameter address (aa bb cc dd), each byte
// holding 7 bits.
type Address [4]byte

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

// DataRequest (RQ1) asks the device to transmit the data at Address, Size
// bytes of it, as a DataSet message.
type DataRequest struct {
	DeviceID uint8
	Address  Address
	Size     int
}

// Bytes serializes the message to its wire representation.
func (m DataRequest) Bytes() []byte {
	size := AddressFromInt(m.Size)
	body := append([]byte{}, m.Address[:]...)
	body = append(body, size[:]...)
	b := []byte{SysExStart, RolandManufacturerID, m.DeviceID & 0x7F, JV1080ModelID, commandRQ1}
	b = append(b, body...)
	b = append(b, Checksum(body), SysExEnd)
	return b
}

func (m DataRequest) String() string {
	return fmt.Sprintf("RQ1 dev=%02X address=%s size=%d", m.DeviceID, m.Address, m.Size)
}

// ParseDataRequest parses an RQ1 message from its wire bytes.
func ParseDataRequest(b []byte) (DataRequest, error) {
	if len(b) < 14 {
		return DataRequest{}, ErrShortMessage
	}
	if b[0] != SysExStart || b[len(b)-1] != SysExEnd {
		return DataRequest{}, fmt.Errorf("midi: not a sysex message: % X", b)
	}
	if b[1] != RolandManufacturerID || b[3] != JV1080ModelID || b[4] != commandRQ1 {
		return DataRequest{}, ErrNotRolandJV
	}
	var m DataRequest
	m.DeviceID = b[2]
	copy(m.Address[:], b[5:9])
	size := Address{b[9], b[10], b[11], b[12]}
	m.Size = size.Addr()
	if Checksum(b[5:13]) != b[13] {
		return m, ErrBadChecksum
	}
	return m, nil
}

// DataSet (DT1) transmits data to be written at Address on the device.
type DataSet struct {
	DeviceID uint8
	Address  Address
	Data     []byte
}

// Bytes serializes the message to its wire representation.
func (m DataSet) Bytes() []byte {
	body := append([]byte{}, m.Address[:]...)
	body = append(body, m.Data...)
	b := []byte{SysExStart, RolandManufacturerID, m.DeviceID & 0x7F, JV1080ModelID, commandDT1}
	b = append(b, body...)
	b = append(b, Checksum(body), SysExEnd)
	return b
}

func (m DataSet) String() string {
	return fmt.Sprintf("DT1 dev=%02X address=%s data=% X", m.DeviceID, m.Address, m.Data)
}

// ParseDataSet parses a DT1 message from its wire bytes.
func ParseDataSet(b []byte) (DataSet, error) {
	if len(b) < 12 {
		return DataSet{}, ErrShortMessage
	}
	if b[0] != SysExStart || b[len(b)-1] != SysExEnd {
		return DataSet{}, fmt.Errorf("midi: not a sysex message: % X", b)
	}
	if b[1] != RolandManufacturerID || b[3] != JV1080ModelID || b[4] != commandDT1 {
		return DataSet{}, ErrNotRolandJV
	}
	var m DataSet
	m.DeviceID = b[2]
	copy(m.Address[:], b[5:9])
	m.Data = append([]byte{}, b[9:len(b)-2]...)
	if Checksum(b[5:len(b)-2]) != b[len(b)-2] {
		return m, ErrBadChecksum
	}
	return m, nil
}

// ParseRoland parses a raw MIDI message into the matching struct. It
// returns an error for messages that are not Roland JV RQ1/DT1 messages.
func ParseRoland(b []byte) (Message, error) {
	if len(b) >= 5 && b[0] == SysExStart && b[1] == RolandManufacturerID && b[3] == JV1080ModelID {
		switch b[4] {
		case commandRQ1:
			m, err := ParseDataRequest(b)
			return m, err
		case commandDT1:
			m, err := ParseDataSet(b)
			return m, err
		}
	}
	return nil, ErrNotRolandJV
}

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
