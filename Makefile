GOOS?=linux
GOARCH?=amd64
MIFY_VERSION=$(shell git describe --tags --always)

SUPPORTED_OS_LIST=linux darwin
SUPPORTED_ARCH_LIST=amd64 arm64

all: build lint

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -ldflags "-X github.com/mify-io/mify/cmd/mify/cmd.MIFY_VERSION=$(MIFY_VERSION)" ./cmd/mify

build-packages:
	for os in $(SUPPORTED_OS_LIST); do \
		for arch in $(SUPPORTED_ARCH_LIST); do \
			echo "making package for os: $$os, arch: $$arch"; \
			GOOS=$$os GOARCH=$$arch go build -v -ldflags "-X github.com/mify-io/mify/cmd/mify/cmd.MIFY_VERSION=$(MIFY_VERSION)" ./cmd/mify; \
			./scripts/makepkg.sh $$os $$arch; \
			echo "done os: $$os, arch: $$arch"; \
		done \
	done

test:
	go test -v ./...

lint:
	go vet ./...
	staticcheck ./...
