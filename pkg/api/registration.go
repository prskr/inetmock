//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/plugins/handler_registry.mock.go -package=plugins_mock
package api

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	pluginFileNamePattern = regexp.MustCompile(`[\w\-]+\.so$`)
)

type Registration func(registry HandlerRegistry) error

type HandlerRegistry interface {
	RegisterHandler(handlerName string, handlerProvider func() ProtocolHandler)
	AvailableHandlers() []string
	HandlerForName(handlerName string) (ProtocolHandler, bool)
}

func NewHandlerRegistry() HandlerRegistry {
	return &handlerRegistry{
		handlers: make(map[string]func() ProtocolHandler),
	}
}

type handlerRegistry struct {
	handlers map[string]func() ProtocolHandler
}

func (h handlerRegistry) AvailableHandlers() (availableHandlers []string) {
	for s := range h.handlers {
		availableHandlers = append(availableHandlers, s)
	}
	return
}

func (h *handlerRegistry) HandlerForName(handlerName string) (instance ProtocolHandler, ok bool) {
	handlerName = strings.ToLower(handlerName)
	var provider func() ProtocolHandler
	if provider, ok = h.handlers[handlerName]; ok {
		instance = provider()
	}
	return
}

func (h *handlerRegistry) RegisterHandler(handlerName string, handlerProvider func() ProtocolHandler) {
	handlerName = strings.ToLower(handlerName)
	if _, exists := h.handlers[handlerName]; exists {
		panic(fmt.Sprintf("handler with name %s is already registered - there's something strange...in the neighborhood", handlerName))
	}
	h.handlers[handlerName] = handlerProvider
}
