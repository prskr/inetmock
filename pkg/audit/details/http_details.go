package details

import (
	"net/http"

	"google.golang.org/protobuf/types/known/anypb"

	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

var (
	httpGoToWireMapping = map[string]auditv1.HTTPMethod{
		http.MethodGet:     auditv1.HTTPMethod_HTTP_METHOD_GET,
		http.MethodHead:    auditv1.HTTPMethod_HTTP_METHOD_HEAD,
		http.MethodPost:    auditv1.HTTPMethod_HTTP_METHOD_POST,
		http.MethodPut:     auditv1.HTTPMethod_HTTP_METHOD_PUT,
		http.MethodPatch:   auditv1.HTTPMethod_HTTP_METHOD_PATCH,
		http.MethodDelete:  auditv1.HTTPMethod_HTTP_METHOD_DELETE,
		http.MethodConnect: auditv1.HTTPMethod_HTTP_METHOD_CONNECT,
		http.MethodOptions: auditv1.HTTPMethod_HTTP_METHOD_OPTIONS,
		http.MethodTrace:   auditv1.HTTPMethod_HTTP_METHOD_TRACE,
	}
	httpWireToGoMapping = map[auditv1.HTTPMethod]string{
		auditv1.HTTPMethod_HTTP_METHOD_GET:     http.MethodGet,
		auditv1.HTTPMethod_HTTP_METHOD_HEAD:    http.MethodHead,
		auditv1.HTTPMethod_HTTP_METHOD_POST:    http.MethodPost,
		auditv1.HTTPMethod_HTTP_METHOD_PUT:     http.MethodPut,
		auditv1.HTTPMethod_HTTP_METHOD_PATCH:   http.MethodPatch,
		auditv1.HTTPMethod_HTTP_METHOD_DELETE:  http.MethodDelete,
		auditv1.HTTPMethod_HTTP_METHOD_CONNECT: http.MethodConnect,
		auditv1.HTTPMethod_HTTP_METHOD_OPTIONS: http.MethodOptions,
		auditv1.HTTPMethod_HTTP_METHOD_TRACE:   http.MethodTrace,
	}
)

type HTTP struct {
	Method  string
	Host    string
	URI     string
	Proto   string
	Headers http.Header
}

func NewHTTPFromWireFormat(entity *auditv1.HTTPDetailsEntity) HTTP {
	headers := http.Header{}
	for name, values := range entity.Headers {
		for idx := range values.Values {
			headers.Add(name, values.Values[idx])
		}
	}

	method := ""
	if mappedMethod, known := httpWireToGoMapping[entity.Method]; known {
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
	method := auditv1.HTTPMethod_HTTP_METHOD_UNSPECIFIED
	if methodValue, known := httpGoToWireMapping[d.Method]; known {
		method = methodValue
	}

	headers := make(map[string]*auditv1.HTTPHeaderValue)

	for k, v := range d.Headers {
		headers[k] = &auditv1.HTTPHeaderValue{
			Values: v,
		}
	}

	protoDetails := &auditv1.HTTPDetailsEntity{
		Method:  method,
		Host:    d.Host,
		Uri:     d.URI,
		Proto:   d.Proto,
		Headers: headers,
	}

	any, err = anypb.New(protoDetails)
	return
}
