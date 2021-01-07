//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/audit/writercloser.mock.go -package=audit_mock

package audit_test

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit"
)

type WriterCloserSyncer interface {
	io.WriteCloser
	Sync() error
}

func Test_eventWriter_Write(t *testing.T) {
	type fields struct {
		order binary.ByteOrder
	}
	type args struct {
		evs []*audit.Event
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Write a single event - little endian",
			fields: fields{
				order: binary.LittleEndian,
			},
			args: args{
				evs: testEvents[:1],
			},
			wantErr: false,
		},
		{
			name: "Write a single event - big endian",
			fields: fields{
				order: binary.BigEndian,
			},
			args: args{
				evs: testEvents[:1],
			},
			wantErr: false,
		},
		{
			name: "Write multiple events - little endian",
			fields: fields{
				order: binary.LittleEndian,
			},
			args: args{
				evs: testEvents,
			},
			wantErr: false,
		},
		{
			name: "Write multiple events - big endian",
			fields: fields{
				order: binary.BigEndian,
			},
			args: args{
				evs: testEvents,
			},
			wantErr: false,
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			writerMock := audit_mock.NewMockWriterCloserSyncer(ctrl)
			calls := make([]*gomock.Call, 0)
			for i := 0; i < len(tt.args.evs); i++ {
				calls = append(calls,
					writerMock.
						EXPECT().
						Write(gomock.Any()).
						Do(func(data []byte) {
							t.Logf("got payload = %s", hex.EncodeToString(data))
							t.Logf("got length %d", tt.fields.order.Uint32(data))
						}),
					writerMock.
						EXPECT().
						Write(gomock.Any()).
						Do(func(data []byte) {
							t.Logf("got payload = %s", hex.EncodeToString(data))
						}),
					writerMock.
						EXPECT().
						Sync(),
				)
			}
			gomock.InOrder(calls...)

			e := audit.NewEventWriter(writerMock, audit.WithWriterByteOrder(tt.fields.order))

			for _, ev := range tt.args.evs {
				if err := e.Write(ev); (err != nil) != tt.wantErr {
					t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
