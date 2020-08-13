package endpoints

import (
	"context"
	"fmt"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/health"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	startupTimeoutDuration = 100 * time.Millisecond
)

type EndpointManager interface {
	RegisteredEndpoints() []Endpoint
	StartedEndpoints() []Endpoint
	CreateEndpoint(name string, multiHandlerConfig config.MultiHandlerConfig) error
	StartEndpoints()
	ShutdownEndpoints()
}

func NewEndpointManager(checker health.Checker, logger logging.Logger) EndpointManager {
	return &endpointManager{
		logger:   logger,
		checker:  checker,
		registry: api.Registry(),
	}
}

type endpointManager struct {
	logger                   logging.Logger
	checker                  health.Checker
	registeredEndpoints      []Endpoint
	properlyStartedEndpoints []Endpoint
	registry                 api.HandlerRegistry
}

func (e endpointManager) RegisteredEndpoints() []Endpoint {
	return e.registeredEndpoints
}

func (e endpointManager) StartedEndpoints() []Endpoint {
	return e.properlyStartedEndpoints
}

func (e *endpointManager) CreateEndpoint(name string, multiHandlerConfig config.MultiHandlerConfig) error {
	for _, handlerConfig := range multiHandlerConfig.HandlerConfigs() {
		if handler, ok := e.registry.HandlerForName(multiHandlerConfig.Handler); ok {
			e.registeredEndpoints = append(e.registeredEndpoints, &endpoint{
				id:      uuid.New(),
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
			_ = e.checker.RegisterCheck(
				endpointComponentName(endpoint),
				health.StaticResultCheckWithMessage(health.HEALTHY, "Successfully started"),
			)
			e.properlyStartedEndpoints = append(e.properlyStartedEndpoints, endpoint)
			endpointLogger.Info("successfully started endpoint")
		} else {
			_ = e.checker.RegisterCheck(
				endpointComponentName(endpoint),
				health.StaticResultCheckWithMessage(health.UNHEALTHY, "failed to start"),
			)
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

	parentCtx, _ := context.WithTimeout(context.Background(), shutdownTimeout)

	perHandlerTimeout := e.shutdownTimePerEndpoint()

	for _, endpoint := range e.properlyStartedEndpoints {
		ctx, _ := context.WithTimeout(parentCtx, perHandlerTimeout)
		endpointLogger := e.logger.With(
			zap.String("endpoint", endpoint.Name()),
		)
		endpointLogger.Info("Triggering shutdown of endpoint")
		go shutdownEndpoint(ctx, endpoint, endpointLogger, &waitGroup)
	}

	waitGroup.Wait()
}

func startEndpoint(ep Endpoint, logger logging.Logger) (success bool) {
	startSuccessful := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Fatal(
					"recovered panic during startup of endpoint",
					zap.Any("recovered", r),
				)
				startSuccessful <- false
			}
		}()

		if err := ep.Start(); err != nil {
			logger.Error(
				"failed to start endpoint",
				zap.Error(err),
			)
			startSuccessful <- false
		} else {
			startSuccessful <- true
		}
	}()

	select {
	case success = <-startSuccessful:
	case <-time.After(startupTimeoutDuration):
		success = false
	}

	return
}

func shutdownEndpoint(ctx context.Context, ep Endpoint, logger logging.Logger, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal(
				"recovered panic during shutdown of endpoint",
				zap.Any("recovered", r),
			)
		}
		wg.Done()
	}()
	if err := ep.Shutdown(ctx); err != nil {
		logger.Error(
			"Failed to shutdown endpoint",
			zap.Error(err),
		)
	}
}

func (e *endpointManager) shutdownTimePerEndpoint() time.Duration {
	return time.Duration((float64(shutdownTimeout) * 0.9) / float64(len(e.properlyStartedEndpoints)))
}

func endpointComponentName(ep Endpoint) string {
	return fmt.Sprintf("endpoint_%s", ep.Name())
}
