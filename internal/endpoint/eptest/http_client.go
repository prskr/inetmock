package eptest

import (
	"net/http"
	"time"
)

func HTTPClientForInMemListener(lis InMemListener) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext:           lis.DialContext,
			DialTLSContext:        lis.DialContext,
			MaxIdleConns:          5,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
