package tls_interceptor

import (
	"fmt"
	"net"
)

type proxyConn struct {
	source net.Conn
	target net.Conn
}

func (p *proxyConn) Close() error {
	var err error
	if targetErr := p.target.Close(); targetErr != nil {
		err = fmt.Errorf("error while closing Target conn: %w", targetErr)
	}
	if sourceErr := p.source.Close(); sourceErr != nil {
		err = fmt.Errorf("error while closing source conn: %w", err)
	}
	return err
}
