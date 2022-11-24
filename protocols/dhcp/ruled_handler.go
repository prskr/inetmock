package dhcp

import (
	"errors"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
	"inetmock.icb4dc0.de/inetmock/internal/state"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	"inetmock.icb4dc0.de/inetmock/protocols"
)

type (
	DHCPv4MessageHandler interface {
		Handle(req, resp *dhcpv4.DHCPv4) error
	}
	RequestFilter interface {
		Matches(msg *dhcpv4.DHCPv4) bool
	}
	FilterChain        []RequestFilter
	HandlerChain       []DHCPv4MessageHandler
	ConditionalHandler struct {
		Handlers HandlerChain
		Chain    FilterChain
	}
	DHCPv4MessageHandlerFunc func(req, resp *dhcpv4.DHCPv4) error
)

var NoOpHandler DHCPv4MessageHandler = DHCPv4MessageHandlerFunc(func(_, _ *dhcpv4.DHCPv4) error {
	return nil
})

func (c FilterChain) Matches(m *dhcpv4.DHCPv4) bool {
	for idx := range c {
		if !c[idx].Matches(m) {
			return false
		}
	}
	return true
}

func (c HandlerChain) Apply(req, resp *dhcpv4.DHCPv4) error {
	for idx := range c {
		if err := c[idx].Handle(req, resp); err != nil {
			return err
		}
	}
	return nil
}

func (f DHCPv4MessageHandlerFunc) Handle(req, resp *dhcpv4.DHCPv4) error {
	return f(req, resp)
}

type RuledHandler struct {
	HandlerName     string
	ProtocolOptions ProtocolOptions
	Logger          logging.Logger
	StateStore      state.KVStore
	handlers        []ConditionalHandler
}

func (h *RuledHandler) RegisterRule(rawRule string) (err error) {
	h.Logger.Debug("Adding routing rule", zap.String("rawRule", rawRule))
	var rule *rules.ChainedResponsePipeline
	if rule, err = rules.Parse[rules.ChainedResponsePipeline](rawRule); err != nil {
		return err
	}

	var conditionalHandler ConditionalHandler

	if conditionalHandler.Chain, err = RequestFiltersForRoutingRule(rule); err != nil {
		return err
	}

	handlerOptions := HandlerOptions{
		Logger:          h.Logger,
		StateStore:      h.StateStore,
		ProtocolOptions: h.ProtocolOptions,
	}
	if conditionalHandler.Handlers, err = HandlerForRoutingRule(rule, handlerOptions); err != nil {
		return err
	}

	h.Logger.Debug("Configure successfully parsed routing rule")
	h.handlers = append(h.handlers, conditionalHandler)

	return nil
}

func (h *RuledHandler) Handle(req, resp *dhcpv4.DHCPv4) error {
	defer prometheus.NewTimer(protocols.RequestDurationHistogram.WithLabelValues("dhcp", h.HandlerName)).ObserveDuration()

	for idx := range h.handlers {
		handler := h.handlers[idx]
		if handler.Chain.Matches(req) {
			if err := handler.Handlers.Apply(req, resp); err != nil {
				return err
			}
			return nil
		}
	}

	if h.ProtocolOptions.Fallback != nil {
		h.Logger.Info("Resolving request with default handler")
		return h.ProtocolOptions.Fallback.Handle(req, resp)
	}

	return errors.New("no matching handler")
}
