package health

import (
	"context"
	"fmt"
	gohttp "net/http"

	"go.uber.org/multierr"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/health/http"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func NewHTTPRuleCheck(name string, client *gohttp.Client, logger logging.Logger, check *rules.Check) (Check, error) {
	var (
		initiator http.Initiator
		chain     http.ValidationChain
		err       error
	)

	if initiator, err = http.InitiatorForRule(check, logger); err != nil {
		return nil, err
	}

	if chain, err = http.ValidatorsForRule(check); err != nil {
		return nil, err
	}

	return NewCheckFunc(name, func(ctx context.Context) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				err = multierr.Append(err, fmt.Errorf("recovered panic in HTTP health check: %v", rec))
			}
		}()

		var resp *gohttp.Response

		if resp, err = initiator.Do(ctx, client); err != nil {
			return err
		}

		return chain.Matches(resp)
	}), nil
}
