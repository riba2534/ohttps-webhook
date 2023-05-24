# 使用 scratch 作为基础镜像
FROM scratch

# 设置工作目录
WORKDIR /root/

# 将你的 Go 应用复制到镜像中
COPY main .

# 设置应用为镜像的启动程序
CMD ["./main"]