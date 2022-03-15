package main

import (
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"gitlab.com/inetmock/inetmock/internal/app"
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/health"
)

var (
	serverApp app.App
	cfg       appConfig
)

type Data struct {
	PCAP      string
	Audit     string
	FakeFiles string
	State     string
}

func (d *Data) setup() (err error) {
	if d.PCAP, err = ensureDataDir(d.PCAP); err != nil {
		return
	}
	if d.Audit, err = ensureDataDir(d.Audit); err != nil {
		return
	}
	var stateDir string
	if stateDir, err = ensureDataDir(filepath.Dir(d.State)); err != nil {
		return
	} else {
		d.State = filepath.Join(stateDir, filepath.Base(d.State))
	}

	if !filepath.IsAbs(d.FakeFiles) {
		d.FakeFiles, err = filepath.Abs(d.FakeFiles)
	}

	return
}

func ensureDataDir(dataDirPath string) (cleanedPath string, err error) {
	const (
		defaultDirPerm = 0o750
	)

	cleanedPath = dataDirPath
	if !filepath.IsAbs(cleanedPath) {
		if cleanedPath, err = filepath.Abs(cleanedPath); err != nil {
			return
		}
	}

	err = os.MkdirAll(cleanedPath, defaultDirPerm)
	return
}

type appConfig struct {
	TLS       cert.Options
	Listeners map[string]endpoint.ListenerSpec
	API       struct {
		Listen string
	}
	Health health.Config
	Data   Data
}

func (c *appConfig) APIURL() *url.URL {
	if u, err := url.Parse(c.API.Listen); err != nil {
		u, _ = url.Parse("tcp://:0")
		return u
	} else {
		return u
	}
}

func main() {
	serverApp = app.NewApp(
		app.Spec{
			Name:        "inetmock",
			LogEncoding: "json",
			Short:       "INetMock is lightweight internet mock",
			Config:      &cfg,
			SubCommands: []*cobra.Command{serveCmd, generateCaCmd},
			Defaults: map[string]any{
				"api.listen":                            "tcp://:0",
				"data.pcap":                             "/var/lib/inetmock/data/pcap",
				"data.audit":                            "/var/lib/inetmock/data/audit",
				"data.state":                            "/var/lib/inetmock/data/state/inetmock.db",
				"tls.curve":                             cert.CurveTypeP256,
				"tls.minTLSVersion":                     cert.TLSVersionTLS10,
				"tls.includeInsecureCipherSuites":       false,
				"tls.validity.server.notBeforeRelative": 168 * time.Hour,
				"tls.validity.server.notAfterRelative":  168 * time.Hour,
				"tls.certCachePath":                     "/tmp",
			},
		},
	)

	serverApp.MustRun()
}
