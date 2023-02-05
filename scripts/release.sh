#!/bin/bash

increment_version() {
    local version="$(echo "$1" | cut -c2-)"
    local delimiter=.
    local array=($(echo "$version" | tr $delimiter '\n'))
    array[$2]=$((array[$2]+1))
    if [ $2 -lt 2 ]; then array[2]=0; fi
    if [ $2 -lt 1 ]; then array[1]=0; fi
    echo v$(local IFS=$delimiter ; echo "${array[*]}")
}

get_version() {
    git describe --tags --abbrev=0
}

usage() {
    echo >&2 "$0 <major|minor|patch>"
    echo >&2 "default version part is patch"
    exit 1
}

case "$1" in
    major) VERSION_PART=0 ;;
    minor) VERSION_PART=1 ;;
    patch) VERSION_PART=2 ;;
    * ) [ -z "$1" ] && VERSION_PART=2 || usage ;;
esac

if [ ! -z "$(git status --porcelain)" ]; then
    echo >&2 "Working directory is not clean, commit the changes before release"
    exit 2
fi

MAIN_BRANCH="upstream/main"

git fetch --all --tags || exit 2

VERSION="$(get_version)"
echo "current version: $VERSION"

NEXT_VERSION="$(increment_version "$VERSION" "$VERSION_PART")"

while true; do
    read -p "next version: $NEXT_VERSION, create the release? [Y/n] " yn
    case $yn in
        [Yy]* ) break ;;
        [Nn]* ) exit ;;
        * ) [ ! -z "$yn" ] && echo "Please answer yes or no" || break ;;
    esac
done

git tag "$NEXT_VERSION" "$MAIN_BRANCH" || exit 2
git push upstream "$NEXT_VERSION" || exit 2

echo "successfully created new release version: $NEXT_VERSION"
