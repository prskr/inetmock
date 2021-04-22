//go:generate mockgen -source=$GOFILE -destination=./mock/app.mock.go -package=mock

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/path"
)

var (
	configFilePath  string
	logLevel        string
	developmentLogs bool
)

type Spec struct {
	Name                    string
	Short                   string
	Config                  interface{}
	IgnoreMissingConfigFile bool
	Defaults                map[string]interface{}
	FlagBindings            map[string]func(flagSet *pflag.FlagSet) *pflag.Flag
	SubCommands             []*cobra.Command
	LateInitTasks           []func(cmd *cobra.Command, args []string) (err error)
}

type App interface {
	Logger() logging.Logger
	Context() context.Context
	RootCommand() *cobra.Command
	MustRun()
	Shutdown()
}

type app struct {
	rootCmd *cobra.Command
	ctx     context.Context
	cancel  context.CancelFunc
	logger  logging.Logger
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
	return a.logger
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

func NewApp(spec Spec) App {
	if spec.Defaults == nil {
		spec.Defaults = make(map[string]interface{})
	}

	ctx, cancel := initAppContext()
	a := &app{
		rootCmd: &cobra.Command{
			Use:          spec.Name,
			Short:        spec.Short,
			SilenceUsage: true,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	lateInitTasks := []func(cmd *cobra.Command, args []string) (err error){
		func(*cobra.Command, []string) (err error) {
			return spec.readConfig(a.rootCmd)
		},
		func(cmd *cobra.Command, args []string) (err error) {
			logging.ConfigureLogging(
				logging.ParseLevel(logLevel),
				developmentLogs,
				map[string]interface{}{
					"cwd":  path.WorkingDirectory(),
					"cmd":  cmd.Name(),
					"args": args,
				},
			)

			if a.logger, err = logging.CreateLogger(); err != nil {
				return
			}
			return
		},
	}

	lateInitTasks = append(lateInitTasks, spec.LateInitTasks...)

	a.rootCmd.AddCommand(completionCmd)
	a.rootCmd.AddCommand(spec.SubCommands...)

	a.rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file that should be used")
	a.rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level to use")
	a.rootCmd.PersistentFlags().BoolVar(&developmentLogs, "development-logs", false, "Enable development mode logs")

	a.rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		for _, initTask := range lateInitTasks {
			if err = initTask(cmd, args); err != nil {
				return
			}
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

func (s Spec) readConfig(rootCmd *cobra.Command) error {
	viperCfg := viper.NewWithOptions()
	viperCfg.SetConfigName("config")
	viperCfg.SetConfigType("yaml")
	viperCfg.AddConfigPath(fmt.Sprintf("/etc/%s/", strings.ToLower(s.Name)))
	viperCfg.AddConfigPath(fmt.Sprintf("$HOME/.%s/", strings.ToLower(s.Name)))
	viperCfg.AddConfigPath(".")
	viperCfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperCfg.SetEnvPrefix("INETMOCK")
	viperCfg.AutomaticEnv()

	if s.FlagBindings != nil {
		for key, selector := range s.FlagBindings {
			if err := viperCfg.BindPFlag(key, selector(rootCmd.Flags())); err != nil {
				return err
			}
		}
	}

	for k, v := range s.Defaults {
		viperCfg.SetDefault(k, v)
	}

	if configFilePath != "" && path.FileExists(configFilePath) {
		viperCfg.SetConfigFile(configFilePath)
	}

	if err := viperCfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !(ok && s.IgnoreMissingConfigFile) {
			return err
		}
	}

	if err := viperCfg.Unmarshal(s.Config); err != nil {
		return err
	}

	return nil
}
