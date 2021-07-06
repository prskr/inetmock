// +build integration
// go:build integration

package mock_test

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"

	"gitlab.com/inetmock/inetmock/internal/test/integration"
)

const (
	charSet         = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	startupTimeout  = 5 * time.Minute
	shutdownTimeout = 5 * time.Second
)

var (
	dnsEndpoint string
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	var (
		code              int
		inetMockContainer testcontainers.Container
		port              = nat.Port("53/udp")
		err               error
		errorHandler      = func(err error) bool {
			if err != nil {
				fmt.Println(err.Error())
				code = 1
				return true
			}
			return false
		}
		terminate = func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			errorHandler(inetMockContainer.Terminate(shutdownCtx))
			cancel()
		}
	)

	defer func() {
		terminate()
		os.Exit(code)
	}()

	startupCtx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	inetMockContainer, err = integration.SetupINetMockContainer(startupCtx, string(port))
	errorHandler(err)
	dnsEndpoint, err = inetMockContainer.PortEndpoint(startupCtx, port, "")
	errorHandler(err)
	cancel()

	if code != 0 {
		return
	}

	code = m.Run()
}

func Benchmark_dnsHandler(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		resolv := resolver(dnsEndpoint)
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
		//nolint:gosec
		buffer.WriteByte(charSet[rand.Intn(len(charSet))])
	}
	return buffer.String()
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
