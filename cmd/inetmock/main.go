package main

import (
	"fmt"

	"gitlab.com/inetmock/inetmock/internal/cmd"
	_ "gitlab.com/inetmock/inetmock/plugins/dns_mock"
	_ "gitlab.com/inetmock/inetmock/plugins/http_mock"
	_ "gitlab.com/inetmock/inetmock/plugins/http_proxy"
	_ "gitlab.com/inetmock/inetmock/plugins/metrics_exporter"
	_ "gitlab.com/inetmock/inetmock/plugins/tls_interceptor"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	cmd.ExecuteServerCommand()
}
