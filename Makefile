PREFIX ?= $(HOME)/.local
VERSION := $(shell git describe --tags --exact-match 2>/dev/null || echo dev)

.PHONY: build test fmt install clean

build:
	mkdir -p bin
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/llm-lint ./cmd/llm-lint

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './project/*')

install: build
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 0755 bin/llm-lint $(DESTDIR)$(PREFIX)/bin/llm-lint

clean:
	rm -rf bin
