package audit

import (
	"reflect"

	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

var _ Details = (*DNS)(nil)

func init() {
	AddMapping(reflect.TypeOf(new(auditv1.EventEntity_Dns)), func(msg *auditv1.EventEntity) Details {
		var dnsDetails *auditv1.DNSDetailsEntity
		if d, ok := msg.ProtocolDetails.(*auditv1.EventEntity_Dns); !ok {
			return nil
		} else {
			dnsDetails = d.Dns
		}

		dns := &DNS{
			OPCode:    dnsDetails.Opcode,
			Questions: make([]DNSQuestion, 0, len(dnsDetails.Questions)),
		}

		for idx := range dnsDetails.Questions {
			q := dnsDetails.Questions[idx]
			dns.Questions = append(dns.Questions, DNSQuestion{
				RRType: q.Type,
				Name:   q.Name,
			})
		}

		return dns
	})
}

type DNSQuestion struct {
	RRType auditv1.ResourceRecordType
	Name   string
}

type DNS struct {
	OPCode    auditv1.DNSOpCode
	Questions []DNSQuestion
}

func (d DNS) AddToMsg(msg *auditv1.EventEntity) {
	details := &auditv1.DNSDetailsEntity{
		Opcode:    d.OPCode,
		Questions: make([]*auditv1.DNSQuestionEntity, 0, len(d.Questions)),
	}

	for idx := range d.Questions {
		q := d.Questions[idx]
		details.Questions = append(details.Questions, &auditv1.DNSQuestionEntity{
			Type: q.RRType,
			Name: q.Name,
		})
	}

	msg.ProtocolDetails = &auditv1.EventEntity_Dns{
		Dns: details,
	}
}
