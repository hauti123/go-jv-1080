# JV-1080 Patch Edit Mode — Menu Structure and SysEx Parameter Map

This document describes the device menu structure of the Patch Edit Mode
(Owner's Manual pp. 42–56) and maps every parameter to its SysEx representation.
Full parameter descriptions are in
[`patch-edit-mode.md`](patch-edit-mode.md); offsets and value ranges come from
the tables *1-3-1 Patch common* and *1-3-2 Patch tone* in
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

All offsets in the Patch common block (*1-3-1).

### Patch name

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Patch name 1–12 | `00 00`–`00 0B` | 32–127 (ASCII) |

### Patch common

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Patch level `<Level>` | `00 2E` | 0–127 |
| Patch pan `<Pan>` | `00 2F` | 0–127 (L64–63R, 64 = center) |
| Analog feel depth `<Analog Feel>` | `00 30` | 0–127 |
| Octave shift `<Octave>` | `00 41` | 0–6 (-3–+3) |
| Stretch tune depth `<Stretch>` | `00 42` | 0–3 (OFF, 1–3) |
| Voice priority `<Priority>` | `00 43` | 0–1 (LAST, LOUDEST) |
| Velocity range switch `<Velocity Range>` | `00 40` | 0–1 (OFF, ON) |
| Default tempo `<Default Tempo>` | `00 2C`–`00 2D` # | 20–250 |
| Structure type 1 & 2 `<Struct>` | `00 44` | 0–9 (1–10) |
| Structure type 3 & 4 `<Struct>` | `00 46` | 0–9 (1–10) |
| Booster level 1 & 2 `<Booster>` | `00 45` | 0–3 (0, +6, +12, +18) |
| Booster level 3 & 4 `<Booster>` | `00 47` | 0–3 (0, +6, +12, +18) |

### Velocity range / Key range (per Tone)

Offsets in each Patch tone block (*1-3-2); the enabling Velocity range switch is
in common (see above).

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Velocity range lower `<Lower>` | `00 0C` | 1–127 |
| Velocity range upper `<Upper>` | `00 0D` | 1–127 |
| Velocity cross fade `<X-Fade>` | `00 0B` | 0–127 |
| Keyboard range lower `<Lower>` | `00 0E` | 0–127 (C-1–G9) |
| Keyboard range upper `<Upper>` | `00 0F` | 0–127 (C-1–G9) |

## Effects

EFX/Chorus/Reverb/EFX-CTRL parameters are in the Patch common block (*1-3-1);
the per-Tone output routing and send levels are in the tone blocks (*1-3-2).

### EFX (Patch EFX type / parameter / out / control)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| EFX type `<Type>` | `00 0C` | 0–39 |
| EFX parameter 1–12 | `00 0D`–`00 18` | 0–127 each |
| EFX output assign `<Output Assign>` | `00 19` | 0–2 (MIX, OUTPUT1, OUTPUT2) |
| EFX output level `<Output Level>` | `00 1A` | 0–127 |
| EFX chorus send level `<Chorus>` | `00 1B` | 0–127 |
| EFX reverb send level `<Reverb>` | `00 1C` | 0–127 |
| EFX control source 1 | `00 1D` | 0–10 (OFF, SYS-CTRL1, SYS-CTRL2, MODULATION, BREATH, FOOT, VOLUME, PAN, EXPRESSION, BENDER, AFTERTOUCH) |
| EFX control depth 1 | `00 1E` | 0–126 (-63–+63) |
| EFX control source 2 | `00 1F` | 0–10 (same list as source 1) |
| EFX control depth 2 | `00 20` | 0–126 (-63–+63) |

### Chorus (Patch chorus)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Chorus rate `<Rat>` | `00 22` | 0–127 |
| Chorus depth `<Dpt>` | `00 23` | 0–127 |
| Chorus pre delay `<Dly>` | `00 24` | 0–127 |
| Chorus feedback `<Fbk>` | `00 25` | 0–127 |
| Chorus level `<Level>` | `00 21` | 0–127 |
| Chorus output assign `<Output>` | `00 26` | 0–2 (MIX, REVERB, MIX+REV) |

### Reverb (Patch reverb)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Reverb type `<Type>` | `00 27` | 0–7 (ROOM1, ROOM2, STAGE1, STAGE2, HALL1, HALL2, DELAY, PAN-DLY) |
| Reverb level `<lev>` | `00 28` | 0–127 |
| Reverb time `<time>` | `00 29` | 0–127 |
| Reverb HF damp `<HF Damp>` | `00 2A` | 0–17 (200, 250, 315, 400, 500, 630, 800, 1000, 1250, 1600, 2000, 2500, 3150, 4000, 5000, 6300, 8000 Hz, BYPASS) |
| Reverb/delay feedback `<fbk>` | `00 2B` | 0–127 |

### Per-Tone output (OUTPUT, per Tone)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Output assign `<Output Assign>` | `00 7D` | 0–3 (MIX, EFX, OUTPUT1, OUTPUT2) |
| Output level `<Output Level>` | `00 7E` | 0–127 |
| Chorus send level `<Chorus>` | `00 7F` | 0–127 |
| Reverb send level `<Reverb>` | `01 00` | 0–127 |

## Control

Key mode, portamento, hold/peak and control sources are Patch common (*1-3-1);
receive switches and the control destinations/depths are per Tone (*1-3-2).

### Key mode & bender (common)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Key assign mode `<assign>` | `00 33` | 0–1 (POLY, SOLO) |
| Solo legato `<legato>` | `00 34` | 0–1 (OFF, ON) |
| Bender range up | `00 31` | 0–12 |
| Bender range down | `00 32` | 0–48 (0 to -48) |

### Portamento (common)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Portamento switch `<sw>` | `00 35` | 0–1 (OFF, ON) |
| Portamento mode `<mode>` | `00 36` | 0–1 (NORMAL, LEGATO) |
| Portamento type `<type>` | `00 37` | 0–1 (RATE, TIME) |
| Portamento start `<start>` | `00 38` | 0–1 (PITCH, NOTE) |
| Portamento time `<tm>` | `00 39` | 0–127 |

### Receive switches (RxSWITCH / DAMPER, per Tone)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Volume control switch `<Volume>` | `00 11` | 0–1 (OFF, ON) |
| Pan control switch `<Pan>` | `00 14` | 0–2 (OFF, CONTINUOUS, KEY-ON) |
| Bender control switch `<Bender>` | `00 13` | 0–1 (OFF, ON) |
| Hold-1 control switch `<Hold-1>` | `00 12` | 0–1 (OFF, ON) |
| Re-damper control switch `<ReDamper>` | `00 10` | 0–1 (OFF, ON) |

### Peak & hold (common)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| EFX control hold/peak `<EfxCtrl>` | `00 3C` | 0–2 (OFF, HOLD, PEAK) |
| Control 1 hold/peak `<Ctrl1>` | `00 3D` | 0–2 (OFF, HOLD, PEAK) |
| Control 2 hold/peak `<Ctrl2>` | `00 3E` | 0–2 (OFF, HOLD, PEAK) |
| Control 3 hold/peak `<Ctrl3>` | `00 3F` | 0–2 (OFF, HOLD, PEAK) |

### Control source (common)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Patch control source 2 `<Control2>` | `00 3A` | 0–15 (OFF, SYS-CTRL1, SYS-CTRL2, MODULATION, BREATH, FOOT, VOLUME, PAN, EXPRESSION, BENDER, AFTERTOUCH, LFO1, LFO2, VELOCITY, KEYFOLLOW, PLAYMATE) |
| Patch control source 3 `<Control3>` | `00 3B` | 0–15 (same list as source 2) |

Control source 1 is fixed to Modulation (CC#1).

### Control 1 to 3 (per Tone)

Destinations 1–4 and depths 1–4 for each of the three controllers. Depth range
0–126 (-63–+63); destination list (0–18): OFF, PCH, CUT, RES, LEV, PAN, MIX,
CHO, REV, PL1, PL2, FL1, FL2, RL1, RL2, pL1, pL2, L1R, L2R.

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Controller 1 destination 1–4 | `00 15`, `00 17`, `00 19`, `00 1B` | 0–18 |
| Controller 1 depth 1–4 | `00 16`, `00 18`, `00 1A`, `00 1C` | 0–126 (-63–+63) |
| Controller 2 destination 1–4 | `00 1D`, `00 1F`, `00 21`, `00 23` | 0–18 |
| Controller 2 depth 1–4 | `00 1E`, `00 20`, `00 22`, `00 24` | 0–126 (-63–+63) |
| Controller 3 destination 1–4 | `00 25`, `00 27`, `00 29`, `00 2B` | 0–18 |
| Controller 3 depth 1–4 | `00 26`, `00 28`, `00 2A`, `00 2C` | 0–126 (-63–+63) |

## Wave

All parameters in this section are per Tone (*1-3-2).

### WAVE

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Tone switch `<Switch>` | `00 00` | 0–1 (OFF, ON) |
| Wave group `<Group>` | `00 01` | 0–2 (INT, PCM card, EXP) |
| Wave group ID | `00 02` | 0–127 (INT: 1 = INT-A, 2 = INT-B; card ID; expansion board no.) |
| Wave number `<Number>` | `00 03`–`00 04` # | 0–254 (wave 1–255) |
| Wave gain `<Gain>` | `00 05` | 0–3 (-6, 0, +6, +12) |

### FXM

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| FXM switch `<Switch>` | `00 06` | 0–1 (OFF, ON) |
| FXM color `<Color>` | `00 07` | 0–3 (1–4) |
| FXM depth `<Depth>` | `00 08` | 0–15 (1–16) |

### Tone delay

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Tone delay mode `<Mode>` | `00 09` | 0–6 (NORMAL, HOLD, KEY-INTERVAL (= PLAY-MATE), CLOCK-SYNC, TAP-SYNC, KEY-OFF-NORMAL, KEY-OFF-DECAY) |
| Tone delay time `<time>` | `00 0A` | 0–127 |

## LFO

All parameters in this section are per Tone (*1-3-2).

### LFO1 / LFO2

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| LFO1 waveform `<form>` | `00 2D` | 0–7 (TRI, SIN, SAW, SQR, TRP, S&H, RND, CHS) |
| LFO1 key trigger `<keytrig>` | `00 2E` | 0–1 (OFF, ON) |
| LFO1 rate `<rate>` | `00 2F` | 0–127 |
| LFO1 level offset `<Offset>` | `00 30` | 0–4 (-100, -50, 0, +50, +100) |
| LFO1 delay time `<Delay>` | `00 31` | 0–127 |
| LFO1 fade mode `<Mode>` | `00 32` | 0–3 (KEY-ON-IN, KEY-ON-OUT, KEY-OFF-IN, KEY-OFF-OUT) |
| LFO1 fade time `<Fade>` | `00 33` | 0–127 |
| LFO1 external sync `<extsync>` | `00 34` | 0–2 (OFF, CLOCK, TAP) |
| LFO2 waveform `<form>` | `00 35` | 0–7 (TRI, SIN, SAW, SQR, TRP, S&H, RND, CHS) |
| LFO2 key trigger `<keytrig>` | `00 36` | 0–1 (OFF, ON) |
| LFO2 rate `<rate>` | `00 37` | 0–127 |
| LFO2 level offset `<Offset>` | `00 38` | 0–4 (-100, -50, 0, +50, +100) |
| LFO2 delay time `<Delay>` | `00 39` | 0–127 |
| LFO2 fade mode `<Mode>` | `00 3A` | 0–3 (KEY-ON-IN, KEY-ON-OUT, KEY-OFF-IN, KEY-OFF-OUT) |
| LFO2 fade time `<Fade>` | `00 3B` | 0–127 |
| LFO2 external sync `<extsync>` | `00 3C` | 0–2 (OFF, CLOCK, TAP) |

### LFO depth 1:2

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Pitch LFO1 depth `<Pitch>` | `00 4E` | 0–126 (-63–+63) |
| Pitch LFO2 depth `<Pitch>` | `00 4F` | 0–126 (-63–+63) |
| Filter LFO1 depth `<TVF>` | `00 63` | 0–126 (-63–+63) |
| Filter LFO2 depth `<TVF>` | `00 64` | 0–126 (-63–+63) |
| Amplitude LFO1 depth `<TVA>` | `00 75` | 0–126 (-63–+63) |
| Amplitude LFO2 depth `<TVA>` | `00 76` | 0–126 (-63–+63) |
| Pan LFO1 depth `<PAN>` | `00 7B` | 0–126 (-63–+63) |
| Pan LFO2 depth `<PAN>` | `00 7C` | 0–126 (-63–+63) |

## Pitch

All parameters in this section are per Tone (*1-3-2).

### PITCH

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Coarse tune `<Coarse>` | `00 3D` | 0–96 (-48–+48) |
| Fine tune `<Fine>` | `00 3E` | 0–100 (-50–+50) |
| Random pitch depth `<Random>` | `00 3F` | 0–30 (0, 1–10, 20–100 step 10, 200–1200 step 100; display 0–1200) |
| Pitch keyfollow `<KeyFlw>` | `00 40` | 0–15 (-100, -70, -50, -30, -10, 0, +10, +20, +30, +40, +50, +70, +100, +120, +150, +200) |

### PCH ENV DPT (pitch envelope depth)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Pitch envelope depth `<Envelope Depth>` | `00 41` | 0–24 (-12–+12) |
| Pitch env velocity sensitivity `<Velocity Sens>` | `00 42` | 0–125 (-100–+150 %) |

### PCH TIME ENV (pitch time envelope)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Pitch env velocity time 1 sens. `<V-T1>` | `00 43` | 0–14 (-100, -70, -50, -40, -30, -20, -10, 0, +10, +20, +30, +40, +50, +70, +100) |
| Pitch env velocity time 4 sens. `<V-T4>` | `00 44` | 0–14 (same 15-step list) |
| Pitch env time keyfollow `<Time Keyfollow>` | `00 45` | 0–14 (same 15-step list) |

### PCH ENVELOPE (pitch envelope)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Pitch env time 1–4 `<T1>`–`<T4>` | `00 46`–`00 49` | 0–127 each |
| Pitch env level 1–4 `<L1>`–`<L4>` | `00 4A`–`00 4D` | 0–126 (-63–+63) each |

## TVF

All parameters in this section are per Tone (*1-3-2).

### FILTER

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Filter type `<typ>` | `00 50` | 0–4 (OFF, LPF, BPF, HPF, PKG) |
| Cutoff frequency `<cut>` | `00 51` | 0–127 |
| Cutoff keyfollow `<keyflw>` | `00 52` | 0–15 (-100, -70, -50, -30, -10, 0, +10, +20, +30, +40, +50, +70, +100, +120, +150, +200) |
| Resonance `<res>` | `00 53` | 0–127 |
| TVF envelope depth `<env dpt>` | `00 55` | 0–126 (-63–+63) |

### TVF VELOCITY

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| TVF env velocity sensitivity `<v-sens>` | `00 57` | 0–125 (-100–+150 %) |
| TVF env velocity curve `<V-Curve>` | `00 56` | 0–6 (1–7) |
| Resonance velocity sensitivity `<V-Resonance>` | `00 54` | 0–125 (-100–+150 %) |

### TVF TIME ENV (TVF time envelope)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Filter env velocity time 1 sens. `<V-T1>` | `00 58` | 0–14 (-100, -70, -50, -40, -30, -20, -10, 0, +10, +20, +30, +40, +50, +70, +100) |
| Filter env velocity time 4 sens. `<V-T4>` | `00 59` | 0–14 (same 15-step list) |
| Filter env time keyfollow `<Time Keyfollow>` | `00 5A` | 0–14 (same 15-step list) |

### TVF ENVELOPE

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Filter env time 1–4 `<T1>`–`<T4>` | `00 5B`–`00 5E` | 0–127 each |
| Filter env level 1–4 `<L1>`–`<L4>` | `00 5F`–`00 62` | 0–127 each |

## TVA

All parameters in this section are per Tone (*1-3-2).

### TVA

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Tone level `<Level>` | `00 65` | 0–127 |
| Tone pan `<Pan>` | `00 77` | 0–127 (L64–63R, 64 = center) |
| Level env velocity sensitivity `<V-Sens>` | `00 6A` | 0–125 (-100–+150 %) |
| Level env velocity curve `<V-Curve>` | `00 69` | 0–6 (1–7) |

### BIAS

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Bias level `<Bias>` | `00 68` | 0–14 (-100, -70, -50, -40, -30, -20, -10, 0, +10, +20, +30, +40, +50, +70, +100) |
| Bias point `<Point>` | `00 67` | 0–127 (C-1–G9) |
| Bias direction `<Direction>` | `00 66` | 0–3 (LOWER, UPPER, LOWER&UPPER, ALL) |

### PAN MODULATION

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Pan keyfollow `<KeyFlw>` | `00 78` | 0–14 (-100, -70, -50, -40, -30, -20, -10, 0, +10, +20, +30, +40, +50, +70, +100) |
| Random pan depth `<Random>` | `00 79` | 0–63 |
| Alternate pan depth `<Alternate>` | `00 7A` | 1–127 (L63–63R) |

### TVA TIME ENV (TVA time envelope)

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Level env velocity time 1 sens. `<V-T1>` | `00 6B` | 0–14 (-100, -70, -50, -40, -30, -20, -10, 0, +10, +20, +30, +40, +50, +70, +100) |
| Level env velocity time 4 sens. `<V-T4>` | `00 6C` | 0–14 (same 15-step list) |
| Level env time keyfollow `<Time Keyfollow>` | `00 6D` | 0–14 (same 15-step list) |

### TVA ENVELOPE

| Parameter | Offset | Valid values (SysEx) |
| --- | --- | --- |
| Level env time 1–4 `<T1>`–`<T4>` | `00 6E`–`00 71` | 0–127 each |
| Level env level 1–3 `<L1>`–`<L3>` | `00 72`–`00 74` | 0–127 each |
