package main

import (
	"gitlab.com/inetmock/inetmock/internal/app"
	dns "gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	http "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/proxy"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/metrics"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/tls/interceptor"
)

var (
	serverApp app.App
)

func main() {
	serverApp = app.NewApp("inetmock", "INetMock is lightweight internet mock").
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
		WithEndpointManager()

	serverApp.MustRun()
}
