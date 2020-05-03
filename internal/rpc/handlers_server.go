package rpc

import (
	"context"
	"github.com/baez90/inetmock/internal/plugins"
)

type handlersServer struct {
	registry plugins.HandlerRegistry
}

func (h *handlersServer) GetHandlers(_ context.Context, _ *GetHandlersRequest) (*GetHandlersResponse, error) {
	return &GetHandlersResponse{
		Handlers: h.registry.AvailableHandlers(),
	}, nil
}
