package netflow

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	packetPolicyNameDrop = "drop"
	packetPolicyNamePass = "pass"
)

var (
	_                      encoding.TextUnmarshaler = (*PacketPolicy)(nil)
	_                      fmt.Stringer             = (*PacketPolicy)(nil)
	packetPolicyString2Val                          = map[string]PacketPolicy{
		packetPolicyNameDrop: PacketPolicyDrop,
		packetPolicyNamePass: PacketPolicyPass,
	}
	packetPolicyVal2String = map[PacketPolicy]string{
		PacketPolicyDrop: packetPolicyNameDrop,
		PacketPolicyPass: packetPolicyNamePass,
	}
)

func PacketPolicyDecodeHook() mapstructure.DecodeHookFuncType {
	return func(from reflect.Type, to reflect.Type, input any) (any, error) {
		if from.Kind() != reflect.String {
			return input, nil
		}

		if to != reflect.TypeOf(PacketPolicyDrop) {
			return input, nil
		}

		var pp PacketPolicy
		if v, ok := input.(string); !ok {
			return input, nil
		} else if err := pp.UnmarshalText([]byte(v)); err != nil {
			return nil, err
		}
		return pp, nil
	}
}

type PacketPolicy uint32

const (
	PacketPolicyDrop PacketPolicy = iota + 1
	PacketPolicyPass
)

func (pp *PacketPolicy) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*pp = PacketPolicyDrop
		return nil
	}

	if val, ok := packetPolicyString2Val[strings.ToLower(string(text))]; ok {
		*pp = val
	} else {
		*pp = PacketPolicyDrop
	}

	return nil
}

func (pp PacketPolicy) String() string {
	if val, ok := packetPolicyVal2String[pp]; ok {
		return val
	} else {
		return "unspecified"
	}
}

func (pp PacketPolicy) XDPAction() XDPAction {
	return XDPAction(pp)
}

func (pp PacketPolicy) RawValue() uint32 {
	return uint32(pp)
}
