package cmd

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
	rpcAPI := rpc.NewINetMockAPI(server)
	logger := server.Logger().Named("inetmock").With(zap.String("command", "serve"))

	for endpointName, endpointHandler := range server.Config().EndpointConfigs() {
		handlerSubConfig := server.Config().Viper().Sub(strings.Join([]string{config.EndpointsKey, endpointName, config.OptionsKey}, "."))
		endpointHandler.Options = handlerSubConfig
		if err := server.EndpointManager().CreateEndpoint(endpointName, endpointHandler); err != nil {
			logger.Warn(
				"error occurred while creating endpoint",
				zap.String("endpointName", endpointName),
				zap.String("handlerName", endpointHandler.Handler),
				zap.Error(err),
			)
		}
	}

	server.EndpointManager().StartEndpoints()
	if err := rpcAPI.StartServer(); err != nil {
		logger.Error(
			"failed to start gRPC API",
			zap.Error(err),
		)
	}

	<-server.Context().Done()

	logger.Info(
		"App context canceled - shutting down",
	)

	rpcAPI.StopServer()
	server.EndpointManager().ShutdownEndpoints()
}
