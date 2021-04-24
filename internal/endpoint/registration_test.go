package endpoint_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	dnsmock "gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	httpmock "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
)

func Test_handlerRegistry_AvailableHandlers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		handlerRegistry       endpoint.HandlerRegistry
		wantAvailableHandlers interface{}
	}{
		{
			name:                  "Empty registry",
			handlerRegistry:       endpoint.NewHandlerRegistry(),
			wantAvailableHandlers: td.Nil(),
		},
		{
			name: "Single handler registered",
			handlerRegistry: func() endpoint.HandlerRegistry {
				registry := endpoint.NewHandlerRegistry()
				_ = httpmock.AddHTTPMock(registry)
				return registry
			}(),
			wantAvailableHandlers: td.Set(endpoint.HandlerReference("http_mock")),
		},
		{
			name: "Multiple handler registered",
			handlerRegistry: func() endpoint.HandlerRegistry {
				registry := endpoint.NewHandlerRegistry()
				_ = httpmock.AddHTTPMock(registry)
				_ = dnsmock.AddDNSMock(registry)
				return registry
			}(),
			wantAvailableHandlers: td.Set(
				endpoint.HandlerReference("dns_mock"),
				endpoint.HandlerReference("http_mock"),
			),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotAvailableHandlers := tt.handlerRegistry.AvailableHandlers()
			td.Cmp(t, gotAvailableHandlers, tt.wantAvailableHandlers)
		})
	}
}

func Test_handlerRegistry_HandlerForName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		handlerRegistry endpoint.HandlerRegistry
		handlerRef      endpoint.HandlerReference
		wantInstance    interface{}
		wantOk          bool
	}{
		{
			name:            "Empty registry",
			handlerRegistry: endpoint.NewHandlerRegistry(),
			handlerRef:      "http_mock",
			wantInstance:    nil,
			wantOk:          false,
		},
		{
			name: "Registry with HTTP mock registered",
			handlerRegistry: func() endpoint.HandlerRegistry {
				registry := endpoint.NewHandlerRegistry()
				_ = httpmock.AddHTTPMock(registry)
				return registry
			}(),
			handlerRef:   "http_mock",
			wantInstance: td.NotNil(),
			wantOk:       true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotInstance, gotOk := tt.handlerRegistry.HandlerForName(tt.handlerRef)
			td.Cmp(t, gotInstance, tt.wantInstance)
			td.Cmp(t, gotOk, tt.wantOk)
		})
	}
}
