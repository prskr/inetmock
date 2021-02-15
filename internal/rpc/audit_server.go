package rpc

import (
	"context"
	"errors"
	"io"
	"os"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type auditServer struct {
	rpc.UnimplementedAuditServer
	logger      logging.Logger
	eventStream audit.EventStream
}

func (a *auditServer) ListSinks(context.Context, *rpc.ListSinksRequest) (*rpc.ListSinksResponse, error) {
	return &rpc.ListSinksResponse{
		Sinks: a.eventStream.Sinks(),
	}, nil
}

func (a *auditServer) WatchEvents(req *rpc.WatchEventsRequest, srv rpc.Audit_WatchEventsServer) (err error) {
	a.logger.Info("watcher attached", zap.String("name", req.WatcherName))
	err = a.eventStream.RegisterSink(srv.Context(), sink.NewGenericSink(req.WatcherName, func(ev audit.Event) {
		if err = srv.Send(ev.ProtoMessage()); err != nil {
			return
		}
	}))

	if err != nil {
		return
	}

	<-srv.Context().Done()
	a.logger.Info("Watcher detached", zap.String("name", req.WatcherName))
	return
}

func (a *auditServer) RegisterFileSink(_ context.Context, req *rpc.RegisterFileSinkRequest) (*rpc.RegisterFileSinkResponse, error) {
	var writer io.WriteCloser
	var flags int

	switch req.OpenMode {
	case rpc.FileOpenMode_APPEND:
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	default:
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	var permissions = os.FileMode(req.Permissions)
	if permissions == 0 {
		permissions = 644
	}

	var err error
	if writer, err = os.OpenFile(req.TargetPath, flags, permissions); err != nil {
		if os.IsPermission(err) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		if os.IsNotExist(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if os.IsTimeout(err) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}
	if err = a.eventStream.RegisterSink(context.Background(), sink.NewWriterSink(req.TargetPath, audit.NewEventWriter(writer))); err != nil {
		if errors.Is(err, audit.ErrSinkAlreadyRegistered) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &rpc.RegisterFileSinkResponse{}, nil
}

func (a *auditServer) RemoveFileSink(_ context.Context, req *rpc.RemoveFileSinkRequest) (*rpc.RemoveFileSinkResponse, error) {
	if gotRemoved := a.eventStream.RemoveSink(req.TargetPath); gotRemoved {
		return &rpc.RemoveFileSinkResponse{}, nil
	}
	return nil, status.Error(codes.NotFound, "file sink with given target path not found")
}
