package main

import (
	"fmt"
	"strconv"
)

// param describes a named parameter: its offset within its block and the
// number of raw bytes it occupies (1, or 2 for nibblized values).
type param struct {
	offset int
	size   int
}

// commonParams maps speaking names to patch common block parameters.
// Values outside the documented range are accepted on purpose (see
// parseValues) to allow special tests.
var commonParams = map[string]param{
	"name":       {0x00, 12}, // 12 character name, write raw bytes
	"efx":        {0x0C, 1},  // EFX type 0-39
	"efxtype":    {0x0C, 1},
	"level":      {0x2E, 1}, // patch level 0-127
	"patchlevel": {0x2E, 1},
	"pan":        {0x2F, 1}, // patch pan 0-127, 64 = center; l13/r13 forms below
	"patchpan":   {0x2F, 1},
}

// toneParams maps speaking names to patch tone block parameters.
var toneParams = map[string]param{
	"switch":     {0x00, 1}, // tone on/off 0-1
	"group":      {0x01, 1}, // wave group type 0-2 (INT, PCM card, EXP)
	"waveGroup":  {0x01, 1},
	"groupid":    {0x02, 1}, // wave group ID
	"wave":       {0x03, 2}, // wave number 1-255, nibblized
	"wavenumber": {0x03, 2},
	"coarse":     {0x3D, 1}, // coarse tune -48..+48
	"fine":       {0x3E, 1}, // fine tune -50..+50
	"cutoff":     {0x51, 1}, // TVF cutoff 0-127
	"resonance":  {0x53, 1}, // TVF resonance 0-127
	"level":      {0x65, 1}, // tone level 0-127
	"pan":        {0x77, 1}, // tone pan 0-127, 64 = center
	"chorus":     {0x7F, 1}, // chorus send 0-127
	"reverb":     {0x80, 1}, // reverb send 0-127
}

// waveParams maps speaking names to wave parameters (tone block).
var waveParams = map[string]param{
	"group":   {0x01, 1}, // wave group type 0-2
	"groupid": {0x02, 1}, // wave group ID
	"number":  {0x03, 2}, // wave number 1-255, nibblized
	"gain":    {0x05, 1}, // wave gain -6..+6 dB (raw 58..70)
}

// tvfParams maps speaking names to TVF parameters (tone block).
var tvfParams = map[string]param{
	"cutoff":    {0x51, 1}, // 0-127
	"resonance": {0x53, 1}, // 0-127
	"envdepth":  {0x54, 1}, // TVF envelope depth -63..+63 (raw 1..127)
	// env: t1 t2 t3 t4 l1 l2 l3 l4 at 0x55..0x5C
}

// tvfEnvOffsets are the offsets of the TVF envelope parameters
// t1, t2, t3, t4, l1, l2, l3, l4 within the tone block.
var tvfEnvOffsets = [8]int{0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C}

func commonParam(name string) (offset, size int, err error) {
	p, ok := commonParams[normalize(name)]
	if !ok {
		return 0, 0, fmt.Errorf("unknown common parameter %q", name)
	}
	return p.offset, p.size, nil
}

func toneParam(name string) (offset, size int, err error) {
	p, ok := toneParams[normalize(name)]
	if !ok {
		return 0, 0, fmt.Errorf("unknown tone parameter %q", name)
	}
	return p.offset, p.size, nil
}

func waveParam(name string) (offset, size int, err error) {
	p, ok := waveParams[normalize(name)]
	if !ok {
		return 0, 0, fmt.Errorf("unknown wave parameter %q (want group, groupid, number, gain)", name)
	}
	return p.offset, p.size, nil
}

func tvfParam(name string) (offset, size int, err error) {
	p, ok := tvfParams[normalize(name)]
	if !ok {
		return 0, 0, fmt.Errorf("unknown tvf parameter %q (want cutoff, resonance, envdepth, env)", name)
	}
	return p.offset, p.size, nil
}

func normalize(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' || c == '_' || c == ' ' {
			continue
		}
		out = append(out, c)
	}
	return string(out)
}

// panValue converts pan notations to a raw pan byte: plain numbers 0-127,
// l1..l63 (left), r1..r63 (right) and center/c.
func panValue(s string) (int, bool) {
	switch s {
	case "center", "c":
		return 64, true
	}
	if len(s) >= 2 && (s[0] == 'l' || s[0] == 'r') {
		n, err := strconv.Atoi(s[1:])
		if err != nil || n < 1 || n > 63 {
			return 0, false
		}
		if s[0] == 'l' {
			return 64 - n, true
		}
		return 64 + n, true
	}
	return 0, false
}
