package main

import (
	"github.com/spf13/viper"
	"net"
	"regexp"
)

const (
	rulesConfigKey            = "rules"
	patternConfigKey          = "pattern"
	responseConfigKey         = "response"
	fallbackStrategyConfigKey = "fallback.strategy"
	fallbackArgsConfigKey     = "fallback.args"
)

type resolverRule struct {
	pattern  *regexp.Regexp
	response net.IP
}

type dnsOptions struct {
	Rules    []resolverRule
	Fallback ResolverFallback
}

func loadFromConfig(config *viper.Viper) dnsOptions {
	options := dnsOptions{}

	anonRules := config.Get(rulesConfigKey).([]interface{})
	for _, rule := range anonRules {
		innerData := rule.(map[interface{}]interface{})
		var err error
		var compiledPattern *regexp.Regexp
		var response net.IP
		if compiledPattern, err = regexp.Compile(innerData[patternConfigKey].(string)); err != nil {
			continue
		}

		if response = net.ParseIP(innerData[responseConfigKey].(string)); response == nil {
			continue
		}

		options.Rules = append(options.Rules, resolverRule{
			pattern:  compiledPattern,
			response: response,
		})
	}

	options.Fallback = CreateResolverFallback(
		config.GetString(fallbackStrategyConfigKey),
		config.Sub(fallbackArgsConfigKey),
	)

	return options
}
