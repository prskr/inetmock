package cmd

import (
	"fmt"
	"os"

	"gitlab.com/inetmock/inetmock/internal/app"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	mock2 "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/proxy"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/metrics"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/tls/interceptor"
)

var (
	server app.App
)

func ExecuteServerCommand() {
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
