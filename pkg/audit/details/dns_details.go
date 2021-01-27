package details

import (
	"google.golang.org/protobuf/types/known/anypb"
)

func NewDNSFromWireFormat(entity *DNSDetailsEntity) DNS {
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
	RRType ResourceRecordType
	Name   string
}

type DNS struct {
	OPCode    DNSOpCode
	Questions []DNSQuestion
}

func (d DNS) MarshalToWireFormat() (any *anypb.Any, err error) {
	detailsEntity := &DNSDetailsEntity{
		Opcode: d.OPCode,
	}

	for _, q := range d.Questions {
		detailsEntity.Questions = append(detailsEntity.Questions, &DNSQuestionEntity{
			Type: q.RRType,
			Name: q.Name,
		})
	}

	any, err = anypb.New(detailsEntity)

	return
}
