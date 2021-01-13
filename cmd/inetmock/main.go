package main

import (
	"gitlab.com/inetmock/inetmock/internal/cmd"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/proxy"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/metrics"
	_ "gitlab.com/inetmock/inetmock/internal/endpoint/handler/tls/interceptor"
)

func main() {
	cmd.ExecuteServerCommand()
}
