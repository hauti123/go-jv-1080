<#
.SYNOPSIS
    Reads every patch parameter of the JV-1080 temporary patch one by one
    with jv1080-cli and writes a markdown table of the results.

.DESCRIPTION
    Iterates over all patch parameters of the Patch Edit Mode structure
    (doc/patch-edit-mode-structure.md): patch common, per-Tone output
    routing, control, wave, LFO, pitch, TVF and TVA parameters. Each
    parameter is read with `jv1080-cli get patch ...`; errors such as
    timeouts are logged (console and log file) but never abort the run.

    The result is a markdown table with the columns:
      - command: the CLI command with everything up to and including
        the word "patch" skipped, e.g. "common level",
        "efx tone 2 reverb", "lfo 1 tone 3 rate"
      - value: the read value, or "error"

    The goal is to validate that each parameter is addressed correctly:
    a wrong offset typically yields a timeout (no answer from the
    device), so every "error" row marks a parameter to investigate.

    Per-Tone parameters are read for all four tones; envelope parameters
    are read through their category with the env marker. The patch name
    (12 bytes) is printed as a hex string by the CLI.

.PARAMETER Cli
    Path to the jv1080-cli binary (default: builds and uses
    ./jv1080-cli.exe via `go build`, or ./jv1080-cli on non-Windows).

.PARAMETER Out
    Name of the MIDI output port (passed as -out to the CLI).

.PARAMETER In
    Name of the MIDI input port (passed as -in to the CLI).

.PARAMETER OutNum
    Number of the MIDI output port (passed as -outnum to the CLI).

.PARAMETER InNum
    Number of the MIDI input port (passed as -innum to the CLI).

.PARAMETER DeviceId
    Device ID in hex (passed as -devid to the CLI, default "10").

.PARAMETER ReportFile
    Path of the markdown report to write (default "patch-parameters.md").

.PARAMETER LogFile
    Path of the error log (default "patch-parameters.log").

.EXAMPLE
    .\scripts\read-patch-parameters.ps1 -Out "SERIES 208i MIDI" -In "SERIES 208i MIDI"

.EXAMPLE
    .\scripts\read-patch-parameters.ps1 -OutNum 0 -InNum 0 -ReportFile report.md
#>
[CmdletBinding()]
param(
    [string]$Cli = "",
    [string]$Out = "SERIES 208i MIDI",
    [string]$In = "SERIES 208i MIDI",
    [int]$OutNum = 0,
    [int]$InNum = 0,
    [string]$DeviceId = "10",
    [string]$ReportFile = "patch-parameters.md",
    [string]$LogFile = "patch-parameters.log"
)

$ErrorActionPreference = "Continue"

# --- locate the CLI -------------------------------------------------------

if (-not $Cli) {
    $repoRoot = Split-Path -Parent $PSScriptRoot
    if ($IsWindows -or $env:OS -eq "Windows_NT") {
        $Cli = Join-Path $repoRoot "jv1080-cli.exe"
    } else {
        $Cli = Join-Path $repoRoot "jv1080-cli"
    }
    if (-not (Test-Path $Cli)) {
        Write-Host "building jv1080-cli..."
        Push-Location $repoRoot
        try {
            go build ./cmd/jv1080-cli
            if ($LASTEXITCODE -ne 0) {
                throw "go build failed with exit code $LASTEXITCODE"
            }
        }
        finally {
            Pop-Location
        }
    }
}
if (-not (Test-Path $Cli)) {
    throw "jv1080-cli not found at '$Cli' (pass -Cli <path>)"
}

$cliArgs = @("-out", $Out, "-in", $In, "-outnum", "$OutNum", "-innum", "$InNum", "-devid", $DeviceId)

# --- parameter list -------------------------------------------------------
# Each entry is the argument list after "get patch". The report shows the
# command with everything up to and including "patch" skipped.

$readCommands = New-Object System.Collections.Generic.List[string[]]

# Patch common.
$commonParams = @(
    "name",
    "efxtype", "efxprm1", "efxprm2", "efxprm3", "efxprm4", "efxprm5",
    "efxprm6", "efxprm7", "efxprm8", "efxprm9", "efxprm10", "efxprm11",
    "efxprm12", "efxout", "efxlevel", "efxchorus", "efxreverb",
    "efxctrl1src", "efxctrl1depth", "efxctrl2src", "efxctrl2depth",
    "choruslevel", "chorusrate", "chorusdepth", "choruspredelay",
    "chorusfeedback", "chorusout",
    "reverbtype", "reverblevel", "reverbtime", "hfdamp", "reverbfeedback",
    "tempo", "level", "pan", "analogfeel", "bendup", "benddown", "assign",
    "legato", "portamento", "portmode", "porttype", "portstart", "porttime",
    "efxhold", "velrange", "octave", "stretch", "priority",
    "struct12", "booster12", "struct34", "booster34"
)
foreach ($p in $commonParams) {
    $readCommands.Add(@("common", $p))
}

# Per-Tone parameters: output routing (efx), receive switches and ranges
# (ctrl/tone), wave, LFO 1-2, pitch, TVF and TVA.
foreach ($tone in 1..4) {
    foreach ($p in @("output", "outputlevel", "chorus", "reverb")) {
        $readCommands.Add(@("efx", "tone", "$tone", $p))
    }
    foreach ($p in @("redamper", "rxvolume", "rxhold", "rxbend", "rxpan")) {
        $readCommands.Add(@("ctrl", "tone", "$tone", $p))
    }
    foreach ($p in @("vellower", "velupper", "velfade", "keylower", "keyupper")) {
        $readCommands.Add(@("tone", "tone", "$tone", $p))
    }
    foreach ($p in @(
        "switch", "group", "groupid", "gain",
        "fxmswitch", "fxmcolor", "fxmdepth", "delaymode", "delaytime"
    )) {
        $readCommands.Add(@("wave", "tone", "$tone", $p))
    }
    foreach ($p in @("wave", "keytrig", "rate", "offset", "delay", "fade", "fadetime", "sync")) {
        $readCommands.Add(@("lfo", "1", "tone", "$tone", $p))
        $readCommands.Add(@("lfo", "2", "tone", "$tone", $p))
    }
    foreach ($p in @("coarse", "fine", "random", "keyfollow")) {
        $readCommands.Add(@("pitch", "tone", "$tone", $p))
    }
    foreach ($p in @(
        "depth", "velsense", "vt1", "vt4", "keyfollow",
        "t1", "t2", "t3", "t4", "l1", "l2", "l3", "l4"
    )) {
        $readCommands.Add(@("pitch", "tone", "$tone", "env", $p))
    }
    foreach ($p in @(
        "filtertype", "cutoff", "cutoffkeyfollow", "resonance", "velres",
        "envdepth", "velcurve", "velsens",
        "vt1", "vt4", "keyfollow", "t1", "t2", "t3", "t4", "l1", "l2", "l3", "l4"
    )) {
        if (@("vt1", "vt4", "keyfollow", "t1", "t2", "t3", "t4", "l1", "l2", "l3", "l4") -contains $p) {
            $readCommands.Add(@("tvf", "tone", "$tone", "env", $p))
        } else {
            $readCommands.Add(@("tvf", "tone", "$tone", $p))
        }
    }
    foreach ($p in @(
        "level", "pan", "envvel", "envcurve",
        "bias", "biaspoint", "biasdir",
        "pankeyfollow", "randompan", "altpan",
        "vt1", "vt4", "keyfollow", "t1", "t2", "t3", "t4", "l1", "l2", "l3"
    )) {
        if (@("vt1", "vt4", "keyfollow", "t1", "t2", "t3", "t4", "l1", "l2", "l3") -contains $p) {
            $readCommands.Add(@("tva", "tone", "$tone", "env", $p))
        } else {
            $readCommands.Add(@("tva", "tone", "$tone", $p))
        }
    }
}

# LFO depth 1:2 (common block, one row per LFO).
foreach ($lfo in 1..2) {
    foreach ($p in @("pitch", "filter", "amp", "pan")) {
        $readCommands.Add(@("lfo", "$lfo", $p))
    }
}

# Control 1-3 destinations and depths (per Tone).
foreach ($ctrl in 1..3) {
    foreach ($tone in 1..4) {
        foreach ($p in @("dst1", "dst2", "dst3", "dst4", "dep1", "dep2", "dep3", "dep4")) {
            $readCommands.Add(@("ctrl", "$ctrl", "tone", "$tone", $p))
        }
    }
}

# Control sources and hold/peak (common block).
foreach ($i in 2..3) {
    $readCommands.Add(@("ctrl", "$i", "src"))
}
foreach ($i in 1..3) {
    $readCommands.Add(@("ctrl", "$i", "hold"))
}

# --- read loop ------------------------------------------------------------

if (Test-Path $LogFile) { Remove-Item $LogFile }
$logLines = New-Object System.Collections.Generic.List[string]
$rows = New-Object System.Collections.Generic.List[object]
$errorCount = 0
$index = 0

foreach ($cmd in $readCommands) {
    $index++
    # Report command: everything up to and including "patch" is skipped.
    $label = ($cmd -join " ")
    $progress = "[{0}/{1}] {2}" -f $index, $readCommands.Count, $label
    Write-Host "$progress" -NoNewline

    $value = $null
    $errorText = $null
    try {
        # Merge stderr into stdout so the SEND/REC log lines and the value
        # line arrive in one stream; suppress the CLI's hex log chatter by
        # keeping only the last value line.
        $output = & $Cli @cliArgs get patch @cmd 2>&1
        $exit = $LASTEXITCODE
        if ($exit -ne 0) {
            # Use the CLI's error message (e.g. the timeout detail) as the
            # error text; fall back to the bare exit code.
            $msgLines = @($output | ForEach-Object { $_.ToString() } |
                Where-Object { $_ -and $_ -notmatch "^(SEND|REC )\s*:" })
            $errorText = if ($msgLines.Count -gt 0) { $msgLines[-1].Trim() } else { "exit code $exit" }
        } else {
            # The value line is the last one that does not start with
            # SEND: or REC :, e.g. "common level: 100" or
            # "tone 3 rate: 60".
            $valueLines = @($output | ForEach-Object { $_.ToString() } |
                Where-Object { $_ -and $_ -notmatch "^(SEND|REC )\s*:" })
            if ($valueLines.Count -gt 0) {
                $last = $valueLines[-1]
                $colon = $last.IndexOf(": ")
                if ($colon -ge 0) {
                    $value = $last.Substring($colon + 2).Trim()
                } else {
                    $value = $last.Trim()
                }
            }
        }
        if ($null -eq $value -and -not $errorText) {
            $errorText = "no answer (timeout)"
        }
    }
    catch {
        $errorText = $_.Exception.Message
    }

    if ($errorText) {
        $errorCount++
        $line = "{0}: {1}" -f $label, $errorText
        $logLines.Add($line)
        Write-Host " -> error ($errorText)" -ForegroundColor Red
        $rows.Add([pscustomobject]@{ Command = $label; Value = "error" })
    } else {
        Write-Host " = $value"
        $rows.Add([pscustomobject]@{ Command = $label; Value = $value })
    }
}

# --- report ---------------------------------------------------------------

$report = New-Object System.Collections.Generic.List[string]
$report.Add("# JV-1080 patch parameter readback")
$report.Add("")
$report.Add("Read $($readCommands.Count) parameters with $([System.IO.Path]::GetFileName($Cli)); $errorCount error(s).")
$report.Add("")
$report.Add("| command | value |")
$report.Add("| --- | --- |")
foreach ($row in $rows) {
    $valueText = "$($row.Value)" -replace "\|", "\\|"
    $report.Add("| {0} | {1} |" -f $row.Command, $valueText)
}
Set-Content -Path $ReportFile -Value ($report -join "`n") -Encoding UTF8

if ($logLines.Count -gt 0) {
    Set-Content -Path $LogFile -Value ($logLines -join "`n") -Encoding UTF8
}

Write-Host ""
Write-Host ("wrote {0} rows to {1} ({2} error(s))" -f $rows.Count, $ReportFile, $errorCount)
if ($logLines.Count -gt 0) {
    Write-Host ("errors logged to {0}" -f $LogFile)
}
