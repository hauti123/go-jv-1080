package message

import "fmt"

// DataRequest (RQ1) asks the device to transmit the data at Address, Size
// bytes of it, as a DataSet message.
type DataRequest struct {
	DeviceID uint8
	Address  Address
	Size     int
}

// NewDataRequest returns a data request for Size bytes at Address.
func NewDataRequest(deviceID uint8, address Address, size int) DataRequest {
	return DataRequest{DeviceID: deviceID, Address: address, Size: size}
}

// Bytes serializes the message to its wire representation.
func (m DataRequest) Bytes() []byte {
	size := AddressFromInt(m.Size)
	body := append([]byte{}, m.Address[:]...)
	body = append(body, size[:]...)
	b := []byte{SysExStart, RolandManufacturerID, m.DeviceID & 0x7F, JV1080ModelID, CommandRQ1}
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
	if b[1] != RolandManufacturerID || b[3] != JV1080ModelID || b[4] != CommandRQ1 {
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
