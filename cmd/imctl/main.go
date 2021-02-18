package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/google/uuid"
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

	currentUser := ""
	if usr, err := user.Current(); err == nil {
		currentUser = usr.Username
	} else {
		currentUser = uuid.New().String()
	}

	hostname := "."
	if hn, err := os.Hostname(); err == nil {
		hostname = hn
	}

	watchEventsCmd.PersistentFlags().StringVar(&listenerName, "listener-name", fmt.Sprintf("%s\\%s is watching", hostname, currentUser), "set listener name - defaults to the current username, if the user cannot be determined a random UUID will be used")
	auditCmd.AddCommand(listSinksCmd, watchEventsCmd, addFileCmd, removeFileCmd, readFileCmd)
	pcapCmd.AddCommand(listAvailableDevicesCmd, listCurrentlyRecordingsCmd, addRecordingCmd, removeCurrentlyActiveRecording)

	cliApp.MustRun()
}

func initGRPCConnection() (err error) {
	dialCtx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	conn, err = grpc.DialContext(dialCtx, inetMockSocketPath, grpc.WithInsecure())
	cancel()

	return
}
