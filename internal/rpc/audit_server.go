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
	"gitlab.com/inetmock/inetmock/pkg/rpc"
)

type auditServer struct {
	rpc.UnimplementedAuditServer
	logger           logging.Logger
	eventStream      audit.EventStream
	auditDataDirPath string
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
	var targetPath = req.TargetPath
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

	return &rpc.RegisterFileSinkResponse{
		ResolvedPath: targetPath,
	}, nil
}

func (a *auditServer) RemoveFileSink(_ context.Context, req *rpc.RemoveFileSinkRequest) (*rpc.RemoveFileSinkResponse, error) {
	if gotRemoved := a.eventStream.RemoveSink(req.TargetPath); gotRemoved {
		return &rpc.RemoveFileSinkResponse{}, nil
	}
	return nil, status.Error(codes.NotFound, "file sink with given target path not found")
}
