package endpoint

import (
	"context"

	"github.com/mitchellh/mapstructure"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type endpointLifecycle struct {
	endpointName string
	ctx          context.Context
	logger       logging.Logger
	certStore    cert.Store
	emitter      audit.Emitter
	uplink       Uplink
	tls          bool
	opts         map[string]interface{}
}

func NewEndpointLifecycleFromContext(
	endpointName string,
	ctx context.Context,
	logger logging.Logger,
	certStore cert.Store,
	emitter audit.Emitter,
	uplink Uplink,
	opts map[string]interface{},
) Lifecycle {
	return &endpointLifecycle{
		endpointName: endpointName,
		ctx:          ctx,
		logger:       logger,
		certStore:    certStore,
		emitter:      emitter,
		uplink:       uplink,
		opts:         opts,
	}
}

func (e endpointLifecycle) Name() string {
	return e.endpointName
}

func (e endpointLifecycle) Uplink() Uplink {
	return e.uplink
}

func (e endpointLifecycle) Logger() logging.Logger {
	return e.logger
}

func (e endpointLifecycle) CertStore() cert.Store {
	return e.certStore
}

func (e endpointLifecycle) Audit() audit.Emitter {
	return e.emitter
}

func (e endpointLifecycle) Context() context.Context {
	return e.ctx
}

func (e endpointLifecycle) UnmarshalOptions(cfg interface{}) error {
	return mapstructure.Decode(e.opts, cfg)
}
