package rpc

import (
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	reqlog "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	v1Health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/internal/pcap"
	"inetmock.icb4dc0.de/inetmock/internal/rpc/middleware"
	"inetmock.icb4dc0.de/inetmock/netflow"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/health"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	rpcv1 "inetmock.icb4dc0.de/inetmock/pkg/rpc/v1"
)

const gracefulShutdownTimeout = 5 * time.Second

type INetMockAPI interface {
	StartServer() error
	StopServer()
}

type inetmockAPI struct {
	lock          sync.Mutex
	url           *url.URL
	server        *grpc.Server
	logger        logging.Logger
	checker       health.Checker
	eventStream   audit.EventStream
	fw            *netflow.Firewall
	nat           *netflow.NAT
	epHost        endpoint.Host
	auditDataDir  string
	pcapDataDir   string
	serverRunning chan struct{}
}

func NewINetMockAPI(
	u *url.URL,
	logger logging.Logger,
	checker health.Checker,
	eventStream audit.EventStream,
	fw *netflow.Firewall,
	nat *netflow.NAT,
	epHost endpoint.Host,
	auditDataDir, pcapDataDir string,
) INetMockAPI {
	return &inetmockAPI{
		url:          u,
		logger:       logger.Named("api"),
		checker:      checker,
		eventStream:  eventStream,
		fw:           fw,
		nat:          nat,
		epHost:       epHost,
		auditDataDir: auditDataDir,
		pcapDataDir:  pcapDataDir,
	}
}

func (i *inetmockAPI) StartServer() (err error) {
	var lis net.Listener
	if lis, err = createListenerFromURL(i.url); err != nil {
		return
	}
	i.server = grpc.NewServer(
		grpc.StreamInterceptor(
			grpcmiddleware.ChainStreamServer(
				recovery.StreamServerInterceptor(),
				prometheus.StreamServerInterceptor,
				reqlog.StreamServerInterceptor(i.logger.ZapLogger()),
			)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			recovery.UnaryServerInterceptor(),
			middleware.ContextErrorConverter,
			prometheus.UnaryServerInterceptor,
			reqlog.UnaryServerInterceptor(i.logger.ZapLogger()),
		)))

	v1Health.RegisterHealthServer(i.server, NewHealthServer(i.checker, 1*time.Second, i.logger))
	rpcv1.RegisterAuditServiceServer(i.server, NewAuditServiceServer(i.logger, i.eventStream, i.auditDataDir))
	rpcv1.RegisterPCAPServiceServer(i.server, NewPCAPServer(i.pcapDataDir, pcap.NewRecorder()))
	rpcv1.RegisterProfilingServiceServer(i.server, NewProfilingServer())
	rpcv1.RegisterEndpointOrchestratorServiceServer(i.server, NewEndpointOrchestratorServer(i.logger, i.epHost))
	rpcv1.RegisterNetFlowControlServiceServer(i.server, NewNetFlowControlServiceServer(i.fw, i.nat))

	reflection.Register(i.server)

	go i.startServerAsync(lis)
	return
}

func (i *inetmockAPI) StopServer() {
	if !i.isRunning() {
		i.logger.Info(
			"Skipping API server shutdown because server is not running",
		)
		return
	}
	gracefulStopChan := make(chan struct{})
	go func() {
		i.server.GracefulStop()
		close(gracefulStopChan)
	}()

	select {
	case <-gracefulStopChan:
	case <-time.After(gracefulShutdownTimeout):
		i.server.Stop()
	}
}

func (i *inetmockAPI) startServerAsync(listener net.Listener) {
	i.lock.Lock()
	i.serverRunning = make(chan struct{})
	i.lock.Unlock()

	defer func() {
		i.lock.Lock()
		close(i.serverRunning)
		i.lock.Unlock()
	}()
	if err := i.server.Serve(listener); err != nil {
		i.logger.Error(
			"failed to start INetMock API",
			zap.Error(err),
		)
	}
}

func (i *inetmockAPI) isRunning() bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	select {
	case _, more := <-i.serverRunning:
		return more
	default:
		return true
	}
}

func createListenerFromURL(u *url.URL) (net.Listener, error) {
	const socketDirectoryPermissions = 0x755

	switch u.Scheme {
	case "unix":
		directory := filepath.Dir(u.Path)
		if _, err := os.Stat(directory); err != nil {
			if err := os.MkdirAll(directory, socketDirectoryPermissions); err != nil {
				return nil, err
			}
		}
		if _, err := os.Stat(u.Path); err == nil {
			if err := os.Remove(u.Path); err != nil {
				return nil, err
			}
		}
		return net.Listen(u.Scheme, u.Path)
	default:
		return net.Listen(u.Scheme, u.Host)
	}
}
