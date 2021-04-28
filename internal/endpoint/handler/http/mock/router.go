package mock

import (
	"io/fs"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	imHttp "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http"
	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type ConditionalHandler struct {
	http.Handler
	Filters []RequestFilter
}

func (c *ConditionalHandler) Matches(req *http.Request) bool {
	for idx := range c.Filters {
		if !c.Filters[idx].Matches(req) {
			return false
		}
	}
	return true
}

type Router struct {
	HandlerName string
	Logger      logging.Logger
	Emitter     audit.Emitter
	FakeFileFS  fs.FS
	handlers    []ConditionalHandler
}

func (r *Router) RegisterRule(rawRule string) error {
	r.Logger.Info("Adding routing rule", zap.String("rawRule", rawRule))
	var err error
	var rule *rules.Routing
	if rule, err = rules.Parse(rawRule); err != nil {
		return err
	}

	var filters []RequestFilter
	if filters, err = RequestFiltersForRoutingRule(rule); err != nil {
		return err
	}

	var handler http.Handler
	if handler, err = HandlerForRoutingRule(rule, r.Logger, r.FakeFileFS); err != nil {
		return err
	}

	r.Logger.Info("Configure successfully parsed routing rule")

	r.handlers = append(r.handlers, ConditionalHandler{
		Handler: handler,
		Filters: filters,
	})

	return nil
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(r.HandlerName))
	defer timer.ObserveDuration()

	r.Emitter.Emit(imHttp.EventFromRequest(request, v1.AppProtocol_APP_PROTOCOL_HTTP))

	for idx := range r.handlers {
		if r.handlers[idx].Matches(request) {
			r.handlers[idx].ServeHTTP(writer, request)
			return
		}
	}

	writer.WriteHeader(http.StatusNotFound)
}
