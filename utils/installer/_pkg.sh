#!/bin/bash

set -e

source ~/.profile

# Install dependencies
go get -u github.com/Azure/azure-sdk-for-go/storage
go get -u github.com/aws/aws-sdk-go
go get -u github.com/golang/protobuf/proto
go get -u github.com/go-ini/ini
go get -u github.com/jmespath/go-jmespath
go get -u github.com/golang/protobuf/protoc-gen-go

