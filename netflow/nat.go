package netflow

import (
	_ "embed"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sync"
	"time"

	manager "github.com/DataDog/ebpf-manager"
	"github.com/google/uuid"
)

//go:embed ebpf/nat.o
var natEBPFProgram []byte

const (
	DefaultConnTrackCleanupWindow = 5 * time.Second
)

type (
	NATConfigKey uint32
	natMaps      struct {
		Config *Map[NATConfigKey, uint32]
		Table  *Map[ConnIdent, ConnMeta]
	}
)

const (
	natConfigKeyCurrentEpoch NATConfigKey = iota
)

func NewNAT(opts ...Option) *NAT {
	n := &NAT{
		programLoader:     EBPFProgramBytesLoader(natEBPFProgram),
		managedInterfaces: make(map[string]*NATInstance),
		errorSink:         noOpErrorSink,
		egressProbeID: manager.ProbeIdentificationPair{
			UID:          uuid.NewString(),
			EBPFSection:  "classifier/egress",
			EBPFFuncName: "egress",
		},
		ingressProbeID: manager.ProbeIdentificationPair{
			UID:          uuid.NewString(),
			EBPFSection:  "classifier/ingress",
			EBPFFuncName: "ingress",
		},
	}

	for i := range opts {
		opts[i].ApplyTo(n)
	}

	return n
}

type NATInstance struct {
	natMaps
	interfaceAddr    netip.Addr
	mgr              *manager.Manager
	epoch            *Epoch
	connTrackCleaner *ConnTrackCleaner
	errorSink        ErrorSink
}

func (n *NATInstance) SetErrorSink(sink ErrorSink) {
	n.errorSink = sink
	n.epoch.ErrorHandler = sink
	n.connTrackCleaner.ErrorHandler = sink
}

func (n *NATInstance) Sync(targets []NATTargetSpec) error {
	desired := make(map[ConnIdent]ConnMeta, len(targets))
	for i := range targets {
		target := targets[i]

		from := ConnIdent{
			Addr:      target.Destination.NetIP(),
			Port:      target.Destination.Port(),
			Transport: target.Destination.Protocol(),
		}

		var to ConnMeta
		switch target.RedirectTo {
		case NATTargetIP:
			to = ConnMeta{
				Addr: target.TranslateTo,
			}

		case NATTargetInterface:
			to = ConnMeta{
				Addr: n.interfaceAddr,
			}
		}

		desired[from] = to
	}

	return n.Table.Sync(desired)
}

func (n *NATInstance) Close() error {
	n.epoch.Stop()
	n.connTrackCleaner.Stop()
	return n.mgr.Stop(manager.CleanAll)
}

func (n *NATInstance) initMaps() (err error) {
	if n.Config, err = mapOfManager[NATConfigKey, uint32](n.mgr, "nat_config"); err != nil {
		return err
	}

	if n.Table, err = mapOfManager[ConnIdent, ConnMeta](n.mgr, "nat_translations"); err != nil {
		return err
	}

	return nil
}

func (n *NATInstance) prepareConnTrackCleaner(interfaceName string, spec NATTableSpec) (err error) {
	const (
		waterMarkMinimum = 0.1
		waterMarkDefault = 0.7
	)

	var connTrackMap *Map[ConnIdent, ConnMeta]
	if connTrackMap, err = mapOfManager[ConnIdent, ConnMeta](n.mgr, "conn_track"); err != nil {
		return err
	}

	connTrackHighWaterMark := spec.ConnTrack.HighWaterMark
	if connTrackHighWaterMark < waterMarkMinimum {
		connTrackHighWaterMark = waterMarkDefault
	}

	n.connTrackCleaner = NewConnTrackCleaner(connTrackMap, n.errorSink, connTrackHighWaterMark, interfaceName)

	connTrackCleanupWindow := spec.ConnTrack.CleanupWindow
	if connTrackCleanupWindow == 0 {
		connTrackCleanupWindow = DefaultConnTrackCleanupWindow
	}

	return n.connTrackCleaner.Start(connTrackCleanupWindow)
}

type NAT struct {
	lock              sync.Mutex
	errorSink         ErrorSink
	programLoader     EBPFProgramLoader
	managedInterfaces map[string]*NATInstance
	ingressProbeID    manager.ProbeIdentificationPair
	egressProbeID     manager.ProbeIdentificationPair
}

func (n *NAT) enableMocking(toMock bool) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if toMock {
		probeID := manager.ProbeIdentificationPair{
			UID:          uuid.New().String(),
			EBPFSection:  "classifier/mock",
			EBPFFuncName: "nat_mock",
		}
		n.ingressProbeID = probeID
		n.egressProbeID = probeID
	}
}

func (n *NAT) SetEBPFProgramLoader(loader EBPFProgramLoader) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.programLoader = loader
}

func (n *NAT) SetErrorSink(sink ErrorSink) {
	if sink == nil {
		return
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	n.errorSink = sink
	for _, inst := range n.managedInterfaces {
		inst.SetErrorSink(sink)
	}
}

func (n *NAT) AttachToInterface(interfaceName string, spec NATTableSpec) (err error) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.managedInterfaces[interfaceName]; ok {
		return nil
	}

	var inst *NATInstance

	if inst, err = n.prepareInstance(interfaceName, n.programLoader); err != nil {
		return err
	}

	if err := inst.initMaps(); err != nil {
		return err
	}

	inst.epoch = NewEpoch(inst.Config, ErrorSinkOption{ErrorSink: inst.errorSink})

	if err := inst.mgr.Start(); err != nil {
		return fmt.Errorf("failed to start NAT manager: %w", err)
	}

	if err := inst.prepareConnTrackCleaner(interfaceName, spec); err != nil {
		return err
	}

	if err := inst.epoch.StartSync(DefaultEpochSyncWindow); err != nil {
		return err
	}

	if err := inst.Sync(spec.Translations); err != nil {
		return err
	}

	n.managedInterfaces[interfaceName] = inst

	return nil
}

func (n *NAT) Close() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	var err error
	for key, inst := range n.managedInterfaces {
		err = errors.Join(err, inst.Close())
		delete(n.managedInterfaces, key)
	}

	return err
}

func (n *NAT) prepareInstance(interfaceName string, loader EBPFProgramLoader) (inst *NATInstance, err error) {
	var primaryAddr netip.Addr
	if primaryAddr, err = determinePrimaryNICAddr(interfaceName); err != nil {
		return nil, err
	}

	inst = &NATInstance{
		interfaceAddr: primaryAddr,
		mgr:           new(manager.Manager),
	}

	probes := []*manager.Probe{
		{
			ProbeIdentificationPair: n.ingressProbeID,
			IfName:                  interfaceName,
			NetworkDirection:        manager.Ingress,
			KeepProgramSpec:         true,
		},
		{
			ProbeIdentificationPair: n.egressProbeID,
			IfName:                  interfaceName,
			NetworkDirection:        manager.Egress,
			KeepProgramSpec:         true,
		},
	}

	requireClone := n.egressProbeID.Matches(n.ingressProbeID)

	if requireClone {
		inst.mgr.Probes = probes[:1]
	} else {
		inst.mgr.Probes = probes
	}

	mgrOpts := manager.Options{
		ConstantEditors: []manager.ConstantEditor{
			{
				Name:              "INTERFACE_IP",
				Value:             ipAddr2int(primaryAddr),
				BTFGlobalConstant: true,
				FailOnMissing:     true,
			},
		},
	}

	if err := inst.mgr.InitWithOptions(loader.LoadProgram(), mgrOpts); err != nil {
		return nil, err
	}

	if requireClone {
		clonedProbe := manager.Probe{
			ProbeIdentificationPair: manager.ProbeIdentificationPair{
				UID:          uuid.NewString(),
				EBPFFuncName: probes[0].EBPFFuncName,
			},
			IfName:           interfaceName,
			NetworkDirection: manager.Egress,
		}
		if err := inst.mgr.CloneProgram(probes[0].UID, &clonedProbe, nil, nil); err != nil {
			return nil, err
		}
	}

	return inst, nil
}

func determinePrimaryNICAddr(nicName string) (primaryAddr netip.Addr, err error) {
	var nic *net.Interface
	if nic, err = net.InterfaceByName(nicName); err != nil {
		return netip.Addr{}, err
	}

	if addrs, err := nic.Addrs(); err != nil {
		return netip.Addr{}, err
	} else {
		for i := range addrs {
			if ipn, ok := addrs[i].(*net.IPNet); ok {
				if ipv4, ok := netip.AddrFromSlice(ipn.IP.To4()); ok {
					return ipv4, nil
				}
			}
		}
	}

	return netip.Addr{}, fmt.Errorf("no IPv4 address found for interface %s", nicName)
}
