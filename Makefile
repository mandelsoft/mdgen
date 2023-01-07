# SPDX-FileCopyrightText: 2023 Mandelsoft
#
# SPDX-License-Identifier: Apache-2.0

NAME                                           = mdgen
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION                                        = $(shell git describe --tags --exact-match 2>/dev/null|| echo "$$(cat $(REPO_ROOT)/VERSION)")
COMMIT                                         = $(shell git rev-parse HEAD)
EFFECTIVE_VERSION                              = $(VERSION)+$(COMMIT)
GIT_TREE_STATE                                 := $(shell [ -z "$$(git status --porcelain 2>/dev/null)" ] && echo clean || echo dirty)


SOURCES := $(shell go list -f '{{$$I:=.Dir}}{{range .GoFiles }}{{$$I}}/{{.}} {{end}}' ./... )
GOPATH                                         := $(shell go env GOPATH)
GEN = $(REPO_ROOT)/gen

NOW         := $(shell date --rfc-3339=seconds | sed 's/ /T/')
BUILD_FLAGS := "-s -w \
 -X github.com/mandelsoft/mdgen/version.gitVersion=$(EFFECTIVE_VERSION) \
 -X github.com/mandelsoft/mdgen/version.gitTreeState=$(GIT_TREE_STATE) \
 -X github.com/mandelsoft/mdgen/version.gitCommit=$(COMMIT) \
 -X github.com/mandelsoft/mdgen/version.buildDate=$(NOW)"

build: ${SOURCES}
	mkdir -p bin
	go build -ldflags $(BUILD_FLAGS) -o bin/mdgen .

.PHONY: test
test: build
	go test ./...
	hack/test
	@rm -rf tmp/test
	@mkdir -p tmp/test
	bin/mdgen src tmp/test
	diff -ur doc tmp/test

.PHONY: all
all: test build doc

.PHONY: cross-build
cross-build: $(GEN)/.crossbuild
$(GEN)/.crossbuild: $(SOURCES) Makefile
	@for i in $(PLATFORMS); do \
    tag=$$(echo $$i | sed -e s:/:-:g); \
    echo GOARCH=$$(basename $$i) GOOS=$$(dirname $$i) go build -ldflags $(BUILD_FLAGS) -o $(GEN)/$(NAME)-$$tag .; \
    GOARCH=$$(basename $$i) GOOS=$$(dirname $$i) go build -ldflags $(BUILD_FLAGS) -o $(GEN)/$(NAME)-$$tag .; \
    done
	@touch $(GEN)/.crossbuild


.PHONY: doc
doc: build
	bin/mdgen src doc

.PHONY: clean
clean:
	rm bin/*
	rm "$(GEN)"/*
	rm -rf tmp/test
