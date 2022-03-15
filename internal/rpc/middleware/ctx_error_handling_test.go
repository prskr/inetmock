package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/inetmock/inetmock/internal/rpc/middleware"
	"gitlab.com/inetmock/inetmock/internal/test"
)

func TestContextErrorConverter(t *testing.T) {
	t.Parallel()
	type args struct {
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
		want      error
	}{
		{
			name: "Not a context error",
			args: args{
				handler: func(context.Context, any) (any, error) {
					return nil, errors.New("there's something strange in the neighborhood")
				},
			},
			wantError: true,
			want:      errors.New("there's something strange in the neighborhood"),
		},
		{
			name: "Not an error at all",
			args: args{
				handler: func(context.Context, any) (any, error) {
					return nil, nil
				},
			},
			wantError: false,
		},
		{
			name: "Context canceled error",
			args: args{
				handler: func(context.Context, any) (any, error) {
					return nil, context.Canceled
				},
			},
			wantError: true,
			want:      status.Error(codes.Canceled, context.Canceled.Error()),
		},
		{
			name: "Context DeadlineExceeded error",
			args: args{
				handler: func(context.Context, any) (any, error) {
					return nil, context.DeadlineExceeded
				},
			},
			wantError: true,
			want:      status.Error(codes.Canceled, context.DeadlineExceeded.Error()),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(test.Context(t), 50*time.Millisecond)
			t.Cleanup(cancel)
			if _, err := middleware.ContextErrorConverter(ctx, nil, new(grpc.UnaryServerInfo), tt.args.handler); err != nil {
				if !tt.wantError && !errors.Is(err, tt.want) {
					t.Errorf("Got error but type was not expected - error = %v", err)
				}
				return
			}
			if tt.wantError {
				t.Error("Expected error but didn't get one")
			}
		})
	}
}
