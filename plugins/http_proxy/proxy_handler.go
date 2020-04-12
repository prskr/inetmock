package main

import (
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type proxyHttpHandler struct {
	options httpProxyOptions
	logger  *zap.Logger
}

/*
TODO implement HTTPS proxy like in TLS interceptor
func (p *proxyHttpHandler) HandleConnect(req string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return &goproxy.ConnectAction{
		Action: goproxy.OkConnect,

	}, ""
}*/

func (p *proxyHttpHandler) Handle(req *http.Request, _ *goproxy.ProxyCtx) (retReq *http.Request, resp *http.Response) {

	retReq = req
	resp = &http.Response{
		Request:          req,
		TransferEncoding: req.TransferEncoding,
		Header:           make(http.Header),
		StatusCode:       http.StatusOK,
	}

	p.logger.Info(
		"Handling request",
		zap.String("source", req.RemoteAddr),
		zap.String("host", req.Host),
		zap.String("method", req.Method),
		zap.String("protocol", req.Proto),
		zap.String("path", req.RequestURI),
		zap.Reflect("headers", req.Header),
	)

	for _, rule := range p.options.Rules {
		if rule.pattern.MatchString(req.URL.Path) {
			if file, err := os.Open(rule.response); err != nil {
				p.logger.Error(
					"failed to open response target file",
					zap.String("resonse", rule.response),
					zap.Error(err),
				)
				continue
			} else {
				resp.Body = file

				if stat, err := file.Stat(); err == nil {
					resp.ContentLength = stat.Size()
				}

				if contentType, err := GetContentType(rule, file); err == nil {
					resp.Header["Content-Type"] = []string{contentType}
				}

				p.logger.Info("returning fake response from rules")
				return req, resp
			}
		}
	}

	if resp, err := p.options.FallbackStrategy.Apply(req); err != nil {
		p.logger.Error(
			"failed to apply fallback strategy",
			zap.Error(err),
		)
	} else {
		p.logger.Info("returning fake response from fallback strategy")
		return req, resp
	}

	p.logger.Info("falling back to proxying request through")
	return req, nil
}

func GetContentType(rule targetRule, file *os.File) (contentType string, err error) {
	if contentType = mime.TypeByExtension(filepath.Ext(rule.response)); contentType != "" {
		return
	}

	var buf [512]byte

	n, _ := io.ReadFull(file, buf[:])

	contentType = http.DetectContentType(buf[:n])
	_, err = file.Seek(0, io.SeekStart)
	return
}
