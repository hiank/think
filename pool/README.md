# Introduce

---

简单介绍下功能，设计思路，部分API说明。

## Pool

- '集成'一个ConnHub
- 拥有两个MessageHub，分别处理读到的Message，需要写入的Message
- 定期清理过期的Conn
- `API: Listen` 监听Conn，起一个Conn读线程，并将Conn加入ConnHub
- `API: Post` 传入待发送的Message

## ConnHub

- 根据key, id 维护Conn。key为服务名，id 为Conn识别号
- map + list 维护Conn，对每次使用的Conn，会移到list 末尾用于提高Conn清理效率
- `API: Upgrade` 清理Conn，清除超时的Conn

## MessageHub

设计这个Hub的目的是限制处理Message 的goruntine 数量

- 使用list 维护待处理的Message
- 限定数量的goruntine 处理Message，每个Message 处理完毕，当前goruntine 检出list 的下一个Message，没有的话结束当前goruntine

## Conn

- 有个唯一ID
