#!/bin/bash

# 用于本地构建项目目录下 proto 文件, 功能如下:
# 1: protoc --go_out=xxx/library/protogo xxx.proto
# 2: cjvalid -f=xxx/library/protogo
outPdProjectPath="library/protogo" # pd 放入的项目路径

function main() {
    curPath=$(pwd)
    echo "当前路径: ${curPath}"

    # 找到 app 所在 path
    tmpIndex=$(awk 'BEGIN{print index("'${curPath}'", "document")}')
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
    echo "protoc --go_out=${goOutPath} $@"
    protoc --go_out=$goOutPath $@

    # 进行 tag 注入
    cjvalid -f=$goOutPath
}

# 执行, 示例 main xxx.proto xxx.proto
main $@
