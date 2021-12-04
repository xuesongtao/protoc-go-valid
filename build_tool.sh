#!/bin/bash

# 根据不同平台进行编译
function goBuild() {
    echo "需要编译的平台选项: 1-Windows 2-Linux 3-Darwin"
    read -p "请输入选项: " platform
    case $platform in
    1)
        GOOS=windows
        ;;
    2)
        GOOS=linux
        ;;
    3)
        GOOS=darwin
        ;;
    *)
        echo "输入不正确"
        exit 1
        ;;
    esac
    CGO_ENABLED=0 GOARCH=amd64 go build -o=protoc-go-valid main.go
}

goBuild
