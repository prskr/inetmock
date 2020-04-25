package main

import (
	"crypto/tls"
	"fmt"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/cert"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

const (
	name = "tls_interceptor"
)

type tlsInterceptor struct {
	options                 tlsOptions
	logger                  logging.Logger
	listener                net.Listener
	certStore               cert.Store
	shutdownRequested       bool
	currentConnectionsCount *sync.WaitGroup
	currentConnections      map[uuid.UUID]*proxyConn
}

func (t *tlsInterceptor) Start(config api.HandlerConfig) (err error) {
	t.options = loadFromConfig(config.Options())
	addr := fmt.Sprintf("%s:%d", config.ListenAddress(), config.Port())

	t.logger = t.logger.With(
		zap.String("address", addr),
		zap.String("target", t.options.redirectionTarget.address()),
	)

	if t.listener, err = tls.Listen("tcp", addr, api.ServicesInstance().CertStore().TLSConfig()); err != nil {
		t.logger.Fatal(
			"failed to create tls listener",
			zap.Error(err),
		)
		err = fmt.Errorf(
			"failed to create tls listener: %w",
			err,
		)
		return
	}

	go t.startListener()
	return
}

func (t *tlsInterceptor) Shutdown() (err error) {
	t.logger.Info("Shutting down TLS interceptor")
	t.shutdownRequested = true
	done := make(chan struct{})
	go func() {
		t.currentConnectionsCount.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(5 * time.Second):
		for _, proxyConn := range t.currentConnections {
			if err = proxyConn.Close(); err != nil {
				t.logger.Error(
					"error while closing remaining proxy connections",
					zap.Error(err),
				)
				err = fmt.Errorf(
					"error while closing remaining proxy connections: %w",
					err,
				)
			}
		}
		return
	}
}

func (t *tlsInterceptor) startListener() {
	for !t.shutdownRequested {
		conn, err := t.listener.Accept()
		if err != nil {
			t.logger.Error(
				"error during accept",
				zap.Error(err),
			)
			continue
		}

		t.currentConnectionsCount.Add(1)
		go t.proxyConn(conn)
	}
}

func (t *tlsInterceptor) proxyConn(conn net.Conn) {
	defer conn.Close()
	defer t.currentConnectionsCount.Done()

	rAddr, err := net.ResolveTCPAddr("tcp", t.options.redirectionTarget.address())
	if err != nil {
		t.logger.Error(
			"failed to resolve proxy target",
			zap.Error(err),
		)
	}

	targetConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		t.logger.Error(
			"failed to connect to proxy target",
			zap.Error(err),
		)
		return
	}
	defer targetConn.Close()

	proxyCon := &proxyConn{
		source: conn,
		target: targetConn,
	}

	conUID := uuid.New()
	t.currentConnections[conUID] = proxyCon
	Pipe(conn, targetConn)
	delete(t.currentConnections, conUID)

	t.logger.Info(
		"connection closed",
		zap.String("remoteAddr", conn.RemoteAddr().String()),
	)
}
