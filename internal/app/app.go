//go:generate mockgen -source=$GOFILE -destination=./mock/app.mock.go -package=mock

package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/path"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var (
	configFilePath  string
	logLevel        string
	developmentLogs bool
)

type contextKey string

const (
	loggerKey          contextKey = "gitlab.com/inetmock/inetmock/app/context/logger"
	configKey          contextKey = "gitlab.com/inetmock/inetmock/app/context/config"
	handlerRegistryKey contextKey = "gitlab.com/inetmock/inetmock/app/context/handlerRegistry"
	healthCheckerKey   contextKey = "gitlab.com/inetmock/inetmock/app/context/healthChecker"
	endpointManagerKey contextKey = "gitlab.com/inetmock/inetmock/app/context/endpointManager"
	certStoreKey       contextKey = "gitlab.com/inetmock/inetmock/app/context/certStore"
	eventStreamKey     contextKey = "gitlab.com/inetmock/inetmock/app/context/eventStream"
)

type App interface {
	api.PluginContext
	EventStream() audit.EventStream
	Config() config.Config
	Checker() health.Checker
	EndpointManager() endpoint.EndpointManager
	HandlerRegistry() api.HandlerRegistry
	Context() context.Context
	RootCommand() *cobra.Command
	MustRun()
	Shutdown()

	// WithCommands adds subcommands to the root command
	// requires nothing
	WithCommands(cmds ...*cobra.Command) App

	// WithHandlerRegistry builds up the handler registry
	// requires nothing
	WithHandlerRegistry(registrations ...api.Registration) App

	// WithHealthChecker adds the health checker mechanism
	// requires nothing
	WithHealthChecker() App

	// WithLogger configures the logging system
	// requires nothing
	WithLogger() App

	// WithEndpointManager creates an endpoint manager instance and adds it to the context
	// requires WithHandlerRegistry, WithHealthChecker and WithLogger
	WithEndpointManager() App

	// WithCertStore initializes the cert store
	// requires WithLogger and WithConfig
	WithCertStore() App

	// WithEventStream adds the audit event stream
	// requires WithLogger
	WithEventStream() App

	// WithConfig loads the config
	// requires nothing
	WithConfig() App

	WithInitTasks(task ...func(cmd *cobra.Command, args []string) (err error)) App
}

type app struct {
	rootCmd       *cobra.Command
	ctx           context.Context
	cancel        context.CancelFunc
	lateInitTasks []func(cmd *cobra.Command, args []string) (err error)
}

func (a *app) MustRun() {
	if err := a.rootCmd.Execute(); err != nil {
		if a.Logger() != nil {
			a.Logger().Error(
				"Failed to run inetmock",
				zap.Error(err),
			)
		} else {
			panic(err)
		}
	}
}

func (a *app) Logger() logging.Logger {
	return a.ctx.Value(loggerKey).(logging.Logger)
}

func (a *app) Config() config.Config {
	return a.ctx.Value(configKey).(config.Config)
}

func (a *app) CertStore() cert.Store {
	return a.ctx.Value(certStoreKey).(cert.Store)
}

func (a *app) Checker() health.Checker {
	return a.ctx.Value(healthCheckerKey).(health.Checker)
}

func (a *app) EndpointManager() endpoint.EndpointManager {
	return a.ctx.Value(endpointManagerKey).(endpoint.EndpointManager)
}

func (a *app) Audit() audit.Emitter {
	return a.ctx.Value(eventStreamKey).(audit.Emitter)
}

func (a *app) EventStream() audit.EventStream {
	return a.ctx.Value(eventStreamKey).(audit.EventStream)
}

func (a *app) HandlerRegistry() api.HandlerRegistry {
	return a.ctx.Value(handlerRegistryKey).(api.HandlerRegistry)
}

func (a *app) Context() context.Context {
	return a.ctx
}

func (a *app) RootCommand() *cobra.Command {
	return a.rootCmd
}

func (a *app) Shutdown() {
	a.cancel()
}

// WithCommands adds subcommands to the root command
// requires nothing
func (a *app) WithCommands(cmds ...*cobra.Command) App {
	a.rootCmd.AddCommand(cmds...)
	return a
}

// WithHandlerRegistry builds up the handler registry
// requires nothing
func (a *app) WithHandlerRegistry(registrations ...api.Registration) App {
	registry := api.NewHandlerRegistry()

	for _, registration := range registrations {
		if err := registration(registry); err != nil {
			panic(err)
		}
	}

	a.ctx = context.WithValue(a.ctx, handlerRegistryKey, registry)

	return a
}

// WithHealthChecker adds the health checker mechanism
// requires nothing
func (a *app) WithHealthChecker() App {
	checker := health.New()
	a.ctx = context.WithValue(a.ctx, healthCheckerKey, checker)
	return a
}

// WithLogger configures the logging system
// requires nothing
func (a *app) WithLogger() App {
	a.lateInitTasks = append(a.lateInitTasks, func(cmd *cobra.Command, args []string) (err error) {
		logging.ConfigureLogging(
			logging.ParseLevel(logLevel),
			developmentLogs,
			map[string]interface{}{
				"cwd":  path.WorkingDirectory(),
				"cmd":  cmd.Name(),
				"args": args,
			},
		)

		var logger logging.Logger
		if logger, err = logging.CreateLogger(); err != nil {
			return
		}
		a.ctx = context.WithValue(a.ctx, loggerKey, logger)
		return
	})
	return a
}

// WithEndpointManager creates an endpoint manager instance and adds it to the context
// requires WithHandlerRegistry, WithHealthChecker and WithLogger
func (a *app) WithEndpointManager() App {
	a.lateInitTasks = append(a.lateInitTasks, func(_ *cobra.Command, _ []string) (err error) {
		epMgr := endpoint.NewEndpointManager(
			a.HandlerRegistry(),
			a.Logger().Named("EndpointManager"),
			a.Checker(),
			a,
		)

		a.ctx = context.WithValue(a.ctx, endpointManagerKey, epMgr)
		return
	})
	return a
}

// WithCertStore initializes the cert store
// requires WithLogger and WithConfig
func (a *app) WithCertStore() App {
	a.lateInitTasks = append(a.lateInitTasks, func(cmd *cobra.Command, args []string) (err error) {
		var certStore cert.Store
		if certStore, err = cert.NewDefaultStore(
			a.Config(),
			a.Logger().Named("CertStore"),
		); err != nil {
			return
		}

		a.ctx = context.WithValue(a.ctx, certStoreKey, certStore)
		return
	})
	return a
}

// WithEventStream adds the audit event stream
// requires WithLogger
func (a *app) WithEventStream() App {
	a.lateInitTasks = append(a.lateInitTasks, func(_ *cobra.Command, _ []string) (err error) {
		var eventStream audit.EventStream
		eventStream, err = audit.NewEventStream(
			a.Logger().Named("EventStream"),
			audit.WithSinkBufferSize(10),
		)
		if err != nil {
			return
		}

		if err = eventStream.RegisterSink(sink.NewLogSink(a.Logger().Named("LogSink"))); err != nil {
			return
		}

		a.ctx = context.WithValue(a.ctx, eventStreamKey, eventStream)
		return
	})

	return a
}

// WithConfig loads the config
// requires nothing
func (a *app) WithConfig() App {
	a.lateInitTasks = append(a.lateInitTasks, func(cmd *cobra.Command, _ []string) (err error) {
		cfg := config.CreateConfig(cmd.Flags())
		if err = cfg.ReadConfig(configFilePath); err != nil {
			return
		}
		a.ctx = context.WithValue(a.ctx, configKey, cfg)
		return
	})

	return a
}

func (a *app) WithInitTasks(task ...func(cmd *cobra.Command, args []string) (err error)) App {
	a.lateInitTasks = append(a.lateInitTasks, task...)
	return a
}

func NewApp(name, short string) App {
	ctx, cancel := initAppContext()
	a := &app{
		rootCmd: &cobra.Command{
			Use:   name,
			Short: short,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	a.rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file that should be used")
	a.rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level to use")
	a.rootCmd.PersistentFlags().BoolVar(&developmentLogs, "development-logs", false, "Enable development mode logs")

	a.rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		for _, initTask := range a.lateInitTasks {
			err = multierr.Append(err, initTask(cmd, args))
		}
		return
	}

	return a
}

func initAppContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signals
		cancel()
	}()

	return ctx, cancel
}
