package jv1080

import (
	"fmt"
	"time"

	"go-jv-1080/internal/midi"
)

// DefaultTimeout is the default time to wait for the device to answer a
// data request.
const DefaultTimeout = 2 * time.Second

// Client is a high level API to control a JV-1080 over a MIDI device
// connection. It supports reading full dumps of the user banks, reading
// single patches and performances, updating single parameters and
// triggering the storing of a changed user patch or performance.
type Client struct {
	dev      *midi.Device
	deviceID uint8
	timeout  time.Duration
}

// NewClient returns a client talking to the device with the given device
// ID (0x10 = factory default, unit 17).
func NewClient(dev *midi.Device, deviceID uint8) *Client {
	return &Client{dev: dev, deviceID: deviceID, timeout: DefaultTimeout}
}

// SetTimeout sets the timeout for data requests.
func (c *Client) SetTimeout(d time.Duration) { c.timeout = d }

// DeviceID returns the device ID the client talks to.
func (c *Client) DeviceID() uint8 { return c.deviceID }

// request sends an RQ1 and waits for the DT1 answer.
func (c *Client) request(addr midi.Address, size int) (midi.DataSet, error) {
	req := midi.DataRequest{DeviceID: c.deviceID, Address: addr, Size: size}
	return c.dev.Request(req, c.timeout, func(ds midi.DataSet) bool {
		return len(ds.Data) == size
	})
}

// sendDataSet writes a data block to the device.
func (c *Client) sendDataSet(addr midi.Address, data []byte) error {
	return c.dev.SendMessage(midi.DataSet{DeviceID: c.deviceID, Address: addr, Data: data})
}

// Identity sends an identity request and returns the device's reply.
func (c *Client) Identity() (midi.IdentityReply, error) {
	return c.dev.RequestIdentity(c.deviceID, c.timeout)
}

// SelectPatch sends bank select (CC0/CC32) and program change to select a
// patch on the given channel (0-15).
func (c *Client) SelectPatch(channel uint8, bank Bank, number int) error {
	msb, lsb := bank.bankSelect()
	if err := c.dev.SendMessage(midi.ControlChange{Channel: channel, Controller: 0, Value: msb}); err != nil {
		return err
	}
	if err := c.dev.SendMessage(midi.ControlChange{Channel: channel, Controller: 32, Value: lsb}); err != nil {
		return err
	}
	return c.dev.SendMessage(midi.ProgramChange{Channel: channel, Program: uint8(clamp(number, 1, 128) - 1)})
}

// SelectPerformance sends bank select and program change to select a
// performance on the control channel (0-15).
func (c *Client) SelectPerformance(channel uint8, bank Bank, number int) error {
	msb, lsb := bank.performanceBankSelect()
	if err := c.dev.SendMessage(midi.ControlChange{Channel: channel, Controller: 0, Value: msb}); err != nil {
		return err
	}
	if err := c.dev.SendMessage(midi.ControlChange{Channel: channel, Controller: 32, Value: lsb}); err != nil {
		return err
	}
	return c.dev.SendMessage(midi.ProgramChange{Channel: channel, Program: uint8(clamp(number, 1, 32) - 1)})
}

// ReadPatch reads a single patch into the patch-mode temporary patch
// buffer and returns it. The number is 1-128.
func (c *Client) ReadPatch(bank Bank, number int) (*Patch, error) {
	if err := c.SelectPatch(0, bank, number); err != nil {
		return nil, err
	}
	time.Sleep(50 * time.Millisecond)
	return c.ReadTemporaryPatch()
}

// ReadTemporaryPatch reads the patch-mode temporary patch (common block
// plus four tones).
func (c *Client) ReadTemporaryPatch() (*Patch, error) {
	area := TemporaryPatchArea()
	p := NewPatch(area)
	addr, size := patchCommonBlock(area)
	ds, err := c.request(addr, size)
	if err != nil {
		return nil, fmt.Errorf("jv1080: read patch common: %w", err)
	}
	copy(p.Common, ds.Data)
	for t := 1; t <= 4; t++ {
		addr, size := patchToneBlock(area, t)
		ds, err := c.request(addr, size)
		if err != nil {
			return nil, fmt.Errorf("jv1080: read patch tone %d: %w", t, err)
		}
		copy(p.Tones[t-1], ds.Data)
	}
	return p, nil
}

// ReadUserPatch reads a user patch (number 1-128) directly from the user
// memory area.
func (c *Client) ReadUserPatch(number int) (*Patch, error) {
	area := userPatchArea(number)
	p := NewPatch(area)
	addr, size := patchCommonBlock(area)
	ds, err := c.request(addr, size)
	if err != nil {
		return nil, fmt.Errorf("jv1080: read user patch %d common: %w", number, err)
	}
	copy(p.Common, ds.Data)
	for t := 1; t <= 4; t++ {
		addr, size := patchToneBlock(area, t)
		ds, err := c.request(addr, size)
		if err != nil {
			return nil, fmt.Errorf("jv1080: read user patch %d tone %d: %w", number, t, err)
		}
		copy(p.Tones[t-1], ds.Data)
	}
	return p, nil
}

// DumpUserPatchBank reads all 128 user patches (the full user patch bank).
func (c *Client) DumpUserPatchBank() ([]*Patch, error) {
	patches := make([]*Patch, 0, 128)
	for n := 1; n <= 128; n++ {
		p, err := c.ReadUserPatch(n)
		if err != nil {
			return nil, err
		}
		patches = append(patches, p)
	}
	return patches, nil
}

// ReadUserPerformance reads a user performance (number 1-32) directly
// from the user memory area.
func (c *Client) ReadUserPerformance(number int) (*Performance, error) {
	area := userPerformanceArea(number)
	perf := NewPerformance(area)
	addr, size := performanceCommonBlock(area)
	ds, err := c.request(addr, size)
	if err != nil {
		return nil, fmt.Errorf("jv1080: read user performance %d common: %w", number, err)
	}
	copy(perf.Common, ds.Data)
	for n := 1; n <= 16; n++ {
		addr, size := performancePartBlock(area, n)
		ds, err := c.request(addr, size)
		if err != nil {
			return nil, fmt.Errorf("jv1080: read user performance %d part %d: %w", number, n, err)
		}
		copy(perf.Parts[n-1], ds.Data)
	}
	return perf, nil
}

// ReadTemporaryPerformance reads the temporary performance.
func (c *Client) ReadTemporaryPerformance() (*Performance, error) {
	area := TemporaryPerformanceArea()
	perf := NewPerformance(area)
	addr, size := performanceCommonBlock(area)
	ds, err := c.request(addr, size)
	if err != nil {
		return nil, fmt.Errorf("jv1080: read temporary performance common: %w", err)
	}
	copy(perf.Common, ds.Data)
	for n := 1; n <= 16; n++ {
		addr, size := performancePartBlock(area, n)
		ds, err := c.request(addr, size)
		if err != nil {
			return nil, fmt.Errorf("jv1080: read temporary performance part %d: %w", n, err)
		}
		copy(perf.Parts[n-1], ds.Data)
	}
	return perf, nil
}

// DumpUserPerformanceBank reads all 32 user performances.
func (c *Client) DumpUserPerformanceBank() ([]*Performance, error) {
	perfs := make([]*Performance, 0, 32)
	for n := 1; n <= 32; n++ {
		p, err := c.ReadUserPerformance(n)
		if err != nil {
			return nil, err
		}
		perfs = append(perfs, p)
	}
	return perfs, nil
}

// WritePatchToTemporary writes a whole patch to the patch-mode temporary
// patch buffer. Blocks larger than 128 bytes are split into packets of at
// most 128 bytes with a 20 ms pause between them, as the MIDI
// implementation requires.
func (c *Client) WritePatchToTemporary(p *Patch) error {
	if err := c.writeBlock(TemporaryPatchArea(), p.Common); err != nil {
		return err
	}
	for t := 1; t <= 4; t++ {
		if err := c.writeBlock(patchToneArea(TemporaryPatchArea(), t), p.Tones[t-1]); err != nil {
			return err
		}
	}
	return nil
}

// WritePerformanceToTemporary writes a whole performance to the temporary
// performance buffer.
func (c *Client) WritePerformanceToTemporary(p *Performance) error {
	if err := c.writeBlock(TemporaryPerformanceArea(), p.Common); err != nil {
		return err
	}
	for n := 1; n <= 16; n++ {
		addr, _ := performancePartBlock(TemporaryPerformanceArea(), n)
		if err := c.writeBlock(addr, p.Parts[n-1]); err != nil {
			return err
		}
	}
	return nil
}

// maxDataSetBytes is the maximum number of data bytes per DT1 message.
const maxDataSetBytes = 128

// writeBlock writes a block of data, splitting it into DT1 packets of at
// most 128 data bytes with a 20 ms interval between them.
func (c *Client) writeBlock(addr midi.Address, data []byte) error {
	for off := 0; off < len(data); off += maxDataSetBytes {
		end := off + maxDataSetBytes
		if end > len(data) {
			end = len(data)
		}
		if err := c.sendDataSet(addr.Add(off), data[off:end]); err != nil {
			return err
		}
		time.Sleep(20 * time.Millisecond)
	}
	return nil
}

// StorePatch triggers storing the patch-mode temporary patch into the
// given user patch number (1-128). The JV-1080 stores the temporary patch
// when it receives a data set for the user patch area whose address and
// size match the whole patch; the recommended way is to write the patch to
// the user area directly.
func (c *Client) StorePatch(number int, p *Patch) error {
	if number < 1 || number > 128 {
		return fmt.Errorf("jv1080: user patch number %d out of range 1-128", number)
	}
	area := userPatchArea(number)
	if err := c.writeBlock(area, p.Common); err != nil {
		return err
	}
	for t := 1; t <= 4; t++ {
		if err := c.writeBlock(patchToneArea(area, t), p.Tones[t-1]); err != nil {
			return err
		}
	}
	return nil
}

// StorePerformance stores a performance into the given user performance
// number (1-32).
func (c *Client) StorePerformance(number int, p *Performance) error {
	if number < 1 || number > 32 {
		return fmt.Errorf("jv1080: user performance number %d out of range 1-32", number)
	}
	area := userPerformanceArea(number)
	if err := c.writeBlock(area, p.Common); err != nil {
		return err
	}
	for n := 1; n <= 16; n++ {
		addr, _ := performancePartBlock(area, n)
		if err := c.writeBlock(addr, p.Parts[n-1]); err != nil {
			return err
		}
	}
	return nil
}

// SetPatchParameter updates a single parameter of a user patch on the
// device. The parameter is identified by its block ("common" or "tone"),
// the tone number (1-4, 0 for common) and the byte offset within the
// block. value holds the raw parameter bytes (one byte, or two nibbles for
// nibblized parameters such as the wave number).
func (c *Client) SetPatchParameter(bank Bank, number int, block string, tone int, offset int, value []byte) error {
	var area midi.Address
	if bank == BankUser {
		area = userPatchArea(number)
	} else {
		area = TemporaryPatchArea()
	}
	var addr midi.Address
	switch block {
	case "common":
		addr = area.Add(offset)
	case "tone":
		if tone < 1 || tone > 4 {
			return fmt.Errorf("jv1080: tone %d out of range 1-4", tone)
		}
		addr = patchToneArea(area, tone).Add(offset)
	default:
		return fmt.Errorf("jv1080: unknown patch block %q", block)
	}
	return c.sendDataSet(addr, value)
}

// SetPerformanceParameter updates a single parameter of a user performance
// on the device, identified by block ("common" or "part"), part number
// (1-16, 0 for common) and byte offset within the block.
func (c *Client) SetPerformanceParameter(number int, block string, part int, offset int, value []byte) error {
	area := userPerformanceArea(number)
	var addr midi.Address
	switch block {
	case "common":
		addr = area.Add(offset)
	case "part":
		if part < 1 || part > 16 {
			return fmt.Errorf("jv1080: part %d out of range 1-16", part)
		}
		paddr, _ := performancePartBlock(area, part)
		addr = paddr.Add(offset)
	default:
		return fmt.Errorf("jv1080: unknown performance block %q", block)
	}
	return c.sendDataSet(addr, value)
}

// NoteOn sends a note on message.
func (c *Client) NoteOn(channel uint8, note, velocity uint8) error {
	return c.dev.SendMessage(midi.NoteOn{Channel: channel, Note: note, Velocity: velocity})
}

// NoteOff sends a note off message (as note on with velocity 0).
func (c *Client) NoteOff(channel uint8, note uint8) error {
	return c.dev.SendMessage(midi.NoteOn{Channel: channel, Note: note, Velocity: 0})
}
