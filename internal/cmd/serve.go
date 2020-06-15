package cmd

import (
	"github.com/baez90/inetmock/internal/endpoints"
	"github.com/baez90/inetmock/internal/rpc"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/health"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/baez90/inetmock/pkg/path"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	endpointManager endpoints.EndpointManager
	serveCmd        = &cobra.Command{
		Use:   "serve",
		Short: "Starts the INetMock server",
		Long:  ``,
		Run:   startINetMock,
	}
)

func onServerInit() (err error) {
	logging.ConfigureLogging(
		logging.ParseLevel(logLevel),
		developmentLogs,
		map[string]interface{}{"cwd": path.WorkingDirectory()},
	)

	logger, _ = logging.CreateLogger()
	config.CreateConfig(serverCmd.Flags())
	appConfig := config.Instance()

	if err = appConfig.ReadConfig(configFilePath); err != nil {
		logger.Error(
			"failed to read config file",
			zap.Error(err),
		)
		return
	}

	if err = api.InitServices(appConfig, logger); err != nil {
		logger.Error(
			"failed to initialize app services",
			zap.Error(err),
		)
		return
	}
	return nil
}

func startINetMock(_ *cobra.Command, _ []string) {
	if err := onServerInit(); err != nil {
		panic(err)
	}
	endpointManager = endpoints.NewEndpointManager(health.CheckerInstance(), logger)
	cfg := config.Instance()
	rpcAPI := rpc.NewINetMockAPI(
		cfg,
		endpointManager,
		api.Registry(),
	)

	for endpointName, endpointHandler := range cfg.EndpointConfigs() {
		handlerSubConfig := cfg.Viper().Sub(strings.Join([]string{config.EndpointsKey, endpointName, config.OptionsKey}, "."))
		endpointHandler.Options = handlerSubConfig
		if err := endpointManager.CreateEndpoint(endpointName, endpointHandler); err != nil {
			logger.Warn(
				"error occurred while creating endpoint",
				zap.String("endpointName", endpointName),
				zap.String("handlerName", endpointHandler.Handler),
				zap.Error(err),
			)
		}
	}

	endpointManager.StartEndpoints()
	if err := rpcAPI.StartServer(); err != nil {
		logger.Error(
			"failed to start gRPC API",
			zap.Error(err),
		)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// block until canceled
	s := <-signalChannel

	logger.Info(
		"got signal to quit",
		zap.String("signal", s.String()),
	)

	rpcAPI.StopServer()
	endpointManager.ShutdownEndpoints()
}
