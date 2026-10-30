# Architecture: GUI, controller and device layers

The JV-1080 editor is split into three layers with strict one-way
dependencies. The GUI widgets never import the MIDI or device packages;
they only know the small `gui.Controller` interface. This keeps the widget
code testable and reusable.

```
┌──────────────────────────────────────────────────────────────────────┐
│  cmd/jv1080-gui (sample app)                                          │
│  • creates the Gio window and the three knobs (level / pan / octave) │
│  • wires knob OnChange callbacks to Controller.Set                   │
│  • forwards Controller.Update notifications to the knobs             │
└───────────────┬──────────────────────────────┬───────────────────────┘
                │ uses                            │ implements
                ▼                                  ▼
┌───────────────────────────────┐   ┌────────────────────────────────────┐
│  internal/gui (widgets)       │   │  internal/gui.JVController         │
│  • Knob gadget (3 styles)     │   │  • gui.Controller implementation   │
│  • 5x7 LED font renderer      │   │  • maps gui.Param → SysEx address  │
│  • dark green MK-II theme     │   │  • decodes unsolicited DT1 into    │
│  • gui.Controller interface   │   │    typed Param updates             │
│  • gui.Param descriptor       │   │  • attached to midi.Device         │
└───────────────┬───────────────┘   │    callbacks                       │
                │ implements          └──────────────┬─────────────────────┘
                ▼                                    │ uses
┌───────────────────────────────┐                  ▼
│  gui.Controller (interface)   │   ┌────────────────────────────────────┐
│  Get / Set / Update           │   │  internal/jv1080 (device client)   │
└───────────────────────────────┘   │  • single-parameter DT1 writes     │
                                    │    (SetPatchParameter / ReadRaw)   │
                                    │  • RQ1/DT1 request-reply reads     │
                                    └──────────────┬─────────────────────┘
                                                   │ uses
                                                   ▼
                                    ┌────────────────────────────────────┐
                                    │  internal/midi                     │
                                    │  • Device: ports, send, receive,   │
                                    │    registered callbacks for        │
                                    │    unsolicited messages            │
                                    │  • message: serializable structs   │
                                    │    (RQ1/DT1/CC/PC/NoteOn/...)      │
                                    └────────────────────────────────────┘
```

## Data flow

**User turns a knob (GUI → device):**

1. `Knob.Layout` detects a drag/scroll, changes its value and calls
   `OnChange(value)`.
2. The sample app calls `Controller.Set(param, []byte{value})`.
3. `JVController` resolves the `gui.Param` (block/tone/offset, CLI name)
   to a SysEx address for the current bank/number via
   `jv1080.PatchParameterAddress`.
4. `jv1080.Client.WriteRaw` sends a single DT1 with exactly one data byte
   — only the changed parameter is transmitted, the device keeps
   processing at a minimum.

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
| `internal/gui` | GUI widgets (`Knob`, 5x7 LED font, theme) plus the `Controller`/`Param` abstraction and the `JVController` implementation that bridges to `internal/jv1080`. |
| `cmd/jv1080-gui` | Sample GUI app (level/pan/octave knobs). |
| `cmd/jv1080-cli` | Command line utility for raw MIDI inspection and patch manipulation. |

## Why the Controller interface

`internal/gui` defines `Controller` (Get/Set/Update) so the widgets can be
developed and tested without any MIDI hardware: a fake controller can drive
the knobs in tests, and the real `JVController` plugs in only in
`cmd/jv1080-gui`. The GUI therefore has a clean, one-way dependency on the
device layer through this narrow interface.
