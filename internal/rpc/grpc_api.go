package rpc

import (
	"github.com/baez90/inetmock/internal/endpoints"
	"github.com/baez90/inetmock/internal/plugins"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/url"
	"time"
)

type INetMockAPI interface {
	StartServer() error
	StopServer()
}

type inetmockAPI struct {
	url             *url.URL
	server          *grpc.Server
	endpointManager endpoints.EndpointManager
	registry        plugins.HandlerRegistry
	logger          logging.Logger
	serverRunning   bool
}

func NewINetMockAPI(
	config config.Config,
	epMgr endpoints.EndpointManager,
	registry plugins.HandlerRegistry,
) INetMockAPI {
	return &inetmockAPI{
		url:             config.APIConfig().ListenURL(),
		endpointManager: epMgr,
		registry:        registry,
	}
}

func (i *inetmockAPI) StartServer() (err error) {
	var lis net.Listener
	if lis, err = createListenerFromURL(i.url); err != nil {
		return
	}
	i.server = grpc.NewServer()

	RegisterHandlersServer(i.server, &handlersServer{
		registry: i.registry,
	})
	RegisterEndpointsServer(i.server, &endpointsServer{
		endpointsManager: i.endpointManager,
	})

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
		lis, err = net.Listen(url.Scheme, url.Path)
	default:
		lis, err = net.Listen(url.Scheme, url.Host)
	}
	return
}
