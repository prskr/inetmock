// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.12.4
// source: rpc/health.proto

package rpc

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type HealthState int32

const (
	HealthState_HEALTHY      HealthState = 0
	HealthState_INITIALIZING HealthState = 1
	HealthState_UNHEALTHY    HealthState = 2
	HealthState_UNKNOWN      HealthState = 3
)

// Enum value maps for HealthState.
var (
	HealthState_name = map[int32]string{
		0: "HEALTHY",
		1: "INITIALIZING",
		2: "UNHEALTHY",
		3: "UNKNOWN",
	}
	HealthState_value = map[string]int32{
		"HEALTHY":      0,
		"INITIALIZING": 1,
		"UNHEALTHY":    2,
		"UNKNOWN":      3,
	}
)

func (x HealthState) Enum() *HealthState {
	p := new(HealthState)
	*p = x
	return p
}

func (x HealthState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HealthState) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_health_proto_enumTypes[0].Descriptor()
}

func (HealthState) Type() protoreflect.EnumType {
	return &file_rpc_health_proto_enumTypes[0]
}

func (x HealthState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HealthState.Descriptor instead.
func (HealthState) EnumDescriptor() ([]byte, []int) {
	return file_rpc_health_proto_rawDescGZIP(), []int{0}
}

type HealthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Components []string `protobuf:"bytes,1,rep,name=components,proto3" json:"components,omitempty"`
}

func (x *HealthRequest) Reset() {
	*x = HealthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_health_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthRequest) ProtoMessage() {}

func (x *HealthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_health_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthRequest.ProtoReflect.Descriptor instead.
func (*HealthRequest) Descriptor() ([]byte, []int) {
	return file_rpc_health_proto_rawDescGZIP(), []int{0}
}

func (x *HealthRequest) GetComponents() []string {
	if x != nil {
		return x.Components
	}
	return nil
}

type ComponentHealth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State   HealthState `protobuf:"varint,1,opt,name=State,proto3,enum=inetmock.rpc.HealthState" json:"State,omitempty"`
	Message string      `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ComponentHealth) Reset() {
	*x = ComponentHealth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_health_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComponentHealth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComponentHealth) ProtoMessage() {}

func (x *ComponentHealth) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_health_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComponentHealth.ProtoReflect.Descriptor instead.
func (*ComponentHealth) Descriptor() ([]byte, []int) {
	return file_rpc_health_proto_rawDescGZIP(), []int{1}
}

func (x *ComponentHealth) GetState() HealthState {
	if x != nil {
		return x.State
	}
	return HealthState_HEALTHY
}

func (x *ComponentHealth) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type HealthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OverallHealthState HealthState                 `protobuf:"varint,1,opt,name=overallHealthState,proto3,enum=inetmock.rpc.HealthState" json:"overallHealthState,omitempty"`
	ComponentsHealth   map[string]*ComponentHealth `protobuf:"bytes,2,rep,name=componentsHealth,proto3" json:"componentsHealth,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_health_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_health_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthResponse.ProtoReflect.Descriptor instead.
func (*HealthResponse) Descriptor() ([]byte, []int) {
	return file_rpc_health_proto_rawDescGZIP(), []int{2}
}

func (x *HealthResponse) GetOverallHealthState() HealthState {
	if x != nil {
		return x.OverallHealthState
	}
	return HealthState_HEALTHY
}

func (x *HealthResponse) GetComponentsHealth() map[string]*ComponentHealth {
	if x != nil {
		return x.ComponentsHealth
	}
	return nil
}

var File_rpc_health_proto protoreflect.FileDescriptor

var file_rpc_health_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x70, 0x63, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0c, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63,
	0x22, 0x2f, 0x0a, 0x0d, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74,
	0x73, 0x22, 0x5c, 0x0a, 0x0f, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x12, 0x2f, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x9f, 0x02, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x49, 0x0a, 0x12, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x6c, 0x6c, 0x48, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19,
	0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x12, 0x6f, 0x76, 0x65, 0x72, 0x61,
	0x6c, 0x6c, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x5e, 0x0a,
	0x10, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73,
	0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x10, 0x63, 0x6f, 0x6d,
	0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x1a, 0x62, 0x0a,
	0x15, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x33, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74,
	0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x2a, 0x48, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x00, 0x12, 0x10, 0x0a,
	0x0c, 0x49, 0x4e, 0x49, 0x54, 0x49, 0x41, 0x4c, 0x49, 0x5a, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12,
	0x0d, 0x0a, 0x09, 0x55, 0x4e, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x02, 0x12, 0x0b,
	0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x03, 0x32, 0x52, 0x0a, 0x06, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x48, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x48, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x12, 0x1b, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x70, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x62, 0x61,
	0x65, 0x7a, 0x39, 0x30, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x72, 0x70,
	0x63, 0x42, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x29, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x65,
	0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x13, 0x49, 0x4e,
	0x65, 0x74, 0x4d, 0x6f, 0x63, 0x6b, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x52, 0x70,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_health_proto_rawDescOnce sync.Once
	file_rpc_health_proto_rawDescData = file_rpc_health_proto_rawDesc
)

func file_rpc_health_proto_rawDescGZIP() []byte {
	file_rpc_health_proto_rawDescOnce.Do(func() {
		file_rpc_health_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_health_proto_rawDescData)
	})
	return file_rpc_health_proto_rawDescData
}

var file_rpc_health_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rpc_health_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_rpc_health_proto_goTypes = []interface{}{
	(HealthState)(0),        // 0: inetmock.rpc.HealthState
	(*HealthRequest)(nil),   // 1: inetmock.rpc.HealthRequest
	(*ComponentHealth)(nil), // 2: inetmock.rpc.ComponentHealth
	(*HealthResponse)(nil),  // 3: inetmock.rpc.HealthResponse
	nil,                     // 4: inetmock.rpc.HealthResponse.ComponentsHealthEntry
}
var file_rpc_health_proto_depIdxs = []int32{
	0, // 0: inetmock.rpc.ComponentHealth.State:type_name -> inetmock.rpc.HealthState
	0, // 1: inetmock.rpc.HealthResponse.overallHealthState:type_name -> inetmock.rpc.HealthState
	4, // 2: inetmock.rpc.HealthResponse.componentsHealth:type_name -> inetmock.rpc.HealthResponse.ComponentsHealthEntry
	2, // 3: inetmock.rpc.HealthResponse.ComponentsHealthEntry.value:type_name -> inetmock.rpc.ComponentHealth
	1, // 4: inetmock.rpc.Health.GetHealth:input_type -> inetmock.rpc.HealthRequest
	3, // 5: inetmock.rpc.Health.GetHealth:output_type -> inetmock.rpc.HealthResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_rpc_health_proto_init() }
func file_rpc_health_proto_init() {
	if File_rpc_health_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_health_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthRequest); i {
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
		file_rpc_health_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComponentHealth); i {
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
		file_rpc_health_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthResponse); i {
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
			RawDescriptor: file_rpc_health_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_health_proto_goTypes,
		DependencyIndexes: file_rpc_health_proto_depIdxs,
		EnumInfos:         file_rpc_health_proto_enumTypes,
		MessageInfos:      file_rpc_health_proto_msgTypes,
	}.Build()
	File_rpc_health_proto = out.File
	file_rpc_health_proto_rawDesc = nil
	file_rpc_health_proto_goTypes = nil
	file_rpc_health_proto_depIdxs = nil
}
