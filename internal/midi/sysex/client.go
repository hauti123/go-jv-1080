package sysex

import (
	"fmt"
	"time"

	"go-jv-1080/internal/midi"
	"go-jv-1080/internal/midi/message"
)

// DefaultTimeout is the default time to wait for the device to answer a
// data request.
const DefaultTimeout = 2 * time.Second

// writeInterval is the pause between DT1 packets when writing a block.
const writeInterval = 20 * time.Millisecond

// Client is a SysEx level API to a JV-1080 over a MIDI device connection:
// downloading patches and performances (single or whole banks) and
// modifying single parameters.
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
func (c *Client) request(addr message.Address, size int) (message.DataSet, error) {
	req := message.NewDataRequest(c.deviceID, addr, size)
	return c.dev.Request(req, c.timeout, func(ds message.DataSet) bool {
		return len(ds.Data) == size
	})
}

// ReadBlock downloads size bytes starting at addr.
func (c *Client) ReadBlock(addr message.Address, size int) ([]byte, error) {
	ds, err := c.request(addr, size)
	if err != nil {
		return nil, err
	}
	return ds.Data, nil
}

// PatchData is the raw data of a patch: the common block plus the four
// tone blocks.
type PatchData struct {
	Area   message.Address
	Common []byte
	Tones  [PatchToneCount][]byte
}

// DownloadPatch downloads a whole patch (common block plus four tones)
// from the given area.
func (c *Client) DownloadPatch(area message.Address) (*PatchData, error) {
	p := &PatchData{Area: area, Common: make([]byte, PatchCommonSize)}
	addr, size := PatchCommonBlock(area)
	data, err := c.ReadBlock(addr, size)
	if err != nil {
		return nil, fmt.Errorf("sysex: read patch common: %w", err)
	}
	copy(p.Common, data)
	for t := 1; t <= PatchToneCount; t++ {
		p.Tones[t-1] = make([]byte, PatchToneSize)
		addr, size := PatchToneBlock(area, t)
		data, err := c.ReadBlock(addr, size)
		if err != nil {
			return nil, fmt.Errorf("sysex: read patch tone %d: %w", t, err)
		}
		copy(p.Tones[t-1], data)
	}
	return p, nil
}

// DownloadTemporaryPatch downloads the patch-mode temporary patch.
func (c *Client) DownloadTemporaryPatch() (*PatchData, error) {
	return c.DownloadPatch(TemporaryPatchArea())
}

// DownloadUserPatch downloads user patch number (1-128).
func (c *Client) DownloadUserPatch(number int) (*PatchData, error) {
	if number < 1 || number > UserPatchCount {
		return nil, fmt.Errorf("sysex: user patch number %d out of range 1-%d", number, UserPatchCount)
	}
	return c.DownloadPatch(UserPatchArea(number))
}

// DownloadUserPatchBank downloads all 128 user patches (the full user
// patch bank).
func (c *Client) DownloadUserPatchBank() ([]*PatchData, error) {
	patches := make([]*PatchData, 0, UserPatchCount)
	for n := 1; n <= UserPatchCount; n++ {
		p, err := c.DownloadUserPatch(n)
		if err != nil {
			return nil, err
		}
		patches = append(patches, p)
	}
	return patches, nil
}

// DownloadCardPatchBank downloads all 128 patches of the data card.
func (c *Client) DownloadCardPatchBank() ([]*PatchData, error) {
	patches := make([]*PatchData, 0, UserPatchCount)
	for n := 1; n <= UserPatchCount; n++ {
		p, err := c.DownloadPatch(CardPatchArea(n))
		if err != nil {
			return nil, err
		}
		patches = append(patches, p)
	}
	return patches, nil
}

// PerformanceData is the raw data of a performance: the common block plus
// the 16 part blocks.
type PerformanceData struct {
	Area   message.Address
	Common []byte
	Parts  [PerformancePartCount][]byte
}

// DownloadPerformance downloads a whole performance (common block plus 16
// parts) from the given area.
func (c *Client) DownloadPerformance(area message.Address) (*PerformanceData, error) {
	p := &PerformanceData{Area: area, Common: make([]byte, PerformanceCommonSize)}
	addr, size := PerformanceCommonBlock(area)
	data, err := c.ReadBlock(addr, size)
	if err != nil {
		return nil, fmt.Errorf("sysex: read performance common: %w", err)
	}
	copy(p.Common, data)
	for n := 1; n <= PerformancePartCount; n++ {
		p.Parts[n-1] = make([]byte, PerformancePartSize)
		addr, size := PerformancePartBlock(area, n)
		data, err := c.ReadBlock(addr, size)
		if err != nil {
			return nil, fmt.Errorf("sysex: read performance part %d: %w", n, err)
		}
		copy(p.Parts[n-1], data)
	}
	return p, nil
}

// DownloadTemporaryPerformance downloads the temporary performance.
func (c *Client) DownloadTemporaryPerformance() (*PerformanceData, error) {
	return c.DownloadPerformance(TemporaryPerformanceArea())
}

// DownloadUserPerformance downloads user performance number (1-32).
func (c *Client) DownloadUserPerformance(number int) (*PerformanceData, error) {
	if number < 1 || number > UserPerformanceCount {
		return nil, fmt.Errorf("sysex: user performance number %d out of range 1-%d", number, UserPerformanceCount)
	}
	return c.DownloadPerformance(UserPerformanceArea(number))
}

// DownloadUserPerformanceBank downloads all 32 user performances.
func (c *Client) DownloadUserPerformanceBank() ([]*PerformanceData, error) {
	perfs := make([]*PerformanceData, 0, UserPerformanceCount)
	for n := 1; n <= UserPerformanceCount; n++ {
		p, err := c.DownloadUserPerformance(n)
		if err != nil {
			return nil, err
		}
		perfs = append(perfs, p)
	}
	return perfs, nil
}

// WriteBlock writes a block of data at addr, splitting it into DT1 packets
// of at most MaxDataSetBytes data bytes with a 20 ms pause between them,
// as the MIDI implementation requires.
func (c *Client) WriteBlock(addr message.Address, data []byte) error {
	for off := 0; off < len(data); off += MaxDataSetBytes {
		end := off + MaxDataSetBytes
		if end > len(data) {
			end = len(data)
		}
		ds := message.NewDataSet(c.deviceID, addr.Add(off), data[off:end])
		if err := c.dev.SendMessage(ds); err != nil {
			return err
		}
		time.Sleep(writeInterval)
	}
	return nil
}

// ParameterAddress returns the address of the parameter at offset within
// the block at area.
func ParameterAddress(area message.Address, offset int) message.Address {
	return area.Add(offset)
}

// PatchParameterAddress returns the address of a single patch parameter.
// block is "common" or "tone"; tone is 1-4 (ignored for "common"); offset
// is the byte offset within the block.
func PatchParameterAddress(area message.Address, block string, tone, offset int) (message.Address, error) {
	switch block {
	case "common":
		return ParameterAddress(area, offset), nil
	case "tone":
		if tone < 1 || tone > PatchToneCount {
			return message.Address{}, fmt.Errorf("sysex: tone %d out of range 1-%d", tone, PatchToneCount)
		}
		return ParameterAddress(PatchToneArea(area, tone), offset), nil
	}
	return message.Address{}, fmt.Errorf("sysex: unknown patch block %q", block)
}

// PerformanceParameterAddress returns the address of a single performance
// parameter. block is "common" or "part"; part is 1-16 (ignored for
// "common"); offset is the byte offset within the block.
func PerformanceParameterAddress(area message.Address, block string, part, offset int) (message.Address, error) {
	switch block {
	case "common":
		return ParameterAddress(area, offset), nil
	case "part":
		if part < 1 || part > PerformancePartCount {
			return message.Address{}, fmt.Errorf("sysex: part %d out of range 1-%d", part, PerformancePartCount)
		}
		addr, _ := PerformancePartBlock(area, part)
		return ParameterAddress(addr, offset), nil
	}
	return message.Address{}, fmt.Errorf("sysex: unknown performance block %q", block)
}

// SetParameter writes a single parameter: value holds the raw parameter
// bytes (one byte, or two nibbles for nibblized parameters such as the
// wave number) at offset within the block at area.
func (c *Client) SetParameter(area message.Address, offset int, value []byte) error {
	ds := message.NewDataSet(c.deviceID, ParameterAddress(area, offset), value)
	return c.dev.SendMessage(ds)
}

// SetPatchParameter updates a single parameter of a patch at area. See
// PatchParameterAddress for the meaning of block, tone and offset.
func (c *Client) SetPatchParameter(area message.Address, block string, tone, offset int, value []byte) error {
	addr, err := PatchParameterAddress(area, block, tone, offset)
	if err != nil {
		return err
	}
	ds := message.NewDataSet(c.deviceID, addr, value)
	return c.dev.SendMessage(ds)
}

// SetPerformanceParameter updates a single parameter of a performance at
// area. See PerformanceParameterAddress for the meaning of block, part and
// offset.
func (c *Client) SetPerformanceParameter(area message.Address, block string, part, offset int, value []byte) error {
	addr, err := PerformanceParameterAddress(area, block, part, offset)
	if err != nil {
		return err
	}
	ds := message.NewDataSet(c.deviceID, addr, value)
	return c.dev.SendMessage(ds)
}
