package details

import (
	"google.golang.org/protobuf/types/known/anypb"

	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

func NewDNSFromWireFormat(entity *auditv1.DNSDetailsEntity) DNS {
	d := DNS{
		OPCode: entity.Opcode,
	}

	for _, q := range entity.Questions {
		d.Questions = append(d.Questions, DNSQuestion{
			RRType: q.Type,
			Name:   q.Name,
		})
	}

	return d
}

type DNSQuestion struct {
	RRType auditv1.ResourceRecordType
	Name   string
}

type DNS struct {
	OPCode    auditv1.DNSOpCode
	Questions []DNSQuestion
}

func (d DNS) MarshalToWireFormat() (any *anypb.Any, err error) {
	detailsEntity := &auditv1.DNSDetailsEntity{
		Opcode: d.OPCode,
	}

	for _, q := range d.Questions {
		detailsEntity.Questions = append(detailsEntity.Questions, &auditv1.DNSQuestionEntity{
			Type: q.RRType,
			Name: q.Name,
		})
	}

	any, err = anypb.New(detailsEntity)

	return
}
