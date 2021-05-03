// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// AuditServiceClient is the client API for AuditService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuditServiceClient interface {
	WatchEvents(ctx context.Context, in *WatchEventsRequest, opts ...grpc.CallOption) (AuditService_WatchEventsClient, error)
	RegisterFileSink(ctx context.Context, in *RegisterFileSinkRequest, opts ...grpc.CallOption) (*RegisterFileSinkResponse, error)
	RemoveFileSink(ctx context.Context, in *RemoveFileSinkRequest, opts ...grpc.CallOption) (*RemoveFileSinkResponse, error)
	ListSinks(ctx context.Context, in *ListSinksRequest, opts ...grpc.CallOption) (*ListSinksResponse, error)
}

type auditServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuditServiceClient(cc grpc.ClientConnInterface) AuditServiceClient {
	return &auditServiceClient{cc}
}

func (c *auditServiceClient) WatchEvents(ctx context.Context, in *WatchEventsRequest, opts ...grpc.CallOption) (AuditService_WatchEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &AuditService_ServiceDesc.Streams[0], "/inetmock.rpc.v1.AuditService/WatchEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &auditServiceWatchEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AuditService_WatchEventsClient interface {
	Recv() (*WatchEventsResponse, error)
	grpc.ClientStream
}

type auditServiceWatchEventsClient struct {
	grpc.ClientStream
}

func (x *auditServiceWatchEventsClient) Recv() (*WatchEventsResponse, error) {
	m := new(WatchEventsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *auditServiceClient) RegisterFileSink(ctx context.Context, in *RegisterFileSinkRequest, opts ...grpc.CallOption) (*RegisterFileSinkResponse, error) {
	out := new(RegisterFileSinkResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.AuditService/RegisterFileSink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auditServiceClient) RemoveFileSink(ctx context.Context, in *RemoveFileSinkRequest, opts ...grpc.CallOption) (*RemoveFileSinkResponse, error) {
	out := new(RemoveFileSinkResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.AuditService/RemoveFileSink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auditServiceClient) ListSinks(ctx context.Context, in *ListSinksRequest, opts ...grpc.CallOption) (*ListSinksResponse, error) {
	out := new(ListSinksResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.AuditService/ListSinks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuditServiceServer is the server API for AuditService service.
// All implementations must embed UnimplementedAuditServiceServer
// for forward compatibility
type AuditServiceServer interface {
	WatchEvents(*WatchEventsRequest, AuditService_WatchEventsServer) error
	RegisterFileSink(context.Context, *RegisterFileSinkRequest) (*RegisterFileSinkResponse, error)
	RemoveFileSink(context.Context, *RemoveFileSinkRequest) (*RemoveFileSinkResponse, error)
	ListSinks(context.Context, *ListSinksRequest) (*ListSinksResponse, error)
	mustEmbedUnimplementedAuditServiceServer()
}

// UnimplementedAuditServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuditServiceServer struct {
}

func (UnimplementedAuditServiceServer) WatchEvents(*WatchEventsRequest, AuditService_WatchEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchEvents not implemented")
}
func (UnimplementedAuditServiceServer) RegisterFileSink(context.Context, *RegisterFileSinkRequest) (*RegisterFileSinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterFileSink not implemented")
}
func (UnimplementedAuditServiceServer) RemoveFileSink(context.Context, *RemoveFileSinkRequest) (*RemoveFileSinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFileSink not implemented")
}
func (UnimplementedAuditServiceServer) ListSinks(context.Context, *ListSinksRequest) (*ListSinksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSinks not implemented")
}
func (UnimplementedAuditServiceServer) mustEmbedUnimplementedAuditServiceServer() {}

// UnsafeAuditServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuditServiceServer will
// result in compilation errors.
type UnsafeAuditServiceServer interface {
	mustEmbedUnimplementedAuditServiceServer()
}

func RegisterAuditServiceServer(s grpc.ServiceRegistrar, srv AuditServiceServer) {
	s.RegisterService(&AuditService_ServiceDesc, srv)
}

func _AuditService_WatchEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchEventsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AuditServiceServer).WatchEvents(m, &auditServiceWatchEventsServer{stream})
}

type AuditService_WatchEventsServer interface {
	Send(*WatchEventsResponse) error
	grpc.ServerStream
}

type auditServiceWatchEventsServer struct {
	grpc.ServerStream
}

func (x *auditServiceWatchEventsServer) Send(m *WatchEventsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _AuditService_RegisterFileSink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterFileSinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuditServiceServer).RegisterFileSink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.AuditService/RegisterFileSink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuditServiceServer).RegisterFileSink(ctx, req.(*RegisterFileSinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuditService_RemoveFileSink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFileSinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuditServiceServer).RemoveFileSink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.AuditService/RemoveFileSink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuditServiceServer).RemoveFileSink(ctx, req.(*RemoveFileSinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuditService_ListSinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuditServiceServer).ListSinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.AuditService/ListSinks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuditServiceServer).ListSinks(ctx, req.(*ListSinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuditService_ServiceDesc is the grpc.ServiceDesc for AuditService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuditService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "inetmock.rpc.v1.AuditService",
	HandlerType: (*AuditServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterFileSink",
			Handler:    _AuditService_RegisterFileSink_Handler,
		},
		{
			MethodName: "RemoveFileSink",
			Handler:    _AuditService_RemoveFileSink_Handler,
		},
		{
			MethodName: "ListSinks",
			Handler:    _AuditService_ListSinks_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchEvents",
			Handler:       _AuditService_WatchEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc/v1/audit.proto",
}