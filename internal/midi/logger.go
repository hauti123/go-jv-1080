package midi

import "go-jv-1080/internal/util"

// LogEnabled reports whether raw MIDI logging for the given direction
// ("RX" = received from the device, "TX" = sent to the device) is
// enabled. It is controlled by the LOG_MIDI_RX / LOG_MIDI_TX environment
// variables and defaults to on (set e.g. LOG_MIDI_RX=0 to disable).
func LogEnabled(direction string) bool {
	switch direction {
	case "RX":
		return util.EnvFlag("LOG_MIDI_RX", true)
	case "TX":
		return util.EnvFlag("LOG_MIDI_TX", true)
	default:
		return true
	}
}
