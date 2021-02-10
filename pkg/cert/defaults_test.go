package cert

import (
	"reflect"
	"testing"
)

func Test_certOptionsDefaulter(t *testing.T) {
	type testCase struct {
		name     string
		arg      GenerationOptions
		expected GenerationOptions
	}
	tests := []testCase{
		{
			name: "Empty options",
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
		{
			name: "Options with country set",
			arg: GenerationOptions{
				CommonName: "CA",
				Country:    []string{"DE"},
			},
			expected: GenerationOptions{
				CommonName:    "CA",
				Country:       []string{"DE"},
				Locality:      []string{"San Francisco"},
				Organization:  []string{"INetMock"},
				StreetAddress: []string{"Golden Gate Bridge"},
				PostalCode:    []string{"94016"},
				Province:      []string{""},
			},
		},
		{
			name: "Options with organization set set",
			arg: GenerationOptions{
				CommonName:   "CA",
				Organization: []string{"inetmock"},
			},
			expected: GenerationOptions{
				CommonName:    "CA",
				Country:       []string{"US"},
				Locality:      []string{"San Francisco"},
				Organization:  []string{"inetmock"},
				StreetAddress: []string{"Golden Gate Bridge"},
				PostalCode:    []string{"94016"},
				Province:      []string{""},
			},
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			if err := applyDefaultGenerationOptions(&tt.arg); err != nil {
				t.Errorf("applyDefaultGenerationOptions() error = %v", err)
			}
			if !reflect.DeepEqual(tt.expected, tt.arg) {
				t.Errorf("Apply defaulter expected=%v got=%v", tt.expected, tt.arg)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
