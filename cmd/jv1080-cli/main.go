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
	// Log all incoming and outgoing raw messages in hex (default on;
	// disable per direction with LOG_MIDI_RX=0 / LOG_MIDI_TX=0).
	logRX := midi.LogEnabled("RX")
	logTX := midi.LogEnabled("TX")
	dev.SetLogger(func(direction string, msg []byte) {
		if direction == "RX" && !logRX || direction == "TX" && !logTX {
			return
		}
		fmt.Printf("%s: [%s]\n", direction, message.Hex(msg))
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

// runPatch dispatches the patch related subcommands: select a patch
// (set patch <bank> <number>) or manipulate/read patch parameters
// (set/get patch common|efx|ctrl|wave|lfo|pitch|tvf|tva ...). Parameters
// are addressed through the device layer's path-like keys
// ("patch/temp/..."); all addressing and value conversion happens there.
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
				p, err := jv1080.ResolveParam(cmd.key)
				if err != nil {
					return err
				}
				value, err := parseParamValue(p, cmd.valArgs)
				if err != nil {
					return err
				}
				if err := c.SetParameter(cmd.key, value); err != nil {
					return err
				}
			}
			return nil
		}
		return getPatchParams(c, cmds)
	}
	return fmt.Errorf("unknown patch target %q (want select, %s)", args[0], strings.Join(patchCategories, ", "))
}

// getPatchParams reads the parameters of parsed patch commands from the
// temporary patch via the device layer's key API and prints them.
func getPatchParams(c *jv1080.Client, cmds []patchCommand) error {
	for _, cmd := range cmds {
		p, err := jv1080.ResolveParam(cmd.key)
		if err != nil {
			return err
		}
		v, err := c.GetParameter(cmd.key)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", p.Key, formatParamValue(p, v))
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

All raw MIDI messages sent and received are logged in hex:
  TX: [F0 41 10 6A ... F7]
  RX: [F0 41 10 6A ... F7]

Parameter values may be given in decimal or hex (0x..). Pan accepts
l1..l64, r1..r63 and center (or a signed value -64..+63); bipolar
parameters (octave, coarse, fine, depths, envelope levels) accept signed
values; nibblized parameters (number, tempo) take a single combined
value. All addressing and value conversion happens in the device layer
(patch/temp/... keys), the CLI only maps its category syntax onto keys.
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
