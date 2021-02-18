package proxy

import (
	"crypto/tls"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"

	imHttp "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type proxyHttpsHandler struct {
	options   httpProxyOptions
	tlsConfig *tls.Config
	emitter   audit.Emitter
}

func (p *proxyHttpsHandler) HandleConnect(_ string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	p.emitter.Emit(imHttp.EventFromRequest(ctx.Req, audit.AppProtocol_HTTP_PROXY))

	return &goproxy.ConnectAction{
		Action: goproxy.ConnectAccept,
		TLSConfig: func(host string, ctx *goproxy.ProxyCtx) (*tls.Config, error) {
			return p.tlsConfig, nil
		},
	}, p.options.Target.host()
}

type proxyHttpHandler struct {
	handlerName string
	options     httpProxyOptions
	logger      logging.Logger
	emitter     audit.Emitter
}

func (p *proxyHttpHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (retReq *http.Request, resp *http.Response) {
	timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(p.handlerName))
	defer timer.ObserveDuration()

	retReq = req
	p.emitter.Emit(imHttp.EventFromRequest(req, audit.AppProtocol_HTTP_PROXY))

	var err error
	var redirectReq *http.Request
	if redirectReq, err = redirectHTTPRequest(p.options.Target.host(), req); err != nil {
		return req, nil
	}
	if resp, err = ctx.RoundTrip(redirectReq); err != nil {
		p.logger.Error(
			"error while doing roundtrip",
			zap.Error(err),
		)
		return req, nil
	}

	return
}

func redirectHTTPRequest(targetHost string, originalRequest *http.Request) (redirectReq *http.Request, err error) {
	redirectReq = new(http.Request)
	if err = copier.Copy(redirectReq, originalRequest); err != nil {
		return
	}
	originalRequest.URL.Host = targetHost
	return
}
