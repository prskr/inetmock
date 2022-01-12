package details

import (
	"google.golang.org/protobuf/types/known/anypb"

	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

func NewDHCPFromWireFormat(entity *auditv1.DHCPDetailsEntity) DHCP {
	d := DHCP{
		HopCount: uint8(entity.HopCount),
		OpCode:   entity.Opcode,
		HWType:   entity.HwType,
	}

	return d
}

type DHCP struct {
	HopCount uint8
	OpCode   auditv1.DHCPOpCode
	HWType   auditv1.DHCPHwType
}

func (d DHCP) MarshalToWireFormat() (any *anypb.Any, err error) {
	detailsEntity := &auditv1.DHCPDetailsEntity{
		HopCount: int32(d.OpCode),
		Opcode:   d.OpCode,
		HwType:   d.HWType,
	}
	any, err = anypb.New(detailsEntity)

	return
}
