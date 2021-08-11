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
	ErrNotADNSPInitiator          = errors.New("the given initiator is not a DNS initiator")
	DefaultResolver      Resolver = new(net.Resolver)

	knownInitiators = map[string]func(args ...rules.Param) (Initiator, error){
		"a":    AorAAAAInitiator,
		"aaaa": AorAAAAInitiator,
		"ptr":  PTRInitiator,
	}
)

type Resolver interface {
	// LookupHost looks up the given host using the local resolver.
	// It returns a slice of that host's addresses.
	LookupHost(ctx context.Context, host string) (addrs []string, err error)

	// LookupAddr performs a reverse lookup for the given address, returning a list
	// of names mapping to that address.
	//
	// The returned names are validated to be properly formatted presentation-format
	// domain names. If the response contains invalid names, those records are filtered
	// out and an error will be returned alongside the the remaining results, if any.
	LookupAddr(ctx context.Context, addr string) (names []string, err error)
}

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
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var host, err = args[0].AsString()
	if err != nil {
		return nil, err
	}
	return InitiatorFunc(func(ctx context.Context, resolver Resolver) (*Response, error) {
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

func PTRInitiator(args ...rules.Param) (Initiator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var ip, err = args[0].AsIP()
	if err != nil {
		return nil, err
	}

	return InitiatorFunc(func(ctx context.Context, resolver Resolver) (*Response, error) {
		var names, err = resolver.LookupAddr(ctx, ip.String())
		if err != nil {
			return nil, err
		}
		return &Response{
			Hosts: names,
		}, nil
	}), nil
}
