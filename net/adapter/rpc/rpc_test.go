package rpc_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net/adapter/rpc"
	"github.com/hiank/think/net/adapter/rpc/pipe"
	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/testdata"
	"github.com/hiank/think/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
)

type tmpREST struct {
	pipe.UnimplementedRestServer
}

func (tr *tmpREST) Get(ctx context.Context, req *anypb.Any) (out *anypb.Any, err error) {
	msg, _ := req.UnmarshalNew()
	if _, ok := msg.(*testdata.G_Example); ok {
		out, _ = anypb.New(&testdata.AnyTest1{Name: "resp"})
	} else {
		err = status.Errorf(codes.InvalidArgument, "request for 'Get' must be a 'G_Example'")
	}
	return
}

//NOTE: if err != nil, out would be nil forever
func (tr *tmpREST) Post(ctx context.Context, v *anypb.Any) (out *emptypb.Empty, err error) {
	msg, _ := v.UnmarshalNew()
	if _, ok := msg.(*testdata.P_Example); !ok {
		err = status.Errorf(codes.InvalidArgument, "resqust for 'Post' must be a 'P_Example'")
	}
	return new(emptypb.Empty), err
}

// func TestWithDefaultListenOption(t *testing.T) {
// 	opt := rpc.Export_defaultListener()
// 	assert.Equal(t, opt.Addr, "11")
// 	assert.DeepEqual(t, opt.Rest, new(pipe.UnimplementedPipeServer))

// 	opt = rpc.Export_withDefaultListenOption(rpc.ListenOption{Rest: &tmpREST{}})
// 	assert.Equal(t, opt.Addr, "")
// 	assert.DeepEqual(t, opt.Rest, &tmpREST{})
// }

func easyDial() pipe.KeepaliveClient {
	cc, _ := grpc.Dial("localhost:10250", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	return pipe.NewKeepaliveClient(cc)
	// any, _ := anypb.New(&testdata.G_Example{Value: "req"})
}

// type convertoCloser func() error

// func (cc convertoCloser) Close() error {
// 	return cc()
// }

func TestListener(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	t.Run("invalid addr", func(t *testing.T) {
		defer func() {
			r := recover()
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.Assert(t, err.Error() != "")
		}()
		rpc.NewListener(ctx, rpc.WithAddress("invalid"))
	})

	// t.Run("ctx cancel", func(t *testing.T) {
	// 	ctx, cancel := context.WithCancel(ctx)
	// 	lis := rpc.NewListener(ctx, rpc.WithAddress(":10250")) //rpc.ListenOption{Addr: ":10250"})
	// 	l := rpc.Export_convertolistener(lis)
	// 	oldCloser := l.Closer
	// 	closeCnt := 0
	// 	// closeDone := make(chan bool, 8)
	// 	l.Closer = convertoCloser(func() error {
	// 		closeCnt++
	// 		// closeDone <- true
	// 		return oldCloser.Close()
	// 	})

	// 	cancel()
	// 	// <-closeDone
	// 	<-time.After(time.Millisecond * 10)
	// 	assert.Equal(t, closeCnt, 1)
	// })

	t.Run("close", func(t *testing.T) {
		// ctx, cancel := context.WithCancel(ctx)
		lis := rpc.NewListener(ctx, rpc.WithAddress(":10250")) //rpc.ListenOption{Addr: ":10250"})
		l := rpc.Export_convertolistener(lis)
		oldCloser := l.Closer
		closeCnt := 0
		// closeDone := make(chan bool, 8)
		l.Closer = run.NewOnceCloser(func() error {
			closeCnt++
			return oldCloser.Close()
		})

		l.Close()
		// <-closeDone
		<-time.After(time.Millisecond * 10)
		assert.Equal(t, closeCnt, 1)
	})

	t.Run("serve stop", func(t *testing.T) {
		l, _, srv := rpc.Export_NewListenerEx(ctx, rpc.WithAddress(":10250")) //rpc.ListenOption{Addr: ":10250"})
		srv.Stop()

		<-time.After(time.Millisecond * 10)
		_, ok := <-l.ChanAccepter
		assert.Equal(t, ok, false, "closed")
	})

	t.Run("listen close", func(t *testing.T) {
		l, lis, _ := rpc.Export_NewListenerEx(ctx, rpc.WithAddress(":10250")) //rpc.ListenOption{Addr: ":10250"})
		// srv.Stop()
		lis.Close()

		// <-time.After(time.Millisecond * 10)
		_, ok := <-l.ChanAccepter
		assert.Equal(t, ok, false, "closed")
	})

	// lis.Close()
	// cancel()
	// <-time.After(time.Second)
	lis := rpc.NewListener(ctx, rpc.WithAddress(":10250"), rpc.WithDefaultKeepaliveServer()) //rpc.ListenOption{Addr: ":10250"})
	pc := easyDial()
	lc, _ := pc.Link(ctx)
	// assert.Assert(t, err != nil)
	h, _ := lc.Header()
	assert.Equal(t, len(h["success"]), 0)

	lc, _ = pc.Link(metadata.NewOutgoingContext(ctx, metadata.Pairs("identity", "111")))
	h, _ = lc.Header()
	assert.Equal(t, len(h["success"]), 1)

	sc, err := lis.Accept()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, sc.Token.Value(box.ContextkeyTokenUid).(string), "111")

	m := box.New(box.WithMessageValue(&testdata.S_Example{Value: "get"}))
	err = lc.Send(m.GetAny())
	assert.Equal(t, err, nil, err)
	// lc.CloseSend()
	m, _ = sc.T.Recv()
	v, _ := m.GetAny().UnmarshalNew()
	assert.Equal(t, v.(*testdata.S_Example).GetValue(), "get")

	m = box.New(box.WithMessageValue(&testdata.S_Example{Value: "resp"}))
	err = sc.T.Send(m)
	assert.Equal(t, err, nil, err)

	amsg, err := lc.Recv()
	assert.Equal(t, err, nil, err)
	v, _ = amsg.UnmarshalNew()
	assert.Equal(t, v.(*testdata.S_Example).GetValue(), "resp")

	err = lc.CloseSend()
	assert.Equal(t, err, nil, err)

	_, err = sc.T.Recv()
	assert.Equal(t, err, io.EOF, err)
}

func TestListenerWithRest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis := rpc.NewListener(ctx, rpc.WithAddress(":10250"), rpc.WithRestServer(&tmpREST{})) //rpc.ListenOption{Addr: ":10250"})
	// pc := easyDial()
	cc, _ := grpc.Dial("localhost:10250", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	rcli := pipe.NewRestClient(cc)
	// lc, _ := pc.Link(ctx)
	// assert.Assert(t, err != nil)
	// h, _ := lc.Header()
	// assert.Equal(t, len(h["success"]), 0)
	ag, _ := anypb.New(&testdata.G_Example{})
	am, err := rcli.Get(ctx, ag)
	assert.Equal(t, err, nil)
	v, _ := am.UnmarshalNew()
	assert.Equal(t, v.(*testdata.AnyTest1).GetName(), "resp")

	ap, _ := anypb.New(&testdata.P_Example{})
	_, err = rcli.Post(ctx, ap)
	assert.Equal(t, err, nil)

	_, err = rcli.Get(ctx, ap)
	assert.Assert(t, err != nil, "must be G_Example")

	_, err = rcli.Post(ctx, ag)
	assert.Assert(t, err != nil, "must be P_Example")

	lis.Close()
	_, err = rcli.Get(ctx, ag)
	assert.Assert(t, err != nil, "serve closed")
}

func TestRestDial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis := rpc.NewListener(ctx, rpc.WithAddress(":10250"), rpc.WithDefaultKeepaliveServer()) //rpc.ListenOption{Addr: ":10250"})

	rc, err := rpc.RestDial(ctx, "localhost:10250")
	assert.Equal(t, err, nil, err)
	ag, _ := anypb.New(&testdata.G_Example{})
	_, err = rc.Get(ctx, ag)
	assert.Assert(t, err != nil, "method Get not implemented")

	ap, _ := anypb.New(&testdata.P_Example{})
	_, err = rc.Post(ctx, ap)
	assert.Assert(t, err != nil, "method Post not implemented")

	lis.Close()
	rpc.NewListener(ctx, rpc.WithAddress(":10250"), rpc.WithRestServer(&tmpREST{}))

	amsg, err := rc.Get(ctx, ag)
	assert.Equal(t, err, nil, "这里要注意，grpc服务关闭后重写启动，不会影响原有的非长连接")
	v1, _ := amsg.UnmarshalNew()
	assert.Equal(t, v1.(*testdata.AnyTest1).GetName(), "resp")

	_, err = rc.Post(ctx, ap)
	assert.Equal(t, err, nil)
}

func TestKeepaliveDial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis := rpc.NewListener(ctx, rpc.WithAddress(":10250")) //rpc.ListenOption{Addr: ":10250"})
	// <-time.After(time.Millisecond * 100)

	_, err := rpc.NewKeepaliveDialer().Dial(ctx, "localhost:10250")
	assert.Equal(t, err, rpc.ErrLinkAuthFailed, err)
	// ag, _ := anypb.New(&testdata.G_Example{})

	dialer := rpc.NewKeepaliveDialer(rpc.WithIdentity("110"))
	_, err = dialer.Dial(ctx, "localhost:10250")
	assert.Equal(t, err, rpc.ErrLinkAuthFailed, err)

	lis.Close()

	lis = rpc.NewListener(ctx, rpc.WithAddress(":10250"), rpc.WithDefaultKeepaliveServer())

	rc, err := dialer.Dial(ctx, "localhost:10250")
	assert.Equal(t, err, nil, err)

	m := box.New(box.WithMessageValue(&testdata.G_Example{Value: "give"}))
	err = rc.Send(m)
	assert.Assert(t, err == nil)

	c, err := lis.Accept()
	assert.Equal(t, err, nil)
	m1, err := c.T.Recv()
	assert.Equal(t, err, nil)
	v1, _ := m1.GetAny().UnmarshalNew()

	assert.Equal(t, v1.(*testdata.G_Example).GetValue(), "give")

	m = box.New(box.WithMessageValue(&testdata.AnyTest1{Name: "ata"}))
	err = c.T.Send(m)
	assert.Equal(t, err, nil, err)

	m1, err = rc.Recv()
	assert.Equal(t, err, nil, err)
	v1, _ = m1.GetAny().UnmarshalNew()
	assert.Equal(t, v1.(*testdata.AnyTest1).GetName(), "ata")

	rc.Close()
	_, err = c.T.Recv()
	assert.Assert(t, err != nil)

	t.Run("server conn closed", func(t *testing.T) {
		rc, err := dialer.Dial(ctx, "localhost:10250")
		assert.Equal(t, err, nil, err)

		c, err := lis.Accept()
		assert.Equal(t, err, nil)

		c.T.Close()
		_, err = rc.Recv()
		assert.Assert(t, err != nil)
	})
}

// func TestDialer(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.TODO())
// 	defer cancel()

// 	lis := rpc.NewListener(ctx, rpc.ListenOption{Addr: ":10250"})

// 	dialer := rpc.NewDialer(ctx)
// 	c, err := dialer.Dial("localhost:10250")
// 	assert.Equal(t, err, nil, err)
// }
