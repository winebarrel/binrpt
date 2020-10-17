SHELL   := /bin/bash
PROGRAM := binrpt
VERSION := v0.2.1
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/$(PROGRAM)

.PHONY: package
package: clean build
	gzip $(PROGRAM) -c > $(PROGRAM)_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum $(PROGRAM)_$(VERSION)_$(GOOS)_$(GOARCH).gz > $(PROGRAM)_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum

.PHONY: clean
clean:
	rm -f $(PROGRAM)
