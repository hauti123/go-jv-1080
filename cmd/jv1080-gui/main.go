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
//	-bank string   patch bank (USER, PRESET A/B/C, GM, CARD A/B, EXP A/B/C/D)
//	-number int    patch number 1-128 in the bank (default 1)
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
	"strings"
	"sync"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"go-jv-1080/internal/gui"
	"go-jv-1080/internal/jv1080"
	"go-jv-1080/internal/midi"
)

var (
	outName  = flag.String("out", "SERIES 208i MIDI", "MIDI output port name")
	inName   = flag.String("in", "SERIES 208i MIDI", "MIDI input port name")
	outNum   = flag.Int("outnum", 0, "MIDI output port number (used when -out is empty)")
	inNum    = flag.Int("innum", 0, "MIDI input port number (used when -in is empty)")
	deviceID = flag.String("devid", "10", "device ID in hex")
	bankName = flag.String("bank", "USER", "patch bank")
	number   = flag.Int("number", 1, "patch number 1-128 in the bank")
)

// parameter offsets in the patch common block, named like in jv1080-cli
// (doc/patch-edit-mode-structure.md).
const (
	offsetLevel  = 0x2E // "level" 0-127
	offsetPan    = 0x2F // "pan" 0-127, 64 = center
	offsetOctave = 0x41 // "octave" 0-6 = -3..+3
)

func main() {
	flag.Parse()
	bank, err := parseBank(*bankName)
	if err != nil {
		log.Fatalf("jv1080-gui: %v", err)
	}
	if *number < 1 || *number > 128 {
		log.Fatalf("jv1080-gui: patch number %d out of range 1-128", *number)
	}

	ui := newUI(bank, *number)

	dev, err := midi.OpenDevice(*outName, *inName, *outNum, *inNum)
	if err != nil {
		log.Printf("jv1080-gui: cannot open MIDI device (%v); running without device", err)
	} else {
		defer dev.Close()
		id, err := strconv.ParseUint(*deviceID, 16, 8)
		if err != nil {
			log.Fatalf("jv1080-gui: invalid device ID %q: %v", *deviceID, err)
		}
		client := jv1080.NewClient(dev, uint8(id))
		ctrl := gui.NewJVController(client, bank, *number)
		ctrl.RegisterCommonParams(
			gui.PatchCommonParam("level", offsetLevel),
			gui.PatchCommonParam("pan", offsetPan),
			gui.PatchCommonParam("octave", offsetOctave),
		)
		ui.setController(ctrl)
		unregister := ctrl.Attach(dev)
		defer unregister()
	}

	go func() {
		w := new(app.Window)
		w.Option(app.Title("JV-1080 Patch Common"), app.Size(unit.Dp(360), unit.Dp(220)))
		if err := ui.loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

// ui holds the GUI state: the three knobs, the device controller and the
// status line. All knob change callbacks run on the GUI goroutine.
type ui struct {
	th     gui.Theme
	level  *gui.Knob
	pan    *gui.Knob
	octave *gui.Knob

	mu         sync.Mutex
	controller gui.Controller
	status     string
	// updates delivers device-originated parameter changes to the GUI
	// goroutine.
	updates chan gui.Param
}

func newUI(bank jv1080.Bank, number int) *ui {
	u := &ui{
		th:      gui.DefaultTheme(),
		updates: make(chan gui.Param, 16),
	}
	u.level = gui.NewKnob(gui.KnobMinMax, "LEVEL", 0, 127)
	u.pan = gui.NewKnob(gui.KnobPan, "PAN", 0, 127)
	u.octave = gui.NewKnob(gui.KnobBipolar, "OCTAVE", -3, 3)
	u.level.OnChange = func(v int) { u.send(gui.PatchCommonParam("level", offsetLevel), v) }
	u.pan.OnChange = func(v int) { u.send(gui.PatchCommonParam("pan", offsetPan), v) }
	u.octave.OnChange = func(v int) { u.send(gui.PatchCommonParam("octave", offsetOctave), v+3) }
	u.status = fmt.Sprintf("%s %d", bank, number)
	return u
}

// setController installs the device controller and registers for updates.
func (u *ui) setController(c gui.Controller) {
	u.mu.Lock()
	u.controller = c
	u.mu.Unlock()
	c.Update(func(p gui.Param, value []byte) {
		select {
		case u.updates <- p:
		default:
		}
	})
	// read the current values from the device once
	for _, p := range []gui.Param{
		gui.PatchCommonParam("level", offsetLevel),
		gui.PatchCommonParam("pan", offsetPan),
		gui.PatchCommonParam("octave", offsetOctave),
	} {
		if v, err := c.Get(p); err == nil && len(v) > 0 {
			u.applyParam(p, int(v[0]))
		}
	}
}

// send writes a changed parameter to the device immediately.
func (u *ui) send(p gui.Param, raw int) {
	u.mu.Lock()
	c := u.controller
	u.mu.Unlock()
	if c == nil {
		u.setStatus("no device connected")
		return
	}
	if err := c.Set(p, []byte{byte(raw)}); err != nil {
		u.setStatus(fmt.Sprintf("write %s failed: %v", p.Name, err))
	}
}

// applyParam sets the knob belonging to p without re-sending it.
func (u *ui) applyParam(p gui.Param, raw int) {
	switch p.Name {
	case "level":
		u.level.Set(raw)
	case "pan":
		u.pan.Set(raw)
	case "octave":
		u.octave.Set(raw - 3)
	}
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
	// apply device-originated updates (drain the channel); each changed
	// parameter is re-read from the device so only that parameter is
	// updated in the UI
	dirty := false
	for {
		select {
		case p := <-u.updates:
			u.mu.Lock()
			c := u.controller
			u.mu.Unlock()
			if c != nil {
				if v, err := c.Get(p); err == nil && len(v) > 0 {
					u.applyParam(p, int(v[0]))
					dirty = true
				}
			}
		default:
			goto drained
		}
	}
drained:
	if dirty || u.level.ConsumeChanges() || u.pan.ConsumeChanges() || u.octave.ConsumeChanges() {
		w.Invalidate()
	}

	paint.Fill(gtx.Ops, u.th.Background)
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.level.Layout(gtx, &u.th) }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.pan.Layout(gtx, &u.th) }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return u.octave.Layout(gtx, &u.th) }),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			u.mu.Lock()
			s := u.status
			u.mu.Unlock()
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				sw := gui.LEDTextWidth(s, gui.FontScale)
				gui.LEDText(gtx.Ops, f32.Point{X: float32(gtx.Constraints.Max.X-int(sw)) / 2, Y: 0}, s, gui.FontScale, u.th.TextDim)
				return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: int(7*gui.FontScale) + gtx.Dp(unit.Dp(8))}}
			})
		}),
	)
}

func parseBank(s string) (jv1080.Bank, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	for _, b := range jv1080.Banks {
		if strings.ReplaceAll(b.String(), " ", "") == strings.ReplaceAll(s, " ", "") {
			return b, nil
		}
	}
	return 0, fmt.Errorf("unknown bank %q (want USER, PRESET A/B/C, GM, CARD A/B, EXP A/B/C/D)", s)
}
