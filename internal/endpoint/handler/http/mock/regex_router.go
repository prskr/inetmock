package mock

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	details "gitlab.com/inetmock/inetmock/pkg/audit/details"
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
	emitter     audit.Emitter
}

func (h *RegexpHandler) Handler(rule targetRule, handler http.Handler) {
	h.routes = append(h.routes, &route{rule, handler})
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

	h.Handler(rule, emittingFileHandler{
		emitter:    h.emitter,
		targetPath: rule.response,
	})
}

type emittingFileHandler struct {
	emitter    audit.Emitter
	targetPath string
}

func (f emittingFileHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	f.emitter.Emit(eventFromRequest(request))
	http.ServeFile(writer, request, f.targetPath)
}

func eventFromRequest(request *http.Request) audit.Event {
	httpDetails := details.HTTP{
		Method:  request.Method,
		Host:    request.Host,
		URI:     request.RequestURI,
		Proto:   request.Proto,
		Headers: request.Header,
	}

	ev := audit.Event{
		Transport:       audit.TransportProtocol_TCP,
		Application:     audit.AppProtocol_HTTP,
		ProtocolDetails: httpDetails,
	}

	ev.SetDestinationIPFromAddr(LocalAddr(request.Context()))
	ev.SetSourceIPFromAddr(RemoteAddr(request.Context()))

	return ev
}
