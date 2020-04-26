module github.com/baez90/inetmock/plugins/tls_interceptor

go 1.14

require (
	github.com/baez90/inetmock v0.0.1
	github.com/google/uuid v1.1.1
	github.com/spf13/viper v1.6.3
	go.uber.org/zap v1.15.0
)

replace github.com/baez90/inetmock v0.0.1 => ../../
