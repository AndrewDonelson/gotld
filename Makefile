# file: Makefile
# description: Cross-platform makefile for building and testing the gotld package

# Define shell to use
SHELL := /bin/sh

# Go commands
GO      := go
GOTEST  := $(GO) test
GOBUILD := $(GO) build

# Output binary name (with OS-specific extension)
ifeq ($(OS),Windows_NT)
	BINARY_NAME := gotld-example.exe
else
	BINARY_NAME := gotld-example
endif

# Directories
EXAMPLE_DIR := ./example
BUILD_DIR   := ./build

# Files
EXAMPLE_MAIN := $(EXAMPLE_DIR)/main.go
BINARY_PATH  := $(BUILD_DIR)/$(BINARY_NAME)

# Make sure build directory exists
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Define all targets as phony
.PHONY: all clean test bench build run lint vet fmt check help

# Default target
all: check test build

# Help target
help:
	@echo "Available targets:"
	@echo "  all    - Run checks, tests, and build"
	@echo "  clean  - Remove build artifacts"
	@echo "  test   - Run tests"
	@echo "  bench  - Run benchmarks"
	@echo "  build  - Build the example application"
	@echo "  run    - Run the example application"
	@echo "  lint   - Run linter"
	@echo "  vet    - Run go vet"
	@echo "  fmt    - Run go fmt"
	@echo "  check  - Run all checks (fmt, vet, lint)"

# Test target
test:
	$(GOTEST) -v ./...

# Benchmarking target
bench:
	$(GOTEST) -bench=. ./...

# Build target
build: $(BUILD_DIR)
	$(GOBUILD) -o $(BINARY_PATH) $(EXAMPLE_MAIN)

# Run target
run: build
	$(BINARY_PATH)

# Clean target
clean:
	rm -rf $(BUILD_DIR)

# Linting target
lint:
	@command -v golint >/dev/null 2>&1 || { echo "golint not installed. Installing..."; go install golang.org/x/lint/golint@latest; }
	golint ./...

# Go vet target
vet:
	$(GO) vet ./...

# Go fmt target
fmt:
	$(GO) fmt ./...

# Check target - runs all checks
check: fmt vet lint