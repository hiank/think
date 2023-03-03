package adapter_test

import (
	"context"
	"io"
	"testing"

	"github.com/hiank/think/auth"
	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter"
	"gotest.tools/v3/assert"
)

var (
	tokenset = auth.NewTokenset(context.Background())
)

type tmpConn struct {
	tk auth.Token
}

func (tc *tmpConn) Token() auth.Token {
	return tc.tk
}

func (tc *tmpConn) Send(*net.Message) error {
	return nil
}

func (tc *tmpConn) Recv() (*net.Message, error) {
	return nil, nil
}

func (tc *tmpConn) Close() error {
	return nil
}

func TestChanAccepter(t *testing.T) {
	ca := make(adapter.ChanAccepter)
	go func() {
		ca <- &tmpConn{tk: tokenset.Derive("test chan")} //net.TmpConn{Token: one.TokenSet().Derive("test chan")}
	}()

	ic, err := ca.Accept()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, ic.Token().ToString(), "test chan")

	close(ca)
	ic, err = ca.Accept()
	assert.Equal(t, err, io.EOF, err)
	assert.Equal(t, ic, nil)
}
