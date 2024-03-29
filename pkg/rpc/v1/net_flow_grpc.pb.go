// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: rpc/v1/net_flow.proto

package rpcv1

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

// NetFlowControlServiceClient is the client API for NetFlowControlService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetFlowControlServiceClient interface {
	ListAvailableNetworkInterfaces(ctx context.Context, in *ListAvailableNetworkInterfacesRequest, opts ...grpc.CallOption) (*ListAvailableNetworkInterfacesResponse, error)
	ListControlledInterfaces(ctx context.Context, in *ListControlledInterfacesRequest, opts ...grpc.CallOption) (*ListControlledInterfacesResponse, error)
	StartPacketFlowControl(ctx context.Context, in *StartPacketFlowControlRequest, opts ...grpc.CallOption) (*StartPacketFlowControlResponse, error)
	StopPacketFlowControl(ctx context.Context, in *StopPacketFlowControlRequest, opts ...grpc.CallOption) (*StopPacketFlowControlResponse, error)
}

type netFlowControlServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNetFlowControlServiceClient(cc grpc.ClientConnInterface) NetFlowControlServiceClient {
	return &netFlowControlServiceClient{cc}
}

func (c *netFlowControlServiceClient) ListAvailableNetworkInterfaces(ctx context.Context, in *ListAvailableNetworkInterfacesRequest, opts ...grpc.CallOption) (*ListAvailableNetworkInterfacesResponse, error) {
	out := new(ListAvailableNetworkInterfacesResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.NetFlowControlService/ListAvailableNetworkInterfaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *netFlowControlServiceClient) ListControlledInterfaces(ctx context.Context, in *ListControlledInterfacesRequest, opts ...grpc.CallOption) (*ListControlledInterfacesResponse, error) {
	out := new(ListControlledInterfacesResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.NetFlowControlService/ListControlledInterfaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *netFlowControlServiceClient) StartPacketFlowControl(ctx context.Context, in *StartPacketFlowControlRequest, opts ...grpc.CallOption) (*StartPacketFlowControlResponse, error) {
	out := new(StartPacketFlowControlResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.NetFlowControlService/StartPacketFlowControl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *netFlowControlServiceClient) StopPacketFlowControl(ctx context.Context, in *StopPacketFlowControlRequest, opts ...grpc.CallOption) (*StopPacketFlowControlResponse, error) {
	out := new(StopPacketFlowControlResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.NetFlowControlService/StopPacketFlowControl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetFlowControlServiceServer is the server API for NetFlowControlService service.
// All implementations must embed UnimplementedNetFlowControlServiceServer
// for forward compatibility
type NetFlowControlServiceServer interface {
	ListAvailableNetworkInterfaces(context.Context, *ListAvailableNetworkInterfacesRequest) (*ListAvailableNetworkInterfacesResponse, error)
	ListControlledInterfaces(context.Context, *ListControlledInterfacesRequest) (*ListControlledInterfacesResponse, error)
	StartPacketFlowControl(context.Context, *StartPacketFlowControlRequest) (*StartPacketFlowControlResponse, error)
	StopPacketFlowControl(context.Context, *StopPacketFlowControlRequest) (*StopPacketFlowControlResponse, error)
	mustEmbedUnimplementedNetFlowControlServiceServer()
}

// UnimplementedNetFlowControlServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNetFlowControlServiceServer struct {
}

func (UnimplementedNetFlowControlServiceServer) ListAvailableNetworkInterfaces(context.Context, *ListAvailableNetworkInterfacesRequest) (*ListAvailableNetworkInterfacesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAvailableNetworkInterfaces not implemented")
}
func (UnimplementedNetFlowControlServiceServer) ListControlledInterfaces(context.Context, *ListControlledInterfacesRequest) (*ListControlledInterfacesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListControlledInterfaces not implemented")
}
func (UnimplementedNetFlowControlServiceServer) StartPacketFlowControl(context.Context, *StartPacketFlowControlRequest) (*StartPacketFlowControlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartPacketFlowControl not implemented")
}
func (UnimplementedNetFlowControlServiceServer) StopPacketFlowControl(context.Context, *StopPacketFlowControlRequest) (*StopPacketFlowControlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopPacketFlowControl not implemented")
}
func (UnimplementedNetFlowControlServiceServer) mustEmbedUnimplementedNetFlowControlServiceServer() {}

// UnsafeNetFlowControlServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetFlowControlServiceServer will
// result in compilation errors.
type UnsafeNetFlowControlServiceServer interface {
	mustEmbedUnimplementedNetFlowControlServiceServer()
}

func RegisterNetFlowControlServiceServer(s grpc.ServiceRegistrar, srv NetFlowControlServiceServer) {
	s.RegisterService(&NetFlowControlService_ServiceDesc, srv)
}

func _NetFlowControlService_ListAvailableNetworkInterfaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAvailableNetworkInterfacesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetFlowControlServiceServer).ListAvailableNetworkInterfaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.NetFlowControlService/ListAvailableNetworkInterfaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetFlowControlServiceServer).ListAvailableNetworkInterfaces(ctx, req.(*ListAvailableNetworkInterfacesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetFlowControlService_ListControlledInterfaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListControlledInterfacesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetFlowControlServiceServer).ListControlledInterfaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.NetFlowControlService/ListControlledInterfaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetFlowControlServiceServer).ListControlledInterfaces(ctx, req.(*ListControlledInterfacesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetFlowControlService_StartPacketFlowControl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartPacketFlowControlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetFlowControlServiceServer).StartPacketFlowControl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.NetFlowControlService/StartPacketFlowControl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetFlowControlServiceServer).StartPacketFlowControl(ctx, req.(*StartPacketFlowControlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetFlowControlService_StopPacketFlowControl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopPacketFlowControlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetFlowControlServiceServer).StopPacketFlowControl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.NetFlowControlService/StopPacketFlowControl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetFlowControlServiceServer).StopPacketFlowControl(ctx, req.(*StopPacketFlowControlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NetFlowControlService_ServiceDesc is the grpc.ServiceDesc for NetFlowControlService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetFlowControlService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "inetmock.rpc.v1.NetFlowControlService",
	HandlerType: (*NetFlowControlServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAvailableNetworkInterfaces",
			Handler:    _NetFlowControlService_ListAvailableNetworkInterfaces_Handler,
		},
		{
			MethodName: "ListControlledInterfaces",
			Handler:    _NetFlowControlService_ListControlledInterfaces_Handler,
		},
		{
			MethodName: "StartPacketFlowControl",
			Handler:    _NetFlowControlService_StartPacketFlowControl_Handler,
		},
		{
			MethodName: "StopPacketFlowControl",
			Handler:    _NetFlowControlService_StopPacketFlowControl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/v1/net_flow.proto",
}
