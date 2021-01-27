package rpc

import (
	"context"

	"gitlab.com/inetmock/inetmock/internal/app"
)

type healthServer struct {
	UnimplementedHealthServer
	app app.App
}

func (h healthServer) GetHealth(_ context.Context, _ *HealthRequest) (resp *HealthResponse, err error) {
	checker := h.app.Checker()
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
