// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: rpc/v1/audit.proto

package rpcv1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	v1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type WatchEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WatcherName string `protobuf:"bytes,1,opt,name=watcher_name,json=watcherName,proto3" json:"watcher_name,omitempty"`
}

func (x *WatchEventsRequest) Reset() {
	*x = WatchEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_audit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchEventsRequest) ProtoMessage() {}

func (x *WatchEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_audit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchEventsRequest.ProtoReflect.Descriptor instead.
func (*WatchEventsRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_audit_proto_rawDescGZIP(), []int{0}
}

func (x *WatchEventsRequest) GetWatcherName() string {
	if x != nil {
		return x.WatcherName
	}
	return ""
}

type WatchEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entity *v1.EventEntity `protobuf:"bytes,1,opt,name=entity,proto3" json:"entity,omitempty"`
}

func (x *WatchEventsResponse) Reset() {
	*x = WatchEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_audit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchEventsResponse) ProtoMessage() {}

func (x *WatchEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_audit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchEventsResponse.ProtoReflect.Descriptor instead.
func (*WatchEventsResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_audit_proto_rawDescGZIP(), []int{1}
}

func (x *WatchEventsResponse) GetEntity() *v1.EventEntity {
	if x != nil {
		return x.Entity
	}
	return nil
}

type RegisterFileSinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetPath string `protobuf:"bytes,1,opt,name=target_path,json=targetPath,proto3" json:"target_path,omitempty"`
}

func (x *RegisterFileSinkRequest) Reset() {
	*x = RegisterFileSinkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_audit_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterFileSinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterFileSinkRequest) ProtoMessage() {}

func (x *RegisterFileSinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_audit_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterFileSinkRequest.ProtoReflect.Descriptor instead.
func (*RegisterFileSinkRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_audit_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterFileSinkRequest) GetTargetPath() string {
	if x != nil {
		return x.TargetPath
	}
	return ""
}

type RegisterFileSinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResolvedPath string `protobuf:"bytes,1,opt,name=resolved_path,json=resolvedPath,proto3" json:"resolved_path,omitempty"`
}

func (x *RegisterFileSinkResponse) Reset() {
	*x = RegisterFileSinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_audit_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterFileSinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterFileSinkResponse) ProtoMessage() {}

func (x *RegisterFileSinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_audit_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterFileSinkResponse.ProtoReflect.Descriptor instead.
func (*RegisterFileSinkResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_audit_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterFileSinkResponse) GetResolvedPath() string {
	if x != nil {
		return x.ResolvedPath
	}
	return ""
}

type RemoveFileSinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetPath string `protobuf:"bytes,1,opt,name=target_path,json=targetPath,proto3" json:"target_path,omitempty"`
}

func (x *RemoveFileSinkRequest) Reset() {
	*x = RemoveFileSinkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_audit_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveFileSinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveFileSinkRequest) ProtoMessage() {}

func (x *RemoveFileSinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_audit_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveFileSinkRequest.ProtoReflect.Descriptor instead.
func (*RemoveFileSinkRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_audit_proto_rawDescGZIP(), []int{4}
}

func (x *RemoveFileSinkRequest) GetTargetPath() string {
	if x != nil {
		return x.TargetPath
	}
	return ""
}

type RemoveFileSinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SinkGotRemoved bool `protobuf:"varint,1,opt,name=sink_got_removed,json=sinkGotRemoved,proto3" json:"sink_got_removed,omitempty"`
}

func (x *RemoveFileSinkResponse) Reset() {
	*x = RemoveFileSinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_audit_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveFileSinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveFileSinkResponse) ProtoMessage() {}

func (x *RemoveFileSinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_audit_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveFileSinkResponse.ProtoReflect.Descriptor instead.
func (*RemoveFileSinkResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_audit_proto_rawDescGZIP(), []int{5}
}

func (x *RemoveFileSinkResponse) GetSinkGotRemoved() bool {
	if x != nil {
		return x.SinkGotRemoved
	}
	return false
}

type ListSinksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListSinksRequest) Reset() {
	*x = ListSinksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_audit_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListSinksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSinksRequest) ProtoMessage() {}

func (x *ListSinksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_audit_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSinksRequest.ProtoReflect.Descriptor instead.
func (*ListSinksRequest) Descriptor() ([]byte, []int) {
	return file_rpc_v1_audit_proto_rawDescGZIP(), []int{6}
}

type ListSinksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sinks []string `protobuf:"bytes,1,rep,name=sinks,proto3" json:"sinks,omitempty"`
}

func (x *ListSinksResponse) Reset() {
	*x = ListSinksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_v1_audit_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListSinksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSinksResponse) ProtoMessage() {}

func (x *ListSinksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_v1_audit_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSinksResponse.ProtoReflect.Descriptor instead.
func (*ListSinksResponse) Descriptor() ([]byte, []int) {
	return file_rpc_v1_audit_proto_rawDescGZIP(), []int{7}
}

func (x *ListSinksResponse) GetSinks() []string {
	if x != nil {
		return x.Sinks
	}
	return nil
}

var File_rpc_v1_audit_proto protoreflect.FileDescriptor

var file_rpc_v1_audit_proto_rawDesc = []byte{
	0x0a, 0x12, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x37, 0x0a, 0x12, 0x57, 0x61, 0x74, 0x63, 0x68, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x77, 0x61, 0x74, 0x63,
	0x68, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x77, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4d, 0x0a, 0x13, 0x57,
	0x61, 0x74, 0x63, 0x68, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x36, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75,
	0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x52, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x3a, 0x0a, 0x17, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f,
	0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x50, 0x61, 0x74, 0x68, 0x22, 0x3f, 0x0a, 0x18, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x64, 0x5f, 0x70,
	0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x6c,
	0x76, 0x65, 0x64, 0x50, 0x61, 0x74, 0x68, 0x22, 0x38, 0x0a, 0x15, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x61, 0x74,
	0x68, 0x22, 0x42, 0x0a, 0x16, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x53,
	0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x73,
	0x69, 0x6e, 0x6b, 0x5f, 0x67, 0x6f, 0x74, 0x5f, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x73, 0x69, 0x6e, 0x6b, 0x47, 0x6f, 0x74, 0x52, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x64, 0x22, 0x12, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x69, 0x6e,
	0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x29, 0x0a, 0x11, 0x4c, 0x69, 0x73,
	0x74, 0x53, 0x69, 0x6e, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x69, 0x6e, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x73,
	0x69, 0x6e, 0x6b, 0x73, 0x32, 0x8a, 0x03, 0x0a, 0x0c, 0x41, 0x75, 0x64, 0x69, 0x74, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x0b, 0x57, 0x61, 0x74, 0x63, 0x68, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x23, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x69, 0x6e, 0x65, 0x74,
	0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x61, 0x74, 0x63,
	0x68, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30,
	0x01, 0x12, 0x67, 0x0a, 0x10, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x46, 0x69, 0x6c,
	0x65, 0x53, 0x69, 0x6e, 0x6b, 0x12, 0x28, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x29, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69,
	0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x61, 0x0a, 0x0e, 0x52, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x6e, 0x6b, 0x12, 0x26, 0x2e, 0x69,
	0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x46, 0x69, 0x6c,
	0x65, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x52, 0x0a,
	0x09, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x69, 0x6e, 0x6b, 0x73, 0x12, 0x21, 0x2e, 0x69, 0x6e, 0x65,
	0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x53, 0x69, 0x6e, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e,
	0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x53, 0x69, 0x6e, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0xb0, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x41, 0x75, 0x64, 0x69, 0x74,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x02, 0x50, 0x01, 0x5a, 0x2d, 0x69, 0x6e, 0x65, 0x74, 0x6d,
	0x6f, 0x63, 0x6b, 0x2e, 0x69, 0x63, 0x62, 0x34, 0x64, 0x63, 0x30, 0x2e, 0x64, 0x65, 0x2f, 0x69,
	0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f,
	0x76, 0x31, 0x3b, 0x72, 0x70, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x52, 0x58, 0xaa, 0x02,
	0x0f, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x0f, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x5c, 0x52, 0x70, 0x63, 0x5c,
	0x56, 0x31, 0xe2, 0x02, 0x1b, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x5c, 0x52, 0x70,
	0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x11, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x3a, 0x3a, 0x52, 0x70, 0x63,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_v1_audit_proto_rawDescOnce sync.Once
	file_rpc_v1_audit_proto_rawDescData = file_rpc_v1_audit_proto_rawDesc
)

func file_rpc_v1_audit_proto_rawDescGZIP() []byte {
	file_rpc_v1_audit_proto_rawDescOnce.Do(func() {
		file_rpc_v1_audit_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_v1_audit_proto_rawDescData)
	})
	return file_rpc_v1_audit_proto_rawDescData
}

var file_rpc_v1_audit_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_rpc_v1_audit_proto_goTypes = []interface{}{
	(*WatchEventsRequest)(nil),       // 0: inetmock.rpc.v1.WatchEventsRequest
	(*WatchEventsResponse)(nil),      // 1: inetmock.rpc.v1.WatchEventsResponse
	(*RegisterFileSinkRequest)(nil),  // 2: inetmock.rpc.v1.RegisterFileSinkRequest
	(*RegisterFileSinkResponse)(nil), // 3: inetmock.rpc.v1.RegisterFileSinkResponse
	(*RemoveFileSinkRequest)(nil),    // 4: inetmock.rpc.v1.RemoveFileSinkRequest
	(*RemoveFileSinkResponse)(nil),   // 5: inetmock.rpc.v1.RemoveFileSinkResponse
	(*ListSinksRequest)(nil),         // 6: inetmock.rpc.v1.ListSinksRequest
	(*ListSinksResponse)(nil),        // 7: inetmock.rpc.v1.ListSinksResponse
	(*v1.EventEntity)(nil),           // 8: inetmock.audit.v1.EventEntity
}
var file_rpc_v1_audit_proto_depIdxs = []int32{
	8, // 0: inetmock.rpc.v1.WatchEventsResponse.entity:type_name -> inetmock.audit.v1.EventEntity
	0, // 1: inetmock.rpc.v1.AuditService.WatchEvents:input_type -> inetmock.rpc.v1.WatchEventsRequest
	2, // 2: inetmock.rpc.v1.AuditService.RegisterFileSink:input_type -> inetmock.rpc.v1.RegisterFileSinkRequest
	4, // 3: inetmock.rpc.v1.AuditService.RemoveFileSink:input_type -> inetmock.rpc.v1.RemoveFileSinkRequest
	6, // 4: inetmock.rpc.v1.AuditService.ListSinks:input_type -> inetmock.rpc.v1.ListSinksRequest
	1, // 5: inetmock.rpc.v1.AuditService.WatchEvents:output_type -> inetmock.rpc.v1.WatchEventsResponse
	3, // 6: inetmock.rpc.v1.AuditService.RegisterFileSink:output_type -> inetmock.rpc.v1.RegisterFileSinkResponse
	5, // 7: inetmock.rpc.v1.AuditService.RemoveFileSink:output_type -> inetmock.rpc.v1.RemoveFileSinkResponse
	7, // 8: inetmock.rpc.v1.AuditService.ListSinks:output_type -> inetmock.rpc.v1.ListSinksResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_v1_audit_proto_init() }
func file_rpc_v1_audit_proto_init() {
	if File_rpc_v1_audit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_v1_audit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchEventsRequest); i {
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
		file_rpc_v1_audit_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchEventsResponse); i {
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
		file_rpc_v1_audit_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterFileSinkRequest); i {
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
		file_rpc_v1_audit_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterFileSinkResponse); i {
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
		file_rpc_v1_audit_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveFileSinkRequest); i {
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
		file_rpc_v1_audit_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveFileSinkResponse); i {
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
		file_rpc_v1_audit_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListSinksRequest); i {
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
		file_rpc_v1_audit_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListSinksResponse); i {
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
			RawDescriptor: file_rpc_v1_audit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_v1_audit_proto_goTypes,
		DependencyIndexes: file_rpc_v1_audit_proto_depIdxs,
		MessageInfos:      file_rpc_v1_audit_proto_msgTypes,
	}.Build()
	File_rpc_v1_audit_proto = out.File
	file_rpc_v1_audit_proto_rawDesc = nil
	file_rpc_v1_audit_proto_goTypes = nil
	file_rpc_v1_audit_proto_depIdxs = nil
}
