package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"gitlab.com/inetmock/inetmock/internal/app"
)

const (
	defaultGRPCTimeout = 5 * time.Second
)

var (
	inetMockSocketPath string
	outputFormat       string
	grpcTimeout        time.Duration
	cliApp             app.App
	conn               *grpc.ClientConn
)

//nolint:lll
func main() {
	healthCmd.AddCommand(generalHealthCmd, containerHealthCmd)

	cliApp = app.NewApp("imctl", "IMCTL is the CLI app to interact with an INetMock server").
		WithCommands(healthCmd, auditCmd, pcapCmd).
		WithInitTasks(func(_ *cobra.Command, _ []string) (err error) {
			return initGRPCConnection()
		}).
		WithLogger()

	cliApp.RootCommand().PersistentFlags().StringVar(&inetMockSocketPath, "socket-path", "unix:///var/run/inetmock.sock", "Path to the INetMock socket file")
	cliApp.RootCommand().PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Output format to use. Possible values: table, json, yaml")
	cliApp.RootCommand().PersistentFlags().DurationVar(&grpcTimeout, "grpc-timeout", defaultGRPCTimeout, "Timeout to connect to the gRPC API")

	cliApp.MustRun()
}

func initGRPCConnection() (err error) {
	dialCtx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	conn, err = grpc.DialContext(dialCtx, inetMockSocketPath, grpc.WithInsecure())
	cancel()

	return
}
