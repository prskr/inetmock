//go:build sudo

package netflow_test

import (
	"testing"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

//nolint:paralleltest // would conflict on same interface
func TestNAT_AttachToInterface(t *testing.T) {
	RemoveMemlock(t)

	type args struct {
		interfaceName string
		spec          netflow.NATTableSpec
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty NAT table",
			args: args{
				interfaceName: "lo",
			},
			wantErr: false,
		},
		{
			name: "Single rule NAT table",
			args: args{
				interfaceName: "lo",
				spec: netflow.NATTableSpec{
					Translations: []netflow.NATTargetSpec{
						{
							Destination: "0.0.0.0:80/tcp",
							RedirectTo:  netflow.NATTargetInterface,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			n := netflow.NewNAT(netflow.WithMockingEnabled(true), ErrorSink(t))

			t.Cleanup(func() {
				if err := n.Close(); err != nil {
					t.Errorf("Failed to close NAT: %v", err)
				}
			})

			if err := n.AttachToInterface(tt.args.interfaceName, tt.args.spec); (err != nil) != tt.wantErr {
				t.Errorf("AttachToInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
