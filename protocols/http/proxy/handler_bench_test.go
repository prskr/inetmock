//go:build integration

package proxy_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"

	"inetmock.icb4dc0.de/inetmock/internal/test"
	"inetmock.icb4dc0.de/inetmock/internal/test/integration"
)

const (
	startupTimeout  = 10 * time.Minute
	shutdownTimeout = 5 * time.Second
)

var (
	availableExtensions = []string{"gif", "html", "ico", "jpg", "png", "txt"}
	proxyHTTPEndpoint   string
	proxyHTTPSEndpoint  string
)

func TestMain(m *testing.M) {
	var (
		inetMockContainer testcontainers.Container
		httpPort          = nat.Port("3128/tcp")
		httpsPort         = nat.Port("3128/tcp")
		code              int
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
	inetMockContainer, err = integration.SetupINetMockContainer(startupCtx, string(httpPort), string(httpsPort))
	errorHandler(err)
	proxyHTTPEndpoint, err = inetMockContainer.PortEndpoint(startupCtx, httpPort, "http")
	errorHandler(err)
	proxyHTTPSEndpoint, err = inetMockContainer.PortEndpoint(startupCtx, httpsPort, "https")
	errorHandler(err)
	cancel()

	if code != 0 {
		return
	}

	code = m.Run()
}

func Benchmark_httpProxy(b *testing.B) {
	type benchmark struct {
		name     string
		endpoint string
	}
	benchmarks := []benchmark{
		{
			name:     "HTTP",
			endpoint: proxyHTTPEndpoint,
		},
		{
			name:     "HTTPS",
			endpoint: proxyHTTPSEndpoint,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			//nolint:gosec // pseudo-random is good enough for tests
			random := rand.New(rand.NewSource(time.Now().Unix()))
			var err error

			var httpClient *http.Client
			if httpClient, err = setupHTTPClient(proxyHTTPEndpoint, proxyHTTPSEndpoint); err != nil {
				return
			}

			time.Sleep(500 * time.Millisecond)

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					extension := availableExtensions[random.Intn(len(availableExtensions))]
					reqURL, _ := url.Parse(fmt.Sprintf("%s/%s.%s", bm.endpoint, test.RandomString(random, 15), extension))
					req := &http.Request{
						Method: http.MethodGet,
						URL:    reqURL,
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

func setupHTTPClient(httpEndpoint, httpsEndpoint string) (*http.Client, error) {
	var (
		repoRoot string
		err      error
	)

	if repoRoot, err = test.DiscoverRepoRoot(); err != nil {
		return nil, err
	}

	var demoCABytes []byte
	if demoCABytes, err = os.ReadFile(filepath.Join(repoRoot, "assets", "demoCA", "ca.pem")); err != nil {
		return nil, err
	}

	rootCaPool := x509.NewCertPool()
	if !rootCaPool.AppendCertsFromPEM(demoCABytes) {
		return nil, errors.New("failed to add CA key")
	}

	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				switch req.URL.Scheme {
				case "http":
					return url.Parse(httpEndpoint)
				case "https":
					return url.Parse(httpsEndpoint)
				default:
					return nil, errors.New("unknown scheme")
				}
			},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, addr)
			},
			//nolint:gosec
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
