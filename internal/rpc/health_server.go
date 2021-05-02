package rpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

var (
	_ v1.HealthServer = (*healthServer)(nil)
)

type healthServer struct {
	v1.UnimplementedHealthServer
	checker health.Checker
}

func NewHealthServer(checker health.Checker) v1.HealthServer {
	return &healthServer{
		checker: checker,
	}
}

func (h healthServer) Check(ctx context.Context, request *v1.HealthCheckRequest) (resp *v1.HealthCheckResponse, err error) {
	var result health.Result
	if result, err = h.checker.Status(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Aborted, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if request.Service != "" {
		known, result := result.CheckResult(request.Service)
		if !known {
			return nil, status.Error(codes.NotFound, request.Service)
		}

		if result == nil {
			return &v1.HealthCheckResponse{
				Status: v1.HealthCheckResponse_SERVING,
			}, nil
		} else {
			return &v1.HealthCheckResponse{
				Status: v1.HealthCheckResponse_NOT_SERVING,
			}, nil
		}
	}

	if result.IsHealthy() {
		return &v1.HealthCheckResponse{
			Status: v1.HealthCheckResponse_SERVING,
		}, nil
	} else {
		return &v1.HealthCheckResponse{
			Status: v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
}

func (h healthServer) Watch(request *v1.HealthCheckRequest, server v1.Health_WatchServer) error {
	var latestStatus v1.HealthCheckResponse_ServingStatus
	if resp, err := h.Check(server.Context(), request); err != nil {
		return err
	} else {
		latestStatus = resp.Status
		if err = server.Send(resp); err != nil {
			return err
		}
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	select {
	case <-server.Context().Done():
		if errors.Is(server.Context().Err(), context.Canceled) {
			return status.Error(codes.Canceled, server.Context().Err().Error())
		}
		if errors.Is(server.Context().Err(), context.DeadlineExceeded) {
			return status.Error(codes.DeadlineExceeded, server.Context().Err().Error())
		}
	case <-ticker.C:
		if resp, err := h.Check(server.Context(), request); err != nil {
			return err
		} else if resp.Status != latestStatus {
			latestStatus = resp.Status
			if err = server.Send(resp); err != nil {
				return err
			}
		}
	}

	return nil
}
