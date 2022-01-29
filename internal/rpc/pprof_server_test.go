package rpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"gitlab.com/inetmock/inetmock/internal/rpc"
	"gitlab.com/inetmock/inetmock/internal/rpc/test"
	tst "gitlab.com/inetmock/inetmock/internal/test"
	rpcv1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

func Test_profilingServer_ProfileDump(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		req     *rpcv1.ProfileDumpRequest
		want    interface{}
		wantErr bool
	}{
		{
			name: "Error - non-existing profile",
			req: &rpcv1.ProfileDumpRequest{
				ProfileName: "asdf",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Dump heap profile, without GC",
			req: &rpcv1.ProfileDumpRequest{
				ProfileName: "heap",
			},
			want: td.Struct(new(rpcv1.ProfileDumpResponse), td.StructFields{
				"ProfileData": td.NotEmpty(),
			}),
		},
		{
			name: "Dump heap profile, with GC",
			req: &rpcv1.ProfileDumpRequest{
				ProfileName:  "heap",
				GcBeforeDump: true,
			},
			want: td.Struct(new(rpcv1.ProfileDumpResponse), td.StructFields{
				"ProfileData": td.NotEmpty(),
			}),
		},
		{
			name: "Dump heap profile, without GC, in legacy format",
			req: &rpcv1.ProfileDumpRequest{
				ProfileName: "heap",
				Debug:       1,
			},
			want: td.Struct(new(rpcv1.ProfileDumpResponse), td.StructFields{
				"ProfileData": td.NotEmpty(),
			}),
		},
		{
			name: "Dump allocs profile",
			req: &rpcv1.ProfileDumpRequest{
				ProfileName: "allocs",
			},
			want: td.Struct(new(rpcv1.ProfileDumpResponse), td.StructFields{
				"ProfileData": td.NotEmpty(),
			}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := setupTestProfilingServer(t)
			ctx, cancel := context.WithTimeout(tst.Context(t), 500*time.Millisecond)
			t.Cleanup(cancel)
			got, err := client.ProfileDump(ctx, tt.req)
			if err != nil {
				if !tt.wantErr {
					td.CmpNoError(t, err)
				}
				return
			}

			td.Cmp(t, got, tt.want)
		})
	}
}

func Test_profilingServer_CPUProfile(t *testing.T) {
	t.Parallel()
	type args struct {
		ctxTimeout time.Duration
		req        *rpcv1.CPUProfileRequest
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Error - duration exceeds timeout",
			args: args{
				ctxTimeout: 50 * time.Millisecond,
				req: &rpcv1.CPUProfileRequest{
					ProfileDuration: durationpb.New(100 * time.Millisecond),
				},
			},
			wantErr: true,
		},
		{
			name: "Collect profile of 1s",
			args: args{
				ctxTimeout: 5 * time.Second,
				req: &rpcv1.CPUProfileRequest{
					ProfileDuration: durationpb.New(1 * time.Second),
				},
			},
			want: td.Struct(new(rpcv1.CPUProfileResponse), td.StructFields{
				"ProfileData": td.NotEmpty(),
			}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := setupTestProfilingServer(t)
			ctx, cancel := context.WithTimeout(tst.Context(t), tt.args.ctxTimeout)
			t.Cleanup(cancel)
			got, err := client.CPUProfile(ctx, tt.args.req)
			if err != nil {
				if !tt.wantErr {
					td.CmpNoError(t, err)
				}
				return
			}

			td.Cmp(t, got, tt.want)
		})
	}
}

func Test_profilingServer_Trace(t *testing.T) {
	t.Parallel()
	type args struct {
		ctxTimeout time.Duration
		req        *rpcv1.TraceRequest
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Error - duration exceeds timeout",
			args: args{
				ctxTimeout: 50 * time.Millisecond,
				req: &rpcv1.TraceRequest{
					TraceDuration: durationpb.New(100 * time.Millisecond),
				},
			},
			wantErr: true,
		},
		{
			name: "Collect profile of 1s",
			args: args{
				ctxTimeout: 5 * time.Second,
				req: &rpcv1.TraceRequest{
					TraceDuration: durationpb.New(1 * time.Second),
				},
			},
			want: td.Struct(new(rpcv1.TraceResponse), td.StructFields{
				"ProfileData": td.NotEmpty(),
			}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := setupTestProfilingServer(t)
			ctx, cancel := context.WithTimeout(tst.Context(t), tt.args.ctxTimeout)
			t.Cleanup(cancel)
			got, err := client.Trace(ctx, tt.args.req)
			if err != nil {
				if !tt.wantErr {
					td.CmpNoError(t, err)
				}
				return
			}

			td.Cmp(t, got, tt.want)
		})
	}
}

func setupTestProfilingServer(t *testing.T) rpcv1.ProfilingServiceClient {
	t.Helper()
	p := rpc.NewProfilingServer()
	srv := test.NewTestGRPCServer(t, func(registrar grpc.ServiceRegistrar) {
		rpcv1.RegisterProfilingServiceServer(registrar, p)
	})

	ctx, cancel := context.WithTimeout(tst.Context(t), 100*time.Millisecond)
	defer cancel()
	conn := srv.Dial(ctx, t)

	return rpcv1.NewProfilingServiceClient(conn)
}
