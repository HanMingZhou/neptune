#!/bin/sh
set -e

echo "================================================"
echo "  Neptune AI Platform - Starting Services"
echo "================================================"

# 启动 nginx (前端)
echo "[1/2] Starting nginx (web frontend) ..."
nginx

# 启动后端服务
echo "[2/2] Starting neptune server ..."
cd /app && ./server
