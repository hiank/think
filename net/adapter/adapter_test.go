package adapter_test

import (
	"io"
	"testing"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter"
	"gotest.tools/v3/assert"
)

func TestChanAccepter(t *testing.T) {
	ca := make(adapter.ChanAccepter)
	go func() {
		ca <- net.IdentityConn{ID: "test chan"}
	}()

	ic, err := ca.Accept()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, ic.ID, "test chan")

	close(ca)
	ic, err = ca.Accept()
	assert.Equal(t, err, io.EOF, err)
	assert.DeepEqual(t, ic, net.IdentityConn{})
}
