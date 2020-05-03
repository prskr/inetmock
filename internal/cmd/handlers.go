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
	getHandlersCmd = &cobra.Command{
		Use:   "get",
		Short: "Get all registered handlers",
		Run:   runGetHandlers,
	}

	handlerCmd = &cobra.Command{
		Use:     "handlers",
		Short:   "handlers is the entrypoint to all kind of commands to interact with handlers",
		Aliases: []string{"handler"},
	}
)

type printableHandler struct {
	Handler string
}

func fromHandlers(hs []string) (handlers []*printableHandler) {
	for idx := range hs {
		handlers = append(handlers, &printableHandler{
			Handler: hs[idx],
		})
	}
	return
}

func runGetHandlers(_ *cobra.Command, _ []string) {
	var err error
	var conn *grpc.ClientConn

	if conn, err = grpc.Dial(inetMockSocketPath, grpc.WithInsecure()); err != nil {
		fmt.Printf("Failed to connecto INetMock socket: %v\n", err)
		os.Exit(10)
	}
	handlersClient := rpc.NewHandlersClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), grpcTimeout)
	var handlersResp *rpc.GetHandlersResponse

	if handlersResp, err = handlersClient.GetHandlers(ctx, &rpc.GetHandlersRequest{}); err != nil {
		fmt.Printf("Failed to get the endpoints: %v", err)
		os.Exit(11)
	}

	writer := format.Writer(outputFormat, os.Stdout)
	if err = writer.Write(fromHandlers(handlersResp.Handlers)); err != nil {
		fmt.Printf("Error occurred during writing response values: %v\n", err)
	}
}
