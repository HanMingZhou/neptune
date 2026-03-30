# ---- 构建阶段 ----
FROM golang:1.24-alpine AS builder

ARG ALPINE_MIRROR=https://mirrors.tuna.tsinghua.edu.cn/alpine
RUN sed -i "s#https://dl-cdn.alpinelinux.org/alpine#${ALPINE_MIRROR}#g" /etc/apk/repositories \
    && apk add --no-cache git

WORKDIR /build
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# ---- 运行阶段 ----
FROM alpine:3.20

ENV GIN_MODE=release
ENV TZ=Asia/Shanghai
ARG ALPINE_MIRROR=https://mirrors.tuna.tsinghua.edu.cn/alpine
RUN sed -i "s#https://dl-cdn.alpinelinux.org/alpine#${ALPINE_MIRROR}#g" /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

WORKDIR /app
COPY --from=builder /build/server ./server
COPY --from=builder /build/config.yaml ./config.yaml
COPY --from=builder /build/resource ./resource/
RUN addgroup -S app -g 10001 \
    && adduser -S -D -H -u 10001 -G app app \
    && mkdir -p /app/log /app/uploads/file /app/resource/excel \
    && chown -R app:app /app

ENV GVA_CONFIG=/app/config.yaml

EXPOSE 8888
USER app:app
ENTRYPOINT ["/app/server"]
