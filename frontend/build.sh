#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
# 默认值设置
REGISTRY="image.ppinfra.com/llm-ops"
IMAGE_NAME="frontend"
PUSH=${PUSH:-false}
VERSION=${VERSION:-$(git describe --tags --always)}
VITE_API_BASE_URL="/api"  # 默认前端 API 路径

# 帮助信息
show_help() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -h, --help        显示帮助信息"
    echo "  -p, --push        构建后推送镜像到仓库"
    echo "  -v, --version     指定版本号 (默认使用git commit hash)"
    echo "  -a, --api-url     指定前端 API 路径 (默认为 /api)"
    echo "Example:"
    echo "  $0 --push         构建并推送镜像"
    echo "  $0 -v 1.0.0       构建指定版本的镜像"
    echo "  $0 -a /custom-api 指定自定义 API 路径"
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
        -a|--api-url)
            VITE_API_BASE_URL="$2"
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

# 创建并使用 buildx builder
docker buildx create --name llm-ops-builder --use || true
docker buildx inspect --bootstrap

# 构建前端镜像
echo "开始构建前端镜像..."
if [ "$PUSH" = true ]; then
    docker buildx build --platform linux/amd64,linux/arm64 \
        -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} \
        -t ${REGISTRY}/${IMAGE_NAME}:latest \
        --build-arg VITE_API_BASE_URL=${VITE_API_BASE_URL} \
        --push .
else
    # 本地构建只构建当前平台
    echo "本地构建模式：仅构建当前平台镜像"
    docker buildx build --platform linux/$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') \
        -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} \
        -t ${REGISTRY}/${IMAGE_NAME}:latest \
        --build-arg VITE_API_BASE_URL=${VITE_API_BASE_URL} \
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