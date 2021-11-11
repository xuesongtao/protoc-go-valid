#!/bin/bash

# 用于本地构建项目目录下 proto 文件, 功能如下:
# 1: protoc --go_out=xxx/library/protogo xxx.proto
# 2: protoc-go-valid -f=xxx/library/protogo
protoFileDirName="document" # proto 存放的目录
outPdProjectPath="library/protogo" # pd 放入的项目路径

function main() {
    curPath=$(pwd)
    echo "当前路径: ${curPath}"

    # 找到 app 所在 path
    tmpIndex=$(awk 'BEGIN{print index("'${curPath}'", "'${protoFileDirName}'")}')
    if [[ $tmpIndex == 0 ]]; then
        echo "${protoFileDirName} is not exists"
        return
    fi
    appDir=${curPath:0:$tmpIndex-2}
    goOutPath="${appDir}/${outPdProjectPath}"
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
    protoc --go_out=$goOutPath $@

    # tag 注入
    for protoFile in $@
    do
        filename=${protoFile%%'.proto'} # 去掉 xxx.proto 的 .proto 
        protoc-go-valid -p="${goOutPath}/${filename}.pb.go"
    done
}

# 执行, 示例 main xxx.proto xxx.proto
main $@