# Easy client-server

## Client

1. operate a `box.Message` `m` use `Auto`
2. unmarshal `m` then get a key for found `Conn` to send `m`
3. cannot found the `Conn`, create one
4. ...

### taskConn

Conn use Tasker to sequential execute `Dial` and `Send box.Message`

- first Task is try dial
- when dial closed

`newTaskConn(ctx context.Context, dial Dial, doneHook func()) Conn` to make a new taskConn

exception enumeration:

- context passed canceled: doneHook will respond

- dial error

### liteConn & initialize

带缓存的`Conn`，无阻塞发送消息。需前置执行`initialize`，以执行连接(连接设置为第一个任务，完成前所有期望发送的消息都会被缓存起来，等待任务的顺序执行)。`initialize`还会开启一个循环`Recv`的goroutine，这个`Recv`时完成连接后得到的`Conn`，而非本身。本身持有的`Recv`为未实现api，因为此调用非异步安全，因此不提供外部访问。循环读取函数作为`initialize`参数传入。

### liteConn: Conn used by connset

- `Conn` with message cache.
- non-blocking `Send` (use `run.Tasker`).
- delay connecting (by `initialize`).
- respond after closing (by `initialize`), `connset` use this feature to remove reference.
- `Recv` unimplemented. `initialize` will start loopRecv after connect success.

### initialize: initialize an empty *liteConn

param:

- `ctx context.Context`: root context, for liteConn's tasker and connect
- `lc *liteConn`: target *liteConn
- `connect Connect`: use to get a new conn
- `loopRecv func(Receiver, box.Token)`: loop recv message from newly obtained conn until recv any error
- `doneHook`: last executed after conn *liteConn closed

<!-- 不要异步调用`Recv`，因为`Receiver`会被重设，异步调用的`Recv`无法保证合理性。而`Close` -->
