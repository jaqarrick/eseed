# Binary name
BINARY_NAME=eseed

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-w -s"

# Test flags
TESTFLAGS=-v

.PHONY: all build test clean run deps lint

all: deps build test

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)

test:
	$(GOTEST) $(TESTFLAGS) ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f *.torrent

run: build
	./$(BINARY_NAME) -input test_file.bin

deps:
	$(GOMOD) download
	$(GOMOD) tidy

lint:
	golangci-lint run

# Development helpers
dev: build
	./$(BINARY_NAME) -input test_file.bin

# Kill any running instances
kill:
	pkill $(BINARY_NAME)

# Create test file
testfile:
	dd if=/dev/urandom of=test_file.bin bs=1M count=10

# Show help
help:
	@echo "Available targets:"
	@echo "  all        - Run deps, build, and test"
	@echo "  build      - Build the binary"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  run        - Build and run with test file"
	@echo "  deps       - Download dependencies"
	@echo "  lint       - Run linter"
	@echo "  dev        - Build and run in development mode"
	@echo "  kill       - Kill any running instances"
	@echo "  testfile   - Create a test file"
	@echo "  help       - Show this help message" 