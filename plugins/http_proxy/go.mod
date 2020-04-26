module github.com/baez90/inetmock/plugins/http_proxy

go 1.14

require (
	github.com/baez90/inetmock v0.0.1
	github.com/spf13/viper v1.6.3
	go.uber.org/zap v1.15.0
	gopkg.in/elazarl/goproxy.v1 v1.0.0-20180725130230-947c36da3153
)

replace github.com/baez90/inetmock v0.0.1 => ../../
