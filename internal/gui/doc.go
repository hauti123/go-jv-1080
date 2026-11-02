// Package gui is the application layer of the JV-1080 editor GUI: the
// window and event loop (App in app.go), the status line and the current
// screen — the patch edit screen (patch.go) with a common tab and one
// tab per tone. It wires the widgets of the widget sub-package to the
// device through the jv1080.Client. The widgets themselves stay
// independent of the MIDI and device layers.
package gui
