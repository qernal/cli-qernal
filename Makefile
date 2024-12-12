# Define the binary name and the source directory
BINARY_NAME := qernal
SRC_DIR := cmd

# Detect the OS and architecture
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Get version, commit, and date
VERSION := $(shell git describe --tags --always --dirty="-dev" 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Define the build target
build:
	@echo "Building for OS: $(GOOS), ARCH: $(GOARCH)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "\
		-X github.com/qernal/cli-qernal/pkg/build.Version=$(VERSION) \
		-X github.com/qernal/cli-qernal/pkg/build.Commit=$(COMMIT) \
		-X github.com/qernal/cli-qernal/pkg/build.Date=$(DATE)" \
		-o $(SRC_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go
	@echo "Build complete: $(SRC_DIR)/$(BINARY_NAME)"

# Define the run target
run: build
	@echo "Running $(SRC_DIR)/$(BINARY_NAME)"
	./$(SRC_DIR)/$(BINARY_NAME)

# Clean target to remove the binary
clean:
	@echo "Cleaning up"
	rm -f $(SRC_DIR)/$(BINARY_NAME)
	@echo "Clean complete"

.PHONY: build run clean
