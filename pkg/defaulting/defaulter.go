package defaulting

import (
	"reflect"
)

type Defaulter func(instance interface{})

type Registry interface {
	Register(t reflect.Type, defaulter ...Defaulter)
	Apply(instance interface{})
}

func New() Registry {
	return &registry{
		defaulters: make(map[reflect.Type][]Defaulter),
	}
}

type registry struct {
	defaulters map[reflect.Type][]Defaulter
}

func (r *registry) Register(t reflect.Type, defaulter ...Defaulter) {
	var given []Defaulter
	if r, ok := r.defaulters[t]; ok {
		given = r
	}

	for _, d := range defaulter {
		given = append(given, d)
	}

	r.defaulters[t] = given
}

func (r *registry) Apply(instance interface{}) {
	if defs, ok := r.defaulters[reflect.TypeOf(instance)]; ok {
		for _, def := range defs {
			def(instance)
		}
	}
}
