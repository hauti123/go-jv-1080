package util

import (
	"os"
	"strings"
)

// EnvFlag reads a boolean flag from an environment variable. If the
// variable is unset or empty, def is returned. The values "0", "false",
// "off" and "no" (case-insensitive) disable the flag; any other value
// enables it.
func EnvFlag(name string, def bool) bool {
	v, ok := os.LookupEnv(name)
	if !ok || v == "" {
		return def
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "0", "false", "off", "no":
		return false
	default:
		return true
	}
}
