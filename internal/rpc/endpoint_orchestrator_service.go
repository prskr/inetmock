package rpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	rpcv1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

var _ rpcv1.EndpointOrchestratorServiceServer = (*endpointOrchestratorServer)(nil)

func NewEndpointOrchestratorServer(
	logger logging.Logger,
	epHost endpoint.Host,
) rpcv1.EndpointOrchestratorServiceServer {
	return &endpointOrchestratorServer{
		UnimplementedEndpointOrchestratorServiceServer: rpcv1.UnimplementedEndpointOrchestratorServiceServer{},
		logger: logger,
		epHost: epHost,
	}
}

type endpointOrchestratorServer struct {
	rpcv1.UnimplementedEndpointOrchestratorServiceServer
	logger logging.Logger
	epHost endpoint.Host
}

func (s *endpointOrchestratorServer) ListAllServingGroups(
	context.Context,
	*rpcv1.ListAllServingGroupsRequest) (*rpcv1.ListAllServingGroupsResponse, error,
) {
	groups := s.epHost.ConfiguredGroups()
	resp := &rpcv1.ListAllServingGroupsResponse{
		Groups: make([]*rpcv1.ListenerGroup, 0, len(groups)),
	}

	for idx := range groups {
		if !groups[idx].Serving {
			continue
		}

		resp.Groups = append(resp.Groups, &rpcv1.ListenerGroup{
			Name:      groups[idx].Name,
			Endpoints: groups[idx].Endpoints,
		})
	}

	return resp, nil
}

func (s *endpointOrchestratorServer) ListAllConfiguredGroups(
	_ context.Context,
	_ *rpcv1.ListAllConfiguredGroupsRequest,
) (*rpcv1.ListAllConfiguredGroupsResponse, error) {
	groups := s.epHost.ConfiguredGroups()
	resp := &rpcv1.ListAllConfiguredGroupsResponse{
		Groups: make([]*rpcv1.ListenerGroup, 0, len(groups)),
	}

	for idx := range groups {
		resp.Groups = append(resp.Groups, &rpcv1.ListenerGroup{
			Name:      groups[idx].Name,
			Endpoints: groups[idx].Endpoints,
		})
	}

	return resp, nil
}

func (s *endpointOrchestratorServer) StartListenerGroup(
	ctx context.Context,
	req *rpcv1.StartListenerGroupRequest,
) (*rpcv1.StartListenerGroupResponse, error) {
	if err := s.epHost.ServeGroup(context.Background(), req.GroupName); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return new(rpcv1.StartListenerGroupResponse), nil
}

func (s *endpointOrchestratorServer) StartAllGroups(
	ctx context.Context,
	_ *rpcv1.StartAllGroupsRequest) (*rpcv1.StartAllGroupsResponse, error) {
	if err := s.epHost.ServeGroups(context.Background()); err != nil {
		s.logger.Error("Failed to start serving groups", zap.Error(err))
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return new(rpcv1.StartAllGroupsResponse), nil
}

func (s *endpointOrchestratorServer) StopListenerGroup(
	ctx context.Context,
	req *rpcv1.StopListenerGroupRequest,
) (*rpcv1.StopListenerGroupResponse, error) {
	if err := s.epHost.ShutdownGroup(ctx, req.GroupName); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return new(rpcv1.StopListenerGroupResponse), nil
}

func (s *endpointOrchestratorServer) StopAllGroups(
	ctx context.Context,
	_ *rpcv1.StopAllGroupsRequest) (*rpcv1.StopAllGroupsResponse, error) {
	if err := s.epHost.Shutdown(ctx); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return new(rpcv1.StopAllGroupsResponse), nil
}

func (s *endpointOrchestratorServer) RestartListenerGroup(
	ctx context.Context,
	req *rpcv1.RestartListenerGroupRequest,
) (*rpcv1.RestartListenerGroupResponse, error) {
	if err := s.epHost.ShutdownGroup(ctx, req.GroupName); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if err := s.epHost.ServeGroup(context.Background(), req.GroupName); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return new(rpcv1.RestartListenerGroupResponse), nil
}

func (s *endpointOrchestratorServer) RestartAllGroups(
	ctx context.Context,
	_ *rpcv1.RestartAllGroupsRequest,
) (*rpcv1.RestartAllGroupsResponse, error) {
	if err := s.epHost.Shutdown(ctx); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	if err := s.epHost.ServeGroups(context.Background()); err != nil {
		s.logger.Error("Failed to start serving groups", zap.Error(err))
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return new(rpcv1.RestartAllGroupsResponse), nil
}
