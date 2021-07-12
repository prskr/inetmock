package endpoint_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/mitchellh/mapstructure"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

type greeter interface {
	Greet(name string)
}

type friendlyGreeter struct {
	ID int
}

func (f friendlyGreeter) Greet(name string) {
	fmt.Printf("Hello %s - nice to meet you!\n", name)
}

type anotherGreeter struct {
	Insult string
}

func (a anotherGreeter) Greet(name string) {
	fmt.Printf("Hi %s - %s\n", name, a.Insult)
}

type testOption struct {
	Greeting string
	Greeter  greeter
}

func Test_OptionByTypeDecoderBuilder_DecodeHook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mappings map[string]reflect.Type
		input    interface{}
		want     interface{}
	}{
		{
			name: "Single mapping",
			mappings: map[string]reflect.Type{
				"friendly": reflect.TypeOf(friendlyGreeter{}),
			},
			input: map[string]interface{}{
				"greeting": "Tom",
				"greeter": map[string]interface{}{
					"type": "friendly",
					"id":   1234,
				},
			},
			want: td.Struct(new(friendlyGreeter), td.StructFields{
				"ID": 1234,
			}),
		},
		{
			name: "Simple mapping with multiple mappings",
			mappings: map[string]reflect.Type{
				"friendly":  reflect.TypeOf(friendlyGreeter{}),
				"insulting": reflect.TypeOf(anotherGreeter{}),
			},
			input: map[string]interface{}{
				"greeting": "Tom",
				"greeter": map[string]interface{}{
					"type":   "insulting",
					"Insult": "now go and fuck yourself!",
				},
			},
			want: td.Struct(new(anotherGreeter), td.StructFields{
				"Insult": "now go and fuck yourself!",
			}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				g greeter
				b = endpoint.OptionByTypeDecoderBuilder{
					OptionType: reflect.TypeOf(&g).Elem(),
					Mappings:   tt.mappings,
				}
				decoderHook   = b.Build()
				out           = new(testOption)
				decoderConfig = &mapstructure.DecoderConfig{
					DecodeHook: decoderHook,
					Result:     out,
				}
				decoder *mapstructure.Decoder
			)

			if d, err := mapstructure.NewDecoder(decoderConfig); err != nil {
				t.Errorf("mapstructure.NewDecoder() error = %v", err)
				return
			} else {
				decoder = d
			}

			if err := decoder.Decode(tt.input); err != nil {
				t.Errorf("Decode() error = %v", err)
				return
			}

			td.Cmp(t, out.Greeter, tt.want)
		})
	}
}
