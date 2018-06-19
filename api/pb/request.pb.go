// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/pb/request.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	RequestID            string                 `protobuf:"bytes,1,opt,name=requestID" json:"requestID,omitempty"`
	DRKeyStoreGet        *Request_DRKeyStoreGet `protobuf:"bytes,2,opt,name=dRKeyStoreGet" json:"dRKeyStoreGet,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_request_2127127062921383, []int{0}
}
func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (dst *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(dst, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetRequestID() string {
	if m != nil {
		return m.RequestID
	}
	return ""
}

func (m *Request) GetDRKeyStoreGet() *Request_DRKeyStoreGet {
	if m != nil {
		return m.DRKeyStoreGet
	}
	return nil
}

type Request_DRKeyStoreGet struct {
	DrKey                []byte   `protobuf:"bytes,1,opt,name=drKey,proto3" json:"drKey,omitempty"`
	MessageNumber        uint64   `protobuf:"varint,2,opt,name=messageNumber" json:"messageNumber,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request_DRKeyStoreGet) Reset()         { *m = Request_DRKeyStoreGet{} }
func (m *Request_DRKeyStoreGet) String() string { return proto.CompactTextString(m) }
func (*Request_DRKeyStoreGet) ProtoMessage()    {}
func (*Request_DRKeyStoreGet) Descriptor() ([]byte, []int) {
	return fileDescriptor_request_2127127062921383, []int{0, 0}
}
func (m *Request_DRKeyStoreGet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request_DRKeyStoreGet.Unmarshal(m, b)
}
func (m *Request_DRKeyStoreGet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request_DRKeyStoreGet.Marshal(b, m, deterministic)
}
func (dst *Request_DRKeyStoreGet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request_DRKeyStoreGet.Merge(dst, src)
}
func (m *Request_DRKeyStoreGet) XXX_Size() int {
	return xxx_messageInfo_Request_DRKeyStoreGet.Size(m)
}
func (m *Request_DRKeyStoreGet) XXX_DiscardUnknown() {
	xxx_messageInfo_Request_DRKeyStoreGet.DiscardUnknown(m)
}

var xxx_messageInfo_Request_DRKeyStoreGet proto.InternalMessageInfo

func (m *Request_DRKeyStoreGet) GetDrKey() []byte {
	if m != nil {
		return m.DrKey
	}
	return nil
}

func (m *Request_DRKeyStoreGet) GetMessageNumber() uint64 {
	if m != nil {
		return m.MessageNumber
	}
	return 0
}

func init() {
	proto.RegisterType((*Request)(nil), "api.Request")
	proto.RegisterType((*Request_DRKeyStoreGet)(nil), "api.Request.DRKeyStoreGet")
}

func init() { proto.RegisterFile("api/pb/request.proto", fileDescriptor_request_2127127062921383) }

var fileDescriptor_request_2127127062921383 = []byte{
	// 165 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x49, 0x2c, 0xc8, 0xd4,
	0x2f, 0x48, 0xd2, 0x2f, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x4e, 0x2c, 0xc8, 0x54, 0xda, 0xc6, 0xc8, 0xc5, 0x1e, 0x04, 0x11, 0x16, 0x92, 0xe1,
	0xe2, 0x84, 0xaa, 0xf0, 0x74, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x42, 0x08, 0x08, 0x39,
	0x70, 0xf1, 0xa6, 0x04, 0x79, 0xa7, 0x56, 0x06, 0x97, 0xe4, 0x17, 0xa5, 0xba, 0xa7, 0x96, 0x48,
	0x30, 0x29, 0x30, 0x6a, 0x70, 0x1b, 0x49, 0xe9, 0x25, 0x16, 0x64, 0xea, 0x41, 0x8d, 0xd0, 0x73,
	0x41, 0x56, 0x11, 0x84, 0xaa, 0x41, 0xca, 0x9b, 0x8b, 0x17, 0x45, 0x5e, 0x48, 0x84, 0x8b, 0x35,
	0xa5, 0xc8, 0x3b, 0xb5, 0x12, 0x6c, 0x19, 0x4f, 0x10, 0x84, 0x23, 0xa4, 0xc2, 0xc5, 0x9b, 0x9b,
	0x5a, 0x5c, 0x9c, 0x98, 0x9e, 0xea, 0x57, 0x9a, 0x9b, 0x94, 0x5a, 0x04, 0xb6, 0x88, 0x25, 0x08,
	0x55, 0x30, 0x89, 0x0d, 0xec, 0x09, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x25, 0x45, 0x47,
	0x0c, 0xdc, 0x00, 0x00, 0x00,
}
