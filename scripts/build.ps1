# PowerShell script for building ghqx with version information
# Usage: .\scripts\build.ps1 [-Version <version>] [-Release]

param(
    [string]$Version = "dev",
    [switch]$Release = $false
)

# Get Git commit hash
$commit = "none"
try {
    $commit = git rev-parse --short HEAD
}
catch {
    Write-Warning "Could not get Git commit hash"
}

# Get current UTC time in RFC3339 format
$buildTime = [DateTime]::UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ")

# Build flags
$ldflags = "-X main.version=$Version -X main.commit=$commit -X main.buildTime=$buildTime"
$buildFlags = "-ldflags `"$ldflags`""

# Ensure output directory exists
$binDir = "bin"
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
}

# Build command
$buildCmd = "go build $buildFlags -o $binDir/ghqx.exe ./cmd/ghqx"

Write-Host "Building ghqx..." -ForegroundColor Green
Write-Host "  Version: $Version"
Write-Host "  Commit: $commit"
Write-Host "  Built at: $buildTime"
Write-Host ""
Write-Host "Command: $buildCmd" -ForegroundColor Gray
Write-Host ""

# Execute build
Invoke-Expression $buildCmd

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build succeeded!" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Binary: $binDir/ghqx.exe"
    
    # Test version output
    Write-Host ""
    Write-Host "Testing version output:" -ForegroundColor Gray
    & "$binDir/ghqx" version
}
else {
    Write-Host "✗ Build failed!" -ForegroundColor Red
    exit 1
}
