package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

var (
	ErrNotADNSPInitiator = errors.New("the given initiator is not a DNS initiator")

	knownInitiators = map[string]func(logger logging.Logger, args ...rules.Param) (Initiator, error){
		"a":    AorAAAAInitiator,
		"aaaa": AorAAAAInitiator,
		"ptr":  PTRInitiator,
	}
)

type Response struct {
	Hosts     []string
	Addresses []net.IP
}

type Initiator interface {
	Do(ctx context.Context, resolver Resolver) (*Response, error)
}

type InitiatorFunc func(ctx context.Context, resolver Resolver) (*Response, error)

func (f InitiatorFunc) Do(ctx context.Context, resolver Resolver) (*Response, error) {
	return f(ctx, resolver)
}

func CheckForRule(rule *rules.Check, logger logging.Logger) (Initiator, error) {
	initiator := rule.Initiator
	if initiator == nil {
		return nil, rules.ErrNoInitiatorDefined
	}

	switch strings.ToLower(initiator.Module) {
	case "dns", "dot", "doh", "doh2":
		if constructor, ok := knownInitiators[strings.ToLower(initiator.Name)]; !ok {
			return nil, fmt.Errorf("%w %s", rules.ErrUnknownInitiator, initiator.Name)
		} else {
			return constructor(logger, initiator.Params...)
		}
	default:
		return nil, ErrNotADNSPInitiator
	}
}

func AorAAAAInitiator(logger logging.Logger, args ...rules.Param) (Initiator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	host, err := args[0].AsString()
	if err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("request_target", host),
	)

	return InitiatorFunc(func(ctx context.Context, resolver Resolver) (*Response, error) {
		logger.Debug("Initiating check")
		if addrs, err := resolver.LookupA(ctx, host); err != nil {
			return nil, err
		} else {
			return &Response{
				Addresses: addrs,
			}, nil
		}
	}), nil
}

func PTRInitiator(logger logging.Logger, args ...rules.Param) (Initiator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var (
		ip  net.IP
		err error
	)

	if ip, err = args[0].AsIP(); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("request_target", ip.String()),
	)

	return InitiatorFunc(func(ctx context.Context, resolver Resolver) (*Response, error) {
		logger.Debug("Setup health initiator")
		names, err := resolver.LookupPTR(ctx, ip.String())
		if err != nil {
			return nil, err
		}
		return &Response{
			Hosts: names,
		}, nil
	}), nil
}
