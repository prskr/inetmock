package rpc

import (
	"net"
	"net/url"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gitlab.com/inetmock/inetmock/internal/pcap"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
)

type INetMockAPI interface {
	StartServer() error
	StopServer()
}

type inetmockAPI struct {
	url           *url.URL
	server        *grpc.Server
	logger        logging.Logger
	checker       health.Checker
	eventStream   audit.EventStream
	auditDataDir  string
	pcapDataDir   string
	serverRunning chan struct{}
}

func NewINetMockAPI(
	u *url.URL,
	logger logging.Logger,
	checker health.Checker,
	eventStream audit.EventStream,
	auditDataDir, pcapDataDir string,
) INetMockAPI {
	return &inetmockAPI{
		url:          u,
		logger:       logger.Named("api"),
		checker:      checker,
		eventStream:  eventStream,
		auditDataDir: auditDataDir,
		pcapDataDir:  pcapDataDir,
	}
}

func (i *inetmockAPI) StartServer() (err error) {
	var lis net.Listener
	if lis, err = createListenerFromURL(i.url); err != nil {
		return
	}
	i.server = grpc.NewServer()

	rpc.RegisterHealthServer(i.server, &healthServer{
		checker: i.checker,
	})

	rpc.RegisterAuditServer(i.server, &auditServer{
		logger:           i.logger,
		eventStream:      i.eventStream,
		auditDataDirPath: i.auditDataDir,
	})

	rpc.RegisterPCAPServer(i.server, NewPCAPServer(i.pcapDataDir, pcap.NewRecorder()))

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
	case <-time.After(5 * time.Second):
		i.server.Stop()
	}
}

func (i *inetmockAPI) startServerAsync(listener net.Listener) {
	i.serverRunning = make(chan struct{})
	defer close(i.serverRunning)
	if err := i.server.Serve(listener); err != nil {
		i.logger.Error(
			"failed to start INetMock API",
			zap.Error(err),
		)
	}
}

func (i *inetmockAPI) isRunning() bool {
	select {
	case _, more := <-i.serverRunning:
		return more
	default:
		return true
	}
}

func createListenerFromURL(u *url.URL) (lis net.Listener, err error) {
	switch u.Scheme {
	case "unix":
		if _, err = os.Stat(u.Path); err == nil {
			if err = os.Remove(u.Path); err != nil {
				return
			}
		}
		lis, err = net.Listen(u.Scheme, u.Path)
	default:
		lis, err = net.Listen(u.Scheme, u.Host)
	}
	return
}
