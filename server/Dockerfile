# ---- 构建阶段 ----
FROM golang:1.25 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# ---- 运行阶段 ----
FROM debian:bookworm-slim

ENV GIN_MODE=release
ENV TZ=Asia/Shanghai
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

WORKDIR /app
COPY --from=builder /build/server ./server
COPY --from=builder /build/config.yaml ./config.yaml
COPY --from=builder /build/resource ./resource/
RUN groupadd --system --gid 10001 app \
    && useradd --system --no-create-home --home-dir /nonexistent --shell /usr/sbin/nologin --uid 10001 --gid 10001 app \
    && mkdir -p /app/log /app/uploads/file /app/resource/excel \
    && chown -R app:app /app

ENV GVA_CONFIG=/app/config.yaml

EXPOSE 8888
USER app:app
ENTRYPOINT ["/app/server"]
