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
// line). The GUI talks to the device only through the gui.Controller
// interface, keeping the widget layer independent of the MIDI layer.
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

	"go-jv-1080/internal/gui"
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

func main() {
	flag.Parse()
	ui := newUI()

	dev, err := midi.OpenDevice(*outName, *inName, *outNum, *inNum)
	if err != nil {
		log.Printf("jv1080-gui: cannot open MIDI device (%v); running without device", err)
	} else {
		defer dev.Close()
		// Log all raw MIDI messages sent and received in hex.
		dev.SetLogger(func(direction string, msg []byte) {
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
		w.Option(app.Title("JV-1080 Patch Common"), app.Size(unit.Dp(360), unit.Dp(220)))
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
	th     gui.Theme
	level  *gui.Knob
	pan    *gui.Knob
	octave *gui.Knob

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
		th:      gui.DefaultTheme(),
		updates: make(chan paramUpdate, 16),
		status:  "patch/temp",
	}
	u.level = gui.NewKnob(gui.KnobMinMax, "LEVEL", keyLevel, 0, 127)
	u.pan = gui.NewKnob(gui.KnobPan, "PAN", keyPan, -64, 63)
	u.octave = gui.NewKnob(gui.KnobBipolar, "OCTAVE", keyOctave, -3, 3)
	u.level.OnChange = func(v int) { u.send(keyLevel, v) }
	u.pan.OnChange = func(v int) { u.send(keyPan, v) }
	u.octave.OnChange = func(v int) { u.send(keyOctave, v) }
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
				cell := gtx.Dp(gui.FontCellDp)
				sz := gui.PixelTextSize(s, cell)
				x := (gtx.Constraints.Max.X - sz.X) / 2
				if x < 0 {
					x = 0
				}
				var d layout.Dimensions
				gui.Offset(gtx, x, 0, func() { d = gui.DrawPixelText(gtx, s, cell) })
				return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: d.Size.Y + gtx.Dp(unit.Dp(10))}}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.level.Layout(gtx, &u.th) }),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.pan.Layout(gtx, &u.th) }),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.octave.Layout(gtx, &u.th) }),
				)
			}),
		)
	})
}
