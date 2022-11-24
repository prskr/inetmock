package netflow

import (
	"encoding/binary"
	"fmt"
)

const (
	fwRuleBinarySize = 8
)

type FirewallRule struct {
	Policy         XDPAction
	MonitorTraffic bool
}

func (f *FirewallRule) UnmarshalBinary(data []byte) error {
	if dataLen := len(data); dataLen < fwRuleBinarySize {
		return fmt.Errorf("expected %d bytes but got %d", fwRuleBinarySize, dataLen)
	}

	f.Policy = XDPAction(binary.LittleEndian.Uint32(data[:4]))
	f.MonitorTraffic = binary.LittleEndian.Uint32(data[4:]) > 0

	return nil
}

func (f FirewallRule) MarshalBinary() (data []byte, err error) {
	data = make([]byte, fwRuleBinarySize)

	if err := f.MarshalBinaryTo(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (f FirewallRule) MarshalBinaryTo(data []byte) (err error) {
	if l := len(data); l != fwRuleBinarySize {
		return fmt.Errorf("rule does not marshal to %d bytes", l)
	}

	binary.LittleEndian.PutUint32(data[:4], uint32(f.Policy))
	if f.MonitorTraffic {
		binary.LittleEndian.PutUint32(data[4:], 1)
	}

	return nil
}

//nolint:unused
type firewallRuleCollection []FirewallRule

//nolint:unused
func (f firewallRuleCollection) MarshalBinary() (data []byte, err error) {
	data = make([]byte, len(f)*fwRuleBinarySize)

	for i := range f {
		if err := f[i].MarshalBinaryTo(data[i*fwRuleBinarySize : (i+1)*fwRuleBinarySize]); err != nil {
			return nil, err
		}
	}

	return data, nil
}
