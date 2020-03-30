# Pool

管理连接，消息的收发。避免并发爆炸

## Quick View

- [Conn](#conn)

- [ConnHub](#connHub)

- [Message](#message)

- [MessageHub](#messageHub)

- [Pool](#pool)

- [Timer](#timer)

## Conn

每个`Conn` 包含一个`token.Token`，用于维护生命周期，当`token.Token` 释放时，`Conn` 将执行相应的资源释放。

## ConnHub

## Message

## MessageHub

## Pool

## Timer