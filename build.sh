#!/usr/bin/env bash

# Exit on error
# set -e

BASE_NAME="aegis-server"
TIME_VERSION=$(TZ=UTC git log -1 --format="%cd" --date=format-local:"%y.%m.%d" | awk -F. '{printf "%d.%d.%d\n", $1, $2, $3}')
VERSION="v${TIME_VERSION}"
SHORT_SHA=$(git rev-parse --short HEAD)
BINARY_NAME="${BASE_NAME}_$(go env GOOS)-$(go env GOARCH)@${VERSION}+${SHORT_SHA}$(go env GOEXE)"
CURRENT_TIME=$(date --rfc-email)
LD_FLAGS="-s -w -extldflags=-static -X '$(go list -m)/buildinfo.compileTime=${CURRENT_TIME}'"
CGO_ENABLED=0 go build -o "${BINARY_NAME}" -tags=osusergo,netgo -trimpath -ldflags "${LD_FLAGS}" ./main

BINARY_MD5=$(md5sum "${BINARY_NAME}" | awk '{print $1}')
BINARY_SHA1=$(sha1sum "${BINARY_NAME}" | awk '{print $1}')
BINARY_SHA256=$(sha256sum "${BINARY_NAME}" | awk '{print $1}')
cat <<EOF > VERSION
VERSION: ${VERSION}+${SHORT_SHA}
BINARY_FILE_NAME: ${BINARY_NAME}
BINARY_FILE_MD5: ${BINARY_MD5}
BINARY_FILE_SHA1: ${BINARY_SHA1}
BINARY_FILE_SHA256: ${BINARY_SHA256}
GO_VERSION: $(go version)
BUILD_TIME: ${CURRENT_TIME}
EOF

echo "编译产物：${BINARY_NAME}"
