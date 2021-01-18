package main

import (
	"fmt"

	"gitlab.com/inetmock/inetmock/internal/app"
	dns "gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	http "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/proxy"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/metrics"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/tls/interceptor"
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

	app.NewApp("inetmock", "INetMock is lightweight internet mock").
		WithHandlerRegistry(
			http.AddHTTPMock,
			dns.AddDNSMock,
			interceptor.AddTLSInterceptor,
			proxy.AddHTTPProxy,
			metrics.AddMetricsExporter).
		WithCommands(serveCmd, generateCaCmd).
		WithConfig().
		WithLogger().
		WithHealthChecker().
		WithCertStore().
		WithEventStream().
		WithEndpointManager().
		MustRun()
}
