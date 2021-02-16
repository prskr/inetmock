package rpc

import (
	"context"

	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
)

type healthServer struct {
	rpc.UnimplementedHealthServer
	checker health.Checker
}

func (h healthServer) GetHealth(_ context.Context, _ *rpc.HealthRequest) (resp *rpc.HealthResponse, err error) {
	result := h.checker.IsHealthy()

	resp = &rpc.HealthResponse{
		OverallHealthState: rpc.HealthState(result.Status),
		ComponentsHealth:   map[string]*rpc.ComponentHealth{}}

	for component, status := range result.Components {
		resp.ComponentsHealth[component] = &rpc.ComponentHealth{
			State:   rpc.HealthState(status.Status),
			Message: status.Message,
		}
	}

	return
}
