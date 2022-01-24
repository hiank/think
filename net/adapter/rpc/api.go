package rpc

import (
	"context"

	"github.com/hiank/think/net/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type REST interface {
	Get(context.Context, *pb.Carrier) (*pb.Carrier, error)
	Post(context.Context, *pb.Carrier) (*emptypb.Empty, error)
}
