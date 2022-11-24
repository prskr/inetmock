package netflow

import (
	"net"
	"net/netip"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var netIPZero = netip.AddrFrom4([4]byte{0, 0, 0, 0})

func IPPortProtoDecodeHook() mapstructure.DecodeHookFuncType {
	return func(from reflect.Type, to reflect.Type, input any) (any, error) {
		if from.Kind() != reflect.String {
			return input, nil
		}

		if to != reflect.TypeOf(IPPortProto("")) {
			return input, nil
		}

		if val, ok := input.(string); !ok {
			return input, nil
		} else {
			return IPPortProto(val), nil
		}
	}
}

type IPPortProto string

func (p IPPortProto) NetIP() netip.Addr {
	val := string(p)
	if len(val) == 0 {
		return netIPZero
	}

	if ipVal, _, found := strings.Cut(val, ":"); !found {
		return netIPZero
	} else if addr, err := netip.ParseAddr(ipVal); err != nil {
		return netIPZero
	} else {
		return addr
	}
}

func (p IPPortProto) IP() net.IP {
	val := string(p)
	if len(val) == 0 {
		return nil
	}

	if ipVal, _, found := strings.Cut(val, ":"); !found {
		return net.IPv4zero
	} else {
		return net.ParseIP(ipVal)
	}
}

func (p IPPortProto) IsWildcardIP() bool {
	return p.IP().Equal(net.IPv4zero)
}

func (p IPPortProto) Protocol() Protocol {
	val := string(p)
	if len(val) == 0 {
		return ProtocolUnspecified
	}

	if ipSplitIdx := strings.Index(val, ":"); ipSplitIdx != -1 {
		val = val[ipSplitIdx+1:]
	}

	var proto Protocol
	if _, protoVal, found := strings.Cut(val, "/"); !found {
		return ProtocolUnspecified
	} else if err := proto.UnmarshalText([]byte(protoVal)); err != nil {
		return ProtocolUnspecified
	} else {
		return proto
	}
}

func (p IPPortProto) Port() uint16 {
	const (
		baseDecimal   = 10
		uint16BitSize = 16
	)
	val := string(p)
	if len(val) == 0 {
		return 0
	}

	if ipSplitIdx := strings.Index(val, ":"); ipSplitIdx != -1 {
		val = val[ipSplitIdx+1:]
	}

	if portVal, _, found := strings.Cut(val, "/"); !found {
		return 0
	} else if i, err := strconv.ParseUint(portVal, baseDecimal, uint16BitSize); err != nil {
		return 0
	} else {
		return uint16(i)
	}
}
