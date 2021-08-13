package health

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/health/dns"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func NewDNSRuleCheck(name string, resolver dns.Resolver, logger logging.Logger, check *rules.Check) (Check, error) {
	switch {
	case name == "":
		return nil, ErrEmptyCheckName
	case resolver == nil:
		return nil, errors.New("passed resolver is nil")
	case check == nil:
		return nil, errors.New("passed check is nil")
	}

	var (
		initiator dns.Initiator
		chain     dns.ValidationChain
		err       error
	)

	if initiator, err = dns.CheckForRule(check, logger); err != nil {
		return nil, err
	}

	if chain, err = dns.ValidatorsForRule(check); err != nil {
		return nil, err
	}

	return NewCheckFunc(name, func(ctx context.Context) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				err = multierr.Append(err, fmt.Errorf("recovered panic during check: %v", rec))
			}
		}()

		if resp, err := initiator.Do(ctx, resolver); err != nil {
			return err
		} else {
			return chain.Matches(resp)
		}
	}), nil
}
