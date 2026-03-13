#!/bin/bash

# 检查是否提供了标签参数
if [ -z "$1" ]; then
    echo "Error: Please provide a tag for the image"
    echo "Usage: $0 <tag>"
    exit 1
fi

# 设置变量
TAG=$1
REGISTRY="image.ppinfra.com/llm-ops"
IMAGE_NAME="llm-admin"
FULL_IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}:${TAG}"
LATEST_IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}:latest"

echo "Building and pushing image: ${FULL_IMAGE_NAME}"

# 确保 data 目录存在
mkdir -p data

# 如果 users.json 不在 data 目录中，创建一个空的 users.json
if [ ! -f data/users.json ]; then
    echo "{}" > data/users.json
fi

# 构建镜像
echo "Building image..."
docker buildx build --platform linux/amd64 --load -t ${FULL_IMAGE_NAME} -t ${LATEST_IMAGE_NAME} .

# 检查构建是否成功
if [ $? -ne 0 ]; then
    echo "Error: Image build failed"
    exit 1
fi

# 推送镜像
echo "Pushing image to registry..."
docker push ${FULL_IMAGE_NAME}
docker push ${LATEST_IMAGE_NAME}

# 检查推送是否成功
if [ $? -ne 0 ]; then
    echo "Error: Image push failed"
    exit 1
fi

echo "Successfully built and pushed: ${FULL_IMAGE_NAME} and ${LATEST_IMAGE_NAME}"

# 清理本地镜像（可选）
read -p "Do you want to remove the local image? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker rmi ${FULL_IMAGE_NAME}
    docker rmi ${LATEST_IMAGE_NAME}
    echo "Local image removed"
fi 