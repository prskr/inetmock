package queue

import (
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
)

func Test_TTL_UpdateTTL(t *testing.T) {
	t.Parallel()
	type fields struct {
		seed []*Entry
	}
	type args struct {
		idxToUpdate int
		newTTL      time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "Already evicted item",
			fields: fields{
				seed: []*Entry{
					{
						Key:     "gitlab.com.",
						Value:   net.IPv4(1, 2, 3, 4),
						timeout: time.Now().Add(-100 * time.Millisecond),
						index:   evictedEntryIndex,
					},
				},
			},
			args: args{
				idxToUpdate: 0,
				newTTL:      100 * time.Millisecond,
			},
			want: td.Struct(&Entry{index: 0}, td.StructFields{}),
		},
		{
			name: "Evicted item without well-known index",
			fields: fields{
				seed: []*Entry{
					{
						Key:     "gitlab.com.",
						Value:   net.IPv4(1, 2, 3, 4),
						timeout: time.Now().Add(-100 * time.Millisecond),
						index:   -42,
					},
				},
			},
			args: args{
				idxToUpdate: 0,
				newTTL:      100 * time.Millisecond,
			},
			want: td.Struct(&Entry{index: 0}, td.StructFields{}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			queue := NewTTLFromSeed(tt.fields.seed)
			e := queue.Get(tt.args.idxToUpdate)
			queue.UpdateTTL(e, tt.args.newTTL)
			td.Cmp(t, e, tt.want)
		})
	}
}
