package rpc

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"inetmock.icb4dc0.de/inetmock/pkg/health"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

var _ v1.HealthServer = (*healthServer)(nil)

type healthServer struct {
	v1.UnimplementedHealthServer
	logger           logging.Logger
	checker          health.Checker
	watchCheckPeriod time.Duration
}

func NewHealthServer(checker health.Checker, watchCheckPeriod time.Duration, logger logging.Logger) v1.HealthServer {
	return &healthServer{
		checker:          checker,
		watchCheckPeriod: watchCheckPeriod,
		logger:           logger,
	}
}

func (h *healthServer) Check(ctx context.Context, request *v1.HealthCheckRequest) (resp *v1.HealthCheckResponse, err error) {
	result := h.checker.Status(ctx)

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
		var fields []zap.Field
		for k, e := range result {
			if e != nil {
				fields = append(fields, zap.NamedError(k, e))
			}
		}
		h.logger.Warn("Health check failed", fields...)
		return &v1.HealthCheckResponse{
			Status: v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
}

func (h *healthServer) Watch(request *v1.HealthCheckRequest, server v1.Health_WatchServer) error {
	var latestStatus v1.HealthCheckResponse_ServingStatus
	if resp, err := h.Check(server.Context(), request); err != nil {
		return err
	} else {
		latestStatus = resp.Status
		if err := server.Send(resp); err != nil {
			return err
		}
	}

	ticker := time.NewTicker(h.watchCheckPeriod)
	defer ticker.Stop()

	for {
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
				if err := server.Send(resp); err != nil {
					return err
				}
			}
		}
	}
}
