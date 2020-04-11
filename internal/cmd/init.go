package cmd

import (
	"github.com/baez90/inetmock/internal/plugins"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/baez90/inetmock/pkg/path"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var (
	appIsInitialized = false
)

func initApp() (err error) {
	if appIsInitialized {
		return
	}
	appIsInitialized = true
	logging.ConfigureLogging(
		logging.ParseLevel(logLevel),
		developmentLogs,
		map[string]interface{}{"cwd": path.WorkingDirectory()},
	)
	logger, _ = logging.CreateLogger()
	registry := plugins.Registry()

	if err = rootCmd.ParseFlags(os.Args); err != nil {
		return
	}

	if err = appConfig.ReadConfig(configFilePath); err != nil {
		logger.Error(
			"unrecoverable error occurred during reading the config file",
			zap.Error(err),
		)
		return
	}

	viperInst := viper.GetViper()
	pluginDir := viperInst.GetString("plugins-directory")
	if err = registry.LoadPlugins(pluginDir); err != nil {
		logger.Error("Failed to load plugins",
			zap.String("pluginsDirectory", pluginDir),
			zap.Error(err),
		)
	}

	pluginsCmd.AddCommand(registry.PluginCommands()...)

	return
}
