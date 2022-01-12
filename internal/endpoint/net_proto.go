package endpoint

type NetProto uint

const (
	NetProtoUDP NetProto = iota + 1
	NetProtoTCP
)

func (p NetProto) String() string {
	switch p {
	case NetProtoTCP:
		return "TCP"
	case NetProtoUDP:
		return "UDP"
	default:
		return ""
	}
}
