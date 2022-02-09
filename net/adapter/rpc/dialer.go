package rpc

import (
	"github.com/hiank/think/net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type dialer struct {
}

func NewDialer() net.Dialer {
	d := &dialer{}
	return d
}

func (d *dialer) Dial(addr string) (out net.IAC, err error) {
	_, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return
	}
	// cc.Connect()
	return net.IAC{}, nil
}
