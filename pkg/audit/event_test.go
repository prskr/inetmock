//nolint:testpackage // testing internals here - needs to be private
package audit

import (
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

func Test_guessDetailsFromApp(t *testing.T) {
	t.Parallel()
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
				any: mustAny(&auditv1.HTTPDetailsEntity{
					Method: auditv1.HTTPMethod_HTTP_METHOD_GET,
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
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := guessDetailsFromApp(tt.args.any)
			td.Cmp(t, got, tt.want)
		})
	}
}
