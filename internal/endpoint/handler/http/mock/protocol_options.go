package mock

import (
	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

type httpOptions struct {
	Rules []string
}

func loadFromConfig(lifecycle endpoint.Lifecycle) (opts httpOptions, err error) {
	err = lifecycle.UnmarshalOptions(&opts)
	return
}
