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
//	set patch common|tone|wave|tvf <param> <value...>
//	set address <hi> <lo> <value...>
//	get patch common|tone <param>
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

// runPatch dispatches the patch related subcommands: select a patch
// (set patch <bank> <number>) or manipulate/read patch parameters
// (set/get patch common|tone|wave|tvf ...).
func runPatch(c *jv1080.Client, verb string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%s patch what? (select, common, tone, wave, tvf) — try: help", verb)
	}
	switch args[0] {
	case "select":
		if verb != "set" {
			return fmt.Errorf("usage: set patch select <bank> <number>")
		}
		return setPatch(c, args[1:])
	case "common":
		if verb == "set" {
			return setCommon(c, args[1:])
		}
		return getCommon(c, args[1:])
	case "tone":
		if verb == "set" {
			return setTone(c, args[1:])
		}
		return getTone(c, args[1:])
	case "wave":
		if verb == "set" {
			return setWave(c, args[1:])
		}
		return fmt.Errorf("get patch wave is not implemented (wave number is nibblized; use: get patch tone 1 wave)")
	case "tvf":
		if verb == "set" {
			return setTVF(c, args[1:])
		}
		return fmt.Errorf("get patch tvf is not implemented (use: get patch tone 1 cutoff)")
	}
	return fmt.Errorf("unknown patch target %q (want select, common, tone, wave, tvf)", args[0])
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
  set patch common <param> <value...>
        set a patch common parameter of the temporary patch
        (name, efx, level, pan, ...)
  set patch tone <n> <param> <value...>
        set a patch tone parameter (tone 1-4): switch, group, groupid,
        wave, coarse, fine, cutoff, resonance, level, pan, chorus, reverb
  set patch wave <param> <value...>
        set a wave parameter of tone 1: group, groupid, number, gain
  set patch tvf <param> <value...>
        set a TVF parameter of tone 1: cutoff, resonance, envdepth, or
        "set patch tvf env <t1> <t2> <t3> <t4> <l1> <l2> <l3> <l4>"
  set address <hi> <lo> <value...>
        write raw bytes to the address given as two 16-bit hex numbers
        (e.g. 0x0000 0x0006 = master tune)
  get address <hi> <lo> [size]
        read raw bytes from an address (default size 1)
  get patch common <param>
        read a patch common parameter of the temporary patch
  get patch tone <n> <param>
        read a patch tone parameter of the temporary patch
  note <name...>
        play notes by name (c1 d3 f#2 bb3 ab3 eb3 ...), 500 ms at full
        velocity with a 500 ms pause between notes
  help
        show this help

All messages to and from the device are logged in hex:
  SEND: [F0 41 10 6A ... F7]
  REC : [F0 41 10 6A ... F7]

Parameter values may be given in decimal or hex (0x..). Values outside
the documented range are accepted deliberately, to allow special tests.
Pan accepts l1..l63, r1..r63 and center; coarse/fine accept signed
values; the wave number is nibblized and takes two nibble bytes.
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

// setCommon writes a patch common parameter of the temporary patch.
func setCommon(c *jv1080.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: set common <param> <value...>")
	}
	off, size, err := commonParam(args[0])
	if err != nil {
		return err
	}
	values, err := parseCommonValues(args[0], args[1:], size)
	if err != nil {
		return err
	}
	return c.SetPatchParameter(jv1080.BankPresetA, 0, "common", 0, off, values)
}

// setTone writes a patch tone parameter of the temporary patch.
func setTone(c *jv1080.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: set tone <n> <param> <value...>")
	}
	tone, err := strconv.Atoi(args[1])
	if err != nil || tone < 1 || tone > 4 {
		return fmt.Errorf("invalid tone %q (want 1-4)", args[1])
	}
	off, size, err := toneParam(args[2])
	if err != nil {
		return err
	}
	values, err := parseToneValues(args[2], args[3:], size)
	if err != nil {
		return err
	}
	return c.SetPatchParameter(jv1080.BankPresetA, 0, "tone", tone, off, values)
}

// setWave writes wave parameters of tone 1 of the temporary patch:
// set wave group <type> <id>, set wave number <n>, set wave gain <db>.
func setWave(c *jv1080.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: set wave <param> <value...>")
	}
	off, size, err := waveParam(args[0])
	if err != nil {
		return err
	}
	values, err := parseValues(args[1:], size)
	if err != nil {
		return err
	}
	return c.SetPatchParameter(jv1080.BankPresetA, 0, "tone", 1, off, values)
}

// setTVF writes TVF parameters of tone 1 of the temporary patch:
// set tvf cutoff <v>, set tvf resonance <v>, set tvf env t1 t2 t3 t4 l1 l2 l3 l4.
func setTVF(c *jv1080.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: set tvf <param> <value...>")
	}
	if strings.EqualFold(args[0], "env") {
		if len(args) != 9 {
			return fmt.Errorf("usage: set tvf env <t1> <t2> <t3> <t4> <l1> <l2> <l3> <l4>")
		}
		values, err := parseValues(args[1:], 1)
		if err != nil {
			return err
		}
		return c.SetTVFEnvelope(jv1080.BankPresetA, 0, 1, values)
	}
	off, size, err := tvfParam(args[0])
	if err != nil {
		return err
	}
	values, err := parseValues(args[1:], size)
	if err != nil {
		return err
	}
	return c.SetPatchParameter(jv1080.BankPresetA, 0, "tone", 1, off, values)
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

// getCommon reads a patch common parameter of the temporary patch.
func getCommon(c *jv1080.Client, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: get common <param>")
	}
	off, size, err := commonParam(args[0])
	if err != nil {
		return err
	}
	data, err := c.ReadRaw(jv1080.TemporaryPatchArea().Add(off), size)
	if err != nil {
		return err
	}
	fmt.Printf("common %s: %s\n", args[0], formatValue(args[0], data))
	return nil
}

// getTone reads a patch tone parameter of the temporary patch.
func getTone(c *jv1080.Client, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: get tone <n> <param>")
	}
	tone, err := strconv.Atoi(args[0])
	if err != nil || tone < 1 || tone > 4 {
		return fmt.Errorf("invalid tone %q (want 1-4)", args[0])
	}
	off, size, err := toneParam(args[1])
	if err != nil {
		return err
	}
	area := jv1080.PatchToneArea(jv1080.TemporaryPatchArea(), tone)
	data, err := c.ReadRaw(area.Add(off), size)
	if err != nil {
		return err
	}
	fmt.Printf("tone %d %s: %s\n", tone, args[1], formatValue(args[1], data))
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
