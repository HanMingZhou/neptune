# ---- 构建阶段 ----
FROM node:20-alpine AS builder

WORKDIR /build
COPY package.json package-lock.json ./
RUN npm ci --registry=https://registry.npmmirror.com

COPY . .

# 生产环境变量
ENV VITE_BASE_API=/aiInfra
RUN npm run build

# ---- 运行阶段 ----
FROM nginx:stable-alpine

# 移除默认配置
RUN rm /etc/nginx/conf.d/default.conf

COPY --from=builder /build/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
