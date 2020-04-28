package plugins

import (
	"github.com/baez90/inetmock/pkg/api"
	"reflect"
	"testing"
)

func Test_handlerRegistry_HandlerForName(t *testing.T) {
	type fields struct {
		handlers map[string]api.PluginInstanceFactory
	}
	type args struct {
		handlerName string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantInstance api.ProtocolHandler
		wantOk       bool
	}{
		{
			name:         "No instance if nothing is registered",
			fields:       fields{},
			args:         args{},
			wantInstance: nil,
			wantOk:       false,
		},
		{
			name: "Nil instance from pseudo factory",
			fields: fields{
				handlers: map[string]api.PluginInstanceFactory{
					"pseudo": func() api.ProtocolHandler {
						return nil
					},
				},
			},
			args: args{
				handlerName: "pseudo",
			},
			wantInstance: nil,
			wantOk:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handlerRegistry{
				handlers: tt.fields.handlers,
			}
			gotInstance, gotOk := h.HandlerForName(tt.args.handlerName)
			if !reflect.DeepEqual(gotInstance, tt.wantInstance) {
				t.Errorf("HandlerForName() gotInstance = %v, want %v", gotInstance, tt.wantInstance)
			}
			if gotOk != tt.wantOk {
				t.Errorf("HandlerForName() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
