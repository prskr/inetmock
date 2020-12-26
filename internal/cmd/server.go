package cmd

import (
	"fmt"
	"os"

	"gitlab.com/inetmock/inetmock/internal/app"
	"gitlab.com/inetmock/inetmock/plugins/dns_mock"
	"gitlab.com/inetmock/inetmock/plugins/http_mock"
	"gitlab.com/inetmock/inetmock/plugins/http_proxy"
	"gitlab.com/inetmock/inetmock/plugins/metrics_exporter"
	"gitlab.com/inetmock/inetmock/plugins/tls_interceptor"
)

var (
	server app.App
)

func ExecuteServerCommand() {
	var err error
	if server, err = app.NewApp(
		http_mock.AddHTTPMock,
		dns_mock.AddDNSMock,
		tls_interceptor.AddTLSInterceptor,
		http_proxy.AddHTTPProxy,
		metrics_exporter.AddMetricsExporter,
	); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	server.
		WithCommands(serveCmd, generateCaCmd).
		MustRun()
}
