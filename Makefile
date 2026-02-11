.PHONY: build clean test install build-release install-release release-snapshot lint fmt deps

# Build variables
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo "unknown")
BUILD_FLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Build the binary (development - no version info)
build:
	go build -o bin/ghqx ./cmd/ghqx

# Build the binary with version information
build-release:
	go build $(BUILD_FLAGS) -o bin/ghqx ./cmd/ghqx

# Install to GOPATH/bin (development)
install:
	go install ./cmd/ghqx

# Install with version information
install-release:
	go install $(BUILD_FLAGS) ./cmd/ghqx

# Release build with GoReleaser (snapshot/local)
release-snapshot:
	goreleaser release --snapshot --rm-dist

# Release build with GoReleaser (requires git tag)
release:
	goreleaser release --clean

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy
