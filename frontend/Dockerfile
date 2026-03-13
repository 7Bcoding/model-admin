# 构建阶段
FROM image.ppinfra.com/llm-ops/node:18-alpine as builder

WORKDIR /app
RUN npm config set registry https://registry.npmmirror.com
COPY package*.json ./
RUN npm install --legacy-peer-deps
RUN npm install -D unplugin-auto-import unplugin-vue-components
ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
# ARG VITE_FEISHU_CALLBACK_URL=https://ops.dev.alpha.dev/feishu/callback
ARG VITE_FEISHU_CALLBACK_URL=https://llm-ops.paigod.work/feishu/callback
ENV VITE_FEISHU_CALLBACK_URL=${VITE_FEISHU_CALLBACK_URL}
COPY . .
RUN npm run build

# 运行阶段
FROM image.ppinfra.com/llm-ops/nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# 创建配置目录
RUN mkdir -p /etc/nginx/templates
# 复制 nginx 配置模板
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
ENV BACKEND_HOST=backend
ENV BACKEND_PORT=8080
ENV BACKEND_beta=http://127.0.0.1:8111
ENV BACKEND_alpha=http://127.0.0.1:8111
ENV FUSION_beta_TOKEN=OGoFOcAB7QTDhfGQToqgxOdjNoV0fxyz
ENV FUSION_alpha_TOKEN=OGoFOcAB7QTDhfGQToqgxOdjNoV0fxyz
EXPOSE 80
CMD ["/bin/sh", "-c", "envsubst '${BACKEND_HOST} ${BACKEND_PORT} ${BACKEND_beta} ${BACKEND_alpha} ${FUSION_beta_TOKEN} ${FUSION_alpha_TOKEN}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]