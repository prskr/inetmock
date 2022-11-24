package netflow

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	natTargetNameUnspecified = "unspecified"
	natTargetNameInterface   = "interface"
	natTargetNameIP          = "ip"
)

var (
	_                   encoding.TextUnmarshaler = (*NATTarget)(nil)
	_                   fmt.Stringer             = (*NATTarget)(nil)
	natTargetString2Val                          = map[string]NATTarget{
		natTargetNameInterface: NATTargetInterface,
		natTargetNameIP:        NATTargetIP,
	}
	natTargetVal2String = map[NATTarget]string{
		NATTargetInterface: natTargetNameInterface,
		NATTargetIP:        natTargetNameIP,
	}
)

func NATTargetDecodingHook() mapstructure.DecodeHookFuncType {
	return func(from reflect.Type, to reflect.Type, input interface{}) (interface{}, error) {
		if from.Kind() != reflect.String {
			return input, nil
		}

		if to != reflect.TypeOf(NATTargetInterface) {
			return input, nil
		}

		var t NATTarget
		if val, ok := input.(string); !ok {
			return input, nil
		} else if err := t.UnmarshalText([]byte(val)); err != nil {
			return nil, err
		}

		return t, nil
	}
}

type NATTarget uint32

const (
	NATTargetInterface NATTarget = iota
	NATTargetIP
)

func (n NATTarget) String() string {
	if val, ok := natTargetVal2String[n]; ok {
		return val
	}
	return natTargetNameUnspecified
}

func (n *NATTarget) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*n = NATTargetInterface
		return nil
	}

	if val, ok := natTargetString2Val[strings.ToLower(string(text))]; ok {
		*n = val
	} else {
		*n = NATTargetInterface
	}

	return nil
}
