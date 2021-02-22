package details

import (
	"net/http"

	"google.golang.org/protobuf/types/known/anypb"

	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

var (
	goToWireMapping = map[string]v1.HTTPMethod{
		http.MethodGet:     v1.HTTPMethod_HTTP_METHOD_GET,
		http.MethodHead:    v1.HTTPMethod_HTTP_METHOD_HEAD,
		http.MethodPost:    v1.HTTPMethod_HTTP_METHOD_POST,
		http.MethodPut:     v1.HTTPMethod_HTTP_METHOD_PUT,
		http.MethodPatch:   v1.HTTPMethod_HTTP_METHOD_PATCH,
		http.MethodDelete:  v1.HTTPMethod_HTTP_METHOD_DELETE,
		http.MethodConnect: v1.HTTPMethod_HTTP_METHOD_CONNECT,
		http.MethodOptions: v1.HTTPMethod_HTTP_METHOD_OPTIONS,
		http.MethodTrace:   v1.HTTPMethod_HTTP_METHOD_TRACE,
	}
	wireToGoMapping = map[v1.HTTPMethod]string{
		v1.HTTPMethod_HTTP_METHOD_GET:     http.MethodGet,
		v1.HTTPMethod_HTTP_METHOD_HEAD:    http.MethodHead,
		v1.HTTPMethod_HTTP_METHOD_POST:    http.MethodPost,
		v1.HTTPMethod_HTTP_METHOD_PUT:     http.MethodPut,
		v1.HTTPMethod_HTTP_METHOD_PATCH:   http.MethodPatch,
		v1.HTTPMethod_HTTP_METHOD_DELETE:  http.MethodDelete,
		v1.HTTPMethod_HTTP_METHOD_CONNECT: http.MethodConnect,
		v1.HTTPMethod_HTTP_METHOD_OPTIONS: http.MethodOptions,
		v1.HTTPMethod_HTTP_METHOD_TRACE:   http.MethodTrace,
	}
)

type HTTP struct {
	Method  string
	Host    string
	URI     string
	Proto   string
	Headers http.Header
}

func NewHTTPFromWireFormat(entity *v1.HTTPDetailsEntity) HTTP {
	headers := http.Header{}
	for name, values := range entity.Headers {
		for idx := range values.Values {
			headers.Add(name, values.Values[idx])
		}
	}

	var method = ""
	if mappedMethod, known := wireToGoMapping[entity.Method]; known {
		method = mappedMethod
	}

	return HTTP{
		Method:  method,
		Host:    entity.Host,
		URI:     entity.Uri,
		Proto:   entity.Proto,
		Headers: headers,
	}
}

func (d HTTP) MarshalToWireFormat() (any *anypb.Any, err error) {
	var method = v1.HTTPMethod_HTTP_METHOD_UNSPECIFIED
	if methodValue, known := goToWireMapping[d.Method]; known {
		method = methodValue
	}

	headers := make(map[string]*v1.HTTPHeaderValue)

	for k, v := range d.Headers {
		headers[k] = &v1.HTTPHeaderValue{
			Values: v,
		}
	}

	protoDetails := &v1.HTTPDetailsEntity{
		Method:  method,
		Host:    d.Host,
		Uri:     d.URI,
		Proto:   d.Proto,
		Headers: headers,
	}

	any, err = anypb.New(protoDetails)
	return
}
