package mock_test

import (
	"net"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	"gitlab.com/inetmock/inetmock/internal/rules"
)

var (
	sampleCIDR = net.IPNet{
		IP:   net.IPv4(192, 168, 0, 0),
		Mask: net.CIDRMask(24, 32),
	}
)

func TestIPHandlerForArgs(t *testing.T) {
	t.Parallel()
	type args struct {
		args []rules.Param
	}
	tests := []struct {
		name    string
		args    args
		want    net.IP
		wantErr bool
	}{
		{
			name: "Handler to return 1.1.1.1",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(1, 1, 1, 1),
					},
				},
			},
			want:    net.IPv4(1, 1, 1, 1),
			wantErr: false,
		},
		{
			name: "Handler to return 192.168.0.1",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 1),
					},
				},
			},
			want:    net.IPv4(192, 168, 0, 1),
			wantErr: false,
		},
		{
			name: "Missing param",
			args: args{
				args: make([]rules.Param, 0),
			},
			wantErr: true,
		},
		{
			name: "Empty param",
			args: args{
				args: []rules.Param{
					{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := mock.IPHandlerForArgs(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("IPHandlerForArgs() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !got.Lookup("").Equal(tt.want) {
				t.Errorf("IPResolver.Lookup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncrementalHandlerForArgs(t *testing.T) {
	t.Parallel()
	type args struct {
		args []rules.Param
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []net.IP
	}{
		{
			name: "Single incremental IP",
			args: args{
				args: []rules.Param{
					{
						CIDR: &rules.CIDR{
							IPNet: mustParseCIDR("192.168.0.0/24"),
						},
					},
				},
			},
			want: []net.IP{
				net.IPv4(192, 168, 0, 1),
			},
		},
		{
			name: "Multiple incremental IPs",
			args: args{
				args: []rules.Param{
					{
						CIDR: &rules.CIDR{
							IPNet: mustParseCIDR("192.168.0.0/24"),
						},
					},
				},
			},
			want: []net.IP{
				net.IPv4(192, 168, 0, 1),
				net.IPv4(192, 168, 0, 2),
				net.IPv4(192, 168, 0, 3),
			},
		},
		{
			name: "Missing param",
			args: args{
				args: make([]rules.Param, 0),
			},
			wantErr: true,
		},
		{
			name: "Empty param",
			args: args{
				args: []rules.Param{
					{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := mock.IncrementalHandlerForArgs(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("IncrementalHandlerForArgs() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			for _, want := range tt.want {
				if ip := got.Lookup(""); ip == nil || !ip.Equal(want) {
					t.Errorf("Lookup() got = %v want = %v", ip, want)
				}
			}
		})
	}
}

func TestRandomHandlerForArgs(t *testing.T) {
	t.Parallel()
	type args struct {
		args []rules.Param
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Simple random IP",
			args: args{
				args: []rules.Param{
					{
						CIDR: &rules.CIDR{
							IPNet: mustParseCIDR("192.168.0.0/24"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing param",
			args: args{
				args: make([]rules.Param, 0),
			},
			wantErr: true,
		},
		{
			name: "Empty param",
			args: args{
				args: []rules.Param{
					{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := mock.RandomHandlerForArgs(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RandomHandlerForArgs() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if ip := got.Lookup(""); !sampleCIDR.Contains(ip) {
				t.Errorf("Expected %v to be in CIDR %s", ip, sampleCIDR.Network())
			}
		})
	}
}
