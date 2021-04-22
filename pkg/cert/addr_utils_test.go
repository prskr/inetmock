package cert

import "testing"

func Test_extractIPFromAddress(t *testing.T) {
	t.Parallel()
	type args struct {
		addr string
	}
	type testCase struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "Get address for IPv4 address",
			want:    "127.0.0.1",
			wantErr: false,
			args: args{
				addr: "127.0.0.1:23492",
			},
		},
		{
			name:    "Get address for IPv6 address",
			want:    "::1",
			wantErr: false,
			args: args{
				addr: "[::1]:23492",
			},
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := extractIPFromAddress(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractIPFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractIPFromAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
