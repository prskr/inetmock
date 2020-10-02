package tls_interceptor

import (
	"fmt"
)

type redirectionTarget struct {
	IPAddress string
	Port      uint16
}

func (rt redirectionTarget) address() string {
	return fmt.Sprintf("%s:%d", rt.IPAddress, rt.Port)
}

type tlsOptions struct {
	Target redirectionTarget
}
