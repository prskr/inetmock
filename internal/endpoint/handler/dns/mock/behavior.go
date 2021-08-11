package mock

import (
	"fmt"
	"net"
	"strings"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
	"gitlab.com/inetmock/inetmock/internal/rules"
)

type HandlerForArgs interface {
	CreateHandler(args ...rules.Param) (dns.IPResolver, error)
}

type HandlerForArgsFunc func(args ...rules.Param) (dns.IPResolver, error)

func (h HandlerForArgsFunc) CreateHandler(args ...rules.Param) (dns.IPResolver, error) {
	return h(args...)
}

var (
	knownResponseHandlers = map[string]HandlerForArgs{
		"ip":          HandlerForArgsFunc(IPHandlerForArgs),
		"incremental": HandlerForArgsFunc(IncrementalHandlerForArgs),
		"random":      HandlerForArgsFunc(RandomHandlerForArgs),
	}
)

func HandlerForRoutingRule(rule *rules.Routing) (dns.IPResolver, error) {
	if rule.Terminator == nil {
		return nil, rules.ErrNoTerminatorDefined
	}

	if handlerForArgs, ok := knownResponseHandlers[strings.ToLower(rule.Terminator.Name)]; !ok {
		return nil, fmt.Errorf("%w %s", rules.ErrUnknownTerminator, rule.Terminator.Name)
	} else {
		return handlerForArgs.CreateHandler(rule.Terminator.Params...)
	}
}

func IPHandlerForArgs(args ...rules.Param) (dns.IPResolver, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var resolver dns.IPResolver
	if ip, err := args[0].AsIP(); err != nil {
		return nil, err
	} else {
		resolver = dns.IPResolverFunc(func(string) net.IP {
			return ip
		})
	}

	return resolver, nil
}

func IncrementalHandlerForArgs(args ...rules.Param) (dns.IPResolver, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var resolver dns.IPResolver
	if cidr, err := args[0].AsCIDR(); err != nil {
		return nil, err
	} else {
		resolver = NewIncrementalIPResolver(cidr.IPNet)
	}

	return resolver, nil
}

func RandomHandlerForArgs(args ...rules.Param) (dns.IPResolver, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var resolver dns.IPResolver
	if cidr, err := args[0].AsCIDR(); err != nil {
		return nil, err
	} else {
		resolver = NewRandomIPResolver(cidr.IPNet)
	}

	return resolver, nil
}
