//go:build integration
// +build integration

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
	"os"
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
	charSet         = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	startupTimeout  = 5 * time.Minute
	shutdownTimeout = 5 * time.Second
)

var (
	availableExtensions = []string{"gif", "html", "ico", "jpg", "png", "txt"}
	httpEndpoint        string
	httpsEndpoint       string
	defaultURLGenerator = func(endpoint string) *url.URL {
		//nolint:gosec
		extension := availableExtensions[rand.Intn(len(availableExtensions))]
		reqURL, _ := url.Parse(fmt.Sprintf("%s/%s.%s", endpoint, randomString(15), extension))
		return reqURL
	}
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	var (
		code              int
		inetMockContainer testcontainers.Container
		httpPort          = nat.Port("80/tcp")
		httpsPort         = nat.Port("443/tcp")
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
	httpEndpoint, err = inetMockContainer.PortEndpoint(startupCtx, httpPort, "http")
	errorHandler(err)
	httpsEndpoint, err = inetMockContainer.PortEndpoint(startupCtx, httpsPort, "https")
	errorHandler(err)
	cancel()

	if code != 0 {
		return
	}

	code = m.Run()
}

func Benchmark_httpHandler(b *testing.B) {
	type benchmark struct {
		name         string
		endpoint     string
		urlGenerator func(endpoint string) *url.URL
	}
	benchmarks := []benchmark{
		{
			name:         "HTTP",
			endpoint:     httpEndpoint,
			urlGenerator: defaultURLGenerator,
		},
		{
			name:     "HTTP - ensure /index.html is handled correctly",
			endpoint: httpEndpoint,
			urlGenerator: func(endpoint string) *url.URL {
				reqURL, _ := url.Parse(fmt.Sprintf("%s/index.html", endpoint))
				return reqURL
			},
		},
		{
			name:         "HTTPS",
			endpoint:     httpsEndpoint,
			urlGenerator: defaultURLGenerator,
		},
	}
	for _, bc := range benchmarks {
		bm := bc
		b.Run(bm.name, func(b *testing.B) {
			var (
				err        error
				httpClient *http.Client
			)

			if httpClient, err = setupHTTPClient(); err != nil {
				return
			}

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					req := &http.Request{
						Method: http.MethodGet,
						URL:    bc.urlGenerator(bm.endpoint),
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
	client := &http.Client{
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
