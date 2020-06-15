package rpc

import (
	"context"
	"github.com/baez90/inetmock/pkg/api"
)

type handlersServer struct {
	registry api.HandlerRegistry
}

func (h *handlersServer) GetHandlers(_ context.Context, _ *GetHandlersRequest) (*GetHandlersResponse, error) {
	return &GetHandlersResponse{
		Handlers: h.registry.AvailableHandlers(),
	}, nil
}
