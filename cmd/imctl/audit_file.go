package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"google.golang.org/grpc"
)

var (
	addFileCmd = &cobra.Command{
		Use:   "addFile",
		Short: "subscribe events to a file",
		Args:  cobra.ExactArgs(1),
		RunE:  runAddFile,
	}

	removeFileCmd = &cobra.Command{
		Use:   "removeFile",
		Short: "remove file subscription",
		Args:  cobra.ExactArgs(1),
		RunE:  runRemoveFile,
	}

	readFileCmd = &cobra.Command{
		Use:   "readFile",
		Short: "reads an audit file and prints the events",
		Args:  cobra.ExactArgs(1),
		RunE:  runReadFile,
	}
)

func runAddFile(_ *cobra.Command, args []string) (err error) {
	var conn *grpc.ClientConn

	if conn, err = grpc.Dial(inetMockSocketPath, grpc.WithInsecure()); err != nil {
		fmt.Printf("Failed to connecto INetMock socket: %v\n", err)
		os.Exit(10)
	}

	auditClient := rpc.NewAuditClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	_, err = auditClient.RegisterFileSink(ctx, &rpc.RegisterFileSinkRequest{TargetPath: args[0]})
	return
}

func runRemoveFile(_ *cobra.Command, args []string) (err error) {
	auditClient := rpc.NewAuditClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	_, err = auditClient.RemoveFileSink(ctx, &rpc.RemoveFileSinkRequest{TargetPath: args[0]})
	return
}

func runReadFile(_ *cobra.Command, args []string) (err error) {
	return
}
