package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

const (
	numberArgsWithoutBody = 1
	numberArgsWithBody    = 2
)

var (
	ErrNotAnHTTPInitiator = errors.New("the given initiator is not an HTTP initiator")

	knownInitiators = map[string]func(logger logging.Logger, args ...rules.Param) (Initiator, error){
		"get":     RequestInitiatorForMethod(http.MethodGet),
		"head":    RequestInitiatorForMethod(http.MethodHead),
		"post":    RequestInitiatorForMethod(http.MethodPost),
		"put":     RequestInitiatorForMethod(http.MethodPut),
		"delete":  RequestInitiatorForMethod(http.MethodDelete),
		"options": RequestInitiatorForMethod(http.MethodOptions),
	}
)

type Initiator interface {
	Do(ctx context.Context, client *http.Client) (resp *http.Response, err error)
}

func InitiatorForRule(rule *rules.Check, logger logging.Logger) (Initiator, error) {
	if rule.Initiator == nil {
		return nil, rules.ErrNoInitiatorDefined
	}

	switch m := strings.ToLower(rule.Initiator.Module); m {
	case "http", "http2":
		if constructor, ok := knownInitiators[strings.ToLower(rule.Initiator.Name)]; !ok {
			return nil, fmt.Errorf("%w %s", rules.ErrUnknownInitiator, rule.Initiator.Name)
		} else {
			return constructor(logger, rule.Initiator.Params...)
		}
	default:
		return nil, fmt.Errorf("%w: %s", ErrNotAnHTTPInitiator, m)
	}
}

type simpleRequest struct {
	logger     logging.Logger
	method     string
	uri        string
	bodyBuffer *bytes.Buffer
	body       []byte
}

func (s *simpleRequest) Do(ctx context.Context, client *http.Client) (resp *http.Response, err error) {
	s.logger.Info("Execute HTTP health check")
	var req *http.Request
	s.bodyBuffer.Reset()
	if _, err = s.bodyBuffer.Write(s.body); err != nil {
		return nil, err
	}
	if req, err = http.NewRequestWithContext(ctx, s.method, s.uri, s.bodyBuffer); err != nil {
		return
	}
	req.Header.Set("Accept-Encoding", "identity")
	return client.Do(req)
}

func RequestInitiatorForMethod(method string) func(logger logging.Logger, params ...rules.Param) (Initiator, error) {
	return func(logger logging.Logger, params ...rules.Param) (Initiator, error) {
		if err := rules.ValidateParameterCount(params, numberArgsWithoutBody); err != nil {
			return nil, err
		}

		var body []byte = nil
		switch method {
		case http.MethodPost, http.MethodPut:
			if len(params) == numberArgsWithBody {
				var err error
				var jsonString string
				if jsonString, err = params[1].AsString(); err == nil {
					body = []byte(jsonString)
				}
			}
		default:
		}

		var err error
		var uri string
		if uri, err = params[0].AsString(); err != nil {
			return nil, err
		}

		logger = logger.With(
			zap.String("method", method),
			zap.String("uri", uri),
		)

		logger.Debug("Setup health initiator")

		return &simpleRequest{
			logger:     logger,
			method:     method,
			uri:        uri,
			body:       body,
			bodyBuffer: bytes.NewBuffer(body),
		}, nil
	}
}
