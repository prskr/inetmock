// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: rpc/v1/endpoint.proto

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

// EndpointOrchestratorServiceClient is the client API for EndpointOrchestratorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EndpointOrchestratorServiceClient interface {
	ListAllServingGroups(ctx context.Context, in *ListAllServingGroupsRequest, opts ...grpc.CallOption) (*ListAllServingGroupsResponse, error)
	ListAllConfiguredGroups(ctx context.Context, in *ListAllConfiguredGroupsRequest, opts ...grpc.CallOption) (*ListAllConfiguredGroupsResponse, error)
	StartListenerGroup(ctx context.Context, in *StartListenerGroupRequest, opts ...grpc.CallOption) (*StartListenerGroupResponse, error)
	StartAllGroups(ctx context.Context, in *StartAllGroupsRequest, opts ...grpc.CallOption) (*StartAllGroupsResponse, error)
	StopListenerGroup(ctx context.Context, in *StopListenerGroupRequest, opts ...grpc.CallOption) (*StopListenerGroupResponse, error)
	StopAllGroups(ctx context.Context, in *StopAllGroupsRequest, opts ...grpc.CallOption) (*StopAllGroupsResponse, error)
	RestartListenerGroup(ctx context.Context, in *RestartListenerGroupRequest, opts ...grpc.CallOption) (*RestartListenerGroupResponse, error)
	RestartAllGroups(ctx context.Context, in *RestartAllGroupsRequest, opts ...grpc.CallOption) (*RestartAllGroupsResponse, error)
}

type endpointOrchestratorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEndpointOrchestratorServiceClient(cc grpc.ClientConnInterface) EndpointOrchestratorServiceClient {
	return &endpointOrchestratorServiceClient{cc}
}

func (c *endpointOrchestratorServiceClient) ListAllServingGroups(ctx context.Context, in *ListAllServingGroupsRequest, opts ...grpc.CallOption) (*ListAllServingGroupsResponse, error) {
	out := new(ListAllServingGroupsResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.EndpointOrchestratorService/ListAllServingGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointOrchestratorServiceClient) ListAllConfiguredGroups(ctx context.Context, in *ListAllConfiguredGroupsRequest, opts ...grpc.CallOption) (*ListAllConfiguredGroupsResponse, error) {
	out := new(ListAllConfiguredGroupsResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.EndpointOrchestratorService/ListAllConfiguredGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointOrchestratorServiceClient) StartListenerGroup(ctx context.Context, in *StartListenerGroupRequest, opts ...grpc.CallOption) (*StartListenerGroupResponse, error) {
	out := new(StartListenerGroupResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.EndpointOrchestratorService/StartListenerGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointOrchestratorServiceClient) StartAllGroups(ctx context.Context, in *StartAllGroupsRequest, opts ...grpc.CallOption) (*StartAllGroupsResponse, error) {
	out := new(StartAllGroupsResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.EndpointOrchestratorService/StartAllGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointOrchestratorServiceClient) StopListenerGroup(ctx context.Context, in *StopListenerGroupRequest, opts ...grpc.CallOption) (*StopListenerGroupResponse, error) {
	out := new(StopListenerGroupResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.EndpointOrchestratorService/StopListenerGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointOrchestratorServiceClient) StopAllGroups(ctx context.Context, in *StopAllGroupsRequest, opts ...grpc.CallOption) (*StopAllGroupsResponse, error) {
	out := new(StopAllGroupsResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.EndpointOrchestratorService/StopAllGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointOrchestratorServiceClient) RestartListenerGroup(ctx context.Context, in *RestartListenerGroupRequest, opts ...grpc.CallOption) (*RestartListenerGroupResponse, error) {
	out := new(RestartListenerGroupResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.EndpointOrchestratorService/RestartListenerGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointOrchestratorServiceClient) RestartAllGroups(ctx context.Context, in *RestartAllGroupsRequest, opts ...grpc.CallOption) (*RestartAllGroupsResponse, error) {
	out := new(RestartAllGroupsResponse)
	err := c.cc.Invoke(ctx, "/inetmock.rpc.v1.EndpointOrchestratorService/RestartAllGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointOrchestratorServiceServer is the server API for EndpointOrchestratorService service.
// All implementations must embed UnimplementedEndpointOrchestratorServiceServer
// for forward compatibility
type EndpointOrchestratorServiceServer interface {
	ListAllServingGroups(context.Context, *ListAllServingGroupsRequest) (*ListAllServingGroupsResponse, error)
	ListAllConfiguredGroups(context.Context, *ListAllConfiguredGroupsRequest) (*ListAllConfiguredGroupsResponse, error)
	StartListenerGroup(context.Context, *StartListenerGroupRequest) (*StartListenerGroupResponse, error)
	StartAllGroups(context.Context, *StartAllGroupsRequest) (*StartAllGroupsResponse, error)
	StopListenerGroup(context.Context, *StopListenerGroupRequest) (*StopListenerGroupResponse, error)
	StopAllGroups(context.Context, *StopAllGroupsRequest) (*StopAllGroupsResponse, error)
	RestartListenerGroup(context.Context, *RestartListenerGroupRequest) (*RestartListenerGroupResponse, error)
	RestartAllGroups(context.Context, *RestartAllGroupsRequest) (*RestartAllGroupsResponse, error)
	mustEmbedUnimplementedEndpointOrchestratorServiceServer()
}

// UnimplementedEndpointOrchestratorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEndpointOrchestratorServiceServer struct {
}

func (UnimplementedEndpointOrchestratorServiceServer) ListAllServingGroups(context.Context, *ListAllServingGroupsRequest) (*ListAllServingGroupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllServingGroups not implemented")
}
func (UnimplementedEndpointOrchestratorServiceServer) ListAllConfiguredGroups(context.Context, *ListAllConfiguredGroupsRequest) (*ListAllConfiguredGroupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllConfiguredGroups not implemented")
}
func (UnimplementedEndpointOrchestratorServiceServer) StartListenerGroup(context.Context, *StartListenerGroupRequest) (*StartListenerGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartListenerGroup not implemented")
}
func (UnimplementedEndpointOrchestratorServiceServer) StartAllGroups(context.Context, *StartAllGroupsRequest) (*StartAllGroupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartAllGroups not implemented")
}
func (UnimplementedEndpointOrchestratorServiceServer) StopListenerGroup(context.Context, *StopListenerGroupRequest) (*StopListenerGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopListenerGroup not implemented")
}
func (UnimplementedEndpointOrchestratorServiceServer) StopAllGroups(context.Context, *StopAllGroupsRequest) (*StopAllGroupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopAllGroups not implemented")
}
func (UnimplementedEndpointOrchestratorServiceServer) RestartListenerGroup(context.Context, *RestartListenerGroupRequest) (*RestartListenerGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartListenerGroup not implemented")
}
func (UnimplementedEndpointOrchestratorServiceServer) RestartAllGroups(context.Context, *RestartAllGroupsRequest) (*RestartAllGroupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartAllGroups not implemented")
}
func (UnimplementedEndpointOrchestratorServiceServer) mustEmbedUnimplementedEndpointOrchestratorServiceServer() {
}

// UnsafeEndpointOrchestratorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EndpointOrchestratorServiceServer will
// result in compilation errors.
type UnsafeEndpointOrchestratorServiceServer interface {
	mustEmbedUnimplementedEndpointOrchestratorServiceServer()
}

func RegisterEndpointOrchestratorServiceServer(s grpc.ServiceRegistrar, srv EndpointOrchestratorServiceServer) {
	s.RegisterService(&EndpointOrchestratorService_ServiceDesc, srv)
}

func _EndpointOrchestratorService_ListAllServingGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAllServingGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointOrchestratorServiceServer).ListAllServingGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.EndpointOrchestratorService/ListAllServingGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointOrchestratorServiceServer).ListAllServingGroups(ctx, req.(*ListAllServingGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointOrchestratorService_ListAllConfiguredGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAllConfiguredGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointOrchestratorServiceServer).ListAllConfiguredGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.EndpointOrchestratorService/ListAllConfiguredGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointOrchestratorServiceServer).ListAllConfiguredGroups(ctx, req.(*ListAllConfiguredGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointOrchestratorService_StartListenerGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartListenerGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointOrchestratorServiceServer).StartListenerGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.EndpointOrchestratorService/StartListenerGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointOrchestratorServiceServer).StartListenerGroup(ctx, req.(*StartListenerGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointOrchestratorService_StartAllGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartAllGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointOrchestratorServiceServer).StartAllGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.EndpointOrchestratorService/StartAllGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointOrchestratorServiceServer).StartAllGroups(ctx, req.(*StartAllGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointOrchestratorService_StopListenerGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopListenerGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointOrchestratorServiceServer).StopListenerGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.EndpointOrchestratorService/StopListenerGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointOrchestratorServiceServer).StopListenerGroup(ctx, req.(*StopListenerGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointOrchestratorService_StopAllGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopAllGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointOrchestratorServiceServer).StopAllGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.EndpointOrchestratorService/StopAllGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointOrchestratorServiceServer).StopAllGroups(ctx, req.(*StopAllGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointOrchestratorService_RestartListenerGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestartListenerGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointOrchestratorServiceServer).RestartListenerGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.EndpointOrchestratorService/RestartListenerGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointOrchestratorServiceServer).RestartListenerGroup(ctx, req.(*RestartListenerGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointOrchestratorService_RestartAllGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestartAllGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointOrchestratorServiceServer).RestartAllGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inetmock.rpc.v1.EndpointOrchestratorService/RestartAllGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointOrchestratorServiceServer).RestartAllGroups(ctx, req.(*RestartAllGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EndpointOrchestratorService_ServiceDesc is the grpc.ServiceDesc for EndpointOrchestratorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EndpointOrchestratorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "inetmock.rpc.v1.EndpointOrchestratorService",
	HandlerType: (*EndpointOrchestratorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAllServingGroups",
			Handler:    _EndpointOrchestratorService_ListAllServingGroups_Handler,
		},
		{
			MethodName: "ListAllConfiguredGroups",
			Handler:    _EndpointOrchestratorService_ListAllConfiguredGroups_Handler,
		},
		{
			MethodName: "StartListenerGroup",
			Handler:    _EndpointOrchestratorService_StartListenerGroup_Handler,
		},
		{
			MethodName: "StartAllGroups",
			Handler:    _EndpointOrchestratorService_StartAllGroups_Handler,
		},
		{
			MethodName: "StopListenerGroup",
			Handler:    _EndpointOrchestratorService_StopListenerGroup_Handler,
		},
		{
			MethodName: "StopAllGroups",
			Handler:    _EndpointOrchestratorService_StopAllGroups_Handler,
		},
		{
			MethodName: "RestartListenerGroup",
			Handler:    _EndpointOrchestratorService_RestartListenerGroup_Handler,
		},
		{
			MethodName: "RestartAllGroups",
			Handler:    _EndpointOrchestratorService_RestartAllGroups_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/v1/endpoint.proto",
}
