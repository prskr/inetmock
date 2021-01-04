package main

import (
	"gitlab.com/inetmock/inetmock/internal/cmd"
	_ "gitlab.com/inetmock/inetmock/plugins/dns_mock"
	_ "gitlab.com/inetmock/inetmock/plugins/http_mock"
	_ "gitlab.com/inetmock/inetmock/plugins/http_proxy"
	_ "gitlab.com/inetmock/inetmock/plugins/metrics_exporter"
	_ "gitlab.com/inetmock/inetmock/plugins/tls_interceptor"
)

func main() {
	cmd.ExecuteServerCommand()
}
