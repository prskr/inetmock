package mock_test

import (
	"testing"

	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
)

func TestConditionHandler_Matches(t *testing.T) {
	t.Parallel()
	type fields struct {
		Filters []mock.RequestFilter
	}
	type args struct {
		req *mdns.Question
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Empty filters - match",
			fields: fields{},
			args: args{
				req: new(mdns.Question),
			},
			want: true,
		},
		{
			name: "Single filter - match",
			fields: fields{
				Filters: []mock.RequestFilter{
					mock.RequestFilterFunc(func(*mdns.Question) bool {
						return true
					}),
				},
			},
			args: args{
				req: new(mdns.Question),
			},
			want: true,
		},
		{
			name: "Single filter - no match",
			fields: fields{
				Filters: []mock.RequestFilter{
					mock.RequestFilterFunc(func(*mdns.Question) bool {
						return false
					}),
				},
			},
			args: args{
				req: new(mdns.Question),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := mock.ConditionHandler{
				Filters: tt.fields.Filters,
			}
			if got := h.Matches(tt.args.req); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
