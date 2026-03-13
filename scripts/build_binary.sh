#!/bin/bash

# 默认只构建当前平台
BUILD_ALL=false
SPECIFIED_PLATFORM=""

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --all)
            BUILD_ALL=true
            shift
            ;;
        --platform)
            SPECIFIED_PLATFORM="$2"
            shift 2
            ;;
        *)
            echo "未知参数: $1"
            exit 1
            ;;
    esac
done

# 编译目标平台
if [ -n "$SPECIFIED_PLATFORM" ]; then
    PLATFORMS="$SPECIFIED_PLATFORM"
elif [ "$BUILD_ALL" = true ]; then
    PLATFORMS="linux/amd64 linux/arm64 darwin/arm64"
else
    # 获取当前系统信息
    CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    CURRENT_ARCH=$(uname -m)
    
    # 转换架构名称
    case ${CURRENT_ARCH} in
        x86_64)
            CURRENT_ARCH="amd64"
            ;;
        aarch64)
            CURRENT_ARCH="arm64"
            ;;
    esac
    
    PLATFORMS="${CURRENT_OS}/${CURRENT_ARCH}"
fi

OUTPUT_DIR="build"

# 清理并创建输出目录
rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}

# 编译不同平台的二进制
for platform in ${PLATFORMS}; do
    # 解析操作系统和架构
    os=${platform%/*}
    arch=${platform#*/}
    
    # 设置输出目录名称
    output_os=$os
    if [ "$os" = "darwin" ]; then
        output_os="mac"
    fi
    
    echo "Building for $os/$arch..."
    
    # 创建输出目录
    mkdir -p ${OUTPUT_DIR}
    
    # 设置编译环境变量
    export GOOS=$os
    export GOARCH=$arch
    export CGO_ENABLED=0
    
    # 编译
    go build -a -installsuffix cgo -o ${OUTPUT_DIR}/main .
    
    if [ $? -eq 0 ]; then
        echo "Successfully built for $os/$arch"
    else
        echo "Failed to build for $os/$arch"
        exit 1
    fi
done

echo "All builds completed!" 