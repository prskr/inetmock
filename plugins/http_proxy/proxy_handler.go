package http_proxy

import (
	"context"
	"crypto/tls"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
	"net/url"
)

type proxyHttpHandler struct {
	handlerName string
	options     httpProxyOptions
	logger      logging.Logger
}

type proxyHttpsHandler struct {
	handlerName string
	tlsConfig   *tls.Config
	logger      logging.Logger
}

func (p *proxyHttpsHandler) HandleConnect(req string, _ *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	totalHttpsRequestCounter.WithLabelValues(p.handlerName).Inc()
	p.logger.Info(
		"Intercepting HTTPS proxy request",
		zap.String("request", req),
	)

	return &goproxy.ConnectAction{
		Action: goproxy.ConnectMitm,
		TLSConfig: func(host string, ctx *goproxy.ProxyCtx) (*tls.Config, error) {
			return p.tlsConfig, nil
		},
	}, ""
}

func (p *proxyHttpHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (retReq *http.Request, resp *http.Response) {
	timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(p.handlerName))
	defer timer.ObserveDuration()
	totalRequestCounter.WithLabelValues(p.handlerName).Inc()

	retReq = req
	p.logger.Info(
		"Handling request",
		zap.String("source", req.RemoteAddr),
		zap.String("host", req.Host),
		zap.String("method", req.Method),
		zap.String("protocol", req.Proto),
		zap.String("path", req.RequestURI),
		zap.Reflect("headers", req.Header),
	)

	var err error
	if resp, err = ctx.RoundTrip(p.redirectHTTPRequest(req)); err != nil {
		p.logger.Error(
			"error while doing roundtrip",
			zap.Error(err),
		)
		return req, nil
	}

	return
}

func (p proxyHttpHandler) redirectHTTPRequest(originalRequest *http.Request) (redirectReq *http.Request) {
	redirectReq = &http.Request{
		Method: originalRequest.Method,
		URL: &url.URL{
			Host:       p.options.Target.host(),
			Path:       originalRequest.URL.Path,
			ForceQuery: originalRequest.URL.ForceQuery,
			Fragment:   originalRequest.URL.Fragment,
			Opaque:     originalRequest.URL.Opaque,
			RawPath:    originalRequest.URL.RawPath,
			RawQuery:   originalRequest.URL.RawQuery,
			User:       originalRequest.URL.User,
		},
		Proto:            originalRequest.Proto,
		ProtoMajor:       originalRequest.ProtoMajor,
		ProtoMinor:       originalRequest.ProtoMinor,
		Header:           originalRequest.Header,
		Body:             originalRequest.Body,
		GetBody:          originalRequest.GetBody,
		ContentLength:    originalRequest.ContentLength,
		TransferEncoding: originalRequest.TransferEncoding,
		Close:            false,
		Host:             originalRequest.Host,
		Form:             originalRequest.Form,
		PostForm:         originalRequest.PostForm,
		MultipartForm:    originalRequest.MultipartForm,
		Trailer:          originalRequest.Trailer,
	}
	redirectReq = redirectReq.WithContext(context.Background())

	return
}
