# 使用官方 Golang 镜像作为构建环境
FROM golang:1.23 AS build

# 设置 CGO_ENABLED 为 0 以进行静态编译
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# 设置工作目录
WORKDIR /app

# 复制代码文件
COPY gpt2.0/gpt2.0.go . 
COPY gpt2.0/.env .
COPY go.mod go.sum ./

# 下载依赖并编译应用程序
RUN go mod download
RUN go build -o server gpt2.0.go

# 使用更小的基础镜像
FROM alpine:latest
# 如果仍需要 CA 证书支持 HTTPS，可以安装它们
RUN apk --no-cache add ca-certificates

# 将编译好的二进制文件复制到最终镜像
COPY --from=build /app/server /server
COPY --from=build /app/.env /app/.env

# 暴露端口
EXPOSE 8080

# 启动服务
CMD ["/server"]
