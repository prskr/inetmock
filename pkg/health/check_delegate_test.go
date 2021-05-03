package health_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/health"
)

func Test_checkDelegate_Name(t *testing.T) {
	t.Parallel()
	type fields struct {
		name           string
		statusDelegate func(ctx context.Context) error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Empty name",
			want: "",
		},
		{
			name: "Any name",
			fields: fields{
				name: "My fancy check",
			},
			want: "My fancy check",
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := health.NewCheckFunc(tt.fields.name, tt.fields.statusDelegate)
			if got := c.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkDelegate_Status(t *testing.T) {
	t.Parallel()
	type fields struct {
		name           string
		statusDelegate func(ctx context.Context) error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "No delegate",
			fields: fields{
				name: "SampleDelegate",
			},
			wantErr: false,
		},
		{
			name: "No error from delegate",
			fields: fields{
				name: "SampleDelegate",
				statusDelegate: func(context.Context) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Error from delegate",
			fields: fields{
				name: "SampleDelegate",
				statusDelegate: func(context.Context) error {
					return errors.New("any kind of error")
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(test.Context(t), 50*time.Millisecond)
			t.Cleanup(cancel)
			c := health.NewCheckFunc(tt.fields.name, tt.fields.statusDelegate)
			if err := c.Status(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Status() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
