#!/bin/bash

# 设置应用名称
APP_NAME="blog_cli"

# 创建 build 目录（如果不存在）
mkdir -p build

# 编译 Windows 版本 (64-bit)
echo "Building for Windows..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}_windows_amd64.exe

# 编译 Linux 版本 (64-bit)
echo "Building for Linux..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/${APP_NAME}_linux_amd64

# 编译 macOS 版本 (64-bit)
echo "Building for macOS..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}_darwin_amd64

# 编译 macOS ARM 版本 (Apple Silicon)
echo "Building for macOS ARM..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/${APP_NAME}_darwin_arm64

echo "Build complete! Check the build directory for the binaries." 