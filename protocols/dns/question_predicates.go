package dns

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

var knownRequestFilters = map[string]func(args ...rules.Param) (QuestionPredicate, error){
	"a":    HostnameQuestionFilter(dns.TypeA),
	"aaaa": HostnameQuestionFilter(dns.TypeAAAA),
}

func QuestionPredicatesForRoutingRule(rule *rules.Routing) (predicates []QuestionPredicate, err error) {
	if rule == nil || rule.Filters == nil || len(rule.Filters.Chain) == 0 {
		return nil, nil
	}

	predicates = make([]QuestionPredicate, 0, len(rule.Filters.Chain))
	for idx := range rule.Filters.Chain {
		if constructor, ok := knownRequestFilters[strings.ToLower(rule.Filters.Chain[idx].Name)]; !ok {
			return nil, fmt.Errorf("%w %s", rules.ErrUnknownFilterMethod, rule.Filters.Chain[idx].Name)
		} else {
			var instance QuestionPredicate
			instance, err = constructor(rule.Filters.Chain[idx].Params...)
			if err != nil {
				return
			}
			predicates = append(predicates, instance)
		}
	}

	return
}

func HostnameQuestionFilter(qType uint16) func(args ...rules.Param) (QuestionPredicate, error) {
	return func(args ...rules.Param) (QuestionPredicate, error) {
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

		return QuestionPredicateFunc(func(req Question) bool {
			if req.Qtype == qType {
				return pattern.MatchString(req.Name)
			}
			return false
		}), nil
	}
}
