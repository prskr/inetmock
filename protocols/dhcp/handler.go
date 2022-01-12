package dhcp

import (
	"context"
	"errors"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/state"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	name               = "dhcp_mock"
	handlerNameLblName = "handler_name"
)

type dhcpHandler struct {
	logger     logging.Logger
	emitter    audit.Emitter
	stateStore state.KVStore
	server     *Server4
}

func (h *dhcpHandler) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	var (
		options ProtocolOptions
		conn    *ipv4.PacketConn
	)
	h.logger = h.logger.With(zap.String("protocol_handler", name))
	if o, err := LoadFromConfig(lifecycle, h.stateStore); err != nil {
		return err
	} else {
		options = o
	}

	if c, err := setupPacketConn(lifecycle.Uplink()); err != nil {
		return err
	} else {
		conn = c
	}

	h.logger = h.logger.With(zap.String("address", lifecycle.Uplink().Addr.String()))
	rh := &RuledHandler{
		HandlerName:     lifecycle.Name(),
		ProtocolOptions: options,
		Logger:          h.logger,
		StateStore:      h.stateStore.WithSuffixes(lifecycle.Name()),
	}

	for idx := range options.Rules {
		rule := options.Rules[idx]
		if err := rh.RegisterRule(rule); err != nil {
			h.logger.Error("Failed to setup rule", zap.String("raw_rule", rule), zap.Error(err))
			return err
		}
	}

	h.server = &Server4{
		Handler: &EmittingHandler{
			Upstream: &FallbackHandler{
				Previous:       rh,
				Logger:         h.logger,
				DefaultOptions: options.Default,
			},
			Emitter: h.emitter,
		},
		Logger: h.logger,
	}

	go h.serve(ctx, conn)

	return nil
}

func (h *dhcpHandler) serve(ctx context.Context, conn *ipv4.PacketConn) {
	if err := h.server.Serve(ctx, conn); err != nil {
		h.logger.Error("Failed to serve")
	}
}

func setupPacketConn(ul endpoint.Uplink) (*ipv4.PacketConn, error) {
	var socketAddr *net.UDPAddr
	if a, ok := ul.Addr.(*net.UDPAddr); ok {
		socketAddr = a
	} else {
		return nil, errors.New("uplink address not an UPD address")
	}

	if updConn, err := server4.NewIPv4UDPConn("", socketAddr); err != nil {
		return nil, err
	} else {
		return ipv4.NewPacketConn(updConn), nil
	}
}
