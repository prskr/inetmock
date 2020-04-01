package plugins

import (
	"fmt"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/spf13/cobra"
	"path/filepath"
	"plugin"
)

var (
	registry HandlerRegistry
)

type HandlerRegistry interface {
	LoadPlugins(pluginsPath string) error
	RegisterHandler(handlerName string, handlerProvider api.PluginInstanceFactory, subCommands ...*cobra.Command)
	HandlerForName(handlerName string) (api.ProtocolHandler, bool)
	PluginCommands() []*cobra.Command
}

type handlerRegistry struct {
	handlers       map[string]api.PluginInstanceFactory
	pluginCommands []*cobra.Command
}

func (h handlerRegistry) PluginCommands() []*cobra.Command {
	return h.pluginCommands
}

func (h *handlerRegistry) HandlerForName(handlerName string) (instance api.ProtocolHandler, ok bool) {
	var provider api.PluginInstanceFactory
	if provider, ok = h.handlers[handlerName]; ok {
		instance = provider()
	}
	return
}

func (h *handlerRegistry) RegisterHandler(handlerName string, handlerProvider api.PluginInstanceFactory, subCommands ...*cobra.Command) {
	if _, exists := h.handlers[handlerName]; exists {
		panic(fmt.Sprintf("plugin %s already registered - there's something strange...in the neighborhood"))
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
	var plugins []string
	if plugins, err = filepath.Glob(fmt.Sprintf("%s%c*.so", pluginsPath, filepath.Separator)); err != nil {
		return
	}

	for _, pluginSo := range plugins {
		if _, err = plugin.Open(pluginSo); err != nil {
			return
		}
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
