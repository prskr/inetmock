package api

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type StaticFilesMiddleware struct {
	compressor http.Handler
	logger     logging.Logger
	http.FileSystem
	MimeTypeOverrides map[string]string
	next              http.Handler
}

func RegisterStaticFileHandlingMiddleware(handler http.Handler, webFS fs.FS, logger logging.Logger) http.Handler {
	fsmw := &StaticFilesMiddleware{
		FileSystem: http.FS(webFS),
		MimeTypeOverrides: map[string]string{
			".wasm": "application/wasm",
		},
		logger: logger,
		next:   handler,
	}

	return fsmw
}

func (m StaticFilesMiddleware) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet, http.MethodHead:
		break
	default:
		m.next.ServeHTTP(resp, req)
		return
	}

	requestPath := req.URL.Path

	if isIndexPath(requestPath) {
		m.next.ServeHTTP(resp, req)
		return
	}

	var (
		f   http.File
		err error
	)
	if f, err = m.FileSystem.Open(strings.TrimLeft(requestPath, "/")); err != nil {
		m.next.ServeHTTP(resp, req)
		return
	}

	defer multierr.AppendInvoke(&err, multierr.Close(f))

	var info fs.FileInfo
	if info, err = f.Stat(); err != nil {
		m.logger.Error("Failed to stat file", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		m.next.ServeHTTP(resp, req)
		return
	}

	extension := path.Ext(info.Name())

	if override, ok := m.MimeTypeOverrides[extension]; ok {
		resp.Header().Set("Content-Type", override)
	}

	http.ServeContent(resp, req, path.Base(info.Name()), info.ModTime(), f)
}
