package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	"gitlab.com/inetmock/inetmock/internal/app"
)

const (
	defaultGRPCTimeout = 5 * time.Second
)

var (
	cliApp app.App
	conn   *grpc.ClientConn
	cfg    config
)

type config struct {
	SocketPath  string
	Format      string
	GRPCTimeout time.Duration
}

func main() {
	healthCmd.AddCommand(generalHealthCmd, containerHealthCmd)
	cliApp = app.NewApp(
		app.Spec{
			Name:        "imctl",
			Short:       "IMCTL is the CLI app to interact with an INetMock server",
			LogEncoding: "console",
			Config:      &cfg,
			SubCommands: []*cobra.Command{healthCmd, auditCmd, pcapCmd, checkCmd},
			LateInitTasks: []func(cmd *cobra.Command, args []string) (err error){
				initGRPCConnection,
			},
			IgnoreMissingConfigFile: true,
			FlagBindings: map[string]func(flagSet *pflag.FlagSet) *pflag.Flag{
				"grpctimeout": func(flagSet *pflag.FlagSet) *pflag.Flag {
					return flagSet.Lookup("grpc-timeout")
				},
				"format": func(flagSet *pflag.FlagSet) *pflag.Flag {
					return flagSet.Lookup("format")
				},
				"socketpath": func(flagSet *pflag.FlagSet) *pflag.Flag {
					return flagSet.Lookup("socket-path")
				},
			},
		},
	)

	cliApp.RootCommand().PersistentFlags().String("socket-path", "unix:///var/run/inetmock/inetmock.sock", "Path to the INetMock socket file")
	cliApp.RootCommand().PersistentFlags().StringP("format", "f", "table", "Output format to use. Possible values: table, json, yaml")
	cliApp.RootCommand().PersistentFlags().Duration("grpc-timeout", defaultGRPCTimeout, "Timeout to connect to the gRPC API")

	cliApp.MustRun()
}

func initGRPCConnection(*cobra.Command, []string) (err error) {
	dialCtx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	conn, err = grpc.DialContext(dialCtx, cfg.SocketPath, grpc.WithInsecure())
	cancel()

	return
}
