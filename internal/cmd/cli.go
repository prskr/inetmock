package cmd

import (
	"context"
	"os"
	"os/signal"
	"os/user"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	cliCmd = &cobra.Command{
		Use:   "",
		Short: "IMCTL is the CLI app to interact with an INetMock server",
	}

	inetMockSocketPath string
	outputFormat       string
	grpcTimeout        time.Duration
	appCtx             context.Context
	appCancel          context.CancelFunc
)

func init() {
	cliCmd.PersistentFlags().StringVar(&inetMockSocketPath, "socket-path", "unix:///var/run/inetmock.sock", "Path to the INetMock socket file")
	cliCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Output format to use. Possible values: table, json, yaml")
	cliCmd.PersistentFlags().DurationVar(&grpcTimeout, "grpc-timeout", 5*time.Second, "Timeout to connect to the gRPC API")

	cliCmd.AddCommand(endpointsCmd, handlerCmd, healthCmd, auditCmd)
	endpointsCmd.AddCommand(getEndpoints)
	handlerCmd.AddCommand(getHandlersCmd)
	healthCmd.AddCommand(generalHealthCmd)
	healthCmd.AddCommand(containerHealthCmd)

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

	appCtx, appCancel = initAppContext()
}

func ExecuteClientCommand() error {
	defer appCancel()
	return cliCmd.Execute()
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
