//go:generate mockgen -source=loading.go -destination=./../../internal/mock/plugins/handler_registry.mock.go -package=plugins_mock
package plugins

import (
	"fmt"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/path"
	"os"
	"path/filepath"
	"plugin"
	"regexp"
	"strings"
)

var (
	registry              HandlerRegistry
	pluginFileNamePattern = regexp.MustCompile(`[\w\-]+\.so$`)
)

type HandlerRegistry interface {
	LoadPlugins(pluginsPath string) error
	AvailableHandlers() []string
	RegisterHandler(handlerName string, handlerProvider api.PluginInstanceFactory)
	HandlerForName(handlerName string) (api.ProtocolHandler, bool)
}

type handlerRegistry struct {
	handlers map[string]api.PluginInstanceFactory
}

func (h handlerRegistry) AvailableHandlers() (availableHandlers []string) {
	for s := range h.handlers {
		availableHandlers = append(availableHandlers, s)
	}
	return
}

func (h *handlerRegistry) HandlerForName(handlerName string) (instance api.ProtocolHandler, ok bool) {
	handlerName = strings.ToLower(handlerName)
	var provider api.PluginInstanceFactory
	if provider, ok = h.handlers[handlerName]; ok {
		instance = provider()
	}
	return
}

func (h *handlerRegistry) RegisterHandler(handlerName string, handlerProvider api.PluginInstanceFactory) {
	handlerName = strings.ToLower(handlerName)
	if _, exists := h.handlers[handlerName]; exists {
		panic(fmt.Sprintf("handler with name %s is already registered - there's something strange...in the neighborhood", handlerName))
	}
	h.handlers[handlerName] = handlerProvider
}

func (h *handlerRegistry) LoadPlugins(pluginsPath string) (err error) {
	if !path.DirExists(pluginsPath) {
		err = fmt.Errorf("plugins path %s does not exist or is not accessible", pluginsPath)
		return
	}
	err = filepath.Walk(pluginsPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && pluginFileNamePattern.MatchString(info.Name()) {
			if _, err := plugin.Open(path); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return
	}

	err = nil
	return
}

func Registry() HandlerRegistry {
	return registry
}

func init() {
	registry = &handlerRegistry{
		handlers: make(map[string]api.PluginInstanceFactory),
	}
}
