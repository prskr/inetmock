package mock_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
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

var (
	availableExtensions = []string{"gif", "html", "ico", "jpg", "png", "txt"}
	defaultUrlGenerator = func(endpoint string) *url.URL {
		extension := availableExtensions[rand.Intn(len(availableExtensions))]
		reqURL, _ := url.Parse(fmt.Sprintf("%s/%s.%s", endpoint, randomString(15), extension))
		return reqURL
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Benchmark_httpHandler(b *testing.B) {
	type benchmark struct {
		name         string
		port         string
		scheme       string
		urlGenerator func(endpoint string) *url.URL
	}
	benchmarks := []benchmark{
		{
			name:         "HTTP",
			port:         "80/tcp",
			scheme:       "http",
			urlGenerator: defaultUrlGenerator,
		},
		{
			name:   "HTTP - ensure /index.html is handled correctly",
			port:   "8080/tcp",
			scheme: "http",
			urlGenerator: func(endpoint string) *url.URL {
				reqURL, _ := url.Parse(fmt.Sprintf("%s/index.html", endpoint))
				return reqURL
			},
		},
		{
			name:         "HTTPS",
			port:         "443/tcp",
			scheme:       "https",
			urlGenerator: defaultUrlGenerator,
		},
	}
	for _, bc := range benchmarks {
		bm := bc
		b.Run(bm.name, func(b *testing.B) {
			var err error
			var endpoint string
			if endpoint, err = setupContainer(b, bm.scheme, bm.port); err != nil {
				b.Fatalf("setupContainer() error = %v", err)
			}

			var httpClient *http.Client
			if httpClient, err = setupHTTPClient(); err != nil {
				return
			}

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					req := &http.Request{
						Method: http.MethodGet,
						URL:    bc.urlGenerator(endpoint),
						Close:  false,
						Host:   "www.inetmock.com",
					}
					if resp, err := httpClient.Do(req); err != nil {
						b.Error(err)
					} else if resp.StatusCode != 200 {
						b.Errorf("Got status code %d", resp.StatusCode)
					}
				}
			})
		})
	}
}

func randomString(length int) (result string) {
	buffer := strings.Builder{}
	for i := 0; i < length; i++ {
		//nolint:gosec
		buffer.WriteByte(charSet[rand.Intn(len(charSet))])
	}
	return buffer.String()
}

func setupContainer(b *testing.B, scheme, port string) (httpEndpoint string, err error) {
	b.Helper()

	startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var inetMockContainer testcontainers.Container
	if inetMockContainer, err = integration.SetupINetMockContainer(startupCtx, b, port); err != nil {
		return
	}

	httpEndpoint, err = inetMockContainer.PortEndpoint(startupCtx, nat.Port(port), scheme)
	return
}

func setupHTTPClient() (*http.Client, error) {
	//nolint:dogsled
	_, fileName, _, _ := runtime.Caller(0)

	var err error
	var repoRoot string
	if repoRoot, err = filepath.Abs(filepath.Join(filepath.Dir(fileName), "..", "..", "..", "..", "..")); err != nil {
		return nil, err
	}

	var demoCABytes []byte
	if demoCABytes, err = ioutil.ReadFile(filepath.Join(repoRoot, "assets", "demoCA", "ca.pem")); err != nil {
		return nil, err
	}

	rootCaPool := x509.NewCertPool()
	if !rootCaPool.AppendCertsFromPEM(demoCABytes) {
		return nil, errors.New("failed to add CA key")
	}

	//nolint:gosec
	var client = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				RootCAs: rootCaPool,
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return client, nil
}
