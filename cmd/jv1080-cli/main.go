// Command jv1080-cli is a command line utility to control a Roland JV-1080:
// switch the panel mode, select patches, manipulate patch parameters, read
// and write raw addresses and play notes. All messages to and from the
// device are logged to the command line in hex.
//
// Usage:
//
//	jv1080-cli [flags] <command> [args]
//
// Flags:
//
//	-out string   MIDI output port name (default "SERIES 208i MIDI")
//	-in string    MIDI input port name (default "SERIES 208i MIDI")
//	-outnum int   MIDI output port number, used when -out is empty (default 0)
//	-innum int    MIDI input port number, used when -in is empty (default 0)
//	-devid string device ID in hex (default "10" = unit 17)
//
// Commands (run "jv1080-cli help" for the full list):
//
//	set panelmode <patch|performance|gm|1|2|3>
//	set patch select <bank> <number>
//	set patch <category> ... (see doc/patch-edit-mode-structure.md)
//	set address <hi> <lo> <value...>
//	get patch <category> ...
//	get address <hi> <lo> [size]
//	note <name...>
//	help
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go-jv-1080/internal/jv1080"
	"go-jv-1080/internal/midi"
	"go-jv-1080/internal/midi/message"
)

var (
	outName  = flag.String("out", "SERIES 208i MIDI", "MIDI output port name")
	inName   = flag.String("in", "SERIES 208i MIDI", "MIDI input port name")
	outNum   = flag.Int("outnum", 0, "MIDI output port number (used when -out is empty)")
	inNum    = flag.Int("innum", 0, "MIDI input port number (used when -in is empty)")
	deviceID = flag.String("devid", "10", "device ID in hex")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: jv1080-cli [flags] <command> [args] (try: help)")
		os.Exit(2)
	}
	// help works without a device.
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printHelp()
		return
	}
	dev, err := midi.OpenDevice(*outName, *inName, *outNum, *inNum)
	if err != nil {
		fatal("cannot open MIDI device: %v", err)
	}
	defer dev.Close()
	// Log all incoming and outgoing messages in hex.
	dev.RegisterCallback(func(msg []byte) {
		fmt.Printf("REC : [%s]\n", message.Hex(msg))
	})
	dev.SetSendLogger(func(msg []byte) {
		fmt.Printf("SEND: [%s]\n", message.Hex(msg))
	})

	id, err := strconv.ParseUint(*deviceID, 16, 8)
	if err != nil {
		fatal("invalid device ID %q: %v", *deviceID, err)
	}
	client := jv1080.NewClient(dev, uint8(id))

	if err := run(client, args); err != nil {
		fatal("%v", err)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "jv1080-cli: "+format+"\n", args...)
	os.Exit(1)
}

func run(c *jv1080.Client, args []string) error {
	switch args[0] {
	case "set":
		return runSet(c, args[1:])
	case "get":
		return runGet(c, args[1:])
	case "note":
		return runNote(c, args[1:])
	case "help", "--help", "-h":
		printHelp()
		return nil
	}
	return fmt.Errorf("unknown command %q (want set, get, note or help)", args[0])
}

func runSet(c *jv1080.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("set what? (panelmode, patch, address) — try: help")
	}
	switch args[0] {
	case "panelmode":
		return setPanelMode(c, args[1:])
	case "patch":
		return runPatch(c, "set", args[1:])
	case "address":
		return setAddress(c, args[1:])
	}
	return fmt.Errorf("unknown set target %q (want panelmode, patch or address)", args[0])
}

func runGet(c *jv1080.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("get what? (patch, address) — try: help")
	}
	switch args[0] {
	case "patch":
		return runPatch(c, "get", args[1:])
	case "address":
		return getAddress(c, args[1:])
	}
	return fmt.Errorf("unknown get target %q (want patch or address)", args[0])
}

// patchCategories lists the patch edit categories of the Patch Edit Mode
// (doc/patch-edit-mode-structure.md) accepted after "set/get patch".
var patchCategories = []string{"common", "efx", "ctrl", "wave", "lfo", "pitch", "tvf", "tva", "tone"}

// patchCommand is a parsed "set/get patch <category> ..." command: the
// resolved parameter plus the index of its first value argument.
type patchCommand struct {
	param   patchParam
	valArgs []string
}

// runPatch dispatches the patch related subcommands: select a patch
// (set patch <bank> <number>) or manipulate/read patch parameters
// (set/get patch common|efx|ctrl|wave|lfo|pitch|tvf|tva ...).
func runPatch(c *jv1080.Client, verb string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%s patch what? (select, %s) — try: help", verb, strings.Join(patchCategories, ", "))
	}
	if args[0] == "select" {
		if verb != "set" {
			return fmt.Errorf("usage: set patch select <bank> <number>")
		}
		return setPatch(c, args[1:])
	}
	if verb != "set" && verb != "get" {
		return fmt.Errorf("unknown verb %q", verb)
	}
	for _, cat := range patchCategories {
		if args[0] != cat {
			continue
		}
		cmds, err := parsePatchParams(cat, args[1:])
		if err != nil {
			return err
		}
		if verb == "set" {
			for _, cmd := range cmds {
				values, err := parsePatchValues(cmd.param, cmd.valArgs)
				if err != nil {
					return err
				}
				if err := c.SetPatchParameter(jv1080.BankPresetA, 0, cmd.param.block, cmd.param.tone, cmd.param.offset, values); err != nil {
					return err
				}
			}
			return nil
		}
		return getPatchParams(c, cmds)
	}
	return fmt.Errorf("unknown patch target %q (want select, %s)", args[0], strings.Join(patchCategories, ", "))
}

// parsePatchParams resolves the category-specific argument list (args are
// the tokens after the category) to the parameter it addresses together
// with the value arguments.
func parsePatchParams(cat string, args []string) ([]patchCommand, error) {
	switch cat {
	case "common":
		if len(args) < 1 {
			return nil, fmt.Errorf("usage: set/get patch common <param> <values...>")
		}
		p, err := lookupCommon(args[0])
		if err != nil {
			return nil, err
		}
		return checkValues(p, args[1:]), nil
	case "efx":
		return parseToneCategory("efx", args)
	case "pitch":
		return parseToneCategory("pitch", args)
	case "tvf":
		return parseToneCategory("tvf", args)
	case "tva":
		return parseToneCategory("tva", args)
	case "ctrl":
		return parseCtrl(args)
	case "wave":
		return parseWave(args)
	case "lfo":
		return parseLFO(args)
	case "tone":
		return parsePatchTone(args)
	}
	return nil, fmt.Errorf("unknown patch category %q", cat)
}

// checkValues validates the value count against the parameter size and
// wraps the parameter in a patchCommand.
func checkValues(p patchParam, vals []string) []patchCommand {
	return []patchCommand{{param: p, valArgs: vals}}
}

// pitchToneNames are the PITCH parameter names that need no env marker.
var pitchToneNames = map[string]bool{
	"coarse": true, "fine": true, "random": true, "keyfollow": true,
}

// efxCommonNames are the EFFECTS category names in the patch common block
// (EFX type/parameters/output/control, patch chorus and patch reverb).
var efxCommonNames = map[string]bool{
	"efxtype": true, "efx": true,
	"efxprm1": true, "efxprm2": true, "efxprm3": true, "efxprm4": true,
	"efxprm5": true, "efxprm6": true, "efxprm7": true, "efxprm8": true,
	"efxprm9": true, "efxprm10": true, "efxprm11": true, "efxprm12": true,
	"efxout": true, "efxlevel": true, "efxchorus": true, "efxreverb": true,
	"efxctrl1src": true, "efxctrl1depth": true, "efxctrl2src": true, "efxctrl2depth": true,
	"choruslevel": true, "chorusrate": true, "chorusdepth": true,
	"choruspredelay": true, "chorusfeedback": true, "chorusout": true,
	"reverbtype": true, "reverblevel": true, "reverbtime": true,
	"hfdamp": true, "reverbfeedback": true,
}

// efxToneNames are the per-Tone OUTPUT parameter names of the EFFECTS
// category.
var efxToneNames = map[string]bool{
	"output": true, "outputlevel": true, "chorus": true, "reverb": true,
}

const efxUsage = "usage: set/get patch efx [tone <1-4>] <param> <val>"
const pitchUsage = "usage: set/get patch pitch [tone <1-4>] [<env>] <param> <val>"
const tvfUsage = "usage: set/get patch tvf [tone <1-4>] [<env>] <param> <val>"
const tvaUsage = "usage: set/get patch tva [tone <1-4>] [<env>] <param> <val>"

// parseToneCategory parses the [tone <n>] [env] <param> <val> form shared
// by the efx, pitch, tva and tvf categories.
func parseToneCategory(cat string, args []string) ([]patchCommand, error) {
	var usage string
	switch cat {
	case "efx":
		usage = efxUsage
	default:
		usage = cat + " " + strings.SplitN(pitchUsage, " ", 4)[3]
	}
	if len(args) < 2 {
		return nil, fmt.Errorf("%s", usage)
	}
	tone := 1
	toneSelected := false
	env := false
	rest := args
	if len(rest) >= 2 && rest[0] == "tone" {
		n, err := strconv.Atoi(rest[1])
		if err != nil || n < 1 || n > 4 {
			return nil, fmt.Errorf("invalid tone %q (want 1-4)", rest[1])
		}
		tone = n
		toneSelected = true
		rest = rest[2:]
	}
	if len(rest) >= 1 && rest[0] == "env" {
		env = true
		rest = rest[1:]
	}
	if len(rest) < 1 {
		return nil, fmt.Errorf("%s", usage)
	}
	name := rest[0]
	if !env && cat == "efx" {
		if tone == 1 && !toneSelected {
			if !efxCommonNames[normalize(name)] {
				return nil, fmt.Errorf("unknown efx parameter %q (common EFX parameters take no tone; per-Tone: output, outputlevel, chorus, reverb)", name)
			}
			return checkValues(patchParam{block: "common", offset: commonParams[normalize(name)].offset, size: commonParams[normalize(name)].size, name: normalize(name)}, rest[1:]), nil
		}
		if !efxToneNames[normalize(name)] {
			return nil, fmt.Errorf("unknown efx tone parameter %q (want output, outputlevel, chorus, reverb; common EFX parameters take no tone: set patch efx <param> <val>)", name)
		}
	}
	switch cat {
	case "pitch":
		if env {
			if _, ok := pitchEnvParams[normalize(name)]; !ok {
				return nil, fmt.Errorf("unknown pitch env parameter %q (want depth, velsense, vt1, vt4, keyfollow, t1-t4, l1-l4)", name)
			}
		} else if !pitchToneNames[normalize(name)] {
			return nil, fmt.Errorf("unknown pitch parameter %q (want coarse, fine, random, keyfollow)", name)
		}
	case "tvf":
		if env {
			if _, ok := tvfEnvParams[normalize(name)]; !ok {
				return nil, fmt.Errorf("unknown tvf env parameter %q (want vt1, vt4, keyfollow, t1-t4, l1-l4)", name)
			}
		} else if _, ok := tvfParams[normalize(name)]; !ok {
			return nil, fmt.Errorf("unknown tvf parameter %q (want filtertype, cutoff, cutoffkeyfollow, resonance, velres, envdepth, velcurve, velsens)", name)
		}
	case "tva":
		if env {
			if _, ok := tvaEnvParams[normalize(name)]; !ok {
				return nil, fmt.Errorf("unknown tva env parameter %q (want vt1, vt4, keyfollow, t1-t4, l1-l3)", name)
			}
		} else if _, ok := tvaParams[normalize(name)]; !ok {
			return nil, fmt.Errorf("unknown tva parameter %q (want level, pan, envvel, envcurve, bias, biaspoint, biasdir, pankeyfollow, randompan, altpan)", name)
		}
	}
	off, size, err := toneCategoryOffset(cat, name, env)
	if err != nil {
		return nil, err
	}
	return checkValues(patchParam{block: "tone", tone: tone, offset: off, size: size, name: normalize(name)}, rest[1:]), nil
}

// toneCategoryOffset resolves a non-env parameter of a tone category to
// its offset and size within the tone block.
func toneCategoryOffset(cat, name string, env bool) (int, int, error) {
	switch cat {
	case "efx":
		p, err := lookupTone(name, 1)
		if err != nil {
			return 0, 0, err
		}
		return p.offset, p.size, nil
	case "pitch":
		if env {
			return pitchEnvParams[normalize(name)].offset, 1, nil
		}
		return toneParams[normalize(name)].offset, toneParams[normalize(name)].size, nil
	case "tvf":
		if env {
			return tvfEnvParams[normalize(name)].offset, 1, nil
		}
		return tvfParams[normalize(name)].offset, tvfParams[normalize(name)].size, nil
	case "tva":
		if env {
			return tvaEnvParams[normalize(name)].offset, 1, nil
		}
		return tvaParams[normalize(name)].offset, tvaParams[normalize(name)].size, nil
	}
	return 0, 0, fmt.Errorf("unknown tone category %q", cat)
}

// parseCtrl parses the CONTROL category forms:
//
//	set/get patch ctrl <param> <val>                     (common, no index)
//	set/get patch ctrl <2-3> src <val>                   (common, index)
//	set/get patch ctrl <1-3> hold <val>                   (common, index)
//	set/get patch ctrl tone <1-4> <param> <val>           (per-Tone RxSWITCH)
//	set/get patch ctrl <1-3> tone <1-4> <dst|dep><1-4> <val>
func parseCtrl(args []string) ([]patchCommand, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("usage: set/get patch ctrl [<1-3>] [tone <1-4>] <param> <val>")
	}
	// Control 1-3 per-Tone form: <1-3> tone <1-4> <dst|dep><1-4> [<val>].
	if len(args) >= 4 && args[1] == "tone" {
		i, err := strconv.Atoi(args[0])
		if err != nil || i < 1 || i > 3 {
			return nil, fmt.Errorf("invalid control index %q (want 1-3)", args[0])
		}
		tone, err := strconv.Atoi(args[2])
		if err != nil || tone < 1 || tone > 4 {
			return nil, fmt.Errorf("invalid tone %q (want 1-4)", args[2])
		}
		name := args[3]
		if !ctrlToneIndexedNames[normalize(name)] {
			return nil, fmt.Errorf("unknown ctrl parameter %q (want dst1-dst4 or dep1-dep4)", name)
		}
		k, err := strconv.Atoi(name[len(name)-1:])
		if err != nil {
			return nil, fmt.Errorf("unknown ctrl parameter %q", name)
		}
		return checkValues(patchParam{block: "tone", tone: tone, offset: ctrlToneOffset(normalize(name), i, k), size: 1, name: normalize(name)}, args[min(5, len(args)):]), nil
	}
	// Common forms first: plain name or <index> hold|src.
	if args[0] != "tone" {
		if len(args) >= 2 && (args[1] == "hold" || args[1] == "src") {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, fmt.Errorf("invalid control index %q", args[0])
			}
			p, err := lookupCommonIndexed(args[1], n)
			if err != nil {
				return nil, err
			}
			return checkValues(p, args[min(3, len(args)):]), nil
		}
		if ctrlToneNames[normalize(args[0])] || ctrlToneIndexedNames[normalize(args[0])] {
			return nil, fmt.Errorf("parameter %q needs a tone selector: set patch ctrl [ctrl <1-3>] tone <1-4> %s <val>", args[0], args[0])
		}
		p, err := lookupCommon(args[0])
		if err != nil {
			return nil, err
		}
		return checkValues(p, args[1:]), nil
	}
	// Per-Tone forms.
	if len(args) < 3 {
		return nil, fmt.Errorf("usage: set/get patch ctrl tone <1-4> <param> [<val>]")
	}
	tone, err := strconv.Atoi(args[1])
	if err != nil || tone < 1 || tone > 4 {
		return nil, fmt.Errorf("invalid tone %q (want 1-4)", args[1])
	}
	name := args[2]
	if !ctrlToneNames[normalize(name)] && !ctrlToneIndexedNames[normalize(name)] {
		return nil, fmt.Errorf("unknown ctrl tone parameter %q (want redamper, rxvolume, rxhold, rxbend, rxpan or dst1-4, dep1-4)", name)
	}
	if !ctrlToneNames[normalize(name)] {
		return nil, fmt.Errorf("parameter %q needs a control index: set patch ctrl <1-3> tone <1-4> %s <val>", name, name)
	}
	p, err := lookupTone(name, tone)
	if err != nil {
		return nil, err
	}
	return checkValues(patchParam{block: "tone", tone: tone, offset: p.offset, size: 1, name: normalize(name)}, args[min(4, len(args)):]), nil
}

// parseWave parses the WAVE category form:
// set/get patch wave [tone <1-4>] <param> <val>.
func parseWave(args []string) ([]patchCommand, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("usage: set/get patch wave [tone <1-4>] <param> <val>")
	}
	tone := 1
	rest := args
	if rest[0] == "tone" {
		n, err := strconv.Atoi(rest[1])
		if err != nil || n < 1 || n > 4 {
			return nil, fmt.Errorf("invalid tone %q (want 1-4)", rest[1])
		}
		tone = n
		rest = rest[2:]
	}
	if len(rest) < 1 {
		return nil, fmt.Errorf("usage: set/get patch wave [tone <1-4>] <param> <val>")
	}
	name := rest[0]
	if !waveNames[normalize(name)] {
		return nil, fmt.Errorf("unknown wave parameter %q (want switch, group, groupid, number, gain, fxmswitch, fxmcolor, fxmdepth, delaymode, delaytime)", name)
	}
	p, err := lookupTone(name, tone)
	if err != nil {
		return nil, err
	}
	return checkValues(p, rest[1:]), nil
}

// parsePatchTone parses the generic per-Tone category form
// set/get patch tone <1-4> <param> <val> for the Velocity range and Key
// range parameters (vellower, velupper, velfade, keylower, keyupper).
func parsePatchTone(args []string) ([]patchCommand, error) {
	if len(args) < 3 || args[0] != "tone" {
		return nil, fmt.Errorf("usage: set/get patch tone <1-4> <param> <val>")
	}
	tone, err := strconv.Atoi(args[1])
	if err != nil || tone < 1 || tone > 4 {
		return nil, fmt.Errorf("invalid tone %q (want 1-4)", args[1])
	}
	name := args[2]
	if !toneRangeNames[normalize(name)] {
		return nil, fmt.Errorf("unknown tone parameter %q (want vellower, velupper, velfade, keylower, keyupper)", name)
	}
	p, err := lookupTone(name, tone)
	if err != nil {
		return nil, err
	}
	return checkValues(p, args[3:]), nil
}

// toneRangeNames are the per-Tone Velocity range and Key range parameter
// names (doc/patch-edit-mode-structure.md).
var toneRangeNames = map[string]bool{
	"vellower": true, "velupper": true, "velfade": true,
	"keylower": true, "keyupper": true,
}

// parseLFO parses the LFO category forms:
//
//	set/get patch lfo <1-2> tone <1-4> <param> <val>   (LFO1/LFO2)
//	set/get patch lfo <1-2> <param> <val>              (LFO depth 1:2)
func parseLFO(args []string) ([]patchCommand, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("usage: set/get patch lfo <1-2> [tone <1-4>] <param> <val>")
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 1 || n > 2 {
		return nil, fmt.Errorf("invalid LFO number %q (want 1 or 2)", args[0])
	}
	rest := args[1:]
	if len(rest) >= 2 && rest[0] == "tone" {
		t, err := strconv.Atoi(rest[1])
		if err != nil || t < 1 || t > 4 {
			return nil, fmt.Errorf("invalid tone %q (want 1-4)", rest[1])
		}
		if len(rest) < 3 {
			return nil, fmt.Errorf("usage: set/get patch lfo <1-2> tone <1-4> <param> <val>")
		}
		name := rest[2]
		p, ok := lfoParams[n][normalize(name)]
		if !ok {
			return nil, fmt.Errorf("unknown lfo parameter %q (want wave, keytrig, rate, offset, delay, fade, fadetime, sync)", name)
		}
		return checkValues(patchParam{block: "tone", tone: t, offset: p.offset, size: p.size, name: normalize(name)}, rest[min(4, len(rest)):]), nil
	}
	if len(rest) < 1 {
		return nil, fmt.Errorf("usage: set/get patch lfo <1-2> <param> <val>")
	}
	name := rest[0]
	p, ok := lfoDepthParams[normalize(name)]
	if !ok {
		return nil, fmt.Errorf("unknown lfo depth parameter %q (want pitch, filter, amp, pan)", name)
	}
	return checkValues(patchParam{block: "tone", tone: 1, offset: p.offset + n - 1, size: p.size, name: normalize(name)}, rest[1:]), nil
}

// parsePatchValues parses the values of a patch parameter, applying
// speaking notations where they exist (pan l13/r13/center, signed
// coarse/fine).
func parsePatchValues(p patchParam, args []string) ([]byte, error) {
	switch p.name {
	case "pan":
		if len(args) == 1 {
			if v, ok := panValue(strings.ToLower(args[0])); ok {
				return []byte{byte(v)}, nil
			}
		}
	case "coarse":
		if len(args) == 1 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				return []byte{byte(v + 48)}, nil
			}
		}
	case "fine":
		if len(args) == 1 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				return []byte{byte(v + 50)}, nil
			}
		}
	}
	return parseValues(args, p.size)
}

// getPatchParams reads the parameters of a parsed patch command from the
// temporary patch and prints them.
func getPatchParams(c *jv1080.Client, cmds []patchCommand) error {
	for _, cmd := range cmds {
		pp := cmd.param
		var area message.Address
		if pp.block == "common" {
			area = jv1080.TemporaryPatchArea().Add(pp.offset)
		} else {
			area = jv1080.PatchToneArea(jv1080.TemporaryPatchArea(), pp.tone).Add(pp.offset)
		}
		data, err := c.ReadRaw(area, pp.size)
		if err != nil {
			return err
		}
		label := pp.name
		if pp.block == "tone" {
			label = fmt.Sprintf("tone %d %s", pp.tone, pp.name)
		}
		fmt.Printf("%s: %s\n", label, formatValue(pp.name, data))
	}
	return nil
}

// printHelp displays all implemented commands.
func printHelp() {
	fmt.Print(helpText)
}

const helpText = `jv1080-cli — control a Roland JV-1080 over MIDI

Usage: jv1080-cli [flags] <command> [args]

Flags:
  -out string    MIDI output port name (default "SERIES 208i MIDI")
  -in string     MIDI input port name (default "SERIES 208i MIDI")
  -outnum int    MIDI output port number, used when -out is empty (default 0)
  -innum int     MIDI input port number, used when -in is empty (default 0)
  -devid string  device ID in hex (default "10" = unit 17)

Commands:
  set panelmode <patch|performance|gm|1|2|3>
        switch the panel mode
  set patch select <bank> <number>
        select a patch; bank: user, pr-a, pr-b, pr-c, gm, card-a, card-b,
        exp-a..d; number 1-128
  set/get patch common <param> <value...>
        set a patch common parameter of the temporary patch:
        name (12 values), efxtype, efxprm1-12, efxout, efxlevel, efxchorus,
        efxreverb, efxctrl1src, efxctrl1depth, efxctrl2src, efxctrl2depth,
        choruslevel, chorusrate, chorusdepth, choruspredelay, chorusfeedback,
        chorusout, reverbtype, reverblevel, reverbtime, hfdamp, reverbfeedback,
        tempo (2 values), level, pan, analogfeel, bendup, benddown, assign,
        legato, portamento, portmode, porttype, portstart, porttime, efxhold,
        velrange, octave, stretch, priority, struct12, booster12, struct34,
        booster34
  set/get patch efx [tone <1-4>] <param> <val>
        EFX: common parameters take no tone selector. Per-Tone output
        routing: output, outputlevel, chorus, reverb
  set/get patch ctrl [<1-3>] [tone <1-4>] <param> <val>
        CONTROL: key mode and bender (assign, legato, bendup, benddown),
        portamento (portamento, portmode, porttype, portstart, porttime),
        control sources (<2-3> src), hold/peak (<1-3> hold),
        per-Tone receive switches (tone <n> redamper, rxvolume, rxhold,
        rxbend, rxpan) and per-Tone control destinations/depths
        (<1-3> tone <n> dst1-4, dep1-4)
  set/get patch wave [tone <1-4>] <param> <val>
        WAVE: switch, group, groupid, number (2 values), gain, fxmswitch,
        fxmcolor, fxmdepth, delaymode, delaytime
  set/get patch lfo <1-2> [tone <1-4>] <param> <val>
        LFO: per-Tone wave, keytrig, rate, offset, delay, fade, fadetime,
        sync; LFO depth 1:2 without tone: pitch, filter, amp, pan
  set/get patch pitch [tone <1-4>] [env] <param> <val>
        PITCH: coarse, fine, random, keyfollow; with env: depth, velsense,
        vt1, vt4, keyfollow, t1-t4, l1-l4
  set/get patch tvf [tone <1-4>] [env] <param> <val>
        TVF: filtertype, cutoff, cutoffkeyfollow, resonance, velres,
        envdepth, velcurve, velsens; with env: vt1, vt4, keyfollow,
        t1-t4, l1-l4
  set/get patch tva [tone <1-4>] [env] <param> <val>
        TVA: level, pan, envvel, envcurve, bias, biaspoint, biasdir,
        pankeyfollow, randompan, altpan; with env: vt1, vt4, keyfollow,
        t1-t4, l1-l3
  set/get patch tone <1-4> <param> <val>
        per-Tone Velocity and Key range: vellower, velupper, velfade,
        keylower, keyupper
  set address <hi> <lo> <value...>
        write raw bytes to the address given as two 16-bit hex numbers
        (e.g. 0x0000 0x0006 = master tune)
  get address <hi> <lo> [size]
        read raw bytes from an address (default size 1)
  note <name...>
        play notes by name (c1 d3 f#2 bb3 ab3 eb3 ...), 500 ms at full
        velocity with a 500 ms pause between notes
  help
        show this help

The patch categories follow the Patch Edit Mode structure described in
doc/patch-edit-mode-structure.md. Parameter names are unique within their
block (common or tone); the same concept uses the same name in both.

All messages to and from the device are logged in hex:
  SEND: [F0 41 10 6A ... F7]
  REC : [F0 41 10 6A ... F7]

Parameter values may be given in decimal or hex (0x..). Values outside
the documented range are accepted deliberately, to allow special tests.
Pan accepts l1..l63, r1..r63 and center; coarse/fine accept signed
values; nibblized parameters (number, tempo) take two nibble bytes.
`

// setPanelMode switches the panel mode: patch (1), performance (2) or GM (3).
func setPanelMode(c *jv1080.Client, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: set panelmode <patch|performance|gm|1|2|3>")
	}
	var mode uint8
	switch strings.ToLower(args[0]) {
	case "patch", "1":
		mode = 1
	case "performance", "2":
		mode = 2
	case "gm", "3":
		mode = 3
	default:
		return fmt.Errorf("unknown panel mode %q (want patch, performance, gm, 1, 2 or 3)", args[0])
	}
	return c.SetPanelMode(mode)
}

// setPatch selects a patch: bank name (user, pr-a, pr-b, pr-c, gm, card-a,
// card-b, exp-a..d) and number (1-128).
func setPatch(c *jv1080.Client, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: set patch select <bank> <number>")
	}
	bank, err := parseBank(args[0])
	if err != nil {
		return err
	}
	number, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid patch number %q", args[1])
	}
	return c.SelectPatch(0, bank, number)
}

func parseBank(s string) (jv1080.Bank, error) {
	switch strings.ToLower(s) {
	case "user":
		return jv1080.BankUser, nil
	case "pr-a", "preset-a", "preseta":
		return jv1080.BankPresetA, nil
	case "pr-b", "preset-b", "presetb":
		return jv1080.BankPresetB, nil
	case "pr-c", "preset-c", "presetc":
		return jv1080.BankPresetC, nil
	case "gm":
		return jv1080.BankGM, nil
	case "card-a", "carda":
		return jv1080.BankCardA, nil
	case "card-b", "cardb":
		return jv1080.BankCardB, nil
	case "exp-a", "expa":
		return jv1080.BankExpA, nil
	case "exp-b", "expb":
		return jv1080.BankExpB, nil
	case "exp-c", "expc":
		return jv1080.BankExpC, nil
	case "exp-d", "expd":
		return jv1080.BankExpD, nil
	}
	return 0, fmt.Errorf("unknown bank %q (want user, pr-a, pr-b, pr-c, gm, card-a, card-b, exp-a..d)", s)
}

// setAddress writes raw bytes to the address given as two 16-bit hex
// numbers (hi lo), e.g. 0x0000 0x0006 for the master tune.
func setAddress(c *jv1080.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: set address <hi> <lo> <value...>")
	}
	addr, err := parseAddress(args[0], args[1])
	if err != nil {
		return err
	}
	values, err := parseValues(args[2:], 1)
	if err != nil {
		return err
	}
	return c.WriteRaw(addr, values)
}

// getAddress reads size bytes (default 1) from the address given as two
// 16-bit hex numbers.
func getAddress(c *jv1080.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: get address <hi> <lo> [size]")
	}
	addr, err := parseAddress(args[0], args[1])
	if err != nil {
		return err
	}
	size := 1
	if len(args) > 2 {
		size, err = strconv.Atoi(args[2])
		if err != nil || size < 1 {
			return fmt.Errorf("invalid size %q", args[2])
		}
	}
	data, err := c.ReadRaw(addr, size)
	if err != nil {
		return err
	}
	fmt.Printf("address %s: [%s]\n", addr, message.Hex(data))
	return nil
}

// formatValue formats a raw parameter value in decimal, applying speaking
// notations where they exist (pan l/r/center, signed coarse/fine).
func formatValue(name string, data []byte) string {
	if len(data) == 0 {
		return ""
	}
	if len(data) == 2 {
		return strconv.Itoa(rawValue(data))
	}
	if len(data) > 2 {
		return message.Hex(data)
	}
	switch normalize(name) {
	case "pan":
		v := int(data[0])
		switch {
		case v == 64:
			return "center (64)"
		case v < 64:
			return fmt.Sprintf("l%d (%d)", 64-v, v)
		default:
			return fmt.Sprintf("r%d (%d)", v-64, v)
		}
	case "coarse":
		return fmt.Sprintf("%d (%d)", int(data[0])-48, data[0])
	case "fine":
		return fmt.Sprintf("%d (%d)", int(data[0])-50, data[0])
	}
	return strconv.Itoa(int(data[0]))
}

// rawValue reads the raw value of a parameter: the plain byte, or the
// combined number for nibblized (2-byte) parameters.
func rawValue(data []byte) int {
	if len(data) == 2 {
		return int(data[0])<<4 | int(data[1])
	}
	return int(data[0])
}

// runNote plays notes given by name (c1 d3 c2 d4 e3 f#2 b3 bb3 ab3 eb3 ...),
// each 500 ms at full velocity with a 500 ms pause between notes.
func runNote(c *jv1080.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: note <name...>")
	}
	for _, name := range args {
		n, err := parseNote(name)
		if err != nil {
			return err
		}
		if err := c.NoteOn(0, n, 127); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
		if err := c.NoteOff(0, n); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// parseNote parses a note name like c1, d3, f#2, bb3, ab3 or eb3 into a
// MIDI note number (c-1 = 0, c4 = 60).
func parseNote(s string) (uint8, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid note %q", s)
	}
	var semitone int
	switch s[0] {
	case 'c':
		semitone = 0
	case 'd':
		semitone = 2
	case 'e':
		semitone = 4
	case 'f':
		semitone = 5
	case 'g':
		semitone = 7
	case 'a':
		semitone = 9
	case 'b':
		semitone = 11
	default:
		return 0, fmt.Errorf("invalid note %q", s)
	}
	i := 1
	for i < len(s) && (s[i] == '#' || s[i] == 'b') {
		if s[i] == '#' {
			semitone++
		} else {
			semitone--
		}
		i++
	}
	octave, err := strconv.Atoi(s[i:])
	if err != nil {
		return 0, fmt.Errorf("invalid note %q", s)
	}
	n := (octave+1)*12 + semitone
	if n < 0 || n > 127 {
		return 0, fmt.Errorf("note %q out of range", s)
	}
	return uint8(n), nil
}

// parseAddress parses two 16-bit hex numbers into a 4-byte address.
// hi holds the first two address bytes, lo the last two, e.g.
// 0x0300 0x002E is the address 03 00 00 2E. Address bytes are 7-bit;
// the high bit of each input byte is dropped.
func parseAddress(hi, lo string) (message.Address, error) {
	h, err := parseNumber(hi)
	if err != nil || h > 0xFFFF {
		return message.Address{}, fmt.Errorf("invalid address hi %q", hi)
	}
	l, err := parseNumber(lo)
	if err != nil || l > 0xFFFF {
		return message.Address{}, fmt.Errorf("invalid address lo %q", lo)
	}
	return message.Address{byte(h>>8) & 0x7F, byte(h) & 0x7F, byte(l>>8) & 0x7F, byte(l) & 0x7F}, nil
}

// parseValues parses decimal or hex values into raw bytes. For size-2
// (nibblized) parameters the values are the two nibble bytes. Values
// outside the documented range are accepted deliberately, to allow special
// tests.
func parseValues(args []string, size int) ([]byte, error) {
	if len(args) != size {
		return nil, fmt.Errorf("want %d value(s), got %d", size, len(args))
	}
	out := make([]byte, size)
	for i, a := range args {
		v, err := parseNumber(a)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q", a)
		}
		out[i] = byte(v)
	}
	return out, nil
}

// parseNumber parses a decimal number or a hex number with 0x prefix.
func parseNumber(s string) (uint64, error) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return strconv.ParseUint(s[2:], 16, 64)
	}
	return strconv.ParseUint(s, 10, 64)
}

// parseCommonValues parses the values of a common parameter, applying
// speaking notations where they exist (pan l13/r13/center).
func parseCommonValues(name string, args []string, size int) ([]byte, error) {
	if p, ok := commonParams[normalize(name)]; ok && p.offset == 0x2F && len(args) == 1 {
		if v, ok := panValue(strings.ToLower(args[0])); ok {
			return []byte{byte(v)}, nil
		}
	}
	return parseValues(args, size)
}

// parseToneValues parses the values of a tone parameter, applying speaking
// notations where they exist (pan l13/r13/center, signed coarse/fine).
func parseToneValues(name string, args []string, size int) ([]byte, error) {
	switch normalize(name) {
	case "pan":
		if len(args) == 1 {
			if v, ok := panValue(strings.ToLower(args[0])); ok {
				return []byte{byte(v)}, nil
			}
		}
	case "coarse":
		if len(args) == 1 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				return []byte{byte(v + 48)}, nil
			}
		}
	case "fine":
		if len(args) == 1 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				return []byte{byte(v + 50)}, nil
			}
		}
	}
	return parseValues(args, size)
}
