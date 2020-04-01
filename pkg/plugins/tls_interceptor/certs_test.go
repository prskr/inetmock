package main

import (
	"crypto/x509"
	"testing"
	"time"
)

type testTimeSource struct {
	nowValue time.Time
}

func (t testTimeSource) UTCNow() time.Time {
	return t.nowValue
}

func Test_certShouldBeRenewed(t *testing.T) {
	type args struct {
		timeSource timeSource
		cert       *x509.Certificate
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Detect cert is expired",
			want: true,
			args: args{
				cert: &x509.Certificate{
					NotAfter:  time.Now().UTC().Add(1 * time.Hour),
					NotBefore: time.Now().UTC().Add(-1 * time.Hour),
				},
				timeSource: testTimeSource{
					nowValue: time.Now().UTC().Add(2 * time.Hour),
				},
			},
		},
		{
			name: "Detect cert should be renewed",
			want: true,
			args: args{
				cert: &x509.Certificate{
					NotAfter:  time.Now().UTC().Add(1 * time.Hour),
					NotBefore: time.Now().UTC().Add(-1 * time.Hour),
				},
				timeSource: testTimeSource{
					nowValue: time.Now().UTC().Add(45 * time.Minute),
				},
			},
		},
		{
			name: "Detect cert shouldn't be renewed",
			want: false,
			args: args{
				cert: &x509.Certificate{
					NotAfter:  time.Now().UTC().Add(1 * time.Hour),
					NotBefore: time.Now().UTC().Add(-1 * time.Hour),
				},
				timeSource: testTimeSource{
					nowValue: time.Now().UTC().Add(25 * time.Minute),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := certShouldBeRenewed(tt.args.timeSource, tt.args.cert); got != tt.want {
				t.Errorf("certShouldBeRenewed() = %v, want %v", got, tt.want)
			}
		})
	}
}
