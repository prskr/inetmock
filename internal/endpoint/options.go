package endpoint

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type OptionByTypeDecoderBuilder struct {
	OptionType reflect.Type
	Mappings   map[string]reflect.Type
}

func (o *OptionByTypeDecoderBuilder) AddMapping(typeName string, targetType reflect.Type) *OptionByTypeDecoderBuilder {
	o.Mappings[typeName] = targetType
	return o
}

func (o OptionByTypeDecoderBuilder) Build() mapstructure.DecodeHookFuncType {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() != reflect.Map {
			return data, nil
		}
		if to != o.OptionType {
			return data, nil
		}

		typeRef := &struct {
			Type string
		}{}

		if err := mapstructure.Decode(data, typeRef); err != nil {
			return nil, err
		}

		if targetType, registered := o.Mappings[typeRef.Type]; registered {
			instance := reflect.New(targetType).Interface()
			if err := mapstructure.Decode(data, instance); err != nil {
				return nil, err
			}
			return instance, nil
		}
		return data, nil
	}
}
