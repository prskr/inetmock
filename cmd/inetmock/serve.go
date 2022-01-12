package main

import (
	"io"
	"io/fs"
	"os"

	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/pcap"
	audit2 "gitlab.com/inetmock/inetmock/internal/pcap/consumers/audit"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"gitlab.com/inetmock/inetmock/internal/state"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	dhcpmock "gitlab.com/inetmock/inetmock/protocols/dhcp"
	"gitlab.com/inetmock/inetmock/protocols/dns/doh"
	dnsmock "gitlab.com/inetmock/inetmock/protocols/dns/mock"
	"gitlab.com/inetmock/inetmock/protocols/http/mock"
	"gitlab.com/inetmock/inetmock/protocols/http/proxy"
	"gitlab.com/inetmock/inetmock/protocols/metrics"
	"gitlab.com/inetmock/inetmock/protocols/pprof"
)

const (
	defaultEventBufferSize = 10
)

var (
	toClose  []io.Closer
	serveCmd = &cobra.Command{
		Use:          "serve",
		Short:        "Starts the INetMock server",
		Long:         ``,
		RunE:         startINetMock,
		SilenceUsage: true,
		PostRunE: func(*cobra.Command, []string) (err error) {
			for idx := range toClose {
				err = multierr.Append(err, toClose[idx].Close())
			}
			return
		},
	}
)

//nolint:gocyclo // central setup point
func startINetMock(_ *cobra.Command, _ []string) error {
	registry := endpoint.NewHandlerRegistry()

	var err error
	appLogger := serverApp.Logger()

	appLogger.Info("Starting the server")

	if err = cfg.Data.setup(); err != nil {
		appLogger.Error("Failed to setup data directories", zap.Error(err))
		return err
	}

	var stateStore state.KVStore
	if stateStore, err = state.NewDefault(state.WithPath(cfg.Data.State)); err != nil {
		appLogger.Error("Failed to setup state store", zap.Error(err))
		return err
	}
	toClose = append(toClose, stateStore)

	if cfg.TLS.CertCachePath, err = ensureDataDir(cfg.TLS.CertCachePath); err != nil {
		appLogger.Error("Failed to setup cert cache directory", zap.Error(err))
	}

	fakeFileFS := os.DirFS(cfg.Data.FakeFiles)

	var certStore cert.Store
	if certStore, err = cert.NewDefaultStore(cfg.TLS, appLogger.Named("CertStore")); err != nil {
		appLogger.Error("Failed to initialize cert store", zap.Error(err))
		return err
	}

	var eventStream audit.EventStream
	if eventStream, err = setupEventStream(appLogger); err != nil {
		return err
	}
	toClose = append(toClose, eventStream)

	var checker health.Checker
	if checker, err = health.NewFromConfig(appLogger.Named("health"), cfg.Health, certStore.TLSConfig()); err != nil {
		appLogger.Error("Failed to setup health checker", zap.Error(err))
		return err
	}

	if err = setupEndpointHandlers(registry, appLogger, eventStream, certStore, stateStore, fakeFileFS, checker); err != nil {
		appLogger.Error("Failed to run registration", zap.Error(err))
		return err
	}

	endpointOrchestrator := endpoint.NewOrchestrator(certStore, registry, appLogger.Named("orchestrator"))
	rpcAPI := rpc.NewINetMockAPI(cfg.APIURL(), appLogger, checker, eventStream, cfg.Data.Audit, cfg.Data.PCAP)

	for name, spec := range cfg.Listeners {
		if spec.Name == "" {
			spec.Name = name
		}
		if err := endpointOrchestrator.RegisterListener(spec); err != nil {
			appLogger.Error("Failed to register listener", zap.Error(err))
			return err
		}
	}

	errChan := endpointOrchestrator.StartEndpoints(serverApp.Context())
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
			switch e := err.(type) {
			case cmux.ErrNotMatched:
				appLogger.Error("Not matched error",
					zap.Bool("temporary", e.Temporary()),
					zap.Bool("timeoutError", e.Timeout()),
					zap.String("error", e.Error()),
				)
			default:
				appLogger.Error("got error from endpoint", zap.Error(err))
			}
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

func setupEndpointHandlers(
	registry endpoint.HandlerRegistry,
	logger logging.Logger,
	emitter audit.Emitter,
	certStore cert.Store,
	stateStore state.KVStore,
	fakeFileFS fs.FS,
	checker health.Checker,
) (err error) {
	mock.AddHTTPMock(registry, logger.Named("http_mock"), emitter, fakeFileFS)
	dnsmock.AddDNSMock(registry, logger.Named("dns_mock"), emitter)
	dhcpmock.AddDHCPMock(registry, logger.Named("dhcp_mock"), emitter, stateStore.WithSuffixes("dhcp_mock"))
	doh.AddDoH(registry, logger.Named("doh_mock"), emitter)
	pprof.AddPprof(registry, logger.Named("pprof"), emitter)
	if err = proxy.AddHTTPProxy(registry, logger.Named("http_proxy"), emitter, certStore); err != nil {
		return
	}
	if err = metrics.AddMetricsExporter(registry, logger.Named("metrics_exporter"), checker); err != nil {
		return
	}
	return nil
}

//nolint:deadcode
func startAuditConsumer(eventStream audit.EventStream) error {
	recorder := pcap.NewRecorder()

	auditConsumer := audit2.NewAuditConsumer("audit", eventStream)

	_, err := recorder.StartRecording(serverApp.Context(), "lo", auditConsumer)
	return err
}
