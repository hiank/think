package rpc_test

import (
	"context"
	"io"
	"testing"

	"github.com/hiank/think/net/adapter/rpc"
	"github.com/hiank/think/net/adapter/rpc/pp"
	"github.com/hiank/think/net/testdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
)

func TestListener(t *testing.T) {
	t.Run("NewListener-panic", func(t *testing.T) {
		defer func(t *testing.T) {
			r := recover()
			assert.Assert(t, r != nil)
		}(t)
		rpc.NewListener(context.Background(), rpc.WithAddr("invalid:port"))
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lis := rpc.NewListener(ctx, rpc.WithAddr(":30222"))
	defer lis.Close() //NOTE: release listen port immediately
	go func(t *testing.T) {
		cc, _ := grpc.Dial("localhost:30222", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		cli := pp.NewPipeClient(cc)
		any, _ := anypb.New(&testdata.AnyTest1{Name: "link"})
		lc, _ := cli.Link(ctx) //cli.Link(metadata.NewOutgoingContext(ctx, metadata.Pairs("identity", "111")))
		// t.Log(lc.Header())
		md, _ := lc.Header()
		assert.Equal(t, len(md.Get("success")), 0)

		lc, _ = cli.Link(metadata.NewOutgoingContext(ctx, metadata.Pairs("identity", "111")))
		md, _ = lc.Header()
		assert.Equal(t, len(md.Get("success")), 1)
		err := lc.Send(any)
		assert.Assert(t, err == nil)
		lc.CloseSend()

		// cc.Close()
	}(t)

	conn, _ := lis.Accept()
	any, err := conn.Recv()
	// t.Log(err)
	assert.Assert(t, err == nil)
	msg, _ := any.Any().UnmarshalNew()
	assert.Equal(t, msg.(*testdata.AnyTest1).GetName(), "link")

	_, err = conn.Recv()
	assert.Equal(t, err, io.EOF)

	// lis.Close()
	// <-time.After(time.Millisecond * 100) //NOTE: wait for listener released
}

type testREST struct {
}

func (tr *testREST) Get(ctx context.Context, req *anypb.Any) (out *anypb.Any, err error) {
	msg, _ := req.UnmarshalNew()
	if _, ok := msg.(*testdata.G_Example); ok {
		out, _ = anypb.New(&testdata.AnyTest1{Name: "resp"})
		// out = &pb.Carrier{Identity: carrier.GetIdentity(), Message: any}
	} else {
		err = status.Errorf(codes.InvalidArgument, "request for 'Get' must be a 'G_Example'")
	}
	return
}

//NOTE: if err != nil, out would be nil forever
func (tr *testREST) Post(ctx context.Context, v *anypb.Any) (out *emptypb.Empty, err error) {
	msg, _ := v.UnmarshalNew()
	if _, ok := msg.(*testdata.P_Example); !ok {
		// any, _ := anypb.New(&testdata.AnyTest2{})
		err = status.Errorf(codes.InvalidArgument, "resqust for 'Post' must be a 'P_Example'")
	}
	out = new(emptypb.Empty)
	return out, err
}

func TestGetPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lis := rpc.NewListener(ctx, rpc.WithAddr(":30222"), rpc.WithREST(&testREST{}))
	defer lis.Close()
	// defer lis.Close()
	// go func(t *testing.T) {
	cc, _ := grpc.Dial("localhost:30222", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	cli := pp.NewPipeClient(cc)
	any, _ := anypb.New(&testdata.G_Example{Value: "req"})

	rlt, err := cli.Get(ctx, any)
	assert.Equal(t, err, nil)
	// assert.Equal(t, rlt.GetIdentity(), uint64(11))
	msg, _ := rlt.UnmarshalNew()
	assert.Equal(t, msg.(*testdata.AnyTest1).GetName(), "resp")

	any2, _ := anypb.New(&testdata.P_Example{Value: "post"})
	rlt, err = cli.Get(ctx, any2)
	assert.Assert(t, rlt == nil)
	assert.Assert(t, err != nil)

	empty, err := cli.Post(ctx, any2)
	assert.Assert(t, empty != nil)
	assert.Assert(t, err == nil)

	empty, err = cli.Post(ctx, any)
	assert.Assert(t, empty == nil)
	assert.Assert(t, err != nil)

	// lis.Close()

	// <-time.After(time.Millisecond * 100) //NOTE: wait for listener released
}

func TestDialer(t *testing.T) {
	// d := rpc.NewDialer()
	// d.Dial("localhost:30222")
}
