package health

import (
	"context"
	"errors"
	"fmt"
	gohttp "net/http"

	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
	"inetmock.icb4dc0.de/inetmock/pkg/health/http"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

func NewHTTPRuleCheck(name string, clients HTTPClientForModule, logger logging.Logger, check *rules.Check) (Check, error) {
	var (
		initiator http.Initiator
		chain     http.ValidationChain
		client    *gohttp.Client
		err       error
	)

	if initiator, err = http.InitiatorForRule(check, logger); err != nil {
		return nil, err
	}

	if chain, err = http.ValidatorsForRule(check); err != nil {
		return nil, err
	}

	if client, err = clients.ClientForModule(check.Initiator.Module); err != nil {
		return nil, err
	}

	return NewCheckFunc(name, func(ctx context.Context) (err error) {
		const maxRetries = 10
		defer func() {
			if rec := recover(); rec != nil {
				err = errors.Join(err, fmt.Errorf("recovered panic in HTTP health check: %v", rec))
			}
		}()

		var resp *gohttp.Response

		for tries := 0; ctx.Err() == nil && tries < maxRetries; tries++ {
			if resp, err = initiator.Do(ctx, client); err != nil {
				logger.Warn("Failed to initiate health check", zap.Error(err))
				if ctx.Err() != nil {
					return err
				}
			} else {
				break
			}
		}

		return chain.Matches(resp)
	}), nil
}
