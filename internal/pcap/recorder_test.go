//go:build linux && sudo
// +build linux,sudo

package pcap_test

import (
	"context"
	"net"
	"sort"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/pcap"
	"gitlab.com/inetmock/inetmock/internal/pcap/consumers"
)

type fakeWriterCloser struct {
	closeHandle func() error
}

func (fakeWriterCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (f fakeWriterCloser) Close() error {
	if f.closeHandle != nil {
		return f.closeHandle()
	}
	return nil
}

func Test_recorder_Subscriptions(t *testing.T) {
	t.Parallel()
	type subscriptionRequest struct {
		Name   string
		Device string
	}
	type testCase struct {
		name              string
		requests          []subscriptionRequest
		wantResult        interface{}
		wantSubscriptions interface{}
	}
	tests := []testCase{
		{
			name:              "Empty",
			wantResult:        td.NotNil(),
			wantSubscriptions: td.Empty(),
		},
		{
			name: "Subscription to loopback",
			requests: []subscriptionRequest{
				{
					Name:   "test",
					Device: "lo",
				},
			},
			wantResult: &pcap.StartRecordingResult{
				ConsumerKey: "lo:test",
			},
			wantSubscriptions: td.Set(pcap.Subscription{
				ConsumerKey:  "lo:test",
				ConsumerName: "test",
			}),
		},
		{
			name: "Multiple subscriptions to loopback",
			requests: []subscriptionRequest{
				{
					Name:   "test",
					Device: "lo",
				},
				{
					Name:   "test2",
					Device: "lo",
				},
			},
			wantResult: td.Any(
				&pcap.StartRecordingResult{
					ConsumerKey: "lo:test",
				},
				&pcap.StartRecordingResult{
					ConsumerKey: "lo:test2",
				},
			),
			wantSubscriptions: td.Set(pcap.Subscription{
				ConsumerKey:  "lo:test",
				ConsumerName: "test",
			}, pcap.Subscription{
				ConsumerKey:  "lo:test2",
				ConsumerName: "test2",
			}),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := pcap.NewRecorder()

			t.Cleanup(func() {
				if err := r.Close(); err != nil {
					t.Errorf("Recorder.Close() error = %v", err)
				}
			})

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			for _, req := range tt.requests {
				if result, err := r.StartRecording(ctx, req.Device, consumers.NewNoOpConsumerWithName(req.Name)); err != nil {
					t.Errorf("StartRecording() error = %v", err)
				} else {
					td.Cmp(t, result, tt.wantResult)
				}
			}

			gotSubscriptions := sortSubscriptions(r.Subscriptions())
			td.Cmp(t, gotSubscriptions, tt.wantSubscriptions)
		})
	}
}

func Test_recorder_StartRecordingWithOptions(t *testing.T) {
	t.Parallel()
	type args struct {
		device   string
		consumer pcap.Consumer
		opts     pcap.RecordingOptions
	}
	type testCase struct {
		name          string
		args          args
		want          interface{}
		wantErr       bool
		recorderSetup func() (recorder pcap.Recorder, err error)
	}
	tests := []testCase{
		{
			name: "Listen to lo",
			recorderSetup: func() (recorder pcap.Recorder, err error) {
				recorder = pcap.NewRecorder()
				return
			},
			args: args{
				device:   "lo",
				consumer: consumers.NewNoOpConsumer(),
				opts: pcap.RecordingOptions{
					Promiscuous: false,
					ReadTimeout: 10 * time.Second,
				},
			},
			want: td.Struct(new(pcap.StartRecordingResult), td.StructFields{
				"ConsumerKey": td.Contains("lo:"),
			}),
			wantErr: false,
		},
		{
			name: "Listen to lo with existing name",
			recorderSetup: func() (recorder pcap.Recorder, err error) {
				recorder = pcap.NewRecorder()
				_, err = recorder.StartRecording(context.Background(), "lo", consumers.NewNoOpConsumerWithName("test"))
				return
			},
			args: args{
				device:   "lo",
				consumer: consumers.NewNoOpConsumerWithName("test"),
				opts: pcap.RecordingOptions{
					Promiscuous: false,
					ReadTimeout: 10 * time.Second,
				},
			},
			want:    td.Nil(),
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var err error
			var recorder pcap.Recorder

			if recorder, err = tt.recorderSetup(); err != nil {
				t.Fatalf("recorderSetup() error = %v", err)
			}

			t.Cleanup(func() {
				if err = recorder.Close(); err != nil {
					t.Errorf("Recorder.Close() error = %v", err)
				}
			})

			var result *pcap.StartRecordingResult
			result, err = recorder.StartRecordingWithOptions(context.Background(), tt.args.device, tt.args.consumer, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartRecordingWithOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
			td.Cmp(t, result, tt.want)
		})
	}
}

func Test_recorder_AvailableDevices(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name    string
		want    interface{}
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Expect lo device",
			want: td.Contains(td.Struct(pcap.Device{}, td.StructFields{
				"Name":        "lo",
				"IPAddresses": td.Contains(net.IPv4(127, 0, 0, 1)),
			})),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			re := pcap.NewRecorder()
			t.Cleanup(func() {
				if err := re.Close(); err != nil {
					t.Errorf("Recorder.Close() error = %v", err)
				}
			})
			gotDevices, err := re.AvailableDevices()
			if (err != nil) != tt.wantErr {
				t.Errorf("AvailableDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotDevices, tt.want)
		})
	}
}

func Test_recorder_StopRecording(t *testing.T) {
	t.Parallel()
	type args struct {
		consumerKey string
	}
	type testCase struct {
		name          string
		args          args
		recorderSetup func(t *testing.T) (recorder pcap.Recorder, err error)
		wantErr       bool
	}
	tests := []testCase{
		{
			name: "Stop non existing recording",
			args: args{
				consumerKey: "lo:test.pcap",
			},
			recorderSetup: func(t *testing.T) (recorder pcap.Recorder, err error) {
				t.Helper()
				return pcap.NewRecorder(), nil
			},
			wantErr: true,
		},
		{
			name: "Stop recording lo:test",
			args: args{
				consumerKey: "lo:test",
			},
			recorderSetup: func(t *testing.T) (recorder pcap.Recorder, err error) {
				t.Helper()
				recorder = pcap.NewRecorder()
				_, err = recorder.StartRecording(context.Background(), "lo", consumers.NewNoOpConsumerWithName("test"))

				return
			},
			wantErr: false,
		},
		{
			name: "Stop recording - ensure closing of writers",
			args: args{
				consumerKey: "lo:test",
			},
			recorderSetup: func(t *testing.T) (recorder pcap.Recorder, err error) {
				t.Helper()
				recorder = pcap.NewRecorder()
				var writerConsumer pcap.Consumer
				gotClosed := false
				writerConsumer, err = consumers.NewWriterConsumer("test", &fakeWriterCloser{
					func() error {
						gotClosed = true
						return nil
					},
				})

				if err != nil {
					return
				}

				t.Cleanup(func() {
					if !gotClosed {
						t.Errorf("writer was not closed")
					}
				})

				_, err = recorder.StartRecording(context.Background(), "lo", writerConsumer)

				return
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var err error
			var r pcap.Recorder
			if r, err = tt.recorderSetup(t); err != nil {
				t.Fatalf("recorderSetup() error = %v", err)
			}

			t.Cleanup(func() {
				if err := r.Close(); err != nil {
					t.Errorf("Recorder.Close() error = %v", err)
				}
			})

			if err := r.StopRecording(tt.args.consumerKey); (err != nil) != tt.wantErr {
				t.Errorf("StopRecording() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_recorder_StopRecordingFromContext(t *testing.T) {
	t.Parallel()
	var err error
	var writerConsumer pcap.Consumer
	gotClosed := false
	writerConsumer, err = consumers.NewWriterConsumer("test", &fakeWriterCloser{
		func() error {
			gotClosed = true
			return nil
		},
	})

	if err != nil {
		return
	}

	t.Cleanup(func() {
		if !gotClosed {
			t.Errorf("writer was not closed")
		}
	})

	recorder := pcap.NewRecorder()

	t.Cleanup(func() {
		if err = recorder.Close(); err != nil {
			t.Errorf("Recorder.Close() error = %v", err)
		}
	})

	recordingCtx, recordingCancel := context.WithCancel(context.Background())
	defer recordingCancel()
	if _, err = recorder.StartRecording(recordingCtx, "lo", writerConsumer); err != nil {
		t.Errorf("StartRecording() error = %v", err)
		return
	}

	if len(recorder.Subscriptions()) < 1 {
		t.Fatal("No subscription even if there is one expected")
	}

	recordingCancel()

	// give the cleanup a bit of time
	time.Sleep(10 * time.Millisecond)

	if len(recorder.Subscriptions()) > 0 {
		t.Errorf("Subscriptions present but none expected: %v", recorder.Subscriptions())
	}
}

func sortSubscriptions(subs []pcap.Subscription) []pcap.Subscription {
	sort.Slice(subs, func(i, j int) bool {
		return subs[i].ConsumerName < subs[j].ConsumerName
	})
	return subs
}
