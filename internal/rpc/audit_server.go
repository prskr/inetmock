package rpc

import (
	"context"
	"io"
	"os"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

type auditServer struct {
	UnimplementedAuditServer
	logger      logging.Logger
	eventStream audit.EventStream
}

func (a *auditServer) ListSinks(context.Context, *ListSinksRequest) (*ListSinksResponse, error) {
	return &ListSinksResponse{
		Sinks: a.eventStream.Sinks(),
	}, nil
}

func (a *auditServer) WatchEvents(req *WatchEventsRequest, srv Audit_WatchEventsServer) (err error) {
	a.logger.Info("watcher attached", zap.String("name", req.WatcherName))
	err = a.eventStream.RegisterSink(srv.Context(), sink.NewGRPCSink(req.WatcherName, func(ev audit.Event) {
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

func (a *auditServer) RegisterFileSink(_ context.Context, req *RegisterFileSinkRequest) (resp *RegisterFileSinkResponse, err error) {
	var writer io.WriteCloser
	var flags int

	switch req.OpenMode {
	case FileOpenMode_APPEND:
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	default:
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	var permissions = os.FileMode(req.Permissions)
	if permissions == 0 {
		permissions = 644
	}

	if writer, err = os.OpenFile(req.TargetPath, flags, permissions); err != nil {
		return
	}
	if err = a.eventStream.RegisterSink(context.Background(), sink.NewWriterSink(req.TargetPath, audit.NewEventWriter(writer))); err != nil {
		return
	}
	resp = &RegisterFileSinkResponse{}
	return
}

func (a *auditServer) RemoveFileSink(_ context.Context, req *RemoveFileSinkRequest) (*RemoveFileSinkResponse, error) {
	gotRemoved := a.eventStream.RemoveSink(req.TargetPath)
	return &RemoveFileSinkResponse{
		SinkGotRemoved: gotRemoved,
	}, nil
}
