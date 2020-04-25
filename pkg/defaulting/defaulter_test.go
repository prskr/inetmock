package defaulting

import (
	"reflect"
	"testing"
)

func Test_registry_Apply(t *testing.T) {
	type sample struct {
		i int
	}
	type fields struct {
		defaulters map[reflect.Type][]Defaulter
	}
	type args struct {
		instance interface{}
	}
	type expect struct {
		result interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect expect
	}{
		{
			name: "Expect setting a sample value",
			fields: fields{
				defaulters: map[reflect.Type][]Defaulter{
					reflect.TypeOf(&sample{}): {func(instance interface{}) {
						switch i := instance.(type) {
						case *sample:
							i.i = 42
						}
					}},
				},
			},
			args: args{
				instance: &sample{},
			},
			expect: expect{
				result: &sample{
					i: 42,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &registry{
				defaulters: tt.fields.defaulters,
			}

			r.Apply(tt.args.instance)

			if !reflect.DeepEqual(tt.expect.result, tt.args.instance) {
				t.Errorf("Apply() expected = %v got %v", tt.args.instance, tt.expect.result)
			}
		})
	}
}

func Test_registry_Register(t *testing.T) {
	type sample struct {
	}
	type fields struct {
		defaulters map[reflect.Type][]Defaulter
	}
	type args struct {
		t         reflect.Type
		defaulter []Defaulter
	}
	type expect struct {
		length int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect expect
	}{
		{
			name: "",
			fields: fields{
				defaulters: make(map[reflect.Type][]Defaulter),
			},
			args: args{
				t: reflect.TypeOf(sample{}),
				defaulter: []Defaulter{func(instance interface{}) {

				}},
			},
			expect: expect{
				length: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &registry{
				defaulters: tt.fields.defaulters,
			}
			r.Register(tt.args.t, tt.args.defaulter...)

			if length := len(r.defaulters); length != tt.expect.length {
				t.Errorf("len(r.defaulters) expect %d got %d", tt.expect.length, length)
			}
		})
	}
}
