package dns

import (
	"fmt"
	"net"
	"strings"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

var knownRuleResolvers = map[string]func(params []rules.Param) (IPResolver, error){
	"ip":          StaticIPResolverForArgs,
	"incremental": IncrementalResolverForArgs,
	"random":      RandomIPResolverForArgs,
}

func ResolverForRule(rule *rules.Routing) (IPResolver, error) {
	if rule == nil || rule.Terminator == nil {
		return nil, rules.ErrNoTerminatorDefined
	}

	if resolverConstructor, ok := knownRuleResolvers[strings.ToLower(rule.Terminator.Name)]; !ok {
		return nil, fmt.Errorf("%w: %s", rules.ErrUnknownTerminator, rule.Terminator.Name)
	} else {
		return resolverConstructor(rule.Terminator.Params)
	}
}

func StaticIPResolverForArgs(args []rules.Param) (IPResolver, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	if ip, err := args[0].AsIP(); err != nil {
		return nil, err
	} else {
		return IPResolverFunc(func(string) net.IP {
			return ip
		}), nil
	}
}

func IncrementalResolverForArgs(args []rules.Param) (IPResolver, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	if cidr, err := args[0].AsCIDR(); err != nil {
		return nil, err
	} else {
		return NewIncrementalIPResolver(cidr.IPNet), nil
	}
}

func RandomIPResolverForArgs(args []rules.Param) (IPResolver, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	if cidr, err := args[0].AsCIDR(); err != nil {
		return nil, err
	} else {
		return NewRandomIPResolver(cidr.IPNet), nil
	}
}
