#!/usr/bin/env bash
APP_NAME="run-flogo-app"
APP_VERSION=`cat VERSION`
LDFLAGS="-s -w -X 'github.com/abhijitWakchaure/run-flogo-app/config.VERSION=${APP_VERSION}'"

export CGO_ENABLED=0

rm -f dist/*

echo "### Building for platform: linux/amd64"
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/${APP_NAME}-linux_amd64
echo "### Building for platform: windows/amd64"
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/${APP_NAME}-windows_amd64.exe
echo "### Building for platform: darwin/amd64"
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/${APP_NAME}-darwin_amd64