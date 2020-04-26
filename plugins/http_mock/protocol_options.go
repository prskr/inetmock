package main

import (
	"github.com/spf13/viper"
	"regexp"
)

const (
	rulesConfigKey    = "rules"
	patternConfigKey  = "pattern"
	responseConfigKey = "response"
)

type targetRule struct {
	pattern  *regexp.Regexp
	response string
}

func (tr targetRule) Pattern() *regexp.Regexp {
	return tr.pattern
}

func (tr targetRule) Response() string {
	return tr.response
}

type httpOptions struct {
	Rules []targetRule
}

func loadFromConfig(config *viper.Viper) (options httpOptions) {
	anonRules := config.Get(rulesConfigKey).([]interface{})

	for _, i := range anonRules {
		innerData := i.(map[interface{}]interface{})

		if rulePattern, err := regexp.Compile(innerData[patternConfigKey].(string)); err == nil {
			options.Rules = append(options.Rules, targetRule{
				pattern:  rulePattern,
				response: innerData[responseConfigKey].(string),
			})
		} else {
			panic(err)
		}
	}

	return
}
