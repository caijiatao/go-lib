// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: push_svr.proto

package push

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Pusher_PushMessage_FullMethodName = "/push.Pusher/PushMessage"
)

// PusherClient is the client API for Pusher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PusherClient interface {
	PushMessage(ctx context.Context, in *PushMessageRequest, opts ...grpc.CallOption) (*PushMessageReply, error)
}

type pusherClient struct {
	cc grpc.ClientConnInterface
}

func NewPusherClient(cc grpc.ClientConnInterface) PusherClient {
	return &pusherClient{cc}
}

func (c *pusherClient) PushMessage(ctx context.Context, in *PushMessageRequest, opts ...grpc.CallOption) (*PushMessageReply, error) {
	out := new(PushMessageReply)
	err := c.cc.Invoke(ctx, Pusher_PushMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PusherServer is the server API for Pusher service.
// All implementations must embed UnimplementedPusherServer
// for forward compatibility
type PusherServer interface {
	PushMessage(context.Context, *PushMessageRequest) (*PushMessageReply, error)
	mustEmbedUnimplementedPusherServer()
}

// UnimplementedPusherServer must be embedded to have forward compatible implementations.
type UnimplementedPusherServer struct {
}

func (UnimplementedPusherServer) PushMessage(context.Context, *PushMessageRequest) (*PushMessageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMessage not implemented")
}
func (UnimplementedPusherServer) mustEmbedUnimplementedPusherServer() {}

// UnsafePusherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PusherServer will
// result in compilation errors.
type UnsafePusherServer interface {
	mustEmbedUnimplementedPusherServer()
}

func RegisterPusherServer(s grpc.ServiceRegistrar, srv PusherServer) {
	s.RegisterService(&Pusher_ServiceDesc, srv)
}

func _Pusher_PushMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PusherServer).PushMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pusher_PushMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PusherServer).PushMessage(ctx, req.(*PushMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Pusher_ServiceDesc is the grpc.ServiceDesc for Pusher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pusher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "push.Pusher",
	HandlerType: (*PusherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushMessage",
			Handler:    _Pusher_PushMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "push_svr.proto",
}
