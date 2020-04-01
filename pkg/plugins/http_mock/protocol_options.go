package main

import "github.com/spf13/viper"

const (
	rulesConfigKey   = "rules"
	patternConfigKey = "pattern"
	targetConfigKey  = "target"
)

type targetRule struct {
	pattern string
	target  string
}

type httpOptions struct {
	Rules []targetRule
}

func loadFromConfig(config *viper.Viper) httpOptions {
	options := httpOptions{}
	anonRules := config.Get(rulesConfigKey).([]interface{})

	for _, i := range anonRules {
		innerData := i.(map[interface{}]interface{})
		options.Rules = append(options.Rules, targetRule{
			pattern: innerData[patternConfigKey].(string),
			target:  innerData[targetConfigKey].(string),
		})
	}

	return options
}
