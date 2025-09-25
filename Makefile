# Variables
BINARY_NAME=diny
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
VERSION_CLEAN := $(shell echo $(VERSION) | sed 's/^v//')
LDFLAGS=-ldflags "-X github.com/dinoDanic/diny/cmd.Version=$(VERSION_CLEAN)"

# Default target
.PHONY: all
all: build

# Build the binary with version injection
.PHONY: build
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Build for production (clean version from tags)
.PHONY: build-release
build-release:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
.PHONY: test
test:
	go test ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Create a new release (tag and build)
.PHONY: release
release:
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version (e.g., 0.1.2): " NEW_VERSION; \
	git tag -a v$$NEW_VERSION -m "Release v$$NEW_VERSION"; \
	git push origin v$$NEW_VERSION; \
	make build-release

# Show current version
.PHONY: version
version:
	@echo $(VERSION_CLEAN)

# Install binary to GOPATH/bin
.PHONY: install
install:
	go install $(LDFLAGS) .