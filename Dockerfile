# 使用 golang 官方的 alpine 镜像作为基础镜像
FROM golang:1.20-alpine

# 在容器内创建一个目录来存放应用代码
RUN mkdir /app

# 将工作目录切换到 /app
WORKDIR /app

# 将本地的代码复制到容器内的工作目录
COPY . .

# 下载依赖并构建应用
RUN go mod download
RUN go build -o main .

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]