<#
.SYNOPSIS
    Runs the JV-1080 integration tests against a connected real device.

.DESCRIPTION
    Sets the environment variables the integration tests read
    (JV1080_INTEGRATION, JV1080_OUT/JV1080_IN or JV1080_OUT_NUM/JV1080_IN_NUM,
    JV1080_DEVICE_ID, JV1080_TEST_PATCH) and runs
    `go test -run Integration ./...`.

    The integration tests write to user patch JV1080_TEST_PATCH (default 113)
    and restore its original content afterwards.

.PARAMETER Out
    Name of the MIDI output port to send to (default "SERIES 208i MIDI").

.PARAMETER In
    Name of the MIDI input port to listen on (default "SERIES 208i MIDI").

.PARAMETER OutNum
    Number of the MIDI output port (only used when -Out is empty).

.PARAMETER InNum
    Number of the MIDI input port (only used when -In is empty).

.PARAMETER DeviceId
    Device ID of the JV-1080 in hex (default "10" = unit 17).

.PARAMETER TestPatch
    User patch slot the store tests write to (default 113). The original
    content is restored after the test.

.PARAMETER Verbose
    Common parameter (added by CmdletBinding): also passes -v to go test.

.EXAMPLE
    .\scripts\run-integration-tests.ps1 -Out "JV-1080" -In "JV-1080"

.EXAMPLE
    .\scripts\run-integration-tests.ps1 -OutNum 1 -InNum 1 -DeviceId 10 -Verbose
#>
[CmdletBinding()]
param(
    [string]$Out = "SERIES 208i MIDI",
    [string]$In = "SERIES 208i MIDI",
    [int]$OutNum = 0,
    [int]$InNum = 0,
    [string]$DeviceId = "10",
    [int]$TestPatch = 113
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
Push-Location $repoRoot
try {
    $env:JV1080_INTEGRATION = "1"
    if ($Out) { $env:JV1080_OUT = $Out } else { Remove-Item Env:JV1080_OUT -ErrorAction SilentlyContinue }
    if ($In)  { $env:JV1080_IN  = $In  } else { Remove-Item Env:JV1080_IN  -ErrorAction SilentlyContinue }
    $env:JV1080_OUT_NUM = "$OutNum"
    $env:JV1080_IN_NUM = "$InNum"
    $env:JV1080_DEVICE_ID = $DeviceId
    $env:JV1080_TEST_PATCH = "$TestPatch"

    Write-Host "Running integration tests against JV-1080 (device ID 0x$DeviceId, test patch $TestPatch)..."
    $goArgs = @("test", "-run", "Integration")
    if ($Verbose) { $goArgs = @("test", "-v", "-run", "Integration") }
    $goArgs += "./..."
    & go @goArgs
    if ($LASTEXITCODE -ne 0) {
        throw "go test failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
    Remove-Item Env:JV1080_INTEGRATION -ErrorAction SilentlyContinue
    Remove-Item Env:JV1080_OUT -ErrorAction SilentlyContinue
    Remove-Item Env:JV1080_IN -ErrorAction SilentlyContinue
    Remove-Item Env:JV1080_OUT_NUM -ErrorAction SilentlyContinue
    Remove-Item Env:JV1080_IN_NUM -ErrorAction SilentlyContinue
    Remove-Item Env:JV1080_DEVICE_ID -ErrorAction SilentlyContinue
    Remove-Item Env:JV1080_TEST_PATCH -ErrorAction SilentlyContinue
}
