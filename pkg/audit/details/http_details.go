package details

import (
	"net/http"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"
)

type HTTP struct {
	Method  string
	Host    string
	URI     string
	Proto   string
	Headers http.Header
}

func NewHTTPFromWireFormat(entity *HTTPDetailsEntity) HTTP {
	headers := http.Header{}
	for name, values := range entity.Headers {
		for idx := range values.Values {
			headers.Add(name, values.Values[idx])
		}
	}

	return HTTP{
		Method:  entity.Method.String(),
		Host:    entity.Host,
		URI:     entity.Uri,
		Proto:   entity.Proto,
		Headers: headers,
	}
}

func (d HTTP) MarshalToWireFormat() (any *anypb.Any, err error) {
	var method = HTTPMethod_GET
	if methodValue, known := HTTPMethod_value[strings.ToUpper(d.Method)]; known {
		method = HTTPMethod(methodValue)
	}

	headers := make(map[string]*HTTPHeaderValue)

	for k, v := range d.Headers {
		headers[k] = &HTTPHeaderValue{
			Values: v,
		}
	}

	protoDetails := &HTTPDetailsEntity{
		Method:  method,
		Host:    d.Host,
		Uri:     d.URI,
		Proto:   d.Proto,
		Headers: headers,
	}

	any, err = anypb.New(protoDetails)
	return
}
