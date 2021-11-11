// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: profile.proto

package proto

import (
	context "context"
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

type UserSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email       string `protobuf:"bytes,1,opt,name=Email,proto3" json:"Email,omitempty"`
	Nickname    string `protobuf:"bytes,2,opt,name=Nickname,proto3" json:"Nickname,omitempty"`
	SmallAvatar string `protobuf:"bytes,3,opt,name=SmallAvatar,proto3" json:"SmallAvatar,omitempty"`
	BigAvatar   string `protobuf:"bytes,4,opt,name=BigAvatar,proto3" json:"BigAvatar,omitempty"`
}

func (x *UserSettings) Reset() {
	*x = UserSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSettings) ProtoMessage() {}

func (x *UserSettings) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSettings.ProtoReflect.Descriptor instead.
func (*UserSettings) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{0}
}

func (x *UserSettings) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserSettings) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *UserSettings) GetSmallAvatar() string {
	if x != nil {
		return x.SmallAvatar
	}
	return ""
}

func (x *UserSettings) GetBigAvatar() string {
	if x != nil {
		return x.BigAvatar
	}
	return ""
}

type UploadSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID         int64         `protobuf:"varint,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Email          string        `protobuf:"bytes,2,opt,name=Email,proto3" json:"Email,omitempty"`
	Nickname       string        `protobuf:"bytes,3,opt,name=Nickname,proto3" json:"Nickname,omitempty"`
	AvatarFilename string        `protobuf:"bytes,4,opt,name=AvatarFilename,proto3" json:"AvatarFilename,omitempty"`
	OldPassword    string        `protobuf:"bytes,5,opt,name=OldPassword,proto3" json:"OldPassword,omitempty"`
	NewPassword    string        `protobuf:"bytes,6,opt,name=NewPassword,proto3" json:"NewPassword,omitempty"`
	OldSettings    *UserSettings `protobuf:"bytes,7,opt,name=OldSettings,proto3" json:"OldSettings,omitempty"`
}

func (x *UploadSettings) Reset() {
	*x = UploadSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadSettings) ProtoMessage() {}

func (x *UploadSettings) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadSettings.ProtoReflect.Descriptor instead.
func (*UploadSettings) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{1}
}

func (x *UploadSettings) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *UploadSettings) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UploadSettings) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *UploadSettings) GetAvatarFilename() string {
	if x != nil {
		return x.AvatarFilename
	}
	return ""
}

func (x *UploadSettings) GetOldPassword() string {
	if x != nil {
		return x.OldPassword
	}
	return ""
}

func (x *UploadSettings) GetNewPassword() string {
	if x != nil {
		return x.NewPassword
	}
	return ""
}

func (x *UploadSettings) GetOldSettings() *UserSettings {
	if x != nil {
		return x.OldSettings
	}
	return nil
}

type ProfileUserID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID int64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *ProfileUserID) Reset() {
	*x = ProfileUserID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileUserID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileUserID) ProtoMessage() {}

func (x *ProfileUserID) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileUserID.ProtoReflect.Descriptor instead.
func (*ProfileUserID) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{2}
}

func (x *ProfileUserID) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type EmptyProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyProfile) Reset() {
	*x = EmptyProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyProfile) ProtoMessage() {}

func (x *EmptyProfile) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyProfile.ProtoReflect.Descriptor instead.
func (*EmptyProfile) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{3}
}

var File_profile_proto protoreflect.FileDescriptor

var file_profile_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x80, 0x01, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x41, 0x76, 0x61, 0x74, 0x61,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x41, 0x76,
	0x61, 0x74, 0x61, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x42, 0x69, 0x67, 0x41, 0x76, 0x61, 0x74, 0x61,
	0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x42, 0x69, 0x67, 0x41, 0x76, 0x61, 0x74,
	0x61, 0x72, 0x22, 0xf7, 0x01, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x65, 0x74,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x14, 0x0a,
	0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x26, 0x0a, 0x0e, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x46,
	0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x6c, 0x64, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x6c,
	0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x4e, 0x65, 0x77,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x4e, 0x65, 0x77, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x2f, 0x0a, 0x0b, 0x4f,
	0x6c, 0x64, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52,
	0x0b, 0x4f, 0x6c, 0x64, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x1f, 0x0a, 0x0d,
	0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x22, 0x0e, 0x0a,
	0x0c, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x32, 0x6d, 0x0a,
	0x07, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x2e, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x53,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x0e, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x0d, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x0f, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x1a, 0x0d, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x00, 0x42, 0x1d, 0x5a, 0x1b,
	0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_profile_proto_rawDescOnce sync.Once
	file_profile_proto_rawDescData = file_profile_proto_rawDesc
)

func file_profile_proto_rawDescGZIP() []byte {
	file_profile_proto_rawDescOnce.Do(func() {
		file_profile_proto_rawDescData = protoimpl.X.CompressGZIP(file_profile_proto_rawDescData)
	})
	return file_profile_proto_rawDescData
}

var file_profile_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_profile_proto_goTypes = []interface{}{
	(*UserSettings)(nil),   // 0: UserSettings
	(*UploadSettings)(nil), // 1: UploadSettings
	(*ProfileUserID)(nil),  // 2: ProfileUserID
	(*EmptyProfile)(nil),   // 3: EmptyProfile
}
var file_profile_proto_depIdxs = []int32{
	0, // 0: UploadSettings.OldSettings:type_name -> UserSettings
	2, // 1: Profile.GetSettings:input_type -> ProfileUserID
	1, // 2: Profile.UpdateSettings:input_type -> UploadSettings
	0, // 3: Profile.GetSettings:output_type -> UserSettings
	3, // 4: Profile.UpdateSettings:output_type -> EmptyProfile
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_profile_proto_init() }
func file_profile_proto_init() {
	if File_profile_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_profile_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserSettings); i {
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
		file_profile_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadSettings); i {
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
		file_profile_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileUserID); i {
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
		file_profile_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyProfile); i {
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
			RawDescriptor: file_profile_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_profile_proto_goTypes,
		DependencyIndexes: file_profile_proto_depIdxs,
		MessageInfos:      file_profile_proto_msgTypes,
	}.Build()
	File_profile_proto = out.File
	file_profile_proto_rawDesc = nil
	file_profile_proto_goTypes = nil
	file_profile_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProfileClient is the client API for Profile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProfileClient interface {
	GetSettings(ctx context.Context, in *ProfileUserID, opts ...grpc.CallOption) (*UserSettings, error)
	UpdateSettings(ctx context.Context, in *UploadSettings, opts ...grpc.CallOption) (*EmptyProfile, error)
}

type profileClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileClient(cc grpc.ClientConnInterface) ProfileClient {
	return &profileClient{cc}
}

func (c *profileClient) GetSettings(ctx context.Context, in *ProfileUserID, opts ...grpc.CallOption) (*UserSettings, error) {
	out := new(UserSettings)
	err := c.cc.Invoke(ctx, "/Profile/GetSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) UpdateSettings(ctx context.Context, in *UploadSettings, opts ...grpc.CallOption) (*EmptyProfile, error) {
	out := new(EmptyProfile)
	err := c.cc.Invoke(ctx, "/Profile/UpdateSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServer is the server API for Profile service.
type ProfileServer interface {
	GetSettings(context.Context, *ProfileUserID) (*UserSettings, error)
	UpdateSettings(context.Context, *UploadSettings) (*EmptyProfile, error)
}

// UnimplementedProfileServer can be embedded to have forward compatible implementations.
type UnimplementedProfileServer struct {
}

func (*UnimplementedProfileServer) GetSettings(context.Context, *ProfileUserID) (*UserSettings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSettings not implemented")
}
func (*UnimplementedProfileServer) UpdateSettings(context.Context, *UploadSettings) (*EmptyProfile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSettings not implemented")
}

func RegisterProfileServer(s *grpc.Server, srv ProfileServer) {
	s.RegisterService(&_Profile_serviceDesc, srv)
}

func _Profile_GetSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Profile/GetSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetSettings(ctx, req.(*ProfileUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_UpdateSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadSettings)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).UpdateSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Profile/UpdateSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).UpdateSettings(ctx, req.(*UploadSettings))
	}
	return interceptor(ctx, in, info, handler)
}

var _Profile_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Profile",
	HandlerType: (*ProfileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSettings",
			Handler:    _Profile_GetSettings_Handler,
		},
		{
			MethodName: "UpdateSettings",
			Handler:    _Profile_UpdateSettings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}