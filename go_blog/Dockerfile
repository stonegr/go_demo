# 第一阶段：构建阶段
FROM golang:1.24.1-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN go build -o main .

# 第二阶段：运行阶段
FROM alpine:latest

# 安装基本工具和证书
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为上海
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

# 暴露端口（根据您的应用需要修改）
EXPOSE 8080

# 运行应用
CMD ["./main"] 