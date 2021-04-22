package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/inetmock/inetmock/internal/format"
	rpcV1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
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

func fromComponentsHealth(componentsHealth map[string]*rpcV1.ComponentHealth) interface{} {
	type printableHealthInfo struct {
		Component string
		State     string
		Message   string
	}

	var componentsInfo = make([]printableHealthInfo, len(componentsHealth))
	var idx int
	for componentName, component := range componentsHealth {
		componentsInfo[idx] = printableHealthInfo{
			Component: componentName,
			State:     component.State.String(),
			Message:   component.Message,
		}
		idx++
	}
	return componentsInfo
}

func getHealthResult() (healthResp *rpcV1.GetHealthResponse, err error) {
	var healthClient = rpcV1.NewHealthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPCTimeout)
	healthResp, err = healthClient.GetHealth(ctx, &rpcV1.GetHealthRequest{})
	cancel()
	return
}

func runGeneralHealth(_ *cobra.Command, _ []string) {
	var healthResp *rpcV1.GetHealthResponse
	var err error

	if healthResp, err = getHealthResult(); err != nil {
		fmt.Printf("Failed to get health information: %v", err)
		os.Exit(1)
	}

	printable := fromComponentsHealth(healthResp.ComponentsHealth)

	writer := format.Writer(cfg.Format, os.Stdout)
	if err = writer.Write(printable); err != nil {
		fmt.Printf("Error occurred during writing response values: %v\n", err)
	}
}

func runContainerHealth(_ *cobra.Command, _ []string) {
	if healthResp, err := getHealthResult(); err != nil {
		fmt.Printf("Failed to get health information: %v", err)
		os.Exit(1)
	} else if healthResp.OverallHealthState != rpcV1.HealthState_HEALTH_STATE_HEALTHY {
		fmt.Println("Overall health state is not healthy")
		os.Exit(1)
	}
}
