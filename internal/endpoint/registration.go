//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/endpoint/handler_registry.mock.go -package=endpoint_mock
package endpoint

import (
	"fmt"
)

type Registration func(registry HandlerRegistry) error

type HandlerRegistry interface {
	RegisterHandler(handlerRef HandlerReference, handlerProvider func() ProtocolHandler)
	AvailableHandlers() []HandlerReference
	HandlerForName(handlerRef HandlerReference) (ProtocolHandler, bool)
}

func NewHandlerRegistry() HandlerRegistry {
	return &handlerRegistry{
		handlers: make(map[HandlerReference]func() ProtocolHandler),
	}
}

type handlerRegistry struct {
	handlers map[HandlerReference]func() ProtocolHandler
}

func (h handlerRegistry) AvailableHandlers() (availableHandlers []HandlerReference) {
	for s := range h.handlers {
		availableHandlers = append(availableHandlers, s)
	}
	return
}

func (h *handlerRegistry) HandlerForName(handlerRef HandlerReference) (instance ProtocolHandler, ok bool) {
	var provider func() ProtocolHandler
	if provider, ok = h.handlers[handlerRef.ToLower()]; ok {
		instance = provider()
	}
	return
}

func (h *handlerRegistry) RegisterHandler(handlerRef HandlerReference, handlerProvider func() ProtocolHandler) {
	handlerRef = handlerRef.ToLower()
	if _, exists := h.handlers[handlerRef]; exists {
		panic(fmt.Sprintf("handler with name %s is already registered - there's something strange...in the neighborhood", handlerRef))
	}
	h.handlers[handlerRef] = handlerProvider
}
