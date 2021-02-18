package audit

import (
	"net/http"
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"gitlab.com/inetmock/inetmock/pkg/audit/details"
)

func Test_guessDetailsFromApp(t *testing.T) {

	mustAny := func(msg proto.Message) *anypb.Any {
		a, err := anypb.New(msg)
		if err != nil {
			panic(a)
		}
		return a
	}

	type args struct {
		any *anypb.Any
	}
	type testCase struct {
		name string
		args args
		want Details
	}
	tests := []testCase{
		{
			name: "HTTP etails",
			args: args{
				any: mustAny(&details.HTTPDetailsEntity{
					Method: details.HTTPMethod_GET,
					Host:   "localhost",
					Uri:    "http://localhost/asdf",
					Proto:  "HTTP",
				}),
			},
			want: details.HTTP{
				Method:  "GET",
				Host:    "localhost",
				URI:     "http://localhost/asdf",
				Proto:   "HTTP",
				Headers: http.Header{},
			},
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			if got := guessDetailsFromApp(tt.args.any); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guessDetailsFromApp() = %v, want %v", got, tt.want)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
