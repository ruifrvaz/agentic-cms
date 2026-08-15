BINARY  := agentic-cms
BINDIR  := installer
PREFIX  ?= /usr/local

# Detect version from git tags, fallback to "dev" if no tags exist. Baked into
# the binary via ldflags; `go install .../agentic-cms@latest` bypasses this
# Makefile, so main.go falls back to the Go module version at runtime instead
# (see resolvedVersion in main.go).
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build build-all install test smoke-test version clean

build:
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINDIR)/$(BINARY) .

# Cross-compile the release assets `agentic-cms update` fetches. Linux only
# for now, matching the project's current supported-platform scope.
build-all:
	mkdir -p $(BINDIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINDIR)/$(BINARY)_linux_amd64 .
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINDIR)/$(BINARY)_linux_arm64 .

install: build
	install -m 0755 $(BINDIR)/$(BINARY) $(PREFIX)/bin/$(BINARY)

test:
	go vet ./...
	go test ./...

smoke-test: build test
	./$(BINDIR)/$(BINARY) --version
	./$(BINDIR)/$(BINARY) --help
	bash scripts/smoke-test-installer.sh ./$(BINDIR)/$(BINARY)

version:
	@echo "$(VERSION)"

clean:
	rm -rf $(BINDIR)
