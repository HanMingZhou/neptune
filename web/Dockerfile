# ---- 构建阶段 ----
FROM node:20-alpine AS builder

WORKDIR /build
COPY package.json package-lock.json ./
RUN npm config set registry https://registry.npmmirror.com \
    && npm config set fetch-retries 5 \
    && npm config set fetch-retry-mintimeout 20000 \
    && npm config set fetch-retry-maxtimeout 120000 \
    && npm config set prefer-offline false \
    && npm cache clean --force \
    && npm install --no-fund --no-audit

COPY . .

# 生产环境变量
ENV VITE_BASE_API=/aiInfra
RUN npm run build

# ---- 运行阶段 ----
FROM nginxinc/nginx-unprivileged:stable-alpine

USER root

# 移除默认配置
RUN rm /etc/nginx/conf.d/default.conf

COPY --from=builder /build/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN chown -R 101:0 /usr/share/nginx/html /etc/nginx/conf.d \
    && chmod -R g=u /usr/share/nginx/html /etc/nginx/conf.d

USER 101

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
