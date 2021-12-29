package api

import (
	"html/template"
	"io/fs"
	"net/http"
	"sort"

	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var indexPaths = []string{
	"/",
	"/index.html",
	"/index.htm",
	"/index",
}

func init() {
	sort.Strings(indexPaths)
}

type IndexHandler struct {
	templates *template.Template
	logger    logging.Logger
}

func RegisterViews(mux *http.ServeMux, viewsFS fs.FS, logger logging.Logger) (err error) {
	h := IndexHandler{
		logger:    logger,
		templates: template.New("static"),
	}

	if h.templates, err = h.templates.ParseFS(viewsFS, "*.gohtml"); err != nil {
		return err
	}

	for _, p := range indexPaths {
		mux.Handle(p, h)
	}

	return nil
}

func (h IndexHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if !isIndexPath(req.URL.Path) {
		return
	}
	if err := h.templates.ExecuteTemplate(resp, "index.gohtml", make(map[string]interface{})); err != nil {
		h.logger.Error("Failed to execute template", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

func isIndexPath(path string) bool {
	idx := sort.SearchStrings(indexPaths, path)
	return idx < len(indexPaths) && indexPaths[idx] == path
}
