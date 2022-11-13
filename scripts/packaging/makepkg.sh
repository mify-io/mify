#!/bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
MIFY_VERSION="$(cd "$REPO_ROOT" && git describe --tags --always)"
MIFY_VERSION=v0.1

build_ubuntu() {
    docker build -t mify-ubuntu-build -f "$REPO_ROOT"/scripts/packaging/Dockerfile-ubuntu "$REPO_ROOT" || exit 2
    docker run -e MIFY_VERSION="$MIFY_VERSION" -v "$REPO_ROOT":/build --rm mify-ubuntu-build
}

build_tar_package() {
    cd "$REPO_ROOT" || exit 2
    GOOS=$1
    GOARCH=$2
    GOOS="$GOOS" GOARCH="$GOARCH" go build ./cmd/mify
    WORKDIR="$REPO_ROOT/build/tar-$GOOS-$GOARCH"
    rm -rf "$WORKDIR" && mkdir -p "$WORKDIR/bin"
    cp "$REPO_ROOT/mify" "$WORKDIR/bin"
    cp "$REPO_ROOT/README.md" "$WORKDIR/"
    cd "$WORKDIR" && tar -cvzf "$REPO_ROOT/build/mify-$GOOS-$GOARCH.tar.gz" .
}

build_mac() {
    cd "$REPO_ROOT" || exit 2
    GOOS=darwin GOARCH=amd64 go build
}

case "$1" in
    ubuntu) build_ubuntu ;;
    mac) build_tar_package darwin amd64 ;;
    mac-arm) build_tar_package darwin arm64 ;;
    linux) build_tar_package linux amd64 ;;
    linux-arm) build_tar_package linux arm64 ;;
    *) echo >&2 "usage: $0 ubuntu" ;;
esac
