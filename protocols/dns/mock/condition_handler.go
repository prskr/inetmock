package mock

import (
	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/protocols/dns"
)

type ConditionHandler struct {
	dns.IPResolver
	Filters []RequestFilter
}

func (h ConditionHandler) Matches(req *mdns.Question) bool {
	for idx := range h.Filters {
		if !h.Filters[idx].Matches(req) {
			return false
		}
	}
	return true
}
