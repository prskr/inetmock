//go:generate mockgen -source=$GOFILE -destination=./mock/app.mock.go -package=mock

package app

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

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
	logEncoding     string
	developmentLogs bool
	randomSeed      int64
)

func RandomSource() rand.Source {
	atomic.CompareAndSwapInt64(&randomSeed, 0, time.Now().UTC().UnixNano())
	return rand.NewSource(atomic.LoadInt64(&randomSeed))
}

type Spec struct {
	Name                    string
	Short                   string
	LogEncoding             string
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
			os.Exit(1)
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
	a.logger.Info("Shutdown initiated")
	a.cancel()
}

func NewApp(spec Spec) App {
	if spec.Defaults == nil {
		spec.Defaults = make(map[string]interface{})
	}

	a := &app{
		rootCmd: &cobra.Command{
			Use:          spec.Name,
			Short:        spec.Short,
			SilenceUsage: true,
		},
	}

	a.ctx, a.cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	lateInitTasks := []func(cmd *cobra.Command, args []string) (err error){
		func(*cobra.Command, []string) (err error) {
			return spec.readConfig(a.rootCmd)
		},
		func(cmd *cobra.Command, args []string) (err error) {
			var cwd string
			if cwd, err = os.Getwd(); err != nil {
				return err
			}
			err = logging.ConfigureLogging(
				logging.WithLevel(logging.ParseLevel(logLevel)),
				logging.WithDevelopment(developmentLogs),
				logging.WithEncoding(spec.LogEncoding),
				logging.WithInitialFields(map[string]interface{}{
					"cwd":  cwd,
					"cmd":  cmd.Name(),
					"args": args,
				}),
			)

			if err != nil {
				return err
			}

			a.logger = logging.CreateLogger()
			return
		},
		func(*cobra.Command, []string) (err error) {
			a.logger.Debug("Random seed", zap.Int64("seed", randomSeed))
			return nil
		},
	}

	lateInitTasks = append(lateInitTasks, spec.LateInitTasks...)

	a.rootCmd.AddCommand(spec.SubCommands...)

	a.rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file that should be used")
	a.rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level to use")
	a.rootCmd.PersistentFlags().StringVar(&logEncoding, "log-encoding", spec.LogEncoding, "Log encoding either 'json' or 'console'")
	a.rootCmd.PersistentFlags().BoolVar(&developmentLogs, "development-logs", false, "Enable development mode logs")
	a.rootCmd.PersistentFlags().Int64Var(
		&randomSeed,
		"random-seed",
		time.Now().UTC().UnixNano(),
		"Seed used for all random instances - defaults to UTC unix timestamp in nanoseconds when the application is started",
	)

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
