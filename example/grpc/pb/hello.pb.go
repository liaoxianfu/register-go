// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/hello.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Req struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Req) Reset()         { *m = Req{} }
func (m *Req) String() string { return proto.CompactTextString(m) }
func (*Req) ProtoMessage()    {}
func (*Req) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9418a89e4940c19, []int{0}
}

func (m *Req) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Req.Unmarshal(m, b)
}
func (m *Req) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Req.Marshal(b, m, deterministic)
}
func (m *Req) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Req.Merge(m, src)
}
func (m *Req) XXX_Size() int {
	return xxx_messageInfo_Req.Size(m)
}
func (m *Req) XXX_DiscardUnknown() {
	xxx_messageInfo_Req.DiscardUnknown(m)
}

var xxx_messageInfo_Req proto.InternalMessageInfo

func (m *Req) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type Resp struct {
	Info                 string   `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Resp) Reset()         { *m = Resp{} }
func (m *Resp) String() string { return proto.CompactTextString(m) }
func (*Resp) ProtoMessage()    {}
func (*Resp) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9418a89e4940c19, []int{1}
}

func (m *Resp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Resp.Unmarshal(m, b)
}
func (m *Resp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Resp.Marshal(b, m, deterministic)
}
func (m *Resp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Resp.Merge(m, src)
}
func (m *Resp) XXX_Size() int {
	return xxx_messageInfo_Resp.Size(m)
}
func (m *Resp) XXX_DiscardUnknown() {
	xxx_messageInfo_Resp.DiscardUnknown(m)
}

var xxx_messageInfo_Resp proto.InternalMessageInfo

func (m *Resp) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func init() {
	proto.RegisterType((*Req)(nil), "pb.req")
	proto.RegisterType((*Resp)(nil), "pb.resp")
}

func init() { proto.RegisterFile("pb/hello.proto", fileDescriptor_b9418a89e4940c19) }

var fileDescriptor_b9418a89e4940c19 = []byte{
	// 123 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x48, 0xd2, 0xcf,
	0x48, 0xcd, 0xc9, 0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x12,
	0xe7, 0x62, 0x2e, 0x4a, 0x2d, 0x14, 0x12, 0xe0, 0x62, 0xce, 0x2d, 0x4e, 0x97, 0x60, 0x54, 0x60,
	0xd4, 0xe0, 0x0c, 0x02, 0x31, 0x95, 0xa4, 0xb8, 0x58, 0x8a, 0x52, 0x8b, 0x0b, 0x84, 0x84, 0xb8,
	0x58, 0x32, 0xf3, 0xd2, 0xf2, 0xa1, 0x52, 0x60, 0xb6, 0x91, 0x1a, 0x17, 0x47, 0x70, 0x62, 0xa5,
	0x07, 0xc8, 0x28, 0x21, 0x29, 0x2e, 0x56, 0xb0, 0x99, 0x42, 0xec, 0x7a, 0x05, 0x49, 0x7a, 0x45,
	0xa9, 0x85, 0x52, 0x1c, 0x10, 0x46, 0x71, 0x41, 0x12, 0x1b, 0xd8, 0x1e, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x1a, 0xbc, 0xfa, 0x52, 0x79, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SayHelloClient is the client API for SayHello service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SayHelloClient interface {
	Hello(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error)
}

type sayHelloClient struct {
	cc *grpc.ClientConn
}

func NewSayHelloClient(cc *grpc.ClientConn) SayHelloClient {
	return &sayHelloClient{cc}
}

func (c *sayHelloClient) Hello(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error) {
	out := new(Resp)
	err := c.cc.Invoke(ctx, "/pb.SayHello/hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SayHelloServer is the server API for SayHello service.
type SayHelloServer interface {
	Hello(context.Context, *Req) (*Resp, error)
}

// UnimplementedSayHelloServer can be embedded to have forward compatible implementations.
type UnimplementedSayHelloServer struct {
}

func (*UnimplementedSayHelloServer) Hello(ctx context.Context, req *Req) (*Resp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}

func RegisterSayHelloServer(s *grpc.Server, srv SayHelloServer) {
	s.RegisterService(&_SayHello_serviceDesc, srv)
}

func _SayHello_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SayHelloServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.SayHello/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SayHelloServer).Hello(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

var _SayHello_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.SayHello",
	HandlerType: (*SayHelloServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "hello",
			Handler:    _SayHello_Hello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/hello.proto",
}
