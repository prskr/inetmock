//go:build sudo

package netflow_test

import (
	"testing"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

//nolint:paralleltest // would conflict with the same interface
func TestFirewall_AttachToInterface(t *testing.T) {
	RemoveMemlock(t)

	type args struct {
		interfaceName string
		fwCfg         netflow.FirewallInterfaceConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty rule set",
			args: args{
				interfaceName: "lo",
				fwCfg: netflow.FirewallInterfaceConfig{
					RemoveMemLock: true,
					DefaultPolicy: netflow.PacketPolicyPass,
					Monitor:       false,
				},
			},
			wantErr: false,
		},
		{
			name: "Non-existing interface",
			args: args{
				interfaceName: "lo1234",
				fwCfg: netflow.FirewallInterfaceConfig{
					RemoveMemLock: true,
					DefaultPolicy: netflow.PacketPolicyPass,
					Monitor:       false,
				},
			},
			wantErr: true,
		},
		{
			name: "Single rule",
			args: args{
				interfaceName: "lo",
				fwCfg: netflow.FirewallInterfaceConfig{
					RemoveMemLock: true,
					DefaultPolicy: netflow.PacketPolicyPass,
					Monitor:       false,
					Rules: []netflow.RuleEntry{
						{
							Policy:      netflow.PacketPolicyDrop,
							Destination: "8080/tcp",
							Monitor:     netflow.BoolP(true),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple rules",
			args: args{
				interfaceName: "lo",
				fwCfg: netflow.FirewallInterfaceConfig{
					RemoveMemLock: true,
					DefaultPolicy: netflow.PacketPolicyPass,
					Monitor:       false,
					Rules: []netflow.RuleEntry{
						{
							Policy:      netflow.PacketPolicyDrop,
							Destination: "80/tcp",
							Monitor:     netflow.BoolP(true),
						},
						{
							Policy:      netflow.PacketPolicyDrop,
							Destination: "8080/tcp",
							Monitor:     netflow.BoolP(true),
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
			f := netflow.NewFirewall(new(packetSinkRecorder), netflow.WithMockingEnabled(true), ErrorSink(t))

			t.Cleanup(func() {
				if err := f.Close(); err != nil {
					t.Errorf("f.Close() err = %v", err)
				}
			})

			if err := f.AttachToInterface(tt.args.interfaceName, tt.args.fwCfg); (err != nil) != tt.wantErr {
				t.Errorf("AttachToInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
