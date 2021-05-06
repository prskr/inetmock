package endpoint

import (
	"context"
	"fmt"

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
	go e.doStart(ctx, lifecycle, logger, startupResult)

	select {
	case err = <-startupResult:
	case <-ctx.Done():
		err = ErrStartupTimeout
	}

	return
}

func (e *Endpoint) doStart(ctx context.Context, lifecycle Lifecycle, logger logging.Logger, startupResult chan error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("Startup error recovered", zap.Any("recovered", r))
			startupResult <- fmt.Errorf("recovered: %v", r)
		}
	}()

	startupResult <- e.Handler.Start(ctx, lifecycle)
}
