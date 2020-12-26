package http_mock

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

type route struct {
	rule    targetRule
	handler http.Handler
}

type RegexpHandler struct {
	handlerName string
	logger      logging.Logger
	routes      []*route
}

func (h *RegexpHandler) Handler(rule targetRule, handler http.Handler) {
	h.routes = append(h.routes, &route{rule, handler})
}

func (h *RegexpHandler) HandleFunc(rule targetRule, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{rule, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(h.handlerName))
	defer timer.ObserveDuration()

	for idx := range h.routes {
		rule := h.routes[idx].rule
		if h.routes[idx].rule.requestMatchTarget.Matches(r, rule.targetKey, rule.pattern) {
			totalRequestCounter.WithLabelValues(h.handlerName, strconv.FormatBool(true)).Inc()
			h.routes[idx].handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	totalRequestCounter.WithLabelValues(h.handlerName, strconv.FormatBool(false)).Inc()
	http.NotFound(w, r)
}

func (h *RegexpHandler) setupRoute(rule targetRule) {
	h.logger.Info(
		"setup routing",
		zap.String("route", rule.Pattern().String()),
		zap.String("response", rule.Response()),
	)

	h.Handler(rule, createHandlerForTarget(h.logger, rule.response))
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
