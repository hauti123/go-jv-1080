# GUI cleanup — review of planned changes and code smells

Status: **implemented** (planned refactor of section 1 applied, plus the
fixes decided in review): widget debug logging is gated by `LOG_GUI_WIDGETS`,
the LED colors have a single definition in the theme, the `Offset`/`offset`
wrapper pair is collapsed, `Envelope.Layout` is split into geometry / event /
draw helpers (the mouse-to-point drag translation is unchanged),
raw MIDI logging is gated by `LOG_MIDI_RX` / `LOG_MIDI_TX` (default on),
all clamp helpers (gui, `internal/jv1080`, `cmd/gui-proto`) are migrated to
`internal/util.Clamp`, `cmd/gui-proto` is kept, and the knob double-click
reset now goes to zero (the center for pan/bipolar knobs). The app window, event loop and
PoC screen moved from `cmd/jv1080-gui/main.go` into `internal/gui` (`app.go`, `poc.go` behind
a `screen` interface, so the PoC can be replaced by the patch edit screen later); the command
itself is thin now. The smell
register below stays recorded; the items not covered above are still
deferred to a later pass.

Scope reviewed: `internal/gui` (all files), its only consumer
`cmd/jv1080-gui`, the duplicated clamp helpers in `internal/jv1080`, and
`cmd/gui-proto` (1222-line prototype). `go test ./internal/gui/...` passes
on `main`; full `go build ./...` is not possible in every environment
(ALSA / pkg-config needed by the MIDI driver and Gio windowing), so
`cmd/*` builds must be verified where those system deps exist.

## 1. Planned refactor

### 1.1 New package `internal/util` for shared helpers

- Add `internal/util/clamp.go` with a generic clamp, e.g.
  `func Clamp[T cmp.Ordered](v, lo, hi T) T` (stdlib `cmp`, Go 1.24).
- Replace the duplicated clamp helpers with it:
  - `internal/gui`: `clampInt` (`internal/gui/knob.go:351`),
    `clampf01` (`internal/gui/envelope.go:140`) → `util.Clamp`;
    delete dead `clampf` (`internal/gui/envelope.go:150`, no callers).
  - `internal/jv1080`: `clamp` / `clampByte` (`internal/jv1080/patch.go:279,289`,
    also used in `client.go:70,80` and `params.go:351,412`) → `util.Clamp`
    (keep a one-line local `clampByte` wrapper there if the ~20 call sites
    should stay unchanged).
  - `cmd/gui-proto/main.go:71` keeps its own copy (prototype, see open
    questions) or is switched to `util.Clamp` as well.
- Replace stdlib re-implementations instead of moving them:
  `itoa` (`internal/gui/knob.go:361`) → `strconv.Itoa`,
  `absF` / `roundF` (`internal/gui/envelope.go:123,130`) → `math.Abs` /
  `math.Round` (same half-away-from-zero semantics).
- Add `internal/util/clamp_test.go`.

### 1.2 Move all widget code to sub-package `internal/gui/widget`

- Move `knob.go`, `envelope.go`, `draw.go`, `font5x7.go`, `theme.go` and the
  three test files from `internal/gui/` to `internal/gui/widget/`,
  renaming the package to `widget`. Tests are white-box (they use
  `angleOf`, `litRange`, `set`, `points`, `pixFont`, …), so they move along
  unchanged apart from the package clause.
- `internal/gui` keeps only a `doc.go` with the package comment (app-level
  GUI code — tabs, bank ribbon — will land there later); widgets live in
  `widget`.
- Public API renames at the call site in `cmd/jv1080-gui/main.go`:
  `gui.Knob` → `widget.Knob`, `gui.NewKnob` → `widget.NewKnob`,
  `gui.Envelope` → `widget.Envelope`, `gui.Theme` → `widget.Theme`,
  `gui.DefaultTheme` → `widget.DefaultTheme`, `gui.FontCellDp`,
  `gui.PixelTextSize`, `gui.DrawPixelText`, `gui.Offset`,
  `gui.KnobMinMax/KnobBipolar/KnobPan` → `widget.*`.
- Update docs: `doc/architecture.md` package table and layer diagram,
  `README.md` package list.

### 1.3 Verification after the move

- `go test ./internal/gui/... ./internal/util/... ./internal/jv1080/...`
  (jv1080 tests that need the MIDI driver may be environment-limited),
  `go vet` on the buildable packages, `gofmt -l .`, and a build of
  `cmd/jv1080-gui` in an environment with ALSA/pkg-config.

## 2. Code smell register (recorded, not fixed in this pass)

| # | Location | Smell | Proposed action (later pass) |
|---|----------|-------|------------------------------|
| S1 | `internal/gui/envelope.go:118,177,251,274,282,327,353,372,384` | 9 `log.Printf` debug calls inside widget code, incl. per-frame logging in `Layout`/`relayout` — a widget library should not log | Remove, or gate behind a debug flag |
| S2 | `internal/gui/draw.go:28` `fillCircle`, `:125` `rotated`, `:102` `drawTicks` (wrapper), `internal/gui/font5x7.go:96` `ledFrame` | Dead code — defined, never called | Delete |
| S3 | `internal/gui/envelope.go:57-58` `envMinX`/`envMaxX` | Dead constants (leftover from gui-proto `clampPoint`) | Delete |
| S4 | `internal/gui/envelope.go:42,72` `defaults`, and type comment "double-click a point resets it to its default" | Field written but never read; documented double-click reset not implemented | Implement reset or remove field + doc line |
| S5 | `internal/gui/envelope.go:138` `Active()`, `internal/gui/knob.go:159` `center()` | Exported/unexported helpers used only by tests | Keep only if needed; otherwise delete and adjust tests |
| S6 | `internal/gui/font5x7.go:58-59` vs `internal/gui/theme.go:23-24` | LED colors defined twice (`ledLit`/`ledDim` vs `ledLitHex`/`ledDimHex`, same values) | Single definition in theme, font uses it |
| S7 | `internal/gui/draw.go:133,136` | Exported `Offset` is a one-line wrapper around unexported `offset` | Collapse to one exported `Offset` |
| S8 | `internal/gui/knob.go:361`, `envelope.go:123,130` | Hand-rolled `itoa`/`absF`/`roundF` duplicate stdlib | `strconv.Itoa`, `math.Abs`, `math.Round` (see 1.1) |
| S9 | `internal/gui` `clampInt`/`clampf`/`clampf01`, `internal/jv1080/patch.go:279,289` `clamp`/`clampByte`, `cmd/gui-proto/main.go:71` `clampf` | Same clamp logic duplicated in 3 packages | `internal/util.Clamp` (see 1.1) |
| S10 | `internal/gui/envelope.go:157+` `Layout` | ~200-line method mixing geometry, drawing and event handling; `Knob` already splits `processEvents`/`draw` | Split `Envelope` the same way |
| S11 | `internal/gui/knob.go` `draw` (radii 0.329/0.371/0.421/0.271/0.207 …), `envelope.go` pads | Magic numbers for geometry | Named constants (low priority) |
| S12 | `internal/gui/knob.go:74` `NewKnob(style, title, key, min, max)` | 5 positional parameters, easy to mix up (e.g. min/max order is silently swapped) | Options struct or named params (API change — decide with reviewer) |
| S13 | `internal/gui/theme.go:1-8` | Package doc comment lives in `theme.go` | Move to `doc.go` in `widget` |
| S14 | `internal/gui/theme.go:13-17` `MinAngle`/`MaxAngle` | Scale geometry only used by the knob | Move to `knob.go` (low priority) |
| S15 | `cmd/jv1080-gui/main.go:281` | `goto drained` to break out of a drain loop | Small helper func instead |
| S16 | `cmd/jv1080-gui/main.go:231` `panUIToKnob` | String pan parsing with `fmt.Sscanf` ignoring the error | Check error / use `strconv.Atoi` |
| S17 | `cmd/jv1080-gui/main.go:96-99` | Every raw MIDI message logged unconditionally via `dev.SetLogger` | Gate behind a `-v`/`-debug` flag |
| S18 | `cmd/gui-proto/main.go` (1222 lines) | Still contains full copies of knob/envelope/draw code and `clampf`, now superseded by `internal/gui/widget` | Delete once the real GUI exists, or keep as style reference (open question) |
| S19 | `doc/architecture.md:78,102` | Two `## Data flow` sections; the second is stale (mentions `JVController.HandleDeviceMessage` / `Controller.Get`, which no longer exist) | Merge into one section matching the current `Client` API |
| S20 | `README.md` | File is UTF-16 LE with CRLF (renders with BOM garbage in some viewers) | Convert to UTF-8 |
| S21 | `internal/data/` | Empty package dir, only `.gitkeep` | Remove or populate |
| S22 | `internal/gui/knob.go` double-click handler (`processEvents`, Press case) | Double-click resets to `Min` even for bipolar/pan knobs (pan resets to L64, not C) | Decide: reset to center for bipolar/pan? |
| S23 | `internal/gui/envelope_drag_test.go` `TestHitTestSelectsPressedHandle` | Test re-implements the hit-test instead of calling it | Extract hit-test into a method and test that |

## 3. Open questions for the reviewer

1. Should `internal/jv1080`'s `clamp`/`clampByte` also migrate to
   `internal/util` in the same pass (they are the same helper duplicated),
   or only the GUI ones?
2. Keep `cmd/gui-proto` as the visual reference (untouched), switch it to
   `internal/gui/widget`, or delete it?
3. `NewKnob`'s positional parameters (S12): change the API now or keep it?
4. Double-click reset semantics for bipolar/pan knobs (S22): `Min` (current)
   or center `0`?
5. README.md conversion from UTF-16 to UTF-8 (S20): include in the cleanup
   PR or separately?
