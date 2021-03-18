package ws

import (
	"bytes"
	"testing"

	"github.com/hiank/think/net/pb"
	td "github.com/hiank/think/net/ws/testdata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type testRWC struct {
	ch chan []byte
}

func (rwc *testRWC) ReadMessage() (messageType int, buf []byte, err error) {
	buf = <-rwc.ch
	return
}

func (rwc *testRWC) WriteMessage(messageType int, buf []byte) error {
	rwc.ch <- buf
	return nil
}

func (rwc *testRWC) Close() error {
	close(rwc.ch)
	return nil
}

func TestConn(t *testing.T) {
	t.Run("Key", func(t *testing.T) {
		c := &conn{ReadWriteCloser: nil, uid: 1001}
		assert.Equal(t, c.Key(), "1001")
	})

	t.Run("Recv", func(t *testing.T) {

		ch := make(chan []byte, 1)
		c := &conn{ReadWriteCloser: &testRWC{ch: ch}}

		anyData, _ := anypb.New(&td.Request{Value: "110"})
		buf, _ := proto.Marshal(anyData)
		ch <- buf

		var val td.Request
		msg, _ := c.Recv()
		msg.GetValue().UnmarshalTo(&val)
		assert.Equal(t, val.GetValue(), "110")
	})

	t.Run("Send", func(t *testing.T) {

		ch := make(chan []byte, 1)
		c := &conn{ReadWriteCloser: &testRWC{ch: ch}}

		anyMsg, _ := anypb.New(&td.Response{Value: "loveWS"})
		c.Send(&pb.Message{Value: anyMsg})

		buf := <-ch
		anyBuf, _ := proto.Marshal(anyMsg)
		assert.Assert(t, bytes.Equal(buf, anyBuf))
	})
}
