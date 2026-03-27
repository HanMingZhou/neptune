# ---- 构建阶段 ----
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

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
RUN apk add --no-cache ca-certificates tzdata \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

WORKDIR /app
COPY --from=builder /build/server .
COPY --from=builder /build/config.yaml ./config.yaml
COPY --from=builder /build/resource ./resource/
RUN mkdir -p /app/log /app/uploads/file

ENV GVA_CONFIG=/app/config.yaml

EXPOSE 8888
ENTRYPOINT ["./server"]
