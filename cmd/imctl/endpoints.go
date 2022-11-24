package main

import (
	"context"

	"github.com/spf13/cobra"

	rpcv1 "inetmock.icb4dc0.de/inetmock/pkg/rpc/v1"
)

var (
	endpointsCmd = &cobra.Command{
		Use:   "endpoints",
		Short: "Manage the endpoint lifecycle of an INetMock instance",
	}

	restartEndpointsCmd = &cobra.Command{
		Use:          "restart",
		Short:        "Restart all endpoints",
		SilenceUsage: true,
		RunE: func(*cobra.Command, []string) error {
			return runRestartEndpoints()
		},
	}
)

func init() {
	endpointsCmd.AddCommand(restartEndpointsCmd)
}

func runRestartEndpoints() error {
	endpointsClient := rpcv1.NewEndpointOrchestratorServiceClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	defer cancel()

	_, err := endpointsClient.RestartAllGroups(ctx, new(rpcv1.RestartAllGroupsRequest))
	return err
}
