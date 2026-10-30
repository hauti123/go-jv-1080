// Package gui contains the widgets of the JV-1080 editor GUI: a compact
// LED-style knob gadget and the shared dark MK-II-like theme. The package
// is deliberately independent of the MIDI/device layer; applications wire
// the widgets to a device through their own controller (see
// cmd/jv1080-gui).
package gui

import "image/color"

// Theme holds the colors of the dark MK-II-style JV-1080 GUI. The accent
// color is green (as on the MK-II Tone Stack display) instead of orange.
type Theme struct {
	Background color.NRGBA // window background
	Surface    color.NRGBA // widget surface
	LEDOn      color.NRGBA // lit LED / active scale segment (green)
	LEDOff     color.NRGBA // unlit scale segment
	Text       color.NRGBA // LED text color
	TextDim    color.NRGBA // titles and units
}

// DefaultTheme returns the dark green MK-II-style theme.
func DefaultTheme() Theme {
	return Theme{
		Background: color.NRGBA{R: 0x10, G: 0x12, B: 0x10, A: 0xff},
		Surface:    color.NRGBA{R: 0x1a, G: 0x1e, B: 0x1a, A: 0xff},
		LEDOn:      color.NRGBA{R: 0x2a, G: 0xff, B: 0x4a, A: 0xff},
		LEDOff:     color.NRGBA{R: 0x24, G: 0x33, B: 0x26, A: 0xff},
		Text:       color.NRGBA{R: 0x9d, G: 0xff, B: 0xb0, A: 0xff},
		TextDim:    color.NRGBA{R: 0x4a, G: 0x7a, B: 0x55, A: 0xff},
	}
}

// FontScale is the scale of the 5x7 LED font used by the widgets. The DC-90
// style prototype draws the font at twice this size; the widgets use half
// of it so that knobs can be stacked compactly.
const FontScale = 1.0
