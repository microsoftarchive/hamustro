#!/bin/bash

set -e

if [[ `uname -s` == "Darwin" ]]; then
  # OSX instllation
  brew install go
  brew install protobuf
  # brew intall lua
  # luarocks install md5
else
  # Linux (Ubuntu & AWS) installation
  cd /tmp/
    wget https://storage.googleapis.com/golang/go1.4.1.linux-amd64.tar.gz
    tar -xf go1.4.1.linux-amd64.tar.gz && rm go1.4.1.linux-amd64.tar.gz
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

# Install wrk network monitoring tool
git clone https://github.com/wg/wrk.git /tmp/wrk
cd /tmp/wrk
  make
  mv wrk /usr/local/bin
cd -
rm -Rf /tmp/wrk

# Install dependencies
go get -u github.com/Azure/azure-sdk-for-go/storage
go get -u github.com/aws/aws-sdk-go
go get -u github.com/golang/protobuf/proto
go get -u github.com/go-ini/ini
go get -u github.com/jmespath/go-jmespath
go get -u github.com/golang/protobuf/protoc-gen-go
