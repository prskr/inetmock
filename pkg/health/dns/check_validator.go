package dns

import (
	"errors"
	"fmt"
	"strings"

	mdns "github.com/miekg/dns"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
)

var (
	ErrResponseNil            = errors.New("response must not be nil")
	ErrResponseEmpty          = errors.New("neither hosts nor addresses are set in the response")
	ErrUnmatchedResolvedHosts = errors.New("resolved hosts do not match")
	ErrUnmatchedResolvedIPs   = errors.New("resolved IPs do not match")

	knownCheckFilters = map[string]func(args ...rules.Param) (Validator, error){
		"notempty":     NotEmptyResponseFilter,
		"resolvedhost": ResolvedHostResponseFilter,
		"resolvedip":   ResolvedIPResponseFilter,
		"incidr":       InCIDRResponseFilter,
	}
)

type Validator interface {
	Matches(resp *Response) error
}

type ValidationChain []Validator

func (c *ValidationChain) Add(v Validator) {
	arr := *c
	arr = append(arr, v)
	*c = arr
}

func (c ValidationChain) Len() int {
	return len([]Validator(c))
}

func (c ValidationChain) Matches(resp *Response) error {
	for idx := range c {
		if err := c[idx].Matches(resp); err != nil {
			return err
		}
	}
	return nil
}

type CheckFilterFunc func(resp *Response) error

func (f CheckFilterFunc) Matches(resp *Response) error {
	return f(resp)
}

func ValidatorsForRule(rule *rules.Check) (filters ValidationChain, err error) {
	if rule.Validators == nil {
		return nil, nil
	}

	filters = make(ValidationChain, 0, len(rule.Validators.Chain))

	for idx := range rule.Validators.Chain {
		validator := rule.Validators.Chain[idx]
		if provider, ok := knownCheckFilters[strings.ToLower(validator.Name)]; !ok {
			return nil, fmt.Errorf("%w: %s", rules.ErrUnknownFilterMethod, validator.Name)
		} else if instance, err := provider(validator.Params...); err != nil {
			return nil, err
		} else {
			filters.Add(instance)
		}
	}

	return
}

func NotEmptyResponseFilter(...rules.Param) (Validator, error) {
	return CheckFilterFunc(func(resp *Response) error {
		switch {
		case resp == nil:
			return ErrResponseNil
		case len(resp.Addresses) == 0 && len(resp.Hosts) == 0:
			return ErrResponseEmpty
		default:
			return nil
		}
	}), nil
}

func ResolvedHostResponseFilter(args ...rules.Param) (Validator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}
	expectedHost, err := args[0].AsString()
	if err != nil {
		return nil, err
	}

	expectedHost = mdns.Fqdn(expectedHost)

	return CheckFilterFunc(func(resp *Response) error {
		switch {
		case resp == nil:
			return ErrResponseNil
		case len(resp.Hosts) == 0:
			return ErrResponseEmpty
		}

		for idx := range resp.Hosts {
			if strings.EqualFold(expectedHost, resp.Hosts[idx]) {
				return nil
			}
		}

		return fmt.Errorf("%w: %s", ErrUnmatchedResolvedHosts, expectedHost)
	}), nil
}

func ResolvedIPResponseFilter(args ...rules.Param) (Validator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	expectedIP, err := args[0].AsIP()
	if err != nil {
		return nil, err
	}

	return CheckFilterFunc(func(resp *Response) error {
		switch {
		case resp == nil:
			return ErrResponseNil
		case len(resp.Addresses) == 0:
			return ErrResponseEmpty
		}

		for idx := range resp.Addresses {
			if resp.Addresses[idx].Equal(expectedIP) {
				return nil
			}
		}
		return fmt.Errorf("%w: %s", ErrUnmatchedResolvedIPs, expectedIP.String())
	}), nil
}

func InCIDRResponseFilter(args ...rules.Param) (Validator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	cidr, err := args[0].AsCIDR()
	if err != nil {
		return nil, err
	}

	return CheckFilterFunc(func(resp *Response) error {
		switch {
		case resp == nil:
			return ErrResponseNil
		case len(resp.Addresses) == 0:
			return ErrResponseEmpty
		}

		for idx := range resp.Addresses {
			if cidr.IPNet.Contains(resp.Addresses[idx]) {
				return nil
			}
		}

		return fmt.Errorf("%w: %s", ErrUnmatchedResolvedIPs, cidr.IPNet.String())
	}), nil
}
