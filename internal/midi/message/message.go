// Package message provides easy-to-understand structs for the individual
// MIDI messages used to communicate with the Roland JV-1080, together with
// serialization from and to raw byte streams.
package message

import (
	"errors"
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

	CommandRQ1 = 0x11 // Data Request 1
	CommandDT1 = 0x12 // Data Set 1
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
