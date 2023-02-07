package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	mdns "github.com/miekg/dns"
)

var (
	RequestPackerPOST RequestPacker = RequestPackerFunc(func(queryURL *url.URL, question *mdns.Msg) (req *http.Request, err error) {
		if data, err := question.Pack(); err != nil {
			return nil, err
		} else {
			if req, err = http.NewRequest(http.MethodPost, queryURL.String(), bytes.NewReader(data)); err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/dns-message")
			return req, nil
		}
	})

	RequestPackerGET RequestPacker = RequestPackerFunc(func(queryURL *url.URL, question *mdns.Msg) (req *http.Request, err error) {
		var (
			data       []byte
			encodedMsg string
		)
		if data, err = question.Pack(); err != nil {
			return nil, err
		}
		encodedMsg = base64.URLEncoding.EncodeToString(data)
		queryValues := queryURL.Query()
		queryValues.Set("dns", encodedMsg)
		queryURL.RawQuery = queryValues.Encode()

		return http.NewRequest(http.MethodGet, queryURL.String(), nil)
	})
)

type RequestPacker interface {
	Pack(queryURL *url.URL, question *mdns.Msg) (req *http.Request, err error)
}

type HTTPClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type RequestPackerFunc func(queryURL *url.URL, question *mdns.Msg) (req *http.Request, err error)

func (f RequestPackerFunc) Pack(queryURL *url.URL, question *mdns.Msg) (req *http.Request, err error) {
	return f(queryURL, question)
}

type HTTPTransport struct {
	Packer RequestPacker
	Client HTTPClient
	Scheme string
	Server string
}

func (h HTTPTransport) RoundTrip(ctx context.Context, question *mdns.Msg) (resp *mdns.Msg, err error) {
	var (
		data      []byte
		parsedURL *url.URL
		req       *http.Request
		httpResp  *http.Response
		client    HTTPClient
		packer    RequestPacker
	)

	if parsedURL, err = url.Parse(fmt.Sprintf("%s://%s/dns-query", h.Scheme, h.Server)); err != nil {
		return nil, err
	}

	if packer = h.Packer; packer == nil {
		packer = RequestPackerPOST
	}

	if req, err = packer.Pack(parsedURL, question); err != nil {
		return nil, err
	}

	if client = h.Client; client == nil {
		client = http.DefaultClient
	}

	req = req.WithContext(ctx)
	if httpResp, err = client.Do(req); err != nil {
		return nil, err
	}

	if httpResp.Body != nil {
		defer func() {
			err = errors.Join(err, httpResp.Body.Close())
		}()
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d does not indicate success", httpResp.StatusCode)
	}

	if data, err = io.ReadAll(httpResp.Body); err != nil {
		return nil, err
	}
	resp = new(mdns.Msg)
	err = resp.Unpack(data)
	return resp, err
}
