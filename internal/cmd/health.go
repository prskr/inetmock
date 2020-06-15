package cmd

import (
	"context"
	"fmt"
	"github.com/baez90/inetmock/internal/format"
	"github.com/baez90/inetmock/internal/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
)

var (
	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "health is the entry point for all health check related commands",
	}

	generalHealthCmd = &cobra.Command{
		Use:   "general",
		Short: "get the health in a more general way i.e. exit code 0 if healthy, exit codes unequal 0 if somethings wrong",
		Long: `
Exit code 1 means the server is still initializing
Exit code 2 means any component is unhealthy
Exit code 10 means an error occurred while opening a connection to the API socket 

The output contains information about each component and it's health state.
`,
		Run: runGeneralHealth,
	}

	containerHealthCmd = &cobra.Command{
		Use:   "container",
		Short: "get the health in a container compatible way i.e. exit code 0 if okay otherwise exit code 1",
		Run:   runContainerHealth,
	}
)

type printableHealthInfo struct {
}

func runGeneralHealth(_ *cobra.Command, _ []string) {
	var err error
	var conn *grpc.ClientConn

	if conn, err = grpc.Dial(inetMockSocketPath, grpc.WithInsecure()); err != nil {
		fmt.Printf("Failed to connecto INetMock socket: %v\n", err)
		os.Exit(10)
	}

	var healthClient = rpc.NewHealthClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), grpcTimeout)
	var healthResp *rpc.HealthResponse

	if healthResp, err = healthClient.GetHealth(ctx, &rpc.HealthRequest{}); err != nil {
		fmt.Printf("Failed to get health information: %v", err)
		os.Exit(1)
	}

	writer := format.Writer(outputFormat, os.Stdout)
	if err = writer.Write(healthResp); err != nil {
		fmt.Printf("Error occurred during writing response values: %v\n", err)
	}
}

func runContainerHealth(_ *cobra.Command, _ []string) {

}
