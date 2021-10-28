#!/bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

rm -f "$REPO_ROOT/mify"
(cd "$REPO_ROOT" || exit 1 && go build ./cmd/mify)
"$REPO_ROOT"/mify "$@"
