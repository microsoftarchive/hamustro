#!/bin/bash

set -e

if [[ `uname -s` == "Darwin" ]]; then
  # OSX installation
  brew install protobuf
else
  # Linux installation
  apt-get install protobuf-compiler
fi
