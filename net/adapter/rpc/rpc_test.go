package rpc_test

import (
	"context"
	"io"
	"testing"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter/rpc"
	"github.com/hiank/think/net/adapter/rpc/pipe"
	"github.com/hiank/think/pbtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
)

type tmpREST struct {
	pipe.UnimplementedRestServer
}

func (tr *tmpREST) Get(ctx context.Context, req *anypb.Any) (out *anypb.Any, err error) {
	msg, _ := req.UnmarshalNew()
	if _, ok := msg.(*pbtest.G_Example); ok {
		out, _ = anypb.New(&pbtest.AnyTest1{Name: "resp"})
	} else {
		err = status.Errorf(codes.InvalidArgument, "request for 'Get' must be a 'G_Example'")
	}
	return
}

// NOTE: if err != nil, out would be nil forever
func (tr *tmpREST) Post(ctx context.Context, v *anypb.Any) (out *emptypb.Empty, err error) {
	msg, _ := v.UnmarshalNew()
	if _, ok := msg.(*pbtest.P_Example); !ok {
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
	cc, _ := grpc.Dial("localhost:30250", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	return pipe.NewKeepaliveClient(cc)
	// any, _ := anypb.New(&pbtest.G_Example{Value: "req"})
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

	lis := rpc.NewListener(ctx, rpc.WithAddress(":30250"), rpc.WithServeKeepalive(rpc.Tokenset))
	defer func() {
		lis.Close()
	}()
	pc := easyDial()
	lc, _ := pc.Link(ctx)
	h, _ := lc.Header()
	assert.Equal(t, len(h["success"]), 0)

	lc, _ = pc.Link(metadata.NewOutgoingContext(ctx, metadata.Pairs("identity", "111")))
	h, _ = lc.Header()
	assert.Equal(t, len(h["success"]), 1)

	sc, err := lis.Accept()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, sc.Token().ToString(), "111")

	m := net.NewMessage(net.WithMessageValue(&pbtest.S_Example{Value: "get"}))
	err = lc.Send(m.Any())
	assert.Equal(t, err, nil, err)
	// lc.CloseSend()
	m, _ = sc.Recv()
	v, _ := m.Any().UnmarshalNew()
	assert.Equal(t, v.(*pbtest.S_Example).GetValue(), "get")

	m = net.NewMessage(net.WithMessageValue(&pbtest.S_Example{Value: "resp"}), net.WithMessageToken(rpc.Tokenset.Derive("any")))
	err = sc.Send(m)
	assert.Equal(t, err, nil, err)

	amsg, err := lc.Recv()
	assert.Equal(t, err, nil, err)
	v, _ = amsg.UnmarshalNew()
	assert.Equal(t, v.(*pbtest.S_Example).GetValue(), "resp")

	err = lc.CloseSend()
	assert.Equal(t, err, nil, err)

	_, err = sc.Recv()
	assert.Equal(t, err, io.EOF, err)
}

func TestListenerWithRest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis := rpc.NewListener(ctx, rpc.WithAddress(":30250"), rpc.WithRestServer(&tmpREST{})) //rpc.ListenOption{Addr: ":30250"})
	// pc := easyDial()
	cc, _ := grpc.Dial("localhost:30250", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	rcli := pipe.NewRestClient(cc)
	// lc, _ := pc.Link(ctx)
	// assert.Assert(t, err != nil)
	// h, _ := lc.Header()
	// assert.Equal(t, len(h["success"]), 0)
	ag, _ := anypb.New(&pbtest.G_Example{})
	am, err := rcli.Get(ctx, ag)
	assert.Equal(t, err, nil)
	v, _ := am.UnmarshalNew()
	assert.Equal(t, v.(*pbtest.AnyTest1).GetName(), "resp")

	ap, _ := anypb.New(&pbtest.P_Example{})
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

	lis := rpc.NewListener(ctx, rpc.WithAddress(":30250"), rpc.WithServeKeepalive(rpc.Tokenset)) //rpc.ListenOption{Addr: ":30250"})

	rc, err := rpc.RestDial(ctx, "localhost:30250")
	assert.Equal(t, err, nil, err)
	ag, _ := anypb.New(&pbtest.G_Example{})
	_, err = rc.Get(ctx, ag)
	assert.Assert(t, err != nil, "method Get not implemented")

	ap, _ := anypb.New(&pbtest.P_Example{})
	_, err = rc.Post(ctx, ap)
	assert.Assert(t, err != nil, "method Post not implemented")

	lis.Close()
	rpc.NewListener(ctx, rpc.WithAddress(":30250"), rpc.WithRestServer(&tmpREST{}))

	amsg, err := rc.Get(ctx, ag)
	assert.Equal(t, err, nil, "这里要注意，grpc服务关闭后重写启动，不会影响原有的非长连接")
	v1, _ := amsg.UnmarshalNew()
	assert.Equal(t, v1.(*pbtest.AnyTest1).GetName(), "resp")

	_, err = rc.Post(ctx, ap)
	assert.Equal(t, err, nil)
}

func TestKeepaliveDial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis := rpc.NewListener(ctx, rpc.WithAddress(":30250")) //rpc.ListenOption{Addr: ":30250"})

	dialer := rpc.NewKeepaliveDialer(rpc.Tokenset.Derive("110"))
	_, err := dialer.Dial(ctx, "localhost:30250")
	assert.Equal(t, err, rpc.ErrLinkAuthFailed, err)

	lis.Close()

	lis = rpc.NewListener(ctx, rpc.WithAddress(":30250"), rpc.WithServeKeepalive(rpc.Tokenset))
	defer func() {
		lis.Close()
	}()

	// dialer := rpc.NewKeepaliveDialer(rpc.Tokenset.Derive("110"))
	rc, err := dialer.Dial(ctx, "localhost:30250")
	assert.Equal(t, err, nil, err)

	// net.NewMessage(net.WithMessageValue(&pbtest.G_Example{Value: "give"}))
	// m := box.New(box.WithMessageValue(&pbtest.G_Example{Value: "give"}))
	err = rc.Send(net.NewMessage(net.WithMessageValue(&pbtest.G_Example{Value: "give"}), net.WithMessageToken(rpc.Tokenset.Derive("110"))))
	assert.Assert(t, err == nil)

	c, err := lis.Accept()
	assert.Equal(t, err, nil)
	m1, err := c.Recv()
	assert.Equal(t, err, nil)
	v1, _ := m1.Any().UnmarshalNew()

	assert.Equal(t, v1.(*pbtest.G_Example).GetValue(), "give")

	// m = box.New(box.WithMessageValue(&pbtest.AnyTest1{Name: "ata"}))
	err = c.Send(net.NewMessage(net.WithMessageValue(&pbtest.AnyTest1{Name: "ata"}), net.WithMessageToken(rpc.Tokenset.Derive("110"))))
	assert.Equal(t, err, nil, err)

	m1, err = rc.Recv()
	assert.Equal(t, err, nil, err)
	v1, _ = m1.Any().UnmarshalNew()
	assert.Equal(t, v1.(*pbtest.AnyTest1).GetName(), "ata")

	rc.Close()
	_, err = c.Recv()
	assert.Assert(t, err != nil)

	t.Run("server conn closed", func(t *testing.T) {
		rc, err := dialer.Dial(ctx, "localhost:30250")
		assert.Equal(t, err, nil, err)

		c, err := lis.Accept()
		assert.Equal(t, err, nil)

		c.Close()
		_, err = rc.Recv()
		assert.Assert(t, err != nil)
	})
}
