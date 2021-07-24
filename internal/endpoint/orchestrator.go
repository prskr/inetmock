package endpoint

import (
	"context"
	"errors"
	"fmt"

	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	ErrStartupTimeout    = errors.New("endpoint did not start in time")
	ErrUnknownHandlerRef = errors.New("no handler for given key registered")
)

type Orchestrator interface {
	RegisterListener(spec ListenerSpec) error
	Endpoints() []Endpoint
	StartEndpoints(ctx context.Context) (errChan chan error)
}

func NewOrchestrator(
	certStore cert.Store,
	registry HandlerRegistry,
	logger logging.Logger,
) Orchestrator {
	return &orchestrator{
		registry:  registry,
		logger:    logger,
		certStore: certStore,
	}
}

type orchestrator struct {
	registry  HandlerRegistry
	logger    logging.Logger
	certStore cert.Store

	endpointListeners []Endpoint
	muxes             []cmux.CMux
}

func (e *orchestrator) RegisterListener(spec ListenerSpec) (err error) {
	for name, s := range spec.Endpoints {
		if handler, registered := e.registry.HandlerForName(s.HandlerRef); registered {
			s.Handler = handler
			spec.Endpoints[name] = s
		} else {
			return fmt.Errorf("%s: %w", s.HandlerRef, ErrUnknownHandlerRef)
		}
	}

	var endpoints []Endpoint
	var muxes []cmux.CMux
	if endpoints, muxes, err = spec.ConfigureMultiplexing(e.certStore.TLSConfig()); err != nil {
		return
	}

	e.endpointListeners = append(e.endpointListeners, endpoints...)
	e.muxes = append(e.muxes, muxes...)

	return
}

func (e orchestrator) Endpoints() []Endpoint {
	return e.endpointListeners
}

func (e *orchestrator) StartEndpoints(ctx context.Context) chan error {
	var errChan = make(chan error)
	for _, epListener := range e.endpointListeners {
		endpointLogger := e.logger.With(
			zap.String("epListener", epListener.name),
		)
		endpointLogger.Debug("Starting endpoint listener")
		lifecycle := NewEndpointLifecycle(
			epListener.name,
			epListener.uplink,
			epListener.Options,
		)

		if err := epListener.Start(ctx, e.logger.With(zap.String("epListener", epListener.name)), lifecycle); err == nil {
			endpointLogger.Debug("Successfully started epListener")
		} else {
			endpointLogger.Error("error occurred during epListener startup - will be skipped for now", zap.Error(err))
		}
	}
	e.logger.Info("Startup of all endpoints completed")

	for _, mux := range e.muxes {
		go func(mux cmux.CMux) {
			mux.HandleError(func(err error) bool {
				errChan <- err
				return true
			})
			if err := mux.Serve(); err != nil && !errors.Is(err, cmux.ErrListenerClosed) {
				errChan <- err
			}
		}(mux)
	}

	return errChan
}
