package endpoint

import (
	"github.com/mitchellh/mapstructure"
)

type endpointLifecycle struct {
	endpointName string
	uplink       Uplink
	opts         map[string]interface{}
}

func NewEndpointLifecycle(
	endpointName string,
	uplink Uplink,
	opts map[string]interface{},
) Lifecycle {
	return &endpointLifecycle{
		endpointName: endpointName,
		uplink:       uplink,
		opts:         opts,
	}
}

func (e *endpointLifecycle) Name() string {
	return e.endpointName
}

func (e *endpointLifecycle) Uplink() Uplink {
	return e.uplink
}

func (e *endpointLifecycle) UnmarshalOptions(cfg interface{}) error {
	return mapstructure.Decode(e.opts, cfg)
}
