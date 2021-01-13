package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"google.golang.org/grpc"
)

var (
	watchEventsCmd = &cobra.Command{
		Use:   "watch",
		Short: "Watch all audit events",
		RunE:  watchAuditEvents,
	}

	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Interact with the audit stream",
	}

	listenerName string
)

func watchAuditEvents(_ *cobra.Command, _ []string) (err error) {
	var conn *grpc.ClientConn

	if conn, err = grpc.Dial(inetMockSocketPath, grpc.WithInsecure()); err != nil {
		fmt.Printf("Failed to connecto INetMock socket: %v\n", err)
		os.Exit(10)
	}

	auditClient := rpc.NewAuditClient(conn)

	var watchClient rpc.Audit_WatchEventsClient
	if watchClient, err = auditClient.WatchEvents(appCtx, &rpc.WatchEventsRequest{WatcherName: listenerName}); err != nil {
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

	<-appCtx.Done()
	err = watchClient.CloseSend()

	return
}
