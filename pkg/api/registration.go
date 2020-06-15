//go:generate mockgen -source=registration.go -destination=./../../internal/mock/plugins/handler_registry.mock.go -package=plugins_mock
package api

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	registry              HandlerRegistry
	pluginFileNamePattern = regexp.MustCompile(`[\w\-]+\.so$`)
)

type handlerRegistry struct {
	handlers map[string]PluginInstanceFactory
}

func (h handlerRegistry) AvailableHandlers() (availableHandlers []string) {
	for s := range h.handlers {
		availableHandlers = append(availableHandlers, s)
	}
	return
}

func (h *handlerRegistry) HandlerForName(handlerName string) (instance ProtocolHandler, ok bool) {
	handlerName = strings.ToLower(handlerName)
	var provider PluginInstanceFactory
	if provider, ok = h.handlers[handlerName]; ok {
		instance = provider()
	}
	return
}

func (h *handlerRegistry) RegisterHandler(handlerName string, handlerProvider PluginInstanceFactory) {
	handlerName = strings.ToLower(handlerName)
	if _, exists := h.handlers[handlerName]; exists {
		panic(fmt.Sprintf("handler with name %s is already registered - there's something strange...in the neighborhood", handlerName))
	}
	h.handlers[handlerName] = handlerProvider
}

func Registry() HandlerRegistry {
	return registry
}

func init() {
	registry = &handlerRegistry{
		handlers: make(map[string]PluginInstanceFactory),
	}
}

type HandlerRegistry interface {
	RegisterHandler(handlerName string, handlerProvider PluginInstanceFactory)
	AvailableHandlers() []string
	HandlerForName(handlerName string) (ProtocolHandler, bool)
}
