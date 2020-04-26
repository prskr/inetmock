package main

import (
	"bytes"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
	"net/http"
)

func (p *httpHandler) setupRoute(rule targetRule) {
	p.logger.Info(
		"setup routing",
		zap.String("route", rule.Pattern().String()),
		zap.String("response", rule.Response()),
	)

	p.router.Handler(rule.Pattern(), createHandlerForTarget(p.logger, rule.response))
}

func createHandlerForTarget(logger logging.Logger, targetPath string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		headerWriter := &bytes.Buffer{}
		request.Header.Write(headerWriter)

		logger.Info(
			"Handling request",
			zap.String("source", request.RemoteAddr),
			zap.String("host", request.Host),
			zap.String("method", request.Method),
			zap.String("protocol", request.Proto),
			zap.String("path", request.RequestURI),
			zap.String("response", targetPath),
			zap.Reflect("headers", request.Header),
		)

		http.ServeFile(writer, request, targetPath)
	})
}
