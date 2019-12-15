#!/bin/sh
echo "### Building for platform: linux/amd64"
env GOOS=linux GOARCH=amd64 go build -o dist/run-flogo-app-linux-amd64
echo "### Building for platform: windows/amd64"
env GOOS=windows GOARCH=amd64 go build -o dist/run-flogo-app-windows-amd64.exe
echo "### Building for platform: darwin/amd64"
env GOOS=darwin GOARCH=amd64 go build -o dist/run-flogo-app-darwin-amd64