package web

import "embed"

//go:embed *.gohtml css/*.css js/*.js wasm/*.wasm
var WebFS embed.FS
