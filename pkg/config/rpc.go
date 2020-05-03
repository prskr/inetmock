package config

import "net/url"

type RPC struct {
	Listen string
}

func (r RPC) ListenURL() (u *url.URL) {
	var err error
	if u, err = url.Parse(r.Listen); err != nil {
		u, _ = url.Parse("tcp://:0")
	}
	return
}
