BINARY  := agentic-cms
PREFIX  ?= /usr/local

.PHONY: build install test clean

build:
	CGO_ENABLED=0 go build -o $(BINARY) .

install: build
	install -m 0755 $(BINARY) $(PREFIX)/bin/$(BINARY)

test:
	go vet ./...
	go test ./...

clean:
	rm -f $(BINARY)
