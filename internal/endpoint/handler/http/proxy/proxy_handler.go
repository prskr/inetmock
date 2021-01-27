package proxy

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
	imHttp "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
)

type proxyHttpHandler struct {
	handlerName string
	options     httpProxyOptions
	logger      logging.Logger
	emitter     audit.Emitter
}

type proxyHttpsHandler struct {
	handlerName string
	tlsConfig   *tls.Config
	logger      logging.Logger
	emitter     audit.Emitter
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
	p.emitter.Emit(imHttp.EventFromRequest(req, audit.AppProtocol_HTTP_PROXY))

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
	redirectReq = redirectReq.WithContext(originalRequest.Context())

	return
}
