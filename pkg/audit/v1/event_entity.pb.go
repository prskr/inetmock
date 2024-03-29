// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: audit/v1/event_entity.proto

package auditv1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TransportProtocol int32

const (
	TransportProtocol_TRANSPORT_PROTOCOL_UNSPECIFIED TransportProtocol = 0
	TransportProtocol_TRANSPORT_PROTOCOL_TCP         TransportProtocol = 1
	TransportProtocol_TRANSPORT_PROTOCOL_UDP         TransportProtocol = 2
)

// Enum value maps for TransportProtocol.
var (
	TransportProtocol_name = map[int32]string{
		0: "TRANSPORT_PROTOCOL_UNSPECIFIED",
		1: "TRANSPORT_PROTOCOL_TCP",
		2: "TRANSPORT_PROTOCOL_UDP",
	}
	TransportProtocol_value = map[string]int32{
		"TRANSPORT_PROTOCOL_UNSPECIFIED": 0,
		"TRANSPORT_PROTOCOL_TCP":         1,
		"TRANSPORT_PROTOCOL_UDP":         2,
	}
)

func (x TransportProtocol) Enum() *TransportProtocol {
	p := new(TransportProtocol)
	*p = x
	return p
}

func (x TransportProtocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransportProtocol) Descriptor() protoreflect.EnumDescriptor {
	return file_audit_v1_event_entity_proto_enumTypes[0].Descriptor()
}

func (TransportProtocol) Type() protoreflect.EnumType {
	return &file_audit_v1_event_entity_proto_enumTypes[0]
}

func (x TransportProtocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransportProtocol.Descriptor instead.
func (TransportProtocol) EnumDescriptor() ([]byte, []int) {
	return file_audit_v1_event_entity_proto_rawDescGZIP(), []int{0}
}

type AppProtocol int32

const (
	AppProtocol_APP_PROTOCOL_UNSPECIFIED    AppProtocol = 0
	AppProtocol_APP_PROTOCOL_DNS            AppProtocol = 1
	AppProtocol_APP_PROTOCOL_HTTP           AppProtocol = 2
	AppProtocol_APP_PROTOCOL_HTTP_PROXY     AppProtocol = 3
	AppProtocol_APP_PROTOCOL_PPROF          AppProtocol = 4
	AppProtocol_APP_PROTOCOL_DNS_OVER_HTTPS AppProtocol = 5
	AppProtocol_APP_PROTOCOL_DHCP           AppProtocol = 6
)

// Enum value maps for AppProtocol.
var (
	AppProtocol_name = map[int32]string{
		0: "APP_PROTOCOL_UNSPECIFIED",
		1: "APP_PROTOCOL_DNS",
		2: "APP_PROTOCOL_HTTP",
		3: "APP_PROTOCOL_HTTP_PROXY",
		4: "APP_PROTOCOL_PPROF",
		5: "APP_PROTOCOL_DNS_OVER_HTTPS",
		6: "APP_PROTOCOL_DHCP",
	}
	AppProtocol_value = map[string]int32{
		"APP_PROTOCOL_UNSPECIFIED":    0,
		"APP_PROTOCOL_DNS":            1,
		"APP_PROTOCOL_HTTP":           2,
		"APP_PROTOCOL_HTTP_PROXY":     3,
		"APP_PROTOCOL_PPROF":          4,
		"APP_PROTOCOL_DNS_OVER_HTTPS": 5,
		"APP_PROTOCOL_DHCP":           6,
	}
)

func (x AppProtocol) Enum() *AppProtocol {
	p := new(AppProtocol)
	*p = x
	return p
}

func (x AppProtocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AppProtocol) Descriptor() protoreflect.EnumDescriptor {
	return file_audit_v1_event_entity_proto_enumTypes[1].Descriptor()
}

func (AppProtocol) Type() protoreflect.EnumType {
	return &file_audit_v1_event_entity_proto_enumTypes[1]
}

func (x AppProtocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AppProtocol.Descriptor instead.
func (AppProtocol) EnumDescriptor() ([]byte, []int) {
	return file_audit_v1_event_entity_proto_rawDescGZIP(), []int{1}
}

type TLSVersion int32

const (
	TLSVersion_TLS_VERSION_UNSPECIFIED TLSVersion = 0
	TLSVersion_TLS_VERSION_TLS10       TLSVersion = 1
	TLSVersion_TLS_VERSION_TLS11       TLSVersion = 2
	TLSVersion_TLS_VERSION_TLS12       TLSVersion = 3
	TLSVersion_TLS_VERSION_TLS13       TLSVersion = 4
)

// Enum value maps for TLSVersion.
var (
	TLSVersion_name = map[int32]string{
		0: "TLS_VERSION_UNSPECIFIED",
		1: "TLS_VERSION_TLS10",
		2: "TLS_VERSION_TLS11",
		3: "TLS_VERSION_TLS12",
		4: "TLS_VERSION_TLS13",
	}
	TLSVersion_value = map[string]int32{
		"TLS_VERSION_UNSPECIFIED": 0,
		"TLS_VERSION_TLS10":       1,
		"TLS_VERSION_TLS11":       2,
		"TLS_VERSION_TLS12":       3,
		"TLS_VERSION_TLS13":       4,
	}
)

func (x TLSVersion) Enum() *TLSVersion {
	p := new(TLSVersion)
	*p = x
	return p
}

func (x TLSVersion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TLSVersion) Descriptor() protoreflect.EnumDescriptor {
	return file_audit_v1_event_entity_proto_enumTypes[2].Descriptor()
}

func (TLSVersion) Type() protoreflect.EnumType {
	return &file_audit_v1_event_entity_proto_enumTypes[2]
}

func (x TLSVersion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TLSVersion.Descriptor instead.
func (TLSVersion) EnumDescriptor() ([]byte, []int) {
	return file_audit_v1_event_entity_proto_rawDescGZIP(), []int{2}
}

type TLSDetailsEntity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version     TLSVersion `protobuf:"varint,1,opt,name=version,proto3,enum=inetmock.audit.v1.TLSVersion" json:"version,omitempty"`
	CipherSuite string     `protobuf:"bytes,2,opt,name=cipher_suite,json=cipherSuite,proto3" json:"cipher_suite,omitempty"`
	ServerName  string     `protobuf:"bytes,3,opt,name=server_name,json=serverName,proto3" json:"server_name,omitempty"`
}

func (x *TLSDetailsEntity) Reset() {
	*x = TLSDetailsEntity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_v1_event_entity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TLSDetailsEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSDetailsEntity) ProtoMessage() {}

func (x *TLSDetailsEntity) ProtoReflect() protoreflect.Message {
	mi := &file_audit_v1_event_entity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSDetailsEntity.ProtoReflect.Descriptor instead.
func (*TLSDetailsEntity) Descriptor() ([]byte, []int) {
	return file_audit_v1_event_entity_proto_rawDescGZIP(), []int{0}
}

func (x *TLSDetailsEntity) GetVersion() TLSVersion {
	if x != nil {
		return x.Version
	}
	return TLSVersion_TLS_VERSION_UNSPECIFIED
}

func (x *TLSDetailsEntity) GetCipherSuite() string {
	if x != nil {
		return x.CipherSuite
	}
	return ""
}

func (x *TLSDetailsEntity) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

type EventEntity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Timestamp       *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Transport       TransportProtocol      `protobuf:"varint,3,opt,name=transport,proto3,enum=inetmock.audit.v1.TransportProtocol" json:"transport,omitempty"`
	Application     AppProtocol            `protobuf:"varint,4,opt,name=application,proto3,enum=inetmock.audit.v1.AppProtocol" json:"application,omitempty"`
	SourceIp        []byte                 `protobuf:"bytes,5,opt,name=source_ip,json=sourceIp,proto3" json:"source_ip,omitempty"`
	DestinationIp   []byte                 `protobuf:"bytes,6,opt,name=destination_ip,json=destinationIp,proto3" json:"destination_ip,omitempty"`
	SourcePort      uint32                 `protobuf:"varint,7,opt,name=source_port,json=sourcePort,proto3" json:"source_port,omitempty"`
	DestinationPort uint32                 `protobuf:"varint,8,opt,name=destination_port,json=destinationPort,proto3" json:"destination_port,omitempty"`
	Tls             *TLSDetailsEntity      `protobuf:"bytes,9,opt,name=tls,proto3" json:"tls,omitempty"`
	// Types that are assignable to ProtocolDetails:
	//	*EventEntity_Http
	//	*EventEntity_Dns
	//	*EventEntity_Dhcp
	//	*EventEntity_NetMon
	ProtocolDetails isEventEntity_ProtocolDetails `protobuf_oneof:"protocol_details"`
}

func (x *EventEntity) Reset() {
	*x = EventEntity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_v1_event_entity_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventEntity) ProtoMessage() {}

func (x *EventEntity) ProtoReflect() protoreflect.Message {
	mi := &file_audit_v1_event_entity_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventEntity.ProtoReflect.Descriptor instead.
func (*EventEntity) Descriptor() ([]byte, []int) {
	return file_audit_v1_event_entity_proto_rawDescGZIP(), []int{1}
}

func (x *EventEntity) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *EventEntity) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *EventEntity) GetTransport() TransportProtocol {
	if x != nil {
		return x.Transport
	}
	return TransportProtocol_TRANSPORT_PROTOCOL_UNSPECIFIED
}

func (x *EventEntity) GetApplication() AppProtocol {
	if x != nil {
		return x.Application
	}
	return AppProtocol_APP_PROTOCOL_UNSPECIFIED
}

func (x *EventEntity) GetSourceIp() []byte {
	if x != nil {
		return x.SourceIp
	}
	return nil
}

func (x *EventEntity) GetDestinationIp() []byte {
	if x != nil {
		return x.DestinationIp
	}
	return nil
}

func (x *EventEntity) GetSourcePort() uint32 {
	if x != nil {
		return x.SourcePort
	}
	return 0
}

func (x *EventEntity) GetDestinationPort() uint32 {
	if x != nil {
		return x.DestinationPort
	}
	return 0
}

func (x *EventEntity) GetTls() *TLSDetailsEntity {
	if x != nil {
		return x.Tls
	}
	return nil
}

func (m *EventEntity) GetProtocolDetails() isEventEntity_ProtocolDetails {
	if m != nil {
		return m.ProtocolDetails
	}
	return nil
}

func (x *EventEntity) GetHttp() *HTTPDetailsEntity {
	if x, ok := x.GetProtocolDetails().(*EventEntity_Http); ok {
		return x.Http
	}
	return nil
}

func (x *EventEntity) GetDns() *DNSDetailsEntity {
	if x, ok := x.GetProtocolDetails().(*EventEntity_Dns); ok {
		return x.Dns
	}
	return nil
}

func (x *EventEntity) GetDhcp() *DHCPDetailsEntity {
	if x, ok := x.GetProtocolDetails().(*EventEntity_Dhcp); ok {
		return x.Dhcp
	}
	return nil
}

func (x *EventEntity) GetNetMon() *NetMonDetailsEntity {
	if x, ok := x.GetProtocolDetails().(*EventEntity_NetMon); ok {
		return x.NetMon
	}
	return nil
}

type isEventEntity_ProtocolDetails interface {
	isEventEntity_ProtocolDetails()
}

type EventEntity_Http struct {
	Http *HTTPDetailsEntity `protobuf:"bytes,20,opt,name=http,proto3,oneof"`
}

type EventEntity_Dns struct {
	Dns *DNSDetailsEntity `protobuf:"bytes,21,opt,name=dns,proto3,oneof"`
}

type EventEntity_Dhcp struct {
	Dhcp *DHCPDetailsEntity `protobuf:"bytes,22,opt,name=dhcp,proto3,oneof"`
}

type EventEntity_NetMon struct {
	NetMon *NetMonDetailsEntity `protobuf:"bytes,23,opt,name=net_mon,json=netMon,proto3,oneof"`
}

func (*EventEntity_Http) isEventEntity_ProtocolDetails() {}

func (*EventEntity_Dns) isEventEntity_ProtocolDetails() {}

func (*EventEntity_Dhcp) isEventEntity_ProtocolDetails() {}

func (*EventEntity_NetMon) isEventEntity_ProtocolDetails() {}

var File_audit_v1_event_entity_proto protoreflect.FileDescriptor

var file_audit_v1_event_entity_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x69,
	0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70,
	0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a,
	0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6e, 0x73, 0x5f, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x61, 0x75, 0x64, 0x69,
	0x74, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x68, 0x63, 0x70, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76,
	0x31, 0x2f, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x6e, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f, 0x01, 0x0a, 0x10, 0x54, 0x4c, 0x53, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x37, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x69,
	0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x4c, 0x53, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x5f, 0x73,
	0x75, 0x69, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x69, 0x70, 0x68,
	0x65, 0x72, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xac, 0x05, 0x0a, 0x0b, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x42, 0x0a, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b,
	0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70,
	0x6f, 0x72, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x52, 0x09, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x40, 0x0a, 0x0b, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x69, 0x6e,
	0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x70, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x52, 0x0b, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x49, 0x70, 0x12, 0x25, 0x0a, 0x0e, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x64,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x70, 0x12, 0x1f, 0x0a, 0x0b,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x29, 0x0a,
	0x10, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x35, 0x0a, 0x03, 0x74, 0x6c, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b,
	0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x4c, 0x53, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x03, 0x74, 0x6c, 0x73, 0x12,
	0x3a, 0x0a, 0x04, 0x68, 0x74, 0x74, 0x70, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x48, 0x00, 0x52, 0x04, 0x68, 0x74, 0x74, 0x70, 0x12, 0x37, 0x0a, 0x03, 0x64,
	0x6e, 0x73, 0x18, 0x15, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d,
	0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x4e, 0x53,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x48, 0x00, 0x52,
	0x03, 0x64, 0x6e, 0x73, 0x12, 0x3a, 0x0a, 0x04, 0x64, 0x68, 0x63, 0x70, 0x18, 0x16, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x24, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75,
	0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x48, 0x43, 0x50, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x73, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x48, 0x00, 0x52, 0x04, 0x64, 0x68, 0x63, 0x70,
	0x12, 0x41, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x5f, 0x6d, 0x6f, 0x6e, 0x18, 0x17, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x26, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64,
	0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x74, 0x4d, 0x6f, 0x6e, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x48, 0x00, 0x52, 0x06, 0x6e, 0x65, 0x74,
	0x4d, 0x6f, 0x6e, 0x42, 0x12, 0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x5f,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2a, 0x6f, 0x0a, 0x11, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x22, 0x0a, 0x1e,
	0x54, 0x52, 0x41, 0x4e, 0x53, 0x50, 0x4f, 0x52, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4f, 0x4c, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x1a, 0x0a, 0x16, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x50, 0x4f, 0x52, 0x54, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c, 0x5f, 0x54, 0x43, 0x50, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16,
	0x54, 0x52, 0x41, 0x4e, 0x53, 0x50, 0x4f, 0x52, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4f, 0x4c, 0x5f, 0x55, 0x44, 0x50, 0x10, 0x02, 0x2a, 0xc5, 0x01, 0x0a, 0x0b, 0x41, 0x70, 0x70,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x1c, 0x0a, 0x18, 0x41, 0x50, 0x50, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x41, 0x50, 0x50, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c, 0x5f, 0x44, 0x4e, 0x53, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11,
	0x41, 0x50, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c, 0x5f, 0x48, 0x54, 0x54,
	0x50, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x41, 0x50, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4f, 0x4c, 0x5f, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x58, 0x59, 0x10, 0x03,
	0x12, 0x16, 0x0a, 0x12, 0x41, 0x50, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c,
	0x5f, 0x50, 0x50, 0x52, 0x4f, 0x46, 0x10, 0x04, 0x12, 0x1f, 0x0a, 0x1b, 0x41, 0x50, 0x50, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c, 0x5f, 0x44, 0x4e, 0x53, 0x5f, 0x4f, 0x56, 0x45,
	0x52, 0x5f, 0x48, 0x54, 0x54, 0x50, 0x53, 0x10, 0x05, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x50, 0x50,
	0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c, 0x5f, 0x44, 0x48, 0x43, 0x50, 0x10, 0x06,
	0x2a, 0x85, 0x01, 0x0a, 0x0a, 0x54, 0x4c, 0x53, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x1b, 0x0a, 0x17, 0x54, 0x4c, 0x53, 0x5f, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11,
	0x54, 0x4c, 0x53, 0x5f, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x4c, 0x53, 0x31,
	0x30, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x4c, 0x53, 0x5f, 0x56, 0x45, 0x52, 0x53, 0x49,
	0x4f, 0x4e, 0x5f, 0x54, 0x4c, 0x53, 0x31, 0x31, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x4c,
	0x53, 0x5f, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x4c, 0x53, 0x31, 0x32, 0x10,
	0x03, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x4c, 0x53, 0x5f, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e,
	0x5f, 0x54, 0x4c, 0x53, 0x31, 0x33, 0x10, 0x04, 0x42, 0xc4, 0x01, 0x0a, 0x15, 0x63, 0x6f, 0x6d,
	0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e,
	0x76, 0x31, 0x42, 0x10, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x48, 0x02, 0x50, 0x01, 0x5a, 0x31, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x63, 0x6b, 0x2e, 0x69, 0x63, 0x62, 0x34, 0x64, 0x63, 0x30, 0x2e, 0x64, 0x65, 0x2f, 0x69, 0x6e,
	0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74,
	0x2f, 0x76, 0x31, 0x3b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x41,
	0x58, 0xaa, 0x02, 0x11, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x41, 0x75, 0x64,
	0x69, 0x74, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x11, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b,
	0x5c, 0x41, 0x75, 0x64, 0x69, 0x74, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x49, 0x6e, 0x65, 0x74,
	0x6d, 0x6f, 0x63, 0x6b, 0x5c, 0x41, 0x75, 0x64, 0x69, 0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x49, 0x6e, 0x65, 0x74,
	0x6d, 0x6f, 0x63, 0x6b, 0x3a, 0x3a, 0x41, 0x75, 0x64, 0x69, 0x74, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_audit_v1_event_entity_proto_rawDescOnce sync.Once
	file_audit_v1_event_entity_proto_rawDescData = file_audit_v1_event_entity_proto_rawDesc
)

func file_audit_v1_event_entity_proto_rawDescGZIP() []byte {
	file_audit_v1_event_entity_proto_rawDescOnce.Do(func() {
		file_audit_v1_event_entity_proto_rawDescData = protoimpl.X.CompressGZIP(file_audit_v1_event_entity_proto_rawDescData)
	})
	return file_audit_v1_event_entity_proto_rawDescData
}

var file_audit_v1_event_entity_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_audit_v1_event_entity_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_audit_v1_event_entity_proto_goTypes = []interface{}{
	(TransportProtocol)(0),        // 0: inetmock.audit.v1.TransportProtocol
	(AppProtocol)(0),              // 1: inetmock.audit.v1.AppProtocol
	(TLSVersion)(0),               // 2: inetmock.audit.v1.TLSVersion
	(*TLSDetailsEntity)(nil),      // 3: inetmock.audit.v1.TLSDetailsEntity
	(*EventEntity)(nil),           // 4: inetmock.audit.v1.EventEntity
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
	(*HTTPDetailsEntity)(nil),     // 6: inetmock.audit.v1.HTTPDetailsEntity
	(*DNSDetailsEntity)(nil),      // 7: inetmock.audit.v1.DNSDetailsEntity
	(*DHCPDetailsEntity)(nil),     // 8: inetmock.audit.v1.DHCPDetailsEntity
	(*NetMonDetailsEntity)(nil),   // 9: inetmock.audit.v1.NetMonDetailsEntity
}
var file_audit_v1_event_entity_proto_depIdxs = []int32{
	2, // 0: inetmock.audit.v1.TLSDetailsEntity.version:type_name -> inetmock.audit.v1.TLSVersion
	5, // 1: inetmock.audit.v1.EventEntity.timestamp:type_name -> google.protobuf.Timestamp
	0, // 2: inetmock.audit.v1.EventEntity.transport:type_name -> inetmock.audit.v1.TransportProtocol
	1, // 3: inetmock.audit.v1.EventEntity.application:type_name -> inetmock.audit.v1.AppProtocol
	3, // 4: inetmock.audit.v1.EventEntity.tls:type_name -> inetmock.audit.v1.TLSDetailsEntity
	6, // 5: inetmock.audit.v1.EventEntity.http:type_name -> inetmock.audit.v1.HTTPDetailsEntity
	7, // 6: inetmock.audit.v1.EventEntity.dns:type_name -> inetmock.audit.v1.DNSDetailsEntity
	8, // 7: inetmock.audit.v1.EventEntity.dhcp:type_name -> inetmock.audit.v1.DHCPDetailsEntity
	9, // 8: inetmock.audit.v1.EventEntity.net_mon:type_name -> inetmock.audit.v1.NetMonDetailsEntity
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_audit_v1_event_entity_proto_init() }
func file_audit_v1_event_entity_proto_init() {
	if File_audit_v1_event_entity_proto != nil {
		return
	}
	file_audit_v1_http_details_proto_init()
	file_audit_v1_dns_details_proto_init()
	file_audit_v1_dhcp_details_proto_init()
	file_audit_v1_netmon_details_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_audit_v1_event_entity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TLSDetailsEntity); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_audit_v1_event_entity_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventEntity); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_audit_v1_event_entity_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*EventEntity_Http)(nil),
		(*EventEntity_Dns)(nil),
		(*EventEntity_Dhcp)(nil),
		(*EventEntity_NetMon)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_audit_v1_event_entity_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_audit_v1_event_entity_proto_goTypes,
		DependencyIndexes: file_audit_v1_event_entity_proto_depIdxs,
		EnumInfos:         file_audit_v1_event_entity_proto_enumTypes,
		MessageInfos:      file_audit_v1_event_entity_proto_msgTypes,
	}.Build()
	File_audit_v1_event_entity_proto = out.File
	file_audit_v1_event_entity_proto_rawDesc = nil
	file_audit_v1_event_entity_proto_goTypes = nil
	file_audit_v1_event_entity_proto_depIdxs = nil
}
