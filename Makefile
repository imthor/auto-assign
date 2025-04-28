# Build variables
BINARY_NAME=autoassigner
VERSION=$(shell git describe --tags --always --dirty || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X autoassigner/version.Version=${VERSION} -X autoassigner/version.BuildTime=${BUILD_TIME} -X autoassigner/version.GitCommit=${GIT_COMMIT}"

# Build targets
.PHONY: all build clean test release

all: clean build

build:
	go build ${LDFLAGS} -o ${BINARY_NAME} main.go

clean:
	rm -f ${BINARY_NAME}
	rm -f ${BINARY_NAME}-*

test:
	go test ./...

# Release targets
release: clean
	@echo "Building release version: ${VERSION}"
	# Build for multiple platforms
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-amd64-${VERSION} main.go
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-arm64-${VERSION} main.go
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-amd64-${VERSION} main.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-arm64-${VERSION} main.go
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-windows-amd64-${VERSION}.exe main.go
	@echo "Release binaries built successfully"

# Development targets
dev: clean
	go build ${LDFLAGS} -o ${BINARY_NAME} main.go

# Install target
install: build
	cp ${BINARY_NAME} /usr/local/bin/

# Help target
help:
	@echo "Available targets:"
	@echo "  all        - Clean and build the binary"
	@echo "  build      - Build the binary"
	@echo "  clean      - Remove built binaries"
	@echo "  test       - Run tests"
	@echo "  release    - Build binaries for multiple platforms"
	@echo "  dev        - Build development binary"
	@echo "  install    - Install binary to /usr/local/bin"
	@echo "  help       - Show this help message"