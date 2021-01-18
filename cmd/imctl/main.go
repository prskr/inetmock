package main

import (
	"context"
	"os/user"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/app"
	"google.golang.org/grpc"
)

var (
	inetMockSocketPath string
	outputFormat       string
	grpcTimeout        time.Duration
	cliApp             app.App
	conn               *grpc.ClientConn
)

func main() {

	endpointsCmd.AddCommand(getEndpoints)
	handlerCmd.AddCommand(getHandlersCmd)
	healthCmd.AddCommand(generalHealthCmd, containerHealthCmd)

	cliApp = app.NewApp("imctl", "IMCTL is the CLI app to interact with an INetMock server").
		WithCommands(endpointsCmd, handlerCmd, healthCmd, auditCmd).
		WithInitTasks(func(_ *cobra.Command, _ []string) (err error) {
			return initGRPCConnection()
		})

	cliApp.RootCommand().PersistentFlags().StringVar(&inetMockSocketPath, "socket-path", "unix:///var/run/inetmock.sock", "Path to the INetMock socket file")
	cliApp.RootCommand().PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Output format to use. Possible values: table, json, yaml")
	cliApp.RootCommand().PersistentFlags().DurationVar(&grpcTimeout, "grpc-timeout", 5*time.Second, "Timeout to connect to the gRPC API")

	currentUser := ""
	if usr, err := user.Current(); err == nil {
		currentUser = usr.Username
	} else {
		currentUser = uuid.New().String()
	}

	watchEventsCmd.PersistentFlags().StringVar(
		&listenerName,
		"listener-name",
		currentUser,
		"set listener name - defaults to the current username, if the user cannot be determined a random UUID will be used",
	)
	auditCmd.AddCommand(watchEventsCmd, addFileCmd, removeFileCmd)

	cliApp.MustRun()

}

func initGRPCConnection() (err error) {
	dialCtx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	conn, err = grpc.DialContext(dialCtx, inetMockSocketPath, grpc.WithInsecure())
	cancel()

	return
}
