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

// patchParam is a patch parameter resolved from a CLI name: the block it
// lives in ("common" or "tone"), the tone number (1-4, 0 for common), the
// offset and size within the block, and the CLI name it was resolved from.
type patchParam struct {
	block  string
	tone   int
	offset int
	size   int
	name   string
}

// commonParams maps the CLI names of doc/patch-edit-mode-structure.md to
// patch common block parameters (*1-3-1). Names must be unique only within
// the block, not across categories. Values outside the documented range
// are accepted on purpose (see parseValues) to allow special tests.
var commonParams = map[string]param{
	"name":           {0x00, 12}, // 12 character name, write raw bytes
	"efxtype":        {0x0C, 1},  // EFX type 0-39
	"efx":            {0x0C, 1},  // legacy alias of efxtype
	"efxprm1":        {0x0D, 1},
	"efxprm2":        {0x0E, 1},
	"efxprm3":        {0x0F, 1},
	"efxprm4":        {0x10, 1},
	"efxprm5":        {0x11, 1},
	"efxprm6":        {0x12, 1},
	"efxprm7":        {0x13, 1},
	"efxprm8":        {0x14, 1},
	"efxprm9":        {0x15, 1},
	"efxprm10":       {0x16, 1},
	"efxprm11":       {0x17, 1},
	"efxprm12":       {0x18, 1},
	"efxout":         {0x19, 1},
	"efxlevel":       {0x1A, 1},
	"efxchorus":      {0x1B, 1},
	"efxreverb":      {0x1C, 1},
	"efxctrl1src":    {0x1D, 1},
	"efxctrl1depth":  {0x1E, 1},
	"efxctrl2src":    {0x1F, 1},
	"efxctrl2depth":  {0x20, 1},
	"choruslevel":    {0x21, 1},
	"chorusrate":     {0x22, 1},
	"chorusdepth":    {0x23, 1},
	"choruspredelay": {0x24, 1},
	"chorusfeedback": {0x25, 1},
	"chorusout":      {0x26, 1},
	"reverbtype":     {0x27, 1},
	"reverblevel":    {0x28, 1},
	"reverbtime":     {0x29, 1},
	"hfdamp":         {0x2A, 1},
	"reverbfeedback": {0x2B, 1},
	"tempo":          {0x2C, 2}, // default tempo 20-250, nibblized
	"level":          {0x2E, 1}, // patch level 0-127
	"patchlevel":     {0x2E, 1}, // legacy alias of level
	"pan":            {0x2F, 1}, // patch pan 0-127, 64 = center; l13/r13 forms
	"patchpan":       {0x2F, 1}, // legacy alias of pan
	"analogfeel":     {0x30, 1},
	"bendup":         {0x31, 1},
	"benddown":       {0x32, 1},
	"assign":         {0x33, 1},
	"legato":         {0x34, 1},
	"portamento":     {0x35, 1},
	"portmode":       {0x36, 1},
	"porttype":       {0x37, 1},
	"portstart":      {0x38, 1},
	"porttime":       {0x39, 1},
	"efxhold":        {0x3C, 1},
	"velrange":       {0x40, 1},
	"octave":         {0x41, 1},
	"stretch":        {0x42, 1},
	"priority":       {0x43, 1},
	"struct12":       {0x44, 1},
	"booster12":      {0x45, 1},
	"struct34":       {0x46, 1},
	"booster34":      {0x47, 1},
}

// commonIndexedParams maps the CLI names of patch common parameters that
// need a control index in the command (hold/peak of control 1-3, patch
// control source 2-3) to their base offset; the index is added to it.
var commonIndexedParams = map[string]struct {
	base  int
	min   int
	max   int
	usage string
}{
	"hold": {0x3C, 1, 3, "set patch ctrl <1-3> hold <val>"},
	"src":  {0x38, 2, 3, "set patch ctrl <2-3> src <val>"},
}

// toneParams maps the CLI names of doc/patch-edit-mode-structure.md to
// patch tone block parameters (*1-3-2).
var toneParams = map[string]param{
	"switch":       {0x00, 1}, // tone on/off 0-1
	"group":        {0x01, 1}, // wave group type 0-2 (INT, PCM card, EXP)
	"waveGroup":    {0x01, 1}, // legacy alias of group
	"groupid":      {0x02, 1}, // wave group ID
	"number":       {0x03, 2}, // wave number 1-255 (raw 0-254), nibblized
	"wavenumber":   {0x03, 2},
	"wave":         {0x03, 2}, // legacy alias of number
	"gain":         {0x05, 1}, // wave gain -6..+6 dB (raw 58..70)
	"fxmswitch":    {0x06, 1},
	"fxmcolor":     {0x07, 1},
	"fxmdepth":     {0x08, 1},
	"delaymode":    {0x09, 1},
	"delaytime":    {0x0A, 1},
	"velfade":      {0x0B, 1},
	"vellower":     {0x0C, 1},
	"velupper":     {0x0D, 1},
	"keylower":     {0x0E, 1},
	"keyupper":     {0x0F, 1},
	"redamper":     {0x10, 1},
	"rxvolume":     {0x11, 1},
	"rxhold":       {0x12, 1},
	"rxbend":       {0x13, 1},
	"rxpan":        {0x14, 1},
	"coarse":       {0x3D, 1}, // coarse tune -48..+48
	"fine":         {0x3E, 1}, // fine tune -50..+50
	"random":       {0x3F, 1},
	"keyfollow":    {0x40, 1}, // pitch keyfollow
	"level":        {0x65, 1}, // tone level 0-127
	"pan":          {0x77, 1}, // tone pan 0-127, 64 = center
	"bias":         {0x68, 1},
	"biaspoint":    {0x67, 1},
	"biasdir":      {0x66, 1},
	"pankeyfollow": {0x78, 1},
	"randompan":    {0x79, 1},
	"altpan":       {0x7A, 1},
	"output":       {0x7D, 1},
	"outputlevel":  {0x7E, 1},
	"chorus":       {0x7F, 1}, // chorus send 0-127
	"reverb":       {0x80, 1}, // reverb send 0-127 (address 01 00)
}

// lfoParams holds the CLI names of the LFO1/LFO2 parameters per LFO
// (index 1 or 2) within a patch tone block.
var lfoParams = [3]map[string]param{
	1: {
		"wave":     {0x2D, 1},
		"keytrig":  {0x2E, 1},
		"rate":     {0x2F, 1},
		"offset":   {0x30, 1},
		"delay":    {0x31, 1},
		"fade":     {0x32, 1},
		"fadetime": {0x33, 1},
		"sync":     {0x34, 1},
	},
	2: {
		"wave":     {0x35, 1},
		"keytrig":  {0x36, 1},
		"rate":     {0x37, 1},
		"offset":   {0x38, 1},
		"delay":    {0x39, 1},
		"fade":     {0x3A, 1},
		"fadetime": {0x3B, 1},
		"sync":     {0x3C, 1},
	},
}

// lfoDepthParams maps the LFO depth 1:2 CLI names (pitch, filter, amp, pan)
// to their base offsets within the patch tone block; the LFO index minus 1
// is added to the base.
var lfoDepthParams = map[string]param{
	"pitch":  {0x4E, 1},
	"filter": {0x63, 1},
	"amp":    {0x75, 1},
	"pan":    {0x7B, 1},
}

// pitchEnvParams maps the CLI names of the pitch envelope parameters (PCH
// ENV DPT, PCH TIME ENV, PCH ENVELOPE) within a patch tone block.
var pitchEnvParams = map[string]param{
	"depth":     {0x41, 1},
	"velsense":  {0x42, 1},
	"vt1":       {0x43, 1},
	"vt4":       {0x44, 1},
	"keyfollow": {0x45, 1},
	"t1":        {0x46, 1},
	"t2":        {0x47, 1},
	"t3":        {0x48, 1},
	"t4":        {0x49, 1},
	"l1":        {0x4A, 1},
	"l2":        {0x4B, 1},
	"l3":        {0x4C, 1},
	"l4":        {0x4D, 1},
}

// tvfParams maps the CLI names of the TVF parameters (FILTER, TVF
// VELOCITY) within a patch tone block.
var tvfParams = map[string]param{
	"filtertype":      {0x50, 1},
	"cutoff":          {0x51, 1},
	"cutoffkeyfollow": {0x52, 1},
	"resonance":       {0x53, 1},
	"velres":          {0x54, 1},
	"envdepth":        {0x55, 1},
	"velcurve":        {0x56, 1},
	"velsens":         {0x57, 1},
}

// tvfEnvParams maps the CLI names of the TVF envelope parameters (TVF TIME
// ENV, TVF ENVELOPE) within a patch tone block.
var tvfEnvParams = map[string]param{
	"vt1":       {0x58, 1},
	"vt4":       {0x59, 1},
	"keyfollow": {0x5A, 1},
	"t1":        {0x5B, 1},
	"t2":        {0x5C, 1},
	"t3":        {0x5D, 1},
	"t4":        {0x5E, 1},
	"l1":        {0x5F, 1},
	"l2":        {0x60, 1},
	"l3":        {0x61, 1},
	"l4":        {0x62, 1},
}

// tvfEnvOffsets are the offsets of the TVF envelope parameters
// t1, t2, t3, t4, l1, l2, l3, l4 within the tone block.
var tvfEnvOffsets = [8]int{0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60, 0x61, 0x62}

// tvaParams maps the CLI names of the TVA parameters (TVA, BIAS, PAN
// MODULATION) within a patch tone block.
var tvaParams = map[string]param{
	"level":        {0x65, 1},
	"pan":          {0x77, 1},
	"envvel":       {0x6A, 1},
	"envcurve":     {0x69, 1},
	"bias":         {0x68, 1},
	"biaspoint":    {0x67, 1},
	"biasdir":      {0x66, 1},
	"pankeyfollow": {0x78, 1},
	"randompan":    {0x79, 1},
	"altpan":       {0x7A, 1},
}

// tvaEnvParams maps the CLI names of the TVA envelope parameters (TVA TIME
// ENV, TVA ENVELOPE) within a patch tone block.
var tvaEnvParams = map[string]param{
	"vt1":       {0x6B, 1},
	"vt4":       {0x6C, 1},
	"keyfollow": {0x6D, 1},
	"t1":        {0x6E, 1},
	"t2":        {0x6F, 1},
	"t3":        {0x70, 1},
	"t4":        {0x71, 1},
	"l1":        {0x72, 1},
	"l2":        {0x73, 1},
	"l3":        {0x74, 1},
}

// ctrlToneNames are the per-Tone receive-switch names of the CONTROL
// category (RxSWITCH / DAMPER).
var ctrlToneNames = map[string]bool{
	"redamper": true, "rxvolume": true, "rxhold": true,
	"rxbend": true, "rxpan": true,
}

// ctrlToneIndexedNames are the per-Tone CONTROL 1-3 destination and depth
// names that need a control index (1-3) and a tone number.
var ctrlToneIndexedNames = map[string]bool{
	"dst1": true, "dst2": true, "dst3": true, "dst4": true,
	"dep1": true, "dep2": true, "dep3": true, "dep4": true,
}

// ctrlToneOffset returns the offset of a CONTROL 1-3 destination or depth
// parameter within a patch tone block: controller i (1-3), destination or
// depth number k (1-4).
func ctrlToneOffset(name string, i, k int) int {
	base := 0x15 + (i-1)*8 + (k-1)*2
	if strings016(name, "dep") {
		return base + 1
	}
	return base
}

func strings016(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// waveNames are the WAVE category parameter names (tone switch, wave
// group/number/gain, FXM, tone delay).
var waveNames = map[string]bool{
	"switch": true, "group": true, "groupid": true, "number": true,
	"wavenumber": true, "wave": true, "gain": true, "fxmswitch": true,
	"fxmcolor": true, "fxmdepth": true, "delaymode": true, "delaytime": true,
}

func lookupCommon(name string) (patchParam, error) {
	p, ok := commonParams[normalize(name)]
	if !ok {
		return patchParam{}, fmt.Errorf("unknown common parameter %q", name)
	}
	return patchParam{block: "common", offset: p.offset, size: p.size, name: normalize(name)}, nil
}

func lookupCommonIndexed(name string, index int) (patchParam, error) {
	c, ok := commonIndexedParams[normalize(name)]
	if !ok {
		return patchParam{}, fmt.Errorf("unknown common parameter %q", name)
	}
	if index < c.min || index > c.max {
		return patchParam{}, fmt.Errorf("control index %d out of range for %q (want %d-%d; usage: %s)", index, name, c.min, c.max, c.usage)
	}
	return patchParam{block: "common", offset: c.base + index, size: 1, name: normalize(name)}, nil
}

func lookupTone(name string, tone int) (patchParam, error) {
	p, ok := toneParams[normalize(name)]
	if !ok {
		return patchParam{}, fmt.Errorf("unknown tone parameter %q", name)
	}
	if tone < 1 || tone > 4 {
		return patchParam{}, fmt.Errorf("invalid tone %d (want 1-4)", tone)
	}
	return patchParam{block: "tone", tone: tone, offset: p.offset, size: p.size, name: normalize(name)}, nil
}

func commonParam(name string) (offset, size int, err error) {
	p, err := lookupCommon(name)
	if err != nil {
		return 0, 0, err
	}
	return p.offset, p.size, nil
}

func toneParam(name string) (offset, size int, err error) {
	p, err := lookupTone(name, 1)
	if err != nil {
		return 0, 0, err
	}
	return p.offset, p.size, nil
}

func waveParam(name string) (offset, size int, err error) {
	if !waveNames[normalize(name)] {
		return 0, 0, fmt.Errorf("unknown wave parameter %q (want switch, group, groupid, number, gain, fxmswitch, fxmcolor, fxmdepth, delaymode, delaytime)", name)
	}
	p, err := lookupTone(name, 1)
	if err != nil {
		return 0, 0, err
	}
	return p.offset, p.size, nil
}

func tvfParam(name string) (offset, size int, err error) {
	p, ok := tvfParams[normalize(name)]
	if !ok {
		return 0, 0, fmt.Errorf("unknown tvf parameter %q (want filtertype, cutoff, cutoffkeyfollow, resonance, velres, envdepth, velcurve, velsens)", name)
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
