package mock

import (
	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
)

type httpOptions struct {
	Rules []string
}

func loadFromConfig(startupSpec *endpoint.StartupSpec) (opts httpOptions, err error) {
	err = startupSpec.UnmarshalOptions(&opts)
	return
}
