# Self-sign the Fin installer using makecert and signtool
# Requires Windows SDK (signtool.exe)

$certName = "bfrovrflw"
$installerPath = "Fin-v1.0.0-Setup.exe"

Write-Host "=== Signing Installer ===" -ForegroundColor Cyan

# Find signtool.exe
$signtoolPaths = @(
    "C:\Program Files (x86)\Windows Kits\10\bin\*\x64\signtool.exe",
    "C:\Program Files (x86)\Windows Kits\10\App Certification Kit\signtool.exe",
    "C:\Program Files\Microsoft SDKs\Windows\*\bin\signtool.exe"
)

$signtool = $null
foreach ($path in $signtoolPaths) {
    $found = Get-Item $path -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($found) {
        $signtool = $found.FullName
        break
    }
}

if (-not $signtool) {
    Write-Host "ERROR: signtool.exe not found. Please install Windows SDK." -ForegroundColor Red
    Write-Host "Download from: https://developer.microsoft.com/windows/downloads/windows-sdk/" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "The installer is built but not signed." -ForegroundColor Yellow
    Write-Host "It will show 'Unknown Publisher' when run." -ForegroundColor Yellow
    exit 1
}

Write-Host "Found signtool: $signtool" -ForegroundColor Green

# Create self-signed certificate using certutil
$certExists = certutil -store My | Select-String -Pattern $certName -Quiet
if (-not $certExists) {
    Write-Host "Creating self-signed certificate..." -ForegroundColor Cyan
    
    # Use PowerShell Certificate provider in Windows PowerShell
    $cert = powershell.exe -Command "New-SelfSignedCertificate -Type CodeSigningCert -Subject 'CN=$certName' -CertStoreLocation Cert:\CurrentUser\My"
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Failed to create certificate" -ForegroundColor Red
        exit 1
    }
    Write-Host "Certificate created successfully" -ForegroundColor Green
} else {
    Write-Host "Certificate already exists" -ForegroundColor Yellow
}

# Sign the installer
Write-Host "Signing $installerPath..." -ForegroundColor Cyan
& $signtool sign /n $certName /t http://timestamp.digicert.com /fd SHA256 $installerPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "Installer signed successfully!" -ForegroundColor Green
} else {
    Write-Host "WARNING: Signing failed. Installer will show 'Unknown Publisher'" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Note ===" -ForegroundColor Yellow
Write-Host "Self-signed certificates will still show a security warning."
Write-Host "To remove the warning, you need a commercial code signing certificate."
Write-Host "Publisher will show as: $certName" -ForegroundColor Cyan

