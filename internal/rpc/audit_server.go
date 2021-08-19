package rpc

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	v1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

var _ v1.AuditServiceServer = (*auditServer)(nil)

type auditServer struct {
	v1.UnimplementedAuditServiceServer
	logger           logging.Logger
	eventStream      audit.EventStream
	auditDataDirPath string
}

func NewAuditServiceServer(logger logging.Logger, eventStream audit.EventStream, auditDataDirPath string) v1.AuditServiceServer {
	return &auditServer{
		logger:           logger,
		eventStream:      eventStream,
		auditDataDirPath: auditDataDirPath,
	}
}

func (a *auditServer) ListSinks(context.Context, *v1.ListSinksRequest) (*v1.ListSinksResponse, error) {
	return &v1.ListSinksResponse{
		Sinks: a.eventStream.Sinks(),
	}, nil
}

func (a *auditServer) WatchEvents(req *v1.WatchEventsRequest, srv v1.AuditService_WatchEventsServer) (err error) {
	logger := a.logger
	logger.Info("watcher attached", zap.String("name", req.WatcherName))
	err = a.eventStream.RegisterSink(srv.Context(), sink.NewGenericSink(req.WatcherName, func(ev audit.Event) {
		if err = srv.Send(&v1.WatchEventsResponse{Entity: ev.ProtoMessage()}); err != nil {
			return
		}
	}))

	if err != nil {
		return
	}

	<-srv.Context().Done()
	logger.Info("Watcher detached", zap.String("name", req.WatcherName))
	return
}

func (a *auditServer) RegisterFileSink(_ context.Context, req *v1.RegisterFileSinkRequest) (*v1.RegisterFileSinkResponse, error) {
	targetPath := req.TargetPath
	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(a.auditDataDirPath, req.TargetPath)
	}

	var writer io.WriteCloser
	var err error
	if writer, err = os.Create(targetPath); err != nil {
		return nil, PathToGRPCError(err)
	}
	if err = a.eventStream.RegisterSink(context.Background(), sink.NewWriterSink(req.TargetPath, audit.NewEventWriter(writer))); err != nil {
		if errors.Is(err, audit.ErrSinkAlreadyRegistered) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &v1.RegisterFileSinkResponse{
		ResolvedPath: targetPath,
	}, nil
}

func (a *auditServer) RemoveFileSink(_ context.Context, req *v1.RemoveFileSinkRequest) (*v1.RemoveFileSinkResponse, error) {
	if gotRemoved := a.eventStream.RemoveSink(req.TargetPath); gotRemoved {
		return &v1.RemoveFileSinkResponse{
			SinkGotRemoved: gotRemoved,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "file sink with given target path not found")
}
