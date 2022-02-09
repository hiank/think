package rpc

import (
	"context"

	"google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Stream interface {
	Send(*anypb.Any) error
	Recv() (*anypb.Any, error)
}

type REST interface {
	Get(context.Context, *anypb.Any) (*anypb.Any, error)
	Post(context.Context, *anypb.Any) (*emptypb.Empty, error)
}
