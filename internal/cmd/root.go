package cmd

import (
	"github.com/baez90/inetmock/internal/plugins"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/baez90/inetmock/pkg/path"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

var (
	logger  logging.Logger
	rootCmd *cobra.Command

	pluginsDirectory string
	configFilePath   string
	logLevel         string
	developmentLogs  bool
)

func init() {
	cobra.OnInitialize(onInit)
	rootCmd = &cobra.Command{
		Use:   "",
		Short: "INetMock is lightweight internet mock",
	}

	rootCmd.PersistentFlags().StringVar(&pluginsDirectory, "plugins-directory", "", "Directory where plugins should be loaded from")
	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file that should be used")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level to use")
	rootCmd.PersistentFlags().BoolVar(&developmentLogs, "development-logs", false, "Enable development mode logs")

	rootCmd.AddCommand(
		serveCmd,
		generateCaCmd,
	)
}

func onInit() {
	logging.ConfigureLogging(
		logging.ParseLevel(logLevel),
		developmentLogs,
		map[string]interface{}{"cwd": path.WorkingDirectory()},
	)

	logger, _ = logging.CreateLogger()
	config.CreateConfig(rootCmd.Flags())
	appConfig := config.Instance()

	if err := appConfig.ReadConfig(configFilePath); err != nil {
		logger.Error(
			"failed to read config file",
			zap.Error(err),
		)
	}

	if err := api.InitServices(appConfig, logger); err != nil {
		logger.Error(
			"failed to initialize app services",
			zap.Error(err),
		)
	}

	registry := plugins.Registry()

	cfg := config.Instance()

	pluginLoadStartTime := time.Now()
	if err := registry.LoadPlugins(cfg.PluginsDir()); err != nil {
		logger.Error("Failed to load plugins",
			zap.String("pluginsDirectory", cfg.PluginsDir()),
			zap.Error(err),
		)
	}
	pluginLoadDuration := time.Since(pluginLoadStartTime)
	logger.Info(
		"loading plugins completed",
		zap.Duration("pluginLoadDuration", pluginLoadDuration),
	)
}

func ExecuteRootCommand() error {
	return rootCmd.Execute()
}
