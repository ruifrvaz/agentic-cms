BINARY  := agentic-cms
BINDIR  := installer
PREFIX  ?= /usr/local

.PHONY: build install test smoke-test clean

build:
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 go build -o $(BINDIR)/$(BINARY) .

install: build
	install -m 0755 $(BINDIR)/$(BINARY) $(PREFIX)/bin/$(BINARY)

test:
	go vet ./...
	go test ./...

smoke-test: build test
	./$(BINDIR)/$(BINARY) --version
	./$(BINDIR)/$(BINARY) --help
	bash scripts/smoke-test-installer.sh ./$(BINDIR)/$(BINARY)

clean:
	rm -rf $(BINDIR)
