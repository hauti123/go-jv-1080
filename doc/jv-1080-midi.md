# Roland JV-1080 — MIDI Implementation
## Transkription der Seiten 148–156 des Owner's Manual
**Quelle:** JV-1080_OM.pdf (Bildschirm-Scan), Seiten 129–156 — Demo-Songs (S. 129) sowie der Abschnitt „MIDI Implementation“ (Roland Exclusive Messages · Channel Voice Messages · Channel Mode Messages · System Real Time · Exclusive Communications · Parameter Address Map · GS · Anhangstabellen · MIDI Implementation Chart).
**Quelle:** JV-1080_OM.pdf (Bildschirm-Scan), Seiten 148–156 — Abschnitt „MIDI Implementation" (Channel Mode Messages · System Real Time · Exclusive Communications · Parameter Address Map · GS · Anhangstabellen · MIDI Implementation Chart).

**Zur Qualität der Vorlage:** Das PDF ist ein leicht unscharfer Scan. Die automatische Texterkennung lieferte für diese Seiten teils stark fehlerhafte Ergebnisse (z. B. `BmH` statt `BnH`, `FBH` statt `FEH`, völlig entstellte Hex-Folgen). Höher aufgelöste Seitenbilder lagen nicht vor, daher wurde stattdessen jede Zeile gegengeprüft: über die MIDI-Norm (Standard-Codes wie F0H/F7H/EOX), über die arithmetische Verifikation **aller** Checksummen-Beispiele und über die identische Parameterkarte des Schwester-Manuals (Roland JV-1010, gleiche Adress-Map, sauberer Textlayer) sowie eine corpus-geprüfte Referenz des JV-1080-SysEx-Formats. Unleserliche bzw. rekonstruierte Stellen sind mit **(1)–(17)** markiert; Auflösung am Ende des Dokuments. Alle Tabellen sind zusätzlich so ausgerichtet, dass sie auch als reiner Text lesbar bleiben.

**Abkürzungen für wiederkehrende Wertelisten** (im Original in jeder Zeile ausgeschrieben; hier zur Lesbarkeit abgekürzt):

- **\*A** (Keyfollow, 16 Werte): −100, −70, −50, −30, −10, 0, +10, +20, +30, +40, +50, +70, +100, +120, +150, +200
- **\*B** (Keyfollow, 15 Werte): −100, −70, −50, −40, −30, −20, −10, 0, +10, +20, +30, +40, +50, +70, +100
- **\*C** (Controller-Ziele, 0–18): OFF, PCH (Pitch), CUT (Cutoff), RES (Resonance), LEV (Level), PAN, MIX, CHO (Chorus), REV (Reverb), PL1/PL2 (Pitch-LFO 1/2), FL1/FL2 (Filter-LFO 1/2), AL1/AL2 (Amp-LFO 1/2), pL1/pL2 (Pan-LFO 1/2), L1R/L2R (LFO-1/2-Rate)
- **\*D** (HF Damp, 0–17): 200, 250, 315, 400, 500, 630, 800, 1000, 1250, 1600, 2000, 2500, 3150, 4000, 5000, 6300, 8000, BYPASS
- **\*R** (Reverb-Typen, 0–7): ROOM1, ROOM2, STAGE1, STAGE2, HALL1, HALL2, DELAY, PAN-DLY
- **\*S** (EFX-Control-Quellen, 0–10): OFF, SYS-CTRL1, SYS-CTRL2, MODULATION, BREATH, FOOT, VOLUME, PAN, EXPRESSION, BENDER, AFTERTOUCH
- **\*P** (Patch-Control-Quellen, 0–15): OFF, SYS-CTRL1, SYS-CTRL2, MODULATION, BREATH, FOOT, VOLUME, PAN, EXPRESSION, BENDER, AFTERTOUCH, LFO1, LFO2, VELOCITY, KEYFOLLOW, PLAYMATE
- **\*M** (Tone Delay Mode, 0–6): NORMAL, HOLD, KEY-INTERVAL, CLOCK-SYNC, TAP-SYNC, KEY-OFF-NORMAL, KEY-OFF-DECAY
- **\*RP** (Random Pitch Depth, 0–30): 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 1100, 1200
- **\*V** (Velocity-Sensitivity, 0–125): −100 … +150 % (3)

**Notation:** Adressen/Größen sind 7-Bit-Hex; `#` vor einer Adresse = durch zwei Datenbytes dargestellter Wert („divided-by-two data"); `--` = „Need not be set" (braucht nicht eingestellt zu werden).


## Seite 129 — The Demo Song Names / Composer Profiles

The following Demo Songs and their composers' profiles are contained in the JV-1080:

### RISE

Music by David Goldblatt. Copyright (C) 1994, David Goldblatt Music.

David Goldblatt is a Los Angeles based composer/keyboardist who has performed and recorded with a broad spectrum of artists such as Tribal Tech with Gary Willis and Scott Henderson, Stanley Clarke, Dizzy Gillespie, Wayne Shorter, Joe Sample, etc. David's compositions were featured on two Tribal Tech CDs and Brandon Fields' latest CD. He is heard on many TV shows and movies, such as Jodie Foster's "Little Man Tate," as a session artist. He has written and arranged for TV shows, jingles, and various productions in all sectors of the music industry, including musical director for various TV shows and a 1989 world tour with Diana Ross, as musical director, keyboardist, and arranger, which included a live CD and a concert video for HBO.

### 1080 Rave

Music by Ryeland Allison. Copyright (C) 1994, Ryeland Allison.

Ryeland arranges simultaneous manifold particle vibrations to proclaim "Resonance." Resonance is to some extent qualified to contented apportion within disassociative continuance, proceeding space. When converged in a synchronous locus, he at once regenerates toward fluidic empathy.

He is honored to introduce this to you.

### StormWarning

Music by Mitsuru Sakaue. Copyright (C) 1994, Roland Corporation.

Mitsuru Sakaue began composing and doing arrangements for commercials and videos while still in school. In particular, his studio work earned for him a solid reputation. Currently, he produces commercial music and jingles for FM stations.

* These demo songs are intended for personal enjoyment/demonstration use only. Any other use (public performance, broadcast, sampling, duplication, transcription, etc.) is strictly prohibited and would constitute a violation of applicable copyright laws.

---

## Seite 130 — Roland Exclusive Messages · Data Format · Address-mapped Data Transfer

# Roland Exclusive messages

### 1. Data Format for Exclusive Messages

Roland's MIDI implementation uses the following data format for all Exclusive messages (type IV):

| Byte    | Description              |
| ------- | ------------------------ |
| F0H     | Exclusive Status         |
| 41H     | Manufacturer ID (Roland) |
| DEV     | Device ID                |
| MDL     | Model ID                 |
| CMD     | Command ID               |
| [BODY]  | Main data                |
| F7H     | End of exclusive         |

### •MIDI status: F0H, F7H

An Exclusive message must be flanked by a pair of status codes, starting with a Manufacturer ID immediately after F0H (MIDI version 1.0).

### •Manufacturer ID: 41H

The Manufacturer ID identifies the manufacturer of a MIDI instrument that sends an Exclusive message. Value 41H represents Roland's Manufacturer ID.

### •Device ID: DEV

The Device ID contains a unique value that identifies individual devices in the implementation of several MIDI instruments. It is usually set to 00H-0FH, a value smaller by one than that of a basic channel, but value 00H-1FH may be used for a device with several basic channels.

### •Model ID: MDL

The Model ID contains a value that identifies one model from another. Different models, however, may share an identical Model ID if they handle similar data.

The Model ID format may contain 00H in one or more places to provide an extended data field. The following are examples of valid Model IDs, each representing a unique model:

```
01H
02H
03H
00H, 01H
00H, 02H
00H, 00H, 01H
```

### •Command ID: CMD

The Command ID indicates the function of an Exclusive message. The Command ID format may contain 00H in one or more places to provide an extended data field. The following are examples of valid Command IDs, each representing a unique function:

```
01H
02H
03H
00H, 01H
00H, 02H
00H, 00H, 01H
```

### •Main data: BODY

This field contains a message to be exchanged across an interface. The exact data size and content will vary with the Model ID and Command ID.

### 2. Address-mapped Data Transfer

Address mapping is a technique for transferring messages conforming to the data format given in Section 1. It assigns a series of memory-resident records - waveform and tone data, switch status, and parameters, for example - to specific locations in a machine-dependent address space, thereby allowing access to data residing at the address a message specifies.

Address-mapped data transfer is therefore independent of models and data categories. This technique allows use of two different transfer procedures: one-way transfer and handshake transfer.

### •One-way transfer procedure (See Section 3 for details.)

This procedure is suited to the transfer of a small amount of data. It sends out an Exclusive message completely independent of the receiving device's status.

Connection Diagram

```text
              One-way transfer
Device A                          Device B
+--------+                       +--------+
| MIDI   |  ==== 1 ====>         | MIDI   |
|  OUT   |                       |  IN    |
+--------+                       +--------+
| MIDI   |         <=== 2 ====   | MIDI   |
|  IN    |                       |  OUT   |
+--------+                       +--------+
```

Connection at point 2 is essential for "Request data" procedures. (See Section 3.)

### •Handshake-transfer procedure (This device does not use this procedure)

This procedure initiates a predetermined transfer sequence (handshaking) across the interface before data transfer takes place. Handshaking ensures that reliability and transfer speed are high enough to handle a large amount of data.

Connection Diagram

```text
              Handshake transfer
Device A                          Device B
+--------+                       +--------+
| MIDI   |  ==== 1 ====>         | MIDI   |
|  OUT   |                       |  IN    |
+--------+                       +--------+
| MIDI   |         <=== 2 ====   | MIDI   |
|  IN    |                       |  OUT   |
+--------+                       +--------+
```

Connection at points 1 and 2 is essential.

### Notes on the above procedures

- There are separate Command IDs for different transfer procedures.
- Devices A and B cannot exchange data unless they use the same transfer procedure, share identical Device ID and Model ID, and are ready for communication.

### 3. One-way Transfer Procedure

This procedure sends out data until it has all been sent and is used when the messages are so short that answerbacks need not be checked.

For longer messages, however, the receiving device must acquire each message in time with the transfer sequence, which inserts 20 milliseconds intervals.

Types of Messages

| Message        | Command ID   |
| -------------- | ------------ |
| Request data 1 | RQ1 (11H)    |
| Data set 1     | DT1 (12H)    |

### •Request data #1: RQ1 (11H)

This message is sent out when there is a need to acquire data from a device at the other end of the interface. It contains data for the address and size that specify designation and length, respectively, of data required.

On receiving an RQ1 message, the remote device checks its memory for the data address and size that satisfy the request.

If it finds them and is ready for communication, the device will transmit a "Data set 1 (DT1)" message, which contains the requested data. Otherwise, the device won't send out anything.

| Byte | Description              |
| ---- | ------------------------ |
| F0H  | Exclusive Status         |
| 41H  | Manufacturer ID (Roland) |
| DEV  | Device ID                |
| MDL  | Model ID                 |
| 11H  | Command ID               |
| aaH  | Address MSB              |
| bbH  | Address                  |
| ccH  | Address                  |
| ddH  | Address LSB              |
| eeH  | Size MSB                 |
| ffH  | Size                     |
| ggH  | Size                     |
| hhH  | Size LSB                 |
| sum  | Check sum                |
| F7H  | End of exclusive         |

---

## Seite 131 — DT1 (12H) · Beispiele der Nachrichten-Transaktion (11)

- The size of the requested data does not indicate the number of bytes that will make up a DT1 message, but represents the address fields where the requested data resides.
- Some models are subject to limitations in data format used for a single transaction. Requested data, for example, may have a limit in length or must be divided into predetermined address fields before it is exchanged across the interface.
- The same number of bytes comprises address and size data, which, however, vary with the Model ID.
- The error-checking process uses a checksum that provides a bit pattern where the last 7 bits are zero when values for an address, size, and that checksum are summed.

### •Data set 1: DT1 (12H)

This message corresponds to the actual data transfer process.

Because every byte in the data is assigned a unique address, a DT1 message can convey the starting address of one or more bits of data as well as a series of data formatted in an address-dependent order.

The MIDI standards inhibit non real-time messages from interrupting an Exclusive one. This fact is inconvenient for devices that support a "soft-thru" function. To maintain compatibility with such devices, Roland has limited the DT1 to 256 bytes so that an excessively long message is sent out in separate 'segments'.

| Byte | Description              |
| ---- | ------------------------ |
| F0H  | Exclusive Status         |
| 41H  | Manufacturer ID (Roland) |
| DEV  | Device ID                |
| MDL  | Model ID                 |
| 12H  | Command ID               |
| aaH  | Address MSB              |
| bbH  | Address                  |
| ccH  | Address                  |
| ddH  | Address LSB              |
| eeH  | Data MSB                 |
| ffH  | Data (variable length)   |
| ggH  | Data                     |
| hhH  | Data LSB                 |
| sum  | Check sum                |
| F7H  | End of exclusive         |

- A DT1 message is capable of providing only the valid data among those specified by an RQ1 message.
- Some models are subject to limitations in data format used for a single transaction. Requested data, for example, may have a limit in length or must be divided into predetermined address fields before it is exchanged across the interface.
- The number of bytes comprising address data varies from one Model ID to another.
- The error-checking process uses a checksum that provides a bit pattern where the last 7 bits are zero when values for an address, size, and that checksum are summed.

### •Example of Message Transactions

### •Device A sending data to Device B

Transfer of a DT1 message is all that takes place.

```text
Device A                          Device B
+--------+                       +--------+
| MIDI   |  === DT1 ===>         | MIDI   |
|  OUT   |                       |  IN    |
+--------+                       +--------+
```

### •Device B requesting data from Device A

Device B sends an RQ1 message to Device A.

Checking the message, Device A sends a DT1 message back to Device B.

```text
Device A                          Device B
+--------+                       +--------+
| MIDI   |  === DT1 ===>         | MIDI   |
|  OUT   |                       |  IN    |
+--------+                       +--------+
| MIDI   |         <=== RQ1 ===  | MIDI   |
|  IN    |                       |  OUT   |
+--------+                       +--------+
```

---
## Seite 132 — MIDI Implementation: How to read the tables · Channel Voice Messages

# [ MULTI TIMBRAL SYNTHESIZER MODULE ]

# Model JV-1080

# MIDI Implementation

Date: JULY 9, 1994 — Version: 1.00 (12)

### ( How to read the tables )

This shows the results which you can get by setting each parameter. "Setting1" and "Setting2" show two different ways of setting to get the same result.

| Parameters to be set       | Level of tone will be changed | Tone parameters can be controlled in real-time |                 | Effect parameters can be controlled in real-time |                 |
| -------------------------- | ----------------------------- | ---------------------------------------------- | --------------- | ------------------------------------------------ | --------------- |
|                            |                               | Setting1                                       | Setting2        | Setting1                                         | Setting2        |
| SYSTEM PARAMETER           |                               |                                                |                 |                                                  |                 |
| System Control Source 1/2  | --                            | --                                             | CC11:EXPRESSION | --                                               | CC11:EXPRESSION |
| Receive Control Change     | ON                            | ON                                             | ON              | ON                                               | ON              |
| Volume Control Source      | VOL&EXP                       | --                                             | --              | --                                               | --              |
| PATCH COMMON PARAMETER     |                               |                                                |                 |                                                  |                 |
| EFX Control Source 1/2     | --                            | --                                             | --              | EXPRESSION                                       | SYS-CTRL1/2     |
| EFX Control Depth 1/2      | --                            | --                                             | --              | Other than 0                                     | Other than 0    |
| Patch Control Source 1/2/3 | --                            | EXPRESSION                                     | SYS-CTRL1/2     | --                                               | --              |
| PATCH TONE PARAMETER       |                               |                                                |                 |                                                  |                 |
| Volume Control Switch      | ON                            | --                                             | --              | --                                               | --              |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                            | Other than OFF                                 | Other than OFF  | --                                               | --              |
| Ctrl 1/2/3 Depth1/2/3/4    | --                            | Other than 0                                   | Other than 0    | --                                               | --              |

# 1. RECEIVE DATA

# ■ Channel Voice Message

### ● Note off

| Status | Second | Third |
| ------ | ------ | ----- |
| 8nH    | kkH    | vvH   |
| 9nH    | kkH    | 00H   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
kk = Note number: 00H - 7FH (0 - 127)
vv = Velocity: 00H - 7FH (0 - 127)

- In the performance mode, receives this message when the MIDI Receive Switch of each part parameter is ON.
- Rhythm part (part 10) receives this message when the envelope mode of a rhythm tone parameter is SUSTAIN.

### ● Note on

| Status | Second | Third |
| ------ | ------ | ----- |
| 9nH    | kkH    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
kk = Note number: 00H - 7FH (0 - 127)
vv = Velocity: 01H - 7FH (1 - 127)

- In the performance mode, receives this message when the MIDI Receive Switch of each part parameter is ON.

### ● Polyphonic key pressure

| Status | Second | Third |
| ------ | ------ | ----- |
| AnH    | kkH    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
kk = Note number: 00H - 7FH (0 - 127)
vv = Pressure value: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set       | Tone parameters can be controlled in real-time |                       | Effect parameters can be controlled in real-time |                       |
| -------------------------- | ---------------------------------------------- | --------------------- | ------------------------------------------------ | --------------------- |
|                            | Setting1                                       | Setting2              | Setting1                                         | Setting2              |
| SYSTEM PARAMETER           |                                                |                       |                                                  |                       |
| System Control Source 1/2  | --                                             | AFTERTOUCH            | --                                               | AFTERTOUCH            |
| Receive Aftertouch         | ON                                             | ON                    | ON                                               | ON                    |
| Aftertouch Source          | CH&POLY or POLY-AFTER                          | CH&POLY or POLY-AFTER | CH&POLY or POLY-AFTER                            | CH&POLY or POLY-AFTER |
| PATCH COMMON PARAMETER     |                                                |                       |                                                  |                       |
| EFX Control Source 1/2     | --                                             | --                    | AFTERTOUCH                                       | SYS-CTRL1/2           |
| EFX Control Depth 1/2      | --                                             | --                    | Other than 0                                     | Other than 0          |
| Patch Control Source 1/2/3 | AFTERTOUCH                                     | SYS-CTRL1/2           | --                                               | --                    |
| PATCH TONE PARAMETER       |                                                |                       |                                                  |                       |
| Ctrl 1/2/3 Dest.1/2/3/4    | Other than OFF                                 | Other than OFF        | --                                               | --                    |
| Ctrl 1/2/3 Depth1/2/3/4    | Other than 0                                   | Other than 0          | --                                               | --                    |

---

## Seite 133 — 2. Performance Mode (Aftertouch) · Control Change: Bank Select · Bank-/Media-Tabellen

### 2. Performance Mode

| Parameters to be set         | Tone parameters can be controlled in real-time |                       | Effect parameters can be controlled in real-time |                       |                       |                       |
| ---------------------------- | ---------------------------------------------- | --------------------- | ------------------------------------------------ | --------------------- | --------------------- | --------------------- |
|                              | Setting1                                       | Setting2              | Setting1                                         | Setting2              | Setting3              | Setting4              |
| SYSTEM PARAMETER             |                                                |                       |                                                  |                       |                       |                       |
| System Control Source 1/2    | --                                             | AFTERTOUCH            | --                                               | --                    | AFTERTOUCH            | AFTERTOUCH            |
| Receive Aftertouch           | ON                                             | ON                    | ON                                               | ON                    | ON                    | ON                    |
| Aftertouch Source            | CH&POLY or POLY-AFTER                          | CH&POLY or POLY-AFTER | CH&POLY or POLY-AFTER                            | CH&POLY or POLY-AFTER | CH&POLY or POLY-AFTER | CH&POLY or POLY-AFTER |
| PERFORMANCE COMMON PARAMETER |                                                |                       |                                                  |                       |                       |                       |
| EFX Source                   | --                                             | --                    | 1 - 16                                           | PERFORM               | 1 - 16                | PERFORM               |
| EFX Control Source 1/2       | --                                             | --                    | --                                               | AFTERTOUCH            | --                    | SYS-CTRL1/2           |
| EFX Control Depth 1/2        | --                                             | --                    | --                                               | Other than 0          | --                    | Other than 0          |
| PERFORMANCE PART PARAMETER   |                                                |                       |                                                  |                       |                       |                       |
| MIDI Receive Switch          | ON                                             | ON                    | ON                                               | ON                    | ON                    | ON                    |
| PATCH COMMON PARAMETER       |                                                |                       |                                                  |                       |                       |                       |
| EFX Control Source 1/2       | --                                             | --                    | AFTERTOUCH                                       | --                    | SYS-CTRL1/2           | --                    |
| EFX Control Depth 1/2        | --                                             | --                    | Other than 0                                     | --                    | Other than 0          | --                    |
| Patch Control Source 1/2/3   | AFTERTOUCH                                     | SYS-CTRL1/2           | --                                               | --                    | --                    | --                    |
| PATCH TONE PARAMETER         |                                                |                       |                                                  |                       |                       |                       |
| Ctrl 1/2/3 Dest.1/2/3/4      | Other than OFF                                 | Other than OFF        | --                                               | --                    | --                    | --                    |
| Ctrl 1/2/3 Depth1/2/3/4      | Other than 0                                   | Other than 0          | --                                               | --                    | --                    | --                    |

### ● Control Change

### ○ Bank select (MSB/LSB)

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 00H    | mmH   |
| BnH    | 20H    | llH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
mm = Upper bytes of bank number: 50H - 54H (80 - 84)
ll = Lower bytes of bank number: 00H - 07H (0 - 7)

### 1. Patch Mode

| Parameters to be set       | Changes Banks upon receiving program change | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |
| -------------------------- | ------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ |
| SYSTEM PARAMETER           |                                             |                                                |                                                  |
| System Control Source 1/2  | --                                          | CC0:BANK-SEL                                   | CC0:BANK-SEL                                     |
| Receive Bank Select        | ON                                          | ON                                             | ON                                               |
| PATCH COMMON PARAMETER     |                                             |                                                |                                                  |
| EFX Control Source 1/2     | --                                          | --                                             | SYS-CTRL1/2                                      |
| EFX Control Depth 1/2      | --                                          | --                                             | Other than 0                                     |
| Patch Control Source 1/2/3 | --                                          | SYS-CTRL1/2                                    | --                                               |
| PATCH TONE PARAMETER       |                                             |                                                |                                                  |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                                          | Other than OFF                                 | --                                               |
| Ctrl 1/2/3 Depth1/2/3/4    | --                                          | Other than 0                                   | --                                               |

`--`: Need not be set

### Bank select and corresponding media (Patches)

| MSB | LSB | Program Change | Media (Patch Number)                           |
| --- | --- | -------------- | ---------------------------------------------- |
| 80  | 0   | 0 - 127        | User (#1 - #128)                               |
| 81  | 0   | 0 - 127        | Preset A (#1 - #128)                           |
| 81  | 1   | 0 - 127        | Preset B (#1 - #128)                           |
| 81  | 2   | 0 - 127        | Preset C (#1 - #128)                           |
| 81  | 3   | 0 - 127        | Preset D (General MIDI Instrument) (#1 - #128) |
| 82  | 0   | 0 - 127        | Data Card (#1 - #128)                          |
| 83  | 0   | 0 - 127        | PCM Card (#1 - #128)                           |
| 84  | 0   | 0 - 127        | Expansion A (#1 - #128)                        |
| 84  | 1   | 0 - 127        | Expansion A (#129 - #256)                      |
| 84  | 2   | 0 - 127        | Expansion B (#1 - #128)                        |
| 84  | 3   | 0 - 127        | Expansion B (#129 - #256)                      |
| 84  | 4   | 0 - 127        | Expansion C (#1 - #128)                        |
| 84  | 5   | 0 - 127        | Expansion C (#129 - #256)                      |
| 84  | 6   | 0 - 127        | Expansion D (#1 - #128)                        |
| 84  | 7   | 0 - 127        | Expansion D (#129 - #256)                      |

### Bank select and corresponding media (Rhythm Sets) (13)

| MSB | LSB | Program Change | Media (Rhythm Set Number)                    |
| --- | --- | -------------- | -------------------------------------------- |
| 80  | 0   | 0 - 1          | User (#1 - #2)                               |
| 81  | 0   | 0 - 1          | Preset A (#1 - #2)                           |
| 81  | 1   | 0 - 1          | Preset B (#1 - #2)                           |
| 81  | 2   | 0 - 1          | Preset C (#1 - #2)                           |
| 81  | 3   | 0 - 1          | Preset D (General MIDI Instrument) (#1 - #2) |
| 82  | 0   | 0 - 1          | Data Card (#1 - #2)                          |
| 83  | 0   | 0 - 127        | PCM Card (#1 - #128)                         |
| 84  | 0   | 0 - 127        | Expansion A (#1 - #128)                      |
| 84  | 1   | 0 - 127        | Expansion A (#129 - #256)                    |
| 84  | 2   | 0 - 127        | Expansion B (#1 - #128)                      |
| 84  | 3   | 0 - 127        | Expansion B (#129 - #256)                    |
| 84  | 4   | 0 - 127        | Expansion C (#1 - #128)                      |
| 84  | 5   | 0 - 127        | Expansion C (#129 - #256)                    |
| 84  | 6   | 0 - 127        | Expansion D (#1 - #128)                      |
| 84  | 7   | 0 - 127        | Expansion D (#129 - #256)                    |

---

## Seite 134 — Bank Select (Performance) · Modulation · Breath

### 2. Performance Mode

| Parameters to be set         | Changes Banks upon receiving the next program change | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |                 |
| ---------------------------- | ---------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ | --------------- |
|                              |                                                      |                                                | Setting1                                         | Setting2        |
| SYSTEM PARAMETER             |                                                      |                                                |                                                  |                 |
| System Control Source 1/2    | --                                                   | CC0:BANK SELECT                                | CC0:BANK SELECT                                  | CC0:BANK SELECT |
| Receive Bank Select          | ON                                                   | ON                                             | ON                                               | ON              |
| PERFORMANCE COMMON PARAMETER |                                                      |                                                |                                                  |                 |
| EFX Source                   | --                                                   | --                                             | 1 - 16                                           | PERFORM         |
| EFX Control Source 1/2       | --                                                   | --                                             | --                                               | SYS-CTRL1/2     |
| EFX Control Depth 1/2        | --                                                   | --                                             | --                                               | Other than 0    |
| PERFORMANCE PART PARAMETER   |                                                      |                                                |                                                  |                 |
| MIDI Receive Switch          | ON                                                   | ON                                             | ON                                               | ON              |
| PATCH COMMON PARAMETER       |                                                      |                                                |                                                  |                 |
| EFX Control Source 1/2       | --                                                   | --                                             | SYS-CTRL1/2                                      | --              |
| EFX Control Depth 1/2        | --                                                   | --                                             | Other than 0                                     | --              |
| Patch Control Source 1/2/3   | --                                                   | SYS-CTRL1/2                                    | --                                               | --              |
| PATCH TONE PARAMETER         |                                                      |                                                |                                                  |                 |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                   | Other than OFF                                 | --                                               | --              |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                   | Other than 0                                   | --                                               | --              |

Bank select and corresponding media are as shown below.

| MSB | LSB | Program Change | Media (Performance Number) |
| --- | --- | -------------- | -------------------------- |
| 80  | 0   | 0 - 31         | User                       |
| 81  | 0   | 0 - 31         | Preset A                   |
| 81  | 1   | 0 - 31         | Preset B                   |
| 82  | 0   | 0 - 31         | Data Card                  |

### ○ Modulation

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 01H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Modulation depth: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set       | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |                |
| -------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | -------------- |
|                            | Setting1                                       | Setting2       | Setting1                                         | Setting2       |
| SYSTEM PARAMETER           |                                                |                |                                                  |                |
| System Control Source 1/2  | --                                             | CC1:MODULATION | --                                               | CC1:MODULATION |
| Receive Modulation         | ON                                             | ON             | ON                                               | ON             |
| PATCH COMMON PARAMETER     |                                                |                |                                                  |                |
| EFX Control Source 1/2     | --                                             | --             | MODULATION                                       | SYS-CTRL1/2    |
| EFX Control Depth 1/2      | --                                             | --             | Other than 0                                     | Other than 0   |
| Patch Control Source 1/2/3 | MODULATION                                     | SYS-CTRL1/2    | --                                               | --             |
| PATCH TONE PARAMETER       |                                                |                |                                                  |                |
| Ctrl 1/2/3 Dest.1/2/3/4    | Other than OFF                                 | Other than OFF | --                                               | --             |
| Ctrl 1/2/3 Depth1/2/3/4    | Other than 0                                   | Other than 0   | --                                               | --             |

### 2. Performance Mode

| Parameters to be set         | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |                |                |
| ---------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ | -------------- | -------------- |
|                              | Setting1                                       | Setting2       | Setting1                                         | Setting2     | Setting3       | Setting4       |
| SYSTEM PARAMETER             |                                                |                |                                                  |              |                |                |
| System Control Source 1/2    | --                                             | CC1:MODULATION | --                                               | --           | CC1:MODULATION | CC1:MODULATION |
| Receive Modulation           | ON                                             | ON             | ON                                               | ON           | ON             | ON             |
| PERFORMANCE COMMON PARAMETER |                                                |                |                                                  |              |                |                |
| EFX Source                   | --                                             | --             | 1 - 16                                           | PERFORM      | 1 - 16         | PERFORM        |
| EFX Control Source 1/2       | --                                             | --             | --                                               | MODULATION   | --             | SYS-CTRL1/2    |
| EFX Control Depth 1/2        | --                                             | --             | --                                               | Other than 0 | --             | Other than 0   |
| PERFORMANCE PART PARAMETER   |                                                |                |                                                  |              |                |                |
| MIDI Receive Switch          | ON                                             | ON             | ON                                               | ON           | ON             | ON             |
| PATCH COMMON PARAMETER       |                                                |                |                                                  |              |                |                |
| EFX Control Source 1/2       | --                                             | --             | MODULATION                                       | --           | SYS-CTRL1/2    | --             |
| EFX Control Depth 1/2        | --                                             | --             | Other than 0                                     | --           | Other than 0   | --             |
| Patch Control Source 1/2/3   | MODULATION                                     | SYS-CTRL1/2    | --                                               | --           | --             | --             |
| PATCH TONE PARAMETER         |                                                |                |                                                  |              |                |                |
| Ctrl 1/2/3 Dest.1/2/3/4      | Other than OFF                                 | Other than OFF | --                                               | --           | --             | --             |
| Ctrl 1/2/3 Depth1/2/3/4      | Other than 0                                   | Other than 0   | --                                               | --           | --             | --             |

`--`: Need not be set

### ● Breath

### ○ Breath

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 02H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Breath: 00H - 7FH (0 - 127)

---
## Seite 135 — Breath (Patch/Performance) · Foot type

### 1. Patch Mode (Breath)

| Parameters to be set       | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |
| -------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ |
|                            | Setting1                                       | Setting2       | Setting1                                         | Setting2     |
| SYSTEM PARAMETER           |                                                |                |                                                  |              |
| System Control Source 1/2  | --                                             | CC2:BREATH     | --                                               | CC2:BREATH   |
| Receive Control Change     | ON                                             | ON             | ON                                               | ON           |
| PATCH COMMON PARAMETER     |                                                |                |                                                  |              |
| EFX Control Source 1/2     | --                                             | --             | BREATH                                           | SYS-CTRL1/2  |
| EFX Control Depth 1/2      | --                                             | --             | Other than 0                                     | Other than 0 |
| Patch Control Source 1/2/3 | BREATH                                         | SYS-CTRL1/2    | --                                               | --           |
| PATCH TONE PARAMETER       |                                                |                |                                                  |              |
| Ctrl 1/2/3 Dest.1/2/3/4    | Other than OFF                                 | Other than OFF | --                                               | --           |
| Ctrl 1/2/3 Depth1/2/3/4    | Other than 0                                   | Other than 0   | --                                               | --           |

`--`: Need not be set

### 2. Performance Mode (Breath)

| Parameters to be set         | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |              |              |
| ---------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ | ------------ | ------------ |
|                              | Setting1                                       | Setting2       | Setting1                                         | Setting2     | Setting3     | Setting4     |
| SYSTEM PARAMETER             |                                                |                |                                                  |              |              |              |
| System Control Source 1/2    | --                                             | CC2:BREATH     | --                                               | --           | CC2:BREATH   | CC2:BREATH   |
| Receive Control Change       | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PERFORMANCE COMMON PARAMETER |                                                |                |                                                  |              |              |              |
| EFX Source                   | --                                             | --             | 1 - 16                                           | PERFORM      | 1 - 16       | PERFORM      |
| EFX Control Source 1/2       | --                                             | --             | --                                               | BREATH       | --           | SYS-CTRL1/2  |
| EFX Control Depth 1/2        | --                                             | --             | --                                               | Other than 0 | --           | Other than 0 |
| PERFORMANCE PART PARAMETER   |                                                |                |                                                  |              |              |              |
| MIDI Receive Switch          | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PATCH COMMON PARAMETER       |                                                |                |                                                  |              |              |              |
| EFX Control Source 1/2       | --                                             | --             | BREATH                                           | --           | SYS-CTRL1/2  | --           |
| EFX Control Depth 1/2        | --                                             | --             | Other than 0                                     | --           | Other than 0 | --           |
| Patch Control Source 1/2/3   | BREATH                                         | SYS-CTRL1/2    | --                                               | --           | --           | --           |
| PATCH TONE PARAMETER         |                                                |                |                                                  |              |              |              |
| Ctrl 1/2/3 Dest.1/2/3/4      | Other than OFF                                 | Other than OFF | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Depth1/2/3/4      | Other than 0                                   | Other than 0   | --                                               | --           | --           | --           |

`--`: Need not be set

### ○ Foot type

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 04H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Foot control: 00H - 7FH (0 - 127)

### 1. Patch Mode (Foot type)

| Parameters to be set       | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |               |
| -------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------- |
|                            | Setting1                                       | Setting2       | Setting1                                         | Setting2      |
| SYSTEM PARAMETER           |                                                |                |                                                  |               |
| System Control Source 1/2  | --                                             | CC4:FOOT-TYPE  | --                                               | CC4:FOOT-TYPE |
| Receive Control Change     | ON                                             | ON             | ON                                               | ON            |
| PATCH COMMON PARAMETER     |                                                |                |                                                  |               |
| EFX Control Source 1/2     | --                                             | --             | FOOT                                             | SYS-CTRL1/2   |
| EFX Control Depth 1/2      | --                                             | --             | Other than 0                                     | Other than 0  |
| Patch Control Source 1/2/3 | FOOT                                           | SYS-CTRL1/2    | --                                               | --            |
| PATCH TONE PARAMETER       |                                                |                |                                                  |               |
| Ctrl 1/2/3 Dest.1/2/3/4    | Other than OFF                                 | Other than OFF | --                                               | --            |
| Ctrl 1/2/3 Depth1/2/3/4    | Other than 0                                   | Other than 0   | --                                               | --            |

`--`: Need not be set

### 2. Performance Mode (Foot type)

| Parameters to be set         | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |               |               |
| ---------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ | ------------- | ------------- |
|                              | Setting1                                       | Setting2       | Setting1                                         | Setting2     | Setting3      | Setting4      |
| SYSTEM PARAMETER             |                                                |                |                                                  |              |               |               |
| System Control Source 1/2    | --                                             | CC4:FOOT-TYPE  | --                                               | --           | CC4:FOOT-TYPE | CC4:FOOT-TYPE |
| Receive Control Change       | ON                                             | ON             | ON                                               | ON           | ON            | ON            |
| PERFORMANCE COMMON PARAMETER |                                                |                |                                                  |              |               |               |
| EFX Source                   | --                                             | --             | 1 - 16                                           | PERFORM      | 1 - 16        | PERFORM       |
| EFX Control Source 1/2       | --                                             | --             | --                                               | FOOT         | --            | SYS-CTRL1/2   |
| EFX Control Depth 1/2        | --                                             | --             | --                                               | Other than 0 | --            | Other than 0  |
| PERFORMANCE PART PARAMETER   |                                                |                |                                                  |              |               |               |
| MIDI Receive Switch          | ON                                             | ON             | ON                                               | ON           | ON            | ON            |
| PATCH COMMON PARAMETER       |                                                |                |                                                  |              |               |               |
| EFX Control Source 1/2       | --                                             | --             | FOOT                                             | --           | SYS-CTRL1/2   | --            |
| EFX Control Depth 1/2        | --                                             | --             | Other than 0                                     | --           | Other than 0  | --            |
| Patch Control Source 1/2/3   | FOOT                                           | SYS-CTRL1/2    | --                                               | --           | --            | --            |
| PATCH TONE PARAMETER         |                                                |                |                                                  |              |               |               |
| Ctrl 1/2/3 Dest.1/2/3/4      | Other than OFF                                 | Other than OFF | --                                               | --           | --            | --            |
| Ctrl 1/2/3 Depth1/2/3/4      | Other than 0                                   | Other than 0   | --                                               | --           | --            | --            |

`--`: Need not be set

---

## Seite 136 — Portamento Time · Volume

### ○ Portamento time

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 05H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Portamento time: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set       | Portamento time of patch common parameter will be changed  | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |
| -------------------------- | ---------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ |
| SYSTEM PARAMETER           |                                                            |                                                |                                                  |
| System Control Source 1/2  | --                                                         | CC5:PORT-TIME                                  | CC5:PORT-TIME                                    |
| Receive Control Change     | ON                                                         | ON                                             | ON                                               |
| PATCH COMMON PARAMETER     |                                                            |                                                |                                                  |
| EFX Control Source 1/2     | --                                                         | --                                             | SYS-CTRL1/2                                      |
| EFX Control Depth 1/2      | --                                                         | --                                             | Other than 0                                     |
| Patch Control Source 1/2/3 | --                                                         | SYS-CTRL1/2                                    | --                                               |
| PATCH TONE PARAMETER       |                                                            |                                                |                                                  |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                                                         | Other than OFF                                 | --                                               |
| Ctrl 1/2/3 Depth1/2/3/4    | --                                                         | Other than 0                                   | --                                               |

`--`: Need not be set

### 2. Performance Mode

| Parameters to be set         | Portamento time of patch common parameter of the part on the receiving channel will be changed   | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |               |
| ---------------------------- | ------------------------------------------------------------------------------------------------ | ---------------------------------------------- | ------------------------------------------------ | ------------- |
|                              |                                                                                                  |                                                | Setting1                                         | Setting2      |
| SYSTEM PARAMETER             |                                                                                                  |                                                |                                                  |               |
| System Control Source 1/2    | --                                                                                               | CC5:PORT-TIME                                  | CC5:PORT-TIME                                    | CC5:PORT-TIME |
| Receive Control Change       | ON                                                                                               | ON                                             | ON                                               | ON            |
| PERFORMANCE COMMON PARAMETER |                                                                                                  |                                                |                                                  |               |
| EFX Source                   | --                                                                                               | --                                             | 1 - 16                                           | PERFORM       |
| EFX Control Source 1/2       | --                                                                                               | --                                             | --                                               | SYS-CTRL1/2   |
| EFX Control Depth 1/2        | --                                                                                               | --                                             | --                                               | Other than 0  |
| PERFORMANCE PART PARAMETER   |                                                                                                  |                                                |                                                  |               |
| MIDI Receive Switch          | ON                                                                                               | ON                                             | ON                                               | ON            |
| PATCH COMMON PARAMETER       |                                                                                                  |                                                |                                                  |               |
| EFX Control Source 1/2       | --                                                                                               | --                                             | SYS-CTRL1/2                                      | --            |
| EFX Control Depth 1/2        | --                                                                                               | --                                             | Other than 0                                     | --            |
| Patch Control Source 1/2/3   | --                                                                                               | SYS-CTRL1/2                                    | --                                               | --            |
| PATCH TONE PARAMETER         |                                                                                                  |                                                |                                                  |               |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                                                               | Other than OFF                                 | --                                               | --            |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                                                               | Other than 0                                   | --                                               | --            |

`--`: Need not be set

### ○ Volume

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 07H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Volume: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set       | Level of tone will be changed | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |
| -------------------------- | ----------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ |
|                            |                               | Setting1                                       | Setting2       | Setting1                                         | Setting2     |
| SYSTEM PARAMETER           |                               |                                                |                |                                                  |              |
| System Control Source 1/2  | --                            | --                                             | CC7:VOLUME     | --                                               | CC7:VOLUME   |
| Receive Volume             | ON                            | ON                                             | ON             | ON                                               | ON           |
| PATCH COMMON PARAMETER     |                               |                                                |                |                                                  |              |
| EFX Control Source 1/2     | --                            | --                                             | --             | VOLUME                                           | SYS-CTRL1/2  |
| EFX Control Depth 1/2      | --                            | --                                             | --             | Other than 0                                     | Other than 0 |
| Patch Control Source 1/2/3 | --                            | VOLUME                                         | SYS-CTRL1/2    | --                                               | --           |
| PATCH TONE PARAMETER       |                               |                                                |                |                                                  |              |
| Volume Control Switch      | ON                            | --                                             | --             | --                                               | --           |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                            | Other than OFF                                 | Other than OFF | --                                               | --           |
| Ctrl 1/2/3 Depth1/2/3/4    | --                            | Other than 0                                   | Other than 0   | --                                               | --           |

`--`: Need not be set

---

## Seite 137 — Volume (Performance) · Pan

### 2. Performance Mode (Volume)

| Parameters to be set         | Tone level of the part on the receiving channel will be changed | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |              |              |
| ---------------------------- | --------------------------------------------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ | ------------ | ------------ |
|                              |                                                                 | Setting1                                       | Setting2       | Setting1                                         | Setting2     | Setting3     | Setting4     |
| SYSTEM PARAMETER             |                                                                 |                                                |                |                                                  |              |              |              |
| System Control Source 1/2    | --                                                              | --                                             | CC7:VOLUME     | --                                               | --           | CC7:VOLUME   | CC7:VOLUME   |
| Receive Volume               | ON                                                              | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PERFORMANCE COMMON PARAMETER |                                                                 |                                                |                |                                                  |              |              |              |
| EFX Source                   | --                                                              | --                                             | --             | 1 - 16                                           | PERFORM      | 1 - 16       | PERFORM      |
| EFX Control Source 1/2       | --                                                              | --                                             | --             | --                                               | VOLUME       | --           | SYS-CTRL1/2  |
| EFX Control Depth 1/2        | --                                                              | --                                             | --             | --                                               | Other than 0 | --           | Other than 0 |
| PERFORMANCE PART PARAMETER   |                                                                 |                                                |                |                                                  |              |              |              |
| MIDI Receive Switch          | ON                                                              | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| Receive Volume               | ON                                                              | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PATCH COMMON PARAMETER       |                                                                 |                                                |                |                                                  |              |              |              |
| EFX Control Source 1/2       | --                                                              | --                                             | --             | VOLUME                                           | --           | SYS-CTRL1/2  | --           |
| EFX Control Depth 1/2        | --                                                              | --                                             | --             | Other than 0                                     | --           | Other than 0 | --           |
| Patch Control Source 1/2/3   | --                                                              | VOLUME                                         | SYS-CTRL1/2    | --                                               | --           | --           | --           |
| PATCH TONE PARAMETER         |                                                                 |                                                |                |                                                  |              |              |              |
| Volume Control Switch        | ON                                                              | --                                             | --             | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                              | Other than OFF                                 | Other than OFF | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                              | Other than 0                                   | Other than 0   | --                                               | --           | --           | --           |

`--`: Need not be set

### ○ Pan

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 0AH    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Pan: 00H - 7FH (0 - 127)

### 1. Patch Mode (Pan)

| Parameters to be set       | Directional localization of tone can be controlled in 128 steps with 0 the leftmost, 64 at the center, 127 the rightmost | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ |
|                            |                                                                                                                          | Setting1                                       | Setting2       | Setting1                                         | Setting2     |
| SYSTEM PARAMETER           |                                                                                                                          |                                                |                |                                                  |              |
| System Control Source 1/2  | --                                                                                                                       | --                                             | CC10:PANPOT    | --                                               | CC10:PANPOT  |
| Receive Control Change     | ON                                                                                                                       | ON                                             | ON             | ON                                               | ON           |
| PATCH COMMON PARAMETER     |                                                                                                                          |                                                |                |                                                  |              |
| EFX Control Source 1/2     | --                                                                                                                       | --                                             | --             | PAN                                              | SYS-CTRL1/2  |
| EFX Control Depth 1/2      | --                                                                                                                       | --                                             | --             | Other than 0                                     | Other than 0 |
| Patch Control Source 1/2/3 | --                                                                                                                       | PAN                                            | SYS-CTRL1/2    | --                                               | --           |
| PATCH TONE PARAMETER       |                                                                                                                          |                                                |                |                                                  |              |
| Pan Control Switch         | CONTINUOUS or KEY-ON                                                                                                     | --                                             | --             | --                                               | --           |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                                                                                                                       | Other than OFF                                 | Other than OFF | --                                               | --           |
| Ctrl 1/2/3 Depth1/2/3/4    | --                                                                                                                       | Other than 0                                   | Other than 0   | --                                               | --           |

### 2. Performance Mode (Pan)

| Parameters to be set         | Directional localization of tone of the part on the receiving channel can be controlled in 128 steps with 0 the leftmost, 64 at the center, 127 the rightmost | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |              |              |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ | ------------ | ------------ |
|                              |                                                                                                                                                               | Setting1                                       | Setting2       | Setting1                                         | Setting2     | Setting3     | Setting4     |
| SYSTEM PARAMETER             |                                                                                                                                                               |                                                |                |                                                  |              |              |              |
| System Control Source 1/2    | --                                                                                                                                                            | --                                             | CC10:PANPOT    | --                                               | --           | CC10:PANPOT  | CC10:PANPOT  |
| Receive Control Change       | ON                                                                                                                                                            | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PERFORMANCE COMMON PARAMETER |                                                                                                                                                               |                                                |                |                                                  |              |              |              |
| EFX Source                   | --                                                                                                                                                            | --                                             | --             | 1 - 16                                           | PERFORM      | 1 - 16       | PERFORM      |
| EFX Control Source 1/2       | --                                                                                                                                                            | --                                             | --             | --                                               | PAN          | --           | SYS-CTRL1/2  |
| EFX Control Depth 1/2        | --                                                                                                                                                            | --                                             | --             | --                                               | Other than 0 | --           | Other than 0 |
| PERFORMANCE PART PARAMETER   |                                                                                                                                                               |                                                |                |                                                  |              |              |              |
| MIDI Receive Switch          | ON                                                                                                                                                            | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PATCH COMMON PARAMETER       |                                                                                                                                                               |                                                |                |                                                  |              |              |              |
| EFX Control Source 1/2       | --                                                                                                                                                            | --                                             | --             | PAN                                              | --           | SYS-CTRL1/2  | --           |
| EFX Control Depth 1/2        | --                                                                                                                                                            | --                                             | --             | Other than 0                                     | --           | Other than 0 | --           |
| Patch Control Source 1/2/3   | --                                                                                                                                                            | PAN                                            | SYS-CTRL1/2    | --                                               | --           | --           | --           |
| PATCH TONE PARAMETER         |                                                                                                                                                               |                                                |                |                                                  |              |              |              |
| Pan Control Switch           | CONTINUOUS or KEY-ON                                                                                                                                          | --                                             | --             | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                                                                                                                            | Other than OFF                                 | Other than OFF | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                                                                                                                            | Other than 0                                   | Other than 0   | --                                               | --           | --           | --           |

`--`: Need not be set

---

## Seite 138 — Expression · Hold 1 (Patch Mode)

### ○ Expression

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 0BH    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Expression: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set       | Level of tone will be changed | Tone parameters can be controlled in real-time |                 | Effect parameters can be controlled in real-time |                 |
| -------------------------- | ----------------------------- | ---------------------------------------------- | --------------- | ------------------------------------------------ | --------------- |
|                            |                               | Setting1                                       | Setting2        | Setting1                                         | Setting2        |
| SYSTEM PARAMETER           |                               |                                                |                 |                                                  |                 |
| System Control Source 1/2  | --                            | --                                             | CC11:EXPRESSION | --                                               | CC11:EXPRESSION |
| Receive Control Change     | ON                            | ON                                             | ON              | ON                                               | ON              |
| Volume Control Source      | VOL&EXP                       | --                                             | --              | --                                               | --              |
| PATCH COMMON PARAMETER     |                               |                                                |                 |                                                  |                 |
| EFX Control Source 1/2     | --                            | --                                             | --              | EXPRESSION                                       | SYS-CTRL1/2     |
| EFX Control Depth 1/2      | --                            | --                                             | --              | Other than 0                                     | Other than 0    |
| Patch Control Source 1/2/3 | --                            | EXPRESSION                                     | SYS-CTRL1/2     | --                                               | --              |
| PATCH TONE PARAMETER       |                               |                                                |                 |                                                  |                 |
| Volume Control Switch      | ON                            | --                                             | --              | --                                               | --              |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                            | Other than OFF                                 | Other than OFF  | --                                               | --              |
| Ctrl 1/2/3 Depth1/2/3/4    | --                            | Other than 0                                   | Other than 0    | --                                               | --              |

`--`: Need not be set

### 2. Performance Mode

| Parameters to be set         | Level of tone of the part on the receiving channel will be changed | Tone parameters can be controlled in real-time |                 | Effect parameters can be controlled in real-time |              |                 |                 |
| ---------------------------- | ------------------------------------------------------------------ | ---------------------------------------------- | --------------- | ------------------------------------------------ | ------------ | --------------- | --------------- |
|                              |                                                                    | Setting1                                       | Setting2        | Setting1                                         | Setting2     | Setting3        | Setting4        |
| SYSTEM PARAMETER             |                                                                    |                                                |                 |                                                  |              |                 |                 |
| System Control Source 1/2    | --                                                                 | --                                             | CC11:EXPRESSION | --                                               | --           | CC11:EXPRESSION | CC11:EXPRESSION |
| Receive Control Change       | ON                                                                 | ON                                             | ON              | ON                                               | ON           | ON              | ON              |
| Volume Control Source        | VOL&EXP                                                            | --                                             | --              | --                                               | --           | --              | --              |
| PERFORMANCE COMMON PARAMETER |                                                                    |                                                |                 |                                                  |              |                 |                 |
| EFX Source                   | --                                                                 | --                                             | --              | 1 - 16                                           | PERFORM      | 1 - 16          | PERFORM         |
| EFX Control Source 1/2       | --                                                                 | --                                             | --              | --                                               | EXPRESSION   | --              | SYS-CTRL1/2     |
| EFX Control Depth 1/2        | --                                                                 | --                                             | --              | --                                               | Other than 0 | --              | Other than 0    |
| PERFORMANCE PART PARAMETER   |                                                                    |                                                |                 |                                                  |              |                 |                 |
| MIDI Receive Switch          | ON                                                                 | ON                                             | ON              | ON                                               | ON           | ON              | ON              |
| PATCH COMMON PARAMETER       |                                                                    |                                                |                 |                                                  |              |                 |                 |
| EFX Control Source 1/2       | --                                                                 | --                                             | --              | EXPRESSION                                       | --           | SYS-CTRL1/2     | --              |
| EFX Control Depth 1/2        | --                                                                 | --                                             | --              | Other than 0                                     | --           | Other than 0    | --              |
| Patch Control Source 1/2/3   | --                                                                 | EXPRESSION                                     | SYS-CTRL1/2     | --                                               | --           | --              | --              |
| PATCH TONE PARAMETER         |                                                                    |                                                |                 |                                                  |              |                 |                 |
| Volume Control Switch        | ON                                                                 | --                                             | --              | --                                               | --           | --              | --              |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                                 | Other than OFF                                 | Other than OFF  | --                                               | --           | --              | --              |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                                 | Other than 0                                   | Other than 0    | --                                               | --           | --              | --              |

`--`: Need not be set

### ○ Hold 1

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 40H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Control value: 00H - 7FH (0 - 127) 0 - 63 = OFF 64 - 127 = ON

### 1. Patch Mode (Hold 1) (14)

| Parameters to be set          | When HOLD-1 is ON, note is kept on | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time | LFO rate is changed to the HOLD-1 ON interval  | Tone delay time is changed to the HOLD-1 ON interval | When HOLD-1 is ON, Effect controller is held | When HOLD-1 is ON, Effect controller is held at its peak level | When HOLD-1 is ON, Tone controller is held | When HOLD-1 is ON, Tone controller is held at its peak level |
| ----------------------------- | ---------------------------------- | ---------------------------------------------- | ------------------------------------------------ | ---------------------------------------------- | ---------------------------------------------------- | -------------------------------------------- | -------------------------------------------------------------- | ------------------------------------------ | ------------------------------------------------------------ |
| SYSTEM PARAMETER              |                                    |                                                |                                                  |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| System Control Source 1/2     | --                                 | CC64:HOLD-1                                    | CC64:HOLD-1                                      | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Receive Hold-1                | ON                                 | ON                                             | ON                                               | ON                                             | ON                                                   | ON                                           | ON                                                             | ON                                         | ON                                                           |
| TAP Control Source            | --                                 | --                                             | --                                               | HOLD-1                                         | HOLD-1                                               | --                                           | --                                                             | --                                         | --                                                           |
| Hold Control Source           | --                                 | --                                             | --                                               | --                                             | --                                                   | HOLD-1                                       | --                                                             | HOLD-1                                     | --                                                           |
| Peak Control Source           | --                                 | --                                             | --                                               | --                                             | --                                                   | --                                           | HOLD-1                                                         | --                                         | HOLD-1                                                       |
| PATCH COMMON PARAMETER        |                                    |                                                |                                                  |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| EFX Control Source 1/2        | --                                 | --                                             | SYS-CTRL1/2                                      | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| EFX Control Depth 1/2         | --                                 | --                                             | Other than 0                                     | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Patch Control Source 1/2/3    | --                                 | SYS-CTRL1/2                                    | --                                               | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| EFX Control Hold/Peak         | --                                 | --                                             | --                                               | --                                             | --                                                   | HOLD                                         | PEAK                                                           | --                                         | --                                                           |
| Control 1/2/3 Hold/Peak       | --                                 | --                                             | --                                               | --                                             | --                                                   | --                                           | --                                                             | HOLD                                       | PEAK                                                         |
| PATCH TONE PARAMETER          |                                    |                                                |                                                  |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| Hold-1 Control Switch         | ON                                 | --                                             | --                                               | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Ctrl 1/2/3 Dest.1/2/3/4       | --                                 | Other than OFF                                 | --                                               | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Ctrl 1/2/3 Depth1/2/3/4       | --                                 | Other than 0                                   | --                                               | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| LFO1/2 External Sync          | --                                 | --                                             | --                                               | TAP                                            | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Tone Delay Mode               | --                                 | --                                             | --                                               | --                                             | TAP-SYNC                                             | --                                           | --                                                             | --                                         | --                                                           |

`--`: Need not be set

---

## Seite 139 — Hold 1 (Performance Mode) · Portamento · Sostenuto

### 2. Performance Mode (Hold 1) (14)

| Parameters to be set          | When HOLD-1 is ON, note of the part on the receiving channel is held | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |              | LFO rate is changed to the HOLD-1 ON interval  | Tone delay time is changed to the HOLD-1 ON interval | When HOLD-1 is ON, Effect controller is held | When HOLD-1 is ON, Effect controller is held at its peak level | When HOLD-1 is ON, Tone controller is held | When HOLD-1 is ON, Tone controller is held at its peak level |
| ----------------------------- | -------------------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ | ------------ | ---------------------------------------------- | ---------------------------------------------------- | -------------------------------------------- | -------------------------------------------------------------- | ------------------------------------------ | ------------------------------------------------------------ |
|                               |                                                                      |                                                | Setting1                                         | Setting2     |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| SYSTEM PARAMETER              |                                                                      |                                                |                                                  |              |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| System Control Source 1/2     | --                                                                   | CC64:HOLD-1                                    | CC64:HOLD-1                                      | CC64:HOLD-1  | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Receive Hold-1                | ON                                                                   | ON                                             | ON                                               | ON           | ON                                             | ON                                                   | ON                                           | ON                                                             | ON                                         | ON                                                           |
| TAP Control Source            | --                                                                   | --                                             | --                                               | --           | HOLD-1                                         | HOLD-1                                               | --                                           | --                                                             | --                                         | --                                                           |
| Hold Control Source           | --                                                                   | --                                             | --                                               | --           | --                                             | --                                                   | HOLD-1                                       | --                                                             | HOLD-1                                     | --                                                           |
| Peak Control Source           | --                                                                   | --                                             | --                                               | --           | --                                             | --                                                   | --                                           | HOLD-1                                                         | --                                         | HOLD-1                                                       |
| PERFORMANCE COMMON PARAMETER  |                                                                      |                                                |                                                  |              |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| EFX Source                    | --                                                                   | --                                             | 1 - 16                                           | PERFORM      | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| EFX Control Source 1/2        | --                                                                   | --                                             | --                                               | SYS-CTRL1/2  | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| EFX Control Depth 1/2         | --                                                                   | --                                             | --                                               | Other than 0 | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| PERFORMANCE PART PARAMETER    |                                                                      |                                                |                                                  |              |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| MIDI Receive Switch           | ON                                                                   | ON                                             | ON                                               | ON           | ON                                             | ON                                                   | ON                                           | ON                                                             | ON                                         | ON                                                           |
| PATCH COMMON PARAMETER        |                                                                      |                                                |                                                  |              |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| EFX Control Source 1/2        | --                                                                   | --                                             | SYS-CTRL1/2                                      | --           | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| EFX Control Depth 1/2         | --                                                                   | --                                             | Other than 0                                     | --           | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Patch Control Source 1/2/3    | --                                                                   | SYS-CTRL1/2                                    | --                                               | --           | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| EFX Control Hold/Peak         | --                                                                   | --                                             | --                                               | --           | --                                             | --                                                   | HOLD                                         | PEAK                                                           | --                                         | --                                                           |
| Control 1/2/3 Hold/Peak       | --                                                                   | --                                             | --                                               | --           | --                                             | --                                                   | --                                           | --                                                             | HOLD                                       | PEAK                                                         |
| PATCH TONE PARAMETER          |                                                                      |                                                |                                                  |              |                                                |                                                      |                                              |                                                                |                                            |                                                              |
| Hold-1 Control Switch         | ON                                                                   | --                                             | --                                               | --           | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Ctrl 1/2/3 Dest.1/2/3/4       | --                                                                   | Other than OFF                                 | --                                               | --           | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Ctrl 1/2/3 Depth1/2/3/4       | --                                                                   | Other than 0                                   | --                                               | --           | --                                             | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| LFO1/2 External Sync          | --                                                                   | --                                             | --                                               | --           | TAP                                            | --                                                   | --                                           | --                                                             | --                                         | --                                                           |
| Tone Delay Mode               | --                                                                   | --                                             | --                                               | --           | --                                             | TAP-SYNC                                             | --                                           | --                                                             | --                                         | --                                                           |

`--`: Need not be set

### ○ Portamento

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 41H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Control value: 00H - 7FH (0 - 127) 0 - 63 = OFF 64 - 127 = ON

### 1. Patch Mode (Portamento)

| Parameters to be set       | Portamento switch of Patch Common parameter is changed its setting | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |
| -------------------------- | ------------------------------------------------------------------ | ---------------------------------------------- | ------------------------------------------------ |
| SYSTEM PARAMETER           |                                                                    |                                                |                                                  |
| System Control Source 1/2  | --                                                                 | CC65:PORTAMENTO                                | CC65:PORTAMENTO                                  |
| Receive Control Change     | ON                                                                 | ON                                             | ON                                               |
| PATCH COMMON PARAMETER     |                                                                    |                                                |                                                  |
| EFX Control Source 1/2     | --                                                                 | --                                             | SYS-CTRL1/2                                      |
| EFX Control Depth 1/2      | --                                                                 | --                                             | Other than 0                                     |
| Patch Control Source 1/2/3 | --                                                                 | SYS-CTRL1/2                                    | --                                               |
| PATCH TONE PARAMETER       |                                                                    |                                                |                                                  |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                                                                 | Other than OFF                                 | --                                               |
| Ctrl 1/2/3 Depth1/2/3/4    | --                                                                 | Other than 0                                   | --                                               |

`--`: Need not be set

### 2. Performance Mode (Portamento)

| Parameters to be set         | Portamento switch of Patch Common parameter of the part on the receiving channel is changed its setting   | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |                 |
| ---------------------------- | --------------------------------------------------------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ | --------------- |
|                              |                                                                                                           |                                                | Setting1                                         | Setting2        |
| SYSTEM PARAMETER             |                                                                                                           |                                                |                                                  |                 |
| System Control Source 1/2    | --                                                                                                        | CC65:PORTAMENTO                                | CC65:PORTAMENTO                                  | CC65:PORTAMENTO |
| Receive Control Change       | ON                                                                                                        | ON                                             | ON                                               | ON              |
| PERFORMANCE COMMON PARAMETER |                                                                                                           |                                                |                                                  |                 |
| EFX Source                   | --                                                                                                        | --                                             | 1 - 16                                           | PERFORM         |
| EFX Control Source 1/2       | --                                                                                                        | --                                             | --                                               | SYS-CTRL1/2     |
| EFX Control Depth 1/2        | --                                                                                                        | --                                             | --                                               | Other than 0    |
| PERFORMANCE PART PARAMETER   |                                                                                                           |                                                |                                                  |                 |
| MIDI Receive Switch          | ON                                                                                                        | ON                                             | ON                                               | ON              |
| PATCH COMMON PARAMETER       |                                                                                                           |                                                |                                                  |                 |
| EFX Control Source 1/2       | --                                                                                                        | --                                             | SYS-CTRL1/2                                      | --              |
| EFX Control Depth 1/2        | --                                                                                                        | --                                             | Other than 0                                     | --              |
| Patch Control Source 1/2/3   | --                                                                                                        | SYS-CTRL1/2                                    | --                                               | --              |
| PATCH TONE PARAMETER         |                                                                                                           |                                                |                                                  |                 |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                                                                        | Other than OFF                                 | --                                               | --              |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                                                                        | Other than 0                                   | --                                               | --              |

`--`: Need not be set

### ○ Sostenuto

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 42H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Control value: 00H - 7FH (0 - 127) 0 - 63 = OFF 64 - 127 = ON

---
## Seite 140 — Sostenuto (Patch/Performance) · Soft Pedal (Message/Patch)

### 1. Patch Mode (Sostenuto) (14)

| Parameters to be set          | When SOSTENUTO is ON, note is kept on | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time | LFO rate is changed to the SOSTENUTO ON interval | Tone delay time is changed to the SOSTENUTO ON interval | When SOSTENUTO is ON, Effect controller is held | When SOSTENUTO is ON, Effect controller is held at its peak level | When SOSTENUTO is ON, Tone controller is held | When SOSTENUTO is ON, Tone controller is held at its peak level |
| ----------------------------- | ------------------------------------- | ---------------------------------------------- | ------------------------------------------------ | ------------------------------------------------ | ------------------------------------------------------- | ----------------------------------------------- | ----------------------------------------------------------------- | --------------------------------------------- | --------------------------------------------------------------- |
| SYSTEM PARAMETER              |                                       |                                                |                                                  |                                                  |                                                         |                                                 |                                                                   |                                               |                                                                 |
| System Control Source 1/2     | --                                    | CC66:SOSTENUTO                                 | CC66:SOSTENUTO                                   | --                                               | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| Receive Sostenuto             | ON                                    | ON                                             | ON                                               | ON                                               | ON                                                      | ON                                              | ON                                                                | ON                                            | ON                                                              |
| TAP Control Source            | --                                    | --                                             | --                                               | SOSTENUTO                                        | SOSTENUTO                                               | --                                              | --                                                                | --                                            | --                                                              |
| Hold Control Source           | --                                    | --                                             | --                                               | --                                               | --                                                      | SOSTENUTO                                       | --                                                                | SOSTENUTO                                     | --                                                              |
| Peak Control Source           | --                                    | --                                             | --                                               | --                                               | --                                                      | --                                              | SOSTENUTO                                                         | --                                            | SOSTENUTO                                                       |
| PATCH COMMON PARAMETER        |                                       |                                                |                                                  |                                                  |                                                         |                                                 |                                                                   |                                               |                                                                 |
| EFX Control Source 1/2        | --                                    | --                                             | SYS-CTRL1/2                                      | --                                               | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| EFX Control Depth 1/2         | --                                    | --                                             | Other than 0                                     | --                                               | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| Patch Control Source 1/2/3    | --                                    | SYS-CTRL1/2                                    | --                                               | --                                               | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| EFX Control Hold/Peak         | --                                    | --                                             | --                                               | --                                               | --                                                      | HOLD                                            | PEAK                                                              | --                                            | --                                                              |
| Control 1/2/3 Hold/Peak       | --                                    | --                                             | --                                               | --                                               | --                                                      | --                                              | --                                                                | HOLD                                          | PEAK                                                            |
| PATCH TONE PARAMETER          |                                       |                                                |                                                  |                                                  |                                                         |                                                 |                                                                   |                                               |                                                                 |
| Ctrl 1/2/3 Dest.1/2/3/4       | --                                    | Other than OFF                                 | --                                               | --                                               | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| Ctrl 1/2/3 Depth1/2/3/4       | --                                    | Other than 0                                   | --                                               | --                                               | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| LFO1/2 External Sync          | --                                    | --                                             | --                                               | TAP                                              | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| Tone Delay Mode               | --                                    | --                                             | --                                               | --                                               | TAP-SYNC                                                | --                                              | --                                                                | --                                            | --                                                              |

`--`: Need not be set

### 2. Performance Mode (Sostenuto) (14)

| Parameters to be set          | When SOSTENUTO is ON, note of the part on the receiving channel is held  | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |                | LFO rate is changed to the SOSTENUTO ON interval   | Tone delay time is changed to the SOSTENUTO ON interval | When SOSTENUTO is ON, Effect controller is held | When SOSTENUTO is ON, Effect controller is held at its peak level | When SOSTENUTO is ON, Tone controller is held | When SOSTENUTO is ON, Tone controller is held at its peak level |
| ----------------------------- | ------------------------------------------------------------------------ | ---------------------------------------------- | ------------------------------------------------ | -------------- | -------------------------------------------------- | ------------------------------------------------------- | ----------------------------------------------- | ----------------------------------------------------------------- | --------------------------------------------- | --------------------------------------------------------------- |
|                               |                                                                          |                                                | Setting1                                         | Setting2       |                                                    |                                                         |                                                 |                                                                   |                                               |                                                                 |
| SYSTEM PARAMETER              |                                                                          |                                                |                                                  |                |                                                    |                                                         |                                                 |                                                                   |                                               |                                                                 |
| System Control Source 1/2     | --                                                                       | CC66:SOSTENUTO                                 | CC66:SOSTENUTO                                   | CC66:SOSTENUTO | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| Receive Sostenuto             | ON                                                                       | ON                                             | ON                                               | ON             | ON                                                 | ON                                                      | ON                                              | ON                                                                | ON                                            | ON                                                              |
| TAP Control Source            | --                                                                       | --                                             | --                                               | --             | SOSTENUTO                                          | SOSTENUTO                                               | --                                              | --                                                                | --                                            | --                                                              |
| Hold Control Source           | --                                                                       | --                                             | --                                               | --             | --                                                 | --                                                      | SOSTENUTO                                       | --                                                                | SOSTENUTO                                     | --                                                              |
| Peak Control Source           | --                                                                       | --                                             | --                                               | --             | --                                                 | --                                                      | --                                              | SOSTENUTO                                                         | --                                            | SOSTENUTO                                                       |
| PERFORMANCE COMMON PARAMETER  |                                                                          |                                                |                                                  |                |                                                    |                                                         |                                                 |                                                                   |                                               |                                                                 |
| EFX Source                    | --                                                                       | --                                             | 1 - 16                                           | PERFORM        | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| EFX Control Source 1/2        | --                                                                       | --                                             | --                                               | SYS-CTRL1/2    | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| EFX Control Depth 1/2         | --                                                                       | --                                             | --                                               | Other than 0   | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| PERFORMANCE PART PARAMETER    |                                                                          |                                                |                                                  |                |                                                    |                                                         |                                                 |                                                                   |                                               |                                                                 |
| MIDI Receive Switch           | ON                                                                       | ON                                             | ON                                               | ON             | ON                                                 | ON                                                      | ON                                              | ON                                                                | ON                                            | ON                                                              |
| PATCH COMMON PARAMETER        |                                                                          |                                                |                                                  |                |                                                    |                                                         |                                                 |                                                                   |                                               |                                                                 |
| EFX Control Source 1/2        | --                                                                       | --                                             | SYS-CTRL1/2                                      | --             | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| EFX Control Depth 1/2         | --                                                                       | --                                             | Other than 0                                     | --             | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| Patch Control Source 1/2/3    | --                                                                       | SYS-CTRL1/2                                    | --                                               | --             | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| EFX Control Hold/Peak         | --                                                                       | --                                             | --                                               | --             | --                                                 | --                                                      | HOLD                                            | PEAK                                                              | --                                            | --                                                              |
| Control 1/2/3 Hold/Peak       | --                                                                       | --                                             | --                                               | --             | --                                                 | --                                                      | --                                              | --                                                                | HOLD                                          | PEAK                                                            |
| PATCH TONE PARAMETER          |                                                                          |                                                |                                                  |                |                                                    |                                                         |                                                 |                                                                   |                                               |                                                                 |
| Ctrl 1/2/3 Dest.1/2/3/4       | --                                                                       | Other than OFF                                 | --                                               | --             | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| Ctrl 1/2/3 Depth1/2/3/4       | --                                                                       | Other than 0                                   | --                                               | --             | --                                                 | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| LFO1/2 External Sync          | --                                                                       | --                                             | --                                               | --             | TAP                                                | --                                                      | --                                              | --                                                                | --                                            | --                                                              |
| Tone Delay Mode               | --                                                                       | --                                             | --                                               | --             | --                                                 | TAP-SYNC                                                | --                                              | --                                                                | --                                            | --                                                              |

`--`: Need not be set

### ○ Soft Pedal

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 43H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Control value: 00H - 7FH (0 - 127) 0 - 63 = OFF 64 - 127 = ON

### 1. Patch Mode (Soft Pedal) (14)

| Parameters to be set          | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time | LFO rate is changed to the SOFT ON interval  | Tone delay time is changed to the SOFT ON interval  | When SOFT is ON, Effect controller is held   | When SOFT is ON, Effect controller is held at its peak level | When SOFT is ON, Tone controller is held  | When SOFT is ON, Tone controller is held at its peak level   |
| ----------------------------- | ---------------------------------------------- | ------------------------------------------------ | -------------------------------------------- | --------------------------------------------------- | -------------------------------------------- | ------------------------------------------------------------ | ----------------------------------------- | ------------------------------------------------------------ |
| SYSTEM PARAMETER              |                                                |                                                  |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| System Control Source 1/2     | CC67:SOFT                                      | CC67:SOFT                                        | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| Receive Soft Pedal            | ON                                             | ON                                               | ON                                           | ON                                                  | ON                                           | ON                                                           | ON                                        | ON                                                           |
| TAP Control Source            | --                                             | --                                               | SOFT                                         | SOFT                                                | --                                           | --                                                           | --                                        | --                                                           |
| Hold Control Source           | --                                             | --                                               | --                                           | --                                                  | SOFT                                         | --                                                           | SOFT                                      | --                                                           |
| Peak Control Source           | --                                             | --                                               | --                                           | --                                                  | --                                           | SOFT                                                         | --                                        | SOFT                                                         |
| PATCH COMMON PARAMETER        |                                                |                                                  |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| EFX Control Source 1/2        | --                                             | SYS-CTRL1/2                                      | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| EFX Control Depth 1/2         | --                                             | Other than 0                                     | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| Patch Control Source 1/2/3    | SYS-CTRL1/2                                    | --                                               | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| EFX Control Hold/Peak         | --                                             | --                                               | --                                           | --                                                  | HOLD                                         | PEAK                                                         | --                                        | --                                                           |
| Control 1/2/3 Hold/Peak       | --                                             | --                                               | --                                           | --                                                  | --                                           | --                                                           | HOLD                                      | PEAK                                                         |
| PATCH TONE PARAMETER          |                                                |                                                  |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| Ctrl 1/2/3 Dest.1/2/3/4       | Other than OFF                                 | --                                               | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| Ctrl 1/2/3 Depth1/2/3/4       | Other than 0                                   | --                                               | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| LFO1/2 External Sync          | --                                             | --                                               | TAP                                          | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| Tone Delay Mode               | --                                             | --                                               | --                                           | TAP-SYNC                                            | --                                           | --                                                           | --                                        | --                                                           |

`--`: Need not be set

---

## Seite 141 — Soft Pedal (Performance) · Hold 2

### 2. Performance Mode (Soft Pedal) (14)

| Parameters to be set          | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |              | LFO rate is changed to the SOFT ON interval  | Tone delay time is changed to the SOFT ON interval  | When SOFT is ON, Effect controller is held   | When SOFT is ON, Effect controller is held at its peak level | When SOFT is ON, Tone controller is held  | When SOFT is ON, Tone controller is held at its peak level   |
| ----------------------------- | ---------------------------------------------- | ------------------------------------------------ | ------------ | -------------------------------------------- | --------------------------------------------------- | -------------------------------------------- | ------------------------------------------------------------ | ----------------------------------------- | ------------------------------------------------------------ |
|                               |                                                | Setting1                                         | Setting2     |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| SYSTEM PARAMETER              |                                                |                                                  |              |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| System Control Source 1/2     | CC67:SOFT                                      | CC67:SOFT                                        | CC67:SOFT    | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| Receive Soft Pedal            | ON                                             | ON                                               | ON           | ON                                           | ON                                                  | ON                                           | ON                                                           | ON                                        | ON                                                           |
| TAP Control Source            | --                                             | --                                               | --           | SOFT                                         | SOFT                                                | --                                           | --                                                           | --                                        | --                                                           |
| Hold Control Source           | --                                             | --                                               | --           | --                                           | --                                                  | SOFT                                         | --                                                           | SOFT                                      | --                                                           |
| Peak Control Source           | --                                             | --                                               | --           | --                                           | --                                                  | --                                           | SOFT                                                         | --                                        | SOFT                                                         |
| PERFORMANCE COMMON PARAMETER  |                                                |                                                  |              |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| EFX Source                    | --                                             | 1 - 16                                           | PERFORM      | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| EFX Control Source 1/2        | --                                             | --                                               | SYS-CTRL1/2  | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| EFX Control Depth 1/2         | --                                             | --                                               | Other than 0 | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| PERFORMANCE PART PARAMETER    |                                                |                                                  |              |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| MIDI Receive Switch           | ON                                             | ON                                               | ON           | ON                                           | ON                                                  | ON                                           | ON                                                           | ON                                        | ON                                                           |
| PATCH COMMON PARAMETER        |                                                |                                                  |              |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| EFX Control Source 1/2        | --                                             | SYS-CTRL1/2                                      | --           | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| EFX Control Depth 1/2         | --                                             | Other than 0                                     | --           | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| Patch Control Source 1/2/3    | SYS-CTRL1/2                                    | --                                               | --           | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| EFX Control Hold/Peak         | --                                             | --                                               | --           | --                                           | --                                                  | HOLD                                         | PEAK                                                         | --                                        | --                                                           |
| Control 1/2/3 Hold/Peak       | --                                             | --                                               | --           | --                                           | --                                                  | --                                           | --                                                           | HOLD                                      | PEAK                                                         |
| PATCH TONE PARAMETER          |                                                |                                                  |              |                                              |                                                     |                                              |                                                              |                                           |                                                              |
| Ctrl 1/2/3 Dest.1/2/3/4       | Other than OFF                                 | --                                               | --           | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| Ctrl 1/2/3 Depth1/2/3/4       | Other than 0                                   | --                                               | --           | --                                           | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| LFO1/2 External Sync          | --                                             | --                                               | --           | TAP                                          | --                                                  | --                                           | --                                                           | --                                        | --                                                           |
| Tone Delay Mode               | --                                             | --                                               | --           | --                                           | TAP-SYNC                                            | --                                           | --                                                           | --                                        | --                                                           |

`--`: Need not be set

### ○ Hold 2

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 45H    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Control value: 00H - 7FH (0 - 127) 0 - 63 = OFF 64 - 127 = ON

### 1. Patch Mode (Hold 2) (14)

| Parameters to be set          | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time | LFO rate is changed to the HOLD-2 ON interval  | Tone delay time is changed to the HOLD-2 ON interval  | When HOLD-2 is ON, Effect controller is held   | When HOLD-2 is ON, Effect controller is held at its peak level  | When HOLD-2 is ON, Tone controller is held  | When HOLD-2 is ON, Tone controller is held at its peak level   |
| ----------------------------- | ---------------------------------------------- | ------------------------------------------------ | ---------------------------------------------- | ----------------------------------------------------- | ---------------------------------------------- | --------------------------------------------------------------- | ------------------------------------------- | -------------------------------------------------------------- |
| SYSTEM PARAMETER              |                                                |                                                  |                                                |                                                       |                                                |                                                                 |                                             |                                                                |
| System Control Source 1/2     | CC69:HOLD-2                                    | CC69:HOLD-2                                      | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                             |
| Receive Hold-2                | ON                                             | ON                                               | ON                                             | ON                                                    | ON                                             | ON                                                              | ON                                          | ON                                                             |
| TAP Control Source            | --                                             | --                                               | HOLD-2                                         | HOLD-2                                                | --                                             | --                                                              | --                                          | --                                                             |
| Hold Control Source           | --                                             | --                                               | --                                             | --                                                    | HOLD-2                                         | --                                                              | HOLD-2                                      | --                                                             |
| Peak Control Source           | --                                             | --                                               | --                                             | --                                                    | --                                             | HOLD-2                                                          | --                                          | HOLD-2                                                         |
| PATCH COMMON PARAMETER        |                                                |                                                  |                                                |                                                       |                                                |                                                                 |                                             |                                                                |
| EFX Control Source 1/2        | --                                             | SYS-CTRL1/2                                      | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                             |
| EFX Control Depth 1/2         | --                                             | Other than 0                                     | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                             |
| Patch Control Source 1/2/3    | SYS-CTRL1/2                                    | --                                               | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                             |
| EFX Control Hold/Peak         | --                                             | --                                               | --                                             | --                                                    | HOLD                                           | PEAK                                                            | --                                          | --                                                             |
| Control 1/2/3 Hold/Peak       | --                                             | --                                               | --                                             | --                                                    | --                                             | --                                                              | HOLD                                        | PEAK                                                           |
| PATCH TONE PARAMETER          |                                                |                                                  |                                                |                                                       |                                                |                                                                 |                                             |                                                                |
| Hold-2 Control Switch         | ON                                             | --                                               | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                             |
| Ctrl 1/2/3 Dest.1/2/3/4       | Other than OFF                                 | --                                               | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                             |
| Ctrl 1/2/3 Depth1/2/3/4       | Other than 0                                   | --                                               | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                             |
| LFO1/2 External Sync          | --                                             | --                                               | TAP                                            | --                                                    | --                                             | --                                                              | --                                          | --                                                             |
| Tone Delay Mode               | --                                             | --                                               | --                                             | TAP-SYNC                                              | --                                             | --                                                              | --                                          | --                                                             |

`--`: Need not be set

### 2. Performance Mode (Hold 2) (14)

| Parameters to be set          | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |              | LFO rate is changed to the HOLD-2 ON interval  | Tone delay time is changed to the HOLD-2 ON interval  | When HOLD-2 is ON, Effect controller is held   | When HOLD-2 is ON, Effect controller is held at its peak level  | When HOLD-2 is ON, Tone controller is held  | When HOLD-2 is ON, Tone controller is held at its peak level    |
| ----------------------------- | ---------------------------------------------- | ------------------------------------------------ | ------------ | ---------------------------------------------- | ----------------------------------------------------- | ---------------------------------------------- | --------------------------------------------------------------- | ------------------------------------------- | --------------------------------------------------------------- |
|                               |                                                | Setting1                                         | Setting2     |                                                |                                                       |                                                |                                                                 |                                             |                                                                 |
| SYSTEM PARAMETER              |                                                |                                                  |              |                                                |                                                       |                                                |                                                                 |                                             |                                                                 |
| System Control Source 1/2     | CC69:HOLD-2                                    | CC69:HOLD-2                                      | CC69:HOLD-2  | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| Receive Hold-2                | ON                                             | ON                                               | ON           | ON                                             | ON                                                    | ON                                             | ON                                                              | ON                                          | ON                                                              |
| TAP Control Source            | --                                             | --                                               | --           | HOLD-2                                         | HOLD-2                                                | --                                             | --                                                              | --                                          | --                                                              |
| Hold Control Source           | --                                             | --                                               | --           | --                                             | --                                                    | HOLD-2                                         | --                                                              | HOLD-2                                      | --                                                              |
| Peak Control Source           | --                                             | --                                               | --           | --                                             | --                                                    | --                                             | HOLD-2                                                          | --                                          | HOLD-2                                                          |
| PERFORMANCE COMMON PARAMETER  |                                                |                                                  |              |                                                |                                                       |                                                |                                                                 |                                             |                                                                 |
| EFX Source                    | --                                             | 1 - 16                                           | PERFORM      | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| EFX Control Source 1/2        | --                                             | --                                               | SYS-CTRL1/2  | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| EFX Control Depth 1/2         | --                                             | --                                               | Other than 0 | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| PERFORMANCE PART PARAMETER    |                                                |                                                  |              |                                                |                                                       |                                                |                                                                 |                                             |                                                                 |
| MIDI Receive Switch           | ON                                             | ON                                               | ON           | ON                                             | ON                                                    | ON                                             | ON                                                              | ON                                          | ON                                                              |
| PATCH COMMON PARAMETER        |                                                |                                                  |              |                                                |                                                       |                                                |                                                                 |                                             |                                                                 |
| EFX Control Source 1/2        | --                                             | SYS-CTRL1/2                                      | --           | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| EFX Control Depth 1/2         | --                                             | Other than 0                                     | --           | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| Patch Control Source 1/2/3    | SYS-CTRL1/2                                    | --                                               | --           | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| EFX Control Hold/Peak         | --                                             | --                                               | --           | --                                             | --                                                    | HOLD                                           | PEAK                                                            | --                                          | --                                                              |
| Control 1/2/3 Hold/Peak       | --                                             | --                                               | --           | --                                             | --                                                    | --                                             | --                                                              | HOLD                                        | PEAK                                                            |
| PATCH TONE PARAMETER          |                                                |                                                  |              |                                                |                                                       |                                                |                                                                 |                                             |                                                                 |
| Hold-2 Control Switch         | ON                                             | --                                               | --           | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| Ctrl 1/2/3 Dest.1/2/3/4       | Other than OFF                                 | --                                               | --           | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| Ctrl 1/2/3 Depth1/2/3/4       | Other than 0                                   | --                                               | --           | --                                             | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| LFO1/2 External Sync          | --                                             | --                                               | --           | TAP                                            | --                                                    | --                                             | --                                                              | --                                          | --                                                              |
| Tone Delay Mode               | --                                             | --                                               | --           | --                                             | TAP-SYNC                                              | --                                             | --                                                              | --                                          | --                                                              |

`--`: Need not be set

---

## Seite 142 — Portamento Control · GP Effect 1 (Reverb)

### ○ Portamento Control

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 54H    | kkH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
kk = Note Number: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set       | The on-note glides to the pitch of the note turned on next | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |
| -------------------------- | ---------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ |
| SYSTEM PARAMETER           |                                                            |                                                |                                                  |
| System Control Source 1/2  | --                                                         | CC84:PORTAMENTO-CTRL                           | CC84:PORTAMENTO-CTRL                             |
| Receive Control Change     | ON                                                         | ON                                             | ON                                               |
| PATCH COMMON PARAMETER     |                                                            |                                                |                                                  |
| Portamento Switch          | ON                                                         | --                                             | --                                               |
| Portamento Mode            | LEGATO                                                     | --                                             | --                                               |
| EFX Control Source 1/2     | --                                                         | --                                             | SYS-CTRL1/2                                      |
| EFX Control Depth 1/2      | --                                                         | --                                             | Other than 0                                     |
| Patch Control Source 1/2/3 | --                                                         | SYS-CTRL1/2                                    | --                                               |
| PATCH TONE PARAMETER       |                                                            |                                                |                                                  |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                                                         | Other than OFF                                 | --                                               |
| Ctrl 1/2/3 Depth1/2/3/4    | --                                                         | Other than 0                                   | --                                               |

`--`: Need not be set

### 2. Performance Mode

| Parameters to be set         | The on-note glides to the pitch of the note turned on next | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |                      |
| ---------------------------- | ---------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ | -------------------- |
|                              |                                                            |                                                | Setting1                                         | Setting2             |
| SYSTEM PARAMETER             |                                                            |                                                |                                                  |                      |
| System Control Source 1/2    | --                                                         | CC84:PORTAMENTO-CTRL                           | CC84:PORTAMENTO-CTRL                             | CC84:PORTAMENTO-CTRL |
| Receive Control Change       | ON                                                         | ON                                             | ON                                               | ON                   |
| PERFORMANCE COMMON PARAMETER |                                                            |                                                |                                                  |                      |
| EFX Source                   | --                                                         | --                                             | 1 - 16                                           | PERFORM              |
| EFX Control Source 1/2       | --                                                         | --                                             | --                                               | SYS-CTRL1/2          |
| EFX Control Depth 1/2        | --                                                         | --                                             | --                                               | Other than 0         |
| PERFORMANCE PART PARAMETER   |                                                            |                                                |                                                  |                      |
| MIDI Receive Switch          | ON                                                         | ON                                             | ON                                               | ON                   |
| PATCH COMMON PARAMETER       |                                                            |                                                |                                                  |                      |
| Portamento Switch            | ON                                                         | --                                             | --                                               | --                   |
| Portamento Mode              | LEGATO                                                     | --                                             | --                                               | --                   |
| EFX Control Source 1/2       | --                                                         | --                                             | SYS-CTRL1/2                                      | --                   |
| EFX Control Depth 1/2        | --                                                         | --                                             | Other than 0                                     | --                   |
| Patch Control Source 1/2/3   | --                                                         | SYS-CTRL1/2                                    | --                                               | --                   |
| PATCH TONE PARAMETER         |                                                            |                                                |                                                  |                      |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                         | Other than OFF                                 | --                                               | --                   |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                         | Other than 0                                   | --                                               | --                   |

`--`: Need not be set

### ○ GP Effect 1 (Reverb) (15)

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 5BH    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Reverb level: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set       | Reverb level is changed | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |
| -------------------------- | ----------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ |
|                            |                         | Setting1                                       | Setting2       | Setting1                                         | Setting2     |
| SYSTEM PARAMETER           |                         |                                                |                |                                                  |              |
| System Control Source 1/2  | CC91:REVERB             | --                                             | --             | --                                               | --           |
| Receive Control Change     | ON                      | ON                                             | ON             | ON                                               | ON           |
| PATCH COMMON PARAMETER     |                         |                                                |                |                                                  |              |
| EFX Control Source 1/2     | --                      | --                                             | --             | REVERB                                           | SYS-CTRL1/2  |
| EFX Control Depth 1/2      | --                      | --                                             | --             | Other than 0                                     | Other than 0 |
| Patch Control Source 1/2/3 | --                      | REVERB                                         | SYS-CTRL1/2    | --                                               | --           |
| PATCH TONE PARAMETER       |                         |                                                |                |                                                  |              |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                      | Other than OFF                                 | Other than OFF | --                                               | --           |
| Ctrl 1/2/3 Depth1/2/3/4    | --                      | Other than 0                                   | Other than 0   | --                                               | --           |

This message, when received in patch mode, will not affect the reverb send level.

`--`: Need not be set

---

## Seite 143 — GP Effect 1 (Performance) · GP Effect 3 (Chorus) · RPN LSB

### 2. Performance Mode (Reverb)

| Parameters to be set         | Reverb send level of the part on the receiving channel will be changed | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |              |              |
| ---------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ | ------------ | ------------ |
|                              |                                                                        | Setting1                                       | Setting2       | Setting1                                         | Setting2     | Setting3     | Setting4     |
| SYSTEM PARAMETER             |                                                                        |                                                |                |                                                  |              |              |              |
| System Control Source 1/2    | --                                                                     | --                                             | --             | --                                               | --           | CC91:REVERB  | CC91:REVERB  |
| Receive Control Change       | ON                                                                     | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PERFORMANCE COMMON PARAMETER |                                                                        |                                                |                |                                                  |              |              |              |
| EFX Source                   | --                                                                     | --                                             | --             | 1 - 16                                           | PERFORM      | 1 - 16       | PERFORM      |
| EFX Control Source 1/2       | --                                                                     | --                                             | --             | --                                               | REVERB       | --           | SYS-CTRL1/2  |
| EFX Control Depth 1/2        | --                                                                     | --                                             | --             | --                                               | Other than 0 | --           | Other than 0 |
| PERFORMANCE PART PARAMETER   |                                                                        |                                                |                |                                                  |              |              |              |
| MIDI Receive Switch          | ON                                                                     | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PATCH COMMON PARAMETER       |                                                                        |                                                |                |                                                  |              |              |              |
| EFX Control Source 1/2       | --                                                                     | --                                             | --             | REVERB                                           | --           | SYS-CTRL1/2  | --           |
| EFX Control Depth 1/2        | --                                                                     | --                                             | --             | Other than 0                                     | --           | Other than 0 | --           |
| Patch Control Source 1/2/3   | --                                                                     | REVERB                                         | SYS-CTRL1/2    | --                                               | --           | --           | --           |
| PATCH TONE PARAMETER         |                                                                        |                                                |                |                                                  |              |              |              |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                                     | Other than OFF                                 | Other than OFF | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                                     | Other than 0                                   | Other than 0   | --                                               | --           | --           | --           |

`--`: Need not be set

### ○ GP Effect 3 (Chorus) (15)

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 5DH    | vvH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Chorus level: 00H - 7FH (0 - 127)

### 1. Patch Mode (Chorus)

| Parameters to be set       | Chorus level is changed | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |
| -------------------------- | ----------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ |
|                            |                         | Setting1                                       | Setting2       | Setting1                                         | Setting2     |
| SYSTEM PARAMETER           |                         |                                                |                |                                                  |              |
| System Control Source 1/2  | CC93:CHORUS             | --                                             | --             | --                                               | --           |
| Receive Control Change     | ON                      | ON                                             | ON             | ON                                               | ON           |
| PATCH COMMON PARAMETER     |                         |                                                |                |                                                  |              |
| EFX Control Source 1/2     | --                      | --                                             | --             | CHORUS                                           | SYS-CTRL1/2  |
| EFX Control Depth 1/2      | --                      | --                                             | --             | Other than 0                                     | Other than 0 |
| Patch Control Source 1/2/3 | --                      | CHORUS                                         | SYS-CTRL1/2    | --                                               | --           |
| PATCH TONE PARAMETER       |                         |                                                |                |                                                  |              |
| Ctrl 1/2/3 Dest.1/2/3/4    | --                      | Other than OFF                                 | Other than OFF | --                                               | --           |
| Ctrl 1/2/3 Depth1/2/3/4    | --                      | Other than 0                                   | Other than 0   | --                                               | --           |

This message, when received in patch mode, will not affect the chorus send level.

`--`: Need not be set

### 2. Performance Mode (Chorus)

| Parameters to be set         | Chorus send level of the part on the receiving channel will be changed | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |              |              |
| ---------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ | ------------ | ------------ |
|                              |                                                                        | Setting1                                       | Setting2       | Setting1                                         | Setting2     | Setting3     | Setting4     |
| SYSTEM PARAMETER             |                                                                        |                                                |                |                                                  |              |              |              |
| System Control Source 1/2    | --                                                                     | --                                             | --             | --                                               | --           | CC93:CHORUS  | CC93:CHORUS  |
| Receive Control Change       | ON                                                                     | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PERFORMANCE COMMON PARAMETER |                                                                        |                                                |                |                                                  |              |              |              |
| EFX Source                   | --                                                                     | --                                             | --             | 1 - 16                                           | PERFORM      | 1 - 16       | PERFORM      |
| EFX Control Source 1/2       | --                                                                     | --                                             | --             | --                                               | CHORUS       | --           | SYS-CTRL1/2  |
| EFX Control Depth 1/2        | --                                                                     | --                                             | --             | --                                               | Other than 0 | --           | Other than 0 |
| PERFORMANCE PART PARAMETER   |                                                                        |                                                |                |                                                  |              |              |              |
| MIDI Receive Switch          | ON                                                                     | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PATCH COMMON PARAMETER       |                                                                        |                                                |                |                                                  |              |              |              |
| EFX Control Source 1/2       | --                                                                     | --                                             | --             | CHORUS                                           | --           | SYS-CTRL1/2  | --           |
| EFX Control Depth 1/2        | --                                                                     | --                                             | --             | Other than 0                                     | --           | Other than 0 | --           |
| Patch Control Source 1/2/3   | --                                                                     | CHORUS                                         | SYS-CTRL1/2    | --                                               | --           | --           | --           |
| PATCH TONE PARAMETER         |                                                                        |                                                |                |                                                  |              |              |              |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                                     | Other than OFF                                 | Other than OFF | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                                     | Other than 0                                   | Other than 0   | --                                               | --           | --           | --           |

`--`: Need not be set

### ○ RPN LSB

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 64H    | llH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
ll = Lower byte of the parameter number specified by RPN

---
## Seite 144 — RPN MSB · Data Entry LSB

### ○ RPN MSB

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 65H    | mmH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
mm = Upper byte of the parameter number specified by RPN.

The following is commonly applied to both RPN MSB and LSB.

### 1. Patch Mode

| Parameters to be set                            | The parameter No. specified by RPN will be changed |
| ----------------------------------------------- | -------------------------------------------------- |
| SYSTEM PARAMETER / Receive Control Change       | ON                                                 |

### 2. Performance Mode

| Parameters to be set                             | Parameter No. specified by RPN, of the part on the receiving channel will be changed |
| ------------------------------------------------ | ------------------------------------------------------------------------------------ |
| SYSTEM PARAMETER / Receive Control Change        | ON                                                                                   |
| PERFORMANCE PART PARAMETER / MIDI Receive Switch | ON                                                                                   |

### ○ Data Entry LSB

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 26H    | llH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
ll = Value for the parameter specified by RPN

### 1. Patch Mode

| Parameters to be set        | LSB of the data for the parameter specified by RPN LSB/MSB will be changed | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |
| --------------------------- | -------------------------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ |
| SYSTEM PARAMETER            |                                                                            |                                                |                                                  |
| System Control Source 1/2   | --                                                                         | CC38:DATA-ENTRY                                | CC38:DATA-ENTRY                                  |
| Receive Control Change      | ON                                                                         | ON                                             | ON                                               |
| PATCH COMMON PARAMETER      |                                                                            |                                                |                                                  |
| EFX Control Source 1/2      | --                                                                         | --                                             | SYS-CTRL1/2                                      |
| EFX Control Depth 1/2       | --                                                                         | --                                             | Other than 0                                     |
| Patch Control Source 1/2/3  | --                                                                         | SYS-CTRL1/2                                    | --                                               |
| PATCH TONE PARAMETER        |                                                                            |                                                |                                                  |
| Ctrl 1/2/3 Dest.1/2/3/4     | --                                                                         | Other than OFF                                 | --                                               |
| Ctrl 1/2/3 Depth1/2/3/4     | --                                                                         | Other than 0                                   | --                                               |

`--`: Need not be set

### 2. Performance Mode

| Parameters to be set         | LSB of the data for parameter specified by RPN LSB/MSB, of the part on the receiving channel will be changed   | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |                 |
| ---------------------------- | -------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ | --------------- |
|                              |                                                                                                                |                                                | Setting1                                         | Setting2        |
| SYSTEM PARAMETER             |                                                                                                                |                                                |                                                  |                 |
| System Control Source 1/2    | --                                                                                                             | CC38:DATA-ENTRY                                | CC38:DATA-ENTRY                                  | CC38:DATA-ENTRY |
| Receive Control Change       | ON                                                                                                             | ON                                             | ON                                               | ON              |
| PERFORMANCE COMMON PARAMETER |                                                                                                                |                                                |                                                  |                 |
| EFX Source                   | --                                                                                                             | --                                             | 1 - 16                                           | PERFORM         |
| EFX Control Source 1/2       | --                                                                                                             | --                                             | --                                               | SYS-CTRL1/2     |
| EFX Control Depth 1/2        | --                                                                                                             | --                                             | --                                               | Other than 0    |
| PERFORMANCE PART PARAMETER   |                                                                                                                |                                                |                                                  |                 |
| MIDI Receive Switch          | ON                                                                                                             | ON                                             | ON                                               | ON              |
| PATCH COMMON PARAMETER       |                                                                                                                |                                                |                                                  |                 |
| EFX Control Source 1/2       | --                                                                                                             | --                                             | SYS-CTRL1/2                                      | --              |
| EFX Control Depth 1/2        | --                                                                                                             | --                                             | Other than 0                                     | --              |
| Patch Control Source 1/2/3   | --                                                                                                             | SYS-CTRL1/2                                    | --                                               | --              |
| PATCH TONE PARAMETER         |                                                                                                                |                                                |                                                  |                 |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                                                                             | Other than OFF                                 | --                                               | --              |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                                                                             | Other than 0                                   | --                                               | --              |

`--`: Need not be set

---

## Seite 145 — Data Entry MSB · Description of RPN

### ○ Data Entry MSB

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 06H    | mmH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
mm = Value for the parameter specified by RPN

### 1. Patch Mode

| Parameters to be set        | MSB of the data for the parameter specified by RPN LSB/MSB will be changed  | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |
| --------------------------- | --------------------------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ |
| SYSTEM PARAMETER            |                                                                             |                                                |                                                  |
| System Control Source 1/2   | --                                                                          | CC6:DATA-ENTRY                                 | CC6:DATA-ENTRY                                   |
| Receive Control Change      | ON                                                                          | ON                                             | ON                                               |
| PATCH COMMON PARAMETER      |                                                                             |                                                |                                                  |
| EFX Control Source 1/2      | --                                                                          | --                                             | SYS-CTRL1/2                                      |
| EFX Control Depth 1/2       | --                                                                          | --                                             | Other than 0                                     |
| Patch Control Source 1/2/3  | --                                                                          | SYS-CTRL1/2                                    | --                                               |
| PATCH TONE PARAMETER        |                                                                             |                                                |                                                  |
| Ctrl 1/2/3 Dest.1/2/3/4     | --                                                                          | Other than OFF                                 | --                                               |
| Ctrl 1/2/3 Depth1/2/3/4     | --                                                                          | Other than 0                                   | --                                               |

`--`: Need not be set

### 2. Performance Mode

| Parameters to be set         | MSB of the data for parameter specified by RPN LSB/MSB, of the part on the receiving channel will be changed  | Tone parameters can be controlled in real-time | Effect parameters can be controlled in real-time |                |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------ | -------------- |
|                              |                                                                                                               |                                                | Setting1                                         | Setting2       |
| SYSTEM PARAMETER             |                                                                                                               |                                                |                                                  |                |
| System Control Source 1/2    | --                                                                                                            | CC6:DATA-ENTRY                                 | CC6:DATA-ENTRY                                   | CC6:DATA-ENTRY |
| Receive Control Change       | ON                                                                                                            | ON                                             | ON                                               | ON             |
| PERFORMANCE COMMON PARAMETER |                                                                                                               |                                                |                                                  |                |
| EFX Source                   | --                                                                                                            | --                                             | 1 - 16                                           | PERFORM        |
| EFX Control Source 1/2       | --                                                                                                            | --                                             | --                                               | SYS-CTRL1/2    |
| EFX Control Depth 1/2        | --                                                                                                            | --                                             | --                                               | Other than 0   |
| PERFORMANCE PART PARAMETER   |                                                                                                               |                                                |                                                  |                |
| MIDI Receive Switch          | ON                                                                                                            | ON                                             | ON                                               | ON             |
| PATCH COMMON PARAMETER       |                                                                                                               |                                                |                                                  |                |
| EFX Control Source 1/2       | --                                                                                                            | --                                             | SYS-CTRL1/2                                      | --             |
| EFX Control Depth 1/2        | --                                                                                                            | --                                             | Other than 0                                     | --             |
| Patch Control Source 1/2/3   | --                                                                                                            | SYS-CTRL1/2                                    | --                                               | --             |
| PATCH TONE PARAMETER         |                                                                                                               |                                                |                                                  |                |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                                                                            | Other than OFF                                 | --                                               | --             |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                                                                            | Other than 0                                   | --                                               | --             |

`--`: Need not be set

### ■ Description of RPN

RPNs (Registered Parameter Numbers) are functions defined by MIDI standard. Each RPN may be used to change parameters of equipment to vary characteristics of tune, performance, etc.

The JV-1080 can recognize the four RPNs: Pitch Bend Sensitivity (RPN#0), Fine Tuning (RPN#1), Coarse Tuning (RPN#2) and RPN Reset (RPN#127).

To effect RPN, first designate the parameter to be controlled using RPN MSB and RPN LSB, and then specify the value of designated parameter in the data entry.

| BnH 65H mmH | BnH 64H llH | BnH 06H xxH      | BnH 26H yyH      |
| ----------- | ----------- | ---------------- | ---------------- |
| (RPN MSB)   | (RPN LSB)   | (Data Entry MSB) | (Data Entry LSB) |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16

| RPN MSB | RPN LSB | Data Entry MSB | Data Entry LSB | Function                                                                                                                                                                                                                                      |
| ------- | ------- | -------------- | -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| mm      | ll      | xx             | yy             |                                                                                                                                                                                                                                               |
| 00H     | 00H     | mmH            | --             | Pitch bend sensitivity: mm : 00H - 0CH (0 - 12 in unit of semitones) ll : Ignored. Up to 1 octave in unit of semitones. * Common to BENDER RANGE UP and BENDER RANGE DOWN. * Rhythm part (part 10) ignores this function.                     |
| 00H     | 01H     | mmH            | llH            | Fine Tuning: mm, ll : 00 00H - 40 00H - 7F 7FH (-8192 x 50/8192 - 0 - +8192 x 50/8192 cent) * In patch mode, sets master tune. * In performance mode, sets fine tune of a part. * When received on the control channel, sets the master tune. |
| 00H     | 02H     | mmH            | --             | Coarse Tuning: mm : 10H - 40H - 70H (-48 - 0 - +48 in unit of semitones) ll : Ignored * Ignored in patch mode. * In performance mode, sets coarse tune of a part.                                                                             |
| 7F      | 7F      | --             | --             | RPN Reset: Cancels the settings made by RPN(s). Internal settings remain unchanged. mm, ll : Ignored                                                                                                                                          |

- RPN is received either MSB first or LSB first.
- Data entry data must be sent MSB first to be correctly received. (LSB is cleared to 0 when MSB is received.)

---

## Seite 146 — Program Change · Channel Pressure

### ○ Program Change

| Status | Second |
| ------ | ------ |
| CnH    | ppH    |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
pp = Program number: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set                      | Changes patches (patch number is the program number plus 1)  |
| ----------------------------------------- | ------------------------------------------------------------ |
| SYSTEM PARAMETER / Receive Program Change | ON                                                           |

### 2. Performance Mode

Changes performance when received on the control channel.

| Parameters to be set                                | Patch of the part on the receiving channel will be changed (the patch number is the program number plus 1)   |
| --------------------------------------------------- | ------------------------------------------------------------------------------------------------------------ |
| SYSTEM PARAMETER / Receive Program Change           | ON                                                                                                           |
| PERFORMANCE PART PARAMETER / MIDI Receive Switch    | ON                                                                                                           |
| PERFORMANCE PART PARAMETER / Receive Program Change | ON                                                                                                           |

### ○ Channel Pressure

| Status | Second |
| ------ | ------ |
| DnH    | vvH    |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
vv = Pressure value: 00H - 7FH (0 - 127)

### 1. Patch Mode

| Parameters to be set        | Tone parameters can be controlled in real-time |                     | Effect parameters can be controlled in real-time |                     |
| --------------------------- | ---------------------------------------------- | ------------------- | ------------------------------------------------ | ------------------- |
|                             | Setting1                                       | Setting2            | Setting1                                         | Setting2            |
| SYSTEM PARAMETER            |                                                |                     |                                                  |                     |
| System Control Source 1/2   | --                                             | AFTERTOUCH          | --                                               | AFTERTOUCH          |
| Receive Aftertouch          | ON                                             | ON                  | ON                                               | ON                  |
| Aftertouch Source           | CH-AFTER or CH&POLY                            | CH-AFTER or CH&POLY | CH-AFTER or CH&POLY                              | CH-AFTER or CH&POLY |
| PATCH COMMON PARAMETER      |                                                |                     |                                                  |                     |
| EFX Control Source 1/2      | --                                             | --                  | AFTERTOUCH                                       | SYS-CTRL1/2         |
| EFX Control Depth 1/2       | --                                             | --                  | Other than 0                                     | Other than 0        |
| Patch Control Source 1/2/3  | AFTERTOUCH                                     | SYS-CTRL1/2         | --                                               | --                  |
| PATCH TONE PARAMETER        |                                                |                     |                                                  |                     |
| Ctrl 1/2/3 Dest.1/2/3/4     | Other than OFF                                 | Other than OFF      | --                                               | --                  |
| Ctrl 1/2/3 Depth1/2/3/4     | Other than 0                                   | Other than 0        | --                                               | --                  |

`--`: Need not be set

### 2. Performance Mode

| Parameters to be set         | Tone parameters can be controlled in real-time |                     | Effect parameters can be controlled in real-time |                     |                     |                     |
| ---------------------------- | ---------------------------------------------- | ------------------- | ------------------------------------------------ | ------------------- | ------------------- | ------------------- |
|                              | Setting1                                       | Setting2            | Setting1                                         | Setting2            | Setting3            | Setting4            |
| SYSTEM PARAMETER             |                                                |                     |                                                  |                     |                     |                     |
| System Control Source 1/2    | --                                             | AFTERTOUCH          | --                                               | --                  | AFTERTOUCH          | AFTERTOUCH          |
| Receive Aftertouch           | ON                                             | ON                  | ON                                               | ON                  | ON                  | ON                  |
| Aftertouch Source            | CH-AFTER or CH&POLY                            | CH-AFTER or CH&POLY | CH-AFTER or CH&POLY                              | CH-AFTER or CH&POLY | CH-AFTER or CH&POLY | CH-AFTER or CH&POLY |
| PERFORMANCE COMMON PARAMETER |                                                |                     |                                                  |                     |                     |                     |
| EFX Source                   | --                                             | --                  | 1 - 16                                           | PERFORM             | 1 - 16              | PERFORM             |
| EFX Control Source 1/2       | --                                             | --                  | --                                               | AFTERTOUCH          | --                  | SYS-CTRL1/2         |
| EFX Control Depth 1/2        | --                                             | --                  | --                                               | Other than 0        | --                  | Other than 0        |
| PERFORMANCE PART PARAMETER   |                                                |                     |                                                  |                     |                     |                     |
| MIDI Receive Switch          | ON                                             | ON                  | ON                                               | ON                  | ON                  | ON                  |
| PATCH COMMON PARAMETER       |                                                |                     |                                                  |                     |                     |                     |
| EFX Control Source 1/2       | --                                             | --                  | AFTERTOUCH                                       | --                  | SYS-CTRL1/2         | --                  |
| EFX Control Depth 1/2        | --                                             | --                  | Other than 0                                     | --                  | Other than 0        | --                  |
| Patch Control Source 1/2/3   | AFTERTOUCH                                     | SYS-CTRL1/2         | --                                               | --                  | --                  | --                  |
| PATCH TONE PARAMETER         |                                                |                     |                                                  |                     |                     |                     |
| Ctrl 1/2/3 Dest.1/2/3/4      | Other than OFF                                 | Other than OFF      | --                                               | --                  | --                  | --                  |
| Ctrl 1/2/3 Depth1/2/3/4      | Other than 0                                   | Other than 0        | --                                               | --                  | --                  | --                  |

`--`: Need not be set

---

## Seite 147 — Pitch Bend Change · Channel Mode Messages

### ○ Pitch Bend Change

| Status | Second | Third |
| ------ | ------ | ----- |
| EnH    | llH    | mmH   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16
ll, mm = Pitch bend change: 00 00H - 40 00H - 7F 7FH [-8192 - 0 - +8192]

### 1. Patch Mode

| Parameters to be set        | Changes pitch of note | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |
| --------------------------- | --------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ |
|                             |                       | Setting1                                       | Setting2       | Setting1                                         | Setting2     |
| SYSTEM PARAMETER            |                       |                                                |                |                                                  |              |
| System Control Source 1/2   | --                    | --                                             | BENDER         | --                                               | BENDER       |
| Receive Bender              | ON                    | ON                                             | ON             | ON                                               | ON           |
| PATCH COMMON PARAMETER      |                       |                                                |                |                                                  |              |
| EFX Control Source 1/2      | --                    | --                                             | --             | BENDER                                           | SYS-CTRL1/2  |
| EFX Control Depth 1/2       | --                    | --                                             | --             | Other than 0                                     | Other than 0 |
| Patch Control Source 1/2/3  | --                    | BENDER                                         | SYS-CTRL1/2    | --                                               | --           |
| PATCH TONE PARAMETER        |                       |                                                |                |                                                  |              |
| Bender Control Switch       | ON                    | --                                             | --             | --                                               | --           |
| Bend Range Upper/Lower      | Other than 0          | --                                             | --             | --                                               | --           |
| Pan Control Switch          | --                    | --                                             | --             | --                                               | --           |
| Ctrl 1/2/3 Dest.1/2/3/4     | --                    | Other than OFF                                 | Other than OFF | --                                               | --           |
| Ctrl 1/2/3 Depth1/2/3/4     | --                    | Other than 0                                   | Other than 0   | --                                               | --           |

`--`: Need not be set

### 2. Performance Mode

| Parameters to be set         | Changes note pitch of the part on the receiving channel  | Tone parameters can be controlled in real-time |                | Effect parameters can be controlled in real-time |              |              |              |
| ---------------------------- | -------------------------------------------------------- | ---------------------------------------------- | -------------- | ------------------------------------------------ | ------------ | ------------ | ------------ |
|                              |                                                          | Setting1                                       | Setting2       | Setting1                                         | Setting2     | Setting3     | Setting4     |
| SYSTEM PARAMETER             |                                                          |                                                |                |                                                  |              |              |              |
| System Control Source 1/2    | --                                                       | --                                             | BENDER         | --                                               | --           | BENDER       | BENDER       |
| Receive Bender               | ON                                                       | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PERFORMANCE COMMON PARAMETER |                                                          |                                                |                |                                                  |              |              |              |
| EFX Source                   | --                                                       | --                                             | --             | 1 - 16                                           | PERFORM      | 1 - 16       | PERFORM      |
| EFX Control Source 1/2       | --                                                       | --                                             | --             | --                                               | BENDER       | --           | SYS-CTRL1/2  |
| EFX Control Depth 1/2        | --                                                       | --                                             | --             | --                                               | Other than 0 | --           | Other than 0 |
| PERFORMANCE PART PARAMETER   |                                                          |                                                |                |                                                  |              |              |              |
| MIDI Receive Switch          | ON                                                       | ON                                             | ON             | ON                                               | ON           | ON           | ON           |
| PATCH COMMON PARAMETER       |                                                          |                                                |                |                                                  |              |              |              |
| EFX Control Source 1/2       | --                                                       | --                                             | --             | BENDER                                           | --           | SYS-CTRL1/2  | --           |
| EFX Control Depth 1/2        | --                                                       | --                                             | --             | Other than 0                                     | --           | Other than 0 | --           |
| Patch Control Source 1/2/3   | --                                                       | BENDER                                         | SYS-CTRL1/2    | --                                               | --           | --           | --           |
| PATCH TONE PARAMETER         |                                                          |                                                |                |                                                  |              |              |              |
| Bender Control Switch        | ON                                                       | --                                             | --             | --                                               | --           | --           | --           |
| Bend Range Upper/Lower       | Other than 0                                             | --                                             | --             | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Dest.1/2/3/4      | --                                                       | Other than OFF                                 | Other than OFF | --                                               | --           | --           | --           |
| Ctrl 1/2/3 Depth1/2/3/4      | --                                                       | Other than 0                                   | Other than 0   | --                                               | --           | --           | --           |

`--`: Need not be set

### ■ Channel Mode Message

### ○ All Sounds Off

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 78H    | 00H   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16

- Turns off all sounding notes on the MIDI channel. However, the state of channel messages does not change.

### ○ Reset All Controllers

| Status | Second | Third |
| ------ | ------ | ----- |
| BnH    | 79H    | 00H   |

n = MIDI channel number: 0H - FH (0 - 15) 0 = ch.1 15 = ch.16

- Upon receiving this message, the JV-1080 changes settings of the controller as follows:

| Controller                          | Settings                                     |
| ----------------------------------- | -------------------------------------------- |
| Modulation                          | 0 (min)                                      |
| Breath                              | 0 (min)                                      |
| Foot                                | 0 (min)                                      |
| Volume                              | 127 (max)                                    |
| Pan                                 | 64 (center)                                  |
| Expression                          | 0 (min) (volume is set at max.)              |
| Hold 1                              | 0 (off)                                      |
| Sostenuto                           | 0 (off)                                      |
| Soft Pedal                          | 0 (off)                                      |
| Hold 2                              | 0 (off)                                      |
| Channel Pressure                    | 0 (min)                                      |
| Polyphonic Pressure                 | 0 (min)                                      |
| Pitch Bend Change                   | 0 (center)                                   |
| RPN                                 | Undefined: does not affect the internal data |
| General purpose system controller 1 | 0 (min)                                      |
| General purpose system controller 2 | 0 (min)                                      |

---

---

## Seite 148 — Channel Mode Messages (Fortsetzung) · System Real Time Messages

### ● All Notes Off

| Status   | Second   | Third    |
| -------- | -------- | -------- |
| BnH      | 7BH      | 00H      |

n = MIDI channel number: 0H – FH (0 – 15)  0 = ch.1  15 = ch.16

- Turns off all MIDI-on notes on the MIDI channel.
  However, sound continues when Hold 1 and/or SOSTENUTO is ON.

### ● OMNI OFF

| Status   | Second   | Third    |
| -------- | -------- | -------- |
| BnH      | 7CH      | 00H      |

n = MIDI channel number: 0H – FH (0 – 15)  0 = ch.1  15 = ch.16

- Serves as All Notes Off.

### ● OMNI ON

| Status   | Second   | Third    |
| -------- | -------- | -------- |
| BnH      | 7DH      | 00H      |

n = MIDI channel number: 0H – FH (0 – 15)  0 = ch.1  15 = ch.16

- Serves as All Notes Off, and not OMNI ON.

### ● MONO

| Status   | Second   | Third    |
| -------- | -------- | -------- |
| BnH      | 7EH      | mmH      |

n = MIDI channel number: 0H – FH (0 – 15)  0 = ch.1  15 = ch.16
mm = Number of MONOs: 00H – 0FH (0 – 15)

- The key assign mode of the patch common parameter is changed to SOLO.
- Serves as All Notes Off and sets the mode to Mode 4 (m = 1).

### ● POLY

| Status   | Second   | Third    |
| -------- | -------- | -------- |
| BnH      | 7FH      | 00H      |

n = MIDI channel number: 0H – FH (0 – 15)  0 = ch.1  15 = ch.16

- The key assign mode of the patch common parameter is changed to POLY.
- Serves as All Notes Off and sets the mode to Mode 3.

### ■ System Real Time Messages

#### ● Active Sensing

| Status   |
| -------- |
| FEH      |

- When JV-1080 receives Active Sensing, it measures time intervals between incoming messages. If the subsequent message will not come within 400 ms after the previous one, JV-1080 turns off all MIDI-on notes, operates as if it receives a Reset All Controllers message, and stops measuring message intervals.

#### ● Timing Clock

| Status   |
| -------- |
| F8H      |

*(MIDI-Clock-Synchronisation; Überschrift im Scan nicht klar erkennbar)*

### 1. Patch Mode

| Parameters to be set                         | Changes LFO rate   | Changes tone delay time   | Changes delay time of effect                                     | Changes the step rate of effect (Flanger)   |
| -------------------------------------------- | ------------------ | ------------------------- | ---------------------------------------------------------------- | ------------------------------------------- |
| SYSTEM PARAMETER: Clock Source               | MIDI               | MIDI                      | MIDI                                                             | MIDI                                        |
| PATCH COMMON PARAMETER: EFX Type             | --                 | --                        | 19.TRIPLE-TAP-DELAY or 20.QUADRUPLE-TAP-DELAY                    | 16.STEP-FLANGER                             |
| PATCH TONE PARAMETER: LFO1/2 External Sync   | CLOCK              | --                        | --                                                               | --                                          |
| PATCH TONE PARAMETER: Tone Delay Mode        | --                 | CLOCK-SYNC                | --                                                               | --                                          |

`--`: Need not be set

---

## Seite 149 — 2. Performance Mode · System Exclusive · Exclusive Communications

### 2. Performance Mode

| Parameters to be set                                  | Changes LFO rate   | Changes tone delay time   | Changes delay time of effect (Setting 1)        | Setting 2                                       | Changes step rate of effect, Flanger (Setting 1)   | Setting 2            |
| ----------------------------------------------------- | ------------------ | ------------------------- | ----------------------------------------------- | ----------------------------------------------- | -------------------------------------------------- | -------------------- |
| SYSTEM PARAMETER: Clock Source                        | MIDI               | MIDI                      | MIDI                                            | MIDI                                            | MIDI                                               | MIDI                 |
| PERFORMANCE COMMON PARAMETER: EFX Type                | --                 | --                        | --                                              | 19.TRIPLE-TAP-DELAY or 20.QUADRUPLE-TAP-DELAY   | --                                                 | 16.STEP-FLANGER      |
| PERFORMANCE COMMON PARAMETER: EFX Source              | --                 | --                        | 1 - 16                                          | PERFORM                                         | 1 - 16                                             | PERFORM              |
| PATCH COMMON PARAMETER: EFX Type                      | --                 | --                        | 19.TRIPLE-TAP-DELAY or 20.QUADRUPLE-TAP-DELAY   | --                                              | 16.STEP-FLANGER                                    | --                   |
| PATCH TONE PARAMETER: LFO1/2 External Sync            | CLOCK              | --                        | --                                              | --                                              | --                                                 | --                   |
| PATCH TONE PARAMETER: Tone Delay Mode                 | --                 | CLOCK-SYNC                | --                                              | --                                              | --                                                 | --                   |

`--`: Need not be set

### 1. TRANSMIT DATA — ■ System Exclusive Message *(Empfang)*

| Status   | Data byte        |
| -------- | ---------------- |
| F0H      | 41H ddH ...eeH   |
| F7H      |                  |

- F0H: System Exclusive
- 41H: Manufacturer ID (Roland)
- dd…ee = Data: 00H – 7FH (0 – 127)
- F7H: EOX (End Of Exclusive)

- The JV-1080 recognizes this message when the receive switch in the system parameter is set to ON.
  For detail, please refer to section 3: Roland exclusive message.

### 2. TRANSMIT DATA — ■ System Exclusive Message *(Senden)*

| Status   | Data byte        |
| -------- | ---------------- |
| F0H      | 41H ddH ...eeH   |
| F7H      |                  |

- F0H: System Exclusive
- 41H: Manufacturer ID (Roland)
- dd…ee = Data: 00H – 7FH (0 – 127)
- F7H: EOX (End Of Exclusive)

For detail, please refer to section 3: Roland exclusive message.

### 3. Exclusive communications

The JV-1080 can send and receive patch parameters etc. using the system exclusive message.
The model ID code of the JV-1080 is 6AH. The device ID code is to be determined by the unit number setting of MIDI function.
The JV-1080 ignores GS exclusive messages other than the scale tune parameter. The model ID of the GS is 42H.

### ■ One way communication

#### ● Request data 1 RQ1 (11H)

| Byte   | Description                   |
| ------ | ----------------------------- |
| F0H    | Exclusive status              |
| 41H    | Manufacturer ID (Roland)      |
| Dev    | Device ID (Dev = UNIT# - 1)   |
| 6AH    | Model ID (JV-1080)            |
| 11H    | Command ID (RQ1)              |
| aaH    | Address MSB                   |
| bbH    | Address                       |
| ccH    | Address                       |
| ddH    | Address LSB                   |
| eeH    | Size MSB                      |
| ffH    | Size                          |
| ggH    | Size                          |
| hhH    | Size LSB                      |
| sum    | Check sum                     |
| F7H    | EOX (End of exclusive)        |

- *Receive only*: the JV-1080 does not send this message.

---

## Seite 150 — DT1 · Parameter Address Map · System (*1-1, *1-1-1, *1-1-2)

### ● Data set 1 DT1 (12H)

#### 1. JV-1080 (MODEL ID = 6AH)

| Byte   | Description                   |
| ------ | ----------------------------- |
| F0H    | Exclusive status              |
| 41H    | Manufacturer ID (Roland)      |
| Dev    | Device ID (Dev = UNIT# - 1)   |
| 6AH    | Model ID (JV-1080)            |
| 12H    | Command ID (DT1)              |
| aaH    | Address MSB                   |
| bbH    | Address                       |
| ccH    | Address                       |
| ddH    | Address LSB                   |
| eeH    | Data (value)                  |
| ffH    | Data (value)                  |
| sum    | Check sum                     |
| F7H    | EOX (End of exclusive)        |

#### 2. GS (MODEL ID = 42H)

| Byte   | Description                   |
| ------ | ----------------------------- |
| F0H    | Exclusive status              |
| 41H    | Manufacturer ID (Roland)      |
| Dev    | Device ID (Dev = UNIT# - 1)   |
| 42H    | Model ID (GS)                 |
| 12H    | Command ID (DT1)              |
| aaH    | Address MSB                   |
| bbH    | Address                       |
| ccH    | Address LSB                   |
| eeH    | Data (value)                  |
| ffH    | Data (value)                  |
| sum    | Check sum                     |
| F7H    | EOX (End of exclusive)        |

- When the device ID is 7FH, the JV-1080 can receive the GS exclusive message even if the unit number is anything.

### 4. Parameter address map (MODEL ID = 6AH)

Address and size are configured in 7 bits, and expressed in hexadecimal.

|                  | MSB          | LSB                                  |
| ---------------- | ------------ | ------------------------------------ |
| Address (bin)    | 0aaa aaaa    | 0bbb bbbb 0ccc cccc 0ddd dddd        |
| Address (hex)    | aaH          | bbH ccH ddH                          |
| Size (bin)       | 0sss ssss    | 0ttt tttt 0uuu uuuu 0vvv vvvv        |
| Size (hex)       | ssH          | ttH uuH vvH                          |

### ■ Parameter base address

All data sent in exclusive message are given particular addresses to identify parameters. These addresses are the sum of the base address and the offset address. Some parameters are defined using multiple offsets. The address included in the message of a data set or a data request must be within the value shown in the table below.

**Note:** A pair of two addresses preceded by the symbol `#` represents a divided-by-two data, e.g. the data ABH (hex) is divided into 0AH and 0BH and sent in that order.

#### 1 JV-1080 (MODEL ID = 6AH) — Base addresses

| Start address   | Description                                       |
| --------------- | ------------------------------------------------- |
| 00 00 00 00     | System *1-1                                       |
| 01 00 00 00     | Temporary performance *1-2                        |
| 02 00 00 00     | Performance mode temporary patch (part 1) *1-3    |
| 02 01 00 00     | Performance mode temporary patch (part 2)         |
| :               |                                                   |
| 02 08 00 00     | Performance mode temporary patch (part 9)         |
| 02 09 00 00     | Temporary rhythm setup *1-4                       |
| 02 0A 00 00     | Performance mode temporary patch (part 11) *1-3   |
| :               |                                                   |
| 02 0F 00 00     | Performance mode temporary patch (part 16)        |
| 03 00 00 00     | Patch mode temporary patch *1-3                   |
| 10 00 00 00     | User performance USER:01 *1-2                     |
| 10 01 00 00     | User performance USER:02                          |
| :               |                                                   |
| 10 1F 00 00     | User performance USER:32                          |
| 10 40 00 00     | User rhythm setup USER:1 *1-4                     |
| 10 41 00 00     | User rhythm setup USER:2                          |
| 11 00 00 00     | User patch USER:001 *1-3                          |
| 11 01 00 00     | User patch USER:002                               |
| :               |                                                   |
| 11 7F 00 00     | User patch USER:128                               |
| 20 00 00 00     | Data card performance CARD:01 *1-2                |
| 20 01 00 00     | Data card performance CARD:02                     |
| :               |                                                   |
| 20 1F 00 00     | Data card performance CARD:32                     |
| 20 40 00 00     | Data card rhythm set CARD:1 *1-4                  |
| 20 41 00 00     | Data card rhythm set CARD:2                       |
| 21 00 00 00     | Data card patch CARD:001 *1-3                     |
| 21 01 00 00     | Data card patch CARD:002                          |
| :               |                                                   |
| 21 7F 00 00     | Data card patch CARD:128                          |

(`:` = im Original fortlaufende Zeilen, hier zusammengefasst)

#### *1-1 System

| Offset   | Description                    |
| -------- | ------------------------------ |
| 00 00    | System common *1-1-1           |
| 10 00    | Part 1 scale tune *1-1-2       |
| 11 00    | Part 2 scale tune              |
| :        |                                |
| 1F 00    | Part 16 scale tune             |
| 20 00    | Patch mode scale tune *1-1-2   |

#### *1-1-1 System common

| Offset   | Size Binary    | Description                    | Value                                                                          |
| -------- | -------------- | ------------------------------ | ------------------------------------------------------------------------------ |
| 00 00    | 0000 00aa      | Panel mode                     | 0 - 2 (PERFORMANCE, PATCH, GM)                                                 |
| 00 01    | 0aaa aaaa      | Performance number             | 0 - 127 (USER:01-USER:32, CARD:01-CARD:32, PR-A:01-PR-A:32, PR-B:01-PR-B:32)   |
| 00 02    | 0000 00aa      | Patch mode patch group         | 0 - 2 (USER, PCM, EXP) (1)                                                     |
| 00 03    | 0aaa aaaa      | Patch mode patch group ID      | 0 - 127                                                                        |
| 00 04    | 0000 aaaa #    | Patch mode patch number        | 0 - 254 (001 - 255)                                                            |
| 00 06    | 0aaa aaaa      | Master tune                    | 0 - 126 (427.4 - 452.6 Hz)                                                     |
| 00 07    | 0000 000a      | Scale tune switch              | 0 - 1 (OFF, ON)                                                                |
| 00 08    | 0000 000a      | EFX switch                     | 0 - 1 (OFF, ON)                                                                |
| 00 09    | 0000 000a      | Chorus switch                  | 0 - 1 (OFF, ON)                                                                |
| 00 0A    | 0000 000a      | Reverb switch                  | 0 - 1 (OFF, ON)                                                                |
| 00 0B    | 0000 000a      | Patch remain                   | 0 - 1 (OFF, ON)                                                                |
| 00 0C    | 0000 000a      | Clock source                   | 0 - 1 (INT, MIDI)                                                              |
| 00 0D    | 0000 0aaa      | Tap control source             | 0 - 4 (OFF, HOLD-1, SOSTENUTO, SOFT, HOLD-2)                                   |
| 00 0E    | 0000 0aaa      | Hold control source            | 0 - 4 (OFF, HOLD-1, SOSTENUTO, SOFT, HOLD-2)                                   |
| 00 0F    | 0000 0aaa      | Peak control source            | 0 - 4 (OFF, HOLD-1, SOSTENUTO, SOFT, HOLD-2)                                   |
| 00 10    | 0000 000a      | Volume control source          | 0 - 1 (VOLUME, VOLUME&EXPRESSION)                                              |
| 00 11    | 0000 00aa      | Aftertouch source              | 0 - 2 (CHANNEL, POLY, CH&POLY)                                                 |
| 00 12    | 0aaa aaaa      | System control source 1        | 0 - 97 (2)                                                                     |
| 00 13    | 0aaa aaaa      | System control source 2        | 0 - 97 (2)                                                                     |
| 00 14    | 0000 000a      | Receive program change         | 0 - 1 (OFF, ON)                                                                |
| 00 15    | 0000 000a      | Receive bank select            | 0 - 1 (OFF, ON)                                                                |
| 00 16    | 0000 000a      | Receive control change         | 0 - 1 (OFF, ON)                                                                |
| 00 17    | 0000 000a      | Receive modulation             | 0 - 1 (OFF, ON)                                                                |
| 00 18    | 0000 000a      | Receive volume                 | 0 - 1 (OFF, ON)                                                                |
| 00 19    | 0000 000a      | Receive hold-1                 | 0 - 1 (OFF, ON)                                                                |
| 00 1A    | 0000 000a      | Receive bender                 | 0 - 1 (OFF, ON)                                                                |
| 00 1B    | 0000 000a      | Receive aftertouch             | 0 - 1 (OFF, ON)                                                                |
| 00 1C    | 000a aaaa      | Control channel                | 0 - 16 (1 - 16, OFF)                                                           |
| 00 1D    | 0000 aaaa      | Patch receive channel          | 0 - 15 (1 - 16)                                                                |
| 00 1E    | 0000 000a      | Rhythm edit source             | 0 - 1 (PANEL, PANEL&MIDI)                                                      |
| 00 1F    | 0000 00aa      | Preview sound mode             | 0 - 2 (SINGLE, CHORD, PHRASE) (4)                                              |
| 00 20    | 0aaa aaaa      | Preview key set 1              | 0 - 127 (C-1 - G9)                                                             |
| 00 21    | 0aaa aaaa      | Preview velocity set 1         | 0 - 127 (OFF, 1 - 127)                                                         |
| 00 22    | 0aaa aaaa      | Preview key set 2              | 0 - 127 (C-1 - G9)                                                             |
| 00 23    | 0aaa aaaa      | Preview velocity set 2         | 0 - 127 (OFF, 1 - 127)                                                         |
| 00 24    | 0aaa aaaa      | Preview key set 3              | 0 - 127 (C-1 - G9)                                                             |
| 00 25    | 0aaa aaaa      | Preview velocity set 3         | 0 - 127 (OFF, 1 - 127)                                                         |
| 00 26    | 0aaa aaaa      | Preview key set 4              | 0 - 127 (C-1 - G9)                                                             |
| 00 27    | 0aaa aaaa      | Preview velocity set 4         | 0 - 127 (OFF, 1 - 127)                                                         |

**Total size: 00 00 00 28**

#### *1-1-2 Scale tune

| Offset   | Size Binary   | Description   | Value                 |
| -------- | ------------- | ------------- | --------------------- |
| 00 00    | 0aaa aaaa     | C             | 0 - 127 (-64 - +63)   |
| 00 01    | 0aaa aaaa     | C#            | 0 - 127 (-64 - +63)   |
| :        |               |               |                       |
| 00 0B    | 0aaa aaaa     | B             | 0 - 127 (-64 - +63)   |

**Total size: 00 00 00 0C**

#### Examples

**Example using RQ1:** To get all the data of the system common, send the following message to the JV-1080.

`F0H 41H 10H 6AH 11H 00H 00H 00H 00H 00H 00H 00H 28H 58H F7H` (6)

**Example using DT1:** To set the Control Channel of the system common to 1, send the following message to the JV-1080.

`F0H 41H 10H 6AH 12H 00H 00H 00H 1CH 00H 64H F7H`

### : Example of exclusive data (9)

*Setting the Reverb Type of the Temporary Performance to DELAY (DT1).*

The Parameter Address Map indicates that the starting address of the Temporary Performance is 01 00 00 00H, that the Performance Common offset address is 00 00H, and that the Reverb Type address is 00 28H. Thus, the address is:

```
01 00 00 00H
        00 00H
+)      00 28H
----------------
01 00 00 28H   (target address)
```

Since DELAY is parameter value 06H, the exclusive message looks like as follows:

```
F0H  41H  Dev  6AH  12H  01H 00H 00H 28H  06H  51H  F7H
(1)  (2)  (3)  (4)  (5)  <--- address ---> data  csum (6)
```

1. F0H: Exclusive status
2. 41H: Manufacturer ID (Roland)
3. Dev: Device ID — the unit number of the system common parameter minus 1. In this example, the unit number is 17: 17 − 1 = 16, which is expressed as 10H in hexadecimal.
4. 6AH: Model ID (JV-1080)
5. 12H: Command ID (DT1)
6. F7H: EOX (End of exclusive)

**Calculation of the check sum:** Roland exclusive messages (RQ1, DT1) are transmitted with a checksum at the end of the data (before F7H) to check that the data was received correctly. The checksum consists of a value whose least significant 7 bits are 0 when the address, size and the checksum itself are added.

```
aaH + bbH + ccH + ddH + eeH + ffH = total
total / 128 = quotient ... remainder
128 - remainder = checksum
```

For the message above:

```
decimal:  01H + 00H + 00H + 28H + 06H = 1 + 0 + 0 + 40 + 6 = 47 (total)
          47 / 128 = 0 (quotient) ... 47 (remainder)
          checksum = 128 - 47 = 81 = 51H
hex:      01H + 00H + 00H + 28H + 06H = 2FH;  80H - 2FH = 51H
```

This means that the message transmitted will be:

```
F0H 41H 10H 6AH 12H 01H 00H 00H 28H 06H 51H F7H
```

**Note:** In hexadecimal, each byte is counted from the right from the lower byte, and data is transmitted in order starting from the byte on the left (the highest byte).

---

## Seite 151 — Performance (*1-2, *1-2-1) · Examples

#### *1-2 Performance

| Offset   | Description                    |
| -------- | ------------------------------ |
| 00 00    | Performance common *1-2-1      |
| 10 00    | Performance part 1             |
| 11 00    | Performance part 2 *1-2-2      |
| :        |                                |
| 1F 00    | Performance part 16            |

#### *1-2-1 Performance common

| Offset   | Size Binary    | Description                  | Value                              |
| -------- | -------------- | ---------------------------- | ---------------------------------- |
| 00 00    | 0aaa aaaa      | Performance name 1           | 32 - 127                           |
| :        |                |                              |                                    |
| 00 0B    | 0aaa aaaa      | Performance name 12          | 32 - 127                           |
| 00 0C    | 0000 aaaa      | EFX:Source                   | 0 - 15 (PERFORM, 1 - 9, 11 - 16)   |
| 00 0D    | 000a aaaa      | EFX:Type                     | 0 - 39                             |
| 00 0E    | 0aaa aaaa      | EFX:Parameter 1              | 0 - 127                            |
| :        |                |                              |                                    |
| 00 19    | 0aaa aaaa      | EFX:Parameter 12             | 0 - 127                            |
| 00 1A    | 0000 00aa      | EFX:Output assign            | 0 - 2 (MIX, OUTPUT1, OUTPUT2)      |
| 00 1B    | 0aaa aaaa      | EFX:Output level             | 0 - 127                            |
| 00 1C    | 0aaa aaaa      | EFX:Chorus send level        | 0 - 127                            |
| 00 1D    | 0aaa aaaa      | EFX:Reverb send level        | 0 - 127                            |
| 00 1E    | 0000 aaaa      | EFX:Control source 1         | 0 - 10 *S                          |
| 00 1F    | 0aaa aaaa      | EFX:Control depth 1          | 0 - 126 (-63 - +63)                |
| 00 20    | 0000 aaaa      | EFX:Control source 2         | 0 - 10 *S                          |
| 00 21    | 0aaa aaaa      | EFX:Control depth 2          | 0 - 126 (-63 - +63)                |
| 00 22    | 0aaa aaaa      | Chorus:Level                 | 0 - 127                            |
| 00 23    | 0aaa aaaa      | Chorus:Rate                  | 0 - 127                            |
| 00 24    | 0aaa aaaa      | Chorus:Depth                 | 0 - 127                            |
| 00 25    | 0aaa aaaa      | Chorus:Pre delay             | 0 - 127                            |
| 00 26    | 0aaa aaaa      | Chorus:Feedback              | 0 - 127                            |
| 00 27    | 0000 00aa      | Chorus:Output assign         | 0 - 2 (MIX, REV, MIX+REV)          |
| 00 28    | 0000 0aaa      | Reverb:type                  | 0 - 7 *R                           |
| 00 29    | 0aaa aaaa      | Reverb:Level                 | 0 - 127                            |
| 00 2A    | 0aaa aaaa      | Reverb:Time                  | 0 - 127                            |
| 00 2B    | 0aaa aaaa      | Reverb:HF damp               | 0 - 17 *D                          |
| 00 2C    | 0aaa aaaa      | Reverb:Feedback              | 0 - 127                            |
| 00 2D    | 0000 aaaa #    | Default tempo                | 20 - 250                           |
| 00 2F    | 0000 000a      | Key range switch             | 0 - 1 (OFF, ON)                    |
| 00 30    | 0aaa aaaa      | Part 1 voice reserve         | 0 - 64                             |
| :        |                |                              |                                    |
| 00 3F    | 0aaa aaaa      | Part 16 voice reserve        | 0 - 64                             |

**Total size: 00 00 00 40**

**Note:** The performance name data returned in response to this request are expressed in ASCII characters of hexadecimal.

**Note:** The sum of voice reserves must be less than or equal to 64.

#### Examples

**Example using RQ1:** To get the name of the user performance USER:01, send the following message to the JV-1080.

`F0H 41H 10H 6AH 11H 10H 00H 00H 00H 00H 00H 00H 0CH 64H F7H` (6)

**Example using DT1:** To set the Reverb type HALL2 of the user performance USER:08, send the following message to the JV-1080.

`F0H 41H 10H 6AH 12H 10H 07H 00H 28H 05H 3CH F7H`

---

## Seite 152 — Performance Part (*1-2-2) · Patch (*1-3, *1-3-1) · Examples

#### *1-2-2 Performance part

| Offset   | Size Binary   | Description                | Value                                       |
| -------- | ------------- | -------------------------- | ------------------------------------------- |
| 00 00    | 0000 000a     | MIDI receive switch        | 0 - 1 (OFF, ON)                             |
| 00 01    | 0000 aaaa     | MIDI channel               | 0 - 15 (1 - 16)                             |
| 00 02    | 0000 00aa     | Patch group                | 0 - 2 (USER, PCM, EXP) (1)                  |
| 00 03    | 0aaa aaaa     | Patch group ID             | 0 - 127                                     |
| 00 04    | 0000 aaaa #   | Patch number               | 0 - 254 (001 - 255)                         |
| 00 06    | 0aaa aaaa     | Part level                 | 0 - 127                                     |
| 00 07    | 0aaa aaaa     | Part pan                   | 0 - 127 (L64 - 63R)                         |
| 00 08    | 0aaa aaaa     | Pitch coarse tune          | 0 - 96 (-48 - +48)                          |
| 00 09    | 0aaa aaaa     | Pitch fine tune            | 0 - 100 (-50 - +50)                         |
| 00 0A    | 0000 00aa     | Output assign              | 0 - 4 (MIX, EFX, OUTPUT1, OUTPUT2, PATCH)   |
| 00 0B    | 0aaa aaaa     | Output level               | 0 - 127                                     |
| 00 0C    | 0aaa aaaa     | Chorus send level          | 0 - 127                                     |
| 00 0D    | 0aaa aaaa     | Reverb send level          | 0 - 127                                     |
| 00 0E    | 0000 000a     | Receive program change     | 0 - 1 (OFF, ON)                             |
| 00 0F    | 0000 000a     | Receive volume             | 0 - 1 (OFF, ON)                             |
| 00 10    | 0000 000a     | Receive hold-1             | 0 - 1 (OFF, ON)                             |
| 00 11    | 0aaa aaaa     | Key range lower            | 0 - 127 (C-1 - G9)                          |
| 00 12    | 0aaa aaaa     | Key range upper            | 0 - 127 (C-1 - G9)                          |

**Total size: 00 00 00 13**

#### Examples

**Example using RQ1:** To get the data of the performance part 3 of the user performance USER:03, send the following message to the JV-1080.

`F0H 41H 10H 6AH 11H 10H 02H 12H 00H 00H 00H 00H 13H 49H F7H`

**Example using DT1:** To set the MIDI receive switch of the performance part 1 of the temporary performance to OFF, send the following message to the JV-1080.

`F0H 41H 10H 6AH 12H 01H 00H 10H 00H 00H 6FH F7H`

#### *1-3 Patch

| Offset   | Description             |
| -------- | ----------------------- |
| 00 00    | Patch common *1-3-1     |
| 10 00    | Patch tone 1 *1-3-2     |
| 12 00    | Patch tone 2 *1-3-2     |
| 14 00    | Patch tone 3 *1-3-2     |
| 16 00    | Patch tone 4 *1-3-2     |

#### *1-3-1 Patch common

| Offset   | Size Binary    | Description                  | Value                            |
| -------- | -------------- | ---------------------------- | -------------------------------- |
| 00 00    | 0aaa aaaa      | Patch name 1                 | 32 - 127                         |
| :        |                |                              |                                  |
| 00 0B    | 0aaa aaaa      | Patch name 12                | 32 - 127                         |
| 00 0C    | 000a aaaa      | EFX:Type                     | 0 - 39                           |
| 00 0D    | 0aaa aaaa      | EFX:Parameter 1              | 0 - 127                          |
| :        |                |                              |                                  |
| 00 18    | 0aaa aaaa      | EFX:Parameter 12             | 0 - 127                          |
| 00 19    | 0000 00aa      | EFX:Output assign            | 0 - 2 (MIX, OUTPUT1, OUTPUT2)    |
| 00 1A    | 0aaa aaaa      | EFX:Output level             | 0 - 127                          |
| 00 1B    | 0aaa aaaa      | EFX:Chorus send level        | 0 - 127                          |
| 00 1C    | 0aaa aaaa      | EFX:Reverb send level        | 0 - 127                          |
| 00 1D    | 0000 aaaa      | EFX:Control source 1         | 0 - 10 *S                        |
| 00 1E    | 0aaa aaaa      | EFX:Control depth 1          | 0 - 126 (-63 - +63)              |
| 00 1F    | 0000 aaaa      | EFX:Control source 2         | 0 - 10 *S                        |
| 00 20    | 0aaa aaaa      | EFX:Control depth 2          | 0 - 126 (-63 - +63)              |
| 00 21    | 0aaa aaaa      | Chorus:Level                 | 0 - 127                          |
| 00 22    | 0aaa aaaa      | Chorus:Rate                  | 0 - 127                          |
| 00 23    | 0aaa aaaa      | Chorus:Depth                 | 0 - 127                          |
| 00 24    | 0aaa aaaa      | Chorus:Pre delay             | 0 - 127                          |
| 00 25    | 0aaa aaaa      | Chorus:Feedback              | 0 - 127                          |
| 00 26    | 0000 00aa      | Chorus:Output assign         | 0 - 2 (MIX, REV, MIX+REV)        |
| 00 27    | 0000 0aaa      | Reverb:type                  | 0 - 7 *R                         |
| 00 28    | 0aaa aaaa      | Reverb:Level                 | 0 - 127                          |
| 00 29    | 0aaa aaaa      | Reverb:Time                  | 0 - 127                          |
| 00 2A    | 0aaa aaaa      | Reverb:HF damp               | 0 - 17 *D                        |
| 00 2B    | 0aaa aaaa      | Reverb:Feedback              | 0 - 127                          |
| 00 2C    | 0000 aaaa #    | Default tempo                | 20 - 250                         |
| 00 2E    | 0aaa aaaa      | Patch level                  | 0 - 127                          |
| 00 2F    | 0aaa aaaa      | Patch pan                    | 0 - 127 (L64 - 63R)              |
| 00 30    | 0aaa aaaa      | Analog feel depth            | 0 - 127                          |
| 00 31    | 0000 aaaa      | Bender range up              | 0 - 12                           |
| 00 32    | 0aaa aaaa      | Bender range down            | 0 - 48 (0 - -48)                 |
| 00 33    | 0000 000a      | Key assign mode              | 0 - 1 (POLY, SOLO)               |
| 00 34    | 0000 000a      | Solo legato                  | 0 - 1 (OFF, ON)                  |
| 00 35    | 0000 000a      | Portamento switch            | 0 - 1 (OFF, ON)                  |
| 00 36    | 0000 000a      | Portamento mode              | 0 - 1 (NORMAL, LEGATO)           |
| 00 37    | 0000 000a      | Portamento type              | 0 - 1 (RATE, TIME)               |
| 00 38    | 0000 000a      | Portamento start             | 0 - 1 (PITCH, NOTE)              |
| 00 39    | 0aaa aaaa      | Portamento time              | 0 - 127                          |
| 00 3A    | 0000 aaaa      | Patch control source 2       | 0 - 15 *P                        |
| 00 3B    | 0000 aaaa      | Patch control source 3       | 0 - 15 *P                        |
| 00 3C    | 0000 00aa      | EFX control hold/peak        | 0 - 2 (OFF, HOLD, PEAK)          |
| 00 3D    | 0000 00aa      | Control 1 hold/peak          | 0 - 2 (OFF, HOLD, PEAK)          |
| 00 3E    | 0000 00aa      | Control 2 hold/peak          | 0 - 2 (OFF, HOLD, PEAK)          |
| 00 3F    | 0000 00aa      | Control 3 hold/peak          | 0 - 2 (OFF, HOLD, PEAK)          |
| 00 40    | 0000 000a      | Velocity range switch        | 0 - 1 (OFF, ON)                  |
| 00 41    | 0000 0aaa      | Octave shift                 | 0 - 6 (-3 - +3)                  |
| 00 42    | 0000 00aa      | Stretch tune depth           | 0 - 3 (OFF, 1 - 3)               |
| 00 43    | 0000 000a      | Voice priority               | 0 - 1 (LAST, LOUDEST)            |
| 00 44    | 0000 aaaa      | Structure type 1 & 2         | 0 - 9 (1 - 10)                   |
| 00 45    | 0000 00aa      | Booster level 1 & 2          | 0 - 3 (0, -6, +12, +18)          |
| 00 46    | 0000 aaaa      | Structure type 3 & 4         | 0 - 9 (1 - 10)                   |
| 00 47    | 0000 00aa      | Booster level 3 & 4          | 0 - 3 (0, -6, +12, +18)          |

**Total size: 00 00 00 48**

#### Examples

**Example using RQ1:** To get the portamento time of the patch mode temporary patch, send the following message to the JV-1080.

`F0H 41H 10H 6AH 11H 03H 00H 00H 39H 00H 00H 00H 01H 43H F7H`

**Example using DT1:** To set the structure type 1 & 2 of the user patch USER:048 to 3, send the following message to the JV-1080.

`F0H 41H 10H 6AH 12H 11H 2FH 00H 44H 02H 7AH F7H`

---

## Seite 153 — Patch Tone (*1-3-2) (10)

#### *1-3-2 Patch tone

| Offset   | Size Binary    | Description                        | Value                                                    |
| -------- | -------------- | ---------------------------------- | -------------------------------------------------------- |
| 00 00    | 0000 000a      | Tone switch                        | 0 - 1 (OFF, ON)                                          |
| 00 01    | 0000 00aa      | Wave group                         | 0 - 2 (INT, PCM, EXP)                                    |
| 00 02    | 0aaa aaaa      | Wave group ID                      | 0 - 127                                                  |
| 00 03    | 0000 aaaa #    | Wave number                        | 0 - 254 (1 - 255)                                        |
| 00 05    | 0000 00aa      | Wave gain                          | 0 - 3 (-6, 0, +6, +12)                                   |
| 00 06    | 0000 000a      | FXM switch                         | 0 - 1 (OFF, ON)                                          |
| 00 07    | 0000 00aa      | FXM color                          | 0 - 3 (1 - 4)                                            |
| 00 08    | 0000 aaaa      | FXM depth                          | 0 - 15 (1 - 16)                                          |
| 00 09    | 0000 0aaa      | Tone delay mode                    | 0 - 6 *M                                                 |
| 00 0A    | 0aaa aaaa      | Tone delay time                    | 0 - 127                                                  |
| 00 0B    | 0aaa aaaa      | Velocity cross fade                | 0 - 127                                                  |
| 00 0C    | 0aaa aaaa      | Velocity range lower               | 1 - 127                                                  |
| 00 0D    | 0aaa aaaa      | Velocity range upper               | 1 - 127                                                  |
| 00 0E    | 0aaa aaaa      | Keyboard range lower               | 0 - 127 (C-1 - G9)                                       |
| 00 0F    | 0aaa aaaa      | Keyboard range upper               | 0 - 127 (C-1 - G9)                                       |
| 00 10    | 0000 000a      | Redamper control switch            | 0 - 1 (OFF, ON)                                          |
| 00 11    | 0000 000a      | Volume control switch              | 0 - 1 (OFF, ON)                                          |
| 00 12    | 0000 000a      | Hold-1 control switch              | 0 - 1 (OFF, ON)                                          |
| 00 13    | 0000 000a      | Bender control switch              | 0 - 1 (OFF, ON)                                          |
| 00 14    | 0000 00aa      | Pan control switch                 | 0 - 2 (OFF, CONTINUOUS, KEY-ON)                          |
| 00 15    | 000a aaaa      | Controller 1 destination 1         | 0 - 18 *C                                                |
| 00 16    | 0aaa aaaa      | Controller 1 depth 1               | 0 - 126 (-63 - +63)                                      |
| 00 17    | 000a aaaa      | Controller 1 destination 2         | 0 - 18 *C                                                |
| 00 18    | 0aaa aaaa      | Controller 1 depth 2               | 0 - 126 (-63 - +63)                                      |
| 00 19    | 000a aaaa      | Controller 1 destination 3         | 0 - 18 *C                                                |
| 00 1A    | 0aaa aaaa      | Controller 1 depth 3               | 0 - 126 (-63 - +63)                                      |
| 00 1B    | 000a aaaa      | Controller 1 destination 4         | 0 - 18 *C                                                |
| 00 1C    | 0aaa aaaa      | Controller 1 depth 4               | 0 - 126 (-63 - +63)                                      |
| 00 1D    | 000a aaaa      | Controller 2 destination 1         | 0 - 18 *C                                                |
| 00 1E    | 0aaa aaaa      | Controller 2 depth 1               | 0 - 126 (-63 - +63)                                      |
| 00 1F    | 000a aaaa      | Controller 2 destination 2         | 0 - 18 *C                                                |
| 00 20    | 0aaa aaaa      | Controller 2 depth 2               | 0 - 126 (-63 - +63)                                      |
| 00 21    | 000a aaaa      | Controller 2 destination 3         | 0 - 18 *C                                                |
| 00 22    | 0aaa aaaa      | Controller 2 depth 3               | 0 - 126 (-63 - +63)                                      |
| 00 23    | 000a aaaa      | Controller 2 destination 4         | 0 - 18 *C                                                |
| 00 24    | 0aaa aaaa      | Controller 2 depth 4               | 0 - 126 (-63 - +63)                                      |
| 00 25    | 000a aaaa      | Controller 3 destination 1         | 0 - 18 *C                                                |
| 00 26    | 0aaa aaaa      | Controller 3 depth 1               | 0 - 126 (-63 - +63)                                      |
| 00 27    | 000a aaaa      | Controller 3 destination 2         | 0 - 18 *C                                                |
| 00 28    | 0aaa aaaa      | Controller 3 depth 2               | 0 - 126 (-63 - +63)                                      |
| 00 29    | 000a aaaa      | Controller 3 destination 3         | 0 - 18 *C                                                |
| 00 2A    | 0aaa aaaa      | Controller 3 depth 3               | 0 - 126 (-63 - +63)                                      |
| 00 2B    | 000a aaaa      | Controller 3 destination 4         | 0 - 18 *C                                                |
| 00 2C    | 0aaa aaaa      | Controller 3 depth 4               | 0 - 126 (-63 - +63)                                      |
| 00 2D    | 0000 0aaa      | LFO1 waveform                      | 0 - 7 (TRI, SIN, SAW, SQR, TRP, S&H, RND, CHS)           |
| 00 2E    | 0000 000a      | LFO1 key sync                      | 0 - 1 (OFF, ON)                                          |
| 00 2F    | 0aaa aaaa      | LFO1 rate                          | 0 - 127                                                  |
| 00 30    | 0000 00aa      | LFO1 offset                        | 0 - 4 (-100, -50, 0, +50, +100)                          |
| 00 31    | 0aaa aaaa      | LFO1 delay time                    | 0 - 127                                                  |
| 00 32    | 0000 00aa      | LFO1 fade mode                     | 0 - 3 (KEY-ON-IN, KEY-ON-OUT, KEY-OFF-IN, KEY-OFF-OUT)   |
| 00 33    | 0aaa aaaa      | LFO1 fade time                     | 0 - 127                                                  |
| 00 34    | 0000 00aa      | LFO1 external sync                 | 0 - 2 (OFF, CLOCK, TAP)                                  |
| 00 35    | 0000 0aaa      | LFO2 waveform                      | 0 - 7 (TRI, SIN, SAW, SQR, TRP, S&H, RND, CHS)           |
| 00 36    | 0000 000a      | LFO2 key sync                      | 0 - 1 (OFF, ON)                                          |
| 00 37    | 0aaa aaaa      | LFO2 rate                          | 0 - 127                                                  |
| 00 38    | 0000 00aa      | LFO2 offset                        | 0 - 4 (-100, -50, 0, +50, +100)                          |
| 00 39    | 0aaa aaaa      | LFO2 delay time                    | 0 - 127                                                  |
| 00 3A    | 0000 00aa      | LFO2 fade mode                     | 0 - 3 (KEY-ON-IN, KEY-ON-OUT, KEY-OFF-IN, KEY-OFF-OUT)   |
| 00 3B    | 0aaa aaaa      | LFO2 fade time                     | 0 - 127                                                  |
| 00 3C    | 0000 00aa      | LFO2 external sync                 | 0 - 2 (OFF, CLOCK, TAP)                                  |
| 00 3D    | 0aaa aaaa      | Pitch coarse tune                  | 0 - 96 (-48 - +48)                                       |
| 00 3E    | 0aaa aaaa      | Pitch fine tune                    | 0 - 100 (-50 - +50)                                      |
| 00 3F    | 000a aaaa      | Random pitch depth                 | 0 - 30 *RP                                               |
| 00 40    | 0000 aaaa      | Pitch keyfollow                    | 0 - 15 *A                                                |
| 00 41    | 000a aaaa      | Pitch envelope depth               | 0 - 24 (-12 - +12)                                       |
| 00 42    | 0aaa aaaa      | Pitch env velocity sensitivity     | 0 - 125 *V                                               |
| 00 43    | 0000 aaaa      | Pitch env velocity time 1 sens.    | 0 - 14 *B                                                |
| 00 44    | 0000 aaaa      | Pitch env velocity time 4 sens.    | 0 - 14 *B                                                |
| 00 45    | 0000 aaaa      | Pitch env time keyfollow           | 0 - 14 *B                                                |
| 00 46    | 0aaa aaaa      | Pitch env time 1                   | 0 - 127                                                  |
| 00 47    | 0aaa aaaa      | Pitch env time 2                   | 0 - 127                                                  |
| 00 48    | 0aaa aaaa      | Pitch env time 3                   | 0 - 127                                                  |
| 00 49    | 0aaa aaaa      | Pitch env time 4                   | 0 - 127                                                  |
| 00 4A    | 0aaa aaaa      | Pitch env level 1                  | 0 - 126 (-63 - +63)                                      |
| 00 4B    | 0aaa aaaa      | Pitch env level 2                  | 0 - 126 (-63 - +63)                                      |
| 00 4C    | 0aaa aaaa      | Pitch env level 3                  | 0 - 126 (-63 - +63)                                      |
| 00 4D    | 0aaa aaaa      | Pitch env level 4                  | 0 - 126 (-63 - +63)                                      |
| 00 4E    | 0aaa aaaa      | Pitch LFO1 depth                   | 0 - 126 (-63 - +63)                                      |
| 00 4F    | 0aaa aaaa      | Pitch LFO2 depth                   | 0 - 126 (-63 - +63)                                      |
| 00 50    | 0000 0aaa      | Filter type                        | 0 - 4 (OFF, LPF, BPF, HPF, PKG)                          |
| 00 51    | 0aaa aaaa      | Cutoff frequency                   | 0 - 127                                                  |
| 00 52    | 0000 aaaa      | Cutoff keyfollow                   | 0 - 15 *A                                                |
| 00 53    | 0aaa aaaa      | Resonance                          | 0 - 127                                                  |
| 00 54    | 0aaa aaaa      | Resonance velocity sensitivity     | 0 - 125 *V                                               |
| 00 55    | 0aaa aaaa      | Filter env depth                   | 0 - 126 (-63 - +63)                                      |
| 00 56    | 0000 0aaa      | Filter env velocity curve          | 0 - 6 (1 - 7)                                            |
| 00 57    | 0aaa aaaa      | Filter env velocity sensitivity    | 0 - 125 *V                                               |
| 00 58    | 0000 aaaa      | Filter env velocity time 1 sens.   | 0 - 14 *B                                                |
| 00 59    | 0000 aaaa      | Filter env velocity time 4 sens.   | 0 - 14 *B                                                |
| 00 5A    | 0000 aaaa      | Filter env time keyfollow          | 0 - 14 *B                                                |
| 00 5B    | 0aaa aaaa      | Filter env time 1                  | 0 - 127                                                  |
| 00 5C    | 0aaa aaaa      | Filter env time 2                  | 0 - 127                                                  |
| 00 5D    | 0aaa aaaa      | Filter env time 3                  | 0 - 127                                                  |
| 00 5E    | 0aaa aaaa      | Filter env time 4                  | 0 - 127                                                  |
| 00 5F    | 0aaa aaaa      | Filter env level 1                 | 0 - 127                                                  |
| 00 60    | 0aaa aaaa      | Filter env level 2                 | 0 - 127                                                  |
| 00 61    | 0aaa aaaa      | Filter env level 3                 | 0 - 127                                                  |
| 00 62    | 0aaa aaaa      | Filter env level 4                 | 0 - 127                                                  |
| 00 63    | 0aaa aaaa      | Filter LFO1 depth                  | 0 - 126 (-63 - +63)                                      |
| 00 64    | 0aaa aaaa      | Filter LFO2 depth                  | 0 - 126 (-63 - +63)                                      |
| 00 65    | 0aaa aaaa      | Tone level                         | 0 - 127                                                  |
| 00 66    | 0000 00aa      | Bias direction                     | 0 - 3 (LOWER, UPPER, LOWER&UPPER, ALL)                   |
| 00 67    | 0aaa aaaa      | Bias position                      | 0 - 127 (C-1 - G9)                                       |
| 00 68    | 0000 aaaa      | Bias level                         | 0 - 14 *B                                                |
| 00 69    | 0000 0aaa      | Level env velocity curve           | 0 - 6 (1 - 7)                                            |
| 00 6A    | 0aaa aaaa      | Level env velocity sensitivity     | 0 - 125 *V                                               |
| 00 6B    | 0000 aaaa      | Level env velocity time 1 sens.    | 0 - 14 *B                                                |
| 00 6C    | 0000 aaaa      | Level env velocity time 4 sens.    | 0 - 14 *B                                                |
| 00 6D    | 0000 aaaa      | Level env time keyfollow           | 0 - 14 *B                                                |
| 00 6E    | 0aaa aaaa      | Level env time 1                   | 0 - 127                                                  |
| 00 6F    | 0aaa aaaa      | Level env time 2                   | 0 - 127                                                  |
| 00 70    | 0aaa aaaa      | Level env time 3                   | 0 - 127                                                  |
| 00 71    | 0aaa aaaa      | Level env time 4                   | 0 - 127                                                  |
| 00 72    | 0aaa aaaa      | Level env level 1                  | 0 - 127                                                  |
| 00 73    | 0aaa aaaa      | Level env level 2                  | 0 - 127                                                  |
| 00 74    | 0aaa aaaa      | Level env level 3                  | 0 - 127                                                  |
| 00 75    | 0aaa aaaa      | Amp LFO1 depth                     | 0 - 126 (-63 - +63)                                      |
| 00 76    | 0aaa aaaa      | Amp LFO2 depth                     | 0 - 126 (-63 - +63)                                      |
| 00 77    | 0aaa aaaa      | Tone pan                           | 0 - 127 (L64 - 63R)                                      |
| 00 78    | 0000 aaaa      | Pan keyfollow                      | 0 - 14 *B                                                |
| 00 79    | 0aaa aaaa      | Random pan depth                   | 0 - 63                                                   |
| 00 7A    | 0aaa aaaa      | Alternate pan depth                | 1 - 127 (L63 - 63R)                                      |
| 00 7B    | 0aaa aaaa      | Pan LFO1 depth                     | 0 - 126 (L63 - 63R)                                      |
| 00 7C    | 0aaa aaaa      | Pan LFO2 depth                     | 0 - 126 (L63 - 63R)                                      |
| 00 7D    | 0000 00aa      | Output assign                      | 0 - 3 (MIX, EFX, OUTPUT1, OUTPUT2)                       |
| 00 7E    | 0aaa aaaa      | Output level                       | 0 - 127                                                  |
| 00 7F    | 0aaa aaaa      | Chorus send level                  | 0 - 127                                                  |
| 01 00    | 0aaa aaaa      | Reverb send level                  | 0 - 127                                                  |

**Total size: 00 00 01 01**

**Note:** If the value of the wave number surpasses the number of waves contained in the corresponding wave group, this message will be ignored.

**Note:** If the value of the velocity range lower is greater than that of the velocity range upper, this message will be ignored.

---

## Seite 154 — Examples · Rhythm Setup (*1-4, *1-4-1, *1-4-2) · Adress-Übersicht

#### Examples

**Example using RQ1:** To get the data of the patch tone 2 of the user patch USER:02, send the following message to the JV-1080.

`F0H 41H 10H 6AH 11H 11H 01H 12H 00H 00H 00H 01H 01H 5AH F7H`

**Example using DT1:** To set the cutoff frequency of the patch tone 3 of the patch mode temporary patch to 100, send the following message to the JV-1080.

`F0H 41H 10H 6AH 12H 03H 00H 14H 51H 64H 34H F7H`

(14H 51H = Cutoff frequency, 64H = 100)

#### *1-4 Rhythm setup

| Offset   | Description                                |
| -------- | ------------------------------------------ |
| 00 00    | Rhythm common *1-4-1                       |
| 23 00    | Rhythm note for key# 35 (B1) *1-4-2        |
| 24 00    | Rhythm note for key# 36 (C2) *1-4-2        |
| :        |                                            |
| 62 00    | Rhythm note for key# 98 (D7) *1-4-2 (8)    |

#### *1-4-1 Rhythm common

| Offset   | Size Binary   | Description           | Value        |
| -------- | ------------- | --------------------- | ------------ |
| 00 00    | 0aaa aaaa     | Rhythm name 1         | 32 - 127     |
| :        |               |                       |              |
| 00 0B    | 0aaa aaaa     | Rhythm name 12        | 32 - 127     |

**Total size: 00 00 00 0C**

#### *1-4-2 Rhythm note (5)

| Offset   | Size Binary   | Description                            | Value                                |
| -------- | ------------- | -------------------------------------- | ------------------------------------ |
| 00 00    | 0000 000a     | Tone switch                            | 0 - 1 (OFF, ON)                      |
| 00 01    | 0000 00aa     | Wave group                             | 0 - 2 (INT, PCM, EXP)                |
| 00 02    | 0aaa aaaa     | Wave group ID                          | 0 - 127                              |
| 00 03    | 0000 aaaa #   | Wave number                            | 0 - 254 (1 - 255)                    |
| 00 05    | 0000 00aa     | Wave gain                              | 0 - 3 (-6, 0, +6, +12)               |
| 00 06    | 0000 aaaa     | Bend range                             | 0 - 12                               |
| 00 07    | 0000 aaaa     | Mute group                             | 0 - 31 (OFF, 1 - 31)                 |
| 00 08    | 0000 000a     | Envelope mode                          | 0 - 1 (NO-SUSTAIN, SUSTAIN)          |
| 00 09    | 0000 000a     | Volume control switch                  | 0 - 1 (OFF, ON)                      |
| 00 0A    | 0000 000a     | Hold-1 control switch                  | 0 - 1 (OFF, ON)                      |
| 00 0B    | 0000 00aa     | Pan control switch                     | 0 - 2 (OFF, CONTINUOUS, KEY-ON)      |
| 00 0C    | 0aaa aaaa     | Coarse tune                            | 0 - 127 (C-1 - G9) (5)               |
| 00 0D    | 0aaa aaaa     | Fine tune                              | 0 - 100 (-50 - +50)                  |
| 00 0E    | 000a aaaa     | Random pitch depth                     | 0 - 30 *RP                           |
| 00 0F    | 000a aaaa     | Pitch env depth                        | 0 - 24 (-12 - +12)                   |
| 00 10    | 0aaa aaaa     | Pitch env velocity sensitivity         | 0 - 125 *V                           |
| 00 11    | 0000 aaaa     | Pitch env velocity time sensitivity    | 0 - 14 *B                            |
| 00 12    | 0aaa aaaa     | Pitch env time 1                       | 0 - 127                              |
| 00 13    | 0aaa aaaa     | Pitch env time 2                       | 0 - 127                              |
| 00 14    | 0aaa aaaa     | Pitch env time 3                       | 0 - 127                              |
| 00 15    | 0aaa aaaa     | Pitch env time 4                       | 0 - 127                              |
| 00 16    | 0aaa aaaa     | Pitch env level 1                      | 0 - 126 (-63 - +63)                  |
| 00 17    | 0aaa aaaa     | Pitch env level 2                      | 0 - 126 (-63 - +63)                  |
| 00 18    | 0aaa aaaa     | Pitch env level 3                      | 0 - 126 (-63 - +63)                  |
| 00 19    | 0aaa aaaa     | Pitch env level 4                      | 0 - 126 (-63 - +63)                  |
| 00 1A    | 0000 0aaa     | Filter type                            | 0 - 4 (OFF, LPF, BPF, HPF, PKG)      |
| 00 1B    | 0aaa aaaa     | Cutoff frequency                       | 0 - 127                              |
| 00 1C    | 0aaa aaaa     | Resonance                              | 0 - 127                              |
| 00 1D    | 0aaa aaaa     | Resonance velocity sensitivity         | 0 - 125 *V                           |
| 00 1E    | 0aaa aaaa     | Filter env depth                       | 0 - 126 (-63 - +63)                  |
| 00 1F    | 0aaa aaaa     | Filter env velocity sensitivity        | 0 - 125 *V                           |
| 00 20    | 0000 aaaa     | Filter env velocity time sensitivity   | 0 - 14 *B                            |
| 00 21    | 0aaa aaaa     | Filter env time 1                      | 0 - 127                              |
| 00 22    | 0aaa aaaa     | Filter env time 2                      | 0 - 127                              |
| 00 23    | 0aaa aaaa     | Filter env time 3                      | 0 - 127                              |
| 00 24    | 0aaa aaaa     | Filter env time 4                      | 0 - 127                              |
| 00 25    | 0aaa aaaa     | Filter env level 1                     | 0 - 127                              |
| 00 26    | 0aaa aaaa     | Filter env level 2                     | 0 - 127                              |
| 00 27    | 0aaa aaaa     | Filter env level 3                     | 0 - 127                              |
| 00 28    | 0aaa aaaa     | Filter env level 4                     | 0 - 127                              |
| 00 29    | 0aaa aaaa     | Tone level                             | 0 - 127                              |
| 00 2A    | 0aaa aaaa     | Level env velocity sensitivity         | 0 - 125 *V                           |
| 00 2B    | 0000 aaaa     | Level env velocity time sensitivity    | 0 - 14 *B                            |
| 00 2C    | 0aaa aaaa     | Level env time 1                       | 0 - 127                              |
| 00 2D    | 0aaa aaaa     | Level env time 2                       | 0 - 127                              |
| 00 2E    | 0aaa aaaa     | Level env time 3                       | 0 - 127                              |
| 00 2F    | 0aaa aaaa     | Level env time 4                       | 0 - 127                              |
| 00 30    | 0aaa aaaa     | Level env level 1                      | 0 - 127                              |
| 00 31    | 0aaa aaaa     | Level env level 2                      | 0 - 127                              |
| 00 32    | 0aaa aaaa     | Level env level 3                      | 0 - 127                              |
| 00 33    | 0aaa aaaa     | Tone pan                               | 0 - 127 (L64 - 63R)                  |
| 00 34    | 0aaa aaaa     | Random pan depth                       | 0 - 63                               |
| 00 35    | 0aaa aaaa     | Alternate pan depth                    | 1 - 127 (L63 - 63R)                  |
| 00 36    | 0000 00aa     | Output assign                          | 0 - 3 (MIX, EFX, OUTPUT1, OUTPUT2)   |
| 00 37    | 0aaa aaaa     | Output level                           | 0 - 127                              |
| 00 38    | 0aaa aaaa     | Chorus send level                      | 0 - 127                              |
| 00 39    | 0aaa aaaa     | Reverb send level                      | 0 - 127                              |

**Total size: 00 00 00 3A**

**Note:** If the value of the wave number surpasses the number of waves contained in the corresponding wave group, this message will be ignored.

#### Examples

**Example using RQ1:** To get all the data of the rhythm note C2 of the temporary rhythm setup, send the following message to the JV-1080.

`F0H 41H 10H 6AH 11H 02H 09H 24H 00H 00H 00H 00H 3AH 17H F7H`

**Example using DT1:** To set the tone switch of the rhythm note D2 of the temporary rhythm setup to OFF, send the following message to the JV-1080. (6)

`F0H 41H 10H 6AH 12H 02H 09H 26H 00H 00H 4FH F7H`

### ■ Address map — Übersicht (Block / Sub block / Reference)

| Start address   | Block                              | Sub block                                   | Reference         |         |
| --------------- | ---------------------------------- | ------------------------------------------- | ----------------- | ------- |
| 00 00 00 00     | System                             | Common                                      | 1-1-1             |         |
|                 |                                    | Scale tune (Part 1 - Part 16, Patch mode)   | 1-1-2             |         |
| 01 00 00 00     | Temporary performance              | Common                                      | 1-2-1             |         |
|                 |                                    | Part 1 - Part 16                            | 1-2-2             |         |
| 02 00 00 00     | Performance mode temporary patch   | Part 1 - Part 9                             | Common            |         |
|                 |                                    |                                             | Tone 1 - Tone 4   | 1-3-2   |
| 02 09 00 00     | Temporary rhythm setup             | Common                                      | 1-4-1             |         |
|                 |                                    | Note# 35 - Note# 98                         | 1-4-2             |         |
| 02 0A 00 00     | Performance mode temporary patch   | Part 11 - Part 16                           | Common            |         |
|                 |                                    |                                             | Tone 1 - Tone 4   | 1-3-2   |
| 03 00 00 00     | Patch mode temporary patch         | Common                                      | 1-3-1             |         |
|                 |                                    | Tone 1 - Tone 4                             | 1-3-2             |         |

(Fortsetzung auf Seite 155; Part 1 - 9 bzw. 11 - 16 haben jeweils die Struktur Common (*1-3-1) + Tone 1 - 4 (*1-3-2))

---

## Seite 155 — Adress-Übersicht (Fortsetzung) · GS (MODEL ID = 42H) · Tabellen A-1 / A-2

### ■ Address map — Übersicht (Fortsetzung)

| Start address   | Block                          | Sub block                | Reference             |         |
| --------------- | ------------------------------ | ------------------------ | --------------------- | ------- |
| 10 00 00 00     | User performance               | USER:01 - USER:32        | Common 1-2-1          |         |
|                 |                                |                          | Part 1 - Part 16      | 1-2-2   |
| 10 40 00 00     | User rhythm setup              | USER:1 - USER:2          | Common 1-4-1          |         |
|                 |                                |                          | Note# 35 - Note# 98   | 1-4-2   |
| 11 00 00 00     | User patch                     | USER:001 - USER:128      | Common 1-3-1          |         |
|                 |                                |                          | Tone 1 - Tone 4       | 1-3-2   |
| 20 00 00 00     | Data card performance          | CARD:01 - CARD:32        | Common 1-2-1          |         |
|                 |                                |                          | Part 1 - Part 16      | 1-2-2   |
| 20 40 00 00     | Data card rhythm setup         | CARD:1 - CARD:2          | Common 1-4-1          |         |
|                 |                                |                          | Note# 35 - Note# 98   | 1-4-2   |
| 21 00 00 00     | Data card patch                | CARD:001 - CARD:128      | Common 1-3-1          |         |
|                 |                                |                          | Tone 1 - Tone 4       | 1-3-2   |

### 2. GS (MODEL ID = 42H)

| Start address   | Description          |
| --------------- | -------------------- |
| 40 10 00        | Scale tune Part 10   |
| 40 11 00        | Scale tune Part 1    |
| 40 12 00        | Scale tune Part 2    |
| :               |                      |
| 40 19 00        | Scale tune Part 9    |
| 40 1A 00        | Scale tune Part 11   |
| :               |                      |
| 40 1F 00        | Scale tune Part 16   |

#### *2-1 Scale tune

| Offset   | Size Binary   | Description   | Value                 |
| -------- | ------------- | ------------- | --------------------- |
| 00 40    | 0aaa aaaa     | C             | 0 - 127 (-64 - +63)   |
| 00 41    | 0aaa aaaa     | C#            | 0 - 127 (-64 - +63)   |
| :        |               |               |                       |
| 00 4B    | 0aaa aaaa     | B             | 0 - 127 (-64 - +63)   |

**Total size: 00 00 00 0C**

**Note:** If you use the GS exclusive message, the start address must be the C note of each part's scale tune. For example, the scale tune C of Part 1 is 40 11 40H.

#### Example using DT1 (7)

To set the scale tune (C - B) of the performance part 1 to ARABIA, send the data as follows:

`F0H 41H 10H 42H 12H 40H 11H 40H 3AH 6DH 3EH 34H 0DH 3EH 68H 3CH 6FH 40H 36H 0FH 73H F7H`

### ● Table A-1: Decimal to Hexadecimal

MIDI uses 7-bit hexadecimal values to indicate data values and the address and size of exclusive messages. The following table shows the correspondence between decimal and hexadecimal numbers.

| D    | H     | D    | H     | D    | H     | D     | H     |
| ---- | ----- | ---- | ----- | ---- | ----- | ----- | ----- |
| 0    | 00H   | 32   | 20H   | 64   | 40H   | 96    | 60H   |
| 1    | 01H   | 33   | 21H   | 65   | 41H   | 97    | 61H   |
| 2    | 02H   | 34   | 22H   | 66   | 42H   | 98    | 62H   |
| 3    | 03H   | 35   | 23H   | 67   | 43H   | 99    | 63H   |
| 4    | 04H   | 36   | 24H   | 68   | 44H   | 100   | 64H   |
| 5    | 05H   | 37   | 25H   | 69   | 45H   | 101   | 65H   |
| 6    | 06H   | 38   | 26H   | 70   | 46H   | 102   | 66H   |
| 7    | 07H   | 39   | 27H   | 71   | 47H   | 103   | 67H   |
| 8    | 08H   | 40   | 28H   | 72   | 48H   | 104   | 68H   |
| 9    | 09H   | 41   | 29H   | 73   | 49H   | 105   | 69H   |
| 10   | 0AH   | 42   | 2AH   | 74   | 4AH   | 106   | 6AH   |
| 11   | 0BH   | 43   | 2BH   | 75   | 4BH   | 107   | 6BH   |
| 12   | 0CH   | 44   | 2CH   | 76   | 4CH   | 108   | 6CH   |
| 13   | 0DH   | 45   | 2DH   | 77   | 4DH   | 109   | 6DH   |
| 14   | 0EH   | 46   | 2EH   | 78   | 4EH   | 110   | 6EH   |
| 15   | 0FH   | 47   | 2FH   | 79   | 4FH   | 111   | 6FH   |
| 16   | 10H   | 48   | 30H   | 80   | 50H   | 112   | 70H   |
| 17   | 11H   | 49   | 31H   | 81   | 51H   | 113   | 71H   |
| 18   | 12H   | 50   | 32H   | 82   | 52H   | 114   | 72H   |
| 19   | 13H   | 51   | 33H   | 83   | 53H   | 115   | 73H   |
| 20   | 14H   | 52   | 34H   | 84   | 54H   | 116   | 74H   |
| 21   | 15H   | 53   | 35H   | 85   | 55H   | 117   | 75H   |
| 22   | 16H   | 54   | 36H   | 86   | 56H   | 118   | 76H   |
| 23   | 17H   | 55   | 37H   | 87   | 57H   | 119   | 77H   |
| 24   | 18H   | 56   | 38H   | 88   | 58H   | 120   | 78H   |
| 25   | 19H   | 57   | 39H   | 89   | 59H   | 121   | 79H   |
| 26   | 1AH   | 58   | 3AH   | 90   | 5AH   | 122   | 7AH   |
| 27   | 1BH   | 59   | 3BH   | 91   | 5BH   | 123   | 7BH   |
| 28   | 1CH   | 60   | 3CH   | 92   | 5CH   | 124   | 7CH   |
| 29   | 1DH   | 61   | 3DH   | 93   | 5DH   | 125   | 7DH   |
| 30   | 1EH   | 62   | 3EH   | 94   | 5EH   | 126   | 7EH   |
| 31   | 1FH   | 63   | 3FH   | 95   | 5FH   | 127   | 7FH   |

**Notes:**

- The decimal value of the MIDI channel, Bank Select and Program Change messages is the value in the table + 1.
- Values which use two bytes are indicated as "aa bbH", and are calculated as: aa bbH = aa x 128 + bb (e.g. the value "12 34H" is 18 x 128 + 52 = 2356). (9)
- Data within exclusive messages uses 7-bit bytes (00H - 7FH, decimal 0 - 127), so a value which is larger than 128 cannot be sent with a single byte; it is divided into two bytes.
- Signed values (-64 - +63) within exclusive messages are indicated as offset binary: 00H = -64, 40H = 0, 7FH = +63.
- Signed values (-8192 - +8191) which use two bytes are indicated as: 00 00H = -8192, 40 00H = 0, 7F 7FH = +8191.

  Example: aa bbH - 40 00H = (aa x 128 + bb) - (64 x 128), e.g. 6A 64H - 40 00H = (106 x 128 + 100) - 8192 = 13668 - 8192 = 5476. (9)

### ● Table A-2: ASCII code

| Character   | H     | Character   | H     | Character   | H     |
| ----------- | ----- | ----------- | ----- | ----------- | ----- |
| SP          | 20H   |             |       |             |       |
| A           | 41H   | a           | 61H   | 1           | 31H   |
| B           | 42H   | b           | 62H   | 2           | 32H   |
| C           | 43H   | c           | 63H   | 3           | 33H   |
| D           | 44H   | d           | 64H   | 4           | 34H   |
| E           | 45H   | e           | 65H   | 5           | 35H   |
| F           | 46H   | f           | 66H   | 6           | 36H   |
| G           | 47H   | g           | 67H   | 7           | 37H   |
| H           | 48H   | h           | 68H   | 8           | 38H   |
| I           | 49H   | i           | 69H   | 9           | 39H   |
| J           | 4AH   | j           | 6AH   | 0           | 30H   |
| K           | 4BH   | k           | 6BH   | -           | 2DH   |
| L           | 4CH   | l           | 6CH   | .           | 2EH   |
| M           | 4DH   | m           | 6DH   | /           | 2FH   |
| N           | 4EH   | n           | 6EH   | :           | 3AH   |
| O           | 4FH   | o           | 6FH   | ;           | 3BH   |
| P           | 50H   | p           | 70H   | <           | 3CH   |
| Q           | 51H   | q           | 71H   | =           | 3DH   |
| R           | 52H   | r           | 72H   | >           | 3EH   |
| S           | 53H   | s           | 73H   | ?           | 3FH   |
| T           | 54H   | t           | 74H   | @           | 40H   |
| U           | 55H   | u           | 75H   |             |       |
| V           | 56H   | v           | 76H   |             |       |
| W           | 57H   | w           | 77H   |             |       |
| X           | 58H   | x           | 78H   |             |       |
| Y           | 59H   | y           | 79H   |             |       |
| Z           | 5AH   | z           | 7AH   |             |       |

**Note:** "SP" is space.

---

## Seite 156 — MIDI Implementation Chart

**MULTI TIMBRAL SYNTHESIZER MODULE** — Model: **JV-1080** — MIDI Implementation Chart — Date: **July 18, 1994** — Version: **1.00**

| Function                                                                                                    | Transmitted             | Recognized                          | Remarks                     |
| ----------------------------------------------------------------------------------------------------------- | ----------------------- | ----------------------------------- | --------------------------- |
| Basic Channel: Default / Changed                                                                            | X / X                   | 1 - 16, OFF / 1 - 16, OFF           | Memorized                   |
| Mode: Default / Messages / Altered                                                                          | X / X / ***             | Mode 3 / *** / Mode 3, 4 (M = 1)    |                             |
| Note Number: True Voice                                                                                     | X / ***                 | 0 - 127 / 0 - 127                   |                             |
| Velocity: Note ON / Note OFF                                                                                | X / X                   | O / O                               |                             |
| After Touch: Key's / Ch's                                                                                   | X / X                   | O *1 / O *1                         |                             |
| Pitch Bend                                                                                                  | X                       | O *1                                | Resolution: 9 bits          |
| Control Change                                                                                              |                         |                                     | (siehe unten)               |
| Prog Change: True #                                                                                         | X / ***                 | O *1 / 0 - 127                      | Program Number 1 - 128      |
| System Exclusive                                                                                            | O                       | O *1                                |                             |
| System Common: Song Pos / Song Sel / Tune                                                                   | X / X / X               | X / X / X                           |                             |
| System Real Time: Clock / Commands                                                                          | X / X                   | O *1 / X                            |                             |
| Aux Messages: All Sound OFF / Reset All Controllers / Local ON/OFF / All Notes OFF / Active Sense / Reset   | X / X / X / X / X / X   | O / O / X / O (123 - 127) / O / X   |                             |

### Control Change (Einzelzeilen)

| Controller   | Transmitted    | Recognized   | Remarks                               |
| ------------ | -------------- | ------------ | ------------------------------------- |
| 0, 32        | X              | O *2         | Bank select                           |
| 1            | X              | O *1         | Modulation                            |
| 2            | X              | O *1         | Breath                                |
| 4            | X              | O *1         | Foot type                             |
| 5            | X              | O *1         | Portamento time                       |
| 6, 38        | X              | O *1         | Data entry                            |
| 7            | X              | O *1         | Volume                                |
| 8            | X              | O *1         | Balance                               |
| 10           | X              | O *1         | Panpot                                |
| 11           | X              | O *1         | Expression                            |
| 64           | X              | O *1         | Hold 1                                |
| 65           | X              | O *1         | Portamento                            |
| 66           | X              | O *1         | Sostenuto                             |
| 67           | X              | O *1         | Soft pedal                            |
| 69           | X              | O *1         | Hold 2                                |
| 84           | X              | O *1         | Portamento control                    |
| 91           | X              | O *1         | General purpose effects 1 (Reverb)    |
| 93           | X              | O *1         | General purpose effects 3 (Chorus)    |
| 100, 101     | X              | O *1         | RPN LSB, MSB                          |

**Notes:**

- **\*1** Can be set to O or X manually and memorized.
- **\*2** Can be changed manually and memorized.

**Mode 1:** OMNI ON, POLY  **Mode 2:** OMNI ON, MONO  **Mode 3:** OMNI OFF, POLY  **Mode 4:** OMNI OFF, MONO
**O:** Yes  **X:** No

---

## Anmerkungen zur Transkription (1)–(17)

- **(1)** Patch group / Patch mode patch group: Im Scan ist die Werteliste unscharf; lesbar sind "USER" und "PCM". Beim JV-1080 existieren die Gruppen USER, PCM (Wellenformen der SR-JV80-Boards) und EXP — im Manual des Schwestergeräts JV-1010 lauten die Namen "USER&PRESET / PCM / EXP". Mit an Sicherheit grenzender Wahrscheinlichkeit gilt hier (USER, PCM, EXP).
- **(2)** System control source 1/2: Wertebereich 0 - 97 ist im Scan lesbar, die zugehörige Werteliste nicht. Nach JV-1010-Manual: 1 - 97 = CC#01 - 05, 07 - 31, 64 - 95, PITCH BEND, AFTERTOUCH. Ob die JV-1080 zusätzlich den Wert 0 (OFF) kennt, ist im Scan nicht erkennbar.
- **(3)** Velocity-Sensitivitäten (0 - 125): Im Scan sind nur Fragmente wie "(+50 - +200)" erkennbar. Nach Schwester-Manual und corpus-geprüfter Referenz des JV-1080-SysEx-Formats entspricht der Bereich 0 - 125 = −100 … +150 %.
- **(4)** Preview sound mode: Im Scan nur "0 - 1 (SINGLE, CHORD)" erkennbar; JV-1010: 0 - 2 (SINGLE, CHORD, PHRASE).
- **(5)** Einige Parameternamen der Rhythm-Note-Tabelle (*1-4-2) sind im Scan unleserlich; sie wurden aus der identischen Parameterkarte des JV-1010 rekonstruiert, insbesondere "Coarse tune" (OCR las "Source key") mit Wert 0 - 127 (C-1 - G9).
- **(6)** Rekonstruierte Beispiel-Bytesequenzen; alle Checksummen wurden nach der Formel `Checksumme = 128 − (Rest der Division (Adresse + Daten/Größe) / 128)` arithmetisch verifiziert. Konkret: RQ1 System common (S. 150) Checksumme 58H (im Scan "5EH" erkennbar); RQ1 Performance-Name USER:01 (S. 151) vollständig rekonstruiert (Checksumme 64H); DT1 Rhythm Note (S. 154) Key Note "D2" aus Adresse 02 09 26 00 und Checksumme 4FH rekonstruiert.
- **(7)** GS-Beispiel "ARABIA" (S. 155): Device-ID im Scan unleserlich (hier 10H = Unit 17 angesetzt); Checksumme aus den gelesenen Daten berechnet: 73H (im Scan "76H/77H"-artig erkennbar); Schlussbyte F7H. Die Datenbytes 3AH 6DH 3EH 34H 0DH 3EH 68H 3CH 6FH 40H 36H 0FH (Skalatune-Werte für C - B, Arabische Stimmung) sind im Scan teilweise unleserlich und wurden plausibel rekonstruiert.
- **(8)** Rhythm note key# 98: Im Scan ist "E7" erkennbar; aus der eigenen Konvention der Tabelle (key# 35 = B1) folgt rechnerisch D7.
- **(9)** Der Beispielblock "Example of exclusive data" (S. 150) sowie die Rechenbeispiele in den Notes zu Tabelle A-1 (S. 155) sind im Scan stark beschädigt. Struktur, Formulierungen und Zahlen wurden nach dem inhaltlich identischen Block im JV-1010-Owner's-Manual rekonstruiert (gleiche Model-ID 6AH, identische Adressen); die im JV-1080-Scan selbst lesbaren Fragmente (Zieladresse 01 00 00 28H, Device-ID-Erklärung "Unit 17 − 1 = 16 = 10H", Formel "128 − remainder") sind eingearbeitet. Alle Rechnungen sind arithmetisch korrekt.
- **(10)** Die Patch-Tone-Tabelle (*1-3-2, S. 152 - 153) ist im Scan nahezu vollständig zerstört (OCR lieferte nur unbrauchbare Zeichenfragmente). Sie wurde vollständig aus der identischen Parameterkarte des JV-1010 rekonstruiert, mit den JV-1080-spezifischen Unterschieden: Tone Delay Mode enthält KEY-INTERVAL (statt PLAYMATE), Output assign (MIX, EFX, OUTPUT1, OUTPUT2) statt (MIX, DIRECT-1, DIRECT-2). Offsets, Bytebreiten und Totalgröße (00 00 01 01, 129 Bytes) stimmen mit der unabhängig corpus-geprüften Referenz des JV-1080-SysEx-Formats überein.

- **(11)** Die mittleren Byte-Zeilen der RQ1-/DT1-Tabellen (S. 130–131: Adressbytes `aaH bbH ccH ddH`, Größen-/Datenbytes `eeH ffH ggH hhH`) sind im Scan teilweise unleserlich; die Platzhalter-Struktur wurde über die vollständig lesbaren Rahmenbytes (F0H 41H devH modelH 11H/12H … F7H), die Identität der beiden Tabellen untereinander und das JV-1010-Manual rekonstruiert. Auch der Text „in 20 milliseconds intervals" (Sendewiederholung bei RQ1) ist nur fragmentarisch lesbar.
- **(12)** Kopfzeile des MIDI-Implementation-Abschnitts (S. 132): Im Scan ist als Datum „JULY 9 1994" erkennbar; der MIDI Implementation Chart (S. 156) des selben Manuals nennt „July 18, 1994". Transkribiert wie gescannt; die Abweichung dürfte ein Scan-/Druckfehler sein.
- **(13)** Bank-Tabelle (S. 133): Die Zeilen „PCM Card" (MSB 83H) und „Expansion A–D" (MSB 84H, LSB 00H–07H) mit Program 0–127 sind identisch mit den entsprechenden Zeilen der Patch-Bank-Tabelle — so im Original gedruckt (vermutlich im Originalmanual übernommen); die Rhythm-Sektion des JV-1080 kennt keine SR-JV80-Bänke in dieser Form. Transkribiert wie gedruckt.
- **(14)** Die Tabellenköpfe und Wertesystematik der Pedal-Familie (Hold 1, Sostenuto, Soft Pedal, Hold 2, S. 137–139) sind im Scan stark verzerrt (z. B. falsche Spaltenzuordnung, „OFF" nicht lesbar). Die Systematik (Setting 1 = nur HOLD, Setting 2 = HOLD + nächstes Pedal, Kontrollquellen OFF/HOLD/SOSTENUTO/SOFT/HOLD-2) wurde durch Abgleich der vier Tabellen untereinander und mit den Modebeschreibungen des Textes rekonstruiert.
- **(15)** Korrigierte Controller-Nummern und Statusbytes (OCR-Fehler in Klammern): General Purpose Effect 1 (Reverb) = BnH 5BH (OCR „58H"), General Purpose Effect 3 (Chorus) = BnH 5DH (OCR „50H"); ferner Statusbytes 8nH/9nH (Note), AnH (Poly Key Pressure), BnH (Control Change), CnH (Program Change), DnH (Channel Pressure), EnH (Pitch Bend); Bank Select 00H/20H, Sostenuto 42H, Soft Pedal 43H, Hold 2 45H, Portamento Control 54H kkH, RPN LSB 64H llH, RPN MSB 65H mmH, Data Entry 06H/26H, All Sounds Off BnH 78H 00H, Reset All Controllers BnH 79H 00H; Pitch Bend überträgt das LSB zuerst (EnH llH mmH).
- **(16)** RPN-Beschreibungsabschnitt (S. 145): OCR las teils „RPN MSB" für LSB-Felder und verstümmelte Wertebereiche. Korrigiert: RPN 0 = Pitch Bend Sensitivity (mm 00H–0CH, gemeinsamer Bereich für BENDER RANGE UP/DOWN), RPN 1 = Fine Tuning (00 00H–40 00H–7F 7FH, −8192×50/8192…0…+8192×50/8192 cent), RPN 2 = Coarse Tuning (10H–40H–70H, −48…0…+48), RPN 127/127 = RPN Reset. Reihenfolge der Nachrichten: RPN wird in MSB- oder LSB-Reihenfolge empfangen; Data Entry muss MSB zuerst senden (LSB wird auf 0 gesetzt).
- **(17)** S. 146–147: korrigierte OCR-Fehler — „Turns off all sounding notes" (All Sounds Off), „Polyphonic Pressure" statt „Poly Key Pressure" im Chart-Kontext, „RPN Undefined" (Tabelle der RPN-Nummern), Pitch Bend Change „0 (center)"; Program Change überträgt die Programmnummer plus 1 (CnH ppH); Channel Pressure = DnH vvH. Die Zeilen „General purpose system controller 1/2" in der Reset-All-Controllers-Tabelle sind die Parameter „System Control Source 1/2" (System-Common-Tabelle *1-1-1, offsets 00 22/00 23); die auffällig spezifische Formulierung ist vermutlich kein OCR-Artefakt und wurde wie gedruckt transkribiert. — Außerdem stehen die Uhr-Tabellen auf S. 148 („1./2. Patch/Performance Mode") im Original unter einer nicht lesbaren Überschrift.
**Verwendete Gegenquellen:** Roland JV-1010 Owner's Manual (identische Parameter-Adress-Map, sauberer Text; ManualsLib), corpus-geprüfte Referenz des JV-1080-Patch-/SysEx-Formats (sigma/integral research documentation) sowie MIDI-Norm (Standard-Statusbytes wie BnH, FEH, F8H, F0H/F7H).