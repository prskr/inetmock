package tls_interceptor

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

const (
	name = "tls_interceptor"
)

type tlsInterceptor struct {
	name                    string
	options                 tlsOptions
	logger                  logging.Logger
	listener                net.Listener
	shutdownRequested       bool
	currentConnectionsCount *sync.WaitGroup
	currentConnections      map[uuid.UUID]*proxyConn
	connectionsMutex        *sync.Mutex
}

func (t *tlsInterceptor) Start(ctx api.PluginContext, config config.HandlerConfig) (err error) {
	t.name = config.HandlerName

	if err = config.Options.Unmarshal(&t.options); err != nil {
		return
	}

	t.logger = t.logger.With(
		zap.String("handler_name", config.HandlerName),
		zap.String("address", config.ListenAddr()),
		zap.String("Target", t.options.Target.address()),
	)

	if t.listener, err = tls.Listen("tcp", config.ListenAddr(), ctx.CertStore().TLSConfig()); err != nil {
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

func (t *tlsInterceptor) Shutdown(ctx context.Context) (err error) {
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
	case <-ctx.Done():
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

		handledRequestCounter.WithLabelValues(t.name).Inc()
		openConnectionsGauge.WithLabelValues(t.name).Inc()
		t.currentConnectionsCount.Add(1)
		go t.proxyConn(conn)
	}
}

func (t *tlsInterceptor) proxyConn(conn net.Conn) {
	timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(t.name))
	defer func() {
		_ = conn.Close()
		t.currentConnectionsCount.Done()
		openConnectionsGauge.WithLabelValues(t.name).Dec()
		timer.ObserveDuration()
	}()

	rAddr, err := net.ResolveTCPAddr("tcp", t.options.Target.address())
	if err != nil {
		t.logger.Error(
			"failed to resolve proxy Target",
			zap.Error(err),
		)
	}

	targetConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		t.logger.Error(
			"failed to connect to proxy Target",
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
	t.storeConnection(conUID, proxyCon)
	Pipe(conn, targetConn)
	t.cleanConnection(conUID)

	switch tlsConn := conn.(type) {
	case *tls.Conn:
		if tlsConn.Handshake() != nil {
			t.logger.Error(
				"error occurred during TLS handshake",
				zap.Error(tlsConn.Handshake()),
			)
		}
	}

	t.logger.Info(
		"connection closed",
		zap.String("remoteAddr", conn.RemoteAddr().String()),
	)
}

func (t *tlsInterceptor) storeConnection(connUUID uuid.UUID, conn *proxyConn) {
	t.connectionsMutex.Lock()
	defer t.connectionsMutex.Unlock()
	t.currentConnections[connUUID] = conn
}

func (t *tlsInterceptor) cleanConnection(connUUID uuid.UUID) {
	t.connectionsMutex.Lock()
	defer t.connectionsMutex.Unlock()
	if _, ok := t.currentConnections[connUUID]; ok {
		delete(t.currentConnections, connUUID)
	}
}
