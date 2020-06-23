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

type FabricDigital struct {
	// 数字身份唯一标识
	Bduid string `protobuf:"bytes,1,opt,name=bduid,proto3" json:"bduid,omitempty"`
	// 区块链网络ID
	NetworkId string `protobuf:"bytes,2,opt,name=network_id,json=networkId,proto3" json:"network_id,omitempty"`
	// 加密后的私钥
	PrivateKey string `protobuf:"bytes,3,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
	// 证书
	Cert string `protobuf:"bytes,4,opt,name=cert,proto3" json:"cert,omitempty"`
	// msp id
	MspId string `protobuf:"bytes,5,opt,name=msp_id,json=mspId,proto3" json:"msp_id,omitempty"`
	// 激活状态
	Status bool `protobuf:"varint,6,opt,name=status,proto3" json:"status,omitempty"`
	// 是否是组织
	IsOrg                bool     `protobuf:"varint,7,opt,name=is_org,json=isOrg,proto3" json:"is_org,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FabricDigital) Reset()         { *m = FabricDigital{} }
func (m *FabricDigital) String() string { return proto.CompactTextString(m) }
func (*FabricDigital) ProtoMessage()    {}
func (*FabricDigital) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{1}
}

func (m *FabricDigital) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FabricDigital.Unmarshal(m, b)
}
func (m *FabricDigital) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FabricDigital.Marshal(b, m, deterministic)
}
func (m *FabricDigital) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FabricDigital.Merge(m, src)
}
func (m *FabricDigital) XXX_Size() int {
	return xxx_messageInfo_FabricDigital.Size(m)
}
func (m *FabricDigital) XXX_DiscardUnknown() {
	xxx_messageInfo_FabricDigital.DiscardUnknown(m)
}

var xxx_messageInfo_FabricDigital proto.InternalMessageInfo

func (m *FabricDigital) GetBduid() string {
	if m != nil {
		return m.Bduid
	}
	return ""
}

func (m *FabricDigital) GetNetworkId() string {
	if m != nil {
		return m.NetworkId
	}
	return ""
}

func (m *FabricDigital) GetPrivateKey() string {
	if m != nil {
		return m.PrivateKey
	}
	return ""
}

func (m *FabricDigital) GetCert() string {
	if m != nil {
		return m.Cert
	}
	return ""
}

func (m *FabricDigital) GetMspId() string {
	if m != nil {
		return m.MspId
	}
	return ""
}

func (m *FabricDigital) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *FabricDigital) GetIsOrg() bool {
	if m != nil {
		return m.IsOrg
	}
	return false
}

type KBaasDigital struct {
	// 数字身份唯一标识
	Bduid string `protobuf:"bytes,1,opt,name=bduid,proto3" json:"bduid,omitempty"`
	// 区块链网络ID
	NetworkId string `protobuf:"bytes,2,opt,name=network_id,json=networkId,proto3" json:"network_id,omitempty"`
	// 加密后的私钥
	PrivateKey string `protobuf:"bytes,3,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
	// 公钥
	PublicKey string `protobuf:"bytes,4,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	// 证书
	Cert string `protobuf:"bytes,5,opt,name=cert,proto3" json:"cert,omitempty"`
	// 激活状态
	Status               bool     `protobuf:"varint,6,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KBaasDigital) Reset()         { *m = KBaasDigital{} }
func (m *KBaasDigital) String() string { return proto.CompactTextString(m) }
func (*KBaasDigital) ProtoMessage()    {}
func (*KBaasDigital) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{2}
}

func (m *KBaasDigital) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KBaasDigital.Unmarshal(m, b)
}
func (m *KBaasDigital) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KBaasDigital.Marshal(b, m, deterministic)
}
func (m *KBaasDigital) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KBaasDigital.Merge(m, src)
}
func (m *KBaasDigital) XXX_Size() int {
	return xxx_messageInfo_KBaasDigital.Size(m)
}
func (m *KBaasDigital) XXX_DiscardUnknown() {
	xxx_messageInfo_KBaasDigital.DiscardUnknown(m)
}

var xxx_messageInfo_KBaasDigital proto.InternalMessageInfo

func (m *KBaasDigital) GetBduid() string {
	if m != nil {
		return m.Bduid
	}
	return ""
}

func (m *KBaasDigital) GetNetworkId() string {
	if m != nil {
		return m.NetworkId
	}
	return ""
}

func (m *KBaasDigital) GetPrivateKey() string {
	if m != nil {
		return m.PrivateKey
	}
	return ""
}

func (m *KBaasDigital) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

func (m *KBaasDigital) GetCert() string {
	if m != nil {
		return m.Cert
	}
	return ""
}

func (m *KBaasDigital) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

type DigitalIdentity struct {
	// 数字身份唯一标识
	Bduid string `protobuf:"bytes,1,opt,name=bduid,proto3" json:"bduid,omitempty"`
	// 服务中心的数字身份地址
	Proxy string `protobuf:"bytes,2,opt,name=proxy,proto3" json:"proxy,omitempty"`
	// 访问的区块链网络ID
	NetworkId string `protobuf:"bytes,3,opt,name=network_id,json=networkId,proto3" json:"network_id,omitempty"`
	// 数字身份
	//
	// Types that are valid to be assigned to Digital:
	//	*DigitalIdentity_Fabric
	//	*DigitalIdentity_Kbaas
	Digital              isDigitalIdentity_Digital `protobuf_oneof:"digital"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *DigitalIdentity) Reset()         { *m = DigitalIdentity{} }
func (m *DigitalIdentity) String() string { return proto.CompactTextString(m) }
func (*DigitalIdentity) ProtoMessage()    {}
func (*DigitalIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{3}
}

func (m *DigitalIdentity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DigitalIdentity.Unmarshal(m, b)
}
func (m *DigitalIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DigitalIdentity.Marshal(b, m, deterministic)
}
func (m *DigitalIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DigitalIdentity.Merge(m, src)
}
func (m *DigitalIdentity) XXX_Size() int {
	return xxx_messageInfo_DigitalIdentity.Size(m)
}
func (m *DigitalIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_DigitalIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_DigitalIdentity proto.InternalMessageInfo

func (m *DigitalIdentity) GetBduid() string {
	if m != nil {
		return m.Bduid
	}
	return ""
}

func (m *DigitalIdentity) GetProxy() string {
	if m != nil {
		return m.Proxy
	}
	return ""
}

func (m *DigitalIdentity) GetNetworkId() string {
	if m != nil {
		return m.NetworkId
	}
	return ""
}

type isDigitalIdentity_Digital interface {
	isDigitalIdentity_Digital()
}

type DigitalIdentity_Fabric struct {
	Fabric *FabricDigital `protobuf:"bytes,4,opt,name=fabric,proto3,oneof"`
}

type DigitalIdentity_Kbaas struct {
	Kbaas *KBaasDigital `protobuf:"bytes,5,opt,name=kbaas,proto3,oneof"`
}

func (*DigitalIdentity_Fabric) isDigitalIdentity_Digital() {}

func (*DigitalIdentity_Kbaas) isDigitalIdentity_Digital() {}

func (m *DigitalIdentity) GetDigital() isDigitalIdentity_Digital {
	if m != nil {
		return m.Digital
	}
	return nil
}

func (m *DigitalIdentity) GetFabric() *FabricDigital {
	if x, ok := m.GetDigital().(*DigitalIdentity_Fabric); ok {
		return x.Fabric
	}
	return nil
}

func (m *DigitalIdentity) GetKbaas() *KBaasDigital {
	if x, ok := m.GetDigital().(*DigitalIdentity_Kbaas); ok {
		return x.Kbaas
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*DigitalIdentity) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*DigitalIdentity_Fabric)(nil),
		(*DigitalIdentity_Kbaas)(nil),
	}
}

func init() {
	proto.RegisterType((*Validate)(nil), "protos.Validate")
	proto.RegisterType((*FabricDigital)(nil), "protos.FabricDigital")
	proto.RegisterType((*KBaasDigital)(nil), "protos.KBaasDigital")
	proto.RegisterType((*DigitalIdentity)(nil), "protos.DigitalIdentity")
}

func init() {
	proto.RegisterFile("auth.proto", fileDescriptor_8bbd6f3875b0e874)
}

var fileDescriptor_8bbd6f3875b0e874 = []byte{
	// 376 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x92, 0xc1, 0x6e, 0xaa, 0x40,
	0x14, 0x86, 0xe5, 0x2a, 0x28, 0x07, 0x6f, 0xee, 0xcd, 0x44, 0xef, 0x25, 0x4d, 0x4c, 0x2d, 0x2b,
	0x17, 0x8d, 0x4d, 0xf4, 0x09, 0xd4, 0xa6, 0xd1, 0xb8, 0x68, 0x42, 0xda, 0x6e, 0xc9, 0xc0, 0x4c,
	0x75, 0x82, 0x02, 0x99, 0x19, 0xda, 0xf2, 0x54, 0x7d, 0x84, 0xee, 0xfa, 0x5c, 0x8d, 0x33, 0xd0,
	0x52, 0x13, 0x97, 0x5d, 0xc1, 0xf9, 0xcf, 0xcf, 0xf0, 0x7f, 0x73, 0x0e, 0x00, 0xce, 0xe5, 0x76,
	0x9c, 0xf1, 0x54, 0xa6, 0xc8, 0x52, 0x0f, 0xe1, 0x2d, 0xa0, 0xf3, 0x80, 0x77, 0x8c, 0x60, 0x49,
	0xd1, 0x05, 0x74, 0x71, 0x14, 0x51, 0x21, 0x02, 0x99, 0xc6, 0x34, 0x71, 0x8d, 0xa1, 0x31, 0xb2,
	0x7d, 0x47, 0x6b, 0x77, 0x07, 0x09, 0xf5, 0xc0, 0x0c, 0x49, 0xce, 0x88, 0xfb, 0x4b, 0xf5, 0x74,
	0xe1, 0xbd, 0x19, 0xf0, 0xfb, 0x06, 0x87, 0x9c, 0x45, 0xd7, 0x6c, 0xc3, 0x24, 0xde, 0x7d, 0xf9,
	0x8c, 0x9a, 0x0f, 0x0d, 0x00, 0x12, 0x2a, 0x9f, 0x53, 0x1e, 0x07, 0x9f, 0x47, 0xd8, 0xa5, 0xb2,
	0x22, 0xe8, 0x1c, 0x9c, 0x8c, 0xb3, 0x27, 0x2c, 0x69, 0x10, 0xd3, 0xc2, 0x6d, 0xaa, 0x3e, 0x94,
	0xd2, 0x9a, 0x16, 0x08, 0x41, 0x2b, 0xa2, 0x5c, 0xba, 0x2d, 0xd5, 0x51, 0xef, 0xa8, 0x0f, 0xd6,
	0x5e, 0x64, 0x87, 0xf3, 0x4c, 0xfd, 0xab, 0xbd, 0xc8, 0x56, 0x04, 0xfd, 0x03, 0x4b, 0x48, 0x2c,
	0x73, 0xe1, 0x5a, 0x43, 0x63, 0xd4, 0xf1, 0xcb, 0xea, 0x60, 0x67, 0x22, 0x48, 0xf9, 0xc6, 0x6d,
	0x2b, 0xdd, 0x64, 0xe2, 0x96, 0x6f, 0xbc, 0x57, 0x03, 0xba, 0xeb, 0x39, 0xc6, 0xe2, 0x67, 0x01,
	0x06, 0x00, 0x59, 0x1e, 0xee, 0x58, 0xa4, 0xfa, 0x1a, 0xc3, 0xd6, 0x4a, 0x9d, 0xcf, 0xac, 0xf1,
	0x9d, 0x00, 0xf1, 0xde, 0x0d, 0xf8, 0x53, 0x86, 0x5d, 0x11, 0x9a, 0x48, 0x26, 0x8b, 0x13, 0xa1,
	0x7b, 0x60, 0x66, 0x3c, 0x7d, 0x29, 0xaa, 0x99, 0xa9, 0xe2, 0x08, 0xa5, 0x79, 0x8c, 0x72, 0x05,
	0xd6, 0xa3, 0x9a, 0xa8, 0x4a, 0xe9, 0x4c, 0xfa, 0x7a, 0x6f, 0xc4, 0xf8, 0xdb, 0x9c, 0x97, 0x0d,
	0xbf, 0xb4, 0xa1, 0x4b, 0x30, 0xe3, 0x10, 0x63, 0xa1, 0xc2, 0x3b, 0x93, 0x5e, 0xe5, 0xaf, 0xdf,
	0xea, 0xb2, 0xe1, 0x6b, 0xd3, 0xdc, 0x86, 0x36, 0xd1, 0xda, 0x64, 0x06, 0x70, 0x2f, 0x28, 0x5f,
	0xd0, 0x44, 0x52, 0x8e, 0xa6, 0xd0, 0x9a, 0xe5, 0x72, 0x8b, 0xfe, 0x56, 0xdf, 0x57, 0xdb, 0x79,
	0xf6, 0xbf, 0x52, 0x8e, 0xa8, 0xbd, 0x46, 0xa8, 0x97, 0x79, 0xfa, 0x11, 0x00, 0x00, 0xff, 0xff,
	0xf9, 0xe1, 0x77, 0xed, 0xe1, 0x02, 0x00, 0x00,
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
