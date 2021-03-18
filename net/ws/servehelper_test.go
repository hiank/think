package ws

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	td "github.com/hiank/think/net/ws/testdata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type testConn struct {
	key string
	net.Sender
	net.Reciver
}

func (tc *testConn) Key() string {
	return tc.key
}

func (tc *testConn) Close() error {
	return nil
}

func testDial(addr string, token string) (*websocket.Conn, *http.Response, error) {

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	return websocket.DefaultDialer.Dial(u.String(), http.Header{"token": {token}})
}

// var testDialer = websocket.Dialer{

// }

func TestNewServeHelper(t *testing.T) {

	helper := NewServeHelper(":10225")
	assert.Assert(t, helper.server != nil)
	assert.Assert(t, helper.upgrader != nil)
	assert.Assert(t, helper.connChan != nil)
}

func TestServeHelperAccept(t *testing.T) {

	helper := NewServeHelper(":10225")

	t.Run("ok", func(t *testing.T) {

		go func() {
			helper.connChan <- &testConn{}
		}()
		conn, err := helper.Accept()
		assert.Assert(t, err == nil, err)
		assert.Assert(t, conn != nil)
	})

	t.Run("close", func(t *testing.T) {

		go func() {
			helper.Close()
		}()

		conn, err := helper.Accept()
		assert.Equal(t, err, io.EOF, err)
		assert.Assert(t, conn == nil)
	})
}

func TestServeHelperAuth(t *testing.T) {

	// assert
	helper := NewServeHelper(":10225")
	uid, pass := helper.auth("any")
	assert.Assert(t, pass)
	assert.Equal(t, uid, uint64(1001))
}

func TestServeHelperListenAndServe(t *testing.T) {

	helper := NewServeHelper(":10225")

	t.Run("close", func(t *testing.T) {

		go func() {
			helper.Close()
		}()

		err := helper.ListenAndServe()
		assert.Equal(t, err, http.ErrServerClosed, err)
	})
}

func TestServeHelperServeHTTP(t *testing.T) {

	helper := NewServeHelper(":10225")
	defer helper.Close()

	go helper.ListenAndServe()

	t.Run("dial ok", func(t *testing.T) {
		conn, response, err := testDial("localhost:10225", "HOPE")
		assert.Assert(t, err == nil, err)
		assert.Assert(t, conn != nil)
		assert.Assert(t, response != nil)
		helper.Accept()
		go conn.Close()
	})

	t.Run("dial close", func(t *testing.T) {

		wait := make(chan bool)
		go func(t *testing.T) {
			clientConn, _, _ := testDial("localhost:10225", "HOPE")
			clientConn.Close()
			close(wait)
		}(t)
		serverConn, err := helper.Accept()
		assert.Assert(t, err == nil, err)

		msg, err := serverConn.Recv()
		assert.Assert(t, msg == nil)
		assert.Assert(t, websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err), err)

		<-wait
	})

	t.Run("send-recv", func(t *testing.T) {

		wait := make(chan bool)
		go func() {
			conn, _, _ := testDial("localhost:10225", "HOPE")
			anyData, _ := anypb.New(&td.Request{Value: "1"})
			buf, _ := proto.Marshal(anyData)
			err := conn.WriteMessage(websocket.BinaryMessage, buf)
			assert.Assert(t, err == nil, err)

			_, buf, err = conn.ReadMessage()
			assert.Assert(t, err == nil, err)
			// ptypes.UnmarshalAny()
			anyData = new(anypb.Any)
			err = proto.Unmarshal(buf, anyData)
			assert.Assert(t, err == nil, err)

			var recvMsg td.Response
			err = anyData.UnmarshalTo(&recvMsg)
			assert.Assert(t, err == nil, err)
			assert.Equal(t, recvMsg.GetValue(), "11")

			conn.Close()
			close(wait)
		}()

		conn, _ := helper.Accept()
		msg, err := conn.Recv()
		assert.Assert(t, err == nil, err)

		var data td.Request
		err = msg.GetValue().UnmarshalTo(&data)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, data.GetValue(), "1")

		anyData, _ := anypb.New(&td.Response{Value: data.GetValue() + "1"})
		// buf, _ := proto.Marshal(anyData)
		err = conn.Send(&pb.Message{Key: "HOPE", Value: anyData})
		assert.Assert(t, err == nil, err)

		<-wait
	})
}
