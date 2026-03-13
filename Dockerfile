FROM test-image.ppinfra.com/test-public/alpine:3.20

WORKDIR /app

# 创建配置目录
RUN mkdir -p /app/config

# 复制启动脚本
COPY scripts/start.sh /app/
RUN chmod +x /app/start.sh

# 二进制文件会在构建时复制
COPY build/main /app/main

EXPOSE 8080

CMD ["/app/start.sh"] 