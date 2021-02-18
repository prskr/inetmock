// +build linux

package pcap_test

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"
	"time"

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
	type subscriptionRequest struct {
		Name   string
		Device string
	}
	type testCase struct {
		name              string
		requests          []subscriptionRequest
		wantSubscriptions []pcap.Subscription
	}
	tests := []testCase{
		{
			name: "Emtpy",
		},
		{
			name: "Subscription to loopback",
			requests: []subscriptionRequest{
				{
					Name:   "test",
					Device: "lo",
				},
			},
			wantSubscriptions: []pcap.Subscription{
				{
					ConsumerKey:  "lo:test",
					ConsumerName: "test",
				},
			},
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
			wantSubscriptions: []pcap.Subscription{
				{
					ConsumerKey:  "lo:test",
					ConsumerName: "test",
				},
				{
					ConsumerKey:  "lo:test2",
					ConsumerName: "test2",
				},
			},
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			r := pcap.NewRecorder()

			t.Cleanup(func() {
				if err := r.Close(); err != nil {
					t.Errorf("Recorder.Close() error = %v", err)
				}
			})

			for _, req := range tt.requests {
				if err := r.StartRecording(context.Background(), req.Device, consumers.NewNoOpConsumerWithName(req.Name)); err != nil {
					t.Errorf("StartRecording() error = %v", err)
				}
			}

			if gotSubscriptions := sortSubscriptions(r.Subscriptions()); !reflect.DeepEqual(gotSubscriptions, tt.wantSubscriptions) {
				t.Errorf("Subscriptions() = %v, want %v", gotSubscriptions, tt.wantSubscriptions)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}

func Test_recorder_StartRecordingWithOptions(t *testing.T) {
	type args struct {
		device   string
		consumer pcap.Consumer
		opts     pcap.RecordingOptions
	}
	type testCase struct {
		name          string
		args          args
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
			wantErr: false,
		},
		{
			name: "Listen to lo with existing name",
			recorderSetup: func() (recorder pcap.Recorder, err error) {
				recorder = pcap.NewRecorder()
				err = recorder.StartRecording(context.Background(), "lo", consumers.NewNoOpConsumerWithName("test"))
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
			wantErr: true,
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			var err error
			var recorder pcap.Recorder

			if recorder, err = tt.recorderSetup(); err != nil {
				t.Fatalf("recorderSetup() error = %v", err)
			}

			t.Cleanup(func() {
				if err := recorder.Close(); err != nil {
					t.Errorf("Recorder.Close() error = %v", err)
				}
			})

			if err = recorder.StartRecordingWithOptions(context.Background(), tt.args.device, tt.args.consumer, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("StartRecordingWithOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}

func Test_recorder_AvailableDevices(t *testing.T) {
	type testCase struct {
		name    string
		mtacher func(got []pcap.Device) error
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Expect lo device",
			mtacher: func(got []pcap.Device) error {
				if len(got) < 1 {
					return errors.New("expected at least one interface")
				}

				foundLoopbackDevice := false
				for _, d := range got {
					foundLoopbackDevice = foundLoopbackDevice || d.Name == "lo"
				}

				if !foundLoopbackDevice {
					return errors.New("didn't find loopback device")
				}

				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if err := tt.mtacher(gotDevices); err != nil {
				t.Errorf("AvailableDevices() matcher error = %v", err)
			}
		})
	}
}

func sortSubscriptions(subs []pcap.Subscription) []pcap.Subscription {
	sort.Slice(subs, func(i, j int) bool {
		return subs[i].ConsumerName < subs[j].ConsumerName
	})
	return subs
}

func Test_recorder_StopRecording(t *testing.T) {
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
				recorder = pcap.NewRecorder()
				err = recorder.StartRecording(context.Background(), "lo", consumers.NewNoOpConsumerWithName("test"))

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
				recorder = pcap.NewRecorder()
				var writerConsumer pcap.Consumer
				gotClosed := false
				writerConsumer, err = consumers.NewWriterConsumer("test", &fakeWriterCloser{
					func() error {
						gotClosed = true
						return nil
					},
				})

				t.Cleanup(func() {
					if !gotClosed {
						t.Errorf("writer was not closed")
					}
				})

				if err != nil {
					return
				}
				err = recorder.StartRecording(context.Background(), "lo", writerConsumer)

				return
			},
			wantErr: false,
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
