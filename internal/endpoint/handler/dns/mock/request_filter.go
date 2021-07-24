package mock

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

var (
	ErrUnknownFilterMethod = errors.New("no filter with the given name is known")

	knownRequestFilters = map[string]func(args ...rules.Param) (RequestFilter, error){
		"a":    HostnameQuestionFilter(dns.TypeA),
		"aaaa": HostnameQuestionFilter(dns.TypeAAAA),
	}
)

type RequestFilter interface {
	Matches(req *dns.Question) bool
}

type RequestFilterFunc func(req *dns.Question) bool

func (r RequestFilterFunc) Matches(req *dns.Question) bool {
	return r(req)
}

func RequestFiltersForRoutingRule(rule *rules.Routing) (filters []RequestFilter, err error) {
	if rule == nil || rule.Filters == nil || len(rule.Filters.Chain) == 0 {
		return nil, nil
	}
	filters = make([]RequestFilter, len(rule.Filters.Chain))
	for idx := range rule.Filters.Chain {
		if constructor, ok := knownRequestFilters[strings.ToLower(rule.Filters.Chain[idx].Name)]; !ok {
			return nil, fmt.Errorf("%w %s", ErrUnknownFilterMethod, rule.Filters.Chain[idx].Name)
		} else {
			var instance RequestFilter
			instance, err = constructor(rule.Filters.Chain[idx].Params...)
			if err != nil {
				return
			}
			filters[idx] = instance
		}
	}
	return
}

func HostnameQuestionFilter(qType uint16) func(args ...rules.Param) (RequestFilter, error) {
	return func(args ...rules.Param) (RequestFilter, error) {
		if err := rules.ValidateParameterCount(args, 1); err != nil {
			return nil, err
		}

		var (
			rawPattern string
			pattern    *regexp.Regexp
			err        error
		)

		if rawPattern, err = args[0].AsString(); err != nil {
			return nil, err
		}
		if pattern, err = regexp.Compile(rawPattern); err != nil {
			return nil, err
		}

		return RequestFilterFunc(func(req *dns.Question) bool {
			// if nil there's nothing to match
			if req == nil {
				return false
			}

			if req.Qtype == qType {
				return pattern.MatchString(req.Name)
			}
			return false
		}), nil
	}
}
