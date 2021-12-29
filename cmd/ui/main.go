//go:build js && wasm
// +build js,wasm

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/Nerzal/tinydom"
	"github.com/tarndt/wasmws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	rpcv1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

var done = make(chan struct{})

func main() {
	window := tinydom.GetWindow()
	location := window.Location()
	log.Printf("Connecting WebSocket")
	grpcProxyUrl := fmt.Sprintf("passthrough:///ws://%s/grpc-proxy", location.Host())
	log.Println(grpcProxyUrl)

	dialCto, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCto,
		grpcProxyUrl,
		grpc.WithContextDialer(wasmws.GRPCDialer),
		grpc.WithDisableRetry(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	)

	if err != nil {
		log.Fatalf("Failed to dial gRPC websocket: %v", err)
	}

	defer conn.Close()

	pcapClient := rpcv1.NewPCAPServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if resp, err := pcapClient.ListAvailableDevices(ctx, new(rpcv1.ListAvailableDevicesRequest)); err != nil {
		log.Printf("Failed to list available devices: %v", err)
		return
	} else {
		log.Printf("Number of results: %d", len(resp.AvailableDevices))
		for _, device := range resp.AvailableDevices {
			log.Println(device.Name)
		}
	}

	auditClient := rpcv1.NewAuditServiceClient(conn)

	watchCtx, cancelWatch := context.WithCancel(context.Background())
	defer cancelWatch()
	stream, err := auditClient.WatchEvents(watchCtx, new(rpcv1.WatchEventsRequest))
	if err != nil {
		log.Printf("Failed to subscribe to events: %v", err)
	} else {
		go followStream(stream)
	}
	<-done
}

func followStream(stream rpcv1.AuditService_WatchEventsClient) {
	for {
		ev, err := stream.Recv()
		if err != nil {
			return
		}
		log.Println(ev.Entity.String())
	}
}
