package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

var (
	cliCmd = &cobra.Command{
		Use:   "",
		Short: "IMCTL is the CLI app to interact with an INetMock server",
	}

	inetMockSocketPath string
	outputFormat       string
	grpcTimeout        time.Duration
)

func init() {
	cliCmd.PersistentFlags().StringVar(&inetMockSocketPath, "socket-path", "./inetmock.sock", "Path to the INetMock socket file")
	cliCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Output format to use. Possible values: table, json, yaml")
	cliCmd.PersistentFlags().DurationVar(&grpcTimeout, "grpc-timeout", 5*time.Second, "Timeout to connect to the gRPC API")

	cliCmd.AddCommand(endpointsCmd, handlerCmd, healthCmd)
	endpointsCmd.AddCommand(getEndpoints)
	handlerCmd.AddCommand(getHandlersCmd)
	healthCmd.AddCommand(containerHealthCmd)
}
