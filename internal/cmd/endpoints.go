package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/format"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"google.golang.org/grpc"
)

var (
	getEndpoints = &cobra.Command{
		Use:   "get",
		Short: "Get all running endpoints",
		RunE:  runGetEndpoints,
	}

	endpointsCmd = &cobra.Command{
		Use:     "endpoints",
		Short:   "endpoints is the entrypoint to all kind of commands to interact with endpoints",
		Aliases: []string{"ep", "endpoint"},
	}
)

type printableEndpoint struct {
	Id            string
	Name          string
	Handler       string
	ListenAddress string
	Port          int
}

func fromEndpoint(ep *rpc.Endpoint) *printableEndpoint {
	return &printableEndpoint{
		Id:            ep.Id,
		Name:          ep.Name,
		Handler:       ep.Handler,
		ListenAddress: ep.ListenAddress,
		Port:          int(ep.Port),
	}
}

func fromEndpoints(eps []*rpc.Endpoint) (out []*printableEndpoint) {
	for idx := range eps {
		out = append(out, fromEndpoint(eps[idx]))
	}
	return
}

func runGetEndpoints(_ *cobra.Command, _ []string) (err error) {
	var conn *grpc.ClientConn

	if conn, err = grpc.Dial(inetMockSocketPath, grpc.WithInsecure()); err != nil {
		fmt.Printf("Failed to connecto INetMock socket: %v\n", err)
		os.Exit(10)
	}
	endpointsClient := rpc.NewEndpointsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	var endpointsResp *rpc.GetEndpointsResponse
	if endpointsResp, err = endpointsClient.GetEndpoints(ctx, &rpc.GetEndpointsRequest{}); err != nil {
		fmt.Printf("Failed to get the endpoints: %v", err)
		os.Exit(11)
	}

	writer := format.Writer(outputFormat, os.Stdout)
	if err = writer.Write(fromEndpoints(endpointsResp.Endpoints)); err != nil {
		fmt.Printf("Error occurred during writing response values: %v\n", err)
	}
	return
}
