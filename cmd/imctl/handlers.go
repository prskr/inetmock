package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/format"
	"gitlab.com/inetmock/inetmock/internal/rpc"
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
	handlersClient := rpc.NewHandlersClient(conn)

	ctx, cancel := context.WithTimeout(appCtx, grpcTimeout)
	defer cancel()
	var err error
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
