#!/bin/bash
# Bash script for building ghqx with version information
# Usage: ./scripts/build.sh [--version <version>] [--release]

set -e

VERSION="dev"
RELEASE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --version)
      VERSION="$2"
      shift 2
      ;;
    --release)
      RELEASE=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Get Git commit hash
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")

# Get current UTC time in RFC3339 format
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

# Build flags
LDFLAGS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME"
BUILD_FLAGS="-ldflags \"$LDFLAGS\""

# Ensure output directory exists
mkdir -p bin

# Build command
BUILD_CMD="go build $BUILD_FLAGS -o bin/ghqx ./cmd/ghqx"

echo "Building ghqx..."
echo "  Version: $VERSION"
echo "  Commit: $COMMIT"
echo "  Built at: $BUILD_TIME"
echo ""
echo "Command: $BUILD_CMD" | sed 's/.*-ldflags/go build -ldflags/'
echo ""

# Execute build
eval "$BUILD_CMD"

if [ $? -eq 0 ]; then
  echo "✓ Build succeeded!"
  echo ""
  echo "  Binary: bin/ghqx"
  echo ""
  echo "Testing version output:"
  ./bin/ghqx version
else
  echo "✗ Build failed!"
  exit 1
fi
