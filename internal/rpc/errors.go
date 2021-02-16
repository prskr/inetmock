package rpc

import (
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PathToGRPCError(err error) error {
	if os.IsPermission(err) {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	if os.IsNotExist(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	if os.IsTimeout(err) {
		return status.Error(codes.DeadlineExceeded, err.Error())
	}

	return status.Error(codes.Unknown, err.Error())
}
