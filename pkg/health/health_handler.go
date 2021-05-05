package health

import (
	"encoding/json"
	"net/http"
)

func NewHealthHandler(checker Checker) http.Handler {
	return &healthHandler{checker: checker}
}

type healthHandler struct {
	checker Checker
}

func (h healthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var result = h.checker.Status(request.Context())

	if !result.IsHealthy() {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusServiceUnavailable)
		var err error
		if err = json.NewEncoder(writer).Encode(result); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
	writer.WriteHeader(http.StatusNoContent)
}
