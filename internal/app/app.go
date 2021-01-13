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
	"go.uber.org/zap"
)

var (
	configFilePath  string
	logLevel        string
	developmentLogs bool
)

type App interface {
	api.PluginContext
	Config() config.Config
	Checker() health.Checker
	EndpointManager() endpoint.EndpointManager
	HandlerRegistry() api.HandlerRegistry
	Context() context.Context
	MustRun()
	Shutdown()
	WithCommands(cmds ...*cobra.Command) App
}

type app struct {
	cfg             config.Config
	rootCmd         *cobra.Command
	rootLogger      logging.Logger
	certStore       cert.Store
	checker         health.Checker
	endpointManager endpoint.EndpointManager
	registry        api.HandlerRegistry
	ctx             context.Context
	cancel          context.CancelFunc
	eventStream     audit.EventStream
}

func (a *app) MustRun() {
	if err := a.rootCmd.Execute(); err != nil {
		if a.rootLogger != nil {
			a.rootLogger.Error(
				"Failed to run inetmock",
				zap.Error(err),
			)
		} else {
			panic(err)
		}
	}
}

func (a app) Logger() logging.Logger {
	return a.rootLogger
}

func (a app) Config() config.Config {
	return a.cfg
}

func (a app) CertStore() cert.Store {
	return a.certStore
}

func (a app) Checker() health.Checker {
	return a.checker
}

func (a app) EndpointManager() endpoint.EndpointManager {
	return a.endpointManager
}

func (a app) Audit() audit.Emitter {
	return a.eventStream
}

func (a app) HandlerRegistry() api.HandlerRegistry {
	return a.registry
}

func (a app) Context() context.Context {
	return a.ctx
}

func (a app) Shutdown() {
	a.cancel()
}

func (a *app) WithCommands(cmds ...*cobra.Command) App {
	a.rootCmd.AddCommand(cmds...)
	return a
}

func NewApp(registrations ...api.Registration) (inetmockApp App, err error) {
	registry := api.NewHandlerRegistry()

	for _, registration := range registrations {
		if err = registration(registry); err != nil {
			return
		}
	}

	ctx, cancel := initAppContext()

	a := &app{
		rootCmd: &cobra.Command{
			Short: "INetMock is lightweight internet mock",
		},
		checker:  health.New(),
		registry: registry,
		ctx:      ctx,
		cancel:   cancel,
	}

	a.rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file that should be used")
	a.rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level to use")
	a.rootCmd.PersistentFlags().BoolVar(&developmentLogs, "development-logs", false, "Enable development mode logs")

	a.rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		logging.ConfigureLogging(
			logging.ParseLevel(logLevel),
			developmentLogs,
			map[string]interface{}{
				"cwd": path.WorkingDirectory(),
			},
		)

		if a.rootLogger, err = logging.CreateLogger(); err != nil {
			return
		}

		a.endpointManager = endpoint.NewEndpointManager(
			a.registry,
			a.Logger().Named("EndpointManager"),
			a.checker,
			a,
		)

		a.cfg = config.CreateConfig(cmd.Flags())

		if err = a.cfg.ReadConfig(configFilePath); err != nil {
			return
		}

		if a.certStore, err = cert.NewDefaultStore(a.cfg, a.rootLogger); err != nil {
			return
		}

		a.eventStream, err = audit.NewEventStream(
			a.Logger().Named("EventStream"),
			audit.WithSinkBufferSize(10),
		)
		if err != nil {
			return
		}

		err = a.eventStream.RegisterSink(sink.NewLogSink(a.Logger().Named("LogSink")))
		return
	}

	return a, nil
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
