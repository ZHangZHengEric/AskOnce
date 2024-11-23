FROM golang:1.22-bullseye as build

ENV GOPROXY=https://goproxy.cn,direct
ENV GIN_MODE=release

RUN sed  -i "s/deb.debian.org/mirrors.aliyun.com/g" /etc/apt/sources.list \
    && sed  -i "s/security.debian.org/mirrors.aliyun.com/g" /etc/apt/sources.list
# 更新APT包列表并安装证书相关工具和时区数据包
RUN apt update && apt install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    upx \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
    && echo Asia/Shanghai > /etc/timezone

# 移动到工作目录：/cm-mgmt
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
# 将代码复制到容器中
COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
# 编译成二进制可执行文件app
RUN go build --tags=kqueue,operator -ldflags "-w -s" -o dispatch


FROM ubuntu:22.04
ENV TZ=Asia/Shanghai
ENV GIN_MODE=release
WORKDIR /app
COPY --from=build /etc/localtime /etc/localtime
COPY --from=build /etc/timezone /etc/timezone
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /build/dispatch .
EXPOSE 8080
ENTRYPOINT ["./dispatch"]
