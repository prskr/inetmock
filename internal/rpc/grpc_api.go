package rpc

import (
	"net"
	"net/url"
	"os"
	"time"

	app2 "gitlab.com/inetmock/inetmock/internal/app"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type INetMockAPI interface {
	StartServer() error
	StopServer()
}

type inetmockAPI struct {
	app           app2.App
	url           *url.URL
	server        *grpc.Server
	logger        logging.Logger
	serverRunning bool
}

func NewINetMockAPI(
	app app2.App,
) INetMockAPI {
	return &inetmockAPI{
		app:    app,
		url:    app.Config().APIConfig().ListenURL(),
		logger: app.Logger().Named("api"),
	}
}

func (i *inetmockAPI) StartServer() (err error) {
	var lis net.Listener
	if lis, err = createListenerFromURL(i.url); err != nil {
		return
	}
	i.server = grpc.NewServer()

	rpc.RegisterHealthServer(i.server, &healthServer{
		app: i.app,
	})

	rpc.RegisterAuditServer(i.server, &auditServer{
		logger:      i.app.Logger(),
		eventStream: i.app.EventStream(),
	})

	reflection.Register(i.server)

	go i.startServerAsync(lis)
	return
}

func (i *inetmockAPI) StopServer() {
	if !i.serverRunning {
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
	i.serverRunning = true
	if err := i.server.Serve(listener); err != nil {
		i.serverRunning = false
		i.logger.Error(
			"failed to start INetMock API",
			zap.Error(err),
		)
	}
}

func createListenerFromURL(url *url.URL) (lis net.Listener, err error) {
	switch url.Scheme {
	case "unix":
		if _, err = os.Stat(url.Path); err == nil {
			if err = os.Remove(url.Path); err != nil {
				return
			}
		}
		lis, err = net.Listen(url.Scheme, url.Path)
	default:
		lis, err = net.Listen(url.Scheme, url.Host)
	}
	return
}
