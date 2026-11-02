# Architecture: GUI and device layers

The JV-1080 editor is split into layers with strict one-way
dependencies. The GUI widgets never import the MIDI or device packages
and never deal with SysEx addresses or offsets: each knob is bound to a
device parameter through a path-like key string (e.g.
"patch/temp/common/level"), which is resolved solely in the device layer
(`internal/jv1080`). This keeps the widget code testable and reusable.

```
┌───────────────────────────────────────────────────────────────┐
│  cmd/jv1080-gui (thin command)                                │
│  • flags, MIDI device, raw MIDI logging (LOG_MIDI_RX/TX)      │
│  • gui.NewApp() → SetDevice(client) → Main()                  │
└──────────────────────────┬────────────────────────────────────┘
                │ uses
                ▼
┌───────────────────────────────────────────────────────────────┐
│  internal/gui (app layer)                                     │
│  • App: Gio window, event loop, status line, current screen   │
│  • PoC screen (poc.go): level/pan/octave knobs + pitch env,  │
│    bound to "patch/temp/..." keys — later the patch edit      │
│    screen; widget.OnChange → client.SetParameter(key, value)  │
│  • client.OnParameterChange → screen.applyParam(key, value)   │
└──────────────────────────┬────────────────────────────────────┘
                │ uses (widget Key strings only, no addresses)
                ▼
│  internal/gui/widget (widgets)│
│  • Knob gadget (3 styles,     │
│    each with a Key field)     │
│  • Envelope editor (4 points, │
│    TimeKeys/LevelKeys)        │
│  • MK-II style drawing +      │
│    5x7 LED font (gui-proto)   │
│  • dark green MK-II theme     │
└───────────────┬───────────────┘
                │ uses
                ▼
┌──────────────────────────────────────────────────────────────┐
│  internal/jv1080 (device client)                              │
│  • key → SysEx address map (params.go, patch/temp/... only)   │
│  • SetParameter(key string, value any) error                  │
│  • GetParameter(key string) (any, error)                      │
│  • OnParameterChange(func(key string, value any))             │
│    decodes unsolicited DT1 into key + UI value                │
│  • UI value ↔ wire conversion (bipolar offset, pan C/L/R)     │
│  • lower-level API unchanged: ReadRaw/WriteRaw, dumps, ...    │
└───────────────┬──────────────────────────────────────────────┘
                │ uses
                ▼
┌──────────────────────────────────────────────────────────────┐
│  internal/midi                                                │
│  • Device: ports, send, receive, registered callbacks for     │
│    unsolicited messages                                       │
│  • message: serializable structs (RQ1/DT1/CC/PC/NoteOn/...)   │
└──────────────────────────────────────────────────────────────┘
```

## Parameter keys

Parameters are addressed with path-like keys; only the patch-mode
temporary patch is supported for now (`patch/user/N/...` comes later):

- `patch/temp/common/<name>` — patch common parameters (*1-3-1)
- `patch/temp/tone/<n>/<name>` — tone parameters directly in the tone
  block, n = 1-4 (wave, ranges, receive switches)
- `patch/temp/tone/<n>/<sub>/<name>` — tone sub-blocks: `efx` (output
  routing), `ctrl<1-3>` (destinations/depths), `lfo<1-2>` (LFO and LFO
  depths), `pitch`, `tvf`, `tva`

All patch parameters of doc/patch-edit-mode-structure.md are registered;
`jv1080.ListParams()` lists them, `jv1080.ResolveParam(key)` resolves one.
Value kinds: raw 0-127, bipolar (signed, stored with an offset), pan
(UI signed -64..+63 or `"C"`/`"L1".."L64"`/`"R1".."R63"`; wire 0-127
with 64 = C), nibblized (two nibble bytes, e.g. `number`, `tempo`) and
the 12-byte `name`.

The mapping from key to SysEx address and the UI-value ↔ wire-byte
conversion live solely in `internal/jv1080` (`params.go`; UI and wire both start at 0 = center for pan, so -1 → wire 63 = L1, +1 → wire 65 = R1),
`Client.SetParameter`/`GetParameter`/`OnParameterChange`). The GUI glue
layer (`internal/gui`), the CLI (`cmd/jv1080-cli`) and the widgets never
deal with addresses, offsets or raw bytes; the CLI only maps its category
syntax (`set patch tvf tone 2 cutoff 100`) onto keys.

## Data flow

**User turns a knob (GUI → device):**

1. `Knob.Layout` detects a drag/scroll, changes its value and calls
   `OnChange(uiValue)`.
2. The app layer (`internal/gui`) calls `client.SetParameter(widget.Key, uiValue)`.
3. `internal/jv1080` resolves the key to a SysEx address and converts
   the UI value to wire bytes (e.g. +2 octave → 5, "L13" → 51).
4. A single DT1 with exactly one data byte is sent — only the changed
   parameter is transmitted.

**Device value changes (device → GUI):**

1. The JV-1080 sends an unsolicited DT1 (front panel edit, another
   editor, a dump echo).
2. `midi.Device` dispatches the raw message to the client's handler.
3. `internal/jv1080` maps the DT1 address back to a key, decodes the
   wire bytes into a UI value and invokes the `OnParameterChange`
   callbacks with `(key, value)`.
4. The app layer calls `widget.Set(value)` on the widget with that key —
   only the one changed knob is updated.


## Data flow

**User turns a knob (GUI → device):**

1. `Knob.Layout` detects a drag/scroll, changes its value and calls
   `OnChange(value)`.
2. The sample app calls `jv1080.Client.SetParameter(key, uiValue)`.
3. The device layer resolves the key to its SysEx address and encodes
   the UI value into wire bytes (`Param.EncodeValue`).
4. A single DT1 with exactly the parameter's data bytes is sent — only
   the changed parameter is transmitted, the device keeps processing at
   a minimum.

**Device value changes (device → GUI):**

1. The JV-1080 sends an unsolicited DT1 (front panel edit, another
   editor, a dump echo).
2. `midi.Device` dispatches the raw message to all registered callbacks,
   including `JVController.HandleDeviceMessage`.
3. The controller checks whether the DT1 address falls into a registered
   parameter of the current patch and, if so, invokes the `Update`
   callbacks with the typed `Param` and raw value.
4. The sample app re-reads only that one parameter (`Controller.Get`)
   and calls `Knob.Set` on the matching knob — the UI reflects the change
   immediately, and no other parameter is touched.

## Packages

| Package | Role |
| --- | --- |
| `internal/midi` | MIDI transport: port enumeration, `Device` with send/receive and callback registration, raw message splitting. |
| `internal/midi/message` | Individual MIDI messages as structs, serializable to/from byte streams (RQ1, DT1, CC, PC, NoteOn, identity). |
| `internal/midi/sysex`, `internal/midi/controlchange` | SysEx-level and control-change helpers on top of `message`. |
| `internal/jv1080` | JV-1080 parameter address map, `Patch`/`Performance` models and the high-level `Client` (dumps, single reads, single-parameter writes, store). |
| `internal/gui` | App layer: window, event loop, status line, PoC screen (`app.go`, `poc.go`); wires widgets to `jv1080.Client`. |
| `internal/gui/widget` | GUI widgets: `Knob` (with `Key` binding), `Envelope` editor, 5x7 LED font, MK-II theme. No device imports. |
| `internal/util` | Shared helpers: generic `Clamp[T]`, env-flag parsing (`EnvFlag`). |
| `cmd/jv1080-gui` | Thin command: flags, MIDI device, raw MIDI logging; starts the app (`internal/gui`). |
| `cmd/jv1080-cli` | Command line utility for raw MIDI inspection and patch manipulation. |
