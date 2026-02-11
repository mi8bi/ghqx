# Building ghqx

This document explains how to build ghqx from source.

## Quick Start

### Using Makefile (Unix/Linux/macOS)

```bash
# Development build (no version info)
make build

# Release build with version information
make build-release VERSION=v0.3.0

# Output binary will be in bin/ghqx
./bin/ghqx version
```

### Using PowerShell Script (Windows)

```powershell
# Development build
.\scripts\build.ps1

# Release build with version information
.\scripts\build.ps1 -Version v0.3.0

# Output binary will be in bin\ghqx
.\bin\ghqx.exe version
```

### Using Bash Script (Unix/Linux/macOS)

```bash
# Development build
./scripts/build.sh

# Release build with version information
./scripts/build.sh --version v0.3.0

# Output binary will be in bin/ghqx
./bin/ghqx version
```

## Build with Version Information

### Using Makefile

The Makefile automatically captures build information:

```bash
make build-release VERSION=v0.3.0
```

This will:
- Set the version to `v0.3.0`
- Capture the Git commit hash
- Record the current UTC timestamp

### Using go build Directly

If you want to build manually with `go build`, use the `-ldflags` option:

**Unix/Linux/macOS (Bash):**
```bash
VERSION=v0.3.0
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

go build -ldflags "-X main.version=$VERSION -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME" \
  -o bin/ghqx ./cmd/ghqx
```

**Windows (PowerShell):**
```powershell
$VERSION="v0.3.0"
$COMMIT=$(git rev-parse --short HEAD)
$BUILD_TIME=$([DateTime]::UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ"))
$LDFLAGS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME"

go build -ldflags "$LDFLAGS" -o bin/ghqx .\cmd\ghqx
```

## Verifying the Build

After building, verify the version output:

```bash
# Default output
./bin/ghqx version
# Output: ghqx v0.3.0

# Detailed output
./bin/ghqx version --verbose
# Output:
# ghqx v0.3.0
# commit: abc123
# built at: 2026-02-11T04:12:00Z
# go version: go1.25.6
```

## Building with GoReleaser

GoReleaser is used for automated multi-platform releases in CI/CD pipelines.

### Prerequisites

```bash
# Install GoReleaser (optional, usually run in CI)
brew install goreleaser  # macOS
# or download from https://goreleaser.com/install/
```

### Local Release Build

```bash
# Build for all platforms (requires git tag)
goreleaser release --snapshot --rm-dist

# Build without creating a GitHub release (dry-run)
goreleaser release --clean --skip=announce
```

### Release Process (via GitHub Actions)

When you push a version tag, GitHub Actions automatically:

1. Checks out the code
2. Builds for multiple platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)
3. Creates checksums
4. Publishes to GitHub Releases

**To trigger a release:**
```bash
git tag v0.3.0
git push origin v0.3.0
```

### GoReleaser Configuration

The main configuration file is [.goreleaser.yaml](.goreleaser.yaml), which defines:

- **Builds**: Multi-platform configurations
- **Archives**: Packaging format (tar.gz for Unix, zip for Windows)
- **Checksums**: SHA256 verification files
- **Release**: GitHub release publishing
- **ldflags**: Automatic version embedding

The old `.goreleaser.linux.yaml`, `.goreleaser.macos.yaml`, and `.goreleaser.windows.yaml` files are no longer needed and can be kept for reference or removed.

### Embedded Version Information

When building with GoReleaser:
- `main.version` is set to the Git tag (e.g., `v0.3.0`)
- `main.commit` is set to the commit hash
- `main.buildTime` is set to the commit timestamp

This happens automatically without manual ldflags configuration.

### Verifying Release Artifacts

```bash
# Check checksums
sha256sum -c ghqx_v0.3.0_checksums.txt

# Test the binary
./ghqx_v0.3.0_linux_amd64/ghqx version
```

---

## Build Variables

The binary embeds the following variables during build:

| Variable | Description | Default |
|----------|-------------|---------|
| `main.version` | Application version (e.g., v0.3.0) | `dev` |
| `main.commit` | Git commit hash | `none` |
| `main.buildTime` | Build timestamp (RFC3339, UTC) | `unknown` |

## Common Issues

### Error: "malformed import path"

If you see an error like:
```
malformed import path "-ldflags": leading dash
```

Make sure you're using the correct go build syntax. The `-ldflags` flag must come **before** the package path:

**Correct:**
```bash
go build -ldflags "..." -o bin/ghqx ./cmd/ghqx
```

**Incorrect:**
```bash
go build ./cmd/ghqx -ldflags "..."  # Wrong order!
```

### Git Commit Not Available

If Git is not available in your environment, the build scripts will use `none` as the commit hash. This is safe and does not affect functionality.

### Time Format Issues on Windows

The PowerShell script uses `[DateTime]::UtcNow.ToString()` to ensure RFC3339 format compatibility across all systems.

## Development Workflow

For development, you can use the simple build target:

```bash
make build
./bin/ghqx version
# Output: ghqx dev  (no version info embedded)
```

Then when releasing, use the release target with the version tag:

```bash
make build-release VERSION=v0.3.0
```

## Dependencies

- Go 1.22 or later (as specified in go.mod)
- Git (optional, for commit hash capture)

## See Also

- [README.md](README.md) - Project overview
- [Makefile](Makefile) - Build automation
- [scripts/build.ps1](scripts/build.ps1) - PowerShell build script
- [scripts/build.sh](scripts/build.sh) - Bash build script
