package message

import "fmt"

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
	b := []byte{SysExStart, RolandManufacturerID, m.DeviceID & 0x7F, JV1080ModelID, CommandDT1}
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
	if b[1] != RolandManufacturerID || b[3] != JV1080ModelID || b[4] != CommandDT1 {
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
