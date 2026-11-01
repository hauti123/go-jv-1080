// Package gui contains the widgets of the JV-1080 editor GUI: a compact
// hardware-style knob gadget and the shared theme. The visual style follows
// the MK-II Tone Stack example of cmd/gui-proto (metal cap + LED arc
// scale) with green instead of amber as the accent color, and the 5x7
// dot-matrix LED font of the DC-90 sample at half its cell size. The
// package is deliberately independent of the MIDI/device layer;
// applications wire the widgets to a device through their own controller
// (see cmd/jv1080-gui).
package gui

import "image/color"

// Scale geometry shared by all knob styles, taken from cmd/gui-proto: the
// scale spans from -135 deg (7 o'clock) to +135 deg (5 o'clock), 0 deg is
// 12 o'clock. Angles are in degrees, 0 = 12 o'clock, increasing clockwise.
const (
	MinAngle = -135.0
	MaxAngle = 135.0
)

// LED palette of the DC-90 sample (cmd/gui-proto).
const (
	ledLitHex = 0x39FF70 // bright green, lit dots
	ledDimHex = 0x0B3315 // unlit dot grid
)

// Theme holds the colors of the MK-II-style JV-1080 GUI, modeled after the
// MK-II Tone Stack panel of cmd/gui-proto with green instead of amber.
type Theme struct {
	Background color.NRGBA // window background
	PanelTop   color.NRGBA // panel gradient top (dark metal)
	PanelBot   color.NRGBA // panel gradient bottom
	PanelEdge  color.NRGBA // panel border
	Title      color.NRGBA // panel title
	Sub        color.NRGBA // panel subtitle
	Label      color.NRGBA // knob title below the scale

	ArcBase  color.NRGBA // unlit part of the LED arc
	Lit      color.NRGBA // lit LED arc / ticks (green)
	TickBase color.NRGBA // unlit ticks
	CapTop   color.NRGBA // knob cap gradient top
	CapBot   color.NRGBA // knob cap gradient bottom
	CapEdge  color.NRGBA // knob cap outline
	Marker   color.NRGBA // indicator line on the cap
}

// DefaultTheme returns the MK-II Tone Stack style theme with green accent.
func DefaultTheme() Theme {
	return Theme{
		Background: rgb(0x0A0A0D),
		PanelTop:   rgb(0x26231F),
		PanelBot:   rgb(0x1B1916),
		PanelEdge:  rgb(0x0E0D0B),
		Title:      rgb(0xD8CDB2),
		Sub:        rgb(0x8D8474),
		Label:      rgb(0xB3A893),
		ArcBase:    rgb(0x34302A),
		Lit:        rgb(ledLitHex),
		TickBase:   rgb(0x6B655C),
		CapTop:     rgb(0x4A4640),
		CapBot:     rgb(0x141310),
		CapEdge:    rgb(0x0A0A08),
		Marker:     rgb(0xF4F1EA),
	}
}

// FontCellDp is the cell size of the 5x7 LED font: half the unit.Dp(2.6)
// used for the knob readout in the DC-90 sample of cmd/gui-proto.
const FontCellDp = 1.3
