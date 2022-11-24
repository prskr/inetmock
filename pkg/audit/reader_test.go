package audit_test

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/pkg/audit"
)

var (
	//nolint:lll
	httpPayloadBytesBigEndian = `000000a7120b088092b8c398feffffff01180120022a047f00000132047f00000138d8fc0140504a3308041224544c535f45434448455f45434453415f574954485f4145535f3235365f4342435f5348411a096c6f63616c686f7374a2014c080112096c6f63616c686f73741a15687474703a2f2f6c6f63616c686f73742f6173646622084854545020312e312a1c0a0641636365707412120a106170706c69636174696f6e2f6a736f6e`
	//nolint:lll
	dnsPayloadBytesBigEndian = `0000003b120b088092b8c398feffffff01180220012a100000000000000000000000000000000132100000000000000000000000000000000138d8fc014050`
)

func mustDecodeHex(hexBytes string) io.Reader {
	b, err := hex.DecodeString(hexBytes)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func Test_eventReader_Read(t *testing.T) {
	t.Parallel()
	type fields struct {
		source io.Reader
	}
	type testCase struct {
		name    string
		fields  fields
		wantEv  *audit.Event
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Read HTTP payload",
			fields: fields{
				source: mustDecodeHex(httpPayloadBytesBigEndian),
			},
			wantEv:  testEvents()[0],
			wantErr: false,
		},
		{
			name: "Read DNS payload",
			fields: fields{
				source: mustDecodeHex(dnsPayloadBytesBigEndian),
			},
			wantEv:  testEvents()[1],
			wantErr: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := audit.NewEventReader(tt.fields.source)
			gotEv, err := e.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			td.Cmp(t, gotEv, tt.wantEv)
		})
	}
}
