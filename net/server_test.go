package net_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/hiank/think/net/rpc"
	"github.com/hiank/think/token"

	"github.com/hiank/think/net"
	"gotest.tools/assert"
)

func TestCloseWS(t *testing.T) {

	exit := make(chan error)

	go func() {
		exit <- net.ServeWS("127.0.0.1")
	}()

	go func() {
		<-time.After(time.Second)
		token.BackgroundLife().Kill()
	}()
	err := <-exit
	assert.Equal(t, err, http.ErrServerClosed)
}

type testK8sHandler struct {
	rpc.IgnoreStream
	rpc.IgnoreGet
	rpc.IgnorePost
}

func TestCloseK8s(t *testing.T) {

	exit := make(chan error)

	go func() {
		exit <- net.ServeRPC("127.0.0.1", &testK8sHandler{})
	}()

	go func() {
		<-time.After(time.Second)
		token.BackgroundLife().Kill()
	}()
	err := <-exit
	assert.Equal(t, err, nil)
}
