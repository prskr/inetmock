package endpoint

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Mapping interface {
	MapTo(in interface{}) (interface{}, error)
}

type MappingFunc func(in interface{}) (interface{}, error)

func (m MappingFunc) MapTo(in interface{}) (interface{}, error) {
	return m(in)
}

type OptionByTypeDecoderBuilder struct {
	OptionType reflect.Type
	Mappings   map[string]Mapping
}

func NewOptionByTypeDecoderBuilderFor(opt interface{}) OptionByTypeDecoderBuilder {
	return OptionByTypeDecoderBuilder{
		OptionType: reflect.TypeOf(opt).Elem(),
		Mappings:   make(map[string]Mapping),
	}
}

func (o *OptionByTypeDecoderBuilder) AddMappingToType(typeName string, targetType reflect.Type) *OptionByTypeDecoderBuilder {
	o.Mappings[typeName] = MappingFunc(func(in interface{}) (interface{}, error) {
		instance := reflect.New(targetType).Interface()
		return instance, mapstructure.Decode(in, instance)
	})
	return o
}

func (o *OptionByTypeDecoderBuilder) AddMappingToProvider(typeName string, provider func() interface{}) *OptionByTypeDecoderBuilder {
	o.Mappings[typeName] = MappingFunc(func(data interface{}) (interface{}, error) {
		instance := provider()
		return instance, mapstructure.Decode(data, instance)
	})
	return o
}

func (o *OptionByTypeDecoderBuilder) AddMappingToMapper(
	typeName string,
	mapper Mapping,
) *OptionByTypeDecoderBuilder {
	o.Mappings[typeName] = mapper
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

		if mapper, registered := o.Mappings[typeRef.Type]; registered {
			if instance, err := mapper.MapTo(data); err != nil {
				return nil, err
			} else {
				return instance, nil
			}
		}
		return data, nil
	}
}
