// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//Message 用于服务间传输的消息，为了绑定Token，所以不直接传输data。服务名可直接通过data 获取，所以不再使用Key字段
type Message struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Data                 *any.Any `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	AimTokens            []string `protobuf:"bytes,3,rep,name=aimTokens,proto3" json:"aimTokens,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Message) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Message) GetAimTokens() []string {
	if m != nil {
		return m.AimTokens
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "Message")
}

func init() {
	proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd)
}

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 167 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0x92, 0x4c, 0xcf, 0xcf, 0x4f, 0xcf,
	0x49, 0xd5, 0x07, 0xf3, 0x92, 0x4a, 0xd3, 0xf4, 0x13, 0xf3, 0x2a, 0x21, 0x52, 0x4a, 0xe9, 0x5c,
	0xec, 0xbe, 0x10, 0xb5, 0x42, 0x22, 0x5c, 0xac, 0x25, 0xf9, 0xd9, 0xa9, 0x79, 0x12, 0x8c, 0x0a,
	0x8c, 0x1a, 0x9c, 0x41, 0x10, 0x8e, 0x90, 0x06, 0x17, 0x4b, 0x4a, 0x62, 0x49, 0xa2, 0x04, 0x93,
	0x02, 0xa3, 0x06, 0xb7, 0x91, 0x88, 0x1e, 0xc4, 0x28, 0x3d, 0x98, 0x51, 0x7a, 0x8e, 0x79, 0x95,
	0x41, 0x60, 0x15, 0x42, 0x32, 0x5c, 0x9c, 0x89, 0x99, 0xb9, 0x21, 0x20, 0x5d, 0xc5, 0x12, 0xcc,
	0x0a, 0xcc, 0x1a, 0x9c, 0x41, 0x08, 0x01, 0x27, 0xe9, 0x28, 0xc9, 0xf4, 0xcc, 0x92, 0x8c, 0xd2,
	0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0x8c, 0xcc, 0xc4, 0xbc, 0x6c, 0xfd, 0x92, 0x8c, 0xcc, 0xbc,
	0x6c, 0xfd, 0x82, 0xa4, 0x24, 0x36, 0xb0, 0x71, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x01,
	0x7c, 0x0b, 0x18, 0xb8, 0x00, 0x00, 0x00,
}
