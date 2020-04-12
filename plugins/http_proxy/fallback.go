package main

import (
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
)

const (
	passthroughStrategyName = "passthrough"
	notFoundStrategyName    = "notfound"
)

var (
	fallbackStrategies map[string]ProxyFallbackStrategy
)

func init() {
	fallbackStrategies = map[string]ProxyFallbackStrategy{
		passthroughStrategyName: &passThroughFallbackStrategy{},
		notFoundStrategyName:    &notFoundFallbackStrategy{},
	}
}

func StrategyForName(name string) ProxyFallbackStrategy {
	if strategy, ok := fallbackStrategies[name]; ok {
		return strategy
	}
	return fallbackStrategies[notFoundStrategyName]
}

type ProxyFallbackStrategy interface {
	Apply(request *http.Request) (*http.Response, error)
}

type passThroughFallbackStrategy struct {
}

func (p passThroughFallbackStrategy) Apply(request *http.Request) (*http.Response, error) {
	return nil, nil
}

type notFoundFallbackStrategy struct {
}

func (n notFoundFallbackStrategy) Apply(request *http.Request) (response *http.Response, err error) {
	response = goproxy.NewResponse(
		request,
		goproxy.ContentTypeText,
		http.StatusNotFound,
		"The requested resource was not found",
	)
	return
}
