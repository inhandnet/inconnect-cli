VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS  = -s -w \
  -X github.com/inhandnet/inconnect-cli/internal/build.Version=$(VERSION) \
  -X github.com/inhandnet/inconnect-cli/internal/build.Commit=$(COMMIT) \
  -X github.com/inhandnet/inconnect-cli/internal/build.Date=$(DATE)

.PHONY: build build-all install test lint fmt clean docs

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/inconnect ./cmd/inconnect

build-all:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/inconnect-linux-amd64       ./cmd/inconnect
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/inconnect-linux-arm64       ./cmd/inconnect
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/inconnect-darwin-amd64      ./cmd/inconnect
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/inconnect-darwin-arm64      ./cmd/inconnect
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/inconnect-windows-amd64.exe ./cmd/inconnect
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/inconnect-windows-arm64.exe ./cmd/inconnect

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/inconnect

test:
	go test ./... -v

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

clean:
	rm -rf bin/

# Generate per-command Markdown reference docs. Pass DOCS_DIR to override the
# output directory (the release workflow points this at the inconnect-skills repo's
# references/commands so a release regenerates the skills' command details).
DOCS_DIR ?= docs/commands
docs:
	go run ./cmd/docgen $(DOCS_DIR)
