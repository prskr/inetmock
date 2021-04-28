package endpoint

import (
	"context"

	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type Endpoint struct {
	Spec
	name   string
	uplink Uplink
}

func (e *Endpoint) Start(ctx context.Context, logger logging.Logger, lifecycle Lifecycle) (err error) {
	startupResult := make(chan error)
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
