package cmd

import (
	"github.com/baez90/inetmock/internal/config"
	"github.com/baez90/inetmock/internal/plugins"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var (
	logger  *zap.Logger
	rootCmd = cobra.Command{
		Use:   "",
		Short: "INetMock is lightweight internet mock",
		Run:   startInetMock,
	}

	configFilePath  string
	logLevel        string
	developmentLogs bool
	handlers        []api.ProtocolHandler
	appConfig       = config.CreateConfig()
)

func init() {
	rootCmd.PersistentFlags().String("plugins-directory", "", "Directory where plugins should be loaded from")
	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file that should be used")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level to use")
	rootCmd.PersistentFlags().BoolVar(&developmentLogs, "development-logs", false, "Enable development mode logs")

	appConfig.InitConfig(rootCmd.PersistentFlags())
}

func startInetMock(cmd *cobra.Command, args []string) {
	registry := plugins.Registry()
	var wg sync.WaitGroup

	//todo introduce endpoint type and move startup and shutdown to this type

	for key, val := range viper.GetStringMap(config.EndpointsKey) {
		handlerSubConfig := viper.Sub(strings.Join([]string{config.EndpointsKey, key, config.OptionsKey}, "."))
		pluginConfig := config.CreateHandlerConfig(val, handlerSubConfig)
		logger.Info(key, zap.Any("value", pluginConfig))

		if handler, ok := registry.HandlerForName(pluginConfig.HandlerName()); ok {
			handlers = append(handlers, handler)
			go startEndpoint(handler, pluginConfig, logger)
			wg.Add(1)
		} else {
			logger.Warn(
				"no matching handler registered",
				zap.String("handler", pluginConfig.HandlerName()),
			)
		}
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// block until canceled
	s := <-signalChannel

	logger.Info(
		"got signal to quit",
		zap.String("signal", s.String()),
	)

	for _, handler := range handlers {
		go shutdownEndpoint(handler, &wg, logger)
	}

	wg.Wait()
}

func startEndpoint(handler api.ProtocolHandler, config config.HandlerConfig, logger *zap.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal(
				"recovered panic during startup of endpoint",
				zap.Any("recovered", r),
			)
		}
	}()
	handler.Start(config)
}

func shutdownEndpoint(handler api.ProtocolHandler, wg *sync.WaitGroup, logger *zap.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal(
				"recovered panic during shutdown of endpoint",
				zap.Any("recovered", r),
			)
		}
	}()
	handler.Shutdown(wg)
}

func ExecuteRootCommand() error {
	if err := initApp(); err != nil {
		return err
	}
	return rootCmd.Execute()
}
