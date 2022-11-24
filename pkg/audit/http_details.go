package audit

import (
	"net/http"
	"reflect"

	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
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
	_ Details = (*HTTP)(nil)
)

func init() {
	AddMapping(reflect.TypeOf(new(auditv1.EventEntity_Http)), func(msg *auditv1.EventEntity) Details {
		var entity *auditv1.HTTPDetailsEntity

		if e, ok := msg.ProtocolDetails.(*auditv1.EventEntity_Http); !ok {
			return nil
		} else {
			entity = e.Http
		}

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

		return &HTTP{
			Method:  method,
			Host:    entity.Host,
			URI:     entity.Uri,
			Proto:   entity.Proto,
			Headers: headers,
		}
	})
}

type HTTP struct {
	Method  string
	Host    string
	URI     string
	Proto   string
	Headers http.Header
}

func (d *HTTP) AddToMsg(msg *auditv1.EventEntity) {
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

	msg.ProtocolDetails = &auditv1.EventEntity_Http{Http: protoDetails}
}
