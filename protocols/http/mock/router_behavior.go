package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var knownResponseHandlers = map[string]func(logger logging.Logger, fakeFileFS fs.FS, args ...rules.Param) (http.Handler, error){
	"file":   FileHandler,
	"status": StatusHandler,
	"json":   JSONHandler,
}

func HandlerForRoutingRule(rule *rules.SingleResponsePipeline, logger logging.Logger, fakeFileFS fs.FS) (http.Handler, error) {
	if rule.Response == nil {
		return nil, rules.ErrNoTerminatorDefined
	}

	if constructor, ok := knownResponseHandlers[strings.ToLower(rule.Response.Name)]; !ok {
		return nil, fmt.Errorf("%w %s", rules.ErrUnknownTerminator, rule.Response.Name)
	} else {
		return constructor(logger, fakeFileFS, rule.Response.Params...)
	}
}

func FileHandler(logger logging.Logger, fakeFileFS fs.FS, args ...rules.Param) (http.Handler, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var (
		filePath string
		err      error
	)
	if filePath, err = args[0].AsString(); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("handler_type", "FileHandler"),
		zap.String("file_path", filePath),
	)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		file, err := fakeFileFS.Open(filePath)
		if err != nil {
			logger.Error("failed to open file to return", zap.Error(err))
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		var rs io.ReadSeeker
		var ok bool
		if rs, ok = file.(io.ReadSeeker); !ok {
			logger.Warn("file returned from FS does not support seeking - returning error")
			http.Error(writer, "internal server error", http.StatusInternalServerError)
		}

		defer func() {
			_ = file.Close()
		}()

		logger.Debug("Returning file response")
		//nolint:gosec
		http.ServeContent(writer, request, path.Base(request.RequestURI), time.Now().Add(-(time.Duration(rand.Int()) * time.Millisecond)), rs)
	}), nil
}

func StatusHandler(logger logging.Logger, _ fs.FS, args ...rules.Param) (http.Handler, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var err error
	var statusCodeToReturn int
	if statusCodeToReturn, err = args[0].AsInt(); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("handler_type", "StatusHandler"),
		zap.Int("code", statusCodeToReturn),
	)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logger.Debug("Returning status code")
		writer.WriteHeader(statusCodeToReturn)
	}), nil
}

func JSONHandler(logger logging.Logger, _ fs.FS, args ...rules.Param) (http.Handler, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var jsonString string
	if s, err := args[0].AsString(); err != nil {
		return nil, err
	} else {
		jsonString = s
	}

	jsonBytes := []byte(jsonString)
	into := make(map[string]any)
	if err := json.Unmarshal(jsonBytes, &into); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("handler_type", "JSONHandler"),
	)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logger.Debug("Returning status code")
		if _, err := writer.Write(jsonBytes); err != nil {
			logger.Warn("Failed to write JSON response", zap.Error(err))
		}
	}), nil
}
