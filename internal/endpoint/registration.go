//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/endpoint/handler_registry.mock.go -package=endpoint_mock
package endpoint

import (
	"fmt"
)

type (
	Registration    func(registry HandlerRegistry) error
	HandlerProvider func() ProtocolHandler
)

type HandlerRegistry interface {
	RegisterHandler(handlerRef HandlerReference, handlerProvider HandlerProvider)
	AvailableHandlers() []HandlerReference
	HandlerForName(handlerRef HandlerReference) (ProtocolHandler, bool)
}

func NewHandlerRegistry() HandlerRegistry {
	return &handlerRegistry{}
}

type handlerRegistry map[HandlerReference]HandlerProvider

func (h *handlerRegistry) AvailableHandlers() (availableHandlers []HandlerReference) {
	for s := range *h {
		availableHandlers = append(availableHandlers, s)
	}
	return
}

func (h *handlerRegistry) HandlerForName(handlerRef HandlerReference) (instance ProtocolHandler, ok bool) {
	var provider func() ProtocolHandler
	handlers := *h
	if provider, ok = handlers[handlerRef.ToLower()]; ok {
		instance = provider()
	}
	return
}

func (h *handlerRegistry) RegisterHandler(handlerRef HandlerReference, handlerProvider HandlerProvider) {
	handlerRef = handlerRef.ToLower()
	handlers := *h
	if _, exists := handlers[handlerRef]; exists {
		panic(fmt.Sprintf("handler with name %s is already registered - there's something strange...in the neighborhood", handlerRef))
	}
	handlers[handlerRef] = handlerProvider
}
