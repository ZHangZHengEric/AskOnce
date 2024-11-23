FROM repo.atomecho.cn/base/golang:1.22
# 将本地的./config目录映射到容器中的/data/config
WORKDIR /app/
COPY  ./main /app/
ENV DOCKER_ENV=online
EXPOSE 8080
ENTRYPOINT ["./main"]