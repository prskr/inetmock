// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.2
// source: audit/v1/http_details.proto

package v1

import (
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

type HTTPMethod int32

const (
	HTTPMethod_HTTP_METHOD_UNSPECIFIED HTTPMethod = 0
	HTTPMethod_HTTP_METHOD_GET         HTTPMethod = 1
	HTTPMethod_HTTP_METHOD_HEAD        HTTPMethod = 2
	HTTPMethod_HTTP_METHOD_POST        HTTPMethod = 3
	HTTPMethod_HTTP_METHOD_PUT         HTTPMethod = 4
	HTTPMethod_HTTP_METHOD_DELETE      HTTPMethod = 5
	HTTPMethod_HTTP_METHOD_CONNECT     HTTPMethod = 6
	HTTPMethod_HTTP_METHOD_OPTIONS     HTTPMethod = 7
	HTTPMethod_HTTP_METHOD_TRACE       HTTPMethod = 8
	HTTPMethod_HTTP_METHOD_PATCH       HTTPMethod = 9
)

// Enum value maps for HTTPMethod.
var (
	HTTPMethod_name = map[int32]string{
		0: "HTTP_METHOD_UNSPECIFIED",
		1: "HTTP_METHOD_GET",
		2: "HTTP_METHOD_HEAD",
		3: "HTTP_METHOD_POST",
		4: "HTTP_METHOD_PUT",
		5: "HTTP_METHOD_DELETE",
		6: "HTTP_METHOD_CONNECT",
		7: "HTTP_METHOD_OPTIONS",
		8: "HTTP_METHOD_TRACE",
		9: "HTTP_METHOD_PATCH",
	}
	HTTPMethod_value = map[string]int32{
		"HTTP_METHOD_UNSPECIFIED": 0,
		"HTTP_METHOD_GET":         1,
		"HTTP_METHOD_HEAD":        2,
		"HTTP_METHOD_POST":        3,
		"HTTP_METHOD_PUT":         4,
		"HTTP_METHOD_DELETE":      5,
		"HTTP_METHOD_CONNECT":     6,
		"HTTP_METHOD_OPTIONS":     7,
		"HTTP_METHOD_TRACE":       8,
		"HTTP_METHOD_PATCH":       9,
	}
)

func (x HTTPMethod) Enum() *HTTPMethod {
	p := new(HTTPMethod)
	*p = x
	return p
}

func (x HTTPMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HTTPMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_audit_v1_http_details_proto_enumTypes[0].Descriptor()
}

func (HTTPMethod) Type() protoreflect.EnumType {
	return &file_audit_v1_http_details_proto_enumTypes[0]
}

func (x HTTPMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HTTPMethod.Descriptor instead.
func (HTTPMethod) EnumDescriptor() ([]byte, []int) {
	return file_audit_v1_http_details_proto_rawDescGZIP(), []int{0}
}

type HTTPHeaderValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *HTTPHeaderValue) Reset() {
	*x = HTTPHeaderValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_v1_http_details_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HTTPHeaderValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPHeaderValue) ProtoMessage() {}

func (x *HTTPHeaderValue) ProtoReflect() protoreflect.Message {
	mi := &file_audit_v1_http_details_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPHeaderValue.ProtoReflect.Descriptor instead.
func (*HTTPHeaderValue) Descriptor() ([]byte, []int) {
	return file_audit_v1_http_details_proto_rawDescGZIP(), []int{0}
}

func (x *HTTPHeaderValue) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type HTTPDetailsEntity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method  HTTPMethod                  `protobuf:"varint,1,opt,name=method,proto3,enum=inetmock.audit.v1.HTTPMethod" json:"method,omitempty"`
	Host    string                      `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
	Uri     string                      `protobuf:"bytes,3,opt,name=uri,proto3" json:"uri,omitempty"`
	Proto   string                      `protobuf:"bytes,4,opt,name=proto,proto3" json:"proto,omitempty"`
	Headers map[string]*HTTPHeaderValue `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *HTTPDetailsEntity) Reset() {
	*x = HTTPDetailsEntity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_v1_http_details_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HTTPDetailsEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPDetailsEntity) ProtoMessage() {}

func (x *HTTPDetailsEntity) ProtoReflect() protoreflect.Message {
	mi := &file_audit_v1_http_details_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPDetailsEntity.ProtoReflect.Descriptor instead.
func (*HTTPDetailsEntity) Descriptor() ([]byte, []int) {
	return file_audit_v1_http_details_proto_rawDescGZIP(), []int{1}
}

func (x *HTTPDetailsEntity) GetMethod() HTTPMethod {
	if x != nil {
		return x.Method
	}
	return HTTPMethod_HTTP_METHOD_UNSPECIFIED
}

func (x *HTTPDetailsEntity) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *HTTPDetailsEntity) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

func (x *HTTPDetailsEntity) GetProto() string {
	if x != nil {
		return x.Proto
	}
	return ""
}

func (x *HTTPDetailsEntity) GetHeaders() map[string]*HTTPHeaderValue {
	if x != nil {
		return x.Headers
	}
	return nil
}

var File_audit_v1_http_details_proto protoreflect.FileDescriptor

var file_audit_v1_http_details_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x5f,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x69,
	0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31,
	0x22, 0x29, 0x0a, 0x0f, 0x48, 0x54, 0x54, 0x50, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0xb3, 0x02, 0x0a, 0x11,
	0x48, 0x54, 0x54, 0x50, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x12, 0x35, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1d, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64,
	0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x69, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x69, 0x12, 0x14,
	0x0a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x4b, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18,
	0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b,
	0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x1a, 0x5e, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x38, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75,
	0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x2a, 0xf7, 0x01, 0x0a, 0x0a, 0x48, 0x54, 0x54, 0x50, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x12, 0x1b, 0x0a, 0x17, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a,
	0x0f, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x47, 0x45, 0x54,
	0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f,
	0x44, 0x5f, 0x48, 0x45, 0x41, 0x44, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x48, 0x54, 0x54, 0x50,
	0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x50, 0x4f, 0x53, 0x54, 0x10, 0x03, 0x12, 0x13,
	0x0a, 0x0f, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x50, 0x55,
	0x54, 0x10, 0x04, 0x12, 0x16, 0x0a, 0x12, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54, 0x48,
	0x4f, 0x44, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x05, 0x12, 0x17, 0x0a, 0x13, 0x48,
	0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45,
	0x43, 0x54, 0x10, 0x06, 0x12, 0x17, 0x0a, 0x13, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54,
	0x48, 0x4f, 0x44, 0x5f, 0x4f, 0x50, 0x54, 0x49, 0x4f, 0x4e, 0x53, 0x10, 0x07, 0x12, 0x15, 0x0a,
	0x11, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x54, 0x52, 0x41,
	0x43, 0x45, 0x10, 0x08, 0x12, 0x15, 0x0a, 0x11, 0x48, 0x54, 0x54, 0x50, 0x5f, 0x4d, 0x45, 0x54,
	0x48, 0x4f, 0x44, 0x5f, 0x50, 0x41, 0x54, 0x43, 0x48, 0x10, 0x09, 0x42, 0x7a, 0x0a, 0x20, 0x63,
	0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x62, 0x61, 0x65, 0x7a, 0x39, 0x30,
	0x2e, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x42,
	0x11, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f, 0x63, 0x6b, 0x2f, 0x69, 0x6e, 0x65, 0x74, 0x6d, 0x6f,
	0x63, 0x6b, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0xaa,
	0x02, 0x15, 0x49, 0x4e, 0x65, 0x74, 0x4d, 0x6f, 0x63, 0x6b, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_audit_v1_http_details_proto_rawDescOnce sync.Once
	file_audit_v1_http_details_proto_rawDescData = file_audit_v1_http_details_proto_rawDesc
)

func file_audit_v1_http_details_proto_rawDescGZIP() []byte {
	file_audit_v1_http_details_proto_rawDescOnce.Do(func() {
		file_audit_v1_http_details_proto_rawDescData = protoimpl.X.CompressGZIP(file_audit_v1_http_details_proto_rawDescData)
	})
	return file_audit_v1_http_details_proto_rawDescData
}

var file_audit_v1_http_details_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_audit_v1_http_details_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_audit_v1_http_details_proto_goTypes = []interface{}{
	(HTTPMethod)(0),           // 0: inetmock.audit.v1.HTTPMethod
	(*HTTPHeaderValue)(nil),   // 1: inetmock.audit.v1.HTTPHeaderValue
	(*HTTPDetailsEntity)(nil), // 2: inetmock.audit.v1.HTTPDetailsEntity
	nil,                       // 3: inetmock.audit.v1.HTTPDetailsEntity.HeadersEntry
}
var file_audit_v1_http_details_proto_depIdxs = []int32{
	0, // 0: inetmock.audit.v1.HTTPDetailsEntity.method:type_name -> inetmock.audit.v1.HTTPMethod
	3, // 1: inetmock.audit.v1.HTTPDetailsEntity.headers:type_name -> inetmock.audit.v1.HTTPDetailsEntity.HeadersEntry
	1, // 2: inetmock.audit.v1.HTTPDetailsEntity.HeadersEntry.value:type_name -> inetmock.audit.v1.HTTPHeaderValue
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_audit_v1_http_details_proto_init() }
func file_audit_v1_http_details_proto_init() {
	if File_audit_v1_http_details_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_audit_v1_http_details_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HTTPHeaderValue); i {
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
		file_audit_v1_http_details_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HTTPDetailsEntity); i {
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
			RawDescriptor: file_audit_v1_http_details_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_audit_v1_http_details_proto_goTypes,
		DependencyIndexes: file_audit_v1_http_details_proto_depIdxs,
		EnumInfos:         file_audit_v1_http_details_proto_enumTypes,
		MessageInfos:      file_audit_v1_http_details_proto_msgTypes,
	}.Build()
	File_audit_v1_http_details_proto = out.File
	file_audit_v1_http_details_proto_rawDesc = nil
	file_audit_v1_http_details_proto_goTypes = nil
	file_audit_v1_http_details_proto_depIdxs = nil
}
