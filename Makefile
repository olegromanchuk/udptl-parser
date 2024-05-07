# Get the latest commit branch, hash, and date
TAG=$(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
BRANCH=$(if $(TAG),$(TAG),$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))
HASH=$(shell git rev-parse --short=7 HEAD 2>/dev/null)
TIMESTAMP=$(shell git log -1 --format=%ct HEAD 2>/dev/null | xargs -I{} date -u -r {} +%Y%m%dT%H%M%S)
GIT_REV=$(shell printf "%s-%s-%s" "$(BRANCH)" "$(HASH)" "$(TIMESTAMP)")
REV=$(if $(filter --,$(GIT_REV)),latest,$(GIT_REV)) # fallback to latest if not in git repo

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run

# Project parameters
BINARY_NAME=udptl-parser
BINARY_MACOS=$(BINARY_NAME)_darwin
CMD_DIR=cmd/udptl-parser/

run:
	cd cmd/udptl-parser/; go build -o udptl-parser main.go
	cd ../../
	./cmd/udptl-parser/udptl-parser

test:
	go test ./cmd/udptl-parser/...
	go test ./...
	golangci-lint run

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_x64 -v $(CMD_DIR)

release:
	goreleaser release --snapshot --skip-publish --clean
.PHONY: run,test,release