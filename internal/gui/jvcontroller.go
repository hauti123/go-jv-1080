package gui

import (
	"fmt"
	"sync"

	"go-jv-1080/internal/jv1080"
	"go-jv-1080/internal/midi"
	"go-jv-1080/internal/midi/message"
)

// JVController is a Controller implementation backed by a jv1080.Client
// connected to a real device. It translates gui.Param descriptors into
// JV-1080 SysEx addresses of the currently selected bank/number and
// dispatches unsolicited DT1 messages from the device to registered
// Update callbacks.
//
// All public methods are safe to call from the GUI goroutine; MIDI I/O
// happens synchronously on the calling goroutine for Get/Set, while
// device-originated updates are delivered on the goroutine that runs
// DeliverDeviceMessages (usually the GUI event loop).
type JVController struct {
	client *jv1080.Client
	bank   jv1080.Bank
	number int // patch 1-128 or performance 1-32

	mu      sync.Mutex
	updates []func(p Param, value []byte)
	// paramByOffset maps common-block byte offsets to Param descriptors
	// for the current number; filled by RegisterCommonParams.
	commonParams map[int]Param
	toneParams   map[int]Param
	partParams   map[int]Param
}

// NewJVController returns a controller for the given client, bank and
// patch/performance number.
func NewJVController(client *jv1080.Client, bank jv1080.Bank, number int) *JVController {
	return &JVController{
		client:       client,
		bank:         bank,
		number:       number,
		commonParams: map[int]Param{},
		toneParams:   map[int]Param{},
		partParams:   map[int]Param{},
	}
}

// RegisterCommonParams tells the controller which common-block offsets
// correspond to which parameters, so that unsolicited DT1 writes to those
// offsets can be reported as typed updates. Call once after creation.
func (c *JVController) RegisterCommonParams(params ...Param) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range params {
		c.commonParams[p.Offset] = p
	}
}

// RegisterToneParams registers tone-block offsets (keyed by tone*0x1000 +
// offset).
func (c *JVController) RegisterToneParams(params ...Param) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range params {
		c.toneParams[p.Index*0x1000+p.Offset] = p
	}
}

// RegisterPartParams registers performance part offsets (keyed by
// part*0x1000 + offset).
func (c *JVController) RegisterPartParams(params ...Param) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range params {
		c.partParams[p.Index*0x1000+p.Offset] = p
	}
}

// address resolves a Param to its SysEx address on the device.
func (c *JVController) address(p Param) (message.Address, error) {
	return jv1080.PatchParameterAddress(c.bank, c.number, p.Block, p.Index, p.Offset)
}

// Get implements Controller.
func (c *JVController) Get(p Param) ([]byte, error) {
	addr, err := c.address(p)
	if err != nil {
		return nil, err
	}
	return c.client.ReadRaw(addr, 1)
}

// Set implements Controller: it writes the single parameter immediately.
func (c *JVController) Set(p Param, value []byte) error {
	addr, err := c.address(p)
	if err != nil {
		return err
	}
	return c.client.WriteRaw(addr, value)
}

// Update implements Controller.
func (c *JVController) Update(cb func(p Param, value []byte)) (unregister func()) {
	c.mu.Lock()
	c.updates = append(c.updates, cb)
	c.mu.Unlock()
	return func() {
		c.mu.Lock()
		for i, u := range c.updates {
			if fmt.Sprintf("%p", u) == fmt.Sprintf("%p", cb) {
				c.updates = append(c.updates[:i], c.updates[i+1:]...)
				break
			}
		}
		c.mu.Unlock()
	}
}

// HandleDeviceMessage inspects one raw incoming MIDI message. When it is
// an unsolicited DT1 data set for a registered parameter of the current
// bank/number, the registered Update callbacks are invoked with the
// decoded parameter. Call this from the MIDI receive path (e.g. from a
// midi.Device callback) for every incoming message.
func (c *JVController) HandleDeviceMessage(raw []byte) {
	ds, err := message.ParseDataSet(raw)
	if err != nil {
		return
	}
	if len(ds.Data) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	base, err := c.baseAddress()
	if err != nil {
		return
	}
	// Only react to writes inside the current patch area. Tone t (1-4)
	// starts at third address byte 0x10+2*(t-1); its block is PatchToneSize
	// bytes long, so tone t spans [0x10+2*(t-1)]*0x80 .. +PatchToneSize.
	off := ds.Address.Addr() - base.Addr()
	var p Param
	var ok bool
	switch {
	case off >= 0 && off < jv1080.PatchCommonSize:
		p, ok = c.commonParams[off]
	default:
		for t := 1; t <= 4; t++ {
			start := (0x10 + 2*(t-1)) * 0x80
			if off >= start && off < start+jv1080.PatchToneSize {
				p, ok = c.toneParams[t*0x1000+(off-start)]
				break
			}
		}
	}
	if !ok {
		return
	}
	for _, cb := range c.updates {
		cb(p, ds.Data)
	}
}

// baseAddress returns the SysEx address of the current patch (common
// block start).
func (c *JVController) baseAddress() (message.Address, error) {
	if c.bank == jv1080.BankUser {
		return jv1080.PatchParameterAddress(c.bank, c.number, "common", 0, 0)
	}
	return jv1080.TemporaryPatchArea(), nil
}

// Attach registers HandleDeviceMessage on the MIDI device so that
// unsolicited parameter changes from the device reach the GUI. It returns
// the unregister function.
func (c *JVController) Attach(dev *midi.Device) (unregister func()) {
	return dev.RegisterCallback(c.HandleDeviceMessage)
}
