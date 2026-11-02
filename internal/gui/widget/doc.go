// Package widget contains the widgets of the JV-1080 editor GUI: a compact
// hardware-style knob gadget and the shared theme. The visual style follows
// the MK-II Tone Stack example of cmd/gui-proto (metal cap + LED arc
// scale) with green instead of amber as the accent color, and the 5x7
// dot-matrix LED font of the DC-90 sample at half its cell size. The
// package is deliberately independent of the MIDI/device layer;
// applications wire the widgets to a device through their own controller
// (see cmd/jv1080-gui).
//
// Debug logging of the widgets (drag handling, relayout, ...) is enabled
// by setting the LOG_GUI_WIDGETS environment variable.
package widget
