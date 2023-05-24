# ohttps-webhook

# 介绍

本项目是 https://ohttps.com/ 的一个 webhook 部署节点，详情参考 https://ohttps.com/docs/cloud/webhook/webhook


# 使用

Docker镜像: https://hub.docker.com/r/riba2534/ohttps-webhook


```
docker run -p 4321:4321 -e TLS_PATH=/path -e CALLBACK_TOKEN=your_callback_token -v /host/path:/path riba2534/ohttps-webhook
```