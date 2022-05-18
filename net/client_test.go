package net_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/testdata"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type tmpKnower struct{}

func (tmpKnower) Identity(uid string) (id string, err error) {
	if uid == "" || uid == "0" {
		return "", fmt.Errorf("invalid uid:%s", uid)
	}
	return uid, nil
}

func (tmpKnower) Uid(id string) (uid string, err error) {
	if id == "" || id == "0" {
		return "", fmt.Errorf("invalid id:%s", id)
	}
	return id, nil
}

func (tmpKnower) ServeAddr(m *box.Message) (addr string, err error) {
	return "tmpServe", nil
}

type tmpDialer struct {
}

func (*tmpDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return nil, nil
}

func TestClientRouteMux(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli := net.NewClient(ctx, nil, nil)
	cli.RouteMux().Handle("AnyTest1", net.HandlerFunc(func(id string, m *box.Message) {

	}))
}

type tmpRest struct{}

func (tmpRest) Get(ctx context.Context, req *anypb.Any) (rsp *anypb.Any, err error) {
	if req == nil {
		err = fmt.Errorf("failed")
	} else {
		rsp, err = anypb.New(&testdata.AnyTest1{Name: "success"})
	}
	return
}

func (tmpRest) Post(ctx context.Context, req *anypb.Any) (out *emptypb.Empty, err error) {
	if req == nil {
		err = fmt.Errorf("failed")
	} else {
		out = &emptypb.Empty{}
	}
	return
}

// func TestClientConn(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	t.Run("unimplement", func(t *testing.T) {
// 		cc := net.NewClientConn()
// 		assert.Equal(t, cc.Close(), net.ErrUnimplementApi)
// 		assert.Equal(t, cc.Send(nil), net.ErrUnimplementApi)
// 		_, err := cc.Recv()
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 		_, err = cc.Get(ctx, nil)
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 		_, err = cc.Post(ctx, nil)
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 	})
// 	t.Run("only conn", func(t *testing.T) {
// 		rpp, spp := make(chan *box.Message, 1), make(chan *box.Message, 1)
// 		cc := net.NewClientConn(net.WithConn(net.Export_newTmpConn(rpp, spp)))
// 		_, err := cc.Get(ctx, nil)
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 		_, err = cc.Post(ctx, nil)
// 		assert.Equal(t, err, net.ErrUnimplementApi)

// 		m, _ := box.New(&testdata.AnyTest1{Name: "a1"})
// 		err = cc.Send(m)
// 		assert.Equal(t, err, nil, err)
// 		m1 := <-spp
// 		v, _ := m1.GetAny().UnmarshalNew()
// 		assert.Equal(t, v.(*testdata.AnyTest1).GetName(), "a1")

// 		m, _ = box.New(&testdata.AnyTest2{Hope: "h2"})
// 		rpp <- m
// 		m2, err := cc.Recv()
// 		assert.Equal(t, err, nil, err)
// 		v, _ = m2.GetAny().UnmarshalNew()
// 		assert.Equal(t, v.(*testdata.AnyTest2).GetHope(), "h2")

// 		err = cc.Close()
// 		assert.Equal(t, err, nil, err)

// 		_, err = cc.Recv()
// 		assert.Equal(t, err, io.EOF)
// 	})
// 	t.Run("only Rest", func(t *testing.T) {
// 		cc := net.NewClientConn(net.WithRest(tmpRest{}))
// 		assert.Equal(t, cc.Close(), net.ErrUnimplementApi)
// 		assert.Equal(t, cc.Send(nil), net.ErrUnimplementApi)
// 		_, err := cc.Recv()
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 		_, err = cc.Get(ctx, new(anypb.Any))
// 		assert.Equal(t, err, nil)
// 		_, err = cc.Post(ctx, new(anypb.Any))
// 		assert.Equal(t, err, nil)
// 	})
// }
