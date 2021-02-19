package rpc

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/inetmock/inetmock/internal/pcap"
	"gitlab.com/inetmock/inetmock/internal/pcap/consumers"
	"gitlab.com/inetmock/inetmock/pkg/rpc"
)

type pcapServer struct {
	rpc.UnimplementedPCAPServer
	recorder    pcap.Recorder
	pcapDataDir string
}

func NewPCAPServer(pcapDataDir string, recorder pcap.Recorder) rpc.PCAPServer {
	return &pcapServer{
		pcapDataDir: pcapDataDir,
		recorder:    recorder,
	}
}

func (p *pcapServer) ListActiveRecordings(context.Context, *rpc.ListRecordingsRequest) (resp *rpc.ListRecordingsResponse, _ error) {
	resp = new(rpc.ListRecordingsResponse)
	subs := p.recorder.Subscriptions()
	for i := range subs {
		resp.Subscriptions = append(resp.Subscriptions, subs[i].ConsumerKey)
	}

	return
}

func (p *pcapServer) ListAvailableDevices(context.Context, *rpc.ListAvailableDevicesRequest) (*rpc.ListAvailableDevicesResponse, error) {
	var err error
	var devs []pcap.Device
	if devs, err = p.recorder.AvailableDevices(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resp = new(rpc.ListAvailableDevicesResponse)
	for i := range devs {
		resp.AvailableDevices = append(resp.AvailableDevices, &rpc.ListAvailableDevicesResponse_PCAPDevice{
			Name:      devs[i].Name,
			Addresses: ipAddressesToBytes(devs[i].IPAddresses),
		})
	}

	return resp, nil
}

func (p *pcapServer) StartPCAPFileRecording(
	_ context.Context,
	req *rpc.StartPCAPFileRecordRequest,
) (*rpc.StartPCAPFileRecordResponse, error) {
	var targetPath = req.TargetPath
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

	if err = p.recorder.StartRecordingWithOptions(context.Background(), req.Device, consumer, opts); err != nil {
		if errors.Is(err, pcap.ErrConsumerAlreadyRegistered) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &rpc.StartPCAPFileRecordResponse{
		ResolvedPath: targetPath,
	}, nil
}

func (p *pcapServer) StopPCAPFileRecord(
	_ context.Context,
	request *rpc.StopPCAPFileRecordRequest,
) (resp *rpc.StopPCAPFileRecordResponse, _ error) {
	resp = new(rpc.StopPCAPFileRecordResponse)
	resp.Removed = p.recorder.StopRecording(request.ConsumerKey) == nil
	return
}

func ipAddressesToBytes(addresses []net.IP) (result [][]byte) {
	for i := range addresses {
		result = append(result, addresses[i])
	}
	return
}
