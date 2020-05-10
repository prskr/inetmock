module github.com/baez90/inetmock/plugins/dns_mock

go 1.14

require (
	github.com/baez90/inetmock v0.0.1
	github.com/miekg/dns v1.1.29
	github.com/spf13/viper v1.7.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71 // indirect
)

replace github.com/baez90/inetmock v0.0.1 => ../../
