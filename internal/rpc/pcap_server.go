package rpc

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"inetmock.icb4dc0.de/inetmock/internal/netutils"
	"inetmock.icb4dc0.de/inetmock/internal/pcap"
	"inetmock.icb4dc0.de/inetmock/internal/pcap/consumers"
	rpcv1 "inetmock.icb4dc0.de/inetmock/pkg/rpc/v1"
)

var _ rpcv1.PCAPServiceServer = (*pcapServer)(nil)

type pcapServer struct {
	rpcv1.UnimplementedPCAPServiceServer
	recorder    pcap.Recorder
	pcapDataDir string
}

func NewPCAPServer(pcapDataDir string, recorder pcap.Recorder) rpcv1.PCAPServiceServer {
	return &pcapServer{
		pcapDataDir: pcapDataDir,
		recorder:    recorder,
	}
}

func (p *pcapServer) ListActiveRecordings(
	context.Context,
	*rpcv1.ListActiveRecordingsRequest,
) (resp *rpcv1.ListActiveRecordingsResponse, _ error) {
	resp = new(rpcv1.ListActiveRecordingsResponse)
	subs := p.recorder.Subscriptions()
	for i := range subs {
		resp.Subscriptions = append(resp.Subscriptions, subs[i].ConsumerKey)
	}

	return
}

func (p *pcapServer) ListAvailableDevices(
	context.Context,
	*rpcv1.ListAvailableDevicesRequest,
) (*rpcv1.ListAvailableDevicesResponse, error) {
	var (
		devs []pcap.Device
		err  error
	)
	if devs, err = p.recorder.AvailableDevices(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := new(rpcv1.ListAvailableDevicesResponse)
	for i := range devs {
		resp.AvailableDevices = append(resp.AvailableDevices, &rpcv1.ListAvailableDevicesResponse_PCAPDevice{
			Name:      devs[i].Name,
			Addresses: netutils.IPAddressesToBytes(devs[i].IPAddresses),
		})
	}

	return resp, nil
}

func (p *pcapServer) StartPCAPFileRecording(
	_ context.Context,
	req *rpcv1.StartPCAPFileRecordingRequest,
) (*rpcv1.StartPCAPFileRecordingResponse, error) {
	targetPath := req.TargetPath
	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(p.pcapDataDir, req.TargetPath)
	}

	var writer io.WriteCloser
	var err error
	if writer, err = os.Create(targetPath); err != nil {
		return nil, PathToGRPCError(err)
	}

	var consumer pcap.Consumer
	if consumer, err = consumers.NewWriterConsumer(req.TargetPath, writer); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	readTimeout := req.ReadTimeout.AsDuration()
	if readTimeout == 0 {
		readTimeout = pcap.DefaultReadTimeout
	}

	opts := pcap.RecordingOptions{
		Promiscuous: req.Promiscuous,
		ReadTimeout: readTimeout,
	}

	var result *pcap.StartRecordingResult
	//nolint:contextcheck // is running independent of request context
	if result, err = p.recorder.StartRecordingWithOptions(context.Background(), req.Device, consumer, opts); err != nil {
		if errors.Is(err, pcap.ErrConsumerAlreadyRegistered) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &rpcv1.StartPCAPFileRecordingResponse{
		ResolvedPath: targetPath,
		ConsumerKey:  result.ConsumerKey,
	}, nil
}

func (p *pcapServer) StopPCAPFileRecording(
	_ context.Context,
	request *rpcv1.StopPCAPFileRecordingRequest,
) (resp *rpcv1.StopPCAPFileRecordingResponse, _ error) {
	resp = new(rpcv1.StopPCAPFileRecordingResponse)
	resp.Removed = p.recorder.StopRecording(request.ConsumerKey) == nil
	return
}
