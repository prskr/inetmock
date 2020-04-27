//go:generate mockgen -source=loading.go -destination=./../../internal/mock/plugins/handler_registry_mock.go -package=plugins_mock
package plugins

import (
	"fmt"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/path"
	"github.com/spf13/cobra"
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
	RegisterHandler(handlerName string, handlerProvider api.PluginInstanceFactory, subCommands ...*cobra.Command)
	HandlerForName(handlerName string) (api.ProtocolHandler, bool)
	PluginCommands() []*cobra.Command
}

type handlerRegistry struct {
	handlers       map[string]api.PluginInstanceFactory
	pluginCommands []*cobra.Command
}

func (h handlerRegistry) AvailableHandlers() (availableHandlers []string) {
	for s := range h.handlers {
		availableHandlers = append(availableHandlers, s)
	}
	return
}

func (h handlerRegistry) PluginCommands() []*cobra.Command {
	return h.pluginCommands
}

func (h *handlerRegistry) HandlerForName(handlerName string) (instance api.ProtocolHandler, ok bool) {
	handlerName = strings.ToLower(handlerName)
	var provider api.PluginInstanceFactory
	if provider, ok = h.handlers[handlerName]; ok {
		instance = provider()
	}
	return
}

func (h *handlerRegistry) RegisterHandler(handlerName string, handlerProvider api.PluginInstanceFactory, subCommands ...*cobra.Command) {
	handlerName = strings.ToLower(handlerName)
	if _, exists := h.handlers[handlerName]; exists {
		panic(fmt.Sprintf("handler with name %s is already registered - there's something strange...in the neighborhood", handlerName))
	}
	h.handlers[handlerName] = handlerProvider

	if len(subCommands) > 0 {
		pluginCmds := &cobra.Command{
			Use: handlerName,
		}
		pluginCmds.AddCommand(subCommands...)
		h.pluginCommands = append(h.pluginCommands, pluginCmds)
	}
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
