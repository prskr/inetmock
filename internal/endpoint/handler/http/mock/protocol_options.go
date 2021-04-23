package mock

import (
	"net/http"
	"path/filepath"
	"regexp"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

type ruleValueSelector func(req *http.Request, targetKey string) string

type httpOptions struct {
	Rules []TargetRule
}

func loadFromConfig(lifecycle endpoint.Lifecycle) (httpOptions, error) {
	type tmpCfg struct {
		Pattern  string
		Response string
		Matcher  string
		Target   string
	}

	tmpRules := struct {
		Rules []tmpCfg
	}{}

	if err := lifecycle.UnmarshalOptions(&tmpRules); err != nil {
		return httpOptions{}, err
	}

	var options httpOptions

	for _, i := range tmpRules.Rules {
		var rulePattern *regexp.Regexp
		var matchTargetValue RequestMatchTarget
		var absoluteResponsePath string
		var parseErr error
		if rulePattern, parseErr = regexp.Compile(i.Pattern); parseErr != nil {
			continue
		}
		if matchTargetValue, parseErr = ParseRequestMatchTarget(i.Matcher); parseErr != nil {
			matchTargetValue = RequestMatchTargetPath
		}

		if absoluteResponsePath, parseErr = filepath.Abs(i.Response); parseErr != nil {
			continue
		}

		options.Rules = append(options.Rules, TargetRule{
			pattern:            rulePattern,
			response:           absoluteResponsePath,
			requestMatchTarget: matchTargetValue,
			targetKey:          i.Target,
		})
	}

	return options, nil
}
