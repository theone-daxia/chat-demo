# chat-demo
一个聊天demo（gin + websocket + mysql + redis + mongodb）

支持的功能
* 用户注册
* 双人单聊
* 历史消息查询

## 运行
必备项
* go环境
* DockerDesktop(分别拉取 redis、mysql、mongo 镜像并运行)

准备好必备项后，即可通过如下命令启动服务
```shell
go mod tidy
go run main.go
```

## 参考项目
* [github.com/CocaineCong/gin-chat-demo](https://github.com/CocaineCong/gin-chat-demo)
* [github.com/gorilla/websocket@v1.5.0/examples/chat](https://github.com/gorilla/websocket)