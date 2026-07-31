BINARY  := agentic-cms
PREFIX  ?= /usr/local

.PHONY: build install test smoke-test clean

build:
	CGO_ENABLED=0 go build -o $(BINARY) .

install: build
	install -m 0755 $(BINARY) $(PREFIX)/bin/$(BINARY)

test:
	go vet ./...
	go test ./...

smoke-test: build test
	./$(BINARY) --version
	./$(BINARY) --help
	bash scripts/smoke-test-installer.sh ./$(BINARY)

clean:
	rm -f $(BINARY)
