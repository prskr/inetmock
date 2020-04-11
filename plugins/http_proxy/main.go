package main

import (
	"github.com/baez90/inetmock/internal/config"
	"go.uber.org/zap"
	"sync"
)

const (
	name = "http_proxy"
)

type httpProxy struct {
	logger *zap.Logger
}

func (h httpProxy) Run(config config.HandlerConfig) {
	panic("implement me")
}

func (h httpProxy) Shutdown(wg *sync.WaitGroup) {
	panic("implement me")
}
