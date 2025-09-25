# Variables
BINARY_NAME=diny
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
VERSION_CLEAN := $(shell echo $(VERSION) | sed 's/^v//')
LDFLAGS=-ldflags "-X github.com/dinoDanic/diny/cmd.Version=$(VERSION_CLEAN)"

# bump tag 
.PHONY: bump
release:
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version (e.g., 0.1.2): " NEW_VERSION; \
	git tag -a v$$NEW_VERSION -m "Release v$$NEW_VERSION"; \
	git push origin v$$NEW_VERSION

# Show current version
.PHONY: version
version:
	@echo $(VERSION_CLEAN)

