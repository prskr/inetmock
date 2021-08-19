package cert

import (
	"github.com/imdario/mergo"
)

var defaultOptions = &GenerationOptions{
	Country:       []string{"US"},
	Locality:      []string{"San Francisco"},
	Organization:  []string{"INetMock"},
	StreetAddress: []string{"Golden Gate Bridge"},
	PostalCode:    []string{"94016"},
	Province:      []string{""},
}

func applyDefaultGenerationOptions(opts *GenerationOptions) error {
	return mergo.Merge(opts, defaultOptions)
}
