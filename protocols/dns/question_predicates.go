package dns

import (
	"fmt"
	"regexp"
	"strings"

	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

var knownRequestFilters = map[string]func(args ...rules.Param) (QuestionPredicate, error){
	"a":    HostnameQuestionFilter(mdns.TypeA),
	"aaaa": HostnameQuestionFilter(mdns.TypeAAAA),
}

func QuestionPredicatesForRoutingRule(rule *rules.SingleResponsePipeline) (predicates []QuestionPredicate, err error) {
	if rule == nil || rule.FilterChain == nil || len(rule.FilterChain.Chain) == 0 {
		return nil, nil
	}

	predicates = make([]QuestionPredicate, 0, len(rule.FilterChain.Chain))
	for idx := range rule.FilterChain.Chain {
		if constructor, ok := knownRequestFilters[strings.ToLower(rule.FilterChain.Chain[idx].Name)]; !ok {
			return nil, fmt.Errorf("%w %s", rules.ErrUnknownFilterMethod, rule.FilterChain.Chain[idx].Name)
		} else {
			var instance QuestionPredicate
			instance, err = constructor(rule.FilterChain.Chain[idx].Params...)
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
