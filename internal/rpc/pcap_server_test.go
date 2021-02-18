package rpc_test

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"google.golang.org/grpc"

	"gitlab.com/inetmock/inetmock/internal/pcap"
	"gitlab.com/inetmock/inetmock/internal/pcap/consumers"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"gitlab.com/inetmock/inetmock/internal/rpc/test"
	rpc2 "gitlab.com/inetmock/inetmock/pkg/rpc"
)

func Test_pcapServer_ListActiveRecordings(t *testing.T) {
	type testCase struct {
		name              string
		recorderSetup     func(t *testing.T) (recorder pcap.Recorder, err error)
		wantSubscriptions []string
		wantErr           bool
	}
	tests := []testCase{
		{
			name: "No subscriptions",
			recorderSetup: func(t *testing.T) (recorder pcap.Recorder, err error) {
				recorder = pcap.NewRecorder()
				return
			},
			wantSubscriptions: nil,
			wantErr:           false,
		},
		{
			name: "Listening to lo interface",
			recorderSetup: func(t *testing.T) (recorder pcap.Recorder, err error) {
				recorder = pcap.NewRecorder()
				err = recorder.StartRecording(context.Background(), "lo", consumers.NewNoOpConsumerWithName("test"))
				return
			},
			wantSubscriptions: []string{"lo:test"},
			wantErr:           false,
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			var err error
			var recorder pcap.Recorder
			if recorder, err = tt.recorderSetup(t); err != nil {
				t.Fatalf("recorderSetup() error = %v", err)
			}

			t.Cleanup(func() {
				if err := recorder.Close(); err != nil {
					t.Errorf("recorder.Close() error = %v", err)
				}
			})

			pcapClient := setupTestPCAPServer(t, recorder)

			gotResp, err := pcapClient.ListActiveRecordings(context.Background(), new(rpc2.ListRecordingsRequest))
			if (err != nil) != tt.wantErr {
				t.Errorf("ListActiveRecordings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotResp == nil) != tt.wantErr {
				t.Errorf("response was nil")
			} else {
				sort.Strings(gotResp.Subscriptions)
				sort.Strings(tt.wantSubscriptions)
				if !reflect.DeepEqual(gotResp.Subscriptions, tt.wantSubscriptions) {
					t.Errorf("ListActiveRecordings() gotResp = %v, want %v", gotResp.Subscriptions, tt.wantSubscriptions)
				}
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}

func setupTestPCAPServer(t *testing.T, recorder pcap.Recorder) rpc2.PCAPClient {
	var err error
	var srv *test.GRPCServer
	p := rpc.NewPCAPServer(t.TempDir(), recorder)
	srv, err = test.NewTestGRPCServer(func(registrar grpc.ServiceRegistrar) {
		rpc2.RegisterPCAPServer(registrar, p)
	})

	if err != nil {
		t.Fatalf("NewTestGRPCServer() error = %v", err)
	}

	if err = srv.StartServer(); err != nil {
		t.Fatalf("StartServer() error = %v", err)
	}

	t.Cleanup(func() {
		srv.StopServer()
	})

	var conn *grpc.ClientConn
	if conn, err = srv.Dial(context.Background(), grpc.WithInsecure()); err != nil {
		t.Fatalf("Dial() error = %v", err)
	}

	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("conn.Close() error = %v", err)
		}
	})

	return rpc2.NewPCAPClient(conn)
}
