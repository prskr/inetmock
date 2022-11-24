package mock

import (
	"io/fs"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type (
	RequestFilter interface {
		Matches(req *http.Request) bool
	}

	FilterChain []RequestFilter

	ConditionalHandler struct {
		http.Handler
		Chain FilterChain
	}
)

func (c FilterChain) Matches(req *http.Request) bool {
	for idx := range c {
		if !c[idx].Matches(req) {
			return false
		}
	}
	return true
}

type Router struct {
	HandlerName string
	Logger      logging.Logger
	FakeFileFS  fs.FS
	handlers    []ConditionalHandler
}

func (r *Router) RegisterRule(rawRule string) error {
	r.Logger.Debug("Adding routing rule", zap.String("rawRule", rawRule))

	var (
		rule *rules.SingleResponsePipeline
		err  error
	)

	if rule, err = rules.Parse[rules.SingleResponsePipeline](rawRule); err != nil {
		return err
	}

	var conditionalHandler ConditionalHandler

	if conditionalHandler.Chain, err = RequestFiltersForRoutingRule(rule); err != nil {
		return err
	}

	if conditionalHandler.Handler, err = HandlerForRoutingRule(rule, r.Logger, r.FakeFileFS); err != nil {
		return err
	}

	r.Logger.Debug("Configure successfully parsed routing rule")
	r.handlers = append(r.handlers, conditionalHandler)

	return nil
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(r.HandlerName))
	defer timer.ObserveDuration()

	for idx := range r.handlers {
		if r.handlers[idx].Chain.Matches(request) {
			r.handlers[idx].ServeHTTP(writer, request)
			return
		}
	}

	writer.WriteHeader(http.StatusNotFound)
}
