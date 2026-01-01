# Makefile for beads project

# Binary name
BINARY_NAME=bd

# Build directory
BUILD_DIR=build

# Install directory
INSTALL_DIR=$(HOME)/bin

# Get version info for ldflags
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
SHORT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags="-X main.Build=$(SHORT_COMMIT) -X main.Commit=$(COMMIT) -X main.Branch=$(BRANCH)"

.PHONY: all build test bench bench-quick clean install help

# Default target
all: build

# Build the bd binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/bd

# Run all tests (skips known broken tests listed in .test-skip)
test:
	@echo "Running tests..."
	@TEST_COVER=1 ./scripts/test.sh

# Run performance benchmarks (10K and 20K issue databases with automatic CPU profiling)
# Generates CPU profile: internal/storage/sqlite/bench-cpu-<timestamp>.prof
# View flamegraph: go tool pprof -http=:8080 <profile-file>
bench:
	@echo "Running performance benchmarks..."
	@echo "This will generate 10K and 20K issue databases and profile all operations."
	@echo "CPU profiles will be saved to internal/storage/sqlite/"
	@echo ""
	go test -bench=. -benchtime=1s -tags=bench -run=^$$ ./internal/storage/sqlite/ -timeout=30m
	@echo ""
	@echo "Benchmark complete. Profile files saved in internal/storage/sqlite/"
	@echo "View flamegraph: cd internal/storage/sqlite && go tool pprof -http=:8080 bench-cpu-*.prof"

# Run quick benchmarks (shorter benchtime for faster feedback)
bench-quick:
	@echo "Running quick performance benchmarks..."
	go test -bench=. -benchtime=100ms -tags=bench -run=^$$ ./internal/storage/sqlite/ -timeout=15m

# Install to ~/bin (symlink to build output)
# This makes `make build` automatically update the human-accessible CLI
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR) (symlink)..."
	@mkdir -p $(INSTALL_DIR)
	@# Remove existing file/symlink and create new symlink
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@ln -sf $(CURDIR)/$(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Linked $(INSTALL_DIR)/$(BINARY_NAME) → $(CURDIR)/$(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts and benchmark profiles
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f internal/storage/sqlite/bench-cpu-*.prof
	rm -f beads-perf-*.prof

# Show help
help:
	@echo "Beads Makefile targets:"
	@echo "  make build        - Build the bd binary"
	@echo "  make test         - Run all tests"
	@echo "  make bench        - Run performance benchmarks (generates CPU profiles)"
	@echo "  make bench-quick  - Run quick benchmarks (shorter benchtime)"
	@echo "  make install      - Install to ~/bin (symlink to build output)"
	@echo "  make clean        - Remove build artifacts and profile files"
	@echo "  make help         - Show this help message"
