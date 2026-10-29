# JV-1080 Patch Edit Mode — Menu Structure and SysEx Parameter Map

This document describes the device menu structure of the Patch Edit Mode
(Owner's Manual pp. 42–56) and maps every parameter to its SysEx representation
plus a suggested name for the `jv1080-cli` tool. Full parameter descriptions
are in [`patch-edit-mode.md`](patch-edit-mode.md); offsets and value ranges
come from the tables *1-3-1 Patch common* and *1-3-2 Patch tone* in
[`jv-1080-midi.md`](jv-1080-midi.md).

## Blocks and addressing

A Patch consists of one **Patch common** block (offset `00 00`, total size
`00 00 00 48`) and four identical **Patch tone** blocks (Tone 1–4 at offsets
`10 00`, `12 00`, `14 00`, `16 00`, total size `00 00 01 01` each). Offsets below
are given relative to the start of the enclosing block (as in the MIDI
implementation); a value marked `#` ("divided-by-two data") is transmitted as
two nibble bytes, high nibble first.

Patch-wide parameters live in the common block; per-Tone parameters live in the
tone blocks and exist once for each of the four Tones.

## CLI naming convention

The suggested names follow the existing `jv1080-cli` parameter style
(`cmd/jv1080-cli/params.go`): short, lowercase, category-free words. They name
only the parameter itself — the command already carries the edit mode and
category, e.g.:

```
jv1080-cli set patch common level <val>
jv1080-cli set patch common pan <val>
jv1080-cli set patch tone 3 cutoff <val>
```

So a name must be unique only within its block (common or tone), not across
categories. Where the same concept exists in both blocks (level, pan), the same
name is used. Numbered series use a suffix (`envt1`…`envt4`, `ctrl1dst1`…).
Names already implemented in `params.go` are kept (`level`, `pan`, `cutoff`,
`resonance`, `group`, `groupid`, `number`/`wavenumber`, `gain`, `coarse`,
`fine`, `switch`, `chorus`, `reverb`).

## Menu structure

```
PATCH EDIT
├── COMMON
│   ├── PATCH NAME ............. name 1–12
│   ├── PATCH COMMON ........... level, pan, analog feel, octave, stretch, priority
│   ├── VELOCITY ................ range switch (common), lower, upper, x-fade (per tone)
│   ├── KEY RANG ................ lower, upper (per tone)
│   └── STRUCT .................. structure, booster (common)
├── EFFECTS
│   ├── OUTPUT .................. output assign/level, chorus/reverb send (per tone)
│   ├── PATCH EFX TYPE .......... EFX type
│   ├── PATCH EFX PRM ........... EFX parameters 1–12
│   ├── PATCH EFX OUT ........... EFX output assign/level, chorus/reverb send
│   ├── PATCH EFX CTRL .......... EFX control sources 1, 2 and depths 1, 2
│   ├── PATCH CHORUS ............ rate, depth, pre delay, feedback, level, output
│   └── PATCH REVERB ............ type, time, level, feedback, HF damp
├── CONTROL
│   ├── KEY MODE & BENDER ....... key assign, legato, bend range up/down
│   ├── PORTAMENTO .............. switch, time, mode, type, start
│   ├── RxSWITCH ................ volume, pan, bender, hold-1, re-damper
│   ├── PEAK & HOLD ............. EFX/Control 1–3 hold/peak
│   ├── CONTROL SOURCE .......... patch control sources 2, 3
│   └── CONTROL 1 to 3 .......... destination 1–4, depth 1–4 (per tone)
├── WAVE ........................ wave group/number/gain, tone switch, FXM, tone delay (per tone)
├── LFO
│   ├── LFO1 / LFO2 ............. waveform, key trigger, rate, ext. sync, fade mode,
│   │                            delay time, fade time, level offset (per tone)
│   └── LFO DEPTH 1:2 ........... pitch, TVF, TVA, pan depth for LFO1/LFO2 (per tone)
├── PITCH
│   ├── PITCH ................... coarse, fine, random, keyfollow (per tone)
│   ├── PCH ENV DPT ............. envelope depth, velocity sensitivity (per tone)
│   ├── PCH TIME ENV ............ V-T1, V-T4, time keyfollow (per tone)
│   └── PCH ENVELOPE ............ T1–T4, L1–L4 (per tone)
├── TVF
│   ├── FILTER .................. type, cutoff, resonance, keyfollow, env depth (per tone)
│   ├── TVF VELOCITY ............ v-sens, v-curve, v-resonance (per tone)
│   ├── TVF TIME ENV ............ V-T1, V-T4, time keyfollow (per tone)
│   └── TVF ENVELOPE ............ T1–T4, L1–L4 (per tone)
└── TVA
    ├── TVA ..................... tone level, tone pan, v-sens, v-curve (per tone)
    ├── BIAS .................... bias level, point, direction (per tone)
    ├── PAN MODULATION .......... keyfollow, random, alternate (per tone)
    ├── TVA TIME ENV ............ V-T1, V-T4, time keyfollow (per tone)
    └── TVA ENVELOPE ............ T1–T4, L1–L3 (per tone)
```

## Common

All offsets in the Patch common block (*1-3-1). Example:
`jv1080-cli set patch common <cli-name> <val>`.

### Patch name

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Patch name 1–12 | `00 00`–`00 0B` | 32–127 (ASCII) | `name` |

### Patch common

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Patch level `<Level>` | `00 2E` | 0–127 | `level` |
| Patch pan `<Pan>` | `00 2F` | 0–127 (L64–63R, 64 = center) | `pan` |
| Analog feel depth `<Analog Feel>` | `00 30` | 0–127 | `analogfeel` |
| Octave shift `<Octave>` | `00 41` | 0–6 (-3–+3) | `octave` |
| Stretch tune depth `<Stretch>` | `00 42` | 0–3 (OFF, 1–3) | `stretch` |
| Voice priority `<Priority>` | `00 43` | 0–1 (LAST, LOUDEST) | `priority` |
| Velocity range switch `<Velocity Range>` | `00 40` | 0–1 (OFF, ON) | `velrange` |
| Default tempo `<Default Tempo>` | `00 2C`–`00 2D` # | 20–250 | `tempo` |
| Structure type 1 & 2 `<Struct>` | `00 44` | 0–9 (1–10) | `struct12` |
| Structure type 3 & 4 `<Struct>` | `00 46` | 0–9 (1–10) | `struct34` |
| Booster level 1 & 2 `<Booster>` | `00 45` | 0–3 (0, +6, +12, +18) | `booster12` |
| Booster level 3 & 4 `<Booster>` | `00 47` | 0–3 (0, +6, +12, +18) | `booster34` |

### Velocity range / Key range (per Tone)

Offsets in each Patch tone block (*1-3-2); the enabling Velocity range switch is
in common (see above). Example: `jv1080-cli set patch tone <n> <cli-name> <val>`.

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Velocity range lower `<Lower>` | `00 0C` | 1–127 | `vellower` |
| Velocity range upper `<Upper>` | `00 0D` | 1–127 | `velupper` |
| Velocity cross fade `<X-Fade>` | `00 0B` | 0–127 | `velfade` |
| Keyboard range lower `<Lower>` | `00 0E` | 0–127 (C-1–G9) | `keylower` |
| Keyboard range upper `<Upper>` | `00 0F` | 0–127 (C-1–G9) | `keyupper` |

## Effects

EFX/Chorus/Reverb/EFX-CTRL parameters are in the Patch common block (*1-3-1);
the per-Tone output routing and send levels are in the tone blocks (*1-3-2).

### EFX (Patch EFX type / parameter / out / control)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| EFX type `<Type>` | `00 0C` | 0–39 | `efxtype` |
| EFX parameter 1–12 | `00 0D`–`00 18` | 0–127 each | `efxprm1`–`efxprm12` |
| EFX output assign `<Output Assign>` | `00 19` | 0–2 (MIX, OUTPUT1, OUTPUT2) | `efxout` |
| EFX output level `<Output Level>` | `00 1A` | 0–127 | `efxlevel` |
| EFX chorus send level `<Chorus>` | `00 1B` | 0–127 | `efxchorus` |
| EFX reverb send level `<Reverb>` | `00 1C` | 0–127 | `efxreverb` |
| EFX control source 1 | `00 1D` | 0–10 (OFF, SYS-CTRL1, SYS-CTRL2, MODULATION, BREATH, FOOT, VOLUME, PAN, EXPRESSION, BENDER, AFTERTOUCH) | `efxctrl1src` |
| EFX control depth 1 | `00 1E` | 0–126 (-63–+63) | `efxctrl1depth` |
| EFX control source 2 | `00 1F` | 0–10 (same list as source 1) | `efxctrl2src` |
| EFX control depth 2 | `00 20` | 0–126 (-63–+63) | `efxctrl2depth` |

### Chorus (Patch chorus)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Chorus rate `<Rat>` | `00 22` | 0–127 | `chorusrate` |
| Chorus depth `<Dpt>` | `00 23` | 0–127 | `chorusdepth` |
| Chorus pre delay `<Dly>` | `00 24` | 0–127 | `choruspredelay` |
| Chorus feedback `<Fbk>` | `00 25` | 0–127 | `chorusfeedback` |
| Chorus level `<Level>` | `00 21` | 0–127 | `choruslevel` |
| Chorus output assign `<Output>` | `00 26` | 0–2 (MIX, REVERB, MIX+REV) | `chorusout` |

### Reverb (Patch reverb)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Reverb type `<Type>` | `00 27` | 0–7 (ROOM1, ROOM2, STAGE1, STAGE2, HALL1, HALL2, DELAY, PAN-DLY) | `reverbtype` |
| Reverb level `<lev>` | `00 28` | 0–127 | `reverblevel` |
| Reverb time `<time>` | `00 29` | 0–127 | `reverbtime` |
| Reverb HF damp `<HF Damp>` | `00 2A` | 0–17 (200–8000 Hz, BYPASS) | `hfdamp` |
| Reverb/delay feedback `<fbk>` | `00 2B` | 0–127 | `reverbfeedback` |

### Per-Tone output (OUTPUT, per Tone)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Output assign `<Output Assign>` | `00 7D` | 0–3 (MIX, EFX, OUTPUT1, OUTPUT2) | `output` |
| Output level `<Output Level>` | `00 7E` | 0–127 | `outputlevel` |
| Chorus send level `<Chorus>` | `00 7F` | 0–127 | `chorus` |
| Reverb send level `<Reverb>` | `01 00` | 0–127 | `reverb` |

## Control

Key mode, portamento, hold/peak and control sources are Patch common (*1-3-1);
receive switches and the control destinations/depths are per Tone (*1-3-2).

### Key mode & bender (common)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Key assign mode `<assign>` | `00 33` | 0–1 (POLY, SOLO) | `assign` |
| Solo legato `<legato>` | `00 34` | 0–1 (OFF, ON) | `legato` |
| Bender range up | `00 31` | 0–12 | `bendup` |
| Bender range down | `00 32` | 0–48 (0 to -48) | `benddown` |

### Portamento (common)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Portamento switch `<sw>` | `00 35` | 0–1 (OFF, ON) | `portamento` |
| Portamento mode `<mode>` | `00 36` | 0–1 (NORMAL, LEGATO) | `portmode` |
| Portamento type `<type>` | `00 37` | 0–1 (RATE, TIME) | `porttype` |
| Portamento start `<start>` | `00 38` | 0–1 (PITCH, NOTE) | `portstart` |
| Portamento time `<tm>` | `00 39` | 0–127 | `porttime` |

### Receive switches (RxSWITCH / DAMPER, per Tone)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Volume control switch `<Volume>` | `00 11` | 0–1 (OFF, ON) | `rxvolume` |
| Pan control switch `<Pan>` | `00 14` | 0–2 (OFF, CONTINUOUS, KEY-ON) | `rxpan` |
| Bender control switch `<Bender>` | `00 13` | 0–1 (OFF, ON) | `rxbend` |
| Hold-1 control switch `<Hold-1>` | `00 12` | 0–1 (OFF, ON) | `rxhold` |
| Re-damper control switch `<ReDamper>` | `00 10` | 0–1 (OFF, ON) | `redamper` |

### Peak & hold (common)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| EFX control hold/peak `<EfxCtrl>` | `00 3C` | 0–2 (OFF, HOLD, PEAK) | `efxhold` |
| Control 1 hold/peak `<Ctrl1>` | `00 3D` | 0–2 (OFF, HOLD, PEAK) | `ctrl1hold` |
| Control 2 hold/peak `<Ctrl2>` | `00 3E` | 0–2 (OFF, HOLD, PEAK) | `ctrl2hold` |
| Control 3 hold/peak `<Ctrl3>` | `00 3F` | 0–2 (OFF, HOLD, PEAK) | `ctrl3hold` |

### Control source (common)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Patch control source 2 `<Control2>` | `00 3A` | 0–15 (OFF, SYS-CTRL1, SYS-CTRL2, MODULATION, BREATH, FOOT, VOLUME, PAN, EXPRESSION, BENDER, AFTERTOUCH, LFO1, LFO2, VELOCITY, KEYFOLLOW, PLAYMATE) | `ctrl2src` |
| Patch control source 3 `<Control3>` | `00 3B` | 0–15 (same list as source 2) | `ctrl3src` |

Control source 1 is fixed to Modulation (CC#1) and has no CLI name.

### Control 1 to 3 (per Tone)

Destinations 1–4 and depths 1–4 for each of the three controllers. Depth range
0–126 (-63–+63); destination list (0–18): OFF, PCH, CUT, RES, LEV, PAN, MIX,
CHO, REV, PL1, PL2, FL1, FL2, RL1, RL2, pL1, pL2, L1R, L2R.

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Controller 1 destination 1–4 | `00 15`, `00 17`, `00 19`, `00 1B` | 0–18 | `ctrl1dst1`–`ctrl1dst4` |
| Controller 1 depth 1–4 | `00 16`, `00 18`, `00 1A`, `00 1C` | 0–126 (-63–+63) | `ctrl1dep1`–`ctrl1dep4` |
| Controller 2 destination 1–4 | `00 1D`, `00 1F`, `00 21`, `00 23` | 0–18 | `ctrl2dst1`–`ctrl2dst4` |
| Controller 2 depth 1–4 | `00 1E`, `00 20`, `00 22`, `00 24` | 0–126 (-63–+63) | `ctrl2dep1`–`ctrl2dep4` |
| Controller 3 destination 1–4 | `00 25`, `00 27`, `00 29`, `00 2B` | 0–18 | `ctrl3dst1`–`ctrl3dst4` |
| Controller 3 depth 1–4 | `00 26`, `00 28`, `00 2A`, `00 2C` | 0–126 (-63–+63) | `ctrl3dep1`–`ctrl3dep4` |

## Wave

All parameters in this section are per Tone (*1-3-2). The wave-related names
match `waveParams` in `cmd/jv1080-cli/params.go`.

### WAVE

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Tone switch `<Switch>` | `00 00` | 0–1 (OFF, ON) | `switch` |
| Wave group `<Group>` | `00 01` | 0–2 (INT, PCM card, EXP) | `group` |
| Wave group ID | `00 02` | 0–127 (INT: 1 = INT-A, 2 = INT-B; card ID; expansion board no.) | `groupid` |
| Wave number `<Number>` | `00 03`–`00 04` # | 0–254 (wave 1–255) | `number` |
| Wave gain `<Gain>` | `00 05` | 0–3 (-6, 0, +6, +12) | `gain` |

### FXM

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| FXM switch `<Switch>` | `00 06` | 0–1 (OFF, ON) | `fxmswitch` |
| FXM color `<Color>` | `00 07` | 0–3 (1–4) | `fxmcolor` |
| FXM depth `<Depth>` | `00 08` | 0–15 (1–16) | `fxmdepth` |

### Tone delay

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Tone delay mode `<Mode>` | `00 09` | 0–6 (NORMAL, HOLD, KEY-INTERVAL (= PLAY-MATE), CLOCK-SYNC, TAP-SYNC, KEY-OFF-NORMAL, KEY-OFF-DECAY) | `delaymode` |
| Tone delay time `<time>` | `00 0A` | 0–127 | `delaytime` |

## LFO

All parameters in this section are per Tone (*1-3-2). LFO1 names carry no
suffix; LFO2 names carry the suffix `2` — except the depths, which always carry
the LFO number (`lfo1…`, `lfo2…`) because their destination is part of the
name.

### LFO1 / LFO2

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| LFO1 waveform `<form>` | `00 2D` | 0–7 (TRI, SIN, SAW, SQR, TRP, S&H, RND, CHS) | `lfowave` |
| LFO1 key trigger `<keytrig>` | `00 2E` | 0–1 (OFF, ON) | `lfokeytrig` |
| LFO1 rate `<rate>` | `00 2F` | 0–127 | `lforate` |
| LFO1 level offset `<Offset>` | `00 30` | 0–4 (-100, -50, 0, +50, +100) | `lfooffset` |
| LFO1 delay time `<Delay>` | `00 31` | 0–127 | `lfodelay` |
| LFO1 fade mode `<Mode>` | `00 32` | 0–3 (KEY-ON-IN, KEY-ON-OUT, KEY-OFF-IN, KEY-OFF-OUT) | `lfofade` |
| LFO1 fade time `<Fade>` | `00 33` | 0–127 | `lfofadetime` |
| LFO1 external sync `<extsync>` | `00 34` | 0–2 (OFF, CLOCK, TAP) | `lfosync` |
| LFO2 waveform `<form>` | `00 35` | 0–7 (TRI, SIN, SAW, SQR, TRP, S&H, RND, CHS) | `lfowave2` |
| LFO2 key trigger `<keytrig>` | `00 36` | 0–1 (OFF, ON) | `lfokeytrig2` |
| LFO2 rate `<rate>` | `00 37` | 0–127 | `lforate2` |
| LFO2 level offset `<Offset>` | `00 38` | 0–4 (-100, -50, 0, +50, +100) | `lfooffset2` |
| LFO2 delay time `<Delay>` | `00 39` | 0–127 | `lfodelay2` |
| LFO2 fade mode `<Mode>` | `00 3A` | 0–3 (KEY-ON-IN, KEY-ON-OUT, KEY-OFF-IN, KEY-OFF-OUT) | `lfofade2` |
| LFO2 fade time `<Fade>` | `00 3B` | 0–127 | `lfofadetime2` |
| LFO2 external sync `<extsync>` | `00 3C` | 0–2 (OFF, CLOCK, TAP) | `lfosync2` |

### LFO depth 1:2

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Pitch LFO1 depth `<Pitch>` | `00 4E` | 0–126 (-63–+63) | `pitchlfo1` |
| Pitch LFO2 depth `<Pitch>` | `00 4F` | 0–126 (-63–+63) | `pitchlfo2` |
| Filter LFO1 depth `<TVF>` | `00 63` | 0–126 (-63–+63) | `filterlfo1` |
| Filter LFO2 depth `<TVF>` | `00 64` | 0–126 (-63–+63) | `filterlfo2` |
| Amplitude LFO1 depth `<TVA>` | `00 75` | 0–126 (-63–+63) | `amplfo1` |
| Amplitude LFO2 depth `<TVA>` | `00 76` | 0–126 (-63–+63) | `amplfo2` |
| Pan LFO1 depth `<PAN>` | `00 7B` | 0–126 (-63–+63) | `panlfo1` |
| Pan LFO2 depth `<PAN>` | `00 7C` | 0–126 (-63–+63) | `panlfo2` |

## Pitch

All parameters in this section are per Tone (*1-3-2).

### PITCH

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Coarse tune `<Coarse>` | `00 3D` | 0–96 (-48–+48) | `coarse` |
| Fine tune `<Fine>` | `00 3E` | 0–100 (-50–+50) | `fine` |
| Random pitch depth `<Random>` | `00 3F` | 0–30 (display 0–1200) | `random` |
| Pitch keyfollow `<KeyFlw>` | `00 40` | 0–15 (-100–+200, 16 steps) | `keyfollow` |

### PCH ENV DPT (pitch envelope depth)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Pitch envelope depth `<Envelope Depth>` | `00 41` | 0–24 (-12–+12) | `envdepth` |
| Pitch env velocity sensitivity `<Velocity Sens>` | `00 42` | 0–125 (-100–+150 %) | `envvel` |

### PCH TIME ENV (pitch time envelope)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Pitch env velocity time 1 sens. `<V-T1>` | `00 43` | 0–14 (15-step list) | `envvt1` |
| Pitch env velocity time 4 sens. `<V-T4>` | `00 44` | 0–14 (same list) | `envvt4` |
| Pitch env time keyfollow `<Time Keyfollow>` | `00 45` | 0–14 (same list) | `envkeyfollow` |

### PCH ENVELOPE (pitch envelope)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Pitch env time 1–4 `<T1>`–`<T4>` | `00 46`–`00 49` | 0–127 each | `envt1`–`envt4` |
| Pitch env level 1–4 `<L1>`–`<L4>` | `00 4A`–`00 4D` | 0–126 (-63–+63) each | `envl1`–`envl4` |

## TVF

All parameters in this section are per Tone (*1-3-2). The cutoff/resonance
names match `tvfParams` in `cmd/jv1080-cli/params.go`.

### FILTER

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Filter type `<typ>` | `00 50` | 0–4 (OFF, LPF, BPF, HPF, PKG) | `filtertype` |
| Cutoff frequency `<cut>` | `00 51` | 0–127 | `cutoff` |
| Cutoff keyfollow `<keyflw>` | `00 52` | 0–15 (-100–+200, 16 steps) | `cutoffkeyfollow` |
| Resonance `<res>` | `00 53` | 0–127 | `resonance` |
| TVF envelope depth `<env dpt>` | `00 55` | 0–126 (-63–+63) | `envdepth` |

### TVF VELOCITY

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| TVF env velocity sensitivity `<v-sens>` | `00 57` | 0–125 (-100–+150 %) | `envvel` |
| TVF env velocity curve `<V-Curve>` | `00 56` | 0–6 (1–7) | `envcurve` |
| Resonance velocity sensitivity `<V-Resonance>` | `00 54` | 0–125 (-100–+150 %) | `resvel` |

### TVF TIME ENV (TVF time envelope)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Filter env velocity time 1 sens. `<V-T1>` | `00 58` | 0–14 (15-step list) | `envvt1` |
| Filter env velocity time 4 sens. `<V-T4>` | `00 59` | 0–14 (same list) | `envvt4` |
| Filter env time keyfollow `<Time Keyfollow>` | `00 5A` | 0–14 (same list) | `envkeyfollow` |

### TVF ENVELOPE

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Filter env time 1–4 `<T1>`–`<T4>` | `00 5B`–`00 5E` | 0–127 each | `envt1`–`envt4` |
| Filter env level 1–4 `<L1>`–`<L4>` | `00 5F`–`00 62` | 0–127 each | `envl1`–`envl4` |

## TVA

All parameters in this section are per Tone (*1-3-2).

### TVA

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Tone level `<Level>` | `00 65` | 0–127 | `level` |
| Tone pan `<Pan>` | `00 77` | 0–127 (L64–63R, 64 = center) | `pan` |
| Level env velocity sensitivity `<V-Sens>` | `00 6A` | 0–125 (-100–+150 %) | `envvel` |
| Level env velocity curve `<V-Curve>` | `00 69` | 0–6 (1–7) | `envcurve` |

### BIAS

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Bias level `<Bias>` | `00 68` | 0–14 (15-step list) | `bias` |
| Bias point `<Point>` | `00 67` | 0–127 (C-1–G9) | `biaspoint` |
| Bias direction `<Direction>` | `00 66` | 0–3 (LOWER, UPPER, LOWER&UPPER, ALL) | `biasdir` |

### PAN MODULATION

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Pan keyfollow `<KeyFlw>` | `00 78` | 0–14 (15-step list) | `pankeyfollow` |
| Random pan depth `<Random>` | `00 79` | 0–63 | `randompan` |
| Alternate pan depth `<Alternate>` | `00 7A` | 1–127 (L63–63R) | `altpan` |

### TVA TIME ENV (TVA time envelope)

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Level env velocity time 1 sens. `<V-T1>` | `00 6B` | 0–14 (15-step list) | `envvt1` |
| Level env velocity time 4 sens. `<V-T4>` | `00 6C` | 0–14 (same list) | `envvt4` |
| Level env time keyfollow `<Time Keyfollow>` | `00 6D` | 0–14 (same list) | `envkeyfollow` |

### TVA ENVELOPE

| Parameter | Offset | Valid values (SysEx) | CLI name |
| --- | --- | --- | --- |
| Level env time 1–4 `<T1>`–`<T4>` | `00 6E`–`00 71` | 0–127 each | `envt1`–`envt4` |
| Level env level 1–3 `<L1>`–`<L3>` | `00 72`–`00 74` | 0–127 each | `envl1`–`envl3` |
