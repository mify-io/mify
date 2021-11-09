#!/bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
PATH=$(go env GOPATH)/bin:$PATH

rm -f "$REPO_ROOT/mify"
(cd "$REPO_ROOT" || exit 1 && go build ./cmd/mify)
dlv exec "$REPO_ROOT"/mify "$@"
