// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: t.proto

package testdata

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

type Test1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Test1) Reset() {
	*x = Test1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_t_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Test1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test1) ProtoMessage() {}

func (x *Test1) ProtoReflect() protoreflect.Message {
	mi := &file_t_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Test1.ProtoReflect.Descriptor instead.
func (*Test1) Descriptor() ([]byte, []int) {
	return file_t_proto_rawDescGZIP(), []int{0}
}

func (x *Test1) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Test2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Age int32 `protobuf:"varint,1,opt,name=age,proto3" json:"age,omitempty"`
}

func (x *Test2) Reset() {
	*x = Test2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_t_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Test2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test2) ProtoMessage() {}

func (x *Test2) ProtoReflect() protoreflect.Message {
	mi := &file_t_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Test2.ProtoReflect.Descriptor instead.
func (*Test2) Descriptor() ([]byte, []int) {
	return file_t_proto_rawDescGZIP(), []int{1}
}

func (x *Test2) GetAge() int32 {
	if x != nil {
		return x.Age
	}
	return 0
}

var File_t_proto protoreflect.FileDescriptor

var file_t_proto_rawDesc = []byte{
	0x0a, 0x07, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1b, 0x0a, 0x05, 0x74, 0x65, 0x73,
	0x74, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x19, 0x0a, 0x05, 0x74, 0x65, 0x73, 0x74, 0x32, 0x12,
	0x10, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x61, 0x67,
	0x65, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x68, 0x69, 0x61, 0x6e, 0x6b, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x6b, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x2f, 0x64, 0x62, 0x2f, 0x64, 0x6f, 0x63, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_t_proto_rawDescOnce sync.Once
	file_t_proto_rawDescData = file_t_proto_rawDesc
)

func file_t_proto_rawDescGZIP() []byte {
	file_t_proto_rawDescOnce.Do(func() {
		file_t_proto_rawDescData = protoimpl.X.CompressGZIP(file_t_proto_rawDescData)
	})
	return file_t_proto_rawDescData
}

var file_t_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_t_proto_goTypes = []interface{}{
	(*Test1)(nil), // 0: test1
	(*Test2)(nil), // 1: test2
}
var file_t_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_t_proto_init() }
func file_t_proto_init() {
	if File_t_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_t_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Test1); i {
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
		file_t_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Test2); i {
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
			RawDescriptor: file_t_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_t_proto_goTypes,
		DependencyIndexes: file_t_proto_depIdxs,
		MessageInfos:      file_t_proto_msgTypes,
	}.Build()
	File_t_proto = out.File
	file_t_proto_rawDesc = nil
	file_t_proto_goTypes = nil
	file_t_proto_depIdxs = nil
}
