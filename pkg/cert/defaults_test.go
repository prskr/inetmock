package cert

import (
	"reflect"
	"testing"
)

func Test_certOptionsDefaulter(t *testing.T) {
	tests := []struct {
		name     string
		arg      GenerationOptions
		expected GenerationOptions
	}{
		{
			name: "",
			arg: GenerationOptions{
				CommonName: "CA",
			},
			expected: GenerationOptions{
				CommonName:    "CA",
				Country:       []string{"US"},
				Locality:      []string{"San Francisco"},
				Organization:  []string{"INetMock"},
				StreetAddress: []string{"Golden Gate Bridge"},
				PostalCode:    []string{"94016"},
				Province:      []string{""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			certOptionsDefaulter(&tt.arg)
			if !reflect.DeepEqual(tt.expected, tt.arg) {
				t.Errorf("Apply defaulter expected=%v got=%v", tt.expected, tt.arg)
			}
		})
	}
}
