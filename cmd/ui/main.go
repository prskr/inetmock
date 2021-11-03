package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

//nolint // playground
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// TODO replace
	host := "localhost"
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("ws://%s/grpc-socket", host))
	if err != nil {
		cancel()
		log.Fatalf("Failed to dial gRPC websocket: %v", err)
	}

	fmt.Printf("Conn target: %s", conn.Target())
}
