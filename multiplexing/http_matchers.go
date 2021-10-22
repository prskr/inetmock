package multiplexing

import (
	"io"
	"net/http"
	"strings"

	"github.com/soheilhy/cmux"
)

const (
	HTTPVersionUnknown HTTPVersion = iota
	HTTPVersion10
	HTTPVersion11
	HTTPVersion2
)

type (
	HTTPVersion    uint8
	RequestPreface struct {
		Version HTTPVersion
		Scheme  string
		Method  string
		Host    string
		Path    string
		Header  http.Header
	}
	RequestMatcher func(req *RequestPreface) bool
)

func HTTP() cmux.Matcher {
	return HTTPMatchAnd(func(req *RequestPreface) bool {
		if req == nil {
			return false
		}
		switch req.Version {
		case HTTPVersionUnknown:
			return false
		case HTTPVersion10, HTTPVersion11, HTTPVersion2:
			fallthrough
		default:
			return true
		}
	})
}

func PathPrefixMatcher(prefix string) RequestMatcher {
	return func(req *RequestPreface) bool {
		return strings.HasPrefix(req.Path, prefix)
	}
}

func HTTPMatchOr(reqMatchers ...RequestMatcher) cmux.Matcher {
	return httpMatch(reqMatchers, false, func(b1, b2 bool) bool {
		return b1 || b2
	})
}

func HTTPMatchAnd(reqMatchers ...RequestMatcher) cmux.Matcher {
	return httpMatch(reqMatchers, true, func(b1, b2 bool) bool {
		return b1 && b2
	})
}

func httpMatch(reqMatchers []RequestMatcher, init bool, fold func(b1, b2 bool) bool) cmux.Matcher {
	return func(reader io.Reader) bool {
		var (
			req *RequestPreface
			err error
		)

		switch v, r := parseHTTPVersion(reader); v {
		case HTTPVersionUnknown:
			return false
		case HTTPVersion10, HTTPVersion11:
			if req, err = parseHTTP1Request(r); err != nil {
				return false
			}
			req.Version = v
		case HTTPVersion2:
			if req, err = parseHTTP2Request(r); err != nil {
				return false
			}
			req.Version = v
		}

		if err != nil {
			return false
		}

		acc := init
		for idx := 0; idx < len(reqMatchers) && acc == init; idx++ {
			acc = fold(acc, reqMatchers[idx](req))
		}
		return acc
	}
}
