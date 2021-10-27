#!/bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
(cd $SCRIPT_DIR || exit 1 && go build ./cmd/mify)
$SCRIPT_DIR/mify "$@"
