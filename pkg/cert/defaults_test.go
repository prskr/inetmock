//nolint:testpackage // testing internals here - needs to be private
package cert

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func Test_certOptionsDefaulter(t *testing.T) {
	t.Parallel()
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
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := applyDefaultGenerationOptions(&tt.arg); err != nil {
				t.Errorf("applyDefaultGenerationOptions() error = %v", err)
			}
			td.Cmp(t, tt.arg, tt.expected)
		})
	}
}
