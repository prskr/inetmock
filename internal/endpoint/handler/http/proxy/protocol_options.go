package proxy

import (
	"fmt"
)

type redirectionTarget struct {
	IPAddress string
	Port      uint16
}

func (rt redirectionTarget) host() string {
	return fmt.Sprintf("%s:%d", rt.IPAddress, rt.Port)
}

type httpProxyOptions struct {
	Target redirectionTarget
}
