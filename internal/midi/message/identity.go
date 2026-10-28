package message

import "fmt"

// IdentityRequest is the universal non-realtime identity request
// (F0 7E dev 06 01 F7).
type IdentityRequest struct {
	DeviceID uint8 // 0x10 for the JV-1080, 0x7F for broadcast
}

// NewIdentityRequest returns an identity request for the given device.
func NewIdentityRequest(deviceID uint8) IdentityRequest {
	return IdentityRequest{DeviceID: deviceID}
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

// NewIdentityReply returns an identity reply.
func NewIdentityReply(deviceID, manufacturer uint8, familyCode, familyNumber uint16, revision [4]uint8) IdentityReply {
	return IdentityReply{DeviceID: deviceID, Manufacturer: manufacturer, FamilyCode: familyCode, FamilyNumber: familyNumber, Revision: revision}
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
