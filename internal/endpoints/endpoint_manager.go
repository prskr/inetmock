package endpoints

import (
	"fmt"
	"github.com/baez90/inetmock/internal/plugins"
	config2 "github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
	"sync"
	"time"
)

type EndpointManager interface {
	RegisteredEndpoints() []Endpoint
	StartedEndpoints() []Endpoint
	CreateEndpoint(name string, multiHandlerConfig config2.MultiHandlerConfig) error
	StartEndpoints()
	ShutdownEndpoints()
}

func NewEndpointManager(logger logging.Logger) EndpointManager {
	return &endpointManager{
		logger:   logger,
		registry: plugins.Registry(),
	}
}

type endpointManager struct {
	logger                   logging.Logger
	registeredEndpoints      []Endpoint
	properlyStartedEndpoints []Endpoint
	registry                 plugins.HandlerRegistry
}

func (e endpointManager) RegisteredEndpoints() []Endpoint {
	return e.registeredEndpoints
}

func (e endpointManager) StartedEndpoints() []Endpoint {
	return e.properlyStartedEndpoints
}

func (e *endpointManager) CreateEndpoint(name string, multiHandlerConfig config2.MultiHandlerConfig) error {
	for _, handlerConfig := range multiHandlerConfig.HandlerConfigs() {
		if handler, ok := e.registry.HandlerForName(multiHandlerConfig.Handler); ok {
			e.registeredEndpoints = append(e.registeredEndpoints, &endpoint{
				name:    name,
				handler: handler,
				config:  handlerConfig,
			})
		} else {
			return fmt.Errorf("no matching handler registered for names %s", multiHandlerConfig.Handler)
		}
	}

	return nil
}

func (e *endpointManager) StartEndpoints() {
	startTime := time.Now()
	for _, endpoint := range e.registeredEndpoints {
		endpointLogger := e.logger.With(
			zap.String("endpoint", endpoint.Name()),
		)
		endpointLogger.Info("Starting endpoint")
		if ok := startEndpoint(endpoint, endpointLogger); ok {
			e.properlyStartedEndpoints = append(e.properlyStartedEndpoints, endpoint)
			endpointLogger.Info("successfully started endpoint")
		} else {
			endpointLogger.Error("error occurred during endpoint startup - will be skipped for now")
		}
	}
	endpointStartupDuration := time.Since(startTime)
	e.logger.Info(
		"Startup of all endpoints completed",
		zap.Duration("startupTime", endpointStartupDuration),
	)
}

func (e *endpointManager) ShutdownEndpoints() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(e.properlyStartedEndpoints))

	for _, endpoint := range e.properlyStartedEndpoints {
		endpointLogger := e.logger.With(
			zap.String("endpoint", endpoint.Name()),
		)
		endpointLogger.Info("Triggering shutdown of endpoint")
		go shutdownEndpoint(endpoint, endpointLogger, &waitGroup)
	}

	waitGroup.Wait()
}

func startEndpoint(ep Endpoint, logger logging.Logger) (success bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal(
				"recovered panic during startup of endpoint",
				zap.Any("recovered", r),
			)
		}
	}()
	if err := ep.Start(); err != nil {
		logger.Error(
			"failed to start endpoint",
			zap.Error(err),
		)
	} else {
		success = true
	}
	return
}

func shutdownEndpoint(ep Endpoint, logger logging.Logger, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal(
				"recovered panic during shutdown of endpoint",
				zap.Any("recovered", r),
			)
		}
		wg.Done()
	}()
	if err := ep.Shutdown(); err != nil {
		logger.Error(
			"Failed to shutdown endpoint",
			zap.Error(err),
		)
	}
}
