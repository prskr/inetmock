package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	v1 "google.golang.org/grpc/health/grpc_health_v1"
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
		Run:          runGeneralHealth,
		SilenceUsage: true,
	}

	containerHealthCmd = &cobra.Command{
		Use:          "container",
		Short:        "get the health in a container compatible way i.e. exit code 0 if okay otherwise exit code 1",
		Run:          runContainerHealth,
		SilenceUsage: true,
	}
)

func getHealthResult() (healthResp *v1.HealthCheckResponse, err error) {
	healthClient := v1.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPCTimeout)
	healthResp, err = healthClient.Check(ctx, new(v1.HealthCheckRequest))
	cancel()
	return
}

func runGeneralHealth(_ *cobra.Command, _ []string) {
	var healthResp *v1.HealthCheckResponse
	var err error

	if healthResp, err = getHealthResult(); err != nil {
		fmt.Printf("Failed to get health information: %v", err)
		os.Exit(1)
	}

	fmt.Println(healthResp.Status.String())
	if healthResp.Status != v1.HealthCheckResponse_SERVING {
		os.Exit(1)
	}
}

func runContainerHealth(_ *cobra.Command, _ []string) {
	if healthResp, err := getHealthResult(); err != nil {
		fmt.Printf("Failed to get health information: %v", err)
		os.Exit(1)
	} else if healthResp.GetStatus() != v1.HealthCheckResponse_SERVING {
		fmt.Println("Overall health state is not healthy")
		os.Exit(1)
	}
}
