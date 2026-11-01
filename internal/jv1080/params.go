package jv1080

import (
	"fmt"
	"strconv"
	"strings"

	"go-jv-1080/internal/midi/message"
)

// Parameter keys are path-like strings addressed to the device layer, e.g.
// "patch/temp/common/level". The first segment selects the memory area
// ("patch"), the second the patch ("temp" = patch-mode temporary patch,
// later "user/1".."user/128" etc.), then the block ("common" or "tone/N"),
// then the parameter name.
//
// The mapping from key to SysEx address and the conversion between UI
// values (int, string, bool) and wire bytes live solely in this package;
// callers never deal with addresses or offsets.

// ParamKind describes how a parameter's UI value maps onto wire bytes.
type ParamKind int

const (
	// KindRaw is a plain 0-127 byte parameter (e.g. level).
	KindRaw ParamKind = iota
	// KindBipolar is a signed value stored with an offset (e.g. octave
	// -3..+3 stored as 0-6).
	KindBipolar
	// KindPan is a pan value: UI value is "C", "L1".."L63", "R1".."R63"
	// (or 0-127 raw); wire value is 0-127 with 64 = center.
	KindPan
)

// Param describes one addressable parameter of the device.
type Param struct {
	Key    string    // path-like key, e.g. "patch/temp/common/level"
	Kind   ParamKind // how UI values convert to wire bytes
	Min    int       // UI value minimum (for KindBipolar/KindPan: signed / raw)
	Max    int       // UI value maximum
	Offset int       // byte offset within the block (device-internal)
	Block  string    // "common" or "tone/N" (device-internal)
	area   message.Address
}

// paramSpec is the static definition of a parameter: its key template and
// wire location within the temporary patch.
type paramSpec struct {
	name   string
	block  string // "common" or "tone/N"
	offset int
	kind   ParamKind
	min    int
	max    int
}

// tempPatchParams are the parameters of the patch-mode temporary patch
// ("patch/temp/...") that are exposed through the key system. The offsets
// and ranges follow doc/patch-edit-mode-structure.md (same names as in
// cmd/jv1080-cli).
var tempPatchParams = []paramSpec{
	{name: "name", block: "common", offset: 0x00, kind: KindRaw, min: 0, max: 0}, // 12 ASCII bytes, handled specially
	{name: "efxtype", block: "common", offset: 0x0C, kind: KindRaw, min: 0, max: 39},
	{name: "level", block: "common", offset: 0x2E, kind: KindRaw, min: 0, max: 127},
	{name: "pan", block: "common", offset: 0x2F, kind: KindPan, min: -64, max: 63}, // UI units -64..+63, 0 = C (JV-1080 scale L64..R63)
	{name: "octave", block: "common", offset: 0x41, kind: KindBipolar, min: -3, max: 3},
	{name: "stretch", block: "common", offset: 0x42, kind: KindRaw, min: 0, max: 3},
	{name: "priority", block: "common", offset: 0x43, kind: KindRaw, min: 0, max: 1},
	{name: "switch", block: "tone/1", offset: 0x00, kind: KindRaw, min: 0, max: 1},
	{name: "level", block: "tone/1", offset: 0x65, kind: KindRaw, min: 0, max: 127},
	{name: "pan", block: "tone/1", offset: 0x77, kind: KindPan, min: -64, max: 63},
	{name: "cutoff", block: "tone/1", offset: 0x51, kind: KindRaw, min: 0, max: 127},
	{name: "resonance", block: "tone/1", offset: 0x53, kind: KindRaw, min: 0, max: 127},
	{name: "coarse", block: "tone/1", offset: 0x3D, kind: KindBipolar, min: -48, max: 48},
	{name: "fine", block: "tone/1", offset: 0x3E, kind: KindBipolar, min: -50, max: 50},
	// tone 4 spot-check: the tone blocks repeat every 2 address bytes
	{name: "cutoff", block: "tone/4", offset: 0x51, kind: KindRaw, min: 0, max: 127},
}

// resolveParam maps a path-like key to its parameter definition and the
// SysEx address of its first byte. Only "patch/temp/..." keys are
// supported for now.
func resolveParam(key string) (Param, error) {
	segs := strings.Split(strings.Trim(key, "/"), "/")
	if len(segs) < 4 || segs[0] != "patch" || segs[1] != "temp" {
		return Param{}, fmt.Errorf("jv1080: unsupported parameter key %q (want patch/temp/<block>/<name> or patch/temp/tone/<n>/<name>)", key)
	}
	var block, name string
	if segs[2] == "tone" {
		if len(segs) != 5 {
			return Param{}, fmt.Errorf("jv1080: unsupported parameter key %q (want patch/temp/tone/<n>/<name>)", key)
		}
		block, name = "tone/"+segs[3], segs[4]
	} else {
		block, name = segs[2], segs[3]
	}
	for _, spec := range tempPatchParams {
		if spec.block == block && spec.name == name {
			area := TemporaryPatchArea()
			var addr message.Address
			switch {
			case block == "common":
				addr = area.Add(spec.offset)
			default:
				// block "tone/N"
				t, err := strconv.Atoi(strings.TrimPrefix(block, "tone/"))
				if err != nil || t < 1 || t > 4 {
					return Param{}, fmt.Errorf("jv1080: invalid tone block %q", block)
				}
				addr = PatchToneArea(area, t).Add(spec.offset)
			}
			key := "patch/temp/" + block + "/" + name
			if strings.HasPrefix(block, "tone/") {
				key = "patch/temp/tone/" + strings.TrimPrefix(block, "tone/") + "/" + name
			}
			return Param{
				Key:    key,
				Kind:   spec.kind,
				Min:    spec.min,
				Max:    spec.max,
				Offset: spec.offset,
				Block:  block,
				area:   addr,
			}, nil
		}
	}
	return Param{}, fmt.Errorf("jv1080: unknown parameter key %q", key)
}

// EncodeValue converts a UI value (int, string, bool) into the raw wire
// bytes for the parameter. For KindBipolar the signed value is shifted by
// -Min; for KindPan strings like "L13"/"R13"/"C" (or plain 0-127 ints) are
// converted to the 0-127 wire value with 64 = center.
func (p Param) EncodeValue(v any) ([]byte, error) {
	switch p.Kind {
	case KindBipolar:
		n, err := toInt(v)
		if err != nil {
			return nil, fmt.Errorf("jv1080: %s: %w", p.Key, err)
		}
		if n < p.Min || n > p.Max {
			return nil, fmt.Errorf("jv1080: %s: value %d out of range %d..%d", p.Key, n, p.Min, p.Max)
		}
		return []byte{byte(n - p.Min)}, nil
	case KindPan:
		raw, err := panToRaw(v)
		if err != nil {
			return nil, fmt.Errorf("jv1080: %s: %w", p.Key, err)
		}
		return []byte{byte(raw)}, nil
		// UI and wire both start at 0 = center: panToRaw maps signed
		// UI units -63..+63 (and "C"/"Lxx"/"Rxx" strings) onto 0..127.
	default:
		if p.Key == "patch/temp/common/name" {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("jv1080: %s: want string value", p.Key)
			}
			b := make([]byte, 12)
			for i := range b {
				if i < len(s) {
					b[i] = clampByte(int(s[i]), 32, 127)
				} else {
					b[i] = ' '
				}
			}
			return b, nil
		}
		n, err := toInt(v)
		if err != nil {
			return nil, fmt.Errorf("jv1080: %s: %w", p.Key, err)
		}
		if n < p.Min || n > p.Max {
			return nil, fmt.Errorf("jv1080: %s: value %d out of range %d..%d", p.Key, n, p.Min, p.Max)
		}
		return []byte{byte(n)}, nil
	}
}

// DecodeValue converts raw wire bytes into the UI value for the parameter
// (int for KindRaw/KindBipolar, string for KindPan and the name).
func (p Param) DecodeValue(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("jv1080: %s: empty data", p.Key)
	}
	switch p.Kind {
	case KindBipolar:
		return int(data[0]) + p.Min, nil
	case KindPan:
		return rawToPan(int(data[0])), nil // returns signed UI units -63..+63, 0 = C
	default:
		if p.Key == "patch/temp/common/name" {
			return decodeName(data), nil
		}
		return int(data[0]), nil
	}
}

// panToRaw converts a pan UI value to the raw 0-127 wire byte. UI and
// wire both start at 0 = center; the JV-1080 scale is L64..L1, C,
// R1..R63, so a plain int is the signed UI value -64..+63 (0 = C,
// -64 = L64) and maps to 64+ui. Strings "C"/"center" (0), "L1".."L64"
// (negative) and "R1".."R63" (positive) are also accepted.
func panToRaw(v any) (int, error) {
	ui, err := panToUI(v)
	if err != nil {
		return 0, err
	}
	if ui < -64 || ui > 63 {
		return 0, fmt.Errorf("pan value %d out of range -64..+63 (L64..R63)", ui)
	}
	return 64 + ui, nil
}

// rawToPan converts a raw 0-127 wire byte to the signed UI pan value
// -64..+63 (0 = C, -64 = L64).
func rawToPan(raw int) int {
	return clamp(raw, 0, 127) - 64
}

// panToUI converts a pan UI value (signed int -64..+63, or "C"/"Lxx"/
// "Rxx" string) to the signed UI value.
func panToUI(v any) (int, error) {
	switch t := v.(type) {
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		switch s {
		case "c", "center":
			return 0, nil
		}
		if len(s) >= 2 && (s[0] == 'l' || s[0] == 'r') {
			n, err := strconv.Atoi(s[1:])
			max := 63
			if s[0] == 'l' {
				max = 64 // left side reaches L64
			}
			if err != nil || n < 1 || n > max {
				return 0, fmt.Errorf("invalid pan %q (want C, L1-L64 or R1-R63)", t)
			}
			if s[0] == 'l' {
				return -n, nil
			}
			return n, nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid pan %q", t)
		}
		return n, nil
	default:
		n, err := toInt(v)
		if err != nil {
			return 0, err
		}
		return n, nil
	}
}

func toInt(v any) (int, error) {
	switch t := v.(type) {
	case int:
		return t, nil
	case int32:
		return int(t), nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		return strconv.Atoi(strings.TrimSpace(t))
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("cannot convert %T to int", v)
}
