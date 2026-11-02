// Command jv1080-gui is the JV-1080 editor GUI: a thin command that opens
// the MIDI device, wires it to the device client and starts the
// application window. The window, layout and the patch edit screen (a
// COMMON tab and one tab per tone with all patch parameters) live in
// internal/gui. Changes made with the widgets are sent to the device
// immediately (single-parameter DT1 writes); changes made on the device
// itself are reflected in the widgets immediately.
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
// Without a connected device the app still starts and shows the widgets;
// they then just do not send anything (an error is shown in the status
// line). The GUI talks to the device only through the jv1080.Client,
// keeping the widget layer independent of the MIDI layer.
//
// Raw MIDI logging is on by default and can be disabled per direction
// with the environment variables LOG_MIDI_RX=0 / LOG_MIDI_TX=0; widget
// debug logging is enabled with LOG_GUI_WIDGETS=1.
package main

import (
	"flag"
	"log"
	"strconv"

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

func main() {
	flag.Parse()
	a := gui.NewApp()

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
		a.SetDevice(jv1080.NewClient(dev, uint8(id)))
	}

	a.Main("JV-1080 Patch Editor", unit.Dp(780), unit.Dp(620))
}
