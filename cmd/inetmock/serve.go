package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"go.uber.org/zap"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the INetMock server",
		Long:  ``,
		RunE:  startINetMock,
	}
)

func startINetMock(_ *cobra.Command, _ []string) (err error) {
	rpcAPI := rpc.NewINetMockAPI(
		serverApp.Config().APIConfig(),
		serverApp.Logger(),
		serverApp.Checker(),
		serverApp.EventStream(),
		serverApp.Config().AuditDataDir(),
		serverApp.Config().PCAPDataDir(),
	)
	logger := serverApp.Logger()

	cfg := serverApp.Config()
	endpointOrchestrator := serverApp.EndpointManager()

	for name, spec := range cfg.ListenerSpecs() {
		if spec.Name == "" {
			spec.Name = name
		}
		if err = endpointOrchestrator.RegisterListener(spec); err != nil {
			logger.Error("Failed to register listener", zap.Error(err))
			return
		}
	}

	errChan := serverApp.EndpointManager().StartEndpoints()
	if err = rpcAPI.StartServer(); err != nil {
		serverApp.Shutdown()
		logger.Error(
			"failed to start gRPC API",
			zap.Error(err),
		)
	}

loop:
	for {
		select {
		case err := <-errChan:
			logger.Error("got error from endpoint", zap.Error(err))
		case <-serverApp.Context().Done():
			break loop
		}
	}

	logger.Info("App context canceled - shutting down")

	rpcAPI.StopServer()
	return
}
