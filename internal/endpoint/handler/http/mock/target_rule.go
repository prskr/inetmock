package mock

import "regexp"

func NewPathTargetRule(pattern, response string) (TargetRule, error) {
	regexpPattern, err := regexp.Compile(pattern)
	if err != nil {
		return TargetRule{}, err
	}
	return TargetRule{
		pattern:            regexpPattern,
		response:           response,
		requestMatchTarget: RequestMatchTargetPath,
	}, nil
}

func NewHeaderTargetRule(headerKey, pattern, response string) (TargetRule, error) {
	regexpPattern, err := regexp.Compile(pattern)
	if err != nil {
		return TargetRule{}, err
	}

	return TargetRule{
		pattern:            regexpPattern,
		response:           response,
		requestMatchTarget: RequestMatchTargetHeader,
		targetKey:          headerKey,
	}, nil
}

func MustPathTargetRule(pattern, response string) TargetRule {
	rule, err := NewPathTargetRule(pattern, response)
	if err != nil {
		panic(err)
	}
	return rule
}

func MustHeaderTargetRule(headerKey, pattern, response string) TargetRule {
	rule, err := NewHeaderTargetRule(headerKey, pattern, response)
	if err != nil {
		panic(err)
	}
	return rule
}

type TargetRule struct {
	pattern            *regexp.Regexp
	response           string
	requestMatchTarget RequestMatchTarget
	targetKey          string
}

func (tr TargetRule) Pattern() *regexp.Regexp {
	return tr.pattern
}

func (tr TargetRule) Response() string {
	return tr.response
}
