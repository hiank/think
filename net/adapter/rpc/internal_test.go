package rpc

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/auth"
	"github.com/hiank/think/net"
	"github.com/hiank/think/pbtest"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

var (
	Tokenset = auth.NewTokenset(context.Background())
)

type tmpSendReciver struct {
	chans chan<- *anypb.Any
	chanr <-chan *anypb.Any
}

func (tsr *tmpSendReciver) Send(am *anypb.Any) error {
	tsr.chans <- am
	return nil
}

func (tsr *tmpSendReciver) Recv() (amsg *anypb.Any, err error) {
	amsg, ok := <-tsr.chanr
	if !ok {
		err = io.EOF
	}
	return
}

func TestConn(t *testing.T) {
	///
	chans, chanr := make(chan *anypb.Any, 1), make(chan *anypb.Any, 1)
	c := &conn{
		tk: Tokenset.Derive("11"),
		sr: &tmpSendReciver{chans: chans, chanr: chanr},
	}
	////
	am, _ := anypb.New(&pbtest.AnyTest1{Name: "anyt1"})
	msg := net.NewMessage(net.WithMessageValue(am), net.WithMessageToken(c.tk.Fork(auth.WithTokenTimeout(time.Millisecond))))
	err := c.Send(msg)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(chans), 1)
	<-chans
	<-time.After(time.Millisecond * 20)
	err = c.Send(msg)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.Equal(t, len(chans), 0)

	////recv
	chanr <- am
	msg, err = c.Recv()
	assert.Equal(t, err, nil)
	assert.Equal(t, msg.Token().ToString(), "11")
	assert.Equal(t, string(msg.Any().MessageName()), "AnyTest1")
	an, _ := msg.Any().UnmarshalNew()
	assert.Equal(t, an.(*pbtest.AnyTest1).GetName(), "anyt1")
	///
	msg = net.NewMessage(net.WithMessageValue(am), net.WithMessageToken(Tokenset.Derive("long")))
	assert.Equal(t, c.Send(msg), nil)
	assert.Equal(t, len(chans), 1)
	<-chans
	c.Token().Close()
	assert.Equal(t, c.Send(msg), context.Canceled)
	assert.Equal(t, len(chans), 0)

	chanr <- am
	_, err = c.Recv()
	assert.Equal(t, err, context.Canceled)
	assert.Equal(t, len(chanr), 1)
	<-chanr
}

func TestFuncLinkServer(t *testing.T) {

}
