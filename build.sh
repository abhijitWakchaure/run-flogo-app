#!/usr/bin/env bash

GEN_DOCS="false"

APP_NAME="run-flogo-app"
APP_VERSION=`cat VERSION`
LDFLAGS="-s -w -X 'github.com/abhijitWakchaure/run-flogo-app/config.VERSION=${APP_VERSION}'"

echo "Building binaries for ${APP_NAME}-${APP_VERSION}"

export CGO_ENABLED=0

rm -rf dist/*

DOC_TAG=$([ "$GEN_DOCS" = "true" ] && echo "-tags docs" || echo "")
[ "$GEN_DOCS" = "true" ] && echo "Using DOC_TAG: ${DOC_TAG}"
echo "### Building for platform: linux/amd64"
GOOS=linux GOARCH=amd64 go build ${DOC_TAG} -ldflags "${LDFLAGS}" -o dist/${APP_NAME}-linux_amd64
echo "### Building for platform: windows/amd64"
GOOS=windows GOARCH=amd64 go build ${DOC_TAG} -ldflags "${LDFLAGS}" -o dist/${APP_NAME}-windows_amd64.exe
echo "### Building for platform: darwin/amd64"
GOOS=darwin GOARCH=amd64 go build ${DOC_TAG} -ldflags "${LDFLAGS}" -o dist/${APP_NAME}-darwin_amd64


echo "### Building for platform: linux/arm64"
GOOS=linux GOARCH=arm64 go build ${DOC_TAG} -ldflags "${LDFLAGS}" -o dist/${APP_NAME}-linux_arm64
echo "### Building for platform: windows/arm64"
GOOS=windows GOARCH=arm64 go build ${DOC_TAG} -ldflags "${LDFLAGS}" -o dist/${APP_NAME}-windows_arm64.exe
echo "### Building for platform: darwin/arm"
GOOS=darwin GOARCH=arm64 go build ${DOC_TAG} -ldflags "${LDFLAGS}" -o dist/${APP_NAME}-darwin_arm64
