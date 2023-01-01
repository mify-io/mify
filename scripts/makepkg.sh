#!/bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
MIFY_VERSION="$(cd "$REPO_ROOT" && git describe --tags --always)"
PKGDIR="$REPO_ROOT/build/pkg"
MIFY_PATH="$REPO_ROOT/mify"
TARGETOS="$1"
TARGETARCH="$2"

if [ ! -e "$MIFY_PATH" ]; then
    echo >&2 "Failed to find mify binary" && exit 2
fi

if [ -z "$TARGETOS" ]; then
    echo >&2 "usage: $0 target-os target-arch" && exit 2
fi

if [ -z "$TARGETARCH" ]; then
    echo >&2 "usage: $0 target-os target-arch" && exit 2
fi

mkdir -p "$PKGDIR"

build_deb() {
    DEBVERSION="$(echo "$MIFY_VERSION" | cut -c2-)"
    WORKDIR="$REPO_ROOT/build/deb"
    rm -rf "$WORKDIR" && mkdir -p "$WORKDIR"/mify/{DEBIAN,usr/bin}

    CONTROL_DATA=$(cat <<HERE
Package: mify
Version: $DEBVERSION
Section: base
Priority: optional
Architecture: all
Maintainer: svc@mify.io
Description: Mify CLI - cloud service generator tool
HERE
    )

    echo "$CONTROL_DATA" > "$WORKDIR/mify/DEBIAN/control"
    cp "$MIFY_PATH" "$WORKDIR/mify/usr/bin"
    cd "$WORKDIR" && dpkg-deb --build mify
    mv mify.deb "$PKGDIR/mify-$MIFY_VERSION-ubuntu-$TARGETARCH.deb"
}

build_tar() {
    WORKDIR="$REPO_ROOT/build/tar-$TARGETOS-$TARGETARCH"
    rm -rf "$WORKDIR" && mkdir -p "$WORKDIR/bin"
    cp "$MIFY_PATH" "$WORKDIR/bin"
    cp "$REPO_ROOT/README.md" "$WORKDIR/"
    cd "$WORKDIR" && tar -cvzf "$PKGDIR/mify-$MIFY_VERSION-$TARGETOS-$TARGETARCH.tar.gz" .
}

case "$TARGETOS" in
    linux) build_tar && build_deb ;;
    darwin) build_tar ;;
    *) echo >&2 "usage: $0 target-os target-arch" && exit 2 ;;
esac
