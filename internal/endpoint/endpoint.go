package endpoint

import (
	"context"
	"time"

	"go.uber.org/zap"
)

const (
	startupTimeoutDuration = 100 * time.Millisecond
)

type Endpoint struct {
	Spec
	name   string
	uplink Uplink
}

func (e *Endpoint) Start(lifecycle Lifecycle) (err error) {
	startupResult := make(chan error)
	ctx, cancel := context.WithTimeout(lifecycle.Context(), startupTimeoutDuration)
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				lifecycle.Logger().Fatal("Startup error recovered", zap.Any("recovered", r))
			}
		}()

		startupResult <- e.Handler.Start(lifecycle)
	}()

	select {
	case err = <-startupResult:
	case <-ctx.Done():
		err = ErrStartupTimeout
	}

	return
}
