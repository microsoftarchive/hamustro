#!/bin/bash

set -e

if [[ `uname -s` == "Darwin" ]]; then
  # OSX instllation
  brew install go
  brew install protobuf
else
  # Linux (Ubuntu & AWS) installation
  cd /tmp/
    curl -O https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz
    tar -xvf go1.5.3.linux-amd64.tar.gz
    mv go /usr/local && mkdir -p /usr/local/gopath
  cd -
  apt-get install protobuf-compiler
fi

# Set environment variables
cd ~/
  echo 'export GOROOT=/usr/local/go' >> .profile
  echo 'export GOPATH=/usr/local/gopath' >> .profile
  echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> .profile
  source .profile
cd -
