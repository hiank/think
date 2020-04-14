# Pool

管理连接，消息的收发。

## Quick View

- [conn](#conn)
- [connHub](#connHub)
- [Pool](#pool)

## conn

每个`Conn` 包含一个`token.Token`，用于维护生命周期，当`token.Token` 释放时，`Conn` 将执行相应的资源释放。

## connHub

维护管理`Conn`，提供`add` `del` `find` `send` 功能

## Pool

提供各类api，用于添加conn，判断conn是否存在，启动conn监听，发送消息

### API

- `NewPool(context.Context) `