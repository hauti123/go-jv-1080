# JV-1080 Patch Edit Mode — Manual Text (Owner's Manual, pp. 42–56)

Source: Roland JV-1080 Owner's Manual, Chapter 3 "Modes and Parameters", section
"1. Patch Edit Mode" (pages 42–56). Text extracted via OCR from the scanned manual;
screen parameter names are shown as `<Name>`, permissible values after the dash.
`[PALETTE]` marks parameters edited on a Palette Edit screen. The raw SysEx offsets
and value ranges for every parameter listed here are in
[`patch-edit-mode-structure.md`](patch-edit-mode-structure.md) and in
[`jv-1080-midi.md`](jv-1080-midi.md) (*1-3-1 Patch common, *1-3-2 Patch tone).

A Patch consists of a Common block and up to four Tones. The manual groups the
parameters into the following pages:

- **COMMON** — settings for the entire Patch (name, level, pan, velocity/key range, structure)
- **EFFECTS** — EFX, Chorus and Reverb added to the Patch, plus per-Tone output routing
- **CONTROL** — controllers, portamento, receive switches, control sources
- **WAVE** — waveform selection, FXM and Tone Delay for each Tone
- **LFO** — the two LFOs and their depths (per Tone)
- **PITCH** — pitch settings of a Tone
- **TVF** — filter settings of a Tone (Time Variant Filter)
- **TVA** — volume/amp settings of a Tone (Time Variant Amplifier)

COMMON, EFFECTS-output and the hold/peak settings are Patch-wide; everything from
WAVE onward (and the per-Tone output/chorus/reverb sends) is set individually for
each of the four Tones.

---

## Making Settings for an Entire Patch (COMMON)

The following parameters are all related to a single Patch.

### PATCH NAME

You can give a Patch a name of up to 12 characters.

Use the [◄]/[►] buttons to move the cursor to the proper position, and then rotate
the [VALUE] knob or press [INC]/[DEC] to select the desired character.

Available characters:

Space, A to Z, a to z, 0 to 9, and symbols including
`! " # % & ' ( ) * + , - . / : ; < = > ? @ [ ] ^ _ | ¥` (per the Owner's Manual).

Pressing [SHIFT] displays the following in the bottom right corner of the screen:

You can execute the following commands by holding down the [SHIFT] button and
pressing the FUNCTION SELECT button that corresponds to the function displayed:

- **TOP:** Press this button to successively go to the first character of the
  character groups ("A," "a," "0," or "+").
- **END:** Press this button to successively go to the final character of the
  character groups ("Z," "z," "9," or "_").
- **DEL:** Pressing this button erases the character under the cursor and shifts
  the following text forward.
- **INS:** Pressing this button inserts a space at the cursor location and shifts
  any following text backward.
- **CHG:** Pressing this button switches the character under the cursor from upper
  to lower case and vice versa.

### PATCH COMMON

**`<Level>` Patch Level — 0 to 127**

This parameter sets the overall volume for the entire Patch.

**`<Pan>` Patch Pan — L64 to 0 to 63R**

This parameter sets the stereo position of the entire Patch: L64 is far left, 0 is
center, and 63R is far right.

> Although a separate Pan setting is made for each of the Tones that make up a
> Patch, the stereo position of these Tones is shifted by this overall Pan setting.

**`<Analog Feel>` Analog Feel Depth — 0 to 127**

This parameter adds a very subtle (yet pleasing) pitch fluctuation to the basic
waveform. This helps create a sound that is more natural in nature. (This is also
called the "1/f fluctuation".)

**`<Octave>` Octave Shift — -3/-2/-1/0/+1/+2/+3**

This parameter changes the pitch of the Patch in one-octave steps. You can raise or
lower the pitch by up to three octaves.

> When Octave Shift is used in the Patch Play screen, this setting also changes
> automatically.

**`<Stretch>` Stretch Tune Depth — OFF/1/2/3**

This parameter changes the pitch using the 'stretch tuning' method typically used
on acoustic pianos. This makes high-range sounds slightly higher in pitch, and
low-range sounds slightly lower in pitch.

A diagram that illustrates the actual changes in pitch is called a "tuning curve."
Changing the tuning curve causes subtle alterations in the resonance of harmonics.

**`<Priority>` Voice Priority — LAST/LOUDEST**

This setting determines the order or priority of sounds to be played when the
maximum number of voices (64) is exceeded.

- **LAST:** The voices played last are given priority, with earlier voices being
  cut sequentially.
- **LOUDEST:** The voices with the highest volume are given priority, with quieter
  voices being cut sequentially.

**`<Velocity Range>` Velocity Range Switch — OFF/ON**

This setting determines whether the Velocity Range setting is enabled ("ON") or
disabled ("OFF").

**`<Default Tempo>` Default Tempo — 20 to 250 bpm**

This parameter determines the tempo clock setting, which controls the parameters
for Tones and Effects.

With the JV-1080, you can control Delay Time (and other time-related parameters)
with the internal tempo clock or with an external device. When no external tempo
clock is used, the tempo set here is used for such control.

- When using the internal tempo clock, set the clock source system parameter
  (p. 67) to "INT".
- The internal tempo clock data is not output through the MIDI OUT port.

### VELOCITY (Velocity Range)

**`<Lower>` Velocity Range Lower — 1 to 127 [PALETTE]**

This parameter determines the maximum range of sound when playing softly. Sound is
heard only when playing velocity exceeds the threshold (value) set here.

> If you try to make this value higher than the setting for Velocity Range Upper,
> the value for Velocity Range Upper changes at the same time.

**`<Upper>` Velocity Range Upper — 1 to 127 [PALETTE]**

This parameter determines the maximum range of sound when playing forcefully.
Sound is heard only when playing velocity does not exceed the value set here.

> If you try to make this value lower than the setting for Velocity Range Lower,
> the value for Velocity Range Lower changes at the same time.

**`<X-Fade>` Velocity Cross Fade Depth — 0 to 127 [PALETTE]**

This parameter sets the fade 'width' between the upper and lower velocity ranges.

### KEY RANG (Key Range)

**`<Lower>` Key Range Lower — C-1 to G9**

This parameter sets the key for the lowest playable sound.

> If you try to make this value higher than the key set for Key Range Upper, the
> value for Key Range Upper changes at the same time.

**`<Upper>` Key Range Upper — C-1 to G9**

This parameter sets the key for the highest playable sound.

> If you try to make this value lower than the key set for Key Range Lower, the
> value for Key Range Lower changes at the same time.

### STRUCT (Structure)

This parameter determines the structure of the Tones in the Patch. Differences in
structure can make the sound vary greatly.

- With Structure, two Tones are edited as a set. This means that it doesn't matter
  whether you press TONE SELECT [1] or [2]; in either case, "12—" appears in the
  left-hand corner of the screen. (With TONE SELECT [3] or [4], "—34" appears.)
- If you switch off one of the Tones in a pair while TYPE 2-10 is selected, the
  other Tone is played as TYPE1.

**`<Struct>` Structure — 1 to 10**

This parameter determines how Tones 1 and 2 or Tones 3 and 4 are combined. The
following combination types are possible.

> Because of space considerations, the characters on the screen are abbreviated.
> Here's what these abbreviations mean:
> **W1/W2:** WG (Wave Generator) 1 or 2 · **F1/F2:** TVF 1 or 2 ·
> **A1/A2:** TVA 1 or 2 · **B:** Booster · **R:** Ring Modulator

- **Type 1:** With this type, Tones 1 and 2 (or 3 and 4) are independent. Use this
  type when you want to preserve PCM sounds or create and combine sounds for each
  Tone.
- **Type 2:** This type combines two filters to enhance filter response. The TVA
  for Tone 1 (or 3) controls the volume balance of the two Tones.
- **Type 3:** This type combines two filters and distorts the waveforms by passing
  them through the booster. The TVA for Tone 1 (or 3) controls the volume balance
  of the two Tones and adjusts the amount of effect that the booster has.
- **Type 4:** This type applies a filter to the mixed sound of Tone 1 (or 3) and
  Tone 2 (or 4), and distorts the waveforms by passing them through the booster.
- **Type 5:** This type combines two filters and boosts the upper harmonics by
  passing the sound through the ring modulator. The TVA for Tone 1 (or 3) controls
  the volume balance of the two Tones and adjusts the depth of the ring modulator.
- **Type 6:** This type combines two filters, boosts the harmonics by passing the
  sound through the ring modulator, and mixes in the sound of Tone 2 (or 4). Because
  the sound from the ring modulator and Tone 2 (or 4) can be mixed, the TVA for
  Tone 1 (or 3) can adjust the amount of ring sound.
- **Type 7:** This type sends the filtered sound of Tone 1 (or 3) and the sound of
  Tone 2 (or 4) through the ring modulator to boost the harmonics.
- **Type 8:** This type sends the filtered sound of Tone 1 (or 3) and the sound of
  Tone 2 (or 4) through the ring modulator, then mixes it with Tone 2 (or 4) and
  filters the result.
- **Type 9:** This type sends the filtered sounds of the Tones through the ring
  modulator to boost the harmonics. The TVA for Tone 1 (or 3) controls the volume
  balance of the two Tones and adjusts the depth of the ring modulator.
- **Type 10:** This type sends the filtered Tones through the ring modulator to
  boost the harmonics, then mixes the result with Tone 2 (or 4). Because the sound
  from the ring modulator and Tone 2 (or 4) can be mixed, the TVA for Tone 1
  (or 3) can adjust the amount of ring sound.

#### What is a Booster?

The booster is a circuit that distorts the input signal.

In addition to using this for distortion, effects similar to PWM (Pulse Width
Modulation; where the harmonic structure changes continuously) can be obtained by
setting the waveform for one Tone (WG1) to a subsonic frequency and shifting the
waveform for the other Tone (WG2) up and down. It is also interesting to amplify
the waveform with Wave Gain (p. 50).

#### What is a Ring Modulator?

By combining the waveforms of the two Tones, a ring modulator can produce many
harmonics (inharmonic partials) not contained in either of the waveforms. (As long
as one of the waveforms is not a sine wave, virtually no frequency components sound
at regular intervals.) Because differences in the pitch of the waveform cause the
harmonic structure to change, a toneless metallic resonation occurs. This works well
when creating metallic timbres, such as for bells.

**`<Booster>` Booster Level — 0 / +6 / +12 / +18 [PALETTE]**

This parameter sets the booster level; a larger value results in greater distortion
of the sound.

---

## Selecting Effects for a Patch (EFFECTS)

> When editing the parameters for an effect, an "s" may appear on the left side of
> the screen. This indicates that the parameters for the effect are being edited
> while the effect's switch (p. 38) is off. Because this makes it impossible to
> check the results of the effect, press the [EFFECTS ON/OFF] button to turn the
> effect on.

### OUTPUT

**`<Output Assign>` Output Assign — MIX/EFX/OUTPUT1/OUTPUT2 [PALETTE]**
**`<Output Level>` Output Level — 0 to 127 [PALETTE]**

Output Assign determines the output destination for each Tone, and Output Level
sets the level of the output signal.

- **MIX:** This parameter sends the original sound to MIX OUT. The sound is also
  simultaneously sent to Chorus and Reverb.
- **EFX:** This parameter sends the original sound to EFX. The sound is also
  simultaneously sent to Chorus and Reverb. The output destination of the sound
  passing through EFX follows the Output Assign setting for PATCH EFX OUT.
- **OUTPUT1:** This setting sends the original sound to OUTPUT1.
- **OUTPUT2:** This setting sends the original sound to OUTPUT2.

> - When OUTPUT1 or OUTPUT2 is selected as the output destination, all settings
>   for Chorus and Reverb are ignored.
> - When TYPE 2 to 10 is selected with Structure (p. 43), the settings for Tone 1
>   (or Tone 3) are ignored.

**`<Chorus>` Chorus Send Level — 0 to 127 [PALETTE]**

This parameter sets the level of the signal sent to Chorus for each Tone.

**`<Reverb>` Reverb Send Level — 0 to 127 [PALETTE]**

This parameter sets the level of the signal sent to Reverb for each Tone.

### PATCH EFX TYPE

**`<Type>` EFX Type**

This parameter determines the type of EFX applied to the Patch.

For a description of the EFX types, check out "Chapter 5 — Multi-Effector EFX
(p. 83)."

### PATCH EFX PRM (Patch EFX Parameter)

This setting selects the parameter for the EFX chosen with EFX Type. The
parameters that can be set vary from one EFX to another.

For a description of EFX parameters, see "Chapter 5 — Multi-Effector EFX (p. 83)."

### PATCH EFX OUT (Patch EFX Output)

This parameter determines how sound will be output when "EFX" is selected with
Output Assign for OUTPUT.

**`<Output Assign>` Output Assign — MIX/OUTPUT1/OUTPUT2**
**`<Output Level>` Output Level — 0 to 127**

Output Assign determines the output destination for the EFX sound, and Output
Level sets the level of the output signal.

- **MIX:** This setting sends the EFX sound to MIX OUT. The sound is also
  simultaneously sent to Chorus and Reverb.
- **OUTPUT1:** This setting sends the EFX sound to OUTPUT1.
- **OUTPUT2:** This setting sends the EFX sound to OUTPUT2.

> When OUTPUT1 or OUTPUT2 is selected as the output destination, all settings for
> Chorus and Reverb are ignored.

**`<Chorus>` Chorus Send Level — 0 to 127**

This parameter sets the signal level sent to Chorus from EFX.

**`<Reverb>` Reverb Send Level — 0 to 127**

This parameter sets the signal level sent to Reverb from EFX.

### PATCH EFX CTRL (Patch EFX Control)

The JV-1080 lets you use any one of a number of MIDI Controllers to vary EFX
parameters in real time.

> The number and type of EFX parameters that can be changed is predetermined
> according to the EFX type. One or two EFX parameters are shown on the top line
> of the screen, and these vary according to the type of EFX that you choose. For
> details on the EFX parameters that can be changed, see "Chapter 5 — Multi-Effector
> EFX (p. 83)."

**EFX control source 1, 2 — EFX control depth 1, 2 (-63 to +63)**

Use EFX Control Source for each of the EFX parameters (that appear on the top line
of the screen) to choose a Controller, and set the degree of action of the
Controller with EFX Control Depth.

You can choose any of the following as the EFX Control Source:

- **OFF:** No Controller is used.
- **SYS-CTRL1:** The Controller set with System Control Source 1, a system
  parameter, is used (p. 68).
- **SYS-CTRL2:** The Controller set with System Control Source 2, a system
  parameter, is used (p. 68).
- **MODULATION:** Modulation (Control Change #1)
- **BREATH:** Breath (Control Change #2)
- **FOOT:** Foot (Control Change #4)
- **VOLUME:** Volume (Control Change #7)
- **PAN:** Pan (Control Change #10)
- **EXPRESSION:** Expression (Control Change #11)
- **BENDER:** Pitch bend
- **AFTERTOUCH:** Aftertouch

> Choose "SYS-CTRL1" or "SYS-CTRL2" if there is no need to use a different
> Controller for each Patch, or if you want to perform control with something other
> than the control changes described here. Use the system's Control Assign 1 page
> to set the Controller (p. 68).

### PATCH CHORUS

**`<Rat>` Chorus Rate — 0 to 127**

This parameter sets the speed of the Chorus effect.

**`<Dpt>` Chorus Depth — 0 to 127**

This parameter sets the depth of the Chorus effect.

**`<Dly>` Pre delay — 0 to 127**

This parameter determines the time interval between when the original sound is
heard and when the effect sound is heard. Larger values result in longer delays
(creating a broader sound).

**`<Fbk>` Chorus Feedback — 0 to 127**

This parameter determines the amount of Chorus sound that is returned (fed back)
to the Chorus unit. Larger values result in more complex Chorus effects.

**`<Level>` Chorus Level — 0 to 127**

This parameter sets the volume of the Chorus sound.

**`<Output>` Chorus Output Assign — MIX/REVERB/MIX+REV**

This setting determines how the Chorus sound is output.

- **MIX:** This setting sends the Chorus sound to MIX OUT.
- **REVERB:** This setting sends the Chorus sound to Reverb.
- **MIX+REV:** This setting sends the Chorus sound to both MIX OUT and Reverb.

### PATCH REVERB

**`<Type>` Reverb Type**

This parameter selects the type of reverberation.

- **ROOM1:** Short, high-density reverb
- **ROOM2:** Short, low-density reverb
- **STAGE1:** Reverb with many late reflections
- **STAGE2:** Reverb with strong initial reflection
- **HALL1:** Sparkling reverb
- **HALL2:** Richly resounding reverb
- **DELAY:** Standard delay
- **PAN-DLY:** A delay that pans (moves) the reflections to left and right.

**`<time>` Reverb Time — 0 to 127**

For a type from "ROOM1" to "HALL2," this parameter sets the reverb time (i.e. how
long the reverb continues). For "DELAY" or "PAN-DLY," this parameter sets the delay
time. Larger values produce a greater sense of spaciousness.

**`<lev>` Reverb Level — 0 to 127**

This parameter sets the volume of the reverb sound.

**`<fbk>` Delay Feedback — 0 to 127**

When "DELAY" or "PAN-DLY" has been selected as the type, this parameter sets the
amount of delayed sound that is returned (fed back) to the Delay unit. Larger
values result in longer delay times.

**`<HF Damp>` High-frequency Damp — 200/250/315/400/500/630/800/1000/1250/1600/2000/2500/3150/4000/5000/6300/8000/BYPASS**

This parameter sets the point (frequency) at which the high-frequency components of
the reverb sound are cut off.

The HF Damp (high-frequency damp) parameter is used to simulate the acoustic
properties of different materials. (For example, glass reflects more high
frequencies than heavy carpeting.)

A lower cutoff frequency results in a "darker" sound, and a higher frequency
produces a "brighter" sound. When set to "BYPASS," no high frequencies are cut.

---

## Using Controllers to Change How Sounds are Played (CONTROL)

This selects the functions of the JV-1080's Controllers.

### KEY MODE & BENDER

**`<assign>` Key Assign Mode — POLY/SOLO**

This parameter determines whether the Patch is played as a polyphonic (POLY) or
monophonic (SOLO) sound. Harmonies can be played when set to "POLY," and when
"SOLO" is selected, only one sound at a time can be played.

**`<legato>` Solo Legato — ON/OFF**

This parameter determines whether the Legato function is used (ON) or not (OFF).
However, the Legato function cannot be used when the Key Assign Mode is set to
"POLY."

> **What is Legato?** When you hold down one key and play another, the Legato
> function preserves any envelope or LFO, and changes only the pitch played. This
> makes it possible to emulate guitar techniques such as hammering-on or
> pulling-off.

**`<bend range>` Bend Range — -48 to 0 (Down) / 0 to 12 (Up)**

Bend Range Down lets you set how far the pitch drops when the Bender lever is moved
to the left (or when the wheel is lowered). You can set a pitch drop of up to four
octaves, in semitone increments. Bend Range Up lets you set how far the pitch rises
when the Bender lever is moved to the right (or when the wheel is raised). You can
set a pitch rise of up to one octave (in semitone increments).

### PORTAMENTO

Portamento is an effect that causes pitch to change smoothly when one played key is
followed by another.

**`<sw>` Portamento Switch — OFF/ON**

This parameter determines whether the Portamento effect is ON or OFF.

Portamento can be used in the SOLO Key Assign Mode to achieve a smooth 'sliding'
effect reminiscent of a violin playing style. Portamento can also be used in the
POLY Key Assign Mode (when playing chordal accompaniments).

**`<tm>` Portamento Time — 0 to 127**

This parameter determines the time required for the pitch to change from one note
to the next when using the Portamento function.

**`<mode>` Portamento Mode — NORMAL/LEGATO**

This parameter selects the Portamento mode. Portamento is always applied when
"NORMAL" is selected, but when set to "LEGATO," Portamento is applied only when you
play in a legato style (with the notes smoothly connected).

**`<type>` Portamento Type — RATE/TIME**

This parameter selects the type of Portamento effect. When set to "RATE," the time
required to move from the first pitch to the second is proportionate to the
distance between the two pitches. When "TIME" is selected, the change takes place
over a fixed amount of time (regardless of the pitch interval).

**`<start>` Portamento Start — PITCH/NOTE**

This parameter sets the point at which the Portamento effect begins.

This setting controls what happens when a new key is played while a change to
another pitch is already in progress.

When set to "PITCH," the new portamento sweep starts at the pitch in effect when
the new key is played. When set to "NOTE," the new portamento sweep starts from the
pitch in effect when the original pitch change is complete.

### RxSWITCH (Receive Switch)

**`<Volume>` Volume Control Switch — OFF/ON [PALETTE]**

This parameter determines whether volume changes (ON) or not (OFF) when MIDI volume
messages for the Tones are received.

> This setting is ignored when the Volume setting for Receive MIDI, a system
> parameter, is set to "OFF" (p. 68).

**`<Pan>` Pan Control Switch — OFF/CONT/KEY-ON [PALETTE]**

This parameter determines how MIDI pan messages are received for each of the Tones.

- **OFF:** The stereo position does not change when a pan message is received.
- **CONT:** Pan messages are received and the stereo position changes.
- **KEY-ON:** The stereo position changes when a key is played according to
  received pan messages. However, even if a new pan message is received while a
  sound is being played, there is no change in the pan position until the next
  time a key is played.

> This setting is ignored when the Control Change setting for Receive MIDI, a
> system parameter, is set to "OFF" (p. 68).

**`<Bender>` Bender Control Switch — OFF/ON [PALETTE]**

This parameter determines whether pitch changes (ON) or not (OFF) when MIDI bender
messages for each of the Tones are received.

> - This setting is ignored when the Bender setting for Receive MIDI, a system
>   parameter, is set to "OFF" (p. 68).
> - Ordinarily, volume messages control volume, pan messages control the stereo
>   position, and bender messages control the pitch being played. With the JV-1080,
>   however, you can use these messages to control other Tone parameters and
>   effects (p. 49). If you attempt this kind of use when the Receive Switches are
>   set to "ON," then the volume, panning, or pitch may change along with the
>   parameter you want to control. To avoid this, be sure to set these switches to
>   "OFF."

### DAMPER

**`<Hold-1 Switch>` Hold 1 Control Switch — OFF/ON [PALETTE]**

This parameter determines whether a sound is held (ON) or not (OFF) when MIDI Hold
1 messages are received for the Tones.

> This setting is ignored when the Hold 1 setting for Receive MIDI, a system
> parameter, is set to "OFF" (p. 68).

**`<ReDamper>` Re-damper Control Switch — OFF/ON [PALETTE]**

This setting determines (for each of the Tones) whether a decaying sound (the key
has been released) is held (ON) or not (OFF) when a Hold 1 message is received.

> When the Hold 1 Control Switch is set to "OFF," no Re-damper effect takes place
> even if the Re-damper Control Switch is set to "ON." This setting is also ignored
> when the Hold 1 setting for Receive MIDI, a system parameter, is set to "OFF"
> (p. 68).

### PEAK & HOLD

With the JV-1080, you can use the hold pedal to sustain not only note messages,
but also Control messages (such as modulation or aftertouch) when a Hold message
is received.

**`<EfxCtrl>` EFX Control Hold/Peak — OFF/HOLD/PEAK**

This setting determines how the value of an EFX parameter (to be changed with the
EFX Control Source; p. 46) is maintained when a Hold message is received.

- **OFF:** Hold messages are not received.
- **HOLD:** The parameter value is maintained when a Hold message is received.
- **PEAK:** The parameter value is revised and maintained each time it changes to
  a new maximum value, even after a Hold message is received.

> When you select the HOLD or PEAK settings, you will also need to specify what
> hold pedal information you will be using by means of these System parameters:
> Hold Control Source, and Peak Control Source (p. 67). Also, you will need to turn
> ON the Hold 1 Control switch and Receive MIDI switch (System parameter) so you
> are set for the hold pedal you are using (p. 68).

**`<Ctrl1>` Control 1 Hold/Peak — OFF/HOLD/PEAK**

This setting determines how the value of a Tone parameter to be changed with the
Patch Control Source 1 (Modulation) is maintained when a Hold message is received.
The setting values are the same as those for the EFX Control Hold/Peak setting.

**`<Ctrl2>` Control 2 Hold/Peak — OFF/HOLD/PEAK**

This setting determines how the value of a Tone parameter to be changed with the
Patch Control Source 2 is maintained when a Hold message is received. The setting
values are the same as those for the EFX Control Hold/Peak setting.

**`<Ctrl3>` Control 3 Hold/Peak — OFF/HOLD/PEAK**

This setting determines how the value of a Tone parameter to be changed with the
Patch Control Source 3 is maintained when a Hold message is received. The setting
values are the same as those for the EFX Control Hold/Peak setting.

### CONTROL SOURCE

The JV-1080 lets you use a variety of MIDI Controllers to change sounds in real
time. This setting selects the Controller used to change the Tone parameters.

> Patch Control Source 1 is set to Modulation (Control Change 1) and cannot be
> changed.

**`<Control2>` Patch Control Source 2**
**`<Control3>` Patch Control Source 3**

The Controllers assigned to Patch Controllers 2 and 3 can be selected from the
following list.

- **OFF:** No Controller is used.
- **SYS-CTRL1:** The Controller set with System Control Source 1, a system
  parameter, is used (p. 68).
- **SYS-CTRL2:** The Controller set with System Control Source 2, a system
  parameter, is used (p. 68).
- **MODULATION:** Modulation (Control Change #1)
- **BREATH:** Breath (Control Change #2)
- **FOOT:** Foot (Control Change #4)
- **VOLUME:** Volume (Control Change #7)
- **PAN:** Pan (Control Change #10)
- **EXPRESSION:** Expression (Control Change #11)
- **BENDER:** Pitch bend
- **AFTERTOUCH:** Aftertouch
- **LFO1:** The rate for LFO1 (The LFOs for Tones 1-4 act with respect to
  parameters within the same Tone.)
- **LFO2:** The rate for LFO2 (The LFOs for Tones 1-4 act with respect to
  parameters within the same Tone.)
- **VELOCITY:** Velocity
- **KEYFOLLOW:** Keyfollow (the parameter value changes according to the position
  on the keyboard, with C4 as "0")
- **PLAY-MATE:** Playmate (the parameter value changes according to the interval
  between Note On and Note Off)

> - Choose "SYS-CTRL1" or "SYS-CTRL2" if there is no need to use a different
>   Controller for each Patch, or if you want to have control with something other
>   than the Control Changes described here. Use the System's Control Assign 1 page
>   to set the Controller (p. 68).
> - When the setting for Receive MIDI, a system parameter, is set to "OFF," then
>   there is no effect when any Controller is used.

### CONTROL 1 to 3

**`<Destination>` Control Destination 1 to 4 [PALETTE]**

This parameter lets you simultaneously set up to four parameters for the Tone
changed by Patch Control Sources 1/2/3.

**`<Depth>` Control Depth 1 to 4 — -63 to +63 [PALETTE]**

This parameter sets the amount of change for any parameters set with Control
Destination.

Assignable Tone parameters (and ranges) are as follows:

| Destination (display) | Meaning | Depth (setting range) |
| --- | --- | --- |
| OFF | No parameters are controlled | — |
| PCH | Pitch | -63 — +63 (*1) |
| CUT | Cutoff frequency | -63 — +63 (*1) |
| RES | Resonance | -63 — +63 (*1) |
| LEV | Level (volume) | -63 — +63 (*1) |
| PAN | Pan (stereo position) | -63 — +63 (*1) |
| MIX | Volume output from MIX OUT | -63 — +63 (*1) |
| CHO | Chorus depth | -63 — +63 (*1) |
| REV | Reverb depth | -63 — +63 (*1) |
| PL1 | Depth of LFO1 applied to pitch | -63 — +63 (*2) |
| PL2 | Depth of LFO2 applied to pitch | -63 — +63 (*2) |
| FL1 | Depth of LFO1 applied to cutoff | -63 — +63 (*2) |
| FL2 | Depth of LFO2 applied to cutoff | -63 — +63 (*2) |
| RL1 | Depth of LFO1 applied to volume | -63 — +63 (*2) |
| RL2 | Depth of LFO2 applied to volume | -63 — +63 (*2) |
| pL1 | Depth of LFO1 applied to panning | -63 — +63 (*2) |
| pL2 | Depth of LFO2 applied to panning | -63 — +63 (*2) |
| L1R | LFO1 rate | -63 — +63 (*3) |
| L2R | LFO2 rate | -63 — +63 (*3) |

> \*1 Change is larger (higher) for positive (+) values and smaller (lower) for
> negative (-) values.
> \*2 The LFO phase is reversed for positive and negative values. In either case,
> however, depth increases as the value moves farther from 0 (zero).
> \*3 The LFO cycle is shorter for positive values and longer for negative values.

---

## Selecting Waveforms (WAVE)

This function lets you select the waveform that serves as the basis for a Tone,
apply effects to the waveform, and control its pitch.

### WAVE

**`<Group>` Wave Group — INT-A or B/CARD/EXP-A to D [PALETTE]**

This parameter determines the memory from which a waveform is selected.

- **INT-A or B:** Selects a waveform stored in the JV-1080.
- **CARD:** Selects a waveform stored on a PCM Card.
- **EXP-A to D:** Selects a waveform stored on Expansion Board A, B, C, or D.

**`<Number>` Wave Number — 1 to 255 [PALETTE]**

Selects the waveform that forms the basis of a Tone. Along with the wave number,
the wave name appears in the display (in parentheses).

**`<Gain>` Wave Gain — -6/0/+6/+12 [PALETTE]**

This parameter changes the gain of a waveform and is effective when adjusting the
level of a Tone (or using an amplified waveform after sending it through the Booster
(p. 44). The value is displayed in decibels (dB).

**`<Switch>` Tone Switch — OFF/ON [PALETTE]**

This parameter determines whether each Tone is to be played (ON) or not (OFF).

> - When TONE SWITCH [1] to [4] is used to switch a Tone on or off, this setting
>   also changes automatically.
> - The number of Tones played can be limited by switching off unused Tones.

### FXM (Frequency Cross Modulation)

FXM combines a specified wave with another waveform to create a third waveform.
This is useful for generating metallic timbres or making special effect sounds.

**`<Switch>` FXM Switch — OFF/ON [PALETTE]**

This parameter determines whether FXM is used (ON) or not (OFF).

**`<Color>` FXM Color — 1 to 4 [PALETTE]**

This parameter determines the ambience of FXM. A larger value results in a
'rougher' sound, and a smaller value makes the sound more metallic.

**`<Depth>` FXM Depth — 1 to 16 [PALETTE]**

This parameter sets the depth of FXM application.

### TONE DELAY

This parameter sets the delay applied to each Tone. Tone Delay is different from the
Effector's delay because it can also be used to change the tone of the delayed
sound and play one-key arpeggios while changing the pitch of each Tone.

The delay time can also be made to synchronize with the internal or external MIDI
clock.

**`<Mode>` Tone Delay Mode [PALETTE]**

Any of the following can be selected for Tone Delay.

- **NORMAL:** The Tone is always played according to the delay time.
- **HOLD:** The delay effect is applied only while the note is on. No Tone is played
  if the note switches off before the end of the delay time.
- **PLAY-MATE:** If the time interval between the previous Note On and the current
  Note On is less than two seconds, the Tone is played with that interval taken as
  the delay time. *(SysEx value 2 is called KEY-INTERVAL in the MIDI implementation,
  see the structure document.)*
- **CLOCK-SYNC:** The Tone is played using the on-board (or external) tempo clock
  (p. 67) setting.
- **TAP-SYNC:** The Tone is played with the interval set by tapping the foot pedal
  (Tap Control Source, p. 67).
- **KEY-OFF-NORMAL:** The Tone is played with the delay time after the note is off.
- **KEY-OFF-DECAY:** The TVA envelope starts from the time the note is on (during
  which time nothing is played), and the Tone is played with the delay time after
  the note is off.

**`<time>` Tone Delay Time — 0 to 127 / 0 to 880 (Note Display) [PALETTE]**

This parameter sets the time interval between a Note On and when sound is actually
heard (or from a Note Off if KEY-OFF-NORMAL or KEY-OFF-DECAY has been selected). If
"PLAY-MATE" has been selected, setting this to "64" causes the interval between the
previous Note On and the current Note On to be taken as the delay time. When set to
"127," an interval approximately twice as long as for "64" is set as the delay time.

> When "CLOCK-SYNC" or "TAP-SYNC" has been selected as the mode, the setting value
> is displayed as a Quarter-note resolution. A note corresponding to this is
> displayed next to the setting value. For instance, if set to 96=↓, the delay time
> would be 0.5 seconds at a tempo of 120.

| Note setting value | Quarter-note resolution |
| --- | --- |
| (Double note) | 768 |
| (Dotted whole note) | 576 |
| (Double-note triplet) | 512 |
| (Whole note) | 384 |
| (Dotted half note) | 288 |
| (Whole-note triplet) | 256 |
| (Half note) | 192 |
| (Dotted quarter note) | 144 |
| (Half-note triplet) | 128 |
| (Quarter note) | 96 |
| (Dotted eighth note) | 72 |
| (Quarter-note triplet) | 64 |
| (Eighth note) | 48 |
| (Dotted sixteenth note) | 36 |
| (Eighth-note triplet) | 32 |
| (Sixteenth note) | 24 |
| (Dotted thirty-second note) | 18 |
| (Sixteenth-note triplet) | 16 |
| (Thirty-second note) | 12 |

> When any of Type 2 to 10 has been selected with Structure (p. 43), the Tone Delay
> setting for Tone 2 is used for TONEs 1 & 2, and the Tone Delay setting for Tone 4
> is used for TONEs 3 & 4.

---

## Adding Undulations to a Sound (LFO)

An LFO (low-frequency oscillator) cyclically changes the pitch, cutoff frequency,
and level to produce modulation effects such as vibrato, wow, and tremolo.

The JV-1080 has two independent LFOs, called LFO1 and LFO2.

### LFO1

**`<form>` LFO Waveform [PALETTE]**

Any of the following LFO waveforms can be selected (per the MIDI implementation:
TRI, SIN, SAW, SQR, TRP, S&H, RND, CHS).

**`<keytrig>` Key Trigger — OFF/ON [PALETTE]**

This parameter determines whether the timing with which the keyboard is played
matches (ON) the LFO cycle or not (OFF).

**`<rate>` LFO Rate — 0 to 127 / 0 to 880 (Note Display)**

This parameter sets the speed of the cycle of the LFO. When "CLOCK" or "TAP" has
been selected with LFO External Sync, the setting value is displayed as a
Quarter-note resolution. A note corresponding to this is displayed next to the
setting value.

See the Tone Delay chart (p. 51) for information on note resolution.

> When "CHS" has been selected for LFO Waveform, LFO Rate has no effect.

**`<extsync>` LFO External Sync — OFF/CLOCK/TAP**

This parameter determines whether the LFO cycle is synchronized with the on-board
(or external) tempo clock (CLOCK), or with the depression rate of the foot pedal
(TAP). No synchronization takes place when set to "OFF."

> - When set to "CLOCK," use the system parameter Clock Source (p. 67) to choose
>   whether the internal or external tempo clock is to be used.
> - When set to "TAP," use the system parameter Tap Control Source (p. 67) to select
>   the foot pedal.

**`<Mode>` Fade Mode [PALETTE]**

This parameter selects how the LFO is applied.

- **ON-IN:** The LFO is applied gradually after the key is played.
- **ON-OUT:** The LFO is applied after the key is played, then gradually fades
  away.
- **OFF-IN:** The LFO is applied gradually after the key is released.
- **OFF-OUT:** The LFO is applied from the time the key is played until it is
  released, and after the release, it fades away gradually.

**`<Delay>` Delay Time — 0 to 127 [PALETTE]**

For ON-IN, this parameter sets the interval from the time the key is played until
the LFO starts to be applied (for ON-OUT, the hold time). For OFF-IN, this parameter
sets the interval from the time the key is released until the LFO starts to be
applied (for OFF-OUT, the hold time).

**`<Fade>` Fade Time — 0 to 127 [PALETTE]**

This parameter sets the interval after the delay time until the LFO amplitude
reaches its maximum (or minimum) value.

**`<Offset>` Level Offset — -100/-50/0/+50/+100 [PALETTE]**

This parameter shifts the LFO waveform up or down from the central values (pitch,
cutoff frequency, and level). Positive values cause the waveform to shift for
undulations above the central values, and negative values cause the waveform to
shift for undulations below the central values.

### LFO 2 Pages

These make settings for LFO2.

> The parameters that can be set are identical to those for LFO1.

### LFO DEPTH 1:2 (LFO Depth)

**`<Pitch>` Pitch LFO Depth 1 or 2 — -63 to +63 [PALETTE]**

This parameter sets the extent of application when LFO1 or LFO2 is applied to the
pitch.

**`<TVF>` Filter LFO Depth 1 or 2 — -63 to +63 [PALETTE]**

This parameter sets the extent of application when LFO1 or LFO2 is applied to the
cutoff frequency.

**`<TVA>` Amplitude LFO Depth 1 or 2 — -63 to +63 [PALETTE]**

This parameter sets the extent of application when LFO1 or LFO2 is applied to the
level (volume).

**`<PAN>` Pan LFO Depth 1 or 2 — -63 to +63 [PALETTE]**

This parameter sets the extent of application when LFO1 or LFO2 is applied to the
pan (stereo position).

> Negative and positive values for depth have opposite effects on the changes in
> pitch and volume. For example, if a positive depth value is set for one Tone and
> a negative value of the same magnitude is set for another Tone, the phase of the
> undulations is reversed. This makes it possible to alternately sound different
> Tones and shift panning cyclically.

---

## Changing the Pitch (PITCH)

This function sets parameters related to the pitch of a Tone.

### PITCH

**`<Coarse>` Coarse Tune — -48 to +48 [PALETTE]**

This parameter raises/lowers the pitch of a Tone by up to four octaves (in semitone
steps).

**`<Fine>` Fine Tune — -50 to +50 [PALETTE]**

This parameter raises/lowers the pitch of a Tone in increments of one cent (1/100th
of a semitone). You can fine-tune the Tone by up to a quarter-tone in either
direction.

**`<Random>` Random Pitch Depth — 0 to 1200 [PALETTE]**

This deliberately causes the pitch of the Tone being played to fluctuate
irregularly. The value is displayed in hundredths of a half-step.

**`<KeyFlw>` Pitch Key Follow — -100/-70/-50/-30/-10/0/+10/+20/+30/+40/+50/+70/+100/+120/+150/+200 [PALETTE]**

This sets how high the pitch rises when the keyboard has been shifted up an octave
(12 keys). When set to "+100," then when shifted upward one octave like an ordinary
keyboard instrument, the pitch also rises by one octave. When set to "+200,"
however, the pitch rises two octaves, and when set to "-100" it drops an octave.
When set to "0," every key plays at the same pitch. This is normally left set at
"+100."

### PCH ENV DPT (Pitch Envelope Depth)

**`<Envelope Depth>` Pitch Envelope Depth — -12 to +12 [PALETTE]**

This sets the extent of the effectiveness of the Pitch Envelope. Larger values,
either positive or negative, result in a greater range of change for the Pitch
Envelope. Negative values reverse the form of the envelope.

**`<Velocity Sens>` Pitch Envelope Velocity Sensitivity — -100 to +150 [PALETTE]**

This changes the level of the Pitch Envelope according to velocity. When the value
is positive, a larger velocity value results in a larger Pitch Envelope level. When
negative, a larger velocity value results in a smaller Pitch Envelope level.

### PCH TIME ENV (Pitch Time Envelope)

**`<V-T1>` Velocity Time 1 Sensitivity — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the Pitch Envelope T1 according to velocity. When the value is
positive, a larger velocity value results in faster change for T1. The change is
slowed when the value is negative.

**`<V-T4>` Velocity Time 4 Sensitivity — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the Pitch Envelope T4 according to the key off velocity. When the value
is positive, a larger velocity value results in faster change for T4. The change is
slowed when the value is negative.

> This effect is not applied when connected to a keyboard that cannot send key off
> velocity.

**`<Time Keyfollow>` Envelope Time Key Follow — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the Pitch Envelope according to the keyboard position. The envelope time
of the C4 key is used as the reference. When the value is positive, keys farther to
the right side of the keyboard (the high range) produce shorter times up to T2 to T4.
Negative values result in longer times.

### PCH ENVELOPE (Pitch Envelope)

**`<T1, T2, T3, or T4>` Pitch Envelope Time 1, 2, 3, or 4 — 0 to 127 [PALETTE]**

This sets the time T1, T2, T3, or T4 for the Pitch Envelope. Larger values result
in longer times until the next pitch is reached (for instance, T2 controls the time
from L1 until L2 is reached).

**`<L1, L2, L3, or L4>` Pitch Envelope Level 1, 2, 3, or 4 — 0 to 127 [PALETTE]**

This sets the level L1, L2, L3, or L4 for the Pitch Envelope. It determines how
much the pitch changes from the reference pitch (set with Coarse Tune or Fine Tune)
at each point. The change is higher than the reference pitch when this value is
positive and lower than the reference pitch when negative.

---

## Changing the Sound (TVF)

You can manipulate Tones with a TVF (Time Variant Filter), changing the brightness,
thickness, and other aspects of the sound.

### FILTER

**`<typ>` Filter Type — OFF/LPF/BPF/HPF/PKG [PALETTE]**

This selects the type of filter. A filter is a function that cuts off a specific
frequency band to change a sound's brightness, thickness, and other qualities. If
you don't want to use a filter, set this to "OFF."

- **LPF (Low-pass Filter):** Components higher than the cutoff frequency are
  eliminated. High-range frequencies are cut off, so the sound is rounder. This is
  the most often-used filter.
- **BPF (Bandpass Filter):** Only components near the cutoff frequency are allowed
  to remain — all others are eliminated. This is good for making highly distinctive
  sounds.
- **HPF (High-pass Filter):** Components lower than the cutoff frequency are
  eliminated. This filter is good for making percussion instruments with distinctive
  high ranges.
- **PKG (Peaking Filter):** The components near the cutoff frequency are
  emphasized. This expresses the special sound of a drum, and can even be used to
  create a wow effect.

**`<cut>` Cutoff Frequency — 0 to 127 [PALETTE]**

This specifies the frequency at which the filter effect starts (the cutoff
frequency) with respect to the frequency components of a waveform.

When using LPF, lowering the cutoff frequency reduces high harmonics for a
"rounder" sound, and a higher cutoff frequency produces a brighter sound.

With BPF, the harmonic components that you hear vary according to the value for the
cutoff frequency. This is good for making highly distinctive sounds.

When using HPF, raising the cutoff frequency reduces the low harmonics, so the
bright components of the sound are emphasized.

With PKG, the emphasized harmonics that are played vary according to the value of
the cutoff frequency.

> The effects of the filters and cutoff frequencies vary greatly from one type of
> waveform to another.

**`<res>` Resonance — 0 to 127 [PALETTE]**

This lifts up the components of the sound near the cutoff frequency, making the
sound more distinctive.

> Setting this value too high may cause vibration and distortion of the sound.

**`<keyflw>` Cutoff Key Follow — -100/-70/-50/-30/-10/0/+10/+20/+30/+40/+50/+70/+100/+120/+150/+200 [PALETTE]**

This makes the cutoff frequency change according to the position of the key played
on the keyboard.

The cutoff frequency for the C4 key is used as the reference. When the value is
positive, keys farther to the right side of the keyboard (the high range) result in
higher cutoff frequencies, and negative values give lower cutoff frequencies.

**`<env dpt>` TVF Envelope Depth — -63 to +63 [PALETTE]**

This sets the extent of the effectiveness of the TVF Envelope (p. 55). Larger
values, either positive or negative, result in a greater range of change for the
TVF Envelope. Negative values reverse the form of the envelope.

### TVF VELOCITY

**`<v-sens>` TVF Envelope Velocity Sensitivity — -100 to +150 [PALETTE]**

This changes the level of the TVF Envelope according to velocity.

When the value is positive, a larger velocity value results in a larger TVF Envelope
level for a brighter sound. When negative, a larger velocity value results in a
smaller TVF Envelope level for a darker sound.

**`<V-Curve>` TVF Envelope Velocity Curve — 1 to 7 [PALETTE]**

This chooses the curve used when changing the cutoff frequency with velocity. The
shape of the curve corresponding to the curve number (1 to 7) appears on the screen.

**`<V-Resonance>` Resonance Velocity Sensitivity — -100 to +150 [PALETTE]**

This changes the amount of resonance applied according to the velocity.

When the value is positive, a larger velocity produces greater resonance. When
negative, resonance is reduced. The effect is not applied when the value is zero.

### TVF TIME ENV (TVF Time Envelope)

**`<V-T1>` Velocity Time 1 Sensitivity — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the TVF Envelope T1 according to velocity. When the value is positive,
a larger velocity value results in faster change for T1. The change is slowed when
the value is negative.

**`<V-T4>` Velocity Time 4 Sensitivity — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the TVF Envelope T4 according to the key off velocity. When the value
is positive, a larger velocity value results in faster change for T4. The change is
slowed when the value is negative.

> This effect is not applied when connected to a keyboard that cannot send key off
> velocity.

**`<Time Keyfollow>` Envelope Time Key Follow — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the TVF Envelope according to the keyboard position. The envelope time
of the C4 key is used as the reference. When the value is positive, keys farther to
the right side of the keyboard (the high range) produce shorter times up to T2 to
T4. Negative values result in longer times.

### TVF ENVELOPE

**`<T1, T2, T3, or T4>` TVF Envelope Time 1, 2, 3, or 4 — 0 to 127 [PALETTE]**

This sets the time T1, T2, T3, or T4 for the TVF Envelope. Larger values result in
longer times until the next cutoff frequency is reached (for instance, T2 controls
the time from L1 until L2 is reached).

**`<L1, L2, L3, or L4>` TVF Envelope Level 1, 2, 3, or 4 — 0 to 127 [PALETTE]**

This sets the level L1, L2, L3, or L4 for the TVF Envelope. It determines how much
the cutoff frequency changes from the reference cutoff frequency (set with Cutoff
Frequency) at each point. The change is higher than the reference cutoff frequency
when this value is positive and lower than the reference cutoff frequency when
negative.

---

## Changing the Volume (TVA)

You can use a TVA (Time Variant Amplifier) to select the change in volume for each
of the Tones.

### TVA

**`<Level>` Tone Level — 0 to 127 [PALETTE]**

This sets the volume of the Tone. It is mainly used to balance the volume with
other Tones.

> The volume of an entire Patch is determined by the Patch Common Patch Level
> (p. 42).

**`<Pan>` Tone Pan — L64 to 0 to 63R**

Sets the panning (localizes sound image) for each of the Tones. L64 is leftmost, 0
is centered, and 63R is rightmost.

**`<V-Sens>` TVA Envelope Velocity Sensitivity — -100 to +150 [PALETTE]**

This changes the level of the TVA Envelope according to velocity.

When the value is positive, a larger velocity value results in a larger TVA
Envelope level for a louder sound. When negative, a larger velocity value results in
a smaller TVA Envelope level for lower volume.

**`<V-Curve>` TVA Envelope Velocity Curve — 1 to 7 [PALETTE]**

This chooses the curve used when changing the volume with velocity. The shape of
the curve corresponding to the curve number (1 to 7) appears on the screen.

### BIAS

**`<Bias>` Bias Level — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

When you want to make the volume change according to the keyboard position, this
sets how much the volume changes as you move away from the bias point. When you set
a positive value, the volume gets louder as you move from the bias point. A negative
value makes the volume grow softer.

> Even when set to a positive value, the level will not exceed the maximum.

**`<Point>` Bias Point — C-1 to G9 [PALETTE]**

This sets the key used as the reference for varying the volume when you want to
change the volume according to the keyboard position.

**`<Direction>` Bias Direction — LOWER/UPPER/LOWER&UPPER/ALL [PALETTE]**

This sets the level for the direction to be changed — left or right — when you want
to change the volume according to the keyboard position.

- **LOWER:** The volume changes when you move to the left (lower register) from the
  bias point.
- **UPPER:** The volume changes when you move to the right (upper register) from
  the bias point.
- **LOWER&UPPER:** The volume changes symmetrically when you move in either
  direction from the bias point.
- **ALL:** The volume changes linearly as you move from the lower register to the
  upper register, with the bias point at the center.

### PAN MODULATION

**`<KeyFlw>` Pan Key Follow — -100/-70/-50/-30/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This makes the panning change according to the position of the key played on the
keyboard. The C4 key is used as the reference. When the value is positive, keys
farther to the right side of the keyboard (the high range) result in panning to the
right, and negative values give panning to the left.

**`<Random>` Random Pan Depth — 0 to 63 [PALETTE]**

Panning changes irregularly each time a key is played. The value you set for this
determines the range of change. There is no change when set to zero.

**`<Alternate>` Alternate Pan Depth — L63 to 0 to R63 [PALETTE]**

The sound image is panned alternatingly to the right and left each time a key is
played. The changes are made with the value set for Tone Pan at the center. There
is no change when set to zero. If you set two Tones respectively to L and R, the
stereo position is swapped each time they are played.

### TVA TIME ENV (TVA Time Envelope)

**`<V-T1>` Velocity Time 1 Sensitivity — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the TVA Envelope T1 according to velocity. When the value is positive,
a larger velocity value results in faster change for T1. The change is slowed when
the value is negative.

**`<V-T4>` Velocity Time 4 Sensitivity — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the TVA Envelope T4 according to the key off velocity. When the value
is positive, a larger velocity value results in faster change for T4. The change is
slowed when the value is negative.

> This effect is not applied when connected to a keyboard that cannot send key off
> velocity.

**`<Time Keyfollow>` Envelope Time Key Follow — -100/-70/-50/-40/-30/-20/-10/0/+10/+20/+30/+40/+50/+70/+100 [PALETTE]**

This varies the TVA Envelope according to the keyboard position. The envelope time
of the C4 key is used as the reference. When the value is positive, keys farther to
the right side of the keyboard (the high range) produce shorter times up to T2 to
T4. Negative values result in longer times.

### TVA ENVELOPE

**`<T1, T2, T3, or T4>` TVA Envelope Time 1, 2, 3, or 4 — 0 to 127 [PALETTE]**

This sets the time T1, T2, T3, or T4 for the TVA Envelope. Larger values result in
longer times until the next volume is reached (for instance, T2 controls the time
from L1 until L2 is reached).

**`<L1, L2, or L3>` TVA Envelope Level 1, 2, or 3 — 0 to 127 [PALETTE]**

This sets the level L1, L2 or L3 for the TVA Envelope. It determines how much the
volume changes from the reference volume (set with Tone Level) at each point.
