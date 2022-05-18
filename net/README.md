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
