package cert

import (
	"github.com/baez90/inetmock/pkg/defaulting"
	"reflect"
)

var (
	certOptionsDefaulter defaulting.Defaulter = func(instance interface{}) {
		switch o := instance.(type) {
		case *GenerationOptions:

			if len(o.Country) < 1 {
				o.Country = []string{"US"}
			}

			if len(o.Locality) < 1 {
				o.Locality = []string{"San Francisco"}
			}

			if len(o.Organization) < 1 {
				o.Organization = []string{"INetMock"}
			}

			if len(o.StreetAddress) < 1 {
				o.StreetAddress = []string{"Golden Gate Bridge"}
			}

			if len(o.PostalCode) < 1 {
				o.PostalCode = []string{"94016"}
			}

			if len(o.Province) < 1 {
				o.Province = []string{""}
			}
		}
	}
)

func init() {
	certOptionsType := reflect.TypeOf(GenerationOptions{})
	defaulters.Register(certOptionsType, certOptionsDefaulter)
}
