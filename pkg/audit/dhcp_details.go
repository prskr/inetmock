package audit

import (
	"reflect"

	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

var _ Details = (*DHCP)(nil)

func init() {
	AddMapping(reflect.TypeOf(new(auditv1.EventEntity_Dhcp)), func(msg *auditv1.EventEntity) Details {
		var entity *auditv1.DHCPDetailsEntity
		if e, ok := msg.ProtocolDetails.(*auditv1.EventEntity_Dhcp); !ok {
			return nil
		} else {
			entity = e.Dhcp
		}

		return &DHCP{
			HopCount: uint8(entity.HopCount),
			OpCode:   entity.Opcode,
			HWType:   entity.HwType,
		}
	})
}

type DHCP struct {
	HopCount uint8
	OpCode   auditv1.DHCPOpCode
	HWType   auditv1.DHCPHwType
}

func (d DHCP) AddToMsg(msg *auditv1.EventEntity) {
	msg.ProtocolDetails = &auditv1.EventEntity_Dhcp{
		Dhcp: &auditv1.DHCPDetailsEntity{
			HopCount: int32(d.OpCode),
			Opcode:   d.OpCode,
			HwType:   d.HWType,
		},
	}
}
