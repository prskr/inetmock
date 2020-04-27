package cmd

import (
	"github.com/baez90/inetmock/internal/endpoints"
	"github.com/baez90/inetmock/pkg/config"
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

func startINetMock(_ *cobra.Command, _ []string) {
	endpointManager = endpoints.NewEndpointManager(logger)
	cfg := config.Instance()
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
