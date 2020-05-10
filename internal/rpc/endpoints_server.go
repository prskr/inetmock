package rpc

import (
	"context"
	"github.com/baez90/inetmock/internal/endpoints"
)

type endpointsServer struct {
	endpointsManager endpoints.EndpointManager
}

func (e endpointsServer) GetEndpoints(_ context.Context, _ *GetEndpointsRequest) (*GetEndpointsResponse, error) {
	eps := rpcEndpointsFromEndpoints(e.endpointsManager.StartedEndpoints())
	return &GetEndpointsResponse{
		Endpoints: *eps,
	}, nil
}

func rpcEndpointsFromEndpoints(eps []endpoints.Endpoint) *[]*Endpoint {
	out := make([]*Endpoint, 0)
	for _, ep := range eps {
		out = append(out, rpcEndpointFromEndpoint(ep))
	}
	return &out
}

func rpcEndpointFromEndpoint(ep endpoints.Endpoint) *Endpoint {
	return &Endpoint{
		Id:            ep.Id().String(),
		Name:          ep.Name(),
		Handler:       ep.Handler(),
		ListenAddress: ep.Listen(),
		Port:          int32(ep.Port()),
	}
}
