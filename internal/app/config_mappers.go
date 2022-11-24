package app

import (
	"net/netip"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func NetIPDecodeHook() mapstructure.DecodeHookFuncType {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() != reflect.String {
			return data, nil
		}

		var (
			result          netip.Addr
			pointerExpected bool
			value           string
		)

		targetType := to
		if targetType.Kind() == reflect.Pointer {
			targetType = targetType.Elem()
			pointerExpected = true
		}

		if targetType != reflect.TypeOf(result) {
			return data, nil
		}

		if v, ok := data.(string); !ok {
			return data, nil
		} else {
			value = v
		}

		if addr, err := netip.ParseAddr(value); err != nil {
			return nil, err
		} else if pointerExpected {
			return &addr, nil
		} else {
			return addr, nil
		}
	}
}
