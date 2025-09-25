# Build and Release Guide

This document explains how to use the Makefile commands for building and releasing `diny`.

## Version Management

The project uses **automated version injection** from git tags. No manual version updates are needed in the code.

### How Versioning Works

- **Clean releases**: Show exact tag version (e.g., `1.2.3`)
- **Development builds**: Show commits ahead + hash (e.g., `1.2.3-5-g1a2b3c4`)  
- **Dirty working tree**: Show `-dirty` suffix (e.g., `1.2.3-dirty`)

## Make Commands

### Building

```bash
# Build the binary with automatic version injection
make build

# Clean build artifacts
make clean

# Build and clean in one command
make clean build
```

### Testing

```bash
# Run all tests
make test
```

### Version Information

```bash
# Show current version (from git tags)
make version

# Check built binary version
./diny --version
```

### Release Workflow

```bash
# Interactive release process
make release
```

This command will:
1. Show current version
2. Prompt for new version number
3. Create and push git tag
4. Build release binary with correct version

### Installation

```bash
# Install to $GOPATH/bin with correct version
make install
```

## Examples

### Development Build
```bash
$ make build
go build -ldflags "-X github.com/dinoDanic/diny/cmd.Version=0.1.1-1-geac875f" -o diny .

$ ./diny --version
diny version 0.1.1-1-geac875f
```

### Clean Release
```bash
$ git tag v0.2.0
$ make build
$ ./diny --version  
diny version 0.2.0
```

### Creating a Release
```bash
$ make release
Current version: 0.1.1-1-geac875f
Enter new version (e.g., 0.1.2): 0.2.0
# Creates tag v0.2.0, pushes to origin, builds release binary
```

## Notes

- Version is automatically synced with git tags
- No need to manually update version in code
- Always use `make build` instead of `go build` for proper versioning
- The `make release` command handles the complete release workflow