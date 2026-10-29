package midi

import (
	"fmt"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // registers the rtmidi driver

	"go-jv-1080/internal/midi/message"
)

// PortInfo describes an available MIDI input or output port.
type PortInfo struct {
	Number int
	Name   string
}

// InputPorts lists the available MIDI input ports.
func InputPorts() ([]PortInfo, error) {
	ins, err := drivers.Ins()
	if err != nil {
		return nil, err
	}
	out := make([]PortInfo, 0, len(ins))
	for _, in := range ins {
		out = append(out, PortInfo{Number: in.Number(), Name: in.String()})
	}
	return out, nil
}

// OutputPorts lists the available MIDI output ports.
func OutputPorts() ([]PortInfo, error) {
	outs, err := drivers.Outs()
	if err != nil {
		return nil, err
	}
	out := make([]PortInfo, 0, len(outs))
	for _, o := range outs {
		out = append(out, PortInfo{Number: o.Number(), Name: o.String()})
	}
	return out, nil
}

// Device is a MIDI device connection to a JV-1080: one output port for
// sending and one input port for receiving. Incoming, unsolicited MIDI
// messages (replies, dumps, identity replies, ...) are dispatched to all
// registered callbacks.
type Device struct {
	out drivers.Out
	in  drivers.In

	mu         sync.RWMutex
	callbacks  []func(msg []byte)
	sendLogger func(msg []byte)
	stop       func()
	done       chan struct{}
}

// SetSendLogger registers a logger that is called with every raw message
// sent to the device (after it was split into individual messages).
func (d *Device) SetSendLogger(logger func(msg []byte)) {
	d.mu.Lock()
	d.sendLogger = logger
	d.mu.Unlock()
}

// MessageHandler is called for every incoming MIDI message. The byte slice
// is only valid for the duration of the call.
type MessageHandler func(msg []byte)

// OpenDevice opens the MIDI output and input port with the given names (or
// numbers when the name is empty) and starts listening for incoming
// messages.
func OpenDevice(outName, inName string, outNum, inNum int) (*Device, error) {
	var out drivers.Out
	var err error
	if outName != "" {
		out, err = drivers.OutByName(outName)
	} else {
		out, err = drivers.OutByNumber(outNum)
	}
	if err != nil {
		return nil, fmt.Errorf("midi: cannot open output port: %w", err)
	}
	var in drivers.In
	if inName != "" {
		in, err = drivers.InByName(inName)
	} else {
		in, err = drivers.InByNumber(inNum)
	}
	if err != nil {
		out.Close()
		return nil, fmt.Errorf("midi: cannot open input port: %w", err)
	}
	d := &Device{out: out, in: in, done: make(chan struct{})}
	if err := d.startListening(); err != nil {
		in.Close()
		out.Close()
		return nil, err
	}
	return d, nil
}

func (d *Device) startListening() error {
	stop, err := d.in.Listen(func(msg []byte, _ int32) {
		d.mu.RLock()
		cbs := append([]func([]byte){}, d.callbacks...)
		d.mu.RUnlock()
		for _, cb := range cbs {
			cb(msg)
		}
	}, drivers.ListenConfig{SysEx: true, SysExBufferSize: 1 << 20})
	if err != nil {
		return fmt.Errorf("midi: cannot listen on input port: %w", err)
	}
	d.stop = stop
	return nil
}

// RegisterCallback registers a handler that is called for every incoming
// MIDI message, including unsolicited ones sent by the device (data set
// replies to requests, identity replies, ...). It returns a function that
// unregisters the handler again.
func (d *Device) RegisterCallback(cb MessageHandler) (unregister func()) {
	d.mu.Lock()
	d.callbacks = append(d.callbacks, cb)
	d.mu.Unlock()
	return func() {
		d.mu.Lock()
		for i, c := range d.callbacks {
			if fmt.Sprintf("%p", c) == fmt.Sprintf("%p", cb) {
				d.callbacks = append(d.callbacks[:i], d.callbacks[i+1:]...)
				break
			}
		}
		d.mu.Unlock()
	}
}

// Send sends a raw MIDI message to the device.
func (d *Device) Send(msg []byte) error {
	d.mu.RLock()
	logger := d.sendLogger
	d.mu.RUnlock()
	if logger != nil {
		logger(msg)
	}
	return d.out.Send(msg)
}

// SendMessage serializes msg and sends it to the device. The byte stream
// is split into the individual MIDI messages first: some drivers (e.g.
// the Windows multimedia backend) reject multi-message streams that are
// not SysEx.
func (d *Device) SendMessage(msg message.Message) error {
	return d.SendMessages(msg.Bytes())
}

// SendMessages splits a byte stream of concatenated MIDI messages and
// sends each message individually.
func (d *Device) SendMessages(b []byte) error {
	for _, m := range message.SplitMessages(b) {
		if err := d.Send(m); err != nil {
			return err
		}
	}
	return nil
}

// Request sends a data request (RQ1) and waits for the device to answer
// with the corresponding data set (DT1) message. The wait ends after
// timeout. The match function decides whether an incoming DT1 is the
// expected answer; when nil, the first DT1 for the requested address is
// accepted.
func (d *Device) Request(req message.DataRequest, timeout time.Duration, match func(message.DataSet) bool) (message.DataSet, error) {
	result := make(chan message.DataSet, 1)
	unregister := d.RegisterCallback(func(raw []byte) {
		ds, err := message.ParseDataSet(raw)
		if err != nil {
			return
		}
		if ds.DeviceID != req.DeviceID || ds.Address != req.Address {
			return
		}
		if match != nil && !match(ds) {
			return
		}
		select {
		case result <- ds:
		default:
		}
	})
	defer unregister()
	if err := d.SendMessage(req); err != nil {
		return message.DataSet{}, err
	}
	select {
	case ds := <-result:
		return ds, nil
	case <-time.After(timeout):
		return message.DataSet{}, fmt.Errorf("midi: timeout waiting for answer to %s", req)
	}
}

// RequestIdentity sends a universal identity request and waits for the
// identity reply.
func (d *Device) RequestIdentity(deviceID uint8, timeout time.Duration) (message.IdentityReply, error) {
	result := make(chan message.IdentityReply, 1)
	unregister := d.RegisterCallback(func(raw []byte) {
		reply, err := message.ParseIdentityReply(raw)
		if err != nil {
			return
		}
		if reply.DeviceID != deviceID {
			return
		}
		select {
		case result <- reply:
		default:
		}
	})
	defer unregister()
	if err := d.SendMessage(message.IdentityRequest{DeviceID: deviceID}); err != nil {
		return message.IdentityReply{}, err
	}
	select {
	case r := <-result:
		return r, nil
	case <-time.After(timeout):
		return message.IdentityReply{}, fmt.Errorf("midi: timeout waiting for identity reply from device %02X", deviceID)
	}
}

// Close stops listening and closes both ports.
func (d *Device) Close() error {
	if d.stop != nil {
		d.stop()
	}
	close(d.done)
	if err := d.in.Close(); err != nil {
		d.out.Close()
		return err
	}
	return d.out.Close()
}
