package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/format"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
)

var (
	listSinksCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "dir"},
		Short:   "List all subscribed sinks",
		RunE:    runListSinks,
	}
	addFileCmd = &cobra.Command{
		Use:     "add-file",
		Aliases: []string{"add"},
		Short:   "subscribe events to a file",
		Args:    cobra.ExactArgs(1),
		RunE:    runAddFile,
	}

	removeFileCmd = &cobra.Command{
		Use:     "remove-file",
		Aliases: []string{"rm", "del"},
		Short:   "remove file subscription",
		Args:    cobra.ExactArgs(1),
		RunE:    runRemoveFile,
	}

	readFileCmd = &cobra.Command{
		Use:     "read-file",
		Aliases: []string{"cat"},
		Short:   "reads an audit file and prints the events",
		Args:    cobra.ExactArgs(1),
		RunE:    runReadFile,
	}
)

func runListSinks(*cobra.Command, []string) (err error) {
	auditClient := rpc.NewAuditClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	var resp *rpc.ListSinksResponse
	if resp, err = auditClient.ListSinks(ctx, new(rpc.ListSinksRequest)); err != nil {
		return
	}

	type printableSink struct {
		Name string
	}

	var sinks []printableSink
	for _, s := range resp.Sinks {
		sinks = append(sinks, printableSink{Name: s})
	}

	writer := format.Writer(outputFormat, os.Stdout)
	err = writer.Write(sinks)
	return
}

func runAddFile(_ *cobra.Command, args []string) (err error) {
	auditClient := rpc.NewAuditClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), grpcTimeout)
	defer cancel()

	var resp *rpc.RegisterFileSinkResponse
	resp, err = auditClient.RegisterFileSink(ctx, &rpc.RegisterFileSinkRequest{TargetPath: args[0]})

	if err != nil {
		return
	}

	cliApp.Logger().Info("Successfully registered file sink", zap.String("targetPath", resp.ResolvedPath))

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
	if len(args) != 1 {
		return errors.New("expected only 1 argument")
	}

	var reader io.ReadCloser
	if reader, err = os.Open(args[0]); err != nil {
		return
	}

	eventReader := audit.NewEventReader(reader)
	var ev audit.Event

	for err == nil {
		if ev, err = eventReader.Read(); err == nil {
			var jsonBytes []byte
			if jsonBytes, err = json.Marshal(ev); err == nil {
				fmt.Println(string(jsonBytes))
			}
		}
	}

	if errors.Is(err, io.EOF) {
		err = nil
	}

	return
}
