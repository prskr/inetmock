package dns

import (
	"net"
	"strconv"
	"strings"
)

const (
	inAddrArpaSuffix = ".in-addr.arpa."
	suffixLength     = len(inAddrArpaSuffix)
	baseDecimal      = 10
	byteLength       = 8
)

func ParseInAddrArpa(inAddrArpa string) net.IP {
	if !strings.HasSuffix(inAddrArpa, inAddrArpaSuffix) {
		return nil
	}

	ip := inAddrArpa[:len(inAddrArpa)-suffixLength]
	ipBytes := strings.Split(ip, ".")

	if len(ipBytes) != net.IPv4len {
		return nil
	}

	bs := make([]byte, 0, net.IPv4len)
	for i := len(ipBytes) - 1; i >= 0; i-- {
		if parsed, err := strconv.ParseUint(ipBytes[i], baseDecimal, byteLength); err != nil {
			return nil
		} else {
			bs = append(bs, byte(parsed))
		}
	}

	return net.IPv4(bs[0], bs[1], bs[2], bs[3])
}
