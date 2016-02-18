#!/bin/bash

set -e

HAMUSTRO_GOPATH="$GOPATH/src/github.com/sub-ninja/hamustro"

if ! [[ -d "$HAMUSTRO_GOPATH" || -L "$HAMUSTRO_GOPATH" ]]; then
  ln -s "`pwd`" "$HAMUSTRO_GOPATH"
fi