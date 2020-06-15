package rpc

import (
	"context"
	"github.com/baez90/inetmock/pkg/health"
)

type healthServer struct {
}

func (h healthServer) GetHealth(_ context.Context, _ *HealthRequest) (resp *HealthResponse, err error) {
	checker := health.CheckerInstance()
	result := checker.IsHealthy()

	resp = &HealthResponse{
		OverallHealthState: HealthState(result.Status),
		ComponentsHealth:   map[string]*ComponentHealth{}}

	for component, status := range result.Components {
		resp.ComponentsHealth[component] = &ComponentHealth{
			State:   HealthState(status.Status),
			Message: status.Message,
		}
	}

	return
}
