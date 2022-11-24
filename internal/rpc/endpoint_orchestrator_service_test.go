package rpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/internal/rpc"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	rpcv1 "inetmock.icb4dc0.de/inetmock/pkg/rpc/v1"
)

func Test_endpointOrchestratorServer_ListAllServingGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		hostSetup func(tb testing.TB) endpoint.Host
		want      any
		wantErr   bool
	}{
		{
			name: "Empty response",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnConfiguredGroups: func() []endpoint.GroupInfo {
						return nil
					},
				}
			},
			want:    td.Struct(new(rpcv1.ListAllServingGroupsResponse), td.StructFields{}),
			wantErr: false,
		},
		{
			name: "Group not serving",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnConfiguredGroups: func() []endpoint.GroupInfo {
						return []endpoint.GroupInfo{
							{
								Name:    "80/tcp",
								Serving: false,
							},
						}
					},
				}
			},
			want:    td.Struct(new(rpcv1.ListAllServingGroupsResponse), td.StructFields{}),
			wantErr: false,
		},
		{
			name: "Single group serving",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnConfiguredGroups: func() []endpoint.GroupInfo {
						return []endpoint.GroupInfo{
							{
								Name:      "80/tcp",
								Serving:   true,
								Endpoints: []string{"plain_http"},
							},
						}
					},
				}
			},
			want: td.Struct(new(rpcv1.ListAllServingGroupsResponse), td.StructFields{
				"Groups": td.Bag(td.Struct(&rpcv1.ListenerGroup{
					Name: "80/tcp",
				}, td.StructFields{
					"Endpoints": td.Bag("plain_http"),
				})),
			}),
			wantErr: false,
		},
		{
			name: "Filter non-serving group",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnConfiguredGroups: func() []endpoint.GroupInfo {
						return []endpoint.GroupInfo{
							{
								Name:      "80/tcp",
								Serving:   true,
								Endpoints: []string{"plain_http"},
							},
							{
								Name:      "443/tcp",
								Serving:   false,
								Endpoints: []string{"https"},
							},
						}
					},
				}
			},
			want: td.Struct(new(rpcv1.ListAllServingGroupsResponse), td.StructFields{
				"Groups": td.Bag(td.Struct(&rpcv1.ListenerGroup{
					Name: "80/tcp",
				}, td.StructFields{
					"Endpoints": td.Bag("plain_http"),
				})),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := rpc.NewEndpointOrchestratorServer(logging.CreateTestLogger(t), tt.hostSetup(t))
			got, err := s.ListAllServingGroups(context.Background(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAllServingGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}

func Test_endpointOrchestratorServer_ListAllConfiguredGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		hostSetup func(tb testing.TB) endpoint.Host
		want      any
		wantErr   bool
	}{
		{
			name: "Empty response",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnConfiguredGroups: func() []endpoint.GroupInfo {
						return nil
					},
				}
			},
			want:    td.Struct(new(rpcv1.ListAllConfiguredGroupsResponse), td.StructFields{}),
			wantErr: false,
		},
		{
			name: "Group not serving",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnConfiguredGroups: func() []endpoint.GroupInfo {
						return []endpoint.GroupInfo{
							{
								Name:      "80/tcp",
								Serving:   false,
								Endpoints: []string{"plain_http"},
							},
						}
					},
				}
			},
			want: td.Struct(new(rpcv1.ListAllConfiguredGroupsResponse), td.StructFields{
				"Groups": td.Bag(td.Struct(&rpcv1.ListenerGroup{
					Name: "80/tcp",
				}, td.StructFields{
					"Endpoints": td.Bag("plain_http"),
				})),
			}),
			wantErr: false,
		},
		{
			name: "Single group serving",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnConfiguredGroups: func() []endpoint.GroupInfo {
						return []endpoint.GroupInfo{
							{
								Name:      "80/tcp",
								Serving:   true,
								Endpoints: []string{"plain_http"},
							},
						}
					},
				}
			},
			want: td.Struct(new(rpcv1.ListAllConfiguredGroupsResponse), td.StructFields{
				"Groups": td.Bag(td.Struct(&rpcv1.ListenerGroup{
					Name: "80/tcp",
				}, td.StructFields{
					"Endpoints": td.Bag("plain_http"),
				})),
			}),
			wantErr: false,
		},
		{
			name: "Filter non-serving group",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnConfiguredGroups: func() []endpoint.GroupInfo {
						return []endpoint.GroupInfo{
							{
								Name:      "80/tcp",
								Serving:   true,
								Endpoints: []string{"plain_http"},
							},
							{
								Name:      "443/tcp",
								Serving:   false,
								Endpoints: []string{"https"},
							},
						}
					},
				}
			},
			want: td.Struct(new(rpcv1.ListAllConfiguredGroupsResponse), td.StructFields{
				"Groups": td.Bag(
					td.Struct(&rpcv1.ListenerGroup{
						Name: "80/tcp",
					}, td.StructFields{
						"Endpoints": td.Bag("plain_http"),
					}),
					td.Struct(&rpcv1.ListenerGroup{
						Name: "443/tcp",
					}, td.StructFields{
						"Endpoints": td.Bag("https"),
					}),
				),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := rpc.NewEndpointOrchestratorServer(logging.CreateTestLogger(t), tt.hostSetup(t))
			got, err := s.ListAllConfiguredGroups(context.Background(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAllConfiguredGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}

func Test_endpointOrchestratorServer_StartListenerGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		hostSetup func(tb testing.TB) endpoint.Host
		req       *rpcv1.StartListenerGroupRequest
		want      any
		wantErr   bool
	}{
		{
			name: "Return err",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()
				return hostMock{
					OnServeGroup: func(ctx context.Context, groupName string) error {
						return errors.New("nope")
					},
				}
			},
			req:     &rpcv1.StartListenerGroupRequest{GroupName: "https"},
			wantErr: true,
		},
		{
			name: "all good",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()
				return hostMock{
					OnServeGroup: func(ctx context.Context, groupName string) error {
						return nil
					},
				}
			},
			req:     &rpcv1.StartListenerGroupRequest{GroupName: "https"},
			wantErr: false,
			want:    td.NotNil(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := rpc.NewEndpointOrchestratorServer(logging.CreateTestLogger(t), tt.hostSetup(t))
			got, err := s.StartListenerGroup(context.Background(), tt.req)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("StartListenerGroup() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}

func Test_endpointOrchestratorServer_StartAllGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		hostSetup func(tb testing.TB) endpoint.Host
		req       *rpcv1.StartListenerGroupRequest
		want      any
		wantErr   bool
	}{
		{
			name: "Return error",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnServeGroups: func(context.Context) error {
						return errors.New("nope")
					},
				}
			},
			req:     nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Return no error",
			hostSetup: func(tb testing.TB) endpoint.Host {
				tb.Helper()

				return hostMock{
					OnServeGroups: func(context.Context) error {
						return nil
					},
				}
			},
			req:     nil,
			want:    td.NotNil(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := rpc.NewEndpointOrchestratorServer(logging.CreateTestLogger(t), tt.hostSetup(t))
			got, err := s.StartAllGroups(context.Background(), nil)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("StartAllGroups() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}

type hostMock struct {
	OnConfiguredGroups func() []endpoint.GroupInfo
	OnServeGroup       func(ctx context.Context, groupName string) error
	OnServeGroups      func(ctx context.Context) error
	OnShutdown         func(ctx context.Context) error
	OnShutdownGroup    func(ctx context.Context, groupName string) error
}

func (m hostMock) ConfiguredGroups() []endpoint.GroupInfo {
	if m.OnConfiguredGroups != nil {
		return m.OnConfiguredGroups()
	}
	return nil
}

func (m hostMock) ServeGroup(ctx context.Context, groupName string) error {
	if m.OnServeGroup != nil {
		return m.OnServeGroup(ctx, groupName)
	}
	return nil
}

func (m hostMock) ServeGroups(ctx context.Context) error {
	if m.OnServeGroups != nil {
		return m.OnServeGroups(ctx)
	}
	return nil
}

func (m hostMock) Shutdown(ctx context.Context) error {
	if m.OnShutdown != nil {
		return m.OnShutdown(ctx)
	}
	return nil
}

func (m hostMock) ShutdownGroup(ctx context.Context, groupName string) error {
	if m.OnShutdownGroup != nil {
		return m.OnShutdownGroup(ctx, groupName)
	}
	return nil
}
