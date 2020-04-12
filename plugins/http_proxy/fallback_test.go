package main

import (
	"reflect"
	"testing"
)

func TestStrategyForName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{
			name: "Test get notfound strategy",
			want: reflect.TypeOf(&notFoundFallbackStrategy{}),
			args: args{
				name: "notfound",
			},
		},
		{
			name: "Test get pass through strategy",
			want: reflect.TypeOf(&passThroughFallbackStrategy{}),
			args: args{
				name: "passthrough",
			},
		},
		{
			name: "Test get fallback strategy notfound because key is not known",
			want: reflect.TypeOf(&notFoundFallbackStrategy{}),
			args: args{
				name: "asdf12234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrategyForName(tt.args.name); reflect.TypeOf(got) != tt.want {
				t.Errorf("StrategyForName() = %v, want %v", got, tt.want)
			}
		})
	}
}
