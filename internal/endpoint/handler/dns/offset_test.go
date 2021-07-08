package dns_test

import (
	"math"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

func TestOffset_Next(t *testing.T) {
	t.Parallel()
	type fields struct {
		CurrentOffset int
	}
	tests := []struct {
		name         string
		fields       fields
		wantOffset   int
		wantOverflow bool
	}{
		{
			name: "Empty offset",
			fields: fields{
				CurrentOffset: 0,
			},
			wantOffset:   0,
			wantOverflow: false,
		},
		{
			name: "Offset of arbitrary number",
			fields: fields{
				CurrentOffset: 1337,
			},
			wantOffset:   1337,
			wantOverflow: false,
		},
		{
			name: "Offset of MaxInt64",
			fields: fields{
				CurrentOffset: math.MaxInt64,
			},
			wantOffset:   0,
			wantOverflow: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t1 *testing.T) {
			t1.Parallel()
			t := td.NewT(t1)
			o := &dns.Offset{
				CurrentOffset: tt.fields.CurrentOffset,
			}
			gotOffset, gotOverflow := o.Next()
			t.Cmp(gotOffset, tt.wantOffset)
			t.Cmp(gotOverflow, tt.wantOverflow)
			if !gotOverflow {
				gotOffset, _ = o.Next()
				t.Cmp(gotOffset, tt.wantOffset+1)
			}
		})
	}
}

func TestOffset_Inc(t *testing.T) {
	t.Parallel()
	type fields struct {
		CurrentOffset int
	}
	type args struct {
		val int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantNewVal   int
		wantOverflow bool
	}{
		{
			name: "Empty offset",
			fields: fields{
				CurrentOffset: 0,
			},
			args: args{
				val: 1,
			},
			wantNewVal:   1,
			wantOverflow: false,
		},
		{
			name: "Any arbitrary Number",
			fields: fields{
				CurrentOffset: 3,
			},
			args: args{
				val: 1337,
			},
			wantNewVal:   1340,
			wantOverflow: false,
		},
		{
			name: "Overflow",
			fields: fields{
				CurrentOffset: math.MaxInt64 - 2,
			},
			args: args{
				val: 4,
			},
			wantNewVal:   2,
			wantOverflow: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t1 *testing.T) {
			t1.Parallel()
			t := td.NewT(t1)
			o := &dns.Offset{
				CurrentOffset: tt.fields.CurrentOffset,
			}
			gotNewVal, gotOverflow := o.Inc(tt.args.val)
			t.Cmp(gotNewVal, tt.wantNewVal)
			t.Cmp(gotOverflow, tt.wantOverflow)
		})
	}
}
