package rpc

import (
	"context"

	"gitlab.com/inetmock/inetmock/pkg/api"
)

type handlersServer struct {
	UnimplementedHandlersServer
	registry api.HandlerRegistry
}

func (h *handlersServer) GetHandlers(_ context.Context, _ *GetHandlersRequest) (*GetHandlersResponse, error) {
	return &GetHandlersResponse{
		Handlers: h.registry.AvailableHandlers(),
	}, nil
}
