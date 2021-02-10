package mock

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"gitlab.com/inetmock/inetmock/internal/test/integration"
)

const (
	charSet = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Benchmark_dnsHandler(b *testing.B) {
	var err error
	var endpoint string
	if endpoint, err = setupContainer(b, "53/udp"); err != nil {
		b.Errorf("setupContainer() error = %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		resolv := resolver(endpoint)
		for pb.Next() {
			lookupCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			_, err := resolv.LookupHost(lookupCtx, fmt.Sprintf("www.%s.com", randomString(8)))
			cancel()
			if err != nil {
				b.Errorf("LookupHost() error = %v", err)
			}
		}
	})

}

func randomString(length int) (result string) {
	buffer := strings.Builder{}
	for i := 0; i < length; i++ {
		buffer.WriteByte(charSet[rand.Intn(len(charSet))])
	}
	return buffer.String()
}

func setupContainer(b *testing.B, port string) (httpEndpoint string, err error) {
	b.Helper()

	startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var inetMockContainer testcontainers.Container
	if inetMockContainer, err = integration.SetupINetMockContainer(startupCtx, b, port); err != nil {
		return
	}

	httpEndpoint, err = inetMockContainer.PortEndpoint(startupCtx, nat.Port(port), "")
	return
}

func resolver(endpoint string) net.Resolver {
	return net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (conn net.Conn, err error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, "udp", endpoint)
		},
	}
}
