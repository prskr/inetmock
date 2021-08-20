package mock_test

import (
	"testing/fstest"
	"time"
)

var (
	defaultHTMLContent = `<html>
<head>
    <title>INetSim default HTML page</title>
</head>
<body>
<p></p>
<p align="center">This is the default HTML page for INetMock HTTP mock protocols.</p>
<p align="center">This file is an HTML document.</p>
</body>
</html>`
	defaultFakeFileFS = fstest.MapFS{
		"default.html": &fstest.MapFile{
			Data:    []byte(defaultHTMLContent),
			ModTime: time.Now().Add(-1337 * time.Second),
		},
	}
)
