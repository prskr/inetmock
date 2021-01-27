package main

import (
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"go.uber.org/zap"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the INetMock server",
		Long:  ``,
		Run:   startINetMock,
	}
)

func startINetMock(_ *cobra.Command, _ []string) {
	rpcAPI := rpc.NewINetMockAPI(serverApp)
	logger := serverApp.Logger().Named("inetmock").With(zap.String("command", "serve"))

	for endpointName, endpointHandler := range serverApp.Config().EndpointConfigs() {
		handlerSubConfig := serverApp.Config().Viper().Sub(strings.Join([]string{config.EndpointsKey, endpointName, config.OptionsKey}, "."))
		endpointHandler.Options = handlerSubConfig
		if err := serverApp.EndpointManager().CreateEndpoint(endpointName, endpointHandler); err != nil {
			logger.Warn(
				"error occurred while creating endpoint",
				zap.String("endpointName", endpointName),
				zap.String("handlerName", endpointHandler.Handler),
				zap.Error(err),
			)
		}
	}

	serverApp.EndpointManager().StartEndpoints()
	if err := rpcAPI.StartServer(); err != nil {
		logger.Error(
			"failed to start gRPC API",
			zap.Error(err),
		)
	}

	<-serverApp.Context().Done()

	logger.Info(
		"App context canceled - shutting down",
	)

	rpcAPI.StopServer()
	serverApp.EndpointManager().ShutdownEndpoints()
}
