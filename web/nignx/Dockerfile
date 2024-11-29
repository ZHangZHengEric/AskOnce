FROM openresty/openresty:alpine AS builder

LABEL stage=openresty

WORKDIR /app

COPY . .
COPY ./nginx/nginx.conf /usr/local/openresty/nginx/conf/nginx.conf

RUN mkdir -p /usr/local/openresty/mylua