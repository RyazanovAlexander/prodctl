BINDIR       := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
BINNAME      ?= prodctl
BUILDTIME    := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# git
LASTTAG     := $(shell git tag --sort=committerdate | tail -1)
GITSHORTSHA := $(shell git rev-parse --short HEAD)

# docker option
DTAG   ?= $(LASTTAG)
DFNAME ?= Dockerfile
DRNAME ?= docker.io
DINAME ?= $(DRNAME)/aryazanov/prodctl

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
#  run

run: build
	$(BINDIR)/$(BINNAME)

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	go build $(GOFLAGS) -tags '$(TAGS)' -o $(BINDIR)/$(BINNAME) .

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
	go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

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
#  image

.PHONY: image
image:
	docker build --build-arg LDFLAGS="$(GOLDFLAGS)" --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) --build-arg BUILD_IMAGE_TAG=$(shell jq '.build."golang-tag"' build-meta.jsonc) -t $(DINAME):$(DTAG) -f ./$(DFNAME) .
	docker push $(DINAME):$(DTAG)

# ------------------------------------------------------------------------------
#  product

.PHONY: product
product:
	docker build -t $(DRNAME)/aryazanov/bundle:$(DTAG) -f ./fakes/.bundle/$(DFNAME) ./fakes/.bundle/
	docker push $(DRNAME)/aryazanov/bundle:$(DTAG)

	docker build -t $(DRNAME)/aryazanov/product:$(DTAG) -f ./fakes/.repositories/cfg.product/$(DFNAME) ./fakes/.repositories/cfg.product
	docker push $(DRNAME)/aryazanov/product:$(DTAG)