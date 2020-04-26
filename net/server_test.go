package net_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/k8s"
	"gotest.tools/assert"
)

func TestCloseWS(t *testing.T) {

	exit := make(chan error)

	go func() {
		exit <- net.ServeWS("127.0.0.1")
	}()

	go func() {
		<- time.After(time.Second)
		net.GetRuntime().Close()	
	}()
	err := <- exit
	assert.Equal(t, err, http.ErrServerClosed)
}


type testK8sHandler struct {
	k8s.IgnoreStream
	k8s.IgnoreGet
	k8s.IgnorePost
}

func TestCloseK8s(t *testing.T) {

	exit := make(chan error)

	go func() {
		exit <- net.ServeK8s("127.0.0.1", &testK8sHandler{})
	}()

	go func() {
		<- time.After(time.Second)
		net.GetRuntime().Close()
	}()
	err := <- exit
	assert.Equal(t, err, nil)
}
