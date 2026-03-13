#!/bin/sh

# 确保配置目录存在
mkdir -p config

# 编译
echo "Building..."
go build -o main .

if [ $? -ne 0 ]; then
    echo "Build failed"
    exit 1
fi

# 使用开发配置启动应用
CONFIG_FILE=./config/config.dev.yaml ./main 