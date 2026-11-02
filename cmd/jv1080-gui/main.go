// Command jv1080-gui is a small sample GUI that controls three patch common
// parameters of the Roland JV-1080: level, pan and octave (named like in
// jv1080-cli). Changes made with the knobs are sent to the device
// immediately (single-parameter DT1 writes); changes made on the device
// itself are reflected in the knobs immediately.
//
// Usage:
//
//	jv1080-gui [flags]
//
// Flags:
//
//	-out string    MIDI output port name (default "SERIES 208i MIDI")
//	-in string     MIDI input port name (default "SERIES 208i MIDI")
//	-outnum int    MIDI output port number, used when -out is empty (default 0)
//	-innum int     MIDI input port number, used when -in is empty (default 0)
//	-devid string  device ID in hex (default "10" = unit 17)
//
// Without a connected device the app still starts and shows the knobs; the
// knobs then just do not send anything (an error is shown in the status
// line). The GUI talks to the device only through the jv1080.Client,
// keeping the widget layer independent of the MIDI layer.
//
// Raw MIDI logging is on by default and can be disabled per direction
// with the environment variables LOG_MIDI_RX=0 / LOG_MIDI_TX=0; widget
// debug logging is enabled with LOG_GUI_WIDGETS=1.
package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"strconv"
	"sync"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"go-jv-1080/internal/gui/widget"
	"go-jv-1080/internal/jv1080"
	"go-jv-1080/internal/midi"
	"go-jv-1080/internal/midi/message"
)

var (
	outName  = flag.String("out", "SERIES 208i MIDI", "MIDI output port name")
	inName   = flag.String("in", "SERIES 208i MIDI", "MIDI input port name")
	outNum   = flag.Int("outnum", 0, "MIDI output port number (used when -out is empty)")
	inNum    = flag.Int("innum", 0, "MIDI input port number (used when -in is empty)")
	deviceID = flag.String("devid", "10", "device ID in hex")
)

// device parameter keys of the temporary patch (patch/temp/...), named
// like in jv1080-cli. The GUI never deals with addresses or offsets; the
// key -> SysEx mapping lives in internal/jv1080.
const (
	keyLevel  = "patch/temp/common/level"
	keyPan    = "patch/temp/common/pan"
	keyOctave = "patch/temp/common/octave"
)

// pitch envelope of tone 1: 4 times + 4 levels (PCH ENVELOPE)
var (
	envTimeKeys  = [4]string{"patch/temp/tone/1/pitch/t1", "patch/temp/tone/1/pitch/t2", "patch/temp/tone/1/pitch/t3", "patch/temp/tone/1/pitch/t4"}
	envLevelKeys = [4]string{"patch/temp/tone/1/pitch/l1", "patch/temp/tone/1/pitch/l2", "patch/temp/tone/1/pitch/l3", "patch/temp/tone/1/pitch/l4"}
)

func main() {
	flag.Parse()
	ui := newUI()

	dev, err := midi.OpenDevice(*outName, *inName, *outNum, *inNum)
	if err != nil {
		log.Printf("jv1080-gui: cannot open MIDI device (%v); running without device", err)
	} else {
		defer dev.Close()
		// Log all raw MIDI messages sent and received in hex (default on;
		// disable per direction with LOG_MIDI_RX=0 / LOG_MIDI_TX=0).
		logRX := midi.LogEnabled("RX")
		logTX := midi.LogEnabled("TX")
		dev.SetLogger(func(direction string, msg []byte) {
			if direction == "RX" && !logRX || direction == "TX" && !logTX {
				return
			}
			log.Printf("MIDI %s: [%s]", direction, message.Hex(msg))
		})
		id, err := strconv.ParseUint(*deviceID, 16, 8)
		if err != nil {
			log.Fatalf("jv1080-gui: invalid device ID %q: %v", *deviceID, err)
		}
		ui.setDevice(jv1080.NewClient(dev, uint8(id)))
	}

	go func() {
		w := new(app.Window)
		w.Option(app.Title("JV-1080 Patch Common"), app.Size(unit.Dp(520), unit.Dp(460)))
		ui.w = w
		if err := ui.loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

// ui holds the GUI state: the three knobs, the device client and the
// status line. All knob change callbacks run on the GUI goroutine. The
// knobs are bound to device parameters by their path-like keys; all
// address/offset handling and value conversion happens in internal/jv1080.
type ui struct {
	th     widget.Theme
	level  *widget.Knob
	pan    *widget.Knob
	octave *widget.Knob
	env    *widget.Envelope

	mu      sync.Mutex
	client  *jv1080.Client
	status  string
	updates chan paramUpdate
	w       *app.Window // set once the window exists, for Invalidate
}

// paramUpdate is a device-originated parameter change, delivered to the
// GUI goroutine.
type paramUpdate struct {
	key   string
	value any
}

func newUI() *ui {
	u := &ui{
		th:      widget.DefaultTheme(),
		updates: make(chan paramUpdate, 16),
		status:  "patch/temp",
	}
	u.level = widget.NewKnob(widget.KnobMinMax, "LEVEL", keyLevel, 0, 127)
	u.pan = widget.NewKnob(widget.KnobPan, "PAN", keyPan, -64, 63)
	u.octave = widget.NewKnob(widget.KnobBipolar, "OCTAVE", keyOctave, -3, 3)
	u.level.OnChange = func(v int) { u.send(keyLevel, v) }
	u.pan.OnChange = func(v int) { u.send(keyPan, v) }
	u.octave.OnChange = func(v int) { u.send(keyOctave, v) }
	u.env = widget.NewEnvelope("PITCH ENV T1", envTimeKeys, envLevelKeys)
	u.env.OnChange = func(timeKey string, timeVal int, levelKey string, levelVal int) {
		u.send(timeKey, timeVal)
		u.send(levelKey, levelVal)
	}
	return u
}

// setDevice installs the device client and registers for parameter
// updates, then reads the current values of the three parameters once.
func (u *ui) setDevice(c *jv1080.Client) {
	u.mu.Lock()
	u.client = c
	u.mu.Unlock()
	c.OnParameterChange(func(key string, value any) {
		// The MIDI receive callback runs on the driver's goroutine;
		// hand the update to the GUI goroutine and wake it up.
		select {
		case u.updates <- paramUpdate{key: key, value: value}:
		default:
		}
		if u.w != nil {
			u.w.Invalidate()
		}
	})
	for _, key := range []string{keyLevel, keyPan, keyOctave} {
		if v, err := c.GetParameter(key); err == nil {
			u.applyParam(key, v)
		}
	}
	// read the pitch envelope once (SetAll so the x-order clamp does
	// not interfere with loading)
	var times, levels [4]int
	for i := 0; i < 4; i++ {
		if v, err := c.GetParameter(envTimeKeys[i]); err == nil {
			if n, ok := v.(int); ok {
				times[i] = n
			}
		}
		if v, err := c.GetParameter(envLevelKeys[i]); err == nil {
			if n, ok := v.(int); ok {
				levels[i] = n
			}
		}
	}
	u.env.SetAll(times, levels)
}

// send writes a changed parameter to the device immediately, addressed by
// its path-like key with the UI value.
func (u *ui) send(key string, value any) {
	u.mu.Lock()
	c := u.client
	u.mu.Unlock()
	if c == nil {
		u.setStatus("no device connected")
		return
	}
	if err := c.SetParameter(key, value); err != nil {
		u.setStatus(fmt.Sprintf("write %s failed: %v", key, err))
	}
}

// applyParam sets the knob bound to key without re-sending it.
func (u *ui) applyParam(key string, value any) {
	switch key {
	case keyLevel:
		if n, ok := value.(int); ok {
			u.level.Set(n)
		}
	case keyPan:
		u.pan.Set(panUIToKnob(value))
	case keyOctave:
		if n, ok := value.(int); ok {
			u.octave.Set(n)
		}
	default:
		for i := 0; i < 4; i++ {
			if key == envTimeKeys[i] || key == envLevelKeys[i] {
				tv, lv := u.env.Point(i)
				if n, ok := value.(int); ok {
					if key == envTimeKeys[i] {
						tv = n
					} else {
						lv = n
					}
				}
				u.env.SetPoint(i, tv, lv)
			}
		}
	}
}

// panUIToKnob converts a pan UI value from the device layer ("C", "L13",
// "R13" or signed -64..+63) into the knob's -64..+63 range.
func panUIToKnob(value any) int {
	switch v := value.(type) {
	case int:
		return v // device layer already returns signed UI units
	case string:
		s := v
		if len(s) >= 2 && (s[0] == 'L' || s[0] == 'R') {
			var n int
			fmt.Sscanf(s[1:], "%d", &n)
			if s[0] == 'L' {
				return -n
			}
			return n
		}
		return 0 // "C"
	}
	return 0
}

func (u *ui) setStatus(s string) {
	u.mu.Lock()
	u.status = s
	u.mu.Unlock()
}

func (u *ui) loop(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			u.layout(gtx, w)
			e.Frame(gtx.Ops)
		}
	}
}

func (u *ui) layout(gtx layout.Context, w *app.Window) layout.Dimensions {
	// apply device-originated updates (drain the channel); the callback
	// already carries the key and the UI value, so only the one changed
	// knob is updated
	dirty := false
	for {
		select {
		case up := <-u.updates:
			u.applyParam(up.key, up.value)
			dirty = true
		default:
			goto drained
		}
	}
drained:
	if dirty {
		w.Invalidate()
	}

	paint.Fill(gtx.Ops, u.th.Background)
	return layout.UniformInset(unit.Dp(18)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				u.mu.Lock()
				s := u.status
				u.mu.Unlock()
				cell := gtx.Dp(widget.FontCellDp)
				sz := widget.PixelTextSize(s, cell)
				x := (gtx.Constraints.Max.X - sz.X) / 2
				if x < 0 {
					x = 0
				}
				var d layout.Dimensions
				widget.Offset(gtx, x, 0, func() { d = widget.DrawPixelText(gtx, s, cell) })
				return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: d.Size.Y + gtx.Dp(unit.Dp(10))}}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.level.Layout(gtx, &u.th) }),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.pan.Layout(gtx, &u.th) }),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.octave.Layout(gtx, &u.th) }),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return u.env.Layout(gtx, &u.th)
				})
			}),
		)
	})
}
