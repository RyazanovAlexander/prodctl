BINDIR       := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
BINNAME      ?= command-executor
BUILDTIME    := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# git
LASTTAG     := $(shell git tag --sort=committerdate | tail -1)
GITSHORTSHA := $(shell git rev-parse --short HEAD)

# docker option
DTAG   ?= $(LASTTAG)
DFNAME ?= Dockerfile
DRNAME ?= docker.io/aryazanov/command-executor

# go option
PKG        := ./...
TESTS      := .
TESTFLAGS  :=
TAGS       :=

GOLDFLAGS += -w
GOLDFLAGS += -s
GOFLAGS   = -ldflags '$(GOLDFLAGS)'

GOOS   := linux
GOARCH := amd64

# Rebuild the buinary if any of these files change
SRC := $(shell find . -type f -name "*.go" -print) go.mod go.sum

# ------------------------------------------------------------------------------
#  init

init:
	minikube start
	minikube docker-env
	minikube -p minikube docker-env | Invoke-Expression
	sudo apt update
	sudo apt install jq

# ------------------------------------------------------------------------------
#  run

run: build
	$(BINDIR)/$(BINNAME)

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -o $(BINDIR)/$(BINNAME) .

# ------------------------------------------------------------------------------
#  install

.PHONY: install
install: build
	@install "$(BINDIR)/$(BINNAME)" "$(INSTALL_PATH)/$(BINNAME)"

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

# ------------------------------------------------------------------------------
#  cover

.PHONY: cover
cover:
	go test -v -coverpkg=./... -coverprofile=profile ./...
	go tool cover -html=profile

# ------------------------------------------------------------------------------
#  clean

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)'

# ------------------------------------------------------------------------------
#  proto

.PHONY: proto
proto:
	@protoc -I ./proto/ ./proto/exec.proto \
            --go_opt=paths=source_relative \
            --go_out=./internal/server/ \
            --go-grpc_out=./internal/server/ \
            --go-grpc_opt=paths=source_relative

	@protoc -I ./proto/ ./proto/exec.proto \
            --go_opt=paths=source_relative \
            --go_out=./fake/.agent/server \
            --go-grpc_out=./fake/.agent/server \
            --go-grpc_opt=paths=source_relative

# ------------------------------------------------------------------------------
#  container

.PHONY: container
container:
	docker build --build-arg LDFLAGS="$(GOLDFLAGS)" --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) --build-arg BUILD_IMAGE_TAG=$(shell jq '.build."golang-tag"' build-meta.jsonc) -t $(DRNAME):$(DTAG) -f ./$(DFNAME) .
	docker push $(DRNAME):$(DTAG)

	docker build --build-arg LDFLAGS="$(GOLDFLAGS)" --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) --build-arg BUILD_IMAGE_TAG=$(shell jq '.build."golang-alpine-tag"' build-meta.jsonc) -t $(DRNAME):$(DTAG)-alpine -f ./$(DFNAME) .
	docker push $(DRNAME):$(DTAG)-alpine

# ------------------------------------------------------------------------------
#  example-echo

# make example name=echo
.PHONY: example
example:
	@skaffold dev -f ./examples/$(name)/skaffold.yaml --no-prune=false --cache-artifacts=false

# make example-delete name=echo
.PHONY: example-delete
example-delete:
	@skaffold delete -f ./examples/$(name)/skaffold.yaml