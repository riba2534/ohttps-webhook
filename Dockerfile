# 使用Go官方的基础镜像
FROM golang:1.20 AS builder

# 设置工作目录
WORKDIR /app

# 将你的go.mod 和 go.sum 文件复制到工作目录中
COPY go.mod go.sum ./

# 下载所有依赖项
RUN go mod download

# 将你的源代码复制到工作目录中
COPY . .

# 编译你的Go应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用 scratch 作为基础镜像
FROM scratch

# 设置工作目录
WORKDIR /root/

# 从builder镜像中，将编译好的Go应用复制到当前镜像中
COPY --from=builder /app/main .

# 设置应用为镜像的启动程序
CMD ["./main"]
