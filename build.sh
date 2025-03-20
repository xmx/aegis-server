#!/usr/bin/env bash

BASE_NAME=$(basename $(pwd))
COMPILE_TIME=$(date --rfc-email)
VERSION=$(TZ="Europe/London" date -d "$(git log -1 --format=%cd --date=iso)"  +"%y.%-m.%-d-%H%M%S")
TARGET_NAME=${BASE_NAME}"-"v${VERSION}$(go env GOEXE)
echo "版本号："${VERSION}
echo "文件名："${TARGET_NAME}

go clean -cache
if [ "$1" = "clean" ]; then
    rm -rf ${BASE_NAME}*
    echo "清理结束"
    exit 0
fi

export CGO_ENABLED=0
LD_FLAGS="-s -w -extldflags=-static -X '$(go list -m)/banner.compileTime=${COMPILE_TIME}'"
go build -o ${TARGET_NAME} -trimpath -v -ldflags "${LD_FLAGS}" ./main

echo "编译结束"
