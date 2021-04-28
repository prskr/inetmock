package mock

import (
	"errors"
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

var (
	ErrNoTerminatorDefined = errors.New("no terminator defined")
	ErrUnknownTerminator   = errors.New("no terminator with the given name is known")

	knownResponseHandlers = map[string]func(logger logging.Logger, fakeFileFS fs.FS, args ...rules.Param) (http.Handler, error){
		"file":   FileHandler,
		"status": StatusHandler,
	}
)

func HandlerForRoutingRule(rule *rules.Routing, logger logging.Logger, fakeFileFS fs.FS) (http.Handler, error) {
	if rule.Terminator == nil {
		return nil, ErrNoTerminatorDefined
	}

	if constructor, ok := knownResponseHandlers[strings.ToLower(rule.Terminator.Name)]; !ok {
		return nil, fmt.Errorf("%w %s", ErrUnknownTerminator, rule.Terminator.Name)
	} else {
		return constructor(logger, fakeFileFS, rule.Terminator.Params...)
	}
}

func FileHandler(logger logging.Logger, fakeFileFS fs.FS, args ...rules.Param) (http.Handler, error) {
	if err := validateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var err error
	var filePath string
	if filePath, err = args[0].AsString(); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("handlerType", "FileHandler"),
		zap.String("filePath", filePath),
	)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		file, err := fakeFileFS.Open(filePath)
		if err != nil {
			logger.Error("failed to open file to return", zap.Error(err))
			http.Error(writer, err.Error(), 500)
			return
		}

		var rs io.ReadSeeker
		var ok bool
		if rs, ok = file.(io.ReadSeeker); !ok {
			logger.Warn("file returned from FS does not support seeking - returning error")
			http.Error(writer, "internal server error", 500)
		}

		defer func() {
			_ = file.Close()
		}()

		logger.Info("Returning file response")
		//nolint:gosec
		http.ServeContent(writer, request, path.Base(request.RequestURI), time.Now().Add(-(time.Duration(rand.Int()) * time.Millisecond)), rs)
	}), nil
}

func StatusHandler(logger logging.Logger, _ fs.FS, args ...rules.Param) (http.Handler, error) {
	if err := validateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var err error
	var statusCodeToReturn int
	if statusCodeToReturn, err = args[0].AsInt(); err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("handlerType", "StatusHandler"),
		zap.Int("code", statusCodeToReturn),
	)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logger.Info("Returning status code")
		writer.WriteHeader(statusCodeToReturn)
	}), nil
}
