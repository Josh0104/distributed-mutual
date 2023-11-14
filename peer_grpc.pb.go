// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: peer.proto

package main

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
	MutualService_Join_FullMethodName             = "/main.MutualService/Join"
	MutualService_StreamFromClient_FullMethodName = "/main.MutualService/StreamFromClient"
	MutualService_Broadcast_FullMethodName        = "/main.MutualService/Broadcast"
	MutualService_GetToken_FullMethodName         = "/main.MutualService/GetToken"
)

// MutualServiceClient is the client API for MutualService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MutualServiceClient interface {
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error)
	StreamFromClient(ctx context.Context, opts ...grpc.CallOption) (MutualService_StreamFromClientClient, error)
	Broadcast(ctx context.Context, in *BroadcastSubscription, opts ...grpc.CallOption) (MutualService_BroadcastClient, error)
	GetToken(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Token, error)
}

type mutualServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMutualServiceClient(cc grpc.ClientConnInterface) MutualServiceClient {
	return &mutualServiceClient{cc}
}

func (c *mutualServiceClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error) {
	out := new(JoinResponse)
	err := c.cc.Invoke(ctx, MutualService_Join_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mutualServiceClient) StreamFromClient(ctx context.Context, opts ...grpc.CallOption) (MutualService_StreamFromClientClient, error) {
	stream, err := c.cc.NewStream(ctx, &MutualService_ServiceDesc.Streams[0], MutualService_StreamFromClient_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &mutualServiceStreamFromClientClient{stream}
	return x, nil
}

type MutualService_StreamFromClientClient interface {
	Send(*Message) error
	CloseAndRecv() (*Message, error)
	grpc.ClientStream
}

type mutualServiceStreamFromClientClient struct {
	grpc.ClientStream
}

func (x *mutualServiceStreamFromClientClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mutualServiceStreamFromClientClient) CloseAndRecv() (*Message, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mutualServiceClient) Broadcast(ctx context.Context, in *BroadcastSubscription, opts ...grpc.CallOption) (MutualService_BroadcastClient, error) {
	stream, err := c.cc.NewStream(ctx, &MutualService_ServiceDesc.Streams[1], MutualService_Broadcast_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &mutualServiceBroadcastClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MutualService_BroadcastClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type mutualServiceBroadcastClient struct {
	grpc.ClientStream
}

func (x *mutualServiceBroadcastClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mutualServiceClient) GetToken(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, MutualService_GetToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MutualServiceServer is the server API for MutualService service.
// All implementations must embed UnimplementedMutualServiceServer
// for forward compatibility
type MutualServiceServer interface {
	Join(context.Context, *JoinRequest) (*JoinResponse, error)
	StreamFromClient(MutualService_StreamFromClientServer) error
	Broadcast(*BroadcastSubscription, MutualService_BroadcastServer) error
	GetToken(context.Context, *Token) (*Token, error)
	mustEmbedUnimplementedMutualServiceServer()
}

// UnimplementedMutualServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMutualServiceServer struct {
}

func (UnimplementedMutualServiceServer) Join(context.Context, *JoinRequest) (*JoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedMutualServiceServer) StreamFromClient(MutualService_StreamFromClientServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamFromClient not implemented")
}
func (UnimplementedMutualServiceServer) Broadcast(*BroadcastSubscription, MutualService_BroadcastServer) error {
	return status.Errorf(codes.Unimplemented, "method Broadcast not implemented")
}
func (UnimplementedMutualServiceServer) GetToken(context.Context, *Token) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetToken not implemented")
}
func (UnimplementedMutualServiceServer) mustEmbedUnimplementedMutualServiceServer() {}

// UnsafeMutualServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MutualServiceServer will
// result in compilation errors.
type UnsafeMutualServiceServer interface {
	mustEmbedUnimplementedMutualServiceServer()
}

func RegisterMutualServiceServer(s grpc.ServiceRegistrar, srv MutualServiceServer) {
	s.RegisterService(&MutualService_ServiceDesc, srv)
}

func _MutualService_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutualServiceServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MutualService_Join_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutualServiceServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MutualService_StreamFromClient_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MutualServiceServer).StreamFromClient(&mutualServiceStreamFromClientServer{stream})
}

type MutualService_StreamFromClientServer interface {
	SendAndClose(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type mutualServiceStreamFromClientServer struct {
	grpc.ServerStream
}

func (x *mutualServiceStreamFromClientServer) SendAndClose(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mutualServiceStreamFromClientServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _MutualService_Broadcast_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BroadcastSubscription)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MutualServiceServer).Broadcast(m, &mutualServiceBroadcastServer{stream})
}

type MutualService_BroadcastServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type mutualServiceBroadcastServer struct {
	grpc.ServerStream
}

func (x *mutualServiceBroadcastServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func _MutualService_GetToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutualServiceServer).GetToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MutualService_GetToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutualServiceServer).GetToken(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

// MutualService_ServiceDesc is the grpc.ServiceDesc for MutualService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MutualService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.MutualService",
	HandlerType: (*MutualServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _MutualService_Join_Handler,
		},
		{
			MethodName: "GetToken",
			Handler:    _MutualService_GetToken_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamFromClient",
			Handler:       _MutualService_StreamFromClient_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Broadcast",
			Handler:       _MutualService_Broadcast_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "peer.proto",
}
