#!/usr/bin/env bash

directory=$(basename $(pwd))                           # 获取当前目录名字
executable=${directory}-$(date +%Y%m%d)$(go env GOEXE) # 生成可执行文件名

compile_time=$(date) # 生成编译时间
if [ $(uname -s) = "Linux" ]; then
    compile_time=$(date --iso-8601=seconds)
fi

go clean -cache
rm -rf ${directory}-*
ld_flags="-s -w -extldflags -static -X 'github.com/xmx/aegis-server/infra/banner.compileTime=${compile_time}'"
go build -o ${executable} -v -trimpath -ldflags "${ld_flags}" ./main/

if [ $? -eq 0 ]; then # 如果存在 upx 则进行压缩。
    if command -v upx &> /dev/null; then
        upx -9 ${executable}
    fi
fi
