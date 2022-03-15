package middleware

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ContextErrorConverter grpc.UnaryServerInterceptor = func(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	resp, err = handler(ctx, req)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		err = status.Error(codes.DeadlineExceeded, err.Error())
	case errors.Is(err, context.Canceled):
		err = status.Errorf(codes.Canceled, err.Error())
	default:
	}

	return
}
