#!/bin/bash

if [ -z "$MIFY_VERSION" ]; then
    echo >&2 "No version" && exit 1
fi

if [ "$PWD" != "/build" ]; then
    echo >&2 "This script should be ran inside docker"
    exit 2
fi

DEBVERSION="$(echo "$MIFY_VERSION" | cut -c2-)"

WORKDIR=$PWD/build/deb

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

go build ./cmd/mify || exit 2

rm -rf "$WORKDIR" && mkdir -p "$WORKDIR"/mify/{DEBIAN,usr/bin}

echo "$CONTROL_DATA" > "$WORKDIR"/mify/DEBIAN/control
cp ./mify "$WORKDIR"/mify/usr/bin
cd "$WORKDIR" && dpkg-deb --build mify
