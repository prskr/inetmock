package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/format"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	rpcV1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

var (
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Interact with the audit API",
	}
	watchEventsCmd = &cobra.Command{
		Use:   "watch",
		Short: "Watch all audit events",
		RunE:  watchAuditEvents,
	}
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

	listenerName string
)

func init() {
	currentUser := ""
	if usr, err := user.Current(); err == nil {
		currentUser = usr.Username
	} else {
		currentUser = uuid.New().String()
	}

	hostname := "."
	if hn, err := os.Hostname(); err == nil {
		hostname = hn
	}

	watchEventsCmd.PersistentFlags().StringVar(
		&listenerName,
		"listener-name",
		fmt.Sprintf("%s\\%s is watching", hostname, currentUser),
		"set listener name - defaults to the current username, if the user cannot be determined a random UUID will be used",
	)
	auditCmd.AddCommand(listSinksCmd, watchEventsCmd, addFileCmd, removeFileCmd, readFileCmd)
}

func watchAuditEvents(_ *cobra.Command, _ []string) (err error) {
	auditClient := rpcV1.NewAuditServiceClient(conn)

	var watchClient rpcV1.AuditService_WatchEventsClient
	if watchClient, err = auditClient.WatchEvents(cliApp.Context(), &rpcV1.WatchEventsRequest{WatcherName: listenerName}); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go func() {
		var resp *rpcV1.WatchEventsResponse
		for resp, err = watchClient.Recv(); err == nil; resp, err = watchClient.Recv() {
			ev := audit.NewEventFromProto(resp.Entity)
			var out []byte
			out, err = json.Marshal(ev)
			if err != nil {
				continue
			}
			fmt.Println(string(out))
		}
	}()

	<-cliApp.Context().Done()
	err = watchClient.CloseSend()

	return
}

func runListSinks(*cobra.Command, []string) (err error) {
	auditClient := rpcV1.NewAuditServiceClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	defer cancel()

	var resp *rpcV1.ListSinksResponse
	if resp, err = auditClient.ListSinks(ctx, new(rpcV1.ListSinksRequest)); err != nil {
		return
	}

	type printableSink struct {
		Name string
	}

	var sinks = make([]printableSink, len(resp.Sinks))
	for i, s := range resp.Sinks {
		sinks[i] = printableSink{Name: s}
	}

	writer := format.Writer(cfg.Format, os.Stdout)
	err = writer.Write(sinks)
	return
}

func runAddFile(_ *cobra.Command, args []string) (err error) {
	auditClient := rpcV1.NewAuditServiceClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	defer cancel()

	var resp *rpcV1.RegisterFileSinkResponse
	resp, err = auditClient.RegisterFileSink(ctx, &rpcV1.RegisterFileSinkRequest{TargetPath: args[0]})

	if err != nil {
		return
	}

	cliApp.Logger().Info("Successfully registered file sink", zap.String("targetPath", resp.ResolvedPath))

	return
}

func runRemoveFile(_ *cobra.Command, args []string) (err error) {
	auditClient := rpcV1.NewAuditServiceClient(conn)
	ctx, cancel := context.WithTimeout(cliApp.Context(), cfg.GRPCTimeout)
	defer cancel()

	_, err = auditClient.RemoveFileSink(ctx, &rpcV1.RemoveFileSinkRequest{TargetPath: args[0]})
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
