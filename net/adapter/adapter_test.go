package adapter_test

import (
	"io"
	"testing"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter"
	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/one"
	"gotest.tools/v3/assert"
)

func TestChanAccepter(t *testing.T) {
	ca := make(adapter.ChanAccepter)
	go func() {
		ca <- net.TokenConn{Token: one.TokenSet().Derive("test chan")}
	}()

	ic, err := ca.Accept()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, ic.Token.Value(box.ContextkeyTokenUid).(string), "test chan")

	close(ca)
	ic, err = ca.Accept()
	assert.Equal(t, err, io.EOF, err)
	assert.DeepEqual(t, ic, net.TokenConn{})
}
