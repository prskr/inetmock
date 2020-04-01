package main

import (
	"bytes"
	"fmt"
	"github.com/baez90/inetmock/internal/config"
	"github.com/baez90/inetmock/pkg/path"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"regexp"
	"sync"
)

const (
	name = "http_mock"
)

type httpHandler struct {
	logger *zap.Logger
	router *RegexpHandler
	server *http.Server
}

func (p *httpHandler) Run(config config.HandlerConfig) {
	options := loadFromConfig(config.Options())
	addr := fmt.Sprintf("%s:%d", config.ListenAddress(), config.Port())
	p.server = &http.Server{Addr: addr, Handler: p.router}
	p.logger = p.logger.With(
		zap.String("address", addr),
	)

	for _, rule := range options.Rules {
		p.setupRoute(rule)
	}

	go p.startServer()
}

func (p *httpHandler) Shutdown(wg *sync.WaitGroup) {
	if err := p.server.Close(); err != nil {
		p.logger.Error(
			"failed to shutdown HTTP server",
			zap.Error(err),
		)
	}

	wg.Done()
}

func (p *httpHandler) startServer() {
	if err := p.server.ListenAndServe(); err != nil {
		p.logger.Error(
			"failed to start http listener",
			zap.Error(err),
		)
	}
}

func (p *httpHandler) setupRoute(rule targetRule) {
	var compiled *regexp.Regexp
	var err error
	if compiled, err = regexp.Compile(rule.pattern); err != nil {
		p.logger.Warn(
			"failed to parse route - skipping",
			zap.String("route", rule.pattern),
			zap.Error(err),
		)
		return
	}
	p.logger.Info(
		"setup routing",
		zap.String("route", compiled.String()),
		zap.String("target", rule.target),
	)

	p.router.Handler(compiled, createHandlerForTarget(p.logger, rule.target))
}

func createHandlerForTarget(logger *zap.Logger, targetPath string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		targetFilePath := filepath.Join(path.WorkingDirectory(), targetPath)

		headerWriter := &bytes.Buffer{}
		request.Header.Write(headerWriter)

		logger.Info(
			"Handling request",
			zap.String("source", request.RemoteAddr),
			zap.String("host", request.Host),
			zap.String("method", request.Method),
			zap.String("protocol", request.Proto),
			zap.String("path", request.RequestURI),
			zap.String("target", targetFilePath),
			zap.Reflect("headers", request.Header),
		)

		http.ServeFile(writer, request, targetFilePath)
	})
}
