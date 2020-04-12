package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/baez90/inetmock/internal/config"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

const (
	name = "tls_interceptor"
)

type tlsInterceptor struct {
	logger                  *zap.Logger
	listener                net.Listener
	certStore               *certStore
	options                 *tlsOptions
	shutdownRequested       bool
	currentConnectionsCount *sync.WaitGroup
	currentConnections      []*proxyConn
}

func (t *tlsInterceptor) Run(config config.HandlerConfig) {
	var err error
	t.options = loadFromConfig(config.Options())
	addr := fmt.Sprintf("%s:%d", config.ListenAddress(), config.Port())

	t.logger = t.logger.With(
		zap.String("address", addr),
		zap.String("target", t.options.redirectionTarget.address()),
	)

	t.certStore = &certStore{
		options:   t.options,
		certCache: make(map[string]*tls.Certificate),
		logger:    t.logger,
	}

	if err = t.certStore.initCaCert(); err != nil {
		t.logger.Error(
			"failed to initialize CA cert",
			zap.Error(err),
		)
	}

	rootCaPool := x509.NewCertPool()
	rootCaPool.AddCert(t.certStore.caCert)

	tlsConfig := &tls.Config{
		GetCertificate: t.getCert,
		RootCAs:        rootCaPool,
	}

	if t.listener, err = tls.Listen("tcp", addr, tlsConfig); err != nil {
		t.logger.Fatal(
			"failed to create tls listener",
			zap.Error(err),
		)
		return
	}

	go t.startListener()
}

func (t *tlsInterceptor) Shutdown(wg *sync.WaitGroup) {
	t.logger.Info("Shutting down TLS interceptor")
	t.shutdownRequested = true
	done := make(chan struct{})
	go func() {
		t.currentConnectionsCount.Wait()
		close(done)
	}()

	select {
	case <-done:
		wg.Done()
	case <-time.After(5 * time.Second):
		for _, proxyConn := range t.currentConnections {
			if err := proxyConn.Close(); err != nil {
				t.logger.Error(
					"error while closing remaining proxy connections",
					zap.Error(err),
				)
			}
		}
		wg.Done()
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

func (t *tlsInterceptor) getCert(info *tls.ClientHelloInfo) (cert *tls.Certificate, err error) {
	var localIp string
	if localIp, err = extractIPFromAddress(info.Conn.LocalAddr().String()); err != nil {
		localIp = "127.0.0.1"
	}
	if cert, err = t.certStore.getCertificate(info.ServerName, localIp); err != nil {
		t.logger.Error(
			"error while resolving certificate",
			zap.String("serverName", info.ServerName),
			zap.String("localAddr", localIp),
			zap.Error(err),
		)
	}

	return
}

func (t *tlsInterceptor) proxyConn(conn net.Conn) {
	defer conn.Close()

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

	t.currentConnections = append(t.currentConnections, &proxyConn{
		source: conn,
		target: targetConn,
	})

	Pipe(conn, targetConn)

	t.currentConnectionsCount.Done()
	t.logger.Info(
		"connection closed",
		zap.String("remoteAddr", conn.RemoteAddr().String()),
	)
}
