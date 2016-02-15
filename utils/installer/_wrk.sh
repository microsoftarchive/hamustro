#!/bin/bash

set -e

if [[ `uname -s` != "Darwin" ]]; then
  # Linux (Ubuntu & AWS) installation
  apt-get install libssl-dev
fi

# Install wrk network monitoring tool
git clone https://github.com/wg/wrk.git /tmp/wrk
cd /tmp/wrk
  make
  mv wrk /usr/local/bin
cd -
rm -Rf /tmp/wrk
