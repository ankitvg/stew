.PHONY: pre-release test build build-all

VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
DIST_DIR ?= dist
LDFLAGS := -X github.com/ankitvg/stew/internal/version.Version=$(VERSION) -X github.com/ankitvg/stew/internal/version.Commit=$(COMMIT) -X github.com/ankitvg/stew/internal/version.Date=$(DATE)

pre-release:
	go test ./...
	go build ./...
	$(MAKE) build VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE)

test:
	go test ./...

build:
	mkdir -p $(DIST_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/stew ./cmd/stew

build-all:
	go build ./...
