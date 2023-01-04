// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: audit/v1/netmon_details.proto

package auditv1

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

type NetMonDetailsEntity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReverseResolvedHost string `protobuf:"bytes,1,opt,name=reverse_resolved_host,json=reverseResolvedHost,proto3" json:"reverse_resolved_host,omitempty"`
}

func (x *NetMonDetailsEntity) Reset() {
	*x = NetMonDetailsEntity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_v1_netmon_details_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetMonDetailsEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetMonDetailsEntity) ProtoMessage() {}

func (x *NetMonDetailsEntity) ProtoReflect() protoreflect.Message {
	mi := &file_audit_v1_netmon_details_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetMonDetailsEntity.ProtoReflect.Descriptor instead.
func (*NetMonDetailsEntity) Descriptor() ([]byte, []int) {
	return file_audit_v1_netmon_details_proto_rawDescGZIP(), []int{0}
}

func (x *NetMonDetailsEntity) GetReverseResolvedHost() string {
	if x != nil {
		return x.ReverseResolvedHost
	}
	return ""
}

var File_audit_v1_netmon_details_proto protoreflect.FileDescriptor

var file_audit_v1_netmon_details_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x6e, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x11, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e,
	0x76, 0x31, 0x22, 0x49, 0x0a, 0x13, 0x4e, 0x65, 0x74, 0x4d, 0x6f, 0x6e, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x32, 0x0a, 0x15, 0x72, 0x65, 0x76,
	0x65, 0x72, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x64, 0x5f, 0x68, 0x6f,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x72, 0x65, 0x76, 0x65, 0x72, 0x73,
	0x65, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x64, 0x48, 0x6f, 0x73, 0x74, 0x42, 0xc6, 0x01,
	0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61,
	0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x12, 0x4e, 0x65, 0x74, 0x6d, 0x6f, 0x6e, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x02, 0x50, 0x01, 0x5a,
	0x31, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x69, 0x63, 0x62, 0x34, 0x64, 0x63,
	0x30, 0x2e, 0x64, 0x65, 0x2f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x75, 0x64, 0x69, 0x74,
	0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x41, 0x58, 0xaa, 0x02, 0x11, 0x49, 0x6e, 0x65, 0x74, 0x6d,
	0x6f, 0x63, 0x6b, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x11, 0x49,
	0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x5c, 0x41, 0x75, 0x64, 0x69, 0x74, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x1d, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x5c, 0x41, 0x75, 0x64, 0x69,
	0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x13, 0x49, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x3a, 0x3a, 0x41, 0x75, 0x64,
	0x69, 0x74, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_audit_v1_netmon_details_proto_rawDescOnce sync.Once
	file_audit_v1_netmon_details_proto_rawDescData = file_audit_v1_netmon_details_proto_rawDesc
)

func file_audit_v1_netmon_details_proto_rawDescGZIP() []byte {
	file_audit_v1_netmon_details_proto_rawDescOnce.Do(func() {
		file_audit_v1_netmon_details_proto_rawDescData = protoimpl.X.CompressGZIP(file_audit_v1_netmon_details_proto_rawDescData)
	})
	return file_audit_v1_netmon_details_proto_rawDescData
}

var file_audit_v1_netmon_details_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_audit_v1_netmon_details_proto_goTypes = []interface{}{
	(*NetMonDetailsEntity)(nil), // 0: inetmock.audit.v1.NetMonDetailsEntity
}
var file_audit_v1_netmon_details_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_audit_v1_netmon_details_proto_init() }
func file_audit_v1_netmon_details_proto_init() {
	if File_audit_v1_netmon_details_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_audit_v1_netmon_details_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetMonDetailsEntity); i {
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
			RawDescriptor: file_audit_v1_netmon_details_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_audit_v1_netmon_details_proto_goTypes,
		DependencyIndexes: file_audit_v1_netmon_details_proto_depIdxs,
		MessageInfos:      file_audit_v1_netmon_details_proto_msgTypes,
	}.Build()
	File_audit_v1_netmon_details_proto = out.File
	file_audit_v1_netmon_details_proto_rawDesc = nil
	file_audit_v1_netmon_details_proto_goTypes = nil
	file_audit_v1_netmon_details_proto_depIdxs = nil
}