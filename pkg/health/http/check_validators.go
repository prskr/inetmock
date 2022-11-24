package http

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"go.uber.org/multierr"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
)

var (
	ErrUnknownCheckFilter = errors.New("no check filter with the given name is known")
	ErrResponseEmpty      = errors.New("response must not be nil")

	knownCheckFilters = map[string]func(args ...rules.Param) (Validator, error){
		"status":   StatusCodeFilter,
		"header":   ResponseHeaderFilter,
		"contains": ResponseBodyContainsFilter,
		"sha256":   ResponseBodyHashSHA256Filter,
	}
	searchBufferPool = &sync.Pool{
		New: func() any {
			data := make([]byte, defaultSearchBufferSize)
			return &data
		},
	}
)

const (
	windowSizeSearchBufferMultiplier = 2
	expectedResponseHeaderParamCount = 2
	defaultSearchBufferSize          = 256 * windowSizeSearchBufferMultiplier
)

func ValidatorsForRule(rule *rules.Check) (chain ValidationChain, err error) {
	if rule.Validators == nil || len(rule.Validators.Chain) == 0 {
		return nil, nil
	}

	filterRules := rule.Validators.Chain
	chain = make(ValidationChain, 0, len(filterRules))

	for idx := range filterRules {
		rawRule := filterRules[idx]
		if constructor, ok := knownCheckFilters[strings.ToLower(rawRule.Name)]; !ok {
			return nil, fmt.Errorf("%w %s", ErrUnknownCheckFilter, rawRule.Name)
		} else {
			var instance Validator
			instance, err = constructor(rawRule.Params...)
			if err != nil {
				return
			}
			chain.Add(instance)
		}
	}
	return
}

type Validator interface {
	Matches(resp *http.Response) error
}

type ValidationChain []Validator

func (c *ValidationChain) Add(v Validator) {
	arr := *c
	arr = append(arr, v)
	*c = arr
}

func (c ValidationChain) Len() int {
	return len([]Validator(c))
}

func (c ValidationChain) Matches(resp *http.Response) error {
	for idx := range c {
		if err := c[idx].Matches(resp); err != nil {
			return err
		}
	}
	return nil
}

type CheckFilterFunc func(resp *http.Response) error

func (c CheckFilterFunc) Matches(resp *http.Response) error {
	return c(resp)
}

func StatusCodeFilter(args ...rules.Param) (Validator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var err error
	var expectedStatusCode int
	if expectedStatusCode, err = args[0].AsInt(); err != nil {
		return nil, err
	}

	return CheckFilterFunc(func(resp *http.Response) error {
		if resp == nil {
			return ErrResponseEmpty
		}
		if resp.StatusCode == expectedStatusCode {
			return nil
		}
		return fmt.Errorf("expected status code %d but got %d", expectedStatusCode, resp.StatusCode)
	}), nil
}

func ResponseHeaderFilter(args ...rules.Param) (Validator, error) {
	if err := rules.ValidateParameterCount(args, expectedResponseHeaderParamCount); err != nil {
		return nil, err
	}

	var err error
	var headerName, expectedValue string
	if headerName, err = args[0].AsString(); err != nil {
		return nil, err
	}

	if expectedValue, err = args[1].AsString(); err != nil {
		return nil, err
	}

	expectedValue = strings.ToLower(expectedValue)

	return CheckFilterFunc(func(resp *http.Response) error {
		if resp == nil {
			return ErrResponseEmpty
		}

		values := resp.Header.Values(headerName)
		if len(values) == 0 {
			return fmt.Errorf("no %s header prsent", headerName)
		}
		for idx := range values {
			if strings.Contains(strings.ToLower(values[idx]), expectedValue) {
				return nil
			}
		}
		return fmt.Errorf("could not match %s: %s", headerName, expectedValue)
	}), nil
}

func ResponseBodyContainsFilter(args ...rules.Param) (Validator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var searchValue []byte
	var searchValueLength int
	if searchString, err := args[0].AsString(); err != nil {
		return nil, err
	} else {
		searchValue = []byte(searchString)
		searchValueLength = len(searchValue)
	}

	return CheckFilterFunc(func(resp *http.Response) (err error) {
		if resp == nil {
			return ErrResponseEmpty
		}
		var searchBuffer []byte
		if searchValueLength < defaultSearchBufferSize {
			searchBuffer = *searchBufferPool.Get().(*[]byte)
			defer func() {
				searchBufferPool.Put(&searchBuffer)
			}()
		} else {
			searchBuffer = make([]byte, searchValueLength*windowSizeSearchBufferMultiplier)
		}

		defer multierr.AppendInvoke(&err, multierr.Close(resp.Body))

		var read, idx int
		more := true
		for more {
			read, err = resp.Body.Read(searchBuffer[searchValueLength:])
			switch {
			case errors.Is(err, nil):
				break
			case errors.Is(err, io.EOF):
				more = false
			default:
				return err
			}

			if idx = bytes.Index(searchBuffer[:searchValueLength+read], searchValue); idx >= 0 {
				return nil
			}
			copy(searchBuffer, searchBuffer[:searchValueLength])
		}
		return errors.New("expected value not found in body")
	}), nil
}

func ResponseBodyHashSHA256Filter(args ...rules.Param) (Validator, error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	var expectedHash []byte
	if hashString, err := args[0].AsString(); err != nil {
		return nil, err
	} else if expectedHash, err = hex.DecodeString(hashString); err != nil {
		return nil, err
	}

	return CheckFilterFunc(func(resp *http.Response) (err error) {
		if resp == nil {
			return ErrResponseEmpty
		}
		defer multierr.AppendInvoke(&err, multierr.Close(resp.Body))

		hash := sha256.New()
		if _, err = io.Copy(hash, resp.Body); err != nil {
			return
		}

		if actual := hash.Sum(nil); !bytes.Equal(actual, expectedHash) {
			return fmt.Errorf("hash values do not match - expected %x got %x", expectedHash, actual)
		}
		return nil
	}), nil
}
