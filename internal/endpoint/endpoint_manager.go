package endpoint

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

const (
	startupTimeoutDuration = 100 * time.Millisecond
)

type EndpointManager interface {
	RegisteredEndpoints() []Endpoint
	StartedEndpoints() []Endpoint
	CreateEndpoint(name string, multiHandlerConfig config.EndpointConfig) error
	StartEndpoints()
	ShutdownEndpoints()
}

func NewEndpointManager(registry api.HandlerRegistry, logging logging.Logger, checker health.Checker, pluginContext api.PluginContext) EndpointManager {
	return &endpointManager{
		registry:      registry,
		logger:        logging,
		checker:       checker,
		pluginContext: pluginContext,
	}
}

type endpointManager struct {
	registry                 api.HandlerRegistry
	logger                   logging.Logger
	checker                  health.Checker
	pluginContext            api.PluginContext
	registeredEndpoints      []Endpoint
	properlyStartedEndpoints []Endpoint
}

func (e endpointManager) RegisteredEndpoints() []Endpoint {
	return e.registeredEndpoints
}

func (e endpointManager) StartedEndpoints() []Endpoint {
	return e.properlyStartedEndpoints
}

func (e *endpointManager) CreateEndpoint(name string, endpointConfig config.EndpointConfig) error {
	for _, handlerConfig := range endpointConfig.HandlerConfigs() {
		if handler, ok := e.registry.HandlerForName(endpointConfig.Handler); ok {
			e.registeredEndpoints = append(e.registeredEndpoints, &endpoint{
				id:      uuid.New(),
				name:    name,
				handler: handler,
				config:  handlerConfig,
			})
		} else {
			return fmt.Errorf("no matching handler registered for names %s", endpointConfig.Handler)
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
		if ok := startEndpoint(endpoint, e.pluginContext, endpointLogger); ok {
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
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(len(e.properlyStartedEndpoints))

	parentCtx, _ := context.WithTimeout(context.Background(), shutdownTimeout)

	perHandlerTimeout := e.shutdownTimePerEndpoint()

	for _, endpoint := range e.properlyStartedEndpoints {
		ctx, _ := context.WithTimeout(parentCtx, perHandlerTimeout)
		endpointLogger := e.logger.With(
			zap.String("endpoint", endpoint.Name()),
		)
		endpointLogger.Info("Triggering shutdown of endpoint")
		go shutdownEndpoint(ctx, endpoint, endpointLogger, waitGroup)
	}

	waitGroup.Wait()
}

func startEndpoint(ep Endpoint, ctx api.PluginContext, logger logging.Logger) (success bool) {
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

		if err := ep.Start(ctx); err != nil {
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
