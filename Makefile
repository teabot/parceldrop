
SHELL=/bin/bash
GITSHORTHASH=$(shell git log -1 --pretty=format:%h)
GO ?= go
VERSION ?= $(GITSHORTHASH)

# Insert build metadata into binary
LDFLAGS := -X github.com/teabot/parceldrop/cmd.ParcelDropVersion=$(VERSION)
LDFLAGS += -X github.com/teabot/parceldrop/cmd.ParcelDropGitCommit=$(GITSHORTHASH)

.PHONY: test
test:
	$(GO) test $(shell go list ./... | grep -v /vendor/)

.PHONY: build
build: dep
	env GOOS=linux GOARCH=arm GOARM=5 $(GO) build -ldflags "$(LDFLAGS)" -o "parceldrop"

.PHONY: package
package: build test
	tar -cvzf parceldrop.tar.gz parceldrop parceldrop.service service.env INSTALL.sh

.PHONY: clean
clean:
	rm -rf parceldrop.tar.gz parceldrop

.PHONY: dep
dep:
	dep ensure -vendor-only