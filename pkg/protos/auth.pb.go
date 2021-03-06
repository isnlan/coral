// 定义项目 API 的 proto 文件 可以同时描述 gRPC 和 HTTP API
// protobuf 文件参考:
//  - https://developers.google.com/protocol-buffers/

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: auth.proto

package protos

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Validate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	Bduid       string `protobuf:"bytes,2,opt,name=bduid,proto3" json:"bduid,omitempty"`
}

func (x *Validate) Reset() {
	*x = Validate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Validate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Validate) ProtoMessage() {}

func (x *Validate) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Validate.ProtoReflect.Descriptor instead.
func (*Validate) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{0}
}

func (x *Validate) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *Validate) GetBduid() string {
	if x != nil {
		return x.Bduid
	}
	return ""
}

var File_auth_proto protoreflect.FileDescriptor

var file_auth_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x43, 0x0a, 0x08, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x21,
	0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x64, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x62, 0x64, 0x75, 0x69, 0x64, 0x32, 0x41, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x43,
	0x65, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x1a,
	0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x44, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_proto_rawDescOnce sync.Once
	file_auth_proto_rawDescData = file_auth_proto_rawDesc
)

func file_auth_proto_rawDescGZIP() []byte {
	file_auth_proto_rawDescOnce.Do(func() {
		file_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_proto_rawDescData)
	})
	return file_auth_proto_rawDescData
}

var file_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_auth_proto_goTypes = []interface{}{
	(*Validate)(nil),        // 0: protos.Validate
	(*DigitalIdentity)(nil), // 1: protos.DigitalIdentity
}
var file_auth_proto_depIdxs = []int32{
	0, // 0: protos.UserCenter.Auth:input_type -> protos.Validate
	1, // 1: protos.UserCenter.Auth:output_type -> protos.DigitalIdentity
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_auth_proto_init() }
func file_auth_proto_init() {
	if File_auth_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Validate); i {
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
			RawDescriptor: file_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_proto_goTypes,
		DependencyIndexes: file_auth_proto_depIdxs,
		MessageInfos:      file_auth_proto_msgTypes,
	}.Build()
	File_auth_proto = out.File
	file_auth_proto_rawDesc = nil
	file_auth_proto_goTypes = nil
	file_auth_proto_depIdxs = nil
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

func (*UnimplementedUserCenterServer) Auth(context.Context, *Validate) (*DigitalIdentity, error) {
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
