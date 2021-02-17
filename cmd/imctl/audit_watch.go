package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
)

var (
	watchEventsCmd = &cobra.Command{
		Use:   "watch",
		Short: "Watch all audit events",
		RunE:  watchAuditEvents,
	}

	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Interact with the audit API",
	}

	listenerName string
)

func watchAuditEvents(_ *cobra.Command, _ []string) (err error) {
	auditClient := rpc.NewAuditClient(conn)

	var watchClient rpc.Audit_WatchEventsClient
	if watchClient, err = auditClient.WatchEvents(cliApp.Context(), &rpc.WatchEventsRequest{WatcherName: listenerName}); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go func() {
		var protoEv *audit.EventEntity
		for protoEv, err = watchClient.Recv(); err == nil; protoEv, err = watchClient.Recv() {
			ev := audit.NewEventFromProto(protoEv)
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
