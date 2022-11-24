package proxy

import (
	"crypto/tls"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/jinzhu/copier"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	"inetmock.icb4dc0.de/inetmock/protocols"
)

type proxyHTTPSHandler struct {
	options   httpProxyOptions
	tlsConfig *tls.Config
}

func (p *proxyHTTPSHandler) HandleConnect(string, *goproxy.ProxyCtx) (resultingAction *goproxy.ConnectAction, redirectTo string) {
	return &goproxy.ConnectAction{
		Action: goproxy.ConnectAccept,
		TLSConfig: func(host string, ctx *goproxy.ProxyCtx) (*tls.Config, error) {
			return p.tlsConfig, nil
		},
	}, p.options.Target.host()
}

type proxyHTTPHandler struct {
	handlerName string
	options     httpProxyOptions
	logger      logging.Logger
}

func (p *proxyHTTPHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (retReq *http.Request, resp *http.Response) {
	defer prometheus.NewTimer(protocols.RequestDurationHistogram.WithLabelValues("http_proxy", p.handlerName)).ObserveDuration()

	retReq = req

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
