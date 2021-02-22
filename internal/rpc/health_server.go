package rpc

import (
	"context"

	"gitlab.com/inetmock/inetmock/pkg/health"
	v1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

var (
	_ v1.HealthServiceServer = (*healthServer)(nil)
)

type healthServer struct {
	v1.UnimplementedHealthServiceServer
	checker health.Checker
}

func (h healthServer) GetHealth(_ context.Context, _ *v1.GetHealthRequest) (resp *v1.GetHealthResponse, err error) {
	result := h.checker.IsHealthy()

	resp = &v1.GetHealthResponse{
		OverallHealthState: v1.HealthState(result.Status),
		ComponentsHealth:   map[string]*v1.ComponentHealth{}}

	for component, status := range result.Components {
		resp.ComponentsHealth[component] = &v1.ComponentHealth{
			State:   v1.HealthState(status.Status),
			Message: status.Message,
		}
	}

	return
}
