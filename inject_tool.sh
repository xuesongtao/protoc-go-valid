#!/bin/bash

# 用于本地构建项目目录下 proto 文件, 功能如下:
# 1: protoc --go_out=xxx/library/protogo xxx.proto
# 2: protoc-go-valid -f=xxx/library/protogo
protoFileDir="test" # proto 存放的目录
outPdProjectDir="test" # 生成 pb 的目录
# protoFileDir="document" # proto 存放的目录
# outPdProjectDir="library/protogo" # 生成 pb 的目录

function checkIsOk() {
    # $1 操作名

    if [[ $? > 0 ]]; then
        echo -e "\"${1}\" is \033[1;31mFailed\033[0m"
        exit 1
    fi
    echo -e "\"${1}\" is \033[1;32mSuccess\033[0m"
}

function main() {
    curPath=$(pwd)
    echo "当前路径: ${curPath}"

    # 找到 app 所在 path
    tmpIndex=$(awk 'BEGIN{print index("'${curPath}'", "'${protoFileDir}'")}')
    if [[ $tmpIndex == 0 ]]; then
        echo "${protoFileDir} is not exists"
        return
    fi
    appDir=${curPath:0:$tmpIndex-2}
    goOutPath="${appDir}/${outPdProjectDir}"
    if [[ ! -d $goOutPath ]]; then
        echo "go out path: ${goOutPath} is not exist"
        return
    fi

    if [[ $# == 0 ]]; then
        echo "targe file is not exist, it is mypro.sh xxx.proto"
        return
    elif [[ $# > 4 ]]; then
        echo "both build max 4 file"
        return
    fi

    # protoc 进行编译 
    protoc --go_out=paths=source_relative:$goOutPath $@
    checkIsOk "protoc"

    # tag 注入
    for protoFile in $@
    do
        filename=${protoFile%%'.proto'} # 去掉 xxx.proto 的 .proto 
        protoc-go-valid -f="${goOutPath}/${filename}.pb.go"
    done
}

# 执行, 示例 main xxx.proto xxx.proto
main $@