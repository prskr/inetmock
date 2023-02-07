package health

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
	"inetmock.icb4dc0.de/inetmock/pkg/health/dns"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

func NewDNSRuleCheck(name string, resolvers ResolverForModule, logger logging.Logger, check *rules.Check) (Check, error) {
	switch {
	case name == "":
		return nil, ErrEmptyCheckName
	case resolvers == nil:
		return nil, errors.New("passed resolver is nil")
	case check == nil:
		return nil, errors.New("passed check is nil")
	}

	var (
		initiator dns.Initiator
		chain     dns.ValidationChain
		resolver  dns.Resolver
		err       error
	)

	if initiator, err = dns.CheckForRule(check, logger); err != nil {
		return nil, err
	}

	if chain, err = dns.ValidatorsForRule(check); err != nil {
		return nil, err
	}

	if resolver, err = resolvers.ResolverForModule(check.Initiator.Module); err != nil {
		return nil, err
	}

	return NewCheckFunc(name, func(ctx context.Context) (err error) {
		const maxRetries = 10
		defer func() {
			if rec := recover(); rec != nil {
				err = errors.Join(err, fmt.Errorf("recovered panic during check: %v", rec))
			}
		}()

		var resp *dns.Response
		for tries := 0; ctx.Err() == nil && tries < maxRetries; tries++ {
			if resp, err = initiator.Do(ctx, resolver); err != nil {
				logger.Warn("Failed to initiate check", zap.Error(err))
				if ctx.Err() != nil {
					return err
				}
			} else if len(resp.Addresses) == 0 && len(resp.Hosts) == 0 {
				logger.Warn("Response empty")
			} else {
				break
			}
		}

		if err != nil {
			return err
		}

		return chain.Matches(resp)
	}), nil
}
