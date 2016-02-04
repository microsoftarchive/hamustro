#!/bin/bash

set -e

# Install wrk network monitoring tool
git clone https://github.com/wg/wrk.git /tmp/wrk
cd /tmp/wrk
  make
  mv wrk /usr/local/bin
cd -
rm -Rf /tmp/wrk
