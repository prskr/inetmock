package mock

import (
	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

type httpOptions struct {
	Rules []string
}

func loadFromConfig(startupSpec *endpoint.StartupSpec) (opts httpOptions, err error) {
	err = startupSpec.UnmarshalOptions(&opts)
	return
}
