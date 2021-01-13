package rpc

import (
	"context"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

type endpointsServer struct {
	UnimplementedEndpointsServer
	endpointsManager endpoint.EndpointManager
}

func (e endpointsServer) GetEndpoints(_ context.Context, _ *GetEndpointsRequest) (*GetEndpointsResponse, error) {
	eps := rpcEndpointsFromEndpoints(e.endpointsManager.StartedEndpoints())
	return &GetEndpointsResponse{
		Endpoints: *eps,
	}, nil
}

func rpcEndpointsFromEndpoints(eps []endpoint.Endpoint) *[]*Endpoint {
	out := make([]*Endpoint, 0)
	for _, ep := range eps {
		out = append(out, rpcEndpointFromEndpoint(ep))
	}
	return &out
}

func rpcEndpointFromEndpoint(ep endpoint.Endpoint) *Endpoint {
	return &Endpoint{
		Id:            ep.Id().String(),
		Name:          ep.Name(),
		Handler:       ep.Handler(),
		ListenAddress: ep.Listen(),
		Port:          int32(ep.Port()),
	}
}
