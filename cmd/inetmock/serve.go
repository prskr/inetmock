package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/pcap"
	"gitlab.com/inetmock/inetmock/internal/pcap/consumers/audit"
	"gitlab.com/inetmock/inetmock/internal/rpc"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the INetMock server",
		Long:  ``,
		RunE:  startINetMock,
	}
)

func startINetMock(_ *cobra.Command, _ []string) error {
	rpcAPI := rpc.NewINetMockAPI(
		serverApp.Config().APIURL(),
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
		if err := endpointOrchestrator.RegisterListener(spec); err != nil {
			logger.Error("Failed to register listener", zap.Error(err))
			return err
		}
	}

	errChan := serverApp.EndpointManager().StartEndpoints()
	if err := rpcAPI.StartServer(); err != nil {
		serverApp.Shutdown()
		logger.Error(
			"failed to start gRPC API",
			zap.Error(err),
		)
	}

	//nolint:gocritic
	/*if err = startAuditConsumer(); err != nil {
		return
	}*/

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
	return nil
}

//nolint:deadcode
func startAuditConsumer() error {
	recorder := pcap.NewRecorder()
	auditConsumer := audit.NewAuditConsumer("audit", serverApp.EventStream())

	return recorder.StartRecording(serverApp.Context(), "lo", auditConsumer)
}
