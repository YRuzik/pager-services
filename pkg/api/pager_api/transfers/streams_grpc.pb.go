// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: transfers/streams.proto

package pager_transfers

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
	PagerStreams_StreamProfile_FullMethodName    = "/com.pager.api.PagerStreams/StreamProfile"
	PagerStreams_StreamChat_FullMethodName       = "/com.pager.api.PagerStreams/StreamChat"
	PagerStreams_StreamChatMember_FullMethodName = "/com.pager.api.PagerStreams/StreamChatMember"
)

// PagerStreamsClient is the client API for PagerStreams service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PagerStreamsClient interface {
	StreamProfile(ctx context.Context, in *ProfileStreamRequest, opts ...grpc.CallOption) (PagerStreams_StreamProfileClient, error)
	StreamChat(ctx context.Context, in *ChatStreamRequest, opts ...grpc.CallOption) (PagerStreams_StreamChatClient, error)
	StreamChatMember(ctx context.Context, in *ChatMemberRequest, opts ...grpc.CallOption) (PagerStreams_StreamChatMemberClient, error)
}

type pagerStreamsClient struct {
	cc grpc.ClientConnInterface
}

func NewPagerStreamsClient(cc grpc.ClientConnInterface) PagerStreamsClient {
	return &pagerStreamsClient{cc}
}

func (c *pagerStreamsClient) StreamProfile(ctx context.Context, in *ProfileStreamRequest, opts ...grpc.CallOption) (PagerStreams_StreamProfileClient, error) {
	stream, err := c.cc.NewStream(ctx, &PagerStreams_ServiceDesc.Streams[0], PagerStreams_StreamProfile_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pagerStreamsStreamProfileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PagerStreams_StreamProfileClient interface {
	Recv() (*TransferObject, error)
	grpc.ClientStream
}

type pagerStreamsStreamProfileClient struct {
	grpc.ClientStream
}

func (x *pagerStreamsStreamProfileClient) Recv() (*TransferObject, error) {
	m := new(TransferObject)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pagerStreamsClient) StreamChat(ctx context.Context, in *ChatStreamRequest, opts ...grpc.CallOption) (PagerStreams_StreamChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &PagerStreams_ServiceDesc.Streams[1], PagerStreams_StreamChat_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pagerStreamsStreamChatClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PagerStreams_StreamChatClient interface {
	Recv() (*TransferObject, error)
	grpc.ClientStream
}

type pagerStreamsStreamChatClient struct {
	grpc.ClientStream
}

func (x *pagerStreamsStreamChatClient) Recv() (*TransferObject, error) {
	m := new(TransferObject)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pagerStreamsClient) StreamChatMember(ctx context.Context, in *ChatMemberRequest, opts ...grpc.CallOption) (PagerStreams_StreamChatMemberClient, error) {
	stream, err := c.cc.NewStream(ctx, &PagerStreams_ServiceDesc.Streams[2], PagerStreams_StreamChatMember_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pagerStreamsStreamChatMemberClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PagerStreams_StreamChatMemberClient interface {
	Recv() (*TransferObject, error)
	grpc.ClientStream
}

type pagerStreamsStreamChatMemberClient struct {
	grpc.ClientStream
}

func (x *pagerStreamsStreamChatMemberClient) Recv() (*TransferObject, error) {
	m := new(TransferObject)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PagerStreamsServer is the server API for PagerStreams service.
// All implementations should embed UnimplementedPagerStreamsServer
// for forward compatibility
type PagerStreamsServer interface {
	StreamProfile(*ProfileStreamRequest, PagerStreams_StreamProfileServer) error
	StreamChat(*ChatStreamRequest, PagerStreams_StreamChatServer) error
	StreamChatMember(*ChatMemberRequest, PagerStreams_StreamChatMemberServer) error
}

// UnimplementedPagerStreamsServer should be embedded to have forward compatible implementations.
type UnimplementedPagerStreamsServer struct {
}

func (UnimplementedPagerStreamsServer) StreamProfile(*ProfileStreamRequest, PagerStreams_StreamProfileServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamProfile not implemented")
}
func (UnimplementedPagerStreamsServer) StreamChat(*ChatStreamRequest, PagerStreams_StreamChatServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamChat not implemented")
}
func (UnimplementedPagerStreamsServer) StreamChatMember(*ChatMemberRequest, PagerStreams_StreamChatMemberServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamChatMember not implemented")
}

// UnsafePagerStreamsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PagerStreamsServer will
// result in compilation errors.
type UnsafePagerStreamsServer interface {
	mustEmbedUnimplementedPagerStreamsServer()
}

func RegisterPagerStreamsServer(s grpc.ServiceRegistrar, srv PagerStreamsServer) {
	s.RegisterService(&PagerStreams_ServiceDesc, srv)
}

func _PagerStreams_StreamProfile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ProfileStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PagerStreamsServer).StreamProfile(m, &pagerStreamsStreamProfileServer{stream})
}

type PagerStreams_StreamProfileServer interface {
	Send(*TransferObject) error
	grpc.ServerStream
}

type pagerStreamsStreamProfileServer struct {
	grpc.ServerStream
}

func (x *pagerStreamsStreamProfileServer) Send(m *TransferObject) error {
	return x.ServerStream.SendMsg(m)
}

func _PagerStreams_StreamChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChatStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PagerStreamsServer).StreamChat(m, &pagerStreamsStreamChatServer{stream})
}

type PagerStreams_StreamChatServer interface {
	Send(*TransferObject) error
	grpc.ServerStream
}

type pagerStreamsStreamChatServer struct {
	grpc.ServerStream
}

func (x *pagerStreamsStreamChatServer) Send(m *TransferObject) error {
	return x.ServerStream.SendMsg(m)
}

func _PagerStreams_StreamChatMember_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChatMemberRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PagerStreamsServer).StreamChatMember(m, &pagerStreamsStreamChatMemberServer{stream})
}

type PagerStreams_StreamChatMemberServer interface {
	Send(*TransferObject) error
	grpc.ServerStream
}

type pagerStreamsStreamChatMemberServer struct {
	grpc.ServerStream
}

func (x *pagerStreamsStreamChatMemberServer) Send(m *TransferObject) error {
	return x.ServerStream.SendMsg(m)
}

// PagerStreams_ServiceDesc is the grpc.ServiceDesc for PagerStreams service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PagerStreams_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.pager.api.PagerStreams",
	HandlerType: (*PagerStreamsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamProfile",
			Handler:       _PagerStreams_StreamProfile_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamChat",
			Handler:       _PagerStreams_StreamChat_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamChatMember",
			Handler:       _PagerStreams_StreamChatMember_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "transfers/streams.proto",
}
