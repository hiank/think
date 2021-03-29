# think

服务框架，protobuf3，通过特定协议名称，自定义的处理方式，自动定向到期望的endpoint上。框架是为kubernetes设计的，目前提供websocket(提供给客户端)，grpc(集群内部使用)，nats(消息中间件，集群内部使用)
典型的连接方式:
