package endpoint_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

var defaultListenerSpec = endpoint.ListenerSpec{Protocol: "tcp", Port: 1234, Unmanaged: false}

func TestListenerGroup_ConfigureEndpoint(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
		le   *endpoint.ListenerEndpoint
	}
	tests := []struct {
		name string
		spec endpoint.ListenerSpec
		args args
		want interface{}
	}{
		{
			name: "nil value - nothing added",
			spec: defaultListenerSpec,
			args: args{
				name: "plain_http",
				le:   nil,
			},
			want: td.Empty(),
		},
		{
			name: "Add listener endpoint",
			spec: defaultListenerSpec,
			args: args{
				name: "plain_http",
				le:   &endpoint.ListenerEndpoint{},
			},
			want: td.Len(1),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lg, err := endpoint.NewListenerGroup(tt.spec)
			if err != nil {
				t.Errorf("endpoint.NewListenerGroup() error = %v", err)
				return
			}
			lg.ConfigureEndpoint(tt.args.name, tt.args.le)
			td.Cmp(t, lg.ConfiguredEndpoints(), tt.want)
		})
	}
}

func TestListenerGroup_GroupByTLS(t *testing.T) {
	t.Parallel()
	type registration struct {
		name string
		le   *endpoint.ListenerEndpoint
	}

	tests := []struct {
		name          string
		spec          endpoint.ListenerSpec
		registrations []registration
		wantPlainGrp  interface{}
		wantTLSGrp    interface{}
		wantErr       bool
	}{
		{
			name: "Empty ListenerGroup",
			spec: defaultListenerSpec,
			wantPlainGrp: td.Struct(new(endpoint.Group), td.StructFields{
				"Handlers": td.Empty(),
			}),
			wantTLSGrp: td.Struct(new(endpoint.Group), td.StructFields{
				"Handlers": td.Empty(),
			}),
			wantErr: false,
		},
		{
			name: "Single plain group",
			spec: defaultListenerSpec,
			registrations: []registration{
				{
					name: "plain_http",
					le: &endpoint.ListenerEndpoint{
						TLS:     false,
						Handler: MultiplexHandlerMock{},
					},
				},
			},
			wantPlainGrp: td.Struct(new(endpoint.Group), td.StructFields{
				"Handlers": td.Len(1),
				"Names":    []string{"plain_http"},
			}),
			wantTLSGrp: td.Struct(new(endpoint.Group), td.StructFields{
				"Handlers": td.Empty(),
			}),
			wantErr: false,
		},
		{
			name: "Mixed plain and TLS group",
			spec: defaultListenerSpec,
			registrations: []registration{
				{
					name: "plain_http",
					le: &endpoint.ListenerEndpoint{
						TLS:     false,
						Handler: MultiplexHandlerMock{},
					},
				},
				{
					name: "https",
					le: &endpoint.ListenerEndpoint{
						TLS:     true,
						Handler: MultiplexHandlerMock{},
					},
				},
			},
			wantPlainGrp: td.Struct(new(endpoint.Group), td.StructFields{
				"Handlers": td.Len(1),
				"Names":    []string{"plain_http"},
			}),
			wantTLSGrp: td.Struct(new(endpoint.Group), td.StructFields{
				"Handlers": td.Len(1),
				"Names":    []string{"https"},
			}),
			wantErr: false,
		},
		{
			name: "Mixed plain and multiple TLS group",
			spec: defaultListenerSpec,
			registrations: []registration{
				{
					name: "plain_http",
					le: &endpoint.ListenerEndpoint{
						TLS:     false,
						Handler: MultiplexHandlerMock{},
					},
				},
				{
					name: "https",
					le: &endpoint.ListenerEndpoint{
						TLS:     true,
						Handler: MultiplexHandlerMock{},
					},
				},
				{
					name: "doh",
					le: &endpoint.ListenerEndpoint{
						TLS:     true,
						Handler: MultiplexHandlerMock{},
					},
				},
			},
			wantPlainGrp: td.Struct(new(endpoint.Group), td.StructFields{
				"Handlers": td.Len(1),
				"Names":    []string{"plain_http"},
			}),
			wantTLSGrp: td.Struct(new(endpoint.Group), td.StructFields{
				"Handlers": td.Len(2),
				"Names":    td.Bag("https", "doh"),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lg, err := endpoint.NewListenerGroup(tt.spec)
			if err != nil {
				t.Errorf("endpoint.NewListenerGroup() error = %v", err)
				return
			}

			for idx := range tt.registrations {
				r := tt.registrations[idx]
				lg.ConfigureEndpoint(r.name, r.le)
			}

			gotPlainGrp, gotTLSGrp, err := lg.GroupByTLS()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupByTLS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotPlainGrp, tt.wantPlainGrp)
			td.Cmp(t, gotTLSGrp, tt.wantTLSGrp)
		})
	}
}

func TestNewListenerGroup(t *testing.T) {
	t.Parallel()
	type args struct {
		spec endpoint.ListenerSpec
	}
	tests := []struct {
		name    string
		args    args
		wantGrp interface{}
		wantErr bool
	}{
		{
			name: "TCP group - empty name",
			args: args{
				spec: endpoint.ListenerSpec{
					Protocol:  "tcp",
					Port:      80,
					Unmanaged: false,
				},
			},
			wantGrp: td.Struct(&endpoint.ListenerGroup{
				Name: "80/tcp",
			}, td.StructFields{
				"Addr": td.Struct(new(net.TCPAddr), td.StructFields{}),
			}),
			wantErr: false,
		},
		{
			name: "TCP group - keep given name",
			args: args{
				spec: endpoint.ListenerSpec{
					Name:      "asdf",
					Protocol:  "tcp",
					Port:      80,
					Unmanaged: false,
				},
			},
			wantGrp: td.Struct(&endpoint.ListenerGroup{
				Name: "asdf",
			}, td.StructFields{
				"Addr": td.Struct(new(net.TCPAddr), td.StructFields{}),
			}),
			wantErr: false,
		},
		{
			name: "UDP group - empty name",
			args: args{
				spec: endpoint.ListenerSpec{
					Protocol:  "udp",
					Port:      53,
					Unmanaged: false,
				},
			},
			wantGrp: td.Struct(&endpoint.ListenerGroup{
				Name: "53/udp",
			}, td.StructFields{
				"Addr": td.Struct(new(net.UDPAddr), td.StructFields{}),
			}),
			wantErr: false,
		},
		{
			name: "Unmanaged TCP group - empty name",
			args: args{
				spec: endpoint.ListenerSpec{
					Protocol:  "udp",
					Port:      53,
					Unmanaged: true,
				},
			},
			wantGrp: td.Struct(&endpoint.ListenerGroup{
				Name:      "53/udp",
				Unmanaged: true,
			}, td.StructFields{
				"Addr": td.Struct(new(net.UDPAddr), td.StructFields{}),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotGrp, err := endpoint.NewListenerGroup(tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewListenerGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotGrp, tt.wantGrp)
		})
	}
}

func TestListenerGroup_Shutdown(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		spec    endpoint.ListenerSpec
		setup   func(tb testing.TB, lg *endpoint.ListenerGroup)
		wantErr bool
	}{
		{
			name:    "Empty group",
			spec:    defaultListenerSpec,
			wantErr: false,
		},
		{
			name: "Stoppable protocol handler - should be stopped on shutdown",
			spec: defaultListenerSpec,
			setup: func(tb testing.TB, lg *endpoint.ListenerGroup) {
				tb.Helper()
				var calledStop bool
				lg.ConfigureEndpoint("asdf", endpoint.NewListenerEndpoint(endpoint.Spec{}, StoppableProtocolHandlerMock{
					OnStop: func(_ context.Context) error {
						calledStop = true
						return nil
					},
				}))
				tb.Cleanup(func() {
					if !calledStop {
						tb.Error("Expected to handler to be stopped")
					}
				})
			},
			wantErr: false,
		},
		{
			name: "Stoppable protocol handler - should be stopped on shutdown",
			spec: defaultListenerSpec,
			setup: func(tb testing.TB, lg *endpoint.ListenerGroup) {
				tb.Helper()
				lg.ConfigureEndpoint("asdf", endpoint.NewListenerEndpoint(endpoint.Spec{}, StoppableProtocolHandlerMock{
					OnStop: func(_ context.Context) error {
						return errors.New("not that bad :-)")
					},
				}))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lg, err := endpoint.NewListenerGroup(tt.spec)
			if err != nil {
				t.Errorf("endpoint.NewListenerGroup() error = %v", err)
				return
			}

			if tt.setup != nil {
				tt.setup(t, lg)
			}

			if err := lg.Shutdown(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
