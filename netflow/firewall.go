package netflow

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"sync"

	manager "github.com/DataDog/ebpf-manager"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/features"
	"github.com/cilium/ebpf/rlimit"
)

var (
	ErrMissingMap                   = errors.New("map is missing")
	ErrCouldNotDetermineMonitorMode = errors.New("could not determine monitor mode")
)

var (
	//go:embed ebpf/firewall.o
	firewallEBPFProgram []byte
	memLockOnce         sync.Once
)

type XDPAction uint32

const (
	XDPActionDrop XDPAction = iota + 1
	XDPActionPass
)

type firewallMaps struct {
	Rules *Map[ConnIdent, FirewallRule]
}

func (m *firewallMaps) initMaps(mgr *manager.Manager) error {
	var err error
	if m.Rules, err = mapOfManager[ConnIdent, FirewallRule](mgr, "firewall_rules"); err != nil {
		return err
	}

	return nil
}

func NewFirewall(packetSink PacketSink, opts ...Option) *Firewall {
	fw := &Firewall{
		programLoader:     EBPFProgramBytesLoader(firewallEBPFProgram),
		packetSink:        packetSink,
		managedInterfaces: make(map[string]*FirewallInstance),
		errorSink:         noOpErrorSink,
	}

	for i := range opts {
		opts[i].ApplyTo(fw)
	}

	return fw
}

type Firewall struct {
	lock              sync.Mutex
	managedInterfaces map[string]*FirewallInstance
	monitorMode       MonitorMode
	errorSink         ErrorSink
	packetSink        PacketSink
	programLoader     EBPFProgramLoader
}

func (f *Firewall) enableMocking(toMock bool) {
	if toMock {
		f.monitorMode = monitorModeMock
	}
}

func (f *Firewall) SetEBPFProgramLoader(loader EBPFProgramLoader) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.programLoader = loader
}

func (f *Firewall) SetErrorSink(sink ErrorSink) {
	if sink == nil {
		return
	}

	f.lock.Lock()
	defer f.lock.Unlock()

	f.errorSink = sink

	for _, mi := range f.managedInterfaces {
		mi.SetErrorSink(sink)
	}
}

func (f *Firewall) AttachToInterface(interfaceName string, fwCfg FirewallInterfaceConfig) (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if _, ok := f.managedInterfaces[interfaceName]; ok {
		return nil
	}

	if fwCfg.RemoveMemLock {
		memLockOnce.Do(func() {
			err = rlimit.RemoveMemlock()
		})
		if err != nil {
			return err
		}
	}

	var (
		mgrOpts *manager.Options
		inst    = &FirewallInstance{
			errorSink: f.errorSink,
		}
	)

	if inst.mgr, mgrOpts, err = f.prepareManager(interfaceName, fwCfg); err != nil {
		return err
	}

	if err = inst.mgr.InitWithOptions(bytes.NewReader(firewallEBPFProgram), *mgrOpts); err != nil {
		return fmt.Errorf("failed to init manager: %w", err)
	}

	//nolint:gocritic // either way there will be some linter not satisfied
	if err = inst.initMaps(inst.mgr); err != nil {
		return err
	}

	if err = inst.mgr.Start(); err != nil {
		return fmt.Errorf("failed to start manager: %w", err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, inst.mgr.Stop(manager.CleanAll))
		}
	}()

	if err := inst.initPacketTransport(f.packetSink, f.monitorMode); err != nil {
		return err
	}

	if err := inst.SyncConfig(fwCfg); err != nil {
		return err
	}

	f.managedInterfaces[interfaceName] = inst

	return nil
}

func (f *Firewall) Close() error {
	f.lock.Lock()
	defer f.lock.Unlock()

	var err error
	for k, mi := range f.managedInterfaces {
		err = errors.Join(err, mi.Close())
		delete(f.managedInterfaces, k)
	}

	return err
}

func (f *Firewall) prepareManager(
	nicName string,
	fwCfg FirewallInterfaceConfig,
) (mgr *manager.Manager, mgrOpts *manager.Options, err error) {
	if f.monitorMode == MonitorModeUnspecified {
		switch err := features.HaveMapType(ebpf.RingBuf); {
		case errors.Is(err, ebpf.ErrNotSupported):
			f.monitorMode = MonitorModePerfEvent
		case err != nil && !errors.Is(err, ebpf.ErrNotSupported):
			return nil, nil, fmt.Errorf("failed to determine map feature compatibility: %w", err)
		default:
			f.monitorMode = MonitorModeRingBuf
		}
	}

	monitorProbeID := manager.ProbeIdentificationPair{
		EBPFSection:  f.monitorMode.Section(),
		EBPFFuncName: f.monitorMode.Function(),
	}

	mgr = &manager.Manager{
		Probes: []*manager.Probe{
			{
				ProbeIdentificationPair: monitorProbeID,
				IfName:                  nicName,
				NetworkDirection:        manager.Ingress,
			},
		},
	}

	mgrOpts = &manager.Options{
		ConstantEditors: []manager.ConstantEditor{
			{
				Name:                     "DEFAULT_POLICY",
				Value:                    fwCfg.DefaultPolicy.XDPAction(),
				BTFGlobalConstant:        true,
				FailOnMissing:            true,
				ProbeIdentificationPairs: []manager.ProbeIdentificationPair{monitorProbeID},
			},
			{
				Name:                     "EMIT_UNMATCHED",
				Value:                    fwCfg.Monitor,
				BTFGlobalConstant:        true,
				FailOnMissing:            true,
				ProbeIdentificationPairs: []manager.ProbeIdentificationPair{monitorProbeID},
			},
		},
	}

	switch f.monitorMode {
	case MonitorModeRingBuf:
		mgrOpts.ExcludedFunctions = append(mgrOpts.ExcludedFunctions, MonitorModePerfEvent.Function(), monitorModeMock.Function())
		mgrOpts.MapSpecEditors = map[string]manager.MapSpecEditor{
			"observed_packets": {
				Type:       ebpf.RingBuf,
				MaxEntries: 1 << 24,
				EditorFlag: manager.EditType | manager.EditMaxEntries,
			},
		}
		return mgr, mgrOpts, nil
	case MonitorModePerfEvent:
		mgrOpts.ExcludedFunctions = append(mgrOpts.ExcludedFunctions, MonitorModeRingBuf.Function(), monitorModeMock.Function())
		return mgr, mgrOpts, nil
	case monitorModeMock:
		mgrOpts.ExcludedFunctions = append(mgrOpts.ExcludedFunctions, MonitorModePerfEvent.Function(), MonitorModeRingBuf.Function())
		return mgr, mgrOpts, nil
	case MonitorModeUnspecified:
		fallthrough
	default:
		return nil, nil, ErrCouldNotDetermineMonitorMode
	}
}

type FirewallInstance struct {
	firewallMaps
	mgr       *manager.Manager
	transport *PacketTransport
	errorSink ErrorSink
}

func (i *FirewallInstance) SyncConfig(fwCfg FirewallInterfaceConfig) error {
	for idx := range fwCfg.Rules {
		entry := fwCfg.Rules[idx]

		id := ConnIdent{
			Addr:      entry.Destination.NetIP(),
			Port:      entry.Destination.Port(),
			Transport: entry.Destination.Protocol(),
		}

		rule := FirewallRule{
			Policy:         entry.Policy.XDPAction(),
			MonitorTraffic: entry.MonitorTraffic(fwCfg.Monitor),
		}

		if err := i.Rules.Put(id, rule); err != nil {
			return err
		}
	}

	return nil
}

func (i *FirewallInstance) SetErrorSink(sink ErrorSink) {
	i.errorSink = sink
}

func (i *FirewallInstance) Close() (err error) {
	if i.transport != nil {
		err = i.transport.Close()
	}

	return errors.Join(err, i.mgr.Stop(manager.CleanAll))
}

func (i *FirewallInstance) initPacketTransport(packetSink PacketSink, mode MonitorMode) error {
	var (
		transportMap *ebpf.Map
		err          error
	)

	if transportMap, err = getMapFromManager(i.mgr, "observed_packets"); err != nil {
		return err
	}

	var reader PacketReader
	switch mode {
	case MonitorModeRingBuf:
		if reader, err = NewRingBufReader(transportMap); err != nil {
			return err
		}
	case MonitorModePerfEvent:
		const defaultPerCPUBufferSize = 8
		if reader, err = NewPerfEventReader(transportMap, defaultPerCPUBufferSize); err != nil {
			return err
		}
	case monitorModeMock:
		return nil
	case MonitorModeUnspecified:
		fallthrough
	default:
		return ErrCouldNotDetermineMonitorMode
	}

	i.transport = NewPacketTransport(reader, packetSink, i.errorSink)
	go i.transport.Start()

	return nil
}
