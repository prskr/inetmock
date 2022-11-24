package rpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"

	"google.golang.org/grpc"

	auditm "inetmock.icb4dc0.de/inetmock/internal/mock/audit"
	"inetmock.icb4dc0.de/inetmock/internal/rpc"
	"inetmock.icb4dc0.de/inetmock/internal/rpc/test"
	tst "inetmock.icb4dc0.de/inetmock/internal/test"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	rpcv1 "inetmock.icb4dc0.de/inetmock/pkg/rpc/v1"
)

const (
	grpcTimeout = 500 * time.Millisecond
)

func Test_auditServer_RemoveFileSink(t *testing.T) {
	t.Parallel()
	type fields struct {
		eventStreamSetup func(t *testing.T) audit.EventStream
	}
	tests := []struct {
		name    string
		req     *rpcv1.RemoveFileSinkRequest
		fields  fields
		want    td.StructFields
		wantErr bool
	}{
		{
			name: "Remove existing file sink - success",
			req: &rpcv1.RemoveFileSinkRequest{
				TargetPath: "test.pcap",
			},
			fields: fields{
				eventStreamSetup: func(t *testing.T) audit.EventStream {
					t.Helper()
					ctrl := gomock.NewController(t)

					es := auditm.NewMockEventStream(ctrl)
					es.
						EXPECT().
						RemoveSink("test.pcap").
						Return(true)

					return es
				},
			},
			want: td.StructFields{
				"SinkGotRemoved": true,
			},
			wantErr: false,
		},
		{
			name: "Remove non-existing file sink - success",
			req: &rpcv1.RemoveFileSinkRequest{
				TargetPath: "test.pcap",
			},
			fields: fields{
				eventStreamSetup: func(t *testing.T) audit.EventStream {
					t.Helper()
					ctrl := gomock.NewController(t)

					es := auditm.NewMockEventStream(ctrl)
					es.
						EXPECT().
						RemoveSink("test.pcap").
						Return(false)

					return es
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testCtx := tst.Context(t)
			logger := logging.CreateTestLogger(t)

			srv := test.NewTestGRPCServer(t, func(registrar grpc.ServiceRegistrar) {
				rpcv1.RegisterAuditServiceServer(registrar, rpc.NewAuditServiceServer(logger, tt.fields.eventStreamSetup(t), t.TempDir()))
			})

			ctx, cancel := context.WithTimeout(testCtx, grpcTimeout)
			conn := srv.Dial(ctx, t)
			cancel()

			client := rpcv1.NewAuditServiceClient(conn)

			ctx, cancel = context.WithTimeout(testCtx, grpcTimeout)
			t.Cleanup(cancel)
			got, err := client.RemoveFileSink(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveFileSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				td.CmpStruct(t, got, new(rpcv1.RemoveFileSinkResponse), tt.want)
			}
		})
	}
}
