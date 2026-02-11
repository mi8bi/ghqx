# CI/CD Pipeline

This document describes the GitHub Actions workflows and release process for ghqx.

## Overview

The project uses GitHub Actions to:

1. **Test** - Run automated tests on pull requests and pushes to main
2. **Release** - Build and publish release artifacts for multiple platforms

## Test Workflow (test.yml)

Runs on:
- Every pull request to `main`
- Every push to `main`

Steps:
1. Checkout code
2. Set up Go environment from `go.mod`
3. Cache Go modules for faster builds
4. Run `go mod tidy` to verify dependencies
5. Run all unit tests with coverage
6. Upload coverage to Codecov (for tracking test coverage trends)

**File**: [.github/workflows/test.yml](.github/workflows/test.yml)

## Release Workflow (release.yml)

Triggered when a version tag is pushed (e.g., `v0.3.0`)

### Prerequisites

The GITHUB_TOKEN secret must be available (automatic in GitHub Actions)

### Release Process

1. **Trigger**: Push a Git tag matching `v*.*.*`
   ```bash
   git tag v0.3.0
   git push origin v0.3.0
   ```

2. **Workflow Steps**:
   - Checkout code with full commit history
   - Set up Go environment
   - Cache dependencies
   - Run GoReleaser with `.goreleaser.yaml`

3. **GoReleaser Output**:
   - Multi-platform binaries:
     - Linux: amd64, arm64 (tar.gz)
     - macOS: amd64, arm64 (tar.gz)
     - Windows: amd64 (zip)
   - SHA256 checksums
   - GitHub Release page with artifacts

**File**: [.github/workflows/release.yml](.github/workflows/release.yml)

## Version Embedding

### How It Works

When GoReleaser builds, it automatically injects version information via `-ldflags`:

```yaml
ldflags:
  - -X main.version={{.Version}}      # Git tag (e.g., v0.3.0)
  - -X main.commit={{.Commit}}        # Full commit hash
  - -X main.buildTime={{.CommitDate}} # Commit timestamp
```

### Verifying Embedded Information

After building a release, you can verify the version:

```bash
./ghqx version
# Output: ghqx v0.3.0

./ghqx version --verbose
# Output:
# ghqx v0.3.0
# commit: abc123def456...
# built at: 2026-02-11T00:00:00Z
# go version: go1.25.6
```

## Manual Testing Before Release

Before pushing a release tag, you can test locally:

### Using Make

```bash
# Build with snapshot version
make build

# Test the build
./bin/ghqx version
```

### Using GoReleaser Snapshot

```bash
# Create multi-platform snapshot release (no GitHub upload)
make release-snapshot

# Binaries will be in dist/ directory
dist/ghqx_*_linux_amd64/ghqx version
```

## Release Checklist

Before pushing a release tag:

- [ ] All tests pass locally and in CI
- [ ] Update CHANGELOG or release notes
- [ ] Bump version in appropriate files if needed
- [ ] Commit and push to main
- [ ] Wait for test workflow to pass
- [ ] Create and push version tag

```bash
# Create annotated tag (recommended)
git tag -a v0.3.0 -m "Release v0.3.0"
git push origin v0.3.0
```

## Post-Release

After the release workflow completes:

1. GitHub Releases page is automatically updated
2. Artifacts are available for download (binaries + checksums)
3. Tags can be installed via `go install`:
   ```bash
   go install github.com/mi8bi/ghqx@v0.3.0
   ```

## Troubleshooting

### Release workflow fails

1. Check [GitHub Actions](https://github.com/mi8bi/ghqx/actions) for error logs
2. Ensure tag format is correct: `v*.*.*`
3. Verify `GITHUB_TOKEN` permissions in repository settings

### GoReleaser issues

- GoReleaser version is pinned to `latest` in the workflow
- See [GoReleaser docs](https://goreleaser.com/) for configuration details
- Test locally with `make release-snapshot`

### Version not embedded correctly

- Check that the tag is annotated: `git tag -a v0.3.0`
- Ensure the workflow uses the correct `.goreleaser.yaml`
- Run `goreleaser build --debug` locally to debug

## Configuration Files

- [.goreleaser.yaml](.goreleaser.yaml) - Main GoReleaser config (all platforms)
- [.github/workflows/test.yml](.github/workflows/test.yml) - Test workflow
- [.github/workflows/release.yml](.github/workflows/release.yml) - Release workflow

## See Also

- [BUILDING.md](BUILDING.md) - Local build instructions
- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Actions Documentation](https://docs.github.com/actions)
