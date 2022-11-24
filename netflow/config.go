package netflow

import (
	"bytes"
	"io"
	"net/netip"
	"time"
)

type (
	EBPFProgramLoader interface {
		LoadProgram() io.ReaderAt
	}

	Service interface {
		SetErrorSink(sink ErrorSink)
		SetEBPFProgramLoader(loader EBPFProgramLoader)
		enableMocking(toMock bool)
	}

	Option interface {
		ApplyTo(svc Service)
	}
)

type EBPFProgramBytesLoader []byte

func (l EBPFProgramBytesLoader) LoadProgram() io.ReaderAt {
	return bytes.NewReader(l)
}

type OptionFunc func(svc Service)

func (f OptionFunc) ApplyTo(svc Service) {
	f(svc)
}

func WithErrorSink(sink ErrorSink) Option {
	return OptionFunc(func(svc Service) {
		svc.SetErrorSink(sink)
	})
}

type ErrorSinkOption struct {
	ErrorSink
}

func (s ErrorSinkOption) ApplyTo(svc Service) {
	svc.SetErrorSink(s.ErrorSink)
}

func (s ErrorSinkOption) Apply(opt *epochOptions) {
	opt.ErrorHandler = s.ErrorSink
}

func WithMockingEnabled(mockingEnabled bool) Option {
	return OptionFunc(func(svc Service) {
		svc.enableMocking(mockingEnabled)
	})
}

func BoolP(val bool) *bool {
	return &val
}

type RuleEntry struct {
	Policy      PacketPolicy
	Destination IPPortProto `mapstructure:"dest"`
	Monitor     *bool
}

func (e RuleEntry) MonitorTraffic(defaultValue bool) bool {
	if e.Monitor == nil {
		return defaultValue
	}
	return *e.Monitor
}

type FirewallInterfaceConfig struct {
	RemoveMemLock bool
	DefaultPolicy PacketPolicy
	Monitor       bool
	Rules         []RuleEntry
}

type NATTargetSpec struct {
	Destination IPPortProto `mapstructure:"dest"`
	RedirectTo  NATTarget
	TranslateTo netip.Addr `mapstructure:"translateTo"`
}

type ConnTrackConfig struct {
	HighWaterMark float64
	CleanupWindow time.Duration
}

type NATTableSpec struct {
	ConnTrack    ConnTrackConfig
	Translations []NATTargetSpec
}
