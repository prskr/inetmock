package endpoint

import (
	"context"
	"time"

	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	startupTimeoutDuration = 100 * time.Millisecond
)

type Endpoint struct {
	Spec
	name   string
	uplink Uplink
}

func (e *Endpoint) Start(ctx context.Context, logger logging.Logger, lifecycle Lifecycle) (err error) {
	startupResult := make(chan error)
	ctx, cancel := context.WithTimeout(ctx, startupTimeoutDuration)
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Fatal("Startup error recovered", zap.Any("recovered", r))
			}
		}()

		startupResult <- e.Handler.Start(ctx, lifecycle)
	}()

	select {
	case err = <-startupResult:
	case <-ctx.Done():
		err = ErrStartupTimeout
	}

	return
}
