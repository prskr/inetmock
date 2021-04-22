package rpc_test

import (
	"context"
	"errors"
	"net"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"gitlab.com/inetmock/inetmock/internal/pcap"
	"gitlab.com/inetmock/inetmock/internal/pcap/consumers"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"gitlab.com/inetmock/inetmock/internal/rpc/test"
	tst "gitlab.com/inetmock/inetmock/internal/test"
	rpcV1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

func Test_pcapServer_ListActiveRecordings(t *testing.T) {
	t.Parallel()
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
				t.Helper()
				recorder = pcap.NewRecorder()
				return
			},
			wantSubscriptions: nil,
			wantErr:           false,
		},
		{
			name: "Listening to lo interface",
			recorderSetup: func(t *testing.T) (recorder pcap.Recorder, err error) {
				t.Helper()
				recorder = pcap.NewRecorder()
				err = recorder.StartRecording(context.Background(), "lo", consumers.NewNoOpConsumerWithName("test"))
				return
			},
			wantSubscriptions: []string{"lo:test"},
			wantErr:           false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var err error
			var recorder pcap.Recorder
			if recorder, err = tt.recorderSetup(t); err != nil {
				t.Errorf("recorderSetup() error = %v", err)
				return
			}

			t.Cleanup(func() {
				if err = recorder.Close(); err != nil {
					t.Errorf("recorder.Close() error = %v", err)
				}
			})

			pcapClient := setupTestPCAPServer(t, recorder)

			gotResp, err := pcapClient.ListActiveRecordings(context.Background(), new(rpcV1.ListActiveRecordingsRequest))
			if (err != nil) != tt.wantErr {
				t.Errorf("ListActiveRecordings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResp == nil {
				if !tt.wantErr {
					t.Errorf("response was nil")
				}
				return
			}

			sort.Strings(gotResp.Subscriptions)
			sort.Strings(tt.wantSubscriptions)
			if !reflect.DeepEqual(gotResp.Subscriptions, tt.wantSubscriptions) {
				t.Errorf("ListActiveRecordings() gotResp = %v, want %v", gotResp.Subscriptions, tt.wantSubscriptions)
			}
		})
	}
}

func Test_pcapServer_ListAvailableDevices(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name    string
		matcher func(devs []*rpcV1.ListAvailableDevicesResponse_PCAPDevice) error
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Ensure that any device was found",
			matcher: func(devs []*rpcV1.ListAvailableDevicesResponse_PCAPDevice) error {
				if len(devs) > 0 {
					return nil
				}
				return errors.New("didn't find loopback device")
			},
			wantErr: false,
		},
		{
			name: "Ensure that any device with an assigned IP was found",
			matcher: func(devs []*rpcV1.ListAvailableDevicesResponse_PCAPDevice) error {
				for _, dev := range devs {
					for _, addr := range dev.Addresses {
						if addr == nil {
							continue
						}
						ip := net.IP(addr)
						if !ip.IsUnspecified() {
							return nil
						}
					}
				}
				return errors.New("didn't find loopback device")
			},
			wantErr: false,
		},
		{
			name: "Ensure that loopback device was found",
			matcher: func(devs []*rpcV1.ListAvailableDevicesResponse_PCAPDevice) error {
				foundLoopback := false
				for _, dev := range devs {
					foundLoopback = foundLoopback || dev.Name == "lo"
					if foundLoopback {
						return nil
					}
				}
				return errors.New("didn't find loopback device")
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var err error
			var recorder = pcap.NewRecorder()

			t.Cleanup(func() {
				if err = recorder.Close(); err != nil {
					t.Errorf("recorder.Close() error = %v", err)
				}
			})

			pcapClient := setupTestPCAPServer(t, recorder)
			got, err := pcapClient.ListAvailableDevices(context.Background(), new(rpcV1.ListAvailableDevicesRequest))
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAvailableDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantErr {
				if err := tt.matcher(got.AvailableDevices); err != nil {
					t.Errorf("ListAvailableDevices() matcher error = %v", err)
				}
			}
		})
	}
}

func Test_pcapServer_StartPCAPFileRecording(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name              string
		req               *rpcV1.StartPCAPFileRecordingRequest
		wantSubscriptions []pcap.Subscription
		wantErr           bool
	}
	tests := []testCase{
		{
			name: "Start a recording on lo interface",
			req: &rpcV1.StartPCAPFileRecordingRequest{
				Device:     "lo",
				TargetPath: "test.pcap",
			},
			wantSubscriptions: []pcap.Subscription{
				{
					ConsumerKey:  "lo:test.pcap",
					ConsumerName: "test.pcap",
				},
			},
			wantErr: false,
		},
		{
			name: "Start a recording on a non-existing interface",
			req: &rpcV1.StartPCAPFileRecordingRequest{
				Device:     uuid.NewString(),
				TargetPath: "test.pcap",
			},
			wantSubscriptions: nil,
			wantErr:           true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var err error
			var recorder = pcap.NewRecorder()

			t.Cleanup(func() {
				if err = recorder.Close(); err != nil {
					t.Errorf("recorder.Close() error = %v", err)
				}
			})

			pcapClient := setupTestPCAPServer(t, recorder)

			if _, err := pcapClient.StartPCAPFileRecording(context.Background(), tt.req); (err != nil) != tt.wantErr {
				t.Errorf("StartPCAPFileRecording() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if currentSubs := recorder.Subscriptions(); !reflect.DeepEqual(currentSubs, tt.wantSubscriptions) {
				t.Errorf("StartPCAPFileRecording() got = %v, want %v", currentSubs, tt.wantSubscriptions)
			}
		})
	}
}

func Test_pcapServer_StopPCAPFileRecord(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name             string
		keyToRemove      string
		recorderSetup    func(t *testing.T) (recorder pcap.Recorder, err error)
		removedRecording bool
		wantErr          bool
	}
	tests := []testCase{
		{
			name:        "Remove non existing recording",
			keyToRemove: "lo:asdf.pcap",
			recorderSetup: func(t *testing.T) (recorder pcap.Recorder, err error) {
				t.Helper()
				recorder = pcap.NewRecorder()
				return
			},
			removedRecording: false,
			wantErr:          false,
		},
		{
			name:        "Remove an existing recording",
			keyToRemove: "lo:test.pcap",
			recorderSetup: func(t *testing.T) (recorder pcap.Recorder, err error) {
				t.Helper()
				recorder = pcap.NewRecorder()
				err = recorder.StartRecording(context.Background(), "lo", consumers.NewNoOpConsumerWithName("test.pcap"))
				return
			},
			removedRecording: true,
			wantErr:          false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var err error
			var recorder pcap.Recorder
			if recorder, err = tt.recorderSetup(t); err != nil {
				t.Errorf("recorderSetup() error = %v", err)
				return
			}

			t.Cleanup(func() {
				if err = recorder.Close(); err != nil {
					t.Errorf("recorder.Close() error = %v", err)
				}
			})

			pcapClient := setupTestPCAPServer(t, recorder)
			gotResp, err := pcapClient.StopPCAPFileRecording(context.Background(), &rpcV1.StopPCAPFileRecordingRequest{
				ConsumerKey: tt.keyToRemove,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("StopPCAPFileRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResp.Removed != tt.removedRecording {
				t.Errorf("StopPCAPFileRecord() removed = %v, want %v", gotResp.Removed, tt.removedRecording)
			}
		})
	}
}

func setupTestPCAPServer(t *testing.T, recorder pcap.Recorder) rpcV1.PCAPServiceClient {
	t.Helper()
	p := rpc.NewPCAPServer(t.TempDir(), recorder)

	srv := test.NewTestGRPCServer(t, func(registrar grpc.ServiceRegistrar) {
		rpcV1.RegisterPCAPServiceServer(registrar, p)
	})

	ctx, cancel := context.WithTimeout(tst.Context(t), 100*time.Millisecond)
	var conn = srv.Dial(ctx, t)
	cancel()

	return rpcV1.NewPCAPServiceClient(conn)
}
