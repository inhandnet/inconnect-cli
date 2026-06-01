VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS  = -s -w \
  -X github.com/inhandnet/ics-cli/internal/build.Version=$(VERSION) \
  -X github.com/inhandnet/ics-cli/internal/build.Commit=$(COMMIT) \
  -X github.com/inhandnet/ics-cli/internal/build.Date=$(DATE)

.PHONY: build build-all install test lint fmt clean

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/ics ./cmd/ics

build-all:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/ics-linux-amd64       ./cmd/ics
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/ics-linux-arm64       ./cmd/ics
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/ics-darwin-amd64      ./cmd/ics
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/ics-darwin-arm64      ./cmd/ics
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/ics-windows-amd64.exe ./cmd/ics
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/ics-windows-arm64.exe ./cmd/ics

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/ics

test:
	go test ./... -v

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

clean:
	rm -rf bin/
