# 多阶段构建: 同时构建 server 和 web, 最终合并到一个容器（仅用于快速体验/演示）

# ---- 1. 构建后端 ----
FROM golang:1.24-alpine AS server-builder

RUN apk add --no-cache git

WORKDIR /build
COPY server/go.mod server/go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download

COPY server/ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# ---- 2. 构建前端 ----
FROM node:20-alpine AS web-builder

WORKDIR /build
COPY web/package.json web/package-lock.json ./
RUN npm ci --registry=https://registry.npmmirror.com

COPY web/ .
ENV VITE_BASE_API=/aiInfra
RUN npm run build

# ---- 3. 运行阶段 ----
FROM alpine:3.20

ENV TZ=Asia/Shanghai
RUN apk add --no-cache ca-certificates tzdata nginx \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && mkdir -p /run/nginx

WORKDIR /app

# 复制后端二进制和资源
COPY --from=server-builder /build/server .
COPY --from=server-builder /build/resource ./resource/

# 复制前端产物
COPY --from=web-builder /build/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY web/nginx.conf /etc/nginx/conf.d/default.conf

# 复制启动脚本
COPY deploy/docker/entrypoint.sh .
RUN chmod +x ./entrypoint.sh

EXPOSE 80 8888

ENTRYPOINT ["./entrypoint.sh"]