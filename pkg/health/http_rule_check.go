package health

import (
	"context"
	gohttp "net/http"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/health/http"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func NewHTTPRuleCheck(name string, client *gohttp.Client, logger logging.Logger, check *rules.Check) (Check, error) {
	var err error
	var initiator http.Initiator
	if initiator, err = http.InitiatorForRule(check, logger); err != nil {
		return nil, err
	}

	var filters []http.Validator
	if filters, err = http.ValidatorsForRule(check); err != nil {
		return nil, err
	}

	return NewCheckFunc(name, func(ctx context.Context) error {
		var err error
		var resp *gohttp.Response
		if resp, err = initiator.Do(ctx, client); err != nil {
			return err
		}

		for idx := range filters {
			filter := filters[idx]
			if err := filter.Matches(resp); err != nil {
				return err
			}
		}

		return nil
	}), nil
}
