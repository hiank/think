# Pool

管理连接，消息的收发。

## Quick View

- [Conn](#conn)
    - [Listen](#listen)
    - [Handle](#handle)
    - [Send](#send)
- [MessageHub](#messageHub)
    - [DoActive](#doActive)
    - [PushWithBack](#pushWithBack)
    - [Push](#push)
- [Pool](#pool)
    - [Listen](#listen)
    - [PostAndWait](#postAndWait)
    - [Post](#post)

## Conn

每个`Conn` 包含一个`token.Token`，用于维护生命周期，当`token.Token` 释放时，`Conn` 将执行相应的资源释放。

### Listen

- **Param**: `MessageHandler`
- **Return**: `error`

循环读取消息，会阻塞

### Handle

- **Param**: `*Message`
- **Return**: `error`

处理消息发送

## connHub

维护管理`Conn`，提供`add` `del` `find` `send` 功能，专门为`Pool`设计的结构，不能独立运作

## Pool

提供各类api，用于添加conn，判断conn是否存在，启动conn监听，发送消息

## MessageHub

处理`Message`，提供几点特性：
1. 可延时执行处理，某些情况下，处理方法可能还未准备好[例如：收到待处理消息时，不存在已建立的k8s连接，此时需要等待连接完成]，提供延迟处理，可以简化调用流程
2. 并发执行处理，并可限制同时处理goroutine数量
3. 当绑定的context关闭后，可以安全退出[主要是处理的goroutine能顺利退出]
4. 激活方法可被多次调用，并实际只有第一次调用是有效的
5. 提供Push 和 PushWithBack 方法，添加待处理消息

### DoActive

激活当前`MessageHub`的处理流程。此方法只会响应第一次调用。
- 第一次调用`DoActive`之前，`MessageHub`启用的是[loop](#loop)方法，监听各种消息
- 第一次嗲用`DoActive`之后，[loop](#loop)会退出，并启用[loopHandle](#loopHandle)方法，监听后续消息

### PushWithBack

- **Param**: `*Message`, `chan<- error`

向`MessageHub`中添加消息请求，并返回处理状态[是否出错]。注意，这个方法一般不会阻塞

### Push

- **Param**: `*Message`

向`MessageHub`中添加消息请求，不会返回处理状态

## Pool

连接池，集中管理Conn

### Listen

- **Param**: `*Token`, `IO`
- **Return**: `error`

构建一个`Conn`并启动监听。这个方法会阻塞。构建得到的`Conn`将存于`Pool`中进行维护，主要是发送消息时会匹配对应的`Conn`，调用其处理方法

### PostAndWait

- **Param**: `*Message`
- **Return**: `error`

发送消息，将消息送到`Pool`的`connHub`中，调用匹配的`Conn`的`Handle`方法。这个方法会阻塞知道处理完成，并返回处理状态[是否出错]

### Post

- **Param**: `*Message`

非阻塞的发送消息方法，将忽略发送状态[调用者不关系发送是否出错]