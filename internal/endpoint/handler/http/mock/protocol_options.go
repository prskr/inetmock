//go:generate go-enum -f $GOFILE --lower --marshal --names
package mock

import (
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/spf13/viper"
)

var (
	ruleValueSelectors = map[RequestMatchTarget]ruleValueSelector{
		RequestMatchTargetHeader: func(req *http.Request, targetKey string) string {
			return req.Header.Get(targetKey)
		},
		RequestMatchTargetPath: func(req *http.Request, _ string) string {
			return req.URL.Path
		},
	}
)

/* ENUM(
Path,
Header
)
*/
type RequestMatchTarget int

func (x RequestMatchTarget) Matches(req *http.Request, targetKey string, regex *regexp.Regexp) bool {
	val := ruleValueSelectors[x](req, targetKey)
	return regex.MatchString(val)
}

type ruleValueSelector func(req *http.Request, targetKey string) string

type targetRule struct {
	pattern            *regexp.Regexp
	response           string
	requestMatchTarget RequestMatchTarget
	targetKey          string
}

func (tr targetRule) Pattern() *regexp.Regexp {
	return tr.pattern
}

func (tr targetRule) Response() string {
	return tr.response
}

type httpOptions struct {
	TLS   bool
	Rules []targetRule
}

func loadFromConfig(config *viper.Viper) (options httpOptions, err error) {
	type tmpCfg struct {
		Pattern  string
		Response string
		Matcher  string
		Target   string
	}

	tmpRules := struct {
		TLS   bool
		Rules []tmpCfg
	}{}

	if err = config.Unmarshal(&tmpRules); err != nil {
		return
	}

	options.TLS = tmpRules.TLS

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

		options.Rules = append(options.Rules, targetRule{
			pattern:            rulePattern,
			response:           absoluteResponsePath,
			requestMatchTarget: matchTargetValue,
			targetKey:          i.Target,
		})
	}

	return
}
