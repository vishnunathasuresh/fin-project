# Fin Compiler Build Script
# Builds fin.exe and creates the NSIS installer

param(
    [switch]$Sign = $false,
    [switch]$Help = $false
)

if ($Help) {
    Write-Host "Fin Build Script" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\build.ps1 [options]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Sign   Sign the installer with a certificate"
    Write-Host "  -Help   Show this help message"
    exit 0
}

$ErrorActionPreference = "Stop"

# Colors
$cyan = [System.ConsoleColor]::Cyan
$green = [System.ConsoleColor]::Green
$red = [System.ConsoleColor]::Red
$yellow = [System.ConsoleColor]::Yellow

function Write-Header {
    param([string]$Message)
    Write-Host ""
    Write-Host "=== $Message ===" -ForegroundColor $cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host $Message -ForegroundColor $green
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "ERROR: $Message" -ForegroundColor $red
}

# Build fin.exe
Write-Header "Building Fin Compiler"
try {
    & go build -o fin.exe ./cmd/fin
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Custom "Go build failed"
        exit 1
    }
    Write-Success "fin.exe built successfully"
} catch {
    Write-Error-Custom "Build failed: $_"
    exit 1
}

# Create installer
Write-Header "Creating Installer with NSIS"
try {
    Push-Location scripts
    & makensis fin_installer.nsi
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Custom "makensis failed"
        Pop-Location
        exit 1
    }
    Write-Success "Installer created successfully"
    Pop-Location
} catch {
    Write-Error-Custom "NSIS failed: $_"
    Pop-Location
    exit 1
}

# Sign installer if requested
if ($Sign) {
    try {
        Push-Location scripts
        & powershell -ExecutionPolicy Bypass -File sign.ps1
        if ($LASTEXITCODE -ne 0) {
            Write-Host "WARNING: Signing failed, but installer is built" -ForegroundColor $yellow
        }
        Pop-Location
    } catch {
        Write-Host "WARNING: Signing error: $_" -ForegroundColor $yellow
    }
}

# Summary
Write-Header "Build Complete"
Write-Host "Installer: scripts\Fin-v1.0.0-Setup.exe"
if (-not $Sign) {
    Write-Host "Tip: Run with -Sign to digitally sign the installer" -ForegroundColor $yellow
}
Write-Host ""
