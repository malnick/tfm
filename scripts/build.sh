#!/bin/bash
set -e

PLATFORM=$(uname | tr [:upper:] [:lower:])
GIT_REF=$(git tag | tail -1)
SOURCE_DIR=$(git rev-parse --show-toplevel)
VERSION=${GIT_REF}
REVISION=$(git rev-parse --short HEAD)
BUILD_DIR="${SOURCE_DIR}/build"
OS=$1

function build_cmd {
	echo "	Version:          ${VERSION}"
	echo "	Revision:         ${REVISION}"
	echo "	Operating System: ${OS}"
	
	GOOS=$OS go build -a -o ${BUILD_DIR}/tfm-${OS}-${GIT_REF}.${REVISION}                               \
        -ldflags "-X main.VERSION=${VERSION} -X main.REVISION=${REVISION}" \
        ${SOURCE_DIR}/cmd/tfm/tfm.go
}

function main {
    build_cmd
    if [ -n "$(which tree)" ]; then
        tree build/
    fi
}

main "$@"
