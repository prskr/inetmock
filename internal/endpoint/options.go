package endpoint

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Mapping interface {
	MapTo(in any) (any, error)
}

type MappingFunc func(in any) (any, error)

func (m MappingFunc) MapTo(in any) (any, error) {
	return m(in)
}

type OptionByTypeDecoderBuilder struct {
	OptionType reflect.Type
	Mappings   map[string]Mapping
}

func NewOptionByTypeDecoderBuilderFor(opt any) OptionByTypeDecoderBuilder {
	return OptionByTypeDecoderBuilder{
		OptionType: reflect.TypeOf(opt).Elem(),
		Mappings:   make(map[string]Mapping),
	}
}

func (o *OptionByTypeDecoderBuilder) AddMappingToType(typeName string, targetType reflect.Type) *OptionByTypeDecoderBuilder {
	o.Mappings[typeName] = MappingFunc(func(in any) (any, error) {
		instance := reflect.New(targetType).Interface()
		return instance, mapstructure.Decode(in, instance)
	})
	return o
}

func (o *OptionByTypeDecoderBuilder) AddMappingToProvider(typeName string, provider func() any) *OptionByTypeDecoderBuilder {
	o.Mappings[typeName] = MappingFunc(func(data any) (any, error) {
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
	return func(from reflect.Type, to reflect.Type, data any) (any, error) {
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
