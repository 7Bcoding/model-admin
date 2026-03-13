#!/bin/bash

# 默认版本
VERSION=${1:-"0.1.0"}
REGISTRY="image.ppinfra.com/llm-ops"
IMAGE_NAME="frontend"

# 帮助信息
show_help() {
    echo "Usage: $0 [version]"
    echo "Options:"
    echo "  version    指定版本号 (默认: 0.1.0)"
    echo "Example:"
    echo "  $0 1.0.0   构建版本 1.0.0"
    echo "  $0         使用默认版本 0.1.0"
}

# 如果第一个参数是 -h 或 --help，显示帮助信息
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

# 进入前端目录
cd frontend

# 确保使用正确的构建器
docker buildx create --name amd64-builder --platform linux/amd64 --use || true
docker buildx inspect --bootstrap

# 构建镜像
echo "Building frontend image version: ${VERSION}..."
docker buildx build \
  --platform linux/amd64 \
  --no-cache \
  --build-arg VITE_API_BASE_URL=/api \
  --build-arg VITE_FEISHU_CALLBACK_URL=https://llm-ops.paigod.work/feishu/callback \
  -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} \
  -t ${REGISTRY}/${IMAGE_NAME}:latest \
  --load \
  .

# 推送镜像
echo "Pushing frontend image..."
docker push ${REGISTRY}/${IMAGE_NAME}:${VERSION}
docker push ${REGISTRY}/${IMAGE_NAME}:latest

echo "Frontend build and push completed!"
echo "- ${REGISTRY}/${IMAGE_NAME}:${VERSION}"
echo "- ${REGISTRY}/${IMAGE_NAME}:latest"