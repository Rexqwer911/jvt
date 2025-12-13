$ErrorActionPreference = "Stop"

Write-Host "Building jvt.exe..."
if (-not (Test-Path "build")) {
    New-Item -ItemType Directory -Force -Path "build" | Out-Null
}

go build -o build/jvt.exe cmd/jvt/main.go
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Building Installer..."
iscc installer/jvt.iss
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Success! Installer created at: installer/Output/jvt-setup.exe"
