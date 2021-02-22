package audit_test

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
	"testing"

	"gitlab.com/inetmock/inetmock/pkg/audit"
)

var (
	//nolint:lll
	httpPayloadBytesBigEndian = `000000e2120b088092b8c398feffffff01180120022a047f00000132047f00000138d8fc0140504a3308041224544c535f45434448455f45434453415f574954485f4145535f3235365f4342435f5348411a096c6f63616c686f73745287010a37747970652e676f6f676c65617069732e636f6d2f696e65746d6f636b2e61756469742e76312e4854545044657461696c73456e74697479124c080112096c6f63616c686f73741a15687474703a2f2f6c6f63616c686f73742f6173646622084854545020312e312a1c0a0641636365707412120a106170706c69636174696f6e2f6a736f6e`
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
			wantEv:  testEvents[0],
			wantErr: false,
		},
		{
			name: "Read DNS payload",
			fields: fields{
				source: mustDecodeHex(dnsPayloadBytesBigEndian),
			},
			wantEv:  testEvents[1],
			wantErr: false,
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			e := audit.NewEventReader(tt.fields.source)
			gotEv, err := e.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(gotEv, *tt.wantEv) {
				t.Errorf("Read() gotEv = %v, want %v", gotEv, tt.wantEv)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
