//go:generate go-enum -f $GOFILE --lower --marshal --names

package mock

import (
	"net/http"
	"regexp"
)

/* ENUM(
Path,
Header
)
*/
type RequestMatchTarget int

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

func (x RequestMatchTarget) Matches(req *http.Request, targetKey string, regex *regexp.Regexp) bool {
	val := ruleValueSelectors[x](req, targetKey)
	return regex.MatchString(val)
}
