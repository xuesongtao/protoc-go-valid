#!/bin/bash

# 编译Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o=protoc-go-valid-windows main.go

# 编译linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o=protoc-go-valid-linux main.go

# 编译darwin
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o=protoc-go-valid-darwin main.go