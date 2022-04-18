// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: pbtest.proto

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
		mi := &file_pbtest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Test1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test1) ProtoMessage() {}

func (x *Test1) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[0]
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
	return file_pbtest_proto_rawDescGZIP(), []int{0}
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

	Hope string `protobuf:"bytes,1,opt,name=hope,proto3" json:"hope,omitempty"`
}

func (x *Test2) Reset() {
	*x = Test2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbtest_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Test2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test2) ProtoMessage() {}

func (x *Test2) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[1]
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
	return file_pbtest_proto_rawDescGZIP(), []int{1}
}

func (x *Test2) GetHope() string {
	if x != nil {
		return x.Hope
	}
	return ""
}

type AnyTest1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *AnyTest1) Reset() {
	*x = AnyTest1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbtest_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnyTest1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnyTest1) ProtoMessage() {}

func (x *AnyTest1) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnyTest1.ProtoReflect.Descriptor instead.
func (*AnyTest1) Descriptor() ([]byte, []int) {
	return file_pbtest_proto_rawDescGZIP(), []int{2}
}

func (x *AnyTest1) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type AnyTest2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hope string `protobuf:"bytes,1,opt,name=hope,proto3" json:"hope,omitempty"`
}

func (x *AnyTest2) Reset() {
	*x = AnyTest2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbtest_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnyTest2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnyTest2) ProtoMessage() {}

func (x *AnyTest2) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnyTest2.ProtoReflect.Descriptor instead.
func (*AnyTest2) Descriptor() ([]byte, []int) {
	return file_pbtest_proto_rawDescGZIP(), []int{3}
}

func (x *AnyTest2) GetHope() string {
	if x != nil {
		return x.Hope
	}
	return ""
}

type MessageTest1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *MessageTest1) Reset() {
	*x = MessageTest1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbtest_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageTest1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageTest1) ProtoMessage() {}

func (x *MessageTest1) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageTest1.ProtoReflect.Descriptor instead.
func (*MessageTest1) Descriptor() ([]byte, []int) {
	return file_pbtest_proto_rawDescGZIP(), []int{4}
}

func (x *MessageTest1) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type S_Example struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *S_Example) Reset() {
	*x = S_Example{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbtest_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S_Example) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S_Example) ProtoMessage() {}

func (x *S_Example) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S_Example.ProtoReflect.Descriptor instead.
func (*S_Example) Descriptor() ([]byte, []int) {
	return file_pbtest_proto_rawDescGZIP(), []int{5}
}

func (x *S_Example) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type G_Example struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *G_Example) Reset() {
	*x = G_Example{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbtest_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *G_Example) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*G_Example) ProtoMessage() {}

func (x *G_Example) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use G_Example.ProtoReflect.Descriptor instead.
func (*G_Example) Descriptor() ([]byte, []int) {
	return file_pbtest_proto_rawDescGZIP(), []int{6}
}

func (x *G_Example) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type P_Example struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *P_Example) Reset() {
	*x = P_Example{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbtest_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *P_Example) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*P_Example) ProtoMessage() {}

func (x *P_Example) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use P_Example.ProtoReflect.Descriptor instead.
func (*P_Example) Descriptor() ([]byte, []int) {
	return file_pbtest_proto_rawDescGZIP(), []int{7}
}

func (x *P_Example) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type TEST_Example struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *TEST_Example) Reset() {
	*x = TEST_Example{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbtest_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TEST_Example) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TEST_Example) ProtoMessage() {}

func (x *TEST_Example) ProtoReflect() protoreflect.Message {
	mi := &file_pbtest_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TEST_Example.ProtoReflect.Descriptor instead.
func (*TEST_Example) Descriptor() ([]byte, []int) {
	return file_pbtest_proto_rawDescGZIP(), []int{8}
}

func (x *TEST_Example) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_pbtest_proto protoreflect.FileDescriptor

var file_pbtest_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x62, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1b,
	0x0a, 0x05, 0x74, 0x65, 0x73, 0x74, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x1b, 0x0a, 0x05, 0x74,
	0x65, 0x73, 0x74, 0x32, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x70, 0x65, 0x22, 0x1e, 0x0a, 0x08, 0x41, 0x6e, 0x79, 0x54,
	0x65, 0x73, 0x74, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x1e, 0x0a, 0x08, 0x61, 0x6e, 0x79, 0x54,
	0x65, 0x73, 0x74, 0x32, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x70, 0x65, 0x22, 0x20, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x54, 0x65, 0x73, 0x74, 0x31, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x21, 0x0a, 0x09, 0x53, 0x5f,
	0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x21, 0x0a,
	0x09, 0x47, 0x5f, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x21, 0x0a, 0x09, 0x50, 0x5f, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x24, 0x0a, 0x0c, 0x54, 0x45, 0x53, 0x54, 0x5f, 0x45, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x30, 0x5a, 0x23, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x69, 0x61, 0x6e, 0x6b, 0x2f, 0x74, 0x68,
	0x69, 0x6e, 0x6b, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61,
	0xaa, 0x02, 0x08, 0x48, 0x75, 0x62, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pbtest_proto_rawDescOnce sync.Once
	file_pbtest_proto_rawDescData = file_pbtest_proto_rawDesc
)

func file_pbtest_proto_rawDescGZIP() []byte {
	file_pbtest_proto_rawDescOnce.Do(func() {
		file_pbtest_proto_rawDescData = protoimpl.X.CompressGZIP(file_pbtest_proto_rawDescData)
	})
	return file_pbtest_proto_rawDescData
}

var file_pbtest_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_pbtest_proto_goTypes = []interface{}{
	(*Test1)(nil),        // 0: test1
	(*Test2)(nil),        // 1: test2
	(*AnyTest1)(nil),     // 2: AnyTest1
	(*AnyTest2)(nil),     // 3: anyTest2
	(*MessageTest1)(nil), // 4: messageTest1
	(*S_Example)(nil),    // 5: S_Example
	(*G_Example)(nil),    // 6: G_Example
	(*P_Example)(nil),    // 7: P_Example
	(*TEST_Example)(nil), // 8: TEST_Example
}
var file_pbtest_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pbtest_proto_init() }
func file_pbtest_proto_init() {
	if File_pbtest_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pbtest_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_pbtest_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_pbtest_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnyTest1); i {
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
		file_pbtest_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnyTest2); i {
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
		file_pbtest_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageTest1); i {
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
		file_pbtest_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S_Example); i {
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
		file_pbtest_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*G_Example); i {
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
		file_pbtest_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*P_Example); i {
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
		file_pbtest_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TEST_Example); i {
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
			RawDescriptor: file_pbtest_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pbtest_proto_goTypes,
		DependencyIndexes: file_pbtest_proto_depIdxs,
		MessageInfos:      file_pbtest_proto_msgTypes,
	}.Build()
	File_pbtest_proto = out.File
	file_pbtest_proto_rawDesc = nil
	file_pbtest_proto_goTypes = nil
	file_pbtest_proto_depIdxs = nil
}
