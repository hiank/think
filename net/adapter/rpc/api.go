package rpc

import (
	"context"

	"github.com/hiank/think/net/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type IHandler interface {
	Get(*pb.Carrier) *pb.Carrier
	Post(*pb.Carrier)
}

type IREST interface {
	Get(context.Context, *pb.Carrier) (*pb.Carrier, error)
	Post(context.Context, *pb.Carrier) (*emptypb.Empty, error)
}

// type I
