# Define the binary name and the source directory
BINARY_NAME := qernal
SRC_DIR := cmd

# Detect the OS and architecture
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Define the build target
build:
	@echo "Building for OS: $(GOOS), ARCH: $(GOARCH)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(SRC_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go
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

