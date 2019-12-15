#!/bin/sh
echo "### Building for platform: linux/amd64"
env GOOS=linux GOARCH=amd64 go build -o dist/linux/run-flogo-app
echo "### Building for platform: windows/amd64"
env GOOS=windows GOARCH=amd64 go build -o dist/windows/run-flogo-app.exe
echo "### Building for platform: darwin/amd64"
env GOOS=darwin GOARCH=amd64 go build -o dist/darwin/run-flogo-app