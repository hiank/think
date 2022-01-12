// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pp

import (
	context "context"
	pb "github.com/hiank/think/net/pb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PipeClient is the client API for Pipe service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PipeClient interface {
	Link(ctx context.Context, opts ...grpc.CallOption) (Pipe_LinkClient, error)
	Get(ctx context.Context, in *pb.Carrier, opts ...grpc.CallOption) (*pb.Carrier, error)
	Post(ctx context.Context, in *pb.Carrier, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type pipeClient struct {
	cc grpc.ClientConnInterface
}

func NewPipeClient(cc grpc.ClientConnInterface) PipeClient {
	return &pipeClient{cc}
}

func (c *pipeClient) Link(ctx context.Context, opts ...grpc.CallOption) (Pipe_LinkClient, error) {
	stream, err := c.cc.NewStream(ctx, &Pipe_ServiceDesc.Streams[0], "/Pipe/Link", opts...)
	if err != nil {
		return nil, err
	}
	x := &pipeLinkClient{stream}
	return x, nil
}

type Pipe_LinkClient interface {
	Send(*pb.Carrier) error
	Recv() (*pb.Carrier, error)
	grpc.ClientStream
}

type pipeLinkClient struct {
	grpc.ClientStream
}

func (x *pipeLinkClient) Send(m *pb.Carrier) error {
	return x.ClientStream.SendMsg(m)
}

func (x *pipeLinkClient) Recv() (*pb.Carrier, error) {
	m := new(pb.Carrier)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pipeClient) Get(ctx context.Context, in *pb.Carrier, opts ...grpc.CallOption) (*pb.Carrier, error) {
	out := new(pb.Carrier)
	err := c.cc.Invoke(ctx, "/Pipe/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pipeClient) Post(ctx context.Context, in *pb.Carrier, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Pipe/Post", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PipeServer is the server API for Pipe service.
// All implementations must embed UnimplementedPipeServer
// for forward compatibility
type PipeServer interface {
	Link(Pipe_LinkServer) error
	Get(context.Context, *pb.Carrier) (*pb.Carrier, error)
	Post(context.Context, *pb.Carrier) (*emptypb.Empty, error)
	mustEmbedUnimplementedPipeServer()
}

// UnimplementedPipeServer must be embedded to have forward compatible implementations.
type UnimplementedPipeServer struct {
}

func (UnimplementedPipeServer) Link(Pipe_LinkServer) error {
	return status.Errorf(codes.Unimplemented, "method Link not implemented")
}
func (UnimplementedPipeServer) Get(context.Context, *pb.Carrier) (*pb.Carrier, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedPipeServer) Post(context.Context, *pb.Carrier) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Post not implemented")
}
func (UnimplementedPipeServer) mustEmbedUnimplementedPipeServer() {}

// UnsafePipeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PipeServer will
// result in compilation errors.
type UnsafePipeServer interface {
	mustEmbedUnimplementedPipeServer()
}

func RegisterPipeServer(s grpc.ServiceRegistrar, srv PipeServer) {
	s.RegisterService(&Pipe_ServiceDesc, srv)
}

func _Pipe_Link_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PipeServer).Link(&pipeLinkServer{stream})
}

type Pipe_LinkServer interface {
	Send(*pb.Carrier) error
	Recv() (*pb.Carrier, error)
	grpc.ServerStream
}

type pipeLinkServer struct {
	grpc.ServerStream
}

func (x *pipeLinkServer) Send(m *pb.Carrier) error {
	return x.ServerStream.SendMsg(m)
}

func (x *pipeLinkServer) Recv() (*pb.Carrier, error) {
	m := new(pb.Carrier)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Pipe_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(pb.Carrier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PipeServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Pipe/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PipeServer).Get(ctx, req.(*pb.Carrier))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pipe_Post_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(pb.Carrier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PipeServer).Post(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Pipe/Post",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PipeServer).Post(ctx, req.(*pb.Carrier))
	}
	return interceptor(ctx, in, info, handler)
}

// Pipe_ServiceDesc is the grpc.ServiceDesc for Pipe service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pipe_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Pipe",
	HandlerType: (*PipeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Pipe_Get_Handler,
		},
		{
			MethodName: "Post",
			Handler:    _Pipe_Post_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Link",
			Handler:       _Pipe_Link_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "net/adapter/rpc/pp/pipe.proto",
}