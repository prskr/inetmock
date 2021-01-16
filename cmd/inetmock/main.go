package main

import (
	"fmt"
	"os"

	"gitlab.com/inetmock/inetmock/internal/app"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	mock2 "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/proxy"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/proxy"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/metrics"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/metrics"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/tls/interceptor"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/tls/interceptor"
	"go.uber.org/zap"
)

var (
	server app.App
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	var err error
	if server, err = app.NewApp(
		mock2.AddHTTPMock,
		mock.AddDNSMock,
		interceptor.AddTLSInterceptor,
		proxy.AddHTTPProxy,
		metrics.AddMetricsExporter,
	); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	server.
		WithCommands(serveCmd, generateCaCmd).
		MustRun()
}
