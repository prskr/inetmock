package dns

import (
	"errors"
	"fmt"
	"strings"

	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

var (
	ErrUnknownCheckFilter     = errors.New("no check filter with the given name is known")
	ErrResponseNil            = errors.New("response must not be nil")
	ErrResponseEmpty          = errors.New("neither hosts nor addresses are set in the response")
	ErrUnmatchedResolvedHosts = errors.New("resolved hosts do not match")
	ErrUnmatchedResolvedIPs   = errors.New("resolved IPs do not match")

	knownCheckFilters = map[string]func(args ...rules.Param) (Validator, error){
		"notempty":     NotEmtpyResponseFilter,
		"resolvedhost": ResolvedHostResponseFilter,
		"resolvedip":   ResolvedIPResponseFilter,
	}
)

type Validator interface {
	Matches(resp *Response) error
}

type ValidationChain []Validator

func (c *ValidationChain) Add(v Validator) {
	var arr = *c
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

	for idx := range rule.Validators.Chain {
		var validator = rule.Validators.Chain[idx]
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

func NotEmtpyResponseFilter(...rules.Param) (Validator, error) {
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
	var expectedHost, err = args[0].AsString()
	if err != nil {
		return nil, err
	}

	if !mdns.IsFqdn(expectedHost) {
		expectedHost = mdns.Fqdn(expectedHost)
	}

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

	var expectedIP, err = args[0].AsIP()
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
