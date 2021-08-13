package health

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	ErrEmptyCheckName     = errors.New("name of the check must not be empty")
	ErrAmbiguousCheckName = errors.New("a check with the same name is already registered")
)

func New() Checker {
	return &checker{
		registeredChecks: map[string]Check{},
	}
}

func NewFromConfig(logger logging.Logger, cfg Config, tlsConfig *tls.Config) (Checker, error) {
	var httpClient = HTTPClient(cfg, tlsConfig)
	var dnsResolver = DNSResolver(cfg)

	var checker = &checker{
		registeredChecks: make(map[string]Check),
	}

	for idx := range cfg.Rules {
		var rawRule = cfg.Rules[idx]
		var check = new(rules.Check)
		if err := rules.Parse(rawRule.Rule, check); err != nil {
			return nil, err
		}
		switch strings.ToLower(check.Initiator.Module) {
		case "":
			return nil, fmt.Errorf("initiator of check '%s' has no module", rawRule.Name)
		case "http":
			if compiledCheck, err := NewHTTPRuleCheck(rawRule.Name, httpClient, logger, check); err != nil {
				return nil, err
			} else if err := checker.AddCheck(compiledCheck); err != nil {
				return nil, err
			}
		case "dns":
			if compiledCheck, err := NewDNSRuleCheck(rawRule.Name, dnsResolver, logger, check); err != nil {
				return nil, err
			} else if err := checker.AddCheck(compiledCheck); err != nil {
				return nil, err
			}
		}
	}

	return checker, nil
}

type Checker interface {
	AddCheck(check Check) error
	Status(ctx context.Context) Result
}

type Check interface {
	Name() string
	Status(ctx context.Context) error
}
