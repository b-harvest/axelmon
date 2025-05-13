# Name of the output binary
BINARY := axelmon

# Derive version (strip leading “v”) and full commit hash from git
VERSION := $(shell git describe --tags --abbrev=0 | sed 's/^v//')
COMMIT  := $(shell git rev-parse --verify HEAD)

# Linker flags to embed version + commit into the binary
LDFLAGS := \
    -X 'bharvest.io/axelmon/version.Version=$(VERSION)' \
    -X 'bharvest.io/axelmon/version.Commit=$(COMMIT)'

# Build flags: enforce module readonly mode and pass in our LDFLAGS
BUILD_FLAGS := -mod=readonly -ldflags "$(LDFLAGS)"

# Go command
GO := go

.PHONY: all tidy build install

# Default target: ensure deps are tidy, then build
all: tidy build

# Tidy up go.mod / go.sum to pull in the correct versions of all deps
tidy:
	@echo "→ running go mod tidy"
	$(GO) mod tidy

# Build the project; output goes to bin/$(BINARY)
build:
	@echo "→ building $(BINARY) (version=$(VERSION) commit=$(COMMIT))"
	@mkdir -p bin
	$(GO) build $(BUILD_FLAGS) -o bin/$(BINARY) .

# Install into your GOBIN (or GOPATH/bin)
install:
	@echo "→ installing $(BINARY) (version=$(VERSION) commit=$(COMMIT))"
	$(GO) install $(BUILD_FLAGS) .
