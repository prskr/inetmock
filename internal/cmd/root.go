package cmd

import (
	"github.com/baez90/inetmock/internal/config"
	"github.com/baez90/inetmock/internal/endpoints"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	logger  logging.Logger
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
	endpointManager endpoints.EndpointManager
)

func init() {
	rootCmd.PersistentFlags().String("plugins-directory", "", "Directory where plugins should be loaded from")
	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file that should be used")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level to use")
	rootCmd.PersistentFlags().BoolVar(&developmentLogs, "development-logs", false, "Enable development mode logs")

	appConfig.InitConfig(rootCmd.PersistentFlags())
}

func startInetMock(cmd *cobra.Command, args []string) {
	for endpointName := range viper.GetStringMap(config.EndpointsKey) {
		handlerSubConfig := viper.Sub(strings.Join([]string{config.EndpointsKey, endpointName}, "."))
		handlerConfig := config.CreateMultiHandlerConfig(handlerSubConfig)
		if err := endpointManager.CreateEndpoint(endpointName, handlerConfig); err != nil {
			logger.Warn(
				"error occurred while creating endpoint",
				zap.String("endpointName", endpointName),
				zap.String("handlerName", handlerConfig.HandlerName()),
				zap.Error(err),
			)
		}
	}

	endpointManager.StartEndpoints()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// block until canceled
	s := <-signalChannel

	logger.Info(
		"got signal to quit",
		zap.String("signal", s.String()),
	)

	endpointManager.ShutdownEndpoints()
}

func ExecuteRootCommand() error {
	if err := initApp(); err != nil {
		return err
	}
	return rootCmd.Execute()
}
