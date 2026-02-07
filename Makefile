.PHONY: build clean test install

# Build the binary
build:
	go build -o bin/ghqx ./cmd/ghqx

# Install to GOPATH/bin
install:
	go install ./cmd/ghqx

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
