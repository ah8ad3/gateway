#!/bin/bash

set -e

# Remove existing binaries
rm -f dist/gateway*

# Check if VERSION variable available
if [ -z "$VERSION" ]; then
  VERSION="0.1.0-alpha"
fi

echo "Building application version $VERSION"
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X github.com/ah8ad3/gateway/cmd.version=${VERSION}" -o "dist/gateway" ${PKG_SRC}

# https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
# all the arch and os for golang

if [ ! -z "${BUILD_DEFAULT}" ]; then
    echo "Only default binary was requested to build"
    exit 0
fi
