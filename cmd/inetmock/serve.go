package main

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net"
	"os"
	"strings"
	"time"

	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/internal/rpc"
	"inetmock.icb4dc0.de/inetmock/internal/state"
	"inetmock.icb4dc0.de/inetmock/netflow"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/audit/sink"
	"inetmock.icb4dc0.de/inetmock/pkg/cert"
	"inetmock.icb4dc0.de/inetmock/pkg/health"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	dhcpmock "inetmock.icb4dc0.de/inetmock/protocols/dhcp"
	"inetmock.icb4dc0.de/inetmock/protocols/dns"
	"inetmock.icb4dc0.de/inetmock/protocols/dns/doh"
	dnsmock "inetmock.icb4dc0.de/inetmock/protocols/dns/mock"
	"inetmock.icb4dc0.de/inetmock/protocols/http/mock"
	"inetmock.icb4dc0.de/inetmock/protocols/http/proxy"
	"inetmock.icb4dc0.de/inetmock/protocols/metrics"
	"inetmock.icb4dc0.de/inetmock/protocols/pprof"
)

const (
	defaultEventBufferSize = 1000
	startGroupsTimeout     = 1 * time.Second
)

var (
	toClose  []io.Closer
	serveCmd = &cobra.Command{
		Use:          "serve",
		Short:        "Starts the INetMock server",
		Long:         ``,
		RunE:         startINetMock,
		SilenceUsage: true,
		PreRunE: func(*cobra.Command, []string) error {
			dns.ConfigureCache(
				dns.WithInitialSize(cfg.Caches.DNS.InitialCapacity),
				dns.WithTTL(cfg.Caches.DNS.TTL),
			)
			return nil
		},
		PostRunE: func(*cobra.Command, []string) (err error) {
			for idx := range toClose {
				err = errors.Join(err, toClose[idx].Close())
			}
			return
		},
	}
)

//nolint:funlen //startup code
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

	setupEndpointHandlers(registry, appLogger, eventStream, certStore, stateStore, fakeFileFS, checker)

	serverBuilder := endpoint.NewServerBuilder(certStore.TLSConfig(), registry, appLogger.Named("orchestrator"))
	srv := serverBuilder.Server()
	srv.ErrorHandler = append(srv.ErrorHandler, endpointErrorHandler(appLogger))

	packetSink := netflow.EmittingPacketSink{
		Lookup:  dns.GlobalCache(),
		Emitter: eventStream,
	}
	sinkOption := netflow.ErrorSinkOption{ErrorSink: netflow.LoggerErrorSink{Logger: appLogger.Named("netflow")}}
	firewall := netflow.NewFirewall(packetSink, sinkOption)
	nat := netflow.NewNAT(sinkOption)

	toClose = append(toClose, firewall, nat)
	rpcAPI := rpc.NewINetMockAPI(
		cfg.APIURL(),
		appLogger,
		checker,
		eventStream,
		firewall,
		nat,
		srv,
		cfg.Data.Audit,
		cfg.Data.PCAP,
	)

	for name, spec := range cfg.Listeners {
		if spec.Name == "" {
			spec.Name = name
		}
		if err := serverBuilder.ConfigureGroup(spec); err != nil {
			appLogger.Error("Failed to register listener", zap.Error(err))
			return err
		}
	}

	startGroupsCtx, startGroupsCancel := context.WithTimeout(serverApp.Context(), startGroupsTimeout)
	if err := srv.ServeGroups(startGroupsCtx); err != nil {
		startGroupsCancel()
		appLogger.Error("Failed to serve listener groups", zap.Error(err))
		return err
	}
	startGroupsCancel()

	srv.ShutdownOnCancel(serverApp.Context())

	if err := rpcAPI.StartServer(); err != nil {
		serverApp.Shutdown()
		appLogger.Error(
			"failed to start gRPC API",
			zap.Error(err),
		)
	}

	if err := initFirewall(firewall, cfg.NetFlow.Firewall); err != nil {
		return err
	}

	if err := initNAT(nat, cfg.NetFlow.NAT); err != nil {
		return err
	}

	appLogger.Info("App startup completed")

	<-serverApp.Context().Done()

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
) {
	mock.AddHTTPMock(registry, logger.Named("http_mock"), emitter, fakeFileFS)
	dnsmock.AddDNSMock(registry, logger.Named("dns_mock"), emitter)
	dhcpmock.AddDHCPMock(registry, logger.Named("dhcp_mock"), emitter, stateStore.WithSuffixes("dhcp_mock"))
	doh.AddDoH(registry, logger.Named("doh_mock"), emitter)
	pprof.AddPprof(registry, logger.Named("pprof"), emitter)
	proxy.AddHTTPProxy(registry, logger.Named("http_proxy"), emitter, certStore)
	metrics.AddMetricsExporter(registry, logger.Named("metrics_exporter"), checker)
}

func endpointErrorHandler(logger logging.Logger) endpoint.ErrorHandler {
	return endpoint.ErrorHandlerFunc(func(err error) {
		var (
			unmatched cmux.ErrNotMatched
			netOp     = new(net.OpError)
		)
		switch {
		case errors.As(err, &unmatched):
			if unmatched.Temporary() {
				return
			}
			logger.Error("Not matched error",
				zap.Bool("timeoutError", unmatched.Timeout()),
				zap.String("error", unmatched.Error()),
			)
		case errors.As(err, &netOp):
			if !strings.EqualFold(netOp.Op, "accept") && !netOp.Temporary() {
				logger.Error("got error from endpoint", zap.Error(err))
			}
		default:
			logger.Error("got error from endpoint", zap.Error(err))
		}
	})
}

func initFirewall(fw *netflow.Firewall, initialCfg map[string]netflow.FirewallInterfaceConfig) error {
	for nic, cfg := range initialCfg {
		if err := fw.AttachToInterface(nic, cfg); err != nil {
			return err
		}
	}
	return nil
}

func initNAT(nat *netflow.NAT, initialCfg map[string]netflow.NATTableSpec) error {
	for nic, spec := range initialCfg {
		if err := nat.AttachToInterface(nic, spec); err != nil {
			return err
		}
	}

	return nil
}
