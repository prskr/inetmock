package proxy_test

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
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Benchmark_httpProxy(b *testing.B) {
	type benchmark struct {
		name   string
		port   string
		scheme string
	}
	benchmarks := []benchmark{
		{
			name:   "HTTP",
			port:   "3128/tcp",
			scheme: "http",
		},
		{
			name:   "HTTPS",
			port:   "3128/tcp",
			scheme: "https",
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var err error
			var endpoint string
			if endpoint, err = setupContainer(b, bm.port); err != nil {
				b.Errorf("setupContainer() error = %v", err)
			}

			var httpClient *http.Client
			if httpClient, err = setupHTTPClient(fmt.Sprintf("http://%s", endpoint), fmt.Sprintf("https://%s", endpoint)); err != nil {
				return
			}

			time.Sleep(500 * time.Millisecond)

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					//nolint:gosec
					extension := availableExtensions[rand.Intn(len(availableExtensions))]

					reqURL, _ := url.Parse(fmt.Sprintf("%s://%s/%s.%s", bm.scheme, endpoint, randomString(15), extension))

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

func randomString(length int) (result string) {
	buffer := strings.Builder{}
	for i := 0; i < length; i++ {
		//nolint:gosec
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

func setupHTTPClient(httpEndpoint, httpsEndpoint string) (*http.Client, error) {
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

	var client = &http.Client{
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
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
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
