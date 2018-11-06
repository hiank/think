package net

import (
	"fmt"
	"net"
	"io"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	tg "github.com/hiank/think/net/protobuf/grpc"
)

// K8sServer represents a server of a k8s.
type K8sServer struct {
	// Config is a WebSocket configuration for new WebSocket connection.
	// Config

	// Handshake is an optional function in WebSocket handshake.
	// For example, you can check, or don't check Origin header.
	// Another example, you can select config.Protocol.
	// Handshake func(*Config, *http.Request) error

	// Handler handles a tg.Request
	Handler
}

// Tran is realize of PipeServer
func (s *K8sServer) Tran(pipe tg.Pipe_TranServer) error {

	return s.serveK8s(pipe)
}

// serveK8s 对每个grpc_client 执行的Tran 操作，启用一个此服务，用于收发数据
func (s *K8sServer) serveK8s(pipe tg.Pipe_TranServer) (err error) {


	ch := make(chan *tg.Response)
	quit := make(chan bool)
	go func () {			

SL:		for {		//NOTE: 发送逻辑处理过的数据到Grpc调用者(grpc_client)

			select {
			case <-quit: break SL
			case tank := <-ch: pipe.Send(tank)
			}
		}
	}()

L:	for {			//NOTE: 从Grpc调用者(grpc_client)收到的数据，使用Handler进行处理。Handler异步执行

		tank, err := pipe.Recv()
		switch {
		case err == io.EOF: break L
		case err != nil: break L
		}
		go s.Handler(tank, ch)
	}
	close(ch)
	close(quit)
	return
}


// Handler is a simple interface for GRPC serve.
type Handler func(*tg.Request, chan *tg.Response)


// Tran is realize of PipeServer
func (h Handler) Tran(pipe tg.Pipe_TranServer) error {

	s := &K8sServer{Handler: h}
	return s.serveK8s(pipe)
}

// ListenAndServeK8s start a PipeServer
func ListenAndServeK8s(server tg.PipeServer, port int) (err error) {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	tg.RegisterPipeServer(grpcServer, server)
	return grpcServer.Serve(lis)
}