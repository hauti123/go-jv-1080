// Package gui is the application layer of the JV-1080 editor GUI: the
// window and event loop (App in app.go), the status line and the current
// screen — today the PoC screen (poc.go), later the patch edit screen.
// It wires the widgets of the widget sub-package to the device through
// the jv1080.Client. The widgets themselves stay independent of the MIDI
// and device layers.
package gui
