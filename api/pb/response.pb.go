// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/pb/response.proto

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

type Response struct {
	RequestID            string                    `protobuf:"bytes,1,opt,name=requestID" json:"requestID,omitempty"`
	Error                string                    `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
	DRKeyStoreGet        *Response_DRKeyStoreGet   `protobuf:"bytes,3,opt,name=dRKeyStoreGet" json:"dRKeyStoreGet,omitempty"`
	DRKeyStoreCount      *Response_DRKeyStoreCount `protobuf:"bytes,4,opt,name=dRKeyStoreCount" json:"dRKeyStoreCount,omitempty"`
	DRKeyStoreAll        *Response_DRKeyStoreAll   `protobuf:"bytes,5,opt,name=dRKeyStoreAll" json:"dRKeyStoreAll,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_response_c70b48d1fa0aa3de, []int{0}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetRequestID() string {
	if m != nil {
		return m.RequestID
	}
	return ""
}

func (m *Response) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *Response) GetDRKeyStoreGet() *Response_DRKeyStoreGet {
	if m != nil {
		return m.DRKeyStoreGet
	}
	return nil
}

func (m *Response) GetDRKeyStoreCount() *Response_DRKeyStoreCount {
	if m != nil {
		return m.DRKeyStoreCount
	}
	return nil
}

func (m *Response) GetDRKeyStoreAll() *Response_DRKeyStoreAll {
	if m != nil {
		return m.DRKeyStoreAll
	}
	return nil
}

type Response_DRKeyStoreGet struct {
	MessageKey           []byte   `protobuf:"bytes,1,opt,name=messageKey,proto3" json:"messageKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response_DRKeyStoreGet) Reset()         { *m = Response_DRKeyStoreGet{} }
func (m *Response_DRKeyStoreGet) String() string { return proto.CompactTextString(m) }
func (*Response_DRKeyStoreGet) ProtoMessage()    {}
func (*Response_DRKeyStoreGet) Descriptor() ([]byte, []int) {
	return fileDescriptor_response_c70b48d1fa0aa3de, []int{0, 0}
}
func (m *Response_DRKeyStoreGet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response_DRKeyStoreGet.Unmarshal(m, b)
}
func (m *Response_DRKeyStoreGet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response_DRKeyStoreGet.Marshal(b, m, deterministic)
}
func (dst *Response_DRKeyStoreGet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response_DRKeyStoreGet.Merge(dst, src)
}
func (m *Response_DRKeyStoreGet) XXX_Size() int {
	return xxx_messageInfo_Response_DRKeyStoreGet.Size(m)
}
func (m *Response_DRKeyStoreGet) XXX_DiscardUnknown() {
	xxx_messageInfo_Response_DRKeyStoreGet.DiscardUnknown(m)
}

var xxx_messageInfo_Response_DRKeyStoreGet proto.InternalMessageInfo

func (m *Response_DRKeyStoreGet) GetMessageKey() []byte {
	if m != nil {
		return m.MessageKey
	}
	return nil
}

type Response_DRKeyStoreCount struct {
	Count                uint64   `protobuf:"varint,1,opt,name=count" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response_DRKeyStoreCount) Reset()         { *m = Response_DRKeyStoreCount{} }
func (m *Response_DRKeyStoreCount) String() string { return proto.CompactTextString(m) }
func (*Response_DRKeyStoreCount) ProtoMessage()    {}
func (*Response_DRKeyStoreCount) Descriptor() ([]byte, []int) {
	return fileDescriptor_response_c70b48d1fa0aa3de, []int{0, 1}
}
func (m *Response_DRKeyStoreCount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response_DRKeyStoreCount.Unmarshal(m, b)
}
func (m *Response_DRKeyStoreCount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response_DRKeyStoreCount.Marshal(b, m, deterministic)
}
func (dst *Response_DRKeyStoreCount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response_DRKeyStoreCount.Merge(dst, src)
}
func (m *Response_DRKeyStoreCount) XXX_Size() int {
	return xxx_messageInfo_Response_DRKeyStoreCount.Size(m)
}
func (m *Response_DRKeyStoreCount) XXX_DiscardUnknown() {
	xxx_messageInfo_Response_DRKeyStoreCount.DiscardUnknown(m)
}

var xxx_messageInfo_Response_DRKeyStoreCount proto.InternalMessageInfo

func (m *Response_DRKeyStoreCount) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type Response_DRKeyStoreAll struct {
	All                  []*Response_DRKeyStoreAll_Key `protobuf:"bytes,1,rep,name=all" json:"all,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *Response_DRKeyStoreAll) Reset()         { *m = Response_DRKeyStoreAll{} }
func (m *Response_DRKeyStoreAll) String() string { return proto.CompactTextString(m) }
func (*Response_DRKeyStoreAll) ProtoMessage()    {}
func (*Response_DRKeyStoreAll) Descriptor() ([]byte, []int) {
	return fileDescriptor_response_c70b48d1fa0aa3de, []int{0, 2}
}
func (m *Response_DRKeyStoreAll) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response_DRKeyStoreAll.Unmarshal(m, b)
}
func (m *Response_DRKeyStoreAll) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response_DRKeyStoreAll.Marshal(b, m, deterministic)
}
func (dst *Response_DRKeyStoreAll) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response_DRKeyStoreAll.Merge(dst, src)
}
func (m *Response_DRKeyStoreAll) XXX_Size() int {
	return xxx_messageInfo_Response_DRKeyStoreAll.Size(m)
}
func (m *Response_DRKeyStoreAll) XXX_DiscardUnknown() {
	xxx_messageInfo_Response_DRKeyStoreAll.DiscardUnknown(m)
}

var xxx_messageInfo_Response_DRKeyStoreAll proto.InternalMessageInfo

func (m *Response_DRKeyStoreAll) GetAll() []*Response_DRKeyStoreAll_Key {
	if m != nil {
		return m.All
	}
	return nil
}

type Response_DRKeyStoreAll_Key struct {
	Key                  []byte            `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	MessageKeys          map[uint64][]byte `protobuf:"bytes,2,rep,name=messageKeys" json:"messageKeys,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Response_DRKeyStoreAll_Key) Reset()         { *m = Response_DRKeyStoreAll_Key{} }
func (m *Response_DRKeyStoreAll_Key) String() string { return proto.CompactTextString(m) }
func (*Response_DRKeyStoreAll_Key) ProtoMessage()    {}
func (*Response_DRKeyStoreAll_Key) Descriptor() ([]byte, []int) {
	return fileDescriptor_response_c70b48d1fa0aa3de, []int{0, 2, 0}
}
func (m *Response_DRKeyStoreAll_Key) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response_DRKeyStoreAll_Key.Unmarshal(m, b)
}
func (m *Response_DRKeyStoreAll_Key) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response_DRKeyStoreAll_Key.Marshal(b, m, deterministic)
}
func (dst *Response_DRKeyStoreAll_Key) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response_DRKeyStoreAll_Key.Merge(dst, src)
}
func (m *Response_DRKeyStoreAll_Key) XXX_Size() int {
	return xxx_messageInfo_Response_DRKeyStoreAll_Key.Size(m)
}
func (m *Response_DRKeyStoreAll_Key) XXX_DiscardUnknown() {
	xxx_messageInfo_Response_DRKeyStoreAll_Key.DiscardUnknown(m)
}

var xxx_messageInfo_Response_DRKeyStoreAll_Key proto.InternalMessageInfo

func (m *Response_DRKeyStoreAll_Key) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Response_DRKeyStoreAll_Key) GetMessageKeys() map[uint64][]byte {
	if m != nil {
		return m.MessageKeys
	}
	return nil
}

func init() {
	proto.RegisterType((*Response)(nil), "api.Response")
	proto.RegisterType((*Response_DRKeyStoreGet)(nil), "api.Response.DRKeyStoreGet")
	proto.RegisterType((*Response_DRKeyStoreCount)(nil), "api.Response.DRKeyStoreCount")
	proto.RegisterType((*Response_DRKeyStoreAll)(nil), "api.Response.DRKeyStoreAll")
	proto.RegisterType((*Response_DRKeyStoreAll_Key)(nil), "api.Response.DRKeyStoreAll.Key")
	proto.RegisterMapType((map[uint64][]byte)(nil), "api.Response.DRKeyStoreAll.Key.MessageKeysEntry")
}

func init() { proto.RegisterFile("api/pb/response.proto", fileDescriptor_response_c70b48d1fa0aa3de) }

var fileDescriptor_response_c70b48d1fa0aa3de = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcd, 0x4a, 0xfb, 0x40,
	0x14, 0xc5, 0x99, 0x26, 0xfd, 0xf3, 0xef, 0x6d, 0x4b, 0xcb, 0x50, 0x61, 0x88, 0x5f, 0xc5, 0x8d,
	0x5d, 0xa5, 0x5a, 0x37, 0xe2, 0x42, 0x08, 0x56, 0x8a, 0x04, 0x37, 0xd7, 0x27, 0x48, 0xf5, 0x22,
	0xc5, 0xb1, 0x19, 0x67, 0xa6, 0x42, 0x9e, 0xc9, 0x37, 0xf2, 0x29, 0x7c, 0x04, 0xc9, 0x44, 0xcd,
	0x07, 0x58, 0x77, 0xb9, 0x27, 0xe7, 0x9c, 0xf9, 0xcd, 0x70, 0x61, 0x27, 0x51, 0xab, 0xa9, 0x5a,
	0x4e, 0x35, 0x19, 0x95, 0xae, 0x0d, 0x85, 0x4a, 0xa7, 0x36, 0xe5, 0x5e, 0xa2, 0x56, 0x47, 0xef,
	0x3e, 0xfc, 0xc7, 0x2f, 0x9d, 0xef, 0x41, 0x47, 0xd3, 0xcb, 0x86, 0x8c, 0xbd, 0x99, 0x0b, 0x36,
	0x66, 0x93, 0x0e, 0x96, 0x02, 0x1f, 0x41, 0x9b, 0xb4, 0x4e, 0xb5, 0x68, 0xb9, 0x3f, 0xc5, 0xc0,
	0x23, 0xe8, 0x3f, 0x60, 0x4c, 0xd9, 0x9d, 0x4d, 0x35, 0x2d, 0xc8, 0x0a, 0x6f, 0xcc, 0x26, 0xdd,
	0xd9, 0x6e, 0x98, 0xa8, 0x55, 0xf8, 0xdd, 0x1c, 0xce, 0xab, 0x16, 0xac, 0x27, 0xf8, 0x02, 0x06,
	0xa5, 0x70, 0x95, 0x6e, 0xd6, 0x56, 0xf8, 0xae, 0x64, 0xff, 0xb7, 0x12, 0x67, 0xc2, 0x66, 0xaa,
	0xce, 0x12, 0x49, 0x29, 0xda, 0xdb, 0x59, 0x22, 0x29, 0xb1, 0x9e, 0x08, 0xa6, 0xd0, 0xaf, 0xb1,
	0xf2, 0x03, 0x80, 0x67, 0x32, 0x26, 0x79, 0xa4, 0x98, 0x32, 0xf7, 0x28, 0x3d, 0xac, 0x28, 0xc1,
	0x31, 0x0c, 0x1a, 0x5c, 0xf9, 0x43, 0xdd, 0xbb, 0x5b, 0xe4, 0x6e, 0x1f, 0x8b, 0x21, 0xf8, 0x60,
	0xd5, 0xea, 0x48, 0x4a, 0x7e, 0x0a, 0x5e, 0x22, 0xa5, 0x60, 0x63, 0x6f, 0xd2, 0x9d, 0x1d, 0x6e,
	0x81, 0x0c, 0x63, 0xca, 0x30, 0xf7, 0x06, 0x6f, 0x0c, 0xbc, 0x98, 0x32, 0x3e, 0x04, 0xef, 0xe9,
	0x07, 0x27, 0xff, 0xe4, 0x08, 0xdd, 0x92, 0xca, 0x88, 0x96, 0x2b, 0x3d, 0xf9, 0xa3, 0x34, 0xbc,
	0x2d, 0x23, 0xd7, 0x6b, 0xab, 0x33, 0xac, 0x96, 0x04, 0x97, 0x30, 0x6c, 0x1a, 0xaa, 0x27, 0xfb,
	0xc5, 0xc9, 0x23, 0x68, 0xbf, 0x26, 0x72, 0x43, 0x6e, 0x2f, 0x7a, 0x58, 0x0c, 0x17, 0xad, 0x73,
	0xb6, 0xfc, 0xe7, 0x16, 0xed, 0xec, 0x33, 0x00, 0x00, 0xff, 0xff, 0x05, 0x18, 0x2d, 0xf6, 0x81,
	0x02, 0x00, 0x00,
}
