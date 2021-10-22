package protocolmock

import (
	"net/http"
	"testing"
	"time"
)

type HTTPClientMockDoCallParams struct {
	Req http.Request
}

type HTTPClientMockDoCallResult struct {
	Resp *http.Response
	Err  error
}

type HTTPClientMockDoCall struct {
	Timestamp time.Time
	Params    HTTPClientMockDoCallParams
	Result    HTTPClientMockDoCallResult
}

type HTTPClientMockCalls struct {
	Do []HTTPClientMockDoCall
}

type HTTPClientMockContext struct {
	HTTPClientMockCalls
	TB testing.TB
}

type HTTPClientMock struct {
	TB    testing.TB
	Calls HTTPClientMockCalls
	OnDo  func(state HTTPClientMockContext, req *http.Request) (*http.Response, error)
}

func (m *HTTPClientMock) Do(req *http.Request) (resp *http.Response, err error) {
	call := HTTPClientMockDoCall{
		Timestamp: time.Now(),
		Params: HTTPClientMockDoCallParams{
			Req: *req,
		},
	}

	defer func() {
		m.Calls.Do = append(m.Calls.Do, call)
	}()

	if m.OnDo != nil {
		ctx := HTTPClientMockContext{
			TB:                  m.TB,
			HTTPClientMockCalls: m.Calls,
		}
		resp, err = m.OnDo(ctx, req)
		call.Result = HTTPClientMockDoCallResult{
			Resp: resp,
			Err:  err,
		}
	}

	return
}
