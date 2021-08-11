package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

const (
	expectedModuleName = "dns"
)

var (
	ErrNotADNSPInitiator = errors.New("the given initiator is not a DNS initiator")

	knownInitiators = map[string]func(args ...rules.Param) (Initiator, error){
		"a":    AorAAAAInitiator,
		"aaaa": AorAAAAInitiator,
	}
)

type Response struct {
	Hosts     []string
	Addresses []net.IP
}

type Initiator interface {
	Do(ctx context.Context, resolver *net.Resolver) (*Response, error)
}

type InitiatorFunc func(ctx context.Context, resolver *net.Resolver) (*Response, error)

func (f InitiatorFunc) Do(ctx context.Context, resolver *net.Resolver) (*Response, error) {
	return f(ctx, resolver)
}

func CheckForRule(rule *rules.Check) (Initiator, error) {
	var initiator = rule.Initiator
	if initiator == nil {
		return nil, rules.ErrNoInitiatorDefined
	}

	if !strings.EqualFold(strings.ToLower(initiator.Module), expectedModuleName) {
		return nil, ErrNotADNSPInitiator
	}

	if constructor, ok := knownInitiators[strings.ToLower(initiator.Name)]; !ok {
		return nil, fmt.Errorf("%w %s", rules.ErrUnknownInitiator, initiator.Name)
	} else {
		return constructor(initiator.Params...)
	}
}

func AorAAAAInitiator(args ...rules.Param) (Initiator, error) {
	var host, err = args[1].AsString()
	if err != nil {
		return nil, err
	}
	return InitiatorFunc(func(ctx context.Context, resolver *net.Resolver) (*Response, error) {
		var addrs, err = resolver.LookupHost(ctx, host)
		if err != nil {
			return nil, err
		}

		var ipAddrs []net.IP
		for idx := range addrs {
			if parsed := net.ParseIP(addrs[idx]); parsed != nil {
				ipAddrs = append(ipAddrs, parsed)
			}
		}

		return &Response{
			Addresses: ipAddrs,
		}, nil
	}), nil
}
