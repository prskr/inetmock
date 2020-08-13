package main

import (
	"github.com/baez90/inetmock/internal/cmd"
	"go.uber.org/zap"
	"os"

	_ "github.com/baez90/inetmock/plugins/dns_mock"
	_ "github.com/baez90/inetmock/plugins/http_mock"
	_ "github.com/baez90/inetmock/plugins/http_proxy"
	_ "github.com/baez90/inetmock/plugins/metrics_exporter"
	_ "github.com/baez90/inetmock/plugins/tls_interceptor"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := cmd.ExecuteServerCommand(); err != nil {
		logger.Error("Failed to run inetmock",
			zap.Error(err),
		)
		os.Exit(1)
	}
}
