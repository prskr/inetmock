package rpc

import (
	"context"
	"github.com/baez90/inetmock/internal/endpoints"
)

type endpointsServer struct {
	endpointsManager endpoints.EndpointManager
}

func (e endpointsServer) GetEndpoints(_ context.Context, _ *GetEndpointsRequest) (*GetEndpointsResponse, error) {
	eps := make([]*Endpoint, 0)
	for _, ep := range e.endpointsManager.StartedEndpoints() {
		eps = append(eps, &Endpoint{
			Id:            ep.Id().String(),
			Name:          ep.Name(),
			Handler:       ep.Handler(),
			ListenAddress: ep.Listen(),
			Port:          int32(ep.Port()),
		})
	}
	return &GetEndpointsResponse{
		Endpoints: eps,
	}, nil
}
