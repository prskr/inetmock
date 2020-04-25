package cert

import (
	"fmt"
	"strings"
)

func extractIPFromAddress(addr string) (ip string, err error) {
	if idx := strings.LastIndex(addr, ":"); idx < 0 {
		err = fmt.Errorf("addr %s does not match expected scheme <ip>:<port>", addr)

	} else {
		/* get IP part of address */
		ip = addr[0:idx]

		/* trim [ ] for IPv6 addresses */
		if ip[0] == '[' {
			ip = ip[1 : len(ip)-1]
		}
	}
	return
}
