// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: rpc/v1/net_flow.proto

package rpcv1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PacketForwardAction int32

const (
	PacketForwardAction_PACKET_FORWARD_ACTION_UNSPECIFIED PacketForwardAction = 0
	PacketForwardAction_PACKET_FORWARD_ACTION_DROP        PacketForwardAction = 1
	PacketForwardAction_PACKET_FORWARD_ACTION_PASS        PacketForwardAction = 2
)

// Enum value maps for PacketForwardAction.
var (
	PacketForwardAction_name = map[int32]string{
		0: "PACKET_FORWARD_ACTION_UNSPECIFIED",
		1: "PACKET_FORWARD_ACTION_DROP",
		2: "PACKET_FORWARD_ACTION_PASS",
	}
	PacketForwardAction_value = map[string]int32{
		"PACKET_FORWARD_ACTION_UNSPECIFIED": 0,
		"PACKET_FORWARD_ACTION_DROP":        1,
		"PACKET_FORWARD_ACTION_PASS":        2,
	}
)

func (x PacketForwardAction) Enum() *PacketForwardAction {
	p := new(PacketForwardAction)
	*p = x
	return p
}

func (x PacketForwardAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PacketForwardAction) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_v1_net_flow_proto_enumTypes[0].Descriptor()
}

func (PacketForwardAction) Type() protoreflect.EnumType {
	return &file_rpc_v1_net_flow_proto_enumTypes[0]
}

func (x PacketForwardAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketForwardAction.Descriptor instead.
func (PacketForwardAction) EnumDescriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{0}
}

type ListControlledInterfacesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListControlledInterfacesRequest) Reset() {
	*x = ListControlledInterfacesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListControlledInterfacesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListControlledInterfacesRequest) ProtoMessage() {}

func (x *ListControlledInterfacesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListControlledInterfacesRequest.ProtoReflect.Descriptor instead.
func (*ListControlledInterfacesRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{0}
}

type ListControlledInterfacesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ControlledInterfaces []string `protobuf:"bytes,1,rep,name=controlled_interfaces,json=controlledInterfaces,proto3" json:"controlled_interfaces,omitempty"`
}

func (x *ListControlledInterfacesResponse) Reset() {
	*x = ListControlledInterfacesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListControlledInterfacesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListControlledInterfacesResponse) ProtoMessage() {}

func (x *ListControlledInterfacesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListControlledInterfacesResponse.ProtoReflect.Descriptor instead.
func (*ListControlledInterfacesResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{1}
}

func (x *ListControlledInterfacesResponse) GetControlledInterfaces() []string {
	if x != nil {
		return x.ControlledInterfaces
	}
	return nil
}

type ListAvailableNetworkInterfacesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListAvailableNetworkInterfacesRequest) Reset() {
	*x = ListAvailableNetworkInterfacesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAvailableNetworkInterfacesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAvailableNetworkInterfacesRequest) ProtoMessage() {}

func (x *ListAvailableNetworkInterfacesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAvailableNetworkInterfacesRequest.ProtoReflect.Descriptor instead.
func (*ListAvailableNetworkInterfacesRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{2}
}

type ListAvailableNetworkInterfacesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AvailableDevices []*ListAvailableNetworkInterfacesResponse_NetworkInterface `protobuf:"bytes,1,rep,name=available_devices,json=availableDevices,proto3" json:"available_devices,omitempty"`
}

func (x *ListAvailableNetworkInterfacesResponse) Reset() {
	*x = ListAvailableNetworkInterfacesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAvailableNetworkInterfacesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAvailableNetworkInterfacesResponse) ProtoMessage() {}

func (x *ListAvailableNetworkInterfacesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAvailableNetworkInterfacesResponse.ProtoReflect.Descriptor instead.
func (*ListAvailableNetworkInterfacesResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{3}
}

func (x *ListAvailableNetworkInterfacesResponse) GetAvailableDevices() []*ListAvailableNetworkInterfacesResponse_NetworkInterface {
	if x != nil {
		return x.AvailableDevices
	}
	return nil
}

type StartPacketFlowControlRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InterfaceName          string   `protobuf:"bytes,1,opt,name=interface_name,json=interfaceName,proto3" json:"interface_name,omitempty"`                              // Name of the network interface to monitor
	WhitelistPorts         []uint32 `protobuf:"varint,2,rep,packed,name=whitelist_ports,json=whitelistPorts,proto3" json:"whitelist_ports,omitempty"`                   // Whitelist of destination ports that are ignored - maximum 32
	SourceIpWhitelist      [][]byte `protobuf:"bytes,3,rep,name=source_ip_whitelist,json=sourceIpWhitelist,proto3" json:"source_ip_whitelist,omitempty"`                // Whitelist of source IPs that are ignore - maximum 20
	DestinationIpWhitelist [][]byte `protobuf:"bytes,4,rep,name=destination_ip_whitelist,json=destinationIpWhitelist,proto3" json:"destination_ip_whitelist,omitempty"` // Whitelist of destination IPs that are ignored - maximum 20
	PortsToIntercept       []uint32 `protobuf:"varint,5,rep,packed,name=ports_to_intercept,json=portsToIntercept,proto3" json:"ports_to_intercept,omitempty"`
	RemoveRlimitMemlock    bool     `protobuf:"varint,7,opt,name=remove_rlimit_memlock,json=removeRlimitMemlock,proto3" json:"remove_rlimit_memlock,omitempty"` // possibly required for kernels < 5.11
}

func (x *StartPacketFlowControlRequest) Reset() {
	*x = StartPacketFlowControlRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartPacketFlowControlRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartPacketFlowControlRequest) ProtoMessage() {}

func (x *StartPacketFlowControlRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartPacketFlowControlRequest.ProtoReflect.Descriptor instead.
func (*StartPacketFlowControlRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{4}
}

func (x *StartPacketFlowControlRequest) GetInterfaceName() string {
	if x != nil {
		return x.InterfaceName
	}
	return ""
}

func (x *StartPacketFlowControlRequest) GetWhitelistPorts() []uint32 {
	if x != nil {
		return x.WhitelistPorts
	}
	return nil
}

func (x *StartPacketFlowControlRequest) GetSourceIpWhitelist() [][]byte {
	if x != nil {
		return x.SourceIpWhitelist
	}
	return nil
}

func (x *StartPacketFlowControlRequest) GetDestinationIpWhitelist() [][]byte {
	if x != nil {
		return x.DestinationIpWhitelist
	}
	return nil
}

func (x *StartPacketFlowControlRequest) GetPortsToIntercept() []uint32 {
	if x != nil {
		return x.PortsToIntercept
	}
	return nil
}

func (x *StartPacketFlowControlRequest) GetRemoveRlimitMemlock() bool {
	if x != nil {
		return x.RemoveRlimitMemlock
	}
	return false
}

type StartPacketFlowControlResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StartPacketFlowControlResponse) Reset() {
	*x = StartPacketFlowControlResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartPacketFlowControlResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartPacketFlowControlResponse) ProtoMessage() {}

func (x *StartPacketFlowControlResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartPacketFlowControlResponse.ProtoReflect.Descriptor instead.
func (*StartPacketFlowControlResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{5}
}

type StopPacketFlowControlRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InterfaceName string `protobuf:"bytes,1,opt,name=interface_name,json=interfaceName,proto3" json:"interface_name,omitempty"`
}

func (x *StopPacketFlowControlRequest) Reset() {
	*x = StopPacketFlowControlRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopPacketFlowControlRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopPacketFlowControlRequest) ProtoMessage() {}

func (x *StopPacketFlowControlRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopPacketFlowControlRequest.ProtoReflect.Descriptor instead.
func (*StopPacketFlowControlRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{6}
}

func (x *StopPacketFlowControlRequest) GetInterfaceName() string {
	if x != nil {
		return x.InterfaceName
	}
	return ""
}

type StopPacketFlowControlResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InterfaceWasControlled bool `protobuf:"varint,1,opt,name=interface_was_controlled,json=interfaceWasControlled,proto3" json:"interface_was_controlled,omitempty"`
}

func (x *StopPacketFlowControlResponse) Reset() {
	*x = StopPacketFlowControlResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopPacketFlowControlResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopPacketFlowControlResponse) ProtoMessage() {}

func (x *StopPacketFlowControlResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopPacketFlowControlResponse.ProtoReflect.Descriptor instead.
func (*StopPacketFlowControlResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{7}
}

func (x *StopPacketFlowControlResponse) GetInterfaceWasControlled() bool {
	if x != nil {
		return x.InterfaceWasControlled
	}
	return false
}

type ListAvailableNetworkInterfacesResponse_NetworkInterface struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Addresses [][]byte `protobuf:"bytes,2,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *ListAvailableNetworkInterfacesResponse_NetworkInterface) Reset() {
	*x = ListAvailableNetworkInterfacesResponse_NetworkInterface{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_net_flow_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAvailableNetworkInterfacesResponse_NetworkInterface) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAvailableNetworkInterfacesResponse_NetworkInterface) ProtoMessage() {}

func (x *ListAvailableNetworkInterfacesResponse_NetworkInterface) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_net_flow_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAvailableNetworkInterfacesResponse_NetworkInterface.ProtoReflect.Descriptor instead.
func (*ListAvailableNetworkInterfacesResponse_NetworkInterface) Descriptor() ([]byte, []int) {
	return file_rpc_v1_net_flow_proto_rawDescGZIP(), []int{3, 0}
}

func (x *ListAvailableNetworkInterfacesResponse_NetworkInterface) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListAvailableNetworkInterfacesResponse_NetworkInterface) GetAddresses() [][]byte {
	if x != nil {
		return x.Addresses
	}
	return nil
}

var File_rpc_v1_net_flow_proto protoreflect.FileDescriptor

var file_rpc_v1_net_flow_proto_rawDesc = []byte{
	0x0a, 0x15, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x65, 0x74, 0x5f, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63,
	0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x22, 0x21, 0x0a, 0x1f, 0x4c, 0x69, 0x73, 0x74,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x57, 0x0a, 0x20, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x64, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x33, 0x0a, 0x15, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x64, 0x5f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x14,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x73, 0x22, 0x27, 0x0a, 0x25, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xe5, 0x01,
	0x0a, 0x26, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x75, 0x0a, 0x11, 0x61, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x48, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61,
	0x62, 0x6c, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x10, 0x61,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a,
	0x44, 0x0a, 0x10, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0xbb, 0x02, 0x0a, 0x1d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x66, 0x61, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x27,
	0x0a, 0x0f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x0e, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69,
	0x73, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x70, 0x5f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0c, 0x52, 0x11, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x70, 0x57, 0x68,
	0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x18, 0x64, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x70, 0x5f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c,
	0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x16, 0x64, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x70, 0x57, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73,
	0x74, 0x12, 0x2c, 0x0a, 0x12, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x5f, 0x74, 0x6f, 0x5f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x63, 0x65, 0x70, 0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x10, 0x70,
	0x6f, 0x72, 0x74, 0x73, 0x54, 0x6f, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x63, 0x65, 0x70, 0x74, 0x12,
	0x32, 0x0a, 0x15, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x5f, 0x72, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x5f, 0x6d, 0x65, 0x6d, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13,
	0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x4d, 0x65, 0x6d, 0x6c,
	0x6f, 0x63, 0x6b, 0x22, 0x20, 0x0a, 0x1e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x45, 0x0a, 0x1c, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61,
	0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x59, 0x0a, 0x1d,
	0x53, 0x74, 0x6f, 0x70, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a,
	0x18, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x77, 0x61, 0x73, 0x5f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x16, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x57, 0x61, 0x73, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x64, 0x2a, 0x7c, 0x0a, 0x13, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x25,
	0x0a, 0x21, 0x50, 0x41, 0x43, 0x4b, 0x45, 0x54, 0x5f, 0x46, 0x4f, 0x52, 0x57, 0x41, 0x52, 0x44,
	0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x50, 0x41, 0x43, 0x4b, 0x45, 0x54, 0x5f,
	0x46, 0x4f, 0x52, 0x57, 0x41, 0x52, 0x44, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44,
	0x52, 0x4f, 0x50, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x50, 0x41, 0x43, 0x4b, 0x45, 0x54, 0x5f,
	0x46, 0x4f, 0x52, 0x57, 0x41, 0x52, 0x44, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x50,
	0x41, 0x53, 0x53, 0x10, 0x02, 0x32, 0x9f, 0x04, 0x0a, 0x15, 0x4e, 0x65, 0x74, 0x46, 0x6c, 0x6f,
	0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x91, 0x01, 0x0a, 0x1e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63,
	0x65, 0x73, 0x12, 0x36, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62,
	0x6c, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61,
	0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x37, 0x2e, 0x69, 0x6e, 0x65,
	0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x7f, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x6c, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x12,
	0x30, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x64,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x31, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c,
	0x65, 0x64, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x79, 0x0a, 0x16, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x2e,
	0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f,
	0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x76, 0x0a, 0x15, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f,
	0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x2d, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d,
	0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xb2, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e,
	0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42,
	0x0c, 0x4e, 0x65, 0x74, 0x46, 0x6c, 0x6f, 0x77, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x02, 0x50,
	0x01, 0x5a, 0x2d, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x69, 0x63, 0x62, 0x34,
	0x64, 0x63, 0x30, 0x2e, 0x64, 0x65, 0x2f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x70, 0x63, 0x76, 0x31,
	0xa2, 0x02, 0x03, 0x49, 0x52, 0x58, 0xaa, 0x02, 0x0f, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63,
	0x6b, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0f, 0x49, 0x6e, 0x65, 0x74, 0x6d,
	0x6f, 0x63, 0x6b, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1b, 0x49, 0x6e, 0x65,
	0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x11, 0x49, 0x6e, 0x65, 0x74, 0x6d,
	0x6f, 0x63, 0x6b, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_v1_net_flow_proto_rawDescOnce sync.Once
	file_rpc_v1_net_flow_proto_rawDescData = file_rpc_v1_net_flow_proto_rawDesc
)

func file_rpc_v1_net_flow_proto_rawDescGZIP() []byte {
	file_rpc_v1_net_flow_proto_rawDescOnce.Do(func() {
		file_rpc_v1_net_flow_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_v1_net_flow_proto_rawDescData)
	})
	return file_rpc_v1_net_flow_proto_rawDescData
}

var file_rpc_v1_net_flow_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rpc_v1_net_flow_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_rpc_v1_net_flow_proto_goTypes = []interface{}{
	(PacketForwardAction)(0),                                        // 0: inetmock.rpc.v1.PacketForwardAction
	(*ListControlledInterfacesRequest)(nil),                         // 1: inetmock.rpc.v1.ListControlledInterfacesRequest
	(*ListControlledInterfacesResponse)(nil),                        // 2: inetmock.rpc.v1.ListControlledInterfacesResponse
	(*ListAvailableNetworkInterfacesRequest)(nil),                   // 3: inetmock.rpc.v1.ListAvailableNetworkInterfacesRequest
	(*ListAvailableNetworkInterfacesResponse)(nil),                  // 4: inetmock.rpc.v1.ListAvailableNetworkInterfacesResponse
	(*StartPacketFlowControlRequest)(nil),                           // 5: inetmock.rpc.v1.StartPacketFlowControlRequest
	(*StartPacketFlowControlResponse)(nil),                          // 6: inetmock.rpc.v1.StartPacketFlowControlResponse
	(*StopPacketFlowControlRequest)(nil),                            // 7: inetmock.rpc.v1.StopPacketFlowControlRequest
	(*StopPacketFlowControlResponse)(nil),                           // 8: inetmock.rpc.v1.StopPacketFlowControlResponse
	(*ListAvailableNetworkInterfacesResponse_NetworkInterface)(nil), // 9: inetmock.rpc.v1.ListAvailableNetworkInterfacesResponse.NetworkInterface
}
var file_rpc_v1_net_flow_proto_depIdxs = []int32{
	9, // 0: inetmock.rpc.v1.ListAvailableNetworkInterfacesResponse.available_devices:type_name -> inetmock.rpc.v1.ListAvailableNetworkInterfacesResponse.NetworkInterface
	3, // 1: inetmock.rpc.v1.NetFlowControlService.ListAvailableNetworkInterfaces:input_type -> inetmock.rpc.v1.ListAvailableNetworkInterfacesRequest
	1, // 2: inetmock.rpc.v1.NetFlowControlService.ListControlledInterfaces:input_type -> inetmock.rpc.v1.ListControlledInterfacesRequest
	5, // 3: inetmock.rpc.v1.NetFlowControlService.StartPacketFlowControl:input_type -> inetmock.rpc.v1.StartPacketFlowControlRequest
	7, // 4: inetmock.rpc.v1.NetFlowControlService.StopPacketFlowControl:input_type -> inetmock.rpc.v1.StopPacketFlowControlRequest
	4, // 5: inetmock.rpc.v1.NetFlowControlService.ListAvailableNetworkInterfaces:output_type -> inetmock.rpc.v1.ListAvailableNetworkInterfacesResponse
	2, // 6: inetmock.rpc.v1.NetFlowControlService.ListControlledInterfaces:output_type -> inetmock.rpc.v1.ListControlledInterfacesResponse
	6, // 7: inetmock.rpc.v1.NetFlowControlService.StartPacketFlowControl:output_type -> inetmock.rpc.v1.StartPacketFlowControlResponse
	8, // 8: inetmock.rpc.v1.NetFlowControlService.StopPacketFlowControl:output_type -> inetmock.rpc.v1.StopPacketFlowControlResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_v1_net_flow_proto_init() }
func file_rpc_v1_net_flow_proto_init() {
	if File_rpc_v1_net_flow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_v1_net_flow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListControlledInterfacesRequest); i {
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
		file_rpc_v1_net_flow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListControlledInterfacesResponse); i {
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
		file_rpc_v1_net_flow_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAvailableNetworkInterfacesRequest); i {
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
		file_rpc_v1_net_flow_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAvailableNetworkInterfacesResponse); i {
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
		file_rpc_v1_net_flow_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartPacketFlowControlRequest); i {
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
		file_rpc_v1_net_flow_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartPacketFlowControlResponse); i {
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
		file_rpc_v1_net_flow_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopPacketFlowControlRequest); i {
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
		file_rpc_v1_net_flow_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopPacketFlowControlResponse); i {
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
		file_rpc_v1_net_flow_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAvailableNetworkInterfacesResponse_NetworkInterface); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_v1_net_flow_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_v1_net_flow_proto_goTypes,
		DependencyIndexes: file_rpc_v1_net_flow_proto_depIdxs,
		EnumInfos:         file_rpc_v1_net_flow_proto_enumTypes,
		MessageInfos:      file_rpc_v1_net_flow_proto_msgTypes,
	}.Build()
	File_rpc_v1_net_flow_proto = out.File
	file_rpc_v1_net_flow_proto_rawDesc = nil
	file_rpc_v1_net_flow_proto_goTypes = nil
	file_rpc_v1_net_flow_proto_depIdxs = nil
}
