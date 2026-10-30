package gui

// Param identifies one controllable parameter of the currently edited
// patch or performance on the JV-1080. It mirrors the CLI parameter names
// of cmd/jv1080-cli (doc/patch-edit-mode-structure.md).
type Param struct {
	// Block is "common" or "tone" (patch) / "part" (performance).
	Block string
	// Tone or Part number (1-4 / 1-16); 0 for the common block.
	Index int
	// Offset is the byte offset of the parameter within its block.
	Offset int
	// Name is the CLI name of the parameter (e.g. "level", "pan",
	// "octave").
	Name string
}

// PatchCommonParam returns the Param descriptor of a patch common block
// parameter with the given CLI name and byte offset, e.g.
// PatchCommonParam("level", 0x2E).
func PatchCommonParam(name string, offset int) Param {
	return Param{Block: "common", Offset: offset, Name: name}
}

// PatchToneParam returns the Param descriptor of a patch tone parameter.
func PatchToneParam(name string, tone, offset int) Param {
	return Param{Block: "tone", Index: tone, Offset: offset, Name: name}
}

// Controller is the device-facing side of the GUI. The GUI layer never
// talks to MIDI directly; it calls Get to read a parameter and Set to
// change one, and receives Update notifications for parameters that
// changed on the device itself (front panel, other software, ...). This
// keeps the GUI widgets completely independent of the device layer.
//
// All methods are called from the GUI goroutine; implementations must
// marshal calls into their own goroutine when talking to the device and
// deliver Update calls back on the GUI goroutine.
type Controller interface {
	// Get reads the current raw value (1 or 2 bytes) of the parameter.
	Get(p Param) ([]byte, error)
	// Set writes the raw value of a single parameter to the device
	// immediately.
	Set(p Param, value []byte) error
	// Update registers a callback that is invoked with the Param and its
	// new raw value whenever the device reports a changed parameter
	// (unsolicited DT1). It returns a function that unregisters it.
	Update(func(p Param, value []byte)) (unregister func())
}
