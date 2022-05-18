package rpc

import (
	"io"

	"github.com/hiank/think/net/adapter/rpc/pipe"
	"google.golang.org/protobuf/types/known/anypb"
)

type SendReciver interface {
	Send(*anypb.Any) error
	Recv() (*anypb.Any, error)
}

type RestClient interface {
	// Get(context.Context, *anypb.Any) (*anypb.Any, error)
	// Post(context.Context, *anypb.Any) (*emptypb.Empty, error)
	pipe.RestClient
	io.Closer
}
