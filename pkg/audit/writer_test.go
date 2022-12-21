//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/audit/writercloser.mock.go -package=audit_mock

package audit_test

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"testing"

	"github.com/golang/mock/gomock"

	audit_mock "inetmock.icb4dc0.de/inetmock/internal/mock/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
)

type WriterCloserSyncer interface {
	io.WriteCloser
	Sync() error
}

func Test_eventWriter_Write(t *testing.T) {
	t.Parallel()
	type args struct {
		evs []*audit.Event
	}
	type testCase struct {
		name    string
		args    args
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Write a single event",
			args: args{
				evs: testEvents()[:1],
			},
			wantErr: false,
		},
		{
			name: "Write multiple events",
			args: args{
				evs: testEvents(),
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			writerMock := audit_mock.NewMockWriterCloserSyncer(ctrl)
			calls := make([]*gomock.Call, 0)
			for i := 0; i < len(tt.args.evs); i++ {
				calls = append(calls,
					writerMock.
						EXPECT().
						Write(gomock.Any()).
						Do(func(data []byte) {
							t.Logf("got payload = %s", hex.EncodeToString(data))
							t.Logf("got length %d", binary.BigEndian.Uint32(data))
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

			e := audit.NewEventWriter(writerMock)

			for _, ev := range tt.args.evs {
				if err := e.Write(ev); (err != nil) != tt.wantErr {
					t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
