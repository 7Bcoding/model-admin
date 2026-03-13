#!/bin/sh

# 替换配置文件中的环境变量
envsubst < /app/config/config.yaml.template > /app/config/config.yaml

# 启动应用
exec ./main 