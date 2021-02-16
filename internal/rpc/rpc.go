package rpc

import "net/url"

type Config struct {
	Listen string
}

func (r Config) ListenURL() (u *url.URL) {
	var err error
	if u, err = url.Parse(r.Listen); err != nil {
		u, _ = url.Parse("tcp://:0")
	}
	return
}
