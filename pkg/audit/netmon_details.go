package audit

import (
	"reflect"

	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

var _ Details = (*NetMon)(nil)

func init() {
	AddMapping(reflect.TypeOf(new(auditv1.EventEntity_NetMon)), func(msg *auditv1.EventEntity) Details {
		var entity *auditv1.NetMonDetailsEntity
		if e, ok := msg.ProtocolDetails.(*auditv1.EventEntity_NetMon); !ok {
			return nil
		} else {
			entity = e.NetMon
		}

		return &NetMon{
			ReverseResolvedDomain: entity.ReverseResolvedHost,
		}
	})
}

type NetMon struct {
	ReverseResolvedDomain string
}

func (n *NetMon) AddToMsg(msg *auditv1.EventEntity) {
	msg.ProtocolDetails = &auditv1.EventEntity_NetMon{
		NetMon: &auditv1.NetMonDetailsEntity{
			ReverseResolvedHost: n.ReverseResolvedDomain,
		},
	}
}
