# 使用 Node.js 18 作为基础镜像
FROM node:18-alpine AS base

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 配置国内镜像源
RUN npm config set registry 'https://registry.npmmirror.com/'

# 安装pnpm
RUN npm install -g pnpm

# 安装依赖
RUN pnpm install

# 复制项目文件到容器
COPY . .

# 构建项目
RUN npm run build

# 使用 nginx 提供静态文件服务
FROM nginx:stable-alpine

# 复制打包后的文件到 nginx 的 html 目录
COPY --from=base /app/dist /usr/share/nginx/html/dist

# 替换默认 nginx 配置
COPY askonce.conf /etc/nginx/conf.d/askonce.conf

# 暴露服务端口
EXPOSE 80

# 启动 nginx 服务
CMD ["nginx", "-g", "daemon off;"]
