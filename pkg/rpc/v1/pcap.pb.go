// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: rpc/v1/pcap.proto

package rpcv1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListAvailableDevicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListAvailableDevicesRequest) Reset() {
	*x = ListAvailableDevicesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAvailableDevicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAvailableDevicesRequest) ProtoMessage() {}

func (x *ListAvailableDevicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAvailableDevicesRequest.ProtoReflect.Descriptor instead.
func (*ListAvailableDevicesRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{0}
}

type ListAvailableDevicesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AvailableDevices []*ListAvailableDevicesResponse_PCAPDevice `protobuf:"bytes,1,rep,name=available_devices,json=availableDevices,proto3" json:"available_devices,omitempty"`
}

func (x *ListAvailableDevicesResponse) Reset() {
	*x = ListAvailableDevicesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAvailableDevicesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAvailableDevicesResponse) ProtoMessage() {}

func (x *ListAvailableDevicesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAvailableDevicesResponse.ProtoReflect.Descriptor instead.
func (*ListAvailableDevicesResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{1}
}

func (x *ListAvailableDevicesResponse) GetAvailableDevices() []*ListAvailableDevicesResponse_PCAPDevice {
	if x != nil {
		return x.AvailableDevices
	}
	return nil
}

type ListActiveRecordingsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListActiveRecordingsRequest) Reset() {
	*x = ListActiveRecordingsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListActiveRecordingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListActiveRecordingsRequest) ProtoMessage() {}

func (x *ListActiveRecordingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListActiveRecordingsRequest.ProtoReflect.Descriptor instead.
func (*ListActiveRecordingsRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{2}
}

type ListActiveRecordingsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subscriptions []string `protobuf:"bytes,1,rep,name=subscriptions,proto3" json:"subscriptions,omitempty"`
}

func (x *ListActiveRecordingsResponse) Reset() {
	*x = ListActiveRecordingsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListActiveRecordingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListActiveRecordingsResponse) ProtoMessage() {}

func (x *ListActiveRecordingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListActiveRecordingsResponse.ProtoReflect.Descriptor instead.
func (*ListActiveRecordingsResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{3}
}

func (x *ListActiveRecordingsResponse) GetSubscriptions() []string {
	if x != nil {
		return x.Subscriptions
	}
	return nil
}

type StartPCAPFileRecordingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Device      string               `protobuf:"bytes,1,opt,name=device,proto3" json:"device,omitempty"`
	TargetPath  string               `protobuf:"bytes,2,opt,name=target_path,json=targetPath,proto3" json:"target_path,omitempty"`
	Promiscuous bool                 `protobuf:"varint,3,opt,name=promiscuous,proto3" json:"promiscuous,omitempty"`
	ReadTimeout *durationpb.Duration `protobuf:"bytes,4,opt,name=read_timeout,json=readTimeout,proto3" json:"read_timeout,omitempty"`
}

func (x *StartPCAPFileRecordingRequest) Reset() {
	*x = StartPCAPFileRecordingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartPCAPFileRecordingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartPCAPFileRecordingRequest) ProtoMessage() {}

func (x *StartPCAPFileRecordingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartPCAPFileRecordingRequest.ProtoReflect.Descriptor instead.
func (*StartPCAPFileRecordingRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{4}
}

func (x *StartPCAPFileRecordingRequest) GetDevice() string {
	if x != nil {
		return x.Device
	}
	return ""
}

func (x *StartPCAPFileRecordingRequest) GetTargetPath() string {
	if x != nil {
		return x.TargetPath
	}
	return ""
}

func (x *StartPCAPFileRecordingRequest) GetPromiscuous() bool {
	if x != nil {
		return x.Promiscuous
	}
	return false
}

func (x *StartPCAPFileRecordingRequest) GetReadTimeout() *durationpb.Duration {
	if x != nil {
		return x.ReadTimeout
	}
	return nil
}

type StartPCAPFileRecordingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResolvedPath string `protobuf:"bytes,1,opt,name=resolved_path,json=resolvedPath,proto3" json:"resolved_path,omitempty"`
	ConsumerKey  string `protobuf:"bytes,2,opt,name=consumer_key,json=consumerKey,proto3" json:"consumer_key,omitempty"`
}

func (x *StartPCAPFileRecordingResponse) Reset() {
	*x = StartPCAPFileRecordingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartPCAPFileRecordingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartPCAPFileRecordingResponse) ProtoMessage() {}

func (x *StartPCAPFileRecordingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartPCAPFileRecordingResponse.ProtoReflect.Descriptor instead.
func (*StartPCAPFileRecordingResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{5}
}

func (x *StartPCAPFileRecordingResponse) GetResolvedPath() string {
	if x != nil {
		return x.ResolvedPath
	}
	return ""
}

func (x *StartPCAPFileRecordingResponse) GetConsumerKey() string {
	if x != nil {
		return x.ConsumerKey
	}
	return ""
}

type StopPCAPFileRecordingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConsumerKey string `protobuf:"bytes,1,opt,name=consumer_key,json=consumerKey,proto3" json:"consumer_key,omitempty"`
}

func (x *StopPCAPFileRecordingRequest) Reset() {
	*x = StopPCAPFileRecordingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopPCAPFileRecordingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopPCAPFileRecordingRequest) ProtoMessage() {}

func (x *StopPCAPFileRecordingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopPCAPFileRecordingRequest.ProtoReflect.Descriptor instead.
func (*StopPCAPFileRecordingRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{6}
}

func (x *StopPCAPFileRecordingRequest) GetConsumerKey() string {
	if x != nil {
		return x.ConsumerKey
	}
	return ""
}

type StopPCAPFileRecordingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Removed bool `protobuf:"varint,1,opt,name=removed,proto3" json:"removed,omitempty"`
}

func (x *StopPCAPFileRecordingResponse) Reset() {
	*x = StopPCAPFileRecordingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopPCAPFileRecordingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopPCAPFileRecordingResponse) ProtoMessage() {}

func (x *StopPCAPFileRecordingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopPCAPFileRecordingResponse.ProtoReflect.Descriptor instead.
func (*StopPCAPFileRecordingResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{7}
}

func (x *StopPCAPFileRecordingResponse) GetRemoved() bool {
	if x != nil {
		return x.Removed
	}
	return false
}

type ListAvailableDevicesResponse_PCAPDevice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Addresses [][]byte `protobuf:"bytes,2,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *ListAvailableDevicesResponse_PCAPDevice) Reset() {
	*x = ListAvailableDevicesResponse_PCAPDevice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_pcap_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAvailableDevicesResponse_PCAPDevice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAvailableDevicesResponse_PCAPDevice) ProtoMessage() {}

func (x *ListAvailableDevicesResponse_PCAPDevice) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_pcap_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAvailableDevicesResponse_PCAPDevice.ProtoReflect.Descriptor instead.
func (*ListAvailableDevicesResponse_PCAPDevice) Descriptor() ([]byte, []int) {
	return file_rpc_v1_pcap_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ListAvailableDevicesResponse_PCAPDevice) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListAvailableDevicesResponse_PCAPDevice) GetAddresses() [][]byte {
	if x != nil {
		return x.Addresses
	}
	return nil
}

var File_rpc_v1_pcap_proto protoreflect.FileDescriptor

var file_rpc_v1_pcap_proto_rawDesc = []byte{
	0x0a, 0x11, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x63, 0x61, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1d, 0x0a, 0x1b, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0xc5, 0x01, 0x0a, 0x1c, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x11, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x5f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x38, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x50,
	0x43, 0x41, 0x50, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x3e, 0x0a, 0x0a, 0x50,
	0x43, 0x41, 0x50, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c,
	0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0x1d, 0x0a, 0x1b, 0x4c,
	0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69,
	0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x44, 0x0a, 0x1c, 0x4c, 0x69,
	0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e,
	0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x0d, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x22, 0xb8, 0x01, 0x0a, 0x1d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x43, 0x41, 0x50, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x70,
	0x72, 0x6f, 0x6d, 0x69, 0x73, 0x63, 0x75, 0x6f, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0b, 0x70, 0x72, 0x6f, 0x6d, 0x69, 0x73, 0x63, 0x75, 0x6f, 0x75, 0x73, 0x12, 0x3c, 0x0a,
	0x0c, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b,
	0x72, 0x65, 0x61, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x68, 0x0a, 0x1e, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x50, 0x43, 0x41, 0x50, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a,
	0x0d, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x64, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x64, 0x50, 0x61,
	0x74, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d,
	0x65, 0x72, 0x4b, 0x65, 0x79, 0x22, 0x41, 0x0a, 0x1c, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x43, 0x41,
	0x50, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65,
	0x72, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e,
	0x73, 0x75, 0x6d, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x22, 0x39, 0x0a, 0x1d, 0x53, 0x74, 0x6f, 0x70,
	0x50, 0x43, 0x41, 0x50, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x6d,
	0x6f, 0x76, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x64, 0x32, 0xea, 0x03, 0x0a, 0x0b, 0x50, 0x43, 0x41, 0x50, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x73, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x2c, 0x2e, 0x69, 0x6e,
	0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x69, 0x6e, 0x65, 0x74,
	0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x73, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x12, 0x2c, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d,
	0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x79, 0x0a,
	0x16, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x43, 0x41, 0x50, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x2e, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50,
	0x43, 0x41, 0x50, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50,
	0x43, 0x41, 0x50, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x76, 0x0a, 0x15, 0x53, 0x74, 0x6f, 0x70,
	0x50, 0x43, 0x41, 0x50, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e,
	0x67, 0x12, 0x2d, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x43, 0x41, 0x50, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2e, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x43, 0x41, 0x50, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0xaf, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63,
	0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x50, 0x63, 0x61, 0x70, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x48, 0x02, 0x50, 0x01, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f, 0x69, 0x6e, 0x65,
	0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31,
	0x3b, 0x72, 0x70, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x52, 0x58, 0xaa, 0x02, 0x0f, 0x49,
	0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x0f, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x1b, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x5c, 0x52, 0x70, 0x63, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x11, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_v1_pcap_proto_rawDescOnce sync.Once
	file_rpc_v1_pcap_proto_rawDescData = file_rpc_v1_pcap_proto_rawDesc
)

func file_rpc_v1_pcap_proto_rawDescGZIP() []byte {
	file_rpc_v1_pcap_proto_rawDescOnce.Do(func() {
		file_rpc_v1_pcap_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_v1_pcap_proto_rawDescData)
	})
	return file_rpc_v1_pcap_proto_rawDescData
}

var file_rpc_v1_pcap_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_rpc_v1_pcap_proto_goTypes = []interface{}{
	(*ListAvailableDevicesRequest)(nil),             // 0: inetmock.rpc.v1.ListAvailableDevicesRequest
	(*ListAvailableDevicesResponse)(nil),            // 1: inetmock.rpc.v1.ListAvailableDevicesResponse
	(*ListActiveRecordingsRequest)(nil),             // 2: inetmock.rpc.v1.ListActiveRecordingsRequest
	(*ListActiveRecordingsResponse)(nil),            // 3: inetmock.rpc.v1.ListActiveRecordingsResponse
	(*StartPCAPFileRecordingRequest)(nil),           // 4: inetmock.rpc.v1.StartPCAPFileRecordingRequest
	(*StartPCAPFileRecordingResponse)(nil),          // 5: inetmock.rpc.v1.StartPCAPFileRecordingResponse
	(*StopPCAPFileRecordingRequest)(nil),            // 6: inetmock.rpc.v1.StopPCAPFileRecordingRequest
	(*StopPCAPFileRecordingResponse)(nil),           // 7: inetmock.rpc.v1.StopPCAPFileRecordingResponse
	(*ListAvailableDevicesResponse_PCAPDevice)(nil), // 8: inetmock.rpc.v1.ListAvailableDevicesResponse.PCAPDevice
	(*durationpb.Duration)(nil),                     // 9: google.protobuf.Duration
}
var file_rpc_v1_pcap_proto_depIdxs = []int32{
	8, // 0: inetmock.rpc.v1.ListAvailableDevicesResponse.available_devices:type_name -> inetmock.rpc.v1.ListAvailableDevicesResponse.PCAPDevice
	9, // 1: inetmock.rpc.v1.StartPCAPFileRecordingRequest.read_timeout:type_name -> google.protobuf.Duration
	0, // 2: inetmock.rpc.v1.PCAPService.ListAvailableDevices:input_type -> inetmock.rpc.v1.ListAvailableDevicesRequest
	2, // 3: inetmock.rpc.v1.PCAPService.ListActiveRecordings:input_type -> inetmock.rpc.v1.ListActiveRecordingsRequest
	4, // 4: inetmock.rpc.v1.PCAPService.StartPCAPFileRecording:input_type -> inetmock.rpc.v1.StartPCAPFileRecordingRequest
	6, // 5: inetmock.rpc.v1.PCAPService.StopPCAPFileRecording:input_type -> inetmock.rpc.v1.StopPCAPFileRecordingRequest
	1, // 6: inetmock.rpc.v1.PCAPService.ListAvailableDevices:output_type -> inetmock.rpc.v1.ListAvailableDevicesResponse
	3, // 7: inetmock.rpc.v1.PCAPService.ListActiveRecordings:output_type -> inetmock.rpc.v1.ListActiveRecordingsResponse
	5, // 8: inetmock.rpc.v1.PCAPService.StartPCAPFileRecording:output_type -> inetmock.rpc.v1.StartPCAPFileRecordingResponse
	7, // 9: inetmock.rpc.v1.PCAPService.StopPCAPFileRecording:output_type -> inetmock.rpc.v1.StopPCAPFileRecordingResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rpc_v1_pcap_proto_init() }
func file_rpc_v1_pcap_proto_init() {
	if File_rpc_v1_pcap_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_v1_pcap_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAvailableDevicesRequest); i {
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
		file_rpc_v1_pcap_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAvailableDevicesResponse); i {
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
		file_rpc_v1_pcap_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListActiveRecordingsRequest); i {
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
		file_rpc_v1_pcap_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListActiveRecordingsResponse); i {
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
		file_rpc_v1_pcap_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartPCAPFileRecordingRequest); i {
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
		file_rpc_v1_pcap_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartPCAPFileRecordingResponse); i {
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
		file_rpc_v1_pcap_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopPCAPFileRecordingRequest); i {
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
		file_rpc_v1_pcap_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopPCAPFileRecordingResponse); i {
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
		file_rpc_v1_pcap_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAvailableDevicesResponse_PCAPDevice); i {
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
			RawDescriptor: file_rpc_v1_pcap_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_v1_pcap_proto_goTypes,
		DependencyIndexes: file_rpc_v1_pcap_proto_depIdxs,
		MessageInfos:      file_rpc_v1_pcap_proto_msgTypes,
	}.Build()
	File_rpc_v1_pcap_proto = out.File
	file_rpc_v1_pcap_proto_rawDesc = nil
	file_rpc_v1_pcap_proto_goTypes = nil
	file_rpc_v1_pcap_proto_depIdxs = nil
}
