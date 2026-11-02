package gui

import (
	"fmt"
	"image"
	"log"
	"sync"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"go-jv-1080/internal/gui/widget"
	"go-jv-1080/internal/jv1080"
)

// screen is one screen of the app (the patch edit screen). A screen
// renders itself and applies device-originated parameter updates to its
// widgets.
type screen interface {
	Layout(gtx layout.Context, th *widget.Theme) layout.Dimensions
	applyParam(key string, value any)
	initDevice(c *jv1080.Client)
}

// App is the JV-1080 editor application: the window, the status line and
// the current screen. All widget change callbacks run on the GUI
// goroutine. The widgets are bound to device parameters by path-like
// keys; all address/offset handling and value conversion happens in
// internal/jv1080.
type App struct {
	th     widget.Theme
	screen screen

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

// NewApp returns the app with the patch edit screen installed.
func NewApp() *App {
	a := &App{
		th:      widget.DefaultTheme(),
		updates: make(chan paramUpdate, 16),
		status:  "patch/temp",
	}
	a.screen = newPatchEditor(a.send)
	return a
}

// Main creates the application window and runs the event loop until the
// window is closed. It must be called from the main goroutine.
func (a *App) Main(title string, width, height unit.Dp) {
	go func() {
		w := new(app.Window)
		w.Option(app.Title(title), app.Size(width, height))
		a.w = w
		if err := a.run(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

// SetDevice installs the device client, registers for parameter updates
// and lets the screen read its current values once.
func (a *App) SetDevice(c *jv1080.Client) {
	a.mu.Lock()
	a.client = c
	a.mu.Unlock()
	c.OnParameterChange(func(key string, value any) {
		// The MIDI receive callback runs on the driver's goroutine;
		// hand the update to the GUI goroutine and wake it up.
		select {
		case a.updates <- paramUpdate{key: key, value: value}:
		default:
		}
		if a.w != nil {
			a.w.Invalidate()
		}
	})
	a.screen.initDevice(c)
}

// send writes a changed parameter to the device immediately, addressed by
// its path-like key with the UI value.
func (a *App) send(key string, value any) {
	a.mu.Lock()
	c := a.client
	a.mu.Unlock()
	if c == nil {
		a.setStatus("no device connected")
		return
	}
	if err := c.SetParameter(key, value); err != nil {
		a.setStatus(fmt.Sprintf("write %s failed: %v", key, err))
	}
}

func (a *App) setStatus(s string) {
	a.mu.Lock()
	a.status = s
	a.mu.Unlock()
}

func (a *App) run(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			a.layout(gtx, w)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) layout(gtx layout.Context, w *app.Window) layout.Dimensions {
	// apply device-originated updates (drain the channel); the callback
	// already carries the key and the UI value, so only the one changed
	// widget is updated
	dirty := false
	for {
		select {
		case up := <-a.updates:
			a.screen.applyParam(up.key, up.value)
			dirty = true
		default:
			goto drained
		}
	}
drained:
	if dirty {
		w.Invalidate()
	}

	paint.Fill(gtx.Ops, a.th.Background)
	return layout.UniformInset(unit.Dp(18)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				a.mu.Lock()
				s := a.status
				a.mu.Unlock()
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
				return a.screen.Layout(gtx, &a.th)
			}),
		)
	})
}
