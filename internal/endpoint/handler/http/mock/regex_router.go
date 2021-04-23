package mock

import (
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	imHttp "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type route struct {
	rule    TargetRule
	handler http.Handler
}

type RegexHandler struct {
	handlerName string
	logger      logging.Logger
	routes      []*route
	emitter     audit.Emitter
}

func NewRegexHandler(name string, logger logging.Logger, emitter audit.Emitter) *RegexHandler {
	return &RegexHandler{
		handlerName: name,
		logger:      logger,
		emitter:     emitter,
	}
}

func (h *RegexHandler) Handler(rule TargetRule, handler http.Handler) {
	h.routes = append(h.routes, &route{rule, handler})
}

func (h *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *RegexHandler) AddRouteRule(rule TargetRule) {
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
	f.emitter.Emit(imHttp.EventFromRequest(request, v1.AppProtocol_APP_PROTOCOL_HTTP))
	file, err := os.Open(f.targetPath)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}
	defer func() {
		_ = file.Close()
	}()
	//nolint:gosec
	http.ServeContent(writer, request, path.Base(request.RequestURI), time.Now().Add(-(time.Duration(rand.Int()) * time.Millisecond)), file)
}
