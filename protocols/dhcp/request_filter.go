package dhcp

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
)

var knownRequestFilters = map[string]func(args ...rules.Param) (RequestFilter, error){
	"matchmac": MatchMACMatcher,
	"exactmac": ExactMACMatcher,
}

type (
	RequestFilterFunc func(msg *dhcpv4.DHCPv4) bool
)

func (f RequestFilterFunc) Matches(msg *dhcpv4.DHCPv4) bool {
	return f(msg)
}

func RequestFiltersForRoutingRule(rule rules.FilteredPipeline) (filters FilterChain, err error) {
	chain := rule.Filters()
	if len(chain) == 0 {
		return nil, nil
	}
	filters = make([]RequestFilter, len(chain))
	for idx := range chain {
		if constructor, ok := knownRequestFilters[strings.ToLower(chain[idx].Name)]; !ok {
			return nil, fmt.Errorf("%w %s", rules.ErrUnknownFilterMethod, chain[idx].Name)
		} else {
			var instance RequestFilter
			instance, err = constructor(chain[idx].Params...)
			if err != nil {
				return
			}
			filters[idx] = instance
		}
	}
	return
}

func MatchMACMatcher(args ...rules.Param) (RequestFilter, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var macMatchRegexp *regexp.Regexp

	if rawRegexp, err := args[0].AsString(); err != nil {
		return nil, err
	} else if exp, err := regexp.Compile(rawRegexp); err != nil {
		return nil, err
	} else {
		macMatchRegexp = exp
	}

	return RequestFilterFunc(func(msg *dhcpv4.DHCPv4) bool {
		return macMatchRegexp.MatchString(msg.ClientHWAddr.String())
	}), nil
}

func ExactMACMatcher(args ...rules.Param) (RequestFilter, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var hwAddr net.HardwareAddr

	if rawAddr, err := args[0].AsString(); err != nil {
		return nil, err
	} else if mac, err := net.ParseMAC(rawAddr); err != nil {
		return nil, err
	} else {
		hwAddr = mac
	}

	return RequestFilterFunc(func(msg *dhcpv4.DHCPv4) bool {
		return bytes.Equal(hwAddr, msg.ClientHWAddr)
	}), nil
}
