package http

import (
	"encoding/json"
	"net/http"

	"gitlab.com/inetmock/inetmock/pkg/health"
)

func NewHealthHandler(checker health.Checker) http.Handler {
	return &healthHandler{checker: checker}
}

type healthHandler struct {
	checker health.Checker
}

func (h healthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	var result health.Result
	if result, err = h.checker.Status(request.Context()); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !result.IsHealthy() {
		var data []byte
		if data, err = json.Marshal(result); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(data)
		writer.WriteHeader(http.StatusServiceUnavailable)
		writer.Header().Set("Content-Type", "application/json")
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}
