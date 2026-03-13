#!/bin/bash

# 默认值设置
REGISTRY="image.ppinfra.com/llm-ops"
IMAGE_NAME="backend"
PUSH=false
VERSION=$(git describe --tags --always)

# 帮助信息
show_help() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -h, --help        显示帮助信息"
    echo "  -p, --push        构建后推送镜像到仓库"
    echo "  -v, --version     指定版本号 (默认使用git commit hash)"
    echo "Example:"
    echo "  $0 --push         构建并推送镜像"
    echo "  $0 -v 1.0.0       构建指定版本的镜像"
}

# 参数解析
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -p|--push)
            PUSH=true
            shift
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        *)
            echo "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
done

# 确保已安装 docker buildx
if ! docker buildx version > /dev/null 2>&1; then
    echo "Error: docker buildx 未安装"
    exit 1
fi

# 创建并��用 buildx builder
docker buildx create --name llm-ops-builder --use || true
docker buildx inspect --bootstrap

# 编译二进制
echo "开始编译二进制文件..."
chmod +x scripts/build_binary.sh
if [ "$PUSH" = true ]; then
    ./scripts/build_binary.sh --platform linux/amd64
else
    ./scripts/build_binary.sh
fi

if [ $? -ne 0 ]; then
    echo "错误：编译失败"
    exit 1
fi

# 构建后端镜像
echo "开始构建后端镜像..."
if [ "$PUSH" = true ]; then
    docker buildx build --platform linux/amd64 \
        -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} \
        -t ${REGISTRY}/${IMAGE_NAME}:latest \
        --no-cache \
        --push .
else
    # 本地构建只构建当前平台
    echo "本地构建模式：仅构建当前平台镜像"
    PLATFORM="linux/$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"
    docker buildx build --platform ${PLATFORM} \
        -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} \
        -t ${REGISTRY}/${IMAGE_NAME}:latest \
        --no-cache \
        --load .
fi

# 检查构建结果
if [ $? -ne 0 ]; then
    echo "错误：构建失败"
    exit 1
fi

echo "构建完成！"
if [ "$PUSH" = true ]; then
    echo "镜像已推送到仓库："
else
    echo "镜像已构建完成："
fi
echo "- ${REGISTRY}/${IMAGE_NAME}:${VERSION}"
echo "- ${REGISTRY}/${IMAGE_NAME}:latest" 