package rpc

import (
	"context"
	"io"
	"os"

	"gitlab.com/inetmock/inetmock/internal/app"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	"go.uber.org/zap"
)

type auditServer struct {
	UnimplementedAuditServer
	app app.App
}

func (a *auditServer) WatchEvents(req *WatchEventsRequest, srv Audit_WatchEventsServer) (err error) {
	a.app.Logger().Info("watcher attached", zap.String("name", req.WatcherName))
	err = a.app.EventStream().RegisterSink(sink.NewGRPCSink(srv.Context(), req.WatcherName, func(ev audit.Event) {
		if err = srv.Send(ev.ProtoMessage()); err != nil {
			return
		}
	}))

	if err != nil {
		return
	}

	<-srv.Context().Done()
	a.app.Logger().Info("Watcher detached", zap.String("name", req.WatcherName))
	return
}

func (a *auditServer) RegisterFileSink(_ context.Context, req *RegisterFileSinkRequest) (resp *RegisterFileSinkResponse, err error) {
	var writer io.WriteCloser
	if writer, err = os.OpenFile(req.TargetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		return
	}
	if err = a.app.EventStream().RegisterSink(sink.NewWriterSink(req.TargetPath, audit.NewEventWriter(writer))); err != nil {
		return
	}
	resp = &RegisterFileSinkResponse{}
	return
}

func (a *auditServer) RemoveFileSink(_ context.Context, req *RemoveFileSinkRequest) (*RemoveFileSinkResponse, error) {
	a.app.EventStream().RemoveSink(req.TargetPath)
	return &RemoveFileSinkResponse{}, nil
}
