// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

package protos

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

type Validate struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	Bduid                string   `protobuf:"bytes,2,opt,name=bduid,proto3" json:"bduid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Validate) Reset()         { *m = Validate{} }
func (m *Validate) String() string { return proto.CompactTextString(m) }
func (*Validate) ProtoMessage()    {}
func (*Validate) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{0}
}

func (m *Validate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Validate.Unmarshal(m, b)
}
func (m *Validate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Validate.Marshal(b, m, deterministic)
}
func (m *Validate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Validate.Merge(m, src)
}
func (m *Validate) XXX_Size() int {
	return xxx_messageInfo_Validate.Size(m)
}
func (m *Validate) XXX_DiscardUnknown() {
	xxx_messageInfo_Validate.DiscardUnknown(m)
}

var xxx_messageInfo_Validate proto.InternalMessageInfo

func (m *Validate) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *Validate) GetBduid() string {
	if m != nil {
		return m.Bduid
	}
	return ""
}

func init() {
	proto.RegisterType((*Validate)(nil), "protos.Validate")
}

func init() {
	proto.RegisterFile("auth.proto", fileDescriptor_8bbd6f3875b0e874)
}

var fileDescriptor_8bbd6f3875b0e874 = []byte{
	// 161 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x2c, 0x2d, 0xc9,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0xc5, 0x52, 0x3c, 0xc9, 0xf9, 0xb9,
	0xb9, 0xf9, 0x79, 0x10, 0x51, 0x25, 0x67, 0x2e, 0x8e, 0xb0, 0xc4, 0x9c, 0xcc, 0x94, 0xc4, 0x92,
	0x54, 0x21, 0x45, 0x2e, 0x9e, 0xc4, 0xe4, 0xe4, 0xd4, 0xe2, 0xe2, 0xf8, 0x92, 0xfc, 0xec, 0xd4,
	0x3c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x6e, 0x88, 0x58, 0x08, 0x48, 0x48, 0x48, 0x84,
	0x8b, 0x35, 0x29, 0xa5, 0x34, 0x33, 0x45, 0x82, 0x09, 0x2c, 0x07, 0xe1, 0x18, 0x39, 0x72, 0x71,
	0x85, 0x16, 0xa7, 0x16, 0x39, 0xa7, 0xe6, 0x95, 0xa4, 0x16, 0x09, 0x19, 0x73, 0xb1, 0x38, 0x96,
	0x96, 0x64, 0x08, 0x09, 0x40, 0xac, 0x28, 0xd6, 0x83, 0x59, 0x20, 0x25, 0x0e, 0x13, 0x71, 0xc9,
	0x4c, 0xcf, 0x2c, 0x49, 0xcc, 0xf1, 0x4c, 0x49, 0xcd, 0x2b, 0xc9, 0x2c, 0xa9, 0x54, 0x62, 0x48,
	0x82, 0xb8, 0xce, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xe3, 0x9f, 0x24, 0x7f, 0xb2, 0x00, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UserCenterClient is the client API for UserCenter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserCenterClient interface {
	Auth(ctx context.Context, in *Validate, opts ...grpc.CallOption) (*DigitalIdentity, error)
}

type userCenterClient struct {
	cc grpc.ClientConnInterface
}

func NewUserCenterClient(cc grpc.ClientConnInterface) UserCenterClient {
	return &userCenterClient{cc}
}

func (c *userCenterClient) Auth(ctx context.Context, in *Validate, opts ...grpc.CallOption) (*DigitalIdentity, error) {
	out := new(DigitalIdentity)
	err := c.cc.Invoke(ctx, "/protos.UserCenter/Auth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserCenterServer is the server API for UserCenter service.
type UserCenterServer interface {
	Auth(context.Context, *Validate) (*DigitalIdentity, error)
}

// UnimplementedUserCenterServer can be embedded to have forward compatible implementations.
type UnimplementedUserCenterServer struct {
}

func (*UnimplementedUserCenterServer) Auth(ctx context.Context, req *Validate) (*DigitalIdentity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Auth not implemented")
}

func RegisterUserCenterServer(s *grpc.Server, srv UserCenterServer) {
	s.RegisterService(&_UserCenter_serviceDesc, srv)
}

func _UserCenter_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Validate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserCenterServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.UserCenter/Auth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserCenterServer).Auth(ctx, req.(*Validate))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserCenter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.UserCenter",
	HandlerType: (*UserCenterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Auth",
			Handler:    _UserCenter_Auth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}
