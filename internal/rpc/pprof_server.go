package rpc

import (
	"context"
	"io"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"time"

	"github.com/valyala/bytebufferpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rpcv1 "gitlab.com/inetmock/inetmock/pkg/rpc/v1"
)

var (
	profilesApplicableToGCRequest = []string{
		"heap",
	}
	cpuProfileCollector = profileDataCollector{
		InitDelegate: pprof.StartCPUProfile,
		StopDelegate: pprof.StopCPUProfile,
	}
	traceCollector = profileDataCollector{
		InitDelegate: trace.Start,
		StopDelegate: trace.Stop,
	}
)

func NewProfilingServer() rpcv1.ProfilingServiceServer {
	return profilingServer{}
}

type profilingServer struct {
	rpcv1.UnimplementedProfilingServiceServer
}

func (profilingServer) ProfileDump(_ context.Context, req *rpcv1.ProfileDumpRequest) (resp *rpcv1.ProfileDumpResponse, err error) {
	profile := pprof.Lookup(req.ProfileName)
	if profile == nil {
		return nil, status.Errorf(codes.NotFound, "Profile of name %s was not found", req.ProfileName)
	}

	if req.GcBeforeDump {
		gcIfRequired(req.ProfileName)
	}

	outBuffer := bytebufferpool.Get()
	defer bytebufferpool.Put(outBuffer)
	if err := profile.WriteTo(outBuffer, int(req.Debug)); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to collect profile: %v", err)
	}

	return &rpcv1.ProfileDumpResponse{
		ProfileData: outBuffer.Bytes(),
	}, nil
}

func (profilingServer) CPUProfile(ctx context.Context, req *rpcv1.CPUProfileRequest) (resp *rpcv1.CPUProfileResponse, err error) {
	profileDuration := req.ProfileDuration.AsDuration()
	if err = durationExceedsTimeout(ctx, profileDuration); err != nil {
		return
	}

	if data, err := cpuProfileCollector.Collect(profileDuration); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to start CPU profiling: %v", err)
	} else {
		return &rpcv1.CPUProfileResponse{ProfileData: data}, nil
	}
}

func (profilingServer) Trace(ctx context.Context, req *rpcv1.TraceRequest) (resp *rpcv1.TraceResponse, err error) {
	traceDuration := req.TraceDuration.AsDuration()
	if err = durationExceedsTimeout(ctx, traceDuration); err != nil {
		return
	}

	if data, err := traceCollector.Collect(traceDuration); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to start CPU tracing: %v", err)
	} else {
		return &rpcv1.TraceResponse{ProfileData: data}, nil
	}
}

func durationExceedsTimeout(ctx context.Context, duration time.Duration) error {
	if deadline, set := ctx.Deadline(); set && deadline.Before(time.Now().Add(duration)) {
		return status.Error(codes.InvalidArgument, "The profile duration exceeds the requests timeout")
	}
	return nil
}

func gcIfRequired(profileName string) {
	if idx := sort.SearchStrings(profilesApplicableToGCRequest, profileName); profilesApplicableToGCRequest[idx] == profileName {
		runtime.GC()
	}
}

type profileDataCollector struct {
	InitDelegate func(writer io.Writer) error
	StopDelegate func()
}

func (c profileDataCollector) Collect(profileDuration time.Duration) ([]byte, error) {
	outBuffer := bytebufferpool.Get()
	defer bytebufferpool.Put(outBuffer)

	if err := c.InitDelegate(outBuffer); err != nil {
		return nil, err
	}
	<-time.After(profileDuration)
	c.StopDelegate()
	return outBuffer.Bytes(), nil
}
