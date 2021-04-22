package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/pcap"
	audit2 "gitlab.com/inetmock/inetmock/internal/pcap/consumers/audit"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	defaultEventBufferSize = 10
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
	registry := endpoint.NewHandlerRegistry()

	var err error
	appLogger := serverApp.Logger()

	if err = cfg.Data.setup(); err != nil {
		appLogger.Error("Failed to setup data directories", zap.Error(err))
		return err
	}

	for _, registration := range registrations {
		if err = registration(registry); err != nil {
			appLogger.Error("Failed to run registration", zap.Error(err))
			return err
		}
	}

	var certStore cert.Store
	if certStore, err = cert.NewDefaultStore(cfg.TLS, appLogger.Named("CertStore")); err != nil {
		appLogger.Error("Failed to initialize cert store", zap.Error(err))
		return err
	}

	var eventStream audit.EventStream
	if eventStream, err = setupEventStream(appLogger); err != nil {
		return err
	}

	var endpointOrchestrator = endpoint.NewOrchestrator(
		serverApp.Context(),
		certStore,
		registry,
		eventStream,
		appLogger,
	)

	checker := health.New()

	rpcAPI := rpc.NewINetMockAPI(
		cfg.APIURL(),
		appLogger,
		checker,
		eventStream,
		cfg.Data.Audit,
		cfg.Data.PCAP,
	)

	for name, spec := range cfg.Listeners {
		if spec.Name == "" {
			spec.Name = name
		}
		if err := endpointOrchestrator.RegisterListener(spec); err != nil {
			appLogger.Error("Failed to register listener", zap.Error(err))
			return err
		}
	}

	errChan := endpointOrchestrator.StartEndpoints()
	if err := rpcAPI.StartServer(); err != nil {
		serverApp.Shutdown()
		appLogger.Error(
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
			appLogger.Error("got error from endpoint", zap.Error(err))
		case <-serverApp.Context().Done():
			break loop
		}
	}

	appLogger.Info("App context canceled - shutting down")

	rpcAPI.StopServer()
	return nil
}

func setupEventStream(appLogger logging.Logger) (audit.EventStream, error) {
	var evenStream audit.EventStream
	var err error
	evenStream, err = audit.NewEventStream(
		appLogger.Named("EventStream"),
		audit.WithSinkBufferSize(defaultEventBufferSize),
	)

	if err != nil {
		appLogger.Error("Failed to initialize event stream", zap.Error(err))
		return nil, err
	}

	if err = evenStream.RegisterSink(serverApp.Context(), sink.NewLogSink(appLogger.Named("LogSink"))); err != nil {
		appLogger.Error("Failed to register log sink to event stream", zap.Error(err))
		return nil, err
	}

	var metricSink audit.Sink
	if metricSink, err = sink.NewMetricSink(); err != nil {
		appLogger.Error("Failed to setup metrics sink", zap.Error(err))
		return nil, err
	}

	if err = evenStream.RegisterSink(serverApp.Context(), metricSink); err != nil {
		appLogger.Error("Failed to register metric sink", zap.Error(err))
		return nil, err
	}

	return evenStream, nil
}

//nolint:deadcode
func startAuditConsumer(eventStream audit.EventStream) error {
	recorder := pcap.NewRecorder()

	auditConsumer := audit2.NewAuditConsumer("audit", eventStream)

	return recorder.StartRecording(serverApp.Context(), "lo", auditConsumer)
}
